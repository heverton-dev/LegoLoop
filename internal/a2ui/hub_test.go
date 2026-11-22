package a2ui

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestHubCreatePublishSubscribeEAction(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("legoloop-concepcao")
	if sess.ID == "" {
		t.Fatal("esperava session_id não vazio")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub, err := hub.Subscribe(ctx, sess.ID)
	if err != nil {
		t.Fatalf("Subscribe falhou: %v", err)
	}

	env := Envelope{
		SessionID:   sess.ID,
		SurfaceID:   "proc:concepcao",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		Components:  []ComponentRef{{ID: "c1", Kind: "Card"}},
	}
	if err := hub.Publish(env); err != nil {
		t.Fatalf("Publish falhou: %v", err)
	}

	select {
	case got := <-sub:
		if got.SurfaceID != "proc:concepcao" || got.MessageType != "createSurface" {
			t.Fatalf("envelope inesperado: %+v", got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout esperando envelope no subscriber")
	}

	act := Action{SessionID: sess.ID, SurfaceID: "proc:concepcao", Origin: "ui", Name: "submit_answer", Payload: map[string]any{"answer": "sim"}}
	if err := hub.SubmitAction(act); err != nil {
		t.Fatalf("SubmitAction falhou: %v", err)
	}

	select {
	case got := <-sess.Actions():
		if got.Name != "submit_answer" || got.Payload["answer"] != "sim" {
			t.Fatalf("action inesperada: %+v", got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout esperando action no canal da sessão")
	}
}

func TestHubSnapshotDeReconexao(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("piloto")

	e1 := Envelope{SessionID: sess.ID, SurfaceID: "proc:p1", ProtocolVer: "v0.9", MessageType: "createSurface", Components: []ComponentRef{{ID: "c1", Kind: "Text"}}}
	e2 := Envelope{SessionID: sess.ID, SurfaceID: "proc:p1", ProtocolVer: "v0.9", MessageType: "updateDataModel", DataModel: map[string]any{"confianca": 95}}
	if err := hub.Publish(e1); err != nil {
		t.Fatalf("e1: %v", err)
	}
	if err := hub.Publish(e2); err != nil {
		t.Fatalf("e2: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub, err := hub.Subscribe(ctx, sess.ID)
	if err != nil {
		t.Fatalf("Subscribe tardio falhou: %v", err)
	}

	select {
	case snap := <-sub:
		if snap.SurfaceID != "proc:p1" || snap.DataModel["confianca"] != 95 {
			t.Fatalf("snapshot incompleto: %+v", snap)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout esperando snapshot de reconexão")
	}
}

func TestHubRejeitaSessaoInexistenteEHostil(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	if _, err := hub.Subscribe(context.Background(), "sess_nao_existe"); !errors.Is(err, ErrSessaoNaoEncontrada) {
		t.Fatalf("esperava ErrSessaoNaoEncontrada, obteve %v", err)
	}
	sess := hub.CreateSession("teste")
	err := hub.Publish(Envelope{
		SessionID:   sess.ID,
		SurfaceID:   "proc:x",
		ProtocolVer: "v0.9",
		MessageType: "updateComponents",
		Components:  []ComponentRef{{ID: "b1", Kind: "RawHTML"}},
	})
	if !errors.Is(err, ErrComponenteProibido) {
		t.Fatalf("esperava ErrComponenteProibido, obteve %v", err)
	}
}

func TestHubCinquentaSubscribersSemVazamento(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("carga")

	const total = 50
	cancels := make([]context.CancelFunc, total)
	var wg sync.WaitGroup
	for i := 0; i < total; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancels[i] = cancel
		ch, err := hub.Subscribe(ctx, sess.ID)
		if err != nil {
			t.Fatalf("sub %d falhou: %v", i, err)
		}
		wg.Add(1)
		go func(c <-chan Envelope) {
			defer wg.Done()
			for range c {
			}
		}(ch)
	}

	for _, c := range cancels {
		c()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("subscribers não encerraram dentro de 500ms")
	}

	if n := hub.SubscriberCount(sess.ID); n != 0 {
		t.Fatalf("esperava 0 subscribers ativos após cancelamento, obteve %d", n)
	}
}
