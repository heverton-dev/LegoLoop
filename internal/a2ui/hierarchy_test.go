package a2ui

import (
	"context"
	"reflect"
	"testing"
)

func TestHierarchyIsolamentoDataModel(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("fabricar")

	pai := Envelope{
		SessionID:   sess.ID,
		SurfaceID:   "proc:fabricar",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		DataModel:   map[string]any{"etapa": "pai", "status": "ativo"},
	}
	filho := Envelope{
		SessionID:   sess.ID,
		SurfaceID:   "proc:fabricar:concepcao",
		ParentID:    "proc:fabricar",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		DataModel:   map[string]any{"etapa": "concepcao", "confianca": 85},
	}
	if err := hub.Publish(pai); err != nil {
		t.Fatalf("Publish pai falhou: %v", err)
	}
	if err := hub.Publish(filho); err != nil {
		t.Fatalf("Publish filho falhou: %v", err)
	}

	updFilho := Envelope{
		SessionID:   sess.ID,
		SurfaceID:   "proc:fabricar:concepcao",
		ProtocolVer: "v0.9",
		MessageType: "updateDataModel",
		DataModel:   map[string]any{"etapa": "espelho", "status": "aguardando"},
	}
	if err := hub.Publish(updFilho); err != nil {
		t.Fatalf("Publish update filho falhou: %v", err)
	}

	envPai, ok := hub.Surface(sess.ID, "proc:fabricar")
	if !ok {
		t.Fatalf("superfície pai não encontrada")
	}
	envFilho, ok := hub.Surface(sess.ID, "proc:fabricar:concepcao")
	if !ok {
		t.Fatalf("superfície filho não encontrada")
	}

	if envPai.DataModel["etapa"] != "pai" || envPai.DataModel["status"] != "ativo" {
		t.Fatalf("DataModel do pai foi contaminado: %+v", envPai.DataModel)
	}
	if _, vazou := envPai.DataModel["confianca"]; vazou {
		t.Fatalf("chave do filho vazou para o pai: %+v", envPai.DataModel)
	}
	if envFilho.DataModel["etapa"] != "espelho" || envFilho.DataModel["confianca"] != 85 {
		t.Fatalf("DataModel do filho incorreto: %+v", envFilho.DataModel)
	}
}

func TestHierarchyTreeAndBreadcrumb(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("fabricar")

	envelopes := []Envelope{
		{
			SessionID:   sess.ID,
			SurfaceID:   "proc:fabricar",
			ProtocolVer: "v0.9",
			MessageType: "createSurface",
		},
		{
			SessionID:   sess.ID,
			SurfaceID:   "proc:fabricar:concepcao",
			ParentID:    "proc:fabricar",
			ProtocolVer: "v0.9",
			MessageType: "createSurface",
		},
		{
			SessionID:   sess.ID,
			SurfaceID:   "proc:fabricar:concepcao:entrevista",
			ParentID:    "proc:fabricar:concepcao",
			ProtocolVer: "v0.9",
			MessageType: "createSurface",
		},
	}
	for _, env := range envelopes {
		if err := hub.Publish(env); err != nil {
			t.Fatalf("Publish %s falhou: %v", env.SurfaceID, err)
		}
	}

	arvore, err := hub.Tree(sess.ID)
	if err != nil {
		t.Fatalf("Tree falhou: %v", err)
	}
	if len(arvore) != 1 || arvore[0].SurfaceID != "proc:fabricar" {
		t.Fatalf("raiz esperada proc:fabricar, obtido: %+v", arvore)
	}
	if len(arvore[0].Children) != 1 || arvore[0].Children[0].SurfaceID != "proc:fabricar:concepcao" {
		t.Fatalf("filho esperado proc:fabricar:concepcao, obtido: %+v", arvore[0].Children)
	}
	netos := arvore[0].Children[0].Children
	if len(netos) != 1 || netos[0].SurfaceID != "proc:fabricar:concepcao:entrevista" {
		t.Fatalf("neto esperado proc:fabricar:concepcao:entrevista, obtido: %+v", netos)
	}

	crumb, err := hub.Breadcrumb(sess.ID, "proc:fabricar:concepcao")
	if err != nil {
		t.Fatalf("Breadcrumb falhou: %v", err)
	}
	esperado := []string{"proc:fabricar", "proc:fabricar:concepcao"}
	if !reflect.DeepEqual(crumb, esperado) {
		t.Fatalf("breadcrumb esperado %v, obtido %v", esperado, crumb)
	}

	crumbNeto, err := hub.Breadcrumb(sess.ID, "proc:fabricar:concepcao:entrevista")
	if err != nil {
		t.Fatalf("Breadcrumb neto falhou: %v", err)
	}
	esperadoNeto := []string{"proc:fabricar", "proc:fabricar:concepcao", "proc:fabricar:concepcao:entrevista"}
	if !reflect.DeepEqual(crumbNeto, esperadoNeto) {
		t.Fatalf("breadcrumb neto esperado %v, obtido %v", esperadoNeto, crumbNeto)
	}
}

func TestHierarchyStopGateFocoAutomatico(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("fabricar")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := hub.Subscribe(ctx, sess.ID)
	if err != nil {
		t.Fatalf("Subscribe falhou: %v", err)
	}

	if err := hub.Publish(Envelope{
		SessionID:   sess.ID,
		SurfaceID:   "proc:fabricar",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
	}); err != nil {
		t.Fatalf("Publish pai falhou: %v", err)
	}
	<-stream

	if got := hub.PendingAttention(sess.ID); got != "" {
		t.Fatalf("esperado nenhum foco pendente inicialmente, obtido %q", got)
	}

	stopGate := Envelope{
		SessionID:   sess.ID,
		SurfaceID:   "proc:fabricar:concepcao",
		ParentID:    "proc:fabricar",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		Components: []ComponentRef{
			{ID: "sg-01", Kind: "LegoStopGatePanel", Props: map[string]any{"titulo": "Stop-Gate 01"}},
		},
	}
	if err := hub.Publish(stopGate); err != nil {
		t.Fatalf("Publish Stop-Gate falhou: %v", err)
	}

	ev := <-stream
	if !ev.RequiresAttention {
		t.Fatalf("evento SSE do Stop-Gate deve marcar RequiresAttention=true")
	}

	filho, ok := hub.Surface(sess.ID, "proc:fabricar:concepcao")
	if !ok || !filho.RequiresAttention {
		t.Fatalf("superfície filha deve estar com RequiresAttention=true, obtido %+v", filho)
	}
	if got := hub.PendingAttention(sess.ID); got != "proc:fabricar:concepcao" {
		t.Fatalf("PendingAttention esperado proc:fabricar:concepcao, obtido %q", got)
	}

	if err := hub.SubmitAction(Action{
		SessionID: sess.ID,
		SurfaceID: "proc:fabricar:concepcao",
		Origin:    "ui",
		Name:      "stop_gate_decision",
		Payload:   map[string]any{"decisao": "1"},
	}); err != nil {
		t.Fatalf("SubmitAction falhou: %v", err)
	}

	if got := hub.PendingAttention(sess.ID); got != "" {
		t.Fatalf("PendingAttention deveria limpar após SubmitAction, obtido %q", got)
	}
	filhoApos, _ := hub.Surface(sess.ID, "proc:fabricar:concepcao")
	if filhoApos.RequiresAttention {
		t.Fatalf("RequiresAttention da superfície filha deveria ser false após SubmitAction")
	}
}
