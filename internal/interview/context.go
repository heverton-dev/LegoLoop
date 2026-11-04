package interview

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// MaxSnippetBytes é o tamanho máximo de cada trecho de documento que vai ao prompt.
const MaxSnippetBytes = 2000

// projectDocs são os documentos lidos do projeto, nesta ordem.
var projectDocs = []string{"AGENTS.md", "CLAUDE.md", "README.md"}

// Snippet é um trecho de documento do projeto. Truncated avisa que o texto foi cortado.
type Snippet struct {
	Name      string
	Text      string
	Truncated bool
}

// Context é o que a entrevista sabe do projeto antes de fazer a primeira pergunta.
type Context struct {
	Goal   string    // objetivo dito pelo usuário
	Docs   []Snippet // trechos de AGENTS.md, CLAUDE.md e README.md
	Pieces []string  // nomes das peças aptas da matriz
}

// Empty diz se não há nada a acrescentar ao prompt.
func (c Context) Empty() bool {
	return strings.TrimSpace(c.Goal) == "" && len(c.Docs) == 0 && len(c.Pieces) == 0
}

// WithContext define o contexto das próximas sessões e devolve o próprio motor.
func (e *Engine) WithContext(c Context) *Engine {
	e.ctx = c
	return e
}

// LoadProjectContext lê os documentos do projeto que existem em dir. Pasta ou arquivo ausente é ignorado.
func LoadProjectContext(dir string) []Snippet {
	var out []Snippet
	for _, name := range projectDocs {
		data, err := readDoc(filepath.Join(dir, name))
		if err != nil || strings.TrimSpace(string(data)) == "" {
			continue
		}
		text, cut := cutUTF8(string(data), MaxSnippetBytes)
		out = append(out, Snippet{Name: name, Text: text, Truncated: cut})
	}
	return out
}

func readDoc(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s é uma pasta", path)
	}
	return os.ReadFile(path)
}

// cutUTF8 corta s em no máximo max bytes, sem partir um caractere ao meio.
func cutUTF8(s string, max int) (string, bool) {
	if len(s) <= max {
		return s, false
	}
	end := max
	for end > 0 && !utf8.RuneStart(s[end]) {
		end--
	}
	return s[:end], true
}

const noOutsideTech = "NÃO sugira tecnologia, ferramenta ou serviço que o contexto do projeto não cite. " +
	"Se o contexto não diz como algo é feito, pergunte ao usuário em vez de supor."

// Render monta o bloco de contexto que se soma ao prompt do sistema.
func (c Context) Render() string {
	var b strings.Builder
	b.WriteString("CONTEXTO DO PROJETO (use-o para fazer perguntas específicas):\n")
	if g := strings.TrimSpace(c.Goal); g != "" {
		fmt.Fprintf(&b, "Objetivo do usuário: %s\n", g)
	}
	for _, d := range c.Docs {
		fmt.Fprintf(&b, "\n--- %s ---\n%s\n", d.Name, d.Text)
		if d.Truncated {
			b.WriteString("(trecho cortado: o documento é maior)\n")
		}
	}
	if len(c.Pieces) > 0 {
		fmt.Fprintf(&b, "\nPeças aptas já catalogadas: %s\n", strings.Join(c.Pieces, ", "))
	}
	b.WriteString("\n" + noOutsideTech)
	return b.String()
}

// systemPrompt é o prompt canônico, acrescido do contexto quando houver.
func (e *Engine) systemPrompt() string {
	if e.ctx.Empty() {
		return SystemPromptCanonical
	}
	return SystemPromptCanonical + "\n\n" + e.ctx.Render()
}
