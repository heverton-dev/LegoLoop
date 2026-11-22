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

	"github.com/heverton-dev/LegoLoop/internal/board"
	"github.com/heverton-dev/LegoLoop/internal/dag"
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
		return board.SwarmStatus{Workers: board.LiveView(raw.Workers), MessageCount: len(raw.Messages)}, nil
	}
	return board.SwarmStatus{Workers: []board.Worker{}, MessageCount: 0}, nil
}

func dispatchSwarmTool(sb *board.SwarmBoard, name string, args map[string]any) (any, bool, error) {
	switch name {
	case "swarm_add":
		w, err := swarmAddWorker(sb, args)
		return w, true, err
	case "swarm_spawn":
		res, err := swarmSpawnWorker(sb, args)
		return res, true, err
	case "swarm_run":
		res, err := swarmRunDAG(sb, args)
		return res, true, err
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

// swarmSpawnWorker inicia o processo real do worker; "command" aceita lista de strings.
func swarmSpawnWorker(sb *board.SwarmBoard, args map[string]any) (any, error) {
	name := strArg(args, "name", "")
	cmd := stringListArg(args, "command")
	if name == "" || len(cmd) == 0 {
		return nil, errors.New("swarm_spawn exige name e command (lista)")
	}
	pid, err := sb.SpawnWorker(context.Background(), name, cmd)
	if err != nil {
		return nil, err
	}
	return map[string]any{"name": name, "pid": pid, "status": board.WorkerWorking}, nil
}

// swarmRunDAG executa o DAG por nivel topologico (workers paralelos por nivel, worktree proprio) e devolve o
// desfecho de cada nivel. Bloqueia ate o fim ou ate o tempo maximo por nivel.
func swarmRunDAG(sb *board.SwarmBoard, args map[string]any) (any, error) {
	path := strArg(args, "dag", "")
	if path == "" {
		return nil, errors.New("swarm_run exige dag (caminho do arquivo DAG)")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p, err := dag.Parse(data)
	if err != nil {
		return nil, err
	}
	levels, skipped, err := board.LevelsFromDAG(p, stringListArg(args, "command"))
	if err != nil {
		return nil, err
	}
	keep, _ := args["keep"].(bool)
	gates, _ := args["gates"].(bool)
	results, runErr := sb.RunLevels(context.Background(), levels, board.RunOptions{
		Mode:           board.WorkspaceMode(strArg(args, "workspace", "share")),
		RepoDir:        strArg(args, "repo", "."),
		Keep:           keep,
		CheckGates:     gates,
		DefaultTimeout: durationArg(args, "timeout", 10*time.Minute),
	})
	out := map[string]any{"dag": p.Name, "levels": results, "skipped": skipped}
	if runErr != nil {
		partial, _ := json.Marshal(out)
		return nil, fmt.Errorf("%w; parcial: %s", runErr, partial)
	}
	return out, nil
}

func stringListArg(args map[string]any, key string) []string {
	raw, _ := args[key].([]any)
	var out []string
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
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
		swarmSpawnDef(),
		swarmRunDef(),
		swarmDef("swarm_post", "Publica uma mensagem no quadro Swarm.", []string{"body"}),
		swarmDef("swarm_read", "Le mensagens do quadro Swarm.", nil),
		swarmDef("swarm_wait", "Aguarda mensagens destinadas a um worker no quadro Swarm.", nil),
		swarmDef("swarm_status", "Consulta o estado do quadro Swarm em modo somente leitura.", nil),
		swarmDef("swarm_kill", "Encerra um worker do quadro Swarm.", []string{"name"}),
		swarmDef("swarm_monitor", "Detecta workers estagnados no quadro Swarm.", nil),
	}
}

func swarmRunDef() map[string]any {
	d := swarmDef("swarm_run", "Executa o DAG por nivel topologico: workers paralelos por nivel, worktree proprio, ordem garantida.", []string{"dag"})
	props := d["inputSchema"].(map[string]any)["properties"].(map[string]any)
	props["dag"] = map[string]any{"type": "string", "description": "Caminho do arquivo DAG"}
	props["command"] = map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Comando global dos workers (o command do no vale antes)"}
	props["repo"] = map[string]any{"type": "string", "description": "Repositorio git dos worktrees"}
	props["workspace"] = map[string]any{"type": "string", "description": "share, branch ou inherit"}
	props["timeout"] = map[string]any{"type": "string", "description": "Tempo maximo por nivel, ex.: 10m"}
	props["keep"] = map[string]any{"type": "boolean", "description": "Mantem workers e worktrees apos o sucesso"}
	props["gates"] = map[string]any{"type": "boolean", "description": "Avalia o gate de cada no no worktree do worker; gate falso reprova o nivel"}
	return d
}

func swarmSpawnDef() map[string]any {
	d := swarmDef("swarm_spawn", "Inicia o processo real de um worker registrado, no seu worktree.", []string{"name", "command"})
	props := d["inputSchema"].(map[string]any)["properties"].(map[string]any)
	props["name"] = map[string]any{"type": "string", "description": "Nome do worker ja registrado"}
	props["command"] = map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Comando e argumentos"}
	return d
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
	specs      map[string]ServerSpec
	catalog    map[string][]CatalogEntry
}

// searchExternal procura nos MCPs externos registrados, sem subir nenhum processo.
func (v *localVault) searchExternal(q string) []ToolMatch {
	var out []ToolMatch
	for srv := range v.specs {
		tools := v.catalog[srv]
		for _, t := range tools {
			if matchQuery(q, srv, t.Name, t.Description) {
				out = append(out, ToolMatch{Server: srv, Name: t.Name, Description: t.Description})
			}
		}
		if len(tools) == 0 && q != "" && matchQuery(q, srv, "", "") {
			out = append(out, ToolMatch{Server: srv, Name: "*", Description: "servidor registrado, ferramentas nao indexadas: rode 'legoloop mcp-registry refresh'"})
		}
	}
	return out
}

func (v *localVault) Search(query string) []ToolMatch {
	q := strings.ToLower(strings.TrimSpace(query))
	out := v.searchExternal(q)
	inst := instalarToolDef()
	instName, _ := inst["name"].(string)
	instDesc, _ := inst["description"].(string)
	if matchQuery(q, servidorLegoloop, instName, instDesc) {
		out = append(out, ToolMatch{Server: servidorLegoloop, Name: instName, Description: instDesc})
	}
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
	ordenarMatches(q, out)
	return out
}

func matchQuery(q, srv, name, desc string) bool {
	return pontuarBusca(q, srv+" "+name+" "+desc) >= minimoDePontos(q)
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
	if lr.router.proxy == nil {
		return "", errors.New("servidor MCP externo nao registrado para execucao direta")
	}
	return lr.router.proxy.CallTool(server, tool, args)
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
