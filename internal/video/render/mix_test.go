package render

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const jsonLoudnormTeste = `[Parsed_loudnorm_0 @ 000001a3]
{
	"input_i" : "-17.78",
	"input_tp" : "-4.21",
	"input_lra" : "2.50",
	"input_thresh" : "-27.92",
	"output_i" : "-14.00",
	"output_tp" : "-1.50",
	"output_lra" : "2.30",
	"output_thresh" : "-24.10",
	"normalization_type" : "dynamic",
	"target_offset" : "-0.02"
}
`

func saidaEbur(i, pico float64) string {
	return fmt.Sprintf(`[Parsed_ebur128_0 @ 0000] Summary:

  Integrated loudness:
    I:         %.1f LUFS
    Threshold: -24.5 LUFS

  Loudness range:
    LRA:         2.1 LU
    Threshold: -34.5 LUFS
    LRA low:   -16.0 LUFS
    LRA high:  -13.9 LUFS

  True peak:
    Peak:       %.1f dBFS
`, i, pico)
}

// execMixFalso grava WAV estereo de 8 kHz com a duracao pedida em -t (ou 1 s).
type execMixFalso struct {
	comandos   [][]string
	erro       error
	fatorDur   float64 // multiplica a duracao das saidas norm*.wav
	lixoNaNorm bool    // grava bytes que nao sao WAV nas saidas norm*.wav
}

func (e *execMixFalso) Executar(_ context.Context, nome string, args ...string) error {
	e.comandos = append(e.comandos, append([]string{nome}, args...))
	if e.erro != nil {
		return e.erro
	}
	dur := 1.0
	for _, a := range args {
		if i := strings.Index(a, "atrim=end_sample="); i >= 0 {
			resto := a[i+len("atrim=end_sample="):]
			fim := strings.IndexFunc(resto, func(r rune) bool { return r < '0' || r > '9' })
			if fim < 0 {
				fim = len(resto)
			}
			n, _ := strconv.Atoi(resto[:fim])
			dur = float64(n) / 48000
		}
	}
	saida := args[len(args)-1]
	if strings.HasPrefix(filepath.Base(saida), "norm") {
		if e.lixoNaNorm {
			return os.WriteFile(saida, []byte("lixo"), 0o644)
		}
		if e.fatorDur != 0 {
			dur *= e.fatorDur
		}
	}
	f := formatoWAV{Canais: 2, Taxa: 8000, Bits: 16}
	return os.WriteFile(saida, montarWAV(f, make([]byte, int(dur*8000)*4)), 0o644)
}

func (e *execMixFalso) com(trecho string) []string {
	var achados []string
	for _, c := range e.comandos {
		if j := strings.Join(c, " "); strings.Contains(j, trecho) {
			achados = append(achados, j)
		}
	}
	return achados
}

type medFalso struct {
	json     string
	medicoes []string
	chamadas []string
	erro     error
}

func (m *medFalso) Saida(_ context.Context, _ string, args ...string) (string, error) {
	juntos := strings.Join(args, " ")
	m.chamadas = append(m.chamadas, juntos)
	if m.erro != nil {
		return "", m.erro
	}
	switch {
	case strings.Contains(juntos, "print_format=json"):
		return m.json, nil
	case strings.Contains(juntos, "ebur128"):
		if len(m.medicoes) == 0 {
			return "", errors.New("sem medicao programada")
		}
		s := m.medicoes[0]
		if len(m.medicoes) > 1 {
			m.medicoes = m.medicoes[1:]
		}
		return s, nil
	}
	return "", nil
}

func projetoMix(t *testing.T, extra map[string]string) (string, Timeline) {
	t.Helper()
	return projetoMixDur(t, 6, extra)
}

