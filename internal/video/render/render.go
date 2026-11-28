package render

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// FPS e a taxa de quadros do render; a tolerancia de sincronia e de 1 quadro.
const (
	FPS             = 30
	ToleranciaFrame = 1.0 / FPS
)

// Resolucao e uma saida do master 16:9.
type Resolucao struct {
	Nome         string
	Largura, Alt int
}

// resolucaoDe devolve a resolucao pelo nome; nome desconhecido devolve o valor zero.
func resolucaoDe(nome string) Resolucao {
	switch nome {
	case "720p":
		return Resolucao{"720p", 1280, 720}
	case "1080p":
		return Resolucao{"1080p", 1920, 1080}
	case "4k":
		return Resolucao{"4k", 3840, 2160}
	}
	return Resolucao{}
}

// RenderDeps isola ffmpeg, ffprobe e a sonda de encoders para testes hermeticos.
type RenderDeps struct {
	Exec  Executor
	Med   Medidor
	Sonda Sondador
}

func normalizarRenderDeps(d RenderDeps) RenderDeps {
	if d.Exec == nil {
		d.Exec = execSistema{}
	}
	if d.Med == nil {
		d.Med = execSistema{}
	}
	if d.Sonda == nil {
		d.Sonda = SondaFFmpeg{Med: d.Med}
	}
	return d
}

func quadrosDaBatida(b Batida) int {
	return max(arredondar(b.FimS*FPS)-arredondar(b.InicioS*FPS), 1)
}

// expressaoCamera devolve o zoompan da camera para uma batida de n quadros.
func expressaoCamera(camera string, n int) (string, error) {
	centroX, centroY := "iw/2-(iw/zoom/2)", "ih/2-(ih/zoom/2)"
	switch camera {
	case "zoom":
		return fmt.Sprintf("z='1+0.12*on/%d':x='%s':y='%s'", n, centroX, centroY), nil
	case "pan":
		return fmt.Sprintf("z='1.12':x='(iw-iw/zoom)*on/%d':y='%s'", n, centroY), nil
	case "shake":
		return fmt.Sprintf("z='1.06':x='%s+14*sin(on*1.7)':y='%s+9*cos(on*2.3)'", centroX, centroY), nil
	case "estatica":
		return fmt.Sprintf("z='1':x='%s':y='%s'", centroX, centroY), nil
	}
	return "", invalido("camera %q desconhecida", camera)
}

// escaparCaminhoFiltro prepara um caminho (Windows inclusive) para uso dentro de um filtro do ffmpeg.
func escaparCaminhoFiltro(p string) string {
	return strings.ReplaceAll(strings.ReplaceAll(p, "\\", "/"), ":", "\\:")
}

type planoRender struct {
	Entradas []string // imagens, relativas ao projeto, na ordem dos -i
	Filtro   string
}

// camadasVisuais separa o cenario das camadas com imagem, em ordem de profundidade.
// Adereco e efeito sem arquivo nao tem imagem recortada e ficam fora da composicao.
func camadasVisuais(b Batida) (cenario string, sobrepostas []CamadaItem, err error) {
	cs := append([]CamadaItem(nil), b.Camadas...)
	sort.SliceStable(cs, func(i, j int) bool { return cs[i].Profundidade < cs[j].Profundidade })
	for _, c := range cs {
		switch {
		case c.Camada == "cenario":
			cenario = c.Arquivo
		case c.Camada != "legenda" && c.Arquivo != "":
			sobrepostas = append(sobrepostas, c)
		}
	}
	if cenario == "" {
		err = entradaInvalida("batida da linha %d sem arquivo de cenario", b.IndiceLinha)
	}
	return
}

