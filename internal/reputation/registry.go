package reputation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// RegistryMatrix guarda todas as notas ranqueadas no cofre canônico.
type RegistryMatrix struct {
	GeneratedAt   string        `json:"generated_at"`
	MatrixVersion string        `json:"matrix_version"`
	TotalPieces   int           `json:"total_pieces"`
	Ratings       []PieceRating `json:"ratings"`

	byKey map[string]*PieceRating
	mu    sync.RWMutex
}

// DefaultRegistryPath devolve o caminho canônico do arquivo no SO.
func DefaultRegistryPath(home string) string {
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	return filepath.Join(home, ".aidd", "canonical", "registry_matrix.json")
}

// LoadRegistry lê e indexa o arquivo de matriz canônica.
func LoadRegistry(filePath string) (*RegistryMatrix, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var reg RegistryMatrix
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	reg.byKey = make(map[string]*PieceRating, len(reg.Ratings)*2)
	for i := range reg.Ratings {
		r := &reg.Ratings[i]
		if r.Name != "" {
			reg.byKey[strings.ToLower(r.Name)] = r
		}
		if r.Ref != "" {
			reg.byKey[strings.ToLower(r.Ref)] = r
			reg.byKey[strings.ToLower(filepath.Base(r.Ref))] = r
		}
	}
	return &reg, nil
}

// Find busca a nota de uma peça por nome ou referência.
func (rm *RegistryMatrix) Find(key string) *PieceRating {
	if rm == nil {
		return nil
	}
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	if rm.byKey == nil {
		return nil
	}
	norm := strings.ToLower(strings.TrimSpace(key))
	if r, ok := rm.byKey[norm]; ok {
		return r
	}
	base := strings.ToLower(filepath.Base(norm))
	if r, ok := rm.byKey[base]; ok {
		return r
	}
	return nil
}
