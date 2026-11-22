package fabricar

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/matriz"
	"github.com/heverton-dev/LegoLoop/internal/processo"
)

// comContratos completa cada peça com o recorte de I/O do manifesto contrato.json ao lado do
// SKILL.md. Peça sem manifesto, ou com manifesto quebrado, segue sem contrato: o comportamento
// atual da aptidão está preservado (AC-E05.1).
func comContratos(es []matriz.Entrada) []matriz.Entrada {
	for i := range es {
		if es[i].ContratoIO != nil {
			continue
		}
		c, err := processo.ContratoDe(contratoAoLado(es[i].Ref))
		if err != nil {
			continue
		}
		es[i].ContratoIO = &matriz.IOContrato{Sai: tiposDeSai(c)}
	}
	return es
}

// contratoAoLado devolve o caminho do manifesto contrato.json ao lado do SKILL.md da peça,
// expandindo o atalho de home quando a Ref da matriz o usa.
func contratoAoLado(ref string) string {
	caminho := ref
	if strings.HasPrefix(caminho, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			caminho = filepath.Join(home, caminho[2:])
		}
	}
	return filepath.Join(filepath.Dir(caminho), "contrato.json")
}

// tiposDeSai lista os tipos de artefato que o contrato da peça promete produzir.
func tiposDeSai(c processo.Contrato) []string {
	out := make([]string, 0, len(c.Sai))
	for _, s := range c.Sai {
		out = append(out, s.Tipo)
	}
	return out
}
