package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/interview"
	"legoloop/internal/provider"
)

const (
	lowReply  = "HIPÓTESE: algo\nCONFIANÇA: 50%\nPERGUNTA: Qual o formato?\nPALPITE: mp3"
	doneReply = "HIPÓTESE: extrator\nCONFIANÇA: 98%\nESPELHO CANÔNICO:\n- Resultado Esperado: extrair áudio em mp3\n- Quem se Beneficia: editores\n- Por que Agora: retrabalho\n- Critério de Sucesso: arquivo gerado\n- Restrições Inegociáveis: sem rede\n- Fora de Escopo: edição de vídeo"
)

// chatLLM responde a SPEC quando o prompt é do especificador e a entrevista com as respostas roteirizadas.
type chatLLM struct {
	replies   []string
	interview int
	user      []string // tudo que o usuário disse, em ordem
}

func (c *chatLLM) Chat(_ context.Context, msgs []interview.Message) (string, error) {
	switch {
	case strings.Contains(msgs[0].Content, "Arquiteto de Especificações"):
		return specFake, nil
	case strings.Contains(msgs[0].Content, "vocabulário fechado do LegoLoop"):
		return `{"natureza": "software", "etapas": ["construir", "verificar"], "casos": []}`, nil
	case strings.Contains(msgs[0].Content, "explicador da Fusão"):
		return `{"explicacoes": {}}`, nil
	}
	c.user = append(c.user, msgs[len(msgs)-1].Content)
	i := min(c.interview, len(c.replies)-1)
	c.interview++
	return c.replies[i], nil
}

type interRig struct {
	*fabRig
	llm *chatLLM
}

// newInterRig prepara o terminal roteirizado: stdin é o que o usuário digitaria, em ordem.
func newInterRig(t *testing.T, stdin string, replies ...string) *interRig {
	t.Helper()
	r := &interRig{fabRig: newFabRig(t), llm: &chatLLM{replies: replies}}
	r.d.Stdin = strings.NewReader(stdin)
	r.d.Discover = func(context.Context) provider.Inventory {
		return provider.Inventory{CLIs: []provider.CLI{{Name: "claude"}}}
	}
	r.d.Fabricar = &FabricarDeps{LLM: r.llm, Matriz: matrizCLI(), Clock: func() time.Time { return fabCLINow }}
	return r
}

func (r *interRig) fabricar(t *testing.T, extra ...string) int {
	t.Helper()
	return run(t, r.d, r.args("novo fluxo", extra...)...)
}

func (r *interRig) wrote(t *testing.T) int { return len(r.files(t)) }

