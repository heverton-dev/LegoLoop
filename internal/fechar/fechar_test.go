package fechar_test

import (
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/fechar"
)

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
	conteudo := []byte(`{"name":"x","version":"1.0.0"}`)
	procs := []struct{ slug, pasta, arquivo string }{
		{"alfa", caminhos.DirLoops, caminhos.ArquivoDAG},
		{"beta", caminhos.DirProcessos, caminhos.ArquivoProcessoDAG},
	}
	for _, p := range procs {
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
	if err := os.WriteFile(novo, []byte(`{"name":"diverge"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	res, _ := fechar.Executar(nil, fechar.Request{Slug: "beta"}, deps)
	if res.Sucesso || !strings.Contains(res.Etapas[3].Detalhe, "divergente") {
		t.Fatalf("esperava divergência detectada no nome novo: %+v", res.Etapas[3])
	}
}

func testarParidadeSHA256(t *testing.T, pasta, arquivo string) {
	dir := t.TempDir()
	dagsDir := filepath.Join(dir, "docs", "dags")
	loopsDir := filepath.Join(dir, pasta, "teste")
	_ = os.MkdirAll(dagsDir, 0o755)
	_ = os.MkdirAll(loopsDir, 0o755)

	dagIso := filepath.Join(dagsDir, "2026-11-09_DAG_TESTE.json")
	loopDag := filepath.Join(loopsDir, arquivo)

	// Caso 1: Idênticos -> Sucesso
	_ = os.WriteFile(dagIso, []byte(`{"name":"teste","version":"1.0.0"}`), 0o644)
	_ = os.WriteFile(loopDag, []byte(`{"name":"teste","version":"1.0.0"}`), 0o644)

	deps := fechar.Deps{RepoDir: dir, SkipBuild: true}
	req := fechar.Request{Slug: "teste"}
	res, err := fechar.Executar(nil, req, deps)
	if err != nil || !res.Sucesso {
		t.Fatalf("esperava sucesso quando DAGs são idênticos: err=%v", err)
	}

	// Caso 2: Divergentes -> Erro na etapa 4
	_ = os.WriteFile(loopDag, []byte(`{"name":"teste-divergente"}`), 0o644)
	resDiv, _ := fechar.Executar(nil, req, deps)
	if resDiv.Sucesso {
		t.Fatalf("esperava falha quando DAGs divergem em SHA-256")
	}
	if resDiv.Etapas[3].Status != "error" || !strings.Contains(resDiv.Etapas[3].Detalhe, "divergente") {
		t.Errorf("etapa 4 deveria falhar com erro de divergência: %+v", resDiv.Etapas[3])
	}
}
