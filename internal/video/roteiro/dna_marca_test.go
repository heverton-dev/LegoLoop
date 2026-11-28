package roteiro

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/redc"
)

func TestAplicarChannelDNAGeraMarcaEPerfilSemPlagio(t *testing.T) {
	dna := redc.ChannelDNA{
		CanalSlug:    "fireship-ref",
		TomDeVoz:     "direto e analitico",
		Arquetipo:    "explicador-rapido",
		RitmoCortes:  "16 cortes/min",
		DesignTokens: map[string]string{"accent": "#00FFAA"},
		MotionTokens: map[string]string{"entrada": "pop", "camera": "in-panR"},
		TopGanchos:   []string{"A promessa que virou gargalo"},
	}
	tokens := redc.BrandTokens{
		NomeMarca:     "Codigo Real",
		CorPrimaria:   "#FFD24A",
		CorSecundaria: "#3DDC84",
		FonteTitulo:   "Bebas Neue",
		FonteCorpo:    "Inter",
		TomPreferido:  "direto com fontes",
	}
	base := Marca{
		SchemaVersion: 1,
		Formatos:      []string{"16:9", "9:16", "1:1"},
		Voz:           VozMarca{Motor: "edge-tts", Voz: "pt-BR-ThalitaNeural", Velocidade: "-5%"},
	}
	dir := t.TempDir()
	marca, perfil, err := GravarBrandSwapDNA(dir, dna, tokens, base)
	if err != nil {
		t.Fatalf("GravarBrandSwapDNA falhou: %v", err)
	}
	if marca.Nome != "Codigo Real" || marca.Cores.Foco != "#FFD24A" || marca.Cores.Foco2 != "#3DDC84" {
		t.Fatalf("cores/nome da marca incorretos: %+v", marca)
	}
	if marca.Fontes.Titulo.Familia != "Bebas Neue" || marca.Fontes.Corpo.Familia != "Inter" {
		t.Fatalf("fontes da marca incorretas: %+v", marca.Fontes)
	}
	if !marca.Vinhetas.Abertura.Ativo || marca.Vinhetas.Abertura.Texto != "CODIGO REAL" {
		t.Fatalf("vinheta de abertura incorreta: %+v", marca.Vinhetas)
	}
	if perfil.CanalReferencia != "fireship-ref" || perfil.RitmoCortes != "16 cortes/min" {
		t.Fatalf("perfil de estilo incorreto: %+v", perfil)
	}
	for _, arq := range []string{"marca.json", "perfil_estilo.json"} {
		if _, err := os.Stat(filepath.Join(dir, arq)); err != nil {
			t.Fatalf("arquivo %s nao gravado: %v", arq, err)
		}
	}
}

func TestAplicarChannelDNARecusaEntradaInvalida(t *testing.T) {
	_, _, err := AplicarChannelDNA(redc.ChannelDNA{}, redc.BrandTokens{NomeMarca: "X"}, Marca{})
	if err == nil || !strings.Contains(err.Error(), "canal_slug") {
		t.Fatalf("esperava erro de canal_slug vazio, obteve: %v", err)
	}
	_, _, err = AplicarChannelDNA(redc.ChannelDNA{CanalSlug: "ok"}, redc.BrandTokens{}, Marca{})
	if err == nil || !strings.Contains(err.Error(), "nome_marca") {
		t.Fatalf("esperava erro de nome_marca vazio, obteve: %v", err)
	}
	m, p, err := AplicarChannelDNA(
		redc.ChannelDNA{CanalSlug: "ok", TomDeVoz: "calmo"},
		redc.BrandTokens{NomeMarca: "Marca Pura"},
		Marca{Logo: LogoMarca{Arquivo: "custom.png"}, Vinhetas: VinhetasMarca{Abertura: VinhetaCfg{Texto: "CUSTOM"}}},
	)
	if err != nil || m.Logo.Arquivo != "custom.png" || m.Vinhetas.Abertura.Texto != "CUSTOM" || p.TomDeVoz != "calmo" {
		t.Fatalf("defaults/preservacao falharam: m=%+v p=%+v err=%v", m, p, err)
	}
	if _, _, err := GravarBrandSwapDNA(t.TempDir(), redc.ChannelDNA{}, redc.BrandTokens{}, Marca{}); err == nil {
		t.Fatalf("esperava erro em GravarBrandSwapDNA invalido")
	}
}
