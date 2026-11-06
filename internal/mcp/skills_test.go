package mcp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkillSearchTaxonomyAndSkipDotted(t *testing.T) {
	// Cria estrutura temporária de teste isolada
	tmpDir := t.TempDir()

	// 1. Skill válida canônica com metadados
	tddDir := filepath.Join(tmpDir, "aidd-tdd")
	if err := os.MkdirAll(tddDir, 0755); err != nil {
		t.Fatalf("erro ao criar dir: %v", err)
	}
	skillContent := `---
name: aidd-tdd
description: Ciclo estrito Red-Green-Refactor com zero stubs e gates de qualidade.
tier: Ouro
natureza: Processo de Engenharia
categoria: build
---
# TDD Body
`
	if err := os.WriteFile(filepath.Join(tddDir, "SKILL.md"), []byte(skillContent), 0644); err != nil {
		t.Fatalf("erro ao criar SKILL.md: %v", err)
	}

	// 2. Skill dentro de pasta .upstream (deve ser IGNORADA)
	dottedDir := filepath.Join(tmpDir, "aidd-pocock", ".upstream", "engineering", "tdd-clone")
	if err := os.MkdirAll(dottedDir, 0755); err != nil {
		t.Fatalf("erro ao criar dotted dir: %v", err)
	}
	cloneContent := `---
name: aidd-tdd-clone
description: Clone redundante de upstream.
tier: Bronze
natureza: Clone
categoria: build
---
# Clone Body
`
	if err := os.WriteFile(filepath.Join(dottedDir, "SKILL.md"), []byte(cloneContent), 0644); err != nil {
		t.Fatalf("erro ao criar SKILL.md no dotted: %v", err)
	}

	engine := NewSkillEngine()
	if err := engine.IndexDirectory(tmpDir); err != nil {
		t.Fatalf("falha ao indexar diretorio: %v", err)
	}

	// Verifica se a pasta dotted foi podada (deve ter exatamente 1 skill no índice, não 2)
	allSkills := engine.List()
	if len(allSkills) != 1 {
		t.Fatalf("esperava exatamente 1 skill indexada (dotted ignorada), obteve %d", len(allSkills))
	}

	// Testa busca por query
	matches := engine.Search("red green refactor", 5)
	if len(matches) == 0 {
		t.Fatalf("esperava encontrar matches para 'red green refactor'")
	}

	top := matches[0]
	if top.ID != "aidd-tdd" {
		t.Errorf("ID esperado 'aidd-tdd', obteve '%s'", top.ID)
	}
	if top.Tier != "Ouro" {
		t.Errorf("Tier esperado 'Ouro', obteve '%s'", top.Tier)
	}
	if top.Natureza != "Processo de Engenharia" {
		t.Errorf("Natureza esperada 'Processo de Engenharia', obteve '%s'", top.Natureza)
	}
	if top.Categoria != "build" {
		t.Errorf("Categoria esperada 'build', obteve '%s'", top.Categoria)
	}
	if top.Resumo == "" {
		t.Errorf("Resumo nao pode ser vazio")
	}
}
