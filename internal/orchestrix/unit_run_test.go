package orchestrix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
	"legoloop/internal/runner"
)

// helperClaudeUnit é o claude falso do UT-04: grava em cwd os arquivos de LEGOLOOP_UNIT_WRITE
// ("caminho=conteúdo", separados por "|") e devolve um resultado JSON.
func helperClaudeUnit() {
	_, _ = io.ReadAll(os.Stdin)
	if dump := os.Getenv("LEGOLOOP_UNIT_ENVDUMP"); dump != "" {
		line := "BASE=" + os.Getenv("ANTHROPIC_BASE_URL") + "\nTOKEN=" + os.Getenv("ANTHROPIC_AUTH_TOKEN") + "\n"
		_ = os.WriteFile(dump, []byte(line), 0o644)
	}
	for _, spec := range strings.Split(os.Getenv("LEGOLOOP_UNIT_WRITE"), "|") {
		rel, content, ok := strings.Cut(spec, "=")
		if !ok {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(rel), 0o755); err == nil {
			if err = os.WriteFile(rel, []byte(content+"\n"), 0o644); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(3)
			}
		}
	}
	out, _ := json.Marshal(map[string]any{"type": "result", "result": "feito\nsegunda linha",
		"usage": map[string]int{"input_tokens": 5, "output_tokens": 2}})
	os.Stdout.Write(out)
	os.Exit(0)
}

func unitOrchestrix(write string) *Orchestrix {
	return &Orchestrix{
		Inventory:      provider.Inventory{CLIs: []provider.CLI{{Name: "claude", Path: os.Args[0]}}},
		ExtraEnv:       []string{"LEGOLOOP_HELPER=claude-unit", "LEGOLOOP_UNIT_WRITE=" + write},
		AllowUnitTasks: true,
	}
}

func unitNode(id, out, gateExpr string, cmd []string) dag.Node {
	return dag.Node{
		ID: id, Type: dag.TypeUnitTask, Command: cmd,
		Routing:    dag.RoutingConfig{Mode: dag.ModeManual, Model: "sonnet"},
		ContextOut: []string{id + ".txt"},
		Output:     dag.HandoffContract{Artifact: out, Gate: gateExpr},
	}
}

func pathExists(p string) bool { _, err := os.Lstat(p); return err == nil }

func TestUnitTaskRefusedWithoutFlag(t *testing.T) {
	repo := newRepoT(t)
	o := unitOrchestrix("n1.txt=x")
	o.AllowUnitTasks = false
	_, err := o.Execute(context.Background(), unitNode("n1", "out/n1.json", "exit_code == 0", nil), repo)
	if err == nil || !strings.Contains(err.Error(), "--allow-unit-tasks") {
		t.Fatalf("err = %v", err)
	}
	if pathExists(filepath.Join(repo, "n1.txt")) {
		t.Fatal("nada deveria ser escrito sem a flag")
	}
}

func TestUnitTaskPassAppliesAndWritesArtifact(t *testing.T) {
	repo := newRepoT(t)
	head := strings.TrimSpace(gitT(t, repo, "rev-parse", "HEAD"))
	n := unitNode("n1", "out/n1.json", "exit_code == 0 && file_exists('n1.txt')", []string{"git", "status", "--short"})
	res, err := unitOrchestrix("n1.txt=ola").Execute(context.Background(), n, repo)
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 || res.Tokens.Input != 5 {
		t.Fatalf("res = %+v", res)
	}
	if got := readT(t, repo, "n1.txt"); got != "ola\n" {
		t.Fatalf("diff não aplicado: %q", got)
	}
	var rep unitReport
	if err := json.Unmarshal([]byte(readT(t, repo, "out/n1.json")), &rep); err != nil {
		t.Fatal(err)
	}
	if !rep.Passed || rep.Node != "n1" || rep.Summary != "feito" || len(rep.Changed) != 1 || rep.Changed[0] != "n1.txt" {
		t.Fatalf("artefato = %+v", rep)
	}
	if got := strings.TrimSpace(gitT(t, repo, "rev-parse", "HEAD")); got != head {
		t.Fatal("a Orchestrix não pode fazer commit")
	}
	if c := worktreeCount(t, repo); c != 1 {
		t.Fatalf("worktree vazou: %d", c)
	}
}

