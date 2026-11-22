package orchestrix

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/runner"
)

func TestHelperMCPServerProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_MCP_SERVER") != "1" {
		return
	}
	mode := os.Getenv("MCP_HELPER_MODE")
	if mode == "sleep" {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
	if mode == "crash_stderr" {
		fmt.Fprintln(os.Stderr, "erro fatal simulado no stderr do servidor MCP")
		os.Exit(2)
	}

	reader := bufio.NewReader(os.Stdin)
	initialized := false
	for {
		line, err := reader.ReadBytes('\n')
		if len(bytes.TrimSpace(line)) > 0 {
			var req struct {
				JSONRPC string         `json:"jsonrpc"`
				ID      any            `json:"id"`
				Method  string         `json:"method"`
				Params  map[string]any `json:"params"`
			}
			if jsonErr := json.Unmarshal(line, &req); jsonErr != nil {
				resp := map[string]any{
					"jsonrpc": "2.0",
					"id":      nil,
					"error":   map[string]any{"code": -32700, "message": "Parse error"},
				}
				writeHelperJSON(resp)
				continue
			}
			switch req.Method {
			case "initialize":
				writeHelperJSON(map[string]any{
					"jsonrpc": "2.0",
					"id":      req.ID,
					"result": map[string]any{
						"protocolVersion": "2024-11-05",
						"capabilities":    map[string]any{"tools": map[string]any{}},
						"serverInfo":      map[string]any{"name": "helper-mcp", "version": "1.0.0"},
					},
				})
			case "notifications/initialized":
				initialized = true
			case "tools/call":
				if !initialized {
					fmt.Fprintln(os.Stderr, "protocolo violado: notifications/initialized nao recebido")
					writeHelperJSON(map[string]any{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"error":   map[string]any{"code": -32600, "message": "not initialized"},
					})
					continue
				}
				name, _ := req.Params["name"].(string)
				args, _ := req.Params["arguments"].(map[string]any)
				if name == "fail_tool" {
					fmt.Fprintln(os.Stderr, "detalhe stderr da ferramenta fail_tool")
					writeHelperJSON(map[string]any{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"result": map[string]any{
							"isError": true,
							"content": []map[string]any{{"type": "text", "text": "falha interna na ferramenta"}},
						},
					})
					continue
				}
				if name == "big_tool" {
					big := strings.Repeat("Z", 75*1024)
					writeHelperJSON(map[string]any{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"result": map[string]any{
							"content": []map[string]any{{"type": "text", "text": big}},
						},
					})
					continue
				}
				msg, _ := args["texto"].(string)
				outText := fmt.Sprintf("mcp_ok:%s:%s", name, msg)
				writeHelperJSON(map[string]any{
					"jsonrpc": "2.0",
					"id":      req.ID,
					"result": map[string]any{
						"content": []map[string]any{{"type": "text", "text": outText}},
					},
				})
			}
		}
		if err != nil {
			break
		}
	}
	os.Exit(0)
}

func writeHelperJSON(v any) {
	b, _ := json.Marshal(v)
	os.Stdout.Write(append(b, '\n'))
}

