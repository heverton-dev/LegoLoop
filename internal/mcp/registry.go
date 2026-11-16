package mcp

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// ServerSpec descreve como subir um servidor MCP externo sob demanda.
type ServerSpec struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
	Cwd     string            `json:"cwd,omitempty"`
}

// CatalogEntry guarda o nome e a descricao de uma ferramenta de um servidor.
type CatalogEntry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type registryFile struct {
	Servers map[string]ServerSpec `json:"servers"`
}

// RegistryPath devolve o caminho canonico do registro de MCPs.
func RegistryPath(home string) string {
	return filepath.Join(home, ".legoloop", "mcp-registry.json")
}

// CatalogPath devolve o caminho do indice de ferramentas dos MCPs.
func CatalogPath(home string) string {
	return filepath.Join(home, ".legoloop", "mcp-catalog.json")
}

// LoadRegistry le o registro; arquivo ausente devolve registro vazio.
func LoadRegistry(path string) (map[string]ServerSpec, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]ServerSpec{}, nil
	}
	if err != nil {
		return nil, err
	}
	var rf registryFile
	if err := json.Unmarshal(data, &rf); err != nil {
		return nil, err
	}
	if rf.Servers == nil {
		rf.Servers = map[string]ServerSpec{}
	}
	return rf.Servers, nil
}

// SaveRegistry grava o registro com permissao restrita (pode conter tokens).
func SaveRegistry(path string, servers map[string]ServerSpec) error {
	return writeJSONFile(path, registryFile{Servers: servers})
}

// MergeRegistry acrescenta ao destino os servidores ausentes, sem sobrescrever.
func MergeRegistry(dst, src map[string]ServerSpec) map[string]ServerSpec {
	if dst == nil {
		dst = map[string]ServerSpec{}
	}
	for name, spec := range src {
		if _, ok := dst[name]; !ok {
			dst[name] = spec
		}
	}
	return dst
}

// LoadCatalog le o indice de ferramentas; ausente devolve vazio.
func LoadCatalog(path string) (map[string][]CatalogEntry, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string][]CatalogEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var cat map[string][]CatalogEntry
	if err := json.Unmarshal(data, &cat); err != nil {
		return nil, err
	}
	return cat, nil
}

// SaveCatalog grava o indice de ferramentas.
func SaveCatalog(path string, cat map[string][]CatalogEntry) error {
	return writeJSONFile(path, cat)
}

func writeJSONFile(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}
