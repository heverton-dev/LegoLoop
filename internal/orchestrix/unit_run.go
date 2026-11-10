package orchestrix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/dag"
	"legoloop/internal/gate"
	"legoloop/internal/provider"
	"legoloop/internal/runner"
)

// unitReport é o conteúdo JSON de output.artifact de um unit_task aprovado.
type unitReport struct {
	Node     string   `json:"node"`
	Model    string   `json:"model"`
	Changed  []string `json:"changed"`
	ExitCode int      `json:"exit_code"`
	Gate     string   `json:"gate"`
	Passed   bool     `json:"passed"`
	Summary  string   `json:"summary,omitempty"`
}

// unitPlan reúne o que o unit_task valida antes de criar a worktree.
type unitPlan struct {
	model  string
	expr   *gate.Expr
	prompt string
}

// prepareUnit valida modelo, gate, raiz do repositório e monta o prompt.
func (o *Orchestrix) prepareUnit(ctx context.Context, n dag.Node, workdir string, main *os.Root) (unitPlan, error) {
	var p unitPlan
	model, err := o.unitModel(n)
	if err != nil {
		return p, err
	}
	expr, err := gate.Compile(n.Output.Gate)
	if err != nil {
		return p, err
	}
	if err := requireRepoRoot(ctx, workdir); err != nil {
		return p, err
	}
	prompt, err := buildUnitPrompt(n, main)
	return unitPlan{model: model, expr: expr, prompt: prompt}, err
}