func projetoMixDur(t *testing.T, dur float64, extra map[string]string) (string, Timeline) {
	t.Helper()
	tl := timelineDe(batidaObj(0, 0, dur/2, "cen_01", "pan", true), batidaObj(1, dur/2, dur, "cen_02", "pan", true))
	b, _ := json.MarshalIndent(tl, "", "  ")
	voz := montarWAV(formatoWAV{Canais: 1, Taxa: 24000, Bits: 16}, make([]byte, int(math.Round(dur*24000))*2))
	arqs := map[string]string{"timeline.json": string(b) + "\n", "audio/voice.wav": string(voz)}
	for k, v := range extra {
		arqs[k] = v
	}
	return projetoFixture(t, arqs), tl
}

func mixCom(ex *execMixFalso, med *medFalso) MixDeps { return MixDeps{Exec: ex, Med: med} }

func lerRelatorioMix(t *testing.T, dir string) RelatorioMix {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, "data", "mix.json"))
	if err != nil {
		t.Fatalf("data/mix.json ausente: %v", err)
	}
	var r RelatorioMix
	if err := json.Unmarshal(b, &r); err != nil {
		t.Fatal(err)
	}
	return r
}

func naoGravouNada(t *testing.T, dir string, timelineAntes []byte) {
	t.Helper()
	for _, rel := range []string{"audio/soundtrack.wav", "data/mix.json"} {
		if existe(dir, rel) {
			t.Fatalf("%s nao pode existir apos falha", rel)
		}
	}
	depois, _ := os.ReadFile(filepath.Join(dir, "timeline.json"))
	if string(depois) != string(timelineAntes) {
		t.Fatal("timeline.json nao pode mudar apos falha")
	}
}

func TestMixAprovaNaPrimeiraComMargemDePico(t *testing.T) {
	dir, _ := projetoMix(t, nil)
	ex := &execMixFalso{}
	med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.1, -2.0)}}
	rel, err := ExecutarMix(context.Background(), dir, mixCom(ex, med))
	if err != nil {
		t.Fatal(err)
	}
	norm := ex.com("loudnorm=")
	if len(norm) != 1 || !strings.Contains(norm[0], ":TP=-2.0:") || !strings.Contains(norm[0], "measured_I=-17.78") || !strings.Contains(norm[0], "measured_TP=-4.21") || !strings.Contains(norm[0], "linear=true") {
		t.Fatalf("normalizacao deve mirar -2,0 dBTP (0,5 dB abaixo do teto) com os valores medidos: %v", norm)
	}
	if len(rel.Tentativas) != 1 || !rel.Tentativas[0].Aprovada || rel.LUFS != -14.1 || rel.TruePeakDBTP != -2.0 || rel.Eventos != 4 {
		t.Fatalf("relatorio inesperado: %+v", rel)
	}
	if got := lerRelatorioMix(t, dir); got.Eventos != 4 || got.SchemaVersion != 1 {
		t.Fatalf("data/mix.json diferente do retorno: %+v", got)
	}
	b, err := os.ReadFile(filepath.Join(dir, "audio", "soundtrack.wav"))
	if err != nil {
		t.Fatal(err)
	}
	if f, dados, err := lerWAV(b); err != nil || math.Abs(duracaoWAV(f, len(dados))-6) > 1e-9 {
		t.Fatalf("soundtrack deve durar exatamente a timeline (6 s): %v", err)
	}
	var tl Timeline
	raw, _ := os.ReadFile(filepath.Join(dir, "timeline.json"))
	if err := json.Unmarshal(raw, &tl); err != nil || tl.Validar() != nil || len(tl.Batidas[0].SFX) != 2 || len(tl.Batidas[1].SFX) != 2 {
		t.Fatalf("timeline.json deve receber os eventos de SFX e validar: %v", err)
	}
}

