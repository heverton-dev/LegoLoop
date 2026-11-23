package triagem

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// CarregarSpecs lê e valida todos os arquivos triagem.json sob processosDir.
func CarregarSpecs(processosDir string) ([]TriagemSpec, error) {
	entries, err := os.ReadDir(processosDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var catalogo []TriagemSpec
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		caminho := filepath.Join(processosDir, e.Name(), "triagem.json")
		spec, ok, err := LerSpec(caminho)
		if err != nil {
			return nil, err
		}
		if ok {
			catalogo = append(catalogo, spec)
		}
	}
	return catalogo, nil
}

// LerSpec lê e valida um único arquivo triagem.json.
func LerSpec(caminho string) (TriagemSpec, bool, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		if os.IsNotExist(err) {
			return TriagemSpec{}, false, nil
		}
		return TriagemSpec{}, false, err
	}
	var spec TriagemSpec
	if err := json.Unmarshal(dados, &spec); err != nil {
		return TriagemSpec{}, false, err
	}
	if err := spec.Validar(); err != nil {
		return TriagemSpec{}, false, err
	}
	return spec, true, nil
}
