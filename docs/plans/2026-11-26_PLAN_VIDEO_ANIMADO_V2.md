---
id: PLAN-2026-11-26_ROUTE_VIDEO_ANIMADO_V2
status: APPROVED
date: 2026-11-26
schema_version: 2
upstream: [2026-11-26_SPEC_VIDEO_ANIMADO_V2]
covers: [AC-E07.1, AC-E07.2, AC-E07.3, AC-E07.4, AC-E07.5, AC-E07.6, AC-E07.7, AC-E07.8, AC-E07.9, AC-E08.1, AC-E08.2, AC-E08.3, AC-E09.1, AC-E09.2, AC-E09.3, AC-E09.4, AC-E09.5, AC-E10.1, AC-E10.2, AC-E10.3, AC-E10.4, AC-E10.5, AC-E10.6, AC-E10.7, AC-E10.8, AC-E11.1, AC-E11.2, AC-E11.3, AC-E11.4, AC-E12.1, AC-E12.2, AC-E12.3, AC-NF.1, AC-NF.2, AC-NF.3]
---

# 🗺️ PLAN: 2026-11-26_SPEC_VIDEO_ANIMADO_V2

> **Status:** Em Planejamento  
> **Data:** 2026-11-26  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-26_EMENDA_INTENT_VIDEO_ANIMADO.md](../intents/2026-11-26_EMENDA_INTENT_VIDEO_ANIMADO.md)  
> **Referência da Especificação:** [docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md](../specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-26_ROUTE_VIDEO_ANIMADO_V2.json](../routes/2026-11-26_ROUTE_VIDEO_ANIMADO_V2.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Peça, sem context bloat)
* **Nível de Supervisão:** `epic`
* **Permissão de Subagentes:** `false` ou `ephemeral` (AGENTS.md, regra 4: agente de harness proibido; subagente efêmero só com as seis condições)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Peça:** `180s`
* **Política de Rollback:** global `abort_processo`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E07** | Épico E07 | 6 | `AC-E07.2`, `AC-E07.6`, `AC-E07.1`, `AC-E07.3`, `AC-E07.4`, `AC-E07.5`, `AC-E07.7`, `AC-E07.8`, `AC-E07.8`, `AC-E07.9` | `micro` |
| **E08** | Épico E08 | 2 | `AC-E08.1`, `AC-E08.2`, `AC-E08.3` | `micro` |
| **E09** | Épico E09 | 4 | `AC-E09.1`, `AC-E09.2`, `AC-E09.3`, `AC-E09.4`, `AC-E09.5` | `micro` |
| **E10** | Épico E10 | 5 | `AC-E10.1`, `AC-E10.2`, `AC-E10.3`, `AC-E10.6`, `AC-E10.8`, `AC-E10.4`, `AC-E10.5`, `AC-E10.7`, `AC-E10.4` | `micro` |
| **E11** | Épico E11 | 4 | `AC-E11.2`, `AC-E11.1`, `AC-E11.3`, `AC-E11.4` | `micro` |
| **E12** | Épico E12 | 4 | `AC-E12.1`, `AC-E12.3`, `AC-E12.2`, `AC-NF.1`, `AC-NF.2`, `AC-NF.3` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 07 — Épico E07

#### Ticket `TICK-01.1` — Criar contratos do pacote inventario: Estilo, Item, Referencia, erros tipados e validacao de tipo
* **Iniciativa:** `V01` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/video/inventario/schemas.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E07.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/video/inventario/schemas_test.go`
* **Outputs:** Artifacts `internal/video/inventario/schemas.go, internal/video/inventario/schemas_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/inventario/ -run TestSchemas`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/video/inventario/schemas.go') && file_exists('internal/video/inventario/schemas_test.go')`
* **Spec Ref:** `AC-E07.2, AC-E07.6`

#### Ticket `TICK-01.2` — Implementar NovoEstilo e GravarItem com gravacao atomica e regra de estilo em rascunho
* **Iniciativa:** `V01` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/video/inventario/store.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E07.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/video/inventario/store_test.go`
* **Outputs:** Artifacts `internal/video/inventario/store.go, internal/video/inventario/store_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/inventario/ -run TestStore`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/video/inventario/store.go') && file_exists('internal/video/inventario/store_test.go')`
* **Spec Ref:** `AC-E07.1, AC-E07.3`