func TestMixGrafoDaMixagemInicial(t *testing.T) {
	dir, _ := projetoMix(t, nil)
	ex := &execMixFalso{}
	med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.0, -2.0)}}
	if _, err := ExecutarMix(context.Background(), dir, mixCom(ex, med)); err != nil {
		t.Fatal(err)
	}
	mix := ex.com("sidechaincompress")
	if len(mix) != 1 {
		t.Fatalf("esperava 1 mixagem inicial: %d", len(mix))
	}
	for _, esperado := range []string{"-stream_loop -1", "pan=stereo|c0=c0|c1=c0", "asplit=2", "adelay=0|0", "adelay=250|250", "adelay=3000|3000", "adelay=3250|3250", "amix=inputs=6:normalize=0", "pcm_f32le"} {
		if !strings.Contains(mix[0], esperado) {
			t.Errorf("mixagem sem %q:\n%s", esperado, mix[0])
		}
	}
}

func TestMixRefazComMaisMargemQuandoOPicoReprova(t *testing.T) {
	dir, _ := projetoMix(t, nil)
	ex := &execMixFalso{}
	med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.0, -1.4), saidaEbur(-14.0, -2.5)}}
	rel, err := ExecutarMix(context.Background(), dir, mixCom(ex, med))
	if err != nil {
		t.Fatal(err)
	}
	if len(rel.Tentativas) != 2 || rel.Tentativas[0].Aprovada || !rel.Tentativas[1].Aprovada {
		t.Fatalf("esperava reprovar a 1a (pico -1,4) e aprovar a 2a: %+v", rel.Tentativas)
	}
	norm := ex.com("loudnorm=")
	if len(norm) != 2 || !strings.Contains(norm[0], ":TP=-2.0:") || !strings.Contains(norm[1], ":TP=-2.5:") {
		t.Fatalf("a 2a tentativa deve ter mais 0,5 dB de margem: %v", norm)
	}
	if rel.TruePeakDBTP != -2.5 {
		t.Fatalf("relatorio deve trazer a medicao aprovada: %+v", rel)
	}
}

func TestMixCompensaOLUFSMedido(t *testing.T) {
	dir, _ := projetoMix(t, nil)
	ex := &execMixFalso{}
	med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.7, -2.0), saidaEbur(-14.0, -2.5)}}
	if _, err := ExecutarMix(context.Background(), dir, mixCom(ex, med)); err != nil {
		t.Fatal(err)
	}
	norm := ex.com("loudnorm=")
	if len(norm) != 2 || !strings.Contains(norm[1], "loudnorm=I=-13.30:") {
		t.Fatalf("LUFS medido -14,7 deve subir o alvo para -13,3: %v", norm)
	}
}

func TestMixEsgotaTentativasENaoGravaNada(t *testing.T) {
	dir, _ := projetoMix(t, nil)
	antes, _ := os.ReadFile(filepath.Join(dir, "timeline.json"))
	ex := &execMixFalso{}
	med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.0, -1.4)}}
	_, err := ExecutarMix(context.Background(), dir, mixCom(ex, med))
	if !errors.Is(err, ErrLoudnessReprovado) || !strings.Contains(err.Error(), "-1.4") {
		t.Fatalf("esperava ErrLoudnessReprovado citando o pico medido: %v", err)
	}
	if n := len(ex.com("loudnorm=")); n != MaxTentativasMix {
		t.Fatalf("esperava %d tentativas, obteve %d", MaxTentativasMix, n)
	}
	naoGravouNada(t, dir, antes)
}

func TestMixSFXSemArquivoBloqueiaSemGravar(t *testing.T) {
	dir, _ := projetoMix(t, map[string]string{
		"render.json": `{"schema_version":1,"sfx_catalogo":"cat.json"}`,
		"cat.json":    `{"movimento_objeto":"audio/sfx/movimento_objeto.wav"}`,
	})
	antes, _ := os.ReadFile(filepath.Join(dir, "timeline.json"))
	_, err := ExecutarMix(context.Background(), dir, mixCom(&execMixFalso{}, &medFalso{json: jsonLoudnormTeste}))
	if !errors.Is(err, ErrSFXSemArquivo) {
		t.Fatalf("evento sem arquivo no catalogo deve bloquear: %v", err)
	}
	naoGravouNada(t, dir, antes)
}

