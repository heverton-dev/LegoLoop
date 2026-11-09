package videomotion_test

import (
	"os"
	"testing"

	"legoloop/internal/videomotion"
)

func TestBridgeGenerationAndFallback(t *testing.T) {
	tempDir := t.TempDir()
	cache := videomotion.NewAssetCache(tempDir)

	// Inicializa a bridge sem chave de API para forçar o fallback gracioso local
	bridge := videomotion.NewStudioAssetBridge("", cache)

	prompt := "cinematic town scene with neon signs"
	seed := "99"
	aspect := videomotion.AspectVertical

	res, err := bridge.ResolveAsset(prompt, seed, aspect)
	if err != nil {
		t.Fatalf("ResolveAsset falhou: %v", err)
	}

	if !res.UsedFallback {
		t.Errorf("esperava UsedFallback == true quando a API key está ausente")
	}

	if res.FromCache {
		t.Errorf("primeira requisição não deve vir do cache")
	}

	if _, err := os.Stat(res.FilePath); os.IsNotExist(err) {
		t.Fatalf("arquivo de fallback não foi gravado no disco: %s", res.FilePath)
	}

	// Segunda chamada com mesmo prompt/seed/aspect deve vir 100% do cache
	res2, err := bridge.ResolveAsset(prompt, seed, aspect)
	if err != nil {
		t.Fatalf("segunda chamada de ResolveAsset falhou: %v", err)
	}

	if !res2.FromCache {
		t.Errorf("segunda requisição deve vir com FromCache == true")
	}
	if res2.FilePath != res.FilePath {
		t.Errorf("caminho esperado igual ao da primeira chamada: %s vs %s", res.FilePath, res2.FilePath)
	}
}

func TestBridgeGeneration_Validation(t *testing.T) {
	bridge := videomotion.NewStudioAssetBridge("", videomotion.NewAssetCache(t.TempDir()))
	_, err := bridge.ResolveAsset("", "1", videomotion.AspectVertical)
	if err == nil {
		t.Fatal("esperava erro para prompt vazio")
	}
}
