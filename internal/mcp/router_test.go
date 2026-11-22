package mcp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
)

// TestToolsViaRouter prova a AC-E09.1: as tools de dois Processos expõem pelo
// router único via STDIO (tools/list + tools/call), com execução real do
// comando declarado e zero servidor MCP por Processo.
func TestToolsViaRouter(t *testing.T) {
	base := t.TempDir()
	home := t.TempDir()
	vault := filepath.Join(home, ".aidd", "canonical")
	prepararProcessoComTools(t, base, "alpha", manifestoEco("eco_alpha"))
	prepararProcessoComTools(t, base, "beta", manifestoEco("eco_beta"))

	r := NewConfiguredRouter(home, vault, base)
	defer r.Close()
	if err := r.CarregarToolsDeProcessos(base); err != nil {
		t.Fatalf("CarregarToolsDeProcessos: %v", err)
	}

	reqs := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05"}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"eco_alpha","arguments":{"texto":"abc"}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"eco_beta","arguments":{"texto":"xyz"}}}`,
	}, "\n") + "\n"
	var out bytes.Buffer
	if err := r.Serve(strings.NewReader(reqs), &out); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	linhas := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(linhas) != 4 {
		t.Fatalf("esperava 4 respostas, obteve %d: %s", len(linhas), out.String())
	}
	verificarEcoNoRouter(t, linhas[1], linhas[2], linhas[3])

	// AC-E09.1: zero servidor MCP por Processo — nenhum externo sobe nem entra no registro.
	if ativos := r.proxy.ativos(); ativos != 0 {
		t.Errorf("servidores MCP externos ativos após tools de Processo: %d", ativos)
	}
	servers, err := LoadRegistry(RegistryPath(home))
	if err != nil {
		t.Fatalf("LoadRegistry: %v", err)
	}
	if len(servers) != 0 {
		t.Errorf("registro de MCPs ganhou servidores por Processo: %v", servers)
	}
}

func verificarEcoNoRouter(t *testing.T, lista, ecoA, ecoB string) {
	t.Helper()
	var env rpcEnvelope
	if err := json.Unmarshal([]byte(lista), &env); err != nil {
		t.Fatalf("unmarshal tools/list: %v", err)
	}
	presentes := map[string]bool{}
	for _, tl := range env.Result.Tools {
		presentes[tl.Name] = true
	}
	for _, nome := range []string{"eco_alpha", "eco_beta"} {
		if !presentes[nome] {
			t.Errorf("tool de Processo %q ausente em tools/list: %v", nome, presentes)
		}
	}
	for linha, esperado := range map[string]string{ecoA: "eco:abc", ecoB: "eco:xyz"} {
		var call rpcEnvelope
		if err := json.Unmarshal([]byte(linha), &call); err != nil {
			t.Fatalf("unmarshal tools/call: %v", err)
		}
		if call.Result.IsError || len(call.Result.Content) == 0 || call.Result.Content[0].Text != esperado {
			t.Errorf("tools/call devia devolver %q, obteve: %+v", esperado, call.Result)
		}
	}
}

func manifestoEco(nome string) string {
	man := map[string]any{"tools": []any{map[string]any{
		"name":        nome,
		"description": "Devolve eco: seguido do texto recebido.",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{"texto": map[string]any{"type": "string"}}},
		"comando":     []string{os.Args[0], "-test.run=TestHelperFerramentaEco"},
		"ambiente":    map[string]string{"LEGOLOOP_HELPER_FERRAMENTA": "1"},
	}}}
	data, err := json.Marshal(man)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func prepararProcessoComTools(t *testing.T, base, nome, manifesto string) {
	t.Helper()
	dir := filepath.Join(caminhos.DirProcessosEm(base), nome)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("montar Processo %s: %v", nome, err)
	}
	if err := os.WriteFile(caminhos.CaminhoToolsDoProcesso(dir), []byte(manifesto), 0o600); err != nil {
		t.Fatalf("gravar tools.json de %s: %v", nome, err)
	}
}

// TestConflitoToolExit prova a AC-E09.2: conflito de nome de tool entre
// Processos recusa o segundo registro e nomeia o conflito com os dois donos.
func TestConflitoToolExit(t *testing.T) {
	base := t.TempDir()
	prepararProcessoComTools(t, base, "alpha", manifestoEco("eco_dup"))
	prepararProcessoComTools(t, base, "beta", manifestoEco("eco_dup"))
	r := NewRouter(nil, nil)
	err := r.CarregarToolsDeProcessos(base)
	if err == nil {
		t.Fatal("segundo registro com conflito de nome devia ser recusado")
	}
	if !errors.Is(err, ErrConflitoTool) {
		t.Fatalf("erro esperado ErrConflitoTool, obteve: %v", err)
	}
	if CodigoDeSaida(err) == 0 {
		t.Error("exit code do conflito devia ser ≠ 0")
	}
	for _, trecho := range []string{"eco_dup", "alpha", "beta"} {
		if !strings.Contains(err.Error(), trecho) {
			t.Errorf("conflito devia nomear %q, obteve: %v", trecho, err)
		}
	}
	if reg, ok := r.toolsProcesso["eco_dup"]; !ok || reg.dono != "alpha" {
		t.Errorf("só o segundo registro devia ser recusado; dono atual: %+v", reg)
	}
}

// TestHelperFerramentaEco é o comando real das tools de teste: ecoa o texto dos
// argumentos recebidos via stdin e encerra sem o resumo do framework.
func TestHelperFerramentaEco(t *testing.T) {
	if os.Getenv("LEGOLOOP_HELPER_FERRAMENTA") != "1" {
		t.Skip("helper de execução de tool: só roda como comando do manifesto")
	}
	var args map[string]any
	_ = json.NewDecoder(os.Stdin).Decode(&args)
	texto, _ := args["texto"].(string)
	fmt.Print("eco:" + texto)
	os.Exit(0)
}

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
