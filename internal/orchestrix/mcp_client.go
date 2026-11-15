package orchestrix

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/runner"
)

const defaultMCPTimeout = 60 * time.Second

// MCPClient executa chamadas JSON-RPC 2.0 sobre STDIO contra um servidor MCP.
type MCPClient struct {
	Command        string
	Args           []string
	Dir            string
	Env            []string
	Timeout        time.Duration
	MaxOutputBytes int
}

type mcpRPCEnvelope struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  struct {
		IsError bool `json:"isError"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"result"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// CallTool realiza handshake initialize + notifications/initialized e invoca tools/call via STDIO.
func (c *MCPClient) CallTool(ctx context.Context, tool string, args map[string]any) (string, error) {
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = defaultMCPTimeout
	}
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd, stdin, stdout, stderr, err := c.startProcess(callCtx)
	if err != nil {
		return "", err
	}
	waited := false
	defer func() {
		if !waited {
			cleanupMCPProcess(cmd, stdin)
		}
	}()

	reader := bufio.NewReader(stdout)
	if err := c.performHandshake(callCtx, stdin, reader); err != nil {
		return "", finishMCPError(cmd, stdin, &waited, stderr, err)
	}
	out, err := c.invokeTool(callCtx, stdin, reader, tool, args)
	if err != nil {
		return "", finishMCPError(cmd, stdin, &waited, stderr, err)
	}
	_ = stdin.Close()
	_ = cmd.Wait()
	waited = true
	return out, nil
}

func (c *MCPClient) startProcess(ctx context.Context) (*exec.Cmd, io.WriteCloser, io.ReadCloser, *limitedBuffer, error) {
	cmd := exec.CommandContext(ctx, c.Command, c.Args...)
	cmd.Dir = c.Dir
	cmd.Env = c.Env
	if len(cmd.Env) == 0 {
		cmd.Env = sterileEnv(os.Getenv)
	}
	cmd.WaitDelay = spawnWaitDelay
	limit := c.MaxOutputBytes
	if limit <= 0 {
		limit = maxResponseBytes
	}
	stderr := &limitedBuffer{limit: limit}
	cmd.Stderr = stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return nil, nil, nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return nil, nil, nil, nil, fmt.Errorf("iniciar servidor MCP: %w%s", err, formatStderr(stderr))
	}
	return cmd, stdin, stdout, stderr, nil
}

func cleanupMCPProcess(cmd *exec.Cmd, stdin io.Closer) {
	if stdin != nil {
		_ = stdin.Close()
	}
	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}
}

func finishMCPError(cmd *exec.Cmd, stdin io.Closer, waited *bool, stderr *limitedBuffer, err error) error {
	if stdin != nil {
		_ = stdin.Close()
	}
	if cmd != nil && !*waited {
		waitCh := make(chan error, 1)
		go func() { waitCh <- cmd.Wait() }()
		select {
		case <-waitCh:
		case <-time.After(500 * time.Millisecond):
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
			<-waitCh
		}
		*waited = true
	}
	return fmt.Errorf("%w%s", err, formatStderr(stderr))
}

func (c *MCPClient) performHandshake(ctx context.Context, stdin io.Writer, reader *bufio.Reader) error {
	initReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo":      map[string]any{"name": "legoloop-orchestrix", "version": "1.0.0"},
		},
	}
	if err := sendJSONLine(stdin, initReq); err != nil {
		return fmt.Errorf("enviar initialize MCP: %w", err)
	}
	resp, err := readRPCResponse(ctx, reader)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf("erro initialize MCP (%d): %s", resp.Error.Code, resp.Error.Message)
	}
	notif := map[string]any{"jsonrpc": "2.0", "method": "notifications/initialized"}
	if err := sendJSONLine(stdin, notif); err != nil {
		return fmt.Errorf("enviar notifications/initialized MCP: %w", err)
	}
	return nil
}

func (c *MCPClient) invokeTool(ctx context.Context, stdin io.Writer, reader *bufio.Reader, tool string, args map[string]any) (string, error) {
	if args == nil {
		args = map[string]any{}
	}
	callReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params":  map[string]any{"name": tool, "arguments": args},
	}
	if err := sendJSONLine(stdin, callReq); err != nil {
		return "", fmt.Errorf("enviar tools/call MCP (%s): %w", tool, err)
	}
	resp, err := readRPCResponse(ctx, reader)
	if err != nil {
		return "", err
	}
	return extractToolOutput(tool, resp)
}

func sendJSONLine(w io.Writer, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(append(b, '\n'))
	return err
}

func readRPCResponse(ctx context.Context, reader *bufio.Reader) (mcpRPCEnvelope, error) {
	for {
		line, err := reader.ReadBytes('\n')
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) > 0 {
			var env mcpRPCEnvelope
			if jsonErr := json.Unmarshal(trimmed, &env); jsonErr != nil {
				return mcpRPCEnvelope{}, fmt.Errorf("resposta JSON-RPC MCP invalida: %w", jsonErr)
			}
			return env, nil
		}
		if err != nil {
			if ctx.Err() != nil {
				return mcpRPCEnvelope{}, fmt.Errorf("timeout ou cancelamento MCP: %w", ctx.Err())
			}
			return mcpRPCEnvelope{}, fmt.Errorf("leitura STDIO MCP encerrada: %w", err)
		}
	}
}

func extractToolOutput(tool string, resp mcpRPCEnvelope) (string, error) {
	if resp.Error != nil {
		return "", fmt.Errorf("erro JSON-RPC na ferramenta %q (%d): %s", tool, resp.Error.Code, resp.Error.Message)
	}
	var b strings.Builder
	for i, c := range resp.Result.Content {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(c.Text)
	}
	text := b.String()
	if resp.Result.IsError {
		return "", fmt.Errorf("ferramenta MCP %q falhou: %s", tool, text)
	}
	return text, nil
}

func formatStderr(stderr *limitedBuffer) string {
	if stderr == nil {
		return ""
	}
	s := strings.TrimSpace(string(stderr.Bytes()))
	if s == "" {
		return ""
	}
	return " [stderr: " + clip(s) + "]"
}

func (o *Orchestrix) runMCP(ctx context.Context, n dag.Node, workdir string) (runner.Result, error) {
	bin, args, err := o.resolveMCPCommand(n)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	tool, toolArgs := o.buildMCPCallArgs(n)
	timeout := time.Duration(n.TimeoutSec) * time.Second
	client := &MCPClient{
		Command:        bin,
		Args:           args,
		Dir:            workdir,
		Env:            append(sterileEnv(os.Getenv), o.ExtraEnv...),
		Timeout:        timeout,
		MaxOutputBytes: o.MaxOutputBytes,
	}
	out, err := client.CallTool(ctx, tool, toolArgs)
	if err != nil {
		return runner.Result{ExitCode: 1, Model: "mcp/" + n.Ref}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	if err := writeMCPArtifactIfNeeded(workdir, n.Output.Artifact, []byte(out)); err != nil {
		return runner.Result{ExitCode: 1, Model: "mcp/" + n.Ref}, err
	}
	return runner.Result{ExitCode: 0, Model: "mcp/" + n.Ref}, nil
}

func (o *Orchestrix) resolveMCPCommand(n dag.Node) (string, []string, error) {
	if len(n.Command) > 0 {
		bin, err := resolveExecutable(n.Command[0])
		if err != nil {
			return "", nil, err
		}
		return bin, o.expandArgs(n.Command[1:]), nil
	}
	bin, err := resolveExecutable("legoloop")
	if err != nil {
		return "", nil, errors.New("nó MCP sem command explícito e binário legoloop não encontrado")
	}
	return bin, []string{"mcp-router"}, nil
}

func (o *Orchestrix) buildMCPCallArgs(n dag.Node) (string, map[string]any) {
	args := make(map[string]any, len(n.Inputs))
	for k, v := range n.Inputs {
		expanded := o.expandArgs([]string{v})
		args[k] = expanded[0]
	}
	tool := strings.TrimSpace(n.Ref)
	if srv, sub, ok := strings.Cut(tool, "/"); ok {
		if len(n.Command) == 0 {
			return "call_mcp_tool", map[string]any{"server": srv, "tool": sub, "arguments": args}
		}
		return sub, args
	}
	return tool, args
}

func writeMCPArtifactIfNeeded(workdir, artifact string, data []byte) error {
	root, err := os.OpenRoot(workdir)
	if err != nil {
		return err
	}
	defer root.Close()
	if !artifactExists(root, artifact) {
		return writeArtifact(root, artifact, data)
	}
	return nil
}
