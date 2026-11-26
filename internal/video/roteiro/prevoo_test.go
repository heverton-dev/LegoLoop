package roteiro

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestPreVoo(t *testing.T) {
	srvOK := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srvOK.Close()

	srvFail := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srvFail.Close()

	catBase := func(url string) CatalogoFontes {
		c := CatalogoPadrao()
		c.Conectores = []ConectorFonte{{ID: "teste", Nome: "Teste", URL: url, Tipo: "json", Perfis: []string{"tecnologia"}}}
		return c
	}

	t.Run("passa com 5 obrigatorios OK e gera apenas AVISO para itens de P2 e P3", func(t *testing.T) {
		dir := t.TempDir()
		acervoDir := filepath.Join(dir, "acervo")
		_ = os.MkdirAll(acervoDir, 0o755)
		deps := DepsPreVoo{
			ChecarLegoloop: func(context.Context) error { return nil },
			ChecarCofreRea: func() bool { return true },
			ChecarOpcional: func(nome string) bool { return false },
			HTTPClient:     srvOK.Client(),
		}
		rel := ExecutarPreVoo(context.Background(), projetoValidoTeste(), catBase(srvOK.URL), dir, acervoDir, deps)
		if rel.Bloqueado {
			t.Fatalf("nao deveria bloquear: %+v", rel.Itens)
		}
		if rel.TotalAvisos() < 5 {
			t.Fatalf("esperava ao menos 5 avisos para itens de P2 e P3, obteve %d", rel.TotalAvisos())
		}
	})

	t.Run("bloqueia quando legoloop version falha", func(t *testing.T) {
		dir := t.TempDir()
		deps := DepsPreVoo{
			ChecarLegoloop: func(context.Context) error { return ErrPreVooBloqueado },
			ChecarCofreRea: func() bool { return true },
			ChecarOpcional: func(string) bool { return true },
			HTTPClient:     srvOK.Client(),
		}
		rel := ExecutarPreVoo(context.Background(), projetoValidoTeste(), catBase(srvOK.URL), dir, dir, deps)
		if !rel.Bloqueado {
			t.Fatalf("esperava bloqueio quando legoloop version falha")
		}
	})

	t.Run("bloqueia quando nenhum conector responde status 200", func(t *testing.T) {
		dir := t.TempDir()
		deps := DepsPreVoo{
			ChecarLegoloop: func(context.Context) error { return nil },
			ChecarCofreRea: func() bool { return true },
			ChecarOpcional: func(string) bool { return true },
			HTTPClient:     srvFail.Client(),
		}
		rel := ExecutarPreVoo(context.Background(), projetoValidoTeste(), catBase(srvFail.URL), dir, dir, deps)
		if !rel.Bloqueado {
			t.Fatalf("esperava bloqueio quando sonda HTTP retorna 503")
		}
	})

	t.Run("bloqueia quando video_referencia esta preenchido e legoloop-video-rea falta no Cofre", func(t *testing.T) {
		dir := t.TempDir()
		p := projetoValidoTeste()
		p.VideoReferencia = "https://youtube.com/watch?v=abc"
		deps := DepsPreVoo{
			ChecarLegoloop: func(context.Context) error { return nil },
			ChecarCofreRea: func() bool { return false },
			ChecarOpcional: func(string) bool { return true },
			HTTPClient:     srvOK.Client(),
		}
		rel := ExecutarPreVoo(context.Background(), p, catBase(srvOK.URL), dir, dir, deps)
		if !rel.Bloqueado {
			t.Fatalf("esperava bloqueio sem legoloop-video-rea")
		}
	})
}
