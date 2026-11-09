package videomotion

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AssetCache gerencia o armazenamento determinístico de ativos indexados por hash.
type AssetCache struct {
	baseDir string
}

// NewAssetCache inicializa o diretório base do cache no disco.
func NewAssetCache(baseDir string) *AssetCache {
	_ = os.MkdirAll(baseDir, 0755)
	return &AssetCache{baseDir: baseDir}
}

// ComputeKey gera o hash SHA-256 unívoco a partir do prompt, seed e proporção de tela.
func (c *AssetCache) ComputeKey(prompt, seed string, aspect AspectRatio) string {
	payload := fmt.Sprintf("%s|%s|%s", strings.TrimSpace(prompt), strings.TrimSpace(seed), aspect)
	h := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(h[:])
}

// GetFilePath devolve o caminho absoluto do arquivo no disco.
func (c *AssetCache) GetFilePath(key, ext string) string {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return filepath.Join(c.baseDir, key+ext)
}

// Get recupera o conteúdo binário caso exista no disco.
func (c *AssetCache) Get(key string) ([]byte, bool) {
	matches, err := filepath.Glob(filepath.Join(c.baseDir, key+".*"))
	if err != nil || len(matches) == 0 {
		return nil, false
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		return nil, false
	}
	return data, true
}

// Put grava os dados binários do ativo no disco.
func (c *AssetCache) Put(key string, data []byte, ext string) error {
	filePath := c.GetFilePath(key, ext)
	return os.WriteFile(filePath, data, 0644)
}
