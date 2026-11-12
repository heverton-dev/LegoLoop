package orchestrix

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
)

const (
	chaveFalsa = "chave-falsa-do-teste-9r"
	baseFalsa  = "https://roteador.exemplo.test"
)

// orchestrixGateway monta o despachante com a chave falsa. O helper grava em dump o ambiente que recebeu.
func orchestrixGateway(dump string) *Orchestrix {
	o := unitOrchestrix("")
	o.RouterKey = chaveFalsa
	o.RouterKeyURL = baseFalsa
	o.Inventory.RouterURL = baseFalsa
	o.ExtraEnv = append(o.ExtraEnv, "LEGOLOOP_UNIT_ENVDUMP="+dump)
	return o
}

func noClaudeVia(id, via, model string) dag.Node {
	n := unitNode(id, "out/"+id+".json", "exit_code == 0", nil)
	n.ContextOut = nil
	n.Routing = dag.RoutingConfig{Mode: dag.ModeManual, Harness: "claude", Via: via, Model: model}
	return n
}

func TestGatewayPorNoSoNoFilho(t *testing.T) {
	t.Setenv("ANTHROPIC_BASE_URL", "")
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "")
	repo := newRepoT(t)
	dir := t.TempDir()
	dumpOficial := filepath.Join(dir, "oficial.txt")
	dumpRouter := filepath.Join(dir, "router.txt")

	resOficial, err := orchestrixGateway(dumpOficial).Execute(context.Background(), noClaudeVia("of", "oficial", "sonnet"), repo)
	if err != nil {
		t.Fatal(err)
	}
	resRouter, err := orchestrixGateway(dumpRouter).Execute(context.Background(), noClaudeVia("rt", "9router", "code-pro"), repo)
	if err != nil {
		t.Fatal(err)
	}

	if got := lerArquivo(t, dumpOficial); got != "BASE=\nTOKEN=\n" {
		t.Fatalf("filho oficial recebeu variáveis do gateway: %q", got)
	}
	if got := lerArquivo(t, dumpRouter); got != "BASE="+baseFalsa+"\nTOKEN="+chaveFalsa+"\n" {
		t.Fatalf("filho 9router sem as variáveis do gateway: %q", got)
	}
	if os.Getenv("ANTHROPIC_BASE_URL") != "" || os.Getenv("ANTHROPIC_AUTH_TOKEN") != "" {
		t.Fatal("o ambiente do processo pai mudou")
	}
	if !strings.HasSuffix(resRouter.Model, "/code-pro") || !strings.HasSuffix(resOficial.Model, "/sonnet") {
		t.Fatalf("modelos: oficial=%q router=%q", resOficial.Model, resRouter.Model)
	}
	semChaveEmSaida(t, repo, resOficial.Notes, resRouter.Notes)
}

// semChaveEmSaida confere que a chave não aparece em notas nem em arquivo algum do repositório.
func semChaveEmSaida(t *testing.T, repo string, notas ...[]string) {
	t.Helper()
	for _, ns := range notas {
		if strings.Contains(strings.Join(ns, "\n"), chaveFalsa) {
			t.Fatalf("a chave apareceu nas notas: %v", ns)
		}
	}
	err := filepath.WalkDir(repo, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if d != nil && d.IsDir() && d.Name() == ".git" {
				return filepath.SkipDir
			}
			return err
		}
		if b, rerr := os.ReadFile(p); rerr == nil && strings.Contains(string(b), chaveFalsa) {
			t.Errorf("a chave apareceu em %s", p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func lerArquivo(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestGatewayRecusaViaRouterSemBaseConfiavel(t *testing.T) {
	o := orchestrixGateway(filepath.Join(t.TempDir(), "d.txt"))
	o.RouterKeyURL = "https://outra-base.exemplo.test"
	_, err := o.Execute(context.Background(), noClaudeVia("rt", "9router", "code-pro"), newRepoT(t))
	if err == nil || !strings.Contains(err.Error(), "base confiável") {
		t.Fatalf("err = %v", err)
	}
}

func TestGatewayNaoVazaChaveNoErro(t *testing.T) {
	o := &Orchestrix{RouterKey: chaveFalsa}
	err := o.semChave(errors.New("falhou com " + chaveFalsa + " no texto"))
	if err == nil || strings.Contains(err.Error(), chaveFalsa) {
		t.Fatalf("err = %v", err)
	}
	if o.semChave(nil) != nil {
		t.Fatal("erro nulo deveria continuar nulo")
	}
}
