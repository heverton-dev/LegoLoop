package render

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

// PalavraTranscrita e uma palavra ouvida no audio, com tempos relativos ao inicio do trecho.
type PalavraTranscrita struct {
	Texto   string
	InicioS float64
	FimS    float64
}

// Transcritor ouve um WAV e devolve as palavras com tempo; e a unica fonte de alinhamento.
type Transcritor interface {
	Transcrever(ctx context.Context, wav []byte, idioma string) ([]PalavraTranscrita, error)
}

// EtapaVoz executa a Etapa 6: sintese, conferencia por Whisper, regravacao seletiva e timings.json.
type EtapaVoz struct {
	Locutor     Locutor
	Transcritor Transcritor
	ReprovasMax int
}

type linhaVoz struct {
	wav      []byte
	palavras []PalavraTranscrita
	wer      float64
	erros    int
	total    int
	regrav   int
}

// Executar so grava audio/ e data/timings.json depois que todas as linhas passam no gate de WER.
func (e EtapaVoz) Executar(ctx context.Context, projDir string, marca roteiro.Marca, script roteiro.Script) (Timings, error) {
	if e.Transcritor == nil {
		return Timings{}, invalido("transcritor ausente")
	}
	p, err := e.Locutor.preparar(ctx, projDir, marca, script)
	if err != nil {
		return Timings{}, err
	}
	linhas := make([]linhaVoz, 0, len(script.Lines))
	for i, l := range script.Lines {
		lv, err := e.aprovarLinha(ctx, p, script.Idioma, i, AplicarPronuncia(l.Text, p.pron))
		if err != nil {
			return Timings{}, err
		}
		linhas = append(linhas, lv)
	}
	global := werGlobal(linhas)
	if global > LimiteWERGlobal {
		return Timings{}, fmt.Errorf("%w: WER global %.1f%% acima de %.0f%%", ErrWERAcimaLimite, global*100, LimiteWERGlobal*100)
	}
	tm := GerarVisemas(montarTimings(marca.Voz.Motor, p.voz, linhas))
	if err := tm.Validar(); err != nil {
		return tm, err
	}
	return tm, gravarEtapaVoz(projDir, marca.Voz.Motor, p.voz, linhas, tm)
}

// aprovarLinha sintetiza e confere; reprovada, regrava somente esta linha ate ReprovasMax vezes.
// Se a regravacao devolve o mesmo audio, o provedor e deterministico e repetir nao ajuda.
func (e EtapaVoz) aprovarLinha(ctx context.Context, p preparo, idioma string, i int, texto string) (linhaVoz, error) {
	var lv linhaVoz
	var anterior []byte
	for tentativa := 0; tentativa <= e.ReprovasMax; tentativa++ {
		novo, err := e.tentar(ctx, p, idioma, i, texto)
		if err != nil {
			return lv, err
		}
		novo.regrav = tentativa
		lv = novo
		if lv.wer <= LimiteWERLinha {
			return lv, nil
		}
		if anterior != nil && bytes.Equal(anterior, lv.wav) {
			return lv, fmt.Errorf("%w: linha %d com WER %.1f%% e audio regravado identico ao anterior", ErrWERAcimaLimite, i, lv.wer*100)
		}
		anterior = lv.wav
	}
	return lv, fmt.Errorf("%w: linha %d com WER %.1f%% apos %d regravacoes", ErrWERAcimaLimite, i, lv.wer*100, e.ReprovasMax)
}

func (e EtapaVoz) tentar(ctx context.Context, p preparo, idioma string, i int, texto string) (linhaVoz, error) {
	if err := ctx.Err(); err != nil {
		return linhaVoz{}, err
	}
	wav, err := p.s.Sintetizar(ctx, texto, p.op)
	if err != nil {
		return linhaVoz{}, fmt.Errorf("diorama-render: linha %d: %w", i, err)
	}
	if _, _, err := lerWAV(wav); err != nil {
		return linhaVoz{}, fmt.Errorf("diorama-render: linha %d: %w", i, err)
	}
	palavras, err := e.Transcritor.Transcrever(ctx, wav, idioma)
	if err != nil {
		return linhaVoz{}, fmt.Errorf("diorama-render: linha %d: transcricao: %w", i, err)
	}
	erros, total := errosPalavras(texto, textoDe(palavras))
	return linhaVoz{wav: wav, palavras: palavras, wer: razaoWER(erros, total), erros: erros, total: total}, nil
}

// werGlobal soma os erros de todas as linhas sobre o total de palavras de referencia.
func werGlobal(linhas []linhaVoz) float64 {
	var erros, total int
	for _, l := range linhas {
		erros += l.erros
		total += l.total
	}
	return razaoWER(erros, total)
}

func textoDe(palavras []PalavraTranscrita) string {
	partes := make([]string, len(palavras))
	for i, p := range palavras {
		partes[i] = p.Texto
	}
	return strings.Join(partes, " ")
}

func montarTimings(provedor, voz string, linhas []linhaVoz) Timings {
	tm := Timings{SchemaVersion: 1, Provedor: provedor, Voz: voz, WERGlobal: werGlobal(linhas)}
	inicio := 0.0
	for i, lv := range linhas {
		f, dados, _ := lerWAV(lv.wav)
		dur := duracaoWAV(f, len(dados))
		tm.Linhas = append(tm.Linhas, LinhaTempo{
			IndiceLinha: i, InicioS: inicio, FimS: inicio + dur, WER: lv.wer, Regravacoes: lv.regrav,
			Palavras: ajustarPalavras(lv.palavras, inicio, dur),
		})
		inicio += dur
	}
	return tm
}

