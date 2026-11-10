package wave

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"legoloop/internal/board"
)

// Regex para detectar referências a arquivos com extensões comuns ou caminhos de diretório.
var filePathRegex = regexp.MustCompile(`[\w\-\.\/]+\.[a-zA-Z0-9]+`)

// ExtractAffectedFiles inspeciona o corpo e o título de um Card para identificar quais arquivos do repositório serão tocados.
func ExtractAffectedFiles(card board.Card) []string {
	seen := make(map[string]struct{})

	// 1. Inspeciona o corpo
	if card.Body != "" {
		scanFilePaths(card.Body, seen)
	}

	// 2. Se nada foi achado no corpo, inspeciona o título como fallback
	if len(seen) == 0 && card.Title != "" {
		scanFilePaths(card.Title, seen)
	}

	var files []string
	for f := range seen {
		files = append(files, f)
	}
	sort.Strings(files)
	return files
}

// HasOverlap verifica se dois conjuntos de arquivos possuem interseção não vazia.
func HasOverlap(filesA, filesB []string) bool {
	setA := make(map[string]struct{}, len(filesA))
	for _, f := range filesA {
		setA[f] = struct{}{}
	}

	for _, f := range filesB {
		if _, ok := setA[f]; ok {
			return true
		}
	}
	return false
}

// scanFilePaths acrescenta a seen cada caminho de arquivo plausível encontrado no texto.
func scanFilePaths(text string, seen map[string]struct{}) {
	for _, m := range filePathRegex.FindAllString(text, -1) {
		cleaned := filepath.ToSlash(strings.Trim(m, "`'\".,;:()[]{}<>"))

		// Ignora números soltos, URLs e artefatos de markdown irrelevantes
		if strings.HasPrefix(cleaned, "http://") || strings.HasPrefix(cleaned, "https://") {
			continue
		}
		if !strings.Contains(cleaned, ".") {
			continue
		}
		if ext := filepath.Ext(cleaned); ext == "" || len(ext) < 2 {
			continue
		}
		// Filtra arquivos típicos de código/documentação
		seen[cleaned] = struct{}{}
	}
}
