package inventario

import (
	"os"
	"path/filepath"
	"testing"
)

// montarArvore cria arquivos a partir de um mapa caminho-conteúdo sob uma raiz.
func montarArvore(t *testing.T, raiz string, arquivos map[string]string) {
	t.Helper()
	for caminho, conteudo := range arquivos {
		alvo := filepath.Join(raiz, filepath.FromSlash(caminho))
		if err := os.MkdirAll(filepath.Dir(alvo), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(alvo, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// contarSob conta arquivos sob um diretório; zero se ele não existe.
func contarSob(t *testing.T, dir string) int {
	t.Helper()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return 0
	}
	return contarArquivos(t, dir)
}

// TestMigraParaArchive prova a migração sem deleção para docs/archive/.
func TestMigraParaArchive(t *testing.T) {
	raiz := t.TempDir()
	montarArvore(t, raiz, map[string]string{
		"AGENTS.md":                           "governança cita docs/reports/global.md",
		"docs/videorea/alvo.json":             "análise de alvo",
		"docs/videorea/nivel/detalhe.md":      "detalhe",
		"docs/sessions/2026-11-01_resumo.md":  "resumo de sessão",
		"docs/reports/historico.md":           "relatório histórico",
		"docs/reports/img/captura.png":        "captura",
		"docs/reports/global.md":              "relatório global",
		"docs/templates/TEMPLATE_X.md":        "lei",
		"docs/vocab/dicionario.json":          "lei",
		"docs/archive/antigo.md":              "já arquivado",
		"docs/intents/2026-11-08_INTENT_X.md": "quinteto fica",
	})
	antes := contarSob(t, filepath.Join(raiz, "docs"))
	res, err := Migrar(raiz)
	if err != nil {
		t.Fatalf("Migrar: %v", err)
	}
	depois := contarSob(t, filepath.Join(raiz, "docs"))
	if antes != depois {
		t.Errorf("migração deletou arquivos: antes=%d depois=%d", antes, depois)
	}
	movidos := []string{
		"docs/archive/videorea/alvo.json",
		"docs/archive/videorea/nivel/detalhe.md",
		"docs/archive/sessions/2026-11-01_resumo.md",
		"docs/archive/reports/historico.md",
		"docs/archive/reports/img/captura.png",
	}
	for _, caminho := range movidos {
		if _, err := os.Stat(filepath.Join(raiz, filepath.FromSlash(caminho))); err != nil {
			t.Errorf("não migrado para archive: %s (%v)", caminho, err)
		}
	}
	permanecem := []string{
		"docs/reports/global.md",
		"docs/templates/TEMPLATE_X.md",
		"docs/vocab/dicionario.json",
		"docs/archive/antigo.md",
		"docs/intents/2026-11-08_INTENT_X.md",
	}
	for _, caminho := range permanecem {
		if _, err := os.Stat(filepath.Join(raiz, filepath.FromSlash(caminho))); err != nil {
			t.Errorf("arquivo movido indevidamente: %s (%v)", caminho, err)
		}
	}
	origens := []string{
		"docs/videorea/alvo.json",
		"docs/sessions/2026-11-01_resumo.md",
		"docs/reports/historico.md",
	}
	for _, caminho := range origens {
		if _, err := os.Stat(filepath.Join(raiz, filepath.FromSlash(caminho))); !os.IsNotExist(err) {
			t.Errorf("origem ainda presente após migração: %s", caminho)
		}
	}
	if len(res.Movidos) != len(movidos) {
		t.Errorf("Movidos=%d; esperado %d", len(res.Movidos), len(movidos))
	}
	if len(res.Protegidos) != 1 || res.Protegidos[0] != "reports/global.md" {
		t.Errorf("Protegidos=%v; esperado [reports/global.md]", res.Protegidos)
	}
}
