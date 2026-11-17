---
id: INTENT-fabricar-web-qa
status: UNDER_REVIEW
date: 2026-11-15
schema_version: 2
upstream: [INTERVIEW-fabricar-web-qa]
covers: [Q1]
---

# 🎯 INTENT: Processo web-qa de auditoria deterministica de UI/UX, acessibilidade, links, rotas e regressao visual para web apps com playwright, obscura, aidd-linter-design, renderlint, axe-core, unlighthouse, lychee e lost-pixel

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-15  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** fabricar-web-qa  
> **Confiança da Entrevista:** 97% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Processo web-qa de auditoria deterministica de UI/UX, acessibilidade, links, rotas e regressao visual para web apps com playwright, obscura, aidd-linter-design, renderlint, axe-core, unlighthouse, lychee e lost-pixel
* **Quem Sofre:** Desenvolvedores, operadores e fluxos autônomos do LegoLoop que auditam web apps e convertem achados de QA em correções governadas.
* **O Porquê:** Garantir auditoria determinística, reprodutível e sem stubs/mocks de UI/UX, acessibilidade, links, rotas, Lighthouse e regressão visual antes da execução de correções em código.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista: Processo web-qa de auditoria deterministica de UI/UX, acessibilidade, links, rotas e regressao visual para web apps com playwright, obscura, aidd-linter-design, renderlint, axe-core, unlighthouse, lychee e lost-pixel
  * **Amanhã:** Processo canônico `web-qa` no LegoLoop (Skill + DAG executável por CLI/Orchestrix) para auditoria determinística de web apps com `playwright`, `obscura`, `aidd-linter-design`, `renderlint`, `axe-core`, `unlighthouse`, `lychee` e `lost-pixel`, gerando relatório consolidado (JSON + Markdown) que serve obrigatoriamente de insumo para entrevista socrática de correção e geração dos artefatos `INTERVIEW`, `INTENT`, `SPEC`, `PLAN`, `ROUTE`, `DAG` e `SWARM` para ciclo de `BUILD`.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Processo canônico `web-qa` no LegoLoop (Skill + DAG executável por CLI/Orchestrix) para auditoria determinística de web apps com `playwright`, `obscura`, `aidd-linter-design`, `renderlint`, `axe-core`, `unlighthouse`, `lychee` e `lost-pixel`, gerando relatório consolidado (JSON + Markdown) que serve obrigatoriamente de insumo para entrevista socrática de correção e geração dos artefatos `INTERVIEW`, `INTENT`, `SPEC`, `PLAN`, `ROUTE`, `DAG` e `SWARM` para ciclo de `BUILD`.
* [ ] Objetivo 2: Preflight bloqueia execução se faltar qualquer das 8 ferramentas; sobe app local via comando+porta ou audita URL alvo; descobre rotas via crawl/sitemap + lista explícita opcional com `storageState` opcional do `playwright` em desktop (`1440x900`) e mobile (`390x844`); reprova com exit code não-zero por link quebrado (`lychee`), erro de rota/console (`playwright`), violação `critical`/`serious` (`axe-core`), erro de layout/design/botão morto (`renderlint`, `aidd-linter-design`, `obscura`), nota abaixo de 90 (`unlighthouse`) ou diff visual (`lost-pixel`); sem baseline no `lost-pixel`, cria snapshot inicial e exige Stop-Gate humano para promover; grava relatório JSON + Markdown em `docs/qa-web/YYYY-MM-DD_<alvo>/` para iniciar a entrevista socrática e cadeia até `BUILD`.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Edição automática cega de código sem entrevista socrática e artefatos completos de governança (`INTERVIEW` até `SWARM`/`BUILD`); testes de carga; testes de penetração.

---

## 3. Mapeamento Preliminar de Peças
Identificação das Skills que o futuro Processo precisará acionar. Tipos permitidos: os tipos de nó do DAG e do ROUTE, descritos no README, seção 6.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | `skills/aidd-spec` | Elaborar a especificação formal com critérios BDD. |
| `PECA-02` | `unit_task` | Código de produção | Implementação da lógica via TDD Red-Green. |
| `PECA-03` | `script` | `go test ./...` | Regressão mecânica e integridade. |
| `PECA-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`  
  *(Define como este Processo acordará no dia a dia quando estiver construído.)*
* **Nível de Risco & Blast Radius:** A definir no Stop-Gate 01
* **Política de Privacidade & Orçamento Cognitivo:** A definir no Stop-Gate 01
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Preflight bloqueia execução se faltar qualquer das 8 ferramentas; sobe app local via comando+porta ou audita URL alvo; descobre rotas via crawl/sitemap + lista explícita opcional com `storageState` opcional do `playwright` em desktop (`1440x900`) e mobile (`390x844`); reprova com exit code não-zero por link quebrado (`lychee`), erro de rota/console (`playwright`), violação `critical`/`serious` (`axe-core`), erro de layout/design/botão morto (`renderlint`, `aidd-linter-design`, `obscura`), nota abaixo de 90 (`unlighthouse`) ou diff visual (`lost-pixel`); sem baseline no `lost-pixel`, cria snapshot inicial e exige Stop-Gate humano para promover; grava relatório JSON + Markdown em `docs/qa-web/YYYY-MM-DD_<alvo>/` para iniciar a entrevista socrática e cadeia até `BUILD`.
* **Restrições Inegociáveis:**
  * Conformidade integral ao `AGENTS.md` (Zero Stubs/Mocks
  * ferramentas reais obrigatórias
  * Graph-First via `codebase-memory-mcp`
  * Orchestrix cirúrgica sem subagentes monolíticos
  * execução estrita via `docs/routes/`).
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.
* **Riscos e Perguntas em Aberto:** Confiança de 97%: as lacunas restantes devem ser tratadas no Stop-Gate 01.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |
| 3 | Preflight bloqueia execução se faltar qualquer das 8 ferramentas; sobe app local via comando+porta ou audita URL alvo; descobre rotas via crawl/sitemap + lista explícita opcional com `storageState` opcional do `playwright` em desktop (`1440x900`) e mobile (`390x844`); reprova com exit code não-zero por link quebrado (`lychee`), erro de rota/console (`playwright`), violação `critical`/`serious` (`axe-core`), erro de layout/design/botão morto (`renderlint`, `aidd-linter-design`, `obscura`), nota abaixo de 90 (`unlighthouse`) ou diff visual (`lost-pixel`); sem baseline no `lost-pixel`, cria snapshot inicial e exige Stop-Gate humano para promover; grava relatório JSON + Markdown em `docs/qa-web/YYYY-MM-DD_<alvo>/` para iniciar a entrevista socrática e cadeia até `BUILD`. | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-15_SPEC_WEB_QA.md`](../specs/2026-11-15_SPEC_WEB_QA.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 97% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