func TestMixEntradasInvalidas(t *testing.T) {
	casos := map[string]struct {
		extra map[string]string
		apaga string
	}{
		"voz ausente":       {apaga: "audio/voice.wav"},
		"timeline ausente":  {apaga: "timeline.json"},
		"timeline invalida": {extra: map[string]string{"timeline.json": "{"}},
		"voz de 3 segundos": {extra: map[string]string{"audio/voice.wav": string(montarWAV(formatoWAV{Canais: 1, Taxa: 24000, Bits: 16}, make([]byte, 24000*2*3)))}},
		"voz com 3 canais":  {extra: map[string]string{"audio/voice.wav": string(montarWAV(formatoWAV{Canais: 3, Taxa: 24000, Bits: 16}, make([]byte, 24000*6*6)))}},
	}
	for nome, c := range casos {
		t.Run(nome, func(t *testing.T) {
			dir, _ := projetoMix(t, c.extra)
			if c.apaga != "" {
				_ = os.Remove(filepath.Join(dir, filepath.FromSlash(c.apaga)))
			}
			antes, _ := os.ReadFile(filepath.Join(dir, "timeline.json"))
			ex := &execMixFalso{}
			if _, err := ExecutarMix(context.Background(), dir, mixCom(ex, &medFalso{json: jsonLoudnormTeste})); err == nil {
				t.Fatal("esperava erro")
			}
			if existe(dir, "audio/soundtrack.wav") || existe(dir, "data/mix.json") {
				t.Fatal("nada pode ser gravado")
			}
			if depois, _ := os.ReadFile(filepath.Join(dir, "timeline.json")); c.apaga != "timeline.json" && string(depois) != string(antes) {
				t.Fatal("timeline.json nao pode mudar")
			}
		})
	}
}

func TestMixFalhasDeFerramenta(t *testing.T) {
	casos := map[string]struct {
		ex  *execMixFalso
		med *medFalso
	}{
		"ffmpeg falha":              {&execMixFalso{erro: errors.New("sem ffmpeg")}, &medFalso{json: jsonLoudnormTeste}},
		"medidor falha":             {&execMixFalso{}, &medFalso{erro: errors.New("sem ffmpeg")}},
		"json de loudnorm invalido": {&execMixFalso{}, &medFalso{json: "sem json"}},
		"medicao ebur128 ilegivel":  {&execMixFalso{}, &medFalso{json: jsonLoudnormTeste, medicoes: []string{"nada util"}}},
	}
	for nome, c := range casos {
		t.Run(nome, func(t *testing.T) {
			dir, _ := projetoMix(t, nil)
			antes, _ := os.ReadFile(filepath.Join(dir, "timeline.json"))
			if _, err := ExecutarMix(context.Background(), dir, mixCom(c.ex, c.med)); err == nil {
				t.Fatal("esperava erro")
			}
			naoGravouNada(t, dir, antes)
		})
	}
	dir, _ := projetoMix(t, nil)
	ctx, cancela := context.WithCancel(context.Background())
	cancela()
	if _, err := ExecutarMix(ctx, dir, mixCom(&execMixFalso{}, &medFalso{})); !errors.Is(err, context.Canceled) {
		t.Fatalf("contexto cancelado deve interromper: %v", err)
	}
}

func TestMixParsers(t *testing.T) {
	m, err := parseEBUR128(saidaEbur(-14.3, -2.1))
	if err != nil || m.LUFS != -14.3 || m.TruePeakDBTP != -2.1 {
		t.Fatalf("ebur128: %+v %v", m, err)
	}
	dois := saidaEbur(-30, -9) + "\n" + saidaEbur(-14.2, -2.2)
	if m, _ := parseEBUR128(dois); m.LUFS != -14.2 {
		t.Fatalf("deve ler o ultimo Summary: %+v", m)
	}
	for _, ruim := range []string{"", "sem resumo", "Summary:\n I: abc LUFS", "Summary:\n I: -14.0 LUFS\n"} {
		if _, err := parseEBUR128(ruim); err == nil {
			t.Errorf("%q deveria falhar", ruim)
		}
	}
	l, err := parseLoudnormJSON(jsonLoudnormTeste)
	if err != nil || l.InputI != "-17.78" || l.InputTP != "-4.21" || l.InputLRA != "2.50" || l.InputThresh != "-27.92" || l.TargetOffset != "-0.02" {
		t.Fatalf("loudnorm json: %+v %v", l, err)
	}
	for _, ruim := range []string{"", "{", "{\"input_i\":\"\"}"} {
		if _, err := parseLoudnormJSON(ruim); err == nil {
			t.Errorf("%q deveria falhar", ruim)
		}
	}
}

