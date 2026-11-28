package motor

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

// ExecutorComando abstrai chamadas externas a python, ffmpeg e ffprobe.
type ExecutorComando interface {
	Rodar(ctx context.Context, nome string, args ...string) ([]byte, error)
}

// ExecutorSistemaComando executa programas reais sem shell.
type ExecutorSistemaComando struct{}

// Rodar invoca o programa e retorna saida combinada.
func (ExecutorSistemaComando) Rodar(ctx context.Context, nome string, args ...string) ([]byte, error) {
	out, err := exec.CommandContext(ctx, nome, args...).CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("%s: %w: %s", nome, err, strings.TrimSpace(string(out)))
	}
	return out, nil
}

// RelatorioMarca descreve o resultado da pos-producao de marca (vinhetas + logo).
type RelatorioMarca struct {
	Saida           string `json:"saida"`
	Formato         string `json:"formato"`
	Largura         int    `json:"largura"`
	Altura          int    `json:"altura"`
	LogoPath        string `json:"logo_path"`
	AberturaAtiva   bool   `json:"abertura_ativa"`
	FechamentoAtivo bool   `json:"fechamento_ativo"`
}

// AplicarMarca aplica vinheta de abertura, logo sobre o video principal e vinheta de fechamento (equivalente a marca.py).
func AplicarMarca(ctx context.Context, projDir, videoEntrada, videoSaida, renderPy string, cmd ExecutorComando) (RelatorioMarca, error) {
	m, err := lerMarcaProjeto(projDir)
	if err != nil {
		return RelatorioMarca{}, err
	}
	w, h, fmtVideo, err := medirResolucaoVideo(ctx, videoEntrada, cmd)
	if err != nil {
		return RelatorioMarca{}, err
	}
	logoPath, err := garantirLogoMarca(projDir, m)
	if err != nil {
		return RelatorioMarca{}, err
	}
	vd := filepath.Join(projDir, ".vinhetas")
	if err := os.MkdirAll(vd, 0o755); err != nil {
		return RelatorioMarca{}, err
	}
	partes, abAtiva, feAtiva, err := renderizarPartesMarca(ctx, projDir, vd, videoEntrada, logoPath, renderPy, fmtVideo, w, h, m, cmd)
	if err != nil {
		return RelatorioMarca{}, err
	}
	if err := concatenarPartesMarca(ctx, partes, videoSaida, cmd); err != nil {
		return RelatorioMarca{}, err
	}
	return RelatorioMarca{Saida: videoSaida, Formato: fmtVideo, Largura: w, Altura: h, LogoPath: logoPath, AberturaAtiva: abAtiva, FechamentoAtivo: feAtiva}, nil
}

func lerMarcaProjeto(projDir string) (roteiro.Marca, error) {
	var m roteiro.Marca
	dados, err := os.ReadFile(filepath.Join(projDir, "marca.json"))
	if err != nil {
		return m, fmt.Errorf("marca.json: %w", err)
	}
	if err := json.Unmarshal(dados, &m); err != nil {
		return m, fmt.Errorf("marca.json invalido: %w", err)
	}
	return m, nil
}

func medirResolucaoVideo(ctx context.Context, video string, cmd ExecutorComando) (int, int, string, error) {
	out, err := cmd.Rodar(ctx, "ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "json", video)
	if err != nil {
		return 0, 0, "", err
	}
	var p struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out, &p); err != nil || len(p.Streams) == 0 {
		return 0, 0, "", fmt.Errorf("ffprobe sem stream de video em %s", video)
	}
	w, h := p.Streams[0].Width, p.Streams[0].Height
	switch {
	case w > h:
		return w, h, "16:9", nil
	case h > w:
		return w, h, "9:16", nil
	default:
		return w, h, "1:1", nil
	}
}

func garantirLogoMarca(projDir string, m roteiro.Marca) (string, error) {
	rel := m.Logo.Arquivo
	if strings.TrimSpace(rel) == "" {
		rel = "marca/logo.png"
	}
	p := filepath.Join(projDir, filepath.FromSlash(rel))
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", err
	}
	img := desenharEmblemaPadrao()
	f, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return p, png.Encode(f, img)
}

func desenharEmblemaPadrao() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 160, 160))
	ouro := color.RGBA{R: 255, G: 210, B: 74, A: 245}
	fundo := color.RGBA{R: 11, G: 14, B: 28, A: 225}
	for y := 0; y < 160; y++ {
		for x := 0; x < 160; x++ {
			dx, dy := x-80, y-80
			d2 := dx*dx + dy*dy
			if d2 <= 74*74 && d2 >= 62*62 {
				img.SetRGBA(x, y, ouro)
			} else if d2 < 62*62 {
				img.SetRGBA(x, y, fundo)
			}
		}
	}
	return img
}

func renderizarPartesMarca(ctx context.Context, projDir, vd, videoIn, logoPath, renderPy, fmtVideo string, w, h int, m roteiro.Marca, cmd ExecutorComando) ([]string, bool, bool, error) {
	var partes []string
	abOk, feOk := false, false
	if m.Vinhetas.Abertura.Ativo && m.Vinhetas.Abertura.DuracaoS > 0 {
		ab, err := gerarVinhetaMotor(ctx, projDir, vd, "abertura", m.Vinhetas.Abertura, m.Cores, fmtVideo, w, h, renderPy, cmd)
		if err != nil {
			return nil, false, false, err
		}
		partes, abOk = append(partes, ab), true
	}
	princ, err := sobreporLogoFFmpeg(ctx, vd, videoIn, logoPath, fmtVideo, w, h, m.Logo, cmd)
	if err != nil {
		return nil, false, false, err
	}
	partes = append(partes, princ)
	if m.Vinhetas.Fechamento.Ativo && m.Vinhetas.Fechamento.DuracaoS > 0 {
		fe, err := gerarVinhetaMotor(ctx, projDir, vd, "fechamento", m.Vinhetas.Fechamento, m.Cores, fmtVideo, w, h, renderPy, cmd)
		if err != nil {
			return nil, false, false, err
		}
		partes, feOk = append(partes, fe), true
	}
	return partes, abOk, feOk, nil
}

