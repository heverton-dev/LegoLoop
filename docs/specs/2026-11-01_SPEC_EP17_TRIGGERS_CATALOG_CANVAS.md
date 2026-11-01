# 🧩 SPEC_EP17: Triggers de Primeira Classe no Catálogo Lateral e Canvas N8N

> **Documento:** `docs/specs/2026-11-01_SPEC_EP17_TRIGGERS_CATALOG_CANVAS.md`  
> **Status:** Aprovado para Decomposição  
> **Data:** 01/11/2026  
> **Autor:** LegoLoop Architecture Core + AIDD  
> **Referência:** Contrato DAG v2 (`internal/dag/types.go`), N8N Trigger Standard

---

## 1. Visão Geral

Esta especificação define a promoção dos **Triggers (Gatilhos)** a cidadãos de primeira classe no **LegoLoop Studio**:
1. Nova categoria retrátil **Triggers** no menu lateral esquerdo, posicionada no topo da sanfona (antes de *Controle*).
2. Peças canônicas para os 5 gatilhos do ecossistema: `manual`, `cron`, `file_change`, `loop_end` e `trigger_board`.
3. Comportamento no Canvas N8N: nó raiz de entrada sem conexões de entrada (`consumes`), com porta de saída (`produces`), estilizado com destaque de nó de início.
4. Inspetor lateral com parâmetros específicos (`cron`, `path`, `loop`) e injeção determinística no contrato DAG v2 durante o `/api/compose`.

---

## 2. Critérios de Aceitação Binários (Gates)

- [ ] **CA-01 (Ontologia Go):** `internal/graph/ontology.go` expõe a função `TriggerBlocks() []Block` com os 5 blocos canônicos (`trigger-manual`, `trigger-cron`, `trigger-file-change`, `trigger-loop-end`, `github-board-trigger`), todos com `Kind` iniciando em `"trigger_"`.
- [ ] **CA-02 (FullCatalog):** `FullCatalog()` inclui `TriggerBlocks()`. `ControlBlocks()` é expurgado de nós do tipo trigger.
- [ ] **CA-03 (Sidebar Treeview):** O menu lateral esquerdo exibe a subpasta `Triggers` como primeiro item do accordion de peças disponíveis, com ícone de raio SVG e badge com o total de gatilhos disponíveis (mínimo 5).
- [ ] **CA-04 (Classificação de Peças):** O frontend `palette.js` classifica automaticamente blocos com `kind.startsWith("trigger_")` na categoria `trigger` com badge `"gatilho"`.
- [ ] **CA-05 (Inserção no Canvas):** O usuário pode clicar ou arrastar um bloco Trigger para o Canvas. O nó renderizado recebe a classe visual `.node--trigger` e não expõe conector de entrada (somente conector de saída).
- [ ] **CA-06 (Restrição de Gatilho Único):** Ao adicionar um nó de Trigger quando outro já existe no Canvas, o sistema substitui o trigger anterior ou notifica com toast determinístico, garantindo a invariante DAG v2 de exatamente 1 trigger por esteira.
- [ ] **CA-07 (Inspetor de Gatilho):** Ao selecionar um nó Trigger, o painel Inspetor renderiza os campos de configuração pertinentes (`cron`, `path` vigiado, ou nome do `loop` anterior).
- [ ] **CA-08 (Serialização DAG v2):** Ao clicar em *Validar*, *Executar* ou *Materializar*, a requisição enviada a `/api/compose` injeta o objeto `trigger` correspondente no payload JSON.
- [ ] **CA-09 (Testes & Recompilação):** Testes unitários do Go (`go test ./...`) e frontend passam 100% sem erros, e o binário `bin/legoloop.exe` é recompilado com sucesso.

---

## 3. Especificação Detalhada dos Blocos de Trigger

| ID do Bloco | Label | Kind | Produz | Artefato Padrão | Gate Padrão | Parâmetros Específicos |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `trigger-manual` | Disparo Manual | `trigger_manual` | `Idea` | `TRIGGER_MANUAL.json` | `file_exists('TRIGGER_MANUAL.json')` | N/A (disparo interativo) |
| `trigger-cron` | Agendador Cron | `trigger_cron` | `Idea` | `TRIGGER_CRON.json` | `file_exists('TRIGGER_CRON.json')` | `cron: "*/15 * * * *"` |
| `trigger-file-change`| Observador de Arquivo | `trigger_file_change` | `Idea` | `TRIGGER_FILE.json` | `file_exists('TRIGGER_FILE.json')` | `path: "docs/intents/"` |
| `trigger-loop-end` | Encadeamento de Loop | `trigger_loop_end` | `Idea` | `TRIGGER_LOOP.json` | `file_exists('TRIGGER_LOOP.json')` | `loop: "<loop_id>"` |
| `github-board-trigger`| Trigger: GitHub Board Kanban | `trigger_board` | `Idea` | `WAVE_PLAN.json` | `file_exists('WAVE_PLAN.json')` | `board: "backlog"` |

---

## 4. Design & Ergonomia Visual

```css
/* Estilo N8N para nós de entrada / Trigger */
.node--trigger {
  border-color: var(--color-amber-500, #f59e0b) !important;
  box-shadow: 0 0 12px rgba(245, 158, 11, 0.2);
}
.node--trigger .node__port--in {
  display: none !important;
}
.node--trigger .node__type-badge {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
}
```
