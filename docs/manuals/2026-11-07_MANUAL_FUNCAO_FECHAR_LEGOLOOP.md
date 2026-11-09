# 📘 MANUAL OPERACIONAL: FUNÇÃO FECHAR NO LEGOLOOP

> **Data:** 2026-11-07  
> **Autoridade:** Ecossistema LegoLoop  
> **Status:** Canônico  
> **Escopo:** Sub-loop universal de canonização e fechamento determinístico em 9 etapas  
> **Fontes de verdade:** `internal/fechar/` · `internal/cli/fechar.go` · `skills/legoloop-fechar/SKILL.md` · `docs/templates/TEMPLATE_REPORT.md`

---

## 1. Visão Geral

A função **`legoloop-fechar`** é o 6º verbo de ciclo de vida do LegoLoop e opera como o nó final compulsório de todas as esteiras de criação e evolução. Ela encadeia deterministicamente as 9 etapas que transformam uma implementação aprovada em um artefato canônico universalmente integrado ao ecossistema:
1. **Manual Didático:** em `docs/manuals/YYYY-MM-DD_MANUAL_FUNCAO_<SLUG>_LEGOLOOP.md`.
2. **Livro Didático LegoPress:** em `docs/books/<slug>/` com auditoria geométrica.
3. **Skill Universal:** em `skills/legoloop-<slug>/SKILL.md` vinculada aos harnesses do SO.
4. **Loop & Grafo DAG:** em `loops/<slug>/pipeline.dag.json` validado.
5. **Recompilação Fresca do Binário:** `bin/legoloop.exe` e `go install ./cmd/legoloop`.
6. **Governança Viva:** `AGENTS.md` e `MEMORY.md`.
7. **Reindexação do Grafo Semântico:** via `codebase-memory-mcp`.
8. **Portões Mecânicos de Qualidade & Commit:** testes, vet, funclen e commit seguro.
9. **Relatório de Conclusão & Telemetria Estrita:** em `docs/reports/YYYY-MM-DD_REPORT_<SLUG>.md` com balanço orçado vs realizado de tempo, tokens e custo USD/BRL.

---

## 2. Comandos e Uso

### 2.1 Pelo Chat do Agente
* `/legoloop-fechar <slug>`
* `/fechar <slug>`

### 2.2 Pela CLI Mecânica
```bash
# Simulação sem efeitos colaterais
legoloop fechar video-motion --dry-run --json

# Fechamento com parâmetros explícitos de telemetria
legoloop fechar video-motion --tempo-min 25 --tokens-in 60000 --tokens-out 15000 --custo-usd 0.40

# Fechamento e canonização real padrão
legoloop fechar video-motion
```
