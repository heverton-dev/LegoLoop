---
id: INTERVIEW-I3
status: UNDER_REVIEW
date: 2026-11-12
schema_version: 2
upstream: [INTERVIEW-I2]
covers: [Q1, Q2, Q3, Q4, Q5, Q6, Q7, Q8, Q9, Q10, Q11, Q12, Q13, Q14]
---

# INTERVIEW: Iniciativa I3, Roteamento honesto

> **Protocolo:** uma pergunta por vez, em prosa. A entrevista só termina com confiança igual ou acima de 95%.
> **Regra de prova:** nenhuma pergunta nasce de inferência. Cada pergunta cita fatos F da seção 3, e cada fato F traz a prova (arquivo e linha, ou comando e saída). Sem prova, o entrevistador verifica no código antes de perguntar. Fato de memória ou de sessão anterior não conta como prova.
> **Confiança atual:** 92% (reaberta em 12/11 depois da refutação do F8; volta a 95% com as respostas das perguntas Q13 e Q14)
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
| Q8 | Quem preenche o roteamento da ROUTE? | F6, F7, F13, F14 | Comando novo em Go (`legoloop rota`): lê o PLAN, o perfil e o inventário, e grava a ROUTE. |
| Q9 | Onde vive o perfil do usuário? | F18, F19, F20 | Global (`~/.aidd/`) e por projeto, com o projeto sobrescrevendo por classe. O padrão vem embutido no binário. |
| Q10 | Como o usuário controla o gasto da bateria? | F15, F16, F21 | Aviso com o teto total estimado e confirmação antes de gastar. Sem confirmação, o harness fica "não homologado". O Claude Code já configurado nunca é bloqueado. |
| Q11 | O que o critério 2 da bateria exige? | F4, F5, F22 | Só a conferência do motor. A bateria testa que o motor detecta e reverte escrita fora do limite. |
| Q12 | Como o campo `via` entra no contrato do DAG? | F24, F25 | Obrigatório no `schema_version: 2`. Mudança de contrato público, com Stop-Gate humano `SG-I3-CONTRATO`. |

Decisões sem pergunta numerada, tomadas pelo usuário durante a entrevista:

* Admissão por prova: o harness só roda nó de código depois de passar na bateria de conformidade. Isso vale para harness proprietário ou aberto.
* A bateria roda na instalação e quando aparece um harness novo no sistema operacional. O resultado vale por par: harness com `via`.
* 
## 2. Confiança por rodada

