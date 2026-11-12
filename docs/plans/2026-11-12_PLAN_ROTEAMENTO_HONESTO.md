---
id: PLAN-I3
status: UNDER_REVIEW
date: 2026-11-12
schema_version: 2
upstream: [SPEC-I3]
covers: [AC-I3-00.1, AC-I3-00.2, AC-I3-00.3, AC-I3-01.1, AC-I3-01.2, AC-I3-01.3, AC-I3-01.4, AC-I3-02.1, AC-I3-02.2, AC-I3-02.3, AC-I3-02.4, AC-I3-02.5, AC-I3-02.6, AC-I3-03.1, AC-I3-03.2, AC-I3-03.3, AC-I3-03.4, AC-I3-04.1, AC-I3-05.1, AC-I3-05.2]
---

# 🗺️ PLAN: Roteamento honesto (Iniciativa I3)

> **Status:** Aguardando aprovação do usuário
> **Data:** 2026-11-12
> **Versão:** 1.0.0
> **Autor:** Claude
> **Entrevista:** [INTERVIEW](../interviews/2026-11-12_INTERVIEW_ROTEAMENTO_HONESTO.md)
> **Intenção:** [INTENT](../intents/2026-11-12_INTENT_ROTEAMENTO_HONESTO.md)
> **Especificação:** [SPEC](../specs/2026-11-12_SPEC_ROTEAMENTO_HONESTO.md)
> **Rota:** [ROUTE](../routes/2026-11-12_ROUTE_ROTEAMENTO_HONESTO.json)
> **DAG:** [DAG](../dags/2026-11-12_DAG_ROTEAMENTO_HONESTO.json)

---

## 1. Decisões travadas (usuário, 2026-11-12)

1. A honestidade da ROUTE é a prioridade 1 (Q1).
2. Em falha de execução, os nós `nano` e `micro` degradam com registro. Os nós de classe grande e `extreme` falham (Q2).
3. O arquivo do usuário declara a preferência. O inventário valida. O 9Router é um eixo `via`, nunca um harness (Q3).
4. O motor comuta o gateway por nó, só no processo filho (Q4).
5. O `call_skill` devolve o texto da skill. O nó do DAG executa (Q5).
6. O motor limita os workers pela memória livre e pelo peso de cada harness (Q6).
7. As raias do agy saem do E05. O alarme de banimento marca o harness como reprovado (Q7).
8. O comando `legoloop rota` em Go preenche o roteamento da ROUTE (Q8).
9. O perfil vive em dois níveis, global e por projeto, com o padrão embutido (Q9).
10. A bateria mostra o teto total e pede confirmação antes de gastar (Q10).
11. A barreira de escrita da bateria é a conferência do motor (Q11).
12. O campo `routing.via` é obrigatório no `schema_version: 2`, com Stop-Gate humano (Q12).
13. O endereço canônico do 9Router é `roteador.signa-docs.com` (Q13).
14. O agy entra na bateria com conta fixa, binário direto, parada no primeiro sinal de ToS e aviso próprio (Q14).
15. Todos os harnesses instalados entram na bateria, e o F16 só vale depois de medido.

## 2. Linha de base

* Todos os arquivos de `docs/dags/` e `docs/routes/` rodam em modo legado (F25). A mudança de contrato não quebra nenhum.
* O nó de código roda só no claude CLI (F1, F5). Por isso esta ROUTE usa `claude` com `haiku`, `sonnet` e `opus`. O `via` entra no contrato no ticket `TICK-I3-00.1`, e o DAG desta iniciativa o usa só depois.
* O inventário atual do `internal/roteamento` detecta apenas claude, codex, gemini, opencode, cursor e kiro (`internal/processo/vinculo.go:11`). Ele não vê mimo, omp nem agy, que estão instalados aqui (F26, F36). O ticket `TICK-I3-01.1` unifica o inventário com o `provider.Inventory`.
* O `SkillEngine` já existe e nada o usa (F35). O ticket `TICK-I3-04.1` o liga ao roteador.
* A SPEC fixa os contratos. O PLAN não os repete.

## 3. Estratégia de execução

