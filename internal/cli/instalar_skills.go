package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/heverton-dev/LegoLoop/internal/skills"
)

type instalarSkillsResult = skills.Relatorio

func cmdInstalarSkills(_ context.Context, args []string, d Deps) int {
	fs := newFlags("instalar-skills", d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	procDir := fs.String("processos-dir", "", "pasta dos Processos (padrão processos)")
	dry := fs.Bool("dry-run", false, "apenas planeja sem gravar")
	eager := fs.Bool("eager", false, "vincula toda skill em toda pasta de harness (padrão: só o índice legoloop-skills; gasta contexto)")
	asJSON := fs.Bool("json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	v, err := resolveVault(*vault)
	if err != nil {
		return fail(d, err)
	}
	res, err := runInstalarSkills(d, v, *procDir, *dry, *eager)
	if err != nil {
		return fail(d, err)
	}
	return reportInstalarSkills(d, res, *asJSON)
}

func runInstalarSkills(d Deps, vault, dirProcessos string, dry, eager bool) (*instalarSkillsResult, error) {
	home := d.Home
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	return skills.Instalar(skills.Opcao{
		Vault: vault, Home: home, DirProcessos: dirProcessos, DryRun: dry, Eager: eager,
	})
}

func instalarSkillsSilencioso(d Deps, vault string) {
	_, _ = runInstalarSkills(d, vault, "", false, false)
}

func reportInstalarSkills(d Deps, res *instalarSkillsResult, asJSON bool) int {
	if asJSON {
		enc := json.NewEncoder(d.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(res)
		return 0
	}
	prefix := "INSTALADO"
	if res.DryRun {
		prefix = "DRY-RUN"
	}
	fmt.Fprintf(d.Stdout, "%s — %d skills universais instaladas no cofre\n", prefix, len(res.Skills))
	fmt.Fprintf(d.Stdout, "Cofre: %s\n", res.Vault)
	if len(res.Harnesses) > 0 {
		fmt.Fprintf(d.Stdout, "Harnesses vinculados: %v\n", res.Harnesses)
	}
	if len(res.Comandos) > 0 {
		fmt.Fprintf(d.Stdout, "Comandos distribuídos: %v\n", res.Comandos)
	}
	return 0
}
