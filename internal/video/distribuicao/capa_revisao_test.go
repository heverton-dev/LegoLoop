package distribuicao

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type execFalso struct {
	chamadas [][]string
	err      error
	gravar   string
	tamanho  int
}

func (e *execFalso) Executar(_ context.Context, nome string, args ...string) error {
	e.chamadas = append(e.chamadas, append([]string{nome}, args...))
	if e.err != nil {
		return e.err
	}
	if e.gravar != "" {
		return os.WriteFile(e.gravar, make([]byte, e.tamanho), 0o644)
	}
	return nil
}

func pedidoCapa(t *testing.T) (string, PedidoCapa) {
	t.Helper()
	dir := t.TempDir()
	tocar(t, dir, "fundos/f.jpg")
	return dir, PedidoCapa{Fundo: "fundos/f.jpg", Linha1: "US$ 2,7 BI", Linha2: "SUMIRAM", Saida: "capas/a.png", Ouro: "#FFD24A", Texto: "#FFFFFF", Logo: "logo.png", Fonte: "f.ttf", TextoDireita: true}
}

func TestGerarCapaMontaOComandoEConfereAArquivo(t *testing.T) {
	dir, p := pedidoCapa(t)
	ex := &execFalso{gravar: filepath.Join(dir, "capas", "a.png"), tamanho: 100}
	_ = os.MkdirAll(filepath.Join(dir, "capas"), 0o755)
	res, err := GerarCapa(context.Background(), ex, dir, p)
	if err != nil || res.Bytes != 100 || res.Aviso != "" {
		t.Fatalf("GerarCapa: %+v %v", res, err)
	}
	args := strings.Join(ex.chamadas[0], " ")
	for _, quer := range []string{"python", "--linha1 US$ 2,7 BI", "--linha2 SUMIRAM", "--logo logo.png", "--fonte f.ttf", "--ouro #FFD24A", "--texto #FFFFFF", "--texto-direita"} {
		if !strings.Contains(args, quer) {
			t.Fatalf("comando sem %q: %s", quer, args)
		}
	}
	ex2 := &execFalso{gravar: filepath.Join(dir, "capas", "a.png"), tamanho: LimiteCapaBytes + 1}
	if res, _ := GerarCapa(context.Background(), ex2, dir, p); res.Aviso == "" {
		t.Fatal("capa acima de 2 MB deve avisar")
	}
}

func TestGerarCapaRecusaPedidoRuim(t *testing.T) {
	dir, p := pedidoCapa(t)
	casos := map[string]func(*PedidoCapa){
		"palavras demais": func(p *PedidoCapa) { p.Linha2 = "UM DOIS TRES QUATRO" },
		"fundo ausente":   func(p *PedidoCapa) { p.Fundo = "fundos/nao.jpg" },
		"saida fora":      func(p *PedidoCapa) { p.Saida = "../x.png" },
		"fundo absoluto":  func(p *PedidoCapa) { p.Fundo = "C:/x.jpg" },
	}
	for nome, mudar := range casos {
		t.Run(nome, func(t *testing.T) {
			q := p
			mudar(&q)
			ex := &execFalso{}
			if _, err := GerarCapa(context.Background(), ex, dir, q); err == nil || len(ex.chamadas) != 0 {
				t.Fatalf("quer erro sem chamar o script: %v %d", err, len(ex.chamadas))
			}
		})
	}
	falha := &execFalso{err: errors.New("python quebrou")}
	if _, err := GerarCapa(context.Background(), falha, dir, p); err == nil {
		t.Fatal("falha do script deve voltar como erro")
	}
	semArquivo := &execFalso{}
	if _, err := GerarCapa(context.Background(), semArquivo, dir, p); !errors.Is(err, ErrPedidoInvalido) {
		t.Fatalf("script que nao grava: %v", err)
	}
}

