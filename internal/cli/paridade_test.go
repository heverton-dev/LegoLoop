package cli_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/cli"
	"github.com/heverton-dev/LegoLoop/internal/graph"
	"github.com/heverton-dev/LegoLoop/internal/provider"
	"github.com/heverton-dev/LegoLoop/internal/studio"
)

const dagParidade = `{"name":"paridade","nodes":[
 {"id":"a","type":"script","command":["go","version"],"timeout_sec":60,
  "routing":{"mode":"manual","model":"none"},
  "output":{"artifact":"GO_VERSION_PARIDADE.txt","gate":"exit_code == 0"}}]}`

// TestParidadeTresVias (TICK-07.4): CLI, harness e Studio executam o MESMO
// processo com o mesmo contrato DAG e o mesmo resultado.
func TestParidadeTresVias(t *testing.T) {
	base := t.TempDir()
	dirProc := filepath.Join(base, caminhos.DirProcessos)
	procDir := filepath.Join(dirProc, "paridade")
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatal(err)
	}
	dagArq := caminhos.CaminhoDAG(procDir)
	if err := os.WriteFile(dagArq, []byte(dagParidade), 0o644); err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(base, "home")
	if err := os.MkdirAll(filepath.Join(home, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	artCLI := viaCLI(t, base, dirProc, home)
	artHarness := viaHarness(t, base, dirProc, home)
	artStudio, dagStudio := viaStudio(t, base, dirProc)
	if shaArquivo(dagArq) != shaArquivo(dagStudio) {
		t.Error("contrato DAG divergente entre a via CLI e a via Studio")
	}
	if artCLI != artHarness || artCLI != artStudio {
		t.Errorf("resultados divergentes: cli=%q harness=%q studio=%q", artCLI, artHarness, artStudio)
	}
}

// viaCLI executa o Processo pelo comando canônico legoloop exec.
func viaCLI(t *testing.T, base, dirProc, home string) string {
	t.Helper()
	out := filepath.Join(base, "out-cli")
	mustMkdir(t, out)
	d := cli.Deps{Stdout: io.Discard, Stderr: io.Discard, Home: home}
	code := cli.Run(context.Background(), []string{
		"exec", "paridade", "--processos-dir", dirProc, "--workdir", out, "--allow-scripts",
	}, d)
	if code != 0 {
		t.Fatalf("via CLI exit=%d", code)
	}
	return lerArtefato(t, out)
}

// viaHarness distribui a skill e executa a instrução que o harness recebeu.
func viaHarness(t *testing.T, base, dirProc, home string) string {
	t.Helper()
	d := cli.Deps{Stdout: io.Discard, Stderr: io.Discard, Home: home}
	vault := filepath.Join(base, "vault")
	code := cli.Run(context.Background(), []string{
		"instalar-skills", "--vault", vault, "--processos-dir", dirProc,
	}, d)
	if code != 0 {
		t.Fatalf("instalar-skills exit=%d", code)
	}
	cmdH := filepath.Join(home, ".claude", "commands", "legoloop-paridade.md")
	dados, err := os.ReadFile(cmdH)
	if err != nil {
		t.Fatalf("comando distribuído ao harness: %v", err)
	}
	m := regexp.MustCompile(`legoloop exec (\S+)`).FindSubmatch(dados)
	if m == nil {
		t.Fatalf("instrução do harness não delega a legoloop exec: %s", dados)
	}
	out := filepath.Join(base, "out-harness")
	mustMkdir(t, out)
	if code := cli.Run(context.Background(), []string{
		"exec", string(m[1]), "--processos-dir", dirProc, "--workdir", out, "--allow-scripts",
	}, d); code != 0 {
		t.Fatalf("via harness exit=%d", code)
	}
	return lerArtefato(t, out)
}

// viaStudio invoca a peça registrada no Studio e devolve artefato e DAG lido.
func viaStudio(t *testing.T, base, dirProc string) (string, string) {
	t.Helper()
	out := filepath.Join(base, "out-studio")
	procOut := filepath.Join(out, caminhos.DirProcessos, "paridade")
	mustMkdir(t, procOut)
	dagStudio := caminhos.CaminhoDAG(procOut)
	copyArquivo(t, caminhos.CaminhoDAG(filepath.Join(dirProc, "paridade")), dagStudio)
	idx, err := graph.NewIndex(graph.DefaultCatalog())
	if err != nil {
		t.Fatal(err)
	}
	srv := &studio.Server{
		Index: idx, ProjectDir: out, Host: "127.0.0.1:7878",
		Discover: func(context.Context) provider.Inventory { return provider.Inventory{} },
	}
	req := httptest.NewRequest("POST", "http://127.0.0.1:7878/api/pecas/paridade/exec",
		strings.NewReader(`{"allow_scripts": true}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("via Studio code=%d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"ok":true`) {
		t.Fatalf("via Studio sem dag_finish ok: %s", rec.Body.String())
	}
	return lerArtefato(t, out), dagStudio
}

func lerArtefato(t *testing.T, workdir string) string {
	t.Helper()
	dados, err := os.ReadFile(filepath.Join(workdir, "GO_VERSION_PARIDADE.txt"))
	if err != nil {
		t.Fatalf("artefato ausente em %s: %v", workdir, err)
	}
	return string(dados)
}

func shaArquivo(caminho string) string {
	dados, _ := os.ReadFile(caminho)
	s := sha256.Sum256(dados)
	return hex.EncodeToString(s[:])
}

func copyArquivo(t *testing.T, origem, destino string) {
	t.Helper()
	dados, err := os.ReadFile(origem)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destino, dados, 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustMkdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}
