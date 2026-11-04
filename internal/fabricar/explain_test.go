package fabricar_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/matriz"
	"legoloop/internal/ste"
)

// cerebro é um cérebro de teste: devolve a resposta dada e guarda os prompts.
type cerebro struct {
	reply   string
	err     error
	prompts []string
}

func (c *cerebro) Chat(_ context.Context, msgs []interview.Message) (string, error) {
	for _, m := range msgs {
		c.prompts = append(c.prompts, m.Content)
	}
	return c.reply, c.err
}

func (c *cerebro) visto() string { return strings.Join(c.prompts, "\n") }

// fusao monta uma Fusão com duas etapas cobertas por peças diferentes e uma lacuna.
func fusao() (fabricar.Merged, matriz.Needs) {
	n := matriz.Needs{Natureza: "software", Etapas: []string{"construir", "verificar", "publicar"}}
	es := []matriz.Entrada{
		skill("tdd", 0.58, []string{"construir"}, sw),
		skill("revisor", 0.90, []string{"verificar"}, sw),
	}
	return fabricar.MergeByFit(n, es, fabricar.MergeOptions{}), n
}

const explicacoes = `{"explicacoes":{"d1":"O tdd conduz a construção por testes.","d2":"O revisor confere o resultado.","d3":"Nenhuma peça publica software."},"ordem":["tdd","revisor"]}`

func TestExplicaContestaCadaDecisaoGanhaExplicacaoEOMotivoFica(t *testing.T) {
	m, n := fusao()
	antes := slices.Clone(m.Decisoes)
	if len(antes) != 3 {
		t.Fatalf("a Fusão de teste deve ter 3 decisões: %+v", antes)
	}
	m.Explain(context.Background(), &cerebro{reply: explicacoes}, n)
	if m.Decisoes[0].Explicacao != "O tdd conduz a construção por testes." || m.Decisoes[2].Explicacao != "Nenhuma peça publica software." {
		t.Errorf("explicações não aplicadas: %+v", m.Decisoes)
	}
	for i := range antes {
		if m.Decisoes[i].Motivo != antes[i].Motivo || m.Decisoes[i].Peca != antes[i].Peca || m.Decisoes[i].Entrou != antes[i].Entrou {
			t.Errorf("o cérebro não pode alterar a decisão: %+v → %+v", antes[i], m.Decisoes[i])
		}
	}
	if len(m.Divergencias) != 0 {
		t.Errorf("mesma ordem não diverge: %v", m.Divergencias)
	}
}

func TestExplicaContestaOrdemDiferenteNaoMudaORankingEViraDivergencia(t *testing.T) {
	m, n := fusao()
	passos, usadas := nomes(m), slices.Clone(m.Used)
	m.Explain(context.Background(), &cerebro{reply: `{"explicacoes":{},"ordem":["revisor","tdd"]}`}, n)
	if got := nomes(m); !slices.Equal(got, passos) {
		t.Errorf("o ranking calculado vale: %v → %v", passos, got)
	}
	if !slices.Equal(m.Used, usadas) {
		t.Error("Used não pode mudar")
	}
	if len(m.Divergencias) != 1 {
		t.Fatalf("esperava 1 divergência: %v", m.Divergencias)
	}
	d := m.Divergencias[0]
	for _, want := range []string{"revisor, tdd", "tdd, revisor", "ranking calculado vale"} {
		if !strings.Contains(d, want) {
			t.Errorf("a divergência não traz %q: %s", want, d)
		}
	}
	if v := ste.Lint(d); len(v) != 0 {
		t.Errorf("fora do perfil STE-PT: %v", v)
	}
}

func TestExplicaContestaPecaForaDoRankingTambemDiverge(t *testing.T) {
	m, n := fusao()
	m.Explain(context.Background(), &cerebro{reply: `{"ordem":["tdd","outra-peca"]}`}, n)
	if len(m.Divergencias) != 1 || !strings.Contains(m.Divergencias[0], "outra-peca") {
		t.Errorf("sugerir peça que a Fusão não escolheu é divergência: %v", m.Divergencias)
	}
	if got := nomes(m); !slices.Equal(got, []string{"tdd", "revisor"}) {
		t.Errorf("Steps mudaram: %v", got)
	}
}

func TestExplicaContestaSemOrdemNaoDiverge(t *testing.T) {
	m, n := fusao()
	m.Explain(context.Background(), &cerebro{reply: `{"explicacoes":{"d1":"ok"}}`}, n)
	if len(m.Divergencias) != 0 || m.Decisoes[0].Explicacao != "ok" {
		t.Errorf("sem ordem proposta não há divergência: %v %+v", m.Divergencias, m.Decisoes[0])
	}
}

func TestExplicaContestaFalhaDoCerebroNaoDerrubaAFusao(t *testing.T) {
	for nome, c := range map[string]*cerebro{
		"erro":                     {err: errors.New("fora do ar")},
		"prosa":                    {reply: "Acho que está tudo certo, sem JSON."},
		"json quebrado":            {reply: `{"explicacoes": {`},
		"vazio":                    {reply: ""},
		"erro com resposta válida": {reply: explicacoes, err: errors.New("cortou no meio")},
	} {
		t.Run(nome, func(t *testing.T) {
			m, n := fusao()
			passos := nomes(m)
			m.Explain(context.Background(), c, n)
			if got := nomes(m); !slices.Equal(got, passos) {
				t.Errorf("Steps mudaram: %v", got)
			}
			for _, d := range m.Decisoes {
				if d.Explicacao != "" {
					t.Errorf("sem resposta válida não há explicação: %+v", d)
				}
			}
			if !slices.ContainsFunc(m.Warnings, func(w string) bool { return strings.Contains(w, "não explicou as decisões") }) {
				t.Errorf("o aviso de falha do cérebro falta: %v", m.Warnings)
			}
			for _, w := range m.Warnings {
				if strings.Contains(w, "não explicou") && len(ste.Lint(w)) != 0 {
					t.Errorf("aviso fora do STE-PT: %q", w)
				}
			}
		})
	}
}

