package processo_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/cli"
	"legoloop/internal/gate"
	"legoloop/internal/paths"
	"legoloop/internal/triagem"
)

// excecoes5Por5 isenta processos do nascimento 5/5, com dono e motivo (AC-E02.3).
var excecoes5Por5 = map[string]string{}

var manifestosComSchema = map[string]string{
	"triagem.json":         "triagem.v1.json",
	triagem.ArquivoMemoria: "memoria_execucao.v1.json",
	"anticorpos.json":      "anticorpos.v1.json",
}

// TestProcessosCincoPorCinco cobre o AC-E02.3: todo Processo gravado nos diretórios
// antigos tem os 5 manifestos do nascimento válidos contra os schemas oficiais e passa
// em legoloop validate, ou é exceção documentada com dono.
func TestProcessosCincoPorCinco(t *testing.T) {
	raiz := filepath.Join("..", "..", paths.DirAntigo)
	entradas, err := os.ReadDir(raiz)
	if err != nil {
		t.Fatalf("ler %s: %v", raiz, err)
	}
	vistos := 0
	for _, e := range entradas {
		if !e.IsDir() {
			continue
		}
		if motivo := excecoes5Por5[e.Name()]; motivo != "" {
			t.Logf("exceção documentada %s: %s", e.Name(), motivo)
			continue
		}
		vistos++
		dir := filepath.Join(raiz, e.Name())
		manifestosValidos(t, dir, e.Name())
		validatePassa(t, dir)
	}
	if vistos == 0 {
		t.Fatal("nenhum Processo encontrado para conferir")
	}
}

// manifestosValidos exige os manifestos do nascimento 5/5 válidos contra os schemas de schemas/.
func manifestosValidos(t *testing.T, dir, nome string) {
	t.Helper()
	for manifesto, esquema := range manifestosComSchema {
		b, err := os.ReadFile(filepath.Join(dir, manifesto))
		if err != nil {
			t.Errorf("%s: manifesto ausente %s (%v)", nome, manifesto, err)
			continue
		}
		schema, err := os.ReadFile(filepath.Join("..", "..", "schemas", esquema))
		if err != nil {
			t.Fatalf("schema %s: %v", esquema, err)
		}
		if !gate.ValidateJSONSchema(schema, b) {
			t.Errorf("%s: %s não valida contra %s", nome, manifesto, esquema)
		}
	}
	if !temDocumentoDoProcesso(dir) {
		t.Errorf("%s: sem documento do Processo (dialeto %s ou %s)", nome, paths.SufixoAntigo, paths.SufixoStudio)
	}
}

// temDocumentoDoProcesso procura o documento do Processo em qualquer dialeto gravado.
func temDocumentoDoProcesso(dir string) bool {
	achou := false
	_ = filepath.WalkDir(dir, func(p string, e os.DirEntry, err error) error {
		if err == nil && !e.IsDir() && (strings.HasSuffix(p, paths.SufixoAntigo) || strings.HasSuffix(p, paths.SufixoStudio)) {
			achou = true
		}
		return nil
	})
	return achou
}

// validatePassa roda o legoloop validate do CLI real sobre o DAG do Processo.
func validatePassa(t *testing.T, dir string) {
	t.Helper()
	var out, errb bytes.Buffer
	dagPath := filepath.Join(dir, caminhos.ArquivoDAG)
	if code := cli.Run(context.Background(), []string{"validate", dagPath}, cli.Deps{Stdout: &out, Stderr: &errb}); code != 0 {
		t.Errorf("legoloop validate falhou em %s (%d): %s", dagPath, code, errb.String())
	}
}
