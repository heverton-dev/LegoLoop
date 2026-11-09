package gate

import (
	"fmt"
	"strings"
)

// MaxSTE100WordsPerSentence define o limite estrito da norma ASD-STE100 para sentenças.
const MaxSTE100WordsPerSentence = 20

var hedgingTerms = []string{
	"talvez",
	"pode ser",
	"quem sabe",
	"possivelmente",
	"provavelmente",
}

// ValidateSTE100 inspeciona um texto procedural garantindo frases com <= 20 palavras e zero hedging.
func ValidateSTE100(text string) ([]string, bool) {
	var violations []string
	sentences := splitSentences(text)

	for _, s := range sentences {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			continue
		}
		if v := checkSentenceWords(trimmed); v != "" {
			violations = append(violations, v)
		}
		if v := checkHedging(trimmed); v != "" {
			violations = append(violations, v)
		}
	}

	return violations, len(violations) == 0
}

func checkSentenceWords(s string) string {
	words := strings.Fields(s)
	if len(words) > MaxSTE100WordsPerSentence {
		return fmt.Sprintf("frase excede %d palavras (encontradas %d): %q", MaxSTE100WordsPerSentence, len(words), s)
	}
	return ""
}

func checkHedging(s string) string {
	lower := strings.ToLower(s)
	for _, term := range hedgingTerms {
		if strings.Contains(lower, term) {
			return fmt.Sprintf("termo de indefinição proibido (hedging): %q em %q", term, s)
		}
	}
	return ""
}

func splitSentences(text string) []string {
	var sentences []string
	clean := strings.ReplaceAll(text, "\r\n", "\n")
	blocks := strings.Split(clean, "\n")

	for _, block := range blocks {
		chunks := strings.FieldsFunc(block, func(r rune) bool {
			return r == '.' || r == '!' || r == '?' || r == ';'
		})
		sentences = append(sentences, chunks...)
	}

	return sentences
}
