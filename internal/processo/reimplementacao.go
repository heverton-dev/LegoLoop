package processo

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// ErrReimplementacao: a peça nova reimplementa o I/O de peça existente do cofre (AC-E04.1).
var ErrReimplementacao = errors.New("processo: reimplementação de peça existente")

// assinaturaIO devolve a assinatura estrutural do I/O: multiconjunto ordenado de tipo e
// obrigação de cada entrada e saída. Compatibilidade estrutural, nunca hash: caminho padrão
// e gate sugerido são detalhe de implementação e não contam.
func assinaturaIO(c Contrato) string {
	partes := make([]string, 0, len(c.Entra)+len(c.Sai))
	for _, e := range c.Entra {
		partes = append(partes, fmt.Sprintf("entra:%s:%t", e.Tipo, e.Opcional))
	}
	for _, s := range c.Sai {
		partes = append(partes, fmt.Sprintf("sai:%s:%t", s.Tipo, s.Promete))
	}
	sort.Strings(partes)
	return strings.Join(partes, "|")
}

// ContratosDoCofre lê o manifesto contrato.json de cada peça do cofre (ao lado do SKILL.md).
func ContratosDoCofre(vault string) ([]Contrato, error) {
	achados, err := filepath.Glob(filepath.Join(vault, "skills", "*", "contrato.json"))
	if err != nil {
		return nil, err
	}
	out := make([]Contrato, 0, len(achados))
	for _, caminho := range achados {
		c, err := ContratoDe(caminho)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

// VerificarReimplementacao reprova a peça nova cujo I/O é estruturalmente igual ao de peça
// existente com outro nome (AC-E04.1), citando a peça existente. Peça de I/O distinto passa
// (AC-E04.2). Atualização da própria peça passa. Nunca grava nada no disco.
func VerificarReimplementacao(vault string, nova Contrato) error {
	existentes, err := ContratosDoCofre(vault)
	if err != nil {
		return err
	}
	alvo := assinaturaIO(nova)
	for _, c := range existentes {
		if c.Nome == nova.Nome {
			continue
		}
		if assinaturaIO(c) == alvo {
			return fmt.Errorf("%w: %s reimplementa o I/O de %s", ErrReimplementacao, nova.Nome, c.Nome)
		}
	}
	return nil
}
