package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/skills"
)

func skillCusto(t *testing.T, base, nome string) {
	t.Helper()
	dir := filepath.Join(base, nome)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	txt := "---\nname: " + nome + "\ndescription: " + strings.Repeat("y", 80) + "\n---\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(txt), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestSkillsCusto prova que o comando mede o agy pelas 3 pastas que o cli.log do agy escaneia
// e que a mesma skill em duas pastas é contada uma vez.
func TestSkillsCusto(t *testing.T) {
	d, out, errb := deps(t)
	g := filepath.Join(d.Home, ".gemini")
	skillCusto(t, filepath.Join(g, "antigravity-cli", "skills"), "alfa")
	skillCusto(t, filepath.Join(g, "config", "skills"), "alfa")
	skillCusto(t, filepath.Join(g, "skills"), "beta.disabled")

	if code := run(t, d, "skills-custo", "--json"); code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errb.String())
	}
	var rel struct {
		Harnesses []skills.CustoHarness `json:"harnesses"`
	}
	if err := json.Unmarshal(out.Bytes(), &rel); err != nil {
		t.Fatalf("saída não é JSON: %v\n%s", err, out.String())
	}
	var agy *skills.CustoHarness
	for i := range rel.Harnesses {
		if rel.Harnesses[i].Harness == "agy-cli" {
			agy = &rel.Harnesses[i]
		}
	}
	if agy == nil {
		t.Fatalf("faltou o harness agy-cli em %+v", rel.Harnesses)
	}
	if agy.Skills != 1 || agy.Duplicadas != 1 || agy.Disabled != 1 || len(agy.Pastas) != 3 {
		t.Fatalf("agy-cli = %+v, quer 1 skill (alfa), 1 duplicada, 1 disabled fora do total, 3 pastas", *agy)
	}
}

func TestSkillsCustoTexto(t *testing.T) {
	d, out, _ := deps(t)
	if code := run(t, d, "skills-custo"); code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out.String(), "agy-cli") || !strings.Contains(out.String(), "tokens") {
		t.Fatalf("saída em texto sem agy-cli/tokens: %q", out.String())
	}
}

// TestSkillsCustoCobreTodosOsHarnesses prova que nenhum harness com pasta de skills própria
// fica fora da medição: um harness eager (cline, 2 skills) e um de layout especial (opencode) aparecem.
func TestSkillsCustoCobreTodosOsHarnesses(t *testing.T) {
	d, out, errb := deps(t)
	skillCusto(t, filepath.Join(d.Home, ".cline", "skills"), "alfa")
	skillCusto(t, filepath.Join(d.Home, ".cline", "skills"), "beta")
	skillCusto(t, filepath.Join(d.Home, ".config", "opencode", "skills"), "gama")
	if code := run(t, d, "skills-custo", "--json"); code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errb.String())
	}
	var rel struct {
		Harnesses []skills.CustoHarness `json:"harnesses"`
	}
	if err := json.Unmarshal(out.Bytes(), &rel); err != nil {
		t.Fatalf("saída não é JSON: %v", err)
	}
	got := map[string]int{}
	for _, h := range rel.Harnesses {
		got[h.Harness] = h.Skills
	}
	if got["cline"] != 2 || got["opencode"] != 1 {
		t.Fatalf("cobertura = %v, quer cline=2 e opencode=1", got)
	}
	for _, nome := range []string{"mimocode", "hermes", "windsurf", "pi", "grok", "crush", "goose", "junie"} {
		if _, ok := got[nome]; !ok {
			t.Errorf("harness %q fora do skills-custo", nome)
		}
	}
}
