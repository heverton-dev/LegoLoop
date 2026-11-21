package interview

import (
	"regexp"
	"strings"
)

var (
	// reRotuloEspelho abre um campo do espelho. Aceita o realce antes ou depois dos dois-pontos
	// (**Campo**: x e **Campo:** x), que é como os cérebros o emitem.
	reRotuloEspelho = regexp.MustCompile(`(?i)^\s*[*-]\s*\*{0,2}(Resultado Esperado|Quem se Beneficia|Por que Agora|Critério de Sucesso|Restrições Inegociáveis|Fora de Escopo)\*{0,2}:\*{0,2}\s*(.*)$`)
	// reFimEspelho encerra o último campo: a confiança ou a pergunta repetidas depois do espelho.
	reFimEspelho = regexp.MustCompile(`(?i)^\s*\*{0,2}(CONFIANÇA|PERGUNTA|PALPITE|HIPÓTESE|PREMONIÇÃO)`)
	reMarcador   = regexp.MustCompile(`^\s*[*-]\s+`)
)

// limpaRealce tira o realce Markdown (no meio e nas pontas) e os espaços das pontas.
func limpaRealce(s string) string {
	return strings.TrimSpace(strings.Trim(strings.ReplaceAll(s, "**", ""), "* \t"))
}

// camposEspelho agrupa as linhas de cada campo do espelho: a linha do rótulo e as seguintes,
// até o próximo rótulo ou o fim do espelho. A chave é o rótulo em minúsculas.
func camposEspelho(raw string) map[string][]string {
	campos := map[string][]string{}
	atual := ""
	for _, linha := range strings.Split(raw, "\n") {
		if m := reRotuloEspelho.FindStringSubmatch(linha); m != nil {
			atual = strings.ToLower(m[1])
			campos[atual] = nil
			linha = m[2]
		} else if atual == "" || reFimEspelho.MatchString(linha) {
			atual = ""
			continue
		} else if strings.TrimSpace(linha) == "" {
			if len(campos[atual]) > 0 {
				atual = "" // linha em branco depois do conteúdo encerra o campo
			}
			continue
		} else {
			linha = reMarcador.ReplaceAllString(linha, "")
		}
		if item := limpaRealce(linha); item != "" {
			campos[atual] = append(campos[atual], item)
		}
	}
	return campos
}

// itensDoCampo separa os itens de uma lista. Várias linhas valem um item por linha (a vírgula faz parte do texto);
// uma linha só vale uma lista separada por vírgulas.
func itensDoCampo(linhas []string) []string {
	if len(linhas) > 1 {
		return linhas
	}
	var out []string
	for _, l := range linhas {
		for _, item := range splitItems(l) {
			if item != "" {
				out = append(out, item)
			}
		}
	}
	return out
}
