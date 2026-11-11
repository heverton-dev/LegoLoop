package studio

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// CA-19.1: Barreira estrita contra emojis em 100% dos arquivos do Studio (HTML, JS, CSS, Templates Go).
// Lei Inquebrável de Design do LegoLoop: ZERO EMOJIS.
func TestZeroEmojisInFrontendAndTemplates(t *testing.T) {
	// Padrão regex cobrindo todas as faixas de emojis do Unicode
	emojiPattern := regexp.MustCompile(`[\x{1F300}-\x{1FAFF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}\x{1F000}-\x{1F2FF}\x{FE00}-\x{FE0F}]`)

	var filesToCheck []string

	// 1. Arquivos estáticos do Studio
	staticDir := filepath.Join("static")
	entries, err := os.ReadDir(staticDir)
	if err != nil {
		t.Fatalf("falha ao ler diretório static: %v", err)
	}
	for _, e := range entries {
		if !e.IsDir() && (strings.HasSuffix(e.Name(), ".html") || strings.HasSuffix(e.Name(), ".js") || strings.HasSuffix(e.Name(), ".css")) {
			filesToCheck = append(filesToCheck, filepath.Join(staticDir, e.Name()))
		}
	}

	// 2. Templates Go
	filesToCheck = append(filesToCheck, "templates.go")

	var violations []string

	for _, file := range filesToCheck {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("não foi possível ler %s: %v", file, err)
			continue
		}
		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			if matches := emojiPattern.FindAllString(line, -1); len(matches) > 0 {
				violations = append(violations, file+":"+string(rune('0'+lineNum+1))+" -> encontrados emojis: "+strings.Join(matches, ", "))
			}
		}
	}

	if len(violations) > 0 {
		t.Fatalf("VIOLAÇÃO DA LEI ZERO EMOJIS! Encontradas %d ocorrências de emojis:\n%s", len(violations), strings.Join(violations, "\n"))
	}
}
