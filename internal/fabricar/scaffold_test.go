package fabricar_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/dag"
	"legoloop/internal/gate"
	"legoloop/internal/paths"
	"legoloop/internal/triagem"
)

// TestNasceCincoPorCinco cobre o AC-E02.1: o fabricar instancia o scaffold 5/5 do
// nascimento do Processo — documento do Processo, DAG, triagem, memória de execução e
// anticorpos — e cada manifesto valida contra o schema oficial de schemas/.
func TestNasceCincoPorCinco(t *testing.T) {
	r := newRig(t, goodPieces())
	res := r.run(t)
	dir := filepath.Join(r.deps.Roots.Vault, paths.DirAntigo, res.CanonicalName)

	comSchema := map[string]string{
		"triagem.json":         "triagem.v1.json",
		triagem.ArquivoMemoria: "memoria_execucao.v1.json",
		"anticorpos.json":      "anticorpos.v1.json",
	}
	for nome, esquema := range comSchema {
		manifesto := lerManifesto(t, filepath.Join(dir, nome))
		schema := lerManifesto(t, filepath.Join("..", "..", "schemas", esquema))
		if !gate.ValidateJSONSchema(schema, manifesto) {
			t.Errorf("%s não valida contra %s", nome, esquema)
		}
	}

	doc := documentoDoProcesso(t, dir, res.CanonicalName)
	schemaDoc := lerManifesto(t, filepath.Join("..", "..", "schemas", "doc_processo.v1.json"))
	if !gate.ValidateJSONSchema(schemaDoc, doc) {
		t.Error("o documento do Processo não valida contra doc_processo.v1.json")
	}

	bruto := lerManifesto(t, filepath.Join(dir, caminhos.ArquivoDAG))
	if _, err := dag.Parse(bruto); err != nil {
		t.Errorf("o DAG do Processo não parseia: %v", err)
	}
	var anticorpos map[string]any
	if err := json.Unmarshal(lerManifesto(t, filepath.Join(dir, "anticorpos.json")), &anticorpos); err != nil {
		t.Fatalf("anticorpos.json ilegível: %v", err)
	}
	if anticorpos["schema_version"] != "anticorpos.v1" || anticorpos["processo"] != res.CanonicalName {
		t.Errorf("anticorpos.json com identidade errada: %v", anticorpos)
	}
}

// documentoDoProcesso devolve o documento do Studio no dialeto gravado (novo ou antigo).
func documentoDoProcesso(t *testing.T, dir, nome string) []byte {
	t.Helper()
	for _, candidato := range []string{nome + paths.SufixoStudio, nome + paths.SufixoAntigo} {
		if b, err := os.ReadFile(filepath.Join(dir, candidato)); err == nil {
			return b
		}
	}
	t.Fatalf("sem documento do Processo em %s", dir)
	return nil
}

func lerManifesto(t *testing.T, caminho string) []byte {
	t.Helper()
	b, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler %s: %v", caminho, err)
	}
	return b
}
