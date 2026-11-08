package videorea

import (
	"os"
	"path/filepath"
)

// TargetBundle encapsula o espaço de trabalho isolado de um vídeo.
type TargetBundle struct {
	Slug               string `json:"slug"`
	BaseDir            string `json:"base_dir"`
	ScratchDir         string `json:"scratch_dir"`
	DossierPath        string `json:"dossier_path"`
	OriginalMatrixPath string `json:"original_matrix_path"`
	CurationPath       string `json:"curation_path"`
	GovernanceDir      string `json:"governance_dir"`
	BookDir            string `json:"book_dir"`
	LoopDir            string `json:"loop_dir"`
	SkillDir           string `json:"skill_dir"`
}

// NewTargetBundle calcula a hierarquia completa de diretórios para o alvo.
func NewTargetBundle(baseDir, slug string) *TargetBundle {
	base := filepath.Join(baseDir, slug)
	return &TargetBundle{
		Slug:               slug,
		BaseDir:            base,
		ScratchDir:         filepath.Join(base, "scratch"),
		DossierPath:        filepath.Join(base, "dossier.json"),
		OriginalMatrixPath: filepath.Join(base, "original_matrix.json"),
		CurationPath:       filepath.Join(base, "curation.json"),
		GovernanceDir:      filepath.Join(base, "governance"),
		BookDir:            filepath.Join(base, "book"),
		LoopDir:            filepath.Join(base, "loop"),
		SkillDir:           filepath.Join(base, "skill"),
	}
}

// Materialize cria a árvore física de pastas no sistema de arquivos.
func (b *TargetBundle) Materialize() error {
	dirs := []string{
		b.BaseDir,
		b.ScratchDir,
		b.GovernanceDir,
		b.BookDir,
		b.LoopDir,
		b.SkillDir,
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}
