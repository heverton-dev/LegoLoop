package mcp

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// SkillItem representa uma skill indexada pelo motor.
type SkillItem struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Tier      string `json:"tier"`
	Natureza  string `json:"natureza"`
	Categoria string `json:"categoria"`
	Resumo    string `json:"resumo"`
	CorpoPath string `json:"corpo_path"`
}

// SkillMatch retorno enriquecido com score de busca.
type SkillMatch struct {
	ID        string  `json:"skill_id"`
	Nome      string  `json:"name"`
	Score     float64 `json:"score"`
	Natureza  string  `json:"natureza"`
	Tier      string  `json:"tier"`
	Categoria string  `json:"categoria"`
	Resumo    string  `json:"resumo"`
}

// SkillEngine gerencia o índice e busca de skills.
type SkillEngine struct {
	skills map[string]SkillItem
}

// NewSkillEngine instancia um novo motor de skills.
func NewSkillEngine() *SkillEngine {
	return &SkillEngine{skills: make(map[string]SkillItem)}
}

// IndexDirectory varre recursivamente e indexa skills, ignorando pastas com ponto.
func (e *SkillEngine) IndexDirectory(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.EqualFold(d.Name(), "SKILL.md") {
			item, parseErr := parseSkillFile(path)
			if parseErr == nil && item.ID != "" {
				e.skills[item.ID] = item
			}
		}
		return nil
	})
}

// List retorna todas as skills indexadas.
func (e *SkillEngine) List() []SkillItem {
	res := make([]SkillItem, 0, len(e.skills))
	for _, item := range e.skills {
		res = append(res, item)
	}
	return res
}

// Search busca skills relevantes com base na query.
func (e *SkillEngine) Search(query string, topK int) []SkillMatch {
	terms := tokenize(query)
	var matches []SkillMatch
	for _, s := range e.skills {
		score := calculateScore(terms, s)
		if score > 0 {
			matches = append(matches, SkillMatch{
				ID:        s.ID,
				Nome:      s.Nome,
				Score:     score,
				Natureza:  s.Natureza,
				Tier:      s.Tier,
				Categoria: s.Categoria,
				Resumo:    s.Resumo,
			})
		}
	}
	return topMatches(matches, topK)
}

func parseSkillFile(path string) (SkillItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return SkillItem{}, err
	}
	defer f.Close()

	item := SkillItem{
		ID:        filepath.Base(filepath.Dir(path)),
		Nome:      filepath.Base(filepath.Dir(path)),
		Tier:      "Prata",
		Natureza:  "Processo Canônico",
		Categoria: "general",
		CorpoPath: path,
	}

	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if inFrontmatter {
			extractMetaLine(line, &item)
		}
	}
	return item, nil
}

func extractMetaLine(line string, item *SkillItem) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return
	}
	key := strings.ToLower(strings.TrimSpace(parts[0]))
	val := strings.TrimSpace(parts[1])
	switch key {
	case "name":
		item.Nome = val
		item.ID = val
	case "description":
		item.Resumo = val
	case "tier":
		item.Tier = val
	case "natureza":
		item.Natureza = val
	case "categoria":
		item.Categoria = val
	}
}

func tokenize(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	var res []string
	for _, w := range words {
		w = strings.Trim(w, ",.-_!?;:()[]")
		if len(w) > 1 {
			res = append(res, w)
		}
	}
	return res
}

func calculateScore(terms []string, item SkillItem) float64 {
	haystack := strings.ToLower(item.ID + " " + item.Nome + " " + item.Resumo + " " + item.Categoria + " " + item.Natureza)
	var matched int
	for _, t := range terms {
		if strings.Contains(haystack, t) {
			matched++
		}
	}
	if len(terms) == 0 {
		return 0
	}
	return float64(matched) / float64(len(terms))
}

func topMatches(matches []SkillMatch, topK int) []SkillMatch {
	if len(matches) <= topK {
		return matches
	}
	return matches[:topK]
}
