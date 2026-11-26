package render

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
)

var fmtTeste = formatoWAV{Canais: 1, Taxa: 24000, Bits: 16}

func wavTeste(valor byte, amostras int) []byte {
	return montarWAV(fmtTeste, bytes.Repeat([]byte{valor}, amostras*2))
}

type sintFalso struct {
	nome     string
	pago     bool
	chamadas int
	textos   []string
	opcoes   []OpcoesVoz
	fn       func(i int, texto string) ([]byte, error)
}

func (s *sintFalso) Nome() string { return s.nome }
func (s *sintFalso) Pago() bool   { return s.pago }
func (s *sintFalso) Sintetizar(_ context.Context, texto string, op OpcoesVoz) ([]byte, error) {
	s.chamadas++
	s.textos = append(s.textos, texto)
	s.opcoes = append(s.opcoes, op)
	if s.fn != nil {
		return s.fn(s.chamadas-1, texto)
	}
	return wavTeste(byte(s.chamadas), 1200*s.chamadas), nil
}

func roteiroTeste() (roteiro.Marca, roteiro.Script) {
	m := roteiro.Marca{Voz: roteiro.VozMarca{Motor: "falso", Voz: "voz-x", Velocidade: "+10%"}}
	s := roteiro.Script{Lines: []roteiro.LinhaRoteiro{{Text: "a conta"}, {Text: "o PIB sobe"}}}
	return m, s
}

func locutorCom(s Sintetizador, gate GateCreditos) Locutor {
	r := NovoRegistro()
	_ = r.Registrar(s)
	return Locutor{Registro: r, Gate: gate}
}

func TestVozRegistro(t *testing.T) {
	r := NovoRegistro()
	if err := r.Registrar(&sintFalso{nome: "a"}); err != nil {
		t.Fatal(err)
	}
	if err := r.Registrar(&sintFalso{nome: "a"}); err == nil {
		t.Fatal("registro duplicado deve falhar")
	}
	if err := r.Registrar(&sintFalso{nome: " "}); err == nil {
		t.Fatal("nome vazio deve falhar")
	}
	if !r.Existe("a") || r.Existe("b") {
		t.Fatal("Existe incorreto")
	}
	if _, err := r.Obter("b"); !errors.Is(err, ErrProvedorDesconhecido) {
		t.Fatalf("esperava ErrProvedorDesconhecido, obteve %v", err)
	}
}

func TestVozSintetizaRoteiro(t *testing.T) {
	s := &sintFalso{nome: "falso"}
	dir := t.TempDir()
	m, sc := roteiroTeste()
	res, err := locutorCom(s, nil).SintetizarRoteiro(context.Background(), dir, m, sc)
	if err != nil {
		t.Fatal(err)
	}
	if res.Provedor != "falso" || res.Voz != "voz-x" || len(res.Linhas) != 2 || len(res.DuracoesS) != 2 {
		t.Fatalf("resultado inesperado: %+v", res)
	}
	if s.opcoes[0].Voz != "voz-x" || s.opcoes[0].Velocidade != "+10%" {
		t.Fatalf("opcoes nao repassadas: %+v", s.opcoes[0])
	}
	b, err := os.ReadFile(filepath.Join(dir, "audio", "voice.wav"))
	if err != nil {
		t.Fatal(err)
	}
	_, dados, err := lerWAV(b)
	if err != nil || len(dados) != (1200+2400)*2 {
		t.Fatalf("voice.wav com tamanho errado: %d err=%v", len(dados), err)
	}
	if dados[0] != 1 || dados[len(dados)-1] != 2 {
		t.Fatal("linhas fora de ordem em voice.wav")
	}
	for _, nome := range res.Linhas {
		if _, err := os.Stat(filepath.Join(dir, nome)); err != nil {
			t.Fatalf("arquivo por linha ausente: %v", err)
		}
	}
	soma := res.DuracoesS[0] + res.DuracoesS[1]
	if soma < 0.149 || soma > 0.151 {
		t.Fatalf("duracao somada incorreta: %f", soma)
	}
}

func TestVozPadraoThalita(t *testing.T) {
	s := &sintFalso{nome: "falso"}
	m, sc := roteiroTeste()
	m.Voz.Voz = ""
	res, err := locutorCom(s, nil).SintetizarRoteiro(context.Background(), t.TempDir(), m, sc)
	if err != nil || res.Voz != "pt-BR-ThalitaMultilingualNeural" {
		t.Fatalf("voz padrao incorreta: %+v err=%v", res, err)
	}
}

