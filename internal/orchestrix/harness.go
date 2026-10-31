package orchestrix

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/provider"
	"legoloop/internal/runner"
)

// maxArgPromptBytes: o agy só recebe o prompt como argumento, e o Windows limita a
// linha de comando a ~32 KB; acima disso o nó falha com mensagem clara.
const maxArgPromptBytes = 24 << 10

// promptFile é o anexo gravado no sandbox para harnesses que aceitam arquivos.
const promptFile = "prompt.md"

const harnessTask = "Execute exatamente a tarefa descrita no arquivo anexado prompt.md."

// harnessSpec descreve como um harness recebe o prompt e como sua saída é lida.
type harnessSpec struct {
	args  func(model, prompt string) ([]string, bool, error)
	parse func(out []byte) (string, int, int, error)
}

func harnessSpecs() map[provider.Source]harnessSpec {
	return map[provider.Source]harnessSpec{
		provider.SourceAgy:      {agyArgs, parseAgy},
		provider.SourceOmp:      {ompArgs, parseOmp},
		provider.SourceMimo:     {opencodeArgs, parseOpencodeEvents},
		provider.SourceOpencode: {opencodeArgs, parseOpencodeEvents},
	}
}

// callHarness roda o harness num sandbox vazio (sem config do repositório) e lê a resposta.
func (o *Orchestrix) callHarness(ctx context.Context, src provider.Source, model, prompt string) (string, runner.TokenUsage, error) {
	spec, ok := harnessSpecs()[src]
	if !ok {
		return "", runner.TokenUsage{}, fmt.Errorf("executor desconhecido %q", src)
	}
	args, withFile, err := spec.args(harnessModel(model), prompt)
	if err != nil {
		return "", runner.TokenUsage{}, err
	}
	sandbox, err := os.MkdirTemp("", "legoloop-orchestrix-*")
	if err != nil {
		return "", runner.TokenUsage{}, err
	}
	defer os.RemoveAll(sandbox)
	if withFile {
		if err := os.WriteFile(filepath.Join(sandbox, promptFile), []byte(prompt), 0o600); err != nil {
			return "", runner.TokenUsage{}, err
		}
	}
	out, code, err := o.spawn(ctx, o.Inventory.CLIPath(string(src)), args, sandbox, nil)
	if err != nil {
		return "", runner.TokenUsage{}, err
	}
	text, in, outTok, perr := spec.parse(out)
	if perr == nil && code != 0 {
		perr = errors.New("resposta recebida, mas o processo falhou")
	}
	if perr != nil {
		return "", runner.TokenUsage{}, fmt.Errorf("%s saiu com código %d: %w", src, code, perr)
	}
	return text, runner.TokenUsage{Input: in, Output: outTok}, nil
}

// harnessModel tira o prefixo ("agy:x" ⇒ "x"); "default" ⇒ "" (modelo padrão do harness).
func harnessModel(model string) string {
	_, m, _ := strings.Cut(model, ":")
	if m == "default" {
		return ""
	}
	return m
}

func withModel(args []string, flag, model string) []string {
	if model == "" {
		return args
	}
	return append(args, flag, model)
}

func agyArgs(model, prompt string) ([]string, bool, error) {
	if len(prompt) > maxArgPromptBytes {
		return nil, false, fmt.Errorf("prompt de %d bytes excede o limite de %d do agy (só aceita prompt por argumento)", len(prompt), maxArgPromptBytes)
	}
	args := []string{"-p=" + prompt, "--output-format", "json", "--disable-slash-commands"}
	return withModel(args, "--model", model), false, nil
}

func ompArgs(model, prompt string) ([]string, bool, error) {
	args := []string{"-p", "--no-session", "--no-tools", "--no-skills", "--no-rules", "--no-extensions", "--mode", "json"}
	return append(withModel(args, "--model", model), "@"+promptFile, harnessTask), true, nil
}

// opencodeArgs serve opencode e mimo (mimocode é um fork com a mesma CLI).
func opencodeArgs(model, prompt string) ([]string, bool, error) {
	args := withModel([]string{"run", harnessTask, "--format", "json", "--pure"}, "-m", model)
	return append(args, "-f", promptFile), true, nil
}

