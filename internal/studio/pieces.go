package studio

import (
	"fmt"
	"path/filepath"
	"strings"

	"legoloop/internal/assets"
	"legoloop/internal/dag"
	"legoloop/internal/gate"
	"legoloop/internal/graph"
)

// maxArtifactLen limita o nome do artefato de saída de uma peça genérica.
const maxArtifactLen = 128

// Resolver encontra itens do inventário pelo id opaco. A resolução acontece só no servidor.
type Resolver interface {
	ByID(id string) (assets.Item, bool)
}

// noInventory é o Resolver de um servidor sem inventário: nenhuma peça nem referência resolve.
type noInventory struct{}

func (noInventory) ByID(string) (assets.Item, bool) { return assets.Item{}, false }

// nodeInfo é o que o compose precisa saber de um nó, venha de um bloco do catálogo ou de uma
// peça do inventário.
type nodeInfo struct {
	Label    string
	Ref      string
	Type     dag.NodeType
	Command  []string
	Artifact string
	Gate     string
	Block    *graph.Block // nil numa peça genérica
}

func (n nodeInfo) generic() bool { return n.Block == nil }

// rootInput nomeia a variável de entrada de um nó sem predecessores.
func (n nodeInfo) rootInput() string {
	if n.Block == nil || len(n.Block.Consumes) == 0 {
		return "input"
	}
	switch n.Block.Consumes[0] {
	case graph.Idea:
		return "idea"
	case graph.URL:
		return "url"
	}
	return "input"
}

func resolveNode(idx *graph.Index, res Resolver, n CompNode) (nodeInfo, string) {
	switch {
	case n.Block != "" && n.Piece != nil:
		return nodeInfo{}, "informe um bloco do catálogo OU uma peça do inventário, não os dois"
	case n.Block != "":
		b, ok := idx.Block(n.Block)
		if !ok {
			return nodeInfo{}, fmt.Sprintf("bloco %q não existe no catálogo", n.Block)
		}
		return nodeInfo{Label: b.Label, Ref: b.Ref, Type: dag.TypeSkill, Artifact: b.Artifact, Gate: b.Gate, Block: &b}, ""
	case n.Piece != nil:
		return resolvePiece(res, *n.Piece)
	}
	return nodeInfo{}, "o nó precisa de um bloco do catálogo ou de uma peça do inventário"
}

func resolvePiece(res Resolver, p CompPiece) (nodeInfo, string) {
	info, msg := pieceByKind(res, p)
	if msg != "" {
		return nodeInfo{}, msg
	}
	art := p.Artifact
	if art == "" && p.Kind == "subloop" {
		if it, ok := res.ByID(p.Ref); ok && it.Meta != nil && it.Meta["artifact"] != "" {
			art = it.Meta["artifact"]
		}
	}
	if msg := artifactIssue(art); msg != "" {
		return nodeInfo{}, msg
	}
	gateExpr := p.Gate
	if gateExpr == "" {
		gateExpr = fmt.Sprintf("file_exists('%s')", art)
	}
	if _, err := gate.Compile(gateExpr); err != nil {
		return nodeInfo{}, fmt.Sprintf("gate da peça inválido: %v", err)
	}
	info.Artifact, info.Gate = art, gateExpr
	return info, ""
}

// pieceByKind resolve o item do inventário no tipo de nó correspondente.
func pieceByKind(res Resolver, p CompPiece) (nodeInfo, string) {
	switch p.Kind {
	case "skill":
		return skillPiece(res, p.Ref)
	case "script":
		return scriptPiece(res, p.Ref)
	case "mcp":
		return nodeInfo{}, "peça MCP ainda não é executável: não há cliente MCP (fase 2)"
	case "subloop":
		return subloopPiece(res, p.Ref)
	}
	return nodeInfo{}, fmt.Sprintf("tipo de peça desconhecido %q", p.Kind)
}

func subloopPiece(res Resolver, ref string) (nodeInfo, string) {
	it, msg := lookup(res, ref, assets.KindLoop)
	if msg != "" {
		return nodeInfo{}, fmt.Sprintf("sub-loop %q não encontrado no acervo (%s)", ref, msg)
	}
	refName := it.Name
	if it.Path != "" {
		refName = strings.TrimSuffix(filepath.Base(it.Path), ".loop.json")
	}
	return nodeInfo{Label: it.Name, Type: dag.TypeSubloop, Ref: refName}, ""
}

