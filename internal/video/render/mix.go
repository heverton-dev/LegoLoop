package render

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/videomotion"
)

// Parametros da Etapa 7. O teto de pico vem de videomotion.TruePeakTetoDBTP (-1,5 dBTP).
const (
	TaxaMix            = 48000
	AlvoLUFS           = -14.0
	ToleranciaLUFS     = 0.5
	MargemPicoDB       = 0.5
	MaxTentativasMix   = 3
	ToleranciaDuracaoS = 0.034
	BGMGanhoPadraoDB   = 3.0
	BGMGanhoMinDB      = -40.0
	BGMGanhoMaxDB      = 12.0
	ganhoSFXdB         = -12
)

// Medicao e o resultado de medir o arquivo final, independente da previsao do loudnorm.
type Medicao struct {
	LUFS         float64
	TruePeakDBTP float64
}

// MedidasLoudnorm guarda, como o ffmpeg imprimiu, os valores medidos na 1a passada.
type MedidasLoudnorm struct {
	InputI, InputTP, InputLRA, InputThresh, TargetOffset string
}

// Medidor roda um programa e devolve stdout e stderr juntos (o ffmpeg imprime as medidas no stderr).
type Medidor interface {
	Saida(ctx context.Context, nome string, args ...string) (string, error)
}

func (execSistema) Saida(ctx context.Context, nome string, args ...string) (string, error) {
	saida, err := exec.CommandContext(ctx, nome, args...).CombinedOutput()
	if err != nil {
		return string(saida), fmt.Errorf("%s: %w", nome, err)
	}
	return string(saida), nil
}

// MixDeps isola o ffmpeg para testes hermeticos.
type MixDeps struct {
	Exec Executor
	Med  Medidor
}

func normalizarMixDeps(d MixDeps) MixDeps {
	if d.Exec == nil {
		d.Exec = execSistema{}
	}
	if d.Med == nil {
		d.Med = execSistema{}
	}
	return d
}

type TentativaMix struct {
	Tentativa    int     `json:"tentativa"`
	AlvoLUFS     float64 `json:"alvo_lufs"`
	AlvoPicoDBTP float64 `json:"alvo_pico_dbtp"`
	LUFS         float64 `json:"lufs"`
	PicoDBTP     float64 `json:"pico_dbtp"`
	Aprovada     bool    `json:"aprovada"`
}

// RelatorioMix e gravado em data/mix.json para auditoria.
type RelatorioMix struct {
	SchemaVersion int            `json:"schema_version"`
	DuracaoS      float64        `json:"duracao_s"`
	Eventos       int            `json:"eventos_sfx"`
	LUFS          float64        `json:"lufs"`
	TruePeakDBTP  float64        `json:"true_peak_dbtp"`
	Tentativas    []TentativaMix `json:"tentativas"`
}

var (
	reLUFS = regexp.MustCompile(`\sI:\s+(-?\d+(?:\.\d+)?) LUFS`)
	rePico = regexp.MustCompile(`Peak:\s+(-?\d+(?:\.\d+)?) dBFS`)
)

func numeroDe(re *regexp.Regexp, texto, nome string) (float64, error) {
	m := re.FindStringSubmatch(texto)
	if m == nil {
		return 0, fmt.Errorf("video-render: %s nao encontrado na saida do ebur128", nome)
	}
	return strconv.ParseFloat(m[1], 64)
}

// parseEBUR128 le o ultimo Summary: LUFS integrado e True Peak.
func parseEBUR128(saida string) (Medicao, error) {
	i := strings.LastIndex(saida, "Summary:")
	if i < 0 {
		return Medicao{}, fmt.Errorf("video-render: saida do ebur128 sem Summary")
	}
	resumo := saida[i:]
	lufs, err := numeroDe(reLUFS, resumo, "LUFS integrado")
	if err != nil {
		return Medicao{}, err
	}
	pico, err := numeroDe(rePico, resumo, "True Peak")
	return Medicao{LUFS: lufs, TruePeakDBTP: pico}, err
}

// parseLoudnormJSON le o bloco JSON impresso por loudnorm com print_format=json.
func parseLoudnormJSON(saida string) (MedidasLoudnorm, error) {
	i, j := strings.LastIndex(saida, "{"), strings.LastIndex(saida, "}")
	if i < 0 || j < i {
		return MedidasLoudnorm{}, fmt.Errorf("video-render: saida do loudnorm sem JSON")
	}
	var b struct {
		I      string `json:"input_i"`
		TP     string `json:"input_tp"`
		LRA    string `json:"input_lra"`
		Thresh string `json:"input_thresh"`
		Offset string `json:"target_offset"`
	}
	if err := json.Unmarshal([]byte(saida[i:j+1]), &b); err != nil {
		return MedidasLoudnorm{}, fmt.Errorf("video-render: JSON do loudnorm invalido: %w", err)
	}
	m := MedidasLoudnorm{InputI: b.I, InputTP: b.TP, InputLRA: b.LRA, InputThresh: b.Thresh, TargetOffset: b.Offset}
	for _, v := range []string{m.InputI, m.InputTP, m.InputLRA, m.InputThresh, m.TargetOffset} {
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			return MedidasLoudnorm{}, fmt.Errorf("video-render: medida %q do loudnorm invalida", v)
		}
	}
	return m, nil
}

