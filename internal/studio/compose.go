package studio

import (
	"encoding/json"
	"fmt"

	"legoloop/internal/dag"
	"legoloop/internal/graph"
	"legoloop/internal/provider"
)

// maxCanvasNodes limita o tamanho de uma esteira vinda do canvas.
const maxCanvasNodes = 200

// Composition é o estado do canvas enviado pelo navegador.
type Composition struct {
	Name    string       `json:"name"`
	Trigger *dag.Trigger `json:"trigger,omitempty"`
	Nodes   []CompNode   `json:"nodes"`
	Edges   []CompEdge   `json:"edges"`
}

// CompNode é um nó posicionado no canvas: instancia um bloco do catálogo (Block, contrato tipado)
// OU uma peça do inventário (Piece, contrato "Qualquer"). Slots e Execution são opcionais.
type CompNode struct {
	ID        string            `json:"id"`
	Block     string            `json:"block,omitempty"`
	Piece     *CompPiece        `json:"piece,omitempty"`
	Routing   dag.RoutingConfig `json:"routing"`
	Slots     *CompSlots        `json:"slots,omitempty"`
	Execution *dag.Execution    `json:"execution,omitempty"`
}

// CompPiece é uma peça do inventário. Ref é o id opaco do item (o navegador nunca envia caminhos).
type CompPiece struct {
	Kind     string `json:"kind"`
	Ref      string `json:"ref"`
	Artifact string `json:"artifact"`
	Gate     string `json:"gate,omitempty"`
}

// CompSlots são os encaixes opcionais da peça. Vazio é uma escolha válida do usuário.
type CompSlots struct {
	Gates   []string  `json:"gates,omitempty"`
	Scripts []SlotCmd `json:"scripts,omitempty"`
	Hooks   struct {
		Pre  []SlotCmd `json:"pre,omitempty"`
		Post []SlotCmd `json:"post,omitempty"`
	} `json:"hooks,omitzero"`
}

// SlotCmd é um comando de slot: ou Ref (item do inventário, resolvido no servidor) ou Command
// (texto livre do usuário, separado em argv sem shell). Exatamente um dos dois.
type SlotCmd struct {
	Ref     string `json:"ref,omitempty"`
	Command string `json:"command,omitempty"`
}

// CompEdge liga a porta de saída de From à porta de entrada de To.
type CompEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Issue é um problema apontado no canvas (por nó, aresta ou global).
type Issue struct {
	Node    string `json:"node,omitempty"`
	Edge    string `json:"edge,omitempty"`
	Message string `json:"message"`
}

// ValidationResult é a resposta de /api/validate e /api/materialize.
type ValidationResult struct {
	OK     bool            `json:"ok"`
	Issues []Issue         `json:"issues"`
	DAG    json.RawMessage `json:"dag,omitempty"`
	Path   string          `json:"path,omitempty"`
}

// compose valida o canvas contra o grafo semântico e o inventário e compila o DAG.
func compose(idx *graph.Index, inv provider.Inventory, res Resolver, c Composition) (ValidationResult, *dag.DAG) {
	if len(c.Nodes) == 0 || len(c.Nodes) > maxCanvasNodes {
		return fail(Issue{Message: fmt.Sprintf("a esteira precisa de 1 a %d nós", maxCanvasNodes)}), nil
	}
	infos, slots, issues := resolveNodes(idx, inv, res, c.Nodes)
	issues = append(issues, checkEdges(idx, infos, c.Edges)...)
	if len(issues) > 0 {
		return fail(issues...), nil
	}
	raw, err := json.Marshal(buildDAG(c, infos, slots))
	if err != nil {
		return fail(Issue{Message: err.Error()}), nil
	}
	d, err := dag.Parse(raw)
	if err != nil {
		return fail(Issue{Message: err.Error()}), nil
	}
	out, _ := d.Marshal()
	return ValidationResult{OK: true, Issues: []Issue{}, DAG: out}, d
}

func fail(issues ...Issue) ValidationResult {
	return ValidationResult{OK: false, Issues: issues}
}

func resolveNodes(idx *graph.Index, inv provider.Inventory, res Resolver, nodes []CompNode) (map[string]nodeInfo, map[string]*dag.Slots, []Issue) {
	infos := map[string]nodeInfo{}
	slots := map[string]*dag.Slots{}
	artifacts := map[string]string{}
	var issues []Issue
	add := func(id, msg string) { issues = append(issues, Issue{Node: id, Message: msg}) }
	for _, n := range nodes {
		info, msg := resolveNode(idx, res, n)
		if msg != "" {
			add(n.ID, msg)
			continue
		}
		if other, dup := artifacts[info.Artifact]; dup {
			add(n.ID, fmt.Sprintf("handoff %s já é produzido por %s", info.Artifact, other))
		}
		artifacts[info.Artifact] = n.ID
		infos[n.ID] = info
		if msg := routingIssue(inv, n.Routing); msg != "" {
			add(n.ID, msg)
		}
		s, msg := resolveSlots(res, n.Slots)
		if msg != "" {
			add(n.ID, msg)
		}
		slots[n.ID] = s
	}
	return infos, slots, issues
}

