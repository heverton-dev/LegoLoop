package processo

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrLacuna: a etapa não tem peça nem processo composto, então não há o que executar nela.
	ErrLacuna = errors.New("processo: etapa sem peça (lacuna)")
	// ErrIncoerente: as projeções do processo não listam as mesmas etapas, na mesma ordem.
	ErrIncoerente = errors.New("processo: as projeções não contam a mesma história")
)

var (
	// reIDEtapa é a mesma regra de id de nó do DAG: o id da etapa é o id do nó.
	reIDEtapa = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)
	// rePeca: o nome da peça vem de uma varredura (dado não confiável) e vira ref de nó e texto de documento.
	rePeca = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)
	// reCompoe: "legoloop-<nome>@<maior.menor.correção>".
	reCompoe = regexp.MustCompile(`^(legoloop-[a-z0-9]+(-[a-z0-9]+)*)@(\d+\.\d+\.\d+)$`)
)

// Ordem devolve as etapas em ordem de dependência: cada etapa vem depois das que ela depende. Entre etapas
// prontas ao mesmo tempo vale a ordem em que foram declaradas. Todas as projeções usam esta ordem.
func (p Processo) Ordem() ([]Etapa, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	feitas := make(map[string]bool, len(p.Etapas))
	out := make([]Etapa, 0, len(p.Etapas))
	for len(out) < len(p.Etapas) {
		avancou := false
		for _, e := range p.Etapas {
			if !feitas[e.ID] && prontas(e, feitas) {
				feitas[e.ID] = true
				out = append(out, e)
				avancou = true
				break
			}
		}
		if !avancou { // Validate já recusa ciclos; isto só protege contra laço infinito
			return nil, invalido("não foi possível ordenar as etapas")
		}
	}
	return out, nil
}

func prontas(e Etapa, feitas map[string]bool) bool {
	for _, d := range e.Depende {
		if !feitas[d] {
			return false
		}
	}
	return true
}

// validarProjecao confere, além do contrato do processo, o que as projeções exigem: ids e peças que viram
// referência segura, caminhos que cabem num texto Markdown, e nenhuma lacuna (lacuna não tem nó executável).
func validarProjecao(p Processo) ([]Etapa, error) {
	ord, err := p.Ordem()
	if err != nil {
		return nil, err
	}
	for _, e := range ord {
		if err := validarEtapa(e); err != nil {
			return nil, err
		}
	}
	return ord, nil
}

func validarEtapa(e Etapa) error {
	if !reIDEtapa.MatchString(e.ID) {
		return invalido("etapa %q: o id deve ter só letras, números, hífen e sublinhado", e.ID)
	}
	if strings.ContainsAny(e.Produz.Caminho, "`\r\n\t ") {
		return invalido("etapa %q: o caminho %q não pode ter espaço, crase nem quebra de linha", e.ID, e.Produz.Caminho)
	}
	switch {
	case e.Peca != "" && e.Compoe != "":
		return invalido("etapa %q: use peça ou processo composto, não os dois", e.ID)
	case e.Peca == "" && e.Compoe == "":
		return invalidoLacuna(e.ID)
	case e.Peca != "" && !rePeca.MatchString(e.Peca):
		return invalido("etapa %q: o nome da peça %q não é seguro", e.ID, e.Peca)
	case e.Compoe != "" && !reCompoe.MatchString(e.Compoe):
		return invalido("etapa %q: composição %q (use legoloop-<nome>@maior.menor.correção)", e.ID, e.Compoe)
	}
	return nil
}

func invalidoLacuna(id string) error {
	return errors.Join(ErrLacuna, invalido("etapa %q sem peça: escolha uma peça antes de gerar o processo", id))
}

// NomeDePecaValido diz se o nome pode virar peça de etapa: ref de nó e texto de documento, sem espaço nem
// caractere de comando. A Fusão usa a mesma regra do contrato para nunca escolher uma peça que ele recusaria.
func NomeDePecaValido(nome string) bool { return rePeca.MatchString(nome) }

// nomeComposto devolve o nome do processo composto ("legoloop-x" de "legoloop-x@1.0.0").
func nomeComposto(compoe string) string { return reCompoe.FindStringSubmatch(compoe)[1] }
