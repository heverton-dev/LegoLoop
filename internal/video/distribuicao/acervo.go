// Package distribuicao implementa o P4 do processo composto legoloop-video-animado:
// revisao, vitrine, acervo, agenda, publicacao via Postiz, comentarios e aprendizado.
package distribuicao

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// Erros tipados do P4.
var (
	ErrAcervoInvalido        = errors.New("video-distribuicao: acervo invalido")
	ErrTemaRepetido          = errors.New("video-distribuicao: tema repetido ou quase igual a um tema do acervo")
	ErrSemAutorizacao        = errors.New("video-distribuicao: item sem autorizacao de publicacao")
	ErrPublicoSemConfirmacao = errors.New("video-distribuicao: visibilidade publica exige confirmacao explicita")
	ErrSemRotuloIA           = errors.New("video-distribuicao: rotulo de conteudo gerado por IA e obrigatorio")
	ErrVitrineInvalida       = errors.New("video-distribuicao: vitrine invalida")
	ErrParecerInvalido       = errors.New("video-distribuicao: parecer invalido")
	ErrRevisaoReprovada      = errors.New("video-distribuicao: revisao nao liberou o video")
	ErrSegredoAusente        = errors.New("video-distribuicao: segredo ausente")
	ErrPostiz                = errors.New("video-distribuicao: falha na API do Postiz")
	ErrPedidoInvalido        = errors.New("video-distribuicao: pedido invalido")
)

// Tipos de pauta e estados do acervo.
const (
	TipoQuente     = "quente"
	TipoPerene     = "perene"
	StatusPronto   = "pronto"
	StatusAgendado = "agendado"
	StatusPublic   = "publicado"
	StatusExpirado = "expirado"

	VisPrivado = "private"
	VisPublico = "public"
)

// Item e um video ou corte aprovado esperando a hora certa.
type Item struct {
	ID           string     `json:"id"`
	Tema         string     `json:"tema"`
	Tipo         string     `json:"tipo"`
	Formatos     []string   `json:"formatos"`
	Canais       []string   `json:"canais"`
	Validade     time.Time  `json:"validade,omitempty"`
	Prioridade   int        `json:"prioridade"`
	Status       string     `json:"status"`
	Agendado     *time.Time `json:"agendado,omitempty"`
	Links        []string   `json:"links,omitempty"`
	Autorizado   bool       `json:"autorizado"`
	Visibilidade string     `json:"visibilidade"`
	Origem       string     `json:"origem,omitempty"`
}

// Acervo e o estoque de itens aprovados (acervo/acervo.json).
type Acervo struct {
	SchemaVersion int    `json:"schema_version"`
	Itens         []Item `json:"itens"`
}

// CaminhoAcervo devolve o caminho do acervo dentro do projeto.
func CaminhoAcervo(projDir string) string {
	return filepath.Join(projDir, "acervo", "acervo.json")
}

// CarregarAcervo le o acervo; arquivo ausente devolve acervo vazio.
func CarregarAcervo(projDir string) (Acervo, error) {
	a := Acervo{SchemaVersion: 1}
	dados, err := os.ReadFile(CaminhoAcervo(projDir))
	if errors.Is(err, os.ErrNotExist) {
		return a, nil
	}
	if err != nil {
		return a, err
	}
	if err := json.Unmarshal(dados, &a); err != nil {
		return Acervo{}, fmt.Errorf("%w: %v", ErrAcervoInvalido, err)
	}
	return a, a.Validar()
}

// Gravar escreve o acervo de forma atomica.
func (a Acervo) Gravar(projDir string) error {
	if err := a.Validar(); err != nil {
		return err
	}
	dados, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	return gravarAtomico(CaminhoAcervo(projDir), append(dados, '\n'))
}

func gravarAtomico(caminho string, dados []byte) error {
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		return err
	}
	tmp := caminho + ".tmp"
	if err := os.WriteFile(tmp, dados, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, caminho)
}