// routingIssue aplica a Regra de Ouro ao canvas: só modelos verdes no inventário.
func routingIssue(inv provider.Inventory, r dag.RoutingConfig) string {
	candidates := []string{r.Model}
	switch r.Mode {
	case dag.ModeAuto:
		candidates = []string{r.Combo}
	case dag.ModeSuggested:
		candidates = append(candidates, r.Fallback)
	}
	for _, m := range candidates {
		if m != "" && inv.Has(m) {
			return ""
		}
	}
	return fmt.Sprintf("nenhum modelo de %v está disponível nesta máquina", candidates)
}

func checkEdges(idx *graph.Index, infos map[string]nodeInfo, edges []CompEdge) []Issue {
	var issues []Issue
	seen := map[CompEdge]bool{}
	for _, e := range edges {
		label := e.From + "→" + e.To
		from, okF := infos[e.From]
		to, okT := infos[e.To]
		switch {
		case seen[e]:
			issues = append(issues, Issue{Edge: label, Message: "conexão duplicada"})
		case !okF || !okT:
			issues = append(issues, Issue{Edge: label, Message: "conexão com nó inexistente"})
		default:
			if msg := edgeIssue(idx, from, to); msg != "" {
				issues = append(issues, Issue{Edge: label, Message: msg})
			}
		}
		seen[e] = true
	}
	return issues
}

// edgeIssue: peça genérica aceita e produz "Qualquer", então liga com tudo; entre blocos tipados
// vale a matriz de compatibilidade do grafo.
func edgeIssue(idx *graph.Index, from, to nodeInfo) string {
	if from.generic() || to.generic() {
		return ""
	}
	if ok, _ := idx.CanConnect(from.Block.ID, to.Block.ID); !ok {
		return fmt.Sprintf("%s produz %s, que %s não consome", from.Label, from.Block.Produces, to.Label)
	}
	return ""
}

func buildDAG(c Composition, infos map[string]nodeInfo, slots map[string]*dag.Slots) dag.DAG {
	incoming := map[string][]string{}
	for _, e := range c.Edges {
		incoming[e.To] = append(incoming[e.To], e.From)
	}
	d := dag.DAG{Schema: "https://legoloop.dev/schemas/dag.v1.json", Name: c.Name, Version: "2.0.0", Trigger: c.Trigger}
	if d.Trigger == nil {
		for _, n := range c.Nodes {
			if info, ok := infos[n.ID]; ok && info.Block != nil {
				switch info.Block.Kind {
				case "trigger_cron":
					d.Trigger = &dag.Trigger{Type: dag.TriggerCron, Cron: "0 * * * *"}
				case "trigger_file_change":
					d.Trigger = &dag.Trigger{Type: dag.TriggerFileChange, Path: "docs"}
				case "trigger_loop_end":
					d.Trigger = &dag.Trigger{Type: dag.TriggerLoopEnd, Loop: "upstream"}
				case "trigger_manual", "trigger_board":
					d.Trigger = &dag.Trigger{Type: dag.TriggerManual}
				}
				if d.Trigger != nil {
					break
				}
			}
		}
		if d.Trigger == nil {
			d.Trigger = &dag.Trigger{Type: dag.TriggerManual}
		}
	}
	for _, n := range c.Nodes {
		info := infos[n.ID]
		inputs := map[string]string{}
		for _, from := range incoming[n.ID] {
			inputs[from] = "file('" + infos[from].Artifact + "')"
		}
		if len(inputs) == 0 {
			if info.Block == nil || len(info.Block.Consumes) > 0 {
				inputs[info.rootInput()] = "var." + info.rootInput()
			}
		}
		d.Nodes = append(d.Nodes, dag.Node{
			ID: n.ID, Type: info.Type, Ref: info.Ref, Command: info.Command, DependsOn: incoming[n.ID],
			Routing: n.Routing, Inputs: inputs, Execution: n.Execution, Slots: slots[n.ID],
			Output: dag.HandoffContract{Artifact: info.Artifact, Gate: info.Gate},
		})
	}
	return d
}
