# 🗺️ PLAN: Auditoria Exaustiva Fase B — LegoLoop Full-Stack

> **Status:** Em Planejamento
> **Data:** 2026-11-02
> **Especificação Base:** [2026-11-02_SPEC_AUDITORIA_EXAUSTIVA_LEGOLOOP.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-02_SPEC_AUDITORIA_EXAUSTIVA_LEGOLOOP.md)

---

## 🏛️ ÉPICO 01: Canvas (Zoom, Pan, Limites, Ctrl+0)

### Node E01.1 — Canvas: Mapear Funções de Zoom
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `[]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:get_code_snippet` |
| **input** | Grafo AST de `internal/studio/static/app.js` |
| **output** | `scratch/E01_1_zoom_map.json` |
| **gate** | `has_keys(artifact, ['function','file','line'])` |
| **spec_ref** | `AC-E01.1` |

### Node E01.2 — Canvas: Mapear Funções de Pan e Limites
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `[]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:trace_path` |
| **input** | Grafo AST de `internal/studio/static/app.js` |
| **output** | `scratch/E01_2_pan_map.json` |
| **gate** | `has_keys(artifact, ['function','file','line','has_clamp'])` |
| **spec_ref** | `AC-E01.2` |

### Node E01.3 — Canvas: Auditar Atalho Ctrl+0 / Bounding-box
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E01.1"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `codebase-memory-mcp:get_code_snippet` |
| **input** | `scratch/E01_1_zoom_map.json` |
| **output** | `scratch/E01_3_fit_analysis.json` |
| **gate** | `field_present(artifact, 'zoom_fit_strategy')` |
| **spec_ref** | `AC-E01.3` |

### Node E01.4 — Canvas: Consolidar Achados E01
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E01.1", "E01.2", "E01.3"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `node` |
| **input** | `scratch/E01_1_zoom_map.json`, `scratch/E01_2_pan_map.json`, `scratch/E01_3_fit_analysis.json` |
| **output** | `scratch/audit_E01_canvas.json` |
| **gate** | `exit_code == 0 && file_exists('scratch/audit_E01_canvas.json')` |
| **spec_ref** | `AC-E01.4` |

---

## 🏛️ ÉPICO 02: Conexões & Arestas

### Node E02.1 — Arestas: Rastrear Fluxo de Criação de Arestas
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E01.4"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:trace_path` |
| **input** | Grafo AST de `internal/studio/static/app.js` |
| **output** | `scratch/E02_1_edge_flow.json` |
| **gate** | `has_keys(artifact, ['start_fn','drag_fn','drop_fn','phantom_edge_guard'])` |
| **spec_ref** | `AC-E02.1` |

### Node E02.2 — Arestas: Smoke Test com Obscura
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E01.4"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `obscura`, `powershell` |
| **input** | Executável `./bin/legoloop.exe studio` |
| **output** | `scratch/E02_2_smoke.png` |
| **gate** | `file_exists('scratch/E02_2_smoke.png')` |
| **spec_ref** | `AC-E02.2` |

### Node E02.3 — Arestas: Interação Dinâmica Playwright
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E02.1", "E02.2"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `playwright`, `node` |
| **input** | Servidor ativo na porta 8080 + `scratch/E02_1_edge_flow.json` |
| **output** | `scratch/E02_3_playwright_log.json` |
| **gate** | `exit_code == 0 && has_keys(artifact, ['phantom_edge','escape_cancel','reconnect'])` |
| **spec_ref** | `AC-E02.3` |

### Node E02.4 — Arestas: Encerramento do Servidor e Consolidação
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E02.3"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `powershell`, `node` |
| **input** | Processo ativo `legoloop.exe` + dados de E02.1 e E02.3 |
| **output** | `scratch/audit_E02_conexoes.json` |
| **gate** | `process_dead('legoloop') && file_exists('scratch/audit_E02_conexoes.json')` |
| **spec_ref** | `AC-E02.4` |

---

## 🏛️ ÉPICO 03: Undo/Redo

### Node E03.1 — Histórico: Estrutura da Pilha e Callers
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E02.4"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:get_code_snippet` |
| **input** | Grafo AST de `internal/studio/static/app.js` |
| **output** | `scratch/E03_1_history_map.json` |
| **gate** | `has_keys(artifact, ['stack_structure','mutations_tracked','limit'])` |
| **spec_ref** | `AC-E03.1` |

