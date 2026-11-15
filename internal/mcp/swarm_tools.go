package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"legoloop/internal/board"
)

func strArg(args map[string]any, key, def string) string {
	if v, ok := args[key].(string); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return def
}

func int64Arg(args map[string]any, key string) int64 {
	if v, ok := args[key].(float64); ok {
		return int64(v)
	}
	return 0
}

func (r *Router) jsonResult(id any, v any) JSONRPCResponse {
	b, err := json.Marshal(v)
	if err != nil {
		return r.toolCallError(id, err.Error())
	}
	return r.toolCallText(id, string(b))
}

func (r *Router) handleSwarmTool(id any, name string, args map[string]any) (JSONRPCResponse, bool) {
	if !strings.HasPrefix(name, "swarm_") {
		return JSONRPCResponse{}, false
	}
	dir := strArg(args, "dir", ".legoloop-swarm")
	if name == "swarm_status" {
		st, err := readSwarmStatusReadOnly(dir)
		if err != nil {
			return r.toolCallError(id, err.Error()), true
		}
		return r.jsonResult(id, st), true
	}
	sb, err := board.OpenSwarmBoard(dir)
	if err != nil {
		return r.toolCallError(id, err.Error()), true
	}
	res, handled, err := dispatchSwarmTool(sb, name, args)
	if !handled {
		return JSONRPCResponse{}, false
	}
	if err != nil {
		return r.toolCallError(id, err.Error()), true
	}
	return r.jsonResult(id, res), true
}

func readSwarmStatusReadOnly(dir string) (board.SwarmStatus, error) {
	clean := filepath.Clean(strings.TrimSpace(dir))
	if clean == "" {
		clean = ".legoloop-swarm"
	}
	proj := clean
	if base := filepath.Base(clean); base == ".legoloop-swarm" || base == ".swarm" {
		proj = filepath.Dir(clean)
	}
	boardDir := filepath.Join(proj, ".legoloop-swarm")
	candidates := []string{
		filepath.Join(boardDir, "swarm.json"),
		filepath.Join(boardDir, "swarm_board.json"),
		filepath.Join(proj, ".swarm", "swarm.json"),
		filepath.Join(proj, ".swarm", "swarm_board.json"),
		filepath.Join(proj, "swarm_board.json"),
	}
	return loadFirstSwarmStatus(candidates)
}

func loadFirstSwarmStatus(candidates []string) (board.SwarmStatus, error) {
	for _, p := range candidates {
		info, err := os.Stat(p)
		if err != nil || info.IsDir() {
			continue
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return board.SwarmStatus{}, err
		}
		var raw struct {
			Workers  []board.Worker  `json:"workers"`
			Messages []board.Message `json:"messages"`
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			return board.SwarmStatus{}, err
		}
		return board.SwarmStatus{Workers: raw.Workers, MessageCount: len(raw.Messages)}, nil
	}
	return board.SwarmStatus{Workers: []board.Worker{}, MessageCount: 0}, nil
}

func dispatchSwarmTool(sb *board.SwarmBoard, name string, args map[string]any) (any, bool, error) {
	switch name {
	case "swarm_add":
		w, err := swarmAddWorker(sb, args)
		return w, true, err
	case "swarm_post":
		m, err := sb.PostMessage(strArg(args, "from", "coordenador"), strArg(args, "to", "@all"), strArg(args, "body", ""))
		return m, true, err
	case "swarm_read":
		msgs, err := sb.ReadMessages(int64Arg(args, "since"), parseSeqsArg(args, "seqs"))
		return msgs, true, err
	case "swarm_wait":
		msgs, err := swarmWaitMessages(sb, args)
		return msgs, true, err
	case "swarm_status":
		st, err := sb.Status()
		return st, true, err
	case "swarm_kill":
		wName := strArg(args, "name", "")
		err := sb.KillWorker(wName, strArg(args, "repo", "."))
		return map[string]string{"status": "killed", "name": wName}, true, err
	case "swarm_monitor":
		stalled, err := sb.CheckStall(durationArg(args, "timeout", 60*time.Second))
		return map[string]any{"stalled": stalled}, true, err
	}
	return nil, false, nil
}

func swarmAddWorker(sb *board.SwarmBoard, args map[string]any) (any, error) {
	return sb.AddWorker(board.AddWorkerRequest{
		Name:          strArg(args, "name", ""),
		Role:          strArg(args, "role", ""),
		Parent:        strArg(args, "parent", ""),
		Model:         strArg(args, "model", ""),
		WorkspaceMode: board.WorkspaceMode(strArg(args, "workspace", "inherit")),
		RepoDir:       strArg(args, "repo", "."),
	})
}

