# SPEC_EP16: LegoLoop Studio — Menu Lateral em Treeview Hierárquico e Sanfona Exclusiva

> **Padrão:** Spec-Driven Development (Addy Osmani)  
> **Referência:** [INTENT_EP16_SIDEBAR_TREEVIEW.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/INTENT_EP16_SIDEBAR_TREEVIEW.md)  
> **Status:** Aprovada  
> **Versão:** 1.0.0

---

## 1. 🎯 Requisitos Funcionais

1. **RF-01 (Confinamento de Viewport):**
   * A `.sidebar` deve conter estritamente seus filhos dentro de `height: calc(100vh - 56px)`.
   * O container de lista ativa deve ter `flex: 1; min-height: 0; overflow-y: auto; position: relative;`. A altura total da lista não pode exceder o espaço visível da tela.
2. **RF-02 (Contadores no Topo):**
   * O botão do accordion de peças deve exibir: `Peças Disponíveis (<total_ativas>)`.
   * O botão do accordion de modelos deve exibir: `Modelos Vivos (<total_ativos>)`.
3. **RF-03 (Treeview de Peças):**
   * A paleta de peças é dividida em subpastas recolhíveis por categoria:
     * `Controle` (3 blocos)
     * `Skills` (únicas consolidadas por SHA-256)
     * `MCPs`
     * `Loops`
     * `Scripts`
     * `Hooks`
     * `Gates`
   * Cada subpasta exibe o seu contador e a seta de expansão (`▾` quando aberta, `▸` quando fechada).
4. **RF-04 (Interação Ágil com Peças):**
   * Clique simples (`click`) em uma peça: Adiciona o nó no centro do Canvas e seleciona-o.
   * Duplo clique (`dblclick`): Abre o Modal Rico com o Raio-X detalhado da peça.
   * Drag-and-drop: Mantido para posicionamento livre no Canvas.
5. **RF-05 (Treeview de Modelos por Harness):**
   * Os modelos vivos são agrupados por provedor/harness:
     * `9Router VPS`
     * `Claude CLI`
     * `Ollama Local`
     * `Gemini / AGY`
   * Cada grupo exibe o status verde e a contagem de modelos ativos.
6. **RF-06 (Sanfona Exclusiva):**
   * Ao expandir "Peças Disponíveis", "Modelos Vivos" fica retraído na base inferior da viewport.
   * Ao expandir "Modelos Vivos", "Peças Disponíveis" retrai para o topo.
7. **RF-07 (Busca Inteligente):**
   * O campo de busca filtra em tempo real peças e modelos, auto-expandindo pastas com ocorrências encontradas.

---

## 2. ✅ Critérios de Aceite Binários

- [ ] **CA-01:** A altura computada de `#palette-list` ou qualquer sublista de peças é $\le 600\text{px}$, eliminando o bug de 208.620px.
- [ ] **CA-02:** O cabeçalho de peças exibe o contador formatado `Peças Disponíveis (N)`.
- [ ] **CA-03:** O cabeçalho de modelos exibe o contador formatado `Modelos Vivos (N)`.
- [ ] **CA-04:** As peças aparecem agrupadas dentro de suas respectivas pastas hierárquicas (`Controle`, `Skills`, `MCPs`, `Scripts`, `Hooks`, `Gates`).
- [ ] **CA-05:** Clicar uma única vez sobre uma peça da paleta adiciona o nó correspondente no Canvas.
- [ ] **CA-06:** Duplo clique na peça abre o modal rico com o Raio-X.
- [ ] **CA-07:** Os modelos vivos são listados em subpastas por Harness (`9Router`, `Claude`, `Ollama`, etc.).
- [ ] **CA-08:** A alternância entre "Peças Disponíveis" e "Modelos Vivos" move os botões sem gerar scroll na janela principal da sidebar.
- [ ] **CA-09:** Zero ocorrências de `innerHTML` (validador `TestNoInnerHTMLInScripts` verde).
- [ ] **CA-10:** Verificação E2E via Playwright/Obscura comprova funcionamento dos cliques e da sanfona sem acumular imagens em disco.