// parametrosTentativa mira 0,5 dB abaixo do teto na 1a tentativa e ganha 0,5 dB de margem por
// tentativa; a partir da 2a compensa o LUFS que saiu diferente do alvo.
func parametrosTentativa(k int, ultimoLUFS, ultimoAlvoI float64) (alvoI, alvoTP float64) {
	alvoTP = videomotion.TruePeakTetoDBTP - MargemPicoDB - 0.5*float64(k)
	if k == 0 {
		return AlvoLUFS, alvoTP
	}
	return ultimoAlvoI + (AlvoLUFS - ultimoLUFS), alvoTP
}

// amostrasDe converte a duracao da timeline em numero de amostras a TaxaMix.
func amostrasDe(dur float64) int { return arredondar(dur * TaxaMix) }

// filtroNorm e a 2a passada do loudnorm. O loudnorm desloca os timestamps de saida em 28 ms, entao
// -t cortaria amostras validas; apad e atrim=end_sample fixam a duracao exata contando amostras.
func filtroNorm(m MedidasLoudnorm, alvoI, alvoTP, dur float64) string {
	return fmt.Sprintf("loudnorm=I=%.2f:TP=%.1f:LRA=11:measured_I=%s:measured_TP=%s:measured_LRA=%s:measured_thresh=%s:offset=%s:linear=true,aresample=%d,apad=pad_dur=1,atrim=end_sample=%d",
		alvoI, alvoTP, m.InputI, m.InputTP, m.InputLRA, m.InputThresh, m.TargetOffset, TaxaMix, amostrasDe(dur))
}

// EntradaMix descreve a mixagem inicial; ArquivosSFX[i] e o audio de Eventos[i].
type EntradaMix struct {
	Duracao     float64
	Voz         string
	CanaisVoz   uint16
	BGM         string
	GanhoBGMdB  float64
	Eventos     []EventoSFX
	ArquivosSFX []string
}

// formatoVoz poe a voz em estereo pleno; o -ac 2 do ffmpeg atenuaria 3 dB num arquivo mono.
func formatoVoz(canais uint16) string {
	if canais == 1 {
		return "pan=stereo|c0=c0|c1=c0"
	}
	return "aformat=channel_layouts=stereo"
}

func filtroMix(e EntradaMix) string {
	partes := []string{
		fmt.Sprintf("[0:a]aresample=%d,%s,asplit=2[voz][vozsc]", TaxaMix, formatoVoz(e.CanaisVoz)),
		fmt.Sprintf("[1:a]aresample=%d,aformat=channel_layouts=stereo,volume=%.1fdB,afade=t=in:d=1,afade=t=out:st=%.3f:d=1.5[bgm0]", TaxaMix, e.GanhoBGMdB, max(e.Duracao-1.5, 0)),
		"[bgm0][vozsc]sidechaincompress=threshold=0.04:ratio=6:attack=20:release=600:makeup=1[bgm]",
	}
	rotulos := "[voz][bgm]"
	for i, ev := range e.Eventos {
		ms := arredondar(ev.TempoS * 1000)
		partes = append(partes, fmt.Sprintf("[%d:a]aresample=%d,aformat=channel_layouts=stereo,volume=%ddB,adelay=%d|%d[s%d]", i+2, TaxaMix, ganhoSFXdB, ms, ms, i))
		rotulos += fmt.Sprintf("[s%d]", i)
	}
	partes = append(partes, fmt.Sprintf("%samix=inputs=%d:normalize=0:duration=longest,atrim=end_sample=%d[mix]", rotulos, 2+len(e.Eventos), amostrasDe(e.Duracao)))
	return strings.Join(partes, ";")
}

// argsMix monta a mixagem inicial: voz, BGM em loop com ducking e um SFX por evento. O corte exato
// e por amostra (atrim); -t, 1 s maior, so impede o loop infinito da BGM.
func argsMix(e EntradaMix, saida string) []string {
	args := []string{"-y", "-hide_banner", "-nostats", "-i", e.Voz, "-stream_loop", "-1", "-i", e.BGM}
	for _, a := range e.ArquivosSFX {
		args = append(args, "-i", a)
	}
	return append(args, "-filter_complex", filtroMix(e), "-map", "[mix]", "-t", fmt.Sprintf("%.6f", e.Duracao+1),
		"-c:a", "pcm_f32le", "-ar", strconv.Itoa(TaxaMix), saida)
}
