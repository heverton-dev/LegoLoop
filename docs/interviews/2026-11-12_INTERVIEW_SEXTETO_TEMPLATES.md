---
id: INTERVIEW-I2
status: aprovada
date: 2026-11-12
schema_version: 2
upstream: []
covers: [Q1, Q2, Q3, Q4, Q5, Q6, Q7, Q8, Q9]
---

# INTERVIEW: Iniciativa I2, Sexteto e templates

> **Protocolo:** uma pergunta por vez, em prosa. Parada com confiança igual ou acima de 95%.
> **Confiança final:** 95%.
> **Programa:** I1 Taxonomia Lei #15, I2 Sexteto e templates, I3 Roteamento honesto, I4 DISSECAR, I5 MATRIZ_PAPEIS, I6 Migração dos legados.

## Respostas

| ID | Pergunta | Decisão do usuário |
| :--- | :--- | :--- |
| Q1 | Quem consome cada template? | O motor valida DAG e ROUTE (campos obrigatórios, falha na carga). O agente e o humano consomem INTERVIEW, INTENT, SPEC, PLAN, REPORT e HANDOFF, com seções fixas conferidas por teste. |
| Q2 | Qual arquivo é a fonte única do roteamento? | A ROUTE é a fonte da decisão. O DAG carrega o roteamento já resolvido (`harness`, `model`, `effort`) por nó. Um teste falha se os dois divergirem. |
| Q3 | Quem classifica o nó (nano, micro, macro, extreme)? | O gerador calcula a classe por regra mensurável e grava a regra no nó. O autor do PLAN sobrescreve com justificativa. `extreme` só por sobrescrita ou `nature` destrutiva, sempre com Stop-Gate. Ressalva: a classe é abstrata e não depende de fornecedor (ver Q4). |
| Q4 | O que fazer quando o inventário não tem o harness ou modelo pedido? | Degradar para o mais próximo disponível, com registro obrigatório da troca na ROUTE e aviso visível. `extreme` nunca degrada em silêncio: sem modelo de topo, o nó vira Stop-Gate humano. Inventário vazio: o gerador falha. O inventário do SO (harness e modelos) é um processo anterior à geração. |
| Q5 | Rastreabilidade entre os artefatos? | Validada pelo motor. IDs `Q`, `R`, `AC`. Front matter com `id`, `status`, `date`, `upstream`, `covers`. O DAG tem `spec_ref` por nó. O validador roda no gate do ticket. |
| Q6 | Quem produz REPORT e HANDOFF? | O motor gera a parte factual no fim de cada rodada. Nenhum agente edita essa parte. O HANDOFF é um recorte do REPORT. O agente escreve só as seções interpretativas, marcadas no front matter. |
| Q7 | Como o motor trata os arquivos legados? | Sem `schema_version`: modo legado com aviso. Com `schema_version: 2`: validação estrita. A migração dos legados vira a iniciativa I6, fora da I2. |
| Q8 | Escopo das peças hook e script? | Revisado pelo usuário: todos os templates ficam agnósticos em três eixos: sistema operacional, harness e modo de execução (`cli`, `agente`, `estudio`). O script ganha variantes por SO ou núcleo em Go. |
| Q9 | Padrão de carregamento das peças skill e MCP? | `load: lazy` é o padrão. `eager` exige justificativa no front matter. A peça MCP perde o `autoApprove` fixo e registra `name`, `description` e `tools` pesquisáveis no cofre (Meta-MCP Router, `internal/mcp/router.go`). Campos `execution_modes`, `os_support` e `harness_support` em todas as peças. A skill declara a classe de capacidade, não um modelo. |

## Fatos medidos durante a entrevista

1. O `TEMPLATE_DAG.json` tem `routing` por nó, sem `harness` nem `executor`. O `TEMPLATE_ROUTE.json` tem `harness` e `executor`.
2. O motor tem executores para `agy`, `omp`, `mimo` e `opencode` em `internal/orchestrix/harness.go`. O `unit_task` aceita só o claude CLI.
3. `DetectarHarnesses` (`internal/processo/vinculo.go`) detecta claude, codex, gemini, opencode, cursor e kiro. Só o vínculo de skills o usa.
4. `TEMPLATE_PIECE_HOOK.json` usa o formato do Claude Code e chama só `powershell`. `TEMPLATE_PIECE_SCRIPT.ps1` exige PowerShell 5.1.
5. `TEMPLATE_PIECE_SKILL.md` usa `suggested_model: code-pro` (alias do 9Router). `TEMPLATE_PIECE_LOOP.json` usa `model: sonnet` e `fallback: code-fast`.
6. O carregamento sob demanda de skills não foi auditado no código. A SPEC inclui a auditoria.
7. A sessão enxerga `mcp_search` e `call_mcp_tool` do roteador. O `TEMPLATE_PIECE_MCP.json` não os reflete.

## Pendências para a SPEC

* Auditar o carregamento sob demanda de skills.
* Definir o perfil de capacidade por classe (nano a extreme) e o formato do inventário do SO.
* Definir o núcleo agnóstico do script (par `.ps1` e `.sh` ou comando Go).
