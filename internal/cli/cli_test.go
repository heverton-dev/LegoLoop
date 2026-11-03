package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/gate"
	"legoloop/internal/provider"
)

func deps(t *testing.T) (Deps, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	var out, errb bytes.Buffer
	return Deps{
		Stdout: &out, Stderr: &errb,
		Discover: func(context.Context) provider.Inventory {
			return provider.Inventory{Ollama: []string{"qwen2.5-coder:7b"}, APIKeys: []string{"anthropic"},
				Errors: []string{"9router: nenhuma base respondeu"}}
		},
	}, &out, &errb
}

func run(t *testing.T, d Deps, args ...string) int {
	t.Helper()
	return Run(context.Background(), args, d)
}

func TestHelpAndUnknown(t *testing.T) {
	d, out, errb := deps(t)
	if code := run(t, d); code != 2 || !strings.Contains(errb.String(), "Uso:") {
		t.Fatalf("sem args: code=%d", code)
	}
	if code := run(t, d, "help"); code != 0 || !strings.Contains(out.String(), "canonize") || !strings.Contains(out.String(), "import") {
		t.Fatalf("help: code=%d", code)
	}
	if code := run(t, d, "explodir"); code != 2 {
		t.Fatalf("comando desconhecido: code=%d", code)
	}
	if code := run(t, d, "version"); code != 0 || !strings.Contains(out.String(), "legoloop") {
		t.Fatal("version")
	}
}

func TestDoctor(t *testing.T) {
	d, out, _ := deps(t)
	if code := run(t, d, "doctor"); code != 0 {
		t.Fatalf("code=%d", code)
	}
	for _, want := range []string{"qwen2.5-coder:7b", "ollama", "anthropic", "9Router", "offline"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("doctor sem %q:\n%s", want, out.String())
		}
	}
	d2, out2, _ := deps(t)
	run(t, d2, "doctor", "--json")
	var inv map[string]any
	if err := json.Unmarshal(out2.Bytes(), &inv); err != nil {
		t.Fatalf("--json inválido: %v", err)
	}
}

func canonFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	body := "---\nname: x\n---\n" + strings.Repeat("conteúdo ", 100)
	for _, p := range []string{"a/SKILL.md", "b/SKILL.md", "c/deep/SKILL.md"} {
		full := filepath.Join(root, filepath.FromSlash(p))
		os.MkdirAll(filepath.Dir(full), 0o755)
		os.WriteFile(full, []byte(body), 0o644)
	}
	return root, filepath.Join(t.TempDir(), "canonical")
}

func TestCanonizeDryRunThenApply(t *testing.T) {
	root, vault := canonFixture(t)
	d, out, _ := deps(t)
	if code := run(t, d, "canonize", "--root", root, "--vault", vault, "--dry-run"); code != 0 {
		t.Fatalf("dry-run code=%d", code)
	}
	if !strings.Contains(out.String(), "DRY-RUN") || !strings.Contains(out.String(), "economia") {
		t.Fatalf("relatório dry-run:\n%s", out.String())
	}
	if _, err := os.Stat(vault); !os.IsNotExist(err) {
		t.Fatal("dry-run escreveu no disco")
	}
	d2, out2, errb := deps(t)
	if code := run(t, d2, "canonize", "--root", root, "--vault", vault, "--apply", "--json"); code != 0 {
		t.Fatalf("apply code=%d stderr=%s", code, errb.String())
	}
	var rep struct {
		Linked     []any `json:"linked"`
		SavedBytes int64 `json:"saved_bytes"`
	}
	if err := json.Unmarshal(out2.Bytes(), &rep); err != nil || len(rep.Linked) != 3 || rep.SavedBytes <= 0 {
		t.Fatalf("apply json: %v %s", err, out2.String())
	}
}

func TestCanonizeFlagErrors(t *testing.T) {
	d, _, _ := deps(t)
	cases := [][]string{
		{"canonize"},
		{"canonize", "--root", "x"},
		{"canonize", "--root", "x", "--dry-run", "--apply"},
		{"canonize", "--bogus"},
		{"canonize", "--root", filepath.Join(t.TempDir(), "nao-existe"), "--dry-run"},
	}
	for _, c := range cases {
		if code := run(t, d, c...); code == 0 {
			t.Errorf("%v deveria falhar", c)
		}
	}
}

