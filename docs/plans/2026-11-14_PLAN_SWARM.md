---
id: PLAN-2026-11-14_ROUTE_SWARM
status: DRAFT
date: 2026-11-14
schema_version: 2
upstream: [2026-11-14_SPEC_SWARM]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E01.4, AC-E01.5]
---

# 🗺️ PLAN: legoloop-swarm

> **Status:** Em Planejamento  
> **Data:** 2026-11-14  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** Não informada na SPEC  
> **Referência da Especificação:** [docs/specs/2026-11-14_SPEC_SWARM.md](../specs/2026-11-14_SPEC_SWARM.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-14_ROUTE_SWARM.json](../routes/2026-11-14_ROUTE_SWARM.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Peça, sem context bloat)
* **Nível de Supervisão:** `epic`
* **Permissão de Subagentes:** `false` ou `ephemeral` (AGENTS.md, regra 4: agente de harness proibido; subagente efêmero só com as seis condições)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Peça:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | legoloop-swarm | 5 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3`, `AC-E01.4`, `AC-E01.5` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — legoloop-swarm

#### Ticket `TICK-01.1` — `e01` — Etapa 1 — peça `legoloop-piece-upgrader`
* **Iniciativa:** `nenhum` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `out/01_etapa.md` · Spec Anchor `docs/specs/2026-11-14_SPEC_SWARM.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `out/01_etapa.md`
* **Verificação (argv, sem shell):** `legoloop validate ~/.aidd/canonical/loops/legoloop-swarm.dag.json`
* **Gate (DSL):** `file_exists('out/01_etapa.md') && word_count('out/01_etapa.md') >= 30 && !has_pattern('out/01_etapa.md', '\{\{')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — `e02` — Etapa 2 — peça `aidd-piece-upgrader`
* **Iniciativa:** `nenhum` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `out/02_etapa.md` · Spec Anchor `docs/specs/2026-11-14_SPEC_SWARM.md#AC-E01.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `out/02_etapa.md`
* **Verificação (argv, sem shell):** `legoloop validate ~/.aidd/canonical/loops/legoloop-swarm.dag.json`
* **Gate (DSL):** `file_exists('out/02_etapa.md') && word_count('out/02_etapa.md') >= 30 && !has_pattern('out/02_etapa.md', '\{\{')`
* **Spec Ref:** `AC-E01.2`

#### Ticket `TICK-01.3` — `e03` — Etapa 3 — peça `legoloop-mega-pipeline`
* **Iniciativa:** `nenhum` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `out/03_etapa.md` · Spec Anchor `docs/specs/2026-11-14_SPEC_SWARM.md#AC-E01.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `out/03_etapa.md`
* **Verificação (argv, sem shell):** `legoloop validate ~/.aidd/canonical/loops/legoloop-swarm.dag.json`
* **Gate (DSL):** `file_exists('out/03_etapa.md') && word_count('out/03_etapa.md') >= 30 && !has_pattern('out/03_etapa.md', '\{\{')`
* **Spec Ref:** `AC-E01.3`

#### Ticket `TICK-01.4` — `e04` — Etapa 4 — peça `aidd-legoloop-mega-pipeline`
* **Iniciativa:** `nenhum` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `out/04_etapa.md` · Spec Anchor `docs/specs/2026-11-14_SPEC_SWARM.md#AC-E01.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `out/04_etapa.md`
* **Verificação (argv, sem shell):** `legoloop validate ~/.aidd/canonical/loops/legoloop-swarm.dag.json`
* **Gate (DSL):** `file_exists('out/04_etapa.md') && word_count('out/04_etapa.md') >= 30 && !has_pattern('out/04_etapa.md', '\{\{')`
* **Spec Ref:** `AC-E01.4`

#### Ticket `TICK-01.5` — `e05` — Etapa 5 — peça `legoloop-video-master`
* **Iniciativa:** `nenhum` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.4"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `out/05_etapa.md` · Spec Anchor `docs/specs/2026-11-14_SPEC_SWARM.md#AC-E01.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `out/05_etapa.md`
* **Verificação (argv, sem shell):** `legoloop validate ~/.aidd/canonical/loops/legoloop-swarm.dag.json`
* **Gate (DSL):** `file_exists('out/05_etapa.md') && word_count('out/05_etapa.md') >= 30 && !has_pattern('out/05_etapa.md', '\{\{')`
* **Spec Ref:** `AC-E01.5`

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test ./...`
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia). `golangci-lint run ./...` apenas se instalado e configurado.
3. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`
4. Emissão do relatório formal em `docs/reports/`
5. Registro de histórico da sessão via `faz-session --id <ID> --format all`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*Plano de execução fatiado com 5 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*
