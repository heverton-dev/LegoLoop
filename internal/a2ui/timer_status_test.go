package a2ui

import (
	"context"
	"testing"
	"time"
)

func TestSessionTimer_PausaERetomadaCongelamPrazo(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("legoloop-concepcao")
	sess.Timer().SetTimeout(120 * time.Millisecond)

	done := make(chan bool, 1)
	go func() {
		_, ok := sess.WaitAction(context.Background(), 0)
		done <- ok
	}()

	time.Sleep(30 * time.Millisecond)
	snap := sess.PauseTimer()
	if !snap.Paused {
		t.Fatalf("esperava timer pausado, obteve: %+v", snap)
	}
	// Espera mais do que o timeout original enquanto pausado: nao pode expirar!
	time.Sleep(160 * time.Millisecond)
	select {
	case ok := <-done:
		t.Fatalf("timer expirou enquanto estava pausado (ok=%v)", ok)
	default:
	}

	snap = sess.ResumeTimer()
	if snap.Paused {
		t.Fatalf("esperava timer retomado, obteve: %+v", snap)
	}
	if err := hub.SubmitAction(Action{SessionID: sess.ID, SurfaceID: "proc:c", Name: "enviar"}); err != nil {
		t.Fatalf("SubmitAction falhou: %v", err)
	}
	select {
	case ok := <-done:
		if !ok {
			t.Fatalf("esperava receber acao apos retomar, mas deu timeout")
		}
	case <-time.After(time.Second):
		t.Fatal("WaitAction nao retornou apos SubmitAction")
	}
}

func TestSessionTimer_ExpiraQuandoNaoPausado(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("legoloop-concepcao")
	_, ok := sess.WaitAction(context.Background(), 50*time.Millisecond)
	if ok {
		t.Fatalf("esperava timeout, mas retornou ok=true")
	}
	snap := sess.Timer().Snapshot()
	if !snap.Expired {
		t.Fatalf("esperava Expired=true apos timeout, obteve %+v", snap)
	}
}

func TestSessionTimer_ConfigurarEEstender(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("legoloop-concepcao")
	sess.ConfigureTimer(600 * time.Second)
	snap := sess.Timer().Snapshot()
	if snap.TimeoutMs != 600000 {
		t.Fatalf("esperava timeout 600000ms, obteve %d", snap.TimeoutMs)
	}
	snap = sess.ExtendTimer(5 * time.Minute)
	if snap.RemainingMs < 300000 {
		t.Fatalf("esperava remaining >= 300000ms apos estender, obteve %d", snap.RemainingMs)
	}
}

func TestSessionStatus_PublicaEAssinaControlEvent(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("legoloop-diorama-ativos")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := hub.SubscribeControl(ctx, sess.ID)
	if err != nil {
		t.Fatalf("SubscribeControl falhou: %v", err)
	}
	st := SessionStatus{
		Etapa:    "execucao",
		Mensagem: "Rodando go test ./internal/diorama/ativos/...",
		Etapas: []EtapaStatus{
			{ID: "e1", Titulo: "Aprovar INTENT, SPEC e PLAN", Status: "done"},
			{ID: "e2", Titulo: "Implementar Etapas 3, 4 e 5 via TDD", Status: "running"},
		},
	}
	if err := hub.PublishStatus(sess.ID, st); err != nil {
		t.Fatalf("PublishStatus falhou: %v", err)
	}

	var recebeuStatus bool
	prazo := time.After(time.Second)
	for !recebeuStatus {
		select {
		case ev := <-ch:
			if ev.Kind == "status" && ev.Status != nil && len(ev.Status.Etapas) == 2 {
				recebeuStatus = true
			}
		case <-prazo:
			t.Fatal("evento de status nao chegou em SubscribeControl")
		}
	}
}

func TestSubmitActionMarcaAguardandoRespostaParaReloadF5(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("legoloop-diorama-ativos")
	if err := hub.Publish(Envelope{
		SessionID:   sess.ID,
		SurfaceID:   "proc:etapa1",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		Components:  []ComponentRef{{ID: "sg1", Kind: "LegoStopGatePanel"}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := hub.SubmitAction(Action{
		SessionID: sess.ID,
		SurfaceID: "proc:etapa1",
		Origin:    "ui",
		Name:      "stop_gate_decision",
	}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrlCh, err := hub.SubscribeControl(ctx, sess.ID)
	if err != nil {
		t.Fatal(err)
	}
	var viuWaiting bool
	for i := 0; i < 2; i++ {
		ev := <-ctrlCh
		if ev.Kind == "status" && ev.Status != nil && ev.Status.AguardandoResposta {
			viuWaiting = true
		}
	}
	if !viuWaiting {
		t.Fatalf("apos SubmitAction, replay no F5 deve informar AguardandoResposta=true")
	}

	if err := hub.Publish(Envelope{
		SessionID:   sess.ID,
		SurfaceID:   "proc:etapa1",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		Components:  []ComponentRef{{ID: "p2", Kind: "Panel"}},
	}); err != nil {
		t.Fatal(err)
	}
	ev2 := <-ctrlCh
	if ev2.Kind != "status" || ev2.Status == nil || ev2.Status.AguardandoResposta {
		t.Fatalf("novo Publish deve limpar AguardandoResposta, obtido: %+v", ev2)
	}
}