func TestUnitTaskOutsideBoundaryFails(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "exit_code == 0", nil)
	_, err := unitOrchestrix("n1.txt=ok|intruso.txt=x|a.txt=mexeu").Execute(context.Background(), n, repo)
	if err == nil || !strings.Contains(err.Error(), "write_boundary") {
		t.Fatalf("err = %v", err)
	}
	for _, f := range []string{"n1.txt", "intruso.txt", "out/n1.json"} {
		if pathExists(filepath.Join(repo, filepath.FromSlash(f))) {
			t.Fatalf("nó reprovado aplicou %s", f)
		}
	}
	if readT(t, repo, "a.txt") != "a\n" {
		t.Fatal("a.txt foi alterado")
	}
	if c := worktreeCount(t, repo); c != 1 {
		t.Fatalf("worktree vazou: %d", c)
	}
}

func TestUnitTaskGateFailureDiscards(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "exit_code == 0", []string{"git", "rev-parse", "--verify", "refs/nao/existe"})
	_, err := unitOrchestrix("n1.txt=ok").Execute(context.Background(), n, repo)
	if err == nil || !strings.Contains(err.Error(), "reprovou") {
		t.Fatalf("err = %v", err)
	}
	if pathExists(filepath.Join(repo, "n1.txt")) || pathExists(filepath.Join(repo, "out", "n1.json")) {
		t.Fatal("gate reprovado não pode aplicar nada")
	}
	if c := worktreeCount(t, repo); c != 1 {
		t.Fatalf("worktree vazou: %d", c)
	}
}

func TestUnitTaskRequiresRepoRoot(t *testing.T) {
	n := unitNode("n1", "out/n1.json", "exit_code == 0", nil)
	if _, err := unitOrchestrix("n1.txt=x").Execute(context.Background(), n, t.TempDir()); err == nil {
		t.Fatal("diretório sem git deveria falhar")
	}
	repo := newRepoT(t)
	if _, err := unitOrchestrix("n1.txt=x").Execute(context.Background(), n, filepath.Join(repo, "docs")); err == nil {
		t.Fatal("subdiretório do repositório deveria falhar")
	}
}

func TestUnitTaskRequiresClaudeModel(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "exit_code == 0", nil)
	n.Routing = dag.RoutingConfig{Mode: dag.ModeManual, Model: "inexistente"}
	if _, err := unitOrchestrix("n1.txt=x").Execute(context.Background(), n, repo); err == nil {
		t.Fatal("modelo fora do inventário deveria falhar")
	}
}

const twoNodeDAG = `{"name":"ut","nodes":[
 {"id":"n1","type":"unit_task","routing":{"mode":"manual","model":"sonnet"},"context_out":%s,
  "output":{"artifact":"out/n1.json","gate":"exit_code == 0 && file_exists('n1.txt')"}},
 {"id":"n2","type":"unit_task","depends_on":["n1"],"routing":{"mode":"manual","model":"sonnet"},"context_out":%s,
  "output":{"artifact":"out/n2.json","gate":"exit_code == 0 && file_exists('n2.txt')"}}]}`

func runTwoNodeDAG(t *testing.T, repo, outs1, outs2 string) runner.Report {
	t.Helper()
	d, err := dag.Parse([]byte(fmt.Sprintf(twoNodeDAG, outs1, outs2)))
	if err != nil {
		t.Fatal(err)
	}
	rep, err := runner.Run(context.Background(), d, runner.Options{
		Workdir: repo, Parallelism: 1, Executor: unitOrchestrix("n1.txt=um|n2.txt=dois")})
	if err != nil {
		t.Fatal(err)
	}
	return rep
}

