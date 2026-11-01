# 🎯 INTENT: Integração do Super-Board (Board-Driven Looping & Guard Rails)

> **Data:** 01/11/2026  
> **Status:** Aprovada (Alinhada com Usuário)  
> **Guarda-Chuva:** Looping Engineering & External Orchestration  
> **Referência Externa:** `https://github.com/EricTechPro/super-board.git`  

---

## 1. Motivação e Oportunidade Observada

O **LegoLoop** possui um motor Go de altíssima performance para execução estéril de DAGs locais (`internal/dag`, `internal/runner/isolated.go`). No entanto:
1. **Falta de Ingestão Externa Reativa:** O acionamento hoje é local (CLI `legoloop run` ou Studio UI). Projetos reais se beneficiam de uma esteira que consome diretamente itens de um quadro Kanban externo (**GitHub Projects v2**), drenando cards sem intervenção manual.
2. **Concorrência Cega a Arquivos:** Ao disparar tarefas paralelas, não há no LegoLoop atual um algoritmo que inspecione antecipadamente o risco de conflito de merge (`file overlap`) entre duas tarefas simultâneas.
3. **Ciclo de Vida de PR Ausente no Worktree:** O isolamento de worktree atual do LegoLoop é puramente efêmero e pontual por nó (`HEAD --detach`). Falta estender esse mecanismo para gerenciar o ciclo de vida completo de uma issue (`Branch ➔ Commits ➔ Draft PR ➔ QA Pass ➔ Merge Gate`).
4. **Vulnerabilidade a Rate Limits e Comandos Destrutivos:** A ausência de *guard hooks* ativos pode permitir que modelos de IA executem exclusões indevidas fora do workspace ou estourem a quota de 5.000 req/h da GitHub API.

---

## 2. Intenção Primária

Incorporar ao ecossistema do **LegoLoop** os mecanismos de orquestração externa, concorrência segura em ondas (`Wave Planning`), ciclo de vida de issues em worktrees e blindagem de segurança (*Guard Hooks* e *Rate Guard*) extraídos da engenharia do repositório `super-board`, mantendo a execução 100% estéril e em Go.

---

## 3. Escopo e Entregáveis

1. **Ingestor de GitHub Projects v2 (`internal/board`):**
   * Leitura de cards na coluna `Ready` e transição de estado (`Building ➔ QA ➔ Review ➔ Done`).
   * Proteção de cota com *Token Bucket* (Exit 79 gracioso em caso de esgotamento de quota).
2. **Planejador de Ondas Anti-Conflito (`internal/wave`):**
   * Extração de arquivos/pastas tocados por cada issue.
   * Particionamento em ondas concorrentes mutuamente disjuntas (zero `file overlap`).
3. **Gerenciador de Ciclo de Vida de Issue (`internal/runner/lifecycle.go`):**
   * Provisionamento de worktree com branch rastreável por issue.
   * Criação automatizada de Draft PR e acionamento de Merge Gate determinístico.
4. **Acervo de Guard Hooks e Peças Canônicas:**
   * Canonização dos hooks `guard-delete-outside`, `guard-protected-push`, `guard-secrets`, `cleanup-wt`.
   * Registro do bloco `Trigger: GitHub Project Kanban` no catálogo de peças do Studio.

---

## 4. Fora de Escopo

* Não clonar o runtime em Node.js/Bash do `super-board` como dependência externa pesada. O core do LegoLoop permanece Go compilado nativo.
* Não criar subagentes de chat que queimam tokens; toda execução segue a Orchestrix Estéril.

---

## 5. Critérios de Sucesso (Definition of Done)

- [ ] Intent formalmente aprovada e referenciada pela Spec Master do épico.
- [ ] Especificação técnica (`docs/specs/2026-11-01_SPEC_SUPER_BOARD_INTEGRATION.md`) validada com critérios de aceite binários.
- [ ] Plano de tickets (`docs/plans/2026-11-01_PLAN_SUPER_BOARD_INTEGRATION.md`) decomposto em fatiamento vertical com TDD estrito.
- [ ] Todos os testes em Go verdes (`go test ./...`) e recompilação do binário finalizada.
