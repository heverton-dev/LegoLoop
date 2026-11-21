package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func gravaSkill(t *testing.T, base, nome, conteudo string) {
	t.Helper()
	dir := filepath.Join(base, nome)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fm(nome string) string {
	return "---\nname: " + nome + "\ndescription: " + strings.Repeat("x", 100) + "\n---\ncorpo\n"
}

// TestCustoGrupo prova a regra do agy: união das pastas deduplicada por nome,
// frontmatter inválido contado à parte e *.disabled fora do total (experimento de 20/11).
func TestCustoGrupo(t *testing.T) {
	home := t.TempDir()
	a := filepath.Join(home, "a")
	b := filepath.Join(home, "b")
	gravaSkill(t, a, "um", fm("um"))
	gravaSkill(t, a, "dois", fm("dois"))
	gravaSkill(t, b, "um", fm("um"))
	gravaSkill(t, b, "velha.disabled", fm("velha"))
	gravaSkill(t, b, "quebrada", "sem frontmatter")

	c := CustoGrupo("agy", []string{a, b})

	if c.Skills != 2 {
		t.Fatalf("skills únicas = %d, quer 2 (um, dois); .disabled o agy ignora", c.Skills)
	}
	if c.Invalidas != 1 {
		t.Fatalf("inválidas = %d, quer 1", c.Invalidas)
	}
	if c.Disabled != 1 || c.TokensDisabled <= 0 {
		t.Fatalf("disabled = %d tokens %d, quer 1 com custo > 0", c.Disabled, c.TokensDisabled)
	}
	if c.Duplicadas != 1 {
		t.Fatalf("duplicadas = %d, quer 1 (um em a e b)", c.Duplicadas)
	}
	if c.Tokens <= 0 || c.Tokens != c.Pastas[0].Tokens {
		t.Fatalf("tokens = %d, quer só um e dois (%d); disabled fica fora do total", c.Tokens, c.Pastas[0].Tokens)
	}
	if len(c.Pastas) != 2 || c.Pastas[0].Skills != 2 || c.Pastas[1].Skills != 2 || c.Pastas[1].Invalidas != 1 {
		t.Fatalf("pastas = %+v, quer 2 válidas em a; 2 válidas e 1 inválida em b (por pasta é bruto; só o total exclui disabled)", c.Pastas)
	}
}

// TestCustoPastaAusente prova que pasta inexistente não derruba a medição.
func TestCustoPastaAusente(t *testing.T) {
	c := CustoGrupo("x", []string{filepath.Join(t.TempDir(), "nao-existe")})
	if c.Skills != 0 || c.Tokens != 0 {
		t.Fatalf("pasta ausente deve custar zero, veio %+v", c)
	}
}
