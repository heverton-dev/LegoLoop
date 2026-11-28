package distribuicao

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/redc"
)

func textosExemploValido() TextosCanal {
	return TextosCanal{
		Nome:     "Tech em Jogo",
		Handle:   "techemjogooficial",
		Promessa: "Como a inteligencia artificial muda o trabalho real, com dados e fontes.",
		Tagline:  "IA e tecnologia com dados reais.",
		Tom:      "Direto, investigativo e sem hype.",
		Pilares:  []string{"Impacto real da IA", "Engenharia de software", "Dados do mercado"},
		IdentidadeVisual: IdentidadeCanal{
			Simbolo: "Circuito dourado dentro de pilula escura",
			Paleta:  map[string]string{"fundo": "#0B0E1C", "ouro": "#FFD24A", "texto": "#FFFFFF"},
			Fontes:  map[string]string{"titulo": "Bebas Neue", "corpo": "Inter"},
		},
		DisponibilidadeNome: map[string]string{"youtube": "a confirmar no cadastro com @techemjogooficial"},
		Redes: map[string]RedeCanal{
			"youtube": {
				NomeExibicao: "Tech em Jogo",
				Handle:       "techemjogooficial",
				Descricao:    "Analises sobre IA e tecnologia com fontes publicas. Imagens e narracao com apoio de IA.",
				Playlists:    []string{"IA no Trabalho", "Bastidores do Codigo"},
			},
			"instagram": {
				NomeExibicao: "Tech em Jogo",
				Handle:       "techemjogooficial",
				Bio:          "IA e tecnologia com numeros e fontes reais.",
			},
			"tiktok": {
				NomeExibicao: "Tech em Jogo",
				Handle:       "techemjogooficial",
				Bio:          "IA e tecnologia com fontes reais.",
			},
			"linkedin": {
				NomeExibicao: "Tech em Jogo",
				Tagline:      "Analises sobre inteligencia artificial e engenharia com dados publicos.",
				Sobre:        "Conteudo informativo sobre tecnologia e IA baseado em fontes nivel A e B.",
			},
		},
		PrimeirasPautas: []PautaCanal{{Titulo: "IA e a semana de 90 horas", Status: "em producao"}},
	}
}

func TestGerarTextosCanalEscreveColarMDComLimitesEAvisos(t *testing.T) {
	dir := t.TempDir()
	tc := textosExemploValido()
	if err := GravarTextosCanal(dir, tc); err != nil {
		t.Fatalf("GravarTextosCanal: %v", err)
	}
	rel, err := GerarTextosCanal(dir)
	if err != nil {
		t.Fatalf("GerarTextosCanal: %v", err)
	}
	if rel.Excedidos != 0 {
		t.Fatalf("nenhum limite deveria exceder: %+v", rel)
	}
	colar, err := os.ReadFile(filepath.Join(dir, "colar.md"))
	if err != nil {
		t.Fatalf("colar.md ausente: %v", err)
	}
	s := string(colar)
	for _, trecho := range []string{"# Kit de canal: Tech em Jogo (@techemjogooficial)", "## YouTube", "## Instagram", "## TikTok", "⚠"} {
		if !strings.Contains(s, trecho) {
			t.Fatalf("colar.md sem %q:\n%s", trecho, s)
		}
	}
}

func TestGerarTextosCanalDetectaLimiteExcedido(t *testing.T) {
	dir := t.TempDir()
	tc := textosExemploValido()
	tk := tc.Redes["tiktok"]
	tk.Bio = strings.Repeat("muito longo ", 12) // > 80 chars
	tc.Redes["tiktok"] = tk
	if err := GravarTextosCanal(dir, tc); err != nil {
		t.Fatalf("GravarTextosCanal: %v", err)
	}
	rel, err := GerarTextosCanal(dir)
	if err == nil || rel.Excedidos == 0 {
		t.Fatalf("esperava erro de limite excedido, obteve rel=%+v err=%v", rel, err)
	}
	colar, _ := os.ReadFile(filepath.Join(dir, "colar.md"))
	if !strings.Contains(string(colar), "**EXCEDE O LIMITE**") {
		t.Fatalf("colar.md deveria marcar EXCEDE O LIMITE")
	}
}

