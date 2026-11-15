---
id: INTENT-HOMOLOGACAO-ECOSISTEMA
status: DRAFT
date: 2026-11-14
schema_version: 2
upstream: [INTERVIEW-HOMOLOGACAO-ECOSISTEMA]
covers: [Q1, Q2, Q3, Q4]
---

# 🎯 INTENT: Auditoria, Correção Integral e Homologação do Ecossistema LegoLoop

> **Status:** DRAFT (Aguardando Stop-Gate Único Humano)  
> **Data:** 2026-11-14  
> **Autor:** Antigravity Coordinator  
> **Iniciativa:** I-A / I-B / I-C — Homologação Integral do Ecossistema LegoLoop  
> **Alvo:** `cmd/legoloop`, `internal/`, `tools/`, `schemas/`, diretório de Processos, `skills/`, `docs/`  
> **Sessão:** `7fb7b443-3368-4296-87f4-ca1316139cab`  
> **Confiança da Entrevista:** 100% (4 decisões arquiteturais respondidas pelo patrocinador)  
> **Template:** TEMPLATE_INTENT.md v2.0.0

---

## 1. Contexto e Motivação
* **O Problema Real:** A auditoria profunda da Etapa 1 revelou falhas críticas nas três frentes do LegoLoop, resultando em média global `4.6/10`:
  * Frente 1 (`4.8/10`): suíte `go test -count=1 ./...` quebrando em 3 pacotes, `taxlinter` com 58 violações, stubs em `internal/motion` e `internal/videorea`, data race no SSE do Studio, brechas de Path Jail (`reveal.go`, `piece_actions.go`, `security.go`), dessincronização `schemas/dag.v2.json` vs `internal/dag`, arquivos com segredos trackeados no git e 9 motores N10 sem chamadores em produção.
  * Frente 2 (`5.2/10`): `mcp-router` iniciado com `NewRouter(nil, nil)` vazio, `fechar` com 4 etapas em stub e telemetria fictícia (`80000`/`65000` tokens), divergências CLI (`usage` vs handlers, `--json` contaminado, `stept status` escrevendo em disco), paridade quebrada de `skills/embed.go` e `sync-harnesses`, e `board` com contratos divergentes.
  * Frente 3 (`3.8/10`): XSS executável via `innerHTML` no Swarm UI, Swarm UI do Studio desconectado dos boards em disco, `#btn-run` quebrado (`compileDAG is not defined`), `+ Novo Processo` retornando HTTP 500 enganoso sem matriz, clipping responsivo em `1440`/`1024`/`390px`, inspetor recolhido irrecuperável, falhas WCAG AA e inconsistência da Peça Gatilho.
* **Quem Sofre:** Desenvolvedores, agentes nos harnesses e operadores do Studio ou CLI que executam Processos, validam DAGs ou coordenam Swarms.
* **O Porquê:** Elevar todas as frentes para nota mínima `>= 9.5/10` (média global `>= 9.5/10`) com zero stubs, segurança reforçada e prova real ponta a ponta antes de qualquer entrega.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Nota global `4.6/10`, testes falhando, `taxlinter` reprovando, XSS e brechas de Path Jail ativas, motores N10 desconectados e Studio incapaz de executar um Processo pela interface.
  * **Amanhã:** Nota `>= 9.5/10` em todas as frentes, suíte 100% verde, `taxlinter` em `0` nos alvos oficiais, 9 motores N10 integrados, executor MCP real, Swarm UI unificado e seguro, e Studio validado em Chrome/Playwright/REA/Obscura.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Curar 100% dos achados críticos, altos, médios e baixos da Etapa 1 nas frentes Go/Arquitetura/Taxonomia, Canais/Motores e UI/UX Studio/Swarm UI.
