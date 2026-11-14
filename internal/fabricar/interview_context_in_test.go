package fabricar_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/matriz"
)

func TestFabricarFromInterview_ConectaIntentComoContextIn(t *testing.T) {
	tempDir := t.TempDir()
	llm := &scriptedLLM{replies: []string{
		`{"confidence": 70, "question": "Qual o objetivo?", "socratic_mirror": {"outcome": "Validar fluxo"}}`,
		`{"confidence": 96, "question": "Pronto", "socratic_mirror": {"outcome": "Validar fluxo", "target_audience": "devs", "constraints": ["zero stubs"], "out_of_scope": ["deploy"]}}`,
	}}
	eng := interview.NewEngine(llm)
	sess, _, err := eng.Start(context.Background(), "processo de teste")
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if _, err := eng.Step(context.Background(), sess.ID, "Validar fluxo"); err != nil {
		t.Fatalf("Step: %v", err)
	}
	if sess.Mirror == nil {
		sess.Mirror = &interview.Mirror{Outcome: "Validar fluxo", Constraints: []string{"zero stubs"}, OutOfScope: []string{"deploy"}}
		sess.Confidence = 96
		sess.Status = interview.StatusCompleted
	}
	intentPath, err := eng.MaterializeIntent(sess.ID, tempDir, "2026-11-14_INTENT_TESTE_CONEXAO.md")
	if err != nil {
		t.Fatalf("MaterializeIntent: %v", err)
	}
	rawIntent, err := os.ReadFile(intentPath)
	if err != nil {
		t.Fatalf("ReadFile intent: %v", err)
	}
	if !strings.Contains(string(rawIntent), "UNDER_REVIEW") {
		t.Fatalf("INTENT esperada em UNDER_REVIEW, obtido:\n%s", string(rawIntent))
	}

	e1 := pecaComCasos("peca-construir", 0.90, []string{"construir"}, []string{"qualquer"}, []string{"Validar fluxo"})
	e2 := pecaComCasos("peca-verificar", 0.90, []string{"verificar"}, []string{"qualquer"}, []string{"Validar fluxo"})
	r := newRig(t, []matriz.Entrada{e1, e2})
	r.req.RawName = "teste-conexao"
	r.req.Goal = "Validar fluxo"

	res, err := fabricar.FabricarFromInterview(context.Background(), r.req, r.deps, sess, intentPath)
	if err != nil {
		t.Fatalf("FabricarFromInterview: %v", err)
	}
	if res == nil {
		t.Fatalf("resultado vazio")
	}

	// Verifica os 4 arquivos padrão em loops/<slug>/ e SPEC, PLAN, ROUTE, DAG em docs/.
	loopDir := filepath.Join(r.deps.Roots.Vault, "loops", res.CanonicalName)
	for _, f := range []string{res.CanonicalName + ".studio.json", "pipeline.dag.json", "triagem.json", "memoria_execucao.json"} {
		if _, err := os.Stat(filepath.Join(loopDir, f)); err != nil {
			t.Errorf("arquivo padrão ausente em loops/: %s (%v)", f, err)
		}
	}
	for _, sub := range []string{"specs", "plans", "routes", "dags"} {
		entries, err := os.ReadDir(filepath.Join(r.deps.Roots.Docs, "docs", sub))
		if err != nil || len(entries) == 0 {
			t.Errorf("artefato ausente em docs/%s: %v", sub, err)
		}
	}

	// Valida o DAG resultante com dag.Parse e verifica context_in apontando para a INTENT.
	dagBytes, err := os.ReadFile(filepath.Join(loopDir, "pipeline.dag.json"))
	if err != nil {
		t.Fatalf("ler pipeline.dag.json: %v", err)
	}
	parsed, err := dag.Parse(dagBytes)
	if err != nil {
		t.Fatalf("dag.Parse falhou: %v", err)
	}
	rawDAG, _ := json.Marshal(parsed)
	if !strings.Contains(string(rawDAG), "INTENT") {
		t.Errorf("context_in do DAG deveria referenciar a INTENT, obtido:\n%s", string(rawDAG))
	}
}

func TestFabricarFromInterview_RejeitaConfiancaMenorQue95(t *testing.T) {
	sess := &interview.Session{ID: "intv_1", Theme: "teste", Confidence: 80}
	r := newRig(t, nil)
	r.req.RawName = "teste-baixa-confianca"
	_, err := fabricar.FabricarFromInterview(context.Background(), r.req, r.deps, sess, "")
	if !errors.Is(err, fabricar.ErrLowConfidence) {
		t.Fatalf("esperado ErrLowConfidence, obtido: %v", err)
	}
}
