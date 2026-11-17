package fechar_test

import (
	"context"
	"encoding/json"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/fechar"
	"legoloop/internal/processo"
	"legoloop/internal/triagem"
)

// dirDialetoAntigo é a pasta de Processos do dialeto antigo; os testes de fechamento a usam por este nome.
var dirDialetoAntigo = caminhos.DirLoops

func TestValidateAndNormalizeSlug(t *testing.T) {
	tests := []struct {
		name      string
		req       fechar.Request
		wantSlug  string
		wantCanon string
		wantErr   bool
	}{
		{
			name:      "slug simples sem prefixo",
			req:       fechar.Request{Slug: "video-motion"},
			wantSlug:  "video-motion",
			wantCanon: "legoloop-video-motion",
			wantErr:   false,
		},
		{
			name:      "slug já com prefixo legoloop-",
			req:       fechar.Request{Slug: "legoloop-video-rea"},
			wantSlug:  "video-rea",
			wantCanon: "legoloop-video-rea",
			wantErr:   false,
		},
		{
			name:      "slug com espacos ao redor",
			req:       fechar.Request{Slug: "  video-motion  "},
			wantSlug:  "video-motion",
			wantCanon: "legoloop-video-motion",
			wantErr:   false,
		},
		{
			name:    "slug vazio deve falhar",
			req:     fechar.Request{Slug: ""},
			wantErr: true,
		},
		{
			name:    "slug com espaco interno deve falhar",
			req:     fechar.Request{Slug: "video motion"},
			wantErr: true,
		},
		{
			name:    "slug com path traversal deve falhar",
			req:     fechar.Request{Slug: "../malicioso"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, canon, err := fechar.NormalizeSlug(tt.req.Slug)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizeSlug(%q) err = %v, wantErr = %v", tt.req.Slug, err, tt.wantErr)
			}
			if !tt.wantErr {
				if slug != tt.wantSlug {
					t.Errorf("slug = %q, want %q", slug, tt.wantSlug)
				}
				if canon != tt.wantCanon {
					t.Errorf("canon = %q, want %q", canon, tt.wantCanon)
				}
				if err := tt.req.Validate(); err != nil {
					t.Errorf("req.Validate() falhou inesperadamente: %v", err)
				}
			}
		})
	}
}

func TestExecutarNoveEtapasEGerarReportComTelemetria(t *testing.T) {
	tmp := t.TempDir()
	semearDocs(t, tmp, "video-motion")
	semearDAG(t, tmp, "video-motion")
	deps := fechar.Deps{
		RepoDir:   tmp,
		VaultDir:  tmp,
		Home:      tmp,
		SkipBuild: true,
	}

	customTelemetria := &fechar.TelemetriaOrcamento{
		TempoEstimadoMin:   20.0,
		TempoGastoMin:      15.5,
		TokensInEstimados:  50000,
		TokensInGastos:     42000,
		TokensOutEstimados: 12000,
		TokensOutGastos:    9500,
		TurnosRealizados:   8,
		ToolCalls:          18,
		CustoEstimadoUSD:   0.33,
		CustoGastoUSD:      0.268,
		TaxaCambioBRL:      5.85,
		Modelos:            []string{"claude-3-7-sonnet"},
	}

	req := fechar.Request{
		Slug:       "video-motion",
		DryRun:     false,
		SkipBuild:  true,
		Telemetria: customTelemetria,
	}

	res, err := fechar.Executar(nil, req, deps)
	if err != nil {
		t.Fatalf("Executar falhou: %v", err)
	}
	if !res.Sucesso {
		t.Fatalf("esperava sucesso no fechamento, obteve false")
	}
	if len(res.Etapas) != 9 {
		t.Fatalf("esperava 9 etapas executadas, obteve %d", len(res.Etapas))
	}
	if res.Etapas[8].Nome != "Relatório de Conclusão & Telemetria" {
		t.Errorf("etapa 9 = %q, want 'Relatório de Conclusão & Telemetria'", res.Etapas[8].Nome)
	}
	if res.ReportPath == "" {
		t.Fatalf("res.ReportPath não deve estar vazio")
	}

	conteudo, err := os.ReadFile(res.ReportPath)
	if err != nil {
		t.Fatalf("falha ao ler relatório gerado em %s: %v", res.ReportPath, err)
	}
	corpo := string(conteudo)
	for _, esperado := range []string{
		"Balanço de Tempo & Cronometria",
		"Balanço de Recursos Computacionais & Tokens",
		"Balanço Financeiro & Custo",
		"20.0 min",
		"15.5 min",
		"42000",
		"9500",
		"claude-3-7-sonnet",
		"Matriz de Qualidade",
	} {
		if !strings.Contains(corpo, esperado) {
			t.Errorf("relatório não contém string esperada %q", esperado)
		}
	}
}

