package fabricar_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"legoloop/internal/interview"
	"legoloop/internal/taxonomy"
)

// marcasDosPrompts identificam os dois prompts do sistema que o fabricar escreve.
var marcasDosPrompts = []string{"vocabulário fechado do LegoLoop", "explicador da Fusão"}

// ecoLLM guarda o prompt do sistema de cada chamada e responde como o routerLLM.
type ecoLLM struct {
	inner   *routerLLM
	prompts map[string]string
}

func (e *ecoLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	if len(msgs) > 0 && msgs[0].Role == "system" {
		for _, marca := range marcasDosPrompts {
			if strings.Contains(msgs[0].Content, marca) {
				e.prompts[marca] = msgs[0].Content
			}
		}
	}
	return e.inner.Chat(ctx, msgs)
}

var (
	reMarca     = regexp.MustCompile(`(?i)legoloop`)
	reCamel     = regexp.MustCompile(`(\p{Ll}|\p{N})(\p{Lu})`)
	reSeparador = regexp.MustCompile(`[^\p{L}\p{N}\n]+`)
	rePalavra   = regexp.MustCompile(`\p{L}+`)
)

// separaAfixos separa o prefixo "sub" e o plural "s" quando o resto é termo proibido, como o taxlinter.
func separaAfixos(w string) string {
	lw := strings.ToLower(w)
	if len(lw) > 3 && strings.HasPrefix(lw, "sub") {
		if _, ok := taxonomy.CheckForbidden(lw[3:]); ok {
			return w[:3] + " " + w[3:]
		}
	}
	if len(lw) > 2 && strings.HasSuffix(lw, "s") {
		if _, ok := taxonomy.CheckForbidden(lw[:len(lw)-1]); ok {
			return w[:len(w)-1] + " s"
		}
	}
	return w
}

// termosProibidos normaliza o texto como o taxlinter (marca fora, identificadores separados) e devolve os termos proibidos.
func termosProibidos(t *testing.T, texto string) []string {
	t.Helper()
	limpo := reMarca.ReplaceAllString(texto, " ")
	limpo = reCamel.ReplaceAllString(limpo, "$1 $2")
	limpo = reSeparador.ReplaceAllString(limpo, " ")
	limpo = rePalavra.ReplaceAllStringFunc(limpo, separaAfixos)
	vs, err := taxonomy.ValidateText(limpo)
	if err != nil {
		t.Fatalf("dicionário da taxonomia: %v", err)
	}
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		out = append(out, v.ForbiddenTerm)
	}
	return out
}

// geraComEco roda o Fabricar real com o LLM falso e devolve os prompts do fabricar e o Processo gerado.
func geraComEco(t *testing.T) (map[string]string, string) {
	t.Helper()
	r := newRig(t, goodPieces())
	eco := &ecoLLM{inner: r.llm, prompts: map[string]string{}}
	r.deps.LLM = eco
	res := r.run(t)
	for _, marca := range marcasDosPrompts {
		if eco.prompts[marca] == "" {
			t.Fatalf("o fabricar não enviou o prompt %q ao cérebro", marca)
		}
	}
	return eco.prompts, read(t, res.Artifacts.ProcessoPath)
}

func TestSaidaLLMSemTermoProibido(t *testing.T) {
	prompts, processo := geraComEco(t)
	for marca, prompt := range prompts {
		if ts := termosProibidos(t, prompt); len(ts) > 0 {
			t.Errorf("o prompt %q tem termo proibido: %v", marca, ts)
		}
	}
	if ts := termosProibidos(t, processo); len(ts) > 0 {
		t.Errorf("o Processo gerado tem termo proibido: %v\n%s", ts, processo)
	}
	t.Run("mutacao", func(t *testing.T) {
		for marca, prompt := range prompts {
			for _, termo := range []string{"fluxo", "subloop", "LoopID", "tarefas"} {
				if len(termosProibidos(t, prompt+"\nUse o "+termo+".")) == 0 {
					t.Errorf("reintroduzir %q no prompt %q deve ser detectado", termo, marca)
				}
			}
		}
	})
}
