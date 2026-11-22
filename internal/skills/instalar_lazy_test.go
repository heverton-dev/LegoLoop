package skills_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/cli"
	"github.com/heverton-dev/LegoLoop/internal/processo"
)

func rodarInstalarArgs(t *testing.T, home, vault, dirProcessos string, extra ...string) int {
	t.Helper()
	var out, errb bytes.Buffer
	d := cli.Deps{Stdout: &out, Stderr: &errb, Home: home}
	args := append([]string{"instalar-skills", "--vault", vault, "--processos-dir", dirProcessos}, extra...)
	code := cli.Run(context.Background(), args, d)
	if code != 0 {
		t.Logf("out=%s err=%s", out.String(), errb.String())
	}
	return code
}

func nomesNaPasta(t *testing.T, dir string) []string {
	t.Helper()
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range ents {
		out = append(out, e.Name())
	}
	sort.Strings(out)
	return out
}

// TestLazyVinculaSoOIndice prova o conserto de 20/11: o padrão deixa só a skill-índice nas pastas
// dos harnesses; as skills ficam no Cofre, onde o skill_search as acha.
func TestLazyVinculaSoOIndice(t *testing.T) {
	_, dirProc, home, vault := montarCenario(t, []string{"claude", "codex", "gemini"})
	if code := rodarInstalarArgs(t, home, vault, dirProc); code != 0 {
		t.Fatalf("exit %d", code)
	}
	for _, h := range processo.DetectarHarnesses(processo.OSFS(), home) {
		got := nomesNaPasta(t, h.Skills)
		if len(got) != 1 || got[0] != processo.IndiceSkills {
			t.Errorf("%s: pasta de skills = %v, quer só [%s]", h.Nome, got, processo.IndiceSkills)
		}
		idx, err := os.ReadFile(filepath.Join(h.Skills, processo.IndiceSkills, "SKILL.md"))
		if err != nil || !strings.Contains(string(idx), "name: "+processo.IndiceSkills) || !strings.Contains(string(idx), "skill_search") || !strings.Contains(string(idx), "call_skill") {
			t.Errorf("%s: índice inválido (%v): %s", h.Nome, err, idx)
		}
	}
	for _, nome := range []string{"legoloop-fabricar", "meu-proc"} {
		if _, err := os.Stat(filepath.Join(vault, "skills", nome, "SKILL.md")); err != nil {
			t.Errorf("%s deve existir no Cofre: %v", nome, err)
		}
	}
}

// TestEagerRestauraOComportamentoAntigo prova a saída de emergência para harness sem MCP.
func TestEagerRestauraOComportamentoAntigo(t *testing.T) {
	_, dirProc, home, vault := montarCenario(t, []string{"claude", "gemini"})
	if code := rodarInstalarArgs(t, home, vault, dirProc, "--eager"); code != 0 {
		t.Fatalf("exit %d", code)
	}
	hs := processo.DetectarHarnesses(processo.OSFS(), home)
	if v := processo.HarnessesVinculados(processo.OSFS(), hs, "meu-proc"); len(v) != 2 {
		t.Errorf("--eager: meu-proc vinculado em %v, quer os 2 harnesses", v)
	}
}

// TestExcecaoEmArquivoFicaVisivel prova a lista de exceções em ~/.legoloop/skills-eager.txt.
func TestExcecaoEmArquivoFicaVisivel(t *testing.T) {
	_, dirProc, home, vault := montarCenario(t, []string{"claude"})
	dir := filepath.Join(home, ".legoloop")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "skills-eager.txt"), []byte("legoloop-fabricar\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := rodarInstalarArgs(t, home, vault, dirProc); code != 0 {
		t.Fatalf("exit %d", code)
	}
	hs := processo.DetectarHarnesses(processo.OSFS(), home)
	got := nomesNaPasta(t, hs[0].Skills)
	want := []string{"legoloop-fabricar", processo.IndiceSkills}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("pasta = %v, quer %v", got, want)
	}
}
