package a2ui

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestConcepcaoEvolucaoCompletaEComponentes(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("legoloop-concepcao")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := hub.Subscribe(ctx, sess.ID)
	if err != nil {
		t.Fatalf("Subscribe falhou: %v", err)
	}

	turno1 := ConcepcaoState{
		SessionID: sess.ID,
		SurfaceID: "proc:concepcao",
		Confianca: 85,
		Meta:      95,
		Turno:     1,
		Pergunta:  "Qual o critério de sucesso principal?",
		Palpite:   "Reduzir tempo de entrega em 30%",
		Lacunas:   []string{"critério de sucesso", "restrições"},
	}
	if err := PublishConcepcao(hub, turno1); err != nil {
		t.Fatalf("PublishConcepcao turno 1 falhou: %v", err)
	}

	turno2 := ConcepcaoState{
		SessionID: sess.ID,
		SurfaceID: "proc:concepcao",
		Confianca: 95,
		Meta:      95,
		Turno:     2,
		Pergunta:  "Confirma o Espelho Canônico?",
		Palpite:   "Aprovado sem ajustes",
		Espelho: map[string]string{
			"resultado":  "Motor A2UI nas 3 vias",
			"usuario":    "Operador do LegoLoop",
			"restricoes": "Zero dependências externas no go.mod",
		},
		IntentPath:         "docs/intents/2026-11-22_INTENT_PILOTO.md",
		AguardandoStopGate: true,
	}
	envs := BuildConcepcaoEnvelopes(turno2)
	if len(envs) == 0 {
		t.Fatalf("BuildConcepcaoEnvelopes retornou vazio")
	}
	cat := DefaultCatalog()
	kinds := make(map[string]bool)
	for _, env := range envs {
		if err := env.Validate(cat); err != nil {
			t.Fatalf("envelope inválido: %v", err)
		}
		if strings.TrimSpace(env.FallbackTxt) == "" {
			t.Fatalf("envelope %s sem FallbackTxt", env.MessageType)
		}
		for _, c := range env.Components {
			kinds[c.Kind] = true
		}
	}

	esperados := []string{
		"LegoProgressMeter",
		"TextField",
		"Button",
		"Card",
		"LegoStopGateCard",
		"LegoArtifactViewer",
	}
	for _, k := range esperados {
		if !kinds[k] {
			t.Errorf("componente esperado ausente: %s", k)
		}
	}

	if err := PublishConcepcao(hub, turno2); err != nil {
		t.Fatalf("PublishConcepcao turno 2 falhou: %v", err)
	}
	if got := hub.PendingAttention(sess.ID); got != "proc:concepcao" {
		t.Errorf("PendingAttention = %q, esperado %q", got, "proc:concepcao")
	}

	select {
	case env := <-stream:
		if env.SessionID != sess.ID {
			t.Errorf("SessionID no SSE = %q, esperado %q", env.SessionID, sess.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout aguardando evento SSE")
	}
}

func TestConcepcaoStopGateOpcoes(t *testing.T) {
	st := ConcepcaoState{
		SessionID:          "sess_1",
		SurfaceID:          "proc:concepcao",
		Confianca:          95,
		Meta:               95,
		Turno:              3,
		AguardandoStopGate: true,
	}
	envs := BuildConcepcaoEnvelopes(st)
	var stopComp *ComponentRef
	for i := range envs {
		for j := range envs[i].Components {
			if envs[i].Components[j].Kind == "LegoStopGateCard" {
				stopComp = &envs[i].Components[j]
			}
		}
	}
	if stopComp == nil {
		t.Fatal("LegoStopGateCard não encontrado")
	}
	opcoes, ok := stopComp.Props["options"].([]string)
	if !ok || len(opcoes) != 3 {
		t.Fatalf("opções de Stop-Gate inválidas: %#v", stopComp.Props["options"])
	}
	if opcoes[0] != "[1] seguir" || opcoes[1] != "[2] personalizar" || opcoes[2] != "[0] cancelar" {
		t.Errorf("opções inesperadas: %v", opcoes)
	}
}

func TestConcepcaoFallbackSemAssinantesSSE(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("legoloop-concepcao")

	if got := hub.SubscriberCount(sess.ID); got != 0 {
		t.Fatalf("SubscriberCount = %d, esperado 0", got)
	}

	done := make(chan error, 1)
	go func() {
		st := ConcepcaoState{
			SessionID:          sess.ID,
			Confianca:          95,
			Meta:               95,
			Turno:              2,
			Pergunta:           "Aprovar INTENT?",
			Palpite:            "Sim",
			IntentPath:         "docs/intents/2026-11-22_INTENT_TEXTO.md",
			AguardandoStopGate: true,
		}
		done <- PublishConcepcao(hub, st)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("PublishConcepcao sem SSE falhou: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("PublishConcepcao travou sem assinantes SSE")
	}

	surf, ok := hub.Surface(sess.ID, DefaultConcepcaoSurfaceID)
	if !ok {
		t.Fatal("superfície não registrada no snapshot sem SSE")
	}
	if !strings.Contains(surf.FallbackTxt, "95%") || !strings.Contains(surf.FallbackTxt, "[1] seguir") {
		t.Errorf("FallbackTxt incompleto: %q", surf.FallbackTxt)
	}
	if !strings.Contains(surf.FallbackTxt, "docs/intents/2026-11-22_INTENT_TEXTO.md") {
		t.Errorf("FallbackTxt sem caminho do INTENT: %q", surf.FallbackTxt)
	}
}

func TestPublishConcepcaoRejeitaPayloadHostilENilHub(t *testing.T) {
	if err := PublishConcepcao(nil, ConcepcaoState{SessionID: "sess_1"}); err != ErrSessaoNaoEncontrada {
		t.Errorf("PublishConcepcao(nil) = %v, esperado ErrSessaoNaoEncontrada", err)
	}

	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("legoloop-concepcao")
	st := ConcepcaoState{
		SessionID: sess.ID,
		Confianca: 50,
		Pergunta:  "<script>alert(1)</script>",
	}
	if err := PublishConcepcao(hub, st); err != ErrPayloadHostil {
		t.Errorf("PublishConcepcao hostil = %v, esperado ErrPayloadHostil", err)
	}
}
