package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/roteamento"
)

const planoRota = `# PLAN

### [TICK-01.1] Primeiro
- **Natureza:** nano
- **Arquivos Afetados:** ` + "`internal/a/a.go`" + `
- **Dependências:** nenhuma
- **Comando de Validação:** go test ./internal/a/...

### [TICK-01.2] Segundo
- **Natureza:** micro
- **Arquivos Afetados:** ` + "`internal/a/b.go`, `internal/a/b_test.go`" + `
- **Dependências:** TICK-01.1
- **Comando de Validação:** go test ./internal/a/...

### [TICK-01.3] Terceiro
- **Natureza:** micro
- **Arquivos Afetados:** ` + "`internal/a/c.go`, `internal/b/c.go`, `internal/c/c.go`, `internal/d/c.go`, `internal/e/c.go`" + `
- **Dependências:** TICK-01.2
- **Comando de Validação:** go test ./...
`

func itemRota(harness, modelo string) roteamento.ItemInv {
	return roteamento.ItemInv{
		Par:     roteamento.Par{Harness: harness, Via: roteamento.ViaOficial},
		Binario: "/bin/" + harness,
		Modelos: []string{modelo},
		Estado:  roteamento.EstadoHomologado,
	}
}

func rodarRota(t *testing.T, inv []roteamento.ItemInv, err error) (int, string, map[string]any) {
	t.Helper()
	dir := t.TempDir()
	plano := filepath.Join(dir, "2026-11-12_PLAN_TESTE.md")
	saida := filepath.Join(dir, "route.json")
	if e := os.WriteFile(plano, []byte(planoRota), 0o644); e != nil {
		t.Fatal(e)
	}
	var out, errb bytes.Buffer
	deps := Deps{Stdout: &out, Stderr: &errb, Rota: &RotaDeps{
		Inventario: func(context.Context) ([]roteamento.ItemInv, error) { return inv, err },
	}}
	code := Run(context.Background(), []string{"rota", "--plan", plano, "--out", saida}, deps)
	dados, e := os.ReadFile(saida)
	if e != nil {
		return code, errb.String(), nil
	}
	var doc map[string]any
	if e := json.Unmarshal(dados, &doc); e != nil {
		t.Fatal(e)
	}
	return code, errb.String(), doc
}

func TestRotaGravaRoteamentoResolvido(t *testing.T) {
	inv := []roteamento.ItemInv{itemRota("claude", "haiku"), itemRota("claude", "sonnet")}
	code, stderr, doc := rodarRota(t, inv, nil)
	if code != 0 {
		t.Fatalf("exit = %d: %s", code, stderr)
	}
	if doc["name"] != "TESTE" || doc["date"] != "2026-11-12" {
		t.Fatalf("nome e data inesperados: %v %v", doc["name"], doc["date"])
	}
	nos := doc["nodes"].([]any)
	if len(nos) != 3 {
		t.Fatalf("nós = %d, esperado 3", len(nos))
	}
	esperado := []struct{ classe, modelo, esforco string }{
		{"nano", "haiku", "low"}, {"micro", "sonnet", "medium"}, {roteamento.ClasseGrande, "sonnet", "high"},
	}
	for i, n := range nos {
		no := n.(map[string]any)
		r := no["routing"].(map[string]any)
		if no["nature"] != esperado[i].classe || no["class_rule"] == "" {
			t.Errorf("nó %d: classe %v regra %v", i, no["nature"], no["class_rule"])
		}
		if r["harness"] != "claude" || r["via"] != "oficial" || r["model"] != esperado[i].modelo || r["effort"] != esperado[i].esforco {
			t.Errorf("nó %d: routing %v", i, r)
		}
	}
}

func TestRotaInventarioVazioSaiComUm(t *testing.T) {
	code, _, doc := rodarRota(t, nil, roteamento.ErrInventarioVazio)
	if code != 1 || doc != nil {
		t.Fatalf("exit = %d, ROUTE = %v", code, doc)
	}
}

func TestRotaSemHomologadoFalha(t *testing.T) {
	inv := []roteamento.ItemInv{itemRota("claude", "haiku")}
	inv[0].Estado = roteamento.EstadoPendente
	code, stderr, doc := rodarRota(t, inv, nil)
	if code != 1 || doc != nil || stderr == "" {
		t.Fatalf("exit = %d, stderr = %q, ROUTE = %v", code, stderr, doc)
	}
}

func TestRotaUsoInvalidoSaiComDois(t *testing.T) {
	var out, errb bytes.Buffer
	if code := Run(context.Background(), []string{"rota"}, Deps{Stdout: &out, Stderr: &errb}); code != 2 {
		t.Fatalf("exit = %d, esperado 2", code)
	}
}