func TestMCPClientSTDIOAndRunnerIntegration(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	baseArgs := []string{"-test.run=TestHelperMCPServerProcess", "--"}
	baseEnv := []string{"GO_WANT_HELPER_MCP_SERVER=1"}

	t.Run("HandshakeAndCallToolSuccess", func(t *testing.T) {
		client := &MCPClient{
			Command: exe,
			Args:    baseArgs,
			Dir:     t.TempDir(),
			Env:     append(sterileEnv(os.Getenv), baseEnv...),
			Timeout: 5 * time.Second,
		}
		out, err := client.CallTool(context.Background(), "echo_mcp", map[string]any{"texto": "ola-mundo"})
		if err != nil {
			t.Fatalf("CallTool falhou: %v", err)
		}
		if out != "mcp_ok:echo_mcp:ola-mundo" {
			t.Fatalf("saida inesperada: %q", out)
		}
	})

	t.Run("PayloadMaiorQue64KiB", func(t *testing.T) {
		client := &MCPClient{
			Command: exe,
			Args:    baseArgs,
			Dir:     t.TempDir(),
			Env:     append(sterileEnv(os.Getenv), baseEnv...),
			Timeout: 5 * time.Second,
		}
		out, err := client.CallTool(context.Background(), "big_tool", nil)
		if err != nil {
			t.Fatalf("CallTool >64KiB falhou: %v", err)
		}
		if len(out) != 75*1024 {
			t.Fatalf("esperava 76800 bytes, obteve %d", len(out))
		}
	})

	t.Run("CapturaStderrEmErroDaFerramentaECrash", func(t *testing.T) {
		client := &MCPClient{
			Command: exe,
			Args:    baseArgs,
			Dir:     t.TempDir(),
			Env:     append(sterileEnv(os.Getenv), baseEnv...),
			Timeout: 5 * time.Second,
		}
		_, err := client.CallTool(context.Background(), "fail_tool", nil)
		if err == nil || !strings.Contains(err.Error(), "detalhe stderr da ferramenta fail_tool") {
			t.Fatalf("esperava erro contendo stderr da ferramenta, obteve: %v", err)
		}

		crashClient := &MCPClient{
			Command: exe,
			Args:    baseArgs,
			Dir:     t.TempDir(),
			Env:     append(sterileEnv(os.Getenv), "GO_WANT_HELPER_MCP_SERVER=1", "MCP_HELPER_MODE=crash_stderr"),
			Timeout: 5 * time.Second,
		}
		_, err = crashClient.CallTool(context.Background(), "echo_mcp", nil)
		if err == nil || !strings.Contains(err.Error(), "erro fatal simulado no stderr") {
			t.Fatalf("esperava erro contendo stderr do crash, obteve: %v", err)
		}
	})

	t.Run("TimeoutCancelaSubprocesso", func(t *testing.T) {
		client := &MCPClient{
			Command: exe,
			Args:    baseArgs,
			Dir:     t.TempDir(),
			Env:     append(sterileEnv(os.Getenv), "GO_WANT_HELPER_MCP_SERVER=1", "MCP_HELPER_MODE=sleep"),
			Timeout: 150 * time.Millisecond,
		}
		start := time.Now()
		_, err := client.CallTool(context.Background(), "echo_mcp", nil)
		if err == nil {
			t.Fatal("esperava erro de timeout, obteve nil")
		}
		if time.Since(start) > 3*time.Second {
			t.Fatalf("timeout demorou demais (%v)", time.Since(start))
		}
	})

	t.Run("RunnerExecutaNoTypeMCPReal", func(t *testing.T) {
		workdir := t.TempDir()
		cmdJSON, _ := json.Marshal(append([]string{exe}, baseArgs...))
		dagSrc := fmt.Sprintf(`{
			"name": "mcp-runner-test",
			"version": "2.0.0",
			"trigger": {"type": "manual"},
			"nodes": [{
				"id": "no_mcp_real",
				"type": "mcp",
				"ref": "echo_mcp",
				"command": %s,
				"inputs": {"texto": "handoff-real-mcp"},
				"routing": {"mode": "manual", "model": "none"},
				"output": {
					"artifact": "out/mcp_result.txt",
					"gate": "has_pattern('out/mcp_result.txt', 'mcp_ok:echo_mcp:handoff-real-mcp') && exit_code == 0"
				}
			}]
		}`, string(cmdJSON))

		d, err := dag.Parse([]byte(dagSrc))
		if err != nil {
			t.Fatalf("dag.Parse falhou: %v", err)
		}

		ox := &Orchestrix{
			ExtraEnv: baseEnv,
		}
		rep, err := runner.Run(context.Background(), d, runner.Options{
			Workdir:  workdir,
			Executor: ox,
		})
		if err != nil {
			t.Fatalf("runner.Run retornou erro: %v", err)
		}
		if !rep.OK || len(rep.Nodes) != 1 || rep.Nodes[0].Status != runner.StatusPassed {
			t.Fatalf("esperava no_mcp_real StatusPassed, obteve: %+v", rep)
		}
		got, err := os.ReadFile(filepath.Join(workdir, "out", "mcp_result.txt"))
		if err != nil {
			t.Fatalf("ler artefato gerado pelo no mcp: %v", err)
		}
		if string(got) != "mcp_ok:echo_mcp:handoff-real-mcp" {
			t.Fatalf("artefato divergente: %q", string(got))
		}
	})
}