const pipelineScript = `{"name":"p","nodes":[
 {"id":"a","type":"script","command":["go","version"],"timeout_sec":60,
  "routing":{"mode":"manual","model":"none"},
  "output":{"artifact":"GO_VERSION.txt","gate":"exit_code == 0 && has_pattern('GO_VERSION.txt', 'go1\\.')"}}]}`

func TestValidateAndRun(t *testing.T) {
	work := t.TempDir()
	p := filepath.Join(work, "pipeline.dag.json")
	os.WriteFile(p, []byte(pipelineScript), 0o644)
	d, out, _ := deps(t)
	if code := run(t, d, "validate", p); code != 0 || !strings.Contains(out.String(), "camada 1") {
		t.Fatalf("validate code=%d out=%s", code, out.String())
	}
	d2, out2, errb := deps(t)
	if code := run(t, d2, "run", p, "--workdir", work, "--allow-scripts"); code != 0 {
		t.Fatalf("run code=%d out=%s err=%s", code, out2.String(), errb.String())
	}
	if !strings.Contains(out2.String(), "passed") || !strings.Contains(out2.String(), "MODELO") || !strings.Contains(out2.String(), "script/go") {
		t.Fatalf("relatório:\n%s", out2.String())
	}
	if b, _ := os.ReadFile(filepath.Join(work, "GO_VERSION.txt")); !strings.Contains(string(b), "go1.") {
		t.Fatal("artefato do script não gravado")
	}
}

func TestRunFailureExitCode(t *testing.T) {
	work := t.TempDir()
	p := filepath.Join(work, "pipeline.dag.json")
	os.WriteFile(p, []byte(strings.Replace(pipelineScript, `go1\\.`, `nunca-casa`, 1)), 0o644)
	d, out, _ := deps(t)
	if code := run(t, d, "run", p, "--workdir", work, "--json", "--var", "idea=x", "--allow-scripts"); code != 1 {
		t.Fatalf("gate reprovado deveria sair 1, code=%d", code)
	}
	var rep map[string]any
	if err := json.Unmarshal(out.Bytes(), &rep); err != nil || rep["ok"] != false {
		t.Fatalf("json: %v %s", err, out.String())
	}
}

func TestValidateAndRunErrors(t *testing.T) {
	d, _, _ := deps(t)
	bad := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(bad, []byte(`{"name":"x","nodes":[]}`), 0o644)
	cases := [][]string{
		{"validate"}, {"validate", bad}, {"validate", "nao-existe.json"},
		{"run"}, {"run", bad}, {"run", bad, "--var", "semigual"},
	}
	for _, c := range cases {
		if code := run(t, d, c...); code == 0 {
			t.Errorf("%v deveria falhar", c)
		}
	}
}

func TestRoute(t *testing.T) {
	d, out, _ := deps(t)
	if code := run(t, d, "route", "--from", "Artifact:URL", "--to", "Artifact:Spec"); code != 0 {
		t.Fatalf("code=%d", code)
	}
	if !strings.Contains(out.String(), "obscura → rea → impeccable → aidd-spec") {
		t.Fatalf("rota:\n%s", out.String())
	}
	if code := run(t, d, "route", "--from", "Artifact:Audit", "--to", "Artifact:URL"); code == 0 {
		t.Fatal("rota impossível deveria falhar")
	}
}

func TestStudioServesAndStops(t *testing.T) {
	d, _, _ := deps(t)
	ready := make(chan string, 1)
	d.Ready = ready
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() { done <- Run(ctx, []string{"studio", "--addr", "127.0.0.1:0", "--project", t.TempDir()}, d) }()
	var addr string
	select {
	case addr = <-ready:
	case <-time.After(5 * time.Second):
		t.Fatal("studio não subiu")
	}
	resp, err := http.Get("http://" + addr + "/api/blocks")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("studio: %v", err)
	}
	resp.Body.Close()
	cancel()
	if code := <-done; code != 0 {
		t.Fatalf("encerramento code=%d", code)
	}
	if code := run(t, d, "studio", "--addr", "0.0.0.0:7878"); code == 0 {
		t.Fatal("bind fora de loopback deveria falhar")
	}
}

func startStudio(t *testing.T, d Deps, args ...string) (addr string, stop func() int) {
	t.Helper()
	ready := make(chan string, 1)
	d.Ready = ready
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() { done <- Run(ctx, append([]string{"studio", "--addr", "127.0.0.1:0"}, args...), d) }()
	select {
	case addr = <-ready:
	case <-time.After(5 * time.Second):
		cancel()
		t.Fatal("studio não subiu")
	}
	return addr, func() int { cancel(); return <-done }
}

