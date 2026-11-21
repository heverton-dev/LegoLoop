package plangen_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/plangen"
)

// TestRoute_AssinaOPlanJaGravadoNaPrimeiraGeracao: a ROUTE assina o PLAN. Na primeira geração o PLAN
// ainda não existia no disco e a assinatura saía zerada, o que desligava a verificação de integridade.
func TestRoute_AssinaOPlanJaGravadoNaPrimeiraGeracao(t *testing.T) {
	res := generated(t)
	_, route := readBoth(t, res)
	var doc struct {
		Upstream plangen.UpstreamSignatures `json:"upstream_signatures"`
	}
	if err := json.Unmarshal([]byte(route), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Upstream.Plan.SHA256 == strings.Repeat("0", 64) {
		t.Fatal("a assinatura do PLAN veio zerada na primeira geração")
	}
	plan, err := os.ReadFile(res.PlanMarkdownPath)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Upstream.Plan.SHA256 != plangen.HashNormalized(plan) {
		t.Fatalf("assinatura do PLAN %s não confere com o arquivo %s", doc.Upstream.Plan.SHA256, plangen.HashNormalized(plan))
	}
}

// TestRoute_AlterarOPlanDepoisQuebraAVerificacao: com a assinatura real, mexer no PLAN é detectado.
func TestRoute_AlterarOPlanDepoisQuebraAVerificacao(t *testing.T) {
	res := generated(t)
	_, route := readBoth(t, res)
	var doc struct {
		Upstream plangen.UpstreamSignatures `json:"upstream_signatures"`
	}
	if err := json.Unmarshal([]byte(route), &doc); err != nil {
		t.Fatal(err)
	}
	base := filepath.Dir(filepath.Dir(filepath.Dir(res.PlanMarkdownPath))) // base/docs/plans/arquivo
	if err := plangen.VerifySignatures(base, doc.Upstream); err != nil {
		t.Fatalf("plano recém-gerado deve verificar: %v", err)
	}
	if err := os.WriteFile(res.PlanMarkdownPath, []byte("adulterado"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := plangen.VerifySignatures(base, doc.Upstream); err == nil {
		t.Fatal("PLAN adulterado deveria falhar na verificação")
	}
}
