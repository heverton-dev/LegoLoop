package motor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/video/render"
)

// ErrCorteInvalido sinaliza intervalo vazio ou duracao fora de 15..59s (Etapa 8b).
var ErrCorteInvalido = errors.New("motor-cortes: corte vertical invalido")

// ItemCorte descreve um Short vertical 9:16 extraido da timeline principal.
type ItemCorte struct {
	ID            string   `json:"id"`
	Titulo        string   `json:"titulo,omitempty"`
	Intervalos    [][2]int `json:"intervalos"`
	CartaoFinal   string   `json:"cartao_final,omitempty"`
	DuracaoCartao float64  `json:"duracao_cartao,omitempty"`
}

// EspecificacaoCortes agrupa os cortes curtos de um projeto.
type EspecificacaoCortes struct {
	CartaoFinal   string      `json:"cartao_final,omitempty"`
	DuracaoCartao float64     `json:"duracao_cartao,omitempty"`
	Cortes        []ItemCorte `json:"cortes"`
}

// ResumoCorte descreve o corte preparado em disco pronto para renderizar em 9:16.
type ResumoCorte struct {
	ID       string  `json:"id"`
	Cenas    int     `json:"cenas"`
	DuracaoS float64 `json:"duracao_s"`
	FalaS    float64 `json:"fala_s"`
	Pasta    string  `json:"pasta"`
}

// SugerirCortes divide uma timeline longa em janelas contiguas de 15 a 55 s (+ cartao CTA de 3.2 s).
func SugerirCortes(tl Timeline) []ItemCorte {
	var itens []ItemCorte
	iniLinha, fimLinha, soma := -1, -1, 0.0
	seq := 1
	for _, s := range tl.Shots {
		if s.Lines[0] < 0 {
			continue
		}
		if iniLinha < 0 {
			iniLinha = s.Lines[0]
		}
		if soma+s.Duration > 52.0 && soma >= 15.0 {
			itens = append(itens, novoItemCorteSugerido(seq, iniLinha, fimLinha))
			seq++
			iniLinha, soma = s.Lines[0], 0
		}
		fimLinha = s.Lines[1]
		soma += s.Duration
	}
	if iniLinha >= 0 && soma >= 15.0 && soma <= 55.5 {
		itens = append(itens, novoItemCorteSugerido(seq, iniLinha, fimLinha))
	}
	return itens
}

func novoItemCorteSugerido(seq, a, b int) ItemCorte {
	return ItemCorte{
		ID:            fmt.Sprintf("short-%02d", seq),
		Intervalos:    [][2]int{{a, b}},
		CartaoFinal:   "Caso completo no canal",
		DuracaoCartao: 3.2,
	}
}

// PrepararCortes monta cortes/<id>/timeline.json (9:16), data/lines.json e audio recortado para cada corte.
func PrepararCortes(ctx context.Context, projDir string, espec EspecificacaoCortes, cmd ExecutorComando) ([]ResumoCorte, error) {
	if len(espec.Cortes) == 0 {
		return nil, fmt.Errorf("%w: lista de cortes vazia", ErrCorteInvalido)
	}
	var tl Timeline
	if err := lerJSONProjeto(projDir, "timeline.json", &tl); err != nil {
		return nil, err
	}
	var tm render.Timings
	if err := lerJSONProjeto(projDir, "data/timings.json", &tm); err != nil {
		return nil, err
	}
	starts := inciosDosShots(tl.Shots)
	resumos := make([]ResumoCorte, 0, len(espec.Cortes))
	for _, c := range espec.Cortes {
		r, err := prepararUmCorte(ctx, projDir, tl, tm, starts, espec, c, cmd)
		if err != nil {
			return nil, err
		}
		resumos = append(resumos, r)
	}
	return resumos, nil
}

func inciosDosShots(shots []Corte) []float64 {
	out := make([]float64, len(shots))
	t := 0.0
	for i, s := range shots {
		out[i] = t
		t += s.Duration
	}
	return out
}