func TestFecharDAG_AuditaParidadeSHA256(t *testing.T) {
	casos := []struct {
		nome, pasta, arquivo string
	}{
		{"dialeto antigo", caminhos.DirLoops, caminhos.ArquivoDAG},
		{"dialeto novo", caminhos.DirProcessos, caminhos.ArquivoProcessoDAG},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			testarParidadeSHA256(t, c.pasta, c.arquivo)
		})
	}
}

// TestParidadeDoisNomes cobre AC-E01.5: dois Processos no mesmo repositório,
// um com o nome antigo do DAG e outro com o novo, passam na auditoria de
// paridade; divergência no nome novo continua sendo detectada.
func TestParidadeDoisNomes(t *testing.T) {
	dir := t.TempDir()
	dagsDir := filepath.Join(dir, "docs", "dags")
	if err := os.MkdirAll(dagsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	conteudo := []byte(`{"name":"x","version":"1.0.0","nodes":[{"id":"01-alfa"}]}`)
	procs := []struct{ slug, pasta, arquivo string }{
		{"alfa", caminhos.DirLoops, caminhos.ArquivoDAG},
		{"beta", caminhos.DirProcessos, caminhos.ArquivoProcessoDAG},
	}
	for _, p := range procs {
		semearDocs(t, dir, p.slug)
		pdir := filepath.Join(dir, p.pasta, p.slug)
		if err := os.MkdirAll(pdir, 0o755); err != nil {
			t.Fatal(err)
		}
		iso := filepath.Join(dagsDir, "2026-11-10_DAG_"+strings.ToUpper(p.slug)+".json")
		for _, f := range []string{iso, filepath.Join(pdir, p.arquivo)} {
			if err := os.WriteFile(f, conteudo, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	deps := fechar.Deps{RepoDir: dir, SkipBuild: true}
	for _, p := range procs {
		res, err := fechar.Executar(nil, fechar.Request{Slug: p.slug}, deps)
		if err != nil || !res.Sucesso {
			t.Fatalf("%s (%s): esperava paridade ok: err=%v res=%+v", p.slug, p.arquivo, err, res)
		}
	}
	novo := filepath.Join(dir, caminhos.DirProcessos, "beta", caminhos.ArquivoProcessoDAG)
	if err := os.WriteFile(novo, []byte(`{"name":"diverge","nodes":[{"id":"01-alfa"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	res, _ := fechar.Executar(nil, fechar.Request{Slug: "beta"}, deps)
	if res.Sucesso || !strings.Contains(res.Etapas[3].Detalhe, "divergente") {
		t.Fatalf("esperava divergência detectada no nome novo: %+v", res.Etapas[3])
	}
}

func testarParidadeSHA256(t *testing.T, pasta, arquivo string) {
	dir := t.TempDir()
	semearDocs(t, dir, "teste")
	dagsDir := filepath.Join(dir, "docs", "dags")
	procDir := filepath.Join(dir, pasta, "teste")
	_ = os.MkdirAll(dagsDir, 0o755)
	_ = os.MkdirAll(procDir, 0o755)

	dagIso := filepath.Join(dagsDir, "2026-11-09_DAG_TESTE.json")
	procDag := filepath.Join(procDir, arquivo)

	// Caso 1: Idênticos -> Sucesso
	dagReal := []byte(`{"name":"teste","version":"1.0.0","nodes":[{"id":"01-alfa"}]}`)
	_ = os.WriteFile(dagIso, dagReal, 0o644)
	_ = os.WriteFile(procDag, dagReal, 0o644)

	deps := fechar.Deps{RepoDir: dir, SkipBuild: true}
	req := fechar.Request{Slug: "teste"}
	res, err := fechar.Executar(nil, req, deps)
	if err != nil || !res.Sucesso {
		t.Fatalf("esperava sucesso quando DAGs são idênticos: err=%v", err)
	}

	// Caso 2: Divergentes -> Erro na etapa 4
	_ = os.WriteFile(procDag, []byte(`{"name":"teste-divergente","nodes":[{"id":"01-alfa"}]}`), 0o644)
	resDiv, _ := fechar.Executar(nil, req, deps)
	if resDiv.Sucesso {
		t.Fatalf("esperava falha quando DAGs divergem em SHA-256")
	}
	if resDiv.Etapas[3].Status != "error" || !strings.Contains(resDiv.Etapas[3].Detalhe, "divergente") {
		t.Errorf("etapa 4 deveria falhar com erro de divergência: %+v", resDiv.Etapas[3])
	}
}

func rodarFecharSimples(t *testing.T, repo, slug string) *fechar.Result {
	t.Helper()
	semearDocs(t, repo, slug)
	deps := fechar.Deps{RepoDir: repo, VaultDir: repo, Home: repo, SkipBuild: true,
		Indexer: func(string) error { return nil }, Portoes: func(string) error { return nil }}
	res, err := fechar.Executar(context.Background(), fechar.Request{Slug: slug}, deps)
	if err != nil || !res.Sucesso {
		t.Fatalf("esperava sucesso: err=%v res=%+v", err, res)
	}
	return res
}

func prepararProcessoCanonico(t *testing.T, repo, canon string) string {
	t.Helper()
	proc := filepath.Join(repo, caminhos.DirLoops, canon)
	if err := os.MkdirAll(proc, 0o755); err != nil {
		t.Fatal(err)
	}
	dag := `{"name":"` + canon + `","version":"1.0.0","nodes":[{"id":"01-alfa"},{"id":"02-beta"}]}`
	if err := os.WriteFile(filepath.Join(proc, caminhos.ArquivoDAG), []byte(dag), 0o644); err != nil {
		t.Fatal(err)
	}
	return proc
}

func TestFecharPrefereDiretoriosCanonicos(t *testing.T) {
	repo := t.TempDir()
	prepararProcessoCanonico(t, repo, "legoloop-foo")
	livro := filepath.Join(repo, "docs", "books", "legoloop-foo")
	if err := os.MkdirAll(livro, 0o755); err != nil {
		t.Fatal(err)
	}
	rodarFecharSimples(t, repo, "foo")
	for _, dup := range []string{filepath.Join(repo, caminhos.DirLoops, "foo"), filepath.Join(repo, "docs", "books", "foo")} {
		if _, err := os.Stat(dup); err == nil {
			t.Errorf("pasta duplicada criada: %s", dup)
		}
	}
	if _, err := os.Stat(filepath.Join(livro, "livro.json")); err != nil {
		t.Errorf("livro.json deveria estar em docs/books/legoloop-foo: %v", err)
	}
}

func TestFecharGeraQuartetoPadrao(t *testing.T) {
	repo := t.TempDir()
	proc := prepararProcessoCanonico(t, repo, "legoloop-bar")
	rodarFecharSimples(t, repo, "bar")
	for _, nome := range []string{caminhos.ArquivoDAG, "triagem.json", "memoria_execucao.json"} {
		if _, err := os.Stat(filepath.Join(proc, nome)); err != nil {
			t.Errorf("arquivo padrão ausente %s: %v", nome, err)
		}
	}
	if _, err := os.Stat(filepath.Join(proc, "legoloop-bar"+caminhos.SufixoLoop)); !os.IsNotExist(err) {
		t.Errorf("o espelho .loop.json não pode mais ser criado: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(proc, "triagem.json"))
	if err != nil || !strings.Contains(string(raw), "02-beta") {
		t.Errorf("triagem.json deve conter IDs dos nós do DAG: err=%v raw=%s", err, raw)
	}
}

func TestFecharDryRunNaoGravaQuarteto(t *testing.T) {
	repo := t.TempDir()
	deps := fechar.Deps{RepoDir: repo, VaultDir: repo, Home: repo, SkipBuild: true,
		Indexer: func(string) error { return nil }, Portoes: func(string) error { return nil }}
	if _, err := fechar.Executar(context.Background(), fechar.Request{Slug: "baz", DryRun: true}, deps); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repo, caminhos.DirLoops)); err == nil {
		t.Errorf("dry-run não deve criar a pasta de Processos")
	}
}

func TestFecharPrefereCanonicoSobreSlugQuandoAmbosExistem(t *testing.T) {
	repo := t.TempDir()
	canon := prepararProcessoCanonico(t, repo, "legoloop-qux")
	if err := os.MkdirAll(filepath.Join(repo, caminhos.DirLoops, "qux"), 0o755); err != nil {
		t.Fatal(err)
	}
	rodarFecharSimples(t, repo, "qux")
	if _, err := os.Stat(filepath.Join(canon, "triagem.json")); err != nil {
		t.Errorf("pasta canônica deveria receber os arquivos padrão: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, caminhos.DirLoops, "qux", caminhos.ArquivoDAG)); err == nil {
		t.Errorf("pasta do slug não deveria receber o DAG")
	}
}

func TestFecharPecaDeUmNoNaoGeraTriagem(t *testing.T) {
	repo := t.TempDir()
	proc := prepararProcessoCanonico(t, repo, "legoloop-uno")
	dag := `{"name":"legoloop-uno","version":"1.0.0","nodes":[{"id":"e01"}]}`
	if err := os.WriteFile(filepath.Join(proc, caminhos.ArquivoDAG), []byte(dag), 0o644); err != nil {
		t.Fatal(err)
	}
	rodarFecharSimples(t, repo, "uno")
	if _, err := os.Stat(filepath.Join(proc, "triagem.json")); err == nil {
		t.Error("um nó só não tem triagem")
	}
}

func TestFecharCompletaATriagemQueNaoCobreTodosOsNos(t *testing.T) {
	repo := t.TempDir()
	proc := prepararProcessoCanonico(t, repo, "legoloop-duo")
	velha := triagem.GerarTriagemPadrao("legoloop-duo", []string{"01-alfa"})
	b, _ := json.Marshal(velha)
	if err := os.WriteFile(filepath.Join(proc, "triagem.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	rodarFecharSimples(t, repo, "duo")
	if err := triagem.Conferir(proc, []string{"01-alfa", "02-beta"}); err != nil {
		t.Errorf("o fechar deve deixar a triagem cobrindo todos os nós: %v", err)
	}
}

// semearProcesso grava o processo do slug no cofre com a natureza pedida.
func semearProcesso(t *testing.T, repo, slug, natureza string) {
	t.Helper()
	p := processo.Processo{
		Nome:     "legoloop-" + slug,
		Versao:   "1.0.0",
		Natureza: natureza,
		Objetivo: "processo de teste da natureza " + natureza,
		CriadoEm: time.Now(),
		Etapas: []processo.Etapa{{
			ID:     "e01",
			Titulo: "Etapa única",
			Tipo:   "construir",
			Gate:   "file_exists('out/e01.md')",
			Produz: processo.Artefato{Caminho: "out/e01.md", Tipo: "markdown"},
		}},
	}
	data, err := p.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	escrever(t, processo.CaminhoNoCofre(repo, p.Nome), string(data))
}

// etapaDe devolve a etapa num (1-based) do resultado.
func etapaDe(t *testing.T, res *fechar.Result, num int) fechar.EtapaStatus {
	t.Helper()
	for _, e := range res.Etapas {
		if e.Etapa == num {
			return e
		}
	}
	t.Fatalf("etapa %d ausente do resultado: %+v", num, res.Etapas)
	return fechar.EtapaStatus{}
}

// TestFecharPorNatureza prova a AC-E03.3: a etapa 5 (Recompilação do Binário) só vale
// para naturezas Go-cêntricas. Processo de natureza midia omite a etapa 5 como não
// aplicável e a cadeia restante executa com sucesso.
func TestFecharPorNatureza(t *testing.T) {
	t.Run("midia com dry-run omite a recompilacao", func(t *testing.T) {
		repo := t.TempDir()
		semearDocs(t, repo, "video-motion")
		semearDAG(t, repo, "video-motion")
		semearProcesso(t, repo, "video-motion", "midia")
		res, err := fechar.Executar(context.Background(),
			fechar.Request{Slug: "video-motion", DryRun: true}, depsFechar(repo))
		if err != nil {
			t.Fatal(err)
		}
		e5 := etapaDe(t, res, 5)
		if e5.Status != "skipped" {
			t.Errorf("etapa 5 deveria ser skipped para natureza midia: %+v", e5)
		}
		if !strings.Contains(e5.Detalhe, "não aplicável") {
			t.Errorf("detalhe da etapa 5 deve dizer não aplicável: %q", e5.Detalhe)
		}
		if !res.Sucesso {
			t.Errorf("cadeia restante deve executar com sucesso: %+v", res.Etapas)
		}
		for _, n := range []int{1, 2, 3, 4, 6, 7, 8, 9} {
			if e := etapaDe(t, res, n); e.Status != "ok" {
				t.Errorf("etapa %d deveria executar normalmente: %+v", n, e)
			}
		}
	})

	t.Run("midia sem dry-run nao chama o Rebuilder", func(t *testing.T) {
		repo := t.TempDir()
		semearDocs(t, repo, "video-motion")
		semearDAG(t, repo, "video-motion")
		semearProcesso(t, repo, "video-motion", "midia")
		rebuilds := 0
		deps := depsFechar(repo)
		deps.SkipBuild = false
		deps.Rebuilder = func(string) error { rebuilds++; return nil }
		res, err := fechar.Executar(context.Background(),
			fechar.Request{Slug: "video-motion"}, deps)
		if err != nil {
			t.Fatal(err)
		}
		if rebuilds != 0 {
			t.Errorf("Rebuilder não pode rodar para natureza midia: %d chamadas", rebuilds)
		}
		if e5 := etapaDe(t, res, 5); e5.Status != "skipped" {
			t.Errorf("etapa 5 deveria ser skipped para natureza midia: %+v", e5)
		}
	})

	t.Run("software mantem a etapa 5", func(t *testing.T) {
		repo := t.TempDir()
		semearDocs(t, repo, "video-motion")
		semearDAG(t, repo, "video-motion")
		semearProcesso(t, repo, "video-motion", "software")
		res, err := fechar.Executar(context.Background(),
			fechar.Request{Slug: "video-motion", DryRun: true}, depsFechar(repo))
		if err != nil {
			t.Fatal(err)
		}
		if e5 := etapaDe(t, res, 5); e5.Status != "ok" {
			t.Errorf("etapa 5 deve continuar valendo para natureza software: %+v", e5)
		}
	})
}
