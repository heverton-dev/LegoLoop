package inventario

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/escrita"
)

// Erros de validação de entrada do Store.
var (
	ErrEstiloInvalido = errors.New("video-inventario: estilo invalido")
	ErrItemInvalido   = errors.New("video-inventario: item invalido")
)

// Store grava estilos e itens sob Raiz/inventario/video/estilos.
type Store struct {
	Raiz string
}

func (s Store) dirEstilos() string {
	return filepath.Join(s.Raiz, "inventario", "video", "estilos")
}

func (s Store) caminhoEstilo(id string) string {
	return filepath.Join(s.dirEstilos(), id, "estilo.json")
}

func (s Store) caminhoItem(it Item) string {
	return filepath.Join(s.dirEstilos(), it.Estilo, "itens", it.Tipo, it.ID+".json")
}

// kebab aceita minúsculas e dígitos em blocos separados por um único hífen.
func kebab(id string) bool {
	if id == "" || id[0] == '-' || id[len(id)-1] == '-' {
		return false
	}
	for i := 0; i < len(id); i++ {
		c := id[i]
		switch {
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9':
		case c == '-' && id[i-1] != '-':
		default:
			return false
		}
	}
	return true
}

func validarEstiloNovo(e Estilo) error {
	switch {
	case !kebab(e.ID):
		return fmt.Errorf("%w: id %q nao e kebab-case", ErrEstiloInvalido, e.ID)
	case e.Renderizador == "":
		return fmt.Errorf("%w: renderizador vazio", ErrEstiloInvalido)
	case e.Status != "rascunho":
		return fmt.Errorf("%w: status %q, esperado rascunho", ErrEstiloInvalido, e.Status)
	case e.PromptMestre == "":
		return fmt.Errorf("%w: prompt_mestre vazio", ErrEstiloInvalido)
	case len(e.Paleta) == 0:
		return fmt.Errorf("%w: paleta vazia", ErrEstiloInvalido)
	case len(e.Referencias) == 0:
		return fmt.Errorf("%w: sem referencias", ErrEstiloInvalido)
	}
	for i, r := range e.Referencias {
		if r.Origem == "" {
			return fmt.Errorf("%w: referencia %d sem origem", ErrEstiloInvalido, i)
		}
	}
	return nil
}

// NovoEstilo valida e grava estilo.json de forma atômica; id repetido retorna ErrEstiloDuplicado.
func (s Store) NovoEstilo(e Estilo) error {
	if err := validarEstiloNovo(e); err != nil {
		return err
	}
	caminho := s.caminhoEstilo(e.ID)
	if _, err := os.Stat(caminho); err == nil {
		return fmt.Errorf("%w: %q", ErrEstiloDuplicado, e.ID)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("video-inventario: consultando %s: %w", caminho, err)
	}
	return gravarJSON(caminho, e)
}

// LerEstilo lê estilo.json; ausência retorna ErrEstiloInexistente.
func (s Store) LerEstilo(id string) (Estilo, error) {
	var e Estilo
	if !kebab(id) {
		return e, fmt.Errorf("%w: %q", ErrEstiloInexistente, id)
	}
	b, err := os.ReadFile(s.caminhoEstilo(id))
	if os.IsNotExist(err) {
		return e, fmt.Errorf("%w: %q", ErrEstiloInexistente, id)
	}
	if err != nil {
		return e, fmt.Errorf("video-inventario: lendo estilo %q: %w", id, err)
	}
	if err := json.Unmarshal(b, &e); err != nil {
		return e, fmt.Errorf("video-inventario: decodificando estilo %q: %w", id, err)
	}
	return e, nil
}

// GravarItem grava o item no estilo; estilo em rascunho recusa item aprovado.
func (s Store) GravarItem(it Item) error {
	if err := ValidarTipo(it.Tipo, false); err != nil {
		return err
	}
	if !kebab(it.ID) {
		return fmt.Errorf("%w: id %q nao e kebab-case", ErrItemInvalido, it.ID)
	}
	e, err := s.LerEstilo(it.Estilo)
	if err != nil {
		return err
	}
	if it.Aprovado && e.Status != "validado" {
		return fmt.Errorf("%w: estilo %q status %q", ErrEstiloNaoValidado, e.ID, e.Status)
	}
	return gravarJSON(s.caminhoItem(it), it)
}

func gravarJSON(caminho string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("video-inventario: codificando %s: %w", caminho, err)
	}
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		return fmt.Errorf("video-inventario: criando pasta de %s: %w", caminho, err)
	}
	return escrita.GravarAtomico(caminho, append(b, '\n'))
}