// O claude falso escreve n1.txt e n2.txt em cada nó; o write_boundary de n1 só permite n1.txt.
func TestUnitTaskDAGOutsideBoundaryFails(t *testing.T) {
	repo := newRepoT(t)
	rep := runTwoNodeDAG(t, repo, `["n1.txt"]`, `["n2.txt"]`)
	if rep.OK || rep.Nodes[0].Status != runner.StatusError || rep.Nodes[1].Status != runner.StatusSkipped {
		t.Fatalf("fora do write_boundary deveria reprovar n1 e pular n2: %+v", rep.Nodes)
	}
	if !strings.Contains(rep.Nodes[0].Error, "write_boundary") {
		t.Fatalf("erro = %q", rep.Nodes[0].Error)
	}
	if pathExists(filepath.Join(repo, "n1.txt")) {
		t.Fatal("n1 reprovado aplicou arquivo")
	}
}

func TestUnitTaskDAGOfTwoPasses(t *testing.T) {
	repo := newRepoT(t)
	rep := runTwoNodeDAG(t, repo, `["n1.txt","n2.txt"]`, `["n1.txt","n2.txt"]`)
	if !rep.OK {
		t.Fatalf("DAG reprovou: %+v", rep.Nodes)
	}
	for _, f := range []string{"n1.txt", "n2.txt", "out/n1.json", "out/n2.json"} {
		if !pathExists(filepath.Join(repo, filepath.FromSlash(f))) {
			t.Errorf("%s ausente", f)
		}
	}
	if c := worktreeCount(t, repo); c != 1 {
		t.Fatalf("worktree vazou: %d", c)
	}
}

// UT-05: o command do nó 2 lê, na worktree, o artefato do nó 1 ainda sem commit.
func unitNode2(cmd []string) dag.Node {
	n := unitNode("n2", "out/n2.json", "exit_code == 0 && file_exists('n2.txt')", cmd)
	n.ContextIn = []string{"n1.txt", "ainda-nao-existe.txt"}
	return n
}

func TestUnitTaskSeedsUncommittedContextIn(t *testing.T) {
	repo := newRepoT(t)
	writeT(t, repo, "n1.txt", "do no 1\n") // sem commit
	n := unitNode2([]string{"git", "hash-object", "n1.txt"})
	if _, err := unitOrchestrix("n2.txt=dois").Execute(context.Background(), n, repo); err != nil {
		t.Fatal(err)
	}
	var rep unitReport
	if err := json.Unmarshal([]byte(readT(t, repo, "out/n2.json")), &rep); err != nil {
		t.Fatal(err)
	}
	if len(rep.Changed) != 1 || rep.Changed[0] != "n2.txt" {
		t.Fatalf("context_in semeado não pode entrar em changed: %v", rep.Changed)
	}
	if readT(t, repo, "n1.txt") != "do no 1\n" || !pathExists(filepath.Join(repo, "n2.txt")) {
		t.Fatal("estado final do workdir inesperado")
	}
}

func TestUnitTaskWithoutSeedCommandWouldFail(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode2([]string{"git", "hash-object", "n1.txt"}) // n1.txt não existe no principal
	if _, err := unitOrchestrix("n2.txt=dois").Execute(context.Background(), n, repo); err == nil {
		t.Fatal("command deveria falhar sem o arquivo")
	}
}

func TestUnitTaskModifyingSeededContextFails(t *testing.T) {
	repo := newRepoT(t)
	writeT(t, repo, "n1.txt", "do no 1\n")
	n := unitNode2(nil)
	_, err := unitOrchestrix("n2.txt=dois|n1.txt=adulterado").Execute(context.Background(), n, repo)
	if err == nil || !strings.Contains(err.Error(), "write_boundary") {
		t.Fatalf("err = %v", err)
	}
	if readT(t, repo, "n1.txt") != "do no 1\n" || pathExists(filepath.Join(repo, "n2.txt")) {
		t.Fatal("nó reprovado alterou o workdir")
	}
}

func TestUnitTaskRejectsGitContextIn(t *testing.T) {
	repo := newRepoT(t)
	for _, f := range []string{".git/config", ".GIT/config", "../fora.txt"} {
		n := unitNode2(nil)
		n.ContextIn = []string{f}
		if _, err := unitOrchestrix("n2.txt=x").Execute(context.Background(), n, repo); err == nil {
			t.Errorf("context_in %q deveria reprovar", f)
		}
	}
}

