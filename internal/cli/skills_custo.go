package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"legoloop/internal/skills"
)

// gruposSkillsCusto devolve, por harness, as pastas que ele escaneia.
// agy-cli: as 3 pastas citadas no ~/.gemini/antigravity-cli/cli.log (linhas de "Failed to parse skill file").
func gruposSkillsCusto(home string) map[string][]string {
	g := filepath.Join(home, ".gemini")
	return map[string][]string{
		"agy-cli": {filepath.Join(g, "antigravity-cli", "skills"), filepath.Join(g, "config", "skills"), filepath.Join(g, "skills")},
		"agy-ide": {filepath.Join(g, "antigravity-ide", "skills")},
		"claude":  {filepath.Join(home, ".claude", "skills")},
		"codex":   {filepath.Join(home, ".codex", "skills")},
		"agents":  {filepath.Join(home, ".agents", "skills")},
	}
}

var ordemSkillsCusto = []string{"agy-cli", "agy-ide", "claude", "codex", "agents"}

type skillsCustoRel struct {
	Harnesses []skills.CustoHarness `json:"harnesses"`
}

func cmdSkillsCusto(_ context.Context, args []string, d Deps) int {
	fs := newFlags("skills-custo", d)
	asJSON := fs.Bool("json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	home := d.Home
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	grupos := gruposSkillsCusto(home)
	rel := skillsCustoRel{}
	for _, nome := range ordemSkillsCusto {
		rel.Harnesses = append(rel.Harnesses, skills.CustoGrupo(nome, grupos[nome]))
	}
	if *asJSON {
		enc := json.NewEncoder(d.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(rel)
		return 0
	}
	imprimirSkillsCusto(d, rel)
	return 0
}

func imprimirSkillsCusto(d Deps, rel skillsCustoRel) {
	fmt.Fprintf(d.Stdout, "%-9s %7s %9s %9s %10s %9s\n", "harness", "skills", "tokens", "invalidas", "duplicadas", "disabled")
	for _, h := range rel.Harnesses {
		fmt.Fprintf(d.Stdout, "%-9s %7d %9d %9d %10d %9d\n", h.Harness, h.Skills, h.Tokens, h.Invalidas, h.Duplicadas, h.Disabled)
	}
	fmt.Fprintln(d.Stdout, "tokens = frontmatter + caminho por skill única, 4 bytes por token (estimativa)")
}