### Node E03.2 — Histórico: Análise de Mousemove vs Mouseup
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E03.1"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `codebase-memory-mcp:get_code_snippet` |
| **input** | `scratch/E03_1_history_map.json` |
| **output** | `scratch/E03_2_timing_analysis.json` |
| **gate** | `field_present(artifact, 'snapshot_event')` |
| **spec_ref** | `AC-E03.2` |

### Node E03.3 — Histórico: Consolidar Achados E03
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E03.2"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `node` |
| **input** | `scratch/E03_1_history_map.json`, `scratch/E03_2_timing_analysis.json` |
| **output** | `scratch/audit_E03_undoredo.json` |
| **gate** | `file_exists('scratch/audit_E03_undoredo.json')` |
| **spec_ref** | `AC-E03.3` |

---

## 🏛️ ÉPICO 04: Teclas de Atalho

### Node E04.1 — Teclas: Mapeamento Global e Guards
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E03.3"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:get_code_snippet` |
| **input** | Grafo AST de `internal/studio/static/app.js` e `ui.js` |
| **output** | `scratch/E04_1_keys_map.json` |
| **gate** | `has_keys(artifact, ['input_guard','escape_targets','shortcuts'])` |
| **spec_ref** | `AC-E04.1` |

### Node E04.2 — Teclas: Testes Playwright (Delete, Escape, Ctrl+Z)
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E03.3"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `playwright`, `node`, `powershell` |
| **input** | Servidor ativo em porta 8080 |
| **output** | `scratch/E04_2_playwright_log.json` |
| **gate** | `exit_code == 0 && has_keys(artifact, ['delete_in_input','escape_modal'])` |
| **spec_ref** | `AC-E04.2` |

### Node E04.3 — Teclas: Encerramento do Servidor e Consolidação
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E04.1", "E04.2"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `powershell`, `node` |
| **input** | Processo ativo + relatórios de E04.1 e E04.2 |
| **output** | `scratch/audit_E04_atalhos.json` |
| **gate** | `process_dead('legoloop') && file_exists('scratch/audit_E04_atalhos.json')` |
| **spec_ref** | `AC-E04.3` |

---

## 🏛️ ÉPICO 05: Paleta & Drag-and-Drop

### Node E05.1 — Paleta: Compensação de Offset de Drop
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E04.3"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:trace_path` |
| **input** | Grafo AST de `palette.js` e `app.js` |
| **output** | `scratch/E05_1_offset_analysis.json` |
| **gate** | `has_keys(artifact, ['formula','zoom_compensated','pan_compensated'])` |
| **spec_ref** | `AC-E05.1` |

### Node E05.2 — Paleta: Virtualização Concorrente
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E04.3"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `codebase-memory-mcp:get_file_outline`, `codebase-memory-mcp:get_code_snippet` |
| **input** | Grafo AST de `virtual.js` e `search.js` |
| **output** | `scratch/E05_2_virtual_analysis.json` |
| **gate** | `field_present(artifact, 'search_glitch_risk')` |
| **spec_ref** | `AC-E05.2` |

### Node E05.3 — Paleta: Validação E2E com Playwright
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E05.1", "E05.2"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `playwright`, `obscura`, `powershell` |
| **input** | Servidor ativo em porta 8080 |
| **output** | `scratch/E05_3_playwright_log.json` |
| **gate** | `exit_code == 0 && file_exists('scratch/E05_3_playwright_log.json')` |
| **spec_ref** | `AC-E05.3` |

### Node E05.4 — Paleta: Encerramento do Servidor e Consolidação
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E05.3"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `powershell`, `node` |
| **input** | Processo ativo + dados de E05.1 a E05.3 |
| **output** | `scratch/audit_E05_paleta.json` |
| **gate** | `process_dead('legoloop') && file_exists('scratch/audit_E05_paleta.json')` |
| **spec_ref** | `AC-E05.4` |

---

## 🏛️ ÉPICO 06: Inspetor

