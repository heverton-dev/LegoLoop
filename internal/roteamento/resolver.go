package roteamento

import "errors"

var ErrInventarioVazio = errors.New("inventario vazio")

type Resolucao struct {
	No            string
	Classe        string
	Usado         Candidato
	StopGate      bool
	Substituicoes []Substituicao
}

func (inv Inventario) tem(c Candidato) bool {
	for _, h := range inv.Harnesses {
		if h.Nome != c.Harness {
			continue
		}
		for _, m := range h.Modelos {
			if m == c.Modelo {
				return true
			}
		}
	}
	return false
}

func Resolver(no, classe string, perfil Perfil, inv Inventario) (Resolucao, error) {
	if len(inv.Harnesses) == 0 {
		return Resolucao{}, ErrInventarioVazio
	}
	res := Resolucao{No: no, Classe: classe}
	candidatos := perfil[classe]
	for i, c := range candidatos {
		if !inv.tem(c) {
			continue
		}
		res.Usado = c
		if i > 0 {
			res.Substituicoes = append(res.Substituicoes, Substituicao{
				No:     no,
				Classe: classe,
				Pedido: candidatos[0],
				Usado:  c,
				Motivo: "primeiro candidato ausente do inventario",
			})
		}
		return res, nil
	}
	if classe == "extreme" {
		res.StopGate = true
	}
	return res, nil
}