func TestVozGateCreditos(t *testing.T) {
	m, sc := roteiroTeste()
	nega := func(context.Context, string) error { return errors.New("sem saldo") }
	libera := func(context.Context, string) error { return nil }
	casos := []struct {
		nome       string
		pago       bool
		gate       GateCreditos
		esperaErro bool
	}{
		{"pago com gate negando", true, nega, true},
		{"pago sem gate", true, nil, true},
		{"pago com gate liberando", true, libera, false},
		{"gratuito ignora o gate", false, nega, false},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			s := &sintFalso{nome: "falso", pago: c.pago}
			_, err := locutorCom(s, c.gate).SintetizarRoteiro(context.Background(), t.TempDir(), m, sc)
			if c.esperaErro {
				if !errors.Is(err, ErrCreditosNegados) || s.chamadas != 0 {
					t.Fatalf("esperava ErrCreditosNegados sem chamada: err=%v chamadas=%d", err, s.chamadas)
				}
				return
			}
			if err != nil || s.chamadas != 2 {
				t.Fatalf("esperava sucesso: err=%v chamadas=%d", err, s.chamadas)
			}
		})
	}
}

func TestVozFalhas(t *testing.T) {
	m, sc := roteiroTeste()
	casos := map[string]func(int, string) ([]byte, error){
		"erro do provedor": func(int, string) ([]byte, error) { return nil, errors.New("rede") },
		"audio invalido":   func(int, string) ([]byte, error) { return []byte("nao e wav"), nil },
		"formato diferente": func(i int, _ string) ([]byte, error) {
			if i == 0 {
				return wavTeste(1, 100), nil
			}
			return montarWAV(formatoWAV{Canais: 2, Taxa: 44100, Bits: 16}, make([]byte, 400)), nil
		},
	}
	for nome, fn := range casos {
		t.Run(nome, func(t *testing.T) {
			dir := t.TempDir()
			_, err := locutorCom(&sintFalso{nome: "falso", fn: fn}, nil).SintetizarRoteiro(context.Background(), dir, m, sc)
			if err == nil {
				t.Fatal("esperava erro")
			}
			if _, e := os.Stat(filepath.Join(dir, "audio", "voice.wav")); e == nil {
				t.Fatal("voice.wav nao deve existir apos falha")
			}
		})
	}
}

func TestVozProvedorEVazioDesconhecido(t *testing.T) {
	m, sc := roteiroTeste()
	m.Voz.Motor = "outro"
	l := locutorCom(&sintFalso{nome: "falso"}, nil)
	if _, err := l.SintetizarRoteiro(context.Background(), t.TempDir(), m, sc); !errors.Is(err, ErrProvedorDesconhecido) {
		t.Fatalf("esperava ErrProvedorDesconhecido, obteve %v", err)
	}
	m.Voz.Motor = "falso"
	sc.Lines = nil
	if _, err := l.SintetizarRoteiro(context.Background(), t.TempDir(), m, sc); err == nil {
		t.Fatal("roteiro sem linhas deve falhar")
	}
	if _, err := (Locutor{}).SintetizarRoteiro(context.Background(), t.TempDir(), m, sc); err == nil {
		t.Fatal("locutor sem registro deve falhar")
	}
}