// Saída em caminho ignorado pelo git (ex.: scratch/) precisa ser detectada e aplicada.
func TestUnitTaskAppliesGitIgnoredOutput(t *testing.T) {
	repo := newRepoT(t)
	writeT(t, repo, ".gitignore", "scratch/\n")
	gitT(t, repo, "add", "-A")
	gitT(t, repo, "commit", "-q", "-m", "ignore")
	n := unitNode("n1", "out/n1.json", "exit_code == 0 && file_exists('scratch/x.json')", nil)
	n.ContextOut = []string{"scratch/x.json"}
	if _, err := unitOrchestrix("scratch/x.json=ok").Execute(context.Background(), n, repo); err != nil {
		t.Fatal(err)
	}
	if got := readT(t, repo, "scratch/x.json"); got != "ok\n" {
		t.Fatalf("saída ignorada não aplicada: %q", got)
	}
}

// Arquivo ignorado fora de allowedWrites não entra na lista nem reprova o nó.
func TestUnitTaskIgnoredOutsideAllowedIsNotApplied(t *testing.T) {
	repo := newRepoT(t)
	writeT(t, repo, ".gitignore", "scratch/\n")
	gitT(t, repo, "add", "-A")
	gitT(t, repo, "commit", "-q", "-m", "ignore")
	n := unitNode("n1", "out/n1.json", "exit_code == 0", nil)
	if _, err := unitOrchestrix("n1.txt=ok|scratch/lixo.json=x").Execute(context.Background(), n, repo); err != nil {
		t.Fatal(err)
	}
	if pathExists(filepath.Join(repo, "scratch", "lixo.json")) {
		t.Fatal("arquivo ignorado fora do boundary foi aplicado")
	}
}

// Gate reprovado precisa mostrar o que o claude respondeu, para o diagnóstico.
func TestUnitTaskGateFailureShowsClaudeSummary(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "file_exists('nao_existe.txt')", nil)
	_, err := unitOrchestrix("n1.txt=ok").Execute(context.Background(), n, repo)
	if err == nil || !strings.Contains(err.Error(), "feito") {
		t.Fatalf("erro deveria citar a resposta do claude: %v", err)
	}
}

// Pasta em context_in: o conteúdo rastreado já está na worktree; o prompt só cita a pasta.
func TestUnitTaskAcceptsDirContextIn(t *testing.T) {
	repo := newRepoT(t)
	for _, f := range []string{"docs/", "docs"} {
		n := unitNode("n1", "out/n1.json", "exit_code == 0 && file_exists('n1.txt')", nil)
		n.ContextIn = []string{f}
		if _, err := unitOrchestrix("n1.txt=ok").Execute(context.Background(), n, repo); err != nil {
			t.Fatalf("context_in %q: %v", f, err)
		}
		_ = os.Remove(filepath.Join(repo, "n1.txt"))
		_ = os.RemoveAll(filepath.Join(repo, "out"))
	}
}

func TestUnitPromptNamesDirWithoutInjecting(t *testing.T) {
	repo := newRepoT(t)
	root, err := os.OpenRoot(repo)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	n := unitNode("n1", "out/n1.json", "exit_code == 0", nil)
	n.ContextIn = []string{"docs/"}
	p, err := buildUnitPrompt(n, root)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(p, "Pasta de contexto: docs") || strings.Contains(p, "b\n") {
		t.Fatalf("prompt = %q", p)
	}
}

func TestUnitTaskRejectsBadDirContextIn(t *testing.T) {
	repo := newRepoT(t)
	linkT(t, filepath.Join(repo, "lnk"), filepath.Join(repo, "docs"))
	for _, f := range []string{".git/", "lnk/", "../fora/", "docs/../../fora/"} {
		n := unitNode("n1", "out/n1.json", "exit_code == 0", nil)
		n.ContextIn = []string{f}
		if _, err := unitOrchestrix("n1.txt=ok").Execute(context.Background(), n, repo); err == nil {
			t.Errorf("pasta %q em context_in deveria reprovar", f)
		}
	}
	if c := worktreeCount(t, repo); c != 1 {
		t.Fatalf("worktree vazou: %d", c)
	}
}

