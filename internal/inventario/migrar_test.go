package inventario

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
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

// caminhoFatia monta o caminho do manifesto de fatia sem citar termo proibido.
func caminhoFatia(slug string) string {
	return filepath.ToSlash(filepath.Join(caminhos.DirProcessos, slug, slug+caminhos.SufixoStudio))
}

// arvoreQuinteto cria a árvore de teste: fatias, donos, órfãos e protegido.
func arvoreQuinteto(t *testing.T, raiz string) {
	t.Helper()
	montarArvore(t, raiz, map[string]string{
		caminhoFatia("legoloop-alfa"):                       `{"upstream_signatures":{"route":"docs/routes/2026-11-01_ROUTE_ALFA.json"}}`,
		caminhoFatia("legoloop-beta"):                       `{"upstream_signatures":{"route":{"path":"docs/routes/2026-11-02_ROUTE_BETA.json","sha256":"x"}}}`,
		caminhoFatia("legoloop-gama"):                       `{"upstream_signatures":{"route":"docs/routes/2026-11-05_ROUTE_DUP.json"}}`,
		caminhoFatia("legoloop-delta"):                      `{"upstream_signatures":{"route":"docs/routes/2026-11-05_ROUTE_DUP.json"}}`,
		"docs/intents/2026-11-01_INTENT_ALFA.md":            "conteudo de intencao alfa",
		"docs/intents/2026-11-03_INTENT_SEM_DONO.md":        "conteudo de intencao orfa",
		"docs/specs/2026-11-02_SPEC_BETA.md":                "conteudo de especificacao beta",
		"docs/plans/2026-11-02_PLAN_BETA.md":                "conteudo de plano beta",
		"docs/plans/2026-11-16_PLAN_MEGA_PLANO_PRODUCAO.md": "demanda global ativa",
		"docs/routes/2026-11-01_ROUTE_ALFA.json":            "{}",
		"docs/routes/2026-11-02_ROUTE_BETA.json":            "{}",
		"docs/routes/2026-11-05_ROUTE_DUP.json":             "{}",
		"docs/routes/2026-11-04_ROUTE_ORFAO.json":           "{}",
		"docs/dags/2026-11-01_DAG_ALFA.json":                "espelho executavel",
	})
}

// TestMigraQuinteto prova a regra canônica: dono e protegido permanecem em
// docs/<tipo>/ da raiz, órfão vai ao arquivo morto com registro, nada é
// deletado e nada é gravado em processos/<slug>/docs/.
func TestMigraQuinteto(t *testing.T) {
	raiz := t.TempDir()
	arvoreQuinteto(t, raiz)
	antes := contarSob(t, raiz)
	res, err := MigrarQuinteto(raiz, []string{"*MEGA_PLANO*"})
	if err != nil {
		t.Fatalf("MigrarQuinteto: %v", err)
	}
	depois := contarSob(t, raiz)
	if depois != antes+1 {
		t.Errorf("registro ausente ou delecao: antes=%d depois=%d", antes, depois)
	}
	conferirMantidos(t, raiz)
	conferirOrfaos(t, raiz, res)
	conferirProtegidos(t, raiz, res)
}

// conferirMantidos assegura donos mantidos em docs/<tipo>/ e nada em processos/<slug>/docs/.
func conferirMantidos(t *testing.T, raiz string) {
	t.Helper()
	mantidos := []string{
		"docs/intents/2026-11-01_INTENT_ALFA.md",
		"docs/routes/2026-11-01_ROUTE_ALFA.json",
		"docs/specs/2026-11-02_SPEC_BETA.md",
		"docs/plans/2026-11-02_PLAN_BETA.md",
		"docs/routes/2026-11-02_ROUTE_BETA.json",
		"docs/routes/2026-11-05_ROUTE_DUP.json",
	}
	for _, origem := range mantidos {
		if _, err := os.Stat(filepath.Join(raiz, filepath.FromSlash(origem))); err != nil {
			t.Errorf("dono não mantido na raiz: %s (%v)", origem, err)
		}
	}
	proibidoDocsDeFatia(t, raiz)
}

