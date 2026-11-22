package a2ui

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestBuildSwarmEnvelopesCaminhoFeliz(t *testing.T) {
	st := SwarmState{
		SessionID: "sess_1",
		Title:     "Execução Swarm",
		Progress:  65,
		Workers: []SwarmWorker{
			{Name: "tick-06-1-go", Role: "backend", State: "working", Model: "gemini"},
			{Name: "tick-06-2-ui", Role: "frontend", State: "waiting", Model: "claude"},
		},
		Messages: []SwarmMessage{
			{Seq: 1, Agent: "coordinator", Text: "Iniciar etapa", Time: "2026-11-22T20:00:00Z"},
			{Seq: 2, Agent: "tick-06-1-go", Text: "Testes escritos", Time: "2026-11-22T20:01:00Z"},
		},
	}

	envs := BuildSwarmEnvelopes(st)
	if len(envs) != 2 {
		t.Fatalf("esperava 2 envelopes, obteve %d", len(envs))
	}
	if envs[0].MessageType != "createSurface" || envs[1].MessageType != "updateDataModel" {
		t.Fatalf("tipos inesperados: %s, %s", envs[0].MessageType, envs[1].MessageType)
	}
	if envs[0].SurfaceID != DefaultSwarmSurfaceID {
		t.Fatalf("surface_id padrão esperado %q, obteve %q", DefaultSwarmSurfaceID, envs[0].SurfaceID)
	}

	cat := DefaultCatalog()
	for i, env := range envs {
		if err := env.Validate(cat); err != nil {
			t.Fatalf("envelope %d inválido: %v", i, err)
		}
		if !strings.Contains(env.FallbackTxt, "Execução Swarm") || !strings.Contains(env.FallbackTxt, "65%") {
			t.Fatalf("fallback pouco descritivo: %q", env.FallbackTxt)
		}
	}

	temBoard := false
	temMeter := false
	for _, c := range envs[0].Components {
		if c.Kind == "LegoSwarmBoard" {
			temBoard = true
		}
		if c.Kind == "LegoProgressMeter" {
			temMeter = true
		}
	}
	if !temBoard || !temMeter {
		t.Fatalf("esperava LegoSwarmBoard (%v) e LegoProgressMeter (%v)", temBoard, temMeter)
	}
}

func TestPublishSwarmCaminhoFelizESessaoInexistente(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("swarm-demo")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := hub.Subscribe(ctx, sess.ID)
	if err != nil {
		t.Fatalf("Subscribe falhou: %v", err)
	}

	st := SwarmState{
		SessionID: sess.ID,
		SurfaceID: "proc:swarm:demo",
		Title:     "Painel Swarm",
		Progress:  80,
		Workers: []SwarmWorker{
			{Name: "w1", Role: "autor", State: "working", Model: "m1"},
		},
		Messages: []SwarmMessage{
			{Seq: 10, Agent: "w1", Text: "Progresso normal", Time: "2026-11-22T20:05:00Z"},
		},
	}

	if err := PublishSwarm(hub, st); err != nil {
		t.Fatalf("PublishSwarm falhou: %v", err)
	}

	recebidos := 0
	for recebidos < 2 {
		select {
		case env := <-stream:
			if env.SurfaceID != "proc:swarm:demo" {
				t.Fatalf("surface_id inesperado: %q", env.SurfaceID)
			}
			recebidos++
		case <-time.After(time.Second):
			t.Fatalf("timeout esperando envelopes SSE (recebidos=%d)", recebidos)
		}
	}

	stInexistente := st
	stInexistente.SessionID = "sess_inexistente"
	if err := PublishSwarm(hub, stInexistente); !errors.Is(err, ErrSessaoNaoEncontrada) {
		t.Fatalf("esperava ErrSessaoNaoEncontrada, obteve %v", err)
	}

	if err := PublishSwarm(nil, st); !errors.Is(err, ErrSessaoNaoEncontrada) {
		t.Fatalf("esperava ErrSessaoNaoEncontrada com hub nil, obteve %v", err)
	}
}

func TestPublishSwarmRejeitaPayloadHostil(t *testing.T) {
	hub := NewHub(DefaultCatalog())
	sess := hub.CreateSession("swarm-xss")

	st := SwarmState{
		SessionID: sess.ID,
		Title:     "Swarm",
		Messages: []SwarmMessage{
			{Seq: 1, Agent: "atacante", Text: "<script>alert(1)</script>", Time: "2026-11-22T20:10:00Z"},
		},
	}
	if err := PublishSwarm(hub, st); !errors.Is(err, ErrPayloadHostil) {
		t.Fatalf("esperava ErrPayloadHostil para mensagem hostil, obteve %v", err)
	}
}

func TestParseSwarmPostAction(t *testing.T) {
	okAct := Action{
		SessionID: "sess_1",
		SurfaceID: DefaultSwarmSurfaceID,
		Origin:    "ui",
		Name:      "swarm_post",
		Payload:   map[string]any{"text": "  @worker-1 revisar testes  "},
	}
	txt, err := ParseSwarmPostAction(okAct)
	if err != nil {
		t.Fatalf("esperava sucesso, obteve erro: %v", err)
	}
	if txt != "@worker-1 revisar testes" {
		t.Fatalf("texto inesperado: %q", txt)
	}

	casosErro := []struct {
		nome string
		act  Action
		want error
	}{
		{
			nome: "nome_invalido",
			act:  Action{Name: "submit_answer", Payload: map[string]any{"text": "oi"}},
			want: ErrAcaoInvalida,
		},
		{
			nome: "payload_nulo",
			act:  Action{Name: "swarm_post"},
			want: ErrAcaoInvalida,
		},
		{
			nome: "texto_vazio",
			act:  Action{Name: "swarm_post", Payload: map[string]any{"text": "   "}},
			want: ErrAcaoInvalida,
		},
		{
			nome: "texto_nao_string",
			act:  Action{Name: "swarm_post", Payload: map[string]any{"text": 123}},
			want: ErrAcaoInvalida,
		},
		{
			nome: "payload_hostil_script",
			act:  Action{Name: "swarm_post", Payload: map[string]any{"text": "<script>alert(1)</script>"}},
			want: ErrPayloadHostil,
		},
		{
			nome: "payload_hostil_javascript",
			act:  Action{Name: "swarm_post", Payload: map[string]any{"text": "javascript:alert(1)"}},
			want: ErrPayloadHostil,
		},
	}

	for _, tc := range casosErro {
		t.Run(tc.nome, func(t *testing.T) {
			if _, err := ParseSwarmPostAction(tc.act); !errors.Is(err, tc.want) {
				t.Fatalf("caso %s: esperava %v, obteve %v", tc.nome, tc.want, err)
			}
		})
	}
}