func skillPiece(res Resolver, ref string) (nodeInfo, string) {
	it, msg := lookup(res, ref, assets.KindSkill)
	if msg != "" {
		return nodeInfo{}, msg
	}
	// o orchestrix resolve <cofre>/<ref>/SKILL.md: vale o nome da pasta, não o do frontmatter
	return nodeInfo{Label: it.Name, Type: dag.TypeSkill, Ref: "skills/" + filepath.Base(filepath.Dir(it.Path))}, ""
}

func scriptPiece(res Resolver, ref string) (nodeInfo, string) {
	it, msg := lookup(res, ref, assets.KindScript)
	if msg != "" {
		return nodeInfo{}, msg
	}
	argv, err := assets.ArgvFor(it)
	if err != nil {
		return nodeInfo{}, fmt.Sprintf("script %q: %v", it.Name, err)
	}
	return nodeInfo{Label: it.Name, Type: dag.TypeScript, Command: argv}, ""
}

// artifactIssue valida o nome do artefato: precisa existir e não pode quebrar o gate padrão.
func artifactIssue(a string) string {
	switch {
	case a == "":
		return "a peça genérica precisa de um nome de artefato de saída"
	case len(a) > maxArtifactLen:
		return fmt.Sprintf("o nome do artefato excede %d caracteres", maxArtifactLen)
	case strings.ContainsAny(a, "'\\\n\r"):
		return "o nome do artefato não pode conter aspas simples, barra invertida ou quebra de linha"
	}
	if err := dag.ValidateArtifact(a); err != nil {
		return "artefato inválido: " + err.Error()
	}
	return ""
}

func lookup(res Resolver, id string, want assets.Kind) (assets.Item, string) {
	it, ok := res.ByID(id)
	if !ok {
		return assets.Item{}, fmt.Sprintf("o item %q não foi encontrado no inventário", id)
	}
	if it.Kind != want {
		return assets.Item{}, fmt.Sprintf("o item %q é do tipo %s, esperado %s", it.Name, it.Kind, want)
	}
	return it, ""
}

// resolveSlots converte os slots do canvas no contrato do DAG. Slot ausente ou vazio é válido.
func resolveSlots(res Resolver, s *CompSlots) (*dag.Slots, string) {
	if s == nil {
		return nil, ""
	}
	out := &dag.Slots{Gates: s.Gates}
	groups := []struct {
		name string
		kind assets.Kind
		in   []SlotCmd
		dst  *[][]string
	}{
		{"scripts", assets.KindScript, s.Scripts, &out.Scripts},
		{"hooks.pre", assets.KindHook, s.Hooks.Pre, &out.Hooks.Pre},
		{"hooks.post", assets.KindHook, s.Hooks.Post, &out.Hooks.Post},
	}
	for _, g := range groups {
		for _, c := range g.in {
			argv, msg := slotArgv(res, c, g.kind)
			if msg != "" {
				return out, fmt.Sprintf("slot %s: %s", g.name, msg)
			}
			*g.dst = append(*g.dst, argv)
		}
	}
	return out, ""
}

func slotArgv(res Resolver, c SlotCmd, kind assets.Kind) ([]string, string) {
	switch {
	case c.Ref != "" && c.Command != "":
		return nil, "informe ref OU command, não os dois"
	case c.Ref != "":
		return refArgv(res, c.Ref, kind)
	case c.Command != "":
		argv, err := assets.SplitCommand(c.Command)
		return argv, errMsg(err)
	}
	return nil, "o item não tem ref nem command"
}

func refArgv(res Resolver, ref string, kind assets.Kind) ([]string, string) {
	it, msg := lookup(res, ref, kind)
	if msg != "" {
		return nil, msg
	}
	if kind == assets.KindHook && assets.ContainsSecret(it.Command) {
		return nil, fmt.Sprintf("o hook %q tem credencial no comando e ela iria para o pipeline.dag.json; use um comando próprio sem segredo", it.Name)
	}
	argv, err := assets.ArgvFor(it)
	return argv, errMsg(err)
}

func errMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
