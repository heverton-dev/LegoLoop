package fabricar

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"legoloop/internal/processo"
)

// ErrExists: já existe um processo com o mesmo nome canônico e Force=false.
var ErrExists = fmt.Errorf("fabricar: %w", processo.ErrExiste)

// Roots, FS e PlannedFile vêm do núcleo: a gravação atômica mora em processo.Store.
type (
	Roots       = processo.Roots
	FS          = processo.FS
	PlannedFile = processo.File
)

// OSFS devolve o FS real, sobre o disco.
func OSFS() FS { return processo.OSFS() }

// MaterializePlan descreve a gravação de um processo.
type MaterializePlan struct {
	Name   string
	Roots  Roots
	Files  []PlannedFile
	Force  bool // sobrescreve após preservar o original como .legoloop-bak
	DryRun bool // calcula tudo e não grava nada

	Vinculos []processo.Vinculo // ligações da skill aos harnesses; criadas junto com a gravação, ou nenhuma
	Link     processo.LinkFunc  // nil = symlink, com cópia do SKILL.md se o sistema não permitir
}

// CanonicalArtifacts consolida os caminhos gerados.
type CanonicalArtifacts struct {
	IntentPath   string   `json:"intent_path"`
	SpecPath     string   `json:"spec_path"`
	PlanPath     string   `json:"plan_path"`
	RoutePath    string   `json:"route_path"`
	TemplatePath string   `json:"template_path"`
	SkillPath    string   `json:"skill_path"`
	LoopPath     string   `json:"loop_path"`
	DagPath      string   `json:"dag_path"`
	SkillLinks   []string `json:"skill_links"` // links nas pastas dos harnesses (no dry-run, os que seriam criados; vazio com --no-link)
}

// Materialize grava todos os arquivos do plano de forma atômica, delegando ao núcleo (processo.Store).
// Qualquer falha desfaz o que foi feito e deixa o estado anterior.
func Materialize(ctx context.Context, plan MaterializePlan, fsys FS) (CanonicalArtifacts, error) {
	if err := ctx.Err(); err != nil {
		return CanonicalArtifacts{}, err
	}
	art := artifactsOf(plan.Files)
	st := processo.Store{Roots: plan.Roots, FS: fsys}
	opts := processo.WriteOpts{Force: plan.Force, DryRun: plan.DryRun}
	if !plan.DryRun && len(plan.Vinculos) > 0 {
		opts.Depois = func() error {
			links, err := processo.Vincular(plan.Vinculos, plan.Link)
			art.SkillLinks = links
			return err
		}
	}
	err := st.Write(ctx, plan.Files, opts)
	switch {
	case err == nil:
		if plan.DryRun {
			art.SkillLinks = linkDests(plan.Vinculos)
		}
		return art, nil
	case errors.Is(err, processo.ErrExiste):
		return art, fmt.Errorf("%w: %s", ErrExists, strings.TrimPrefix(err.Error(), processo.ErrExiste.Error()+": "))
	}
	return CanonicalArtifacts{}, err
}

func artifactsOf(files []PlannedFile) CanonicalArtifacts {
	var a CanonicalArtifacts
	for _, f := range files {
		switch f.Kind {
		case "intent":
			a.IntentPath = f.Path
		case "spec":
			a.SpecPath = f.Path
		case "plan":
			a.PlanPath = f.Path
		case "route":
			a.RoutePath = f.Path
		case "template":
			a.TemplatePath = f.Path
		case "skill":
			a.SkillPath = f.Path
		case "loop":
			a.LoopPath = f.Path
		case "dag":
			a.DagPath = f.Path
		}
	}
	return a
}