#### Ticket `TICK-01.3` — Implementar Buscar por estilo, tipo e tags, heranca de estilos e Verificar por hash SHA-256
* **Iniciativa:** `V01` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/video/inventario/busca.go, internal/video/inventario/hash.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E07.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/video/inventario/busca_test.go`
* **Outputs:** Artifacts `internal/video/inventario/busca.go, internal/video/inventario/busca_test.go, internal/video/inventario/hash.go, internal/video/inventario/hash_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/inventario/ -run TestBusca|TestHash|TestHeranca`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/video/inventario/busca.go') && file_exists('internal/video/inventario/busca_test.go') && file_exists('internal/video/inventario/hash.go') && file_exists('internal/video/inventario/hash_test.go')`
* **Spec Ref:** `AC-E07.4, AC-E07.5, AC-E07.7`

#### Ticket `TICK-01.4` — Implementar ValidarEstilo com arquivo de aprovacao humana e o comando video-inventario
* **Iniciativa:** `V01` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/video/inventario/validar.go, internal/cli/video_inventario.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E07.8` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/video/inventario/validar_test.go`
* **Outputs:** Artifacts `internal/video/inventario/validar.go, internal/video/inventario/validar_test.go, internal/cli/video_inventario.go, internal/cli/video_inventario_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/inventario/ ./internal/cli/ -run TestValidarEstilo|TestVideoInventario`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/video/inventario/validar.go') && file_exists('internal/video/inventario/validar_test.go') && file_exists('internal/cli/video_inventario.go') && file_exists('internal/cli/video_inventario_test.go')`
* **Spec Ref:** `AC-E07.8`

#### Ticket `TICK-02.1` — Criar o estilo molde-de-pilula com referencias, paleta, luz e personagens de teste para validacao do usuario
* **Iniciativa:** `V01` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.4"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `inventario/video/estilos/molde-de-pilula/estilo.json, inventario/video/estilos/molde-de-pilula/personagens.json` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E07.8` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `inventario/video/estilos/molde-de-pilula/estilo.json, inventario/video/estilos/molde-de-pilula/personagens.json`
* **Verificação (argv, sem shell):** `go run ./cmd/legoloop video-inventario verificar --estilo molde-de-pilula`
* **Gate (DSL):** `exit_code == 0 && file_exists('inventario/video/estilos/molde-de-pilula/estilo.json') && file_exists('inventario/video/estilos/molde-de-pilula/personagens.json')`
* **Spec Ref:** `AC-E07.8`

#### Ticket `TICK-02.2` — Adicionar genero e faixa etaria ao Item, BuscarPersonagens e tres personagens de teste (mulher, idoso, crianca)
* **Iniciativa:** `V01` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/video/inventario/schemas.go, internal/video/inventario/genero.go, inventario/video/estilos/molde-de-pilula/personagens.json` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E07.9` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/video/inventario/genero_test.go`
* **Outputs:** Artifacts `internal/video/inventario/schemas.go, internal/video/inventario/genero.go, internal/video/inventario/genero_test.go, inventario/video/estilos/molde-de-pilula/personagens.json`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/inventario/ -run TestGenero|TestBuscarPersonagens` · `go run ./cmd/legoloop video-inventario verificar --estilo molde-de-pilula`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/video/inventario/schemas.go') && file_exists('internal/video/inventario/genero.go') && file_exists('internal/video/inventario/genero_test.go') && file_exists('inventario/video/estilos/molde-de-pilula/personagens.json')`
* **Spec Ref:** `AC-E07.9`

### 🛑 STOP-GATE HUMANO (`SG-E07`)
* **Gate (DSL):** `prompt('Evidências do Épico 07 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 08 — Épico E08

#### Ticket `TICK-03.1` — Ler modo_producao e estilo no projeto.json v2 e bloquear estilo nao validado no pre-voo
* **Iniciativa:** `V01` | **Épico:** `E08` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2","SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/cenarios.go, internal/diorama/ativos/prevoo.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E08.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/cenarios_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/cenarios.go, internal/diorama/ativos/cenarios_test.go, internal/diorama/ativos/prevoo.go, internal/diorama/ativos/prevoo_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestPreVoo`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/cenarios.go') && file_exists('internal/diorama/ativos/cenarios_test.go') && file_exists('internal/diorama/ativos/prevoo.go') && file_exists('internal/diorama/ativos/prevoo_test.go')`
* **Spec Ref:** `AC-E08.1`

