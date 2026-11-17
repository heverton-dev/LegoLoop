package fabricar

import (
	"encoding/json"

	"legoloop/internal/processo"
)

// vacinaAnticorpo é uma entrada de anticorpo: assinatura do erro, o erro, o antídoto e as reincidências.
type vacinaAnticorpo struct {
	Assinatura   string `json:"assinatura"`
	Erro         string `json:"erro"`
	Antidoto     string `json:"antidoto"`
	Recorrencias int    `json:"recorrencias"`
}

// anticorposV1 é o manifesto anticorpos.v1 do Processo: vacinas globais (infra/LLM) e de domínio.
type anticorposV1 struct {
	SchemaVersion string            `json:"schema_version"`
	Processo      string            `json:"processo"`
	Globais       []vacinaAnticorpo `json:"globais"`
	Dominio       []vacinaAnticorpo `json:"dominio"`
}

// AnticorposPadrao renderiza o anticorpos.json do nascimento 5/5: sem vacinas registradas,
// pronto para receber as de infra/LLM (globais) e as do domínio do Processo.
func AnticorposPadrao(nome string) ([]byte, error) {
	return json.MarshalIndent(anticorposV1{
		SchemaVersion: "anticorpos.v1",
		Processo:      nome,
		Globais:       []vacinaAnticorpo{},
		Dominio:       []vacinaAnticorpo{},
	}, "", "  ")
}

// anticorposDe renderiza o anticorpos.json ligado ao Processo recém-projetado.
func anticorposDe(p processo.Processo) ([]byte, error) {
	return AnticorposPadrao(p.Nome)
}