func swarmWaitMessages(sb *board.SwarmBoard, args map[string]any) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), durationArg(args, "timeout", 30*time.Second))
	defer cancel()
	return sb.WaitMessages(ctx, strArg(args, "target", "@all"), int64Arg(args, "since"), 100*time.Millisecond)
}

func parseSeqsArg(args map[string]any, key string) []int64 {
	raw := strArg(args, key, "")
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out []int64
	for _, p := range strings.Split(raw, ",") {
		if n, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			out = append(out, n)
		}
	}
	return out
}

func durationArg(args map[string]any, key string, def time.Duration) time.Duration {
	if s := strArg(args, key, ""); s != "" {
		if d, err := time.ParseDuration(s); err == nil && d > 0 {
			return d
		}
	}
	if n := int64Arg(args, key); n > 0 {
		return time.Duration(n) * time.Millisecond
	}
	return def
}

func swarmToolDefs() []map[string]any {
	return []map[string]any{
		swarmDef("swarm_add", "Registra um worker no quadro Swarm.", []string{"role"}),
		swarmDef("swarm_post", "Publica uma mensagem no quadro Swarm.", []string{"body"}),
		swarmDef("swarm_read", "Le mensagens do quadro Swarm.", nil),
		swarmDef("swarm_wait", "Aguarda mensagens destinadas a um worker no quadro Swarm.", nil),
		swarmDef("swarm_status", "Consulta o estado do quadro Swarm em modo somente leitura.", nil),
		swarmDef("swarm_kill", "Encerra um worker do quadro Swarm.", []string{"name"}),
		swarmDef("swarm_monitor", "Detecta workers estagnados no quadro Swarm.", nil),
	}
}

func swarmDef(name, desc string, req []string) map[string]any {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"dir": map[string]any{"type": "string", "description": "Diretorio do projeto ou quadro Swarm"},
		},
	}
	if len(req) > 0 {
		schema["required"] = req
	}
	return map[string]any{"name": name, "description": desc, "inputSchema": schema}
}

type localVault struct {
	home       string
	vault      string
	projectDir string
	skills     *SkillEngine
}

func (v *localVault) Search(query string) []ToolMatch {
	q := strings.ToLower(strings.TrimSpace(query))
	var out []ToolMatch
	for _, td := range swarmToolDefs() {
		name, _ := td["name"].(string)
		desc, _ := td["description"].(string)
		if matchQuery(q, "swarm", name, desc) {
			out = append(out, ToolMatch{Server: "swarm", Name: name, Description: desc})
		}
	}
	if v.skills != nil {
		for _, s := range v.skills.List() {
			if matchQuery(q, "skills", s.ID, s.Resumo+" "+s.Nome) {
				out = append(out, ToolMatch{Server: "skills", Name: s.ID, Description: s.Resumo})
			}
		}
	}
	return out
}

func matchQuery(q, srv, name, desc string) bool {
	if q == "" {
		return true
	}
	hay := strings.ToLower(srv + " " + name + " " + desc)
	return strings.Contains(hay, q)
}

type localRunner struct {
	router *Router
}

func (lr *localRunner) Run(server, tool string, args map[string]any) (string, error) {
	if args == nil {
		args = map[string]any{}
	}
	if server == "swarm" || strings.HasPrefix(tool, "swarm_") {
		resp, ok := lr.router.handleSwarmTool(1, tool, args)
		if !ok {
			return "", fmt.Errorf("ferramenta swarm desconhecida: %s", tool)
		}
		return extractResponseText(resp)
	}
	if server == "skills" {
		if _, ok := args["skill_id"]; !ok && tool != "skill_search" && tool != "call_skill" {
			args["skill_id"] = tool
		}
		resp := lr.router.callSkill(1, args)
		return extractResponseText(resp)
	}
	return "", errors.New("servidor MCP externo nao registrado para execucao direta")
}

func extractResponseText(resp JSONRPCResponse) (string, error) {
	m, ok := resp.Result.(map[string]any)
	if !ok {
		return "", errors.New("resposta MCP em formato inesperado")
	}
	isErr, _ := m["isError"].(bool)
	items, _ := m["content"].([]map[string]any)
	var text string
	if len(items) > 0 {
		text, _ = items[0]["text"].(string)
	}
	if isErr {
		return "", errors.New(text)
	}
	return text, nil
}
