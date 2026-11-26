package inventario

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// elo é um estilo da cadeia de herança; herdado indica que só itens herdáveis valem.
type elo struct {
	id      string
	herdado bool
}

// cadeia segue Herda a partir do estilo pedido; ciclo retorna ErrHerancaCiclica.
func (s Store) cadeia(estilo string) ([]elo, error) {
	vistos := map[string]bool{}
	var out []elo
	for atual := estilo; atual != ""; {
		if vistos[atual] {
			return nil, fmt.Errorf("%w: %q volta a %q", ErrHerancaCiclica, estilo, atual)
		}
		vistos[atual] = true
		e, err := s.LerEstilo(atual)
		if err != nil {
			return nil, err
		}
		out = append(out, elo{id: atual, herdado: atual != estilo})
		atual = e.Herda
	}
	return out, nil
}

// lerItens lê os itens de um estilo; tipo vazio lê todos os tipos.
func (s Store) lerItens(estilo, tipo string) ([]Item, error) {
	if tipo == "" {
		tipo = "*"
	}
	padrao := filepath.Join(s.dirEstilos(), estilo, "itens", tipo, "*.json")
	arquivos, err := filepath.Glob(padrao)
	if err != nil {
		return nil, fmt.Errorf("video-inventario: listando %s: %w", padrao, err)
	}
	itens := make([]Item, 0, len(arquivos))
	for _, a := range arquivos {
		b, err := os.ReadFile(a)
		if err != nil {
			return nil, fmt.Errorf("video-inventario: lendo item %s: %w", a, err)
		}
		var it Item
		if err := json.Unmarshal(b, &it); err != nil {
			return nil, fmt.Errorf("video-inventario: decodificando item %s: %w", a, err)
		}
		itens = append(itens, it)
	}
	return itens, nil
}

// casa aceita item aprovado, não obsoleto, herdável quando herdado e com todas as tags.
func casa(it Item, herdado bool, tags []string) bool {
	if !it.Aprovado || it.Obsoleto || (herdado && !it.Herdavel) {
		return false
	}
	tem := make(map[string]bool, len(it.Tags))
	for _, t := range it.Tags {
		tem[t] = true
	}
	for _, t := range tags {
		if !tem[t] {
			return false
		}
	}
	return true
}

// Buscar devolve itens do estilo e herdáveis dos ancestrais que casam com tipo e tags, ordenados por id.
func (s Store) Buscar(estilo, tipo string, tags []string) ([]Item, error) {
	elos, err := s.cadeia(estilo)
	if err != nil {
		return nil, err
	}
	out := []Item{}
	for _, e := range elos {
		itens, err := s.lerItens(e.id, tipo)
		if err != nil {
			return nil, err
		}
		for _, it := range itens {
			if casa(it, e.herdado, tags) {
				out = append(out, it)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