func prepararUmCorte(ctx context.Context, projDir string, tl Timeline, tm render.Timings, starts []float64, espec EspecificacaoCortes, c ItemCorte, cmd ExecutorComando) (ResumoCorte, error) {
	novos, janelas, err := selecionarShotsEJanelas(projDir, tl.Shots, starts, c)
	if err != nil {
		return ResumoCorte{}, err
	}
	ctaTxt, ctaDur := resolverCartaoFinal(espec, c)
	novos = append(novos, Corte{ID: "cta", Type: "card", Lines: [2]int{-1, -1}, Text: ctaTxt, Duration: ctaDur, CPS: 10, Chars: []Boneco{}, Layers: []map[string]any{}})
	total := somarDuracoesCortes(novos)
	if total < 15.0 || total > 59.0 {
		return ResumoCorte{}, fmt.Errorf("%w: %s dura %.1fs (exigido 15..59s)", ErrCorteInvalido, c.ID, total)
	}
	pasta := filepath.Join(projDir, "cortes", c.ID)
	if err := persistirArtefatosDoCorte(ctx, projDir, pasta, tl, tm, novos, janelas, total, cmd); err != nil {
		return ResumoCorte{}, err
	}
	return ResumoCorte{ID: c.ID, Cenas: len(novos) - 1, DuracaoS: total, FalaS: total - ctaDur, Pasta: pasta}, nil
}

func selecionarShotsEJanelas(projDir string, shots []Corte, starts []float64, c ItemCorte) ([]Corte, [][2]float64, error) {
	if strings.TrimSpace(c.ID) == "" || len(c.Intervalos) == 0 {
		return nil, nil, fmt.Errorf("%w: id ou intervalos vazios", ErrCorteInvalido)
	}
	var novos []Corte
	var janelas [][2]float64
	for _, iv := range c.Intervalos {
		iniIdx, fimIdx := -1, -1
		for i, s := range shots {
			if s.Lines[0] >= iv[0] && s.Lines[1] <= iv[1] && s.Lines[0] >= 0 {
				if iniIdx < 0 {
					iniIdx = i
				}
				fimIdx = i
				cp := s
				if cp.BG != "" && !filepath.IsAbs(cp.BG) {
					cp.BG = filepath.ToSlash(filepath.Join(projDir, cp.BG))
				}
				novos = append(novos, cp)
			}
		}
		if iniIdx < 0 {
			return nil, nil, fmt.Errorf("%w: %s sem cortes nas linhas %d-%d", ErrCorteInvalido, c.ID, iv[0], iv[1])
		}
		janelas = append(janelas, [2]float64{starts[iniIdx], starts[fimIdx] + shots[fimIdx].Duration})
	}
	return novos, janelas, nil
}

func resolverCartaoFinal(espec EspecificacaoCortes, c ItemCorte) (string, float64) {
	txt := c.CartaoFinal
	if txt == "" {
		txt = espec.CartaoFinal
	}
	if txt == "" {
		txt = "Caso completo no canal"
	}
	dur := c.DuracaoCartao
	if dur <= 0 {
		dur = espec.DuracaoCartao
	}
	if dur <= 0 {
		dur = 3.2
	}
	return txt, dur
}

func somarDuracoesCortes(cs []Corte) float64 {
	s := 0.0
	for _, c := range cs {
		s += c.Duration
	}
	return s
}

