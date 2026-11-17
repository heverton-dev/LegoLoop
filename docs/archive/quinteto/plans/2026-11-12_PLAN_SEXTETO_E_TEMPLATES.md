---
id: PLAN-I2
status: APPROVED
date: 2026-11-12
schema_version: 2
upstream: [SPEC-I2]
covers: [AC-I2-00.1, AC-I2-00.2, AC-I2-00.3, AC-I2-01.1, AC-I2-01.2, AC-I2-01.3, AC-I2-01.4, AC-I2-02.1, AC-I2-02.2, AC-I2-02.3, AC-I2-02.4, AC-I2-02.5, AC-I2-03.1, AC-I2-03.2, AC-I2-04.1, AC-I2-04.2, AC-I2-04.3, AC-I2-04.4, AC-I2-04.5, AC-I2-04.6, AC-I2-05.1]
---

# 🗺️ PLAN: Sexteto e templates agnósticos (Iniciativa I2)

> **Status:** Aguardando aprovação do usuário
> **Data:** 2026-11-12
> **Versão:** 1.0.0
> **Autor:** Claude
> **Entrevista:** [INTERVIEW](../interviews/2026-11-12_INTERVIEW_SEXTETO_TEMPLATES.md)
> **Intenção:** [INTENT](../intents/2026-11-12_INTENT_SEXTETO_E_TEMPLATES.md)
> **Especificação:** [SPEC](../specs/2026-11-12_SPEC_SEXTETO_E_TEMPLATES.md)
> **Rota:** [ROUTE](../routes/2026-11-12_ROUTE_SEXTETO_E_TEMPLATES.json)
> **DAG:** [DAG](../dags/2026-11-12_DAG_SEXTETO_E_TEMPLATES.json)

---

## 1. Decisões travadas (usuário, 2026-11-12)

1. O motor valida DAG e ROUTE. Agente e humano consomem os documentos com seções fixas conferidas por teste (Q1).
2. A ROUTE é a fonte da decisão de roteamento. O DAG carrega o roteamento resolvido (Q2).
3. O gerador classifica o nó por regra mensurável. O autor sobrescreve com justificativa (Q3).
4. O gerador degrada para o candidato mais próximo e registra a troca. `extreme` nunca degrada em silêncio (Q4).
5. A rastreabilidade Q, R, AC, nó e `spec_ref` é validada pelo motor (Q5).
6. O motor gera a parte factual do REPORT. O HANDOFF recorta o REPORT (Q6).
7. Arquivo sem `schema_version` roda em modo legado. A migração é a iniciativa I6 (Q7).
8. Todo template é agnóstico em SO, harness e modo de execução (`cli`, `agente`, `estudio`) (Q8).
9. As peças skill e MCP carregam `lazy` por padrão (Q9).

## 2. Linha de base

* O DAG da Lei #15 e os DAGs de LegoPress e VideoREA não têm `schema_version`. Eles rodam em modo legado.
* O `unit_task` roda só modelos do claude CLI. O roteamento desta iniciativa usa `claude` com `haiku`, `sonnet` e `opus`.
* A SPEC fixa os contratos. O PLAN não os repete.

## 3. Estratégia de execução

* **Engine:** `orchestrix_hermetic`. Um ticket por contexto limpo. Sem subagentes.
* **Modelo e esforço por classe:** `nano` Haiku `low`. `micro` Sonnet `medium`. `macro` Sonnet `medium` ou `high`. `extreme` Opus `high` com Stop-Gate.
* **Classificação:** cada nó grava `class_rule`. Esta ROUTE foi classificada pela regra da SPEC seção 4.4. O `TICK-I2-00.3` tem sobrescrita para `extreme`, com justificativa.
* **Paralelismo:** `--parallel 1` com `GOFLAGS=-p=2`. Os tickets de `docs/templates` editam os mesmos arquivos e rodam em cadeia.
* **Gate padrão:** `go run ./tools/ticketgate` com `-pkgs` e `-tax` do nó. O `TICK-I2-00.2` acrescenta `-trace`.
* **Rollback:** `stash_changes` por ticket. Commit por nó só dos arquivos do `context_out`. Sem push.
* **Custo:** parar e chamar o usuário se um nó passar de US$ 1,50.
* **Dados reais:** nenhum ticket toca o Cofre real `~/.aidd`.

## 4. Mapa de épicos

| Épico | Nome | Nós | Classe dominante |
| :--- | :--- | ---: | :--- |
| **E00** | Fundação | 4 | `micro` e `macro` |
| **E01** | Roteamento por capacidade | 4 | `micro` e `macro` |
| **E02** | Templates do sexteto | 5 | `micro` e `macro` |
| **E03** | REPORT e HANDOFF | 2 | `micro` e `macro` |
| **E04** | Peças agnósticas | 6 | `micro` |
| **E05** | Fechamento | 2 | `nano` |

Ordem de execução: E00, E01, E02 (02.1 a 02.4), E03, E04, 02.5, E05.

## 5. Tickets

