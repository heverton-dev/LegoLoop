package board

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func findProjectBoardDir(projDir string) string {
	candidates := []string{
		filepath.Join(projDir, ".legoloop-swarm", "swarm.json"),
		filepath.Join(projDir, ".legoloop-swarm", "swarm_board.json"),
		filepath.Join(projDir, ".swarm", "swarm.json"),
		filepath.Join(projDir, ".swarm", "swarm_board.json"),
		filepath.Join(projDir, "swarm_board.json"),
	}
	for _, p := range candidates {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return filepath.Join(projDir, ".legoloop-swarm")
		}
	}
	return ""
}

// DiscoverBoards localiza quadros Swarm em .legoloop-swarm/ e caminhos legados.
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
		if bdir := findProjectBoardDir(filepath.Join(root, e.Name())); bdir != "" {
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
		st, err := sb.Status()
		if err != nil {
			continue
		}
		msgs, err := sb.ReadMessages(0, nil)
		if err != nil {
			continue
		}
		out[name] = map[string]any{"status": st, "messages": msgs}
	}
	return json.Marshal(out)
}