#### Ticket `TICK-03.2` — Validar cenas.json v2 e elenco.json contra o inventario e as partes suportadas do renderizador
* **Iniciativa:** `V01` | **Épico:** `E08` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1","SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/schemas.go, internal/diorama/ativos/storyboard.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E08.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/schemas_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/schemas.go, internal/diorama/ativos/schemas_test.go, internal/diorama/ativos/storyboard.go, internal/diorama/ativos/storyboard_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestStoryboard|TestElencoCast`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/schemas.go') && file_exists('internal/diorama/ativos/schemas_test.go') && file_exists('internal/diorama/ativos/storyboard.go') && file_exists('internal/diorama/ativos/storyboard_test.go')`
* **Spec Ref:** `AC-E08.2, AC-E08.3`

### 🛑 STOP-GATE HUMANO (`SG-E08`)
* **Gate (DSL):** `prompt('Evidências do Épico 08 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 09 — Épico E09

#### Ticket `TICK-04.1` — Etapa 3 vetorial: montar elenco.json com cast a partir do inventario, sem Flow, e gravar itens novos como nao aprovados
* **Iniciativa:** `V01` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/elenco.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E09.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/elenco_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/elenco.go, internal/diorama/ativos/elenco_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestElencoVetorial`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/elenco.go') && file_exists('internal/diorama/ativos/elenco_test.go')`
* **Spec Ref:** `AC-E09.1, AC-E09.2`

#### Ticket `TICK-04.2` — Manter o caminho de fichas e rembg para estilos de recorte sem regressao
* **Iniciativa:** `V01` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/fichas.go, internal/diorama/ativos/drivers.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E09.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/fichas_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/fichas.go, internal/diorama/ativos/drivers.go, internal/diorama/ativos/fichas_test.go, internal/diorama/ativos/drivers_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestFichas|TestDrivers|TestElenco`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/fichas.go') && file_exists('internal/diorama/ativos/drivers.go') && file_exists('internal/diorama/ativos/fichas_test.go') && file_exists('internal/diorama/ativos/drivers_test.go')`
* **Spec Ref:** `AC-E09.3`

#### Ticket `TICK-04.3` — Etapa 4: consultar o inventario antes de gerar fundos e registrar fundo novo com prompt, tags e hash
* **Iniciativa:** `V01` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/cenarios.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E09.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/cenarios_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/cenarios.go, internal/diorama/ativos/cenarios_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestCenariosInventario`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/cenarios.go') && file_exists('internal/diorama/ativos/cenarios_test.go')`
* **Spec Ref:** `AC-E09.4`

#### Ticket `TICK-04.4` — Marcar entrega suspeita do Flow (sem confirmacao de renomeacao, hash cruzado, fora do padrao de estilo)
* **Iniciativa:** `V01` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/suspeito.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E09.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/suspeito_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/suspeito.go, internal/diorama/ativos/suspeito_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestSuspeito`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/suspeito.go') && file_exists('internal/diorama/ativos/suspeito_test.go')`
* **Spec Ref:** `AC-E09.5`

### 🛑 STOP-GATE HUMANO (`SG-E09`)
* **Gate (DSL):** `prompt('Evidências do Épico 09 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 10 — Épico E10

#### Ticket `TICK-05.1` — Criar motor-local com LICENCA.json, ignorar no git e teste de governanca da licenca
* **Iniciativa:** `V01` | **Épico:** `E10` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2","SG-E09"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `.gitignore, motor-local/LICENCA.json` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E10.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/video/motor/licenca_test.go`
* **Outputs:** Artifacts `.gitignore, motor-local/LICENCA.json, internal/video/motor/licenca_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/motor/ -run TestLicenca`
* **Gate (DSL):** `exit_code == 0 && file_exists('.gitignore') && file_exists('motor-local/LICENCA.json') && file_exists('internal/video/motor/licenca_test.go')`
* **Spec Ref:** `AC-E10.1`