// filtroBatida monta o trecho de uma batida (camera, sprites) e devolve o proximo indice de entrada.
func filtroBatida(b Batida, i, proxima int) (partes, entradas []string, novaProxima int, err error) {
	cenario, sobre, err := camadasVisuais(b)
	if err != nil {
		return nil, nil, proxima, err
	}
	n := quadrosDaBatida(b)
	cam, err := expressaoCamera(b.Camera, n)
	if err != nil {
		return nil, nil, proxima, err
	}
	partes = append(partes, fmt.Sprintf("[%d:v]scale=2880:1620:force_original_aspect_ratio=increase,crop=2880:1620,zoompan=%s:d=%d:s=%dx%d:fps=%d,setsar=1[b%d_0]",
		proxima, cam, n, LarguraBase, AlturaBase, FPS, i))
	entradas = append(entradas, cenario)
	proxima++
	for k, c := range sobre {
		partes = append(partes,
			fmt.Sprintf("[%d:v]scale=%d:%d[s%d_%d]", proxima, c.Caixa.W, c.Caixa.H, i, k),
			fmt.Sprintf("[b%d_%d][s%d_%d]overlay=%d:%d[b%d_%d]", i, k, i, k, c.Caixa.X, c.Caixa.Y, i, k+1))
		entradas = append(entradas, c.Arquivo)
		proxima++
	}
	partes = append(partes, fmt.Sprintf("[b%d_%d]null[v%d]", i, len(sobre), i))
	return partes, entradas, proxima, nil
}

// montarFiltroRender monta o filter_complex do master: batidas, concat, legenda ASS e escala final.
func montarFiltroRender(tl Timeline, res Resolucao, ass string) (planoRender, error) {
	var p planoRender
	var partes []string
	rotulos := ""
	proxima := 0
	for i, b := range tl.Batidas {
		ps, es, prox, err := filtroBatida(b, i, proxima)
		if err != nil {
			return p, err
		}
		partes, proxima = append(partes, ps...), prox
		p.Entradas = append(p.Entradas, es...)
		rotulos += fmt.Sprintf("[v%d]", i)
	}
	partes = append(partes, fmt.Sprintf("%sconcat=n=%d:v=1:a=0[vc]", rotulos, len(tl.Batidas)),
		fmt.Sprintf("[vc]ass=filename='%s',scale=%d:%d:flags=lanczos,format=yuv420p[v]", escaparCaminhoFiltro(ass), res.Largura, res.Alt))
	p.Filtro = strings.Join(partes, ";")
	return p, nil
}

func opcoesEncoder(enc string) []string {
	switch enc {
	case "h264_nvenc":
		return []string{"-preset", "p5", "-rc", "vbr", "-cq", "19", "-b:v", "0"}
	case "h264_qsv":
		return []string{"-global_quality", "19"}
	case "h264_amf":
		return []string{"-quality", "quality", "-rc", "cqp", "-qp_i", "19", "-qp_p", "19"}
	}
	return []string{"-preset", "medium", "-crf", "18"}
}

// argsRender monta o comando do ffmpeg; imagens e soundtrack ja vem como caminhos absolutos.
func argsRender(imagens []string, soundtrack, filtro, saida, enc string) []string {
	args := []string{"-y", "-hide_banner", "-nostats"}
	for _, im := range imagens {
		args = append(args, "-i", im)
	}
	args = append(args, "-i", soundtrack, "-filter_complex", filtro, "-map", "[v]", "-map", strconv.Itoa(len(imagens))+":a", "-c:v", enc)
	args = append(args, opcoesEncoder(enc)...)
	return append(args, "-r", strconv.Itoa(FPS), "-pix_fmt", "yuv420p", "-c:a", "aac", "-b:a", "192k", "-movflags", "+faststart", saida)
}

var reVersao = regexp.MustCompile(`^v(\d+)(_[a-z0-9]+)?\.mp4$`)

