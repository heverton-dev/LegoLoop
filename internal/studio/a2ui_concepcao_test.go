package studio

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
	"github.com/heverton-dev/LegoLoop/internal/interview"
)

type stubInterviewLLM struct {
	resps []string
	idx   int
}

func (f *stubInterviewLLM) Chat(_ context.Context, _ []interview.Message) (string, error) {
	if f.idx >= len(f.resps) {
		return f.resps[len(f.resps)-1], nil
	}
	r := f.resps[f.idx]
	f.idx++
	return r, nil
}

func TestA2UIConcepcaoAoVivoNoStudio(t *testing.T) {
	s, _ := newServer(t)
	s.Interview = interview.NewEngine(&stubInterviewLLM{
		resps: []string{
			"HIPÓTESE: O usuário quer um painel de notas.\nCONFIANÇA: 85%\nPERGUNTA: Qual o público-alvo principal?\nPALPITE: Desenvolvedores locais",
			"HIPÓTESE: Painel local de notas.\nCONFIANÇA: 97%\nESPELHO:\n- Resultado esperado: Painel local\n- Quem se beneficia: Desenvolvedores\n- Por que agora: Organização diária\n- Critério de sucesso: Testes verdes\n- Restrições inegociáveis: Zero nuvem\n- Fora de escopo: Mobile",
		},
	})
	h := s.Handler()

	rec := do(t, h, http.MethodGet, "/concepcao", "")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("GET /concepcao esperava 303, obteve %d body=%s", rec.Code, rec.Body.String())
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "/sessao/") {
		t.Fatalf("Location inesperado: %q", loc)
	}
	sessID := strings.TrimPrefix(loc, "/sessao/")

	act1 := do(t, h, http.MethodPost, "/api/sessao/"+sessID+"/action", `{"surface_id":"proc:concepcao","origin":"ui","name":"submit_answer","payload":{"answer":"Criar painel de notas"}}`)
	if act1.Code != http.StatusAccepted {
		t.Fatalf("turno 1 code=%d body=%s", act1.Code, act1.Body.String())
	}
	surf1, ok := s.A2UIHub().Surface(sessID, a2ui.DefaultConcepcaoSurfaceID)
	if !ok || surf1.DataModel["confianca"] != 85 {
		t.Fatalf("esperava confiança 85 na superfície após turno 1, obteve ok=%v dm=%+v", ok, surf1.DataModel)
	}

	act2 := do(t, h, http.MethodPost, "/api/sessao/"+sessID+"/action", `{"surface_id":"proc:concepcao","origin":"ui","name":"submit_answer","payload":{"answer":"Desenvolvedores locais"}}`)
	if act2.Code != http.StatusAccepted {
		t.Fatalf("turno 2 code=%d body=%s", act2.Code, act2.Body.String())
	}
	surf2, ok := s.A2UIHub().Surface(sessID, a2ui.DefaultConcepcaoSurfaceID)
	if !ok || surf2.DataModel["confianca"] != 97 || !surf2.RequiresAttention {
		t.Fatalf("esperava confiança 97 e RequiresAttention=true no Stop-Gate 01, obteve ok=%v att=%v dm=%+v", ok, surf2.RequiresAttention, surf2.DataModel)
	}

	act3 := do(t, h, http.MethodPost, "/api/sessao/"+sessID+"/action", `{"surface_id":"proc:concepcao","origin":"ui","name":"stop_gate_decision","payload":{"decision":"1"}}`)
	if act3.Code != http.StatusAccepted {
		t.Fatalf("stop-gate code=%d body=%s", act3.Code, act3.Body.String())
	}
	surf3, ok := s.A2UIHub().Surface(sessID, a2ui.DefaultConcepcaoSurfaceID)
	intentPath, _ := surf3.DataModel["intent_path"].(string)
	if !ok || strings.TrimSpace(intentPath) == "" {
		t.Fatalf("esperava intent_path materializado na superfície após Stop-Gate 01, obteve dm=%+v", surf3.DataModel)
	}
}
