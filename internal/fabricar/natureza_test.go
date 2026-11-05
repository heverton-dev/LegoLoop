package fabricar_test

import (
	stdctx "context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/fabricar"
	"legoloop/internal/gate"
	"legoloop/internal/interview"
	"legoloop/internal/matriz"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

var agora = time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)

func bg() stdctx.Context { return stdctx.Background() }

const objetivoPrecos = "coletar preços de 20 lojas e entregar uma planilha"

func espelhoPrecos() interview.Mirror {
	return interview.Mirror{Outcome: objetivoPrecos, User: "o time de compras", WhyNow: "comparar fornecedores",
		Success: "uma planilha com 20 lojas e seus preços", Constraints: []string{"só fontes públicas"}, OutOfScope: []string{"negociar preços"}}
}

func TestNaturezaAbertaCerebroTraduzOEspelhoEmNecessidades(t *testing.T) {
	c := &cerebro{reply: `{"natureza":"dados","etapas":["coletar","transformar","entregar"],"casos":["coleta-de-dados"]}`}
	n, avisos, _ := fabricar.Necessidades(bg(), c, objetivoPrecos, espelhoPrecos(), matriz.DefaultVocab())
	if n.Natureza != "dados" || !slices.Equal(n.Etapas, []string{"coletar", "transformar", "entregar"}) || !slices.Equal(n.Casos, []string{"coleta-de-dados"}) {
		t.Errorf("necessidades = %+v", n)
	}
	if len(avisos) != 0 {
		t.Errorf("sem problema não há aviso: %v", avisos)
	}
	p := c.visto()
	for _, want := range []string{objetivoPrecos, "só fontes públicas", "coletar", "software", "coleta-de-dados", "NÃO invente", "dados"} {
		if !strings.Contains(p, want) {
			t.Errorf("o prompt não traz %q", want)
		}
	}
}

func TestNaturezaAbertaTermosForaDoVocabularioSaem(t *testing.T) {
	c := &cerebro{reply: `{"natureza":"Culinária","etapas":["coletar","dançar","COLETAR","entregar"],"casos":["tdd","inventado"]}`}
	n, avisos, _ := fabricar.Necessidades(bg(), c, objetivoPrecos, espelhoPrecos(), matriz.DefaultVocab())
	if n.Natureza != "dados" {
		t.Errorf("natureza inventada cai na inferência pelas palavras do objetivo: %q", n.Natureza)
	}
	if !slices.Equal(n.Etapas, []string{"coletar", "entregar"}) || !slices.Equal(n.Casos, []string{"tdd"}) {
		t.Errorf("só termos do vocabulário, sem repetir: %+v", n)
	}
	all := strings.Join(avisos, "\n")
	for _, want := range []string{"Culinária", "dançar", "inventado"} {
		if !strings.Contains(all, want) {
			t.Errorf("o aviso deve citar o termo descartado %q: %v", want, avisos)
		}
	}
	for _, a := range avisos {
		if v := ste.Lint(a); len(v) != 0 {
			t.Errorf("aviso fora do perfil STE-PT: %q → %v", a, v)
		}
	}
}

func TestNaturezaAbertaNaturezaIgnoraMaiusculasEEspacos(t *testing.T) {
	c := &cerebro{reply: `{"natureza":"  Dados ","etapas":[" Coletar "]}`}
	n, _, _ := fabricar.Necessidades(bg(), c, "x", interview.Mirror{}, matriz.DefaultVocab())
	if n.Natureza != "dados" || !slices.Equal(n.Etapas, []string{"coletar"}) {
		t.Errorf("%+v", n)
	}
}

// TestNaturezaAbertaCerebroCaidoParaAFabricacao: cérebro fora do ar ou expirado não vira inferência silenciosa;
// a fabricação para com o erro, que mantém a causa (o prazo vencido continua reconhecível).
func TestNaturezaAbertaCerebroCaidoParaAFabricacao(t *testing.T) {
	for nome, causa := range map[string]error{"fora do ar": errors.New("fora do ar"), "prazo": stdctx.DeadlineExceeded} {
		t.Run(nome, func(t *testing.T) {
			_, _, err := fabricar.Necessidades(bg(), &cerebro{err: causa}, objetivoPrecos, espelhoPrecos(), matriz.DefaultVocab())
			if !errors.Is(err, causa) {
				t.Errorf("esperava o erro do cérebro, veio %v", err)
			}
		})
	}
}

