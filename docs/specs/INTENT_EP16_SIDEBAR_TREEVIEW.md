# INTENT_EP16: Menu Lateral Esquerdo — Arquitetura em Treeview Hierárquico, Sanfona Exclusiva e Confinamento de Viewport

> **Status:** Aprovado para Execução  
> **Data:** 31/10/2026  
> **Caso de Teste do Mega-Pipeline:** Refatoração Estrutural da Sidebar do LegoLoop Studio

---

### 1. 📌 Problemas Diagnosticados & Foco da Intenção

1. **Colapso Estrutural da Viewport (Bug dos 208.620px):**
   * A lista de peças (`#palette-list`) foi colocada dentro de um container sem restrição flexível (`height: auto`). A virtualização calculou que o viewport tinha 208 mil pixels e renderizou 4.346 nós de uma só vez, quebrando o scroll e tornando as peças invisíveis no navegador.
2. **Poluição Vertical dos Chips/Pillows:**
   * 8 botões de categoria no topo quebraram em 3 linhas sobrepostas, consumindo quase 100px verticais de área nobre em uma barra de 320px.
3. **Modelos Vivos Esmagados e sem Agrupamento:**
   * 225 modelos da VPS 9Router ficaram misturados em uma única coluna sem filtro, sem respiro e sem identificação clara de harness.
4. **Clique Simples Inerte:**
   * O usuário clica em uma peça e nada acontece (o código exigia exclusivamente duplo clique ou drag).

---

### 2. 🏛️ Escopo da Solução: Treeview em Sanfona Inteligente

A Sidebar passa a adotar a estrutura consagrada de **Treeview / Sanfona Hierárquica**:

```
┌────────────────────────────────────────────────────────┐
│ [Catálogo]                                        [◂]  │
├────────────────────────────────────────────────────────┤
│ 🔍 Buscar peça ou modelo...                        (/) │
├────────────────────────────────────────────────────────┤
│ ▼ PEÇAS DISPONÍVEIS (769 únicas · 6.122 no SO)        │
│   ▸ Controle (3)                                       │
│   ▼ Skills (769 agrupadas)                             │
│     📦 AIDD Grill        Idea ➔ Briefing               │
│     📦 AIDD Spec         Briefing ➔ Spec               │
│   ▸ MCPs (28)                                          │
│   ▸ Loops (0)                                          │
│   ▸ Scripts (282)                                      │
│   ▸ Hooks (20)                                         │
│   ▸ Gates (6)                                          │
├────────────────────────────────────────────────────────┤
│ ► MODELOS VIVOS (225 ativos)                           │
│   ▸ 9Router VPS (218)                                  │
│   ▸ Claude CLI (3)                                     │
│   ▸ Ollama Local (1)                                   │
│   ▸ Gemini CLI (3)                                     │
└────────────────────────────────────────────────────────┘
```

1. **Cabeçalhos com Contadores Vivos:**
   * `Peças Disponíveis (N ativas)`
   * `Modelos Vivos (N ativos)`
2. **Comportamento Sanfona Exclusiva:**
   * Clicar em "Peças Disponíveis" expande o catálogo de peças e empurra "Modelos Vivos" para o rodapé da viewport (`calc(100vh - 56px)`).
   * Clicar em "Modelos Vivos" retrai as peças e expande os provedores e modelos.
   * A rolagem acontece **estritamente dentro da pasta aberta**, com viewport contido e virtualização de 60fps.
3. **Interação com Peças:**
   * 1 clique simples $\to$ Adiciona a peça no centro do Canvas com feedback visual.
   * Duplo clique $\to$ Abre o Modal Rico de Raio-X.
   * Arrastar $\to$ Solta o nó na coordenada exata do ponteiro.
4. **Modelos Agrupados por Provedor/Harness:**
   * Categorizados em subpastas: `9Router VPS`, `Claude CLI`, `Ollama`, `Gemini`.