func TestMixParametrosDasTentativas(t *testing.T) {
	if i, tp := parametrosTentativa(0, 0, 0); i != -14 || tp != -2.0 {
		t.Fatalf("1a tentativa: I=%v TP=%v", i, tp)
	}
	if i, tp := parametrosTentativa(1, -14.7, -14); math.Abs(i-(-13.3)) > 1e-9 || tp != -2.5 {
		t.Fatalf("2a tentativa: I=%v TP=%v", i, tp)
	}
	if i, tp := parametrosTentativa(2, -13.9, -13.3); math.Abs(i-(-13.4)) > 1e-9 || tp != -3.0 {
		t.Fatalf("3a tentativa: I=%v TP=%v", i, tp)
	}
}

func TestMixMedidorDoSistema(t *testing.T) {
	var m Medidor = execSistema{}
	saida, err := m.Saida(context.Background(), "go", "version")
	if err != nil || !strings.Contains(saida, "go version") {
		t.Fatalf("go version: %q %v", saida, err)
	}
	if _, err := m.Saida(context.Background(), "go", "comando-inexistente-xyz"); err == nil {
		t.Fatal("comando invalido deve falhar")
	}
}

func TestMixVozEstereoNaoUsaPan(t *testing.T) {
	estereo := montarWAV(formatoWAV{Canais: 2, Taxa: 24000, Bits: 16}, make([]byte, 24000*4*6))
	dir, _ := projetoMix(t, map[string]string{"audio/voice.wav": string(estereo)})
	ex := &execMixFalso{}
	med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.0, -2.0)}}
	if _, err := ExecutarMix(context.Background(), dir, mixCom(ex, med)); err != nil {
		t.Fatal(err)
	}
	mix := ex.com("sidechaincompress")
	if len(mix) != 1 || strings.Contains(mix[0], "pan=stereo") || !strings.Contains(mix[0], "[0:a]aresample=48000,aformat=channel_layouts=stereo,asplit=2") {
		t.Fatalf("voz estereo nao deve passar por pan: %v", mix)
	}
}

func TestMixUsaBGMDoRenderJSON(t *testing.T) {
	dir, _ := projetoMix(t, map[string]string{"render.json": `{"schema_version":1,"bgm":"musica/minha.wav"}`, "musica/minha.wav": string(wavTeste(1, 100))})
	ex := &execMixFalso{}
	med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.0, -2.0)}}
	if _, err := ExecutarMix(context.Background(), dir, mixCom(ex, med)); err != nil {
		t.Fatal(err)
	}
	mix := ex.com("sidechaincompress")
	if len(mix) != 1 || !strings.Contains(mix[0], "-stream_loop -1 -i "+filepath.Join(dir, "musica", "minha.wav")) {
		t.Fatalf("deve usar a BGM do projeto: %v", mix)
	}
	if existe(dir, "audio/bgm.wav") {
		t.Fatal("nao deve sintetizar a BGM padrao quando o projeto traz a sua")
	}
}

