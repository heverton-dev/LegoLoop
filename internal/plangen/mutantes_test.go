package plangen

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestMutarAssembleCanonEpicETerminals(t *testing.T) {
	if got := canonEpic("", "TICK-02.1"); got != "E02" {
		t.Fatalf("canonEpic vazio com ticket: obtido %q, esperado E02", got)
	}
	if got := canonEpic("", "SEM-NUMERO"); got != "E01" {
		t.Fatalf("canonEpic sem número: obtido %q, esperado E01", got)
	}
	nodes := []routeNode{
		{ID: "A"},
		{ID: "B", DependsOn: []string{"A"}},
	}
	got := terminals(nodes)
	if !slices.Equal(got, []string{"B"}) {
		t.Fatalf("terminals: obtido %v, esperado [B]", got)
	}
}

func TestMutarAssembleResolveDepsEUncovered(t *testing.T) {
	tickets := []Ticket{
		{ID: "TICK-01.1", DependsOn: []string{"INEXISTENTE", "TICK-01.1"}, SpecCriteria: []string{"AC-1"}},
		{ID: "TICK-01.2", DependsOn: []string{"TICK-01.1"}},
	}
	warns := resolveDeps(tickets)
	if len(warns) != 2 {
		t.Fatalf("resolveDeps: esperados 2 avisos, obtido %d (%v)", len(warns), warns)
	}
	if len(tickets[0].DependsOn) != 0 {
		t.Fatalf("dependências inválidas não descartadas: %v", tickets[0].DependsOn)
	}
	if !slices.Equal(tickets[1].DependsOn, []string{"TICK-01.1"}) {
		t.Fatalf("dependência válida perdida: %v", tickets[1].DependsOn)
	}
	if cov := uncovered(tickets, []string{"AC-1"}); len(cov) != 0 {
		t.Fatalf("critério coberto reportado como descoberto: %v", cov)
	}
}

func TestMutarConcurrencyLimitesESkill(t *testing.T) {
	validos := []ConcurrencyStrategy{
		{Mode: "serial", MaxParallelSubagents: 1, HarnessCapability: ""},
		{Mode: "parallel_clean", MaxParallelSubagents: 8, HarnessCapability: "subagent_native"},
	}
	for _, c := range validos {
		if err := ValidateConcurrency(c); err != nil {
			t.Fatalf("estratégia válida rejeitada (%+v): %v", c, err)
		}
	}
	inv := ConcurrencyStrategy{Mode: "serial", MaxParallelSubagents: 2, HarnessCapability: "desconhecido"}
	if err := ValidateConcurrency(inv); err == nil {
		t.Fatal("esperado erro para harness_capability desconhecido")
	}
}

func TestMutarMetaPathsENomes(t *testing.T) {
	tmp := t.TempDir()
	proj := filepath.Join(tmp, "proj")
	ps := pathSet{
		baseDir:   proj,
		specPath:  filepath.Join(proj, "docs", "specs", "s.md"),
		planPath:  filepath.Join(proj, "docs", "plans", "p.md"),
		routePath: filepath.Join(proj, "docs", "routes", "r.json"),
	}
	if got := ps.rel(filepath.Join(tmp, "fora.md")); got != "fora.md" {
		t.Fatalf("rel fora da raiz: obtido %q, esperado fora.md", got)
	}
	if got := ps.link("docs/specs/s.md"); got != "../specs/s.md" {
		t.Fatalf("link relativo: obtido %q, esperado ../specs/s.md", got)
	}
	specComIntent := "# 📋 SPEC: Tema\n> **Referência da Intenção:** docs/intents/i.md\n"
	mCom := newMeta(ps, specComIntent, "2026-11-21", SupervisionEpic)
	if mCom.IntentLink != "../intents/i.md" {
		t.Fatalf("IntentLink com intent: obtido %q", mCom.IntentLink)
	}
	mSem := newMeta(ps, "# 📋 SPEC: Tema\n", "2026-11-21", SupervisionEpic)
	if mSem.IntentLink != "" {
		t.Fatalf("IntentLink sem intent deve ser vazio, obtido %q", mSem.IntentLink)
	}
	mZero := planMeta{Epics: map[int]string{0: "Zero"}}
	if got := mZero.epicName("SEM-NUM"); got != "Épico SEM-NUM" {
		t.Fatalf("epicName com 0: obtido %q", got)
	}
}

func TestMutarPlanMDHeaderEpicsEContagem(t *testing.T) {
	var b1, b2 strings.Builder
	writePlanHeader(&b1, planMeta{Theme: "T", IntentRef: "docs/intents/i.md", IntentLink: "../intents/i.md"})
	if !strings.Contains(b1.String(), "[docs/intents/i.md](../intents/i.md)") {
		t.Fatalf("header com intent inesperado: %s", b1.String())
	}
	writePlanHeader(&b2, planMeta{Theme: "T"})
	if !strings.Contains(b2.String(), "Não informada na SPEC") {
		t.Fatalf("header sem intent inesperado: %s", b2.String())
	}
	nodes := []routeNode{
		{ID: "T1", Type: tipoUnidade, epic: "E01", Nature: "nano"},
		{ID: "T2", Type: tipoUnidade, epic: "E01", Nature: "micro"},
		{ID: "SG1", Type: "stop_gate", epic: "E01"},
		{ID: "T3", Type: tipoUnidade, epic: "E02", Nature: "nano"},
	}
	if got := epicsOf(nodes); !slices.Equal(got, []string{"E01", "E02"}) {
		t.Fatalf("epicsOf: obtido %v, esperado [E01 E02]", got)
	}
	total, _, nature := epicSummary(nodes, "E01")
	if total != 2 || nature != "micro" {
		t.Fatalf("epicSummary(E01): obtido total=%d nature=%q, esperado 2 e micro", total, nature)
	}
	if got := verifyText(nil); !strings.Contains(got, "nenhuma") {
		t.Fatalf("verifyText(nil): obtido %q", got)
	}
	if got := verifyText([][]string{{"go", "test"}}); got != "`go test`" {
		t.Fatalf("verifyText comando: obtido %q", got)
	}
	if got := countTickets(nodes); got != 3 {
		t.Fatalf("countTickets: obtido %d, esperado 3", got)
	}
}
