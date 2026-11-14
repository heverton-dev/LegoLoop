package fabricar

import (
	"fmt"
	"legoloop/internal/processo"
	"slices"
)

// resolverCascataLacunas aplica a cascata de 3 fontes para etapas sem peça na Matriz:
// 1) Matriz do Cofre/Processo (já resolvida quando Peca ou Compoe != "").
// 2) Sistema Operacional (ferramentas/skills locais, se disponíveis).
// 3) GitHub/Fontes externas: não bloqueia; registra sugestão em Warnings e CuratedPieces.
func resolverCascataLacunas(p *processo.Processo, c *consulta) {
	for i := range p.Etapas {
		e := &p.Etapas[i]
		if e.Peca != "" || e.Compoe != "" {
			continue
		}
		nome := nomeEtapaLacuna(*e)
		sugestao := sugerirFonteExterna(nome)
		e.Peca = sugestao
		aviso := fmt.Sprintf("lacuna na etapa %q: peça ausente no cofre e no SO; sugestão externa (GitHub/Fontes): %s", nome, sugestao)
		c.m.Warnings = append(c.m.Warnings, aviso)
		if !slices.ContainsFunc(c.m.Used, func(dp DiscoveredPiece) bool { return dp.Name == sugestao }) {
			c.m.Used = append(c.m.Used, DiscoveredPiece{Name: sugestao, Kind: PieceKind("externa")})
		}
	}
}

func nomeEtapaLacuna(e processo.Etapa) string {
	if e.Tipo != "" {
		return e.Tipo
	}
	return e.Titulo
}

func sugerirFonteExterna(etapa string) string {
	return fmt.Sprintf("curadoria-externa-%s", etapa)
}
