# 🗺️ PLAN: Quinteto de Governança ISO e Harmonização de Grafos Executáveis

> **Status:** APPROVED (Atualizado com Paridade Total Cognitiva)  
> **Data:** 2026-11-09  
> **Versão:** 1.1.0  
> **Autor:** LegoLoop Core Team & Antigravity  
> **Iniciativa:** I05 — Governança Pentafásica e Harmonia de Grafos Executáveis  
> **Referência da Intenção:** [docs/intents/2026-11-09_INTENT_QUINTETO_GOVERNANCA_ISO.md](../intents/2026-11-09_INTENT_QUINTETO_GOVERNANCA_ISO.md)  
> **Referência da Especificação:** [docs/specs/2026-11-09_SPEC_QUINTETO_GOVERNANCA_ISO.md](../specs/2026-11-09_SPEC_QUINTETO_GOVERNANCA_ISO.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-09_ROUTE_QUINTETO_GOVERNANCA_ISO.json](../routes/2026-11-09_ROUTE_QUINTETO_GOVERNANCA_ISO.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (AGENTS.md, regra 4: agentes de harness proibidos; execução estéril)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `1` (execução sequencial segura de integridade)
* **Timeout Padrão por Bloco:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `discard_uncommitted`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Núcleo de Projeção & Fábrica | 3 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | Composição com Quinteto | 1 | `AC-E02.1` | `nano` |
| **E03** | Motor Fechar & Auditoria SHA-256 | 1 | `AC-E03.1`, `AC-E03.2` | `micro` |
| **E04** | Templates de Governança & Regras | 1 | `AC-E04.1` | `nano` |
| **E05** | Saneamento do Acervo de Loops | 2 | `AC-E05.1`, `AC-E05.2` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Núcleo de Projeção & Fábrica (`internal/processo`, `internal/fabricar`, `internal/dag`)

#### Ticket TICK-01.1: Suporte a `DAGNome` e `CheckParity` em `internal/processo`
* **Nome Canônico:** `processo_projecoes_dagnome`
* **Natureza:** `micro`
* **Complexidade Cognitiva:** `simples`
* **Roteamento de Modelo:** `code-fast` (fallback: `flash`)
* **Dependências:** Nenhuma (`[]`)
* **Entradas de Contexto:** `docs/specs/2026-11-09_SPEC_QUINTETO_GOVERNANCA_ISO.md`
* **Saídas Contratadas:** `internal/processo/projecoes.go`, `internal/processo/projecoes_test.go`
* **Costura de Teste:** `TestProjetarEmiteDAGNome` e `TestCheckParity`
* **Comando de Verificação:** `go test -v ./internal/processo -run "TestProjetarEmiteDAGNome|TestCheckParity"`
* **Critério Coberto:** `AC-E01.1`
* **Stop-Gate:** `null`

#### Ticket TICK-01.2: Persistência de `docs/dags/` na esteira `fabricar`
* **Nome Canônico:** `fabricar_persist_dag_iso`
* **Natureza:** `micro`
* **Complexidade Cognitiva:** `simples`
* **Roteamento de Modelo:** `code-fast` (fallback: `flash`)
* **Dependências:** `["TICK-01.1"]`
* **Entradas de Contexto:** `internal/processo/projecoes.go`
* **Saídas Contratadas:** `internal/fabricar/camada4.go`, `internal/fabricar/materialize.go`, `internal/fabricar/materialize_test.go`
* **Costura de Teste:** `TestMaterialize_PersisteDAGIso`
* **Comando de Verificação:** `go test -v ./internal/fabricar -run TestMaterialize_PersisteDAGIso`
* **Critério Coberto:** `AC-E01.2`
* **Stop-Gate:** `null`

#### Ticket TICK-01.3: Paridade Total de Metadados Cognitivos ROUTE ↔ DAG
* **Nome Canônico:** `dag_paridade_cognitiva_route`
* **Natureza:** `micro`
* **Complexidade Cognitiva:** `moderada`
* **Roteamento de Modelo:** `code-fast` (fallback: `flash`)
* **Dependências:** `["TICK-01.2"]`
* **Entradas de Contexto:** `internal/dag/types.go`, `docs/routes/2026-11-09_ROUTE_QUINTETO_GOVERNANCA_ISO.json`
* **Saídas Contratadas:** `internal/dag/types.go`, `docs/templates/TEMPLATE_DAG.json`, `docs/dags/2026-11-09_DAG_QUINTETO_GOVERNANCA_ISO.json`
* **Costura de Teste:** `TestDAG_PreservaContratosCognitivos`
* **Comando de Verificação:** `go test -v ./internal/dag -run TestParse`
* **Critério Coberto:** `AC-E01.3`
* **Stop-Gate:** `null`

---

### Épico 02 — Composição com Quinteto (`internal/compor`)

#### Ticket TICK-02.1: Emissão de DAG ISO durante a composição de sub-loops
* **Nome Canônico:** `compor_regera_dag_iso`
* **Natureza:** `nano`
* **Complexidade Cognitiva:** `simples`
* **Roteamento de Modelo:** `code-fast` (fallback: `flash`)
* **Dependências:** `["TICK-01.2"]`
* **Entradas de Contexto:** `internal/compor/compor.go`
* **Saídas Contratadas:** `internal/compor/compor.go`, `internal/compor/compor_test.go`
* **Costura de Teste:** `TestCompor_RegeraDAGIso`
* **Comando de Verificação:** `go test -v ./internal/compor -run TestCompor_RegeraDAGIso`
* **Critério Coberto:** `AC-E02.1`
* **Stop-Gate:** `null`

---

### Épico 03 — Motor Fechar & Auditoria SHA-256 (`internal/fechar`)

#### Ticket TICK-03.1: Auditoria mecânica de paridade SHA-256 no `fecharDAG`
* **Nome Canônico:** `fechar_audit_dag_parity`
* **Natureza:** `micro`
* **Complexidade Cognitiva:** `moderada`
* **Roteamento de Modelo:** `code-fast` (fallback: `flash`)
* **Dependências:** `["TICK-01.1"]`
* **Entradas de Contexto:** `internal/fechar/fechar.go`
* **Saídas Contratadas:** `internal/fechar/fechar.go`, `internal/fechar/fechar_test.go`
* **Costura de Teste:** `TestFecharDAG_AuditaParidadeSHA256`
* **Comando de Verificação:** `go test -v ./internal/fechar -run TestFecharDAG_AuditaParidadeSHA256`
* **Critério Coberto:** `AC-E03.1`, `AC-E03.2`
* **Stop-Gate:** `null`

---

### Épico 04 — Templates de Governança & Regras (`docs/templates/`, `AGENTS.md`)

#### Ticket TICK-04.1: Conformidade do Quinteto nos templates e governança viva
* **Nome Canônico:** `templates_conformidade_quinteto`
* **Natureza:** `nano`
* **Complexidade Cognitiva:** `simples`
* **Roteamento de Modelo:** `code-fast` (fallback: `flash`)
* **Dependências:** `["TICK-01.1"]`
* **Entradas de Contexto:** `docs/templates/templates.go`, `AGENTS.md`
* **Saídas Contratadas:** `docs/templates/templates.go`, `docs/templates/templates_test.go`, `AGENTS.md`
* **Costura de Teste:** `TestTemplatesConformidadeQuinteto`
* **Comando de Verificação:** `go test -v ./docs/templates -run TestTemplatesConformidade`
* **Critério Coberto:** `AC-E04.1`
* **Stop-Gate:** `null`

---

### Épico 05 — Saneamento do Acervo de Loops (`loops/`, `docs/dags/`)

#### Ticket TICK-05.1: Saneamento e emissão de DAG ISO para `legoloop-video-rea`
* **Nome Canônico:** `saneamento_video_rea_dag`
* **Natureza:** `micro`
* **Complexidade Cognitiva:** `simples`
* **Roteamento de Modelo:** `code-fast` (fallback: `flash`)
* **Dependências:** `["TICK-01.2", "TICK-03.1"]`
* **Entradas de Contexto:** `loops/video-rea/pipeline.dag.json`
* **Saídas Contratadas:** `loops/legoloop-video-rea/pipeline.dag.json`, `docs/dags/2026-11-09_DAG_LEGOLOOP_VIDEO_REA.json`
* **Costura de Teste:** `TestVideoREADAG_SemURLHardcoded`
* **Comando de Verificação:** `go test -v ./internal/videorea -run TestDetectPlatform`
* **Critério Coberto:** `AC-E05.1`
* **Stop-Gate:** `null`

#### Ticket TICK-05.2: Normalização de prefixo `legoloop-*` no acervo de loops
* **Nome Canônico:** `normalizacao_prefixo_loops`
* **Natureza:** `micro`
* **Complexidade Cognitiva:** `simples`
* **Roteamento de Modelo:** `code-fast` (fallback: `flash`)
* **Dependências:** `["TICK-05.1"]`
* **Entradas de Contexto:** `loops/`
* **Saídas Contratadas:** `loops/legoloop-video-motion/`, `loops/legoloop-video-edition/`, `loops/legoloop-fechar/`, `docs/dags/`
* **Costura de Teste:** `TestAcervoLoops_PrefixoCanonico`
* **Comando de Verificação:** `go test -v ./internal/loop -run TestListarLoops`
* **Critério Coberto:** `AC-E05.2`
* **Stop-Gate:** `null`
