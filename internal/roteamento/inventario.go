package roteamento

import (
	"runtime"

	"legoloop/internal/processo"
)

type Candidato struct{ Harness, Modelo, Esforco string }

type Perfil map[string][]Candidato

type HarnessInv struct {
	Nome    string
	Binario string
	Modelos []string
}

type Inventario struct {
	SO        string
	Harnesses []HarnessInv
}

type Substituicao struct {
	No, Classe    string
	Pedido, Usado Candidato
	Motivo        string
}

func (inv Inventario) Vazio() bool { return len(inv.Harnesses) == 0 }

func LerInventario(fsys processo.FS, home string, binario func(string) string, modelos map[string][]string) Inventario {
	inv := Inventario{SO: runtime.GOOS}
	for _, h := range processo.DetectarHarnesses(fsys, home) {
		item := HarnessInv{Nome: h.Nome, Modelos: modelos[h.Nome]}
		if binario != nil {
			item.Binario = binario(h.Nome)
		}
		inv.Harnesses = append(inv.Harnesses, item)
	}
	return inv
}
