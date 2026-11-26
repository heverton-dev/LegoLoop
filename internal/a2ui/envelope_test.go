package a2ui

import (
	"errors"
	"testing"
)

func TestEnvelopeValidateAceitaCatalogoConfiavel(t *testing.T) {
	cat := DefaultCatalog()
	env := Envelope{
		SessionID:   "sess_1",
		SurfaceID:   "proc:concepcao:entrevista",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		Components: []ComponentRef{
			{ID: "c1", Kind: "Panel", Props: map[string]any{"title": "Pergunta 1"}},
			{ID: "c2", Kind: "LegoProgressMeter", Props: map[string]any{"value": 85, "target": 95}},
			{ID: "c3", Kind: "LegoStopGatePanel", Props: map[string]any{"label": "Stop-Gate 01"}},
		},
		FallbackTxt: "Pergunta 1",
	}
	if err := env.Validate(cat); err != nil {
		t.Fatalf("esperava envelope valido, obteve erro: %v", err)
	}
}

func TestEnvelopeValidateRejeitaitemForaDoCatalogo(t *testing.T) {
	cat := DefaultCatalog()
	env := Envelope{
		SessionID:   "sess_1",
		SurfaceID:   "proc:concepcao",
		ProtocolVer: "v0.9",
		MessageType: "updateComponents",
		Components: []ComponentRef{
			{ID: "x1", Kind: "RawHTML", Props: map[string]any{"html": "<b>oi</b>"}},
		},
	}
	err := env.Validate(cat)
	if !errors.Is(err, ErrItemProibido) {
		t.Fatalf("esperava ErrItemProibido, obteve: %v", err)
	}
}

func TestEnvelopeValidateRejeitaPayloadHostil(t *testing.T) {
	cat := DefaultCatalog()
	casos := []struct {
		nome  string
		props map[string]any
	}{
		{"tag script", map[string]any{"text": "<script>alert(1)</script>"}},
		{"url javascript", map[string]any{"href": "javascript:alert(1)"}},
		{"evento inline", map[string]any{"label": "<img src=x onerror=alert(1)>"}},
		{"aninhado", map[string]any{"nested": map[string]any{"v": "<iframe src='x'>"}}},
		{"lista", map[string]any{"items": []any{"ok", "data:text/html,<svg onload=1>"}}},
	}
	for _, tc := range casos {
		t.Run(tc.nome, func(t *testing.T) {
			env := Envelope{
				SessionID:   "sess_1",
				SurfaceID:   "proc:concepcao",
				ProtocolVer: "v0.9",
				MessageType: "updateComponents",
				Components:  []ComponentRef{{ID: "c1", Kind: "Text", Props: tc.props}},
			}
			err := env.Validate(cat)
			if !errors.Is(err, ErrPayloadHostil) {
				t.Fatalf("caso %q: esperava ErrPayloadHostil, obteve %v", tc.nome, err)
			}
		})
	}
}

func TestEnvelopeValidateRejeitaSurfaceIDInvalido(t *testing.T) {
	cat := DefaultCatalog()
	invalidos := []string{"", "   ", "proc com espaco", "../escape", "a/../b", "proc?x=1"}
	for _, sid := range invalidos {
		env := Envelope{
			SessionID:   "sess_1",
			SurfaceID:   sid,
			ProtocolVer: "v0.9",
			MessageType: "createSurface",
		}
		err := env.Validate(cat)
		if !errors.Is(err, ErrSurfaceIDInvalido) {
			t.Fatalf("surface_id %q: esperava ErrSurfaceIDInvalido, obteve %v", sid, err)
		}
	}
}

func TestEnvelopeValidateRejeitaMessageTypeESessionInvalidos(t *testing.T) {
	cat := DefaultCatalog()
	envSemSessao := Envelope{SurfaceID: "proc:ok", ProtocolVer: "v0.9", MessageType: "createSurface"}
	if !errors.Is(envSemSessao.Validate(cat), ErrSessaoInvalida) {
		t.Fatalf("esperava ErrSessaoInvalida, obteve: %v", envSemSessao.Validate(cat))
	}
	envTipoRuim := Envelope{SessionID: "sess_1", SurfaceID: "proc:ok", ProtocolVer: "v0.9", MessageType: "execJS"}
	if !errors.Is(envTipoRuim.Validate(cat), ErrMensagemInvalida) {
		t.Fatalf("esperava ErrMensagemInvalida, obteve: %v", envTipoRuim.Validate(cat))
	}
}

func TestEnvelopeValidateNaoBloqueiaTextoBenignoComIgual(t *testing.T) {
	cat := DefaultCatalog()
	textos := []string{
		"sessionid=sess_3",
		"jsonoutput=true",
		"actionurl=http://127.0.0.1",
		"estilo_visual=codigo",
	}
	for _, txt := range textos {
		env := Envelope{
			SessionID:   "sess_1",
			SurfaceID:   "proc:e2e",
			ProtocolVer: "v0.9",
			MessageType: "createSurface",
			Components:  []ComponentRef{{ID: "c1", Kind: "Text", Props: map[string]any{"text": txt}}},
		}
		if err := env.Validate(cat); err != nil {
			t.Fatalf("texto benigno %q nao deveria ser rejeitado: %v", txt, err)
		}
	}
}
