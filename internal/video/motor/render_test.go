package motor

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// executorFalso simula o motor: grava render.mp4 e um arquivo por cena em cache/.
type executorFalso struct {
	rodadas int
	variar  bool
	falha   error
}

func (e *executorFalso) Executar(dir string) error {
	if e.falha != nil {
		return e.falha
	}
	e.rodadas++
	conteudo := "cena-fixa"
	if e.variar {
		conteudo = fmt.Sprintf("cena-%d", e.rodadas)
	}
	for _, nome := range []string{"l000.mp4", "l001.mp4"} {
		if err := escreverConteudo(filepath.Join(dir, DirCache, nome), conteudo); err != nil {
			return err
		}
	}
	return escreverConteudo(filepath.Join(dir, ArquivoRender), "mp4")
}

func escreverConteudo(caminho, conteudo string) error {
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		return err
	}
	return os.WriteFile(caminho, []byte(conteudo), 0o644)
}

// sondaFalsa devolve a mídia medida sem ffprobe.
type sondaFalsa struct {
	midia Midia
	lido  string
}

func (s *sondaFalsa) Medir(mp4 string) (Midia, error) {
	s.lido = mp4
	return s.midia, nil
}

func midiaReferencia() Midia {
	return Midia{Largura: 1080, Altura: 1920, FPS: 30, TemAudio: true, DuracaoS: 6.71}
}

// AC-E10.4: o MP4 tem 1080x1920, áudio, duração da voz (± 1 quadro) e zero quadros em branco.
func TestRenderReferenciaVerticalComVoz(t *testing.T) {
	dir := t.TempDir()
	s := &sondaFalsa{midia: midiaReferencia()}
	m, err := Renderizar(PedidoRender{Dir: dir, DuracaoVozS: 6.71, FPS: 30}, &executorFalso{}, s)
	if err != nil || m.Largura != 1080 || m.Altura != 1920 {
		t.Fatalf("Renderizar: %v (%+v)", err, m)
	}
	if s.lido != filepath.Join(dir, ArquivoRender) {
		t.Fatalf("sonda leu %q, quer %q", s.lido, filepath.Join(dir, ArquivoRender))
	}
}

// AC-E10.4: formato, áudio, duração ou quadro em branco fora do contrato retorna erro tipado.
func TestRenderRecusaMidiaForaDoContrato(t *testing.T) {
	casos := map[string]struct {
		mudar func(*Midia)
		quer  error
	}{
		"horizontal":       {func(m *Midia) { m.Largura, m.Altura = 1920, 1080 }, ErrRenderFormato},
		"24 fps":           {func(m *Midia) { m.FPS = 24 }, ErrRenderFormato},
		"sem audio":        {func(m *Midia) { m.TemAudio = false }, ErrRenderSemAudio},
		"2 quadros a mais": {func(m *Midia) { m.DuracaoS = 6.71 + 2.0/30 }, ErrRenderDuracao},
		"quadro branco":    {func(m *Midia) { m.TrechosBrancos = 1 }, ErrRenderQuadroBranco},
	}
	for nome, c := range casos {
		m := midiaReferencia()
		c.mudar(&m)
		_, err := Renderizar(PedidoRender{Dir: t.TempDir(), DuracaoVozS: 6.71, FPS: 30}, &executorFalso{}, &sondaFalsa{midia: m})
		if !errors.Is(err, c.quer) {
			t.Errorf("%s: erro %v, quer %v", nome, err, c.quer)
		}
	}
	m := midiaReferencia()
	m.DuracaoS = 6.71 + 0.9/30
	if err := ValidarMidia(m, 6.71, 30); err != nil {
		t.Fatalf("diferença menor que 1 quadro deveria passar: %v", err)
	}
}

// AC-E10.4: falha do motor ou MP4 ausente retorna ErrRenderMotor.
func TestRenderFalhaDoMotor(t *testing.T) {
	p := PedidoRender{Dir: t.TempDir(), DuracaoVozS: 6.71, FPS: 30}
	s := &sondaFalsa{midia: midiaReferencia()}
	if _, err := Renderizar(p, &executorFalso{falha: errors.New("chromium caiu")}, s); !errors.Is(err, ErrRenderMotor) {
		t.Fatalf("falha do motor: erro %v", err)
	}
	if _, err := Renderizar(p, executorVazio{}, s); !errors.Is(err, ErrRenderMotor) {
		t.Fatalf("MP4 ausente: erro %v", err)
	}
}