func TestMixSaidasInconsistentesNaoGravamNada(t *testing.T) {
	casos := map[string]struct {
		ex     *execMixFalso
		trecho error
	}{
		"soundtrack mais curto que a timeline": {&execMixFalso{fatorDur: 0.5}, ErrRenderInvalido},
		"soundtrack que nao e wav":             {&execMixFalso{lixoNaNorm: true}, ErrAudioInvalido},
	}
	for nome, c := range casos {
		t.Run(nome, func(t *testing.T) {
			dir, _ := projetoMix(t, nil)
			antes, _ := os.ReadFile(filepath.Join(dir, "timeline.json"))
			med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.0, -2.0)}}
			if _, err := ExecutarMix(context.Background(), dir, mixCom(c.ex, med)); !errors.Is(err, c.trecho) {
				t.Fatalf("esperava %v, obteve %v", c.trecho, err)
			}
			naoGravouNada(t, dir, antes)
		})
	}
}

func TestMixMedicaoFinalIndisponivelEntradaSemBatidas(t *testing.T) {
	dir, _ := projetoMix(t, nil)
	antes, _ := os.ReadFile(filepath.Join(dir, "timeline.json"))
	if _, err := ExecutarMix(context.Background(), dir, mixCom(&execMixFalso{}, &medFalso{json: jsonLoudnormTeste})); err == nil {
		t.Fatal("sem medicao do arquivo final deve falhar")
	}
	naoGravouNada(t, dir, antes)
	vazia, _ := projetoMix(t, map[string]string{"timeline.json": `{"schema_version":1}`})
	if _, err := ExecutarMix(context.Background(), vazia, mixCom(&execMixFalso{}, &medFalso{})); !errors.Is(err, ErrEntradaTimeline) {
		t.Fatalf("timeline sem batidas deve falhar com erro tipado: %v", err)
	}
}

func TestMixFalhaAoGravarNaoDeixaNada(t *testing.T) {
	dir, _ := projetoMix(t, map[string]string{"data": "arquivo no lugar da pasta"})
	antes, _ := os.ReadFile(filepath.Join(dir, "timeline.json"))
	med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.0, -2.0)}}
	if _, err := ExecutarMix(context.Background(), dir, mixCom(&execMixFalso{}, med)); err == nil {
		t.Fatal("falha ao gravar deve ser erro")
	}
	naoGravouNada(t, dir, antes)
	sobras, _ := filepath.Glob(filepath.Join(dir, "*", "*.tmp"))
	mais, _ := filepath.Glob(filepath.Join(dir, "*.tmp"))
	if len(sobras)+len(mais) != 0 {
		t.Fatalf("nao pode sobrar temporario: %v %v", sobras, mais)
	}
}

func TestMixDepsPadrao(t *testing.T) {
	d := normalizarMixDeps(MixDeps{})
	if d.Exec == nil || d.Med == nil {
		t.Fatal("sem dependencias informadas, usa o executor do sistema")
	}
}

