package roteamento

import (
	"errors"
	"fmt"
)

const (
	ClasseNano    = "nano"
	ClasseMicro   = "micro"
	ClasseGrande  = "macro"
	ClasseExtreme = "extreme"
)

const (
	maxArquivosNode = 12
	maxTokensNode   = 32000
	tokensGrande    = 16000
	tokensMicro     = 8000
)

var ErrNoGrande = errors.New("no acima do limite de tamanho: exige divisao")

type EntradaClasse struct {
	ArquivosSaida  int
	TokensEntrada  int
	Pacotes        int
	Destrutivo     bool
	Sobrescrita    string
	JustificativaS string
}

type ResultadoClasse struct {
	Classe              string
	ClassRule           string
	ClassOverrideReason string
}

func regra(classe, texto string) (ResultadoClasse, error) {
	return ResultadoClasse{Classe: classe, ClassRule: classe + ": " + texto}, nil
}

func Classificar(e EntradaClasse) (ResultadoClasse, error) {
	if e.ArquivosSaida > maxArquivosNode || e.TokensEntrada > maxTokensNode {
		return ResultadoClasse{}, ErrNoGrande
	}
	if e.Sobrescrita != "" {
		return classeSobrescrita(e)
	}
	switch {
	case e.Destrutivo:
		return regra(ClasseExtreme, "nature destrutiva")
	case e.ArquivosSaida > 4:
		return regra(ClasseGrande, "mais de 4 arquivos em context_out")
	case e.TokensEntrada > tokensGrande:
		return regra(ClasseGrande, "entrada acima de 16000 tokens")
	case e.Pacotes > 1:
		return regra(ClasseGrande, "toca mais de 1 pacote")
	case e.ArquivosSaida >= 2:
		return regra(ClasseMicro, "de 2 a 4 arquivos em context_out")
	case e.TokensEntrada >= tokensMicro:
		return regra(ClasseMicro, "entrada de 8000 a 16000 tokens")
	}
	return regra(ClasseNano, "no maximo 1 arquivo e entrada abaixo de 8000 tokens")
}

func classeSobrescrita(e EntradaClasse) (ResultadoClasse, error) {
	switch e.Sobrescrita {
	case ClasseNano, ClasseMicro, ClasseGrande, ClasseExtreme:
	default:
		return ResultadoClasse{}, fmt.Errorf("classe de sobrescrita invalida: %q", e.Sobrescrita)
	}
	if e.JustificativaS == "" {
		return ResultadoClasse{}, errors.New("sobrescrita exige justificativa")
	}
	return ResultadoClasse{
		Classe:              e.Sobrescrita,
		ClassRule:           e.Sobrescrita + ": sobrescrita explicita",
		ClassOverrideReason: e.JustificativaS,
	}, nil
}
