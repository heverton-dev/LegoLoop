package remover

import (
	"strings"

	"legoloop/internal/ste"
)

func init() {
	for id, msg := range map[string]string{
		"remover.cartao.acao":    "Remover o processo %s do cofre.",
		"remover.cartao.aviso":   "Esta ação não tem volta. Os arquivos apagados não vão para a lixeira.",
		"remover.cartao.passo_a": "Apagar %d itens do cofre: a skill, o loop, o DAG e o processo.",
		"remover.cartao.passo_b": "Apagar %d vínculos nos harnesses: %s.",
		"remover.cartao.efeito":  "Os documentos do projeto em docs/ ficam como estão.",
		"remover.cartao.desfaz":  "Se algo falhar, o processo continua no cofre. Rode legoloop remover de novo.",
	} {
		ste.Default.Register(id, msg)
	}
}

// cartao explica o efeito antes de apagar. O aviso de ação irreversível vem antes dos passos.
func cartao(p Plano) ste.Cartao {
	c := ste.Cartao{
		Acao:         ste.Msg("remover.cartao.acao", p.Nome),
		Aviso:        ste.Msg("remover.cartao.aviso"),
		Passos:       []string{ste.Msg("remover.cartao.passo_a", len(p.Arquivos))},
		Efeitos:      []string{ste.Msg("remover.cartao.efeito")},
		SeDerErrado:  ste.Msg("remover.cartao.desfaz"),
		Irreversivel: true,
	}
	if len(p.Links) > 0 {
		c.Passos = append(c.Passos, ste.Msg("remover.cartao.passo_b", len(p.Links), strings.Join(p.Harnesses, ", ")))
	}
	return c
}
