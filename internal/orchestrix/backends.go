package orchestrix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/runner"
)

// maxResponseBytes limita a resposta de qualquer backend.
const maxResponseBytes = 32 << 20

// routerTokenOverhead: o normalizador do 9Router soma +2000 FALSOS a prompt_tokens
// (código do servidor, 03/10). Medido em 31/10 na VPS com Llama 3.1 idêntico ao Ollama
// local: excesso fixo total = 2.405, dos quais ~405 são tokens REAIS injetados (saver +
// template) e ficam na conta; a contagem por token é exata (docs/references/achados_9router.md).
const routerTokenOverhead = 2000

func (o *Orchestrix) callRouter(ctx context.Context, combo, prompt string) (string, runner.TokenUsage, error) {
	req := map[string]any{"model": combo, "stream": false,
		"messages": []map[string]string{{"role": "user", "content": prompt}}}
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			Prompt     int `json:"prompt_tokens"`
			Completion int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := o.postJSON(ctx, o.Inventory.RouterURL+"/v1/chat/completions", o.routerCredential(), req, &resp); err != nil {
		return "", runner.TokenUsage{}, err
	}
	if len(resp.Choices) == 0 {
		return "", runner.TokenUsage{}, errors.New("9router: resposta sem choices")
	}
	in := max(resp.Usage.Prompt-routerTokenOverhead, 0)
	usage := runner.TokenUsage{Input: in, Output: resp.Usage.Completion, ReportedInput: resp.Usage.Prompt}
	return resp.Choices[0].Message.Content, usage, nil
}

// routerCredential só libera a chave para a base explicitamente confiável.
func (o *Orchestrix) routerCredential() string {
	if o.RouterKeyURL != "" && o.Inventory.RouterURL == o.RouterKeyURL {
		return o.RouterKey
	}
	return ""
}

func (o *Orchestrix) callOllama(ctx context.Context, model, prompt string) (string, runner.TokenUsage, error) {
	req := map[string]any{"model": model, "prompt": prompt, "stream": false}
	var resp struct {
		Response   string `json:"response"`
		PromptEval int    `json:"prompt_eval_count"`
		Eval       int    `json:"eval_count"`
	}
	if err := o.postJSON(ctx, o.OllamaURL+"/api/generate", "", req, &resp); err != nil {
		return "", runner.TokenUsage{}, err
	}
	return resp.Response, runner.TokenUsage{Input: resp.PromptEval, Output: resp.Eval}, nil
}

func (o *Orchestrix) postJSON(ctx context.Context, url, bearer string, body, dst any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	client := o.HTTP
	if client == nil {
		client = &http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status HTTP %d", resp.StatusCode)
	}
	return json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes)).Decode(dst)
}

// callClaude abre um `claude -p` efêmero e isolado: cwd é um diretório temporário vazio
// (o repositório nunca injeta .claude/settings.json, hooks ou CLAUDE.md), só configurações
// do usuário, modo seguro, sem ferramentas, sem MCP, sem sessão persistida; prompt via stdin.
func (o *Orchestrix) callClaude(ctx context.Context, model, prompt string) (string, runner.TokenUsage, error) {
	sandbox, err := os.MkdirTemp("", "legoloop-orchestrix-*")
	if err != nil {
		return "", runner.TokenUsage{}, err
	}
	defer os.RemoveAll(sandbox)
	args := []string{"-p", "--model", model, "--output-format", "json", "--no-session-persistence",
		"--strict-mcp-config", "--safe-mode", "--setting-sources", "user", "--tools", "", "--disable-slash-commands"}
	out, code, err := o.spawn(ctx, o.Inventory.CLIPath("claude"), args, sandbox, strings.NewReader(prompt))
	if err != nil {
		return "", runner.TokenUsage{}, err
	}
	r := parseClaudeOutput(out)
	if code != 0 || r.isError {
		return "", r.usage, fmt.Errorf("claude saiu com código %d: %s", code, r.reason())
	}
	return r.text, r.usage, nil
}

// claudeResult é o JSON de `claude -p --output-format json` (ou texto puro, se não for JSON).
type claudeResult struct {
	text     string
	usage    runner.TokenUsage
	isError  bool
	terminal string
}

// reason resume o motivo da falha em no máximo 300 caracteres.
func (r claudeResult) reason() string {
	msg := strings.TrimSpace(r.text)
	if r.terminal != "" {
		msg = r.terminal + ": " + msg
	}
	if len(msg) > 300 {
		msg = msg[:300] + "…"
	}
	if msg == "" {
		return "sem detalhes na saída"
	}
	return msg
}

