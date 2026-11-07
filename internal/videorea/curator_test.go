package videorea_test

import (
	"strings"
	"testing"

	"legoloop/internal/videorea"
)

func TestCurateAgainstMatrix(t *testing.T) {
	dossier := &videorea.ForensicsDossier{
		Title:           "How To Make Insane Video Editing",
		ToolsIdentified: []string{"Higgsfield", "After Effects", "Whisper", "Claude"},
		HiddenTools:     []string{"FFmpeg CLI"},
	}

	report, err := videorea.CurateAgainstMatrix(dossier)
	if err != nil {
		t.Fatalf("CurateAgainstMatrix falhou: %v", err)
	}

	if report.OriginalCost == "" || report.OriginalCost == "R$ 0,00/mes" {
		t.Errorf("esperava custo original positivo para ferramentas comerciais, obteve %s", report.OriginalCost)
	}

	if report.SovereignCost != "R$ 0,00/mes" {
		t.Errorf("esperava custo soberano de R$ 0,00/mes, obteve %s", report.SovereignCost)
	}

	if _, ok := report.Replacements["Higgsfield"]; !ok {
		t.Errorf("esperava substituição para Higgsfield na curadoria")
	}
}

func TestCanonicalNaming(t *testing.T) {
	titles := []string{
		"How To Make Insane Video Editing With Opus 5.5",
		"Criando Talking Head Reels no CapCut",
		"Kinetic Typography Tutorial 2026",
	}

	for _, title := range titles {
		slugs := videorea.SuggestCanonicalNames(title)
		if len(slugs) != 3 {
			t.Errorf("esperava 3 sugestões de nomes, obteve %d", len(slugs))
		}
		for _, s := range slugs {
			if !strings.HasPrefix(s, "legoloop-") {
				t.Errorf("slug %q não possui o prefixo obrigatório legoloop-*", s)
			}
		}
	}
}
