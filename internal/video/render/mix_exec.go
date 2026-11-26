package render

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/videomotion"
)

// preparoMix reune o que a mixagem precisa depois de validar todas as entradas.
type preparoMix struct {
	tl      Timeline
	entrada EntradaMix
}

func lerTimelineSalva(projDir string) (Timeline, error) {
	var tl Timeline
	if err := lerJSON(projDir, "timeline.json", &tl, true); err != nil {
		return tl, entradaInvalida("timeline.json: %v", err)
	}
	if err := tl.Validar(); err != nil {
		return tl, fmt.Errorf("%w: %w", ErrEntradaTimeline, err)
	}
	return tl, nil
}

// conferirVoz exige WAV de 1 ou 2 canais com a mesma duracao da timeline.
func conferirVoz(projDir string, dur float64) (uint16, error) {
	b, err := lerBytes(projDir, "audio/voice.wav")
	if err != nil {
		return 0, bloqueio("audio/voice.wav", err)
	}
	f, dados, err := lerWAV(b)
	if err != nil {
		return 0, err
	}
	if f.Canais != 1 && f.Canais != 2 {
		return 0, audioInvalido(fmt.Sprintf("voice.wav com %d canais; aceito 1 ou 2", f.Canais))
	}
	if d := duracaoWAV(f, len(dados)); math.Abs(d-dur) > ToleranciaDuracaoS {
		return 0, fmt.Errorf("%w: voice.wav dura %.3f s e a timeline %.3f s", ErrRenderInvalido, d, dur)
	}
	return f.Canais, nil
}

func prepararMix(ctx context.Context, projDir string, c ContratosRender, deps MixDeps) (preparoMix, error) {
	var p preparoMix
	tl, err := lerTimelineSalva(projDir)
	if err != nil {
		return p, err
	}
	canais, err := conferirVoz(projDir, tl.DuracaoS)
	if err != nil {
		return p, err
	}
	if err := GarantirAtivosSonoros(ctx, projDir, c.Config, deps.Exec); err != nil {
		return p, err
	}
	cat, err := CarregarCatalogo(projDir, c.Config)
	if err != nil {
		return p, err
	}
	if p.tl, err = GerarEventosSFX(projDir, tl, cat); err != nil {
		return p, err
	}
	p.entrada, err = montarEntradaMix(projDir, c.Config, p.tl, canais)
	return p, err
}

func montarEntradaMix(projDir string, cfg ConfigRender, tl Timeline, canais uint16) (EntradaMix, error) {
	e := EntradaMix{Duracao: tl.DuracaoS, CanaisVoz: canais, GanhoBGMdB: BGMGanhoPadraoDB}
	if cfg.BGMGanhoDB != nil {
		e.GanhoBGMdB = *cfg.BGMGanhoDB
	}
	bgm := cfg.BGM
	if bgm == "" {
		bgm = "audio/bgm.wav"
	}
	var err error
	if e.Voz, err = resolverSeguro(projDir, "audio/voice.wav"); err != nil {
		return e, err
	}
	if e.BGM, err = resolverSeguro(projDir, bgm); err != nil {
		return e, bloqueio("bgm", err)
	}
	for _, b := range tl.Batidas {
		for _, ev := range b.SFX {
			abs, err := resolverSeguro(projDir, ev.Arquivo)
			if err != nil {
				return e, err
			}
			e.Eventos = append(e.Eventos, ev)
			e.ArquivosSFX = append(e.ArquivosSFX, abs)
		}
	}
	return e, nil
}

func medirLoudnorm(ctx context.Context, med Medidor, wav string) (MedidasLoudnorm, error) {
	saida, err := med.Saida(ctx, "ffmpeg", "-hide_banner", "-nostats", "-i", wav, "-af", "loudnorm=I=-14:TP=-2.0:LRA=11:print_format=json", "-f", "null", "-")
	if err != nil {
		return MedidasLoudnorm{}, err
	}
	return parseLoudnormJSON(saida)
}

// medirEBUR128 mede o arquivo final com ebur128 (true peak ligado), independente do loudnorm.
func medirEBUR128(ctx context.Context, med Medidor, wav string) (Medicao, error) {
	saida, err := med.Saida(ctx, "ffmpeg", "-hide_banner", "-nostats", "-i", wav, "-af", "ebur128=peak=true", "-f", "null", "-")
	if err != nil {
		return Medicao{}, err
	}
	return parseEBUR128(saida)
}

func conforme(m Medicao) error {
	return videomotion.ValidateLoudnessConformity(videomotion.LoudnessMetrics{OutputIntegratedLUFS: m.LUFS, OutputTruePeakDBFS: m.TruePeakDBTP}, AlvoLUFS, ToleranciaLUFS)
}

