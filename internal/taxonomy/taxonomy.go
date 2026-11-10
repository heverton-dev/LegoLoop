package taxonomy

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

//go:embed dictionary.json
var dictionaryData []byte

type TermEntry struct {
	ApprovedTerm      string   `json:"approved_term"`
	GrammaticalClass  string   `json:"grammatical_class"`
	StrictDefinition  string   `json:"strict_definition"`
	ForbiddenSynonyms []string `json:"forbidden_synonyms"`
	ContextDomain     string   `json:"context_domain"`
}

type Dictionary struct {
	Version     string      `json:"version"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Terms       []TermEntry `json:"terms"`
}

type TaxonomyViolation struct {
	ForbiddenTerm string `json:"forbidden_term"`
	ApprovedTerm  string `json:"approved_term"`
	ContextDomain string `json:"context_domain"`
	Line          int    `json:"line"`
	Column        int    `json:"column"`
}

var (
	dictOnce    sync.Once
	cachedDict  *Dictionary
	synonymMap  map[string]TermEntry
	dictLoadErr error
)

func loadDictionary() (*Dictionary, map[string]TermEntry, error) {
	dictOnce.Do(func() {
		var dict Dictionary
		if err := json.Unmarshal(dictionaryData, &dict); err != nil {
			dictLoadErr = fmt.Errorf("falha ao decodificar dicionario: %w", err)
			return
		}
		cachedDict = &dict
		synonymMap = make(map[string]TermEntry)
		for _, entry := range dict.Terms {
			for _, syn := range entry.ForbiddenSynonyms {
				synonymMap[strings.ToLower(strings.TrimSpace(syn))] = entry
			}
		}
	})
	return cachedDict, synonymMap, dictLoadErr
}

// CheckForbidden verifica se o termo fornecido é um sinônimo proibido.
func CheckForbidden(term string) (string, bool) {
	_, sMap, err := loadDictionary()
	if err != nil {
		return "", false
	}
	entry, ok := sMap[strings.ToLower(strings.TrimSpace(term))]
	if !ok {
		return "", false
	}
	return entry.ApprovedTerm, true
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func isWordBoundary(runes []rune, start, length int) bool {
	if start > 0 && isWordRune(runes[start-1]) {
		return false
	}
	end := start + length
	if end < len(runes) && isWordRune(runes[end]) {
		return false
	}
	return true
}

func scanLineForEntry(lineRunes []rune, lowerLine string, entry TermEntry, lineNum int) []TaxonomyViolation {
	var violations []TaxonomyViolation
	for _, syn := range entry.ForbiddenSynonyms {
		synLower := strings.ToLower(syn)
		synLen := utf8.RuneCountInString(syn)
		pos := 0
		for {
			idx := strings.Index(lowerLine[pos:], synLower)
			if idx < 0 {
				break
			}
			matchBytePos := pos + idx
			runeStart := utf8.RuneCountInString(lowerLine[:matchBytePos])
			if isWordBoundary(lineRunes, runeStart, synLen) {
				violations = append(violations, TaxonomyViolation{
					ForbiddenTerm: syn,
					ApprovedTerm:  entry.ApprovedTerm,
					ContextDomain: entry.ContextDomain,
					Line:          lineNum,
					Column:        runeStart + 1,
				})
			}
			pos = matchBytePos + len(synLower)
		}
	}
	return violations
}

// ValidateText analisa o texto e identifica todas as violações taxonômicas.
func ValidateText(content string) ([]TaxonomyViolation, error) {
	dict, _, err := loadDictionary()
	if err != nil {
		return nil, err
	}
	var allViolations []TaxonomyViolation
	lines := strings.Split(content, "\n")
	for lineIdx, line := range lines {
		lineNum := lineIdx + 1
		lineRunes := []rune(line)
		lowerLine := strings.ToLower(line)
		for _, entry := range dict.Terms {
			v := scanLineForEntry(lineRunes, lowerLine, entry, lineNum)
			allViolations = append(allViolations, v...)
		}
	}
	return allViolations, nil
}
