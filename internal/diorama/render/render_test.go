package render

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/diorama/ativos"
)

// tlRender: duas batidas de 1 s (pan e zoom), cenario e sprite com arquivo, legenda so na primeira.
func tlRender() Timeline {
	mk := func(i int, ini, fim float64, cam string, pal []PalavraTempo) Batida {
		return Batida{IndiceLinha: i, InicioS: ini, FimS: fim, CenarioID: "cen_01", Camera: cam, Camadas: []CamadaItem{
			{ID: "cenario:cen_01", Camada: "cenario", Profundidade: 0, Arquivo: "flow/cen_01.png", Caixa: ativos.BoundingBox{W: 1920, H: 1080}},
			{ID: "o", Camada: "aderecos", Profundidade: 1, Caixa: ativos.BoundingBox{X: 10, Y: 10, W: 50, H: 50}},
			{ID: "p", Camada: "sprite", Profundidade: 2, SpriteID: "narrador_pose_1", Arquivo: "sprites/narrador_pose_1.png", Caixa: ativos.BoundingBox{X: 100, Y: 50, W: 200, H: 400}},
			{ID: "legenda:x", Camada: "legenda", Profundidade: 3, Palavras: pal, Caixa: ativos.BoundingBox{W: 8, H: 2}},
		}}
	}
	pal := []PalavraTempo{{Texto: "Olá", InicioS: 0, FimS: 0.4}, {Texto: "mundo", InicioS: 0.4, FimS: 0.9}}
	return Timeline{SchemaVersion: 1, DuracaoS: 2, Batidas: []Batida{mk(0, 0, 1, "pan", pal), mk(1, 1, 2, "zoom", nil)}}
}

func TestRenderASSKaraoke(t *testing.T) {
	ass := GerarASS(tlRender())
	for _, trecho := range []string{"PlayResX: 1920", "PlayResY: 1080", "Dialogue: 0,0:00:00.00,0:00:00.90,Karaoke", "{\\k40}Olá {\\k50}mundo"} {
		if !strings.Contains(ass, trecho) {
			t.Errorf("ASS sem %q:\n%s", trecho, ass)
		}
	}
	if strings.Count(ass, "Dialogue:") != 1 {
		t.Errorf("so a batida com palavras gera Dialogue:\n%s", ass)
	}
}

func TestRenderASSQuebraLinhaELacunas(t *testing.T) {
	var ps []PalavraTempo
	for i := 0; i < 14; i++ {
		ps = append(ps, PalavraTempo{Texto: "palavra", InicioS: float64(i) * 0.5, FimS: float64(i)*0.5 + 0.3})
	}
	tl := tlRender()
	tl.Batidas[0].Camadas[3].Palavras = ps
	ass := GerarASS(tl)
	if !strings.Contains(ass, "\\N") {
		t.Errorf("texto longo deve quebrar em 2 linhas:\n%s", ass)
	}
	if !strings.Contains(ass, "{\\k30}palavra {\\k20}") {
		t.Errorf("a lacuna entre palavras deve virar tempo vazio:\n%s", ass)
	}
}

func TestRenderFiltroCameraSpriteEConcat(t *testing.T) {
	p, err := montarFiltroRender(tlRender(), Resolucao{"1080p", 1920, 1080}, "/p/legenda.ass")
	if err != nil {
		t.Fatal(err)
	}
	for _, trecho := range []string{
		"d=30:s=1920x1080:fps=30", "zoompan=z='1.12'", "zoompan=z='1+0.12*on/30'",
		"scale=200:400", "overlay=100:50", "concat=n=2:v=1:a=0", "scale=1920:1080", "format=yuv420p[v]",
	} {
		if !strings.Contains(p.Filtro, trecho) {
			t.Errorf("filtro sem %q:\n%s", trecho, p.Filtro)
		}
	}
	if strings.Contains(p.Filtro, "scale=50:50") {
		t.Errorf("adereco sem arquivo nao pode virar overlay:\n%s", p.Filtro)
	}
	if want := []string{"flow/cen_01.png", "sprites/narrador_pose_1.png", "flow/cen_01.png", "sprites/narrador_pose_1.png"}; strings.Join(p.Entradas, ",") != strings.Join(want, ",") {
		t.Errorf("entradas: %v", p.Entradas)
	}
}

