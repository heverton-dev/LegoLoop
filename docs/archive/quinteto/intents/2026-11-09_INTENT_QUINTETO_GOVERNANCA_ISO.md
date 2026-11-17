# 🎯 INTENT: Quinteto de Governança ISO no LegoLoop

> **Status:** APPROVED (Atualizado com Paridade Total Cognitiva)  
> **Data:** 2026-11-09  
> **Autor:** LegoLoop Team & Antigravity  
> **Iniciativa:** I05 — Governança Pentafásica e Harmonia de Grafos Executáveis  
> **Alvo:** `internal/processo/`, `internal/fabricar/`, `internal/compor/`, `internal/fechar/`, `internal/dag/`, `internal/videorea/`, `docs/templates/`, `loops/`, `docs/dags/`  
> **Sessão:** 17bc79c3-324a-4324-88ca-3914b46981e6  
> **Confiança da Entrevista:** 99% (validado em auditoria de integridade com o operador)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Existência de grave dissonância estrutural entre o planejamento em rotas (`ROUTE.json`) e os grafos executáveis de tempo de execução (`pipeline.dag.json`). Faltava o 5º elemento documental arquivado cronologicamente em `docs/dags/YYYY-MM-DD_DAG_<SLUG>.json`. Além disso, identificou-se o risco crítico da **"Casca Oca" (Santo do Pau Oco)**: se o DAG for gerado apenas com comandos mecânicos de teste de sistema operacional (`go test`), ele perde toda a riqueza semântica e cognitiva da ROUTE (`context_in`, `context_out`, `test_seam`, `token_budget`, `isolation`, `tooling`), tornando impossível para um agente ou motor de inferência (CLI/Studio/LLM) reproduzir ou executar o pipeline autonomamente.
* **Quem Sofre:** Desenvolvedores, agentes autônomos em contexto estéril (Orchestrix) e o motor de execução `legoloop run`, que falha ou opera podado ao rodar grafos sem as amarras cognitivas originais.
* **O Porquê:** O LegoLoop exige o Princípio da Harmonia Estrita. A Skill (prosa), a Rota (plano de voo) e o DAG (grafo executável) devem ser espelhos fiéis e completos. O DAG não pode ser uma versão podada ou mutilada da ROUTE; deve carregar paridade integral de contratos e metadados.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** A governança formal em `docs/` parava no `ROUTE.json`. Quando projetado, o DAG sofria perda de parâmetros cognitivos, virando um mero runner de terminal sem instruções de inferência.
  * **Amanhã:** O **Quinteto de Governança ISO** (`[INTENT -> SPEC -> PLAN -> ROUTE -> DAG]`) é lei em todo o ecossistema com **Paridade Total de Metadados**. O arquivo `docs/dags/YYYY-MM-DD_DAG_<SLUG>.json` é gerado atomicamente contendo todos os contratos de isolamento, tokens, tooling, costuras de teste e contexto da ROUTE, projetado para `loops/` com validação de hash SHA-256 no fechamento.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [x] Objetivo 1: Instituir o 5º elemento `docs/dags/YYYY-MM-DD_DAG_<SLUG>.json` com **Paridade Total de Contratos Cognitivos** em relação a `ROUTE.json`, eliminando cascas ocas e sustentando execução autônoma integral por agentes e CLI.