// proximaVersao devolve N de vN.mp4: um acima da maior versao existente (variantes contam).
func proximaVersao(dir string) (int, error) {
	if st, err := os.Stat(dir); err == nil && !st.IsDir() {
		return 0, fmt.Errorf("%s nao e uma pasta", dir)
	}
	entradas, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	maior := 0
	for _, e := range entradas {
		if m := reVersao.FindStringSubmatch(e.Name()); m != nil {
			n, _ := strconv.Atoi(m[1])
			maior = max(maior, n)
		}
	}
	return maior + 1, nil
}

// ---- gate de integridade (AC-E05.3) ----

func renderReprovado(formato string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrRenderInvalido, fmt.Sprintf(formato, args...))
}

func duracoesDoMP4(ctx context.Context, med Medidor, mp4 string) (video, audio float64, err error) {
	saida, err := med.Saida(ctx, "ffprobe", "-v", "error", "-show_entries", "stream=codec_type,duration", "-of", "json", mp4)
	if err != nil {
		return 0, 0, renderReprovado("ffprobe falhou: %v", err)
	}
	var r struct {
		Streams []struct {
			Tipo    string `json:"codec_type"`
			Duracao string `json:"duration"`
		} `json:"streams"`
	}
	if err := json.Unmarshal([]byte(saida), &r); err != nil {
		return 0, 0, renderReprovado("saida do ffprobe ilegivel: %v", err)
	}
	achou := map[string]bool{}
	for _, s := range r.Streams {
		d, err := strconv.ParseFloat(s.Duracao, 64)
		if err != nil {
			return 0, 0, renderReprovado("duracao %q da trilha %s ilegivel", s.Duracao, s.Tipo)
		}
		switch s.Tipo {
		case "video":
			video = d
		case "audio":
			audio = d
		}
		achou[s.Tipo] = true
	}
	if !achou["video"] || !achou["audio"] {
		return 0, 0, renderReprovado("o MP4 precisa de trilha de video e de audio")
	}
	return video, audio, nil
}

// quadroEmBranco procura quadros pretos; com negar, procura quadros brancos.
func quadroEmBranco(ctx context.Context, med Medidor, mp4 string, negar bool) (bool, error) {
	filtro := "blackdetect=d=0.03:pix_th=0.10"
	if negar {
		filtro = "negate," + filtro
	}
	saida, err := med.Saida(ctx, "ffmpeg", "-hide_banner", "-nostats", "-i", mp4, "-an", "-vf", filtro, "-f", "null", "-")
	if err != nil {
		return false, renderReprovado("deteccao de quadros em branco falhou: %v", err)
	}
	return strings.Contains(saida, "black_start"), nil
}

// ConferirRender exige audio e video com a mesma duracao (tolerancia de 1 quadro) e nenhum quadro em branco.
func ConferirRender(ctx context.Context, med Medidor, mp4 string) error {
	v, a, err := duracoesDoMP4(ctx, med, mp4)
	if err != nil {
		return err
	}
	if math.Abs(v-a) > ToleranciaFrame+1e-6 {
		return renderReprovado("video dura %.3f s e audio %.3f s (tolerancia %.3f s)", v, a, ToleranciaFrame)
	}
	for _, negar := range []bool{false, true} {
		achou, err := quadroEmBranco(ctx, med, mp4, negar)
		if err != nil {
			return err
		}
		if achou {
			cor := map[bool]string{false: "preto", true: "branco"}[negar]
			return renderReprovado("o render tem quadros em %s", cor)
		}
	}
	return nil
}

// ---- execucao ----

type resumoEntradas struct {
	provedor string
	lufs, tp float64
}

func existeArquivo(projDir, rel string) (string, error) {
	abs, err := resolverSeguro(projDir, rel)
	if err != nil {
		return "", bloqueio(rel, err)
	}
	if st, err := os.Stat(abs); err != nil || st.IsDir() {
		return "", bloqueio(rel, fmt.Errorf("arquivo ausente"))
	}
	return abs, nil
}

