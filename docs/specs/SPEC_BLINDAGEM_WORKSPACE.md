# 📋 SPEC: Blindagem Canônica do Workspace e Harness Engineering

> **Referência:** `docs/specs/INTENT_BLINDAGEM_WORKSPACE.md`  
> **Status:** Especificação com Critérios Binários (Pass/Fail)  
> **Metodologia:** SDD (Spec-Driven Development) & Harness Engineering  

---

## 1. Escopo Técnico & Arquitetura da Blindagem

A blindagem opera em duas frentes complementares:

1. **Camada de Grafo de Conhecimento (Graph-First):**
   * Servidor: `codebase-memory-mcp`
   * Projeto: `LegoLoop`
   * Caminho raiz: `C:\Users\trcnologia\Desktop\LegoLoop`
   * Ação: Criação e indexação completa da AST de Go, HTML, JS e CSS do repositório.

2. **Camada de Governança de Agente (`AGENTS.md`):**
   * **Lei Graph-First (Inviolável):** Proibido usar `read_file` em arquivos inteiros (> 100 linhas) ou comandos cegos de busca sem antes consultar o grafo semântico.
   * **Rule 10 (Silent Executor & Token Cap):**
     * Resposta padrão: $\le 300$ tokens.
     * Resposta técnica/arquitetural: $\le 600$ tokens.
     * Proibição de repetição de código já gravado no disco ou resumo prolixo de diffs.
   * **Proibição Estrita de Polling:** Banido o uso de `manage_task(status)` em loop de espera. O harness deve aguardar passivamente a notificação de *Reactive Wakeup*.
   * **Correção Canônica de Caminhos:** Atualização de ponteiros para `docs/specs/SPEC.md` e `docs/plans/PLAN.md`.

---

## 2. Critérios Binários de Aceitação (Pass/Fail)

### Bloco A: Grafo Semântico (`codebase-memory-mcp`)
- [ ] **CA-GRAFO-01:** O comando de indexação no `codebase-memory-mcp` para o projeto `LegoLoop` encerra com status de sucesso (`status == "indexed"` ou `"complete"`).
- [ ] **CA-GRAFO-02:** A ferramenta `index_status` ou `detect_changes` confirma cobertura do repositório com arquivos mapeados.
- [ ] **CA-GRAFO-03:** Uma consulta semântica de teste (`search_code` ou `get_code_snippet` para o struct `Orchestrix`) retorna o trecho exato de código sem ler o arquivo inteiro via `read_file`.

### Bloco B: Governança do Agente (`AGENTS.md`)
- [ ] **CA-GOV-01:** O arquivo `AGENTS.md` contém a seção **Rule 10 (Silent Executor)** com os tetos expressos de 300 e 600 tokens.
- [ ] **CA-GOV-02:** O arquivo `AGENTS.md` contém a proibição explícita de polling em background tasks (`manage_task`).
- [ ] **CA-GOV-03:** O arquivo `AGENTS.md` contém a **Lei Graph-First** exigindo consulta ao grafo antes de varreduras amplas.
- [ ] **CA-GOV-04:** Os links operacionais apontam para `docs/specs/SPEC.md` e `docs/plans/PLAN.md`, sem links quebrados na raiz.

### Bloco C: Integridade e Não-Regressão
- [ ] **CA-INT-01:** O comando `go test ./...` executa e encerra com código `0` e todos os pacotes aprovados.
- [ ] **CA-INT-02:** `git status -s` reporta apenas arquivos modificados/adicionados rastreados, sem artefatos temporários pendentes.

---

## 3. Protocolo de Transição (Gate para o Planejamento)

A esteira só avança para a confecção do plano de micro-tickets (`PLAN_BLINDAGEM_WORKSPACE.md`) após a aprovação formal desta especificação.