// helperClaudeRetry falha na primeira tentativa (só escreve tentativa1.txt) e acerta quando o
// prompt traz a falha da tentativa anterior.
func helperClaudeRetry() {
	in, _ := io.ReadAll(os.Stdin)
	file := "tentativa1.txt"
	if strings.Contains(string(in), "Tentativa anterior falhou") {
		file = "n1.txt"
	}
	_ = os.WriteFile(file, []byte("ok\n"), 0o644)
	out, _ := json.Marshal(map[string]any{"type": "result", "result": "feito",
		"usage": map[string]int{"input_tokens": 5, "output_tokens": 2}})
	os.Stdout.Write(out)
	os.Exit(0)
}

func TestUnitTaskRetriesWithGateFailure(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "exit_code == 0 && file_exists('n1.txt')", nil)
	n.ContextOut = []string{"n1.txt", "tentativa1.txt"}
	o := unitOrchestrix("")
	o.ExtraEnv = []string{"LEGOLOOP_HELPER=claude-retry"}
	res, err := o.Execute(context.Background(), n, repo)
	if err != nil {
		t.Fatalf("a segunda tentativa deveria passar: %v", err)
	}
	if res.Tokens.Input != 10 {
		t.Errorf("tokens devem somar as tentativas: %+v", res.Tokens)
	}
	if got := readT(t, repo, "n1.txt"); got != "ok\n" {
		t.Errorf("n1.txt = %q", got)
	}
}

func TestUnitTaskStopsAfterMaxRetries(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "file_exists('nunca.txt')", nil)
	res, err := unitOrchestrix("n1.txt=x").Execute(context.Background(), n, repo)
	if err == nil || !strings.Contains(err.Error(), "reprovou") {
		t.Fatalf("err = %v", err)
	}
	if want := 5 * (1 + maxUnitRetries); res.Tokens.Input != want {
		t.Errorf("tokens = %d, esperado %d (1 + %d retentativas)", res.Tokens.Input, want, maxUnitRetries)
	}
	if pathExists(filepath.Join(repo, "n1.txt")) {
		t.Fatal("nó reprovado aplicou arquivo")
	}
	if !pathExists(filepath.Join(repo, "scratch", "gates", "n1.failed.patch")) {
		t.Errorf("deveria ter salvo patch das tentativas que falharam")
	}
}

func TestUnitTaskHygieneCleansCommentAndFormatsGo(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "exit_code == 0 && file_exists('pkg/code.go')", []string{"git", "status", "--short"})
	n.ContextOut = []string{"pkg/code.go"}
	comentario := "Esse lo" + "op de teste"
	claudeOutput := "pkg/code.go=package pkg\n\n// " + comentario + " faz algo\nfunc Hello(){\nfmt.Println(\"hi\")\n}\n"
	res, err := unitOrchestrix(claudeOutput).Execute(context.Background(), n, repo)
	if err != nil {
		t.Fatalf("execução falhou: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit code = %d", res.ExitCode)
	}
	applied := readT(t, repo, "pkg/code.go")
	if strings.Contains(applied, comentario) {
		t.Errorf("comentário com termo proibido deveria ter sido removido pela higiene")
	}
	if !strings.Contains(applied, "\tfmt.Println(\"hi\")") {
		t.Errorf("código Go deveria estar formatado com gofmt")
	}
}

func TestUnitTaskGateImprovementRetainsPartial(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "exit_code == 0", []string{os.Args[0], "gate-step"})
	n.ContextOut = []string{"edited.txt"}
	o := unitOrchestrix("edited.txt=nova versao")
	o.ExtraEnv = append(o.ExtraEnv, "GATE_MODE=improve")
	res, err := o.Execute(context.Background(), n, repo)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if res.Status != runner.StatusPartial || res.ExitCode != 1 {
		t.Fatalf("esperava StatusPartial e exit code 1, obteve: %+v", res)
	}
	if !pathExists(filepath.Join(repo, "edited.txt")) {
		t.Fatal("edição melhorada deveria ter sido mantida no repositório")
	}
	var rep unitReport
	if err := json.Unmarshal([]byte(readT(t, repo, "out/n1.json")), &rep); err != nil {
		t.Fatal(err)
	}
	if rep.Passed || rep.Status != "parcial" {
		t.Fatalf("relatório deveria indicar parcial: %+v", rep)
	}
}

