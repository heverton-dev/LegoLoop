package inventario

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Erros da validação humana de estilo (AC-E07.8).
var (
	ErrAprovacaoAusente  = errors.New("video-inventario: arquivo de aprovacao humana ausente")
	ErrAprovacaoInvalida = errors.New("video-inventario: arquivo de aprovacao humana invalido")
)

// Aprovacao é o Stop-Gate humano gravado em arquivo pelo usuário.
type Aprovacao struct {
	Estilo           string   `json:"estilo"`
	AprovadoPor      string   `json:"aprovado_por"`
	Data             string   `json:"data,omitempty"`
	PersonagensTeste []string `json:"personagens_teste"`
	Nota             string   `json:"nota,omitempty"`
}

// Validacao registra a data e o arquivo de aprovação ao lado de estilo.json.
type Validacao struct {
	Estilo           string   `json:"estilo"`
	Versao           string   `json:"versao"`
	ValidadoEm       string   `json:"validado_em"`
	ArquivoAprovacao string   `json:"arquivo_aprovacao"`
	AprovadoPor      string   `json:"aprovado_por"`
	PersonagensTeste []string `json:"personagens_teste"`
}

func (s Store) caminhoValidacao(id string) string {
	return filepath.Join(s.dirEstilos(), id, "validacao.json")
}

// LerAprovacao lê e confere o arquivo de aprovação humana do estilo id.
func LerAprovacao(caminho, id string) (Aprovacao, error) {
	var a Aprovacao
	if caminho == "" {
		return a, fmt.Errorf("%w: caminho vazio", ErrAprovacaoAusente)
	}
	b, err := os.ReadFile(caminho)
	if err != nil {
		return a, fmt.Errorf("%w: %s: %v", ErrAprovacaoAusente, caminho, err)
	}
	if err := json.Unmarshal(b, &a); err != nil {
		return a, fmt.Errorf("%w: %s: %v", ErrAprovacaoInvalida, caminho, err)
	}
	switch {
	case a.Estilo != id:
		return a, fmt.Errorf("%w: estilo %q, esperado %q", ErrAprovacaoInvalida, a.Estilo, id)
	case a.AprovadoPor == "":
		return a, fmt.Errorf("%w: aprovado_por vazio", ErrAprovacaoInvalida)
	case len(a.PersonagensTeste) == 0:
		return a, fmt.Errorf("%w: sem personagens de teste aprovados", ErrAprovacaoInvalida)
	}
	return a, nil
}

// ValidarEstilo muda o estilo de rascunho para validado só com o arquivo de aprovação humana.
func (s Store) ValidarEstilo(id, arquivoAprovacao string, agora time.Time) (Estilo, error) {
	e, err := s.LerEstilo(id)
	if err != nil {
		return e, err
	}
	if len(e.Referencias) == 0 {
		return e, fmt.Errorf("%w: estilo %q sem referencias", ErrEstiloInvalido, id)
	}
	a, err := LerAprovacao(arquivoAprovacao, id)
	if err != nil {
		return e, err
	}
	v := Validacao{
		Estilo: id, Versao: e.Versao, ValidadoEm: agora.UTC().Format(time.RFC3339),
		ArquivoAprovacao: arquivoAprovacao, AprovadoPor: a.AprovadoPor, PersonagensTeste: a.PersonagensTeste,
	}
	if err := gravarJSON(s.caminhoValidacao(id), v); err != nil {
		return e, err
	}
	e.Status = "validado"
	if err := gravarJSON(s.caminhoEstilo(id), e); err != nil {
		return e, err
	}
	return e, nil
}

// LerValidacao lê validacao.json do estilo.
func (s Store) LerValidacao(id string) (Validacao, error) {
	var v Validacao
	b, err := os.ReadFile(s.caminhoValidacao(id))
	if err != nil {
		return v, fmt.Errorf("video-inventario: lendo validacao de %q: %w", id, err)
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return v, fmt.Errorf("video-inventario: decodificando validacao de %q: %w", id, err)
	}
	return v, nil
}

// VerificarEstilo confere o hash de todos os itens do estilo; devolve os erros por item.
func (s Store) VerificarEstilo(id string) ([]error, error) {
	if _, err := s.LerEstilo(id); err != nil {
		return nil, err
	}
	itens, err := s.lerItens(id, "")
	if err != nil {
		return nil, err
	}
	var falhas []error
	for _, it := range itens {
		if err := s.Verificar(it); err != nil {
			falhas = append(falhas, err)
		}
	}
	return falhas, nil
}