func gerarVinhetaMotor(ctx context.Context, projDir, vd, nome string, cfg roteiro.VinhetaCfg, cores roteiro.CoresMarca, fmtVideo string, w, h int, renderPy string, cmd ExecutorComando) (string, error) {
	foco := cores.Foco
	if foco == "" {
		foco = "#FFD24A"
	}
	cont := cores.Contorno
	if cont == "" {
		cont = "#111111"
	}
	tl := map[string]any{
		"fps": 30, "formato": fmtVideo,
		"shots": []map[string]any{{
			"id": nome, "lines": [2]int{0, 0}, "duration": cfg.DuracaoS, "fx": []string{"bokeh", "dust"},
			"cam":    []QuadroCamera{{T: 0, X: 960, Y: 540, Z: 1.0}, {T: cfg.DuracaoS, X: 960, Y: 540, Z: 1.08}},
			"layers": []map[string]any{{"kind": "text", "text": cfg.Texto, "x": 960.0, "y": 560.0, "size": 140.0, "font": "brush", "color": foco, "stroke": cont, "sw": 12.0, "start": 0.15, "in": "pop", "hud": true, "glow": 14}},
		}},
	}
	tp := filepath.Join(vd, nome+".json")
	b, _ := json.Marshal(tl)
	if err := os.WriteFile(tp, b, 0o644); err != nil {
		return "", err
	}
	crua := filepath.Join(vd, nome+".mp4")
	marcaPath := filepath.Join(projDir, "marca.json")
	if _, err := cmd.Rodar(ctx, "python", "-I", renderPy, tp, "--formato", fmtVideo, "--marca", marcaPath, "--legendas", "nao", "--sem-checagem", "--out", crua); err != nil {
		return "", err
	}
	norm := filepath.Join(vd, nome+"_s.mp4")
	_, err := cmd.Rodar(ctx, "ffmpeg", "-y", "-v", "error", "-i", crua, "-f", "lavfi", "-i", "anullsrc=r=48000:cl=stereo", "-shortest", "-vf", fmt.Sprintf("scale=%d:%d,fps=30", w, h), "-c:v", "libx264", "-crf", "19", "-pix_fmt", "yuv420p", "-c:a", "aac", "-ar", "48000", norm)
	return norm, err
}

func sobreporLogoFFmpeg(ctx context.Context, vd, videoIn, logoPath, fmtVideo string, w, h int, lg roteiro.LogoMarca, cmd ExecutorComando) (string, error) {
	pct := lg.LarguraPct
	if pct <= 0 {
		pct = 8
	}
	op := lg.Opacidade
	if op <= 0 {
		op = 0.85
	}
	fator := 1.5
	if fmtVideo != "16:9" {
		fator = 2.6
	}
	lw := max(48, int(float64(w)*float64(pct)/100.0*fator))
	mg := int(float64(w) * 0.03)
	px, py := posicaoOverlayLogo(lg.Posicao, fmtVideo, h, mg)
	mainOut := filepath.Join(vd, "principal.mp4")
	fc := fmt.Sprintf("[1]scale=%d:-1,format=rgba,colorchannelmixer=aa=%.2f[lg];[0:v][lg]overlay=%s:%s,fps=30[v]", lw, op, px, py)
	_, err := cmd.Rodar(ctx, "ffmpeg", "-y", "-v", "error", "-i", videoIn, "-i", logoPath, "-filter_complex", fc, "-map", "[v]", "-map", "0:a?", "-c:v", "libx264", "-crf", "19", "-pix_fmt", "yuv420p", "-c:a", "aac", "-ar", "48000", "-ac", "2", mainOut)
	return mainOut, err
}

func posicaoOverlayLogo(pos, fmtVideo string, h, mg int) (string, string) {
	px, py := fmt.Sprintf("W-w-%d", mg), fmt.Sprintf("%d", mg)
	switch pos {
	case "canto_superior_esquerdo":
		px, py = fmt.Sprintf("%d", mg), fmt.Sprintf("%d", mg)
	case "canto_inferior_direito":
		px, py = fmt.Sprintf("W-w-%d", mg), fmt.Sprintf("H-h-%d", mg)
	case "canto_inferior_esquerdo":
		px, py = fmt.Sprintf("%d", mg), fmt.Sprintf("H-h-%d", mg)
	}
	if fmtVideo == "9:16" && !strings.Contains(py, "H-h") {
		py = fmt.Sprintf("%d", int(float64(h)*0.16))
	}
	return px, py
}

func concatenarPartesMarca(ctx context.Context, partes []string, saida string, cmd ExecutorComando) error {
	if err := os.MkdirAll(filepath.Dir(saida), 0o755); err != nil {
		return err
	}
	var args []string
	args = append(args, "-y", "-v", "error")
	var streams strings.Builder
	for i, p := range partes {
		args = append(args, "-i", p)
		fmt.Fprintf(&streams, "[%d:v][%d:a]", i, i)
	}
	fc := fmt.Sprintf("%sconcat=n=%d:v=1:a=1[v][a]", streams.String(), len(partes))
	args = append(args, "-filter_complex", fc, "-map", "[v]", "-map", "[a]", "-c:v", "libx264", "-crf", "19", "-pix_fmt", "yuv420p", "-c:a", "aac", "-b:a", "192k", saida)
	_, err := cmd.Rodar(ctx, "ffmpeg", args...)
	return err
}