* **Engine:** `orchestrix_hermetic`. Um ticket por contexto limpo. Sem subagentes.
* **Modelo e esforço por classe:** `nano` Haiku `low`. `micro` Sonnet `medium`. Classe grande Sonnet `high`. `extreme` Opus `high` com Stop-Gate.
* **Classificação:** cada nó grava `class_rule`. O pacote `internal/roteamento` classificou esta ROUTE pela regra da SPEC-I2 seção 4.4. A resolução rodou contra o inventário desta máquina (claude com `haiku`, `sonnet` e `opus`). Nenhum nó recebeu classe à mão. O `TICK-I3-00.1` é `extreme` porque muda o contrato público.
* **Paralelismo:** `--parallel 1` com `GOFLAGS=-p=2`. Os tickets de uma mesma pasta editam os mesmos arquivos e rodam em cadeia.
* **Gate padrão:** `go run ./tools/ticketgate` com `-pkgs` e `-tax` do nó. O ticket `TICK-I3-00.3` acrescenta a opção `-interview`.
* **Rollback:** `stash_changes` por ticket. Commit por nó só dos arquivos do `context_out`. Sem push.
* **Custo:** parar e chamar o usuário se um nó passar de US$ 1,50.
* **Dados reais:** nenhum ticket toca o Cofre real `~/.aidd`. Os testes usam diretórios temporários, harnesses falsos e servidor HTTP falso.
* **Chamadas reais de modelo:** só o `SG-I3-HOMOLOGACAO`, executado pelo usuário, com confirmação do teto.

## 4. Mapa de épicos

| Épico | Nome | Nós | Classe dominante |
| :--- | :--- | ---: | :--- |
| **E00** | Contrato e prova | 4 | grande, `extreme` e `micro` |
| **E01** | Inventário, perfil e gerador | 4 | `micro` |
| **E02** | Execução honesta | 6 | `micro` e grande |
| **E03** | Homologação | 4 | `micro` e grande |
| **E04** | Skill-router | 1 | `micro` |
| **E05** | Fechamento | 3 | `nano` |

Total: 22 nós, dos quais 3 são Stop-Gates humanos.

## 5. Tickets

### Épico 00 — Contrato e prova

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I3-00.3` | Regra de prova da entrevista e opção `-interview` do ticketgate | grande Sonnet `high` | [] | `AC-I3-00.3` |
| `SG-I3-CONTRATO` | Stop-Gate Humano (SG-I3-CONTRATO) | `nano` humano | `TICK-I3-00.3` | `-` |
| `TICK-I3-00.1` | `routing.via` obrigatório no contrato v2 do DAG | `extreme` Opus `high` | `SG-I3-CONTRATO` | `AC-I3-00.1` |
| `TICK-I3-00.2` | Paridade ROUTE e DAG com `via` e harness `direto` | `micro` Sonnet `medium` | `TICK-I3-00.1` | `AC-I3-00.2` |

### Épico 01 — Inventário, perfil e gerador

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I3-01.1` | Inventário unificado por par harness e `via` | `micro` Sonnet `medium` | `TICK-I3-00.1` | `AC-I3-01.1` |
| `TICK-I3-01.2` | Perfil em dois níveis com padrão embutido | `micro` Sonnet `medium` | `TICK-I3-01.1` | `AC-I3-01.2` |
| `TICK-I3-01.3` | Resolvedor com `via` e política de falha por classe | `micro` Sonnet `medium` | `TICK-I3-01.2` | `AC-I3-01.3` |
| `TICK-I3-01.4` | Comando `legoloop rota` | `micro` Sonnet `medium` | `TICK-I3-01.3`, `TICK-I3-00.2` | `AC-I3-01.4` |

### Épico 02 — Execução honesta

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I3-02.1` | Nó de código em harness homologado | `micro` Sonnet `medium` | `TICK-I3-01.3` | `AC-I3-02.1` |
| `TICK-I3-02.2` | Comutação do gateway por nó | `micro` Sonnet `medium` | `TICK-I3-02.1` | `AC-I3-02.2` |
| `TICK-I3-02.3` | Despacho com degradação por classe e API direta | `micro` Sonnet `medium` | `TICK-I3-02.2` | `AC-I3-02.3` |
| `TICK-I3-02.4` | Limite de workers por memória livre | grande Sonnet `high` | `TICK-I3-02.3` | `AC-I3-02.4` |
| `TICK-I3-02.5` | Detecção de conflito de arquivo entre nós | `micro` Sonnet `medium` | `TICK-I3-02.4` | `AC-I3-02.5` |
| `TICK-I3-02.6` | REPORT factual com harness e `via` reais | `micro` Sonnet `medium` | `TICK-I3-02.3` | `AC-I3-02.6` |

### Épico 03 — Homologação

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I3-03.1` | Bateria de conformidade com medição de memória | `micro` Sonnet `medium` | `TICK-I3-01.1`, `TICK-I3-02.1` | `AC-I3-03.1` |
| `TICK-I3-03.2` | Armazenamento da homologação e comando `homologar` | grande Sonnet `high` | `TICK-I3-03.1`, `TICK-I3-01.4` | `AC-I3-03.2` |
| `TICK-I3-03.3` | Condições do agy e alarme de banimento | grande Sonnet `high` | `TICK-I3-03.2` | `AC-I3-03.3` |
| `TICK-I3-03.4` | Sonda e correção do provider 9Router com backup | grande Sonnet `high` | `TICK-I3-03.2` | `AC-I3-03.4` |

