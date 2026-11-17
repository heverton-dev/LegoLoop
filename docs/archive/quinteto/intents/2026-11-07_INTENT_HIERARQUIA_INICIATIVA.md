# 🎯 INTENT: HIERARQUIA_INICIATIVA_LEGOLOOP

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-07  
> **Autor:** LegoLoop Mega Pipeline & Antigravity  
> **Alvo:** docs/templates, internal/plangen, internal/specgen, skills/aidd-legoloop-mega-pipeline  
> **Sessão:** 8b70e2fa-9b0e-4604-9f78-939da83d0d55  
> **Confiança da Entrevista:** 98% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** A matriz atual de granularidade do LegoLoop opera com 3 níveis (Épico ➔ Micro-ticket ➔ Nano-ticket). Para demandas de grande porte corporativo ou multi-sistemas (ex.: mudanças conjuntas em Binário Go + Studio Web + CLI + Infra), a falta de um nível agregador superior (Iniciativa) força a criação de dezenas de épicos planos ou demandas fragmentadas sem visão sistêmica unificada.
* **Quem Sofre:** Arquitetos de software, orquestradores de pipelines complexos e desenvolvedores que precisam coordenar múltiplos subsistemas dependentes em uma mesma esteira de entrega.
* **O Porquê:** Introduzir o nível canônico **Iniciativa (`I01`, `I02`...)** acima do **Épico (`E01`, `E02`...)**, de forma estritamente opcional e com 100% de retrocompatibilidade, permite estruturar portfólios e mega-demandas sem quebrar os fluxos menores já existentes.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** O nível máximo de agrupamento é o Épico (`E01`). Demandas amplas precisam ser divididas artificialmente ou saturar o escopo de um único épico.
  * **Amanhã:** O sistema suporta formalmente a hierarquia de 4 níveis: **Iniciativa (`I01`) ➔ Épico (`E01`) ➔ Micro-ticket (`TICK-01.1`) ➔ Nano-ticket (`TICK-01.2`)**. Em demandas simples, a Iniciativa é opcional (default transparente); em demandas complexas, permite governança por frentes de engenharia.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Padronizar a taxonomia canônica oficial do nível superior como **Iniciativa** (`I01`, `I02`...).
* [ ] Objetivo 2: Atualizar os templates canônicos SSOT (`docs/templates/TEMPLATE_INTENT.md`, `TEMPLATE_SPEC.md`, `TEMPLATE_PLAN.md`, `TEMPLATE_ROUTE.json`) para suportar a hierarquia com Iniciativas como campo opcional/retrocompatível.
* [ ] Objetivo 3: Adaptar os motores Go (`internal/plangen` e geradores correlatos) para processar Iniciativas sem quebrar rotas e planos baseados puramente em Épicos.
* [ ] Objetivo 4: Atualizar a documentação normativa e a skill mestre (`skills/aidd-legoloop-mega-pipeline/SKILL.md` e `AGENTS.md`) com a matriz de 4 níveis.
* [ ] Objetivo 5: Garantir integridade de 100% dos testes unitários existentes e conformidade de linters (`go test ./...`, `go vet ./...`, `funclen`).

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não tornaremos o campo Iniciativa obrigatório em demandas antigas ou simples (retrocompatibilidade inegociável).
* ⛔ Não renomearemos nem removeremos o termo "Épico" ou "Micro-ticket".
* ⛔ Não adicionaremos camadas burocráticas intermediárias (como "Feature" ou "História").
* ⛔ Não alteraremos o mecanismo básico de despacho de execução isolada (Orchestrix hermética).

---

## 3. Mapeamento Preliminar de Blocos e Subloops
Identificação das capacidades que a futura esteira precisará acionar:

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | aidd-spec | Elaborar especificação formal (SPEC) com critérios BDD para a taxonomia e parsing de Iniciativas. |
| `BLOCO-02` | `skill` | aidd-plan | Gerar o PLAN estruturado e ROUTE.json correspondente com os tickets de engenharia. |
| `BLOCO-03` | `unit_task` | Código de templates e Go | Atualização de `docs/templates/` e evolução dos parsers/montadores em `internal/plangen/`. |
| `BLOCO-04` | `script` | `go test ./...` | Validação de regressão mecânica e integridade de todos os pacotes. |
| `BLOCO-05` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória do binário do LegoLoop. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (solicitação explícita via chat ou comando `/legoloop`).
* **Nível de Risco & Blast Radius:** `Baixo (workspace local)` — As alterações são cumulativas e retrocompatíveis; testes automatizados validam cada nó.
* **Política de Privacidade & Orçamento Cognitivo:** `Nuvem autorizada` · Teto: modelos de alto raciocínio para arquitetura/SPEC e modelos rápidos para patches.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  1. Suíte de testes `go test -count=1 ./...` executando 100% verde.
  2. Testes de `docs/templates/templates_test.go` passando sem divergência de embeds.
  3. Binário `bin/legoloop.exe` recompilado e apto a processar planos com e sem Iniciativas.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  1. O termo eleito é **Iniciativa** (`I01`), sem conflitos de vocabulário.
  2. Iniciativa é um metadado opcional: planos legados sem `I01` continuam sendo válidos e executáveis.
* **Riscos e Perguntas em Aberto:** Nenhum risco impeditivo identificado.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo em `plangen` e `templates_test.go` validando Iniciativas | `go test -count=1 ./internal/plangen/... ./docs/templates/...` |
| 3 | Binário compilado sem avisos ou erros. | Geração de artefato executável fresco | `go build -o bin/legoloop.exe ./cmd/legoloop` |
| 4 | Templates canônicos atualizados e consistentes com o schema. | Teste de igualdade de embeds | `go test -v -run TestEmbeddedTemplatesMatchFiles ./docs/templates/...` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md`](../specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 98% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