### Node E06.1 — Inspetor: Mapeamento de Debounce e Reatividade
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E05.4"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `codebase-memory-mcp:get_file_outline`, `codebase-memory-mcp:search_code` |
| **input** | Grafo AST de `inspector.js` e `app.js` |
| **output** | `scratch/E06_1_inspector_map.json` |
| **gate** | `has_keys(artifact, ['debounce_status','card_sync'])` |
| **spec_ref** | `AC-E06.1` |

### Node E06.2 — Inspetor: Teste de Rajadas e Race Conditions
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E05.4"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `playwright`, `powershell` |
| **input** | Servidor ativo em porta 8080 |
| **output** | `scratch/E06_2_playwright_log.json` |
| **gate** | `exit_code == 0 && field_present(artifact, 'requests_fired')` |
| **spec_ref** | `AC-E06.2` |

### Node E06.3 — Inspetor: Encerramento e Consolidação
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E06.1", "E06.2"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `powershell`, `node` |
| **input** | Processo ativo + relatórios de E06.1 e E06.2 |
| **output** | `scratch/audit_E06_inspetor.json` |
| **gate** | `process_dead('legoloop') && file_exists('scratch/audit_E06_inspetor.json')` |
| **spec_ref** | `AC-E06.3` |

---

## 🏛️ ÉPICO 07: Inventário

### Node E07.1 — Inventário: Auditoria de Segurança innerHTML
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E06.3"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:get_code_snippet` |
| **input** | Grafo AST de `components/lego-item-modal.js` |
| **output** | `scratch/E07_1_modal_security.json` |
| **gate** | `has_keys(artifact, ['injection_points','risk_level'])` |
| **spec_ref** | `AC-E07.1` |

### Node E07.2 — Inventário: Navegação e Pesquisa E2E
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E06.3"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `playwright`, `obscura`, `powershell` |
| **input** | Servidor ativo em porta 8080 |
| **output** | `scratch/E07_2_playwright_log.json` |
| **gate** | `exit_code == 0 && has_keys(artifact, ['tabs_verified','accented_search'])` |
| **spec_ref** | `AC-E07.2` |

### Node E07.3 — Inventário: Encerramento e Consolidação
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E07.1", "E07.2"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `powershell`, `node` |
| **input** | Processo ativo + relatórios de E07.1 e E07.2 |
| **output** | `scratch/audit_E07_inventario.json` |
| **gate** | `process_dead('legoloop') && file_exists('scratch/audit_E07_inventario.json')` |
| **spec_ref** | `AC-E07.3` |

---

## 🏛️ ÉPICO 08: Templates

### Node E08.1 — Templates: Duplicação de Nomes e Fallback
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E07.3"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:get_code_snippet` |
| **input** | Código fonte `internal/studio/templates.go` e `app.js` |
| **output** | `scratch/E08_1_templates_map.json` |
| **gate** | `has_keys(artifact, ['cause_duplication','icon_fallback'])` |
| **spec_ref** | `AC-E08.1` |

### Node E08.2 — Templates: Teste E2E de Confirmação de Carga
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E07.3"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `playwright`, `powershell` |
| **input** | Servidor ativo em porta 8080 |
| **output** | `scratch/E08_2_playwright_log.json` |
| **gate** | `exit_code == 0 && field_present(artifact, 'dialog_prompted')` |
| **spec_ref** | `AC-E08.2` |

### Node E08.3 — Templates: Encerramento e Consolidação
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E08.1", "E08.2"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `powershell`, `node` |
| **input** | Processo ativo + relatórios de E08.1 e E08.2 |
| **output** | `scratch/audit_E08_templates.json` |
| **gate** | `process_dead('legoloop') && file_exists('scratch/audit_E08_templates.json')` |
| **spec_ref** | `AC-E08.3` |

---

## 🏛️ ÉPICO 09: Modal de Execução (Run)

### Node E09.1 — Run: Rastreamento do Botão Run e Validate
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E08.3"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:trace_path` |
| **input** | Grafo AST de `internal/studio/static/app.js` |
| **output** | `scratch/E09_1_run_flow.json` |
| **gate** | `has_keys(artifact, ['validate_callers','unhandled_return'])` |
| **spec_ref** | `AC-E09.1` |

