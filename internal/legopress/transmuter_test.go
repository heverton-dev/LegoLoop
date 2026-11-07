package legopress_test

import (
	"strings"
	"testing"

	"legoloop/internal/legopress"
)

func TestTransmuterAndSync(t *testing.T) {
	// 1. Validar Purificação Didática contra IA-Slop (AC-E04.2)
	transmuter := legopress.NewTransmuter()
	rawInput := `Com certeza! Vou ficar muito feliz em ajudar você nessa jornada.
Aqui está a explicação passo a passo com certeza absoluta:
O motor em Go opera com goroutines leves e contexto estéril.
Espero que isso tenha ajudado! Se precisar de mais alguma coisa, estou à disposição.`

	cleanText := transmuter.Sanitize(rawInput)

	if strings.Contains(cleanText, "Com certeza!") {
		t.Errorf("esperava remocao de saudacao IA-slop, mas encontrou: %s", cleanText)
	}
	if strings.Contains(cleanText, "Espero que isso tenha ajudado") {
		t.Errorf("esperava remocao de despedida prolixa, mas encontrou: %s", cleanText)
	}
	if !strings.Contains(cleanText, "O motor em Go opera com goroutines leves e contexto estéril.") {
		t.Errorf("esperava conteudo essencial preservado, obteve: %s", cleanText)
	}

	// 2. Validar Sincronizacao por Git Diff (AC-E04.1)
	syncEngine := legopress.NewSyncEngine()
	manifest := legopress.Manifest{
		Title: "LegoLoop Book",
		Parts: []legopress.Part{
			{ID: "P01", Title: "Macro", File: "partes/01-macro.md"},
			{ID: "P04", Title: "Micro", File: "partes/04-micro-o-motor-e-governanca.md"},
		},
	}

	diffs := []string{
		"internal/dag/runner.go",
		"internal/stept/stept.go",
	}

	staleParts := syncEngine.DetectStaleParts(diffs, manifest)
	if len(staleParts) == 0 {
		t.Fatal("esperava detectar partes defasadas a partir de diffs em Go, obteve zero")
	}

	foundMicro := false
	for _, p := range staleParts {
		if strings.Contains(p.File, "04-micro") {
			foundMicro = true
			break
		}
	}
	if !foundMicro {
		t.Errorf("esperava '04-micro' marcado como defasado por diffs no motor Go: %+v", staleParts)
	}
}
