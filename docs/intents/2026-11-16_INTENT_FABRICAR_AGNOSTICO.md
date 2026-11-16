---
id: INTENT-FABRICAR-AGNOSTICO
status: DRAFT
date: 2026-11-16
schema_version: 2
upstream: [REPORT-AUDITORIA-VERBOS]
covers: [R1, R2, R3, R4, R5]
---

# 🎯 INTENT: Fabricar Agnóstico — Contrato de I/O e Molde Universal por Natureza

> **Status:** DRAFT
> **Data:** 2026-11-16
> **Autor:** audit-verbos (enxame de auditoria de arquitetura)
> **Iniciativa:** I01 — Agnosticidade Universal dos Verbos de Processo
> **Alvo:** internal/fabricar, internal/processo, internal/matriz, internal/refinar, internal/compor, internal/fechar, internal/listar
> **Sessão:** ab938316-779d-4aa8-9045-fb0d43336880
> **Confiança da Entrevista:** 95% (baseada na auditoria mecânica de código; ver ressalva no item 6)
> **Template:** TEMPLATE_INTENT.md v2.0.0

---

## 1. Contexto e Motivação

* **O Problema Real:** Os 7 verbos de processo violam a premissa de agnosticidade universal: deveriam criar/montar qualquer peça ou processo e na prática usam moldes fixos. Evidência: [2026-16_REPORT_AUDITORIA_VERBOS.md](../reports/2026-11-16_REPORT_AUDITORIA_VERBOS.md) — `fabricar` tem molde fixo de etapas (vocabulário fechado embutido `internal/matriz/vocab.default.json:3-13`), I/O derivado do molde (`internal/fabricar/natureza.go:188-201`, artefato sempre `out/NN_<etapa>.md`) e gate por tipo de artefato (`internal/processo/gates.go:36-48`), não por contrato da peça. Dry-run reproduz 5 LACUNA e gates genéricos.
* **Quem Sofre:** Todo operador que pede um processo fora do molde das 9 etapas canônicas (ex.: mídia, pesquisa, operação): a fabricação gera lacunas e artefatos genéricos que não refletem o que a peça produz.
* **O Porquê:** Sem contrato de I/O na peça, a etapa (e não a peça) define o que existe. Resolver agora converte o molde de restrição em default retrocompatível e destrava os 7 verbos.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** etapas vêm de vocabulário fechado embutido; artefato é `out/NN_<etapa>.<ext>` fixo; gate é deduzido da extensão; aptidão é casamento de palavras-chave; refinar não muda I/O; fechar tem 9 etapas hardcoded Go-cêntricas.
  * **Amanhã:** peça declara `entra`/`sai`/`efeitos` em manifesto `contrato.json`; etapa, gate e aptidão derivam do contrato; nós aceitam 3 naturezas (peça | subprocesso | processo); vocabulário cresce por acréscimo sem recompilar; fechar e refinar parametrizáveis por natureza.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Contrato de I/O declarado por peça (`contrato.json`: entra, sai, efeitos, gate_sugerido) com validação formal e retrocompatibilidade com o molde atual.
* [ ] Objetivo 2: Molde agnóstico por natureza de nó (peça | subprocesso | processo) com derivação de etapa, artefato e gate a partir do contrato.
* [ ] Objetivo 3: Cadeia Sexteto completa (INTERVIEW, INTENT, SPEC, PLAN, ROUTE/DAG, REPORT) exigível e gerada pelo `fabricar` para qualquer natureza.
* [ ] Objetivo 4: Gate anti-reimplementação que barra peça nova redeclarando I/O de peça existente no cofre.
* [ ] Objetivo 5: Varredura dos outros 6 verbos (matriz, listar, refinar, remover, compor, fechar) para consumir o contrato.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Não reimplementaremos peças curadas existentes (`aidd-spec`, `aidd-tickets` etc.): o gate anti-reimplementação existe justamente para isso.
* ⛔ Não migraremos processos v1 do cofre; a retrocompatibilidade preserva o molde atual como default.

---

## 3. Mapeamento Preliminar de Peças

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `unit_no` | Código de produção (Go) `internal/processo`, `internal/fabricar` | Contrato de I/O, derivação de etapa/gate e nós por natureza, via TDD Red-Green. |
| `PECA-02` | `unit_no` | Código de produção (Go) `internal/refinar`, `internal/compor`, `internal/fechar` | Verbos derivados consumindo o contrato. |
| `PECA-03` | `script` | `go test ./...` | Regressão mecânica e integridade. |
| `PECA-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis

* **Gatilho Operacional de Disparo (Trigger):** `manual`
* **Nível de Risco & Blast Radius:** `Alto (exige workspace git_worktree e rollback stash_changes)` — toca o núcleo `internal/processo` e o motor `fabricar`.
* **Política de Privacidade & Orçamento Cognitivo:** `Teto: modelos Flash/Fast` com exceção pontual autorizada para TDD do núcleo. Dados sensíveis ou segredos envolvidos: nenhum.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** `legoloop fabricar --dry-run` de um pedido de natureza "midia" com peça de contrato declarado gera etapa com artefato e gate vindos do contrato (não `out/NN_<etapa>.md`) e zero LACUNA para a etapa coberta pela peça; `go test ./...` verde.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  1. O formato `contrato.json` pode conviver com o `.processo.json` atual sem migração de cofre.
  2. O vocabulário fechado vira semente editável por acréscimo, sem quebrar `matriz.Needs` existente.
  3. A cadeia Sexteto pode ser exigida pelo `fabricar` sem exigir LLM offline (a INTERVIEW pode ser o registro da entrevista já feita).
  4. Os consumidores do contrato são os 7 verbos; não há consumidor externo do `.processo.json`.
* **Riscos e Perguntas em Aberto:**
  1. `fechar` Go-cêntrico (etapa "Recompilação do Binário") vira opcional por natureza — decisão de produto pendente de confirmação.
  2. Gate anti-reimplementação: o critério de "mesmo I/O" deve ser por hash do contrato ou por compatibilidade estrutural? (proposta: compatibilidade estrutural de `entra`/`sai`).

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Peça com `contrato.json` gera etapa com I/O e gate derivados do contrato. | Teste unitário Red-Green no motor | `go test ./internal/fabricar/...` exit 0 |
| 2 | Peça sem contrato mantém o comportamento atual (retrocompatível). | Teste de regressão com fixture atual | `go test ./...` exit 0 |
| 3 | Os 7 verbos aceitam nó de natureza peça, subprocesso e processo. | Testes por verbo + dry-run real | `go test ./internal/refinar/... ./internal/compor/... ./internal/fechar/...` |
| 4 | Peça nova com I/O igual a peça existente é reprovada pelo gate anti-reimplementação. | Teste do gate com dupla fixture | `go test ./internal/...` exit 0 |
| 5 | Cadeia Sexteto exigida e gerada para qualquer natureza. | Verificação dos artefatos gerados no dry-run | saída do `fabricar --dry-run` lista os 6 artefatos |

---

## 6. Handoff de Governança

* **Próximo Documento Mandatório:** [`docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md`](../specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (Stop-Gate) no chat. Ressalva: a confiança de 95% vem da auditoria mecânica de código (provas arquivo:linha no relatório); não houve entrevista Socrática formal — se preferir, o Stop-Gate 01 pode incluir uma rodada de perguntas antes de aprovar.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**
*A intenção consolidada com 95% de certeza factual atende à sua expectativa?*
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*