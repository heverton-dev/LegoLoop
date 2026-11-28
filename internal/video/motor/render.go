package motor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Erros tipados do render de referência (AC-E10.4, AC-E10.5, AC-E10.7).
var (
	ErrRenderMotor        = errors.New("motor-render: o motor falhou ou não gerou o MP4")
	ErrRenderFormato      = errors.New("motor-render: MP4 fora de 1080x1920 a 30 quadros por segundo")
	ErrRenderSemAudio     = errors.New("motor-render: MP4 sem trilha de áudio")
	ErrRenderDuracao      = errors.New("motor-render: duração do MP4 difere da voz em mais de 1 quadro")
	ErrRenderQuadroBranco = errors.New("motor-render: MP4 com quadro em branco")
	ErrCacheDivergente    = errors.New("motor-render: cache por cena diverge entre renders")
	ErrBocaParada         = errors.New("motor-render: boca parada durante a fala")
	ErrBocaForaDaFala     = errors.New("motor-render: boca mexendo fora da fala")
)

const (
	// ArquivoRender é o MP4 final que o motor grava na pasta do projeto.
	ArquivoRender = "render.mp4"
	// DirCache guarda um arquivo por cena renderizada.
	DirCache = "cache"
	// larguraVertical e alturaVertical definem o quadro 9:16.
	larguraVertical = 1080
	alturaVertical  = 1920
	// limiarBoca é a diferença média mínima que separa dois estados de boca.
	limiarBoca = 0.05
)

// Executor roda o motor sobre dir/timeline.json e grava dir/render.mp4 e dir/cache/.
type Executor interface {
	Executar(dir string) error
}

// ExecutorPython roda o script de render do motor com python.
type ExecutorPython struct {
	Script string
	Args   []string
}

// Executar chama python <script> <dir> <args...> e devolve a saída no erro.
func (e ExecutorPython) Executar(dir string) error {
	args := append([]string{e.Script, dir}, e.Args...)
	saida, err := exec.Command("python", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(saida)))
	}
	return nil
}

// Midia é a medida de um MP4 renderizado.
type Midia struct {
	Largura        int
	Altura         int
	FPS            float64
	TemAudio       bool
	DuracaoS       float64
	TrechosBrancos int
}

// Sonda mede um MP4.
type Sonda interface {
	Medir(mp4 string) (Midia, error)
}

// PedidoRender descreve o render de uma pasta de projeto no formato pedido (padrao 9:16).
type PedidoRender struct {
	Dir         string
	Formato     string
	DuracaoVozS float64
	FPS         int
}

// DimensoesFormato devolve largura e altura de saida para 16:9, 9:16 (padrao quando vazio) e 1:1.
func DimensoesFormato(formato string) (int, int, error) {
	switch strings.TrimSpace(formato) {
	case "", "9:16":
		return larguraVertical, alturaVertical, nil
	case "16:9":
		return 1920, 1080, nil
	case "1:1":
		return 1080, 1080, nil
	default:
		return 0, 0, fmt.Errorf("%w: formato %q invalido", ErrRenderFormato, formato)
	}
}

// Renderizar roda o motor, mede o MP4 e confere dimensoes do formato, audio, duracao e zero brancos.
func Renderizar(p PedidoRender, ex Executor, s Sonda) (Midia, error) {
	if p.FPS <= 0 {
		p.FPS = fpsPadrao
	}
	if _, _, err := DimensoesFormato(p.Formato); err != nil {
		return Midia{}, err
	}
	if err := ex.Executar(p.Dir); err != nil {
		return Midia{}, fmt.Errorf("%w: %v", ErrRenderMotor, err)
	}
	mp4 := filepath.Join(p.Dir, ArquivoRender)
	if _, err := os.Stat(mp4); err != nil {
		return Midia{}, fmt.Errorf("%w: %v", ErrRenderMotor, err)
	}
	m, err := s.Medir(mp4)
	if err != nil {
		return Midia{}, fmt.Errorf("%w: sonda: %v", ErrRenderMotor, err)
	}
	return m, ValidarMidiaFormato(m, p.Formato, p.DuracaoVozS, p.FPS)
}

