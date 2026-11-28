package distribuicao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Vereditos validos do revisor sem contexto.
const (
	VereditoEnviar    = "ENVIAR"
	VereditoNaoEnviar = "NAO ENVIAR"
)

// QuadrosPorFolha e o numero de quadros de cada folha de contato (6x4).
const QuadrosPorFolha = 24

// RegrasRevisor devolve as regras que o revisor sem contexto confere (P7 do processo original).
func RegrasRevisor() []string {
	return []string{
		"corte ou movimento de camera a cada 3 s (2 a 2,5 s em 9:16 e 1:1)",
		"imagem coerente com a frase dita naquele segundo",
		"numeros iguais ao roteiro",
		"texto legivel no celular",
		"nada importante cortado pela borda",
		"personagem com o mesmo rosto e a mesma roupa em todo o video",
		"sem maos ou membros extras",
		"sem texto falso nas imagens",
		"marca aplicada (fontes, cores, logo)",
		"areas seguras do formato respeitadas",
		"volume perto de -14 LUFS",
	}
}

// ChecklistRevisor devolve a autorrevisao obrigatoria "sem cara de IA" antes do revisor.
func ChecklistRevisor() []string {
	return []string{"maos", "texto falso", "olhar vazio", "pele de plastico", "estilo que muda entre cenas"}
}

var reTempo = regexp.MustCompile(`^\d{1,3}:\d{2}$`)

// Problema aponta um defeito por versao, corte e minutagem.
type Problema struct {
	Versao   int    `json:"versao"`
	Corte    string `json:"corte"`
	Tempo    string `json:"tempo"`
	Problema string `json:"problema"`
	Correcao string `json:"correcao"`
}

// Parecer e a resposta do revisor sem contexto para uma versao.
type Parecer struct {
	SchemaVersion int        `json:"schema_version"`
	Versao        int        `json:"versao"`
	Veredito      string     `json:"veredito"`
	Problemas     []Problema `json:"problemas"`
	Revisor       string     `json:"revisor"`
}

// Validar confere veredito, minutagem e a regra: NAO ENVIAR exige ao menos um problema.
func (p Parecer) Validar() error {
	if p.Versao < 1 || strings.TrimSpace(p.Revisor) == "" {
		return fmt.Errorf("%w: versao e revisor sao obrigatorios", ErrParecerInvalido)
	}
	if p.Veredito != VereditoEnviar && p.Veredito != VereditoNaoEnviar {
		return fmt.Errorf("%w: veredito %q", ErrParecerInvalido, p.Veredito)
	}
	if p.Veredito == VereditoNaoEnviar && len(p.Problemas) == 0 {
		return fmt.Errorf("%w: NAO ENVIAR exige a lista de problemas", ErrParecerInvalido)
	}
	for _, pr := range p.Problemas {
		if !reTempo.MatchString(pr.Tempo) || strings.TrimSpace(pr.Corte) == "" || strings.TrimSpace(pr.Problema) == "" || strings.TrimSpace(pr.Correcao) == "" {
			return fmt.Errorf("%w: problema sem corte, minutagem mm:ss, descricao ou correcao", ErrParecerInvalido)
		}
	}
	return nil
}

// CaminhoParecer devolve review/vN_parecer.json.
func CaminhoParecer(projDir string, versao int) string {
	return filepath.Join(projDir, "review", fmt.Sprintf("v%d_parecer.json", versao))
}

// GravarParecer grava o parecer e recusa sobrescrever: nenhuma versao revisada e apagada.
func GravarParecer(projDir string, p Parecer) error {
	p.SchemaVersion = 1
	if err := p.Validar(); err != nil {
		return err
	}
	caminho := CaminhoParecer(projDir, p.Versao)
	if _, err := os.Stat(caminho); err == nil {
		return fmt.Errorf("%w: v%d ja tem parecer", ErrParecerInvalido, p.Versao)
	}
	dados, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return gravarAtomico(caminho, append(dados, '\n'))
}

// Liberado devolve nil somente quando o parecer da versao e ENVIAR.
func Liberado(projDir string, versao int) error {
	dados, err := os.ReadFile(CaminhoParecer(projDir, versao))
	if err != nil {
		return fmt.Errorf("%w: v%d sem parecer", ErrRevisaoReprovada, versao)
	}
	var p Parecer
	if err := json.Unmarshal(dados, &p); err != nil {
		return fmt.Errorf("%w: %v", ErrParecerInvalido, err)
	}
	if err := p.Validar(); err != nil {
		return err
	}
	if p.Veredito != VereditoEnviar {
		return fmt.Errorf("%w: v%d recebeu %s", ErrRevisaoReprovada, versao, p.Veredito)
	}
	return nil
}

var reVersao = regexp.MustCompile(`^v(\d+)\.mp4$`)

