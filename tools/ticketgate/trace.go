package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/artefato"
)

func expandTrace(root string, alvos []string) ([]string, error) {
	var arquivos []string
	for _, a := range alvos {
		base := filepath.Join(root, a)
		info, err := os.Stat(base)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			arquivos = append(arquivos, base)
			continue
		}
		err = filepath.WalkDir(base, func(p string, d fs.DirEntry, err error) error {
			ext := strings.ToLower(filepath.Ext(p))
			if err == nil && !d.IsDir() && (ext == ".md" || ext == ".json") {
				arquivos = append(arquivos, p)
			}
			return err
		})
		if err != nil {
			return nil, err
		}
	}
	return arquivos, nil
}

func traceStep(c config) step {
	return step{Name: "trace " + strings.Join(c.Trace, " "), Run: func() (string, error) {
		arquivos, err := expandTrace(c.Root, c.Trace)
		if err != nil {
			return "", err
		}
		cadeia, err := artefato.CarregarCadeia(arquivos)
		if err != nil {
			return "", err
		}
		if orfaos := cadeia.Orfaos(); len(orfaos) > 0 {
			return "órfãos: " + strings.Join(orfaos, ", "), errors.New("cadeia de rastreabilidade quebrada")
		}
		return "", nil
	}}
}