func parseAgy(out []byte) (string, int, int, error) {
	var r struct {
		Status   string `json:"status"`
		Response string `json:"response"`
		Usage    struct {
			Input  int `json:"input_tokens"`
			Output int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(out), &r); err != nil {
		return "", 0, 0, fmt.Errorf("saída do agy não é JSON: %s", clip(string(out)))
	}
	if r.Status != "SUCCESS" {
		return "", 0, 0, fmt.Errorf("agy %s: %s", r.Status, clip(r.Response))
	}
	return r.Response, r.Usage.Input, r.Usage.Output, nil
}

type ompMessage struct {
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
	StopReason   string `json:"stopReason"`
	ErrorMessage string `json:"errorMessage"`
}

// parseOmp lê o evento agent_end; o omp sai com 0 mesmo em erro, então stopReason decide.
func parseOmp(out []byte) (string, int, int, error) {
	var last *ompMessage
	for _, line := range ndjson(out) {
		var ev struct {
			Type     string       `json:"type"`
			Messages []ompMessage `json:"messages"`
		}
		if json.Unmarshal(line, &ev) != nil || ev.Type != "agent_end" {
			continue
		}
		for i := range ev.Messages {
			if ev.Messages[i].Role == "assistant" {
				last = &ev.Messages[i]
			}
		}
	}
	if last == nil {
		return "", 0, 0, fmt.Errorf("omp não emitiu agent_end: %s", clip(string(out)))
	}
	if last.StopReason == "error" {
		return "", 0, 0, fmt.Errorf("omp: %s", last.ErrorMessage)
	}
	var b strings.Builder
	for _, c := range last.Content {
		if c.Type == "text" {
			b.WriteString(c.Text)
		}
	}
	return nonEmpty(b.String(), last.Usage.Input, last.Usage.Output, "omp")
}

type opencodeEvent struct {
	Type string `json:"type"`
	Part struct {
		Text   string `json:"text"`
		Tokens struct {
			Input  int `json:"input"`
			Output int `json:"output"`
		} `json:"tokens"`
	} `json:"part"`
	Error json.RawMessage `json:"error"`
}

// parseOpencodeEvents concatena os eventos "text" e soma os tokens de "step_finish".
func parseOpencodeEvents(out []byte) (string, int, int, error) {
	var b strings.Builder
	in, outTok := 0, 0
	for _, line := range ndjson(out) {
		var ev opencodeEvent
		if json.Unmarshal(line, &ev) != nil {
			continue
		}
		switch ev.Type {
		case "text":
			b.WriteString(ev.Part.Text)
		case "step_finish":
			in, outTok = in+ev.Part.Tokens.Input, outTok+ev.Part.Tokens.Output
		case "error":
			return "", 0, 0, fmt.Errorf("%s", eventErrorMessage(ev.Error))
		}
	}
	return nonEmpty(b.String(), in, outTok, "harness")
}

// eventErrorMessage extrai error.data.message, error.message ou o JSON bruto.
func eventErrorMessage(raw json.RawMessage) string {
	var e struct {
		Message string `json:"message"`
		Data    struct {
			Message string `json:"message"`
		} `json:"data"`
	}
	_ = json.Unmarshal(raw, &e)
	for _, m := range []string{e.Data.Message, e.Message} {
		if m != "" {
			return m
		}
	}
	return clip(string(raw))
}

func ndjson(out []byte) [][]byte {
	var lines [][]byte
	sc := bufio.NewScanner(bytes.NewReader(out))
	sc.Buffer(make([]byte, 0, 64<<10), maxResponseBytes)
	for sc.Scan() {
		if l := bytes.TrimSpace(sc.Bytes()); len(l) > 0 {
			lines = append(lines, append([]byte(nil), l...))
		}
	}
	return lines
}

func nonEmpty(text string, in, out int, who string) (string, int, int, error) {
	if strings.TrimSpace(text) == "" {
		return "", 0, 0, fmt.Errorf("%s não devolveu texto", who)
	}
	return text, in, out, nil
}

func clip(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 300 {
		return s[:300] + "…"
	}
	if s == "" {
		return "sem detalhes"
	}
	return s
}