func TestNaturezaAbertaSemCerebroInferePelasPalavras(t *testing.T) {
	for nome, c := range map[string]*cerebro{
		"prosa":    {reply: "não sei dizer"},
		"quebrado": {reply: `{"natureza":`},
	} {
		t.Run(nome, func(t *testing.T) {
			n, avisos, _ := fabricar.Necessidades(bg(), c, objetivoPrecos, espelhoPrecos(), matriz.DefaultVocab())
			if n.Natureza != "dados" || !slices.Equal(n.Etapas, []string{"coletar", "entregar"}) {
				t.Errorf("a inferência determinística deve dar dados, coletar e entregar: %+v", n)
			}
			if !slices.ContainsFunc(avisos, func(a string) bool { return strings.Contains(a, "palavras do objetivo") }) {
				t.Errorf("o aviso de que a natureza veio das palavras falta: %v", avisos)
			}
		})
	}
	n, _, _ := fabricar.Necessidades(bg(), nil, objetivoPrecos, espelhoPrecos(), matriz.DefaultVocab())
	if n.Natureza != "dados" {
		t.Errorf("sem cérebro (nil) também infere: %+v", n)
	}
}

func TestNaturezaAbertaInferenciaOutrasNaturezas(t *testing.T) {
	for objetivo, quero := range map[string]string{
		"construir uma API com TDD":                "software",
		"escrever o manual de uso do produto":      "documentacao",
		"editar um vídeo de apresentação":          "midia",
		"pesquisar concorrentes e fazer benchmark": "pesquisa",
		"fazer uma coisa qualquer":                 "qualquer",
	} {
		n, avisos, _ := fabricar.Necessidades(bg(), nil, objetivo, interview.Mirror{}, matriz.DefaultVocab())
		if n.Natureza != quero {
			t.Errorf("%q → %q, queria %q", objetivo, n.Natureza, quero)
		}
		if quero == "qualquer" && !slices.ContainsFunc(avisos, func(a string) bool { return strings.Contains(a, "qualquer") }) {
			t.Errorf("sem pista alguma o aviso deve dizer que vale qualquer natureza: %v", avisos)
		}
	}
}

func TestNaturezaAbertaMistoEscolheAMaisCitada(t *testing.T) {
	n, _, _ := fabricar.Necessidades(bg(), nil, "gerar um csv, uma planilha e um dataset a partir da API", interview.Mirror{}, matriz.DefaultVocab())
	if n.Natureza != "dados" {
		t.Errorf("3 pistas de dados contra 1 de software: %q", n.Natureza)
	}
}

func TestNaturezaAbertaTextoDoEspelhoTambemConta(t *testing.T) {
	n, _, _ := fabricar.Necessidades(bg(), nil, "fazer o trabalho", interview.Mirror{Success: "uma planilha pronta"}, matriz.DefaultVocab())
	if n.Natureza != "dados" {
		t.Errorf("o critério de sucesso do espelho também dá pista: %q", n.Natureza)
	}
}

func TestNaturezaAbertaArtefatoPorEtapaENatureza(t *testing.T) {
	for _, c := range []struct {
		etapa, natureza string
		i               int
		tipo, caminho   string
	}{
		{"coletar", "dados", 1, "dados", "out/01_coletar.csv"},
		{"transformar", "dados", 2, "dados", "out/02_transformar.csv"},
		{"entregar", "dados", 3, "dados", "out/03_entregar.csv"},
		{"documentar", "dados", 4, "markdown", "out/04_documentar.md"},
		{"verificar", "dados", 5, "markdown", "out/05_verificar.md"},
		{"construir", "software", 2, "markdown", "out/02_construir.md"},
		{"coletar", "software", 1, "markdown", "out/01_coletar.md"},
		{"", "dados", 1, "markdown", "out/01_etapa.md"},
	} {
		a := fabricar.ArtefatoDe(c.etapa, c.natureza, c.i)
		if a.Tipo != c.tipo || a.Caminho != c.caminho || !a.Promete {
			t.Errorf("ArtefatoDe(%q,%q,%d) = %+v", c.etapa, c.natureza, c.i, a)
		}
	}
}

