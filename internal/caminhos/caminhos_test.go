package caminhos

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// referenciasPermitidas lista os únicos arquivos que ainda podem citar o
// caminho legado "loops" após o rename para "processos" (AC-E06.1). Cada
// entrada tem a razão documentada. Qualquer menção fora desta lista é
// referência viva órfã e reprova o contrato.
var referenciasPermitidas = map[string]string{
	"internal/caminhos/nomes.go":                                                   "literal do dialeto antigo: leitura compat e guia de migração",
	"internal/caminhos/caminhos_test.go":                                           "este teste examina o próprio literal",
	"internal/cli/legado.go":                                                       "guia de migração do caminho legado (AC-E06.3)",
	"internal/cli/legado_test.go":                                                  "teste do guia de migração",
	"internal/studio/static/app.js":                                                "fallback de leitura do caminho legado no Studio",
	"tools/taxlinter/linter_test.go":                                               "corpus do linter de taxonomia: string de teste, não caminho",
	"docs/reports/2026-11-17_REPORT_E06_RENAME_PROCESSOS.md":                       "registro histórico do rename: comparação entre caminho antigo e novo (auditoria do AUDITOR)",
	"docs/archive/quinteto/plans/2026-11-10_PLAN_NORMALIZACAO_TAXONOMICA_LEI15.md": "menção histórica ao literal de código do Quinteto (registro arquivado)",
	"docs/archive/quinteto/specs/2026-11-10_SPEC_NORMALIZACAO_TAXONOMICA_LEI15.md": "menção histórica ao literal de código do Quinteto (registro arquivado)",
}

var extensoesTexto = map[string]bool{
	".go": true, ".md": true, ".json": true, ".sh": true, ".ps1": true,
	".js": true, ".html": true, ".css": true, ".yaml": true, ".yml": true,
	".toml": true, ".txt": true,
}

var dirsIgnorados = map[string]bool{
	".git": true, ".worktrees": true, "bin": true, "dist": true,
	"node_modules": true, ".codebase-memory": true, ".freebuff": true,
	".renderlint": true, ".legoloop-swarm": true, ".claude": true,
	".mimocode": true,
}

func ehArquivoDeTexto(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return true
	}
	return extensoesTexto[ext]
}

// contemReferenciaLegada diz se a linha cita o caminho legado "loops".
func contemReferenciaLegada(linha string) bool {
	for _, padrao := range []string{"loops/", `loops\`, `"loops"`, `'loops'`} {
		if strings.Contains(linha, padrao) {
			return true
		}
	}
	return false
}

// TestReferenciasProcessos prova que zero referência viva ao caminho legado
// sobra no repositório: só os arquivos da lista de exceção documentada podem
// citar "loops" (AC-E06.1).
func TestReferenciasProcessos(t *testing.T) {
	raiz := filepath.Join("..", "..")
	var violacoes []string
	err := filepath.WalkDir(raiz, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if dirsIgnorados[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !ehArquivoDeTexto(path) {
			return nil
		}
		rel, err := filepath.Rel(raiz, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if _, ok := referenciasPermitidas[rel]; ok {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, linha := range strings.Split(string(b), "\n") {
			if contemReferenciaLegada(linha) {
				violacoes = append(violacoes, rel+":"+itoa(i+1))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("varredura: %v", err)
	}
	if len(violacoes) > 0 {
		t.Fatalf("referências vivas ao caminho legado fora da exceção documentada (%d):\n%s",
			len(violacoes), strings.Join(violacoes, "\n"))
	}
}

// TestRenameProcessos prova que o rename físico aconteceu e que o runtime
// resolve processos/ sem acionar o dialeto antigo (AC-E06.1, TICK-06.2).
func TestRenameProcessos(t *testing.T) {
	raiz := filepath.Join("..", "..")
	if _, err := os.Stat(filepath.Join(raiz, DirLoops)); !os.IsNotExist(err) {
		t.Fatalf("diretório legado %q ainda existe na raiz (err=%v)", DirLoops, err)
	}
	proc := filepath.Join(raiz, DirProcessos)
	info, err := os.Stat(proc)
	if err != nil || !info.IsDir() {
		t.Fatalf("diretório %q ausente na raiz: %v", DirProcessos, err)
	}
	entradas, err := os.ReadDir(proc)
	if err != nil {
		t.Fatalf("leitura de %s: %v", DirProcessos, err)
	}
	processos := 0
	for _, e := range entradas {
		if !e.IsDir() {
			continue
		}
		pasta := filepath.Join(proc, e.Name())
		if _, err := os.Stat(CaminhoDAGLeitura(pasta)); err == nil {
			processos++
		}
	}
	if processos < 26 {
		t.Fatalf("processos com DAG em %s: %d; esperado ao menos 26", DirProcessos, processos)
	}
	dirs := DirsLeitura(raiz)
	if len(dirs) == 0 || dirs[0] != proc {
		t.Fatalf("runtime não resolve %s primeiro: %v", DirProcessos, dirs)
	}
	dag := CaminhoDAGLeitura(filepath.Join(proc, "legoloop-fechar"))
	if _, err := os.Stat(dag); err != nil {
		t.Fatalf("runtime não resolve o DAG de legoloop-fechar em %s: %v", DirProcessos, err)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf []byte
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
