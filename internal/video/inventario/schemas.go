// Package inventario guarda estilos e itens reutilizáveis do vídeo animado.
package inventario

import (
	"errors"
	"fmt"
)

// Erros tipados do inventário (SPEC 2026-11-26_SPEC_VIDEO_ANIMADO_V2 §4).
var (
	ErrEstiloDuplicado   = errors.New("video-inventario: estilo ja existe")
	ErrEstiloInexistente = errors.New("video-inventario: estilo inexistente")
	ErrEstiloNaoValidado = errors.New("video-inventario: estilo nao validado nao aceita item aprovado")
	ErrTipoInvalido      = errors.New("video-inventario: tipo de item invalido")
	ErrHashDivergente    = errors.New("video-inventario: hash do arquivo diverge")
	ErrHerancaCiclica    = errors.New("video-inventario: heranca de estilos em ciclo")
	ErrParteNaoSuportada = errors.New("video-inventario: parte nao suportada pelo renderizador")
	ErrGeneroInvalido    = errors.New("video-inventario: genero invalido")
	ErrFaixaInvalida     = errors.New("video-inventario: faixa etaria invalida")
)

// Estilo descreve um padrão visual e o renderizador que o desenha.
type Estilo struct {
	SchemaVersion    int          `json:"schema_version"`
	ID               string       `json:"id"` // kebab-case
	Nome             string       `json:"nome"`
	Status           string       `json:"status"`            // "rascunho" | "validado"
	Versao           string       `json:"versao"`            // semver
	Herda            string       `json:"herda,omitempty"`   // id de outro estilo
	Renderizador     string       `json:"renderizador"`      // ex.: "molde-pilula"
	PartesSuportadas []string     `json:"partes_suportadas"` // partes que o renderizador desenha
	PromptMestre     string       `json:"prompt_mestre"`     // prompt de estilo dos cenários do Flow
	Paleta           []string     `json:"paleta"`
	Luz              *Luz         `json:"luz,omitempty"`
	Referencias      []Referencia `json:"referencias"`
}

// Luz descreve a luz de recorte do estilo.
type Luz struct {
	Rim        bool    `json:"rim"`
	Cor        string  `json:"cor,omitempty"`
	Quantidade float64 `json:"quantidade,omitempty"`
}

// Referencia aponta uma imagem, vídeo ou link de origem do estilo.
type Referencia struct {
	Tipo   string `json:"tipo"`   // "imagem" | "video" | "link"
	Origem string `json:"origem"` // caminho ou URL
	Nota   string `json:"nota,omitempty"`
}

// Item é uma peça do inventário de um estilo ou da área comum.
type Item struct {
	SchemaVersion int            `json:"schema_version"`
	ID            string         `json:"id"`
	Tipo          string         `json:"tipo"`
	Estilo        string         `json:"estilo"`
	EstiloVersao  string         `json:"estilo_versao"`
	Tags          []string       `json:"tags"`       // tema, época, canal
	Parametros    map[string]any `json:"parametros"` // cores, forma, bloco cast
	Arquivo       string         `json:"arquivo,omitempty"`
	Hash          string         `json:"hash,omitempty"` // SHA-256 do arquivo
	Origem        string         `json:"origem"`         // projeto e data
	Aprovado      bool           `json:"aprovado"`
	Genero        string         `json:"genero,omitempty"`       // personagem: "masculino" | "feminino" | "neutro"
	FaixaEtaria   string         `json:"faixa_etaria,omitempty"` // personagem: "crianca" | "adulto" | "idoso"
	Herdavel      bool           `json:"herdavel,omitempty"`
	Obsoleto      bool           `json:"obsoleto,omitempty"`
}

// tipoPermitido evita variavel global mutavel: valida o tipo por switch.
func tipoPermitido(tipo string, comum bool) bool {
	if comum {
		switch tipo {
		case "sfx", "trilha", "voz":
			return true
		}
		return false
	}
	switch tipo {
	case "pele", "traje", "penteado", "adorno", "aderecos", "pose", "humor", "personagem", "cenario", "objeto":
		return true
	}
	return false
}

// ValidarTipo aceita tipos de estilo fora de comum/ e só sfx, trilha e voz em comum/.
func ValidarTipo(tipo string, comum bool) error {
	if !tipoPermitido(tipo, comum) {
		return fmt.Errorf("%w: %q (comum=%t)", ErrTipoInvalido, tipo, comum)
	}
	return nil
}
