package taxonomy

import (
	"errors"
	"strings"
)

// Erros sentinela do subsistema de papeis e taxonomia.
var (
	ErrPapelInvalido               = errors.New("taxonomy: papel operacional desconhecido ou fora do catalogo")
	ErrGuildaInvalida              = errors.New("taxonomy: guilda operacional desconhecida")
	ErrPapelPrimarioObrigatorio    = errors.New("matriz: peca deve possuir papel_primario valido")
	ErrNenhumaFerramentaEncontrada = errors.New("fabricar: nenhuma ferramenta compativel encontrada no cofre")
)

// Guilda representa uma divisao operacional da fabrica.
type Guilda string

const (
	GuildaGovernanca Guilda = "governanca"
	GuildaReversa    Guilda = "reversa"
	GuildaConstrucao Guilda = "construcao"
	GuildaInfra      Guilda = "infra"
	GuildaAudio      Guilda = "audiovisual"
	GuildaQualidade  Guilda = "qualidade"
	GuildaPublicacao Guilda = "publicacao"
)

// Papel representa um cracha canonico de execucao.
type Papel struct {
	Slug        string `json:"slug"`
	Nome        string `json:"nome"`
	Guilda      Guilda `json:"guilda"`
	Descricao   string `json:"descricao"`
	Ferramentas string `json:"ferramentas_tipicas"`
}

// CatalogoCanonico contem os crachas oficiais das 7 guildas.
var CatalogoCanonico = map[string]Papel{
	"patrocinador":        {Slug: "patrocinador", Nome: "Patrocinador", Guilda: GuildaGovernanca, Descricao: "Aprovacao de regras e Stop-Gates", Ferramentas: "Chat UI"},
	"arquiteto":           {Slug: "arquiteto", Nome: "Arquiteto", Guilda: GuildaGovernanca, Descricao: "Requisitos EARS e schemas", Ferramentas: "Templates v2"},
	"planejador":          {Slug: "planejador", Nome: "Planejador", Guilda: GuildaGovernanca, Descricao: "Divisao em tickets atomicos", Ferramentas: "PLAN v2"},
	"orquestrador":        {Slug: "orquestrador", Nome: "Orquestrador", Guilda: GuildaGovernanca, Descricao: "Despacho do DAG e isolamento esteril", Ferramentas: "LegoLoop CLI"},
	"dissecador-web":      {Slug: "dissecador-web", Nome: "Dissecador Web", Guilda: GuildaReversa, Descricao: "Extracao de tokens e rotas SaaS", Ferramentas: "REA Web"},
	"dissecador-mobile":   {Slug: "dissecador-mobile", Nome: "Dissecador Mobile", Guilda: GuildaReversa, Descricao: "Descompilacao de APKs", Ferramentas: "Jeb / Apktool"},
	"dissecador-midia":    {Slug: "dissecador-midia", Nome: "Dissecador Midia", Guilda: GuildaReversa, Descricao: "Forense de video e transcricao", Ferramentas: "Video-REA"},
	"engenheiro-frontend": {Slug: "engenheiro-frontend", Nome: "Engenheiro Frontend", Guilda: GuildaConstrucao, Descricao: "Interfaces Web e Next.js", Ferramentas: "React / Tailwind"},
	"engenheiro-mobile":   {Slug: "engenheiro-mobile", Nome: "Engenheiro Mobile", Guilda: GuildaConstrucao, Descricao: "Aplicativos moveis", Ferramentas: "Android / iOS"},
	"engenheiro-backend":  {Slug: "engenheiro-backend", Nome: "Engenheiro Backend", Guilda: GuildaConstrucao, Descricao: "APIs e logica de dominio", Ferramentas: "Go / REST"},
	"engenheiro-tdd":      {Slug: "engenheiro-tdd", Nome: "Engenheiro TDD", Guilda: GuildaConstrucao, Descricao: "Red-Green-Refactor e testes", Ferramentas: "Go Test / Vitest"},
	"engenheiro-cloud":    {Slug: "engenheiro-cloud", Nome: "Engenheiro Cloud", Guilda: GuildaInfra, Descricao: "VPS, Edge e Docker", Ferramentas: "Cloudflare / Ops"},
	"engenheiro-secops":   {Slug: "engenheiro-secops", Nome: "Engenheiro SecOps", Guilda: GuildaInfra, Descricao: "Path-jail e segredos", Ferramentas: "Security Gates"},
	"roteirista-copy":     {Slug: "roteirista-copy", Nome: "Roteirista Copy", Guilda: GuildaAudio, Descricao: "Roteirizacao estruturada", Ferramentas: "Craft Story"},
	"ingestor-midia":      {Slug: "ingestor-midia", Nome: "Ingestor Midia", Guilda: GuildaAudio, Descricao: "Download e FFMPEG", Ferramentas: "yt-dlp / Whisper"},
	"animador-motion":     {Slug: "animador-motion", Nome: "Animador Motion", Guilda: GuildaAudio, Descricao: "Motion procedural em codigo", Ferramentas: "Remotion / Three.js"},
	"pos-produtor":        {Slug: "pos-produtor", Nome: "Pos Produtor", Guilda: GuildaAudio, Descricao: "Cortes de silencio e RMS", Ferramentas: "Auto-cutter"},
	"auditor-qualidade":   {Slug: "auditor-qualidade", Nome: "Auditor Qualidade", Guilda: GuildaQualidade, Descricao: "Linters e mutacao de teste", Ferramentas: "Taxlinter / Funclen"},
	"redator-tecnico":     {Slug: "redator-tecnico", Nome: "Redator Tecnico", Guilda: GuildaPublicacao, Descricao: "Manuais ASD-STE100", Ferramentas: "Manual Creator"},
	"compilador-livro":    {Slug: "compilador-livro", Nome: "Compilador Livro", Guilda: GuildaPublicacao, Descricao: "Livros didaticos e PDF", Ferramentas: "LegoPress"},
}

var guildasOrdem = []Guilda{
	GuildaGovernanca,
	GuildaReversa,
	GuildaConstrucao,
	GuildaInfra,
	GuildaAudio,
	GuildaQualidade,
	GuildaPublicacao,
}

// ValidarPapel verifica se o slug pertence ao catalogo canonico.
func ValidarPapel(slug string) error {
	_, err := ObterPapel(slug)
	return err
}

// ObterPapel retorna o cracha correspondente ao slug ou ErrPapelInvalido.
func ObterPapel(slug string) (Papel, error) {
	s := strings.TrimSpace(slug)
	p, ok := CatalogoCanonico[s]
	if !ok {
		return Papel{}, ErrPapelInvalido
	}
	return p, nil
}

// ValidarGuilda verifica se a guilda pertence ao catalogo oficial.
func ValidarGuilda(g Guilda) error {
	norm := Guilda(strings.TrimSpace(string(g)))
	for _, oficial := range guildasOrdem {
		if norm == oficial {
			return nil
		}
	}
	return ErrGuildaInvalida
}

// ListarGuildas retorna a lista ordenada das 7 guildas oficiais.
func ListarGuildas() []Guilda {
	out := make([]Guilda, len(guildasOrdem))
	copy(out, guildasOrdem)
	return out
}
