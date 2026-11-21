package renomear

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func planejarDocs(op Opcoes) ([]Alteracao, error) {
	info, err := os.Stat(op.Docs)
	if err != nil || !info.IsDir() {
		return nil, nil
	}
	var alts []Alteracao
	err = filepath.WalkDir(op.Docs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if ignorarPastaDoc(op.Docs, path, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !arquivoDocElegivel(op.Raiz, path) {
			return nil
		}
		dados, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if substituirReferencias(string(dados), op.Antigo, op.Novo) != string(dados) {
			alts = append(alts, Alteracao{Tipo: "documento", Origem: path, Destino: path, Detalhe: "atualizar referências"})
		}
		return nil
	})
	return alts, err
}

func aplicarDocs(op Opcoes, alts []Alteracao, bk *backupCtx) error {
	for _, a := range alts {
		if a.Tipo != "documento" {
			continue
		}
		dados, err := os.ReadFile(a.Origem)
		if err != nil {
			return err
		}
		atualizado := substituirReferencias(string(dados), op.Antigo, op.Novo)
		if atualizado == string(dados) {
			continue
		}
		if err := bk.salvar(a.Origem, "atualizar_documento"); err != nil {
			return err
		}
		if err := os.WriteFile(a.Origem, []byte(atualizado), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func ignorarPastaDoc(baseDocs, path, nome string) bool {
	if nome == ".git" || nome == ".backups" || nome == "node_modules" || nome == "bin" {
		return true
	}
	if nome == "archive" {
		return true
	}
	rel, err := filepath.Rel(baseDocs, path)
	if err == nil {
		norm := filepath.ToSlash(rel)
		if norm == "archive" || strings.HasPrefix(norm, "archive/") || strings.Contains(norm, "/archive/") {
			return true
		}
	}
	return false
}

func arquivoDocElegivel(raiz, path string) bool {
	absPath, _ := filepath.Abs(path)
	absReg, _ := filepath.Abs(caminhoRegistro(raiz))
	if strings.EqualFold(absPath, absReg) {
		return false
	}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md", ".json", ".txt", ".yaml", ".yml":
		return true
	default:
		return false
	}
}

// substituirReferencias substitui apenas referências estruturais de Processo,
// evitando corromper palavras comuns ou nomes de verbos da CLI.
func substituirReferencias(texto, antigo, novo string) string {
	esc := regexp.QuoteMeta(antigo)
	// Fronteira direita: não pode ser seguido por caractere de identificador/slug.
	fronteira := `([^a-zA-Z0-9_-]|$)`
	padroes := []string{
		`(processos[/\\])` + esc + fronteira,
		`(skills[/\\])` + esc + fronteira,
		`(legoloop\s+(?:exec|sessoes|anticorpos\s+(?:status|sync))\s+)` + esc + fronteira,
		`((?i:processo)\s+` + "`?" + `)` + esc + `(` + "`?" + `[^a-zA-Z0-9_-]|$)`,
		`("(?:ref|processo|slug|padrao_slug)"\s*:\s*")` + esc + `(")`,
		`(\b)` + esc + `(@[0-9])`,
	}
	out := texto
	for _, pat := range padroes {
		re := regexp.MustCompile(pat)
		out = re.ReplaceAllString(out, "${1}"+novo+"${2}")
	}
	return out
}
