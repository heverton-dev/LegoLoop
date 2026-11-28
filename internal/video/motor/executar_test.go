package motor

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/ativos"
	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
	"github.com/heverton-dev/LegoLoop/internal/video/render"
)

// wavMono24k monta um WAV PCM mono de 24 kHz com o numero de amostras pedido.
func wavMono24k(amostras int) []byte {
	dados := bytes.Repeat([]byte{1, 0}, amostras)
	var b bytes.Buffer
	b.WriteString("RIFF")
	_ = binary.Write(&b, binary.LittleEndian, uint32(36+len(dados)))
	b.WriteString("WAVEfmt ")
	for _, v := range []any{uint32(16), uint16(1), uint16(1), uint32(24000), uint32(48000), uint16(2), uint16(16)} {
		_ = binary.Write(&b, binary.LittleEndian, v)
	}
	b.WriteString("data")
	_ = binary.Write(&b, binary.LittleEndian, uint32(len(dados)))
	b.Write(dados)
	return b.Bytes()
}

func gravar(t *testing.T, caminho string, v any) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	dados, ok := v.([]byte)
	if !ok {
		var err error
		if dados, err = json.Marshal(v); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(caminho, dados, 0o644); err != nil {
		t.Fatal(err)
	}
}

// projetoMotor monta um projeto minimo: 1 cena, 1 personagem, 1 fundo codigo-desenhado e voz de 3 s.
func projetoMotor(t *testing.T) (string, inventario.Store) {
	t.Helper()
	dir := t.TempDir()
	st := inventario.Store{Raiz: t.TempDir()}
	est := inventario.Estilo{SchemaVersion: 1, ID: "molde-teste", Nome: "Molde", Status: "rascunho", Versao: "1.0.0",
		Renderizador: "molde-pilula", PromptMestre: "p", Paleta: []string{"#fff"},
		Referencias: []inventario.Referencia{{Origem: "teste"}}}
	if err := st.NovoEstilo(est); err != nil {
		t.Fatal(err)
	}
	est.Status = "validado"
	gravar(t, filepath.Join(st.Raiz, "inventario", "video", "estilos", est.ID, "estilo.json"), est)
	item := inventario.Item{SchemaVersion: 1, Aprovado: true, ID: "escritorio", Tipo: "cenario", Estilo: "molde-teste",
		Parametros: map[string]any{"kind": "office"}}
	if err := st.GravarItem(item); err != nil {
		t.Fatal(err)
	}
	gravar(t, filepath.Join(dir, "projeto.json"), map[string]any{"estilo_visual": "molde-teste"})
	gravar(t, filepath.Join(dir, "cenas.json"), ativos.Storyboard{SchemaVersion: 2, Cenas: []ativos.Cena{
		{IndiceLinha: 0, CenarioID: "escritorio", ElementosMoveis: []ativos.ElementoMovel{
			{ID: "m1", Tipo: "personagem", PersonagemID: "ana", Fala: true}}}}})
	gravar(t, filepath.Join(dir, "elenco.json"), ElencoPonte{Personagens: []PersonagemPonte{{ID: "ana", Cast: map[string]any{"skin": "#E6BE8E"}}}})
	gravar(t, filepath.Join(dir, "data", "timings.json"), render.Timings{SchemaVersion: 1, Provedor: "t", Voz: "v",
		Linhas: []render.LinhaTempo{{IndiceLinha: 0, InicioS: 0, FimS: 3}}})
	gravar(t, filepath.Join(dir, "audio", "voice.wav"), wavMono24k(72000))
	return dir, st
}

func TestCarregarEntradaLeContratosDoProjeto(t *testing.T) {
	dir, st := projetoMotor(t)
	e, err := CarregarEntrada(dir, st)
	if err != nil {
		t.Fatalf("CarregarEntrada: %v", err)
	}
	if len(e.Cenas.Cenas) != 1 || len(e.Fundos) != 1 || e.Estilo.ID != "molde-teste" || e.FPS != 30 {
		t.Fatalf("entrada incompleta: %+v", e)
	}
	if e.DuracaoVozS < 2.99 || e.DuracaoVozS > 3.01 || e.Audio != "audio/soundtrack.wav" {
		t.Fatalf("voz %.3f s, audio %q", e.DuracaoVozS, e.Audio)
	}
}

