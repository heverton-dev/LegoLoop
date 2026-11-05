package fabricar

import (
	"path/filepath"
	"testing"

	"legoloop/internal/matriz"
)

// TestCaminhoDoRefNoFormatoDaMatriz: os ref da matriz real são "~/caminho" e "arquivo#nome"; o caminho
// resolvido fica dentro do HOME, sem o "#nome".
func TestCaminhoDoRefNoFormatoDaMatriz(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	casos := map[string]string{
		"~/.claude/skills/x/SKILL.md":   filepath.Join(home, ".claude", "skills", "x", "SKILL.md"),
		`~\.claude\skills\x\SKILL.md`:   filepath.Join(home, `\.claude\skills\x\SKILL.md`),
		"~/.claude.json#designi-mcp":    filepath.Join(home, ".claude.json"),
		"~/.claude/settings.json#Pre:*": filepath.Join(home, ".claude", "settings.json"),
		"~":                             home,
	}
	for ref, want := range casos {
		if got := caminhoDoRef(ref); filepath.Clean(got) != filepath.Clean(want) {
			t.Errorf("caminhoDoRef(%q) = %q, quero %q", ref, got, want)
		}
	}
}

// TestConfiaveisAceitaSkillDoHomeEJulgaSoSkills: skill no HOME (disco fixo) fica; MCP e hook nunca viram
// aviso de mídia removível; skill fora de disco fixo sai com aviso.
func TestConfiaveisAceitaSkillDoHomeEJulgaSoSkills(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	es := []matriz.Entrada{
		{Ref: "~/.claude/skills/x/SKILL.md", Nome: "x", Tipo: "skill"},
		{Ref: "~/.claude.json#mcp", Nome: "mcp", Tipo: "mcp"},
		{Ref: "/media/usb/hook.json#h", Nome: "hook-usb", Tipo: "hook"},
		{Ref: "/media/usb/s/SKILL.md", Nome: "usb", Tipo: "skill"},
	}
	kept, cortes := confiaveis(es, false)
	if len(kept) != 3 || len(cortes) != 1 || kept[0].Nome != "x" {
		t.Fatalf("esperava manter x, mcp e hook-usb e cortar usb: kept=%v cortes=%v", nomesDe(kept), cortes)
	}
	if kept, cortes := confiaveis(es, true); len(kept) != 4 || len(cortes) != 0 {
		t.Errorf("com trust nada sai: kept=%v cortes=%v", nomesDe(kept), cortes)
	}
}

func nomesDe(es []matriz.Entrada) []string {
	out := make([]string, len(es))
	for i, e := range es {
		out[i] = e.Nome
	}
	return out
}