### Épico 04 — Skill-router

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I3-04.1` | `skill_search` e `call_skill` no roteador | `micro` Sonnet `medium` | `TICK-I3-00.3` | `AC-I3-04.1` |

### Épico 05 — Fechamento

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I3-05.1` | Fechamento: build fresco, go test global e taxlinter dos templates | `nano` Haiku `low` | `TICK-I3-00.2`, `TICK-I3-01.4`, `TICK-I3-02.5`, `TICK-I3-02.6`, `TICK-I3-03.3`, `TICK-I3-03.4`, `TICK-I3-04.1` | `AC-I3-05.1` |
| `SG-I3-HOMOLOGACAO` | Stop-Gate Humano (SG-I3-HOMOLOGACAO) | `nano` humano | `TICK-I3-05.1` | `AC-I3-05.2` |
| `SG-I3-FINAL` | Stop-Gate Humano (SG-I3-FINAL) | `nano` humano | `SG-I3-HOMOLOGACAO` | `-` |

#### 🛑 STOP-GATE HUMANO — `SG-I3-CONTRATO`
* **Gate:** `prompt('Aprova incluir routing.via como campo obrigatório do contrato público do DAG v2?')`
* Antes do `TICK-I3-00.1`. Só o usuário aprova.

#### 🛑 STOP-GATE HUMANO — `SG-I3-HOMOLOGACAO`
* **Gate:** o usuário roda `legoloop homologar` com confirmação do teto e `legoloop homologar --provider` com `--corrigir-provider`. Ele aprova o resultado de claude, opencode, mimo, omp e agy nesta máquina.
* Só o usuário aprova. É o único ponto com chamada real de modelo.

#### 🛑 STOP-GATE HUMANO — `SG-I3-FINAL`
* **Gate:** `prompt('Aprova o fechamento da I3?')`
* Só o usuário aprova.

## 6. Riscos

| Risco | Efeito | Mitigação |
| :--- | :--- | :--- |
| Dois nós do mesmo pacote entram em conflito (tipos duplicados, nomes de pacote distintos), como na I2. | O gate de cada nó passa e o conjunto falha. | Rodar `go test -count=1 ./...` depois de cada rodada, antes do commit. Os nós do mesmo pacote rodam em cadeia. |
| O nó `TICK-I3-00.1` muda o contrato público. | DAGs existentes deixam de validar. | O `SG-I3-CONTRATO` vem antes. Nenhum arquivo atual declara `schema_version` (F25). O teste roda `legoloop validate` nos DAGs de `docs/dags/`. |
| O inventário atual não vê mimo, omp e agy (F36). | A bateria ignora harnesses instalados. | O `TICK-I3-01.1` une o `provider.Inventory`. O teste simula os cinco harnesses. |
| A bateria gasta dinheiro e memória. | Custo inesperado ou kill por memória. | Aviso do teto total, confirmação, conferência da memória livre e `--non-interactive` sem chamada. |
| O agy responde com ToS durante a bateria. | Piora do estado da conta. | Conta fixa, binário direto, parada no primeiro sinal e marca "reprovado". |
| Corrigir a configuração de um harness quebra a configuração do usuário. | O harness deixa de abrir. | Backup `.legoloop-bak`, troca só do endereço, releitura e confirmação por arquivo. |
| O peso de 750 MB para harness sem medição é conservador e não medido. | Workers abaixo do necessário. | A bateria substitui o valor pela medição (F16). |
| Memória insuficiente na rodada. | Kill da execução. | `--parallel 1`. Reinício só a pedido do usuário. |
| Três falhas com a mesma causa. | Perda de tempo e custo. | Parar e chamar o usuário. |

## 7. Custo estimado

* Nós `nano`: US$ 0,05 a 0,15. Nós `micro`: US$ 0,30 a 0,60. Nós de classe grande: US$ 0,60 a 1,20. Nó `extreme`: US$ 1,50 a 3,00.
* Total estimado: US$ 12 a 20 (12 nós `micro`, 5 de classe grande, 1 `extreme`, 1 `nano`). A Lei #15 mediu US$ 0,20 a 1,22 por nó, e o ticket de lacunas da I2 mediu US$ 0,06 a 0,09.
* A bateria real (`SG-I3-HOMOLOGACAO`) tem custo próprio, mostrado antes da confirmação.
