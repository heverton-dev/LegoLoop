package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/routedag"
)

const specDoFabricarPlano = `---
id: SPEC-TESTE
---

# SPEC: Sistema de Teste

> **Referência da Intenção:** [docs/intents/2026-11-21_INTENT_TESTE.md](../intents/2026-11-21_INTENT_TESTE.md)  

## 5. Critérios

* [ ] ` + "`AC-E01.1`" + `: **primeiro critério**
* [ ] ` + "`AC-E01.2`" + `: **segundo critério**
`

const ticketsDoFabricarPlano = `[
  {"id":"TICK-01.1","title":"Criar Peça um","epic_id":"E01","nature":"micro","depends_on":[],"target_files":["mod1.go"],"spec_criteria":["AC-E01.1"],"validation_cmd":"go test ./mod1/..."},
  {"id":"TICK-01.2","title":"Criar Peça dois","epic_id":"E01","nature":"nano","depends_on":["TICK-01.1"],"target_files":["mod2.go"],"spec_criteria":["AC-E01.2"],"validation_cmd":"go test ./mod2/..."}
]`

func prepararAmbienteFabricarPlano(t *testing.T, tickets string) (raiz, spec, ticketsFile string) {
	t.Helper()
	raiz = t.TempDir()
	spec = filepath.Join(raiz, "docs", "specs", "2026-11-21_SPEC_TESTE.md")
	ticketsFile = filepath.Join(raiz, "tickets.json")
	intencao := filepath.Join(raiz, "docs", "intents", "2026-11-21_INTENT_TESTE.md")
	arquivos := map[string]string{
		spec:        specDoFabricarPlano,
		ticketsFile: tickets,
		intencao:    "---\nid: INTENT-TESTE\nstatus: APPROVED\n---\n# INTENT\n",
	}
	for arq, conteudo := range arquivos {
		if err := os.MkdirAll(filepath.Dir(arq), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(arq, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return raiz, spec, ticketsFile
}

func rodarFabricar(args ...string) (int, string) {
	var out, errb bytes.Buffer
	code := cmdFabricar(context.Background(), args, Deps{Stdout: &out, Stderr: &errb})
	return code, out.String() + errb.String()
}

func TestFabricarPlano(t *testing.T) {
	raiz, spec, ticketsFile := prepararAmbienteFabricarPlano(t, ticketsDoFabricarPlano)
	args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21"}

	code, saida := rodarFabricar(args...)
	if code != 0 {
		t.Fatalf("fabricar plano falhou com código %d: %s", code, saida)
	}

	planPath := filepath.Join(raiz, "docs", "plans", "2026-11-21_PLAN_TESTE.md")
	routePath := filepath.Join(raiz, "docs", "routes", "2026-11-21_ROUTE_TESTE.json")
	dagPath := filepath.Join(raiz, "docs", "dags", "2026-11-21_DAG_TESTE.json")

	for _, p := range []string{planPath, routePath, dagPath} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("arquivo esperado não foi gerado: %s (%v)", p, err)
		}
	}

	rawRoute, err := os.ReadFile(routePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := routedag.VerificarCadeia(rawRoute, raiz); err != nil {
		t.Fatalf("routedag.VerificarCadeia falhou na ROUTE gerada: %v", err)
	}

	rawDAG, err := os.ReadFile(dagPath)
	if err != nil {
		t.Fatal(err)
	}
	var dagObj map[string]any
	if err := json.Unmarshal(rawDAG, &dagObj); err != nil {
		t.Fatalf("DAG gerado não é JSON válido: %v", err)
	}
	if nodes, ok := dagObj["nodes"].([]any); !ok || len(nodes) < 2 {
		t.Fatalf("DAG gerado não possui nodes esperados: %v", dagObj["nodes"])
	}
}

func TestFabricarPlanoJSON(t *testing.T) {
	raiz, spec, ticketsFile := prepararAmbienteFabricarPlano(t, ticketsDoFabricarPlano)
	args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21", "--json"}

	var out, errb bytes.Buffer
	code := cmdFabricar(context.Background(), args, Deps{Stdout: &out, Stderr: &errb})
	if code != 0 {
		t.Fatalf("fabricar plano --json falhou: %d: %s", code, errb.String())
	}
	var res fabricarPlanoResult
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatalf("saída não é JSON do resultado: %v\nsaída=%s", err, out.String())
	}
	if !strings.Contains(res.PlanMarkdownPath, "PLAN_TESTE") ||
		!strings.Contains(res.RouteJSONPath, "ROUTE_TESTE") ||
		!strings.Contains(res.DagJSONPath, "DAG_TESTE") {
		t.Fatalf("caminhos inesperados no resultado JSON: %+v", res)
	}
}

func TestFabricarPlanoRecusa(t *testing.T) {
	raiz, spec, ticketsFile := prepararAmbienteFabricarPlano(t, ticketsDoFabricarPlano)
	intencao := filepath.Join(raiz, "docs", "intents", "2026-11-21_INTENT_TESTE.md")

	t.Run("status under review", func(t *testing.T) {
		if err := os.WriteFile(intencao, []byte("---\nid: INTENT-TESTE\nstatus: UNDER_REVIEW\n---\n# INTENT\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21"}
		code, saida := rodarFabricar(args...)
		if code != 2 {
			t.Fatalf("esperava código 2, obteve %d: %s", code, saida)
		}
		if !strings.Contains(saida, "INTENT") {
			t.Fatalf("saída esperada citando INTENT, obteve: %s", saida)
		}
		if !strings.Contains(saida, "legoloop-concepcao") {
			t.Fatalf("saída esperada citando legoloop-concepcao, obteve: %s", saida)
		}
		for _, pasta := range []string{"plans", "routes", "dags"} {
			if _, err := os.Stat(filepath.Join(raiz, "docs", pasta)); !os.IsNotExist(err) {
				t.Fatalf("pasta docs/%s não deveria ter sido criada na recusa", pasta)
			}
		}
	})

	t.Run("intent inexistente", func(t *testing.T) {
		if err := os.Remove(intencao); err != nil {
			t.Fatal(err)
		}
		args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21"}
		code, saida := rodarFabricar(args...)
		if code != 2 {
			t.Fatalf("esperava código 2, obteve %d: %s", code, saida)
		}
		if !strings.Contains(saida, "INTENT") {
			t.Fatalf("saída esperada citando INTENT, obteve: %s", saida)
		}
		if !strings.Contains(saida, "legoloop-concepcao") {
			t.Fatalf("saída esperada citando legoloop-concepcao, obteve: %s", saida)
		}
		for _, pasta := range []string{"plans", "routes", "dags"} {
			if _, err := os.Stat(filepath.Join(raiz, "docs", pasta)); !os.IsNotExist(err) {
				t.Fatalf("pasta docs/%s não deveria ter sido criada na recusa", pasta)
			}
		}
	})
}

func TestFabricarPlanoEstrito(t *testing.T) {
	ticketsIncompletos := `[
  {"id":"TICK-01.1","title":"Criar Peça um","epic_id":"E01","nature":"micro","depends_on":[],"target_files":["mod1.go"],"spec_criteria":["AC-E01.1"],"validation_cmd":"go test ./mod1/..."}
]`

	t.Run("sem flag falha com 1 e nao grava arquivos", func(t *testing.T) {
		raiz, spec, ticketsFile := prepararAmbienteFabricarPlano(t, ticketsIncompletos)
		args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21"}

		code, saida := rodarFabricar(args...)
		if code != 1 {
			t.Fatalf("esperava código 1, obteve %d: %s", code, saida)
		}
		if !strings.Contains(saida, "AC-E01.2") {
			t.Fatalf("esperava aviso citando AC-E01.2, obteve: %s", saida)
		}

		planPath := filepath.Join(raiz, "docs", "plans", "2026-11-21_PLAN_TESTE.md")
		routePath := filepath.Join(raiz, "docs", "routes", "2026-11-21_ROUTE_TESTE.json")
		dagPath := filepath.Join(raiz, "docs", "dags", "2026-11-21_DAG_TESTE.json")

		for _, p := range []string{planPath, routePath, dagPath} {
			if _, err := os.Stat(p); !os.IsNotExist(err) {
				t.Fatalf("arquivo %s não deveria existir após reprovação estrita", p)
			}
		}
	})

	t.Run("com permissivo grava 3 artefatos emite aviso e sai 0", func(t *testing.T) {
		raiz, spec, ticketsFile := prepararAmbienteFabricarPlano(t, ticketsIncompletos)
		args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21", "--permissivo"}

		var out, errb bytes.Buffer
		code := cmdFabricar(context.Background(), args, Deps{Stdout: &out, Stderr: &errb})
		if code != 0 {
			t.Fatalf("esperava código 0 com --permissivo, obteve %d: %s", code, errb.String())
		}
		if !strings.Contains(errb.String(), "AC-E01.2") {
			t.Fatalf("esperava aviso em stderr citando AC-E01.2, obteve: %s", errb.String())
		}

		planPath := filepath.Join(raiz, "docs", "plans", "2026-11-21_PLAN_TESTE.md")
		routePath := filepath.Join(raiz, "docs", "routes", "2026-11-21_ROUTE_TESTE.json")
		dagPath := filepath.Join(raiz, "docs", "dags", "2026-11-21_DAG_TESTE.json")

		for _, p := range []string{planPath, routePath, dagPath} {
			if _, err := os.Stat(p); err != nil {
				t.Fatalf("arquivo esperado não foi gerado com --permissivo: %s (%v)", p, err)
			}
		}
	})
}

func TestFabricarPlanoAtomico(t *testing.T) {
	t.Run("ciclo de dependencias entre tickets", testFabricarPlanoAtomicoCiclo)
	t.Run("falha na etapa do dag remove artefatos gravados", testFabricarPlanoAtomicoFalhaDAG)
}

func testFabricarPlanoAtomicoCiclo(t *testing.T) {
	ticketsCiclo := `[
  {"id":"TICK-01.1","title":"Peça 1","epic_id":"E01","nature":"micro","depends_on":["TICK-01.2"],"target_files":["m1.go"],"spec_criteria":["AC-E01.1"],"validation_cmd":"go test ./..."},
  {"id":"TICK-01.2","title":"Peça 2","epic_id":"E01","nature":"nano","depends_on":["TICK-01.1"],"target_files":["m2.go"],"spec_criteria":["AC-E01.2"],"validation_cmd":"go test ./..."}
]`
	raiz, spec, ticketsFile := prepararAmbienteFabricarPlano(t, ticketsCiclo)
	args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21"}

	code, saida := rodarFabricar(args...)
	if code != 1 {
		t.Fatalf("esperava código 1, obteve %d: %s", code, saida)
	}
	verificarArtefatosInexistentes(t, raiz)
}

func testFabricarPlanoAtomicoFalhaDAG(t *testing.T) {
	raiz, spec, ticketsFile := prepararAmbienteFabricarPlano(t, ticketsDoFabricarPlano)
	dagsDir := filepath.Join(raiz, "docs", "dags")
	if err := os.WriteFile(dagsDir, []byte("bloqueio"), 0o644); err != nil {
		t.Fatal(err)
	}
	args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21"}

	code, saida := rodarFabricar(args...)
	if code != 1 {
		t.Fatalf("esperava código 1, obteve %d: %s", code, saida)
	}
	verificarArtefatosInexistentes(t, raiz)
}

func verificarArtefatosInexistentes(t *testing.T, raiz string) {
	t.Helper()
	planPath := filepath.Join(raiz, "docs", "plans", "2026-11-21_PLAN_TESTE.md")
	routePath := filepath.Join(raiz, "docs", "routes", "2026-11-21_ROUTE_TESTE.json")
	dagPath := filepath.Join(raiz, "docs", "dags", "2026-11-21_DAG_TESTE.json")
	progDir := filepath.Join(raiz, "docs", "progress", "2026-11-21_plan-teste")

	for _, p := range []string{planPath, routePath, dagPath, progDir} {
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() || (err == nil && p == progDir) {
			t.Fatalf("artefato não deveria existir no disco: %s", p)
		}
	}
}

func TestLegoloopPlanoRemovido(t *testing.T) {
	var out, errb bytes.Buffer
	code := Run(context.Background(), []string{"plano"}, Deps{Stdout: &out, Stderr: &errb})
	if code != 2 {
		t.Fatalf("esperava exit code 2 para comando plano removido, obteve %d", code)
	}
	if !strings.Contains(errb.String(), `comando desconhecido "plano"`) {
		t.Fatalf("mensagem inesperada de erro: %s", errb.String())
	}
	verificarAusenciaLegoloopPlano(t)
}

func verificarAusenciaLegoloopPlano(t *testing.T) {
	t.Helper()
	root := repoRoot(t)
	alvos := []string{
		filepath.Join(root, "AGENTS.md"),
		filepath.Join(root, "README.md"),
		filepath.Join(root, "docs", "prompts"),
		filepath.Join(root, "internal"),
	}
	for _, alvo := range alvos {
		varrerLegoloopPlano(t, alvo)
	}
}

func varrerLegoloopPlano(t *testing.T, caminho string) {
	t.Helper()
	termo := "legoloop" + " " + "plano"
	err := filepath.Walk(caminho, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		conteudo, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(conteudo), termo) {
			t.Errorf("encontrada referencia proibida em %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("falha ao varrer %s: %v", caminho, err)
	}
}