type executorVazio struct{}

func (executorVazio) Executar(string) error { return nil }

// AC-E10.5: dois renders da mesma timeline geram o mesmo hash por cena.
func TestRenderDeterminismoDoCache(t *testing.T) {
	ex := &executorFalso{}
	hashes, err := ConferirDeterminismo(t.TempDir(), ex)
	if err != nil || ex.rodadas != 2 || len(hashes) != 2 {
		t.Fatalf("determinismo: %v (rodadas %d, hashes %v)", err, ex.rodadas, hashes)
	}
	if _, err := ConferirDeterminismo(t.TempDir(), &executorFalso{variar: true}); !errors.Is(err, ErrCacheDivergente) {
		t.Fatalf("cache variável deveria divergir, erro %v", err)
	}
	if _, err := ConferirDeterminismo(t.TempDir(), executorVazio{}); !errors.Is(err, ErrCacheDivergente) {
		t.Fatalf("cache vazio deveria falhar, erro %v", err)
	}
}

// quadro desenha um fundo cinza e a boca na região indicada com a altura dada.
func quadro(boca image.Rectangle, abertura int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := range 64 {
		for x := range 64 {
			img.Set(x, y, color.RGBA{180, 150, 120, 255})
		}
	}
	for y := boca.Min.Y; y < boca.Min.Y+abertura && y < boca.Max.Y; y++ {
		for x := boca.Min.X; x < boca.Max.X; x++ {
			img.Set(x, y, color.RGBA{60, 10, 20, 255})
		}
	}
	return img
}

// AC-E10.7: trecho falado mostra 2+ estados de boca; trecho sem fala mostra 1.
func TestRenderBocaAcompanhaFala(t *testing.T) {
	boca := image.Rect(24, 40, 40, 50)
	fala := []image.Image{quadro(boca, 1), quadro(boca, 8), quadro(boca, 4), quadro(boca, 1), quadro(boca, 8)}
	silencio := []image.Image{quadro(boca, 1), quadro(boca, 1), quadro(boca, 1)}
	if n := EstadosBoca(fala, boca); n < 2 {
		t.Fatalf("fala com %d estado(s) de boca, quer 2 ou mais", n)
	}
	if n := EstadosBoca(silencio, boca); n != 1 {
		t.Fatalf("silêncio com %d estados de boca, quer 1", n)
	}
	if err := ValidarBoca(fala, silencio, boca); err != nil {
		t.Fatalf("ValidarBoca: %v", err)
	}
	if err := ValidarBoca(silencio, silencio, boca); !errors.Is(err, ErrBocaParada) {
		t.Fatalf("fala parada: erro %v", err)
	}
	if err := ValidarBoca(fala, fala, boca); !errors.Is(err, ErrBocaForaDaFala) {
		t.Fatalf("boca mexendo sem fala: erro %v", err)
	}
}

// AC-E10.4: a sonda real mede um MP4 9:16 sintético com áudio e acha o trecho branco.
func TestRenderSondaFFmpegReal(t *testing.T) {
	for _, bin := range []string{"ffmpeg", "ffprobe"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s ausente: %v", bin, err)
		}
	}
	dir := t.TempDir()
	cheio := gerarClipe(t, dir, "cheio.mp4", "red")
	m, err := SondaFFmpeg{}.Medir(cheio)
	if err != nil {
		t.Fatalf("Medir: %v", err)
	}
	if err := ValidarMidia(m, 2, 30); err != nil {
		t.Fatalf("clipe vermelho 9:16 deveria passar: %v (%+v)", err, m)
	}
	branco, err := SondaFFmpeg{}.Medir(gerarClipe(t, dir, "branco.mp4", "white"))
	if err != nil || branco.TrechosBrancos == 0 {
		t.Fatalf("clipe branco sem trecho branco detectado: %v (%+v)", err, branco)
	}
}

func gerarClipe(t *testing.T, dir, nome, cor string) string {
	t.Helper()
	saida := filepath.Join(dir, nome)
	cmd := exec.Command("ffmpeg", "-v", "error", "-y",
		"-f", "lavfi", "-i", "color=c="+cor+":s=1080x1920:r=30:d=2",
		"-f", "lavfi", "-i", "sine=frequency=440:duration=2",
		"-c:v", "libx264", "-pix_fmt", "yuv420p", "-c:a", "aac", "-shortest", saida)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("ffmpeg gerar clipe: %v: %s", err, out)
	}
	return saida
}