func TestExplicaContestaAceitaJSONCercadoDeTexto(t *testing.T) {
	m, n := fusao()
	reply := "Claro! Segue:\n```json\n" + explicacoes + "\n```\nQualquer dúvida, avise."
	m.Explain(context.Background(), &cerebro{reply: reply}, n)
	if m.Decisoes[1].Explicacao != "O revisor confere o resultado." {
		t.Errorf("JSON dentro de cerca de código deve valer: %+v", m.Decisoes[1])
	}
}

func TestExplicaContestaTextoDoCerebroEDadoLimpoELimitado(t *testing.T) {
	m, n := fusao()
	longo := strings.Repeat("ç", 500)
	reply := `{"explicacoes":{"d1":"linha 1\nlinha 2\t  fim","d2":"` + longo + `","dZZ":"id que não existe","d3":42}}`
	m.Explain(context.Background(), &cerebro{reply: reply}, n)
	if got := m.Decisoes[0].Explicacao; got != "linha 1 linha 2 fim" {
		t.Errorf("quebras e espaços devem virar um espaço só: %q", got)
	}
	if got := []rune(m.Decisoes[1].Explicacao); len(got) != 300 {
		t.Errorf("explicação passa do limite: %d caracteres", len(got))
	}
	if m.Decisoes[2].Explicacao != "" {
		t.Errorf("valor que não é texto é ignorado: %q", m.Decisoes[2].Explicacao)
	}
}

func TestExplicaContestaPromptPedeSoExplicacaoETrataTextoComoDado(t *testing.T) {
	m, n := fusao()
	c := &cerebro{reply: explicacoes}
	m.Explain(context.Background(), c, n)
	p := c.visto()
	for _, want := range []string{"tdd", "revisor", "construir", "verificar", "publicar", "software", "d1", "d3", "NÃO escolhe", "dados"} {
		if !strings.Contains(p, want) {
			t.Errorf("o prompt não traz %q", want)
		}
	}
	for _, d := range m.Decisoes {
		if d.Motivo != "" && !strings.Contains(p, d.Motivo) {
			t.Errorf("o prompt deve trazer o motivo calculado: %q", d.Motivo)
		}
	}
}

func TestExplicaContestaSemDecisoesNaoChamaOCerebro(t *testing.T) {
	c := &cerebro{reply: explicacoes}
	var m fabricar.Merged
	m.Explain(context.Background(), c, matriz.Needs{})
	if len(c.prompts) != 0 {
		t.Error("sem decisões não há o que explicar")
	}
}

func TestExplicaContestaCancelamentoNaoDerruba(t *testing.T) {
	m, n := fusao()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	m.Explain(ctx, &cerebro{err: ctx.Err()}, n)
	if len(m.Steps) == 0 {
		t.Error("a Fusão deve continuar íntegra")
	}
}

func TestExplicaContestaMotivoIndependeDoCerebro(t *testing.T) {
	a, n := fusao()
	b, _ := fusao()
	a.Explain(context.Background(), &cerebro{reply: explicacoes}, n)
	b.Explain(context.Background(), &cerebro{reply: `{"explicacoes":{"d1":"outra"},"ordem":["revisor"]}`}, n)
	for i := range a.Decisoes {
		if a.Decisoes[i].Motivo != b.Decisoes[i].Motivo || a.Decisoes[i].Peca != b.Decisoes[i].Peca {
			t.Errorf("cérebros diferentes não podem mudar o Motivo nem a peça: %+v vs %+v", a.Decisoes[i], b.Decisoes[i])
		}
	}
}

func TestExplicaContestaAvisaQuandoExplicaSoParte(t *testing.T) {
	m, n := fusao()
	m.Explain(context.Background(), &cerebro{reply: `{"explicacoes":{"d1":"ok"}}`}, n)
	if !slices.ContainsFunc(m.Warnings, func(w string) bool { return strings.Contains(w, "não explicou 2 decisões") }) {
		t.Errorf("faltou o aviso das 2 decisões sem explicação: %v", m.Warnings)
	}
	m2, n2 := fusao()
	m2.Explain(context.Background(), &cerebro{reply: explicacoes}, n2)
	if slices.ContainsFunc(m2.Warnings, func(w string) bool { return strings.Contains(w, "não explicou") }) {
		t.Errorf("tudo explicado, sem aviso: %v", m2.Warnings)
	}
}

func TestExplicaContestaItensEmBrancoNaOrdemSaoIgnorados(t *testing.T) {
	m, n := fusao()
	m.Explain(context.Background(), &cerebro{reply: `{"ordem":["tdd","","   ",7,"revisor"]}`}, n)
	if len(m.Divergencias) != 0 {
		t.Errorf("itens em branco ou que não são texto não contam: %v", m.Divergencias)
	}
}
