package roteamento

import (
	"encoding/json"
	"fmt"
)

type CandidatoVia struct {
	Harness string `json:"harness"`
	Via     Via    `json:"via"`
	Modelo  string `json:"modelo"`
	Esforco string `json:"esforco"`
}

type GatewayPerfil struct {
	Habilitado *bool `json:"habilitado,omitempty"`
}

type PerfilArquivo struct {
	SchemaVersion int                       `json:"schema_version"`
	Gateways      map[string]GatewayPerfil  `json:"gateways,omitempty"`
	Candidatos    map[string][]CandidatoVia `json:"candidatos"`
}

func PerfilPadrao() PerfilArquivo {
	oficial := func(modelo, esforco string) []CandidatoVia {
		return []CandidatoVia{{Harness: "claude", Via: ViaOficial, Modelo: modelo, Esforco: esforco}}
	}
	return PerfilArquivo{
		SchemaVersion: 1,
		Candidatos: map[string][]CandidatoVia{
			ClasseNano:    oficial("haiku", "low"),
			ClasseMicro:   oficial("sonnet", "low"),
			ClasseGrande:  oficial("sonnet", "high"),
			ClasseExtreme: oficial("opus", "high"),
		},
	}
}

func lerPerfil(dados []byte) (PerfilArquivo, error) {
	var p PerfilArquivo
	if err := json.Unmarshal(dados, &p); err != nil {
		return PerfilArquivo{}, fmt.Errorf("perfil invalido: %w", err)
	}
	return p, nil
}

func sobrepor(base, forte PerfilArquivo) PerfilArquivo {
	for classe, lista := range forte.Candidatos {
		base.Candidatos[classe] = lista
	}
	for nome, gw := range forte.Gateways {
		if base.Gateways == nil {
			base.Gateways = map[string]GatewayPerfil{}
		}
		if gw.Habilitado != nil {
			base.Gateways[nome] = gw
		}
	}
	return base
}

func (p PerfilArquivo) viaHabilitada(via Via) bool {
	if via != Via9Router {
		return true
	}
	gw, ok := p.Gateways[string(Via9Router)]
	return !ok || gw.Habilitado == nil || *gw.Habilitado
}

func (p PerfilArquivo) Elegiveis(store map[string]Homologacao) map[string][]CandidatoVia {
	out := map[string][]CandidatoVia{}
	for classe, lista := range p.Candidatos {
		var ok []CandidatoVia
		for _, c := range lista {
			par := Par{Harness: c.Harness, Via: c.Via}
			h, achou := store[par.Chave()]
			if p.viaHabilitada(c.Via) && achou && h.Estado == EstadoHomologado {
				ok = append(ok, c)
			}
		}
		out[classe] = ok
	}
	return out
}

func CarregarPerfil(global, projeto []byte, store map[string]Homologacao) (map[string][]CandidatoVia, error) {
	perfil := PerfilPadrao()
	for _, dados := range [][]byte{global, projeto} {
		if len(dados) == 0 {
			continue
		}
		forte, err := lerPerfil(dados)
		if err != nil {
			return nil, err
		}
		perfil = sobrepor(perfil, forte)
	}
	return perfil.Elegiveis(store), nil
}
