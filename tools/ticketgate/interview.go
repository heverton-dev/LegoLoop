package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/artefato"
)

func interviewStep(root string, alvos []string) step {
	return step{Name: "interview " + strings.Join(alvos, " "), Run: func() (string, error) {
		arquivos, err := expandTrace(root, alvos)
		if err != nil {
			return "", err
		}
		var falhas []string
		for _, a := range arquivos {
			if strings.ToLower(filepath.Ext(a)) != ".md" {
				continue
			}
			dados, err := os.ReadFile(a)
			if err != nil {
				return "", err
			}
			for _, f := range artefato.ValidarEntrevista(string(dados)) {
				falhas = append(falhas, filepath.Base(a)+": "+f)
			}
		}
		if len(falhas) > 0 {
			return strings.Join(falhas, "\n"), errors.New("entrevista sem prova")
		}
		return "", nil
	}}
}
