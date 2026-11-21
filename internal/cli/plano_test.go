package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const specDoPlano = "---\nid: SPEC-X\n---\n\n# SPEC: Tema\n\n> **Referência da Intenção:** [docs/intents/2026-11-21_INTENT_TEMA.md](../intents/2026-11-21_INTENT_TEMA.md)  \n\n## 5. Critérios\n\n* [ ] `AC-E01.1`: **um**\n* [ ] `AC-E01.2`: **dois**\n"

const ticketsDoPlano = `[
 {"id":"TICK-01.1","title":"Criar a peça um","epic_id":"E01","nature":"micro","depends_on":[],"target_files":["a.go"],"spec_criteria":["AC-E01.1"],"validation_cmd":"go test ./a/..."},
 {"id":"TICK-01.2","title":"Criar a peça dois","epic_id":"E01","nature":"nano","depends_on":["TICK-01.1"],"target_files":["b.go"],"spec_criteria":["AC-E01.2"],"validation_cmd":"go test ./b/..."}
]`

func prepararPlano(t *testing.T, tickets string) (raiz, spec, js string) {
	t.Helper()
	raiz = t.TempDir()
	spec = filepath.Join(raiz, "docs", "specs", "2026-11-21_SPEC_TEMA.md")
	js = filepath.Join(raiz, "tickets.json")
	intencao := filepath.Join(raiz, "docs", "intents", "2026-11-21_INTENT_TEMA.md")
	for arq, txt := range map[string]string{spec: specDoPlano, js: tickets, intencao: "---\nid: INTENT-X\nstatus: APPROVED\n---\n# INTENT\n"} {
		if err := os.MkdirAll(filepath.Dir(arq), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(arq, []byte(txt), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return raiz, spec, js
}

func rodarPlano(args ...string) (int, string) {
	var out, errb bytes.Buffer
	code := cmdPlano(context.Background(), args, Deps{Stdout: &out, Stderr: &errb})
	return code, out.String() + errb.String()
}

func TestPlanoRenderizaPlanERouteDeterministicos(t *testing.T) {
	raiz, spec, js := prepararPlano(t, ticketsDoPlano)
	args := []string{"--raiz", raiz, "--spec", spec, "--tickets", js, "--data", "2026-11-21"}
	if code, saida := rodarPlano(args...); code != 0 {
		t.Fatalf("plano saiu %d: %s", code, saida)
	}
	plano := filepath.Join(raiz, "docs", "plans", "2026-11-21_PLAN_TEMA.md")
	rota := filepath.Join(raiz, "docs", "routes", "2026-11-21_ROUTE_TEMA.json")
	p1, r1 := readFile(t, plano), readFile(t, rota)
	if !strings.Contains(p1, "TICK-01.2") || !strings.Contains(r1, "TICK-01.1") {
		t.Fatalf("PLAN ou ROUTE sem os Tickets.\nplan=%s", p1)
	}
	if code, _ := rodarPlano(args...); code != 0 || readFile(t, plano) != p1 || readFile(t, rota) != r1 {
		t.Fatal("a mesma entrada deve gerar os mesmos bytes")
	}
}

func TestPlanoAvisaCriterioSemTicketEEstritoReprova(t *testing.T) {
	soUm := `[{"id":"TICK-01.1","title":"Criar um","epic_id":"E01","nature":"micro","spec_criteria":["AC-E01.1"]}]`
	raiz, spec, js := prepararPlano(t, soUm)
	base := []string{"--raiz", raiz, "--spec", spec, "--tickets", js, "--data", "2026-11-21"}
	code, saida := rodarPlano(base...)
	if code != 0 || !strings.Contains(saida, "AC-E01.2") {
		t.Fatalf("sem --estrito: exit %d e aviso esperado sobre AC-E01.2, veio: %s", code, saida)
	}
	if code, _ := rodarPlano(append(base, "--estrito")...); code != 1 {
		t.Fatalf("com --estrito o critério sem Ticket deve sair com 1, saiu %d", code)
	}
}

func TestPlanoExigeSpecETicketsEAceitaSoJSONValido(t *testing.T) {
	if code, _ := rodarPlano("--spec", "x.md"); code != 2 {
		t.Fatalf("sem --tickets deveria sair com 2, saiu %d", code)
	}
	raiz, spec, js := prepararPlano(t, "{nao e json")
	if code, _ := rodarPlano("--raiz", raiz, "--spec", spec, "--tickets", js); code != 1 {
		t.Fatalf("JSON inválido deveria sair com 1, saiu %d", code)
	}
}

func TestPlanoRecusaSpecSemIntencaoAprovada(t *testing.T) {
	raiz, spec, js := prepararPlano(t, ticketsDoPlano)
	args := []string{"--raiz", raiz, "--spec", spec, "--tickets", js, "--data", "2026-11-21"}
	intencao := filepath.Join(raiz, "docs", "intents", "2026-11-21_INTENT_TEMA.md")
	if err := os.WriteFile(intencao, []byte("---\nid: INTENT-X\nstatus: UNDER_REVIEW\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, saida := rodarPlano(args...); code != 2 || !strings.Contains(saida, "INTENT") {
		t.Fatalf("INTENT não aprovada deveria sair com 2 e citar a INTENT, saiu %d: %s", code, saida)
	}
	if err := os.Remove(intencao); err != nil {
		t.Fatal(err)
	}
	if code, _ := rodarPlano(args...); code != 2 {
		t.Fatalf("INTENT ausente deveria sair com 2, saiu %d", code)
	}
	if _, err := os.Stat(filepath.Join(raiz, "docs", "plans")); err == nil {
		t.Fatal("a recusa não pode deixar PLAN nem ROUTE no disco")
	}
}
