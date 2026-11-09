package videomotion_test

import (
	"testing"

	"legoloop/internal/videomotion"
)

func TestKineticWordFlipTiming(t *testing.T) {
	words := []videomotion.WordItem{
		{Index: 0, Word: "Inteligência", StartMs: 200, EndMs: 600, IsPunchWord: false},
		{Index: 1, Word: "ARTIFICIAL", StartMs: 650, EndMs: 1200, IsPunchWord: true},
		{Index: 2, Word: "autônoma", StartMs: 1250, EndMs: 1700, IsPunchWord: false},
	}

	engine, err := videomotion.NewKineticEngine(words, videomotion.KineticThemeConfig{
		PunchFontFamily:   "Anton-Regular",
		RegularFontFamily: "Didot-Serif",
		AccentColorHex:    "#E50914",
		DefaultColorHex:   "#FFFFFF",
	})
	if err != nil {
		t.Fatalf("NewKineticEngine falhou: %v", err)
	}

	// No tempo t=400ms: palavra ativa é "Inteligência" (palavra 0, regular)
	stateAt400 := engine.GetKineticStateAt(400)
	if stateAt400.ActiveWordIndex != 0 {
		t.Errorf("esperava palavra ativa 0 em t=400ms, obteve %d", stateAt400.ActiveWordIndex)
	}
	if stateAt400.ActiveWord.Word != "Inteligência" {
		t.Errorf("palavra ativa esperada 'Inteligência', obteve '%s'", stateAt400.ActiveWord.Word)
	}
	if stateAt400.FontFamily != "Didot-Serif" {
		t.Errorf("fonte esperada 'Didot-Serif', obteve '%s'", stateAt400.FontFamily)
	}

	// No tempo t=800ms: palavra ativa é "ARTIFICIAL" (palavra 1, punch word em vermelho Anton)
	stateAt800 := engine.GetKineticStateAt(800)
	if stateAt800.ActiveWordIndex != 1 {
		t.Errorf("esperava palavra ativa 1 em t=800ms, obteve %d", stateAt800.ActiveWordIndex)
	}
	if stateAt800.FontFamily != "Anton-Regular" {
		t.Errorf("fonte esperada 'Anton-Regular', obteve '%s'", stateAt800.FontFamily)
	}
	if stateAt800.ColorHex != "#E50914" {
		t.Errorf("cor esperada '#E50914', obteve '%s'", stateAt800.ColorHex)
	}
	if !stateAt800.IsPunch {
		t.Errorf("esperava IsPunch == true em t=800ms")
	}

	// Grid em perspectiva: inclinação e translação contínua
	gridOffset := engine.PerspectiveGridOffsetAt(1000)
	if gridOffset <= 0 {
		t.Errorf("esperava offset do grid > 0 em t=1000ms, obteve %f", gridOffset)
	}
}

func TestKineticWordFlipTiming_EmptyWords(t *testing.T) {
	_, err := videomotion.NewKineticEngine(nil, videomotion.KineticThemeConfig{})
	if err == nil {
		t.Fatal("esperava erro para lista vazia de palavras")
	}
}
