---
id: INTERVIEW-I3
status: DRAFT
date: 2026-11-12
schema_version: 2
upstream: [INTERVIEW-I2]
covers: [Q1, Q2, Q3, Q4, Q5, Q6, Q7]
---

# INTERVIEW: Iniciativa I3, Roteamento honesto

> **Protocolo:** uma pergunta por vez, em prosa. A entrevista só termina com confiança igual ou acima de 95%.
> **Regra de prova:** nenhuma pergunta nasce de inferência. Cada pergunta cita fatos F da seção 3, e cada fato F traz a prova (arquivo e linha, ou comando e saída). Sem prova, o entrevistador verifica no código antes de perguntar. Fato de memória ou de sessão anterior não conta como prova.
> **Confiança atual:** 83%
> **Iniciativa:** I3 Roteamento honesto
> **Template:** TEMPLATE_INTERVIEW.md v2.3.0

## 1. Respostas

| ID | Pergunta | Fatos | Decisão do usuário |
| :--- | :--- | :--- | :--- |
| Q1 | Qual problema a I3 resolve primeiro? | F1, F2 | A honestidade da ROUTE é a prioridade 1. O custo e a resiliência vêm depois. |
| Q2 | O que o motor faz quando o harness do DAG falha na execução? | F2, F5 | Depende da classe. Nós `nano` e `micro` degradam com registro. Nós de classe grande e `extreme` falham. |
| Q3 | De onde vem a lista de modelos e esforços de cada harness? | F5, F6, F7, F8 | Arquivo do usuário com validação pelo inventário. O 9Router entra como `via`, não como harness. |
| Q4 | Quem comuta o gateway quando o nó pede `claude` com `via: 9router`? | F9 | O motor, por nó, só no processo filho. O motor não toca o marcador global do `gw`. |
| Q5 | O que o `call_skill` devolve? | F10 | O texto da skill. O nó do DAG executa. |
| Q6 | Quem define quantos nós rodam ao mesmo tempo? | F11 | O motor, por medição de memória. O peso por harness vem da bateria de conformidade. O `--parallel` é o teto. |
| Q7 | O que fazer com as raias do agy (TICK-05.3)? | F12 | Remover. Guardar o alarme de banimento como "harness reprovado" no inventário. |

Decisões sem pergunta numerada, tomadas pelo usuário durante a entrevista:

* Admissão por prova: o harness só roda nó de código depois de passar na bateria de conformidade. Isso vale para harness proprietário ou aberto.
* A bateria roda na instalação e quando aparece um harness novo no sistema operacional. O resultado vale por par: harness com `via`.
* O gerador do roteamento (passo 2 da cascata) fica pendente da pergunta Q8.

## 2. Confiança por rodada

| Rodada | Perguntas feitas | Confiança | Motivo da variação |
| :--- | :--- | :--- | :--- |
| 1 a 3 | Q1, Q2 e a Q3 | 40% | O problema, a falha em execução e a origem dos modelos. |
| 4 a 6 | Q3 refinada, Q4, admissão por prova | 60% | O 9Router é `via` e não harness. A bateria entra no plano. |
| 7 a 9 | Q5, Q6, Q7 | 83% | Skill sob demanda, paralelismo e fim das raias do agy. |
| 10 | Q8 | aberta | Refeita com prova depois da correção do usuário. |

## 3. Fatos medidos durante a entrevista

