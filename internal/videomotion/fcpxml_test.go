package videomotion_test

import (
	"encoding/xml"
	"testing"

	"legoloop/internal/videomotion"
)

func TestExportFCPXMLTimeline(t *testing.T) {
	clips := []videomotion.TimelineClip{
		{Name: "intro_hook", FilePath: "intro.mp4", DurationSec: 3.5},
		{Name: "broll_concept", FilePath: "concept.mp4", DurationSec: 4.2},
		{Name: "cta_outro", FilePath: "cta.mp4", DurationSec: 2.8},
	}

	spec := videomotion.FCPXMLSpec{
		ProjectName: "LegoLoop Video Project",
		FrameRate:   25,
		Width:       1080,
		Height:      1920,
		Clips:       clips,
	}

	xmlData, err := videomotion.GenerateFCPXML(spec)
	if err != nil {
		t.Fatalf("GenerateFCPXML falhou: %v", err)
	}

	xmlStr := string(xmlData)

	// 1. Validar cabeçalho DTD e versão FCPXML
	if !containsStr(xmlStr, "<!DOCTYPE fcpxml>") || !containsStr(xmlStr, `<fcpxml version="1.9">`) {
		t.Errorf("cabeçalho DTD ou versão fcpxml ausente:\n%s", xmlStr)
	}

	// 2. Validar que é um XML sintaticamente válido
	var parsed struct {
		XMLName xml.Name `xml:"fcpxml"`
		Version string   `xml:"version,attr"`
	}
	if err := xml.Unmarshal(xmlData, &parsed); err != nil {
		t.Fatalf("XML gerado é sintaticamente inválido: %v", err)
	}
	if parsed.Version != "1.9" {
		t.Errorf("versão esperada 1.9, obteve %s", parsed.Version)
	}

	// 3. Validar presença dos clipes no spine
	if !containsStr(xmlStr, `name="intro_hook"`) || !containsStr(xmlStr, `name="broll_concept"`) {
		t.Errorf("clipes ausentes no XML gerado")
	}
}

func TestExportFCPXMLTimeline_EmptyClips(t *testing.T) {
	spec := videomotion.FCPXMLSpec{
		ProjectName: "Empty",
		FrameRate:   25,
		Clips:       nil,
	}
	_, err := videomotion.GenerateFCPXML(spec)
	if err == nil {
		t.Fatal("esperava erro para lista vazia de clipes")
	}
}
