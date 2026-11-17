# 🎯 INTENT: Refatoração Completa da Lei #15 (ASD-STE100) e Orchestrix sobre Claude Code

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-10  
> **Autor:** Claude (a pedido do usuário)  
> **Alvo:** `internal/`, `cmd/`, `loops/`, `skills/`, `docs/templates/`, `schemas/`, `tools/taxlinter`  
> **Sessão:** 9520cb21-e300-42c1-83f2-c6fab90c3188  
> **Base das decisões:** decisões D1 a D8 confirmadas pelo usuário na sessão de 2026-11-10  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação

A Lei #15 (`AGENTS.md` item 12) proíbe "loop", "pipeline", "esteira", "fluxo", "bloco" e outros sinônimos. O `taxlinter` reporta 412 violações (piso, medido em 2026-11-10). Os geradores (`fabricar`, `compor`), os templates, os prompts do LLM, a CLI e as skills ainda emitem os termos proibidos. Cada Processo novo cria violações novas.

O usuário também quer que o Orchestrix execute os tickets sobre o Claude Code, com subagentes efêmeros, roteamento por harness e modelo, paralelismo e execução assíncrona.

## 2. Escopo Inegociável

### O que ESTÁ no Escopo:
* [ ] Normalizar código, schema, nomes de arquivo, rotas, templates, prompts, CLI, skills e documentação para a Lei #15.
* [ ] Migrar os contratos em disco com leitor retrocompatível. Os loops existentes (15 em `loops/` e 6 no Cofre) continuam válidos.
* [ ] Implementar o Orchestrix sobre Claude Code: executor de subagente efêmero, roteador, raias do `agy`, modelo por natureza, paralelismo e telemetria.
* [ ] Ampliar o `taxlinter` para medir o valor real das violações.

### O que NÃO ESTÁ no Escopo:
* Renomear o projeto. A marca continua **LegoLoop**.
* Alterar a VPS do 9Router (outro agente cuida).
* Renomear Processos que são dado do usuário (`mega-pipeline`, `nova-esteira`).
* Remover o leitor dos nomes antigos.
* Commit e push sem pedido do usuário.

## 3. Mapeamento Preliminar de Peças e Processos

* **Peças existentes reaproveitadas:** `internal/dag`, `internal/gate`, `internal/runner`, `internal/orchestrix`, `internal/plangen`, `internal/fechar`, `internal/processo`, `tools/taxlinter`, `tools/funclen`.
* **Peças novas:** `internal/golden` (teste dourado), `internal/migrar` (`legoloop migrar`), `internal/studiodoc` (renome de `internal/loop`), `tools/ticketgate` (bateria única de verificação), caminho único de resolução de nomes.
* **Processo executor:** este plano roda como DAG de 34 tickets e 5 Stop-Gates (`docs/routes/2026-11-10_ROUTE_NORMALIZACAO_TAXONOMICA_LEI15.json`).

## 4. Governança Operacional e Restrições Inegociáveis

* Subagente efêmero permitido sob seis condições (contexto só dos Inputs, ferramentas listadas, um ticket por subagente, saída em disco mais Gate, sem sub-subagentes, Write Boundary `strict`). Fork e workflow livre continuam proibidos.
* Ordem de roteamento: 9Router, Antigravity, Mimocode, OpenCode, Claude Code, Ollama. O ticket sobrescreve com justificativa.
* `agy` em raias fixas: conta da própria sessão, ticket fixo na raia, 429 leva ao próximo harness, alarme "Terms of Service" desliga todas as raias.
* Modelo por natureza: `nano` Haiku, `micro` e `macro` Sonnet, Opus só com a marca `extreme` e Stop-Gate.
* Nenhum ticket toca o Cofre real sem Stop-Gate e backup.
* Stop-Gates são humanos. Nenhum agente os aprova.

## 5. Critérios de Sucesso de Alto Nível

* `go run ./tools/taxlinter internal cmd loops skills docs/templates` retorna exit 0 com o linter ampliado.
* Os 15 loops de `loops/` e os 6 do Cofre validam antes e depois de cada ticket.
* Um Processo novo gerado pelo `fabricar` não contém termo proibido.
* Um Processo de duas etapas roda de ponta a ponta com o binário real e LLM real, com roteamento e telemetria no REPORT.

## 6. Handoff de Governança

* Consome: decisões D1 a D8. Produz: SPEC, PLAN e ROUTE.
* Quinteto: INTENT (este arquivo), SPEC, PLAN, ROUTE. O DAG espelho de `docs/dags/` é gerado depois da aprovação da ROUTE.

---

## 7. Registro de Execução (2026-11-12)

* O Épico 03 terminou em 2026-11-11 e o `SG-E03` foi aprovado em 2026-11-12. Os Épicos 00 a 02 estão implementados. Os Épicos 04 a 06 não começaram.
* A execução diferiu da intenção em quatro pontos: roteamento de modelo só sugerido, `--parallel 1`, três etapas por agente com shell e gate do E03 com limite de violações.
* O quinteto continua vigente. Falta a entrevista: não existe entrevista da Lei #15 no repositório.
* Detalhes, commits e custo: [REPORT de execução](../reports/2026-11-12_REPORT_EXECUCAO_LEI15.md).
