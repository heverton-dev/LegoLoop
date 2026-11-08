package redc

import (
	"errors"
)

var (
	ErrCanalInvalido    = errors.New("redc: identificador do canal inválido")
	ErrVideosInsuf      = errors.New("redc: lote de vídeos insuficiente (mínimo 1 vídeo)")
	ErrDadosCorrompidos = errors.New("redc: dados brutos de métricas corrompidos")
	ErrTemaObrigatorio  = errors.New("redc: tema do usuário é obrigatório para geração")
)

// VideoEntry representa a análise individual de um vídeo virall.
type VideoEntry struct {
	ID          string  `json:"id"`
	Titulo      string  `json:"titulo"`
	Views       int64   `json:"views"`
	Likes       int64   `json:"likes"`
	Comentarios int64   `json:"comentarios"`
	DuracaoSeg  int     `json:"duracao_seg"`
	Gancho3s    string  `json:"gancho_3s"`
	Metodo      string  `json:"metodo"`
	ClusterTema string  `json:"cluster_tema"`
	Transcricao string  `json:"transcricao"`
	Score       float64 `json:"score,omitempty"`
}

// ChannelDNA encapsula a inteligência extraída dos top vídeos virais.
type ChannelDNA struct {
	CanalSlug    string            `json:"canal_slug"`
	TomDeVoz     string            `json:"tom_de_voz"`
	Arquetipo    string            `json:"arquetipo"`
	RitmoCortes  string            `json:"ritmo_cortes"`
	DesignTokens map[string]string `json:"design_tokens"`
	MotionTokens map[string]string `json:"motion_tokens"`
	TopGanchos   []string          `json:"top_ganchos"`
}

// BrandTokens define os atributos de personalização da marca do usuário.
type BrandTokens struct {
	NomeMarca     string `json:"nome_marca"`
	CorPrimaria   string `json:"cor_primaria"`
	CorSecundaria string `json:"cor_secundaria"`
	FonteTitulo   string `json:"fonte_titulo"`
	FonteCorpo    string `json:"fonte_corpo"`
	TomPreferido  string `json:"tom_preferido"`
}

// GenerationRequest parâmetros para geração de novo vídeo parametrizado.
type GenerationRequest struct {
	CanalSlug string      `json:"canal_slug"`
	Tema      string      `json:"tema"`
	Marca     BrandTokens `json:"marca"`
}

// ChannelReport relatório resumido das operações executadas pelo motor.
type ChannelReport struct {
	CanalSlug       string   `json:"canal_slug"`
	TotalVideos     int      `json:"total_videos"`
	AmostraDeepREDC int      `json:"amostra_deep_redc"`
	MetodosUnicos   int      `json:"metodos_unicos"`
	DiretorioSaida  string   `json:"diretorio_saida"`
	ArquivosCriados []string `json:"arquivos_criados"`
}
