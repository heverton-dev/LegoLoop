# 📋 SPEC_EP23: Sub-grafos em Mega-Pipelines (Drill-Down, Breadcrumbs e Copy-on-Write)

> **Status:** Proposta para Execução  
> **Data:** 2026-11-01  
> **Épico Alvo:** ÉPICO 23 (Navegação em Camadas e Sub-grafos no Studio Canvas)  
> **Referência:** `docs/intents/2026-11-01_INTENT_ONTOLOGIA_ASSETS_SO_E_MEGA_PIPELINE_TEMPLATES.md`

---

## 1. Visão Geral e Arquitetura

O **ÉPICO 23** expande o Looping Studio para operar com **Mega-Pipelines Hierárquicos**, permitindo que qualquer esteira ou template complexo (como o Mega Pipeline de 5 loops) seja composto e navegado em múltiplos níveis de abstração:

1. **Macro-Cards Visuais para Sub-Loops (`.node--macro`):**
   * Nós com `piece.kind === "subloop"` ou `block.kind === "loop"` recebem estilização distinta de cartão composto (*Macro-Card*), com indicação visual de camadas aninhadas, badge "sub-loop" e contagem de nós internos quando disponível.

2. **Navegação Drill-Down por Duplo Clique (`dblclick`):**
   * O duplo clique sobre o Macro-Card ou acionamento pelo Inspector armazena o snapshot da esteira atual na pilha de navegação (`state.navStack`).
   * Faz a busca do sub-grafo via `/api/loops/{id}` e renderiza seus nós e arestas no Canvas com auto-layout limpo.

3. **Barra de Breadcrumbs no Canvas (`#canvas-breadcrumbs`):**
   * Barra horizontal localizada acima da superfície do canvas com a trilha de navegação (ex: `Principal > Sub-Loop: L4 Construção TDD`).
   * Clicar em qualquer segmento anterior desempilha o estado (`popNavStack`), restaurando o grafo do nível pai sem perda de posições ou conexões.
   * Totalmente estilizada com separadores vetoriais SVG e Zero Emojis.

4. **Política Copy-on-Write (Fork Versionado):**
   * Prevenção de corrupção acidental em templates compartilhados.
   * Suporte a disparo de duplicação/fork (`POST /api/loops/{id}/duplicate`) antes de mutações destrutivas.

---

## 2. Critérios de Aceitação Binários (Gates Mecânicos)

* [ ] **CA-23.1 (Macro-Cards no Canvas):** Nós do tipo sub-loop renderizam a classe visual `.node--macro` com badge identificador e evento de duplo clique associado.
* [ ] **CA-23.2 (Barra de Breadcrumbs no Topo do Canvas):** `#canvas-breadcrumbs` presente no HTML (`index.html`), oculta quando no nível raiz e ativada durante drill-down.
* [ ] **CA-23.3 (Pilha de Navegação Drill-Down e Retorno):** O duplo clique em um nó sub-loop salva o grafo pai em `state.navStack`, carrega o grafo filho e o clique no breadcrumb pai restaura o grafo original.
* [ ] **CA-23.4 (Ação no Inspetor para Acesso Direto):** O painel do Inspetor oferece o botão "Abrir Sub-grafo" para nós compostos.
* [ ] **CA-23.5 (Preservação da Barreira Zero Emojis):** Nenhum emoji em breadcrumbs, macro-cards ou botões de navegação.
* [ ] **CA-23.6 (Zero Regressões & 100% Testes Verdes):** Toda a suíte passa (`go test ./...`) e o binário é recompilado.
