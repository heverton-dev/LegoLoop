package matriz

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"legoloop/internal/assets"
	"legoloop/internal/ste"
)

func init() {
	ste.Default.Register("matriz.ausente", "A matriz de peças não existe. Rode legoloop matriz gerar.")
	ste.Default.Register("matriz.formato_antigo", "A matriz de peças está no formato antigo e não tem perfil de uso. Rode legoloop matriz gerar.")
	ste.Default.Register("matriz.defasada", "A matriz de peças está desatualizada: %s. Rode legoloop matriz gerar.")
}

// Stale diz o quanto a matriz gravada está atrás das peças que existem hoje.
type Stale struct {
	Ausente       bool     `json:"ausente"`        // não há registro
	FormatoAntigo bool     `json:"formato_antigo"` // registro v1: sem hash nem perfil de uso
	Novas         []string `json:"novas"`          // Ref de peças que a matriz não conhece
	Alteradas     []string `json:"alteradas"`      // Ref de peças cujo hash mudou
	Removidas     []string `json:"removidas"`      // Ref de peças que a matriz tem e o disco não
	Retentar      []string `json:"retentar"`       // Ref de peças que falharam ou esperam o cérebro
}

// Defasada diz se rodar `matriz gerar` mudaria alguma coisa.
func (s Stale) Defasada() bool {
	return s.Ausente || s.FormatoAntigo || len(s.Novas)+len(s.Alteradas)+len(s.Removidas)+len(s.Retentar) > 0
}

// Staleness compara a matriz com as peças atuais sem reprocessar conteúdo: usa o sha256 que a varredura
// já trouxe e só lê o arquivo de script, gate e loop (pequeno) para o hash. Não chama nota nem cérebro.
// reg == nil significa que não há matriz.
func Staleness(ctx context.Context, reg *Registro, items []assets.Item, d Deps) (Stale, error) {
	if reg == nil {
		return Stale{Ausente: true}, nil
	}
	if len(reg.Entradas) == 0 && len(reg.Ratings) > 0 {
		return Stale{FormatoAntigo: true}, nil
	}
	conhecidas := make(map[string]Entrada, len(reg.Entradas))
	for _, e := range reg.Entradas {
		conhecidas[e.Ref] = e
	}
	var s Stale
	for _, it := range filtrar(items, d) {
		if err := ctx.Err(); err != nil {
			return Stale{}, err
		}
		r, hash, _, _ := identidade(it, d)
		prev, ok := conhecidas[r]
		delete(conhecidas, r)
		switch {
		case !ok:
			s.Novas = append(s.Novas, r)
		case prev.Hash != hash:
			s.Alteradas = append(s.Alteradas, r)
		case prev.Retentar:
			s.Retentar = append(s.Retentar, r)
		}
	}
	for r := range conhecidas {
		s.Removidas = append(s.Removidas, r)
	}
	slices.Sort(s.Novas)
	slices.Sort(s.Alteradas)
	slices.Sort(s.Removidas)
	slices.Sort(s.Retentar)
	return s, nil
}

// Mensagem é o aviso em STE-PT para o usuário; vazio quando a matriz está em dia.
func (s Stale) Mensagem() string {
	switch {
	case s.Ausente:
		return ste.Msg("matriz.ausente")
	case s.FormatoAntigo:
		return ste.Msg("matriz.formato_antigo")
	case !s.Defasada():
		return ""
	}
	var partes []string
	if n := len(s.Novas); n > 0 {
		partes = append(partes, plural(n, "nova", "novas"))
	}
	if n := len(s.Alteradas); n > 0 {
		partes = append(partes, plural(n, "alterada", "alteradas"))
	}
	if n := len(s.Removidas); n > 0 {
		partes = append(partes, plural(n, "removida", "removidas"))
	}
	if n := len(s.Retentar); n > 0 {
		partes = append(partes, fmt.Sprintf("%d esperando o cérebro", n))
	}
	return ste.Msg("matriz.defasada", strings.Join(partes, ", "))
}

func plural(n int, um, varios string) string {
	if n == 1 {
		return "1 " + um
	}
	return fmt.Sprintf("%d %s", n, varios)
}

// Previa diz quantos conteúdos distintos `matriz gerar` recatalogaria agora: todos, se não há matriz ou ela é
// da v1; senão os das peças novas, alteradas ou à espera do cérebro. Cópias idênticas contam uma vez só (o
// perfil é compartilhado) e remover uma peça não dá trabalho de catalogação.
func Previa(ctx context.Context, reg *Registro, items []assets.Item, d Deps) (int, error) {
	s, err := Staleness(ctx, reg, items, d)
	if err != nil {
		return 0, err
	}
	precisa := map[string]bool{}
	todas := s.Ausente || s.FormatoAntigo
	for _, r := range append(append(slices.Clone(s.Novas), s.Alteradas...), s.Retentar...) {
		precisa[r] = true
	}
	distintos := map[string]bool{}
	for _, it := range filtrar(items, d) {
		r, hash, _, _ := identidade(it, d)
		if todas || precisa[r] {
			distintos[chaveHash(string(it.Kind), hash)] = true
		}
	}
	return len(distintos), nil
}
