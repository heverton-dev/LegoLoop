# 🎯 INTENT: Taxonomia Completa de Peças Híbridas, Triggers e Governança de Slots Canônicos

> **Status:** Aprovado pelo Usuário (Stop-Gate 01 Concluído)  
> **Data:** 2026-11-02  
> **Autor:** Antigravity AI + Usuário  
> **Alvo:** `internal/assets/`, `internal/dag/`, `web/`, motor de execução e Looping Studio  

---

## 1. Contexto e Motivação
O ecossistema LegoLoop possui 8 tipos canônicos de ativos (`internal/assets/types.go`): **Skill, MCP, Script, Hook, Gate, Harness, Binary e Loop**, além da entidade de disparo **Trigger** (`internal/dag/types.go`).

Historicamente, algumas peças operavam apenas embutidas em slots (`Gates`, `Hooks`), outras apenas no grafo macro (`Skills`, `Subloops`), e outras como ambos (`Scripts`). Além disso, peças prontas vindas de packs do sistema (ex.: Impeccable) chegam com slots pré-configurados que não devem ser editados diretamente no acervo core.

Esta intenção estabelece a **taxonomia formal e exaustiva** de quais peças são elegíveis a operar no **Modelo Híbrido (*Dual-Role*)**, quais são estritamente unitárias/estruturais, como os **Triggers** se comportam como a cabeça mandatória de todo loop, e a governança de **Imutabilidade Canônica com Extensibilidade** (Slot Override vs. Fork).

---

## 2. Taxonomia Canônica: Quem é Híbrido e Quem Não É

### Grupo A: Peças Híbridas (*Dual-Role* — Podem ser Nó Externo ou Componente Interno)

| Peça (`Kind`) | Papel como Componente Interno (Slot / Middleware) | Papel como Peça Externa (Nó de 1ª Classe no DAG) | Justificativa Arquitetural |
| :--- | :--- | :--- | :--- |
| **Gate** | Validação atômica de saída do nó (`slots.gates`). | **Stop-Gate ou Gate de Decisão** no DAG (`type: "gate"`): pausa para aprovação humana ou bifurcação condicional. | Um gate pode ser um assert local de schema ou um marco de aprovação de toda uma fase. |
| **Hook** | Gatilhos de ciclo de vida atômicos (`pre`/`post` execução daquele nó específico). | **Event Listener de Pipeline** (`type: "hook"`): interceptador de eventos do grafo (`on_pipeline_error`, `on_pipeline_finish`). | Pode proteger o ambiente de um único nó ou reagir a eventos globais do fluxo. |
| **Script** | Comando utilitário embutido (`slots.scripts`). | **Worker de Execução** (`type: "script"`): comando determinístico com dependências, I/O e telemetria. | Já existente e funcional no motor atual. |
| **Skill** | **Sub-Skill / Persona Auxiliar:** Injetada no prompt de outro nó como helper de regras ou tom de voz. | **Nó de IA (Padrão):** O processador cognitivo principal que despacha prompts para modelos de LLM. | Permite composição modular de personas dentro de um nó executor. |
| **MCP** | **Tool Provider de IA:** Servidor MCP vinculado ao nó para que a LLM invoque tools contextuais. | **Nó de Tool Call Direta:** Disparo de uma ferramenta MCP sem chamar LLM (execução mecânica determinística). | Economiza tokens quando a operação exige apenas a execução crua de uma tool do MCP. |
| **Binary** | Utilitário local invocado em background por scripts do nó. | **Nó Executável:** Chamada direta de binário compilado no fluxo do pipeline com flags de linha de comando. | Permite tratar ferramentas de sistema como passos independentes de compilação ou empacotamento. |

---

### Grupo B: Peças Estritamente Não-Híbridas (Unicidade Estrutural)

