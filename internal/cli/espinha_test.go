package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/espinha"
)

func TestCmdEspinhaVerificarRepositorioReal(t *testing.T) {
	raiz := filepath.Join("..", "..")
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), []string{"espinha", "verificar", "--raiz", raiz, "--json"}, Deps{
		Stdout: &out,
		Stderr: &errBuf,
	})
	var vs []espinha.Violacao
	if err := json.Unmarshal(out.Bytes(), &vs); err != nil {
		t.Fatalf("JSON invalido: %v\nsaida: %s", err, out.String())
	}
	for _, v := range vs {
		if v.Regra != "processo-incompleto" && v.Regra != "comando-exclusivo-de-so" {
			t.Errorf("violacao inesperada no repo real: %s", v.String())
		}
	}
	if len(vs) > 0 && code != 2 {
		t.Fatalf("esperado exit 2 com dividas conhecidas E06/E07, obtido %d", code)
	}
	if len(vs) == 0 && code != 0 {
		t.Fatalf("esperado exit 0 com 0 violacoes, obtido %d", code)
	}
}

func TestCmdEspinhaVerificarComViolacaoTextoEJSON(t *testing.T) {
	raiz := t.TempDir()
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), []string{"espinha", "verificar", "--raiz", raiz}, Deps{
		Stdout: &out,
		Stderr: &errBuf,
	})
	if code != 2 {
		t.Fatalf("esperado exit 2 com violacao, obtido %d", code)
	}
	if !strings.Contains(out.String(), "regra: registro-invalido; alvo: espinha/registro.json; motivo: ") {
		t.Fatalf("formato de saida inesperado: %q", out.String())
	}

	out.Reset()
	errBuf.Reset()
	code = Run(context.Background(), []string{"espinha", "verificar", "--raiz", raiz, "--json"}, Deps{
		Stdout: &out,
		Stderr: &errBuf,
	})
	if code != 2 {
		t.Fatalf("esperado exit 2 com --json, obtido %d", code)
	}
	var vs []espinha.Violacao
	if err := json.Unmarshal(out.Bytes(), &vs); err != nil {
		t.Fatalf("JSON invalido: %v\nsaida: %s", err, out.String())
	}
	if len(vs) == 0 || vs[0].Regra != "registro-invalido" {
		t.Fatalf("violacoes JSON inesperadas: %+v", vs)
	}
}
