package videomotion_test

import (
	"path/filepath"
	"testing"

	"legoloop/internal/videomotion"
)

func TestAssetCacheHitMiss(t *testing.T) {
	tempDir := t.TempDir()
	cache := videomotion.NewAssetCache(tempDir)

	prompt := "futuristic glowing town square at night"
	seed := "42"
	aspect := videomotion.AspectVertical

	key := cache.ComputeKey(prompt, seed, aspect)
	if len(key) != 64 { // SHA-256 tem 64 caracteres hexadecimais
		t.Fatalf("chave esperada com 64 caracteres hexadecimais, obteve %d (%s)", len(key), key)
	}

	// 1. Teste de Miss inicial
	_, found := cache.Get(key)
	if found {
		t.Errorf("esperava miss para chave inexistente no cache")
	}

	// 2. Armazenar ativo sintético fictício
	fakeImageData := []byte("fake-png-binary-data-for-broll")
	err := cache.Put(key, fakeImageData, ".png")
	if err != nil {
		t.Fatalf("cache.Put falhou: %v", err)
	}

	// 3. Teste de Hit subsequente
	cachedData, found := cache.Get(key)
	if !found {
		t.Fatalf("esperava hit após Put no cache")
	}
	if string(cachedData) != string(fakeImageData) {
		t.Errorf("dados recuperados divergem dos dados armazenados")
	}

	// 4. Teste de caminho de arquivo no disco
	expectedPath := filepath.Join(tempDir, key+".png")
	actualPath := cache.GetFilePath(key, ".png")
	if actualPath != expectedPath {
		t.Errorf("caminho esperado %s, obteve %s", expectedPath, actualPath)
	}
}

func TestAssetCache_DeterministicKey(t *testing.T) {
	cache := videomotion.NewAssetCache(t.TempDir())
	k1 := cache.ComputeKey("test prompt", "100", videomotion.AspectHorizontal)
	k2 := cache.ComputeKey("test prompt", "100", videomotion.AspectHorizontal)
	k3 := cache.ComputeKey("test prompt", "101", videomotion.AspectHorizontal)

	if k1 != k2 {
		t.Errorf("chaves para mesmos parâmetros devem ser idênticas: %s vs %s", k1, k2)
	}
	if k1 == k3 {
		t.Errorf("chaves para seeds diferentes devem divergir: %s vs %s", k1, k3)
	}
}
