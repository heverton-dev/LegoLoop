package fechar

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	// ErrSlugInvalido indica que o identificador do processo é inválido.
	ErrSlugInvalido = errors.New("fechar: slug do processo inválido")
	// ErrArtefatoAusente indica que um arquivo essencial não foi encontrado.
	ErrArtefatoAusente = errors.New("fechar: artefato obrigatório ausente")
	// ErrPortaoFalhou indica que a validação de portão mecânico falhou.
	ErrPortaoFalhou = errors.New("fechar: portão de qualidade reprovado")
)

var reSlugLimpo = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Request representa os parâmetros de disparo do fechamento.
type Request struct {
	Slug      string `json:"slug"`
	Titulo    string `json:"titulo,omitempty"`
	Autor     string `json:"autor,omitempty"`
	DryRun    bool   `json:"dry_run"`
	SkipBuild bool   `json:"skip_build"`
}

// Validate confere se a request cumpre o contrato básico antes de executar.
func (r Request) Validate() error {
	_, _, err := NormalizeSlug(r.Slug)
	return err
}

// NormalizeSlug limpa e padroniza o slug (ex: "video-motion" -> "video-motion", "legoloop-video-motion").
func NormalizeSlug(raw string) (string, string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", "", ErrSlugInvalido
	}
	s = strings.TrimPrefix(s, "legoloop-")
	if !reSlugLimpo.MatchString(s) {
		return "", "", ErrSlugInvalido
	}
	canon := "legoloop-" + s
	return s, canon, nil
}

// EtapaStatus registra o resultado atômico de uma das 8 etapas de fechamento.
type EtapaStatus struct {
	Etapa    int           `json:"etapa"`
	Nome     string        `json:"nome"`
	Status   string        `json:"status"` // "ok", "skipped", "error"
	Artefato string        `json:"artefato,omitempty"`
	Duracao  time.Duration `json:"duracao_ms"`
	Detalhe  string        `json:"detalhe,omitempty"`
}

// Result resume a execução do pipeline de fechamento.
type Result struct {
	Slug      string        `json:"slug"`
	CanonName string        `json:"canonical_name"`
	Sucesso   bool          `json:"sucesso"`
	Etapas    []EtapaStatus `json:"etapas"`
	Duracao   time.Duration `json:"duracao_total_ms"`
	Timestamp time.Time     `json:"timestamp"`
}
