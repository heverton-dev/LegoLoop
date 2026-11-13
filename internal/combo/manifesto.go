package combo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ModeloManifesto representa um modelo no manifesto do 9Router.
type ModeloManifesto struct {
	ID       string  `json:"id"`
	Provedor string  `json:"provedor"`
	Latencia float64 `json:"latencia"`
}

// ComboManifesto define a estrutura de um combo no manifesto exportado.
type ComboManifesto struct {
	Nome             string            `json:"nome"`
	Finalidade       string            `json:"finalidade"`
	MaxTokens        int               `json:"max_tokens"`
	ExigeFerramentas bool              `json:"exige_ferramentas"`
	ExigeImagem      bool              `json:"exige_imagem"`
	Modelos          []ModeloManifesto `json:"modelos"`
}

// Manifesto9Router define o envelope raiz do manifesto para o 9Router.
type Manifesto9Router struct {
	Combos []ComboManifesto `json:"combos"`
}

func resolverCaminhoManifesto(caminho string) string {
	if caminho != "" {
		return caminho
	}
	dir, err := os.Getwd()
	if err == nil {
		curr := dir
		for {
			if _, err := os.Stat(filepath.Join(curr, "go.mod")); err == nil {
				return filepath.Join(curr, "scratch", "combos_9router_manifest.json")
			}
			parent := filepath.Dir(curr)
			if parent == curr {
				break
			}
			curr = parent
		}
	}
	return filepath.Join("scratch", "combos_9router_manifest.json")
}

func construirManifesto() Manifesto9Router {
	nomes := ListarCombos()
	manifesto := Manifesto9Router{Combos: make([]ComboManifesto, 0, len(nomes))}
	for _, nome := range nomes {
		c := CatalogoCombos[nome]
		modelos := make([]ModeloManifesto, 0, len(c.Membros))
		for _, m := range c.Membros {
			modelos = append(modelos, ModeloManifesto{
				ID:       m.ID,
				Provedor: m.Provedor,
				Latencia: m.LatenciaMedia100k,
			})
		}
		manifesto.Combos = append(manifesto.Combos, ComboManifesto{
			Nome:             c.Nome,
			Finalidade:       c.Finalidade,
			MaxTokens:        c.MaxTokens,
			ExigeFerramentas: c.ExigeFerramentas,
			ExigeImagem:      c.ExigeImagem,
			Modelos:          modelos,
		})
	}
	return manifesto
}

func gravarArquivoAtomico(destino string, dados []byte) error {
	dir := filepath.Dir(destino)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("manifesto: falha ao criar diretorio: %w", err)
	}
	tmpFile, err := os.CreateTemp(dir, "manifesto_*.tmp")
	if err != nil {
		return fmt.Errorf("manifesto: falha ao criar arquivo temporario: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	if _, err := tmpFile.Write(dados); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("manifesto: falha ao gravar conteudo: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("manifesto: falha ao fechar arquivo temporario: %w", err)
	}
	if err := os.Rename(tmpPath, destino); err != nil {
		return fmt.Errorf("manifesto: falha ao renomear atomico: %w", err)
	}
	return nil
}

// GerarManifesto9Router exporta os 6 combos certificados em formato JSON
// aceito pelo gateway 9Router, gravando atomicamente no destino indicado.
func GerarManifesto9Router(caminhoSaida string) error {
	destino := resolverCaminhoManifesto(caminhoSaida)
	manifesto := construirManifesto()
	dados, err := json.MarshalIndent(manifesto, "", "  ")
	if err != nil {
		return fmt.Errorf("manifesto: falha ao serializar JSON: %w", err)
	}
	dados = append(dados, '\n')
	return gravarArquivoAtomico(destino, dados)
}
