package render

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

type sondaFalsa struct {
	lista []string
	erro  error
}

func (s sondaFalsa) Encoders(context.Context) ([]string, error) { return s.lista, s.erro }

func TestEncoderEscolhaPorOrdem(t *testing.T) {
	casos := []struct {
		lista []string
		quer  string
	}{
		{[]string{"libx264", "h264_amf", "h264_qsv", "h264_nvenc"}, "h264_nvenc"},
		{[]string{"libx264", "h264_amf", "h264_qsv"}, "h264_qsv"},
		{[]string{"libx264", "h264_amf"}, "h264_amf"},
		{[]string{"libx264"}, "libx264"},
		{nil, "libx264"},
	}
	for _, c := range casos {
		if got := EscolherEncoder(context.Background(), sondaFalsa{lista: c.lista}); got != c.quer {
			t.Errorf("lista %v: esperava %s, obteve %s", c.lista, c.quer, got)
		}
	}
}

func TestEncoderSondaComErroCaiNoLibx264(t *testing.T) {
	if got := EscolherEncoder(context.Background(), sondaFalsa{erro: errors.New("sem ffmpeg")}); got != EncoderFallback {
		t.Fatalf("erro na sonda deve cair em libx264, obteve %s", got)
	}
}

type medEncFalso struct {
	saida string
	err   error
}

func (m medEncFalso) Saida(context.Context, string, ...string) (string, error) {
	return m.saida, m.err
}

type execEncFalso struct {
	usados []string
	falham map[string]bool
}

func (e *execEncFalso) Executar(_ context.Context, _ string, args ...string) error {
	enc := args[0]
	e.usados = append(e.usados, enc)
	if e.falham[enc] {
		return errors.New("falha de " + enc)
	}
	return nil
}

func montarEnc(enc string) []string { return []string{enc} }

func TestEncoderHardwareQueFalhaRepeteComLibx264(t *testing.T) {
	ex := &execEncFalso{falham: map[string]bool{"h264_nvenc": true}}
	usado, err := ExecutarComFallback(context.Background(), ex, "h264_nvenc", montarEnc)
	if err != nil || usado != "libx264" {
		t.Fatalf("esperava libx264 sem erro, obteve %s, %v", usado, err)
	}
	if strings.Join(ex.usados, ",") != "h264_nvenc,libx264" {
		t.Fatalf("ordem das tentativas: %v", ex.usados)
	}
}

func TestEncoderHardwareQueFuncionaNaoRepete(t *testing.T) {
	ex := &execEncFalso{}
	usado, err := ExecutarComFallback(context.Background(), ex, "h264_qsv", montarEnc)
	if err != nil || usado != "h264_qsv" || len(ex.usados) != 1 {
		t.Fatalf("obteve %s, %v, %v", usado, err, ex.usados)
	}
}

func TestEncoderLibx264NaoTemOutroFallback(t *testing.T) {
	ex := &execEncFalso{falham: map[string]bool{"libx264": true}}
	usado, err := ExecutarComFallback(context.Background(), ex, "libx264", montarEnc)
	if err == nil || usado != "libx264" || len(ex.usados) != 1 {
		t.Fatalf("libx264 deve falhar uma vez so: %s, %v, %v", usado, err, ex.usados)
	}
}

func TestEncoderFallbackTambemFalhaReportaAmbosErros(t *testing.T) {
	ex := &execEncFalso{falham: map[string]bool{"h264_amf": true, "libx264": true}}
	_, err := ExecutarComFallback(context.Background(), ex, "h264_amf", montarEnc)
	if err == nil || !strings.Contains(err.Error(), "h264_amf") || !strings.Contains(err.Error(), "libx264") {
		t.Fatalf("erro deve citar os dois encoders: %v", err)
	}
}

func TestEncoderContextoCanceladoNaoRepete(t *testing.T) {
	ctx, cancela := context.WithCancel(context.Background())
	cancela()
	ex := &execEncFalso{}
	if _, err := ExecutarComFallback(ctx, ex, "h264_nvenc", montarEnc); !errors.Is(err, context.Canceled) || len(ex.usados) != 0 {
		t.Fatalf("contexto cancelado: %v %v", err, ex.usados)
	}
}

func TestEncoderParseDaListaDoFFmpeg(t *testing.T) {
	saida := "Encoders:\n V..... libx264              libx264 H.264\n V....D h264_nvenc           NVIDIA NVENC H.264\n V..... h264_amf             AMD AMF H.264\n A....D aac                  AAC\n"
	got := parseEncoders(saida)
	if strings.Join(got, ",") != "libx264,h264_nvenc,h264_amf,aac" {
		t.Fatalf("parse: %v", got)
	}
}

func TestEncoderSondaFFmpegFalsa(t *testing.T) {
	s := SondaFFmpeg{Med: medEncFalso{saida: "Encoders:\n V..... h264_qsv  Intel QSV\n"}}
	l, err := s.Encoders(context.Background())
	if err != nil || len(l) != 1 || l[0] != "h264_qsv" {
		t.Fatalf("sonda: %v %v", l, err)
	}
	s = SondaFFmpeg{Med: medEncFalso{err: errors.New("x")}}
	if _, err := s.Encoders(context.Background()); err == nil {
		t.Fatal("erro do ffmpeg deve propagar")
	}
}

func TestEncoderSondaComFFmpegReal(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg ausente")
	}
	l, err := SondaFFmpeg{}.Encoders(context.Background())
	if err != nil || len(l) == 0 {
		t.Fatalf("sonda real: %v %v", l, err)
	}
	if got := EscolherEncoder(context.Background(), SondaFFmpeg{}); got == "" {
		t.Fatal("encoder vazio")
	}
}