* [ ] Objetivo 2: Integrar os 9 motores N10 em produção, implementar executor MCP STDIO real e substituir stubs de mídia (`motion` e `videorea`) por FFmpeg/Playwright reais ou erro explícito sem placeholder.
* [ ] Objetivo 3: Unificar o Swarm UI e board em `.legoloop-swarm/swarm.json` (com leitura retrocompatível), zero XSS/CSRF, file-lock cross-processo, spawn/kill real de PIDs e gestão de modelo na interface.
* [ ] Objetivo 4: Zerar violações do `taxlinter` nos alvos oficiais (`limiteTaxonomia = 0`), migrar combos legados para combos honestos certificados, limpar arquivos com segredos do git e homologar com nota `>= 9.5/10` e Fechamento Canônico Universal (Lei #11).

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não adicionaremos dependências externas novas no `go.mod` (100% stdlib Go + frontend embutido + acervo local).
* ⛔ Não alteraremos o cofre real `~/.aidd` durante testes automatizados nem executaremos `taskkill /IM legoloop.exe`.
* ⛔ Não reescreveremos documentos históricos antigos em `docs/reports/` ou `docs/sessions/` fora dos pares `ROUTE`/`DAG` divergentes e dos alvos oficiais do `taxlinter`.

---

## 3. Mapeamento Preliminar de Peças
Identificação das Peças e Skills do repositório e do sistema que o Processo orquestrará na Etapa 4, Etapa 5 e Etapa 6.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `unit_no` | `internal/dag`, `schemas/`, `internal/graph` | Saneamento de schemas, DAGs, paridade `ROUTE` e regra de Gatilho. |
| `PECA-02` | `unit_no` | `internal/runner`, `internal/orchestrix`, `internal/brain`, `internal/combo` | Integração dos 9 motores N10, Path Jail, executor MCP STDIO e resiliência HTTP. |
| `PECA-03` | `unit_no` | `internal/motion`, `internal/videorea` | Substituição de stubs por FFmpeg/Playwright reais ou erro explícito. |
| `PECA-04` | `unit_no` | `internal/board`, `internal/cli`, `internal/fechar`, `internal/mcp` | Unificação do Swarm Board, `mcp-router` funcional, CLI alinhada e `fechar` real. |
| `PECA-05` | `unit_no` | `internal/studio`, `internal/studio/static` | Segurança XSS/CSRF/Path Jail, correção E2E do Studio e design system WCAG AA. |
| `PECA-06` | `script` | `tools/taxlinter`, `tools/funclen`, `go test ./...` | Gates mecânicos contínuos por nível topológico e verificação cega (Etapa 5). |
| `PECA-07` | `skill` | `skills/legoloop-fechar`, `playwright`, `chrome-devtools-mcp`, `rea`, `obscura` | Auditoria visual E2E e Fechamento Canônico Universal (Lei #11). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (executado via Swarm nativo em `workspace="share"` por nível topológico após o Stop-Gate único).
* **Nível de Risco & Blast Radius:** `Alto (exige workspace isolado por worker no Swarm e portões mecânicos antes de integrar na branch principal)`.
* **Política de Privacidade & Orçamento Cognitivo:** `Nuvem autorizada` com teto controlado por `token_budget` nos nós do DAG. Dados sensíveis envolvidos: `.env-empresa`, `.env-full`, `.env-pessoal` e `credential_gemini_backup.json` na raiz (nunca expor, copiar em worktree temporário nem enviar para logs/respostas), além da remoção dos rastros `RELATORIO_AUDITORIA_GEMINI4.md` e `AUDITORIA_FONTE.txt` do controle de versão.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** `go test -count=1 ./...`, `go vet ./...`, `gofmt -l .`, `go run ./tools/funclen internal cmd tools` e `go run ./tools/taxlinter` nos alvos oficiais do AGENTS.md com exit code `0`, `git diff --exit-code go.mod go.sum` limpo, todos os DAGs válidos via `bin/legoloop.exe validate`, prova E2E de `mcp-router`, `fechar`, `instalar-skills`, Studio Web e Swarm UI com nota `>= 9.5/10` em cada frente.
* **Premissas Assumidas (a confirmar no Stop-Gate Único):**
  * Os arquivos `.env-*` e `credential_gemini_backup.json` da raiz permanecem protegidos no `.gitignore` local e bloqueados contra cópia em isolamentos temporários, enquanto os arquivos com segredos parciais trackeados (`RELATORIO_AUDITORIA_GEMINI4.md` e `AUDITORIA_FONTE.txt`) serão removidos do repositório.
  * Todos os executáveis externos opcionais (`ffmpeg`, `playwright`, `claude`, `agy`, `mimo`) quando ausentes no ambiente de teste produzem erro explícito tipado e acionável, nunca placeholder silencioso.
* **Riscos e Perguntas em Aberto:** Nenhuma pendência aberta após a entrevista socrática (`100%` de clareza). O único risco operacional é conflito de arquivo entre workers paralelos do Swarm, mitigado pela partição estrita de arquivos por nível topológico na `ROUTE` e no `DAG`.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério possui comprovação mecânica objetiva.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Zero falhas de suíte, análise estática, tamanho de função e formatação. | Execução após integração final | `go test -count=1 ./...`, `go vet ./...`, `gofmt -l .`, `go run ./tools/funclen internal cmd tools` (`exit 0`) |
| 2 | Zero violações taxonômicas ASD-STE100 nos alvos oficiais. | Execução com `limiteTaxonomia = 0` | `go run ./tools/taxlinter` nos alvos oficiais (`exit 0`) |
| 3 | Zero dependências externas novas no projeto Go. | Diff limpo de manifesto Go | `git diff --exit-code go.mod go.sum` (`exit 0`) |
| 4 | Nota `>= 9.5/10` verificada por auditores cegos independentes nas 3 frentes. | Relatório Etapa 5 + provas Playwright/DevTools/REA/Obscura | `docs/reports/2026-11-14_REPORT_HOMOLOGACAO_ECOSISTEMA.md` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md`](../specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md)
* **Gate de Transição:** Este Intent integra o Quinteto ISO submetido ao Stop-Gate único humano antes de qualquer modificação em código de produção na Etapa 4.

---

🛑 **STOP-GATE ÚNICO DE AUDITORIA HUMANA:**  
*A intenção consolidada com 100% de certeza factual atende à sua expectativa?*  
*A aprovação deste documento ocorre em conjunto com a SPEC, PLAN, ROUTE e DAG no Stop-Gate único humano.*
