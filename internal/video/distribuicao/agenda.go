package distribuicao

import (
	"fmt"
	"slices"
	"sort"
	"time"
)

// Acoes do plano de agenda.
const (
	AcaoAgendar = "agendar"
	AcaoExpirar = "expirar"
	AcaoPular   = "pular"
)

// ConfigAgenda define horarios (UTC), tetos e janelas da agenda.
type ConfigAgenda struct {
	Horarios            []string      `json:"horarios"`
	TetoDiario          int           `json:"teto_diario"`
	CadenciaPerene      time.Duration `json:"cadencia_perene"`
	EstoqueMinimoPerene int           `json:"estoque_minimo_perene"`
	JanelaQuente        time.Duration `json:"janela_quente"`
	LimiteQuente        time.Duration `json:"limite_quente"`
	DiasAFrente         int           `json:"dias_a_frente"`
}

// Decisao e o que o plano manda fazer com um item.
type Decisao struct {
	ID     string    `json:"id"`
	Acao   string    `json:"acao"`
	Quando time.Time `json:"quando,omitempty"`
	Motivo string    `json:"motivo,omitempty"`
}

// Plano reune as decisoes e os alertas de estoque.
type Plano struct {
	Decisoes []Decisao `json:"decisoes"`
	Alertas  []string  `json:"alertas,omitempty"`
}

func (c ConfigAgenda) comPadroes() ConfigAgenda {
	if len(c.Horarios) == 0 {
		c.Horarios = []string{"14:00", "19:00"}
	}
	padroes := []struct {
		alvo *int
		val  int
	}{{&c.TetoDiario, 2}, {&c.EstoqueMinimoPerene, 3}, {&c.DiasAFrente, 14}}
	for _, p := range padroes {
		if *p.alvo <= 0 {
			*p.alvo = p.val
		}
	}
	for _, d := range []struct {
		alvo *time.Duration
		val  time.Duration
	}{{&c.CadenciaPerene, 24 * time.Hour}, {&c.JanelaQuente, 24 * time.Hour}, {&c.LimiteQuente, 48 * time.Hour}} {
		if *d.alvo <= 0 {
			*d.alvo = d.val
		}
	}
	return c
}

// horarios gera os horarios livres futuros em ordem crescente.
func (c ConfigAgenda) horarios(agora time.Time) ([]time.Time, error) {
	var out []time.Time
	inicio := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, time.UTC)
	for d := 0; d <= c.DiasAFrente; d++ {
		for _, h := range c.Horarios {
			t, err := time.Parse("15:04", h)
			if err != nil {
				return nil, fmt.Errorf("%w: horario %q", ErrAcervoInvalido, h)
			}
			slot := inicio.AddDate(0, 0, d).Add(time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute)
			if slot.After(agora) {
				out = append(out, slot)
			}
		}
	}
	slices.SortFunc(out, func(a, b time.Time) int { return a.Compare(b) })
	return out, nil
}

func dia(t time.Time) string { return t.UTC().Format("2006-01-02") }

// estadoAgenda guarda a ocupacao por dia e o ultimo perene agendado.
type estadoAgenda struct {
	pordia       map[string]int
	ultimoPerene time.Time
	cfg          ConfigAgenda
	slots        []time.Time
	agora        time.Time
	decisoes     []Decisao
}

func novoEstado(a Acervo, cfg ConfigAgenda, agora time.Time) (*estadoAgenda, error) {
	slots, err := cfg.horarios(agora)
	if err != nil {
		return nil, err
	}
	e := &estadoAgenda{pordia: map[string]int{}, cfg: cfg, slots: slots, agora: agora}
	for _, it := range a.Itens {
		if it.Status != StatusAgendado || it.Agendado == nil {
			continue
		}
		e.pordia[dia(*it.Agendado)]++
		if it.Tipo == TipoPerene && it.Agendado.After(e.ultimoPerene) {
			e.ultimoPerene = *it.Agendado
		}
	}
	return e, nil
}

// proximo acha o primeiro horario livre que respeita teto diario e cadencia.
func (e *estadoAgenda) proximo(it Item) (time.Time, bool) {
	for _, s := range e.slots {
		if e.pordia[dia(s)] >= e.cfg.TetoDiario {
			continue
		}
		if it.Tipo == TipoPerene && !e.ultimoPerene.IsZero() && s.Before(e.ultimoPerene.Add(e.cfg.CadenciaPerene)) {
			continue
		}
		return s, true
	}
	return time.Time{}, false
}

