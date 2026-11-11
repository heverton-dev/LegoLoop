package interview

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

// recLLM registra cada lista de mensagens recebida e responde sempre com a mesma resposta.
type recLLM struct {
	calls [][]Message
	reply string
}

func (r *recLLM) Chat(_ context.Context, msgs []Message) (string, error) {
	r.calls = append(r.calls, append([]Message(nil), msgs...))
	return r.reply, nil
}

const lowTurn = "HIPÓTESE: refinar o build\nCONFIANÇA: 40%\nPERGUNTA: Qual etapa falha?\nPALPITE: compilar"

func writeDoc(t *testing.T, dir, name, text string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

func allText(msgs []Message) string {
	var b strings.Builder
	for _, m := range msgs {
		b.WriteString(m.Content + "\n")
	}
	return b.String()
}

func TestContextoPromptCarregaProjetoPecasEObjetivo(t *testing.T) {
	dir := t.TempDir()
	writeDoc(t, dir, "AGENTS.md", "# Regras\nGo, stdlib only. Sem dependências externas.")
	llm := &recLLM{reply: lowTurn}
	eng := NewEngine(llm).WithContext(Context{
		Goal:   "refinar a etapa de BUILD do mega-processo",
		Docs:   LoadProjectContext(dir),
		Pieces: []string{"legoloop-tdd", "go-build-check"},
	})
	if _, _, err := eng.Start(context.Background(), "refinar o build"); err != nil {
		t.Fatal(err)
	}
	if len(llm.calls) != 1 {
		t.Fatalf("esperava 1 chamada, houve %d", len(llm.calls))
	}
	got := allText(llm.calls[0])
	for _, want := range []string{"Go, stdlib only", "AGENTS.md", "legoloop-tdd", "go-build-check", "refinar a etapa de BUILD do mega-processo"} {
		if !strings.Contains(got, want) {
			t.Errorf("o prompt não contém %q:\n%s", want, got)
		}
	}
}

func TestContextoPromptProibeTecnologiaForaDoProjeto(t *testing.T) {
	llm := &recLLM{reply: lowTurn}
	eng := NewEngine(llm).WithContext(Context{Goal: "x", Docs: []Snippet{{Name: "AGENTS.md", Text: "Go, stdlib only"}}})
	if _, _, err := eng.Start(context.Background(), "x"); err != nil {
		t.Fatal(err)
	}
	sys := llm.calls[0][0]
	if sys.Role != "system" || !strings.Contains(sys.Content, SystemPromptCanonical) {
		t.Fatalf("a primeira mensagem deve manter o prompt canônico:\n%s", sys.Content)
	}
	if !strings.Contains(sys.Content, "NÃO sugira tecnologia, ferramenta ou serviço que o contexto do projeto não cite") {
		t.Errorf("faltou a regra contra tecnologia fora do projeto:\n%s", sys.Content)
	}
}

func TestContextoVazioNaoMudaOPrompt(t *testing.T) {
	llm := &recLLM{reply: lowTurn}
	if _, _, err := NewEngine(llm).Start(context.Background(), "x"); err != nil {
		t.Fatal(err)
	}
	if got := llm.calls[0][0].Content; got != SystemPromptCanonical {
		t.Errorf("sem contexto o prompt do sistema deve ficar igual ao canônico, veio:\n%s", got)
	}
}

func TestContextoLoadLeNaOrdemEIgnoraOQueFalta(t *testing.T) {
	dir := t.TempDir()
	writeDoc(t, dir, "CLAUDE.md", "claude")
	writeDoc(t, dir, "AGENTS.md", "agents")
	if err := os.Mkdir(filepath.Join(dir, "README.md"), 0o755); err != nil { // pasta com nome de doc não vale
		t.Fatal(err)
	}
	got := LoadProjectContext(dir)
	if len(got) != 2 || got[0].Name != "AGENTS.md" || got[1].Name != "CLAUDE.md" {
		t.Fatalf("ordem ou conteúdo inesperado: %+v", got)
	}
	if len(LoadProjectContext(filepath.Join(dir, "nao-existe"))) != 0 {
		t.Error("pasta inexistente deve dar contexto vazio, sem erro")
	}
}

func TestContextoLoadCortaSemQuebrarCaractere(t *testing.T) {
	dir := t.TempDir()
	writeDoc(t, dir, "AGENTS.md", "x"+strings.Repeat("ç", MaxSnippetBytes)) // o limite cai no meio de um "ç"
	got := LoadProjectContext(dir)
	if len(got) != 1 {
		t.Fatalf("esperava 1 trecho, veio %d", len(got))
	}
	if len(got[0].Text) > MaxSnippetBytes {
		t.Errorf("trecho com %d bytes passa do limite %d", len(got[0].Text), MaxSnippetBytes)
	}
	if !utf8.ValidString(got[0].Text) {
		t.Error("o corte quebrou um caractere UTF-8")
	}
	if !got[0].Truncated {
		t.Error("o trecho cortado deve marcar Truncated")
	}
}

func TestContextoRenderMarcaCorte(t *testing.T) {
	out := Context{Docs: []Snippet{{Name: "AGENTS.md", Text: "abc", Truncated: true}}}.Render()
	if !strings.Contains(out, "trecho cortado") {
		t.Errorf("o prompt deve avisar que o trecho foi cortado:\n%s", out)
	}
}

func TestContextoLoadIgnoraDocumentoVazio(t *testing.T) {
	dir := t.TempDir()
	writeDoc(t, dir, "AGENTS.md", "   ")
	writeDoc(t, dir, "CLAUDE.md", "regra")
	got := LoadProjectContext(dir)
	if len(got) != 1 || got[0].Name != "CLAUDE.md" {
		t.Errorf("documento só com espaços não deve entrar no prompt: %+v", got)
	}
}
