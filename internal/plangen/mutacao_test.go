package plangen

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"legoloop/internal/interview"
)

func TestMutacao_GateSafeEReadAndHash(t *testing.T) {
	if gateSafe("") || gateSafe("/abs/a.go") || gateSafe("../fora.go") || gateSafe("a'b.go") || gateSafe("a:b.go") {
		t.Fatal("gateSafe deveria recusar caminhos vazios, absolutos, com traversal ou caracteres proibidos")
	}
	if !gateSafe("internal/plangen/plangen.go") {
		t.Fatal("gateSafe deveria aceitar caminho relativo seguro")
	}
	if _, err := ReadAndHash(filepath.Join(t.TempDir(), "nao_existe.md")); err == nil {
		t.Fatal("ReadAndHash deveria retornar erro para arquivo inexistente")
	}
}

func TestMutacao_NormalizeParseTicketsEDeps(t *testing.T) {
	if _, err := normalizeTickets([]Ticket{{ID: ""}}, SupervisionEpic); err == nil {
		t.Fatal("normalizeTickets deveria rejeitar ID vazio")
	}
	if _, err := normalizeTickets([]Ticket{{ID: "TICK-01.1"}, {ID: "TICK-01.1"}}, SupervisionEpic); err == nil {
		t.Fatal("normalizeTickets deveria rejeitar ID duplicado")
	}
	norm, err := normalizeTickets([]Ticket{{ID: "TICK-01.1"}, {ID: "TICK-01.2"}}, SupervisionEpic)
	if err != nil || len(norm) != 2 {
		t.Fatalf("normalizeTickets falhou para IDs validos: err=%v len=%d", err, len(norm))
	}
	if parseDeps("") != nil || parseDeps("Nenhuma") != nil {
		t.Fatal("parseDeps deveria retornar nil para vazio ou Nenhuma")
	}
	if deps := parseDeps("TICK-01.1"); !reflect.DeepEqual(deps, []string{"TICK-01.1"}) {
		t.Fatalf("parseDeps inesperado: %v", deps)
	}
	mdDup := "### [TICK-01.1] Um\n- **Modo:** sequential\n### [TICK-01.1] Dois\n"
	if _, err := parseTickets(mdDup, SupervisionEpic); err == nil {
		t.Fatal("parseTickets deveria rejeitar ticket duplicado")
	}
	mdOk := "### [TICK-01.1] Um\n- **Modo:** sequential\n### [TICK-01.2] Dois\n"
	tks, err := parseTickets(mdOk, SupervisionEpic)
	if err != nil || len(tks) != 2 || tks[0].ExecutionMode != "sequential" || tks[1].ExecutionMode != "tdd_red_green" {
		t.Fatalf("parseTickets modos inesperados: err=%v tks=%+v", err, tks)
	}
}

type llmFalsoMut struct{ resp string }

func (f llmFalsoMut) Chat(_ context.Context, _ []interview.Message) (string, error) {
	return f.resp, nil
}