func lerResumo(projDir string) (resumoEntradas, error) {
	var r resumoEntradas
	var tm struct {
		Provedor string `json:"provedor"`
	}
	var mx struct {
		LUFS float64 `json:"lufs"`
		TP   float64 `json:"true_peak_dbtp"`
	}
	if err := lerJSON(projDir, "data/timings.json", &tm, false); err != nil {
		return r, bloqueio("data/timings.json", err)
	}
	if err := lerJSON(projDir, "data/mix.json", &mx, false); err != nil {
		return r, bloqueio("data/mix.json", err)
	}
	return resumoEntradas{provedor: tm.Provedor, lufs: mx.LUFS, tp: mx.TP}, nil
}

// imagensAbsolutas confere e resolve todas as imagens da timeline antes de qualquer render.
func imagensAbsolutas(projDir string, rel []string) ([]string, error) {
	abs := make([]string, len(rel))
	for i, r := range rel {
		var err error
		if abs[i], err = existeArquivo(projDir, r); err != nil {
			return nil, err
		}
	}
	return abs, nil
}

// ---- render pelo motor portado (AC-E10.4) ----

// MotorRender roda o motor portado sobre timeline.json e confere o MP4 9:16 contra a voz.
type MotorRender interface {
	Renderizar(ctx context.Context, projDir string, vozS float64, fps int) error
}

// RelatorioMotor resume o que o render Go reconfirmou antes de chamar o motor.
type RelatorioMotor struct {
	ProvedorVoz  string
	DuracaoVozS  float64
	LUFS         float64
	TruePeakDBTP float64
	EventosSFX   int
	Cortes       int
}

// DuracaoVoz le a duracao de audio/voice.wav do projeto.
func DuracaoVoz(projDir string) (float64, error) { return duracaoDaVoz(projDir) }

func duracaoDaVoz(projDir string) (float64, error) {
	b, err := lerBytes(projDir, "audio/voice.wav")
	if err != nil {
		return 0, bloqueio("audio/voice.wav", err)
	}
	f, dados, err := lerWAV(b)
	if err != nil {
		return 0, err
	}
	return duracaoWAV(f, len(dados)), nil
}

// reconfirmarSFX exige o arquivo de troca_cena do catalogo quando a nova timeline tem cortes.
func reconfirmarSFX(projDir string, tl TimelineMotor) (int, error) {
	var c ContratosRender
	if err := carregarConfig(projDir, &c); err != nil {
		return 0, err
	}
	cat, err := CarregarCatalogo(projDir, c.Config)
	if err != nil {
		return 0, err
	}
	n := tl.eventosTrocaCena()
	if _, err := resolverSFX(projDir, cat, "troca_cena"); n > 0 && err != nil {
		return 0, err
	}
	return n, nil
}

// reconfirmarMix exige data/mix.json da mesma voz, no alvo de loudness e com os SFX da nova timeline.
func reconfirmarMix(projDir string, vozS float64, eventos int) (RelatorioMix, error) {
	var mx RelatorioMix
	if err := lerJSON(projDir, "data/mix.json", &mx, false); err != nil {
		return mx, bloqueio("data/mix.json", err)
	}
	if math.Abs(mx.DuracaoS-vozS) > ToleranciaFrame+1e-6 {
		return mx, renderReprovado("mix dura %.3f s e a voz %.3f s", mx.DuracaoS, vozS)
	}
	if err := conforme(Medicao{LUFS: mx.LUFS, TruePeakDBTP: mx.TruePeakDBTP}); err != nil {
		return mx, renderReprovado("mix fora do alvo: %v", err)
	}
	if mx.Eventos != eventos {
		return mx, renderReprovado("mix tem %d eventos de SFX e a nova timeline %d", mx.Eventos, eventos)
	}
	return mx, nil
}

func nomeDaSaida(n int, res string, primaria bool) string {
	if primaria {
		return fmt.Sprintf("v%d.mp4", n)
	}
	return fmt.Sprintf("v%d_%s.mp4", n, res)
}