func TestRenderFiltroCamerasConhecidas(t *testing.T) {
	casos := map[string]string{"zoom": "1+0.12*on/30", "pan": "z='1.12'", "shake": "sin(on*1.7)", "estatica": "z='1'"}
	for cam, quer := range casos {
		tl := tlRender()
		tl.Batidas[0].Camera = cam
		p, err := montarFiltroRender(tl, Resolucao{"720p", 1280, 720}, "a.ass")
		if err != nil || !strings.Contains(p.Filtro, quer) || !strings.Contains(p.Filtro, "scale=1280:720") {
			t.Errorf("camera %s: %v\n%s", cam, err, p.Filtro)
		}
	}
	tl := tlRender()
	tl.Batidas[0].Camera = "giro"
	if _, err := montarFiltroRender(tl, Resolucao{"720p", 1280, 720}, "a.ass"); !errors.Is(err, ErrSchemaInvalido) {
		t.Errorf("camera desconhecida deve falhar: %v", err)
	}
}

func TestRenderQuadrosSomamADuracao(t *testing.T) {
	tl := Timeline{SchemaVersion: 1, DuracaoS: 52.872, Batidas: []Batida{{InicioS: 0, FimS: 1.016}, {InicioS: 1.016, FimS: 52.872}}}
	soma := 0
	for _, b := range tl.Batidas {
		soma += quadrosDaBatida(b)
	}
	if soma != 1586 {
		t.Fatalf("quadros somam %d e a duracao pede 1586", soma)
	}
}

func TestRenderEscapaCaminhoDoFiltro(t *testing.T) {
	if got := escaparCaminhoFiltro(`C:\proj\renders\a b.ass`); got != `C\:/proj/renders/a b.ass` {
		t.Fatalf("escape: %s", got)
	}
}

