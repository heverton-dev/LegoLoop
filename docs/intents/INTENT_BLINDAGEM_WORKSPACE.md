# 🎯 INTENT: Blindagem Canônica do Workspace e Harness Engineering

> **Data:** 31/10/2026  
> **Status:** Proposta de Intenção (Aguardando Aprovação para Spec)  
> **Guarda-Chuva:** Looping Engineering & Harness Governance  

---

## 1. Motivação e Problema Observado

Durante a validação prática da Orchestrix no LegoLoop, constatamos um contraste crítico:
* **No motor do LegoLoop:** O pipeline rodou de forma estéril, rápida e econômica (1.836 tokens de entrada).
* **No chat do Antigravity (Harness Interativo):** Ocorreu um sangramento de **145.000 tokens de input**, provocado por:
  1. Ausência de indexação do repositório em um grafo semântico (`codebase-memory-mcp`), forçando buscas e leituras de arquivos inteiros.
  2. Polling ativo em background tasks (`manage_task(status)` em loop) ignorando o mecanismo nativo de *Reactive Wakeup*.
  3. Ausência de limiares numéricos de saída nas respostas do chat (Rule 10 do ecossistema de referência).
  4. Caminhos de governança no `AGENTS.md` apontando para arquivos movidos na raiz.

---

## 2. Intenção Primária

Transformar o workspace atual do **LegoLoop** em um ambiente **100% blindado contra desperdício de tokens e amnésia de contexto**, transplantando com fidelidade cirúrgica as leis de contenção comprovadas no `proj--ecossistema-aidd`.

---

## 3. Escopo e Entregáveis

1. **Grafo Semântico Local (Graph-First — Lei #14):**
   * Indexar integralmente o repositório LegoLoop no `codebase-memory-mcp`.
   * Tornar obrigatória a consulta ao grafo antes de qualquer leitura integral de arquivos.
2. **Atualização da Lei de Conduta (`AGENTS.md`):**
   * Inserir a **Rule 10 (Silent Executor / Token Cap):** Respostas normais $\le 300$ tokens, técnicas $\le 600$ tokens.
   * Proibir terminantemente o polling em background tasks, exigindo espera passiva pelo *Reactive Wakeup*.
   * Corrigir links canônicos de governança para `docs/specs/SPEC.md` e `docs/plans/PLAN.md`.
3. **Validação de Não-Regressão:**
   * Garantir que `go test ./...` permaneça 100% verde com zero stubs e zero warnings.

---

## 4. Fora de Escopo

* Alterações na lógica de execução do motor Go em `internal/` (que já está estável e testado).
* Instalação ou alteração de binários no sistema operacional fora do diretório do projeto.

---

## 5. Critérios de Sucesso (Definition of Done da Intenção)

- [ ] Repositório `LegoLoop` indexado e com status ativo no `codebase-memory-mcp`.
- [ ] `AGENTS.md` atualizado com as leis de Graph-First, teto de tokens de resposta e banimento de polling.
- [ ] Handoff pronto para a confecção da **`SPEC_BLINDAGEM_WORKSPACE.md`** com critérios binários Pass/Fail.