| ID | Afirmação | Prova (arquivo:linha ou comando) | Resultado |
| :--- | :--- | :--- | :--- |
| F1 | O nó de código traduz `code-fast` e `code-pro` para `sonnet` e `flash` para `haiku`. | `internal/orchestrix/unit_run.go:321` | VERIFICADO |
| F2 | O `dispatch` do motor tem o caso `SourceRouter` (chamada HTTP direta ao 9Router) e um caso padrão para outros harnesses. | `internal/orchestrix/orchestrix.go:106` e `:108` | VERIFICADO |
| F3 | A chamada direta devolve só texto: `callRouter` lê `choices[0].message.content`. | `internal/orchestrix/backends.go:29` a `:50` | VERIFICADO |
| F4 | O omp roda com `--no-tools` e o opencode roda com `--pure` no motor. | `internal/orchestrix/harness.go:105` e `:114` | VERIFICADO |
| F5 | O caminho do claude no nó de código fecha a lista de ferramentas, define o esforço e limita o gasto com `--max-budget-usd`. | `internal/orchestrix/unit_claude.go:43` | VERIFICADO |
| F6 | `LerInventario` recebe os modelos de fora (`modelos map[string][]string`). Nada no código os fornece. | `internal/roteamento/inventario.go:32` | VERIFICADO |
| F7 | Nenhum arquivo importa `legoloop/internal/roteamento`. | `grep -rln '"legoloop/internal/roteamento"' --include=*.go .` sem resultado | VERIFICADO |
| F8 | O mimo e o omp têm o 9Router como provider configurado. O opencode também. | `~/.config/mimocode/mimocode.jsonc:32`, `~/.omp/agent/models.yml` (provider `aidd9r`), `~/.config/opencode/opencode.json` (provider `9router`) | VERIFICADO |
| F9 | O 9Router muda o Claude Code por variáveis `ANTHROPIC_BASE_URL` e token só no processo filho. | `gateway-claude/REGISTRO.md:18` | VERIFICADO |
| F10 | Existe o roteador de MCP (`mcp_search`). Não existe `skill_search` nem `call_skill` em Go. | `internal/mcp/router.go:140`; `grep -rn "skill_search\|call_skill" --include=*.go .` sem resultado | VERIFICADO |
| F11 | O DAG da Lei #15 permite 2 workers. | `docs/dags/2026-11-10_DAG_NORMALIZACAO_TAXONOMICA_LEI15.json:18` | VERIFICADO |
| F12 | O PLAN define `TICK-05.3` como raias do `agy`, com desligamento global ao ver "Terms of Service". | `docs/plans/2026-11-10_PLAN_NORMALIZACAO_TAXONOMICA_LEI15.md:175` | VERIFICADO |
| F13 | O gerador `gen.mjs` não existe no repositório e nunca entrou no histórico do git. | `find . -name gen.mjs` sem resultado; `git log --all -- '*gen.mjs'` sem resultado | VERIFICADO |
| F14 | A `tools/route2dag` só projeta a ROUTE em DAG. Ela não escolhe harness, modelo nem esforço. | `tools/route2dag/main.go:1` a `:8` | VERIFICADO |
| F15 | O instalador atual é `go install`. Existem os comandos manuais `doctor` e `instalar-skills`. Não existe bateria de conformidade de harness. | `README.md:89`, `:92` e `:104` | VERIFICADO |
| F16 | O mimo usa cerca de 750 MB por processo, e dois processos do agy derrubam a máquina. | Registro de sessão anterior (memória). Não reproduzido nesta sessão. | NÃO REPRODUZIDO |
| F17 | O Claude Code exige comutação para usar os combos do 9Router. | Ver F9. Falta provar que os outros harnesses rodam o nó de código com o 9Router. A bateria prova isso. | PARCIAL |

## 4. Pendências para a SPEC

* Q8: quem preenche o roteamento da ROUTE (passo 2 da cascata). Opções: comando novo em Go, agente com validador, ou dentro da `route2dag`.
* Medir o peso de memória de cada harness (F16) dentro da bateria de conformidade.
* Provar na bateria quais harnesses rodam o nó de código com o 9Router (F17).
* Definir o formato do arquivo de perfil do usuário e o campo `via` do `Candidato`.
* Unir o `provider.Inventory` do motor ao `internal/roteamento` (F7).
* Criar o `skill-router` (F10) sobre o mesmo cofre do roteador de MCP.
* Absorver o E05: manter 05.1, 05.5 e 05.6 adaptados. Substituir 05.2 e 05.4 pelas decisões Q2 e pela classe por regra. Remover 05.3.
* Criar a bateria de conformidade e o comando de homologação de harness, com teto de gasto e aviso antes de rodar.
