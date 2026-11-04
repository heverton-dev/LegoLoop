package ste

import (
	"strings"
	"testing"
)

// frase gera uma frase com exatamente n palavras.
func frase(n int) string {
	return strings.TrimSpace(strings.Repeat("Grave ", n)) + "."
}

func regras(vs []Violation) []string {
	var out []string
	for _, v := range vs {
		out = append(out, v.Regra)
	}
	return out
}

func temRegra(vs []Violation, regra string) bool {
	for _, v := range vs {
		if v.Regra == regra {
			return true
		}
	}
	return false
}

func TestSTE_LimiteDePalavras(t *testing.T) {
	cases := []struct {
		nome    string
		msg     string
		violado bool
	}{
		{"20 palavras passam", frase(20), false},
		{"21 palavras reprovam", frase(21), true},
		{"duas frases curtas passam", frase(15) + " " + frase(15), false},
		{"id com ponto não divide a frase", "Veja AC-E01.1 e AC-E01.2 agora.", false},
		{"hífen solto não conta como palavra", "Grave - agora.", false},
	}
	for _, c := range cases {
		t.Run(c.nome, func(t *testing.T) {
			if got := temRegra(Lint(c.msg), "frase-longa"); got != c.violado {
				t.Fatalf("frase-longa = %v, quer %v (%v)", got, c.violado, regras(Lint(c.msg)))
			}
		})
	}
}

func TestSTE_PalavrasAmbiguas(t *testing.T) {
	for _, w := range []string{"talvez", "Geralmente", "pode ser", "às vezes", "provavelmente"} {
		msg := "O sistema " + w + " grava o arquivo."
		if !temRegra(Lint(msg), "ambigua") {
			t.Errorf("%q deveria violar 'ambigua'", w)
		}
	}
	if temRegra(Lint("O sistema grava o arquivo."), "ambigua") {
		t.Error("frase sem ambiguidade foi reprovada")
	}
	// "talvezinho" não é a palavra proibida: o casamento é por palavra inteira.
	if temRegra(Lint("O talvezinho grava."), "ambigua") {
		t.Error("casou parte de palavra")
	}
}

func TestSTE_TermoForaDoGlossario(t *testing.T) {
	vs := Lint("Grave a esteira no cofre.")
	if !temRegra(vs, "termo") {
		t.Fatalf("'esteira' deveria violar 'termo': %v", vs)
	}
	if !strings.Contains(vs[0].Msg, "processo") {
		t.Errorf("a violação deve citar o termo preferido: %q", vs[0].Msg)
	}
	if temRegra(Lint("Grave o processo no cofre."), "termo") {
		t.Error("o termo preferido foi reprovado")
	}
}

func cartaoExemplo() Cartao {
	return Cartao{
		Acao:        `gravar o processo "legoloop-build-engine".`,
		Passos:      []string{"Grave 8 arquivos no cofre e no projeto.", "Ligue a skill a 9 harnesses."},
		Efeitos:     []string{"8 arquivos novos.", "9 links novos."},
		SeDerErrado: "O comando desfaz tudo.",
	}
}

func TestSTE_CartaoDeExemploPassa(t *testing.T) {
	if vs := LintCartao(cartaoExemplo()); len(vs) != 0 {
		t.Fatalf("o cartão de exemplo da SPEC deve passar: %v", vs)
	}
}

func TestSTE_AcaoSemAvisoDeRisco(t *testing.T) {
	c := cartaoExemplo()
	c.Irreversivel = true
	if !temRegra(LintCartao(c), "sem-aviso") {
		t.Fatal("ação irreversível sem aviso deve violar 'sem-aviso'")
	}
	c.Aviso = "Esta ação remove o processo. Não há volta."
	if vs := LintCartao(c); len(vs) != 0 {
		t.Fatalf("com aviso o cartão deve passar: %v", vs)
	}
}

func TestSTE_CartaoSemComoDesfazer(t *testing.T) {
	c := cartaoExemplo()
	c.SeDerErrado = ""
	if !temRegra(LintCartao(c), "sem-desfazer") {
		t.Fatal("cartão sem 'se der errado' deve violar 'sem-desfazer'")
	}
}

func TestSTE_CartaoLintaOsTextos(t *testing.T) {
	c := cartaoExemplo()
	c.Passos = append(c.Passos, "Talvez grave mais um arquivo.")
	if !temRegra(LintCartao(c), "ambigua") {
		t.Fatal("o texto dos passos também é verificado")
	}
}

func TestSTE_RenderPoeOAvisoAntesDosPassos(t *testing.T) {
	c := cartaoExemplo()
	c.Irreversivel = true
	c.Aviso = "Esta ação remove o processo."
	out := c.Render()
	a, p := strings.Index(out, "AVISO"), strings.Index(out, "O QUE VAI ACONTECER")
	if a < 0 || p < 0 || a > p {
		t.Fatalf("o aviso deve vir antes dos passos (aviso=%d passos=%d):\n%s", a, p, out)
	}
	for _, quer := range []string{"[1] seguir", "[2] personalizar", "[0] cancelar", "SE DER ERRADO"} {
		if !strings.Contains(out, quer) {
			t.Errorf("faltou %q:\n%s", quer, out)
		}
	}
}

func TestSTE_ParseDecisao(t *testing.T) {
	cases := map[string]Decisao{"1": Seguir, "s": Seguir, "sim": Seguir, "2": Personalizar, "p": Personalizar, "0": Cancelar, "n": Cancelar, "": Cancelar}
	for in, quer := range cases {
		got, ok := ParseDecisao(in)
		if !ok || got != quer {
			t.Errorf("ParseDecisao(%q) = %v, %v; quer %v", in, got, ok, quer)
		}
	}
	if _, ok := ParseDecisao("talvez"); ok {
		t.Error("entrada desconhecida não pode virar decisão")
	}
}

func TestSTE_CatalogoReprovaMensagemRuim(t *testing.T) {
	c := NewCatalog()
	c.Register("ok", "Grave o processo.")
	c.Register("ruim", "Talvez grave a esteira agora.")
	vs := c.Violations()
	if len(vs) == 0 || !strings.Contains(vs[0].Msg, "ruim") {
		t.Fatalf("o catálogo deve apontar o id da mensagem ruim: %v", vs)
	}
	if got := c.Msg("ok"); got != "Grave o processo." {
		t.Errorf("Msg(ok) = %q", got)
	}
	if got := c.Msg("naoexiste"); !strings.Contains(got, "naoexiste") {
		t.Errorf("id desconhecido deve aparecer na saída, veio %q", got)
	}
}

func TestSTE_CatalogoDoProdutoNaoTemViolacao(t *testing.T) {
	if vs := Default.Violations(); len(vs) != 0 {
		t.Fatalf("o catálogo do produto tem violações: %v", vs)
	}
}
