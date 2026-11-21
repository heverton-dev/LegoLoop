package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/runner"
)

var artefatosTresScripts = []string{"n1.txt", "n2.txt", "n3.txt"}

func registrarProcessoDeterminismo(t *testing.T, slug string) string {
	t.Helper()
	processosDir := filepath.Join(t.TempDir(), caminhos.DirProcessos)
	procDir := filepath.Join(processosDir, slug)
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatalf("criar diretório do processo: %v", err)
	}
	if err := os.WriteFile(caminhos.CaminhoDAG(procDir), []byte(dagTresScriptsDeterminismo), 0o644); err != nil {
		t.Fatalf("gravar DAG do processo: %v", err)
	}
	return processosDir
}

func sha256Arquivo(t *testing.T, caminho string) string {
	t.Helper()
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler artefato %s: %v", caminho, err)
	}
	if len(dados) == 0 {
		t.Fatalf("artefato %s está vazio", caminho)
	}
	soma := sha256.Sum256(dados)
	return hex.EncodeToString(soma[:])
}

func hashesArtefatos(t *testing.T, workDir string, nomes []string) map[string]string {
	t.Helper()
	out := make(map[string]string, len(nomes))
	for _, nome := range nomes {
		out[nome] = sha256Arquivo(t, filepath.Join(workDir, nome))
	}
	return out
}

func executarExecDeterminismo(t *testing.T, bin, processosDir, slug, root string) (runner.Report, string, map[string]string) {
	t.Helper()
	router := binRouter(t, 100)
	work := filepath.Join(root, "work")
	vault := filepath.Join(root, "vault")
	repPath := filepath.Join(root, "report.md")
	res := runBinaryEm(t, bin, router, root, work,
		"exec", slug, "--processos-dir", processosDir,
		"--workdir", work, "--skills", vault,
		"--allow-scripts", "--json", "--report", repPath)
	if res.code != 0 {
		t.Fatalf("legoloop exec saiu %d\nstdout:\n%s\nstderr:\n%s", res.code, res.out, res.errb)
	}
	var rep runner.Report
	if err := json.Unmarshal([]byte(res.out), &rep); err != nil {
		t.Fatalf("json inválido do relatório: %v\nsaída:\n%s", err, res.out)
	}
	factual, err := os.ReadFile(repPath)
	if err != nil {
		t.Fatalf("ler relatório factual: %v", err)
	}
	return rep, string(factual), hashesArtefatos(t, work, artefatosTresScripts)
}

// TestExecDuasVezesRelatorioIdentico (AC-E01.3): executar `legoloop exec <slug>`
// duas vezes sobre um DAG de três nós script registrado como Processo gera
// artefatos com o mesmo SHA-256 e relatórios idênticos, descontadas durações.
func TestExecDuasVezesRelatorioIdentico(t *testing.T) {
	bin := buildBinary(t)
	const slug = "determinismo-tres-scripts"
	processosDir := registrarProcessoDeterminismo(t, slug)

	rep1, factual1, hashes1 := executarExecDeterminismo(t, bin, processosDir, slug, t.TempDir())
	rep2, factual2, hashes2 := executarExecDeterminismo(t, bin, processosDir, slug, t.TempDir())

	if !rep1.OK || !rep2.OK || len(rep1.Nodes) != 3 || len(rep2.Nodes) != 3 {
		t.Fatalf("esperava 3 nós aprovados nas duas execuções: rep1=%+v rep2=%+v", rep1, rep2)
	}
	if !reflect.DeepEqual(hashes1, hashes2) {
		t.Fatalf("hashes SHA-256 divergiram:\nexec1: %v\nexec2: %v", hashes1, hashes2)
	}
	t.Logf("hashes SHA-256 idênticos: %v", hashes1)

	n1, n2 := normalizarRelatorioRun(rep1), normalizarRelatorioRun(rep2)
	if !reflect.DeepEqual(n1, n2) {
		t.Fatalf("relatórios divergiram (descontadas as durações):\nexec1: %+v\nexec2: %+v", n1, n2)
	}
	if factual1 != factual2 {
		t.Fatalf("relatórios factuais divergiram:\nexec1:\n%s\nexec2:\n%s", factual1, factual2)
	}
}
