package render

import (
	"context"
	"errors"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func existe(dir, rel string) bool {
	_, err := os.Stat(filepath.Join(dir, filepath.FromSlash(rel)))
	return err == nil
}

func TestAtivosSonorosSintetizaOQueFalta(t *testing.T) {
	dir := t.TempDir()
	ex := &execFalso{}
	if err := GarantirAtivosSonoros(context.Background(), dir, ConfigPadrao(), ex); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"audio/bgm.wav", "audio/sfx/troca_cena.wav", "audio/sfx/movimento_objeto.wav", "audio/sfx/destaque.wav"} {
		if !existe(dir, rel) {
			t.Errorf("%s deveria existir", rel)
		}
	}
	if len(ex.comandos) != 4 {
		t.Fatalf("esperava 4 chamadas ao ffmpeg, obteve %d", len(ex.comandos))
	}
	juntos := ""
	for _, c := range ex.comandos {
		if c[0] != "ffmpeg" {
			t.Fatalf("so o ffmpeg pode ser chamado: %v", c[0])
		}
		juntos += strings.Join(c, " ") + "\n"
	}
	for _, esperado := range []string{"-f lavfi", "anoisesrc", "seed=7", "sine=f=660", "sine=f=1320", "amix=inputs=4", "pcm_s16le"} {
		if !strings.Contains(juntos, esperado) {
			t.Errorf("comando de sintese sem %q:\n%s", esperado, juntos)
		}
	}
	entradas, _ := filepath.Glob(filepath.Join(dir, "audio", "*", "*.tmp.wav"))
	mais, _ := filepath.Glob(filepath.Join(dir, "audio", "*.tmp.wav"))
	entradas = append(entradas, mais...)
	if len(entradas) != 0 {
		t.Fatalf("nao pode sobrar arquivo temporario: %v", entradas)
	}
}

func TestAtivosSonorosNaoRegeraNemSobrescreve(t *testing.T) {
	dir := t.TempDir()
	ex := &execFalso{}
	_ = GarantirAtivosSonoros(context.Background(), dir, ConfigPadrao(), ex)
	caminho := filepath.Join(dir, "audio", "bgm.wav")
	if err := os.WriteFile(caminho, []byte("do usuario"), 0o644); err != nil {
		t.Fatal(err)
	}
	ex.comandos = nil
	if err := GarantirAtivosSonoros(context.Background(), dir, ConfigPadrao(), ex); err != nil {
		t.Fatal(err)
	}
	if len(ex.comandos) != 0 {
		t.Fatalf("nada a gerar quando tudo existe: %v", ex.comandos)
	}
	if b, _ := os.ReadFile(caminho); string(b) != "do usuario" {
		t.Fatal("arquivo existente nao pode ser sobrescrito")
	}
}

func TestAtivosSonorosRespeitaRenderJSON(t *testing.T) {
	dir := t.TempDir()
	ex := &execFalso{}
	bgm := filepath.Join(dir, "minha.wav")
	if err := os.WriteFile(bgm, wavTeste(1, 100), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := ConfigRender{SchemaVersion: 1, BGM: "minha.wav", SFXCatalogo: "cat.json"}
	if err := GarantirAtivosSonoros(context.Background(), dir, cfg, ex); err != nil || len(ex.comandos) != 0 {
		t.Fatalf("com bgm e catalogo proprios nada e sintetizado: %v %v", err, ex.comandos)
	}
	cfg.BGM = "ausente.wav"
	if err := GarantirAtivosSonoros(context.Background(), dir, cfg, ex); !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("bgm configurado e ausente deve bloquear: %v", err)
	}
	cfg.BGM = "../fora.wav"
	if err := GarantirAtivosSonoros(context.Background(), dir, cfg, ex); !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("bgm fora da raiz deve bloquear: %v", err)
	}
}

func TestAtivosSonorosFalhas(t *testing.T) {
	dir := t.TempDir()
	if err := GarantirAtivosSonoros(context.Background(), dir, ConfigPadrao(), &execFalso{erro: errors.New("sem ffmpeg")}); err == nil {
		t.Fatal("falha do ffmpeg deve propagar")
	}
	if existe(dir, "audio/bgm.wav") || existe(dir, "audio/sfx/troca_cena.wav") {
		t.Fatal("falha nao pode deixar arquivo parcial")
	}
	ctx, cancela := context.WithCancel(context.Background())
	cancela()
	if err := GarantirAtivosSonoros(ctx, dir, ConfigPadrao(), &execFalso{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("contexto cancelado deve interromper: %v", err)
	}
}

func TestAtivosSonorosPastaBloqueada(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "audio"), []byte("arquivo"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := GarantirAtivosSonoros(context.Background(), dir, ConfigPadrao(), &execFalso{}); err == nil {
		t.Fatal("audio como arquivo deve falhar")
	}
}

// Teste com o ffmpeg real: valida que os comandos de sintese produzem WAV com a duracao esperada.
func TestAtivosSonorosComFFmpegReal(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg ausente")
	}
	dir := t.TempDir()
	if err := GarantirAtivosSonoros(context.Background(), dir, ConfigPadrao(), nil); err != nil {
		t.Fatal(err)
	}
	esperado := map[string]float64{"audio/bgm.wav": 20, "audio/sfx/troca_cena.wav": 0.7, "audio/sfx/movimento_objeto.wav": 0.2, "audio/sfx/destaque.wav": 0.6}
	for rel, dur := range esperado {
		b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		f, dados, err := lerWAV(b)
		if err != nil || f.Canais != 2 || f.Taxa != 48000 || math.Abs(duracaoWAV(f, len(dados))-dur) > 0.05 {
			t.Fatalf("%s: formato %+v duracao %.3f (esperado %.1f) err=%v", rel, f, duracaoWAV(f, len(dados)), dur, err)
		}
	}
}