func persistirArtefatosDoCorte(ctx context.Context, projDir, pasta string, tl Timeline, tm render.Timings, novos []Corte, janelas [][2]float64, total float64, cmd ExecutorComando) error {
	for _, sub := range []string{"audio", "data"} {
		if err := os.MkdirAll(filepath.Join(pasta, sub), 0o755); err != nil {
			return err
		}
	}
	if err := recortarAudioJanelas(ctx, filepath.Join(projDir, "audio", "voice.wav"), janelas, total, filepath.Join(pasta, "audio", "voice.wav"), cmd); err != nil {
		return err
	}
	stOrig := filepath.Join(projDir, "audio", "soundtrack.wav")
	audioRef := "audio/voice.wav"
	if _, err := os.Stat(stOrig); err == nil {
		if err := recortarAudioJanelas(ctx, stOrig, janelas, total, filepath.Join(pasta, "audio", "soundtrack.wav"), cmd); err != nil {
			return err
		}
		audioRef = "audio/soundtrack.wav"
	}
	if err := gravarLinesCorte(pasta, tm, janelas, total); err != nil {
		return err
	}
	saidaTL := map[string]any{"fps": tl.FPS, "formato": "9:16", "audio": audioRef, "music": tl.Music, "cast": tl.Cast, "shots": novos}
	b, _ := json.MarshalIndent(saidaTL, "", "  ")
	return os.WriteFile(filepath.Join(pasta, "timeline.json"), b, 0o644)
}

func recortarAudioJanelas(ctx context.Context, src string, janelas [][2]float64, total float64, saida string, cmd ExecutorComando) error {
	var partes, ins []string
	for i, j := range janelas {
		d := j[1] - j[0]
		partes = append(partes, fmt.Sprintf("[0:a]atrim=%.3f:%.3f,asetpts=PTS-STARTPTS,afade=t=in:d=0.04,afade=t=out:st=%.3f:d=0.06[p%d]", j[0], j[1], max(0.0, d-0.06), i))
		ins = append(ins, fmt.Sprintf("[p%d]", i))
	}
	fc := fmt.Sprintf("%s;%sconcat=n=%d:v=0:a=1,apad=whole_dur=%.3f[o]", strings.Join(partes, ";"), strings.Join(ins, ""), len(janelas), total)
	_, err := cmd.Rodar(ctx, "ffmpeg", "-y", "-v", "error", "-i", src, "-filter_complex", fc, "-map", "[o]", "-ar", "48000", "-ac", "2", saida)
	return err
}

func gravarLinesCorte(pasta string, tm render.Timings, janelas [][2]float64, total float64) error {
	var linhas []map[string]any
	off := 0.0
	for _, j := range janelas {
		for _, l := range tm.Linhas {
			if l.FimS > j[0]+0.05 && l.InicioS < j[1]-0.05 {
				linhas = append(linhas, map[string]any{
					"i": len(linhas), "text": textoDasPalavras(l.Palavras),
					"start": max(0.0, off+l.InicioS-j[0]), "end": off + min(l.FimS, j[1]) - j[0],
				})
			}
		}
		off += j[1] - j[0]
	}
	b, _ := json.MarshalIndent(map[string]any{"duration": total, "lines": linhas}, "", "  ")
	return os.WriteFile(filepath.Join(pasta, "data", "lines.json"), b, 0o644)
}

func textoDasPalavras(ps []render.PalavraTempo) string {
	ws := make([]string, len(ps))
	for i, p := range ps {
		ws[i] = p.Texto
	}
	return strings.Join(ws, " ")
}

// RenderizarCortes renderiza cada corte preparado em 9:16 com legendas queimadas e marca.
func RenderizarCortes(ctx context.Context, projDir string, resumos []ResumoCorte, renderPy string, cmd ExecutorComando) ([]string, error) {
	saidaDir, err := filepath.Abs(filepath.Join(projDir, "renders", "shorts"))
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(saidaDir, 0o755); err != nil {
		return nil, err
	}
	marcaPath := filepath.Join(projDir, "marca.json")
	var gerados []string
	for _, r := range resumos {
		out := filepath.Join(saidaDir, r.ID+".mp4")
		args := []string{renderPy, r.Pasta, "--formato", "9:16", "--legendas", "sim", "--sem-checagem", "--out", out}
		if _, err := os.Stat(marcaPath); err == nil {
			args = append(args, "--marca", marcaPath)
		}
		if _, err := cmd.Rodar(ctx, "python", args...); err != nil {
			return nil, err
		}
		gerados = append(gerados, out)
	}
	return gerados, nil
}