func TestCarregarEntradaRecusaProjetoIncompleto(t *testing.T) {
	casos := map[string]string{
		"sem projeto.json": "projeto.json",
		"sem cenas.json":   "cenas.json",
		"sem elenco.json":  "elenco.json",
		"sem timings.json": "data/timings.json",
		"sem voice.wav":    "audio/voice.wav",
	}
	for nome, arq := range casos {
		t.Run(nome, func(t *testing.T) {
			dir, st := projetoMotor(t)
			if err := os.Remove(filepath.Join(dir, arq)); err != nil {
				t.Fatal(err)
			}
			if _, err := CarregarEntrada(dir, st); err == nil {
				t.Fatal("quer erro, obteve nil")
			}
		})
	}
	t.Run("estilo vazio", func(t *testing.T) {
		dir, st := projetoMotor(t)
		gravar(t, filepath.Join(dir, "projeto.json"), map[string]any{"estilo_visual": ""})
		if _, err := CarregarEntrada(dir, st); !errors.Is(err, ErrEstiloAusente) {
			t.Fatalf("quer ErrEstiloAusente, obteve %v", err)
		}
	})
	t.Run("estilo fora do inventario", func(t *testing.T) {
		dir, _ := projetoMotor(t)
		if _, err := CarregarEntrada(dir, inventario.Store{Raiz: t.TempDir()}); err == nil {
			t.Fatal("quer erro, obteve nil")
		}
	})
	t.Run("json invalido", func(t *testing.T) {
		dir, st := projetoMotor(t)
		gravar(t, filepath.Join(dir, "cenas.json"), []byte("{"))
		if _, err := CarregarEntrada(dir, st); err == nil {
			t.Fatal("quer erro, obteve nil")
		}
	})
}

type checadorOK struct{ saida string }

func (c checadorOK) Checar(string) (string, error) { return c.saida, nil }

type motorFalso struct {
	chamadas int
	err      error
}

func (m *motorFalso) Renderizar(context.Context, string, float64, int) error {
	m.chamadas++
	return m.err
}

// Executar grava timeline.json e, so com CHECK OK, chama o motor; sem CHECK OK nao chama.
func TestExecutarGeraTimelineEChamaOMotorSoComCheckOK(t *testing.T) {
	dir, st := projetoMotor(t)
	m := &motorFalso{}
	_, err := Executar(context.Background(), dir, st, checadorOK{saida: "nada"}, m)
	if !errors.Is(err, ErrCheckTimeline) || m.chamadas != 0 {
		t.Fatalf("sem CHECK OK: err=%v chamadas=%d", err, m.chamadas)
	}
	if _, err := os.Stat(filepath.Join(dir, "timeline.json")); err != nil {
		t.Fatalf("timeline.json nao gravada: %v", err)
	}
}

func TestExecutarPropagaFalhaDeCarga(t *testing.T) {
	dir, st := projetoMotor(t)
	_ = os.Remove(filepath.Join(dir, "cenas.json"))
	m := &motorFalso{}
	if _, err := Executar(context.Background(), dir, st, checadorOK{saida: "CHECK OK"}, m); err == nil || m.chamadas != 0 {
		t.Fatalf("quer erro de carga sem chamar o motor: err=%v chamadas=%d", err, m.chamadas)
	}
}

type execFalso struct{ dirs []string }

func (e *execFalso) Executar(dir string) error { e.dirs = append(e.dirs, dir); return nil }

// RenderMotor adapta Renderizar ao contrato render.MotorRender.
func TestRenderMotorAdaptaORenderizar(t *testing.T) {
	dir := t.TempDir()
	gravar(t, filepath.Join(dir, ArquivoRender), []byte("mp4"))
	ex := &execFalso{}
	bom := Midia{Largura: 1080, Altura: 1920, FPS: 30, TemAudio: true, DuracaoS: 3}
	var mr render.MotorRender = RenderMotor{Ex: ex, Sonda: &sondaFalsa{midia: bom}}
	if err := mr.Renderizar(context.Background(), dir, 3, 30); err != nil || len(ex.dirs) != 1 || ex.dirs[0] != dir {
		t.Fatalf("err=%v dirs=%v", err, ex.dirs)
	}
	ruim := Midia{Largura: 1920, Altura: 1080, FPS: 30, TemAudio: true, DuracaoS: 3}
	mr = RenderMotor{Ex: ex, Sonda: &sondaFalsa{midia: ruim}}
	if err := mr.Renderizar(context.Background(), dir, 3, 30); err == nil {
		t.Fatal("quer erro para 16:9, obteve nil")
	}
}