func TestVozPronuncia(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pronuncia.json"), []byte(`{"PIB":"pê i bê"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	s := &sintFalso{nome: "falso"}
	m, sc := roteiroTeste()
	if _, err := locutorCom(s, nil).SintetizarRoteiro(context.Background(), dir, m, sc); err != nil {
		t.Fatal(err)
	}
	if s.textos[1] != "o pê i bê sobe" || s.textos[0] != "a conta" {
		t.Fatalf("pronuncia nao aplicada so a palavra inteira: %q", s.textos)
	}
	if got := AplicarPronuncia("PIBinho e PIB", map[string]string{"PIB": "X"}); got != "PIBinho e X" {
		t.Fatalf("deve trocar so palavra inteira: %q", got)
	}
	if err := os.WriteFile(filepath.Join(dir, "pronuncia.json"), []byte(`{`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := locutorCom(s, nil).SintetizarRoteiro(context.Background(), dir, m, sc); err == nil {
		t.Fatal("pronuncia.json invalido deve falhar")
	}
}

type execFalso struct {
	comandos [][]string
	erro     error
	semSaida bool
}

func (e *execFalso) Executar(_ context.Context, nome string, args ...string) error {
	e.comandos = append(e.comandos, append([]string{nome}, args...))
	if e.erro != nil {
		return e.erro
	}
	if e.semSaida {
		return nil
	}
	for _, a := range args {
		if v, ok := strings.CutPrefix(a, "--write-media="); ok {
			return os.WriteFile(v, []byte("mp3"), 0o644)
		}
	}
	if nome == "ffmpeg" {
		return os.WriteFile(args[len(args)-1], wavTeste(7, 480), 0o644)
	}
	return nil
}

func TestVozEdgeTTS(t *testing.T) {
	ex := &execFalso{}
	r := RegistroPadrao(ex)
	s, err := r.Obter("edge-tts")
	if err != nil || s.Pago() || s.Nome() != "edge-tts" {
		t.Fatalf("edge-tts deve estar registrado e ser gratuito: %v", err)
	}
	b, err := s.Sintetizar(context.Background(), "--texto perigoso", OpcoesVoz{Voz: "v1", Velocidade: "+5%"})
	if err != nil {
		t.Fatal(err)
	}
	if _, dados, err := lerWAV(b); err != nil || len(dados) != 960 {
		t.Fatalf("wav inesperado: %d err=%v", len(dados), err)
	}
	if len(ex.comandos) != 2 || ex.comandos[0][0] != "edge-tts" || ex.comandos[1][0] != "ffmpeg" {
		t.Fatalf("comandos inesperados: %v", ex.comandos)
	}
	juntos := strings.Join(ex.comandos[0], " ")
	for _, esperado := range []string{"--voice=v1", "--rate=+5%", "--text=--texto perigoso"} {
		if !strings.Contains(juntos, esperado) {
			t.Fatalf("argumento %q ausente (texto nao pode virar opcao): %s", esperado, juntos)
		}
	}
}

func TestVozEdgeTTSFalhas(t *testing.T) {
	for nome, ex := range map[string]*execFalso{
		"executor falha":  {erro: errors.New("sem binario")},
		"sem arquivo wav": {semSaida: true},
	} {
		t.Run(nome, func(t *testing.T) {
			s, _ := RegistroPadrao(ex).Obter("edge-tts")
			if _, err := s.Sintetizar(context.Background(), "oi", OpcoesVoz{Voz: "v"}); err == nil {
				t.Fatal("esperava erro")
			}
		})
	}
	if RegistroPadrao(nil) == nil {
		t.Fatal("registro padrao com executor nulo deve usar o do sistema")
	}
}

func TestWAV(t *testing.T) {
	comLista := append([]byte(nil), wavTeste(3, 10)...)
	comLista = append(comLista[:12], append([]byte("LIST\x04\x00\x00\x00abcd"), comLista[12:]...)...)
	if _, dados, err := lerWAV(comLista); err != nil || len(dados) != 20 {
		t.Fatalf("chunk extra deve ser ignorado: %d err=%v", len(dados), err)
	}
	for nome, b := range map[string][]byte{"curto": []byte("RIFF"), "sem data": wavTeste(1, 1)[:36], "nao pcm": append(append([]byte(nil), wavTeste(1, 1)[:20]...), append([]byte{3, 0}, wavTeste(1, 1)[22:]...)...)} {
		if _, _, err := lerWAV(b); !errors.Is(err, ErrAudioInvalido) {
			t.Errorf("%s: esperava ErrAudioInvalido, obteve %v", nome, err)
		}
	}
	if _, err := juntarWAV(nil); !errors.Is(err, ErrAudioInvalido) {
		t.Fatalf("juntar vazio deve falhar: %v", err)
	}
	if d := duracaoWAV(fmtTeste, 48000); d != 1 {
		t.Fatalf("duracao 1s esperada: %f", d)
	}
	if d := duracaoWAV(formatoWAV{}, 10); d != 0 {
		t.Fatalf("formato zerado deve dar 0: %f", d)
	}
}

func TestVozExecutorSistema(t *testing.T) {
	var ex Executor = execSistema{}
	if err := ex.Executar(context.Background(), "go", "version"); err != nil {
		t.Fatalf("go version deve rodar: %v", err)
	}
	if err := ex.Executar(context.Background(), "programa-que-nao-existe-xyz"); err == nil {
		t.Fatal("programa inexistente deve falhar")
	}
}

func TestVozGravacaoForaDaRaiz(t *testing.T) {
	if err := gravarSeguro(t.TempDir(), "../fora.wav", []byte("x")); err == nil {
		t.Fatal("gravacao fora da raiz deve falhar")
	}
	if _, err := carregarPronuncia(t.TempDir()); err != nil {
		t.Fatalf("pronuncia ausente nao e erro: %v", err)
	}
}
