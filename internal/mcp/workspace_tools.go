package mcp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Servidores embutidos do router, acessíveis por call_mcp_tool sem subir processo externo:
//   - ws:       bash, read_file, write_file, edit_file (o trabalho de shell e arquivos);
//   - legoloop: swarm_* (o quadro Swarm).
const (
	servidorWs       = "ws"
	servidorLegoloop = "legoloop"
	limiteSaidaWs    = 30000
	timeoutBashWs    = 120 * time.Second
)

// descricaoGateway anuncia os servidores embutidos em poucos bytes: é a única descrição que o modelo recebe.
const descricaoGateway = "ws: bash, read_file, write_file, edit_file. legoloop: skill_search, call_skill, mcp_search, instalar, swarm_*. Erro mostra o uso."

// gatewayToolDefs é a única ferramenta exposta em modo gateway (piso de contexto).
func gatewayToolDefs() []map[string]any {
	return []map[string]any{{
		"name":        "call_mcp_tool",
		"description": descricaoGateway,
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"server":    map[string]any{"type": "string"},
				"tool":      map[string]any{"type": "string"},
				"arguments": map[string]any{"type": "object"},
			},
		},
	}}
}

// runEmbutido despacha os servidores embutidos; ok=false deixa o executor externo agir.
func (r *Router) runEmbutido(id any, srv, tool string, args map[string]any) (JSONRPCResponse, bool) {
	switch srv {
	case servidorLegoloop:
		return r.runLegoloop(id, tool, args)
	case servidorWs:
		out, err := runWs(tool, args)
		if err != nil {
			return r.toolCallError(id, err.Error()), true
		}
		return r.toolCallText(id, out), true
	}
	return JSONRPCResponse{}, false
}

// runLegoloop atende o servidor embutido legoloop: buscas, tools de Processo e swarm_*.
func (r *Router) runLegoloop(id any, tool string, args map[string]any) (JSONRPCResponse, bool) {
	switch tool {
	case "skill_search":
		return r.callSkillSearch(id, args), true
	case "call_skill":
		return r.callSkill(id, args), true
	case "mcp_search":
		return r.callMcpSearch(id, args), true
	case "instalar":
		return r.callInstalar(id, args), true
	}
	if _, ok := r.toolsProcesso[tool]; ok {
		return r.chamarToolProcesso(id, tool, args), true
	}
	if resp, ok := r.handleSessaoTool(id, tool, args); ok {
		return resp, true
	}
	return r.handleSwarmTool(id, tool, args)
}

func runWs(tool string, args map[string]any) (string, error) {
	switch tool {
	case "bash":
		return wsBash(args)
	case "read_file":
		return wsRead(args)
	case "write_file":
		return wsWrite(args)
	case "edit_file":
		return wsEdit(args)
	}
	return "", fmt.Errorf("ws: ferramenta desconhecida %q (use bash, read_file, write_file, edit_file)", tool)
}

func shellDisponivel() string {
	if p, err := exec.LookPath("bash"); err == nil {
		return p
	}
	return "sh"
}

func cortarSaida(s string) string {
	if len(s) <= limiteSaidaWs {
		return s
	}
	return s[:limiteSaidaWs] + fmt.Sprintf("\n[saída cortada: %d bytes]", len(s)-limiteSaidaWs)
}

func wsBash(args map[string]any) (string, error) {
	cmdTxt := strArg(args, "command", "")
	if cmdTxt == "" {
		return "", errors.New("ws.bash: uso arguments={command[,cwd]}")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutBashWs)
	defer cancel()
	cmd := exec.CommandContext(ctx, shellDisponivel(), "-c", cmdTxt)
	cmd.Dir = strArg(args, "cwd", "")
	out, err := cmd.CombinedOutput()
	saida := cortarSaida(string(out))
	var saiu *exec.ExitError
	if errors.As(err, &saiu) {
		return "", fmt.Errorf("exit %d: %s", saiu.ExitCode(), saida)
	}
	if err != nil {
		return "", fmt.Errorf("ws.bash: %v", err)
	}
	return saida, nil
}

func wsRead(args map[string]any) (string, error) {
	b, err := os.ReadFile(strArg(args, "path", ""))
	if err != nil {
		return "", fmt.Errorf("ws.read_file: uso arguments={path}: %v", err)
	}
	return cortarSaida(string(b)), nil
}

func wsWrite(args map[string]any) (string, error) {
	path := strArg(args, "path", "")
	if path == "" {
		return "", errors.New("ws.write_file: uso arguments={path,content}")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("ws.write_file: %v", err)
	}
	conteudo := strArg(args, "content", "")
	if err := os.WriteFile(path, []byte(conteudo), 0o644); err != nil {
		return "", fmt.Errorf("ws.write_file: %v", err)
	}
	return fmt.Sprintf("gravado %d bytes em %s", len(conteudo), path), nil
}

func wsEdit(args map[string]any) (string, error) {
	path, antigo, novo := strArg(args, "path", ""), strArg(args, "old", ""), strArg(args, "new", "")
	if antigo == "" {
		return "", errors.New("ws.edit_file: uso arguments={path,old,new[,replace_all]}")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("ws.edit_file: %v", err)
	}
	n := strings.Count(string(b), antigo)
	if n == 0 {
		return "", fmt.Errorf("ws.edit_file: trecho não encontrado em %s", path)
	}
	todos, _ := args["replace_all"].(bool)
	if n > 1 && !todos {
		return "", fmt.Errorf("ws.edit_file: %d ocorrências em %s; detalhe o trecho ou use replace_all", n, path)
	}
	if err := os.WriteFile(path, []byte(strings.ReplaceAll(string(b), antigo, novo)), 0o644); err != nil {
		return "", fmt.Errorf("ws.edit_file: %v", err)
	}
	return fmt.Sprintf("editado %s (%d substituição)", path, n), nil
}
