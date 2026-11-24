package roteiro

import (
	"fmt"
	"strings"
)

type ViolacaoRoteiro struct {
	Linha        int    `json:"linha"`
	Natureza     string `json:"natureza"`
	RegraViolada string `json:"regra_violada"`
}

type RelatorioValidacao struct {
	Aprovado  bool              `json:"aprovado"`
	Violacoes []ViolacaoRoteiro `json:"violacoes,omitempty"`
}

// ValidarRoteiro aplica as regras epistemicas universais para as 4 naturezas e bloqueia nivel C.
func ValidarRoteiro(sc Script, p Projeto, cat CatalogoFontes) RelatorioValidacao {
	var rel RelatorioValidacao
	if len(sc.Lines) == 0 {
		rel.Violacoes = append(rel.Violacoes, ViolacaoRoteiro{Linha: 0, RegraViolada: "roteiro sem linhas"})
		return rel
	}
	mapa := indexarFontes(sc.Sources, p, cat)
	for i, linha := range sc.Lines {
		validarLinha(i, linha, mapa, &rel)
	}
	rel.Aprovado = len(rel.Violacoes) == 0
	return rel
}

func indexarFontes(fontes []Fonte, p Projeto, cat CatalogoFontes) map[string]ClassificacaoFonte {
	mapa := make(map[string]ClassificacaoFonte, len(fontes))
	for _, f := range fontes {
		mapa[f.ID] = ClassificarFonte(f.URL, p, cat)
	}
	return mapa
}

func validarLinha(idx int, l LinhaRoteiro, mapa map[string]ClassificacaoFonte, rel *RelatorioValidacao) {
	if strings.TrimSpace(l.Text) == "" {
		addViolacao(rel, idx, l.Natureza, "texto da linha vazio")
		return
	}
	switch l.Natureza {
	case "narrativa_transicao":
		return
	case "fato_consolidado":
		validarFatoConsolidado(idx, l, mapa, rel)
	case "evidencia_preliminar":
		validarEvidenciaPreliminar(idx, l, mapa, rel)
	case "alegacao_terceiro":
		validarAlegacaoTerceiro(idx, l, mapa, rel)
	default:
		addViolacao(rel, idx, l.Natureza, fmt.Sprintf("natureza invalida %q", l.Natureza))
	}
}

func validarFatoConsolidado(idx int, l LinhaRoteiro, mapa map[string]ClassificacaoFonte, rel *RelatorioValidacao) {
	a, b, preprints, invalidas := contarFontesLinha(l.Fontes, mapa)
	if invalidas > 0 {
		addViolacao(rel, idx, l.Natureza, "fonte nivel C ou inexistente bloqueada como prova")
		return
	}
	if preprints > 0 {
		addViolacao(rel, idx, l.Natureza, "preprint exige natureza evidencia_preliminar")
		return
	}
	if a < 1 && b < 2 {
		addViolacao(rel, idx, l.Natureza, "fato_consolidado exige >=1 fonte A ou >=2 fontes B distintas")
	}
}

func validarEvidenciaPreliminar(idx int, l LinhaRoteiro, mapa map[string]ClassificacaoFonte, rel *RelatorioValidacao) {
	a, b, preprints, invalidas := contarFontesLinha(l.Fontes, mapa)
	if invalidas > 0 || (a+b+preprints) == 0 {
		addViolacao(rel, idx, l.Natureza, "evidencia_preliminar exige fonte valida A, B ou preprint")
		return
	}
	if strings.TrimSpace(l.Escopo) == "" && !l.Preprint && strings.TrimSpace(l.Amostra) == "" {
		addViolacao(rel, idx, l.Natureza, "evidencia_preliminar exige escopo, preprint ou amostra declarados")
	}
}

func validarAlegacaoTerceiro(idx int, l LinhaRoteiro, mapa map[string]ClassificacaoFonte, rel *RelatorioValidacao) {
	a, b, preprints, invalidas := contarFontesLinha(l.Fontes, mapa)
	if invalidas > 0 || (a+b+preprints) == 0 {
		addViolacao(rel, idx, l.Natureza, "alegacao_terceiro exige fonte valida A ou B")
		return
	}
	if strings.TrimSpace(l.Atribuicao) == "" {
		addViolacao(rel, idx, l.Natureza, "alegacao_terceiro exige atribuicao explicita a quem afirmou")
	}
}

func contarFontesLinha(ids []string, mapa map[string]ClassificacaoFonte) (int, int, int, int) {
	var a, b, preprints, invalidas int
	dominiosB := map[string]bool{}
	for _, id := range ids {
		cls, ok := mapa[id]
		if !ok || cls.Nivel == "C" {
			invalidas++
			continue
		}
		if cls.Preprint {
			preprints++
			continue
		}
		if cls.Nivel == "A" {
			a++
		} else if cls.Nivel == "B" && !dominiosB[cls.Dominio] {
			dominiosB[cls.Dominio] = true
			b++
		}
	}
	return a, b, preprints, invalidas
}

func addViolacao(rel *RelatorioValidacao, idx int, natureza, regra string) {
	rel.Violacoes = append(rel.Violacoes, ViolacaoRoteiro{
		Linha: idx, Natureza: natureza, RegraViolada: regra,
	})
}
