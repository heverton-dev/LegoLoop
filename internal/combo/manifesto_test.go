package combo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGerarManifesto9Router(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "manifesto_test.json")

	err := GerarManifesto9Router(outPath)
	if err != nil {
		t.Fatalf("erro inesperado ao gerar manifesto: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("falha ao ler arquivo gerado: %v", err)
	}

	var parsed struct {
		Combos []struct {
			Nome             string `json:"nome"`
			Finalidade       string `json:"finalidade"`
			MaxTokens        int    `json:"max_tokens"`
			ExigeFerramentas bool   `json:"exige_ferramentas"`
			ExigeImagem      bool   `json:"exige_imagem"`
			Modelos          []struct {
				ID       string  `json:"id"`
				Provedor string  `json:"provedor"`
				Latencia float64 `json:"latencia"`
			} `json:"modelos"`
		} `json:"combos"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("JSON invalido: %v", err)
	}

	if len(parsed.Combos) != 6 {
		t.Fatalf("esperava 6 combos, obteve %d", len(parsed.Combos))
	}

	mapa := make(map[string]int)
	for _, c := range parsed.Combos {
		mapa[c.Nome] = c.MaxTokens
		if len(c.Modelos) == 0 {
			t.Errorf("combo %s nao possui modelos", c.Nome)
		}
	}

	esperados := map[string]int{
		"t-agente-128k":     131072,
		"t-chat-128k":       131072,
		"t-raciocinio-128k": 131072,
		"t-visao-128k":      131072,
		"t-codigo-128k":     131072,
		"t-micro-8k":        8192,
	}

	for k, maxTokens := range esperados {
		got, ok := mapa[k]
		if !ok {
			t.Errorf("combo esperado %s nao encontrado no manifesto", k)
		} else if got != maxTokens {
			t.Errorf("combo %s: esperado max_tokens=%d, obteve %d", k, maxTokens, got)
		}
	}
}

func TestGerarManifesto9Router_CaminhoDefault(t *testing.T) {
	defaultPath := filepath.Join("..", "..", "scratch", "combos_9router_manifest.json")
	_ = os.Remove(defaultPath)

	err := GerarManifesto9Router("")
	if err != nil {
		t.Fatalf("erro ao gerar manifesto no caminho default: %v", err)
	}

	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		t.Fatalf("arquivo default nao foi criado em %s", defaultPath)
	}
}
