package fabricar_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/fabricar"
	"github.com/heverton-dev/LegoLoop/internal/gate"
	"github.com/heverton-dev/LegoLoop/internal/processo"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

// TestGovernancaDocsSoNaRaiz prova a regra de governança: a escrita canônica do Sexteto
// é docs/ da raiz do projeto, com prefixo ISO. O espelho em processos/<slug>/docs/ é
// proibido (Lei 15: zero documento redundante — a fatia não repete o Sexteto).
func TestGovernancaDocsSoNaRaiz(t *testing.T) {
	r := newRig(t, goodPieces())
	res := r.run(t)
	dir := filepath.Join(r.deps.Roots.Vault, caminhos.DirProcessos, res.CanonicalName)
	for _, manifesto := range []string{"triagem.json", triagem.ArquivoMemoria, "anticorpos.json", caminhos.ArquivoDAG} {
		if _, err := os.Stat(filepath.Join(dir, manifesto)); err != nil {
			t.Errorf("manifesto ausente na fatia: %s (%v)", manifesto, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "docs")); err == nil {
		t.Error("espelho proibido: processos/<slug>/docs/ não pode existir; a escrita canônica é docs/ da raiz")
	}
	for _, sub := range []string{"intents", "specs", "plans", "routes", "reports", "manuals"} {
		subDir := filepath.Join(r.deps.Roots.Docs, "docs", sub)
		entradas, err := os.ReadDir(subDir)
		if err != nil || len(entradas) == 0 {
			t.Errorf("Sexteto incompleto: docs/%s da raiz sem documento (%v)", sub, err)
		}
	}
}

// TestStudioDerivadoDoDAG prova que o <slug>.studio.json gravado pelo fabricar é a
// projeção fiel dos nós do DAG canônico e passa no gate de paridade.
func TestStudioDerivadoDoDAG(t *testing.T) {
	r := newRig(t, goodPieces())
	res := r.run(t)
	dir := filepath.Join(r.deps.Roots.Vault, caminhos.DirProcessos, res.CanonicalName)
	dagB := lerManifesto(t, filepath.Join(dir, caminhos.ArquivoDAG))
	stB := lerManifesto(t, filepath.Join(dir, res.CanonicalName+caminhos.SufixoStudio))
	if err := processo.ParidadeStudio(dagB, stB); err != nil {
		t.Errorf("studio.json não é derivado do DAG: %v", err)
	}
}

// TestNasceCincoPorCinco cobre o AC-E02.1: o fabricar instancia o scaffold 5/5 do
// nascimento do Processo — documento do Processo, DAG, triagem, memória de execução e
// anticorpos — e cada manifesto valida contra o schema oficial de schemas/.
func TestNasceCincoPorCinco(t *testing.T) {
	r := newRig(t, goodPieces())
	res := r.run(t)
	dir := filepath.Join(r.deps.Roots.Vault, caminhos.DirProcessos, res.CanonicalName)

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
	for _, candidato := range []string{nome + caminhos.SufixoStudio, nome + caminhos.SufixoAntigo} {
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

func TestResolvePorSO(t *testing.T) {
	variantes := []string{"deploy.sh", "deploy.ps1"}
	casos := []struct {
		nome, so string
		esperado string
	}{
		{"windows escolhe ps1 mesmo com sh primeiro", "windows", "deploy.ps1"},
		{"linux escolhe sh", "linux", "deploy.sh"},
		{"darwin escolhe sh", "darwin", "deploy.sh"},
		{"sem variante compativel devolve vazio", "windows", ""},
		{"so desconhecido devolve vazio", "plan9", ""},
	}
	for _, c := range casos {
		candidatos := variantes
		if c.esperado == "" {
			candidatos = []string{"deploy.sh"}
			if c.so == "plan9" {
				candidatos = variantes
			}
		}
		got := fabricar.ResolveParaSO(candidatos, c.so)
		if got != c.esperado {
			t.Errorf("%s: ResolveParaSO(%v, %q) = %q, esperado %q", c.nome, candidatos, c.so, got, c.esperado)
		}
	}
	got := fabricar.ResolvePorSO(variantes)
	esperado := "deploy.ps1"
	if runtime.GOOS != "windows" {
		esperado = "deploy.sh"
	}
	if got != esperado {
		t.Errorf("SO detectado %s: ResolvePorSO devolveu %q, esperado %q", runtime.GOOS, got, esperado)
	}
}
