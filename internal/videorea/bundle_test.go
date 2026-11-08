package videorea_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/videorea"
)

func TestTargetBundlePaths(t *testing.T) {
	baseDir := filepath.Join("docs", "videorea")
	slug := "legoloop-video-tutorial-3d"

	bundle := videorea.NewTargetBundle(baseDir, slug)
	if bundle.Slug != slug {
		t.Fatalf("esperava slug %q, obteve %q", slug, bundle.Slug)
	}

	expectedBase := filepath.Join(baseDir, slug)
	if bundle.BaseDir != expectedBase {
		t.Errorf("BaseDir esperado %q, obteve %q", expectedBase, bundle.BaseDir)
	}
	if bundle.ScratchDir != filepath.Join(expectedBase, "scratch") {
		t.Errorf("ScratchDir incorreto: %s", bundle.ScratchDir)
	}
	if bundle.GovernanceDir != filepath.Join(expectedBase, "governance") {
		t.Errorf("GovernanceDir incorreto: %s", bundle.GovernanceDir)
	}
	if bundle.BookDir != filepath.Join(expectedBase, "book") {
		t.Errorf("BookDir incorreto: %s", bundle.BookDir)
	}
	if bundle.LoopDir != filepath.Join(expectedBase, "loop") {
		t.Errorf("LoopDir incorreto: %s", bundle.LoopDir)
	}
	if bundle.OriginalMatrixPath != filepath.Join(expectedBase, "original_matrix.json") {
		t.Errorf("OriginalMatrixPath incorreto: %s", bundle.OriginalMatrixPath)
	}
	if bundle.SkillDir != filepath.Join(expectedBase, "skill") {
		t.Errorf("SkillDir incorreto: %s", bundle.SkillDir)
	}
}

func TestTargetBundleMaterializeAndProject(t *testing.T) {
	tempRoot := t.TempDir()
	baseDir := filepath.Join(tempRoot, "docs", "videorea")
	slug := "legoloop-video-auto-cut"

	report := &videorea.CurationReport{
		SuggestedName: slug,
		OriginalCost:  "R$ 250,00/mes",
		SovereignCost: "R$ 0,00/mes",
		Replacements: map[string]string{
			"Premiere": "FFmpeg + Remotion",
		},
		Approved: true,
	}

	bundle, err := videorea.MaterializeAndProject(report, baseDir, tempRoot)
	if err != nil {
		t.Fatalf("MaterializeAndProject falhou: %v", err)
	}

	// 1. Verificar diretórios internos do bundle
	dirsToCheck := []string{
		bundle.ScratchDir,
		bundle.GovernanceDir,
		bundle.BookDir,
		bundle.LoopDir,
		bundle.SkillDir,
	}
	for _, d := range dirsToCheck {
		if fi, err := os.Stat(d); os.IsNotExist(err) || !fi.IsDir() {
			t.Errorf("diretório do bundle não foi criado: %s", d)
		}
	}

	// 2. Verificar arquivos essenciais no bundle
	filesToCheck := []string{
		filepath.Join(bundle.BookDir, "livro.json"),
		filepath.Join(bundle.BookDir, "partes", "00-frontmatter.md"),
		filepath.Join(bundle.BookDir, "partes", "01-fundacao.md"),
		filepath.Join(bundle.LoopDir, "pipeline.dag.json"),
		filepath.Join(bundle.LoopDir, "fixture_3s.mp4"),
		filepath.Join(bundle.SkillDir, "SKILL.md"),
		bundle.OriginalMatrixPath,
	}
	for _, f := range filesToCheck {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("arquivo essencial do bundle não existe: %s", f)
		}
	}

	// 3. Verificar Projeção Canônica na raiz do LegoLoop
	canonicalDAG := filepath.Join(tempRoot, "loops", slug, "pipeline.dag.json")
	if _, err := os.Stat(canonicalDAG); os.IsNotExist(err) {
		t.Errorf("projeção canônica do DAG não encontrada: %s", canonicalDAG)
	}

	canonicalSkill := filepath.Join(tempRoot, "skills", slug, "SKILL.md")
	if _, err := os.Stat(canonicalSkill); os.IsNotExist(err) {
		t.Errorf("projeção canônica da Skill não encontrada: %s", canonicalSkill)
	}

	canonicalBook := filepath.Join(tempRoot, "docs", "books", slug, "livro.json")
	if _, err := os.Stat(canonicalBook); os.IsNotExist(err) {
		t.Errorf("projeção canônica do Livro não encontrada: %s", canonicalBook)
	}
}

func TestMaterializeAndProject_GravaMatrizOriginalCompleta(t *testing.T) {
	tempRoot := t.TempDir()
	baseDir := filepath.Join(tempRoot, "docs", "videorea")
	slug := "legoloop-minimax-h3-local"

	report := &videorea.CurationReport{
		SuggestedName: slug,
		VideoTitle:    "Adeus IAs pagas: Nova IA gera videos em PC fraco",
		OriginalCost:  "R$ 0,00/mes",
		SovereignCost: "R$ 0,00/mes",
		OriginalMatrix: &videorea.OriginalMatrix{
			VideoTitle:       "Adeus IAs pagas: Nova IA gera videos em PC fraco",
			TotalMonthlyCost: "R$ 0,00/mes",
			Tools: []videorea.OriginalTool{
				{Name: "ComfyUI", IsFree: true, CostTier: "R$ 0,00 (Gratuito/Aberto)"},
				{Name: "MiniMax H3", IsFree: true, CostTier: "R$ 0,00 (Gratuito/Aberto)"},
				{Name: "FLUX", IsFree: true, CostTier: "R$ 0,00 (Gratuito/Aberto)"},
			},
			OriginalSteps: []string{"Execucao local via ComfyUI Portable"},
		},
		Approved: true,
	}

	bundle, err := videorea.MaterializeAndProject(report, baseDir, tempRoot)
	if err != nil {
		t.Fatalf("MaterializeAndProject falhou: %v", err)
	}

	data, err := os.ReadFile(bundle.OriginalMatrixPath)
	if err != nil {
		t.Fatalf("lendo original_matrix.json: %v", err)
	}
	var matrix videorea.OriginalMatrix
	if err := json.Unmarshal(data, &matrix); err != nil {
		t.Fatalf("parse original_matrix.json: %v", err)
	}
	if len(matrix.Tools) != 3 {
		t.Errorf("esperava 3 ferramentas gravadas, obteve %d (%s)", len(matrix.Tools), data)
	}
	if matrix.VideoTitle != report.VideoTitle {
		t.Errorf("video_title do arquivo esperado %q, obteve %q", report.VideoTitle, matrix.VideoTitle)
	}
}