func TestRenderProximaVersao(t *testing.T) {
	dir := t.TempDir()
	if n, err := proximaVersao(dir); err != nil || n != 1 {
		t.Fatalf("pasta vazia: %d %v", n, err)
	}
	for _, f := range []string{"v1.mp4", "v2_720p.mp4", "legenda.srt", "vx.mp4"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if n, err := proximaVersao(dir); err != nil || n != 3 {
		t.Fatalf("deve pular v1 e v2 (inclusive variantes): %d %v", n, err)
	}
	if _, err := proximaVersao(filepath.Join(dir, "v1.mp4")); err == nil {
		t.Fatal("pasta que e arquivo deve falhar")
	}
}

// ---- gate de integridade ----

type medRenderFalso struct {
	probe  string
	preto  string
	branco string
	erros  map[string]error
}

func (m *medRenderFalso) Saida(_ context.Context, nome string, args ...string) (string, error) {
	j := strings.Join(args, " ")
	switch {
	case nome == "ffprobe":
		return m.probe, m.erros["probe"]
	case strings.Contains(j, "negate,blackdetect"):
		return m.branco, m.erros["branco"]
	case strings.Contains(j, "blackdetect"):
		return m.preto, m.erros["preto"]
	}
	return "", nil
}

func probe(v, a string) string {
	return fmt.Sprintf(`{"streams":[{"codec_type":"video","duration":%q},{"codec_type":"audio","duration":%q}]}`, v, a)
}

func TestRenderGateAprovaDuracoesIguaisESemQuadrosEmBranco(t *testing.T) {
	for _, par := range [][2]string{{"2.000000", "2.000000"}, {"2.000000", "2.030000"}, {"1.967000", "2.000000"}} {
		if err := ConferirRender(context.Background(), &medRenderFalso{probe: probe(par[0], par[1])}, "a.mp4"); err != nil {
			t.Errorf("%v dentro de 1 quadro deve passar: %v", par, err)
		}
	}
}

func TestRenderGateReprova(t *testing.T) {
	casos := map[string]*medRenderFalso{
		"audio mais curto": {probe: probe("2.000000", "1.900000")},
		"video mais curto": {probe: probe("1.900000", "2.000000")},
		"quadro preto":     {probe: probe("2", "2"), preto: "[blackdetect @ 0x1] black_start:0.5 black_end:0.6"},
		"quadro branco":    {probe: probe("2", "2"), branco: "[blackdetect @ 0x1] black_start:1 black_end:1.2"},
		"sem trilha audio": {probe: `{"streams":[{"codec_type":"video","duration":"2"}]}`},
		"sem trilha video": {probe: `{"streams":[{"codec_type":"audio","duration":"2"}]}`},
		"probe ilegivel":   {probe: "lixo"},
		"duracao ilegivel": {probe: probe("N/A", "2")},
		"probe falhou":     {erros: map[string]error{"probe": errors.New("x")}},
		"deteccao preto":   {probe: probe("2", "2"), erros: map[string]error{"preto": errors.New("x")}},
		"deteccao branco":  {probe: probe("2", "2"), erros: map[string]error{"branco": errors.New("x")}},
	}
	for nome, m := range casos {
		if err := ConferirRender(context.Background(), m, "a.mp4"); !errors.Is(err, ErrRenderInvalido) {
			t.Errorf("%s: esperava ErrRenderInvalido, obteve %v", nome, err)
		}
	}
}

// ---- execucao do render ----

type execRenderFalso struct {
	comandos [][]string
	falham   map[string]bool
}

func (e *execRenderFalso) Executar(_ context.Context, _ string, args ...string) error {
	e.comandos = append(e.comandos, args)
	j := strings.Join(args, " ")
	for enc := range e.falham {
		if strings.Contains(j, "-c:v "+enc+" ") {
			return errors.New("falha de " + enc)
		}
	}
	return os.WriteFile(args[len(args)-1], []byte("mp4"), 0o644)
}

func projetoRender(t *testing.T, extra map[string]string) string {
	t.Helper()
	b, _ := json.Marshal(tlRender())
	arqs := map[string]string{
		"timeline.json": string(b), "audio/soundtrack.wav": "wav", "flow/cen_01.png": "png",
		"sprites/narrador_pose_1.png": "png",
		"data/mix.json":               `{"schema_version":1,"duracao_s":2,"eventos_sfx":0,"lufs":-14,"true_peak_dbtp":-2,"tentativas":[]}`,
		"data/timings.json":           `{"schema_version":1,"provedor":"edge-tts","voz":"pt-BR-AntonioNeural","linhas":[]}`,
	}
	for k, v := range extra {
		arqs[k] = v
	}
	return projetoFixture(t, arqs)
}

func depsRender(ex *execRenderFalso, lista ...string) RenderDeps {
	return RenderDeps{Exec: ex, Med: &medRenderFalso{probe: probe("2.000000", "2.000000")}, Sonda: sondaFalsa{lista: lista}}
}

func TestRenderGravaMasterV1ERelatorio(t *testing.T) {
	dir := projetoRender(t, nil)
	ex := &execRenderFalso{}
	rel, err := ExecutarRender(context.Background(), dir, depsRender(ex, "libx264", "h264_nvenc"))
	if err != nil {
		t.Fatal(err)
	}
	if !existe(dir, "renders/v1.mp4") || existe(dir, "renders/legenda.ass") == false {
		t.Fatalf("esperava renders/v1.mp4 e renders/legenda.ass")
	}
	if rel.Encoder != "h264_nvenc" || rel.ProvedorVoz != "edge-tts" || rel.LUFS != -14 || rel.TruePeakDBTP != -2 {
		t.Fatalf("relatorio: %+v", rel)
	}
	if rel.Renders["1080p"] != "renders/v1.mp4" || len(rel.Renders) != 1 {
		t.Fatalf("renders: %v", rel.Renders)
	}
	if err := rel.Validar(); err != nil {
		t.Fatal(err)
	}
	var nomes []string
	entradas, _ := os.ReadDir(filepath.Join(dir, "renders"))
	for _, e := range entradas {
		nomes = append(nomes, e.Name())
	}
	if strings.Join(nomes, ",") != "legenda.ass,v1.mp4" {
		t.Fatalf("so os arquivos finais podem ficar em renders/: %v", nomes)
	}
}

func TestRenderResolucoesMultiplasEQuatroKSoSeListada(t *testing.T) {
	dir := projetoRender(t, map[string]string{"render.json": `{"schema_version":1,"resolucoes":["1080p","720p","4k"]}`})
	ex := &execRenderFalso{}
	rel, err := ExecutarRender(context.Background(), dir, depsRender(ex))
	if err != nil {
		t.Fatal(err)
	}
	quer := map[string]string{"1080p": "renders/v1.mp4", "720p": "renders/v1_720p.mp4", "4k": "renders/v1_4k.mp4"}
	for r, caminho := range quer {
		if rel.Renders[r] != caminho || !existe(dir, caminho) {
			t.Errorf("%s: %v", r, rel.Renders)
		}
	}
	if len(ex.comandos) != 3 {
		t.Fatalf("esperava 3 execucoes, houve %d", len(ex.comandos))
	}
	padrao := projetoRender(t, nil)
	rel, err = ExecutarRender(context.Background(), padrao, depsRender(&execRenderFalso{}))
	if err != nil || existe(padrao, "renders/v1_4k.mp4") || rel.Renders["4k"] != "" {
		t.Fatalf("4k so quando listado: %v %v", rel.Renders, err)
	}
}

func TestRenderNaoSobrescreveVersaoAnterior(t *testing.T) {
	dir := projetoRender(t, map[string]string{"renders/v1.mp4": "antigo"})
	rel, err := ExecutarRender(context.Background(), dir, depsRender(&execRenderFalso{}))
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(dir, "renders", "v1.mp4"))
	if string(b) != "antigo" || rel.Renders["1080p"] != "renders/v2.mp4" || !existe(dir, "renders/v2.mp4") {
		t.Fatalf("v1 deve ficar intacto e sair v2: %q %v", b, rel.Renders)
	}
}

