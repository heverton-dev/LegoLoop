package graph

import (
	"container/heap"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/gate"
)

var blockIDRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

// Index é o índice imutável de compatibilidade (porta de saída → portas de entrada).
type Index struct {
	blocks    []Block
	byID      map[string]Block
	consumers map[ArtifactType][]Block
}

// NewIndex valida o catálogo e indexa consumidores por tipo (ordenados por ID).
func NewIndex(catalog []Block) (*Index, error) {
	idx := &Index{byID: map[string]Block{}, consumers: map[ArtifactType][]Block{}}
	for _, b := range catalog {
		if err := validateBlock(b); err != nil {
			return nil, err
		}
		if _, dup := idx.byID[b.ID]; dup {
			return nil, fmt.Errorf("graph: peça duplicada %q", b.ID)
		}
		idx.byID[b.ID] = b
		idx.blocks = append(idx.blocks, b)
		for _, t := range b.Consumes {
			idx.consumers[t] = append(idx.consumers[t], b)
		}
	}
	sort.Slice(idx.blocks, func(i, j int) bool { return idx.blocks[i].ID < idx.blocks[j].ID })
	for t := range idx.consumers {
		sort.Slice(idx.consumers[t], func(i, j int) bool { return idx.consumers[t][i].ID < idx.consumers[t][j].ID })
	}
	return idx, nil
}

func validateBlock(b Block) error {
	if !blockIDRe.MatchString(b.ID) {
		return fmt.Errorf("graph: id de peça inválido %q", b.ID)
	}
	isTrigger := IsTriggerKind(b.Kind)
	if isTrigger && len(b.Consumes) > 0 {
		return fmt.Errorf("graph: gatilho %q não aceita portas de entrada", b.ID)
	}
	if b.Cost <= 0 || (len(b.Consumes) == 0 && !isTrigger) {
		return fmt.Errorf("graph: peça %q exige custo > 0 e ao menos uma entrada", b.ID)
	}
	for _, t := range append([]ArtifactType{b.Produces}, b.Consumes...) {
		if _, err := ParseArtifactType(string(t)); err != nil {
			return fmt.Errorf("graph: peça %q: %w", b.ID, err)
		}
	}
	if _, err := gate.Compile(b.Gate); err != nil {
		return fmt.Errorf("graph: peça %q: %w", b.ID, err)
	}
	return nil
}

// Blocks devolve o catálogo ordenado por ID.
func (idx *Index) Blocks() []Block { return slices.Clone(idx.blocks) }

func (idx *Index) resolveBlock(id string) (Block, bool) {
	if b, ok := idx.byID[id]; ok {
		return b, true
	}
	alias := "trigger-" + strings.ReplaceAll(string(dag.TriggerFimLegado), "_", "-")
	if id == alias {
		b, ok := idx.byID["trigger-processo-fim"]
		return b, ok
	}
	return Block{}, false
}

// Block busca uma peça por ID (incluindo aliases legados de gatilho).
func (idx *Index) Block(id string) (Block, bool) {
	return idx.resolveBlock(id)
}

// Consumers lista as peças que aceitam o tipo como entrada.
func (idx *Index) Consumers(t ArtifactType) []Block { return slices.Clone(idx.consumers[t]) }

func canConnectBlocks(a, b Block) bool {
	if IsTriggerKind(b.Kind) {
		return false
	}
	if IsTriggerKind(a.Kind) {
		return true
	}
	return slices.Contains(b.Consumes, a.Produces)
}

// CanConnect informa se a saída de `from` encaixa numa entrada de `to`.
func (idx *Index) CanConnect(from, to string) (bool, error) {
	a, ok := idx.resolveBlock(from)
	if !ok {
		return false, fmt.Errorf("graph: peça %q inexistente", from)
	}
	b, ok := idx.resolveBlock(to)
	if !ok {
		return false, fmt.Errorf("graph: peça %q inexistente", to)
	}
	return canConnectBlocks(a, b), nil
}

func (idx *Index) nonTriggerTargets() []string {
	targets := []string{}
	for _, b := range idx.blocks {
		if !IsTriggerKind(b.Kind) {
			targets = append(targets, b.ID)
		}
	}
	return targets
}

// Matrix devolve, para cada peça, os IDs das peças às quais ela pode se conectar.
func (idx *Index) Matrix() map[string][]string {
	m := make(map[string][]string, len(idx.blocks))
	for _, b := range idx.blocks {
		if IsTriggerKind(b.Kind) {
			m[b.ID] = idx.nonTriggerTargets()
			continue
		}
		targets := []string{}
		for _, c := range idx.consumers[b.Produces] {
			if !IsTriggerKind(c.Kind) {
				targets = append(targets, c.ID)
			}
		}
		m[b.ID] = targets
	}
	return m
}

// Route calcula (Dijkstra) a sequência de peças de menor custo que transforma `from` em `to`.
// Empates são resolvidos pelo ID da peça, garantindo resultado determinístico.
func (idx *Index) Route(from, to ArtifactType) ([]Block, int, error) {
	for _, t := range []ArtifactType{from, to} {
		if _, err := ParseArtifactType(string(t)); err != nil {
			return nil, 0, err
		}
	}
	dist := map[ArtifactType]int{from: 0}
	via := map[ArtifactType]edge{}
	pq := &typeQueue{{t: from}}
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(typeItem)
		if cur.t == to {
			return unwind(from, to, via), cur.d, nil
		}
		if cur.d > dist[cur.t] {
			continue
		}
		idx.relax(cur, dist, via, pq)
	}
	return nil, 0, errors.New("graph: nenhuma rota entre " + string(from) + " e " + string(to))
}

// edge registra qual peça produziu um tipo e a partir de qual tipo de entrada.
type edge struct {
	block Block
	from  ArtifactType
}

func (idx *Index) relax(cur typeItem, dist map[ArtifactType]int, via map[ArtifactType]edge, pq *typeQueue) {
	for _, b := range idx.consumers[cur.t] {
		nd := cur.d + b.Cost
		if old, seen := dist[b.Produces]; !seen || nd < old {
			dist[b.Produces] = nd
			via[b.Produces] = edge{block: b, from: cur.t}
			heap.Push(pq, typeItem{t: b.Produces, d: nd})
		}
	}
}

func unwind(from, to ArtifactType, via map[ArtifactType]edge) []Block {
	var path []Block
	for t := to; t != from; t = via[t].from {
		path = append(path, via[t].block)
	}
	slices.Reverse(path)
	return path
}

type typeItem struct {
	t ArtifactType
	d int
}

type typeQueue []typeItem

func (q typeQueue) Len() int { return len(q) }
func (q typeQueue) Less(i, j int) bool {
	if q[i].d != q[j].d {
		return q[i].d < q[j].d
	}
	return q[i].t < q[j].t
}
func (q typeQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *typeQueue) Push(x any)   { *q = append(*q, x.(typeItem)) }
func (q *typeQueue) Pop() any {
	old := *q
	it := old[len(old)-1]
	*q = old[:len(old)-1]
	return it
}
