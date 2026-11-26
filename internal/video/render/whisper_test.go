package render

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

// transcritorFalso devolve a fala certa do ultimo texto sintetizado, salvo nos textos marcados como ruins.
type transcritorFalso struct {
	s        *sintFalso
	ruins    map[string]int
	hip      map[string]string
	erro     error
	chamadas int
}

func (t *transcritorFalso) Transcrever(_ context.Context, _ []byte, _ string) ([]PalavraTranscrita, error) {
	t.chamadas++
	if t.erro != nil {
		return nil, t.erro
	}
	texto := t.s.textos[len(t.s.textos)-1]
	palavras := strings.Fields(strings.ToLower(texto))
	if h, ok := t.hip[texto]; ok {
		palavras = strings.Fields(strings.ToLower(h))
	}
	if t.ruins[texto] > 0 {
		t.ruins[texto]--
		palavras = []string{"xxx", "yyy"}
	}
	saida := make([]PalavraTranscrita, len(palavras))
	for i, p := range palavras {
		saida[i] = PalavraTranscrita{Texto: p, InicioS: 0.01 * float64(i), FimS: 0.01*float64(i) + 0.005}
	}
	return saida, nil
}

func etapaCom(s *sintFalso, tr *transcritorFalso, max int) EtapaVoz {
	return EtapaVoz{Locutor: locutorCom(s, nil), Transcritor: tr, ReprovasMax: max}
}

