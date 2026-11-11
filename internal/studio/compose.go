package studio

import (
	"encoding/json"
	"fmt"

	"legoloop/internal/dag"
	"legoloop/internal/graph"
	"legoloop/internal/provider"
)

// maxCanvasNodes limita o tamanho de um Processo vindo do canvas.
const maxCanvasNodes = 200

// Composition é o estado do canvas enviado pelo navegador.
type Composition struct {
	Name    string       `json:"name"`
	Trigger *dag.Trigger `json:"trigger,omitempty"`
	Nodes   []CompNode   `json:"nodes"`
	Edges   []CompEdge   `json:"edges"`
}

// CompNode é um nó posicionado no canvas: Block (contrato tipado) ou Piece (contrato "Qualquer").
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

// Problema é uma falha apontada no canvas (por nó, aresta ou global).
type Problema struct {
	Node    string `json:"node,omitempty"`
	Edge    string `json:"edge,omitempty"`
	Message string `json:"message"`
}

// ValidationResult é a resposta de /api/validate e /api/materialize.
type ValidationResult struct {
	OK     bool            `json:"ok"`
	Issues []Problema      `json:"issues"`
	DAG    json.RawMessage `json:"dag,omitempty"`
	Path   string          `json:"path,omitempty"`
}

// compose valida o canvas contra o grafo semântico e o inventário e compila o DAG.
func compose(idx *graph.Index, inv provider.Inventory, res Resolver, c Composition) (ValidationResult, *dag.DAG) {
	if len(c.Nodes) == 0 || len(c.Nodes) > maxCanvasNodes {
		return fail(Problema{Message: fmt.Sprintf("o Processo precisa de 1 a %d nós", maxCanvasNodes)}), nil
	}
	infos, slots, problemas := resolveNodes(idx, inv, res, c.Nodes)
	problemas = append(problemas, checkEdges(idx, infos, c.Edges)...)
	if len(problemas) > 0 {
		return fail(problemas...), nil
	}
	raw, err := json.Marshal(buildDAG(c, infos, slots))
	if err != nil {
		return fail(Problema{Message: err.Error()}), nil
	}
	d, err := dag.Parse(raw)
	if err != nil {
		return fail(Problema{Message: err.Error()}), nil
	}
	out, _ := d.Marshal()
	return ValidationResult{OK: true, Issues: []Problema{}, DAG: out}, d
}

func fail(problemas ...Problema) ValidationResult {
	return ValidationResult{OK: false, Issues: problemas}
}

func resolveNodes(idx *graph.Index, inv provider.Inventory, res Resolver, nodes []CompNode) (map[string]nodeInfo, map[string]*dag.Slots, []Problema) {
	infos := map[string]nodeInfo{}
	slots := map[string]*dag.Slots{}
	artifacts := map[string]string{}
	var problemas []Problema
	add := func(id, msg string) { problemas = append(problemas, Problema{Node: id, Message: msg}) }
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
		if msg := routingProblem(inv, n.Routing); msg != "" {
			add(n.ID, msg)
		}
		s, msg := resolveSlots(res, n.Slots)
		if msg != "" {
			add(n.ID, msg)
		}
		slots[n.ID] = s
	}
	return infos, slots, problemas
}

// routingProblem aplica a Regra de Ouro ao canvas: só modelos verdes no inventário.
func routingProblem(inv provider.Inventory, r dag.RoutingConfig) string {
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

func checkEdges(idx *graph.Index, infos map[string]nodeInfo, edges []CompEdge) []Problema {
	var problemas []Problema
	seen := map[CompEdge]bool{}
	for _, e := range edges {
		label := e.From + "→" + e.To
		from, okF := infos[e.From]
		to, okT := infos[e.To]
		switch {
		case seen[e]:
			problemas = append(problemas, Problema{Edge: label, Message: "conexão duplicada"})
		case !okF || !okT:
			problemas = append(problemas, Problema{Edge: label, Message: "conexão com nó inexistente"})
		default:
			if msg := edgeProblem(idx, from, to); msg != "" {
				problemas = append(problemas, Problema{Edge: label, Message: msg})
			}
		}
		seen[e] = true
	}
	return problemas
}

// edgeProblem: peça genérica liga com tudo; entre peças tipadas vale a matriz de compatibilidade do grafo.
func edgeProblem(idx *graph.Index, from, to nodeInfo) string {
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
		d.Trigger = inferTrigger(c, infos)
	}
	for _, n := range c.Nodes {
		info := infos[n.ID]
		d.Nodes = append(d.Nodes, dag.Node{
			ID: n.ID, Type: info.Type, Ref: info.Ref, Command: info.Command, DependsOn: incoming[n.ID],
			Routing: n.Routing, Inputs: nodeInputs(n.ID, info, incoming[n.ID], infos), Execution: n.Execution, Slots: slots[n.ID],
			Output: dag.HandoffContract{Artifact: info.Artifact, Gate: info.Gate},
		})
	}
	return d
}

// inferTrigger deduz o gatilho do primeiro nó de gatilho da composição; padrão manual.
func inferTrigger(c Composition, infos map[string]nodeInfo) *dag.Trigger {
	for _, n := range c.Nodes {
		info, ok := infos[n.ID]
		if !ok || info.Block == nil {
			continue
		}
		if t := triggerForKind(info.Block.Kind); t != nil {
			return t
		}
	}
	return &dag.Trigger{Type: dag.TriggerManual}
}

func triggerForKind(kind string) *dag.Trigger {
	switch kind {
	case "trigger_cron":
		return &dag.Trigger{Type: dag.TriggerCron, Cron: "0 * * * *"}
	case "trigger_file_change":
		return &dag.Trigger{Type: dag.TriggerFileChange, Path: "docs"}
	case "trigger_loop_end":
		return &dag.Trigger{Type: dag.TriggerLoopEnd, Loop: "upstream"}
	case "trigger_manual", "trigger_board":
		return &dag.Trigger{Type: dag.TriggerManual}
	}
	return nil
}

// nodeInputs liga cada aresta de entrada ao artefato do nó de origem; sem arestas, usa a entrada raiz.
func nodeInputs(id string, info nodeInfo, from []string, infos map[string]nodeInfo) map[string]string {
	inputs := map[string]string{}
	for _, f := range from {
		inputs[f] = "file('" + infos[f].Artifact + "')"
	}
	if len(inputs) == 0 && (info.Block == nil || len(info.Block.Consumes) > 0) {
		inputs[info.rootInput()] = "var." + info.rootInput()
	}
	return inputs
}