// Teste real: roda o script embutido com Pillow, quando python e Pillow existem.
func TestGerarCapaReal(t *testing.T) {
	if err := exec.Command("python", "-c", "import PIL").Run(); err != nil {
		t.Skip("python com Pillow ausente")
	}
	dir := t.TempDir()
	fundo := filepath.Join(dir, "fundo.png")
	if err := exec.Command("python", "-c", "from PIL import Image; Image.new('RGB',(1600,900),(20,30,60)).save(r'"+fundo+"')").Run(); err != nil {
		t.Fatal(err)
	}
	res, err := GerarCapa(context.Background(), ExecSistema{}, dir, PedidoCapa{Fundo: "fundo.png", Linha1: "2,7 BI", Linha2: "SUMIRAM", Saida: "capa.png"})
	if err != nil || res.Bytes == 0 {
		t.Fatalf("capa real: %+v %v", res, err)
	}
	var w, h int
	out, err := exec.Command("python", "-c", "from PIL import Image; im=Image.open(r'"+filepath.Join(dir, "capa.png")+"'); print(im.width, im.height)").Output()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fmt.Sscan(string(out), &w, &h); err != nil || w != 1280 || h != 720 {
		t.Fatalf("capa deve ter 1280x720: %q", out)
	}
}

func TestExecSistema(t *testing.T) {
	if err := (ExecSistema{}).Executar(context.Background(), "programa-que-nao-existe-xyz"); err == nil {
		t.Fatal("programa inexistente deve falhar")
	}
	if err := (ExecSistema{}).Executar(context.Background(), "go", "version"); err != nil {
		t.Fatalf("go version: %v", err)
	}
}

func parecerOK() Parecer {
	return Parecer{Versao: 1, Veredito: VereditoNaoEnviar, Revisor: "agente-isolado", Problemas: []Problema{{Versao: 1, Corte: "l003", Tempo: "00:12", Problema: "texto cortado", Correcao: "mover para a area segura"}}}
}

func TestParecerValidar(t *testing.T) {
	casos := map[string]func(*Parecer){
		"sem revisor":      func(p *Parecer) { p.Revisor = "" },
		"versao zero":      func(p *Parecer) { p.Versao = 0 },
		"veredito":         func(p *Parecer) { p.Veredito = "TALVEZ" },
		"nao enviar vazio": func(p *Parecer) { p.Problemas = nil },
		"tempo ruim":       func(p *Parecer) { p.Problemas[0].Tempo = "12s" },
		"sem correcao":     func(p *Parecer) { p.Problemas[0].Correcao = "" },
		"sem corte":        func(p *Parecer) { p.Problemas[0].Corte = " " },
		"sem descricao":    func(p *Parecer) { p.Problemas[0].Problema = "" },
	}
	for nome, mudar := range casos {
		t.Run(nome, func(t *testing.T) {
			p := parecerOK()
			mudar(&p)
			if err := p.Validar(); !errors.Is(err, ErrParecerInvalido) {
				t.Fatalf("quer ErrParecerInvalido, obteve %v", err)
			}
		})
	}
	p := parecerOK()
	p.Veredito, p.Problemas = VereditoEnviar, nil
	if err := p.Validar(); err != nil {
		t.Fatalf("ENVIAR sem problemas e valido: %v", err)
	}
}

func TestGravarParecerELiberar(t *testing.T) {
	dir := t.TempDir()
	if err := Liberado(dir, 1); !errors.Is(err, ErrRevisaoReprovada) {
		t.Fatalf("sem parecer: %v", err)
	}
	if err := GravarParecer(dir, parecerOK()); err != nil {
		t.Fatal(err)
	}
	if err := GravarParecer(dir, parecerOK()); !errors.Is(err, ErrParecerInvalido) {
		t.Fatalf("nao pode sobrescrever: %v", err)
	}
	if err := Liberado(dir, 1); !errors.Is(err, ErrRevisaoReprovada) {
		t.Fatalf("NAO ENVIAR nao libera: %v", err)
	}
	ok := Parecer{Versao: 2, Veredito: VereditoEnviar, Revisor: "r"}
	if err := GravarParecer(dir, ok); err != nil {
		t.Fatal(err)
	}
	if err := Liberado(dir, 2); err != nil {
		t.Fatalf("ENVIAR deve liberar: %v", err)
	}
	if err := GravarParecer(dir, Parecer{}); !errors.Is(err, ErrParecerInvalido) {
		t.Fatalf("parecer vazio: %v", err)
	}
	_ = os.WriteFile(CaminhoParecer(dir, 3), []byte("{"), 0o644)
	if err := Liberado(dir, 3); !errors.Is(err, ErrParecerInvalido) {
		t.Fatalf("json quebrado: %v", err)
	}
	_ = os.WriteFile(CaminhoParecer(dir, 4), []byte(`{"versao":4,"veredito":"X","revisor":"r"}`), 0o644)
	if err := Liberado(dir, 4); !errors.Is(err, ErrParecerInvalido) {
		t.Fatalf("parecer invalido gravado: %v", err)
	}
}

