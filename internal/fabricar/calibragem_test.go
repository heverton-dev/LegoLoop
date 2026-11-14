package fabricar_test

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/fabricar"
	"legoloop/internal/matriz"
)

// pecaComCasos é uma skill com ref no formato real da matriz ("~/caminho/SKILL.md").
func pecaComCasos(nome string, score float64, etapas, naturezas, casos []string) matriz.Entrada {
	e := skill(nome, score, etapas, naturezas)
	e.Ref = "~/.agents/skills/" + nome + "/SKILL.md"
	e.Perfil.CasosDeUso = casos
	return e
}

func temAviso(m fabricar.Merged, partes ...string) bool {
	return slices.ContainsFunc(m.Warnings, func(w string) bool {
		for _, p := range partes {
			if !strings.Contains(w, p) {
				return false
			}
		}
		return true
	})
}

var dados = []string{"dados"}

func TestCalibragem(t *testing.T) {
	t.Run("peça de aptidão 0 não entra: a etapa vira lacuna com motivo", func(t *testing.T) {
		generica := pecaComCasos("super-qa", 0.95, []string{"coletar"}, []string{"qualquer"}, nil)
		n := matriz.Needs{Natureza: "dados", Etapas: []string{"coletar"}, Casos: []string{"coleta-de-dados"}}
		m := fabricar.MergeByFit(n, []matriz.Entrada{generica}, fabricar.MergeOptions{})
		if slices.Contains(nomes(m), "super-qa") {
			t.Fatalf("peça que só casa etapa e natureza coringa não pode entrar: %v", nomes(m))
		}
		if !slices.Equal(m.Lacunas, []string{"coletar"}) {
			t.Fatalf("Lacunas = %v", m.Lacunas)
		}
		if d := decisaoDe(t, m, "coletar"); d.Entrou || !strings.Contains(d.Motivo, "coleta-de-dados") {
			t.Errorf("o motivo da lacuna deve citar o caso que ninguém cobre: %+v", d)
		}
	})

	t.Run("peça que cobre o caso da etapa entra com aptidão 1", func(t *testing.T) {
		generica := pecaComCasos("super-qa", 0.95, []string{"coletar"}, []string{"qualquer"}, nil)
		coletor := pecaComCasos("coletor", 0.30, []string{"coletar"}, dados, []string{"coleta-de-dados"})
		n := matriz.Needs{Natureza: "dados", Etapas: []string{"coletar"}, Casos: []string{"coleta-de-dados"}}
		m := fabricar.MergeByFit(n, []matriz.Entrada{generica, coletor}, fabricar.MergeOptions{})
		if got := nomes(m); !slices.Equal(got, []string{"coletor"}) {
			t.Fatalf("a aptidão decide, não o Score: %v", got)
		}
		if d := decisaoDe(t, m, "coletar"); d.Aptidao != 1 {
			t.Errorf("aptidão = %v, quer 1", d.Aptidao)
		}
	})

	t.Run("o caso conta só na etapa a que pertence", func(t *testing.T) {
		// Se os dois casos contassem nas duas etapas, construtor e qa empatariam em 0.5 em construir
		// e o Score (0.90) daria construir ao qa.
		construtor := pecaComCasos("construtor-tdd", 0.30, []string{"construir"}, sw, []string{"tdd"})
		qa := pecaComCasos("qa", 0.90, []string{"construir", "verificar"}, sw, []string{"verificacao-de-qualidade"})
		n := matriz.Needs{Natureza: "software", Etapas: []string{"construir", "verificar"}, Casos: []string{"tdd", "verificacao-de-qualidade"}}
		m := fabricar.MergeByFit(n, []matriz.Entrada{construtor, qa}, fabricar.MergeOptions{})
		if got := nomes(m); !slices.Equal(got, []string{"construtor-tdd", "qa"}) {
			t.Fatalf("construir casa tdd e verificar casa qualidade: %v", got)
		}
		for _, etapa := range []string{"construir", "verificar"} {
			if d := decisaoDe(t, m, etapa); d.Aptidao != 1 {
				t.Errorf("%s: aptidão = %v, quer 1 (só o caso da etapa conta)", etapa, d.Aptidao)
			}
		}
		if temAviso(m, "não contou") {
			t.Errorf("todo caso pertence a uma etapa pedida; nada a avisar: %v", m.Warnings)
		}
	})

	t.Run("etapa sem caso de uso pedido vira lacuna com motivo", func(t *testing.T) {
		coletor := pecaComCasos("coletor", 0.30, []string{"coletar"}, dados, []string{"coleta-de-dados"})
		transformador := pecaComCasos("transformador", 0.99, []string{"transformar"}, dados, nil)
		n := matriz.Needs{Natureza: "dados", Etapas: []string{"coletar", "transformar"}, Casos: []string{"coleta-de-dados"}}
		m := fabricar.MergeByFit(n, []matriz.Entrada{coletor, transformador}, fabricar.MergeOptions{})
		if got := nomes(m); !slices.Equal(got, []string{"coletor"}) {
			t.Fatalf("sem caso pedido o Score não pode escolher sozinho: %v", got)
		}
		if !slices.Equal(m.Lacunas, []string{"transformar"}) {
			t.Fatalf("Lacunas = %v", m.Lacunas)
		}
		if d := decisaoDe(t, m, "transformar"); !strings.Contains(d.Motivo, "nenhum caso de uso pedido") {
			t.Errorf("motivo = %q", d.Motivo)
		}
		if !temAviso(m, "Lacuna", "transformar") {
			t.Errorf("a lacuna deve ser destacada: %v", m.Warnings)
		}
	})

	t.Run("caso pedido que não pertence a nenhuma etapa pedida é avisado", func(t *testing.T) {
		construtor := pecaComCasos("construtor-tdd", 0.30, []string{"construir"}, sw, []string{"tdd"})
		n := matriz.Needs{Natureza: "software", Etapas: []string{"construir"}, Casos: []string{"tdd", "publicacao"}}
		m := fabricar.MergeByFit(n, []matriz.Entrada{construtor}, fabricar.MergeOptions{})
		if got := nomes(m); !slices.Equal(got, []string{"construtor-tdd"}) {
			t.Fatalf("passos = %v", got)
		}
		if !temAviso(m, "publicacao") {
			t.Errorf("o caso sem etapa pedida deve ser avisado: %v", m.Warnings)
		}
	})

	t.Run("sem nenhum caso pedido a etapa basta, como antes", func(t *testing.T) {
		construtor := pecaComCasos("construtor", 0.30, []string{"construir"}, sw, nil)
		m := fabricar.MergeByFit(matriz.Needs{Natureza: "software", Etapas: []string{"construir"}}, []matriz.Entrada{construtor}, fabricar.MergeOptions{})
		if got := nomes(m); !slices.Equal(got, []string{"construtor"}) || len(m.Lacunas) != 0 {
			t.Errorf("passos = %v, lacunas = %v", got, m.Lacunas)
		}
	})

	t.Run("a etapa do caso vem do vocabulário das opções", func(t *testing.T) {
		v := matriz.DefaultVocab()
		v.CasosDeUso = append(v.CasosDeUso, matriz.Termo{Nome: "raspagem", Etapas: []string{"coletar"}})
		raspador := pecaComCasos("raspador", 0.30, []string{"coletar"}, dados, []string{"raspagem"})
		n := matriz.Needs{Natureza: "dados", Etapas: []string{"coletar"}, Casos: []string{"raspagem"}}
		m := fabricar.MergeByFit(n, []matriz.Entrada{raspador}, fabricar.MergeOptions{Vocab: v})
		if got := nomes(m); !slices.Equal(got, []string{"raspador"}) {
			t.Errorf("passos = %v, avisos = %v", got, m.Warnings)
		}
	})
}

