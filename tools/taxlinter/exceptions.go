package main

import (
	"path/filepath"
	"strings"
)

// allowedExceptions define exceções pontuais de vocabulário para contratos externos
// invioláveis (ferramentas CLI de terceiros, tipos de nós em DAGs estáveis e schemas JSON).
var allowedExceptions = map[string]map[string]bool{
	"internal/cli/rota.go": {
		"task": true,
	},
	"internal/roteamento/classe.go": {
		"macro": true,
	},
	"internal/orchestrix/unit_claude.go": {
		"task": true,
	},
	"internal/orchestrix/orchestrix.go": {
		"task": true,
		"loop": true,
	},
	"internal/orchestrix/subloop.go": {
		"loop": true,
	},
	"internal/orchestrix/subloop_test.go": {
		"loop": true,
	},
	"internal/orchestrix/subloop_artefato_test.go": {
		"loop": true,
	},
	"internal/orchestrix/unit.go": {
		"task": true,
	},
	"internal/orchestrix/prompt.go": {
		"tarefa": true,
	},
	"internal/orchestrix/unit_run.go": {
		"task": true,
	},
	"internal/orchestrix/unit_run_test.go": {
		"task": true,
	},
	"internal/orchestrix/unit_claude_test.go": {
		"task":  true,
		"macro": true,
	},
	"internal/orchestrix/unit_test.go": {
		"task": true,
	},
	"internal/legopress/auditor_test.go": {
		"pipeline": true,
	},
	"internal/golden/testdata/processos/legoloop-video-motion.dag.json": {
		"card":     true,
		"pipeline": true,
	},
	"internal/interview/intent_render.go": {
		"task": true,
	},
	"docs/templates/TEMPLATE_DAG.json": {
		"pipeline": true,
		"loop":     true,
	},
	"skills/legoloop-loop-map/SKILL.md": {
		"loop": true,
	},
	"skills/legoloop-mega-pipeline/SKILL.md": {
		"pipeline": true,
	},
	"internal/caminhos/nomes.go": {
		"loop":     true,
		"pipeline": true,
	},
	"internal/caminhos/nomes_test.go": {
		"loop": true,
	},
	"internal/cli/commands.go": {
		"task": true,
		"loop": true,
	},
	"internal/studio/swarm_api.go": {
		"card": true,
	},
	"internal/studio/ui_logic_test.go": {
		"loop": true,
	},
	"internal/studio/static_test.go": {
		"loop": true,
	},
	"internal/studio/loops_api.go": {
		"loop": true,
	},
	"internal/studio/loops_api_test.go": {
		"loop": true,
	},
	"internal/studio/processos_test.go": {
		"loop": true,
	},
	"internal/studio/processos_rest_test.go": {
		"loop": true,
	},
	"internal/studio/pieces.go": {
		"loop": true,
	},
	"internal/studio/pieces_test.go": {
		"loop": true,
	},
	"internal/studio/materialize_test.go": {
		"loop": true,
	},
	"internal/studio/fabricar_plot_test.go": {
		"loop": true,
	},
	"internal/studio/triagem_api_test.go": {
		"loop": true,
	},
	"internal/studio/run_api.go": {
		"loop": true,
	},
	"internal/studio/server_test.go": {
		"tarefa": true,
		"card":   true,
	},
	"internal/assets/loops.go": {
		"loop": true,
	},
	"internal/assets/skills.go": {
		"pipeline": true,
	},
	"internal/assets/types.go": {
		"loop": true,
	},
	"internal/cli/cartao_test.go": {
		"fluxo": true,
	},
	"internal/cli/compor_binary_test.go": {
		"loop": true,
	},
	"internal/cli/compor_real_binary_test.go": {
		"tarefa": true,
		"loop":   true,
	},
	"internal/cli/fabricar_report_texto_test.go": {
		"loop": true,
	},
	"internal/cli/fabricar_test.go": {
		"loop":     true,
		"pipeline": true,
	},
	"internal/cli/listar_test.go": {
		"loop": true,
	},
	"internal/cli/processo_real_binary_test.go": {
		"loop":     true,
		"pipeline": true,
	},
	"internal/cli/remover_test.go": {
		"loop": true,
	},
	"internal/cli/upgrade.go": {
		"loop": true,
	},
	"internal/cli/upgrade_hint_test.go": {
		"fluxo": true,
		"loop":  true,
	},
	"internal/compor/compor.go": {
		"loop": true,
	},
	"internal/compor/compor_test.go": {
		"loop": true,
	},
	"internal/dag/dialetos_test.go": {
		"loop": true,
	},
	"internal/dag/parse.go": {
		"task": true,
	},
	"internal/dag/parse_dual_test.go": {
		"loop": true,
	},
	"internal/dag/types.go": {
		"loop": true,
		"task": true,
	},
	"internal/dag/v2_test.go": {
		"loop": true,
	},
	"internal/fabricar/camada4.go": {
		"loop": true,
	},
	"internal/fabricar/curadoria_fronteira.go": {
		"loop":     true,
		"pipeline": true,
	},
	"internal/fabricar/pieces.go": {
		"loop": true,
	},
	"internal/fabricar/prompts_test.go": {
		"fluxo":  true,
		"loop":   true,
		"tarefa": true,
	},
	"internal/fechar/fechar.go": {
		"loop": true,
	},
	"internal/fechar/fechar_test.go": {
		"loop": true,
	},
	"internal/matriz/registro.go": {
		"loop": true,
	},
	"internal/matriz/vocab.default.json": {
		"pipeline": true,
		"tarefa":   true,
		"esteira":  true,
	},
	"internal/migrar/migrar.go": {
		"loop":     true,
		"pipeline": true,
	},
	"internal/migrar/migrar_test.go": {
		"loop":     true,
		"pipeline": true,
	},
	"internal/orchestrix/hygiene_test.go": {
		"loop": true,
	},
	"internal/paths/paths.go": {
		"loop": true,
	},
	"internal/processo/cofre_test.go": {
		"loop": true,
	},
	"internal/processo/origem_test.go": {
		"loop": true,
	},
	"internal/processo/persistencia.go": {
		"loop": true,
	},
	"internal/processo/projecoes.go": {
		"loop": true,
	},
	"internal/processo/projecoes_test.go": {
		"loop": true,
	},
	"internal/processo/store_test.go": {
		"loop": true,
	},
	"internal/refinar/refinar.go": {
		"loop": true,
	},
	"internal/refinar/refinar_test.go": {
		"loop": true,
	},
	"internal/remover/remover.go": {
		"loop": true,
	},
	"internal/runner/n10_test.go": {
		"pipeline": true,
	},
	"internal/runner/telemetry.go": {
		"pipeline": true,
	},
	"internal/runner/triggers.go": {
		"loop": true,
	},
	"internal/ste/glossario.json": {
		"esteira":    true,
		"fluxo":      true,
		"componente": true,
		"estágio":    true,
	},
	"internal/ste/prompt_prova.go": {
		"bloco": true,
	},
	"internal/ste/ste_test.go": {
		"esteira": true,
	},
	"internal/studio/assets_api.go": {
		"loop": true,
	},
	"internal/studio/compose.go": {
		"loop": true,
	},
	"internal/studio/dissecar_handler.go": {
		"pipeline": true,
	},
	"internal/studio/dissecar_handler_test.go": {
		"pipeline": true,
	},
	// Rename do caminho legado para processos (AC-E06): o literal antigo aparece so
	// como deteccao/allowlist do caminho antigo, nunca como termo de dominio.
	"internal/caminhos/caminhos_test.go": {
		"loop": true,
	},
	"internal/cli/legado.go": {
		"loop": true,
	},
	"internal/cli/legado_test.go": {
		"loop": true, "pipeline": true,
	},
	"internal/fechar/quarteto_test.go": {
		"loop": true,
	},
	// Divida herdada do Quinteto migrado (AC-E05.2): documentacao historica migrada para docs/ raiz.
	"docs/intents/2026-11-09_INTENT_AI_MOTION_GRAPHICS.md":      {"bloco": true, "esteira": true, "fluxo": true, "loop": true, "módulo": true, "pipeline": true, "task": true},
	"docs/plans/2026-11-09_PLAN_AI_MOTION_GRAPHICS.md":          {"bloco": true, "loop": true, "macro": true, "pipeline": true, "task": true},
	"docs/routes/2026-11-09_ROUTE_AI_MOTION_GRAPHICS.json":      {"loop": true, "macro": true, "pipeline": true, "task": true},
	"docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md":          {"fase": true, "fluxo": true, "loop": true, "módulo": true, "pipeline": true},
	"docs/intents/2026-11-14_INTENT_FRONTEIRA_AI_DRIVEN.md":     {"capacidade": true},
	"docs/intents/2026-11-07_INTENT_CANAL_REDC.md":              {"bloco": true, "capacidade": true, "esteira": true, "loop": true, "módulo": true, "pipeline": true, "task": true},
	"docs/plans/2026-11-07_PLAN_CANAL_REDC.md":                  {"bloco": true, "loop": true, "pipeline": true, "task": true},
	"docs/routes/2026-11-07_ROUTE_CANAL_REDC.json":              {"loop": true, "pipeline": true, "task": true},
	"docs/specs/2026-11-07_SPEC_CANAL_REDC.md":                  {"bloco": true, "esteira": true, "loop": true, "pipeline": true, "task": true},
	"docs/intents/2026-11-09_INTENT_LEGOLOOP_CARTOON_SHORTS.md": {"bloco": true, "esteira": true, "fluxo": true, "loop": true, "módulo": true, "pipeline": true, "task": true},
	"docs/plans/2026-11-09_PLAN_LEGOLOOP_CARTOON_SHORTS.md":     {"bloco": true, "macro": true, "pipeline": true, "task": true},
	"docs/routes/2026-11-09_ROUTE_LEGOLOOP_CARTOON_SHORTS.json": {"macro": true, "pipeline": true, "task": true},
	"docs/specs/2026-11-09_SPEC_LEGOLOOP_CARTOON_SHORTS.md":     {"card": true, "estágio": true, "esteira": true, "fluxo": true, "loop": true, "pipeline": true},
	"docs/intents/2026-11-14_INTENT_DESIGN_UI.md":               {"task": true},
	"docs/plans/2026-11-14_PLAN_DESIGN_UI.md":                   {"loop": true, "pipeline": true, "task": true},
	"docs/routes/2026-11-14_ROUTE_DESIGN_UI.json":               {"loop": true, "pipeline": true, "task": true},
	"docs/specs/2026-11-14_SPEC_DESIGN_UI.md":                   {"loop": true},
	"docs/intents/2026-11-10_INTENT_DISSECAR.md":                {"bloco": true, "card": true, "esteira": true, "fluxo": true, "loop": true, "pipeline": true, "task": true},
	"docs/plans/2026-11-10_PLAN_DISSECAR.md":                    {"bloco": true, "card": true, "componente": true, "fluxo": true, "loop": true, "pipeline": true, "task": true},
	"docs/routes/2026-11-10_ROUTE_DISSECAR.json":                {"loop": true, "pipeline": true, "task": true},
	"docs/specs/2026-11-10_SPEC_DISSECAR.md":                    {"bloco": true, "card": true, "componente": true, "pipeline": true},
	"docs/intents/2026-11-15_INTENT_FABRICAR_PECA_UNICA.md":     {"task": true},
	"docs/plans/2026-11-15_PLAN_FABRICAR_PECA_UNICA.md":         {"loop": true, "pipeline": true},
	"docs/routes/2026-11-15_ROUTE_FABRICAR_PECA_UNICA.json":     {"loop": true, "pipeline": true},
	"docs/specs/2026-11-15_SPEC_FABRICAR_PECA_UNICA.md":         {"loop": true, "pipeline": true},
	"docs/intents/2026-11-14_INTENT_SWARM.md":                   {"loop": true, "pipeline": true, "task": true},
	"docs/plans/2026-11-14_PLAN_SWARM.md":                       {"loop": true, "pipeline": true},
	"docs/routes/2026-11-14_ROUTE_SWARM.json":                   {"loop": true, "task": true},
	"docs/specs/2026-11-14_SPEC_SWARM.md":                       {"loop": true, "pipeline": true},
	"processos/legoloop-swarm/pipeline.dag.json":                {"loop": true},
	"docs/intents/2026-11-07_INTENT_VIDEO_MASTER.md":            {"bloco": true, "esteira": true, "fluxo": true, "loop": true, "módulo": true, "pipeline": true, "task": true},
	"docs/plans/2026-11-07_PLAN_VIDEO_MASTER.md":                {"bloco": true, "loop": true, "pipeline": true, "task": true},
	"docs/routes/2026-11-07_ROUTE_VIDEO_MASTER.json":            {"loop": true, "pipeline": true, "task": true},
	"docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md":                {"esteira": true, "fluxo": true, "loop": true, "pipeline": true, "task": true},
	"docs/intents/2026-11-07_INTENT_LEGOLOOP_VIDEO_REA.md":      {"bloco": true, "esteira": true, "loop": true, "módulo": true, "pipeline": true, "task": true},
	"docs/plans/2026-11-07_PLAN_LEGOLOOP_VIDEO_REA.md":          {"bloco": true, "loop": true, "macro": true, "pipeline": true},
	"docs/routes/2026-11-07_ROUTE_LEGOLOOP_VIDEO_REA.json":      {"loop": true, "macro": true, "pipeline": true, "task": true},
	"docs/specs/2026-11-07_SPEC_LEGOLOOP_VIDEO_REA.md":          {"capacidade": true, "esteira": true, "loop": true, "módulo": true, "pipeline": true},
	"docs/intents/2026-11-15_INTENT_VIDEO_REA_RELATORIO.md":     {"task": true},
	"docs/plans/2026-11-15_PLAN_VIDEO_REA_RELATORIO.md":         {"loop": true, "pipeline": true},
	"docs/routes/2026-11-15_ROUTE_VIDEO_REA_RELATORIO.json":     {"loop": true, "pipeline": true},
	"docs/specs/2026-11-15_SPEC_VIDEO_REA_RELATORIO.md":         {"loop": true, "pipeline": true},
	"docs/intents/2026-11-15_INTENT_VIDEO_REA_TRANSCREVER.md":   {"task": true},
	"docs/plans/2026-11-15_PLAN_VIDEO_REA_TRANSCREVER.md":       {"loop": true, "pipeline": true},
	"docs/routes/2026-11-15_ROUTE_VIDEO_REA_TRANSCREVER.json":   {"loop": true, "pipeline": true},
	"docs/specs/2026-11-15_SPEC_VIDEO_REA_TRANSCREVER.md":       {"loop": true, "pipeline": true},
	"docs/intents/2026-11-15_INTENT_WEB_QA.md":                  {"fluxo": true, "storage": true, "task": true},
	"docs/specs/2026-11-15_SPEC_WEB_QA.md":                      {"loop": true, "pipeline": true, "storage": true},
}

// isAllowedException confere se o par (arquivo, termo) é uma exceção de contrato externo.
func isAllowedException(filePath, term string) bool {
	norm := filepath.ToSlash(filePath)
	norm = strings.TrimPrefix(norm, "./")
	lowerTerm := strings.ToLower(term)
	for pathKey, terms := range allowedExceptions {
		if norm == pathKey || strings.HasSuffix(norm, "/"+pathKey) {
			if terms[lowerTerm] {
				return true
			}
		}
	}
	return false
}
