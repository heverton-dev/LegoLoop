package refinar

import (
	"regexp"
	"slices"
	"strings"

	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

// reIDNomeado aceita o id da etapa: minúsculas, números, _ ou -, começando com letra. Vale e04 e renderizar_video.
var reIDNomeado = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

func init() {
	ste.Default.Register("refinar.etapa_existe", "A etapa %s já existe em %s. Para mudar uma etapa existente, tire --adicionar.")
	ste.Default.Register("refinar.adicionar_sem_peca", "A etapa nova %s precisa de uma peça. Informe com --peca.")
	ste.Default.Register("refinar.depende_desconhecido", "A etapa %s depende de %s, que não existe. As etapas são: %s.")
	ste.Default.Register("refinar.id_invalido", "O id %s não vale. Use letras minúsculas, números, _ ou -, começando com letra, como e04 ou renderizar_video.")
}

// adicionarEtapa põe uma etapa nova no fim do processo, ligada a Depende (padrão: a última etapa).
// A etapa nova produz out/<id>.md e ganha o gate de arquivo desse artefato.
func adicionarEtapa(p processo.Processo, m Mudanca) (processo.Processo, processo.Etapa, error) {
	id, peca, casos := strings.TrimSpace(m.Etapa), strings.TrimSpace(m.Peca), limpar(m.Casos)
	ids := idsDe(p.Etapas)
	switch {
	case !idValido(id):
		return p, processo.Etapa{}, erroMudanca("refinar.id_invalido", id)
	case slices.Contains(ids, id):
		return p, processo.Etapa{}, erroMudanca("refinar.etapa_existe", id, p.Nome)
	case peca == "":
		return p, processo.Etapa{}, erroMudanca("refinar.adicionar_sem_peca", id)
	case len(casos) == 0:
		return p, processo.Etapa{}, erroMudanca("refinar.sem_caso", id)
	}
	deps, err := dependenciasDe(p, m.Depende, ids)
	if err != nil {
		return p, processo.Etapa{}, err
	}
	art, gate, err := artefatoDaMudanca(m)
	if err != nil {
		return p, processo.Etapa{}, err
	}
	titulo := strings.TrimSpace(m.Titulo)
	if titulo == "" {
		titulo = "Etapa " + id
	}
	e := processo.Etapa{ID: id, Titulo: titulo, Peca: peca, Casos: casos, Natureza: m.Natureza, Produz: art, Gate: gate, Depende: deps}
	for _, d := range deps {
		e.Consome = append(e.Consome, p.Etapas[indiceDe(p.Etapas, d)].Produz.Caminho)
	}
	p.Etapas = append(slices.Clone(p.Etapas), e)
	return p, processo.Etapa{}, nil
}

// artefatoDaMudanca devolve o artefato da etapa nova e o gate que o cobre. O artefato vem
// do Sai pedido; sem gate pedido, o gate deriva do artefato.
func artefatoDaMudanca(m Mudanca) (processo.Artefato, string, error) {
	art := processo.Artefato{Caminho: "out/" + strings.TrimSpace(m.Etapa) + ".md", Tipo: "markdown", Promete: true}
	if m.Sai != nil {
		art = *m.Sai
	}
	gate := m.Gate
	if gate == "" {
		var err error
		if gate, err = processo.GateParaArtefato(art); err != nil {
			return art, "", err
		}
	}
	return art, gate, nil
}

func idsDe(etapas []processo.Etapa) []string {
	ids := make([]string, len(etapas))
	for i, e := range etapas {
		ids[i] = e.ID
	}
	return ids
}

func idValido(id string) bool { return reIDNomeado.MatchString(id) }

// dependenciasDe confere cada id pedido; sem pedido, liga na última etapa.
func dependenciasDe(p processo.Processo, pedidas, ids []string) ([]string, error) {
	deps := limpar(pedidas)
	if len(deps) == 0 {
		if len(ids) == 0 {
			return nil, nil
		}
		return []string{ids[len(ids)-1]}, nil
	}
	for _, d := range deps {
		if !slices.Contains(ids, d) {
			return nil, erroMudanca("refinar.depende_desconhecido", "a etapa nova", d, strings.Join(ids, ", "))
		}
	}
	return deps, nil
}