func fusaoDados(t *testing.T) (matriz.Needs, fabricar.Merged) {
	t.Helper()
	n, _, _ := fabricar.Necessidades(bg(), nil, objetivoPrecos, espelhoPrecos(), matriz.DefaultVocab())
	es := []matriz.Entrada{
		skill("coletor-web", 0.7, []string{"coletar"}, []string{"dados"}),
		skill("exportador", 0.6, []string{"entregar"}, []string{"dados"}),
	}
	return n, fabricar.MergeByFit(n, es, fabricar.MergeOptions{})
}

func TestNaturezaAbertaProcessoDeMontaEtapasEncadeadasComGateDeConteudo(t *testing.T) {
	n, m := fusaoDados(t)
	p, err := fabricar.ProcessoDe("legoloop-precos", objetivoPrecos, n, m, agora)
	if err != nil {
		t.Fatal(err)
	}
	if p.Natureza != "dados" || p.Objetivo != objetivoPrecos || len(p.Etapas) != 2 || p.Versao != "1.0.0" {
		t.Fatalf("processo = %+v", p)
	}
	e0, e1 := p.Etapas[0], p.Etapas[1]
	if e0.ID != "e01" || e0.Peca != "coletor-web" || e0.Tipo != "coletar" || e1.ID != "e02" || e1.Peca != "exportador" || !slices.Equal(e1.Depende, []string{"e01"}) || len(e0.Depende) != 0 {
		t.Errorf("etapas: %+v %+v", e0, e1)
	}
	g, _ := processo.GateParaArtefato(e0.Produz)
	if e0.Gate != g || !strings.Contains(e0.Gate, "line_count") {
		t.Errorf("o gate da etapa de dados confere registros: %q", e0.Gate)
	}
	if err := p.Validate(); err != nil {
		t.Errorf("o processo montado deve ser válido: %v", err)
	}
	if !slices.Equal(e1.Consome, []string{e0.Produz.Caminho}) {
		t.Errorf("a etapa consome o artefato da anterior: %v", e1.Consome)
	}
}

func TestNaturezaAbertaProcessoDeLacunaVaiParaOProjetar(t *testing.T) {
	n, _, _ := fabricar.Necessidades(bg(), nil, objetivoPrecos, espelhoPrecos(), matriz.DefaultVocab())
	m := fabricar.MergeByFit(n, []matriz.Entrada{skill("coletor-web", 0.7, []string{"coletar"}, []string{"dados"})}, fabricar.MergeOptions{})
	p, err := fabricar.ProcessoDe("legoloop-precos", objetivoPrecos, n, m, agora)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Etapas) != 2 || p.Etapas[1].Peca != "" {
		t.Fatalf("a etapa sem peça entra com Peca vazia: %+v", p.Etapas)
	}
	if _, err := processo.Projetar(p, agora); !errors.Is(err, processo.ErrLacuna) {
		t.Errorf("a lacuna é decisão do usuário, não do motor: %v", err)
	}
}

func TestNaturezaAbertaProcessoDeSemEtapasPedidasUsaOsPassosDaFusao(t *testing.T) {
	m := fabricar.Merged{Steps: []fabricar.Step{{Piece: "aidd-spec"}, {Piece: "aidd-tickets"}}}
	p, err := fabricar.ProcessoDe("legoloop-x", "fazer x", matriz.Needs{Natureza: "qualquer"}, m, agora)
	if err != nil || len(p.Etapas) != 2 || p.Etapas[0].Peca != "aidd-spec" || p.Etapas[1].Peca != "aidd-tickets" {
		t.Errorf("%+v %v", p, err)
	}
}

var (
	// proibidosAC são as expressões que o AC-E02.11 proíbe no INTENT e na ROUTE de um processo de dados.
	proibidosAC = []string{"go test", "go build", "zero stubs"}
	// proibidosTexto vale também para os textos escritos (INTENT, SKILL, SPEC): TDD é método de software.
	proibidosTexto = append(slices.Clone(proibidosAC), "tdd")
)

func semGo(t *testing.T, nome, texto string, proibidos ...string) {
	t.Helper()
	low := strings.ToLower(texto)
	for _, p := range proibidos {
		if strings.Contains(low, p) {
			t.Errorf("%s não pode falar de %q num processo de dados", nome, p)
		}
	}
}