func TestRenderHardwareQueFalhaCaiNoLibx264(t *testing.T) {
	dir := projetoRender(t, nil)
	ex := &execRenderFalso{falham: map[string]bool{"h264_nvenc": true}}
	rel, err := ExecutarRender(context.Background(), dir, depsRender(ex, "h264_nvenc"))
	if err != nil || rel.Encoder != "libx264" || !existe(dir, "renders/v1.mp4") {
		t.Fatalf("fallback: %+v %v", rel, err)
	}
}

func TestRenderGateReprovadoNaoGravaVersao(t *testing.T) {
	dir := projetoRender(t, nil)
	deps := depsRender(&execRenderFalso{})
	deps.Med = &medRenderFalso{probe: probe("2.000000", "1.000000")}
	_, err := ExecutarRender(context.Background(), dir, deps)
	if !errors.Is(err, ErrRenderInvalido) {
		t.Fatalf("esperava ErrRenderInvalido: %v", err)
	}
	entradas, _ := os.ReadDir(filepath.Join(dir, "renders"))
	if len(entradas) != 0 {
		t.Fatalf("render reprovado nao pode deixar arquivos: %v", entradas)
	}
}

func TestRenderEntradasAusentesBloqueiam(t *testing.T) {
	for _, falta := range []string{"audio/soundtrack.wav", "flow/cen_01.png", "sprites/narrador_pose_1.png", "timeline.json", "data/mix.json", "data/timings.json"} {
		dir := projetoRender(t, nil)
		if err := os.Remove(filepath.Join(dir, filepath.FromSlash(falta))); err != nil {
			t.Fatal(err)
		}
		ex := &execRenderFalso{}
		if _, err := ExecutarRender(context.Background(), dir, depsRender(ex)); err == nil || len(ex.comandos) != 0 {
			t.Errorf("sem %s deve bloquear antes de renderizar: %v", falta, err)
		}
	}
}

func TestRenderFalhaDoFFmpegNaoDeixaTemporario(t *testing.T) {
	dir := projetoRender(t, nil)
	ex := &execRenderFalso{falham: map[string]bool{"libx264": true}}
	if _, err := ExecutarRender(context.Background(), dir, depsRender(ex)); err == nil {
		t.Fatal("falha do ffmpeg deve propagar")
	}
	if entradas, _ := os.ReadDir(filepath.Join(dir, "renders")); len(entradas) != 0 {
		t.Fatalf("sobrou arquivo: %v", entradas)
	}
}

func TestRenderContextoCancelado(t *testing.T) {
	ctx, cancela := context.WithCancel(context.Background())
	cancela()
	if _, err := ExecutarRender(ctx, projetoRender(t, nil), depsRender(&execRenderFalso{})); !errors.Is(err, context.Canceled) {
		t.Fatalf("esperava cancelamento: %v", err)
	}
}

func pngCor(t *testing.T, caminho string, w, h int, c color.RGBA) {
	t.Helper()
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			im.Set(x, y, c)
		}
	}
	f, err := os.Create(caminho)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, im); err != nil {
		t.Fatal(err)
	}
}

func TestRenderComFFmpegReal(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg ausente")
	}
	dir := projetoRender(t, map[string]string{"render.json": `{"schema_version":1,"resolucoes":["720p"]}`})
	pngCor(t, filepath.Join(dir, "flow", "cen_01.png"), 640, 360, color.RGBA{200, 120, 60, 255})
	pngCor(t, filepath.Join(dir, "sprites", "narrador_pose_1.png"), 100, 200, color.RGBA{30, 90, 200, 255})
	wav := montarWAV(formatoWAV{Canais: 2, Taxa: 48000, Bits: 16}, make([]byte, 96000*4))
	if err := os.WriteFile(filepath.Join(dir, "audio", "soundtrack.wav"), wav, 0o644); err != nil {
		t.Fatal(err)
	}
	rel, err := ExecutarRender(context.Background(), dir, RenderDeps{})
	if err != nil {
		t.Fatalf("render real: %v", err)
	}
	if !existe(dir, "renders/v1.mp4") || rel.Encoder == "" {
		t.Fatalf("sem mp4: %+v", rel)
	}
}