// tentarNormalizacao roda a 2a passada da tentativa k e mede o arquivo resultante.
func tentarNormalizacao(ctx context.Context, deps MixDeps, tmp, pre string, med MedidasLoudnorm, dur float64, k int, ultLUFS, ultAlvo float64) (TentativaMix, string, error) {
	alvoI, alvoTP := parametrosTentativa(k, ultLUFS, ultAlvo)
	saida := filepath.Join(tmp, fmt.Sprintf("norm%d.wav", k))
	args := []string{"-y", "-i", pre, "-af", filtroNorm(med, alvoI, alvoTP, dur), "-c:a", "pcm_s16le", "-ar", fmt.Sprint(TaxaMix), saida}
	t := TentativaMix{Tentativa: k + 1, AlvoLUFS: alvoI, AlvoPicoDBTP: alvoTP}
	if err := deps.Exec.Executar(ctx, "ffmpeg", args...); err != nil {
		return t, "", err
	}
	m, err := medirEBUR128(ctx, deps.Med, saida)
	if err != nil {
		return t, "", err
	}
	t.LUFS, t.PicoDBTP = m.LUFS, m.TruePeakDBTP
	t.Aprovada = conforme(m) == nil
	return t, saida, nil
}

// normalizarVerificado so devolve audio que passou na medicao independente do arquivo final.
func normalizarVerificado(ctx context.Context, deps MixDeps, p preparoMix) (RelatorioMix, []byte, error) {
	rel := RelatorioMix{SchemaVersion: 1, DuracaoS: p.entrada.Duracao, Eventos: len(p.entrada.Eventos)}
	tmp, err := os.MkdirTemp("", "mix-*")
	if err != nil {
		return rel, nil, err
	}
	defer os.RemoveAll(tmp)
	pre := filepath.Join(tmp, "premix.wav")
	if err := deps.Exec.Executar(ctx, "ffmpeg", argsMix(p.entrada, pre)...); err != nil {
		return rel, nil, err
	}
	med, err := medirLoudnorm(ctx, deps.Med, pre)
	if err != nil {
		return rel, nil, err
	}
	var ultimo TentativaMix
	for k := 0; k < MaxTentativasMix; k++ {
		t, arq, err := tentarNormalizacao(ctx, deps, tmp, pre, med, p.entrada.Duracao, k, ultimo.LUFS, ultimo.AlvoLUFS)
		if err != nil {
			return rel, nil, err
		}
		rel.Tentativas = append(rel.Tentativas, t)
		if t.Aprovada {
			rel.LUFS, rel.TruePeakDBTP = t.LUFS, t.PicoDBTP
			b, err := finalizarSoundtrack(arq, p.entrada.Duracao)
			return rel, b, err
		}
		ultimo = t
	}
	return rel, nil, fmt.Errorf("%w: %d tentativas sem conformidade; ultimo: %v", ErrLoudnessReprovado, MaxTentativasMix,
		conforme(Medicao{LUFS: ultimo.LUFS, TruePeakDBTP: ultimo.PicoDBTP}))
}

func finalizarSoundtrack(arquivo string, dur float64) ([]byte, error) {
	b, err := os.ReadFile(arquivo)
	if err != nil {
		return nil, err
	}
	f, dados, err := lerWAV(b)
	if err != nil {
		return nil, err
	}
	if d := duracaoWAV(f, len(dados)); math.Abs(d-dur) > ToleranciaDuracaoS {
		return nil, fmt.Errorf("%w: soundtrack dura %.3f s e a timeline %.3f s", ErrRenderInvalido, d, dur)
	}
	return b, nil
}

func gravarSaidasMix(projDir string, tl Timeline, soundtrack []byte, rel RelatorioMix) error {
	mixJSON, err := json.MarshalIndent(rel, "", "  ")
	if err != nil {
		return err
	}
	tlJSON, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}
	return gravarEmLote(projDir, []saidaArquivo{
		{"audio/soundtrack.wav", soundtrack},
		{"data/mix.json", append(mixJSON, '\n')},
		{"timeline.json", append(tlJSON, '\n')},
	})
}

// ExecutarMix e a Etapa 7: valida entradas, gera eventos de SFX, mixa com ducking, normaliza com
// margem de pico e so grava soundtrack.wav, timeline.json e data/mix.json se a medicao independente
// do arquivo final aprovar (-14 +- 0,5 LUFS e pico <= -1,5 dBTP).
func ExecutarMix(ctx context.Context, projDir string, deps MixDeps) (RelatorioMix, error) {
	if err := ctx.Err(); err != nil {
		return RelatorioMix{}, err
	}
	deps = normalizarMixDeps(deps)
	c, err := CarregarContratos(projDir)
	if err != nil {
		return RelatorioMix{}, err
	}
	p, err := prepararMix(ctx, projDir, c, deps)
	if err != nil {
		return RelatorioMix{}, err
	}
	rel, soundtrack, err := normalizarVerificado(ctx, deps, p)
	if err != nil {
		return rel, err
	}
	return rel, gravarSaidasMix(projDir, p.tl, soundtrack, rel)
}