// ValidarMidia confere 1080x1920 (9:16), fps, áudio, duração da voz (± 1 quadro) e zero brancos.
func ValidarMidia(m Midia, vozS float64, fps int) error {
	return ValidarMidiaFormato(m, "9:16", vozS, fps)
}

// ValidarMidiaFormato confere as dimensoes do formato (16:9, 9:16 ou 1:1), fps, audio, duracao e zero brancos.
func ValidarMidiaFormato(m Midia, formato string, vozS float64, fps int) error {
	w, h, err := DimensoesFormato(formato)
	if err != nil {
		return err
	}
	if m.Largura != w || m.Altura != h || math.Abs(m.FPS-float64(fps)) > 0.01 {
		return fmt.Errorf("%w: %dx%d a %.2f fps (esperado %dx%d)", ErrRenderFormato, m.Largura, m.Altura, m.FPS, w, h)
	}
	if !m.TemAudio {
		return ErrRenderSemAudio
	}
	if math.Abs(m.DuracaoS-vozS) > 1.0/float64(fps) {
		return fmt.Errorf("%w: MP4 %.3fs, voz %.3fs", ErrRenderDuracao, m.DuracaoS, vozS)
	}
	if m.TrechosBrancos > 0 {
		return fmt.Errorf("%w: %d trecho(s)", ErrRenderQuadroBranco, m.TrechosBrancos)
	}
	return nil
}

// HashesCache devolve o SHA-256 de cada arquivo de dir/cache/, pelo caminho relativo.
func HashesCache(dir string) (map[string]string, error) {
	raiz := filepath.Join(dir, DirCache)
	hashes := map[string]string{}
	err := filepath.WalkDir(raiz, func(caminho string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		dados, err := os.ReadFile(caminho)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(raiz, caminho)
		if err != nil {
			return err
		}
		soma := sha256.Sum256(dados)
		hashes[filepath.ToSlash(rel)] = hex.EncodeToString(soma[:])
		return nil
	})
	return hashes, err
}

// CompararCaches exige os mesmos arquivos de cena com os mesmos hashes.
func CompararCaches(a, b map[string]string) error {
	if len(a) == 0 || len(a) != len(b) {
		return fmt.Errorf("%w: %d contra %d arquivo(s) de cena", ErrCacheDivergente, len(a), len(b))
	}
	for nome, hash := range a {
		if b[nome] != hash {
			return fmt.Errorf("%w: cena %q", ErrCacheDivergente, nome)
		}
	}
	return nil
}

// ConferirDeterminismo renderiza a mesma timeline duas vezes e compara o cache (AC-E10.5).
func ConferirDeterminismo(dir string, ex Executor) (map[string]string, error) {
	var rodadas [2]map[string]string
	for i := range rodadas {
		if err := ex.Executar(dir); err != nil {
			return nil, fmt.Errorf("%w: rodada %d: %v", ErrRenderMotor, i+1, err)
		}
		hashes, err := HashesCache(dir)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		rodadas[i] = hashes
	}
	if err := CompararCaches(rodadas[0], rodadas[1]); err != nil {
		return nil, err
	}
	return rodadas[0], nil
}

// EstadosBoca conta os estados distintos da região da boca por diferença de imagem.
func EstadosBoca(quadros []image.Image, regiao image.Rectangle) int {
	var estados []image.Image
	for _, q := range quadros {
		novo := true
		for _, e := range estados {
			if diferencaMedia(q, e, regiao) <= limiarBoca {
				novo = false
				break
			}
		}
		if novo {
			estados = append(estados, q)
		}
	}
	return len(estados)
}

// diferencaMedia devolve a diferença RGB média entre 0 e 1 dentro da região.
func diferencaMedia(a, b image.Image, regiao image.Rectangle) float64 {
	soma, n := 0.0, 0
	for y := regiao.Min.Y; y < regiao.Max.Y; y++ {
		for x := regiao.Min.X; x < regiao.Max.X; x++ {
			r1, g1, b1, _ := a.At(x, y).RGBA()
			r2, g2, b2, _ := b.At(x, y).RGBA()
			soma += absDif(r1, r2) + absDif(g1, g2) + absDif(b1, b2)
			n += 3
		}
	}
	if n == 0 {
		return 0
	}
	return soma / float64(n) / 0xffff
}