func TestGravacaoEmLote(t *testing.T) {
	dir := t.TempDir()
	ok := []saidaArquivo{{"a/um.txt", []byte("1")}, {"b/dois.txt", []byte("2")}}
	if err := gravarEmLote(dir, ok); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "b", "dois.txt")); string(b) != "2" || existe(dir, "a/um.txt.tmp") {
		t.Fatal("todos gravados e sem temporarios")
	}
	if err := gravarEmLote(dir, []saidaArquivo{{"a/um.txt", []byte("novo")}}); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "a", "um.txt")); string(b) != "novo" {
		t.Fatal("deve substituir arquivo existente")
	}
	if err := os.WriteFile(filepath.Join(dir, "c"), []byte("arquivo"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := gravarEmLote(dir, []saidaArquivo{{"d/tres.txt", []byte("3")}, {"c/quatro.txt", []byte("4")}})
	if err == nil || existe(dir, "d/tres.txt") || existe(dir, "d/tres.txt.tmp") {
		t.Fatalf("se o 2o falhar, o 1o nao pode ficar gravado nem como temporario: %v", err)
	}
	if err := gravarEmLote(dir, []saidaArquivo{{"../fora.txt", []byte("x")}}); err == nil {
		t.Fatal("caminho fora da raiz deve falhar")
	}
}

func TestTimelineFalhaAoGravarNaoDeixaNada(t *testing.T) {
	dir := projetoTimeline(t, map[string]string{"renders": "arquivo no lugar da pasta"})
	if _, err := ExecutarTimeline(context.Background(), dir); err == nil {
		t.Fatal("falha ao gravar deve ser erro")
	}
	if existe(dir, "timeline.json") {
		t.Fatal("timeline.json nao pode existir se o SRT nao foi gravado")
	}
}

func vozSenoide(segundos float64) []byte {
	n := int(math.Round(24000 * segundos))
	b := make([]byte, n*2)
	for i := 0; i < n; i++ {
		v := int16(8000 * math.Sin(2*math.Pi*440*float64(i)/24000))
		binary.LittleEndian.PutUint16(b[2*i:], uint16(v))
	}
	return montarWAV(formatoWAV{Canais: 1, Taxa: 24000, Bits: 16}, b)
}

// Teste com o ffmpeg real: a duracao do soundtrack deve coincidir com a da timeline em ate 1 ms.
func TestMixComFFmpegReal(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg ausente")
	}
	dir, _ := projetoMix(t, map[string]string{"audio/voice.wav": string(vozSenoide(6))})
	rel, err := ExecutarMix(context.Background(), dir, MixDeps{})
	if err != nil {
		t.Fatalf("mixagem real: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "audio", "soundtrack.wav"))
	if err != nil {
		t.Fatal(err)
	}
	f, dados, err := lerWAV(b)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(dados) / int(f.Canais*2); n != 288000 {
		t.Fatalf("soundtrack tem %d amostras e a timeline pede 288000", n)
	}
	if rel.LUFS < -14.5 || rel.LUFS > -13.5 || rel.TruePeakDBTP > -1.5 {
		t.Fatalf("medicao fora do alvo: %.1f LUFS, %.1f dBTP", rel.LUFS, rel.TruePeakDBTP)
	}
}

func TestMixCortaPorAmostraENaoPorTempo(t *testing.T) {
	dir, _ := projetoMix(t, nil)
	ex := &execMixFalso{}
	med := &medFalso{json: jsonLoudnormTeste, medicoes: []string{saidaEbur(-14.0, -2.0)}}
	if _, err := ExecutarMix(context.Background(), dir, mixCom(ex, med)); err != nil {
		t.Fatal(err)
	}
	norm := ex.com("loudnorm=")
	if len(norm) != 1 || !strings.Contains(norm[0], "apad=pad_dur=1,atrim=end_sample=288000") || strings.Contains(norm[0], " -t ") {
		t.Fatalf("o loudnorm desloca os timestamps em 28 ms, entao o corte deve ser por amostra e nunca por -t: %v", norm)
	}
	mix := ex.com("sidechaincompress")
	if len(mix) != 1 || !strings.Contains(mix[0], "duration=longest,atrim=end_sample=288000[mix]") || !strings.Contains(mix[0], " -t 7.000000 ") {
		t.Fatalf("mixagem corta por amostra e usa -t so como trava de seguranca: %v", mix)
	}
}

// Duracao que nao e multiplo de 0,1 s: reproduz o defeito de 28 ms observado no demo.
func TestMixComFFmpegRealDuracaoQuebrada(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg ausente")
	}
	dir, _ := projetoMixDur(t, 5.872, map[string]string{"audio/voice.wav": string(vozSenoide(5.872))})
	if _, err := ExecutarMix(context.Background(), dir, MixDeps{}); err != nil {
		t.Fatalf("mixagem real: %v", err)
	}
	b, _ := os.ReadFile(filepath.Join(dir, "audio", "soundtrack.wav"))
	f, dados, err := lerWAV(b)
	if err != nil {
		t.Fatal(err)
	}
	quer := int(math.Round(5.872 * 48000))
	if n := len(dados) / int(f.Canais*2); n != quer {
		t.Fatalf("soundtrack tem %d amostras e a timeline pede %d (diferenca de %.1f ms)", n, quer, float64(n-quer)/48)
	}
}
