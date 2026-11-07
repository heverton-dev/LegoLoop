package plangen

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileSignature representa o caminho e o hash de um artefato upstream.
type FileSignature struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// UpstreamSignatures agrega a assinatura encadeada dos 3 artefatos anteriores da Tríade.
type UpstreamSignatures struct {
	Intent *FileSignature `json:"intent,omitempty"`
	Spec   FileSignature  `json:"spec"`
	Plan   FileSignature  `json:"plan"`
}

// HashNormalized calcula o SHA-256 após normalizar quebras de linha Windows para LF.
func HashNormalized(data []byte) string {
	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

// ReadAndHash lê o arquivo e devolve o hash SHA-256 com quebras de linha normalizadas.
func ReadAndHash(fullPath string) (string, error) {
	b, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return HashNormalized(b), nil
}

// FileSignatureFor calcula a assinatura de um arquivo relativo se ele existir.
func FileSignatureFor(baseDir, relPath string) *FileSignature {
	if relPath == "" {
		return nil
	}
	full := filepath.Join(baseDir, filepath.FromSlash(relPath))
	hash, err := ReadAndHash(full)
	if err != nil {
		return &FileSignature{
			Path:   filepath.ToSlash(relPath),
			SHA256: strings.Repeat("0", 64),
		}
	}
	return &FileSignature{
		Path:   filepath.ToSlash(relPath),
		SHA256: hash,
	}
}

// VerifySignatures valida os hashes dos arquivos upstream contra o disco.
func VerifySignatures(baseDir string, sig UpstreamSignatures) error {
	checks := []FileSignature{sig.Spec, sig.Plan}
	if sig.Intent != nil {
		checks = append(checks, *sig.Intent)
	}
	for _, c := range checks {
		if c.SHA256 == strings.Repeat("0", 64) {
			continue
		}
		full := filepath.Join(baseDir, filepath.FromSlash(c.Path))
		got, err := ReadAndHash(full)
		if err != nil {
			return fmt.Errorf("plangen: falha ao ler upstream %s: %w", c.Path, err)
		}
		if got != c.SHA256 {
			return fmt.Errorf("plangen: integridade criptográfica violada (mismatch) em %s", c.Path)
		}
	}
	return nil
}
