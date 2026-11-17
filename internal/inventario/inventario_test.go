// Package inventario classifica, migra e protege a documentação de docs/.
// Segue a Lei 6/15: docs/templates/ e docs/vocab/ são lei e não se movem.
package inventario

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// raizRepo devolve a raiz do repositório que contém este teste.
func raizRepo(t *testing.T) string {
	t.Helper()
	_, arquivo, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("sem localização do arquivo de teste")
	}
	raiz := filepath.Clean(filepath.Join(filepath.Dir(arquivo), "..", ".."))
	if _, err := os.Stat(filepath.Join(raiz, "go.mod")); err != nil {
		t.Fatalf("raiz do repositório inválida: %v", err)
	}
	return raiz
}

// contarArquivos devolve a quantidade de arquivos sob um diretório.
func contarArquivos(t *testing.T, dir string) int {
	t.Helper()
	total := 0
	err := filepath.WalkDir(dir, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			total++
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir: %v", err)
	}
	return total
}

// TestClassificaTodosArquivos prova a cobertura total do inventário de docs/.
func TestClassificaTodosArquivos(t *testing.T) {
	docs := filepath.Join(raizRepo(t), "docs")
	entradas, err := Varrer(docs)
	if err != nil {
		t.Fatalf("Varrer: %v", err)
	}
	if len(entradas) == 0 {
		t.Fatal("varredura vazia")
	}
	for _, e := range entradas {
		if !classeValida(e.Classe) {
			t.Errorf("arquivo sem classificação válida: %s (classe %q)", e.Caminho, e.Classe)
		}
		if e.Classe != ClassePendente && e.Classe != ClasseRever && e.Dono == "" {
			t.Errorf("arquivo sem dono: %s (classe %s)", e.Caminho, e.Classe)
		}
	}
	if total := contarArquivos(t, docs); len(entradas) != total {
		t.Errorf("contagem divergente: entradas=%d arquivos=%d", len(entradas), total)
	}
	vistos := make(map[string]bool, len(entradas))
	for _, e := range entradas {
		if vistos[e.Caminho] {
			t.Errorf("caminho duplicado na varredura: %s", e.Caminho)
		}
		vistos[e.Caminho] = true
	}
	saida := filepath.Join(t.TempDir(), "INVENTARIO.md")
	if err := GerarRelatorio(entradas, saida); err != nil {
		t.Fatalf("GerarRelatorio: %v", err)
	}
	conteudo, err := os.ReadFile(saida)
	if err != nil {
		t.Fatalf("leitura do relatório: %v", err)
	}
	texto := string(conteudo)
	for _, e := range entradas {
		if !strings.Contains(texto, e.Caminho) {
			t.Errorf("arquivo sem classificação no relatório: %s", e.Caminho)
		}
	}
}

