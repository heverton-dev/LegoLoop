//go:build linux

package descoberta

import (
	"os"
	"path/filepath"
)

type provedorLinux struct{}

// ProvedorPadrao devolve o provedor real somente leitura para Linux.
func ProvedorPadrao() Provedor {
	return provedorLinux{}
}

func (provedorLinux) SO() string { return "linux" }

func (provedorLinux) ListarUnidades() []Unidade {
	var out []Unidade
	out = append(out, listarMontagens("/mnt", TipoFixa)...)
	out = append(out, listarMontagens("/media", TipoRemovivel)...)
	return out
}

func listarMontagens(base string, tipo TipoUnidade) []Unidade {
	entradas, err := os.ReadDir(base)
	if err != nil {
		return nil
	}
	var out []Unidade
	for _, e := range entradas {
		if !e.IsDir() {
			continue
		}
		out = append(out, Unidade{
			Caminho: filepath.Join(base, e.Name()),
			Tipo:    tipo,
		})
	}
	return out
}
