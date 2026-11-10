// Package suggest calcula, de forma determinística e só a partir do inventário, o que cabe numa
// peça: gates, scripts, hooks (pré e pós) e peças vizinhas. Nenhuma sugestão vem de LLM e nada
// que não exista na máquina é sugerido.
package suggest

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"legoloop/internal/assets"
	"legoloop/internal/dag"
	"legoloop/internal/graph"
	"legoloop/internal/reputation"
)

const (
	// MaxItems limita cada lista de scripts e hooks sugeridos.
	MaxItems = 50
	// MaxNeighbors limita as peças de catálogo sugeridas como vizinhas.
	MaxNeighbors = 12
	// MaxMCPNeighbors limita os MCPs do inventário sugeridos como vizinhos.
	MaxMCPNeighbors = 8
	// MaxLoopNeighbors limita os loops do acervo sugeridos como vizinhos.
	MaxLoopNeighbors = 8
)

// Context descreve a peça para a qual se pedem sugestões. Block é o bloco tipado do catálogo
// (nil para peça genérica, cuja saída é "Qualquer"); Artifact é o nome do artefato de saída.
type Context struct {
	Block    *graph.Block
	Kind     string
	Artifact string
}

// Gate é uma expressão de gate sugerida, já instanciada com o artefato da peça.
type Gate struct {
	Expr  string `json:"expr"`
	Label string `json:"label"`
}

// Neighbor é uma peça sugerida para ligar depois desta.
type Neighbor struct {
	Kind      string                  `json:"kind"` // block | mcp
	ID        string                  `json:"id"`
	Name      string                  `json:"name"`
	Installed bool                    `json:"installed"`
	Rating    *reputation.PieceRating `json:"rating,omitempty"`
}

// Result agrupa as sugestões. Todas as listas são vazias (nunca nil) quando não há o que sugerir.
type Result struct {
	Gates        []Gate
	Scripts      []assets.Item
	HooksPre     []assets.Item
	HooksPost    []assets.Item
	SkippedHooks int
	Neighbors    []Neighbor
}

// Suggest monta as sugestões para a peça a partir do catálogo e do inventário.
func Suggest(ctx Context, idx *graph.Index, snap assets.Snapshot) Result {
	pre, post, skipped := hookSuggestions(snap.Of(assets.KindHook))
	return Result{
		Gates:        gateSuggestions(ctx),
		Scripts:      scriptSuggestions(snap.Of(assets.KindScript)),
		HooksPre:     pre,
		HooksPost:    post,
		SkippedHooks: skipped,
		Neighbors:    neighborSuggestions(ctx, idx, snap),
	}
}

// ---------- Gates ----------

func gateSuggestions(ctx Context) []Gate {
	a := ctx.Artifact
	if a == "" || strings.ContainsAny(a, "'\\\n\r") {
		return []Gate{}
	}
	cands := []Gate{{fmt.Sprintf("file_exists('%s')", a), "Arquivo existe"}}
	if strings.EqualFold(filepath.Ext(a), ".md") {
		cands = append(cands,
			Gate{fmt.Sprintf("word_count('%s') > 100", a), "Mais de 100 palavras"},
			Gate{fmt.Sprintf("line_count('%s') > 10", a), "Mais de 10 linhas"})
		if ctx.Block == nil || ctx.Block.Produces == graph.Spec {
			cands = append(cands,
				Gate{fmt.Sprintf("has_binary_criteria('%s')", a), "Tem critérios binários (checklist)"},
				Gate{fmt.Sprintf("has_acceptance_criteria('%s')", a), "Tem seção de critérios de aceitação"})
		}
	}
	own := ""
	if ctx.Block != nil {
		own = ctx.Block.Gate
	}
	out := []Gate{}
	for _, g := range cands {
		if own == "" || !strings.Contains(own, g.Expr) { // o gate próprio do bloco já cobre
			out = append(out, g)
		}
	}
	return out
}

// ---------- Scripts e hooks ----------

// runnable diz se o item vira um comando que o contrato aceita (sem shell, executável válido).
func runnable(it assets.Item) bool {
	if it.Kind == assets.KindHook && assets.ContainsSecret(it.Command) {
		return false // a credencial iria para o arquivo DAG
	}
	argv, err := assets.ArgvFor(it)
	return err == nil && dag.ValidateCommand(argv) == nil
}

