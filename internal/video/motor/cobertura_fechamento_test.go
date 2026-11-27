package motor

import (
	"errors"
	"github.com/heverton-dev/LegoLoop/internal/video/ativos"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

func TestPoseOuIdle(t *testing.T) {
	if poseOuIdle("") != "idle" {
		t.Fatal("pose vazia deveria virar idle")
	}
	if poseOuIdle("gesture") != "gesture" {
		t.Fatal("pose definida nao deveria mudar")
	}
}

func TestDirecao(t *testing.T) {
	casos := []struct {
		olhar string
		x     float64
		quer  int
	}{
		{"esquerda", 500, -1}, {"direita", 500, 1}, {"", 500, 1}, {"", 1500, -1},
	}
	for _, c := range casos {
		if got := direcao(c.olhar, c.x); got != c.quer {
			t.Errorf("direcao(%q,%v) = %d, quer %d", c.olhar, c.x, got, c.quer)
		}
	}
}

func TestLuzDoEstiloSemLuz(t *testing.T) {
	if luzDoEstilo(inventario.Estilo{}) != nil {
		t.Fatal("estilo sem luz nao deveria gerar rim")
	}
	if luzDoEstilo(inventario.Estilo{Luz: &inventario.Luz{Rim: false}}) != nil {
		t.Fatal("luz sem rim nao deveria gerar rim")
	}
}

func TestLerTaxa(t *testing.T) {
	casos := []struct {
		taxa string
		quer float64
	}{
		{"30/1", 30}, {"30", 30}, {"x/1", 0}, {"30/0", 0}, {"30/x", 0},
	}
	for _, c := range casos {
		if got := lerTaxa(c.taxa); got != c.quer {
			t.Errorf("lerTaxa(%q) = %v, quer %v", c.taxa, got, c.quer)
		}
	}
}

func TestExecutorPythonSucessoEErro(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(t.TempDir(), "ok.py")
	if err := os.WriteFile(script, []byte("import sys\nprint(sys.argv[1])\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := (ExecutorPython{Script: script}).Executar(dir); err != nil {
		t.Fatalf("esperava sucesso, veio %v", err)
	}
	if err := (ExecutorPython{Script: filepath.Join(t.TempDir(), "nao-existe.py")}).Executar(dir); err == nil {
		t.Fatal("esperava erro para script inexistente")
	}
}

func TestLuzDoEstiloComCorExplicita(t *testing.T) {
	l := luzDoEstilo(inventario.Estilo{Luz: &inventario.Luz{Rim: true, Cor: "#111111"}})
	if l == nil || l.Rim != "#111111" {
		t.Fatalf("luz = %+v, queria rim #111111", l)
	}
	l2 := luzDoEstilo(inventario.Estilo{Luz: &inventario.Luz{Rim: true}})
	if l2 == nil || l2.Rim != "#FFD27A" {
		t.Fatalf("luz = %+v, queria rim padrao #FFD27A", l2)
	}
}

func TestValidarProporcaoCasos(t *testing.T) {
	casos := []struct {
		nome string
		prop any
		ok   bool
	}{
		{"nao e mapa", "x", false},
		{"sem cabeca", map[string]any{"corpo": 0.5}, false},
		{"sem corpo", map[string]any{"cabeca": 1.2}, false},
		{"cabeca <= 1", map[string]any{"cabeca": 1.0, "corpo": 0.5}, false},
		{"corpo <= 0", map[string]any{"cabeca": 1.2, "corpo": 0.0}, false},
		{"corpo >= 1", map[string]any{"cabeca": 1.2, "corpo": 1.0}, false},
		{"valido", map[string]any{"cabeca": 1.2, "corpo": 0.5}, true},
	}
	for _, c := range casos {
		err := validarProporcao(c.prop)
		if (err == nil) != c.ok {
			t.Errorf("%s: err=%v, queria ok=%v", c.nome, err, c.ok)
		}
	}
}

type checadorComErro struct{}

func (checadorComErro) Checar(dir string) (string, error) { return "", os.ErrPermission }

func TestGerarChecadorComErro(t *testing.T) {
	e, dir := entradaDemo(t)
	if err := Gerar(e, dir, checadorComErro{}); !errors.Is(err, ErrCheckTimeline) {
		t.Fatalf("err = %v, queria ErrCheckTimeline", err)
	}
}

func TestGerarEntradaInvalidaNaoChamaChecador(t *testing.T) {
	e, dir := entradaDemo(t)
	e.Timings.Linhas = nil
	if err := Gerar(e, dir, &checadorFalso{saida: "CHECK OK"}); err == nil {
		t.Fatal("esperava erro de Montar propagado por Gerar")
	}
}

type sondaComErro struct{}

func (sondaComErro) Medir(string) (Midia, error) { return Midia{}, os.ErrInvalid }

func TestRenderizarSondaComErro(t *testing.T) {
	p := PedidoRender{Dir: t.TempDir(), DuracaoVozS: 6.71, FPS: 30}
	if _, err := Renderizar(p, &executorFalso{}, sondaComErro{}); !errors.Is(err, ErrRenderMotor) {
		t.Fatalf("sonda com erro: %v, queria ErrRenderMotor", err)
	}
}

func TestConferirDeterminismoExecutorFalha(t *testing.T) {
	if _, err := ConferirDeterminismo(t.TempDir(), &executorFalso{falha: errors.New("caiu")}); !errors.Is(err, ErrRenderMotor) {
		t.Fatalf("executor com falha: %v, queria ErrRenderMotor", err)
	}
}

func TestDependenciasPadraoRodarReal(t *testing.T) {
	d := Dependencias{}.comPadroes()
	if err := d.Rodar("python", "-c", "pass"); err != nil {
		t.Fatalf("Rodar padrao com python -c pass deveria funcionar: %v", err)
	}
	if err := d.Rodar("comando-que-nao-existe-de-jeito-nenhum"); err == nil {
		t.Fatal("Rodar padrao com comando inexistente deveria falhar")
	}
}

func TestMontarFPSPadrao(t *testing.T) {
	e, dir := entradaDemo(t)
	e.FPS = 0
	tl, err := Montar(e, dir)
	if err != nil || tl.FPS != fpsPadrao {
		t.Fatalf("fps = %d, err = %v, queria %d sem erro", tl.FPS, err, fpsPadrao)
	}
}

func TestFinsEmQuadrosIndiceNegativo(t *testing.T) {
	e, _ := entradaDemo(t)
	e.Cenas.Cenas[0].IndiceLinha = -1
	if _, err := finsEmQuadros(e); !errors.Is(err, ErrLinhaSemTiming) {
		t.Fatalf("indice negativo: err = %v, queria ErrLinhaSemTiming", err)
	}
}

func TestFinsEmQuadrosSemQuadrosNoMeio(t *testing.T) {
	e, _ := entradaDemo(t)
	// insere uma 3a cena repetindo a linha 0 no meio, antes da ultima (linha 1, com override da voz).
	original := e.Cenas.Cenas
	e.Cenas.Cenas = []ativos.Cena{original[0], original[0], original[1]}
	if _, err := finsEmQuadros(e); !errors.Is(err, ErrLinhaSemTiming) {
		t.Fatalf("linha repetida no meio: err = %v, queria ErrLinhaSemTiming", err)
	}
}
