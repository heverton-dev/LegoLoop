package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"legoloop/internal/processo"
	"legoloop/skills"
)

type instalarSkillsResult struct {
	Vault     string   `json:"vault"`
	Skills    []string `json:"skills"`
	Harnesses []string `json:"harnesses"`
	Vinculos  []string `json:"vinculos"`
	DryRun    bool     `json:"dry_run"`
	Status    string   `json:"status"`
}

func cmdInstalarSkills(_ context.Context, args []string, d Deps) int {
	fs := newFlags("instalar-skills", d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	dry := fs.Bool("dry-run", false, "apenas planeja sem gravar")
	asJSON := fs.Bool("json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	v, err := resolveVault(*vault)
	if err != nil {
		return fail(d, err)
	}
	res, err := runInstalarSkills(d, v, *dry)
	if err != nil {
		return fail(d, err)
	}
	return reportInstalarSkills(d, res, *asJSON)
}

func runInstalarSkills(d Deps, vault string, dry bool) (*instalarSkillsResult, error) {
	home := d.Home
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	hs := processo.DetectarHarnesses(processo.OSFS(), home)
	hNames := make([]string, 0, len(hs))
	for _, h := range hs {
		hNames = append(hNames, h.Nome)
	}
	allSkills, err := skills.TodasAsSkills()
	if err != nil {
		return nil, fmt.Errorf("carregar skills nativas: %w", err)
	}
	skNames, vinculos, err := aplicarSkills(vault, hs, allSkills, dry)
	if err != nil {
		return nil, err
	}
	return &instalarSkillsResult{
		Vault: vault, Skills: skNames, Harnesses: hNames, Vinculos: vinculos, DryRun: dry, Status: "ok",
	}, nil
}

func aplicarSkills(vault string, hs []processo.Harness, allSkills map[string][]byte, dry bool) ([]string, []string, error) {
	var skNames, vinculos []string
	for _, fn := range skills.FuncoesNativas {
		skNames = append(skNames, fn)
		dirCofre := filepath.Join(vault, "skills", fn)
		if !dry {
			if err := gravarSkillNoCofre(dirCofre, allSkills[fn]); err != nil {
				return nil, nil, err
			}
		}
		vs := processo.PlanejarVinculos(hs, fn, dirCofre)
		for _, v := range vs {
			vinculos = append(vinculos, v.Destino)
		}
		if !dry && len(vs) > 0 {
			if _, err := processo.Vincular(vs, nil); err != nil {
				return nil, nil, fmt.Errorf("vincular %s: %w", fn, err)
			}
		}
	}
	return skNames, vinculos, nil
}

func gravarSkillNoCofre(dirCofre string, data []byte) error {
	if err := os.MkdirAll(dirCofre, 0o755); err != nil {
		return fmt.Errorf("criar pasta do cofre %s: %w", dirCofre, err)
	}
	dest := filepath.Join(dirCofre, "SKILL.md")
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("gravar %s: %w", dest, err)
	}
	return nil
}

func instalarSkillsSilencioso(d Deps, vault string) {
	_, _ = runInstalarSkills(d, vault, false)
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
	return 0
}
