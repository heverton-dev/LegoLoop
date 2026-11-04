package cli

import (
	"fmt"
	"io"
	"strings"

	"legoloop/internal/fabricar"
	"legoloop/internal/ste"
)

// aprovar mostra o cartão do plano único e conduz as três decisões. Enter vazio cancela. "Personalizar" abre os
// parâmetros e volta ao cartão. Depois de "seguir" vêm as pausas extras de cada ação de risco.
func (p *prompter) aprovar(plano fabricar.Plano, replan func([]string) fabricar.Plano) (bool, error) {
	invalidas := 0
	for invalidas < maxInvalid {
		card := plano.Cartao()
		fmt.Fprintf(p.w, "\n%s", card.Render())
		in, err := p.line("> ")
		if err != nil {
			return false, err
		}
		dec, ok := ste.ParseDecisao(in)
		switch {
		case !ok:
			invalidas++
			fmt.Fprintln(p.w, "Entrada inválida: responda 1, 2 ou 0.")
		case dec == ste.Seguir:
			return p.pausasExtras(plano)
		case dec == ste.Personalizar:
			if plano, err = p.personalizar(card, plano, replan); err != nil {
				return false, err
			}
		default:
			return false, nil
		}
	}
	return false, errTooManyInvalid
}

// personalizar pergunta o novo valor de cada parâmetro (Enter mantém) e devolve o plano recalculado.
func (p *prompter) personalizar(card ste.Cartao, plano fabricar.Plano, replan func([]string) fabricar.Plano) (fabricar.Plano, error) {
	for _, par := range card.Parametros {
		fmt.Fprintf(p.w, "%s\n", par.Ajuda)
		in, err := p.line(fmt.Sprintf("%s [%s]: ", par.Nome, par.Valor))
		if err != nil {
			return plano, err
		}
		if in != "" && par.Nome == "harnesses" {
			plano = replan(parseNomes(in))
		}
	}
	return plano, nil
}

// parseNomes separa nomes por vírgula ou espaço. "nenhum" vale lista vazia.
func parseNomes(in string) []string {
	var out []string
	for _, f := range strings.FieldsFunc(in, func(r rune) bool { return r == ',' || r == ' ' }) {
		if strings.EqualFold(f, "nenhum") {
			return nil
		}
		out = append(out, f)
	}
	return out
}

// pausasExtras pede uma confirmação própria para cada ação de risco do plano. Qualquer "não" cancela tudo.
func (p *prompter) pausasExtras(plano fabricar.Plano) (bool, error) {
	pausas := []struct {
		ativa bool
		texto string
	}{
		{len(plano.Sobrescreve) > 0, fmt.Sprintf("Trocar %d arquivos que já existem? [s/N]: ", len(plano.Sobrescreve))},
		{len(plano.Vinculos) > 0, fmt.Sprintf("Ligar a skill a %d harnesses? [s/N]: ", len(plano.Vinculos))},
		{len(plano.Removivel) > 0, fmt.Sprintf("Usar %d peças de mídia removível? [s/N]: ", len(plano.Removivel))},
	}
	for _, x := range pausas {
		if !x.ativa {
			continue
		}
		if ok, err := p.yesNo(x.texto); err != nil || !ok {
			return false, err
		}
	}
	return true, nil
}

// aprovarSemPerguntar é o modo --yes: imprime o cartão e segue só no que ele mostrou.
func aprovarSemPerguntar(w io.Writer) fabricar.ApproveFunc {
	return func(plano fabricar.Plano, _ func([]string) fabricar.Plano) (bool, error) {
		fmt.Fprintf(w, "\n%sSegue sem perguntar, por causa do --yes.\n", plano.Cartao().Texto())
		return true, nil
	}
}
