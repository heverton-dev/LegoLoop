package studio

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
	"testing"
)

// CA-19.1: Barreira estrita contra emojis em 100% dos arquivos do Studio (HTML, JS, CSS, Templates Go).
// Lei Inquebrável de Design do LegoLoop: ZERO EMOJIS.
func TestZeroEmojisInFrontendAndTemplates(t *testing.T) {
	// Padrão regex cobrindo todas as faixas de emojis e setas dingbat do Unicode
	emojiPattern := regexp.MustCompile(`[\x{1F300}-\x{1FAFF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}\x{1F000}-\x{1F2FF}\x{2B00}-\x{2BFF}\x{FE00}-\x{FE0F}]`)

	var violations []string
	checkContent := func(name, content string) {
		for lineNum, line := range strings.Split(content, "\n") {
			if matches := emojiPattern.FindAllString(line, -1); len(matches) > 0 {
				violations = append(violations, fmt.Sprintf("%s:%d -> encontrados emojis: %s", name, lineNum+1, strings.Join(matches, ", ")))
			}
		}
	}

	// 1. Arquivos estáticos do Studio (percorre recursivamente toda a árvore static/)
	err := fs.WalkDir(staticFS, "static", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return walkErr
		}
		if !(strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css")) {
			return nil
		}
		b, rerr := fs.ReadFile(staticFS, path)
		if rerr != nil {
			return rerr
		}
		checkContent(path, string(b))
		return nil
	})
	if err != nil {
		t.Fatalf("falha ao percorrer staticFS: %v", err)
	}

	// 2. Templates Go e JSON
	for _, file := range []string{"templates.go", "templates.json"} {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("não foi possível ler %s: %v", file, err)
			continue
		}
		checkContent(file, string(content))
	}

	if len(violations) > 0 {
		t.Fatalf("VIOLAÇÃO DA LEI ZERO EMOJIS! Encontradas %d ocorrências de emojis:\n%s", len(violations), strings.Join(violations, "\n"))
	}
}
