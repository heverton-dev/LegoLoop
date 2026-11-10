---
name: legoloop-n10-foundry
description: "Use quando precisar operar, auditar ou evoluir a Fundição Autônoma de Software LegoLoop com os 7 Protocolos Industriais de Missão Crítica (BLUF/STANAG, ASD-STE100, EARS, FMEA, ICAO Readback, Toyota Jidoka/Andon Cord, SBAR) e os 6 Motores Globais de Fronteira N10 (Anti-Reward Hacking, Merge AST, Sistema Imunológico antibodies.db, FinOps Especulativo, Caixa-Preta OTel GenAI e Fuzzing Nativo)."
---

# ♾️ LegoLoop N10 Foundry — A Fundição Autônoma de Software

Esta skill governa a operação de engenharia em nível N10 da fábrica de software LegoLoop.  
Ela aplica determinismo rigoroso em todas as Etapas da cascata de valor: **INTENT ➔ SPEC ➔ PLAN ➔ ROUTE ➔ BUILD (TDD) ➔ HOMOLOGAÇÃO**.

---

## 1. As 7 Regras de Ouro dos Protocolos Industriais

1. **NATO STANAG 2014 & BLUF:** Toda instrução segue a ordem de operação em 5 partes (Situação, Missão, Execução, Logística, Comando).
2. **ASD-STE100:** Frases com no máximo 20 palavras em procedimentos, voz ativa e zero termos de indefinição (*hedging*).
3. **EARS Requisitos:** Toda regra de negócio escrita em templates `WHEN...`, `WHILE...` ou `IF... THEN...`.
4. **FMEA Pré-Voo:** Matriz de riscos gerada antes de qualquer mutação. Severidade $\ge 3$ exige teste Poka-Yoke prévio.
5. **ICAO Readback/Hearback:** Ações destrutivas exibem cartão com `TARGET`, `MUTATION`, `BLAST RADIUS` e `ROLLBACK`, exigindo confirmação `HEARBACK APPROVED`.
6. **Toyota Jidoka (Andon Cord):** Limite máximo de 2 tentativas por Nó. Na segunda falha com mesma assinatura, puxa o Andon Cord.
7. **Escalonamento SBAR:** Falhas de Nó emitem relatório estruturado em `[S] SITUATION`, `[B] BACKGROUND`, `[A] ASSESSMENT` e `[R] RECOMMENDATION`.

---

## 2. Os 6 Motores de Fronteira Industrial (Nível N10)

1. **Motor I — Anti-Reward Hacking:** Agentes Blue-Team nunca recebem permissão de escrita em `*_test.go`.
2. **Motor II — Compilador AST em Worktree:** Mesclagem determinística de código Go via AST (`go/ast`), prevenindo conflitos textuais.
3. **Motor III — Sistema Imunológico antibodies.db:** Gravação contínua de triplas SBAR e injeção de contexto preventivo.
4. **Motor IV — Cascade Speculative FinOps Routing:** Despacho para modelos econômicos com hot-swap para modelos pesados em falhas.
5. **Motor V — Caixa-Preta OpenTelemetry:** Medição de *First-Pass Yield (FPY)* e persistência de `flight_record.json`.
6. **Motor VI — Property-Based Testing & Fuzzing Nativo:** Exigência de testes `testing.F` em todos os parsers e crawlers.

---

## 3. Comandos Canônicos

```bash
# Validar contrato e proveniência semântica:
legoloop validate processo.dag.json

# Executar com isolamento Orchestrix e Andon Cord:
legoloop run processo.dag.json --parallel 2 --json

# Abrir o Studio Visual:
legoloop studio --addr 127.0.0.1:7878
```
