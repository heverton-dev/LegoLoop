package roteiro

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type RegistroTema struct {
	Slug      string `json:"slug"`
	Titulo    string `json:"titulo"`
	Data      string `json:"data"`
	TipoPauta string `json:"tipo_pauta,omitempty"`
}

type AcervoTemas struct {
	SchemaVersion int            `json:"schema_version"`
	Temas         []RegistroTema `json:"temas"`
}

// NormalizarSlug gera uma chave canonica estavel para deduplicacao de temas.
func NormalizarSlug(texto string) string {
	var b strings.Builder
	hifen := false
	for _, r := range strings.ToLower(strings.TrimSpace(texto)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			hifen = false
			continue
		}
		if !hifen && b.Len() > 0 {
			b.WriteByte('-')
			hifen = true
		}
	}
	return strings.TrimRight(b.String(), "-")
}

// Contem verifica se um titulo ou slug ja existe no acervo.
func (a AcervoTemas) Contem(titulo string) bool {
	alvo := NormalizarSlug(titulo)
	if alvo == "" {
		return false
	}
	for _, t := range a.Temas {
		if NormalizarSlug(t.Slug) == alvo || NormalizarSlug(t.Titulo) == alvo {
			return true
		}
	}
	return false
}

// CarregarAcervoTemas le acervo/temas.json, retornando acervo vazio se o arquivo nao existir.
func CarregarAcervoTemas(caminho string) (AcervoTemas, error) {
	dados, err := os.ReadFile(caminho)
	if os.IsNotExist(err) {
		return AcervoTemas{SchemaVersion: 1}, nil
	}
	if err != nil {
		return AcervoTemas{}, err
	}
	var ac AcervoTemas
	if err := json.Unmarshal(dados, &ac); err != nil {
		return AcervoTemas{}, err
	}
	if ac.SchemaVersion == 0 {
		ac.SchemaVersion = 1
	}
	return ac, nil
}

// SalvarAcervoTemas grava acervo/temas.json de forma atomica (temp + rename).
func SalvarAcervoTemas(caminho string, ac AcervoTemas) error {
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		return err
	}
	if ac.SchemaVersion == 0 {
		ac.SchemaVersion = 1
	}
	dados, err := json.MarshalIndent(ac, "", "  ")
	if err != nil {
		return err
	}
	tmp := caminho + ".tmp"
	if err := os.WriteFile(tmp, append(dados, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, caminho)
}
