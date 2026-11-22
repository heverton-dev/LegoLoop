package studio

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/combo"
	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/fabricar"
	"github.com/heverton-dev/LegoLoop/internal/graph"
)

func TestFabricarErrSemMatrizAndComposeTrigger(t *testing.T) {
	t.Run("fabricar_err_sem_matriz", verifyFabricarSemMatriz)
	t.Run("compose_trigger_e_combos", verifyComposeTriggerAndCombos)
	t.Run("compose_trigger_alvo_upstream", verifyTriggerDefaultsAndLegado)
}

func verifyFabricarSemMatriz(t *testing.T) {
	t.Helper()
	s, _ := newServer(t)
	vault := t.TempDir()
	s.Fabricar = func(ctx context.Context, req fabricar.Request) (*fabricar.Result, error) {
		semMatriz := func() (fabricar.FonteMatriz, error) {
			return fabricar.FonteMatriz{}, fabricar.ErrSemMatriz
		}
		return fabricar.Fabricar(ctx, req, fabricar.Deps{
			LLM:    &specOnlyLLM{},
			Matriz: semMatriz,
			Roots:  fabricar.Roots{Vault: vault, Docs: s.ProjectDir},
		})
	}
	code, body := postFabricar(t, s, okBody)
	if code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, esperava 422; corpo: %s", code, body)
	}
	var errResp map[string]string
	if err := json.Unmarshal([]byte(body), &errResp); err != nil {
		t.Fatalf("json inválido: %v", err)
	}
	if errResp["code"] != "sem_matriz" || !strings.Contains(errResp["error"], "legoloop matriz gerar") {
		t.Fatalf("resposta inesperada para ErrSemMatriz: %+v", errResp)
	}
}

func verifyComposeTriggerAndCombos(t *testing.T) {
	t.Helper()
	s := newFullCatalogServer(t)
	for _, nomeCombo := range combo.ListarCombos() {
		payload := buildTriggerFlowPayload("trigger-processo-fim", nomeCombo, false)
		res := postValidate(t, s, payload)
		if !res.OK {
			t.Fatalf("combo %s deveria validar com gatilho -> interview: %+v", nomeCombo, res.Problemas)
		}
	}
	badPayload := buildTriggerFlowPayload("trigger-processo-fim", "t-codigo-128k", true)
	badRes := postValidate(t, s, badPayload)
	if badRes.OK || len(badRes.Problemas) == 0 {
		t.Fatalf("conexão de entrada em gatilho deveria falhar: %+v", badRes)
	}
}

func verifyTriggerDefaultsAndLegado(t *testing.T) {
	t.Helper()
	s := newFullCatalogServer(t)
	if trg := triggerForKind(graph.KindTriggerProcessoFim); trg == nil || trg.Processo != "upstream" {
		t.Fatalf("triggerForKind processo_fim inesperado: %+v", trg)
	}
	if trg := triggerForKind(graph.KindTriggerFimLegado); trg == nil || !hasTriggerTarget(trg, "upstream") {
		t.Fatalf("triggerForKind legado inesperado: %+v", trg)
	}
	comp := Composition{
		Name:    "processo-gatilho-padrao",
		Trigger: &dag.Trigger{Type: dag.TriggerProcessoFim},
		Nodes: []CompNode{
			{ID: "t1", Block: "trigger-processo-fim", Routing: dag.RoutingConfig{Mode: dag.ModeAuto, Combo: "t-codigo-128k"}},
			{ID: "n1", Block: "interview", Routing: dag.RoutingConfig{Mode: dag.ModeAuto, Combo: "t-chat-128k"}},
		},
		Edges: []CompEdge{{From: "t1", To: "n1"}},
	}
	res, parsed := compose(s.Index, s.inventory(t.Context()), s.resolver(t.Context()), comp)
	if !res.OK || parsed == nil {
		t.Fatalf("composição com processo_fim vazio deveria assumir upstream: %+v", res.Problemas)
	}
	eff := parsed.EffectiveTrigger()
	if !hasTriggerTarget(&eff, "upstream") {
		t.Fatalf("alvo padrão esperado 'upstream', obteve %+v", eff)
	}
}

func hasTriggerTarget(trg *dag.Trigger, target string) bool {
	raw, err := json.Marshal(trg)
	if err != nil {
		return false
	}
	return strings.Contains(string(raw), `"`+target+`"`)
}

func newFullCatalogServer(t *testing.T) *Server {
	t.Helper()
	s, _ := newServer(t)
	fullIdx, err := graph.NewIndex(graph.FullCatalog())
	if err != nil {
		t.Fatalf("falha ao criar catálogo completo: %v", err)
	}
	s.Index = fullIdx
	return s
}

func buildTriggerFlowPayload(trigBlock, comboName string, reverseEdge bool) string {
	edge := `{"from":"t1","to":"n1"}`
	if reverseEdge {
		edge = `{"from":"n1","to":"t1"}`
	}
	return `{"name":"e2e-trigger","nodes":[` +
		`{"id":"t1","block":"` + trigBlock + `","routing":{"mode":"auto","combo":"` + comboName + `"}},` +
		`{"id":"n1","block":"interview","routing":{"mode":"auto","combo":"` + comboName + `"}}],` +
		`"edges":[` + edge + `]}`
}

func postValidate(t *testing.T, s *Server, payload string) ValidationResult {
	t.Helper()
	rec := do(t, s.Handler(), "POST", "/api/validate", payload)
	var res ValidationResult
	decode(t, rec, &res)
	return res
}
