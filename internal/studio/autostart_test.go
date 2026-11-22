package studio

import (
	"net/http"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

func TestStartEmbeddedServeSessaoSemStudioExterno(t *testing.T) {
	hub := a2ui.NewHub(a2ui.DefaultCatalog())
	baseURL, stop, err := StartEmbedded(t.TempDir(), hub)
	if err != nil {
		t.Fatalf("StartEmbedded falhou: %v", err)
	}
	defer func() {
		if cErr := stop(); cErr != nil {
			t.Errorf("stop falhou: %v", cErr)
		}
	}()
	if !strings.HasPrefix(baseURL, "http://127.0.0.1:") {
		t.Fatalf("baseURL inesperada: %q", baseURL)
	}

	sess := hub.CreateSession("legoloop-concepcao")
	resp, err := http.Get(baseURL + "/sessao/" + sess.ID)
	if err != nil {
		t.Fatalf("GET /sessao falhou: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 em /sessao, obteve %d", resp.StatusCode)
	}
}
