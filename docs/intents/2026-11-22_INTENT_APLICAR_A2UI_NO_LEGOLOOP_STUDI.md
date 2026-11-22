---
id: INTENT-intv_1795370552989667600
status: APPROVED
date: 2026-11-22
schema_version: 2
upstream: [INTERVIEW-intv_1795370552989667600]
covers: [Q1]
---

# 🎯 INTENT: Motor Único de Interface Declarativa A2UI nas 3 Vias do LegoLoop

> **Status:** APPROVED  
> **Data:** 2026-11-22  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** `internal/studio`, `internal/runner`, `internal/interview`, `internal/mcp`, `cmd/legoloop`  
> **Sessão:** intv_1795370552989667600  
> **Confiança da Entrevista:** 98% (mínimo de 95% para sair de DRAFT; Espelho Canônico reconciliado em `docs/intents/ESPELHO_CANONICO.md`)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Hoje o LegoLoop exige interações em texto linear na via Terminal (CLI) e na via Agente de IA / Slash Command, enquanto o botão Executar do Studio (`#btn-run` -> `POST /api/run`) abre apenas um modal de log SSE textual e o Swarm mantém uma interface web separada de monitoramento.
* **Quem Sofre:** O usuário e operadores do LegoLoop em Windows, Linux e macOS, usando qualquer harness (Antigravity, Claude Code, OpenCode, MimoCode ou CLI pura).
* **O Porquê:** Construir um motor único de interface declarativa baseado no protocolo aberto A2UI (`@a2ui/lit` embutido no binário Go via `go:embed`) faz com que todos os Processos do LegoLoop herdem automaticamente uma interface gráfica rica nas 3 vias, unificando interação e monitoramento e eliminando a UI separada do Swarm na Fase 2.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Três experiências fragmentadas — texto linear em CLI/Chat, Canvas + modal de log SSE no Studio, e UI separada de monitoramento do Swarm.
  * **Amanhã:** Apenas duas interfaces complementares — **Studio Canvas** (construção de Processos/DAGs) e **Interface Interativa A2UI** (execução, interação e monitoramento) acessível por uma URL única de sessão nas 3 vias, com Design System do LegoLoop, suporte a macroprocessos via `surfaceId` e fallback textual garantido.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1 (Fase 1 — Motor Único A2UI nas 3 Vias): Criar sessão interativa centralizada no servidor local do LegoLoop com SSE de saída (`GET /api/sessao/{id}/stream`), entrada de ações (`POST /api/sessao/{id}/action`), corrida `stdin` vs. navegador na CLI, espera estruturada via MCP no Chat/Slash Command e abertura direta pelo botão Executar no Studio.
* [ ] Objetivo 2 (Fase 1 — Catálogo LegoLoop, Design System e Macroprocessos): Integrar `@a2ui/lit` embutido via `go:embed` com Design System próprio do LegoLoop, combinando catálogo básico A2UI e componentes de domínio (DAG, progresso/confiança, cartão Stop-Gate, visualizador de artefatos e logs), com isolamento de macroprocessos por `surfaceId` (árvore lateral + breadcrumb + foco automático em Stop-Gate).
* [ ] Objetivo 3 (Fase 1 — Projeto Piloto `legoloop-concepcao`): Validar o motor ponta a ponta no Processo `legoloop-concepcao` nas 3 vias (perguntas dinâmicas, barra de confiança, Espelho Canônico editável, Stop-Gate 01 clicável e preview de `INTENT.md`) com fallback textual íntegro.
* [ ] Objetivo 4 (Fase 2 — Unificação e Remoção da UI Antiga do Swarm): Plugar o monitoramento do Swarm (`native_swarm` / `swarm_api.go`) como superfície A2UI na mesma interface interativa, provar paridade funcional por testes e remover a UI web antiga e separada do Swarm (preservando o motor/backend do Swarm).

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON ou o formato do DAG sem spec formal prévia e retrocompatibilidade.
* ⛔ Não criaremos protocolo proprietário de UI fora do padrão aberto A2UI.
* ⛔ Não reescreveremos o Studio Canvas inteiro em React, Angular ou Flutter.
* ⛔ Não removeremos o motor/backend do Swarm (`legoloop swarm`, quadro de mensagens, workers), apenas sua UI web antiga após paridade provada na Fase 2.
* ⛔ Não integraremos dashboards externos fora do LegoLoop.

---