func TestNaturezaAbertaPrecosViraProcessoDeDadosSemGo(t *testing.T) {
	n, m := fusaoDados(t)
	if n.Natureza != "dados" {
		t.Fatalf("natureza = %q", n.Natureza)
	}
	p, err := fabricar.ProcessoDe("legoloop-precos", objetivoPrecos, n, m, agora)
	if err != nil {
		t.Fatal(err)
	}
	pr := projetarOK(t, p)
	intent := interview.RenderIntent(interview.IntentView{Theme: objetivoPrecos, Date: "2026-11-04", SessionID: "s", Confidence: 98,
		Natureza: n.Natureza, Mirror: ptr(espelhoPrecos()), SpecRel: "docs/specs/x_SPEC_Y.md"})
	semGo(t, "o INTENT", intent, proibidosTexto...)
	semGo(t, "a ROUTE", string(pr.Route), proibidosAC...) // a chave tdd_enforced da ROUTE é do plangen
	semGo(t, "a SKILL", pr.Skill, proibidosTexto...)
	semGo(t, "a SPEC", pr.Spec, proibidosTexto...)
	if !strings.Contains(intent, "planilha") {
		t.Error("o INTENT deve carregar o espelho")
	}
	if !temGateDeDados(t, pr.Route) {
		t.Errorf("a ROUTE precisa de ao menos um gate que confira o conteúdo dos dados:\n%s", pr.Route)
	}
}

func projetarOK(t *testing.T, p processo.Processo) *processo.Projecoes {
	t.Helper()
	pr, err := processo.Projetar(p, agora)
	if err != nil {
		t.Fatalf("Projetar: %v", err)
	}
	return pr
}

func ptr[T any](v T) *T { return &v }

// temGateDeDados lê os gates dos tickets da ROUTE: algum deles confere registros ou esquema?
func temGateDeDados(t *testing.T, route []byte) bool {
	t.Helper()
	var d struct {
		Nodes []struct {
			Type   string
			Output struct{ Gate string }
		}
	}
	if err := json.Unmarshal(route, &d); err != nil {
		t.Fatal(err)
	}
	achou := false
	for _, n := range d.Nodes {
		if n.Type != "unit_task" {
			continue
		}
		e, err := gate.Compile(n.Output.Gate)
		if err != nil {
			t.Fatalf("gate da ROUTE não compila: %q → %v", n.Output.Gate, err)
		}
		if e.ChecksContent() && (strings.Contains(n.Output.Gate, "line_count") || strings.Contains(n.Output.Gate, "schema_valid")) {
			achou = true
		}
	}
	return achou
}

func TestNaturezaAbertaSoftwareContinuaComTextoDeGo(t *testing.T) {
	for _, nat := range []string{"", "software"} {
		intent := interview.RenderIntent(interview.IntentView{Theme: "x", Date: "2026-11-04", SessionID: "s", Confidence: 98,
			Natureza: nat, Mirror: ptr(espelhoPrecos()), SpecRel: "docs/specs/x_SPEC_Y.md"})
		for _, want := range []string{"go test", "go build", "TDD Red-Green"} {
			if !strings.Contains(intent, want) {
				t.Errorf("natureza %q: o INTENT de software deve manter %q", nat, want)
			}
		}
	}
}

func TestNaturezaAbertaSkillDeSoftwareMantemZeroStubs(t *testing.T) {
	p := processoDeSoftware()
	pr := projetarOK(t, p)
	if !strings.Contains(pr.Skill, "Zero Stubs") {
		t.Error("a SKILL de software mantém a regra Zero Stubs")
	}
}

func processoDeSoftware() processo.Processo {
	p := processo.Processo{Nome: "legoloop-api", Versao: "1.0.0", Natureza: "software", Objetivo: "construir uma API", CriadoEm: agora}
	e := processo.Etapa{ID: "e01", Titulo: "Construir", Tipo: "construir", Peca: "aidd-tdd",
		Produz: processo.Artefato{Caminho: "out/01_construir.md", Tipo: "markdown", Promete: true}}
	e.Gate, _ = processo.GateParaArtefato(e.Produz)
	p.Etapas = []processo.Etapa{e}
	return p
}

// pecasDeDados são as peças boas, mas declaradas para dados: o processo de dados tem peça em cada etapa.
func pecasDeDados() []matriz.Entrada {
	es := goodPieces()
	for i := range es {
		es[i].Perfil.Naturezas = []string{"dados"}
	}
	return es
}

