package mcp

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRouterHandshakeAndListTools(t *testing.T) {
	// Red-Green Cycle: Simula um cliente MCP enviando initialize e tools/list via stdio
	initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05"}}` + "\n"
	listReq := `{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n"
	in := bytes.NewBufferString(initReq + listReq)
	out := &bytes.Buffer{}

	router := NewRouter(nil, nil)
	err := router.Serve(in, out)
	if err != nil {
		t.Fatalf("Serve falhou com erro: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) < 2 {
		t.Fatalf("esperava ao menos 2 respostas JSON-RPC, recebeu %d: %s", len(lines), out.String())
	}

	// 1. Verifica resposta do initialize
	var initResp JSONRPCResponse
	if err := json.Unmarshal([]byte(lines[0]), &initResp); err != nil {
		t.Fatalf("erro ao decodificar initialize: %v", err)
	}
	if initResp.ID != float64(1) {
		t.Errorf("ID esperado 1, obteve %v", initResp.ID)
	}

	// 2. Verifica resposta do tools/list: deve conter ESTREITAMENTE 2 ferramentas
	var listResp struct {
		JSONRPC string `json:"jsonrpc"`
		ID      any    `json:"id"`
		Result  struct {
			Tools []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"tools"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(lines[1]), &listResp); err != nil {
		t.Fatalf("erro ao decodificar tools/list: %v", err)
	}

	tools := listResp.Result.Tools
	if len(tools) != 2 {
		t.Fatalf("esperava exatamente 2 ferramentas no handshake inicial, recebeu %d", len(tools))
	}
	names := []string{tools[0].Name, tools[1].Name}
	if names[0] != "mcp_search" || names[1] != "call_mcp_tool" {
		t.Errorf("ferramentas esperadas ['mcp_search', 'call_mcp_tool'], obteve %v", names)
	}
}

func TestRouterSearchAndCallTool(t *testing.T) {
	// Cria cofre simulado em memória
	vault := &fakeVault{
		servers: map[string][]ToolDef{
			"rea": {
				{Name: "decompile_jeb_item", Description: "Decompiles a Java item using JEB"},
				{Name: "inspect_native_instruction", Description: "Inspects assembly instructions"},
			},
			"playwright": {
				{Name: "browser_click", Description: "Click an element on page"},
			},
		},
	}
	runner := &fakeRunner{
		output: "Resultado da descompilacao com sucesso",
	}

	searchReq := `{"jsonrpc":"2.0","id":10,"method":"tools/call","params":{"name":"mcp_search","arguments":{"query":"decompile"}}}` + "\n"
	callReq := `{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"call_mcp_tool","arguments":{"server":"rea","tool":"decompile_jeb_item","arguments":{"item":"Main"}}}}` + "\n"
	in := bytes.NewBufferString(searchReq + callReq)
	out := &bytes.Buffer{}

	router := NewRouter(vault, runner)
	if err := router.Serve(in, out); err != nil {
		t.Fatalf("falha ao servir requisicoes: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) < 2 {
		t.Fatalf("esperava 2 respostas, recebeu %d: %s", len(lines), out.String())
	}

	// 1. Validar resposta de busca: deve encontrar apenas decompile_jeb_item
	var searchResp struct {
		ID     any `json:"id"`
		Result struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(lines[0]), &searchResp); err != nil {
		t.Fatalf("erro ao decodificar searchResp: %v", err)
	}
	if len(searchResp.Result.Content) == 0 || !strings.Contains(searchResp.Result.Content[0].Text, "decompile_jeb_item") {
		t.Errorf("esperava encontrar 'decompile_jeb_item', obteve: %+v", searchResp.Result)
	}
	if strings.Contains(searchResp.Result.Content[0].Text, "browser_click") {
		t.Errorf("ferramenta irrelevante 'browser_click' nao deveria aparecer na busca")
	}

	// 2. Validar execucao da ferramenta
	var callResp struct {
		ID     any `json:"id"`
		Result struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(lines[1]), &callResp); err != nil {
		t.Fatalf("erro ao decodificar callResp: %v", err)
	}
	if len(callResp.Result.Content) == 0 || !strings.Contains(callResp.Result.Content[0].Text, "Resultado da descompilacao") {
		t.Errorf("esperava output do runner, obteve: %+v", callResp.Result)
	}
}

type fakeVault struct {
	servers map[string][]ToolDef
}

func (f *fakeVault) Search(query string) []ToolMatch {
	var matches []ToolMatch
	q := strings.ToLower(query)
	for srv, tools := range f.servers {
		for _, tool := range tools {
			if strings.Contains(strings.ToLower(tool.Name), q) || strings.Contains(strings.ToLower(tool.Description), q) {
				matches = append(matches, ToolMatch{
					Server:      srv,
					Name:        tool.Name,
					Description: tool.Description,
				})
			}
		}
	}
	return matches
}

func (f *fakeVault) GetServerCommand(srv string) (string, []string, error) {
	return "cmd.exe", []string{"/c", "echo", "fake"}, nil
}

type fakeRunner struct {
	output string
}

func (f *fakeRunner) Run(server, tool string, args map[string]any) (string, error) {
	return f.output, nil
}