func parseClaudeOutput(out []byte) claudeResult {
	var r struct {
		Result   string `json:"result"`
		IsError  bool   `json:"is_error"`
		Terminal string `json:"terminal_reason"`
		Usage    struct {
			Input  int `json:"input_tokens"`
			Output int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(out, &r); err != nil {
		return claudeResult{text: string(out)}
	}
	return claudeResult{text: r.Result, isError: r.IsError, terminal: r.Terminal,
		usage: runner.TokenUsage{Input: r.Usage.Input, Output: r.Usage.Output}}
}

// runScript executa o comando SEM shell (args literais, imunes a injeção) e só com
// consentimento explícito. Se o script não gravar o artefato, o stdout vira o artefato.
func (o *Orchestrix) runScript(ctx context.Context, n dag.Node, workdir string) (runner.Result, error) {
	if !o.AllowScripts {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q é script (%s); rode com --allow-scripts se confia neste pipeline", n.ID, n.Command[0])
	}
	bin, err := resolveExecutable(n.Command[0])
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	args := o.expandArgs(n.Command[1:])
	out, code, err := o.spawn(ctx, bin, args, workdir, nil)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	root, err := os.OpenRoot(workdir)
	if err != nil {
		return runner.Result{}, err
	}
	defer root.Close()
	if !artifactExists(root, n.Output.Artifact) {
		err = writeArtifact(root, n.Output.Artifact, out)
	}
	return runner.Result{ExitCode: code, Model: "script/" + n.Command[0]}, err
}

func resolveExecutable(name string) (string, error) {
	if name == "legoloop" {
		if exe, err := os.Executable(); err == nil && strings.Contains(strings.ToLower(exe), "legoloop") {
			return exe, nil
		}
	}
	return exec.LookPath(name)
}

func (o *Orchestrix) expandArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}
	out := make([]string, len(args))
	for i, a := range args {
		out[i] = os.Expand(a, func(k string) string {
			if v, ok := o.Vars[k]; ok {
				return v
			}
			return os.Getenv(k)
		})
	}
	return out
}

// spawnWaitDelay: após cancelamento/saída, quanto esperar por netos que seguram os pipes.
const spawnWaitDelay = 2 * time.Second

// spawn roda um subprocesso efêmero com ambiente estéril e stdout limitado.
func (o *Orchestrix) spawn(ctx context.Context, bin string, args []string, dir string, stdin io.Reader) ([]byte, int, error) {
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Env = append(sterileEnv(os.Getenv), o.ExtraEnv...)
	cmd.Stdin = stdin
	cmd.WaitDelay = spawnWaitDelay
	stdout := &limitedBuffer{limit: o.MaxOutputBytes}
	if stdout.limit <= 0 {
		stdout.limit = maxResponseBytes
	}
	cmd.Stdout = stdout
	err := cmd.Run()
	if stdout.overflow {
		return nil, -1, fmt.Errorf("saída do subprocesso excede o limite de %d bytes", stdout.limit)
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && ctx.Err() == nil {
		return stdout.Bytes(), exitErr.ExitCode(), nil
	}
	if err != nil {
		return nil, -1, err
	}
	return stdout.Bytes(), 0, nil
}

// sterileEnv repassa só o mínimo para o SO e o harness localizarem config/autenticação.
// NoDefaultCurrentDirectoryInExePath impede o cmd.exe (shims .cmd) de achar binários no cwd.
func sterileEnv(getenv func(string) string) []string {
	keys := []string{"PATH", "PATHEXT", "SYSTEMROOT", "SYSTEMDRIVE", "WINDIR", "COMSPEC",
		"HOME", "USERPROFILE", "APPDATA", "LOCALAPPDATA", "TEMP", "TMP", "TMPDIR", "LANG"}
	env := []string{"NoDefaultCurrentDirectoryInExePath=1"}
	for _, k := range keys {
		if v := getenv(k); v != "" {
			env = append(env, k+"="+v)
		}
	}
	return env
}

// limitedBuffer descarta o excesso SEM devolver erro: o filho nunca bloqueia no pipe cheio.
// O buffer NÃO é embutido de propósito: bytes.Buffer expõe ReadFrom, e io.Copy usaria
// esse atalho ignorando o limite deste Write.
type limitedBuffer struct {
	buf      bytes.Buffer
	limit    int
	overflow bool
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.overflow || b.buf.Len()+len(p) > b.limit {
		b.overflow = true
		return len(p), nil
	}
	return b.buf.Write(p)
}

func (b *limitedBuffer) Bytes() []byte { return b.buf.Bytes() }