// Validar confere cada item e a unicidade dos ids.
func (a Acervo) Validar() error {
	vistos := map[string]bool{}
	for _, it := range a.Itens {
		if err := it.Validar(); err != nil {
			return err
		}
		if vistos[it.ID] {
			return fmt.Errorf("%w: id %q repetido", ErrAcervoInvalido, it.ID)
		}
		vistos[it.ID] = true
	}
	return nil
}

// Validar confere os campos obrigatorios de um item.
func (it Item) Validar() error {
	switch {
	case strings.TrimSpace(it.ID) == "" || strings.TrimSpace(it.Tema) == "":
		return fmt.Errorf("%w: id e tema sao obrigatorios", ErrAcervoInvalido)
	case it.Tipo != TipoQuente && it.Tipo != TipoPerene:
		return fmt.Errorf("%w: item %q com tipo %q", ErrAcervoInvalido, it.ID, it.Tipo)
	case !slices.Contains([]string{StatusPronto, StatusAgendado, StatusPublic, StatusExpirado}, it.Status):
		return fmt.Errorf("%w: item %q com status %q", ErrAcervoInvalido, it.ID, it.Status)
	case it.Visibilidade != VisPrivado && it.Visibilidade != VisPublico:
		return fmt.Errorf("%w: item %q com visibilidade %q", ErrAcervoInvalido, it.ID, it.Visibilidade)
	case len(it.Canais) == 0:
		return fmt.Errorf("%w: item %q sem canais", ErrAcervoInvalido, it.ID)
	case it.Tipo == TipoQuente && it.Validade.IsZero():
		return fmt.Errorf("%w: item quente %q exige validade", ErrAcervoInvalido, it.ID)
	}
	return nil
}

// Adicionar valida o item e aplica a regra de originalidade contra os temas existentes.
func (a *Acervo) Adicionar(it Item) error {
	if err := it.Validar(); err != nil {
		return err
	}
	for _, o := range a.Itens {
		if o.ID == it.ID {
			return fmt.Errorf("%w: id %q repetido", ErrAcervoInvalido, it.ID)
		}
		if it.Origem == "" && TemasSemelhantes(o.Tema, it.Tema) {
			return fmt.Errorf("%w: %q parece com %q", ErrTemaRepetido, it.Tema, o.Tema)
		}
	}
	a.Itens = append(a.Itens, it)
	return nil
}

// Autorizar grava ou revoga a autorizacao de publicacao; revogar so vale antes da publicacao.
func (a *Acervo) Autorizar(id string, valor bool) error {
	it := a.Buscar(id)
	if it == nil {
		return fmt.Errorf("%w: item %q inexistente", ErrAcervoInvalido, id)
	}
	if it.Status == StatusPublic {
		return fmt.Errorf("%w: item %q ja publicado", ErrAcervoInvalido, id)
	}
	it.Autorizado = valor
	return nil
}

// Buscar devolve um ponteiro para o item ou nil.
func (a *Acervo) Buscar(id string) *Item {
	for i := range a.Itens {
		if a.Itens[i].ID == id {
			return &a.Itens[i]
		}
	}
	return nil
}

// TemasSemelhantes e verdadeiro quando os temas sao iguais ou compartilham 80% ou mais das palavras.
func TemasSemelhantes(a, b string) bool {
	pa, pb := palavras(a), palavras(b)
	if len(pa) == 0 || len(pb) == 0 {
		return false
	}
	inter := 0
	for p := range pa {
		if pb[p] {
			inter++
		}
	}
	uniao := len(pa) + len(pb) - inter
	return float64(inter)/float64(uniao) >= 0.8
}

func palavras(s string) map[string]bool {
	out := map[string]bool{}
	for _, p := range strings.Fields(strings.ToLower(s)) {
		p = strings.Trim(p, ".,;:!?\"'()-")
		if p != "" {
			out[p] = true
		}
	}
	return out
}
