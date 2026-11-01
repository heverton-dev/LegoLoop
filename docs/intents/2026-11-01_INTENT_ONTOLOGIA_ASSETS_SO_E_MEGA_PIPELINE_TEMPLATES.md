# 🎯 INTENT: Ontologia Universal, Canvas Studio Impeccable (Zero Emojis), Novas Peças & Mega-Pipelines

> **Status:** Proposta Expandida Consolidada (STOP-GATE HUMANO 01)  
> **Data:** 2026-11-01  
> **Autor:** Antigravity AI + Usuário  
> **Alvo:** `internal/assets/`, `internal/graph/`, `internal/studio/`, `internal/studio/static/`, Templates, CLI

---

## 1. Contexto e Motivação Factual

O LegoLoop está evoluindo para se tornar o **Ambiente de Engenharia e Operação Contínua Definitivo (IDE Local + Runtime Autônomo)**.
Para cumprir essa missão, identificamos lacunas factuais críticas no motor e na interface:

1. **Canvas do Studio Inoperante:** O Canvas atual em `internal/studio/static/` apresenta quebra generalizada de eventos: scroll/wheel não responde adequadamente, atalhos de teclado (shortcuts) inativos, captura de ponteiro (pointer events/drag) instável e cliques de seleção nos cards não funcionando.
2. **Violação da Regra de Design (Presença de Emojis):** O frontend faz uso deliberado de emojis (ex: `📋`, `🗄️`, `📸`, `🌱`, `📦`), o que quebra a estética corporativa e a identidade do produto.
3. **Assimetria de Contrato:** Ativos do SO (skills em `~/.gemini/config/skills/`, scripts e MCPs) são lidos apenas com metadados básicos, sem inferência de portas ontológicas (`Consumes`, `Produces`, `Gate`).
4. **Ausência de Peças Essenciais de Controle, Segurança e Entrega:** Faltam peças críticas para fluxos reais de produção, tais como o **Stop-Gate Humano (`human-gate`)**, o **Cofre de Variáveis (`env-vault`)**, automação Git (`commit`, `merge`) e implantação (`deploy-vps`, `deploy-vercel`, `deploy-docker`).
5. **Dificuldade de Navegação em Loops Compostos:** Falta suporte a sub-grafos com navegação drill-down e breadcrumbs em Mega-Pipelines.

---

## 2. Leis Inquebráveis e Diretrizes Críticas

> [!CAUTION]
> **LEI INQUEBRÁVEL DE DESIGN — ZERO EMOJIS NO FRONTEND:**  
> A aplicação LegoLoop **NÃO PODE**, em hipótese alguma, utilizar emojis em botões, abas, cards, títulos ou badges.  
> Todos os emojis devem ser 100% erradicados e substituídos por um **Sistema Canônico de Ícones Vetoriais SVG** (limpos, nítidos, padronizados e sem dependências pesadas).

> [!IMPORTANT]
> **CONTRATO GLOBAL DE ATIVOS DO SO:**  
> Novas peças e skills criadas ou importadas de terceiros (ex: via GitHub) devem ser **obrigatoriamente salvas a nível global** (`~/.gemini/config/skills/`), passando por **canonização e hashing SHA-256 imediata** (`legoloop canonize`).

---

## 3. Matriz de Novas Peças e Taxonomia Estrutural de Cores

### 3.1. Código de Cores por Natureza da Peça no Canvas
| Natureza da Peça | Papel no Canvas | Cor Semântica | Hex | Exemplo |
|---|---|---|---|---|
| ⚡ **Triggers** | Gatilhos de partida do fluxo | Laranja Elétrico | `#F97316` | `cron`, `webhook`, `file-watcher` |
| 🎛️ **Controles** | Decisão, auditoria humana e fluxo | Âmbar Ouro | `#F59E0B` | `human-gate` (Stop-Gate), `if-else`, `router` |
| 🧠 **Skills** | Agentes cognitivos e raciocínio LLM | Roxo Real | `#8B5CF6` | `aidd-spec`, `aidd-tdd`, `aidd-grill` |
| 🔌 **MCPs** | Protocolo de ferramentas externas | Ciano Técnico | `#06B6D4` | `chrome-devtools`, `supabase`, `filesystem` |
| 🔄 **Loops** | Sub-grafos e esteiras aninhadas | Azul Royal | `#2563EB` | `mega-pipeline`, `sub-loop-tdd` |
| 📜 **Scripts** | Execução determinística de código | Verde Esmeralda | `#10B981` | `python.py`, `powershell.ps1`, `bash.sh` |
| 🪝 **Hooks** | Interceptadores de ciclo de vida | Rosa Magenta | `#EC4899` | `pre-commit`, `post-build`, `guard` |
| 🛑 **Gates** | Validações mecânicas de saída | Vermelho Carmim | `#EF4444` | `go-test`, `linter`, `budget-breaker` |

