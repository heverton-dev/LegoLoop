package skills

import (
	"os"
	"path/filepath"

	"legoloop/internal/processo"
)

// conteudoIndice é a skill-índice do modo lazy (~150 tokens).
const conteudoIndice = `---
name: legoloop-skills
description: Índice das skills do LegoLoop. Use quando precisar de uma skill que não está listada; busque por skill_search e leia por call_skill.
---

# Skills sob demanda (lazy)

As skills ficam no Cofre, fora do contexto. Carregue só a que a tarefa pedir.

1. Busque: ferramenta ` + "`skill_search`" + ` do MCP ` + "`legoloop`" + ` com ` + "`query`" + ` (palavras da tarefa) e ` + "`top_k`" + ` 5.
2. Escolha o ` + "`skill_id`" + ` mais relevante do resultado.
3. Leia: ferramenta ` + "`call_skill`" + ` com ` + "`skill_id`" + `. O retorno é o SKILL.md completo; siga-o.

Regras:
- Não liste o Cofre inteiro. Busque pela tarefa.
- Se skill_search voltar vazio, reformule a query com outros termos antes de desistir.
- Os comandos legoloop (fabricar, fechar, listar, matriz) continuam disponíveis na CLI.
`

// garantirIndice grava a skill-índice na pasta de skills de cada harness (cópia, nunca link).
func garantirIndice(op Opcao, hs []processo.Harness) ([]string, error) {
	var out []string
	for _, h := range hs {
		dest := filepath.Join(h.Skills, processo.IndiceSkills, "SKILL.md")
		out = append(out, filepath.Dir(dest))
		if op.DryRun {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(dest, []byte(conteudoIndice), 0o644); err != nil {
			return nil, err
		}
	}
	return out, nil
}
