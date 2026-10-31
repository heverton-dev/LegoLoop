# SPEC_EP15: LegoLoop Studio — Canvas N8N Enterprise, Controle de Fluxo, Ergonomia & Micro-Governança

> **Status:** Aprovada  
> **Versão:** 1.0.0  
> **Data:** 31/10/2026  
> **Autor:** LegoLoop Master Architecture  
> **Referência Master:** [SPEC.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/SPEC.md)

---

## 1. 🎯 Objetivo Executivo

Elevar o **LegoLoop Studio** ao padrão da indústria inspirado no **n8n**, eliminando todos os atritos ergonômicos e operacionais identificados na auditoria profunda:
1. **Navegabilidade & Ergonomia:** Sidebar esquerda (320px) com agrupamento por SHA-256 no catálogo do canvas; Inspector direito (280px) retrátil; recolhimento fluido via CSS Grid.
2. **Motor de Canvas N8N:** Cursor padrão (`default`/`pointer`), modo Pan ativado exclusivamente ao segurar a barra de **Espaço** (movimentação infinita $X/Y \in (-\infty, +\infty)$), e desconexão ágil de arestas via clique no botão `✕`.
3. **Nós Estruturais de Controle de Fluxo:** Introdução de nós lógicos (`IF / Branch` com portas `True`/`False`, `Loop / Recursão` com portas `Iterar`/`Concluído`, e `Merge / Barreira`).
4. **Animação Vetorial de Fluxo:** Arestas com animação luminosa em tempo real demonstrando a direção e o status de execução das peças.
5. **Micro-Governança Integrada:** Suporte a **Regras (`Rules`)**, Hooks (pré/pós), Gates e Scripts com badges no card e separação visual no Inspector entre slots **Nativos** (de fábrica) e **Customizados**.
6. **Modal Rico Enterprise com Dimensão Fixa:** Janela com tamanho estável (880px × 640px), ações diretas no cabeçalho (*"Canonizar Peça"* e *"Abrir no Explorer"*), diagrama **IPO** (*Input $\to$ Processamento $\to$ Output*), filtros por categoria e paginação interna na sub-aba de peças compatíveis, e botão de abrir pasta por ocorrência.
7. **Produtividade Visual Avançada:** Histórico Desfazer/Refazer (`Ctrl+Z`/`Ctrl+Y`), Auto-Layout em 1 clique, Snapping magnético no grid (20px), Menu de Contexto (botão direito) e Exportação do diagrama em PNG/SVG.

---

## 2. 🏛️ Arquitetura e Contratos de Dados

### 2.1 Backend Go (`internal/graph`, `internal/studio`, `internal/assets`)

1. **Flexibilização da Ontologia (`ontology.go`):**
   * `ParseArtifactType(s string)` aceita tanto nomes curtos (`URL`, `Spec`, `Code`, etc.) quanto prefixados (`Artifact:URL`, `Artifact:Spec`), corrigindo o erro 422 na rota ótima.
2. **Novos Tipos de Nós Estruturais (`Block.Kind`):**
   * `Kind: "flow_if"`: 1 entrada, 2 saídas (`true`, `false`), campo `Condition`.
   * `Kind: "flow_loop"`: 1 entrada, 2 saídas (`loop_body`, `loop_done`), campos `MaxIterations` e `StopCondition`.
   * `Kind: "flow_merge"`: N entradas, 1 saída sincronizada.
3. **Endpoint de Abertura no SO (`POST /api/assets/reveal`):**
   * Recebe `{ "path": "<caminho_absoluto_ou_com_til>" }`.
   * Valida segurança contra path traversal e chama nativamente no Windows:
     `exec.Command("explorer.exe", "/select,"+resolvedPath).Start()`.
4. **Endpoint de Canonização de Peça Única (`POST /api/canonize/piece`):**
   * Recebe `{ "sha256": "<hash>", "paths": [...] }`.
   * Executa a unificação no cofre central apenas para as duplicatas do item selecionado.

### 2.2 Frontend (`app.js`, `style.css`, `index.html`, `inventario.html`, `inventory.js`)

1. **Layout Reativo no Workspace:**
   ```css
   .workspace {
     display: grid;
     grid-template-columns: var(--sidebar-w, 320px) minmax(0, 1fr) var(--inspector-w, 280px);
     transition: grid-template-columns 200ms cubic-bezier(0.4, 0, 0.2, 1);
   }
   .workspace.is-sidebar-collapsed { --sidebar-w: 48px; }
   .workspace.is-inspector-collapsed { --inspector-w: 0px; }
   ```
2. **Canvas N8N Spacebar Pan:**
   * `state.isSpacePressed`: Alterna cursor para `grab` e habilita Pan 2D contínuo sem limites de quadrante.
   * Mouse pointer padrão: `default` para permitir clique direto e preciso no botão `✕` da aresta.
3. **Animação Vetorial de Arestas:**
   * CSS `@keyframes flowDash { to { stroke-dashoffset: -20px; } }` em arestas com classe `.edge--running` ou `.edge--animated`.
4. **Modal de Tamanho Fixo:**
   * `width: min(880px, 92vw); height: min(640px, 88vh);` fixados, com áreas internas roláveis independentes.

---

## 3. ✅ Critérios de Aceite Binários

- [ ] **CA-01:** O endpoint `/api/route?from=URL&to=Spec` responde HTTP 200 com a rota calculada sem erro 422.
- [ ] **CA-02:** A barra lateral esquerda exibe largura inicial de 320px e seus cards de catálogo têm duplicatas agrupadas por SHA-256 com badge de cópias.
- [ ] **CA-03:** Clicar no botão `◂` da sidebar reduz a coluna para 48px de forma fluida.
- [ ] **CA-04:** O inspector direito possui botão de recolher (`▸`), reduzindo sua coluna para 0px ao ser fechado.
- [ ] **CA-05:** No canvas, o cursor sobre o fundo é `default`; pressionar a barra de **Espaço** transforma o cursor em `grab` e possibilita Pan livre em $X$ e $Y$.
- [ ] **CA-06:** Clicar no botão `✕` no centro da aresta remove a conexão instantaneamente.
- [ ] **CA-07:** Nós de controle (`IF`, `Loop`, `Merge`) estão disponíveis na paleta e podem ser adicionados ao canvas com suas respectivas portas múltiplas.
- [ ] **CA-08:** As arestas conectadas exibem animação vetorial luminosa indicando o sentido do fluxo.
- [ ] **CA-09:** O card do nó no canvas exibe os badges discretos no rodapé (`rules`, `pre`, `gates`, `post`, `scripts`).
- [ ] **CA-10:** No Inspector, as seções de governança diferenciam claramente slots `[Nativo]` de slots `[Custom]`.
- [ ] **CA-11:** O Modal de Peça abre com dimensões fixas (880px × 640px), sem alterar seu tamanho durante a alternância entre abas.
- [ ] **CA-12:** A sub-aba "Visão Geral" do modal exibe o diagrama visual **IPO** e os badges dos artefatos entregues.
- [ ] **CA-13:** A sub-aba "Peças Compatíveis" contém botões de filtro por categoria e paginação interna.
- [ ] **CA-14:** A sub-aba "Ocorrências no SO" contém botão individual "Abrir Pasta" que aciona `POST /api/assets/reveal`.
- [ ] **CA-15:** Atalhos `Ctrl + Z` e `Ctrl + Y` desfazem e refazem alterações de nós e arestas no canvas.
- [ ] **CA-16:** `go test ./...` passa 100% verde e `funclen` mantém 0 funções com mais de 35 linhas.