// O motor usa o vocabulário do cofre, não o embutido: se o cofre diz que tdd é de verificar, a etapa construir
// fica sem caso pedido e vira lacuna (com o padrão, alfa entraria em construir).
func TestCalibragemMotorUsaOVocabularioDoCofre(t *testing.T) {
	r := newRig(t, goodPieces())
	v := matriz.DefaultVocab()
	for i, c := range v.CasosDeUso {
		if c.Nome == "tdd" {
			v.CasosDeUso[i].Etapas = []string{"verificar"}
		}
	}
	r.deps.Matriz = func() (fabricar.FonteMatriz, error) {
		return fabricar.FonteMatriz{Entradas: goodPieces(), Vocab: v, GeradaEm: fabNow.Add(-48 * time.Hour)}, nil
	}
	res, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
	if err != nil {
		t.Fatalf("cascata de 3 fontes não deve bloquear lacuna em construir: %v", err)
	}
	if !slices.ContainsFunc(res.Warnings, func(w string) bool { return strings.Contains(w, "construir") }) {
		t.Errorf("esperava aviso de lacuna em construir pelo vocabulário do cofre: %v", res.Warnings)
	}
}

func TestCalibragemVocabulario(t *testing.T) {
	v := matriz.DefaultVocab()
	if v.Versao < 2 {
		t.Errorf("o vocabulário mudou de formato (etapas por caso): versão = %d", v.Versao)
	}
	for _, c := range v.CasosDeUso {
		etapas := v.EtapasDoCaso(c.Nome)
		if len(etapas) == 0 {
			t.Errorf("o caso %q não diz a que etapa pertence", c.Nome)
		}
		for _, e := range etapas {
			if !v.Has(matriz.CatEtapa, e) {
				t.Errorf("o caso %q aponta a etapa %q, que não está no vocabulário", c.Nome, e)
			}
		}
	}
	for caso, etapa := range map[string]string{"tdd": "construir", "coleta-de-dados": "coletar", "especificacao": "planejar",
		"verificacao-de-qualidade": "verificar", "publicacao": "publicar", "entrevista-de-requisitos": "clarificar"} {
		if !slices.Contains(v.EtapasDoCaso(caso), etapa) {
			t.Errorf("%s deve pertencer a %s: %v", caso, etapa, v.EtapasDoCaso(caso))
		}
	}
	// Vocabulário de cofre antigo (sem etapas nos casos) cai no padrão; caso novo sem etapa fica sem etapa.
	antigo := matriz.Vocab{Versao: 1, CasosDeUso: []matriz.Termo{{Nome: "tdd"}, {Nome: "novo-caso"}}}
	if got := antigo.EtapasDoCaso("tdd"); !slices.Equal(got, []string{"construir"}) {
		t.Errorf("cofre antigo: tdd = %v", got)
	}
	if got := antigo.EtapasDoCaso("novo-caso"); len(got) != 0 {
		t.Errorf("caso sem etapa declarada: %v", got)
	}
}