// ajustarPalavras desloca para a linha do video e garante ordem, sem sobreposicao e dentro da janela.
func ajustarPalavras(brutas []PalavraTranscrita, inicio, dur float64) []PalavraTempo {
	fim, cursor := inicio+dur, inicio
	var saida []PalavraTempo
	for _, p := range brutas {
		texto := strings.TrimSpace(p.Texto)
		if texto == "" {
			continue
		}
		ini := math.Min(math.Max(inicio+p.InicioS, cursor), fim)
		f := math.Min(math.Max(inicio+p.FimS, ini), fim)
		saida = append(saida, PalavraTempo{Texto: texto, InicioS: ini, FimS: f})
		cursor = f
	}
	return saida
}

func gravarEtapaVoz(projDir, provedor, voz string, linhas []linhaVoz, tm Timings) error {
	trechos := make([][]byte, len(linhas))
	for i, lv := range linhas {
		trechos[i] = lv.wav
	}
	if _, err := gravarVoz(projDir, ResultadoVoz{Provedor: provedor, Voz: voz}, trechos); err != nil {
		return err
	}
	b, err := json.MarshalIndent(tm, "", "  ")
	if err != nil {
		return err
	}
	return gravarSeguro(projDir, "data/timings.json", append(b, '\n'))
}

func palavrasDe(s string) []string {
	return strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

// errosPalavras devolve a distancia de palavras e o total de palavras da referencia.
func errosPalavras(ref, hip string) (erros, total int) {
	r, h := palavrasDe(ref), palavrasDe(hip)
	return distanciaPalavras(r, h), len(r)
}

// razaoWER limita a 1; referencia vazia so passa se nao houver erro.
func razaoWER(erros, total int) float64 {
	if total == 0 {
		if erros == 0 {
			return 0
		}
		return 1
	}
	return math.Min(1, float64(erros)/float64(total))
}

// CalcularWER mede a distancia de palavras entre referencia e hipotese, limitada a 1.
// Acentos contam: "publica" e "pública" sao palavras diferentes.
func CalcularWER(ref, hip string) float64 {
	return razaoWER(errosPalavras(ref, hip))
}

func distanciaPalavras(a, b []string) int {
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		atual := make([]int, len(b)+1)
		atual[0] = i
		for j := 1; j <= len(b); j++ {
			custo := 1
			if a[i-1] == b[j-1] {
				custo = 0
			}
			atual[j] = min(prev[j]+1, atual[j-1]+1, prev[j-1]+custo)
		}
		prev = atual
	}
	return prev[len(b)]
}

// Capturador roda um programa e devolve sua saida padrao; injetavel nos testes.
type Capturador interface {
	Capturar(ctx context.Context, nome string, args ...string) ([]byte, error)
}

func (execSistema) Capturar(ctx context.Context, nome string, args ...string) ([]byte, error) {
	saida, err := exec.CommandContext(ctx, nome, args...).Output()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return nil, fmt.Errorf("%s: %w: %s", nome, err, strings.TrimSpace(string(ee.Stderr)))
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", nome, err)
	}
	return saida, nil
}

// scriptFasterWhisper imprime JSON com ASCII escapado para sobreviver ao console do Windows.
const scriptFasterWhisper = `import sys, json
from faster_whisper import WhisperModel
m = WhisperModel(sys.argv[3], device="cpu", compute_type="int8")
segs, _ = m.transcribe(sys.argv[1], language=sys.argv[2], word_timestamps=True)
print(json.dumps([{"w": w.word, "s": w.start, "e": w.end} for s in segs for w in (s.words or [])]))`

// FasterWhisper e o Transcritor real, executado em CPU pelo modulo Python faster_whisper.
type FasterWhisper struct {
	Exec   Capturador
	Modelo string
}

func (f FasterWhisper) Transcrever(ctx context.Context, wav []byte, idioma string) ([]PalavraTranscrita, error) {
	dir, err := os.MkdirTemp("", "whisper-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	arq := filepath.Join(dir, "trecho.wav")
	if err := os.WriteFile(arq, wav, 0o600); err != nil {
		return nil, err
	}
	modelo := f.Modelo
	if modelo == "" {
		modelo = "small"
	}
	saida, err := f.Exec.Capturar(ctx, "python", "-c", scriptFasterWhisper, arq, idiomaBase(idioma), modelo)
	if err != nil {
		return nil, err
	}
	return lerPalavrasWhisper(saida)
}

func idiomaBase(idioma string) string {
	base, _, _ := strings.Cut(strings.TrimSpace(idioma), "-")
	return strings.ToLower(base)
}

func lerPalavrasWhisper(saida []byte) ([]PalavraTranscrita, error) {
	var brutas []struct {
		W string  `json:"w"`
		S float64 `json:"s"`
		E float64 `json:"e"`
	}
	if err := json.Unmarshal(saida, &brutas); err != nil {
		return nil, fmt.Errorf("diorama-render: saida do faster-whisper invalida: %w", err)
	}
	palavras := make([]PalavraTranscrita, 0, len(brutas))
	for _, b := range brutas {
		if b.S < 0 || b.E < b.S {
			return nil, fmt.Errorf("diorama-render: tempo invalido do faster-whisper para %q", b.W)
		}
		palavras = append(palavras, PalavraTranscrita{Texto: strings.TrimSpace(b.W), InicioS: b.S, FimS: b.E})
	}
	return palavras, nil
}
