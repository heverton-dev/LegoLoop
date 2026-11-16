package processo

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"legoloop/internal/triagem"
)

// TriagemDoProcesso devolve a triagem que o processo exige, ou nil se não exige. Com mais de um nó a triagem
// precisa existir e cobrir todas as etapas: a ausente nasce, a existente ganha só as regras que faltam.
// Só vale no layout do fabricar, em que o processo tem pasta própria no cofre (a pasta do DAG leva o nome dele).
func TriagemDoProcesso(dagPath string, p Processo) (*File, error) {
	pasta := filepath.Dir(dagPath)
	if !triagem.Exigida(len(p.Etapas)) || filepath.Base(pasta) != p.Nome {
		return nil, nil
	}
	ids := make([]string, len(p.Etapas))
	for i, e := range p.Etapas {
		ids[i] = e.ID
	}
	atual, err := lerTriagemExistente(pasta)
	if err != nil {
		return nil, err
	}
	dados, err := json.MarshalIndent(triagem.Completar(atual, p.Nome, ids), "", "  ")
	if err != nil {
		return nil, err
	}
	return &File{Kind: "triagem", Path: filepath.Join(pasta, triagem.ArquivoTriagem), Data: append(dados, '\n')}, nil
}

func lerTriagemExistente(pasta string) (triagem.TriagemSpec, error) {
	var s triagem.TriagemSpec
	b, err := os.ReadFile(filepath.Join(pasta, triagem.ArquivoTriagem))
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, err
	}
	return s, json.Unmarshal(b, &s)
}
