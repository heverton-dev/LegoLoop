package mcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrRequisicaoInvalida  = errors.New("mcp: requisicao json-rpc invalida")
	ErrMetodoNaoEncontrado = errors.New("mcp: metodo nao encontrado")
)

// ToolDef define os metadados de uma ferramenta no cofre.
type ToolDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema,omitempty"`
}

// ToolMatch representa o resultado de uma busca no cofre.
type ToolMatch struct {
	Server      string `json:"server"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// VaultRegistry define o contrato para consultar MCPs e ferramentas.
type VaultRegistry interface {
	Search(query string) []ToolMatch
}

// ProcessRunner define o contrato para disparar a execucao de ferramenta no MCP real.
type ProcessRunner interface {
	Run(server, tool string, args map[string]any) (string, error)
}

// JSONRPCRequest modelo padrao da requisicao JSON-RPC 2.0.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse modelo padrao de resposta JSON-RPC 2.0.
type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// Router gerencia o ciclo de vida do Meta-MCP Router via STDIO.
type Router struct {
	vault     VaultRegistry
	runner    ProcessRunner
	skills    *SkillEngine
	skillRoot string
	swarm     bool
}

// NewRouter instancia um novo Router.
func NewRouter(vault VaultRegistry, runner ProcessRunner) *Router {
	return &Router{vault: vault, runner: runner}
}

// EnableSwarm habilita a exposicao de ferramentas swarm_* em tools/list.
func (r *Router) EnableSwarm() {
	r.swarm = true
}

// NewConfiguredRouter instancia um Router conectado ao cofre, skills reais e ferramentas Swarm.
func NewConfiguredRouter(home, vault, projectDir string) *Router {
	engine, skillRoot := loadRealSkillEngine(vault, projectDir)
	lv := &localVault{home: home, vault: vault, projectDir: projectDir, skills: engine}
	r := &Router{vault: lv, skills: engine, skillRoot: skillRoot, swarm: true}
	r.runner = &localRunner{router: r}
	return r
}

func loadRealSkillEngine(vault, projectDir string) (*SkillEngine, string) {
	engine := NewSkillEngine()
	vaultSkills := filepath.Join(vault, "skills")
	if info, err := os.Stat(vaultSkills); err == nil && info.IsDir() {
		_ = engine.IndexDirectory(vaultSkills)
		if len(engine.List()) > 0 {
			return engine, vault
		}
	}
	projSkills := filepath.Join(projectDir, "skills")
	if abs, err := filepath.Abs(projSkills); err == nil {
		if info, stErr := os.Stat(abs); stErr == nil && info.IsDir() {
			_ = engine.IndexDirectory(abs)
			return engine, abs
		}
	}
	if vault != "" {
		return engine, vault
	}
	return engine, projectDir
}

// Serve processa as mensagens de entrada linha a linha via STDIO sem limite de 64KiB.
func (r *Router) Serve(in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadBytes('\n')
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) > 0 {
			if wErr := r.processLine(trimmed, out); wErr != nil {
				return wErr
			}
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (r *Router) processLine(line []byte, out io.Writer) error {
	var req JSONRPCRequest
	if err := json.Unmarshal(line, &req); err != nil || req.Method == "" {
		return writeJSONRPC(out, JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      nil,
			Error:   map[string]any{"code": -32700, "message": "Parse error"},
		})
	}
	if req.ID == nil && strings.HasPrefix(req.Method, "notifications/") {
		return nil
	}
	return writeJSONRPC(out, r.handleRequest(req))
}

func writeJSONRPC(out io.Writer, resp JSONRPCResponse) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = out.Write(append(data, '\n'))
	return err
}

// handleRequest despacha o metodo JSON-RPC apropriado.
func (r *Router) handleRequest(req JSONRPCRequest) JSONRPCResponse {
	switch req.Method {
	case "initialize":
		return r.handleInitialize(req.ID)
	case "notifications/initialized":
		return JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{}}
	case "tools/list":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  map[string]any{"tools": r.minimalTools()},
		}
	case "tools/call":
		return r.handleToolCall(req)
	default:
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: map[string]any{
				"code":    -32601,
				"message": "Method not found",
			},
		}
	}
}

func (r *Router) handleInitialize(id any) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "legoloop-mcp-router", "version": "1.0.0"},
		},
	}
}

type toolCallParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

func (r *Router) handleToolCall(req JSONRPCRequest) JSONRPCResponse {
	var params toolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return r.toolCallError(req.ID, fmt.Sprintf("invalid params: %v", err))
	}
	switch params.Name {
	case "mcp_search":
		return r.callMcpSearch(req.ID, params.Arguments)
	case "call_mcp_tool":
		return r.callMcpTool(req.ID, params.Arguments)
	case "skill_search":
		return r.callSkillSearch(req.ID, params.Arguments)
	case "call_skill":
		return r.callSkill(req.ID, params.Arguments)
	default:
		if resp, ok := r.handleSwarmTool(req.ID, params.Name, params.Arguments); ok {
			return resp
		}
		return r.toolCallError(req.ID, fmt.Sprintf("unknown tool: %s", params.Name))
	}
}

func (r *Router) callMcpSearch(id any, args map[string]any) JSONRPCResponse {
	if r.vault == nil {
		return r.toolCallText(id, "Cofre de ferramentas vazio.")
	}
	query, _ := args["query"].(string)
	matches := r.vault.Search(query)
	if len(matches) == 0 {
		return r.toolCallText(id, fmt.Sprintf("Nenhuma ferramenta encontrada para '%s'.", query))
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Encontradas %d ferramentas no cofre:\n", len(matches)))
	for _, m := range matches {
		b.WriteString(fmt.Sprintf("- Servidor [%s] Ferramenta '%s': %s\n", m.Server, m.Name, m.Description))
	}
	return r.toolCallText(id, b.String())
}

func (r *Router) callMcpTool(id any, args map[string]any) JSONRPCResponse {
	if r.runner == nil {
		return r.toolCallError(id, "Nenhum executor de processo configurado.")
	}
	srv, _ := args["server"].(string)
	tool, _ := args["tool"].(string)
	toolArgs, _ := args["arguments"].(map[string]any)
	out, err := r.runner.Run(srv, tool, toolArgs)
	if err != nil {
		return r.toolCallError(id, fmt.Sprintf("Falha ao executar %s:%s: %v", srv, tool, err))
	}
	return r.toolCallText(id, out)
}

func (r *Router) toolCallText(id any, text string) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]any{
			"content": []map[string]any{{"type": "text", "text": text}},
		},
	}
}

func (r *Router) toolCallError(id any, msg string) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]any{
			"isError": true,
			"content": []map[string]any{{"type": "text", "text": msg}},
		},
	}
}

func (r *Router) minimalTools() []map[string]any {
	tools := r.baseTools()
	if r.skills != nil {
		tools = append(tools, skillToolDefs()...)
	}
	if r.swarm || r.vault != nil || r.runner != nil {
		tools = append(tools, swarmToolDefs()...)
	}
	return tools
}

func (r *Router) baseTools() []map[string]any {
	return []map[string]any{
		{
			"name":        "mcp_search",
			"description": "Busca ferramentas e schemas no cofre de MCPs sob demanda para evitar inchaço de contexto.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{"type": "string", "description": "Termo de busca funcional (ex: decompile, screenshot, postgres)"},
				},
				"required": []string{"query"},
			},
		},
		{
			"name":        "call_mcp_tool",
			"description": "Executa uma ferramenta de um servidor MCP real sob demanda.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"server":    map[string]any{"type": "string", "description": "Nome do servidor MCP no cofre (ex: rea, playwright)"},
					"tool":      map[string]any{"type": "string", "description": "Nome da ferramenta real"},
					"arguments": map[string]any{"type": "object", "description": "Parametros para a ferramenta"},
				},
				"required": []string{"server", "tool"},
			},
		},
	}
}
