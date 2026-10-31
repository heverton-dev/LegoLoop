# 📋 Plano de Implementação — Varredura Cascata Hiper-Profunda do SO

> **Spec:** [docs/specs/SPEC_SCANNER_CASCATA_PROFUNDA.md](docs/specs/SPEC_SCANNER_CASCATA_PROFUNDA.md)  
> **Método:** Fatias verticais (tracer bullet) · TDD estrito (Red → Green) · Critérios binários  
> **Gate de cada ticket:** `go test ./...` verde, `funclen` ≤ 35 linhas, cobertura ≥ 85%, zero stubs.

---

## Épico 13: Varredura Cascata Hiper-Profunda

- [x] **TICK-1301 (Micro):** Coletor Concorrente de Projetos (`internal/assets/discover.go`).  
  *Critérios:* Varre raízes configuradas (`Desktop`, `Projects`, `Workspace`, `repos`, `src`, `dev`, `code`, `Documents`) até profundidade 4. Identifica marcadores (`.git`, `package.json`, `go.mod`, `.claude/`, `.gemini/`). Poda estritamente `node_modules`, `.git`, `vendor`, `.venv`, `.cache`.  
  *Teste:* `internal/assets/discover_test.go` com estrutura de diretórios aninhada e podas comprovadas.

- [x] **TICK-1302 (Nano):** Camada Scaffolds & Templates (`internal/assets/scaffolds.go`).  
  *Critérios:* Varre `~/.aidd/scaffolds/**` e `~/.aidd/canonical/**` extraindo skills, hooks, scripts e loops com `scope: "scaffold"`.  
  *Teste:* `internal/assets/scaffolds_test.go` validando escopo e integridade dos itens.

- [x] **TICK-1303 (Micro):** Cascata de Hooks e MCPs em todos os Projetos (`internal/assets/hooks.go`, `internal/assets/mcps.go`).  
  *Critérios:* Indexa hooks de todos os projetos descobertos (`.claude/settings.json`, `settings.local.json`) e MCPs (`.mcp.json`, `.claude.json`), atribuindo corretamente escopos `project` e `local`. Mantém regras de sanitização SEC-SCAN-1/2.  
  *Teste:* `internal/assets/hooks_cascade_test.go` e `internal/assets/mcps_cascade_test.go`.

- [x] **TICK-1304 (Nano):** Cascata de Scripts e Loops em Projetos e Scaffolds (`internal/assets/scripts.go`, `internal/assets/loops.go`).  
  *Critérios:* Coleta scripts em pastas de projeto (`scripts/`, `tools/`, `bin/`) e loops commitados nos projetos (`*.loop.json`).  
  *Teste:* `internal/assets/scripts_cascade_test.go`.

- [x] **TICK-1305 (Nano):** Expansão de Binários e Harnesses do Usuário (`internal/assets/system.go`).  
  *Critérios:* Varre diretórios de binários de linguagens e ferramentas do usuário (`~/.cargo/bin`, `~/go/bin`, `AppData/Local/Programs`, `~/.local/bin`).  
  *Teste:* `internal/assets/system_test.go`.

- [x] **TICK-1306 (Micro):** Verificação de Performance, Auditoria de Segurança e E2E no Studio.  
  *Critérios:* Suíte completa verde; benchmark de re-varredura com cache hit ≤ 150 ms; sem vazamento de memória ou caminhos de segredo; visualização na rota `/inventario` com contadores reais das 5 camadas.
