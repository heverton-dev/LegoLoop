---
id: INTENT-AUDITORIA-SWARM-01
status: DRAFT
date: 2026-11-14
schema_version: 2
upstream: [SESSION-d8076846]
covers: [Q1]
---

# 🎯 INTENT: Auditoria Profunda de Código e UI/UX em Modo Swarm (4 Frentes)

> **Status:** DRAFT  
> **Data:** 2026-11-14  
> **Autor:** Antigravity (Coordinator)  
> **Iniciativa:** Auditoria Profunda Swarm  
> **Alvo:** `internal/`, `cmd/legoloop/`, `internal/studio/`, `loops/`, `schemas/`  
> **Sessão:** `2b93726e-e3b1-4f72-aa2a-d2ddcd8ca025`  
> **Confiança da Entrevista:** 96%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Após a entrega das iniciativas I5, I6, I7 e da Frente B do Studio Web, o repositório acumulou alterações rápidas em Go, UI/UX e manifestos DAG que precisam de auditoria profunda independente antes de novas evoluções.
* **Quem Sofre:** Desenvolvedores e operadores do LegoLoop Studio, sujeitos a regressões silenciosas de interface, dívida técnica ou desvios de contrato.
* **O Porquê:** Executar uma auditoria paralela em 4 frentes via Swarm reduz o tempo de inspeção e garante cobertura completa sem poluir o contexto único.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Verificação fragmentada por demandas isoladas, sem visão consolidada de saúde de código, UI/UX e contratos.
  * **Amanhã:** Relatório único auditável com achados classificados por severidade, evidências reais e plano de remediação priorizado.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Auditar código Go (`internal/`, `cmd/`) quanto a complexidade, `funclen`, `taxlinter`, `go vet` e cobertura de testes.
* [ ] Objetivo 2: Auditar UI/UX do Studio Web (`internal/studio/`) quanto a acessibilidade, responsividade, contraste e estados vazios.
* [ ] Objetivo 3: Auditar contratos JSON, schemas e 23 DAGs quanto à conformidade com `schema_version: 2`.
* [ ] Objetivo 4: Consolidar os achados das 4 frentes em relatório executivo com evidências reproduzíveis.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Não aplicaremos correções destrutivas durante a fase de auditoria (somente leitura e diagnóstico).

---

## 3. Mapeamento Preliminar de Peças

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `script` | `go test ./... && go vet ./...` | Auditoria mecânica de código Go e testes. |
| `PECA-02` | `skill` | `skills/aidd-uiux` | Auditoria de UI/UX, contraste e acessibilidade no Studio Web. |
| `PECA-03` | `script` | `legoloop validate` | Auditoria de schemas e integridade dos DAGs. |
| `PECA-04` | `stop_gate` | Revisão Humana | Aprovação do relatório consolidado e priorização de correções. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`
* **Nível de Risco & Blast Radius:** `Baixo (workspace local, somente leitura na auditoria)`
* **Política de Privacidade & Orçamento Cognitivo:** `Teto: modelos Flash/Fast para leitura e Pro para consolidação`. Dados sensíveis: `nenhum (nunca ler ~/.aidd/secrets)`.
* **Contrato de Evidência Factual:** Relatório Markdown em `docs/reports/2026-11-14_REPORT_AUDITORIA_PROFUNDA_SWARM.md` com comandos e saídas reais.
* **Premissas Assumidas:** Os 4 workers do Swarm executam em contexto estéril e escrevem apenas seus artefatos de saída.
* **Riscos e Perguntas em Aberto:** Nenhum bloqueio conhecido para iniciar a leitura.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as 4 frentes produzem artefatos de achados sem stubs. | Arquivos JSON/MD presentes em `docs/reports/` | `Test-Path docs/reports/2026-11-14_REPORT_AUDITORIA_PROFUNDA_SWARM.md` |
| 2 | Suíte Go e linters executados com saída real registrada. | Exit code e logs capturados | `go test -count=1 ./...` |
| 3 | Consolidação final apresenta matriz de severidade e próximos passos. | Seção de achados preenchida | Leitura do relatório final |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-14_SPEC_AUDITORIA_PROFUNDA_SWARM.md`](../specs/2026-11-14_SPEC_AUDITORIA_PROFUNDA_SWARM.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (Stop-Gate) no chat.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 96% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se avançar para a execução da SPEC.*