// proibidoDocsDeFatia falha se algum processos/<slug>/docs/ recebe escrita.
func proibidoDocsDeFatia(t *testing.T, raiz string) {
	t.Helper()
	base := filepath.Join(raiz, caminhos.DirProcessos)
	_ = filepath.WalkDir(base, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && d.Name() == "docs" {
			t.Errorf("escrita proibida em processos/<slug>/docs/: %s", filepath.ToSlash(p))
		}
		return nil
	})
}

// conferirOrfaos assegura órfãos no arquivo morto e registro completo (AC-E05.4).
func conferirOrfaos(t *testing.T, raiz string, res ResultadoQuinteto) {
	t.Helper()
	esperados := map[string]string{
		"docs/archive/quinteto/intents/2026-11-03_INTENT_SEM_DONO.md": "docs/intents/2026-11-03_INTENT_SEM_DONO.md",
		"docs/archive/quinteto/routes/2026-11-04_ROUTE_ORFAO.json":    "docs/routes/2026-11-04_ROUTE_ORFAO.json",
	}
	for destino, origem := range esperados {
		if _, err := os.Stat(filepath.Join(raiz, filepath.FromSlash(destino))); err != nil {
			t.Errorf("orfao nao arquivado: %s (%v)", destino, err)
		}
		if _, err := os.Stat(filepath.Join(raiz, filepath.FromSlash(origem))); !os.IsNotExist(err) {
			t.Errorf("origem de orfao ainda presente: %s", origem)
		}
	}
	if len(res.Orfaos) != len(esperados) {
		t.Fatalf("Orfaos=%d; esperado %d", len(res.Orfaos), len(esperados))
	}
	registro, err := OrfaosDe(raiz)
	if err != nil {
		t.Fatalf("OrfaosDe: %v", err)
	}
	if len(registro) != len(esperados) {
		t.Fatalf("registro com %d entradas; esperado %d", len(registro), len(esperados))
	}
	for _, o := range registro {
		caminho := "docs/" + o.Caminho
		destino := "docs/" + o.Destino
		if o.Motivo != "sem_dono" || esperados[destino] != caminho {
			t.Errorf("registro incompleto: %+v", o)
		}
	}
}

// conferirProtegidos assegura protegido na raiz e espelho executável intocado.
func conferirProtegidos(t *testing.T, raiz string, res ResultadoQuinteto) {
	t.Helper()
	intactos := []string{
		"docs/plans/2026-11-16_PLAN_MEGA_PLANO_PRODUCAO.md",
		"docs/dags/2026-11-01_DAG_ALFA.json",
	}
	for _, caminho := range intactos {
		if _, err := os.Stat(filepath.Join(raiz, filepath.FromSlash(caminho))); err != nil {
			t.Errorf("arquivo intocado saiu do lugar: %s (%v)", caminho, err)
		}
	}
	esperados := map[string]bool{
		"plans/2026-11-16_PLAN_MEGA_PLANO_PRODUCAO.md": true,
		"intents/2026-11-01_INTENT_ALFA.md":            true,
		"routes/2026-11-01_ROUTE_ALFA.json":            true,
		"specs/2026-11-02_SPEC_BETA.md":                true,
		"plans/2026-11-02_PLAN_BETA.md":                true,
		"routes/2026-11-02_ROUTE_BETA.json":            true,
		"routes/2026-11-05_ROUTE_DUP.json":             true,
	}
	if len(res.Protegidos) != len(esperados) {
		t.Errorf("Protegidos=%v; esperado %d mantidos na raiz", res.Protegidos, len(esperados))
	}
	for _, p := range res.Protegidos {
		if !esperados[p] {
			t.Errorf("Protegido inesperado: %s", p)
		}
	}
}