// TestClassificaRegras prova a regra de classificação por caminho.
func TestClassificaRegras(t *testing.T) {
	casos := []struct {
		rel    string
		classe string
		dono   string
	}{
		{"templates/TEMPLATE_X.md", ClasseImplementado, "lei-6-15"},
		{"vocab/2026-11-10_ASD_STE100_CANONICAL_DICTIONARY.json", ClasseImplementado, "lei-6-15"},
		{"intents/2026-11-08_INTENT_VIDEO_BUNDLE.md", ClasseImplementado, "VIDEO_BUNDLE"},
		{"specs/2026-11-08_SPEC_VIDEO_BUNDLE.md", ClasseImplementado, "VIDEO_BUNDLE"},
		{"plans/sem-padrao-iso.md", ClassePendente, ""},
		{"routes/2026-11-08_ROUTE_VIDEO_BUNDLE.json", ClasseImplementado, "VIDEO_BUNDLE"},
		{"dags/2026-11-08_DAG_VIDEO_BUNDLE.md", ClasseImplementado, "VIDEO_BUNDLE"},
		{"manuals/2026-11-09_MANUAL_FUNCAO_N10_LEGOLOOP.md", ClasseImplementado, "N10_LEGOLOOP"},
		{"books/meu-livro/livro.json", ClasseImplementado, "meu-livro"},
		{"fichas/e01.md", ClasseImplementado, "e01"},
		{"reports/2026-11-06_REPORT_STEPT_SWITCH.md", ClasseUltrapassado, "STEPT_SWITCH"},
		{"videorea/analise/alvo.json", ClasseUltrapassado, "video-rea"},
		{"sessions/2026-11-01_resumo.md", ClasseUltrapassado, "sessao"},
		{"archive/antigo.md", ClasseUltrapassado, "arquivo-morto"},
		{"qa-web/relatorio.md", ClasseImplementado, "web-qa"},
		{"prompts/prompt.md", ClasseRever, ""},
		{"examples/exemplo.md", ClasseRever, ""},
		{"interviews/entrevista.md", ClasseRever, ""},
		{"references/referencia.md", ClasseRever, ""},
		{"solto-na-raiz.md", ClasseRever, ""},
	}
	for _, c := range casos {
		got := Classificar(c.rel)
		if got.Classe != c.classe || got.Dono != c.dono {
			t.Errorf("Classificar(%q) = {%s,%s}; esperado {%s,%s}",
				c.rel, got.Classe, got.Dono, c.classe, c.dono)
		}
	}
}

// TestRefinamentosDoVarrer prova os refinamentos de classe na varredura.
func TestRefinamentosDoVarrer(t *testing.T) {
	raiz := t.TempDir()
	docs := filepath.Join(raiz, "docs")
	arquivos := map[string]string{
		"AGENTS.md":                              "governança viva cita docs/reports/global.md",
		"docs/reports/global.md":                 "relatório global",
		"docs/reports/2026-11-01_REPORT_ALFA.md": "histórico",
		"docs/reports/2026-11-09_REPORT_ALFA.md": "histórico novo",
		"docs/intents/2026-11-08_INTENT_BETA.md": "sem fechamento",
		"docs/plans/2026-11-08_PLAN_BETA.md":     "sem fechamento",
		"docs/intents/2026-11-08_INTENT_GAMA.md": "com fechamento",
		"docs/reports/2026-11-08_REPORT_GAMA.md": "fechamento",
	}
	for caminho, conteudo := range arquivos {
		alvo := filepath.Join(raiz, filepath.FromSlash(caminho))
		if err := os.MkdirAll(filepath.Dir(alvo), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(alvo, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	entradas, err := Varrer(docs)
	if err != nil {
		t.Fatalf("Varrer: %v", err)
	}
	porCaminho := make(map[string]Entrada, len(entradas))
	for _, e := range entradas {
		porCaminho[e.Caminho] = e
	}
	esperados := map[string]struct {
		classe string
		dono   string
	}{
		"reports/global.md":                 {ClasseImplementado, "governanca-viva"},
		"reports/2026-11-01_REPORT_ALFA.md": {ClasseSuperado, "ALFA"},
		"reports/2026-11-09_REPORT_ALFA.md": {ClasseUltrapassado, "ALFA"},
		"intents/2026-11-08_INTENT_BETA.md": {ClassePendente, "BETA"},
		"plans/2026-11-08_PLAN_BETA.md":     {ClassePendente, "BETA"},
		"intents/2026-11-08_INTENT_GAMA.md": {ClasseImplementado, "GAMA"},
		"reports/2026-11-08_REPORT_GAMA.md": {ClasseUltrapassado, "GAMA"},
	}
	for caminho, want := range esperados {
		got, ok := porCaminho[caminho]
		if !ok {
			t.Errorf("caminho ausente da varredura: %s", caminho)
			continue
		}
		if got.Classe != want.classe || got.Dono != want.dono {
			t.Errorf("Varrer[%s] = {%s,%s}; esperado {%s,%s}",
				caminho, got.Classe, got.Dono, want.classe, want.dono)
		}
	}
}