* [x] Objetivo 2: Atualizar os motores nucleares em Go (`internal/dag/`, `internal/processo/`, `internal/fabricar/`, `internal/compor/`, `internal/fechar/`, `internal/videorea/`) para gerar, projetar, interpretar e auditar paridade SHA-256 do Quinteto.
* [x] Objetivo 3: Oficializar `TEMPLATE_DAG.json` estendido com paridade de schema na suíte de testes em `docs/templates/templates_test.go` e atualizar as leis do `AGENTS.md` (Regras #6 e #10).
* [x] Objetivo 4: Definir a Skill Canônica Universal como a manifestação em prosa estrita do Quinteto, auditada por linter de nós.
* [x] Objetivo 5: Sanear ativamente todos os loops do acervo (`video-rea`, `video-motion`, `video-edition`, `fechar`, etc.), eliminando nós defeituosos e padronizando os nomes sob o prefixo canônico compulsório `legoloop-*`.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto (stdlib Go pura).
* ⛔ Não utilizaremos engenharia reversa via REA MCP sobre código Go aberto: a auditoria é via AST e compilador nativo Go.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/aidd-spec` | Elaborar especificação formal com critérios BDD para o Quinteto e Paridade Total. |
| `BLOCO-02` | `unit_task` | `internal/dag`, `internal/processo`, `internal/fabricar` | Suporte a metadados cognitivos no DAG, gerador e projetor do 5º elemento com TDD. |
| `BLOCO-03` | `unit_task` | `internal/compor`, `internal/fechar` | Ajuste da recomposição e auditoria mecânica de paridade SHA-256. |
| `BLOCO-04` | `unit_task` | `docs/templates`, `loops/` | Inclusão de `TEMPLATE_DAG.json` na suíte e saneamento dos loops `legoloop-*`. |
| `BLOCO-05` | `script` | `go test ./...` | Regressão mecânica e integridade total (100% verde). |
| `BLOCO-06` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (invocado via CLI `legoloop fabricar`, chat `/legoloop-fabricar` ou Studio).
* **Nível de Risco & Blast Radius:** Médio (afeta os motores centrais de geração de arquivos, schema do DAG e persistência de loops).
* **Política de Privacidade & Orçamento Cognitivo:** Exclusivamente local, stdlib Go pura. Modelos cognitivos roteados por complexidade (`code-fast` / `flash` para geração mecânica; `sonnet` para inferência de contratos).
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**  
  1. Diretório `docs/dags/` existente e povoado com arquivos ISO válidos e ricos em parâmetros cognitivos.  
  2. Suíte de testes `go test -count=1 ./...` passando com zero falhas.  
  3. `legoloop fechar` auditando paridade mecânica entre `docs/dags/` e `loops/`.  
  4. Loop `legoloop-video-rea` saneado sem URLs *hardcoded*.
* **Premissas Assumidas:**  
  * O schema `dag.v2.json` / `dag.v1.json` absorve ou mapeia transparentemente todos os blocos contratuais da ROUTE sem quebra de retrocompatibilidade.
  * A paridade entre `ROUTE.json` e `DAG.json` é auditada deterministicamente por equivalência de metadados e SHA-256 no fechamento.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Paridade Total ROUTE ↔ DAG comprovada sem omissão de parâmetros cognitivos. | Comparação estrutural dos nós de rota e DAG | Inspeção e testes automatizados de parse |
| 3 | Conformidade do Quinteto nos templates automatizada. | Teste de templates aprovando os 5 artefatos | `go test -v ./docs/templates -run TestTemplatesConformidade` |
| 4 | Emissão automática de `docs/dags/YYYY-MM-DD_DAG_<SLUG>.json` na fabricação e composição. | Testes de unidade em `internal/fabricar` e `internal/compor` | `go test -v ./internal/fabricar ./internal/compor` |
| 5 | Auditoria de paridade SHA-256 no fechamento. | Teste de fechamento rejeitando divergência DAG vs loop | `go test -v ./internal/fechar -run TestFecharDAG` |
| 6 | Normalização de nomenclatura `legoloop-*` em todo o acervo de loops. | Listagem de loops e validação mecânica | `legoloop listar` exibindo prefixo uniforme |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-09_SPEC_QUINTETO_GOVERNANCA_ISO.md`](../specs/2026-11-09_SPEC_QUINTETO_GOVERNANCA_ISO.md)
* **Gate de Transição:** Este Intent foi aprovado e refinado pelo operador humano para assegurar a Paridade Total dos Contratos Cognitivos.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*Status: APROVADO pelo operador em 2026-11-09.*
