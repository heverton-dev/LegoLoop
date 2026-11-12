package homologacao

import (
	"io/fs"
	"os"
	"path/filepath"
)

const (
	EstadoHomologado = "homologado"
	EstadoReprovado  = "reprovado"
)

func ForaDoLimite(dir string, permitidos []string) []string {
	ok := map[string]bool{}
	for _, p := range permitidos {
		ok[filepath.ToSlash(filepath.Clean(p))] = true
	}
	var fora []string
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if !ok[rel] {
			fora = append(fora, rel)
		}
		return nil
	})
	return fora
}

var reverterArquivo = os.Remove

func Reverter(dir string, lista []string) error {
	var primeiro error
	for _, rel := range lista {
		alvo := filepath.Join(dir, filepath.FromSlash(rel))
		if err := reverterArquivo(alvo); err != nil && primeiro == nil {
			primeiro = err
		}
		if pai := filepath.Dir(alvo); pai != filepath.Clean(dir) {
			_ = reverterArquivo(pai)
		}
	}
	return primeiro
}
