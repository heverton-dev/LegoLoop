package ste

import (
	"strings"
	"testing"
)

// TestLintEvasaoDeProva cobre a lei da prova: fato sem prova se verifica, achado se corrige.
// Frase que declara pendencia ou ausencia de verificacao reprova.
func TestLintEvasaoDeProva(t *testing.T) {
	casos := []string{
		"Não verificado.",
		"Nao verificado: falta rodar o teste.",
		"Isso fica como pendência.",
		"Não consegui rodar o comando.",
		"Ainda falta validar o schema.",
		"A próxima etapa é revisar.",
		"Isso fica para depois.",
		"Falta verificar o gate.",
		"Não foi possível verificar o resultado.",
		"Não testei a mudança.",
	}
	for _, c := range casos {
		if !temRegra(Lint(c), "evasao-de-prova") {
			t.Errorf("frase %q deveria reprovar com evasao-de-prova", c)
		}
	}
}

// TestLintEvasaoLimpo garante que frase com prova nao reprova.
func TestLintEvasaoLimpo(t *testing.T) {
	casos := []string{
		"O teste passou (go test ./internal/ste: ok).",
		"O worker terminou com exit 0 (swarm_run).",
		"A espinha não tem violações (espinha verificar: 0).",
	}
	for _, c := range casos {
		if temRegra(Lint(c), "evasao-de-prova") {
			t.Errorf("frase %q não deveria reprovar", c)
		}
	}
}

// TestPromptProvaProibeEvasao garante que o texto distribuido aos harnesses proibe a evasao
// e nao manda mais escrever "Não verificado".
func TestPromptProvaProibeEvasao(t *testing.T) {
	for _, p := range []Pausa{PausaNunca, PausaSempre, PausaDestrutivo} {
		txt := promptProva(p)
		if strings.Contains(txt, `escreva "Não verificado"`) {
			t.Errorf("prompt ainda manda escrever Não verificado (pausa %v)", p)
		}
		if !strings.Contains(txt, "Proibido") || !strings.Contains(txt, "verifique") {
			t.Errorf("prompt sem a proibicao de evasao e a ordem de verificar (pausa %v)", p)
		}
	}
}

// TestLintEvasaoIgnoraCitacao garante que citar a frase proibida entre aspas ou crases nao reprova:
// explicar a lei exige citar o que ela proibe.
func TestLintEvasaoIgnoraCitacao(t *testing.T) {
	casos := []string{
		`A lei proibe a frase "Não verificado" em qualquer resposta.`,
		"O termo `pendência` está proibido.",
		"A frase “ainda falta validar” é evasão.",
	}
	for _, c := range casos {
		if temRegra(Lint(c), "evasao-de-prova") {
			t.Errorf("citação %q não deveria reprovar", c)
		}
	}
}
