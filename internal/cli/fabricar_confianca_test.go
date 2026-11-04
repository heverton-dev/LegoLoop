package cli

import (
	"strings"
	"testing"
)

// answers repete a resposta "1" (aceitar o palpite) n vezes.
func answers(n int) string { return strings.Repeat("1\n", n) }

// A entrevista que não chega a 95% em 12 turnos pergunta se deve gerar a SPEC assim mesmo.
func TestConfiancaHonestaCLIPerguntaAntesDaSpec(t *testing.T) {
	t.Run("confirmar segue e grava o INTENT em rascunho", func(t *testing.T) {
		r := newInterRig(t, "\n"+answers(12)+"s\ns\n", lowReply)
		if code := r.fabricar(t); code != 0 {
			t.Fatalf("exit %d\n%s\n%s", code, r.out, r.errb)
		}
		for _, want := range []string{"Confiança: 50%", "A entrevista parou em 50%", "Lacuna: quem se beneficia", "SPEC mesmo assim"} {
			if !strings.Contains(r.out.String()+r.errb.String(), want) {
				t.Errorf("saída sem %q:\n%s\n%s", want, r.out, r.errb)
			}
		}
		intent := firstGlob(t, r.doc+"/docs/intents/*.md")
		if !strings.Contains(intent, "**Status:** DRAFT") || !strings.Contains(intent, "50%") {
			t.Errorf("o INTENT deve sair DRAFT com a confiança medida:\n%s", intent)
		}
	})

	t.Run("recusar não gera SPEC e não grava", func(t *testing.T) {
		r := newInterRig(t, "\n"+answers(12)+"n\n", lowReply)
		if code := r.fabricar(t); code != 1 {
			t.Fatalf("exit %d, esperado 1", code)
		}
		if r.wrote(t) != 0 {
			t.Errorf("gravou %d arquivos", r.wrote(t))
		}
	})

	t.Run("acima da meta não pergunta", func(t *testing.T) {
		r := newInterRig(t, "\n1\ns\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 0 {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
		if strings.Contains(r.out.String(), "SPEC mesmo assim") {
			t.Error("com 98% não há o que confirmar")
		}
	})
}
