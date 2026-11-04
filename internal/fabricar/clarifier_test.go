package fabricar_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/ste"
)

const turn70 = "HIPÓTESE: extrator\nCONFIANÇA: 70%\nPERGUNTA: Qual o formato?\nPALPITE: mp3"

// interviewRig prepara uma entrevista interativa que nunca passa de 70%.
func interviewRig(t *testing.T, replies ...string) (*rig, *int) {
	t.Helper()
	r := newRig(t, goodPieces())
	if len(replies) == 0 {
		replies = []string{turn70}
	}
	r.llm.interviewReplies = replies
	r.req.Interactive = true
	asked := 0
	r.deps.Ask = func(string, string) (string, error) { asked++; return "mp3", nil }
	return r, &asked
}

func filesIn(t *testing.T, roots ...string) []string {
	t.Helper()
	var out []string
	for _, root := range roots {
		_ = filepath.Walk(root, func(p string, i os.FileInfo, err error) error {
			if err == nil && !i.IsDir() {
				out = append(out, p)
			}
			return nil
		})
	}
	return out
}

func TestConfiancaHonestaLimiteDeTurnosGeraDraftComConfiancaMedida(t *testing.T) {
	r, asked := interviewRig(t)
	r.deps.ConfirmDraft = func(fabricar.Clarification) (bool, error) { return true, nil }
	res := r.run(t)
	if *asked != 12 {
		t.Fatalf("a entrevista devia usar os 12 turnos, usou %d", *asked)
	}
	intent := read(t, res.Artifacts.IntentPath)
	for _, want := range []string{"**Status:** DRAFT", "**Confiança da Entrevista:** 70%"} {
		if !strings.Contains(intent, want) {
			t.Errorf("INTENT sem %q:\n%s", want, intent)
		}
	}
	if strings.Contains(intent, "**Confiança da Entrevista:** 0%") {
		t.Error("a confiança não pode cair para 0%")
	}
}

func TestConfiancaHonestaRelatorioAvisaELista(t *testing.T) {
	r, _ := interviewRig(t)
	r.deps.ConfirmDraft = func(fabricar.Clarification) (bool, error) { return true, nil }
	res := r.run(t)
	all := strings.Join(res.Warnings, "\n")
	for _, want := range []string{"70%", "95%", "quem se beneficia", "por que agora", "critério de sucesso", "restrições inegociáveis", "fora de escopo"} {
		if !strings.Contains(all, want) {
			t.Errorf("os avisos não trazem %q:\n%s", want, all)
		}
	}
	for _, w := range res.Warnings {
		if v := ste.Lint(w); len(v) != 0 {
			t.Errorf("aviso fora do perfil STE-PT: %q → %v", w, v)
		}
	}
}

func TestConfiancaHonestaSemConfirmacaoNaoGeraSpec(t *testing.T) {
	r, _ := interviewRig(t) // ConfirmDraft ausente
	_, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
	if !errors.Is(err, fabricar.ErrLowConfidence) {
		t.Fatalf("esperava ErrLowConfidence, veio %v", err)
	}
	if r.llm.specCalls != 0 {
		t.Errorf("a SPEC foi pedida %d vez(es) sem confirmação", r.llm.specCalls)
	}
	if got := filesIn(t, r.deps.Roots.Vault, r.deps.Roots.Docs); len(got) != 0 {
		t.Errorf("gravou sem confirmação: %v", got)
	}
}

func TestConfiancaHonestaRecusaCancelaSemSpec(t *testing.T) {
	r, _ := interviewRig(t)
	r.deps.ConfirmDraft = func(fabricar.Clarification) (bool, error) { return false, nil }
	_, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
	if !errors.Is(err, fabricar.ErrAborted) {
		t.Fatalf("esperava ErrAborted, veio %v", err)
	}
	if r.llm.specCalls != 0 || len(filesIn(t, r.deps.Roots.Vault, r.deps.Roots.Docs)) != 0 {
		t.Error("recusar não pode gerar SPEC nem gravar")
	}
}

func TestConfiancaHonestaConfirmacaoRecebeAMedida(t *testing.T) {
	r, _ := interviewRig(t)
	var got fabricar.Clarification
	r.deps.ConfirmDraft = func(c fabricar.Clarification) (bool, error) { got = c; return true, nil }
	r.run(t)
	if got.Confidence != 70 || got.Reached || len(got.Gaps) == 0 {
		t.Errorf("a confirmação deve receber a medida e as lacunas: %+v", got)
	}
}

func TestConfiancaHonestaAcimaDaMetaNaoPerguntaNemAvisa(t *testing.T) {
	r, _ := interviewRig(t, turn70, doneReply)
	r.deps.ConfirmDraft = func(fabricar.Clarification) (bool, error) {
		t.Error("com 95% ou mais não há o que confirmar")
		return false, nil
	}
	res := r.run(t)
	for _, w := range res.Warnings {
		if strings.Contains(w, "meta") || strings.Contains(w, "Lacuna") {
			t.Errorf("aviso de confiança baixa indevido: %q", w)
		}
	}
}

func TestConfiancaHonestaModoAutonomoNaoExigeConfirmacao(t *testing.T) {
	r := newRig(t, goodPieces()) // sem Interactive: nenhuma entrevista aconteceu
	if _, err := fabricar.Fabricar(context.Background(), r.req, r.deps); err != nil {
		t.Fatalf("o modo autônomo não tem entrevista a confirmar: %v", err)
	}
}
