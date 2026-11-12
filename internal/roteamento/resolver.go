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

type SubstituicaoVia struct {
	No, Classe    string
	Pedido, Usado CandidatoVia
	Motivo        string
}

type ResolucaoVia struct {
	No            string
	Classe        string
	Usado         CandidatoVia
	Politica      string
	StopGate      bool
	Substituicoes []SubstituicaoVia
}

func disponivelVia(c CandidatoVia, inv []ItemInv) bool {
	par := Par{Harness: c.Harness, Via: c.Via}
	for _, item := range inv {
		if item.Par != par || item.Estado != EstadoHomologado {
			continue
		}
		for _, m := range item.Modelos {
			if m == c.Modelo {
				return true
			}
		}
	}
	return false
}

func ResolverVia(no, classe string, elegiveis map[string][]CandidatoVia, inv []ItemInv) (ResolucaoVia, error) {
	if len(inv) == 0 {
		return ResolucaoVia{}, ErrInventarioVazio
	}
	res := ResolucaoVia{No: no, Classe: classe}
	candidatos := elegiveis[classe]
	for i, c := range candidatos {
		if !disponivelVia(c, inv) {
			continue
		}
		res.Usado = c
		if i > 0 {
			res.Substituicoes = append(res.Substituicoes, SubstituicaoVia{
				No:     no,
				Classe: classe,
				Pedido: candidatos[0],
				Usado:  c,
				Motivo: "primeiro candidato ausente do inventario",
			})
		}
		return res, nil
	}
	res.Politica = PoliticaDeFalha(classe)
	res.StopGate = classe == ClasseExtreme
	return res, nil
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