### 3.2. Catálogo de Novas Peças Estratégicas
1. **`human-gate` (Stop-Gate Humano):** Pausa o loop, alerta o usuário e exige aprovação visual (com feedback/rejeição) para prosseguir.
2. **`env-vault` (Cofre de Segredos & Variáveis):** Armazena chaves de API e variáveis de ambiente com segurança, injetando-as sob demanda em MCPs e Scripts sem expô-las no Canvas ou no DAG JSON.
3. **`git-commit`, `git-merge`, `git-tag-release`:** Automação de versionamento convencional e esteira Git local.
4. **`deploy-vps`, `deploy-vercel`, `deploy-docker`:** Peças de entrega contínua com parâmetros configuráveis no Inspector.
5. **`hermetic-wall` & `fanout-matrix` / `fanin-consolidator`:** Controle de esterilidade de contexto e paralelização DAG.

---

## 4. Arquitetura do Canvas e Telemetria

1. **Auditoria Forense (Loop REA):**
   * Engenharia reversa em `internal/studio/static/app.js`, `slots.js`, `ui.js` para restaurar o pipeline de eventos de mouse, zoom, scroll e atalhos de teclado.
2. **Design Pass (Loop Impeccable):**
   * Erradicação de 100% dos emojis e substituição por biblioteca de SVGs inline canônicos.
   * Aplicação de tokens de espaçamento, tipografia monoespaçada/sans corporativa e alto contraste.
3. **Telemetria Dupla em Tempo Real:**
   * **Por Card:** Indicador de tempo de execução (`⏱ 1.8s`), consumo de tokens (`In: 1.4k / Out: 320`) e status do gate.
   * **No Topo do Canvas:** Header global consolidado exibindo:
     * Tempo Total de Execução da Esteira.
     * Tokens de Entrada Totais ($\sum \text{Input}$).
     * Tokens de Saída Totais ($\sum \text{Output}$).
     * Custo e modelo consolidado.
4. **Mega-Pipelines & Sub-Loops:**
   * Visualização como cartões compostos (`Macro-Cards`).
   * Drill-down via duplo clique com barra breadcrumb de navegação de volta.
   * Política de edição Copy-on-Write (Fork versionado sem corromper o template original).

---

## 5. Estrutura dos Épicos de Implementação

* **🏛️ ÉPICO 18 — Engenharia Reversa (REA) & Restauração Funcional do Canvas:**
  * Diagnóstico e conserto dos event listeners: scroll/pan, zoom, atalhos de teclado, seleção e arrasto de cards.
* **🏛️ ÉPICO 19 — Design Pass Impeccable (Regra Zero Emojis & Sistema de Ícones SVG):**
  * Eliminação total de emojis no frontend e backend; implantação do sistema canônico de ícones SVG e paleta de cores ontológica por tipo de peça.
* **🏛️ ÉPICO 20 — Inferência de Contrato Ontológico de Ativos do SO & Peça `human-gate`:**
  * Leitura e inferência de portas em `SKILL.md`, scripts e MCPs; implementação da peça e modal do Stop-Gate Humano.
* **🏛️ ÉPICO 21 — Cofre de Credenciais (`env-vault`) & Peças de Operações Git e Deploy:**
  * Implementação da peça de injeção de variáveis de ambiente e peças parametrizáveis de Git e Deploy (VPS, Vercel, Docker).
* **🏛️ ÉPICO 22 — Telemetria em Tempo Real (Por Card e Header Global no Canvas):**
  * Medição e exibição de tempo, tokens de input, tokens de output e custos agregados.
* **🏛️ ÉPICO 23 — Sub-grafos em Mega-Pipelines (Drill-Down, Breadcrumbs e Templates de 5 Loops):**
  * Suporte a loops aninhados, navegação em camadas e template do Mega Pipeline no Studio.
* **🏛️ ÉPICO 24 — CRUD em Camadas, Fork/Versionamento e Importador Global do GitHub com Canonização HDSM:**
  * Importação via chat/URL direto para `~/.gemini/config/skills/`, canonização automática com SHA-256 e cópia protegida.

---

## 6. STOP-GATE DE AUDITORIA HUMANA 01

Esta especificação de intenção cobre integralmente todas as suas exigências funcionais, estéticas (Zero Emojis) e de robustez sistêmica.

Com a sua aprovação, fechamos o **Loop 1 (Concepção & Alinhamento)** e avançamos para a **Tríade de Planejamento (`SPEC` + `PLAN` + `ROUTE`)** iniciando a execução pelo **ÉPICO 18 (Engenharia Reversa e Conserto do Canvas)** e **ÉPICO 19 (Design Pass Impeccable / Zero Emojis)**!