func absDif(a, b uint32) float64 {
	return math.Abs(float64(a) - float64(b))
}

// ValidarBoca exige 2+ estados na fala e 1 estado no silêncio (AC-E10.7).
func ValidarBoca(fala, silencio []image.Image, regiao image.Rectangle) error {
	if n := EstadosBoca(fala, regiao); n < 2 {
		return fmt.Errorf("%w: %d estado(s)", ErrBocaParada, n)
	}
	if n := EstadosBoca(silencio, regiao); n != 1 {
		return fmt.Errorf("%w: %d estados", ErrBocaForaDaFala, n)
	}
	return nil
}

// SondaFFmpeg mede o MP4 com ffprobe e procura trechos pretos ou brancos com ffmpeg.
type SondaFFmpeg struct{}

// sondaFFprobe espelha a saída JSON do ffprobe.
type sondaFFprobe struct {
	Streams []struct {
		Tipo    string `json:"codec_type"`
		Largura int    `json:"width"`
		Altura  int    `json:"height"`
		Taxa    string `json:"r_frame_rate"`
	} `json:"streams"`
	Formato struct {
		Duracao string `json:"duration"`
	} `json:"format"`
}

// Medir lê dimensão, fps, áudio, duração e trechos em branco do MP4.
func (SondaFFmpeg) Medir(mp4 string) (Midia, error) {
	saida, err := exec.Command("ffprobe", "-v", "error", "-show_entries",
		"stream=codec_type,width,height,r_frame_rate:format=duration", "-of", "json", mp4).Output()
	if err != nil {
		return Midia{}, fmt.Errorf("ffprobe: %w", err)
	}
	m, err := lerFFprobe(saida)
	if err != nil {
		return Midia{}, err
	}
	m.TrechosBrancos, err = contarBrancos(mp4)
	return m, err
}

// lerFFprobe converte o JSON do ffprobe em Midia.
func lerFFprobe(dados []byte) (Midia, error) {
	var p sondaFFprobe
	if err := json.Unmarshal(dados, &p); err != nil {
		return Midia{}, fmt.Errorf("ffprobe json: %w", err)
	}
	var m Midia
	for _, s := range p.Streams {
		switch s.Tipo {
		case "video":
			m.Largura, m.Altura, m.FPS = s.Largura, s.Altura, lerTaxa(s.Taxa)
		case "audio":
			m.TemAudio = true
		}
	}
	dur, err := strconv.ParseFloat(p.Formato.Duracao, 64)
	if err != nil {
		return Midia{}, fmt.Errorf("ffprobe duração %q: %w", p.Formato.Duracao, err)
	}
	m.DuracaoS = dur
	return m, nil
}

// lerTaxa converte "30/1" em 30.
func lerTaxa(taxa string) float64 {
	num, den, achou := strings.Cut(taxa, "/")
	n, errN := strconv.ParseFloat(num, 64)
	if errN != nil {
		return 0
	}
	if !achou {
		return n
	}
	d, errD := strconv.ParseFloat(den, 64)
	if errD != nil || d == 0 {
		return 0
	}
	return n / d
}

// contarBrancos conta trechos pretos e, com a imagem negada, trechos brancos com duracao >= 1.2s.
func contarBrancos(mp4 string) (int, error) {
	total := 0
	for _, filtro := range []string{"blackdetect=d=1.2:pix_th=0.10", "negate,blackdetect=d=1.2:pix_th=0.10"} {
		saida, err := exec.Command("ffmpeg", "-v", "info", "-nostats", "-i", mp4,
			"-vf", filtro, "-an", "-f", "null", "-").CombinedOutput()
		if err != nil {
			return 0, fmt.Errorf("ffmpeg %s: %w: %s", filtro, err, strings.TrimSpace(string(saida)))
		}
		total += strings.Count(string(saida), "black_start:")
	}
	return total, nil
}
