# 🎯 INTENT: GOVERNANCA_ROUTE_V2_CONCORRENCIA

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-07  
> **Autor:** Antigravity & Engenharia LegoLoop  
> **Iniciativa:** I01 — Governança e Resiliência da Tríade Executável  
> **Alvo:** `docs/templates/`, `internal/plangen`, `internal/fabricar`, `internal/validator`, `internal/orchestrator`  
> **Sessão:** 6353f637-cb0a-41a1-92ee-8c1400e75faa  
> **Confiança da Entrevista:** 98% (requisito >= 95% atendido)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação

* **O Problema Real:** O contrato atual da ROUTE (`route.v2.json`) possui parâmetros genéricos de execução (`max_parallel_workers` como inteiro solto, `rollback_on_failure` como string rasa, sem teto de tokens de raciocínio e sem validação criptográfica entre os artefatos da Tríade). Além disso, não reflete a nova hierarquia canônica (**Iniciativa $\rightarrow$ Épico $\rightarrow$ Micro/Nano Tickets**) nem desacopla formalmente o **espaço físico de arquivo** (`workspace`) do **motor executor** (`executor`), abrindo brechas para concorrência descontrolada, queima de tokens por subagentes com contexto herdado e divergência silenciosa de arquivos upstream.
* **Quem Sofre:** Desenvolvedores e agentes autônomos que operam pipelines complexos no LegoLoop, sofrendo com travamentos por falta de teto de CPU/RAM, erros de rate limit (429) em APIs e desvios de integridade quando um arquivo de SPEC ou PLAN é alterado no meio do ciclo.
* **O Porquê:** Blindar a esteira com determinismo absoluto, garantindo que execuções assíncronas paralelas ocorram apenas via subagentes efêmeros limpos sob teto estrito de hardware/harness, com integridade do quarteto validada por SHA-256 e fail-fast cirúrgico.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Parâmetros soltos no cabeçalho da rota, ausência de verificação de hash entre `INTENT`, `SPEC`, `PLAN` e `ROUTE`, sem controle de teto de subagentes efêmeros e rollback simplista.
  * **Amanhã:** Cabeçalho da ROUTE v2.1 tipado e estrito, com verificação de cadeia criptográfica SHA-256 (Merkle encadeada), matriz ortogonal de isolamento (`workspace` $\times$ `executor`), bloco de concorrência com teto de subagentes paralelos e política de rollback granular.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Atualizar o template e schema canônico da ROUTE (`docs/templates/route.template.json`) para a versão `2.1.0`, contemplando assinaturas SHA-256 upstream, bloco `concurrency` com teto de subagentes, matriz ortogonal de isolamento (`workspace` vs. `executor`), `budget` de tokens e `rollback_policy` estruturada.
* [ ] Objetivo 2: Implementar no motor Go (`internal/plangen`, `internal/fabricar`) a geração determinística do bloco `upstream_signatures` (calculando SHA-256 de `INTENT`, `SPEC` e `PLAN`) e o preenchimento dos novos campos de governança.
* [ ] Objetivo 3: Criar validador de integridade no runtime do LegoLoop para confrontar os hashes do quarteto no boot da execução, bloqueando desvios silenciosos.
* [ ] Objetivo 4: Garantir retrocompatibilidade com rotas v2.0 existentes (tratando campos ausentes com defaults seguros sem quebrar o ecossistema).
* [ ] Objetivo 5: Validar a suíte inteira via TDD com 100% de conformidade aos portões de qualidade do repositório (`go test`, `funclen`, `gofmt`).

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além da biblioteca padrão do Go (`crypto/sha256`, `encoding/json`, etc.).
* ⛔ Não implementaremos containers Docker reais ou sandboxes remotas nesta etapa; apenas tiparemos as opções canônicas no schema/templates para suportá-las nativamente.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `aidd-spec` / `spec-driven-development` | Elaborar especificação formal (SPEC v2.1) com critérios BDD de validação de rota. |
| `BLOCO-02` | `unit_task` | Structs & Templates (`docs/templates/`) | Atualização dos templates canônicos e schema formal com testes em `templates_test.go`. |
| `BLOCO-03` | `unit_task` | Motor de Hashing & Validação (Go) | Implementação em Go do gerador de hash SHA-256 encadeado e validador de concorrência via TDD. |
| `BLOCO-04` | `script` | `go test -count=1 ./...` | Regressão mecânica e integridade de toda a suíte. |
| `BLOCO-05` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis

* **Gatilho Operacional de Disparo (Trigger):** `manual` (via `legoloop fabricar`, `legoloop orchestrate` ou pipelines da Tríade).
* **Nível de Risco & Blast Radius:** `Baixo (workspace local)` — alterações circunscritas a templates e motores internos de planejamento/validação.
* **Política de Privacidade & Orçamento Cognitivo:** Nuvem autorizada. Sem dados sensíveis ou segredos expostos.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Teste unitário comprovando que rotas geradas contêm os hashes exatos dos upstreams e que a rota aborta com erro binário se qualquer arquivo upstream for adulterado.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  1. A validação de hash dos upstreams deve ser tolerante a finais de linha (`CRLF` vs. `LF`), normalizando antes de calcular o SHA-256 para não quebrar entre Windows e Linux.
  2. O teto de subagentes paralelos padrão será `2`, configurável pelo usuário ou pelo harness.
* **Riscos e Perguntas em Aberto:** Nenhum risco impeditivo identificado.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test -count=1 ./...` com exit code 0 |
| 2 | Templates canônicos e schemas atualizados passam nos testes de conformidade. | Execução de teste de templates | `go test -v ./docs/templates/...` com exit code 0 |
| 3 | Cálculo e verificação de SHA-256 do quarteto validado em testes unitários. | Testes em `internal/plangen` e `internal/validator` | Teste de mutação de arquivo falhando com erro de integridade |
| 4 | Funções Go enxutas e conformes à regra de governança. | Linter de tamanho de função | `go run ./tools/funclen internal` (nenhuma função > 35 linhas) |
| 5 | Binário local recompilado com sucesso. | Build do executável | `go build -o bin/legoloop.exe ./cmd/legoloop` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-07_SPEC_ROUTE_V2_GOVERNANCE.md`](../specs/2026-11-07_SPEC_ROUTE_V2_GOVERNANCE.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 98% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
