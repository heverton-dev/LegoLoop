package matriz

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"legoloop/internal/reputation"
)

func entrada(ref string, score float64, p *PerfilDeUso) Entrada {
	return Entrada{Ref: ref, Nome: ref, Tipo: "skill", Hash: "h-" + ref, Perfil: p,
		Rating: reputation.PieceRating{Ref: ref, TotalScore: score, Tier: reputation.TierPrata}}
}

func perfil(etapas, naturezas []string, casos ...string) *PerfilDeUso {
	return &PerfilDeUso{Etapas: etapas, Naturezas: naturezas, CasosDeUso: casos, Origem: OrigemDeclarado}
}

// matriz42 é a matriz do AC-E02.6: aidd-tdd, aidd-site-dominator e 40 outras que não servem a software.
func matriz42(scoreSite float64, tdd *PerfilDeUso) []Entrada {
	site := entrada("aidd-site-dominator", scoreSite, perfil([]string{"publicar"}, []string{"midia"}))
	site.Perfil.QuandoNaoUsar = []string{"software"}
	es := []Entrada{entrada("aidd-tdd", 0.58, tdd), site}
	for i := 0; i < 40; i++ {
		es = append(es, entrada(fmt.Sprintf("outra-%02d", i), 0.5, perfil([]string{"documentar"}, []string{"documentacao"})))
	}
	return es
}

func tddPadrao() *PerfilDeUso {
	return perfil([]string{"construir", "verificar"}, []string{"software"}, "tdd")
}

var precisaSoftware = Needs{Natureza: "software", Etapas: []string{"construir", "verificar"}}

func achar(t *testing.T, r []Ranked, ref string) Ranked {
	t.Helper()
	for _, x := range r {
		if x.Entrada.Ref == ref {
			return x
		}
	}
	t.Fatalf("%s não está no resultado", ref)
	return Ranked{}
}

func escolhidas(r []Ranked) []string {
	var out []string
	for _, x := range r {
		if !x.Excluida {
			out = append(out, x.Entrada.Ref)
		}
	}
	return out
}

func TestRankEscolhePeloCasoDeUsoEExcluiComMotivo(t *testing.T) {
	r := Rank(precisaSoftware, matriz42(0.90, tddPadrao()), 0)
	if !slices.Contains(escolhidas(r), "aidd-tdd") {
		t.Fatalf("aidd-tdd devia estar entre as escolhidas: %v", escolhidas(r))
	}
	site := achar(t, r, "aidd-site-dominator")
	if !site.Excluida || !strings.Contains(site.Motivo, "quando_nao_usar") || !strings.Contains(site.Motivo, "software") {
		t.Errorf("aidd-site-dominator devia sair com o motivo: %+v", site)
	}
	if site.Aptidao != 0 {
		t.Errorf("peça excluída tem aptidão 0, veio %v", site.Aptidao)
	}
}

func TestRankScoreNaoMudaOResultado(t *testing.T) {
	a := Rank(precisaSoftware, matriz42(0.90, tddPadrao()), 0)
	b := Rank(precisaSoftware, matriz42(1.0, tddPadrao()), 0)
	if !slices.Equal(escolhidas(a), escolhidas(b)) {
		t.Errorf("o Score da peça excluída alterou as escolhidas: %v vs %v", escolhidas(a), escolhidas(b))
	}
	if !achar(t, b, "aidd-site-dominator").Excluida {
		t.Error("Score 1.0 não pode reabilitar a peça excluída")
	}
}

func TestRankMudarOPerfilTiraAPeca(t *testing.T) {
	r := Rank(precisaSoftware, matriz42(0.90, perfil([]string{"publicar"}, []string{"software"}, "tdd")), 0)
	tdd := achar(t, r, "aidd-tdd")
	if !tdd.Excluida || !strings.Contains(tdd.Motivo, "etapa") {
		t.Errorf("com etapas [publicar] a peça sai por etapa: %+v", tdd)
	}
}

func TestRankPerfilNilETierToxicoNuncaEntram(t *testing.T) {
	semPerfil := entrada("sem-perfil", 0.95, nil)
	toxica := entrada("toxica", 0.99, tddPadrao())
	toxica.Rating.Tier = reputation.TierToxica
	r := Rank(precisaSoftware, []Entrada{semPerfil, toxica, entrada("boa", 0.1, tddPadrao())}, 0)
	for _, x := range r {
		if x.Entrada.Ref == "sem-perfil" || x.Entrada.Ref == "toxica" {
			t.Errorf("%s não pode aparecer no resultado: %+v", x.Entrada.Ref, x)
		}
	}
	if got := escolhidas(r); len(got) != 1 || got[0] != "boa" {
		t.Errorf("escolhidas = %v", got)
	}
}