func TestGerarTextosCanalDeDNAPreencheQuatroRedesComTransparenciaIA(t *testing.T) {
	dna := redc.ChannelDNA{
		CanalSlug:   "canal-ref",
		TomDeVoz:    "direto e investigativo",
		RitmoCortes: "16 cortes/min",
		TopGanchos:  []string{"Por que a conta nao fecha?"},
	}
	tokens := redc.BrandTokens{
		NomeMarca:     "Tech em Jogo",
		CorPrimaria:   "#FFD24A",
		CorSecundaria: "#3DDC84",
		FonteTitulo:   "Bebas Neue",
		FonteCorpo:    "Inter",
	}
	tc, err := GerarTextosCanalDeDNA(dna, tokens, "techemjogooficial", "tecnologia e inteligencia artificial")
	if err != nil {
		t.Fatalf("GerarTextosCanalDeDNA: %v", err)
	}
	dir := t.TempDir()
	if err := GravarTextosCanal(dir, tc); err != nil {
		t.Fatalf("GravarTextosCanal: %v", err)
	}
	rel, err := GerarTextosCanal(dir)
	if err != nil || rel.Excedidos != 0 {
		t.Fatalf("textos gerados de DNA devem caber nos limites: rel=%+v err=%v", rel, err)
	}
}

func TestGerarImagensCanalInvocaScriptEConfereSaidas(t *testing.T) {
	dir := t.TempDir()
	if _, err := GerarImagensCanal(context.Background(), &execCanalFalso{canalDir: dir}, dir, ""); err == nil {
		t.Fatalf("esperava erro sem textos.json")
	}
	if err := GravarTextosCanal(dir, textosExemploValido()); err != nil {
		t.Fatalf("GravarTextosCanal: %v", err)
	}
	if _, err := GerarImagensCanal(context.Background(), &execCanalFalso{canalDir: dir}, dir, ""); err == nil {
		t.Fatalf("esperava erro sem base/perfil ou base/banner")
	}
	tocar(t, dir, "base/perfil.png")
	tocar(t, dir, "base/banner.png")
	logoDest := filepath.Join(dir, "marca", "logo.png")
	ex := &execCanalFalso{canalDir: dir}
	res, err := GerarImagensCanal(context.Background(), ex, dir, logoDest)
	if err != nil {
		t.Fatalf("GerarImagensCanal: %v", err)
	}
	if len(res.Arquivos) < 4 || res.LogoSincronizado != logoDest {
		t.Fatalf("esperava imagens e logo sincronizado, obteve: %+v", res)
	}
	if _, err := GerarTextosCanalDeDNA(redc.ChannelDNA{}, redc.BrandTokens{}, "", ""); err == nil {
		t.Fatalf("esperava erro em GerarTextosCanalDeDNA vazio")
	}
	if err := GravarTextosCanal(dir, TextosCanal{}); err == nil {
		t.Fatalf("esperava erro em TextosCanal vazio")
	}
}

type execCanalFalso struct {
	canalDir string
}

func (e *execCanalFalso) Executar(_ context.Context, _ string, _ ...string) error {
	out := filepath.Join(e.canalDir, "imagens")
	_ = os.MkdirAll(out, 0o755)
	for _, n := range []string{"perfil_1080.png", "youtube_perfil_800.png", "youtube_marca_dagua_150.png", "youtube_banner_2560x1440.png", "previa_perfil_98px_ampliada.png", "previa_banner_area_segura.png"} {
		if err := os.WriteFile(filepath.Join(out, n), []byte("png"), 0o644); err != nil {
			return err
		}
	}
	return nil
}
