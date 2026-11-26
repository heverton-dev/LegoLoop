package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/cli"
	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

func semearEstiloCLI(t *testing.T) (string, string) {
	t.Helper()
	raiz := t.TempDir()
	s := inventario.Store{Raiz: raiz}
	e := inventario.Estilo{
		SchemaVersion: 1, ID: "molde-de-pilula", Nome: "Molde", Status: "rascunho", Versao: "0.1.0",
		Renderizador: "molde-pilula", PromptMestre: "diorama", Paleta: []string{"#000000"},
		Referencias: []inventario.Referencia{{Tipo: "imagem", Origem: "ref.png"}},
	}
	if err := s.NovoEstilo(e); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(inventario.Aprovacao{
		Estilo: "molde-de-pilula", AprovadoPor: "usuario", PersonagensTeste: []string{"rei"},
	})
	if err != nil {
		t.Fatal(err)
	}
	arq := filepath.Join(raiz, "aprovacao.json")
	if err := os.WriteFile(arq, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return raiz, arq
}

func rodarVideoInventario(args ...string) (int, string, string) {
	var out, errb bytes.Buffer
	d := cli.Deps{Stdout: &out, Stderr: &errb}
	code := cli.Run(context.Background(), append([]string{"video-inventario"}, args...), d)
	return code, out.String(), errb.String()
}

func TestVideoInventarioHelp(t *testing.T) {
	code, out, _ := rodarVideoInventario("--help")
	if code != 0 || !strings.Contains(out, "validar") || !strings.Contains(out, "verificar") {
		t.Fatalf("help: code=%d out=%q", code, out)
	}
	if code, _, _ := rodarVideoInventario("desconhecido"); code != 2 {
		t.Fatalf("subcomando desconhecido: code=%d, esperado 2", code)
	}
}

func TestVideoInventarioValidar(t *testing.T) {
	raiz, arq := semearEstiloCLI(t)
	code, out, errb := rodarVideoInventario("validar", "--raiz", raiz, "--estilo", "molde-de-pilula", "--aprovacao", arq, "--json")
	if code != 0 {
		t.Fatalf("validar: code=%d out=%q err=%q", code, out, errb)
	}
	if !strings.Contains(out, `"status": "validado"`) {
		t.Fatalf("saida sem status validado: %s", out)
	}
	e, err := inventario.Store{Raiz: raiz}.LerEstilo("molde-de-pilula")
	if err != nil || e.Status != "validado" {
		t.Fatalf("estilo gravado: status=%q err=%v", e.Status, err)
	}
}

func TestVideoInventarioValidarSemAprovacao(t *testing.T) {
	raiz, _ := semearEstiloCLI(t)
	ausente := filepath.Join(raiz, "nao-existe.json")
	code, out, _ := rodarVideoInventario("validar", "--raiz", raiz, "--estilo", "molde-de-pilula", "--aprovacao", ausente, "--json")
	if code != 1 || !strings.Contains(out, "aprovacao") {
		t.Fatalf("sem aprovacao: code=%d out=%q", code, out)
	}
	e, err := inventario.Store{Raiz: raiz}.LerEstilo("molde-de-pilula")
	if err != nil || e.Status != "rascunho" {
		t.Fatalf("status=%q err=%v, esperado rascunho", e.Status, err)
	}
	if code, _, _ := rodarVideoInventario("validar", "--raiz", raiz); code != 2 {
		t.Fatalf("sem --estilo: code=%d, esperado 2", code)
	}
}

func TestVideoInventarioVerificar(t *testing.T) {
	raiz, _ := semearEstiloCLI(t)
	code, out, errb := rodarVideoInventario("verificar", "--raiz", raiz, "--estilo", "molde-de-pilula", "--json")
	if code != 0 || !strings.Contains(out, `"ok": true`) {
		t.Fatalf("verificar: code=%d out=%q err=%q", code, out, errb)
	}
	s := inventario.Store{Raiz: raiz}
	it := inventario.Item{SchemaVersion: 1, ID: "rei", Tipo: "personagem", Estilo: "molde-de-pilula", Arquivo: "falta.png", Hash: "abc"}
	if err := s.GravarItem(it); err != nil {
		t.Fatal(err)
	}
	code, out, _ = rodarVideoInventario("verificar", "--raiz", raiz, "--estilo", "molde-de-pilula", "--json")
	if code != 1 || !strings.Contains(out, "hash") {
		t.Fatalf("hash divergente: code=%d out=%q", code, out)
	}
	if code, _, _ := rodarVideoInventario("verificar", "--raiz", raiz, "--estilo", "inexistente"); code != 1 {
		t.Fatalf("estilo inexistente: code=%d, esperado 1", code)
	}
}
