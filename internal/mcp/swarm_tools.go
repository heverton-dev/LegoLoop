package mcp

import (
	"context"
	"encoding/json"
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
	sb, err := board.OpenSwarmBoard(strArg(args, "dir", ".legoloop-swarm"))
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

func dispatchSwarmTool(sb *board.SwarmBoard, name string, args map[string]any) (any, bool, error) {
	switch name {
	case "swarm_add":
		w, err := swarmAddWorker(sb, args)
		return w, true, err
	case "swarm_post":
		m, err := sb.PostMessage(strArg(args, "from", "coordenador"), strArg(args, "to", "@all"), strArg(args, "body", ""))
		return m, true, err
	case "swarm_read":
		msgs, err := sb.ReadMessages(int64Arg(args, "since"), nil)
		return msgs, true, err
	case "swarm_wait":
		msgs, err := swarmWaitMessages(sb, args)
		return msgs, true, err
	case "swarm_status":
		st, err := sb.Status()
		return st, true, err
	case "swarm_kill":
		err := sb.KillWorker(strArg(args, "name", ""), strArg(args, "repo", "."))
		return map[string]string{"status": "killed"}, true, err
	case "swarm_monitor":
		stalled, err := sb.CheckStall(60 * time.Second)
		return map[string]any{"stalled": stalled}, true, err
	}
	return nil, false, nil
}

func swarmAddWorker(sb *board.SwarmBoard, args map[string]any) (any, error) {
	return sb.AddWorker(board.AddWorkerRequest{
		Name:          strArg(args, "name", ""),
		Role:          strArg(args, "role", ""),
		Model:         strArg(args, "model", ""),
		WorkspaceMode: board.WorkspaceMode(strArg(args, "workspace", "inherit")),
		RepoDir:       strArg(args, "repo", "."),
	})
}

func swarmWaitMessages(sb *board.SwarmBoard, args map[string]any) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return sb.WaitMessages(ctx, strArg(args, "target", "@all"), int64Arg(args, "since"), 100*time.Millisecond)
}
