package studio

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"legoloop/internal/brain"
	"legoloop/internal/provider"
)

// O 9Router deixa de ser requisito no Studio: só com Ollama a fábrica real roda de ponta a ponta.
func TestStudioCerebro_OllamaSemRouterFabrica(t *testing.T) {
	var modelo atomic.Value
	modelo.Store("")
	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Model string }
		_ = json.NewDecoder(r.Body).Decode(&req)
		modelo.Store(req.Model)
		_ = json.NewEncoder(w).Encode(map[string]any{"response": studioSpec})
	}))
	defer ollama.Close()

	s, _, _ := newAssetsServer(t)
	s.Discover = func(context.Context) provider.Inventory { return provider.Inventory{Ollama: []string{"qwen2.5"}} }
	s.BrainOpts = []brain.Option{brain.WithOllamaURL(ollama.URL)}
	s.FabricarScanner = studioScan{}

	code, body := postFabricar(t, s, `{"name":"Extrator","goal":"extrair","category":"build","model":"qwen2.5"}`)
	if code != http.StatusOK {
		t.Fatalf("status %d: %s", code, body)
	}
	if modelo.Load().(string) != "qwen2.5" {
		t.Errorf("o Ollama deveria conduzir a fabricação (modelo=%v)", modelo.Load())
	}
	n := 0
	_ = filepath.WalkDir(s.ProjectDir, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			n++
		}
		return nil
	})
	if n != 5 {
		t.Errorf("devia gravar os 5 documentos no projeto, gravou %d", n)
	}
}

func TestStudioCerebro_ErrosDoCerebroViramStatusComMensagemPropria(t *testing.T) {
	cases := []struct {
		err    error
		status int
		diz    string
	}{
		{brain.ErrModeloDesconhecido, http.StatusUnprocessableEntity, "não está disponível"},
		{brain.ErrSemCerebro, http.StatusServiceUnavailable, "Ollama"},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		writeFabricarError(rec, c.err)
		if rec.Code != c.status || !strings.Contains(rec.Body.String(), c.diz) {
			t.Errorf("%v -> %d %s; esperado %d com %q", c.err, rec.Code, rec.Body.String(), c.status, c.diz)
		}
		if strings.Contains(rec.Body.String(), "127.0.0.1") {
			t.Errorf("a mensagem não pode vazar detalhes locais: %s", rec.Body.String())
		}
	}
}