### Épico 00 — Fundação: front matter, rastreabilidade e contrato v2

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I2-00.1` | Front matter dos artefatos | `micro` Sonnet `medium` | [] | `AC-I2-00.1` |
| `TICK-I2-00.2` | Validador de rastreabilidade e opção -trace do ticketgate | `macro` Sonnet `high` | `TICK-I2-00.1` | `AC-I2-00.2` |
| `SG-I2-CONTRATO` | Stop-Gate Humano (SG-I2-CONTRATO) | `nano` humano | `TICK-I2-00.1`, `TICK-I2-00.2` | `-` |
| `TICK-I2-00.3` | Contrato schema_version 2 e modo legado do DAG | `extreme` Opus `high` | `SG-I2-CONTRATO` | `AC-I2-00.3` |

### Épico 01 — Roteamento por capacidade

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I2-01.1` | Classe de capacidade por regra | `micro` Sonnet `medium` | [] | `AC-I2-01.1` |
| `TICK-I2-01.2` | Inventário do SO (harnesses, binários, modelos) | `micro` Sonnet `medium` | `TICK-I2-01.1` | `AC-I2-01.2` |
| `TICK-I2-01.3` | Resolvedor com degradação registrada | `micro` Sonnet `medium` | `TICK-I2-01.2` | `AC-I2-01.3` |
| `TICK-I2-01.4` | Paridade ROUTE e DAG no route2dag | `macro` Sonnet `high` | `TICK-I2-01.3`, `TICK-I2-00.3` | `AC-I2-01.4` |

### Épico 02 — Templates do sexteto

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I2-02.1` | TEMPLATE_INTERVIEW | `micro` Sonnet `medium` | `TICK-I2-00.2` | `AC-I2-02.1` |
| `TICK-I2-02.2` | TEMPLATE_INTENT sem termos proibidos e com front matter | `macro` Sonnet `medium` | `TICK-I2-02.1` | `AC-I2-02.2` |
| `TICK-I2-02.3` | TEMPLATE_SPEC e TEMPLATE_PLAN com IDs R e AC | `macro` Sonnet `medium` | `TICK-I2-02.2` | `AC-I2-02.3` |
| `TICK-I2-02.4` | TEMPLATE_ROUTE e TEMPLATE_DAG completos | `macro` Sonnet `high` | `TICK-I2-02.3`, `TICK-I2-01.4` | `AC-I2-02.4` |
| `TICK-I2-02.5` | Testes TestTemplatesDoSexteto e TestPecasAgnosticas | `micro` Sonnet `medium` | `TICK-I2-04.6` | `AC-I2-02.5` |

### Épico 03 — REPORT e HANDOFF

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I2-03.1` | TEMPLATE_REPORT factual e interpretativo e TEMPLATE_HANDOFF | `micro` Sonnet `medium` | `TICK-I2-02.4` | `AC-I2-03.1` |
| `TICK-I2-03.2` | Gerador factual do REPORT | `macro` Sonnet `high` | `TICK-I2-00.1` | `AC-I2-03.2` |

### Épico 04 — Peças agnósticas

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I2-04.1` | Auditoria do carregamento sob demanda de skills | `micro` Sonnet `medium` | [] | `AC-I2-04.1` |
| `TICK-I2-04.2` | TEMPLATE_PIECE_SKILL agnóstico e sob demanda | `micro` Sonnet `medium` | `TICK-I2-04.1`, `TICK-I2-03.1` | `AC-I2-04.2` |
| `TICK-I2-04.3` | TEMPLATE_PIECE_MCP com registro no cofre | `micro` Sonnet `medium` | `TICK-I2-04.2` | `AC-I2-04.3` |
| `TICK-I2-04.4` | TEMPLATE_PIECE_HOOK neutro por harness e SO | `micro` Sonnet `medium` | `TICK-I2-04.3` | `AC-I2-04.4` |
| `TICK-I2-04.5` | TEMPLATE_PIECE_SCRIPT em par ps1 e sh | `micro` Sonnet `medium` | `TICK-I2-04.4` | `AC-I2-04.5` |
| `TICK-I2-04.6` | TEMPLATE_PIECE_GATE e nova peça de Processo agnósticas | `micro` Sonnet `medium` | `TICK-I2-04.5` | `AC-I2-04.6` |

### Épico 05 — Fechamento

| Ticket | Ação | Classe / Modelo | Dep. | AC |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-I2-05.1` | Fechamento: build fresco, go test global e taxlinter dos templates | `nano` Haiku `low` | `TICK-I2-02.5`, `TICK-I2-03.2`, `TICK-I2-01.4` | `AC-I2-05.1` |
| `SG-I2-FINAL` | Stop-Gate Humano (SG-I2-FINAL) | `nano` humano | `TICK-I2-05.1` | `-` |

#### 🛑 STOP-GATE HUMANO — `SG-I2-CONTRATO`
* **Gate:** `prompt('Aprova mudar o contrato público do DAG (schema_version 2 com modo legado)?')`
* Antes do `TICK-I2-00.3`. Só o usuário aprova.

#### 🛑 STOP-GATE HUMANO — `SG-I2-FINAL`
* **Gate:** `prompt('Aprova o fechamento da I2?')`
* Só o usuário aprova.

## 6. Riscos

| Risco | Efeito | Mitigação |
| :--- | :--- | :--- |
| Os testes existentes de `docs/templates` falham com o template novo. | O gate do nó fica vermelho. | Cada nó de template inclui `templates_test.go` no `context_out`. |
| O nome `schema_version` colide com `$schema`. | Confusão no parser. | A SPEC seção 1, premissa 8, fixa a regra. |
| O inventário de modelos varia por harness. | O resolvedor erra o candidato. | O contrato mínimo está na SPEC seção 4.5. O teste usa inventário simulado. |
| Memória insuficiente na rodada. | Kill da execução. | `--parallel 1`. Reinício só a pedido do usuário. |
| Três falhas com a mesma causa. | Perda de tempo e custo. | Parar e chamar o usuário. |

## 7. Custo estimado

* Nós `nano`: US$ 0,05 a 0,15. Nós `micro`: US$ 0,30 a 0,60. Nós `macro`: US$ 0,60 a 1,20. Nó `extreme`: US$ 1,50 a 3,00.
* Total estimado: US$ 10 a 18. A Lei #15 mediu US$ 0,20 a 1,22 por nó.