| Rodada | Perguntas feitas | Confiança | Motivo da variação |
| :--- | :--- | :--- | :--- |
| 1 a 3 | Q1, Q2 e a Q3 | 40% | O problema, a falha em execução e a origem dos modelos. |
| 4 a 6 | Q3 refinada, Q4, admissão por prova | 60% | O 9Router é `via` e não harness. A bateria entra no plano. |
| 7 a 9 | Q5, Q6, Q7 | 83% | Skill sob demanda, paralelismo e fim das raias do agy. |
| 10 a 12 | Q8, Q9, Q10 | 90% | Gerador em Go, perfil em dois níveis e controle do gasto da bateria. Q8 refeita com prova depois da correção do usuário. |
| 13 a 14 | Q11, Q12 | 95% | Barreira de escrita e contrato `via`. Critério de parada atingido. |
| 15 | Respostas às premissas do INTENT | 92% | O usuário exigiu todos os harnesses instalados, a medição do F16 e a verificação do 9Router. A verificação refutou o F8. Perguntas Q13 e Q14 abertas. |

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
| F8 | O opencode, o mimo e o omp usam o 9Router como provider configurado, sem comutação. | O `omp` usa o 9Router em `~/.omp/agent/models.yml:184`. O `opencode.json` não tem provider `9router` (providers: cloudflare, freellmapi, mimo, omniroute). O `omniroute` do opencode e o `ai-studio` do mimo apontam para `roteador.vpsconexao.org/v1` (F27). | REFUTADO para opencode e mimo. VERIFICADO para o omp. |
| F9 | O 9Router muda o Claude Code por variáveis `ANTHROPIC_BASE_URL` e token só no processo filho. | `gateway-claude/REGISTRO.md:18` | VERIFICADO |
| F10 | Existe o roteador de MCP (`mcp_search`). Não existe `skill_search` nem `call_skill` em Go. | `internal/mcp/router.go:140`; `grep -rn "skill_search\|call_skill" --include=*.go .` sem resultado | VERIFICADO |
| F11 | O DAG da Lei #15 permite 2 workers. | `docs/dags/2026-11-10_DAG_NORMALIZACAO_TAXONOMICA_LEI15.json:18` | VERIFICADO |
| F12 | O PLAN define `TICK-05.3` como raias do `agy`, com desligamento global ao ver "Terms of Service". | `docs/plans/2026-11-10_PLAN_NORMALIZACAO_TAXONOMICA_LEI15.md:175` | VERIFICADO |
| F13 | O gerador `gen.mjs` não existe no repositório e nunca entrou no histórico do git. | `find . -name gen.mjs` sem resultado; `git log --all -- '*gen.mjs'` sem resultado | VERIFICADO |
| F14 | A `tools/route2dag` só projeta a ROUTE em DAG. Ela não escolhe harness, modelo nem esforço. | `tools/route2dag/main.go:1` a `:8` | VERIFICADO |
| F15 | O instalador atual é `go install`. Existem os comandos manuais `doctor` e `instalar-skills`. Não existe bateria de conformidade de harness. | `README.md:89`, `:92` e `:104` | VERIFICADO |
| F16 | O mimo usa cerca de 750 MB por processo, e dois processos do agy derrubam a máquina. | Registro de sessão anterior (memória). Não reproduzido nesta sessão. | NÃO REPRODUZIDO |
| F18 | O scanner de skills mescla raízes globais e raízes do projeto. | `internal/assets/skills.go:22` a `:35` | VERIFICADO |
| F19 | O código lê configuração do usuário em `~/.aidd`. | `internal/provider/env.go:16` | VERIFICADO |
| F20 | O cofre canônico fica em `~/.aidd/canonical`. | `internal/canon/linker.go:41` | VERIFICADO |
| F21 | O teto de gasto por chamada é US$ 0,80 para `micro` e `nano`, e US$ 2,50 para os demais. | `internal/orchestrix/unit_claude.go:53` a `:58` | VERIFICADO |
| F22 | O motor confere o `write_boundary` depois da execução, sem depender do harness. | `internal/orchestrix/unit.go:42`; `unit_run.go:134`; `unit_seed.go:105` | VERIFICADO |
| F24 | O schema v2 aceita só a versão 2 e exige `routing.harness`, `routing.model` e `routing.effort`. O campo `via` não existe. | `internal/dag/schema_v2.go:23`, `:28` e `:69` | VERIFICADO |
| F25 | Nenhum dos 50 arquivos de `docs/dags/` e `docs/routes/` declara `schema_version`. | `grep -l '"schema_version"' docs/dags/*.json docs/routes/*.json` sem resultado | VERIFICADO |
| F26 | Harnesses instalados nesta máquina: claude, opencode, mimo, omp e agy. Ausentes: codex, gemini, cursor e kiro. | `command -v` de cada nome | VERIFICADO |
| F27 | O opencode (provider `omniroute`) e o mimo (provider `ai-studio`) apontam para `roteador.vpsconexao.org/v1`, com um único modelo, `ai-studio`. O modelo padrão do opencode é `opencode/big-pickle`. | `~/.config/opencode/opencode.json:207`, `:50`; `~/.config/mimocode/mimocode.jsonc:32` | VERIFICADO |
| F28 | O host `roteador.vpsconexao.org` não resolve. | `curl: (6) Could not resolve host`; `nslookup`: `Non-existent domain` | VERIFICADO |
| F29 | O host `9router.vpsconexao.org` roda a versão 0.5.95 (a última é 0.5.116) e expõe 8 combos, sem os `free-*`. O omp usa esse host e declara dois combos extras (as entradas das linhas 200 e 204 do arquivo), que não existem nesse host nem no de F30. | `curl https://9router.vpsconexao.org/api/version`; `GET /v1/models` com a chave; `~/.omp/agent/models.yml:184`, `:200`, `:204` | VERIFICADO |
| F30 | O host `roteador.signa-docs.com` roda a versão 0.5.116, sem atualização pendente, e expõe 17 combos, com os `free-*`. | `curl https://roteador.signa-docs.com/api/version`; `GET /v1/models` com a chave | VERIFICADO |
| F31 | O PLAN e o REPORT da I1 registram o agy como banido por ToS desde 2026-11-10. | `docs/reports/2026-11-12_REPORT_EXECUCAO_LEI15.md:89` | VERIFICADO |
| F17 | O Claude Code exige comutação para usar os combos do 9Router. | Ver F9. Falta provar que os outros harnesses rodam o nó de código com o 9Router. A bateria prova isso. | PARCIAL |

Decisões do usuário sobre as premissas do INTENT, em 12/11:

* A bateria homologa TODOS os harnesses instalados nesta máquina (F26), não apenas um.
* Sem reprodução e sem mensuração, o F16 não vale. A medição de memória entra como etapa obrigatória.
* O LegoLoop verifica se cada harness usa o 9Router com a API atual, e corrige o que estiver errado (F27 a F30).

Perguntas abertas:

* Q13: qual é o endereço canônico do 9Router?
* Q14: o agy (F31) entra na bateria?

## 4. Pendências para a SPEC

* Criar o comando `legoloop rota` (Q8) sobre o pacote `internal/roteamento`.
* Criar o validador que recusa pergunta de entrevista sem fato citado.
* Adicionar `routing.via` ao contrato v2 e criar o Stop-Gate `SG-I3-CONTRATO` (Q12).
* Definir o perfil em dois níveis e o perfil padrão embutido (Q9).
* Medir o peso de memória de cada harness instalado (F16, obrigatório) dentro da bateria de conformidade.
* Corrigir os providers do 9Router nos harnesses (F27 a F30) e criar a sonda de conectividade e de versão do provider.
* Provar na bateria quais harnesses rodam o nó de código com o 9Router (F17).
* Definir o formato do arquivo de perfil do usuário e o campo `via` do `Candidato`.
* Unir o `provider.Inventory` do motor ao `internal/roteamento` (F7).
* Criar o `skill-router` (F10) sobre o mesmo cofre do roteador de MCP.
* Absorver o E05: manter 05.1, 05.5 e 05.6 adaptados. Substituir 05.2 e 05.4 pelas decisões Q2 e pela classe por regra. Remover 05.3.
* Criar a bateria de conformidade e o comando de homologação de harness, com teto de gasto e aviso antes de rodar.