### Node E09.2 — Run: Streaming de Logs no Backend
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E08.3"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `codebase-memory-mcp:get_file_outline`, `codebase-memory-mcp:get_code_snippet` |
| **input** | Código fonte `internal/studio/run_api.go` |
| **output** | `scratch/E09_2_stream_protocol.json` |
| **gate** | `field_present(artifact, 'protocol_type')` |
| **spec_ref** | `AC-E09.2` |

### Node E09.3 — Run: Testes Playwright de Execução e Badges
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E09.1", "E09.2"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `playwright`, `powershell` |
| **input** | Servidor ativo em porta 8080 |
| **output** | `scratch/E09_3_playwright_log.json` |
| **gate** | `exit_code == 0 && has_keys(artifact, ['badge_states','abort_handling'])` |
| **spec_ref** | `AC-E09.3` |

### Node E09.4 — Run: Encerramento e Consolidação
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E09.3"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `powershell`, `node` |
| **input** | Processo ativo + relatórios de E09.1 a E09.3 |
| **output** | `scratch/audit_E09_run.json` |
| **gate** | `process_dead('legoloop') && file_exists('scratch/audit_E09_run.json')` |
| **spec_ref** | `AC-E09.4` |

---

## 🏛️ ÉPICO 10: Cofre Canônico (Vault)

### Node E10.1 — Vault: Mapeamento de Handlers e Null Assets
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E09.4"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `codebase-memory-mcp:search_code`, `codebase-memory-mcp:get_code_snippet` |
| **input** | Grafo AST de `assets_api.go` |
| **output** | `scratch/E10_1_vault_map.json` |
| **gate** | `has_keys(artifact, ['handlers','status_code_on_nil'])` |
| **spec_ref** | `AC-E10.1` |

### Node E10.2 — Vault: Testes Playwright de Dry-Run e Apply
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `parallel` |
| **depends_on** | `["E09.4"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `playwright`, `powershell` |
| **input** | Servidor ativo em porta 8080 |
| **output** | `scratch/E10_2_playwright_log.json` |
| **gate** | `exit_code == 0 && field_present(artifact, 'ui_feedback_present')` |
| **spec_ref** | `AC-E10.2` |

### Node E10.3 — Vault: Encerramento e Consolidação
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E10.1", "E10.2"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `powershell`, `node` |
| **input** | Processo ativo + relatórios de E10.1 e E10.2 |
| **output** | `scratch/audit_E10_vault.json` |
| **gate** | `process_dead('legoloop') && file_exists('scratch/audit_E10_vault.json')` |
| **spec_ref** | `AC-E10.3` |

---

## 🏛️ ÉPICO 11: Consolidação Final

### Node E11.1 — Síntese: Unificação e Deduplicação dos Achados
| Campo | Valor |
|---|---|
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E10.3"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `node` |
| **input** | Todos os `scratch/audit_E*.json` |
| **output** | `scratch/matriz_bugs_consolidada.json` |
| **gate** | `file_exists('scratch/matriz_bugs_consolidada.json')` |
| **spec_ref** | `AC-E11.1` |

### Node E11.2 — Síntese: Emissão do Relatório Master Exaustivo
| Campo | Valor |
|---|---|
| **nature** | `epic` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E11.1"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | Síntese de documentação |
| **input** | `scratch/matriz_bugs_consolidada.json` |
| **output** | `docs/reports/2026-11-02_REPORT_AUDITORIA_EXAUSTIVA.md` |
| **gate** | `file_exists('docs/reports/2026-11-02_REPORT_AUDITORIA_EXAUSTIVA.md')` |
| **spec_ref** | `AC-E11.2` |

### Node E11.3 — Limpeza: Eliminação de Resíduos Temporários
| Campo | Valor |
|---|---|
| **nature** | `nano` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E11.2"]` |
| **model** | `9router/code-nano` · effort `low` · fallback `flash-lite` |
| **tooling** | `powershell` |
| **input** | Lista de arquivos temporários em `scratch/` |
| **output** | `scratch/cleanup.log` |
| **gate** | `exit_code == 0` |
| **spec_ref** | `AC-E11.3` |

---

### 🛑 STOP-GATE DE AUDITORIA HUMANA — FASE B COMPLETA
* [ ] Relatório master gerado com matriz P0-P3 completa
* [ ] Todos os processos de background finalizados
* [ ] Aval do usuário para avançar à Fase C (Correções TDD)
