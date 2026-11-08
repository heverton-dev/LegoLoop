package videorea

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MaterializeAndProject cria o bundle isolado e projeta pontes na raiz canônica.
func MaterializeAndProject(report *CurationReport, baseDir, repoRoot string) (*TargetBundle, error) {
	if report == nil || !report.Approved {
		return nil, ErrGateRejected
	}
	bundle := NewTargetBundle(baseDir, report.SuggestedName)
	if err := bundle.Materialize(); err != nil {
		return nil, err
	}
	if err := writeBundleArtifacts(bundle, report); err != nil {
		return nil, err
	}
	if err := projectCanonical(bundle, repoRoot); err != nil {
		return nil, err
	}
	return bundle, nil
}

func writeBundleArtifacts(b *TargetBundle, r *CurationReport) error {
	if err := writeLocalGovernance(b, r); err != nil {
		return err
	}
	if err := writeLocalBook(b, r); err != nil {
		return err
	}
	if err := writeLocalLoop(b, r); err != nil {
		return err
	}
	if err := writeLocalOriginalMatrix(b, r); err != nil {
		return err
	}
	return writeLocalSkill(b, r)
}

func writeLocalOriginalMatrix(b *TargetBundle, r *CurationReport) error {
	matrix := r.OriginalMatrix
	if matrix == nil {
		matrix = &OriginalMatrix{
			VideoTitle:       r.VideoTitle,
			TotalMonthlyCost: r.OriginalCost,
			Tools:            []OriginalTool{},
			OriginalSteps:    []string{},
		}
	}
	data, err := json.MarshalIndent(matrix, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.OriginalMatrixPath, append(data, '\n'), 0o644)
}

