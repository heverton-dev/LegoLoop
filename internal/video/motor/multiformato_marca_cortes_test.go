package motor

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/render"
	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

func TestRenderMultiformatoAceita16x9_9x16_1x1(t *testing.T) {
	casos := []struct {
		fmt  string
		w, h int
	}{
		{"16:9", 1920, 1080},
		{"9:16", 1080, 1920},
		{"1:1", 1080, 1080},
	}
	for _, c := range casos {
		dir := t.TempDir()
		s := &sondaFalsa{midia: Midia{Largura: c.w, Altura: c.h, FPS: 30, TemAudio: true, DuracaoS: 6.71}}
		m, err := Renderizar(PedidoRender{Dir: dir, Formato: c.fmt, DuracaoVozS: 6.71, FPS: 30}, &executorFalso{}, s)
		if err != nil || m.Largura != c.w || m.Altura != c.h {
			t.Fatalf("formato %s falhou: m=%+v err=%v", c.fmt, m, err)
		}
	}
	sErr := &sondaFalsa{midia: Midia{Largura: 1080, Altura: 1920, FPS: 30, TemAudio: true, DuracaoS: 6.71}}
	if _, err := Renderizar(PedidoRender{Dir: t.TempDir(), Formato: "16:9", DuracaoVozS: 6.71, FPS: 30}, &executorFalso{}, sErr); !errors.Is(err, ErrRenderFormato) {
		t.Fatalf("esperava ErrRenderFormato quando 16:9 recebe 1080x1920, obteve: %v", err)
	}
	if _, _, err := DimensoesFormato("4:3"); !errors.Is(err, ErrRenderFormato) {
		t.Fatalf("esperava ErrRenderFormato para 4:3, obteve: %v", err)
	}
}

type cmdGravadorFalso struct {
	comandos [][]string
}

func (c *cmdGravadorFalso) Rodar(_ context.Context, nome string, args ...string) ([]byte, error) {
	c.comandos = append(c.comandos, append([]string{nome}, args...))
	if nome == "ffprobe" {
		return []byte(`{"streams":[{"width":1920,"height":1080}],"format":{"duration":"120.0"}}`), nil
	}
	// Grava o ultimo argumento se parecer caminho .mp4 ou .wav
	if len(args) > 0 {
		ult := args[len(args)-1]
		if strings.HasSuffix(ult, ".mp4") || strings.HasSuffix(ult, ".wav") {
			_ = os.MkdirAll(filepath.Dir(ult), 0o755)
			_ = os.WriteFile(ult, []byte("midia"), 0o644)
		}
		for i := 0; i+1 < len(args); i++ {
			if args[i] == "--out" {
				_ = os.MkdirAll(filepath.Dir(args[i+1]), 0o755)
				_ = os.WriteFile(args[i+1], []byte("midia"), 0o644)
			}
		}
	}
	return []byte("ok"), nil
}

func TestAplicarMarcaGeraLogoProvisorioVinhetasEConcatena(t *testing.T) {
	dir := t.TempDir()
	entrada := filepath.Join(dir, "render.mp4")
	saida := filepath.Join(dir, "renders", "v1.mp4")
	_ = os.WriteFile(entrada, []byte("mp4"), 0o644)
	marca := roteiro.Marca{
		SchemaVersion: 1,
		Nome:          "Tech em Jogo",
		Formatos:      []string{"16:9"},
		Cores:         roteiro.CoresMarca{Foco: "#FFD24A", Contorno: "#111111"},
		Logo:          roteiro.LogoMarca{Arquivo: "marca/logo.png", Posicao: "canto_superior_direito", Opacidade: 0.85, LarguraPct: 8},
		Vinhetas: roteiro.VinhetasMarca{
			Abertura:   roteiro.VinhetaCfg{Ativo: true, DuracaoS: 1.5, Texto: "TECH EM JOGO"},
			Fechamento: roteiro.VinhetaCfg{Ativo: true, DuracaoS: 3.0, Texto: "Inscreva-se"},
		},
	}
	b, _ := json.MarshalIndent(marca, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "marca.json"), b, 0o644)
	cmd := &cmdGravadorFalso{}
	rel, err := AplicarMarca(context.Background(), dir, entrada, saida, "motor-local/render.py", cmd)
	if err != nil {
		t.Fatalf("AplicarMarca: %v", err)
	}
	if rel.Formato != "16:9" || !rel.AberturaAtiva || !rel.FechamentoAtivo {
		t.Fatalf("relatorio de marca inesperado: %+v", rel)
	}
	if _, err := os.Stat(filepath.Join(dir, "marca", "logo.png")); err != nil {
		t.Fatalf("logo provisorio nao criado: %v", err)
	}
	if _, err := os.Stat(saida); err != nil {
		t.Fatalf("video de saida com marca nao criado: %v", err)
	}
}

func TestPrepararERenderizarCortesVerticais9x16(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, "audio"), 0o755)
	_ = os.MkdirAll(filepath.Join(dir, "data"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, "audio", "voice.wav"), []byte("wav"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "audio", "soundtrack.wav"), []byte("wav"), 0o644)
	tl := Timeline{
		FPS:   30,
		Audio: "audio/soundtrack.wav",
		Cast:  map[string]map[string]any{"dev": {"skin": "#E6BE8E"}},
		Shots: []Corte{
			{ID: "l000", Type: "scene", Lines: [2]int{0, 0}, Set: map[string]any{"kind": "server"}, Duration: 10.0},
			{ID: "l001", Type: "scene", Lines: [2]int{1, 1}, Set: map[string]any{"kind": "office"}, Duration: 12.0},
			{ID: "l002", Type: "scene", Lines: [2]int{2, 2}, Set: map[string]any{"kind": "boardroom"}, Duration: 11.0},
		},
	}
	btl, _ := json.MarshalIndent(tl, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "timeline.json"), btl, 0o644)
	tm := render.Timings{SchemaVersion: 1, Linhas: []render.LinhaTempo{
		{IndiceLinha: 0, InicioS: 0, FimS: 10},
		{IndiceLinha: 1, InicioS: 10, FimS: 22},
		{IndiceLinha: 2, InicioS: 22, FimS: 33},
	}}
	btm, _ := json.MarshalIndent(tm, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "data", "timings.json"), btm, 0o644)

	sugeridos := SugerirCortes(tl)
	if len(sugeridos) == 0 {
		t.Fatalf("SugerirCortes deveria sugerir ao menos 1 corte de 15..59s")
	}
	cmd := &cmdGravadorFalso{}
	resumos, err := PrepararCortes(context.Background(), dir, EspecificacaoCortes{Cortes: sugeridos}, cmd)
	if err != nil || len(resumos) == 0 {
		t.Fatalf("PrepararCortes: resumos=%+v err=%v", resumos, err)
	}
	gerados, err := RenderizarCortes(context.Background(), dir, resumos, "motor-local/render.py", cmd)
	if err != nil || len(gerados) != len(resumos) {
		t.Fatalf("RenderizarCortes: gerados=%+v err=%v", gerados, err)
	}
}