// runUnitTask executa um unit_task: worktree efêmera, claude com ferramentas fechadas,
// write_boundary, command, gate do nó e, só se o gate passar, aplicação no workdir principal.
func (o *Orchestrix) runUnitTask(ctx context.Context, n dag.Node, workdir string) (runner.Result, error) {
	fail := func(err error) (runner.Result, error) {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	if !o.AllowUnitTasks {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q é unit_task; rode com --allow-unit-tasks se confia neste pipeline", n.ID)
	}
	main, err := os.OpenRoot(workdir)
	if err != nil {
		return runner.Result{}, err
	}
	defer main.Close()
	plan, err := o.prepareUnit(ctx, n, workdir, main)
	if err != nil {
		return fail(err)
	}
	var res runner.Result
	var report unitReport
	err = RunInWorktree(ctx, workdir, func(w *Worktree) error {
		var err error
		if res, report, err = o.unitInWorktree(ctx, n, w, plan); err != nil {
			return err
		}
		return w.Apply(report.Changed)
	})
	if err != nil {
		return res, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return res, err
	}
	return res, writeArtifact(main, n.Output.Artifact, data)
}

// unitInWorktree roda claude e command na worktree e decide o gate. Não escreve no principal.
func (o *Orchestrix) unitInWorktree(ctx context.Context, n dag.Node, w *Worktree, plan unitPlan) (runner.Result, unitReport, error) {
	seeded, err := seedContext(w, n)
	if err != nil {
		return runner.Result{}, unitReport{}, err
	}
	prompt, err := o.promptWithRed(ctx, n, w.Dir, plan.prompt)
	if err != nil {
		return runner.Result{}, unitReport{}, err
	}
	res := runner.Result{Model: "unit_task/" + plan.model}
	base := prompt
	for attempt := 0; ; attempt++ {
		report, failure, tokens, err := o.unitAttempt(ctx, n, w, plan, prompt, seeded)
		res.Tokens = addUsage(res.Tokens, tokens)
		res.ExitCode = report.ExitCode
		if err != nil {
			return res, unitReport{}, err
		}
		if failure == "" {
			return res, report, nil
		}
		if attempt >= maxUnitRetries {
			return res, unitReport{}, fmt.Errorf("%s (após %d tentativa(s))", failure, attempt+1)
		}
		dropProvisionalArtifact(w, n)
		prompt = retryPrompt(base, failure)
	}
}

// dropProvisionalArtifact remove da worktree o artefato provisório que o gate gravou.
// Sem isso a próxima tentativa o enxergaria como alteração fora do write_boundary.
func dropProvisionalArtifact(w *Worktree, n dag.Node) {
	if n.Output.Artifact == "" {
		return
	}
	if root, err := os.OpenRoot(w.Dir); err == nil {
		defer root.Close()
		_ = root.Remove(filepath.FromSlash(n.Output.Artifact))
	}
}

// maxUnitRetries é o número de novas tentativas do claude, na mesma worktree, após falha do gate.
const maxUnitRetries = 2

// addUsage soma o consumo de tokens de duas chamadas.
func addUsage(a, b runner.TokenUsage) runner.TokenUsage {
	return runner.TokenUsage{Input: a.Input + b.Input, Output: a.Output + b.Output, ReportedInput: a.ReportedInput + b.ReportedInput}
}

// retryPrompt repete o prompt base e anexa a falha da tentativa anterior.
func retryPrompt(base, failure string) string {
	return base + "\n## Tentativa anterior falhou\n\n" + failure + "\n\nCorrija os arquivos permitidos. As suas edições anteriores continuam na pasta de trabalho.\n"
}

// unitAttempt roda o claude, o command e o gate uma vez. A falha do gate volta em failure
// (texto vazio = gate aprovado). Erro de boundary, link ou execução volta em err e não retenta.
func (o *Orchestrix) unitAttempt(ctx context.Context, n dag.Node, w *Worktree, plan unitPlan, prompt string, seeded map[string][]byte) (unitReport, string, runner.TokenUsage, error) {
	text, usage, err := o.callClaudeTools(ctx, plan.model, prompt, w.Dir, nil)
	if err != nil {
		return unitReport{}, "", usage, err
	}
	changed, err := unitChanged(ctx, n, w, seeded)
	if err != nil {
		return unitReport{}, "", usage, err
	}
	cmdOut, code, err := o.unitCommand(ctx, n, w.Dir)
	if err != nil {
		return unitReport{}, "", usage, err
	}
	report := unitReport{Node: n.ID, Model: plan.model, Changed: changed, ExitCode: code,
		Gate: n.Output.Gate, Passed: true, Summary: firstLine(text)}
	ok, err := o.unitGate(w, n, plan.expr, report)
	if err != nil {
		return report, "", usage, fmt.Errorf("gate %q: %w", n.Output.Gate, err)
	}
	if !ok {
		return report, fmt.Sprintf("gate %q reprovou (exit_code=%d); worktree descartada; claude respondeu: %s; saída do command: %s",
			n.Output.Gate, code, firstLine(text), tailText(cmdOut, maxCmdTail)), usage, nil
	}
	return report, "", usage, nil
}

// unitGate avalia o gate com a worktree como raiz. O relatório provisório entra na worktree
// só para o gate poder exigir o artefato; a lista de alterados já está fechada e não o inclui.
func (o *Orchestrix) unitGate(w *Worktree, n dag.Node, expr *gate.Expr, report unitReport) (bool, error) {
	root, err := os.OpenRoot(w.Dir)
	if err != nil {
		return false, err
	}
	defer root.Close()
	if data, err := json.Marshal(report); err == nil && n.Output.Artifact != "" {
		if err := writeArtifact(root, n.Output.Artifact, data); err != nil {
			return false, err
		}
	}
	return expr.Eval(gate.Env{Root: root, ExitCode: report.ExitCode, Vars: o.Vars})
}

// unitCommand roda o command do nó na worktree, sem shell. Sem command o código de saída é 0.
func (o *Orchestrix) unitCommand(ctx context.Context, n dag.Node, dir string) ([]byte, int, error) {
	if len(n.Command) == 0 {
		return nil, 0, nil
	}
	bin, err := resolveExecutable(n.Command[0])
	if err != nil {
		return nil, 0, err
	}
	return o.spawn(ctx, bin, o.expandArgs(n.Command[1:]), dir, nil)
}

// maxRedBytes limita a saída vermelha do command que entra no prompt (fica a cauda).
const maxRedBytes = 20000

// promptWithRed roda o command na worktree antes do claude (Red do TDD) e anexa a falha ao prompt.
func (o *Orchestrix) promptWithRed(ctx context.Context, n dag.Node, dir, prompt string) (string, error) {
	out, code, err := o.unitCommand(ctx, n, dir)
	if err != nil {
		return "", err
	}
	return withRedOutput(prompt, out, code), nil
}

// withRedOutput anexa a saída do command ao prompt só quando ele falhou e imprimiu algo.
func withRedOutput(prompt string, out []byte, code int) string {
	text := tailText(out, maxRedBytes)
	if code == 0 || text == "" {
		return prompt
	}
	return prompt + "\n## Saída da verificação (falha atual)\n\n```\n" + text + "\n```\n"
}

// maxCmdTail limita a saída do command que entra na mensagem de erro do gate.
const maxCmdTail = 1500

// tailText devolve os últimos max bytes do texto, sem espaços nas pontas.
func tailText(b []byte, max int) string {
	s := strings.TrimSpace(string(b))
	if len(s) > max {
		s = "..." + s[len(s)-max:]
	}
	return s
}

// unitModel escolhe o primeiro modelo do roteamento que o claude CLI executa.
func (o *Orchestrix) unitModel(n dag.Node) (string, error) {
	r := n.Routing
	r.Model, r.Fallback = unitAlias(r.Model), unitAlias(r.Fallback)
	routes, err := o.routes(r)
	if err != nil {
		return "", err
	}
	for _, rt := range routes {
		if rt.src == provider.SourceClaudeCLI {
			return rt.model, nil
		}
	}
	return "", errors.New("unit_task exige um modelo do claude CLI no roteamento")
}

// unitModelAliases traduz os combos do 9Router usados nos DAGs para modelos do claude CLI.
// O unit_task só executa claude CLI com ferramentas fechadas; combo do 9Router não roda aqui.
var unitModelAliases = map[string]string{"code-fast": "sonnet", "code-pro": "opus", "flash": "haiku"}

// unitAlias devolve o modelo do claude CLI para um alias; nome sem alias volta igual.
func unitAlias(m string) string {
	if a, ok := unitModelAliases[m]; ok {
		return a
	}
	return m
}

// requireRepoRoot exige que workdir seja a raiz do repositório git: a worktree espelha a raiz.
func requireRepoRoot(ctx context.Context, workdir string) error {
	out, err := gitRun(ctx, workdir, "rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("workdir não é repositório git: %w", err)
	}
	top := strings.TrimSpace(string(out))
	a, errA := os.Stat(filepath.FromSlash(top))
	b, errB := os.Stat(workdir)
	if errA != nil || errB != nil || !os.SameFile(a, b) {
		return fmt.Errorf("workdir %q não é a raiz do repositório git (%s)", workdir, top)
	}
	return nil
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(s), "\n")
	if len(line) > 200 {
		line = line[:200]
	}
	return line
}