func writeLocalGovernance(b *TargetBundle, r *CurationReport) error {
	intent := fmt.Sprintf("# INTENT: %s\n\nCusto Soberano: %s\n", r.SuggestedName, r.SovereignCost)
	spec := fmt.Sprintf("# SPEC: %s\n\nSubstituicoes: %v\n", r.SuggestedName, r.Replacements)
	plan := fmt.Sprintf("# PLAN: %s\n", r.SuggestedName)
	route := fmt.Sprintf(`{"name": %q, "version": "1.0.0"}`, r.SuggestedName)

	items := map[string]string{
		filepath.Join(b.GovernanceDir, "INTENT.md"):  intent,
		filepath.Join(b.GovernanceDir, "SPEC.md"):    spec,
		filepath.Join(b.GovernanceDir, "PLAN.md"):    plan,
		filepath.Join(b.GovernanceDir, "ROUTE.json"): route,
	}
	for p, content := range items {
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func writeLocalBook(b *TargetBundle, r *CurationReport) error {
	partesDir := filepath.Join(b.BookDir, "partes")
	if err := os.MkdirAll(partesDir, 0o755); err != nil {
		return err
	}
	if err := writeBookParts(partesDir, b.Slug, r); err != nil {
		return err
	}
	return writeBookManifest(b.BookDir, b.Slug, r)
}

func writeBookParts(partesDir, slug string, r *CurationReport) error {
	parts := map[string]string{
		"00-frontmatter.md":         buildBookFrontmatter(slug, r),
		"01-fundacao.md":            buildBookFundacao(slug, r),
		"02-matriz-original.md":     buildBookMatriz(r),
		"03-guias-autorais.md":      buildBookGuias(slug),
		"04-arquitetura-dag.md":     buildBookDAG(slug),
		"05-operacao-e-execucao.md": buildBookOperacao(slug),
	}
	for filename, content := range parts {
		if err := os.WriteFile(filepath.Join(partesDir, filename), []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func writeBookManifest(bookDir, slug string, r *CurationReport) error {
	manifest := fmt.Sprintf(`{
  "versao_manifesto": "1.0",
  "titulo": "Pipeline Canônico de %s",
  "nome_base": %q,
  "pasta_partes": "partes",
  "template": "ativos/livro.typst",
  "partes": [
    "00-frontmatter.md",
    "01-fundacao.md",
    "02-matriz-original.md",
    "03-guias-autorais.md",
    "04-arquitetura-dag.md",
    "05-operacao-e-execucao.md"
  ],
  "revisoes": [],
  "hash_fonte": %q,
  "ultima_compilacao": "2026-11-08"
}
`, r.SuggestedName, slug, slug)
	return os.WriteFile(filepath.Join(bookDir, "livro.json"), []byte(manifest), 0o644)
}

func buildBookFrontmatter(slug string, r *CurationReport) string {
	return fmt.Sprintf(`---
title: "Pipeline Canônico de %s"
subtitle: "Engenharia Reversa Audiovisual e Governança Soberana LegoLoop"
author:
  - Ecossistema LegoLoop
date: "08/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "MANUAL OPERACIONAL DIDÁTICO · VIDEO AS CODE"
cover-title: "%s"
cover-accent: "Engenharia Reversa e Adoção Soberana"
cover-caption: "LIVRO DIDÁTICO LEGOPRESS"
---
# Introdução
Manual técnico do pipeline %s compilado pelo ecossistema LegoLoop.
`, r.SuggestedName, slug, slug)
}

func buildBookFundacao(slug string, r *CurationReport) string {
	return fmt.Sprintf(`# 1. Fundação e Princípios Soberanos
Este volume documenta a engenharia reversa metódica do processo %s.
Custo Original Identificado: %s.
Custo Soberano LegoLoop: %s.
`, slug, r.OriginalCost, r.SovereignCost)
}

func buildBookMatriz(r *CurationReport) string {
	return fmt.Sprintf(`# 2. Matriz Original e Substituições
Substituições identificadas na curadoria:
%v
`, r.Replacements)
}

func buildBookGuias(slug string) string {
	return fmt.Sprintf(`# 3. Guias Autorais e Metodologia
Detalhamento dos princípios de edição, sincronia e governança estrita de %s.
`, slug)
}

func buildBookDAG(slug string) string {
	return fmt.Sprintf(`# 4. Arquitetura do Pipeline DAG
Estrutura topológica, isolamento em scratch e stop-gates de revisão humana para %s.
`, slug)
}

func buildBookOperacao(slug string) string {
	return fmt.Sprintf(`# 5. Operação e Execução
Comandos de execução via CLI legoloop, triggers cognitivos e telemetria de produção.
`)
}

func writeLocalLoop(b *TargetBundle, r *CurationReport) error {
	dagJSON := fmt.Sprintf(`{"name": %q, "version": "1.0.0", "nodes": []}`, r.SuggestedName)
	if err := os.WriteFile(filepath.Join(b.LoopDir, "pipeline.dag.json"), []byte(dagJSON), 0o644); err != nil {
		return err
	}
	_, err := SynthesizeGoldenFixture(3.0, b.LoopDir)
	return err
}

func writeLocalSkill(b *TargetBundle, r *CurationReport) error {
	skillMD := fmt.Sprintf(`---
name: %s
description: Use quando precisar executar ou orquestrar o pipeline %s no ecossistema LegoLoop.
tags: [audiovisual, video-rea, pipeline, soberano]
---
# Skill Universal: %s

## Visão Geral
Pipeline operacional gerado via engenharia reversa com governança LegoLoop v2.2.0.

## Comandos e Triggers
- Chat: /legoloop-%s
- CLI: legoloop run %s
`, b.Slug, r.SuggestedName, b.Slug, b.Slug, b.Slug)
	return os.WriteFile(filepath.Join(b.SkillDir, "SKILL.md"), []byte(skillMD), 0o644)
}

func projectCanonical(b *TargetBundle, repoRoot string) error {
	if err := copyFile(filepath.Join(b.LoopDir, "pipeline.dag.json"), filepath.Join(repoRoot, "loops", b.Slug, "pipeline.dag.json")); err != nil {
		return err
	}
	if err := copyFile(filepath.Join(b.SkillDir, "SKILL.md"), filepath.Join(repoRoot, "skills", b.Slug, "SKILL.md")); err != nil {
		return err
	}
	return projectCanonicalBook(b.BookDir, filepath.Join(repoRoot, "docs", "books", b.Slug))
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func projectCanonicalBook(srcDir, dstDir string) error {
	if err := copyFile(filepath.Join(srcDir, "livro.json"), filepath.Join(dstDir, "livro.json")); err != nil {
		return err
	}
	srcPartes := filepath.Join(srcDir, "partes")
	dstPartes := filepath.Join(dstDir, "partes")
	entries, err := os.ReadDir(srcPartes)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			if err := copyFile(filepath.Join(srcPartes, entry.Name()), filepath.Join(dstPartes, entry.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}