// ProximaVersao devolve o numero da proxima versao em renders/ (v1 quando nao ha nenhuma).
func ProximaVersao(projDir string) (int, error) {
	entradas, err := os.ReadDir(filepath.Join(projDir, "renders"))
	if errors.Is(err, os.ErrNotExist) {
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

// Medidor devolve a duracao de um video em segundos.
type Medidor func(ctx context.Context, video string) (float64, error)

// DuracaoFFprobe mede a duracao pelo ffprobe.
func DuracaoFFprobe(ctx context.Context, video string) (float64, error) {
	saida, err := exec.CommandContext(ctx, "ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "csv=p=0", video).Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe: %w", err)
	}
	return strconv.ParseFloat(strings.TrimSpace(string(saida)), 64)
}

// GerarFolhas cria review/sheet_NN.jpg com um quadro a cada passo segundos (24 por folha).
func GerarFolhas(ctx context.Context, ex Executor, medir Medidor, projDir, video string, passo float64) ([]string, error) {
	if passo <= 0 {
		return nil, fmt.Errorf("%w: passo deve ser maior que zero", ErrPedidoInvalido)
	}
	if err := existe(projDir, video); err != nil {
		return nil, err
	}
	dur, err := medir(ctx, filepath.Join(projDir, filepath.FromSlash(video)))
	if err != nil {
		return nil, err
	}
	n := int(math.Ceil(dur / (QuadrosPorFolha * passo)))
	var folhas []string
	for k := 1; k <= n; k++ {
		rel := fmt.Sprintf("review/sheet_%02d.jpg", k)
		inicio := float64(k-1) * QuadrosPorFolha * passo
		if err := os.MkdirAll(filepath.Join(projDir, "review"), 0o755); err != nil {
			return folhas, err
		}
		args := argsFolha(filepath.Join(projDir, filepath.FromSlash(video)), filepath.Join(projDir, filepath.FromSlash(rel)), inicio, passo)
		if err := ex.Executar(ctx, "ffmpeg", args...); err != nil {
			return folhas, err
		}
		folhas = append(folhas, rel)
	}
	return folhas, nil
}

func argsFolha(video, saida string, inicio, passo float64) []string {
	vf := fmt.Sprintf("fps=1/%g,scale=384:-1,tile=6x4", passo)
	return []string{"-y", "-v", "error", "-ss", fmt.Sprintf("%g", inicio), "-t", fmt.Sprintf("%g", QuadrosPorFolha*passo), "-i", video, "-vf", vf, "-frames:v", "1", saida}
}

// PacoteRevisor e o que o revisor sem contexto recebe: video, folhas, roteiro e regras. Nada mais.
type PacoteRevisor struct {
	Versao    int      `json:"versao"`
	Video     string   `json:"video"`
	Folhas    []string `json:"folhas"`
	Roteiro   string   `json:"roteiro"`
	Regras    []string `json:"regras"`
	Checklist []string `json:"checklist"`
	PassoS    float64  `json:"passo_s"`
}

// GravarPacote grava review/vN_pacote.json; exige renders/vN.mp4, script.json e as folhas.
func GravarPacote(projDir string, versao int, folhas []string, passo float64) (string, error) {
	video := fmt.Sprintf("renders/v%d.mp4", versao)
	for _, rel := range append([]string{video, "script.json"}, folhas...) {
		if err := existe(projDir, rel); err != nil {
			return "", err
		}
	}
	if len(folhas) == 0 {
		return "", fmt.Errorf("%w: pacote sem folhas de contato", ErrPedidoInvalido)
	}
	p := PacoteRevisor{Versao: versao, Video: video, Folhas: folhas, Roteiro: "script.json", Regras: RegrasRevisor(), Checklist: ChecklistRevisor(), PassoS: passo}
	dados, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	caminho := filepath.Join(projDir, "review", fmt.Sprintf("v%d_pacote.json", versao))
	return caminho, gravarAtomico(caminho, append(dados, '\n'))
}

// Liberacao e o registro auditavel de que a versao foi liberada para a vitrine.
type Liberacao struct {
	Versao     int       `json:"versao"`
	Veredito   string    `json:"veredito"`
	LiberadoEm time.Time `json:"liberado_em"`
}

// CaminhoLiberacao devolve review/vN_liberada.json.
func CaminhoLiberacao(projDir string, versao int) string {
	return filepath.Join(projDir, "review", fmt.Sprintf("v%d_liberada.json", versao))
}

// Liberar confere o parecer ENVIAR e grava o registro de liberacao da versao.
func Liberar(projDir string, versao int, agora time.Time) error {
	if err := Liberado(projDir, versao); err != nil {
		return err
	}
	dados, err := json.MarshalIndent(Liberacao{Versao: versao, Veredito: VereditoEnviar, LiberadoEm: agora.UTC()}, "", "  ")
	if err != nil {
		return err
	}
	return gravarAtomico(CaminhoLiberacao(projDir, versao), append(dados, '\n'))
}

var reVideoVersao = regexp.MustCompile(`^renders/v(\d+)\.mp4$`)

// ExigirLiberada recusa a vitrine de um video cuja versao nao foi liberada pela revisao.
func ExigirLiberada(projDir, video string) error {
	m := reVideoVersao.FindStringSubmatch(filepath.ToSlash(video))
	if m == nil {
		return fmt.Errorf("%w: video %q fora do padrao renders/vN.mp4", ErrRevisaoReprovada, video)
	}
	n, _ := strconv.Atoi(m[1])
	if _, err := os.Stat(CaminhoLiberacao(projDir, n)); err != nil {
		return fmt.Errorf("%w: v%d sem registro de liberacao", ErrRevisaoReprovada, n)
	}
	return Liberado(projDir, n)
}