| Peça (`Kind` / Entidade) | Papel Exclusivo | Por que NÃO Pode ser Híbrida? |
| :--- | :--- | :--- |
| **Loop / Subloop** | **Nó Macro Composto:** Representa um pipeline DAG encapsulado com sua própria topologia e isolamento. | Possui camadas, ciclo de vida e estado próprios; embuti-lo como "slot" de um nó atômico violaria o princípio de esterilidade do Orchestrix. |
| **Trigger** | **Cabeça de Entrada do Loop:** Gatilho universal que dispara a esteira (`manual`, `cron`, `file_change`, `loop_end`). | **Todo loop possui obrigatoriamente um Trigger.** Ele é o iniciador da esteira e define como/quando o DAG acorda; não é uma peça interna de nós subsequentes. |
| **Harness** | **Ambiente / Runtime de Orquestração:** (ex.: Claude Code, AGY, OpenCode). | É a infraestrutura de despacho e tráfego que executa os modelos; não é nó de fluxo nem slot, mas sim o canal subjacente. |

---

## 3. Diretriz dos Triggers: Cabeça Universal de Todo Loop
* **Todo Loop tem um Trigger:** Nenhum DAG opera sem ponto de partida explícito.
* **Tipos Canônicos:**
  * `manual`: Disparo pontual via CLI (`legoloop run`) ou UI (botão "Run" do Studio).
  * `cron`: Agendamento temporal determinístico.
  * `file_change`: Reatividade a alterações no sistema de arquivos do projeto.
  * `loop_end`: Encadeamento reativo disparado após o encerramento com sucesso de outro loop.
* **Visualização no Studio:** O Trigger é renderizado visualmente como o nó inicial `#trigger-start`, conectado aos nós de entrada do DAG.

---

## 4. Governança de Slots Herdados de Peças Canônicas (Ex.: Impeccable)

As peças canônicas do acervo do LegoLoop são **imutáveis (read-only)**. O usuário possui duas formas legítimas de personalização:

1. **`slots_override` no DAG (Sem Alterar Arquivos de Acervo):**
   * No `pipeline.dag.json`, o usuário pode adicionar (`add_gates`, `add_scripts`) ou suprimir (`remove_gates`, `remove_scripts`) slots da peça referenciada sem tocar no arquivo canônico.
   * Garante que melhorias upstream do LegoLoop continuem sendo recebidas automaticamente.
2. **Fork da Peça (*Clone-to-Craft*):**
   * O usuário clona a peça para o escopo local (`.legoloop/skills/minha-copia/`).
   * A cópia local torna-se um ativo desvinculado, 100% editável pelo usuário.

---

## 5. Épicos — IDs Canônicos (Pivot de Rastreabilidade)

| ID | Nome | Descrição Resumida |
| :--- | :--- | :--- |
| **E01** | **Taxonomia & Expansão de Tipos no DAG** | Atualização de `NodeType` (`gate`, `hook`, `mcp`, `binary`), unificação da semântica dos nós e validação estrita em `internal/dag/`. |
| **E02** | **Padronização Universal de Triggers** | Garantir que todo loop instancie e valide seu Trigger canônico, com representação formal de I/O de entrada. |
| **E03** | **Governança de Slots & Mecanismo de Override** | Implementação de `slots_override` na struct de `Node` e suporte a Fork de peças canônicas. |
| **E04** | **Adaptação Visual no Looping Studio** | Representação do nó Trigger como início visual, nós de Gates/Hooks/MCPs no canvas e painel de overrides no Inspetor. |

---

## 6. Decisões Consolidadas com o Usuário
1. **Comportamento do Nó de Gate Externo:**
   * **Opção (C) Adotada:** Suportar **ambas as modalidades**.
     * **Modalidade 1 (Stop-Gate Binário):** Pausa para aprovação humana com validação binária (`PASS` prossegue, `FAIL` pausa/bloqueia).
     * **Modalidade 2 (Bifurcação Condicional):** Roteia dinamicamente o fluxo para caminhos alternativos (`on_pass` ➔ Nó A, `on_fail` ➔ Nó B).
2. **Precedência do Override & Uso Canônico:**
   * **Padrão Nativo (Sem Override):** Se `slots_override` for nulo ou omitido no nó do DAG, o nó consome **100% da definição canônica original da peça** intacta, sem nenhuma alteração.
   * **Com Override Explícito:** Somente quando o usuário declarar `slots_override` (ex.: `add_gates`, `remove_gates`, `add_scripts`), as modificações pontuais sobrepõem os slots canônicos especificamente para aquela instância.

---

🛑 **Stop-Gate 01:** Concluído e Aprovado pelo Usuário. Autorizado avanço para a Tríade de Planejamento (SPEC + PLAN + ROUTE).