func TestRankNaturezaCoringaEExplicitaDesempata(t *testing.T) {
	coringa := entrada("coringa", 0.99, perfil([]string{"construir"}, nil)) // natureza vazia = qualquer
	qualquer := entrada("qualquer", 0.98, perfil([]string{"construir"}, []string{"qualquer"}))
	explicita := entrada("explicita", 0.10, perfil([]string{"construir"}, []string{"software"}))
	midia := entrada("midia", 0.99, perfil([]string{"construir"}, []string{"midia"}))
	r := Rank(Needs{Natureza: "software", Etapas: []string{"construir"}}, []Entrada{coringa, qualquer, explicita, midia}, 0)
	if got := escolhidas(r); !slices.Equal(got, []string{"explicita", "coringa", "qualquer"}) {
		t.Errorf("a natureza explícita desempata, mesmo com Score menor: %v", got)
	}
	if m := achar(t, r, "midia"); !m.Excluida || !strings.Contains(m.Motivo, "natureza") {
		t.Errorf("peça só de mídia sai por natureza: %+v", m)
	}
}

func TestRankAptidaoEhFracaoDosCasosPedidos(t *testing.T) {
	two := entrada("dois", 0.1, perfil([]string{"construir"}, nil, "tdd", "refatorar"))
	one := entrada("um", 0.9, perfil([]string{"construir"}, nil, "tdd"))
	none := entrada("zero", 0.99, perfil([]string{"construir"}, nil, "outro"))
	n := Needs{Etapas: []string{"construir"}, Casos: []string{"TDD", "refatorar"}} // maiúsculas casam
	r := Rank(n, []Entrada{one, none, two}, 0)
	if got := escolhidas(r); !slices.Equal(got, []string{"dois", "um", "zero"}) {
		t.Errorf("ordem por aptidão: %v", got)
	}
	if a := achar(t, r, "dois").Aptidao; a != 1 {
		t.Errorf("2 de 2 casos = 1, veio %v", a)
	}
	if a := achar(t, r, "um").Aptidao; a != 0.5 {
		t.Errorf("1 de 2 casos = 0.5, veio %v", a)
	}
	r = Rank(n, []Entrada{one, none, two}, 0.5)
	if z := achar(t, r, "zero"); !z.Excluida || !strings.Contains(z.Motivo, "mínimo") {
		t.Errorf("abaixo do mínimo sai com motivo: %+v", z)
	}
	if achar(t, r, "um").Excluida {
		t.Error("aptidão igual ao mínimo ainda passa")
	}
}

func TestRankSemCasosPedidosTodasAptasTemAptidaoUm(t *testing.T) {
	r := Rank(Needs{Etapas: []string{"construir"}}, []Entrada{entrada("a", 0.5, perfil([]string{"construir"}, nil))}, 1)
	if x := achar(t, r, "a"); x.Excluida || x.Aptidao != 1 {
		t.Errorf("sem casos pedidos a aptidão é 1: %+v", x)
	}
}

func TestRankDesempateScoreDepoisRef(t *testing.T) {
	p := func() *PerfilDeUso { return perfil([]string{"construir"}, []string{"software"}) }
	es := []Entrada{entrada("b", 0.5, p()), entrada("a", 0.5, p()), entrada("c", 0.9, p())}
	r := Rank(Needs{Natureza: "software", Etapas: []string{"construir"}}, es, 0)
	if got := escolhidas(r); !slices.Equal(got, []string{"c", "a", "b"}) {
		t.Errorf("Score desc e depois Ref asc: %v", got)
	}
}

func TestRankEhPura(t *testing.T) {
	es := matriz42(0.9, tddPadrao())
	antes := slices.Clone(es)
	a := Rank(precisaSoftware, es, 0)
	b := Rank(precisaSoftware, es, 0)
	if fmt.Sprint(a) != fmt.Sprint(b) {
		t.Error("duas chamadas iguais devem dar o mesmo resultado")
	}
	for i := range es {
		if es[i].Ref != antes[i].Ref {
			t.Fatal("Rank reordenou a entrada")
		}
	}
}

func TestRankQuandoNaoUsarAgeComOsCasosPedidos(t *testing.T) {
	e := entrada("x", 0.5, perfil([]string{"construir"}, nil, "tdd"))
	e.Perfil.QuandoNaoUsar = []string{"Prototipo"}
	r := Rank(Needs{Etapas: []string{"construir"}, Casos: []string{"prototipo"}}, []Entrada{e}, 0)
	if x := achar(t, r, "x"); !x.Excluida || !strings.Contains(x.Motivo, "quando_nao_usar") {
		t.Errorf("quando_nao_usar também vale para os casos pedidos: %+v", x)
	}
}

func TestRankMotivoDizOQueCasou(t *testing.T) {
	r := Rank(Needs{Natureza: "software", Etapas: []string{"construir", "verificar"}, Casos: []string{"tdd"}}, matriz42(0.9, tddPadrao()), 0)
	m := achar(t, r, "aidd-tdd").Motivo
	for _, want := range []string{"software", "construir", "verificar", "tdd"} {
		if !strings.Contains(m, want) {
			t.Errorf("o motivo da escolha não cita %q: %s", want, m)
		}
	}
}

func TestRankMinimoMaiorQueUmSemCasosNaoQuebra(t *testing.T) {
	r := Rank(Needs{Etapas: []string{"construir"}}, []Entrada{entrada("a", 0.5, perfil([]string{"construir"}, nil))}, 1.5)
	if x := achar(t, r, "a"); !x.Excluida || strings.Contains(x.Motivo, "NaN") || strings.Contains(x.Motivo, "Inf") {
		t.Errorf("motivo inválido: %+v", x)
	}
}
