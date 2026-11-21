//go:build darwin

package descoberta

import (
	"os"
	"path/filepath"
)

type provedorDarwin struct{}

// ProvedorPadrao devolve o provedor real somente leitura para macOS.
func ProvedorPadrao() Provedor {
	return provedorDarwin{}
}

func (provedorDarwin) SO() string { return "darwin" }

func (provedorDarwin) ListarUnidades() []Unidade {
	entradas, err := os.ReadDir("/Volumes")
	if err != nil {
		return nil
	}
	var out []Unidade
	for _, e := range entradas {
		caminho := filepath.Join("/Volumes", e.Name())
		out = append(out, Unidade{
			Caminho: caminho,
			Tipo:    classificarVolumeDarwin(caminho),
		})
	}
	return out
}

func classificarVolumeDarwin(caminho string) TipoUnidade {
	destino, err := os.Readlink(caminho)
	if err == nil && destino == "/" {
		return TipoFixa
	}
	return TipoRemovivel
}