#### Ticket `TICK-05.2` — Implementar o pre-voo do motor: python, Playwright, Chromium, ffmpeg e motor-local, com ErrMotorAusente
* **Iniciativa:** `V01` | **Épico:** `E10` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1","SG-E09"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/video/motor/prevoo.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E10.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/video/motor/prevoo_test.go`
* **Outputs:** Artifacts `internal/video/motor/prevoo.go, internal/video/motor/prevoo_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/motor/ -run TestPreVoo`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/video/motor/prevoo.go') && file_exists('internal/video/motor/prevoo_test.go')`
* **Spec Ref:** `AC-E10.2`

#### Ticket `TICK-05.3` — Implementar a ponte cenas.json v2 + elenco + timings + inventario para timeline.json, com luz de contorno do estilo
* **Iniciativa:** `V01` | **Épico:** `E10` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.2","TICK-04.1","SG-E09"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/video/motor/ponte.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E10.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/video/motor/ponte_test.go`
* **Outputs:** Artifacts `internal/video/motor/ponte.go, internal/video/motor/ponte_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/motor/ -run TestPonte`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/video/motor/ponte.go') && file_exists('internal/video/motor/ponte_test.go')`
* **Spec Ref:** `AC-E10.3, AC-E10.6, AC-E10.8`

#### Ticket `TICK-05.4` — Renderizar a cena de referencia pelo motor: 9:16, voz, determinismo e boca acompanhando a fala
* **Iniciativa:** `V01` | **Épico:** `E10` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.3","SG-E09"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/video/motor/render.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E10.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/video/motor/render_test.go`
* **Outputs:** Artifacts `internal/video/motor/render.go, internal/video/motor/render_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/motor/ -run TestRender`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/video/motor/render.go') && file_exists('internal/video/motor/render_test.go')`
* **Spec Ref:** `AC-E10.4, AC-E10.5, AC-E10.7`

#### Ticket `TICK-06.1` — Ligar o render Go ao motor e reconfirmar voz, mix, SFX e cortes da SPEC v1.4.0 sobre a nova timeline
* **Iniciativa:** `V01` | **Épico:** `E10` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.4","SG-E09"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/render.go, internal/diorama/render/render_exec.go, internal/diorama/render/timeline.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E10.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/render_test.go`
* **Outputs:** Artifacts `internal/diorama/render/render.go, internal/diorama/render/render_exec.go, internal/diorama/render/timeline.go, internal/diorama/render/render_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/render.go') && file_exists('internal/diorama/render/render_exec.go') && file_exists('internal/diorama/render/timeline.go') && file_exists('internal/diorama/render/render_test.go')`
* **Spec Ref:** `AC-E10.4`

### 🛑 STOP-GATE HUMANO (`SG-E10`)
* **Gate (DSL):** `prompt('Evidências do Épico 10 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 11 — Épico E11

#### Ticket `TICK-07.1` — Mover internal/diorama para internal/video e atualizar imports sem regressao
* **Iniciativa:** `V01` | **Épico:** `E11` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.1","SG-E10"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/video/roteiro/, internal/video/ativos/, internal/video/render/` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E11.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `internal/video/roteiro/, internal/video/ativos/, internal/video/render/`
* **Verificação (argv, sem shell):** `go build ./...` · `go vet ./...` · `go test -count=1 ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/video/roteiro/') && file_exists('internal/video/ativos/') && file_exists('internal/video/render/')`
* **Spec Ref:** `AC-E11.2`

#### Ticket `TICK-07.2` — Registrar os comandos video-roteiro, video-ativos e video-render com alias diorama-* e aviso de depreciacao
* **Iniciativa:** `V01` | **Épico:** `E11` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.1","SG-E10"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/cli.go, internal/cli/video_alias.go` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E11.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/video_alias_test.go`
* **Outputs:** Artifacts `internal/cli/cli.go, internal/cli/video_alias.go, internal/cli/video_alias_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/cli/ -run TestVideoAlias`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/cli.go') && file_exists('internal/cli/video_alias.go') && file_exists('internal/cli/video_alias_test.go')`
* **Spec Ref:** `AC-E11.1`

