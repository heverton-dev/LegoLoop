package board

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DiscoverBoards localiza diretórios de quadros Swarm ativos.
func DiscoverBoards(root, singleDir, procName string) map[string]string {
	boards := make(map[string]string)
	if singleDir != "" {
		boards[procName] = singleDir
		return boards
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return boards
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		bdir := filepath.Join(root, e.Name(), ".swarm")
		if _, err := os.Stat(filepath.Join(bdir, "swarm.json")); err == nil {
			boards[e.Name()] = bdir
		}
	}
	return boards
}

// BoardsPayload serializa o snapshot de todos os quadros descobertos.
func BoardsPayload(root, singleDir, procName string) ([]byte, error) {
	boards := DiscoverBoards(root, singleDir, procName)
	out := make(map[string]any, len(boards))
	for name, bdir := range boards {
		sb, err := OpenSwarmBoard(bdir)
		if err != nil {
			continue
		}
		st, _ := sb.Status()
		msgs, _ := sb.ReadMessages(0, nil)
		out[name] = map[string]any{"status": st, "messages": msgs}
	}
	return json.Marshal(out)
}