func TestMutacao_MaterializeEGenerateErroBaseDirInvalido(t *testing.T) {
	tmp := t.TempDir()
	arq := filepath.Join(tmp, "nao_dir")
	if err := os.WriteFile(arq, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := materialize(arq, "docs/specs/SPEC.md", "# SPEC", SupervisionEpic, []Ticket{{ID: "TICK-01.1"}}, time.Now())
	if err == nil || !strings.Contains(err.Error(), "falha ao criar") {
		t.Fatalf("materialize deveria falhar em newPathSet (falha ao criar), obtido: %v", err)
	}
	specDir := filepath.Join(tmp, "docs", "specs")
	intentDir := filepath.Join(tmp, "docs", "intents")
	_ = os.MkdirAll(specDir, 0o755)
	_ = os.MkdirAll(intentDir, 0o755)
	intentPath := filepath.Join(intentDir, "INTENT_X.md")
	_ = os.WriteFile(intentPath, []byte("---\nstatus: APPROVED\n---\n# INTENT\n"), 0o600)
	specPath := filepath.Join(specDir, "SPEC_X.md")
	_ = os.WriteFile(specPath, []byte("# SPEC\n**Referência da Intenção:** docs/intents/INTENT_X.md\n"), 0o600)
	_ = os.WriteFile(filepath.Join(tmp, "docs", "plans"), []byte("bloqueio"), 0o600)
	gen := NewGenerator(llmFalsoMut{resp: "### [TICK-01.1] Teste\n- **Épico:** E01\n"})
	if _, err := gen.Generate(context.Background(), tmp, specPath, SupervisionEpic); err == nil || !strings.Contains(err.Error(), "falha ao criar") {
		t.Fatalf("Generate deveria propagar erro de materialize, obtido: %v", err)
	}
}

func TestMutacao_RoutePoliticasETokenBudget(t *testing.T) {
	if !defaultRollbackPolicy().PreserveArtifacts {
		t.Fatal("PreserveArtifacts deveria ser true")
	}
	if !defaultStrategy(SupervisionEpic).TDDEnforced {
		t.Fatal("TDDEnforced deveria ser true")
	}
	if iso := localIsolation("local_process", ""); iso.WriteBoundary != "strict" {
		t.Fatalf("WriteBoundary vazio deveria virar strict, obtido %q", iso.WriteBoundary)
	}
	if iso := localIsolation("local_process", "none"); iso.WriteBoundary != "none" {
		t.Fatalf("WriteBoundary explicito deveria ser mantido, obtido %q", iso.WriteBoundary)
	}
	if rp := defaultRetryPolicy(true); rp.MaxRetries != 0 || rp.FallbackModelOnRetry {
		t.Fatalf("stopGate retry policy inesperada: %+v", rp)
	}
	if rp := defaultRetryPolicy(false); rp.MaxRetries != 1 || !rp.FallbackModelOnRetry {
		t.Fatalf("unit retry policy inesperada: %+v", rp)
	}
	bNano, rNano := tokenBudgetFor("nano")
	bMicro, rMicro := tokenBudgetFor("micro")
	if bNano.MaxOutputTokens != 300 || rNano.Model != "code-fast" || bMicro.MaxOutputTokens != 600 || rMicro.Model != "code-pro" {
		t.Fatalf("tokenBudgetFor inesperado: nano=%+v micro=%+v", bNano, bMicro)
	}
}

func TestMutacao_RouteNodesSourceFilesEIniciativa(t *testing.T) {
	meta := planMeta{SpecRel: "docs/specs/S.md", EpicInitiative: map[int]int{1: 0, 2: 4}}
	if got := resolveInitiative(Ticket{EpicID: "E01"}, meta); got != nil {
		t.Fatalf("iniciativa 0 deveria retornar nil, obtido %v", *got)
	}
	if got := resolveInitiative(Ticket{EpicID: "E02"}, meta); got == nil || *got != "I04" {
		t.Fatalf("iniciativa 4 deveria retornar I04, obtido %v", got)
	}
	if got := sourceFiles([]string{"a_test.go"}); !reflect.DeepEqual(got, []string{"a_test.go"}) {
		t.Fatalf("sourceFiles so com testes inesperado: %v", got)
	}
	if got := sourceFiles([]string{"a.go", "a_test.go"}); !reflect.DeepEqual(got, []string{"a.go"}) {
		t.Fatalf("sourceFiles misto inesperado: %v", got)
	}
	tn := ticketNode(Ticket{ID: "TICK-01.1", Nature: "micro", ExecutionMode: "tdd_red_green"}, meta)
	if tn.RetryPolicy.MaxRetries != 1 || !tn.RetryPolicy.FallbackModelOnRetry {
		t.Fatalf("ticketNode retry inesperado: %+v", tn.RetryPolicy)
	}
	hg := humanGate("SG-E01", "ok?", []string{"TICK-01.1"}, "E01")
	if hg.Isolation.WriteBoundary != "none" || hg.RetryPolicy.MaxRetries != 0 || hg.RetryPolicy.FallbackModelOnRetry {
		t.Fatalf("humanGate inesperado: iso=%+v retry=%+v", hg.Isolation, hg.RetryPolicy)
	}
	fg := finalGate([]string{"SG-E01"})
	if fg.Isolation.WriteBoundary != "strict" {
		t.Fatalf("finalGate deveria ter WriteBoundary strict, obtido %q", fg.Isolation.WriteBoundary)
	}
}
