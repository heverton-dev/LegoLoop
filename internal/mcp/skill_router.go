package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const skillTopKPadrao = 5

func (r *Router) SetSkills(engine *SkillEngine, root string) {
	r.skills = engine
	r.skillRoot = root
}

func (r *Router) callSkillSearch(id any, args map[string]any) JSONRPCResponse {
	if r.skills == nil {
		return r.toolCallError(id, "Cofre de skills nao configurado.")
	}
	query, _ := args["query"].(string)
	topK := skillTopKPadrao
	if v, ok := args["top_k"].(float64); ok && v >= 1 {
		topK = int(v)
	}
	matches := r.skills.Search(query, topK)
	if matches == nil {
		matches = []SkillMatch{}
	}
	data, err := json.Marshal(matches)
	if err != nil {
		return r.toolCallError(id, fmt.Sprintf("falha ao serializar: %v", err))
	}
	return r.toolCallText(id, string(data))
}

func (r *Router) callSkill(id any, args map[string]any) JSONRPCResponse {
	if r.skills == nil {
		return r.toolCallError(id, "Cofre de skills nao configurado.")
	}
	skillID, _ := args["skill_id"].(string)
	item, ok := r.skills.skills[skillID]
	if !ok {
		return r.toolCallError(id, fmt.Sprintf("skill nao encontrada: %s", skillID))
	}
	if !dentroDoCofre(r.skillRoot, item.CorpoPath) {
		return r.toolCallError(id, "caminho da skill fora do cofre")
	}
	body, err := os.ReadFile(item.CorpoPath)
	if err != nil {
		return r.toolCallError(id, fmt.Sprintf("falha ao ler skill %s: %v", skillID, err))
	}
	return r.toolCallText(id, string(body))
}

func dentroDoCofre(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func skillToolDefs() []map[string]any {
	return []map[string]any{
		{
			"name":        "skill_search",
			"description": "Busca skills no cofre por pontuacao de relevancia.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{"type": "string", "description": "Termos de busca"},
					"top_k": map[string]any{"type": "integer", "description": "Maximo de resultados"},
				},
				"required": []string{"query"},
			},
		},
		{
			"name":        "call_skill",
			"description": "Devolve o texto do SKILL.md de uma skill do cofre.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"skill_id": map[string]any{"type": "string", "description": "Identificador da skill"},
				},
				"required": []string{"skill_id"},
			},
		},
	}
}
