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

func TestCurateAgainstMatrix_CarregaMatrizOriginal(t *testing.T) {
	dossier := &videorea.ForensicsDossier{
		Title:           "Adeus IAs pagas: Nova IA gera videos em PC fraco",
		ToolsIdentified: []string{"ComfyUI", "MiniMax H3", "FLUX"},
	}

	report, err := videorea.CurateAgainstMatrix(dossier)
	if err != nil {
		t.Fatalf("CurateAgainstMatrix falhou: %v", err)
	}

	if report.VideoTitle != dossier.Title {
		t.Errorf("video_title esperado %q, obteve %q", dossier.Title, report.VideoTitle)
	}
	if report.OriginalMatrix == nil {
		t.Fatalf("esperava original_matrix no relatório")
	}
	if len(report.OriginalMatrix.Tools) != 3 {
		t.Errorf("esperava 3 ferramentas na matriz, obteve %d", len(report.OriginalMatrix.Tools))
	}
}

func findingByName(r *videorea.CurationReport, name string) *videorea.Finding {
	for i := range r.Findings {
		if r.Findings[i].Name == name {
			return &r.Findings[i]
		}
	}
	return nil
}

func TestCurateAgainstMatrix_ClassificaTodosOsAchados(t *testing.T) {
	dossier := &videorea.ForensicsDossier{
		Title:           "Analise de video",
		ToolsIdentified: []string{"Higgsfield", "ComfyUI", "Kling"},
		HiddenTools:     []string{"FFmpeg CLI"},
		Concepts:        []string{"DAG acíclico"},
		Processes:       []string{"ingestão e corte"},
		Repositories:    []string{"github.com/obra/superpowers"},
	}
	report, err := videorea.CurateAgainstMatrix(dossier)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Findings) != 7 {
		t.Fatalf("achados classificados = %d, esperado 7: %+v", len(report.Findings), report.Findings)
	}
	for _, f := range report.Findings {
		switch f.Class {
		case videorea.ClassDeve, videorea.ClassPode, videorea.ClassNaoPode:
		default:
			t.Errorf("%s sem classe válida: %q", f.Name, f.Class)
		}
		if strings.TrimSpace(f.Why) == "" || strings.TrimSpace(f.How) == "" {
			t.Errorf("%s sem porquê ou sem como: %+v", f.Name, f)
		}
	}
	if f := findingByName(report, "Higgsfield"); f == nil || f.Class != videorea.ClassNaoPode {
		t.Errorf("ferramenta paga deve ser NAO_PODE: %+v", f)
	}
	if f := findingByName(report, "Kling"); f == nil || f.Class != videorea.ClassPode {
		t.Errorf("ferramenta sem licença conhecida deve ser PODE: %+v", f)
	}
	if f := findingByName(report, "github.com/obra/superpowers"); f == nil || f.Class != videorea.ClassPode || f.Kind != "repositorio" {
		t.Errorf("repositório deve ser PODE: %+v", f)
	}
}

func TestCurateAgainstMatrix_FerramentaGratisSegueOEncaixe(t *testing.T) {
	dossier := &videorea.ForensicsDossier{Title: "x", ToolsIdentified: []string{"ComfyUI"}}
	report, err := videorea.CurateAgainstMatrix(dossier)
	if err != nil {
		t.Fatal(err)
	}
	f := findingByName(report, "ComfyUI")
	if f == nil || len(report.AdoptedTools) != 1 {
		t.Fatalf("esperava 1 ferramenta adotada e seu achado: %+v", report.Findings)
	}
	ouro := report.AdoptedTools[0].Diagnosis.FitStatus == videorea.FitOuro
	if ouro != (f.Class == videorea.ClassDeve) {
		t.Errorf("DEVE deve coincidir com FIT_OURO: classe %s, fit %s", f.Class, report.AdoptedTools[0].Diagnosis.FitStatus)
	}
}