func TestWhisperAprovaNaPrimeira(t *testing.T) {
	s := &sintFalso{nome: "falso"}
	dir := t.TempDir()
	m, sc := roteiroTeste()
	tm, err := etapaCom(s, &transcritorFalso{s: s}, 3).Executar(context.Background(), dir, m, sc)
	if err != nil {
		t.Fatal(err)
	}
	if s.chamadas != 2 || tm.Provedor != "falso" || tm.Voz != "voz-x" || len(tm.Linhas) != 2 {
		t.Fatalf("resultado inesperado: chamadas=%d %+v", s.chamadas, tm)
	}
	if tm.WERGlobal != 0 || tm.Linhas[0].WER != 0 || tm.Linhas[0].Regravacoes != 0 {
		t.Fatalf("linha 0 deveria passar de primeira: %+v", tm.Linhas[0])
	}
	if tm.Linhas[1].InicioS != tm.Linhas[0].FimS || tm.Linhas[1].Palavras[0].InicioS < tm.Linhas[1].InicioS {
		t.Fatalf("palavras da linha 1 devem usar deslocamento da linha: %+v", tm.Linhas[1])
	}
	if err := tm.Validar(); err != nil {
		t.Fatalf("timings deve validar: %v", err)
	}
	if len(tm.Linhas[0].Visemas) == 0 || len(tm.Linhas[1].Visemas) == 0 {
		t.Fatalf("a Etapa 6 deve gravar visemas em cada linha: %+v", tm.Linhas)
	}
	b, err := os.ReadFile(filepath.Join(dir, "data", "timings.json"))
	if err != nil {
		t.Fatal(err)
	}
	var lido Timings
	if err := json.Unmarshal(b, &lido); err != nil || len(lido.Linhas) != 2 {
		t.Fatalf("timings.json ilegivel: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "audio", "voice.wav")); err != nil {
		t.Fatalf("voice.wav ausente: %v", err)
	}
}

func TestWhisperRegravaSoALinhaReprovada(t *testing.T) {
	s := &sintFalso{nome: "falso"}
	tr := &transcritorFalso{s: s, ruins: map[string]int{"o PIB sobe": 2}}
	m, sc := roteiroTeste()
	tm, err := etapaCom(s, tr, 3).Executar(context.Background(), t.TempDir(), m, sc)
	if err != nil {
		t.Fatal(err)
	}
	if s.chamadas != 4 {
		t.Fatalf("esperava 1 sintese da linha 0 e 3 da linha 1, obteve %d", s.chamadas)
	}
	if tm.Linhas[0].Regravacoes != 0 || tm.Linhas[1].Regravacoes != 2 {
		t.Fatalf("regravacoes incorretas: %d e %d", tm.Linhas[0].Regravacoes, tm.Linhas[1].Regravacoes)
	}
	if tm.Linhas[1].WER > LimiteWERLinha {
		t.Fatalf("WER final acima do limite: %f", tm.Linhas[1].WER)
	}
}

func TestWhisperEsgotaRegravacoes(t *testing.T) {
	s := &sintFalso{nome: "falso"}
	tr := &transcritorFalso{s: s, ruins: map[string]int{"o PIB sobe": 99}}
	dir := t.TempDir()
	m, sc := roteiroTeste()
	_, err := etapaCom(s, tr, 3).Executar(context.Background(), dir, m, sc)
	if !errors.Is(err, ErrWERAcimaLimite) || !strings.Contains(err.Error(), "linha 1") || !strings.Contains(err.Error(), "100.0%") {
		t.Fatalf("esperava ErrWERAcimaLimite citando linha e WER: %v", err)
	}
	if s.chamadas != 5 {
		t.Fatalf("esperava 1 + 4 sinteses, obteve %d", s.chamadas)
	}
	for _, rel := range []string{"audio/voice.wav", "data/timings.json"} {
		if _, e := os.Stat(filepath.Join(dir, filepath.FromSlash(rel))); e == nil {
			t.Fatalf("%s nao deve existir apos reprovacao", rel)
		}
	}
}

const (
	texto10 = "um dois tres quatro cinco seis sete oito nove dez"
	hip10   = "um dois tres quatro cinco seis sete oito nove dezz"
	texto20 = "a b c d e f g h i j k l m n o p q r s t"
)

func roteiroDe(a, b string) (roteiro.Marca, roteiro.Script) {
	m, _ := roteiroTeste()
	return m, roteiro.Script{Lines: []roteiro.LinhaRoteiro{{Text: a}, {Text: b}}}
}

func TestWhisperLinhaEntre5e15Passa(t *testing.T) {
	s := &sintFalso{nome: "falso"}
	tr := &transcritorFalso{s: s, hip: map[string]string{texto10: hip10}}
	m, sc := roteiroDe(texto10, texto20)
	tm, err := etapaCom(s, tr, 3).Executar(context.Background(), t.TempDir(), m, sc)
	if err != nil {
		t.Fatal(err)
	}
	if s.chamadas != 2 || tm.Linhas[0].Regravacoes != 0 || math.Abs(tm.Linhas[0].WER-0.1) > 0.001 {
		t.Fatalf("10%% por linha deve passar sem regravar: chamadas=%d %+v", s.chamadas, tm.Linhas[0])
	}
	if math.Abs(tm.WERGlobal-1.0/30.0) > 0.001 {
		t.Fatalf("WER global deve ser erros totais sobre palavras totais: %f", tm.WERGlobal)
	}
}

func TestWhisperWERGlobalReprova(t *testing.T) {
	s := &sintFalso{nome: "falso"}
	tr := &transcritorFalso{s: s, hip: map[string]string{texto10: hip10}}
	dir := t.TempDir()
	m, sc := roteiroDe(texto10, texto10)
	_, err := etapaCom(s, tr, 3).Executar(context.Background(), dir, m, sc)
	if !errors.Is(err, ErrWERAcimaLimite) || !strings.Contains(err.Error(), "global") || !strings.Contains(err.Error(), "10.0%") {
		t.Fatalf("esperava ErrWERAcimaLimite citando o WER global: %v", err)
	}
	for _, rel := range []string{"audio/voice.wav", "data/timings.json"} {
		if _, e := os.Stat(filepath.Join(dir, filepath.FromSlash(rel))); e == nil {
			t.Fatalf("%s nao deve existir apos reprovacao global", rel)
		}
	}
}

func TestWhisperParaEmAudioIdentico(t *testing.T) {
	constante := func(int, string) ([]byte, error) { return wavTeste(5, 1200), nil }
	s := &sintFalso{nome: "falso", fn: constante}
	tr := &transcritorFalso{s: s, ruins: map[string]int{"o PIB sobe": 99}}
	m, sc := roteiroTeste()
	_, err := etapaCom(s, tr, 3).Executar(context.Background(), t.TempDir(), m, sc)
	if !errors.Is(err, ErrWERAcimaLimite) || !strings.Contains(err.Error(), "identico") || !strings.Contains(err.Error(), "linha 1") {
		t.Fatalf("esperava parada por audio identico: %v", err)
	}
	if s.chamadas != 3 {
		t.Fatalf("esperava 1 sintese da linha 0 e 2 da linha 1 (a regravacao identica encerra), obteve %d", s.chamadas)
	}
}

func TestWhisperSemRegravacoesPermitidas(t *testing.T) {
	s := &sintFalso{nome: "falso"}
	tr := &transcritorFalso{s: s, ruins: map[string]int{"a conta": 1}}
	m, sc := roteiroTeste()
	if _, err := etapaCom(s, tr, 0).Executar(context.Background(), t.TempDir(), m, sc); !errors.Is(err, ErrWERAcimaLimite) || s.chamadas != 1 {
		t.Fatalf("com limite 0 a primeira reprova deve bloquear: err=%v chamadas=%d", err, s.chamadas)
	}
}

func TestWhisperTrocaDeProvedorMantemContrato(t *testing.T) {
	m, sc := roteiroTeste()
	var resultados []Timings
	for _, nome := range []string{"p1", "p2"} {
		s := &sintFalso{nome: nome}
		m.Voz.Motor = nome
		tm, err := etapaCom(s, &transcritorFalso{s: s}, 3).Executar(context.Background(), t.TempDir(), m, sc)
		if err != nil {
			t.Fatal(err)
		}
		resultados = append(resultados, tm)
	}
	a, b := resultados[0], resultados[1]
	if a.Provedor == b.Provedor || a.SchemaVersion != b.SchemaVersion || len(a.Linhas) != len(b.Linhas) {
		t.Fatalf("contrato deve ser o mesmo com provedores distintos: %+v %+v", a, b)
	}
	ja, _ := json.Marshal(a.Linhas[0].Palavras)
	jb, _ := json.Marshal(b.Linhas[0].Palavras)
	if string(ja) != string(jb) {
		t.Fatalf("alinhamento deve vir so do transcritor: %s vs %s", ja, jb)
	}
}

func TestWhisperErros(t *testing.T) {
	m, sc := roteiroTeste()
	t.Run("transcritor falha", func(t *testing.T) {
		s := &sintFalso{nome: "falso"}
		_, err := etapaCom(s, &transcritorFalso{s: s, erro: errors.New("modelo ausente")}, 3).Executar(context.Background(), t.TempDir(), m, sc)
		if err == nil || !strings.Contains(err.Error(), "linha 0") {
			t.Fatalf("esperava erro citando a linha: %v", err)
		}
	})
	t.Run("sintese falha", func(t *testing.T) {
		s := &sintFalso{nome: "falso", fn: func(int, string) ([]byte, error) { return nil, errors.New("rede") }}
		if _, err := etapaCom(s, &transcritorFalso{s: s}, 3).Executar(context.Background(), t.TempDir(), m, sc); err == nil {
			t.Fatal("esperava erro")
		}
	})
	t.Run("gate de creditos nega", func(t *testing.T) {
		s := &sintFalso{nome: "falso", pago: true}
		e := etapaCom(s, &transcritorFalso{s: s}, 3)
		if _, err := e.Executar(context.Background(), t.TempDir(), m, sc); !errors.Is(err, ErrCreditosNegados) || s.chamadas != 0 {
			t.Fatalf("esperava ErrCreditosNegados sem chamada: %v", err)
		}
	})
	t.Run("provedor desconhecido", func(t *testing.T) {
		s := &sintFalso{nome: "falso"}
		mm := m
		mm.Voz.Motor = "outro"
		if _, err := etapaCom(s, &transcritorFalso{s: s}, 3).Executar(context.Background(), t.TempDir(), mm, sc); !errors.Is(err, ErrProvedorDesconhecido) {
			t.Fatalf("esperava ErrProvedorDesconhecido: %v", err)
		}
	})
	t.Run("contexto cancelado", func(t *testing.T) {
		s := &sintFalso{nome: "falso"}
		ctx, cancela := context.WithCancel(context.Background())
		cancela()
		if _, err := etapaCom(s, &transcritorFalso{s: s}, 3).Executar(ctx, t.TempDir(), m, sc); !errors.Is(err, context.Canceled) {
			t.Fatalf("esperava context.Canceled: %v", err)
		}
	})
	t.Run("sem transcritor", func(t *testing.T) {
		s := &sintFalso{nome: "falso"}
		e := EtapaVoz{Locutor: locutorCom(s, nil), ReprovasMax: 3}
		if _, err := e.Executar(context.Background(), t.TempDir(), m, sc); err == nil {
			t.Fatal("esperava erro sem transcritor")
		}
	})
}

func TestCalcularWER(t *testing.T) {
	casos := []struct {
		ref, hip string
		quer     float64
	}{
		{"a conta publica", "a conta publica", 0},
		{"A conta, PÚblica!", "a conta publica", 0.3333},
		{"a conta publica sobe hoje", "a conta publica sobe amanha", 0.2},
		{"", "", 0},
		{"", "algo", 1},
		{"oi", "oi oi oi oi", 1},
		{"quatro palavras aqui mesmo", "", 1},
		{"pê i bê", "PÊ I BÊ", 0},
	}
	for _, c := range casos {
		if got := CalcularWER(c.ref, c.hip); math.Abs(got-c.quer) > 0.001 {
			t.Errorf("WER(%q, %q) = %f, quer %f", c.ref, c.hip, got, c.quer)
		}
	}
}

func TestAjustarPalavras(t *testing.T) {
	brutas := []PalavraTranscrita{
		{Texto: "a", InicioS: -0.2, FimS: 0.4},
		{Texto: "  ", InicioS: 0.4, FimS: 0.5},
		{Texto: "b", InicioS: 0.3, FimS: 0.6},
		{Texto: "c", InicioS: 0.9, FimS: 5},
	}
	got := ajustarPalavras(brutas, 10, 1)
	if len(got) != 3 {
		t.Fatalf("palavra vazia deve sair: %+v", got)
	}
	cursor := 10.0
	for _, p := range got {
		if p.InicioS < cursor || p.FimS < p.InicioS || p.FimS > 11 {
			t.Fatalf("palavra fora da janela ou sobreposta: %+v", p)
		}
		cursor = p.FimS
	}
}

type capFalso struct {
	saida    string
	erro     error
	args     []string
	conteudo []byte
}

func (c *capFalso) Capturar(_ context.Context, nome string, args ...string) ([]byte, error) {
	c.args = append([]string{nome}, args...)
	c.conteudo, _ = os.ReadFile(args[2])
	return []byte(c.saida), c.erro
}

func TestFasterWhisper(t *testing.T) {
	cap := &capFalso{saida: `[{"w":" olá","s":0.1,"e":0.4},{"w":"mundo","s":0.5,"e":0.9}]`}
	fw := FasterWhisper{Exec: cap}
	wav := wavTeste(1, 100)
	palavras, err := fw.Transcrever(context.Background(), wav, "pt-BR")
	if err != nil {
		t.Fatal(err)
	}
	if len(palavras) != 2 || palavras[0].Texto != "olá" || palavras[1].FimS != 0.9 {
		t.Fatalf("palavras inesperadas: %+v", palavras)
	}
	if cap.args[0] != "python" || cap.args[1] != "-c" || cap.args[4] != "pt" || cap.args[5] != "small" {
		t.Fatalf("argumentos inesperados: %v", cap.args[:1])
	}
	if string(cap.conteudo) != string(wav) {
		t.Fatal("o WAV deve estar em arquivo temporario durante a chamada")
	}
	fw.Modelo = "base"
	_, _ = fw.Transcrever(context.Background(), wav, "en")
	if cap.args[4] != "en" || cap.args[5] != "base" {
		t.Fatalf("idioma e modelo devem ser repassados: %v", cap.args[4:])
	}
}

func TestFasterWhisperFalhas(t *testing.T) {
	casos := map[string]*capFalso{
		"executor falha": {erro: errors.New("python ausente")},
		"json invalido":  {saida: "Traceback..."},
		"tempo invalido": {saida: `[{"w":"a","s":2,"e":1}]`},
	}
	for nome, c := range casos {
		t.Run(nome, func(t *testing.T) {
			if _, err := (FasterWhisper{Exec: c}).Transcrever(context.Background(), wavTeste(1, 10), "pt"); err == nil {
				t.Fatal("esperava erro")
			}
		})
	}
}

func TestCapturadorSistema(t *testing.T) {
	var c Capturador = execSistema{}
	saida, err := c.Capturar(context.Background(), "go", "env", "GOOS")
	if err != nil || strings.TrimSpace(string(saida)) == "" {
		t.Fatalf("go env deve responder: %q %v", saida, err)
	}
	if _, err := c.Capturar(context.Background(), "go", "comando-inexistente-xyz"); err == nil {
		t.Fatal("comando invalido deve falhar")
	}
	if _, err := c.Capturar(context.Background(), "programa-que-nao-existe-xyz"); err == nil {
		t.Fatal("programa inexistente deve falhar")
	}
}