#### Ticket `TICK-07.3` — Migrar skills, processos, espinha e PROCESSOS_CANONIZADOS para os nomes legoloop-video-* com redirecionamento
* **Iniciativa:** `V01` | **Épico:** `E11` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.2","SG-E10"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `skills/legoloop-video-roteiro/SKILL.md, skills/legoloop-video-ativos/SKILL.md, espinha/registro.json, docs/PROCESSOS_CANONIZADOS.md` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E11.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `skills/legoloop-video-roteiro/SKILL.md, skills/legoloop-video-ativos/SKILL.md, espinha/registro.json, docs/PROCESSOS_CANONIZADOS.md`
* **Verificação (argv, sem shell):** `go run ./cmd/legoloop espinha verificar`
* **Gate (DSL):** `exit_code == 0 && file_exists('skills/legoloop-video-roteiro/SKILL.md') && file_exists('skills/legoloop-video-ativos/SKILL.md') && file_exists('espinha/registro.json') && file_exists('docs/PROCESSOS_CANONIZADOS.md')`
* **Spec Ref:** `AC-E11.3`

#### Ticket `TICK-07.4` — Aplicar a nota de redirecionamento nos documentos arquivados sem alterar o conteudo
* **Iniciativa:** `V01` | **Épico:** `E11` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.3","SG-E10"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/intents/, docs/specs/, docs/reports/` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E11.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/intents/, docs/specs/, docs/reports/`
* **Verificação (argv, sem shell):** `git diff --stat -- docs/intents docs/specs docs/reports`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/intents/') && file_exists('docs/specs/') && file_exists('docs/reports/')`
* **Spec Ref:** `AC-E11.4`

### 🛑 STOP-GATE HUMANO (`SG-E11`)
* **Gate (DSL):** `prompt('Evidências do Épico 11 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 12 — Épico E12

#### Ticket `TICK-08.1` — Renderizar a cena 1 do diorama-demo pelo processo e colher a aprovacao visual do usuario
* **Iniciativa:** `V01` | **Épico:** `E12` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.3","SG-E11"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `projetos/video-animado/diorama-demo/renders/` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E12.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `projetos/video-animado/diorama-demo/renders/`
* **Verificação (argv, sem shell):** `go run ./cmd/legoloop video-render executar --projeto projetos/video-animado/diorama-demo --cena 1`
* **Gate (DSL):** `exit_code == 0 && file_exists('projetos/video-animado/diorama-demo/renders/')`
* **Spec Ref:** `AC-E12.1`

#### Ticket `TICK-08.2` — Criar o estilo palito herdando do molde-de-pilula e validar com o usuario
* **Iniciativa:** `V01` | **Épico:** `E12` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.1","SG-E11"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `inventario/video/estilos/palito/estilo.json` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E12.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `inventario/video/estilos/palito/estilo.json`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/video/inventario/ -run TestHeranca`
* **Gate (DSL):** `exit_code == 0 && file_exists('inventario/video/estilos/palito/estilo.json')`
* **Spec Ref:** `AC-E12.3`

#### Ticket `TICK-08.3` — Reproduzir o First Brands pelo processo: timeline completa com CHECK OK e ritmo no alvo
* **Iniciativa:** `V01` | **Épico:** `E12` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.1","SG-E11"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `projetos/video-animado/first-brands/timeline.json` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-E12.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `projetos/video-animado/first-brands/timeline.json`
* **Verificação (argv, sem shell):** `go run ./cmd/legoloop video-render timeline --projeto projetos/video-animado/first-brands`
* **Gate (DSL):** `exit_code == 0 && file_exists('projetos/video-animado/first-brands/timeline.json')`
* **Spec Ref:** `AC-E12.2`

#### Ticket `TICK-09.1` — Fechamento: cobertura, dependencias, portabilidade, espinha e legoloop fechar 5 de 5
* **Iniciativa:** `V01` | **Épico:** `E12` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.2","TICK-08.3","SG-E11"]`
* **Isolamento:** Workspace `local` · Executor `native_swarm` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/reports/` · Spec Anchor `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md#AC-NF.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/reports/`
* **Verificação (argv, sem shell):** `go run ./cmd/legoloop fechar legoloop-video-animado`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/')`
* **Spec Ref:** `AC-NF.1, AC-NF.2, AC-NF.3`

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
*Plano de execução fatiado com 25 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*
