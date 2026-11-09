package gate

import "testing"

func TestGateTrivialChecksContent(t *testing.T) {
	for _, c := range []struct {
		src  string
		want bool
	}{
		{"file_exists('a.md')", false},
		{"exit_code == 0", false},
		{"exit_code == 0 && file_exists('a.md')", false},
		{"!file_exists('a.md')", false},
		{"file_exists('a.md') && word_count('a.md') >= 30", true},
		{"file_exists('a.md') && 30 <= word_count('a.md')", true}, // o número pode vir à esquerda
		{"word_count('a.md') >= 30 && file_exists('a.md')", true},
		{"file_exists('a.md') && !has_pattern('a.md', '\\{\\{')", true},
		{"file_exists('a.md') && has_binary_criteria('a.md')", true},
		{"file_exists('a.md') && has_acceptance_criteria('a.md')", true},
		{"file_exists('a.md') && line_count('a.md') >= 2", true},
		{"file_exists('a.json') && schema_valid('s.json', 'a.json')", true},
		{"file_exists('a.go') && no_stubs('a.go')", true},
		{"prompt('Aprova?')", true},
		{"interactive('Aprova?')", true},
		// "||" só prova conteúdo se os dois lados provam
		{"file_exists('a.md') || word_count('a.md') >= 30", false},
		{"word_count('a.md') >= 30 || has_pattern('a.md', 'x')", true},
		// parênteses não escondem a existência
		{"(file_exists('a.md'))", false},
		{"(file_exists('a.md') && exit_code == 0) && file_exists('b.md')", false},
		{"file_exists('a.md') && (exit_code == 0 || word_count('a.md') >= 30)", false},
		{"file_exists('a.md') && (word_count('a.md') >= 30 || has_pattern('a.md', 'x'))", true},
	} {
		e, err := Compile(c.src)
		if err != nil {
			t.Fatalf("%s não compila: %v", c.src, err)
		}
		if got := e.ChecksContent(); got != c.want {
			t.Errorf("ChecksContent(%q) = %v, queria %v", c.src, got, c.want)
		}
	}
}