// O comando studio liga o inventário de ativos: /api/assets responde com dados reais do HOME informado.
func TestStudioWiresAssetInventory(t *testing.T) {
	d, _, _ := deps(t)
	home := t.TempDir()
	d.Home = home
	skill := filepath.Join(home, ".claude", "skills", "alpha", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(skill), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skill, []byte("---\nname: alpha\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	extra := t.TempDir()
	addr, stop := startStudio(t, d, "--project", t.TempDir(), "--canon-root", extra)
	resp, err := http.Get("http://" + addr + "/api/assets")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("/api/assets: %v %v", err, resp)
	}
	var body struct {
		Counts map[string]int `json:"counts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if body.Counts["skill"] != 1 || body.Counts["gate"] != len(gate.Catalog()) {
		t.Errorf("counts = %v", body.Counts)
	}
	if code := stop(); code != 0 {
		t.Fatalf("encerramento code=%d", code)
	}
}

func TestStudioRejectsInvalidCanonRoot(t *testing.T) {
	d, _, errb := deps(t)
	missing := filepath.Join(t.TempDir(), "nao-existe")
	// Com prazo: se a validação sumir, o servidor subiria e o teste falharia em vez de travar.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []string{"studio", "--addr", "127.0.0.1:0", "--project", t.TempDir(), "--canon-root", missing}
	if code := Run(ctx, args, d); code == 0 {
		t.Fatal("--canon-root inexistente deveria falhar")
	}
	if !strings.Contains(errb.String(), "canon-root") {
		t.Errorf("erro deveria citar a flag: %q", errb.String())
	}
}

func TestHumanBytes(t *testing.T) {
	cases := map[int64]string{512: "512 B", 2048: "2.0 KiB", 5 << 20: "5.0 MiB", 3 << 30: "3.0 GiB"}
	for in, want := range cases {
		if got := humanBytes(in); got != want {
			t.Errorf("humanBytes(%d) = %s, want %s", in, got, want)
		}
	}
}

// SEC-01: sem --allow-scripts o nó script não executa e o motivo é explicado.
func TestRunRefusesScriptsWithoutConsent(t *testing.T) {
	work := t.TempDir()
	p := filepath.Join(work, "pipeline.dag.json")
	os.WriteFile(p, []byte(pipelineScript), 0o644)
	d, out, _ := deps(t)
	if code := run(t, d, "run", p, "--workdir", work); code != 1 {
		t.Fatalf("script rodou sem consentimento: code=%d", code)
	}
	if !strings.Contains(out.String(), "--allow-scripts") {
		t.Fatalf("relatório não explica o bloqueio:\n%s", out.String())
	}
	if _, err := os.Stat(filepath.Join(work, "GO_VERSION.txt")); !os.IsNotExist(err) {
		t.Fatal("comando executou mesmo bloqueado")
	}
}

// SEC-11: pipeline acima do limite é recusado antes de ser carregado inteiro.
func TestValidateRejectsHugeFile(t *testing.T) {
	big := filepath.Join(t.TempDir(), "big.json")
	os.WriteFile(big, make([]byte, 2<<20), 0o644)
	d, _, errb := deps(t)
	if code := run(t, d, "validate", big); code == 0 || !strings.Contains(errb.String(), "excede") {
		t.Fatalf("arquivo gigante aceito: %s", errb.String())
	}
}

// SEC-04: --include-git é aceito e repositórios git ficam fora por padrão.
func TestCanonizeIncludeGitFlag(t *testing.T) {
	root, vault := canonFixture(t)
	os.MkdirAll(filepath.Join(root, "a", ".git"), 0o755)
	d, out, _ := deps(t)
	run(t, d, "canonize", "--root", root, "--vault", vault, "--dry-run", "--json")
	var p struct {
		Groups []struct {
			Links []string `json:"links"`
		} `json:"groups"`
	}
	json.Unmarshal(out.Bytes(), &p)
	if len(p.Groups) != 1 || len(p.Groups[0].Links) != 2 {
		t.Fatalf("repo git deveria ficar fora: %s", out.String())
	}
	d2, out2, _ := deps(t)
	run(t, d2, "canonize", "--root", root, "--vault", vault, "--dry-run", "--json", "--include-git")
	json.Unmarshal(out2.Bytes(), &p)
	if len(p.Groups[0].Links) != 3 {
		t.Fatalf("--include-git ignorado: %s", out2.String())
	}
}

// R3-1: raiz pulada por estar num repositório git é avisada, nunca silenciosa.
func TestCanonizeWarnsSkippedGitRoot(t *testing.T) {
	root, vault := canonFixture(t)
	os.MkdirAll(filepath.Join(root, ".git"), 0o755)
	d, _, errb := deps(t)
	run(t, d, "canonize", "--root", filepath.Join(root, "a"), "--vault", vault, "--dry-run")
	if !strings.Contains(errb.String(), "--include-git") {
		t.Fatalf("raiz pulada sem aviso: %q", errb.String())
	}
}

// TICK-1203: legoloop forge abre o Studio e dispara a primeira varredura.
func TestForgeServesAndTriggersInitialScan(t *testing.T) {
	d, out, _ := deps(t)
	d.Home = t.TempDir()
	ready := make(chan string, 1)
	d.Ready = ready
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	done := make(chan int, 1)
	go func() { done <- Run(ctx, []string{"forge", "--addr", "127.0.0.1:0", "--project", t.TempDir()}, d) }()
	var addr string
	select {
	case addr = <-ready:
	case <-time.After(15 * time.Second):
		t.Fatal("forge não subiu a tempo")
	}
	resp, err := http.Get("http://" + addr + "/api/blocks")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("forge /api/blocks: err=%v status=%v", err, resp.StatusCode)
	}
	resp.Body.Close()
	cancel()
	<-done
	if !strings.Contains(out.String(), "Looping Studio") && !strings.Contains(out.String(), "http://") {
		t.Fatalf("saída do forge deve informar URL: %s", out.String())
	}
}

func TestRunInteractiveFlags(t *testing.T) {
	work := t.TempDir()
	p := filepath.Join(work, "pipeline.dag.json")
	dagJSON := `{
		"name": "cli-interactive-test",
		"nodes": [{
			"id": "node_ask",
			"type": "script",
			"command": ["go", "version"],
			"routing": {"mode": "manual", "model": "test"},
			"output": {"artifact": "out.txt", "gate": "interactive('Aprovar acao?') && exit_code == 0"}
		}]
	}`
	if err := os.WriteFile(p, []byte(dagJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1. --yes deve aprovar automaticamente
	d1, _, _ := deps(t)
	if code := run(t, d1, "run", p, "--workdir", work, "--allow-scripts", "--yes"); code != 0 {
		t.Fatalf("--yes deveria aprovar, code=%d", code)
	}

	// 2. -y (alias) deve aprovar
	d2, _, _ := deps(t)
	if code := run(t, d2, "run", p, "--workdir", work, "--allow-scripts", "-y"); code != 0 {
		t.Fatalf("-y deveria aprovar, code=%d", code)
	}

	// 3. --non-interactive deve rejeitar
	d3, _, _ := deps(t)
	if code := run(t, d3, "run", p, "--workdir", work, "--allow-scripts", "--non-interactive"); code != 1 {
		t.Fatalf("--non-interactive deveria falhar com 1, code=%d", code)
	}

	// 4. Stdin com "s\n" deve aprovar
	d4, out4, _ := deps(t)
	d4.Stdin = strings.NewReader("s\n")
	if code := run(t, d4, "run", p, "--workdir", work, "--allow-scripts"); code != 0 {
		t.Fatalf("stdin 's' deveria aprovar, code=%d\nout: %s", code, out4.String())
	}

	// 5. Stdin com "n\n" deve rejeitar
	d5, _, _ := deps(t)
	d5.Stdin = strings.NewReader("n\n")
	if code := run(t, d5, "run", p, "--workdir", work, "--allow-scripts"); code != 1 {
		t.Fatalf("stdin 'n' deveria rejeitar com 1, code=%d", code)
	}
}

func TestImportCLIValidation(t *testing.T) {
	d, _, errb := deps(t)

	// Sem argumentos
	if code := run(t, d, "import"); code != 2 {
		t.Fatalf("import sem args: code=%d", code)
	}
	if !strings.Contains(errb.String(), "informe a URL") {
		t.Fatalf("mensagem esperada no stderr: %s", errb.String())
	}

	// URL com domínio não-github
	d2, _, errb2 := deps(t)
	if code := run(t, d2, "import", "https://notgithub.com/skill"); code != 1 {
		t.Fatalf("import com url inválida: code=%d", code)
	}
	if !strings.Contains(errb2.String(), "apenas URLs do github.com") {
		t.Fatalf("mensagem esperada sobre github: %s", errb2.String())
	}
}
