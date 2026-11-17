# 🎯 INTENT: Refatoração da Skill `aidd-legoloop-mega-pipeline`

> **Status:** Aprovado pelo Usuário  
> **Data de Aprovação:** 2026-11-01  
> **Autor:** Antigravity AI + Usuário  
> **Alvo:** `C:\Users\trcnologia\.gemini\config\skills\aidd-legoloop-mega-pipeline\SKILL.md`

---

## 1. Contexto e Motivação

Durante a operação prática do agente e da esteira, foi constatada uma vulnerabilidade comportamental crítica:
* **Afobação Cognitiva do Agente:** O agente pulou fases essenciais (concepção, entrevista socrática, geração de intenção formal, rotas explícitas) e executou ações de código diretamente na sessão.
* **Falta de Guardrails Conversacionais na Skill:** A skill atual descreve a esteira teórica de 5 sub-loops, mas **não impõe travas mecânicas no comportamento conversacional do agente**. O agente não possuía regras invioláveis de parada obrigatória (*Stop-Gates Interativos*) entre os loops.
* **Falta de Granularidade por Épico e Auditoria Humana:** Ausência de protocolo explícito que exija que o plano seja decomposto em Épicos, cada um com seus micro/nano-tickets roteados por modelo cognitivo, com auditoria humana mandatória entre os épicos antes da liberação do próximo bloco de build.

---

## 2. Objetivos Principais

1. **Blindagem contra Afobação Cognitiva:**
   * Transformar o `SKILL.md` em um manual com **Stop-Gates Humanos Inegociáveis** (o agente é fisicamente proibido pelo prompt da skill de gerar código sem aprovação explícita do `INTENT`, da `SPEC` e do `PLANO`).
2. **Decomposição Canônica em Épicos & Micro/Nano-Tickets:**
   * Formalizar o modelo de planejamento hierárquico:
     * **Épico (Macro):** Visão arquitetural e fronteiras de domínio.
     * **Micro-ticket:** Unidade atômica de entrega (1 arquivo/módulo + teste TDD).
     * **Nano-ticket:** Ajustes cirúrgicos (< 4k tokens: linters, asserts, tipagens, tokens CSS).
   * Roteamento cognitivo explícito de cada ticket para o modelo ideal (Sonnet, Pro, Flash, Flash-Lite, Ollama).
3. **Protocolo de Build por Épico com Gate de Auditoria Humana:**
   * Execução estrita bloco a bloco.
   * Conclusão de um Épico ➔ Parada imediata ➔ Apresentação de sumário factual ➔ Aguardo da validação humana ➔ Próximo Épico.
4. **Governança Estéril e Clamping:**
   * Integração das lições de Harness Engineering: diretrizes para saídas estruturadas, contenção de prolixidade e respeito aos tokens de raciocínio (*effort clamping*).

---

## 3. Escopo da Refatoração da Skill

### O que entra:
* Reestruturação completa de [`SKILL.md`](file:///C:/Users/trcnologia/.gemini/config/skills/aidd-legoloop-mega-pipeline/SKILL.md).
* Definição explícita dos 5 Sub-Loops com **Gates de Entrada, Execução e Saída Humana**.
* Tabela de roteamento de modelos com indicação de esforço (`--effort low|medium|high`) e tipo de tarefa.
* Templates padronizados de artefatos: `INTENT.md`, `SPEC.md`, `PLAN.md` (com épicos e tickets), `REPORT.md`.
* Protocolo operacional estrito para o chat (o que o agente PODE e NÃO PODE fazer em cada etapa).

### O que NÃO entra nesta etapa:
* Modificações no código Go do motor LegoLoop ([`internal/`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal)). Esta etapa ocorrerá na fase seguinte, sendo orquestrada pela própria skill já refatorada.

---

## 4. Questões para Validação com o Usuário

Para fecharmos este `INTENT` e avançarmos para a `SPEC`:

1. **Formato dos Épicos e Micro/Nano-Tickets:** Deseja que a skill mantenha os tickets embutidos no próprio `PLAN.md` estruturado em markdown ou prefere a geração de arquivos individuais na pasta `docs/tickets/TICK-*.md`?
2. **Roteamento de Modelos Padrão:** Na sua máquina, confirmamos o uso prioritário de `Claude Sonnet` / `Gemini Pro` para Épicos e `Gemini Flash` / `Flash-Lite` para Micro/Nano-tickets?
3. **Critério de Parada:** Concorda que, se o agente tentar gerar código antes do `PLAN.md` aprovado, isso deve ser tratado como quebra grave de protocolo na skill?