## 3. Mapeamento Preliminar de Peças
Identificação das Skills que o futuro Processo precisará acionar. Tipos permitidos: os tipos de nó do DAG e do ROUTE, descritos no README, seção 6.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | `a2ui` (Cofre canônico) | Referência oficial do protocolo A2UI e renderizador `@a2ui/lit`. |
| `PECA-02` | `skill` | `skills/aidd-spec` | Elaborar a especificação formal com critérios BDD (Fase 1 e Fase 2). |
| `PECA-03` | `unit_no` | `internal/studio`, `internal/runner`, `internal/interview`, `internal/mcp` | Sessão interativa SSE + `POST /action`, corrida `stdin`/UI, MCP wait, catálogo LegoLoop e piloto `legoloop-concepcao`. |
| `PECA-04` | `unit_no` | `internal/studio/static/` (`@a2ui/lit` + `go:embed`) | Cockpit A2UI, Design System LegoLoop, árvore/breadcrumb de macroprocessos e superfície Swarm (Fase 2). |
| `PECA-05` | `script` | `go test -count=1 ./...` | Regressão mecânica, hardening de segurança e paridade antes da remoção da UI antiga do Swarm. |
| `PECA-06` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (qualquer invocação de Processo nas 3 vias).
* **Nível de Risco & Blast Radius:** `Médio` (altera Studio, runner de sessão, entrevista e MCP; exige hardening contra injeção e preservação integral do fallback textual).
* **Política de Privacidade & Orçamento Cognitivo:** `Local` (servidor em `127.0.0.1`, assets `@a2ui/lit` embutidos via `go:embed` sem dependência de CDN/internet em runtime; nenhum segredo trafega em payloads A2UI).
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  * Piloto `legoloop-concepcao` executado de ponta a ponta nas 3 vias com interação real no navegador e fallback textual comprovado.
  * Macroprocesso com `subprocesso` navegável por árvore/breadcrumb e isolado por `surfaceId`.
  * Na Fase 2, paridade de monitoramento do Swarm provada por testes antes da remoção da UI antiga.
* **Restrições Inegociáveis:**
  * 1. Zero HTML ou JavaScript executável vindo de agentes (somente JSON A2UI validado sobre catálogo confiável).
  * 2. Fallback textual garantido nas 3 vias (terminal e chat continuam funcionando sem abrir navegador).
  * 3. Isolamento por bloco e por `surfaceId` em macroprocessos (Lei 4, contexto estéril).
  * 4. Design System próprio do LegoLoop no host (`@a2ui/lit` + `go:embed`, sem React).
  * 5. TDD estrito (Red-Green-Refactor) e governança viva (Leis 5, 8, 12, 16 e 18).
* **Premissas Assumidas (confirmadas no Espelho Canônico):**
  * Piloto obrigatório em `legoloop-concepcao` antes de generalizar para a Fase 2.
  * Unificação das interfaces do LegoLoop para 2 superfícies (Studio Canvas + Interface Interativa A2UI), removendo a UI antiga do Swarm somente na Fase 2 após paridade provada.
* **Riscos e Perguntas em Aberto:**
  * Empacotamento offline do bundle `@a2ui/lit` + `@a2ui/web_core` dentro do `go:embed` (a detalhar na SPEC).
  * Precedência de execução frente às iniciativas I01, I02 e I03 (já liberadas pós-I04).

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test -count=1 ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Testes novos falhando antes do código | Relatório de testes e `tools/cobertura` |
| 3 | Zero HTML/JS executável aceito em payloads A2UI. | Testes de hardening rejeitando payload malicioso | `go test -count=1 -run Security ./internal/studio/...` |
| 4 | Piloto `legoloop-concepcao` funciona nas 3 vias com corrida `stdin`/UI e fallback textual. | Teste E2E automatizado + execução real | `go test -count=1 ./internal/studio/... ./internal/interview/...` |
| 5 | Fase 2: paridade do Swarm em A2UI comprovada antes da remoção da UI antiga. | Testes de contrato SSE e ações do Swarm verdes | `go test -count=1 -run Swarm ./internal/studio/...` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md`](../specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md)
* **Gate de Transição:** Esta INTENT só é considerada aprovada após a validação humana explícita (*Stop-Gate 01*) no chat, mudando `status: UNDER_REVIEW` para `status: APPROVED` e gravando `docs/intents/INTENT_APROVADA.md`, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 98% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se mudar para `APPROVED` e gerar a SPEC (`TEMPLATE_SPEC.md`).*
