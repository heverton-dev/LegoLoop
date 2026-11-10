package videorea

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

// ErrBundleEscape sinaliza slug que escapa da raiz do bundle.
var ErrBundleEscape = errors.New("videorea: slug escapa da raiz do bundle")

// SanitizeBundlePath confina o slug dentro de baseDir (Path Jail).
func SanitizeBundlePath(baseDir, slug string) (string, error) {
	if slug == "" || slug == "." || strings.ContainsAny(slug, `/\`) || strings.Contains(slug, "..") || filepath.IsAbs(slug) || filepath.VolumeName(slug) != "" {
		return "", fmt.Errorf("%w: %q", ErrBundleEscape, slug)
	}
	root := filepath.Clean(baseDir)
	target := filepath.Clean(filepath.Join(root, slug))
	if !strings.HasPrefix(target, root+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: %q", ErrBundleEscape, slug)
	}
	return target, nil
}

// NewTargetBundle calcula a hierarquia completa de diretórios para o alvo, confinada em baseDir.
func NewTargetBundle(baseDir, slug string) (*TargetBundle, error) {
	base, err := SanitizeBundlePath(baseDir, slug)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

// Materialize cria a árvore física de pastas no sistema de arquivos.
// Percorre todo campo string cujo nome termina em "Dir" (BaseDir, ScratchDir etc.).
func (b *TargetBundle) Materialize() error {
	v := reflect.ValueOf(b).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() != reflect.String || !strings.HasSuffix(f.Name, "Dir") {
			continue
		}
		if err := os.MkdirAll(v.Field(i).String(), 0o755); err != nil {
			return err
		}
	}
	return nil
}
