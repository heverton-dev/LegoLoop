package main

import (
	"path/filepath"
	"strings"
)

// allowedExceptions define exceções pontuais de vocabulário para contratos externos
// invioláveis (ferramentas CLI de terceiros, tipos de nós em DAGs estáveis e schemas JSON).
var allowedExceptions = map[string]map[string]bool{
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
	"internal/golden/testdata/loops/legoloop-video-motion.dag.json": {
		"card":     true,
		"pipeline": true,
	},
	"internal/golden/testdata/loops/mega-pipeline.loop.json": {
		"pipeline": true,
	},
	"internal/interview/intent_render.go": {
		"task": true,
	},
	"docs/templates/TEMPLATE_DAG.json": {
		"pipeline": true,
		"loop":     true,
	},
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