func TestUnitTaskGateWorseDiscards(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "exit_code == 0", []string{os.Args[0], "gate-step"})
	n.ContextOut = []string{"edited.txt"}
	o := unitOrchestrix("edited.txt=piorou")
	o.ExtraEnv = append(o.ExtraEnv, "GATE_MODE=worse")
	_, err := o.Execute(context.Background(), n, repo)
	if err == nil || !strings.Contains(err.Error(), "reprovou") {
		t.Fatalf("esperava erro de reprovação do gate, obteve: %v", err)
	}
	if pathExists(filepath.Join(repo, "edited.txt")) {
		t.Fatal("edição que piorou não pode ser mantida")
	}
	if !pathExists(filepath.Join(repo, "scratch", "gates", "n1.failed.patch")) {
		t.Fatal("deveria ter salvo failed.patch")
	}
}

func TestUnitTaskGateBrokenBuildDiscards(t *testing.T) {
	repo := newRepoT(t)
	writeT(t, repo, "go.mod", "module testbroken\n\ngo 1.24\n")
	gitT(t, repo, "add", "go.mod")
	gitT(t, repo, "commit", "-q", "-m", "add go.mod")
	n := unitNode("n1", "out/n1.json", "exit_code == 0", []string{os.Args[0], "gate-step"})
	n.ContextOut = []string{"edited.txt", "broken.go"}
	o := unitOrchestrix("edited.txt=ok|broken.go=package main\ninvalid syntax !!!")
	o.ExtraEnv = append(o.ExtraEnv, "GATE_MODE=improve")
	_, err := o.Execute(context.Background(), n, repo)
	if err == nil {
		t.Fatal("esperava erro devido a build quebrado")
	}
	if pathExists(filepath.Join(repo, "edited.txt")) || pathExists(filepath.Join(repo, "broken.go")) {
		t.Fatal("código com build quebrado não pode ser mantido")
	}
}

func TestUnitTaskRetryReceivesPrioritizedOutput(t *testing.T) {
	var lines []string
	for i := 1; i <= 250; i++ {
		lines = append(lines, fmt.Sprintf("file_%d.go:%d:1: violacao taxonomica", i, i))
	}
	full := strings.Join(lines, "\n")
	formatted := formatGateOutput([]byte(full), 1000)
	if !strings.Contains(formatted, "file_1.go:1:1:") {
		t.Errorf("formatação priorizada deveria manter as primeiras violações, obteve:\n%s", formatted)
	}
	if !strings.Contains(formatted, "adicionais omitidas") {
		t.Errorf("esperava indicação de violações omitidas ao estourar o limite")
	}
}

func TestUnitTaskDiagnosticModePrompt(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "exit_code == 0", []string{"git", "status"})
	n.Verification = map[string]any{
		"diagnostic_command": []string{os.Args[0], "gate-step"},
	}
	o := unitOrchestrix("n1.txt=ok")
	o.ExtraEnv = append(o.ExtraEnv, "GATE_MODE=diag")
	prompt, score, err := o.prepareInitialPrompt(context.Background(), n, repo, "# Base Prompt")
	if err != nil {
		t.Fatalf("erro ao preparar prompt inicial: %v", err)
	}
	if score.violations != 3 {
		t.Errorf("esperava 3 violações detectadas, obteve: %+v", score)
	}
	if !strings.Contains(prompt, "Diagnóstico do Gate (violações encontradas)") {
		t.Errorf("prompt deveria conter seção de diagnóstico: %s", prompt)
	}
	if !strings.Contains(prompt, "Corrija todas as violações listadas") {
		t.Errorf("prompt deveria conter instrução de correção e releitura: %s", prompt)
	}
}

func TestUnitTaskPromptQualityInstruction(t *testing.T) {
	root, err := os.OpenRoot(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	n := unitNode("n1", "out/n1.json", "exit_code == 0", nil)
	prompt, err := buildUnitPrompt(n, root)
	if err != nil {
		t.Fatal(err)
	}
	want := "Corrija todas as violações listadas, inclusive em comentários e em arquivos _test.go, e releia a lista no fim."
	if !strings.Contains(prompt, want) {
		t.Fatalf("prompt base deve conter instrução curta de qualidade:\n%s", prompt)
	}
}
