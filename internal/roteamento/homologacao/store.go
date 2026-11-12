package homologacao

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const EstadoPendente = "pendente"

type Entrada struct {
	Estado    string          `json:"estado"`
	Data      string          `json:"data"`
	PesoMB    int             `json:"peso_mb"`
	Criterios map[string]bool `json:"criterios"`
	Motivo    string          `json:"motivo"`
}

type Store map[string]Entrada

func Chave(harness, via string) string { return harness + "|" + via }

func Caminho(home string) string {
	return filepath.Join(home, ".aidd", "homologacao.json")
}

func Ler(caminho string) (Store, error) {
	dados, err := os.ReadFile(caminho)
	if errors.Is(err, os.ErrNotExist) {
		return Store{}, nil
	}
	if err != nil {
		return nil, err
	}
	s := Store{}
	if err := json.Unmarshal(dados, &s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s Store) Gravar(caminho string) error {
	dados, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(caminho)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "homologacao-*.tmp")
	if err != nil {
		return err
	}
	_, err = tmp.Write(append(dados, '\n'))
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), caminho); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}
