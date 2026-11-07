package fechar_test

import (
	"os"
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
