package inventario

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// HashArquivo devolve o SHA-256 em hexadecimal do arquivo.
func HashArquivo(caminho string) (string, error) {
	f, err := os.Open(caminho)
	if err != nil {
		return "", fmt.Errorf("video-inventario: abrindo %s: %w", caminho, err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("video-inventario: lendo %s: %w", caminho, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Verificar confere o arquivo pesado do item pelo hash; ausente ou alterado retorna ErrHashDivergente.
func (s Store) Verificar(it Item) error {
	if it.Arquivo == "" {
		return nil
	}
	caminho := it.Arquivo
	if !filepath.IsAbs(caminho) {
		caminho = filepath.Join(s.Raiz, caminho)
	}
	h, err := HashArquivo(caminho)
	if err != nil {
		return fmt.Errorf("%w: item %q: %v", ErrHashDivergente, it.ID, err)
	}
	if h != it.Hash {
		return fmt.Errorf("%w: item %q: esperado %s, obtido %s", ErrHashDivergente, it.ID, it.Hash, h)
	}
	return nil
}