func (e *estadoAgenda) decidir(it Item) {
	s, ok := e.proximo(it)
	switch {
	case !ok:
		e.decisoes = append(e.decisoes, Decisao{ID: it.ID, Acao: AcaoPular, Motivo: "sem horario livre na janela"})
	case it.Tipo == TipoQuente && s.After(e.agora.Add(e.cfg.LimiteQuente)):
		e.decisoes = append(e.decisoes, Decisao{ID: it.ID, Acao: AcaoExpirar, Motivo: "quente fora da janela maxima de 48 h"})
	case !it.Validade.IsZero() && s.After(it.Validade):
		e.decisoes = append(e.decisoes, Decisao{ID: it.ID, Acao: AcaoExpirar, Motivo: "horario livre depois da validade"})
	default:
		motivo := ""
		if it.Tipo == TipoQuente && s.After(e.agora.Add(e.cfg.JanelaQuente)) {
			motivo = "quente acima de 24 h (limite de 48 h)"
		}
		e.pordia[dia(s)]++
		if it.Tipo == TipoPerene && s.After(e.ultimoPerene) {
			e.ultimoPerene = s
		}
		e.decisoes = append(e.decisoes, Decisao{ID: it.ID, Acao: AcaoAgendar, Quando: s, Motivo: motivo})
	}
}

// ordenar poe quentes antes de perenes; quentes por validade, depois prioridade e id.
func ordenar(itens []Item) {
	sort.SliceStable(itens, func(i, j int) bool {
		a, b := itens[i], itens[j]
		if a.Tipo != b.Tipo {
			return a.Tipo == TipoQuente
		}
		if a.Tipo == TipoQuente && !a.Validade.Equal(b.Validade) {
			return a.Validade.Before(b.Validade)
		}
		if a.Prioridade != b.Prioridade {
			return a.Prioridade > b.Prioridade
		}
		return a.ID < b.ID
	})
}

// Planejar decide expirar, agendar ou pular cada item; nao altera o acervo.
func Planejar(a Acervo, cfg ConfigAgenda, agora time.Time) (Plano, error) {
	cfg = cfg.comPadroes()
	if err := a.Validar(); err != nil {
		return Plano{}, err
	}
	e, err := novoEstado(a, cfg, agora)
	if err != nil {
		return Plano{}, err
	}
	var fila []Item
	for _, it := range a.Itens {
		switch {
		case it.Status == StatusPublic || it.Status == StatusExpirado:
			continue
		case (it.Status == StatusPronto || it.Status == StatusAgendado) && !it.Validade.IsZero() && it.Validade.Before(agora):
			e.decisoes = append(e.decisoes, Decisao{ID: it.ID, Acao: AcaoExpirar, Motivo: "validade vencida"})
		case it.Status == StatusPronto && !it.Autorizado:
			e.decisoes = append(e.decisoes, Decisao{ID: it.ID, Acao: AcaoPular, Motivo: "sem autorizacao de publicacao"})
		case it.Status == StatusPronto:
			fila = append(fila, it)
		}
	}
	ordenar(fila)
	for _, it := range fila {
		e.decidir(it)
	}
	return Plano{Decisoes: e.decisoes, Alertas: alertasEstoque(a, e.decisoes, cfg)}, nil
}

// alertasEstoque avisa quando o estoque de perenes cai abaixo do minimo.
func alertasEstoque(a Acervo, ds []Decisao, cfg ConfigAgenda) []string {
	expirados := map[string]bool{}
	for _, d := range ds {
		if d.Acao == AcaoExpirar {
			expirados[d.ID] = true
		}
	}
	n := 0
	for _, it := range a.Itens {
		if it.Tipo == TipoPerene && (it.Status == StatusPronto || it.Status == StatusAgendado) && !expirados[it.ID] {
			n++
		}
	}
	if n < cfg.EstoqueMinimoPerene {
		return []string{fmt.Sprintf("estoque de perenes baixo: %d de %d", n, cfg.EstoqueMinimoPerene)}
	}
	return nil
}

// Aplicar grava no acervo as decisoes do plano.
func (a *Acervo) Aplicar(p Plano) error {
	for _, d := range p.Decisoes {
		it := a.Buscar(d.ID)
		if it == nil {
			return fmt.Errorf("%w: item %q inexistente", ErrAcervoInvalido, d.ID)
		}
		switch d.Acao {
		case AcaoAgendar:
			q := d.Quando
			it.Status, it.Agendado = StatusAgendado, &q
		case AcaoExpirar:
			it.Status, it.Agendado = StatusExpirado, nil
		}
	}
	return nil
}

// GerarReels cria o item do Instagram para todo curto publico do YouTube e devolve quantos criou.
func (a *Acervo) GerarReels() int {
	n := 0
	for _, it := range slices.Clone(a.Itens) {
		id := it.ID + "-reel"
		curto := slices.Contains(it.Formatos, "9:16") && slices.Contains(it.Canais, "youtube")
		vivo := it.Status == StatusPronto || it.Status == StatusAgendado
		if !curto || !vivo || it.Visibilidade != VisPublico || a.Buscar(id) != nil {
			continue
		}
		a.Itens = append(a.Itens, Item{ID: id, Tema: it.Tema, Tipo: it.Tipo, Formatos: []string{"9:16"}, Canais: []string{"instagram"},
			Validade: it.Validade, Prioridade: it.Prioridade, Status: StatusPronto, Autorizado: it.Autorizado, Visibilidade: VisPublico, Origem: it.ID})
		n++
	}
	return n
}