func TestProximaVersao(t *testing.T) {
	dir := t.TempDir()
	if n, err := ProximaVersao(dir); err != nil || n != 1 {
		t.Fatalf("sem renders: %d %v", n, err)
	}
	for _, rel := range []string{"renders/v1.mp4", "renders/v3.mp4", "renders/v3_9x16.mp4", "renders/legenda.srt"} {
		tocar(t, dir, rel)
	}
	if n, err := ProximaVersao(dir); err != nil || n != 4 {
		t.Fatalf("quer 4: %d %v", n, err)
	}
	if _, err := ProximaVersao("pasta\x00invalida"); err == nil {
		t.Fatal("caminho invalido deve falhar")
	}
}

func TestGerarFolhas(t *testing.T) {
	dir := t.TempDir()
	tocar(t, dir, "renders/v1.mp4")
	medir := func(context.Context, string) (float64, error) { return 100, nil }
	ex := &execFalso{}
	folhas, err := GerarFolhas(context.Background(), ex, medir, dir, "renders/v1.mp4", 2)
	if err != nil || len(folhas) != 3 || folhas[0] != "review/sheet_01.jpg" || len(ex.chamadas) != 3 {
		t.Fatalf("100 s a 2 s por quadro = 3 folhas: %v %v", folhas, err)
	}
	if !strings.Contains(strings.Join(ex.chamadas[1], " "), "-ss 48") {
		t.Fatalf("segunda folha comeca em 48 s: %v", ex.chamadas[1])
	}
	if _, err := GerarFolhas(context.Background(), ex, medir, dir, "renders/v1.mp4", 0); !errors.Is(err, ErrPedidoInvalido) {
		t.Fatalf("passo zero: %v", err)
	}
	if _, err := GerarFolhas(context.Background(), ex, medir, dir, "renders/v9.mp4", 1); !errors.Is(err, ErrVitrineInvalida) {
		t.Fatalf("video ausente: %v", err)
	}
	erroMedir := func(context.Context, string) (float64, error) { return 0, errors.New("ffprobe falhou") }
	if _, err := GerarFolhas(context.Background(), ex, erroMedir, dir, "renders/v1.mp4", 1); err == nil {
		t.Fatal("erro do medidor deve voltar")
	}
	if _, err := GerarFolhas(context.Background(), &execFalso{err: errors.New("ffmpeg falhou")}, medir, dir, "renders/v1.mp4", 1); err == nil {
		t.Fatal("erro do ffmpeg deve voltar")
	}
	bloqueio := t.TempDir()
	tocar(t, bloqueio, "renders/v1.mp4")
	_ = os.WriteFile(filepath.Join(bloqueio, "review"), []byte("x"), 0o644)
	if _, err := GerarFolhas(context.Background(), ex, medir, bloqueio, "renders/v1.mp4", 1); err == nil {
		t.Fatal("review que e arquivo deve falhar")
	}
}

func TestDuracaoFFprobe(t *testing.T) {
	if _, err := DuracaoFFprobe(context.Background(), filepath.Join(t.TempDir(), "nao-existe.mp4")); err == nil {
		t.Fatal("arquivo inexistente deve falhar")
	}
}

func TestGravarPacote(t *testing.T) {
	dir := t.TempDir()
	for _, rel := range []string{"renders/v1.mp4", "script.json", "review/sheet_01.jpg"} {
		tocar(t, dir, rel)
	}
	caminho, err := GravarPacote(dir, 1, []string{"review/sheet_01.jpg"}, 2)
	if err != nil {
		t.Fatal(err)
	}
	dados, _ := os.ReadFile(caminho)
	for _, quer := range []string{`"video": "renders/v1.mp4"`, `"roteiro": "script.json"`, "sem maos ou membros extras", "maos", `"passo_s": 2`} {
		if !strings.Contains(string(dados), quer) {
			t.Fatalf("pacote sem %q:\n%s", quer, dados)
		}
	}
	if _, err := GravarPacote(dir, 1, nil, 2); !errors.Is(err, ErrPedidoInvalido) {
		t.Fatalf("sem folhas: %v", err)
	}
	if _, err := GravarPacote(dir, 7, []string{"review/sheet_01.jpg"}, 2); !errors.Is(err, ErrVitrineInvalida) {
		t.Fatalf("video da versao ausente: %v", err)
	}
}