// A natureza que o motor conhece chega ao INTENT que ele grava.
func TestNaturezaAbertaMotorGravaIntentSemGoParaDados(t *testing.T) {
	clar := fabricar.Clarification{Theme: objetivoPrecos, Mirror: espelhoPrecos(), Confidence: 98, Reached: true, Interviewed: true, Natureza: "dados", Turns: 2}
	r := newRig(t, pecasDeDados())
	r.deps.Interview = &clar
	intent := read(t, r.run(t).Artifacts.IntentPath)
	semGo(t, "o INTENT gravado", intent, proibidosTexto...)
	if !strings.Contains(intent, "planilha") {
		t.Error("o INTENT gravado deve carregar o espelho")
	}
	clar.Natureza = ""
	r2 := newRig(t, goodPieces())
	r2.deps.Interview = &clar
	if got := read(t, r2.run(t).Artifacts.IntentPath); !strings.Contains(got, "go test") {
		t.Error("sem natureza conhecida o INTENT mantém o texto de software, como antes")
	}
}

func TestNaturezaAbertaInferenciaDeCasosPelasPalavras(t *testing.T) {
	n, _, _ := fabricar.Necessidades(bg(), nil, "fazer a coleta de dados e uma análise de dados", interview.Mirror{}, matriz.DefaultVocab())
	for _, quero := range []string{"coleta-de-dados", "analise-de-dados"} {
		if !slices.Contains(n.Casos, quero) {
			t.Errorf("os casos de uso vêm das palavras do pedido: falta %q em %v", quero, n.Casos)
		}
	}
	if slices.Contains(n.Casos, "tdd") {
		t.Errorf("caso sem palavra no pedido não entra: %v", n.Casos)
	}
}

func TestNaturezaAbertaRestricoesEForaDeEscopoTambemDaoPista(t *testing.T) {
	n, _, _ := fabricar.Necessidades(bg(), nil, "fazer o trabalho", interview.Mirror{Constraints: []string{"entregar como csv"}}, matriz.DefaultVocab())
	if n.Natureza != "dados" {
		t.Errorf("a restrição do espelho também conta: %q", n.Natureza)
	}
	n, _, _ = fabricar.Necessidades(bg(), nil, "fazer o trabalho", interview.Mirror{OutOfScope: []string{"montar o dataset"}}, matriz.DefaultVocab())
	if n.Natureza != "dados" {
		t.Errorf("o fora de escopo do espelho também conta: %q", n.Natureza)
	}
	n, _, _ = fabricar.Necessidades(bg(), nil, "fazer o trabalho", interview.Mirror{WhyNow: "precisamos do manual", User: "equipe"}, matriz.DefaultVocab())
	if n.Natureza != "documentacao" {
		t.Errorf("por que agora também conta: %q", n.Natureza)
	}
	n, _, _ = fabricar.Necessidades(bg(), nil, "fazer o trabalho", interview.Mirror{User: "analistas de planilha"}, matriz.DefaultVocab())
	if n.Natureza != "dados" {
		t.Errorf("quem se beneficia também conta: %q", n.Natureza)
	}
}

func TestNaturezaAbertaProcessoDeNaturezaVaziaViraQualquer(t *testing.T) {
	m := fabricar.Merged{Steps: []fabricar.Step{{Piece: "aidd-spec"}}}
	p, err := fabricar.ProcessoDe("legoloop-x", "fazer x", matriz.Needs{}, m, agora)
	if err != nil || p.Natureza != "qualquer" {
		t.Errorf("sem natureza o processo vale qualquer: %q %v", p.Natureza, err)
	}
	if err := p.Validate(); err != nil {
		t.Errorf("o processo precisa de natureza para validar: %v", err)
	}
}

func TestNaturezaAbertaProcessoDeUsaSoAPecaQueEntrou(t *testing.T) {
	m := fabricar.Merged{Decisoes: []fabricar.Decisao{
		{Etapa: "coletar", Peca: "recusada", Entrou: false},
		{Etapa: "coletar", Peca: "escolhida", Entrou: true},
	}}
	p, err := fabricar.ProcessoDe("legoloop-x", "coletar", matriz.Needs{Natureza: "dados", Etapas: []string{"coletar"}}, m, agora)
	if err != nil || len(p.Etapas) != 1 || p.Etapas[0].Peca != "escolhida" {
		t.Errorf("só a decisão que entrou vale: %+v %v", p.Etapas, err)
	}
}
