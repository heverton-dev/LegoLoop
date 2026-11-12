---
id: INTERVIEW-{{ID_DA_INICIATIVA}}
status: DRAFT
date: {{AAAA-MM-DD}}
schema_version: 2
upstream: []
covers: [Q1]
---

# INTERVIEW: {{TITULO_DA_ENTREVISTA}}

> **Protocolo:** uma pergunta por vez, em prosa. A entrevista só termina com confiança igual ou acima de 95%.
> **Regra de prova:** nenhuma pergunta nasce de inferência. Cada pergunta cita fatos F da seção 3, e cada fato F traz a prova (arquivo e linha, ou comando e saída). Sem prova, o entrevistador verifica no código antes de perguntar. Fato de memória ou de sessão anterior não conta como prova.
> **Confiança final:** {{NN}}%
> **Iniciativa:** {{NOME_OU_ID_DA_INICIATIVA}}
> **Template:** TEMPLATE_INTERVIEW.md v2.3.0

## 1. Respostas

| ID | Pergunta | Fatos | Decisão do usuário |
| :--- | :--- | :--- | :--- |
| Q1 | {{PERGUNTA_1}} | {{F1, F2}} | {{DECISAO_1}} |
| Q2 | {{PERGUNTA_2}} | {{F3}} | {{DECISAO_2}} |

## 2. Confiança por rodada

| Rodada | Perguntas feitas | Confiança | Motivo da variação |
| :--- | :--- | :--- | :--- |
| 1 | {{Q1, Q2}} | {{NN}}% | {{O que ficou claro e o que continua aberto.}} |

## 3. Fatos medidos durante a entrevista

| ID | Afirmação | Prova (arquivo:linha ou comando) | Resultado |
| :--- | :--- | :--- | :--- |
| F1 | {{AFIRMACAO}} | {{ARQUIVO:LINHA_OU_COMANDO}} | {{VERIFICADO ou REFUTADO}} |

## 4. Pendências para a SPEC

* {{PENDENCIA_1}}
