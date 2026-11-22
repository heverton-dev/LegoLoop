package interview_test

import (
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/interview"
)

const ideiaLonga = "Iniciativa 4 do macroprojeto: corrigir a espinha dorsal. DECISÕES DO USUÁRIO: escopo das 14 lacunas, registro da espinha, descoberta por hash, vocabulário amarrado, CI nos três SOs, mutação em duas camadas e outras decisões longas que não cabem num título de documento."

func intentDeTeste(theme string, mirror *interview.Mirror) string {
	return interview.RenderIntent(interview.IntentView{
		Theme: theme, Date: "2026-11-21", SessionID: "intv_1", Confidence: 96, Mirror: mirror,
		SpecRel: "../specs/2026-11-21_SPEC_X.md",
	})
}

func espelhoDeTeste() *interview.Mirror {
	return &interview.Mirror{
		Outcome: "dois Processos", User: "o operador", WhyNow: "a trava existe",
		Success: "prova um; prova dois; prova três", Constraints: []string{"sem LLM"}, OutOfScope: []string{"formato do DAG"},
	}
}

// Achado da I04: o título da INTENT repetia a mensagem inteira que o usuário enviou ao iniciar a entrevista.
func TestTituloDaIntentNaoRepeteAIdeiaLonga(t *testing.T) {
	md := intentDeTeste(ideiaLonga, espelhoDeTeste())
	for _, l := range strings.Split(md, "\n") {
		if strings.HasPrefix(l, "# 🎯 INTENT:") && len([]rune(l)) > 130 {
			t.Fatalf("título longo demais (%d caracteres): %q", len([]rune(l)), l)
		}
	}
}

func TestIdeiaCompletaFicaSoEmOProblemaReal(t *testing.T) {
	md := intentDeTeste(ideiaLonga, espelhoDeTeste())
	if n := strings.Count(md, "outras decisões longas que não cabem num título"); n != 1 {
		t.Fatalf("a ideia completa deveria aparecer uma vez (em O Problema Real), apareceu %d vezes", n)
	}
	if !strings.Contains(md, "* **O Problema Real:** "+ideiaLonga) {
		t.Fatal("O Problema Real deve guardar a ideia completa")
	}
}

func TestTituloCurtoFicaIntacto(t *testing.T) {
	md := intentDeTeste("Extrair áudio de vídeo", espelhoDeTeste())
	if !strings.Contains(md, "# 🎯 INTENT: Extrair áudio de vídeo\n") {
		t.Fatal("título curto não pode ser alterado")
	}
}

// Os critérios de sucesso viram um subitem por critério, e não um texto único separado por ponto e vírgula.
func TestCriteriosDeSucessoEmSubitens(t *testing.T) {
	md := intentDeTeste("tema", espelhoDeTeste())
	if strings.Contains(md, "prova um; prova dois") {
		t.Fatal("critérios ainda estão colados por ponto e vírgula")
	}
	for _, c := range []string{"  * prova um\n", "  * prova dois\n", "  * prova três\n"} {
		if !strings.Contains(md, c) {
			t.Fatalf("falta o subitem %q", c)
		}
	}
}