func TestRunFabricarInteractive(t *testing.T) {
	t.Run("caminho feliz: palpite aceito e confirmado", func(t *testing.T) {
		r := newInterRig(t, "1\ns\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 0 {
			t.Fatalf("exit %d\n%s\n%s", code, r.out, r.errb)
		}
		for _, want := range []string{"[1] Aceitar o palpite: mp3", "[1] seguir", "SKILL.md", "legoloop-novo-fluxo"} {
			if !strings.Contains(r.out.String(), want) {
				t.Errorf("saída sem %q:\n%s", want, r.out)
			}
		}
		if r.wrote(t) != 9 {
			t.Errorf("devia gravar 9 arquivos, gravou %d", r.wrote(t))
		}
		if strings.Contains(r.out.String(), "Drives detectados") {
			t.Errorf("o fabricar não pergunta drives: as peças vêm da matriz:\n%s", r.out)
		}
		intent := firstGlob(t, filepath.Join(r.doc, "docs", "intents", "*.md"))
		if !strings.Contains(intent, "extrair áudio em mp3") {
			t.Errorf("a INTENT deve nascer do espelho da entrevista:\n%s", intent)
		}
	})

	t.Run("resposta livre vai para a entrevista", func(t *testing.T) {
		r := newInterRig(t, "2\nSó mp3 e wav\ns\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 0 {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
		if !strings.Contains(strings.Join(r.llm.user, "|"), "Só mp3 e wav") {
			t.Errorf("a resposta livre não chegou à entrevista: %v", r.llm.user)
		}
	})

	t.Run("responder N nao grava e sai com 1", func(t *testing.T) {
		r := newInterRig(t, "1\nn\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 1 {
			t.Fatalf("exit %d, esperado 1", code)
		}
		if r.wrote(t) != 0 || !strings.Contains(r.errb.String(), "cancelado") {
			t.Errorf("gravou=%d stderr=%s", r.wrote(t), r.errb)
		}
	})

	t.Run("Enter vazio na confirmacao vale N", func(t *testing.T) {
		r := newInterRig(t, "1\n\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 1 || r.wrote(t) != 0 {
			t.Fatalf("exit %d gravou=%d", code, r.wrote(t))
		}
	})

	t.Run("EOF antes da confirmacao aborta com 1 sem gravar", func(t *testing.T) {
		r := newInterRig(t, "1\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 1 || r.wrote(t) != 0 {
			t.Fatalf("exit %d gravou=%d", code, r.wrote(t))
		}
	})

	t.Run("EOF no meio da entrevista aborta com 1", func(t *testing.T) {
		r := newInterRig(t, "", lowReply, doneReply)
		if code := r.fabricar(t); code != 1 || r.wrote(t) != 0 {
			t.Fatalf("exit %d gravou=%d", code, r.wrote(t))
		}
	})

	t.Run("opcao 0 cancela a entrevista", func(t *testing.T) {
		r := newInterRig(t, "0\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 1 || r.wrote(t) != 0 {
			t.Fatalf("exit %d gravou=%d", code, r.wrote(t))
		}
	})

	t.Run("tres entradas invalidas seguidas abortam", func(t *testing.T) {
		r := newInterRig(t, "9\nabc\nzzz\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 1 || r.wrote(t) != 0 {
			t.Fatalf("exit %d gravou=%d", code, r.wrote(t))
		}
		if n := strings.Count(r.out.String(), "Qual o formato?"); n != 3 {
			t.Errorf("a pergunta devia ser reapresentada até 3 vezes, apareceu %d", n)
		}
		if !strings.Contains(r.errb.String(), "inválida") {
			t.Errorf("o erro deve explicar:\n%s", r.errb)
		}
	})

	t.Run("duas invalidas e depois valida segue em frente", func(t *testing.T) {
		r := newInterRig(t, "9\nabc\n1\ns\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 0 {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
	})

	t.Run("12 turnos sem chegar a 95 por cento pergunta sobre a SPEC e depois confirma a gravacao", func(t *testing.T) {
		in := strings.Repeat("1\n", 12) + "s\ns\n"
		r := newInterRig(t, in, lowReply)
		if code := r.fabricar(t); code != 0 {
			t.Fatalf("exit %d: %s\n%s", code, r.errb, r.out)
		}
		if r.llm.interview != 13 {
			t.Errorf("1 início + 12 turnos = 13 chamadas, houve %d", r.llm.interview)
		}
		if !strings.Contains(r.out.String(), "SPEC mesmo assim") {
			t.Errorf("faltou perguntar se gera a SPEC com confiança baixa:\n%s", r.out)
		}
	})

	t.Run("dry-run interativo nao pergunta a confirmacao nem grava", func(t *testing.T) {
		r := newInterRig(t, "1\n", lowReply, doneReply)
		if code := r.fabricar(t, "--dry-run"); code != 0 || r.wrote(t) != 0 {
			t.Fatalf("exit %d gravou=%d: %s", code, r.wrote(t), r.errb)
		}
	})

	t.Run("sem stdin falha com mensagem clara", func(t *testing.T) {
		r := newInterRig(t, "", lowReply)
		r.d.Stdin = nil
		if code := r.fabricar(t); code != 1 || !strings.Contains(r.errb.String(), "terminal") {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
	})
}

func firstGlob(t *testing.T, pattern string) string {
	t.Helper()
	m, err := filepath.Glob(pattern)
	if err != nil || len(m) == 0 {
		t.Fatalf("nenhum arquivo para %s", pattern)
	}
	b, err := os.ReadFile(m[0])
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