func scriptSuggestions(items []assets.Item) []assets.Item {
	out := []assets.Item{}
	for _, it := range items {
		if runnable(it) {
			out = append(out, it)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if (a.Scope == "project") != (b.Scope == "project") {
			return a.Scope == "project"
		}
		if an, bn := strings.ToLower(a.Name), strings.ToLower(b.Name); an != bn {
			return an < bn
		}
		return a.Path < b.Path
	})
	return capItems(out)
}

// hookSlots diz em quais slots um evento de hook faz sentido: eventos "Pre*" e de início só
// antes da peça, eventos "Post*" e de fim só depois; os demais valem nos dois.
func hookSlots(event string) (pre, post bool) {
	switch event {
	case "PreToolUse", "SessionStart", "UserPromptSubmit", "PreCompact":
		return true, false
	case "PostToolUse", "Stop", "SubagentStop", "SessionEnd":
		return false, true
	}
	return true, true
}

func hookSuggestions(items []assets.Item) (pre, post []assets.Item, skipped int) {
	pre, post = []assets.Item{}, []assets.Item{}
	for _, it := range items {
		if !runnable(it) {
			skipped++
			continue
		}
		p, q := hookSlots(it.Meta["event"])
		if p {
			pre = append(pre, it)
		}
		if q {
			post = append(post, it)
		}
	}
	return capItems(sortByName(pre)), capItems(sortByName(post)), skipped
}

func sortByName(items []assets.Item) []assets.Item {
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if an, bn := strings.ToLower(a.Name), strings.ToLower(b.Name); an != bn {
			return an < bn
		}
		return a.Scope < b.Scope
	})
	return items
}

func capItems(items []assets.Item) []assets.Item {
	if len(items) > MaxItems {
		return items[:MaxItems]
	}
	return items
}

// ---------- Vizinhos ----------

func neighborSuggestions(ctx Context, idx *graph.Index, snap assets.Snapshot) []Neighbor {
	installed := map[string]bool{}
	for _, s := range snap.Of(assets.KindSkill) {
		installed[strings.ToLower(s.Name)] = true
	}
	var blocks []graph.Block
	if ctx.Block != nil {
		blocks = idx.Consumers(ctx.Block.Produces)
	} else {
		blocks = idx.Blocks()
	}
	typed := []Neighbor{}
	for _, b := range blocks {
		rating := reputation.ScorePiece(b.Ref, b.Label, b.Gate, 1, 1)
		typed = append(typed, Neighbor{
			Kind:      "block",
			ID:        b.ID,
			Name:      b.Label,
			Installed: installed[strings.ToLower(path.Base(b.Ref))],
			Rating:    &rating,
		})
	}
	sort.SliceStable(typed, func(i, j int) bool { return neighborLess(typed[i], typed[j]) })
	if len(typed) > MaxNeighbors {
		typed = typed[:MaxNeighbors]
	}
	out := append(typed, mcpNeighbors(snap.Of(assets.KindMCP))...)
	return append(out, loopNeighbors(snap.Of(assets.KindLoop))...)
}

// neighborLess ordena instalados primeiro, depois maior Score de Reputação, depois nome.
func neighborLess(a, b Neighbor) bool {
	if a.Installed != b.Installed {
		return a.Installed
	}
	if sa, sb := neighborScore(a), neighborScore(b); sa != sb {
		return sa > sb
	}
	return strings.ToLower(a.Name) < strings.ToLower(b.Name)
}

func neighborScore(n Neighbor) float64 {
	if n.Rating == nil {
		return 0.0
	}
	return n.Rating.TotalScore
}

func mcpNeighbors(mcps []assets.Item) []Neighbor {
	sorted := sortByName(append([]assets.Item(nil), mcps...))
	out := []Neighbor{}
	seen := map[string]bool{}
	for _, m := range sorted {
		if seen[m.Name] || len(out) >= MaxMCPNeighbors {
			continue
		}
		seen[m.Name] = true
		out = append(out, Neighbor{Kind: "mcp", ID: m.ID, Name: m.Name, Installed: true})
	}
	return out
}

func loopNeighbors(loops []assets.Item) []Neighbor {
	sorted := sortByName(append([]assets.Item(nil), loops...))
	out := []Neighbor{}
	seen := map[string]bool{}
	for _, l := range sorted {
		if seen[l.Name] || len(out) >= MaxLoopNeighbors {
			continue
		}
		seen[l.Name] = true
		out = append(out, Neighbor{Kind: "loop", ID: l.ID, Name: l.Name, Installed: true})
	}
	return out
}
