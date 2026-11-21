package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"legoloop/internal/skills"
)

func cmdSkillsMigrar(_ context.Context, args []string, d Deps) int {
	fs := newFlags("skills-migrar", d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	dry := fs.Bool("dry-run", false, "apenas conta o que mudaria")
	apply := fs.Bool("apply", false, "move as skills eager de cada harness para o Cofre (com backup)")
	asJSON := fs.Bool("json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *dry == *apply {
		fmt.Fprintln(d.Stderr, "Uso: legoloop skills-migrar (--dry-run|--apply) [--vault DIR] [--json]")
		return 2
	}
	v, err := resolveVault(*vault)
	if err != nil {
		return fail(d, err)
	}
	home := d.Home
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	rel, err := skills.Migrar(skills.MigrarOpcao{Home: home, Vault: v, Data: time.Now().Format("20060102"), Aplicar: *apply})
	if err != nil {
		return fail(d, err)
	}
	return imprimirSkillsMigrar(d, rel, *apply, *asJSON)
}

func imprimirSkillsMigrar(d Deps, rel skills.MigrarRelatorio, aplicado, asJSON bool) int {
	if asJSON {
		enc := json.NewEncoder(d.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(rel)
		return 0
	}
	fmt.Fprintf(d.Stdout, "%-10s %7s %9s %9s\n", "harness", "mover", "ao cofre", "invalidas")
	for _, h := range rel.Harnesses {
		fmt.Fprintf(d.Stdout, "%-10s %7d %9d %9d\n", h.Nome, h.Mover, h.Copiar, h.Invalidas)
	}
	if aplicado {
		fmt.Fprintf(d.Stdout, "backup com manifesto em %s\n", rel.Backup)
	}
	return 0
}
