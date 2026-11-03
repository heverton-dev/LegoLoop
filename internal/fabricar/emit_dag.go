package fabricar

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"legoloop/internal/dag"
	"legoloop/internal/loop"
)

// ErrNoSteps: o processo não tem nenhum passo.
var ErrNoSteps = errors.New("fabricar: o processo não tem passos")

// Step é um passo do processo fabricado: uma skill que vira um nó do grafo.
type Step struct {
	Piece    string // nome da skill (sanitizado antes de virar ref)
	Artifact string // vazio = out/NN_<peça>.md
	Gate     string // vazio = file_exists('<artifact>')
}

// ProcessSpec é tudo que os emissores precisam. Eles são funções puras: nada é lido do disco nem executado.
type ProcessSpec struct {
	Name        string // legoloop-<slug>
	Description string
	Goal        string
	Model       string // modelo sugerido; vazio = sonnet
	Produces    string // tipo do grafo semântico; vazio = Artifact:Spec
	Steps       []Step
}

const (
	dagSchemaURL = "https://legoloop.dev/schemas/dag.v1.json"
	defaultModel = "sonnet"
)

var unsafeRef = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// refName deixa só [A-Za-z0-9._-]. O nome vem de uma peça varrida (dado não confiável): nunca vira comando.
func refName(piece string) string {
	s := strings.Trim(unsafeRef.ReplaceAllString(piece, "-"), "-.")
	if s == "" {
		return "peca"
	}
	if len(s) > 40 {
		s = strings.Trim(s[:40], "-.")
	}
	return s
}

func idPart(piece string) string { return strings.ReplaceAll(refName(piece), ".", "-") }

func (ps ProcessSpec) validate() error {
	if !strings.HasPrefix(ps.Name, namePrefix) || loop.ValidateID(ps.Name) != nil {
		return fmt.Errorf("%w: %q", ErrInvalidName, ps.Name)
	}
	if len(ps.Steps) == 0 {
		return ErrNoSteps
	}
	return nil
}

func (ps ProcessSpec) model() string {
	if ps.Model == "" {
		return defaultModel
	}
	return ps.Model
}

// artifactOf devolve o artefato do passo i (com o padrão out/NN_<peça>.md).
func (ps ProcessSpec) artifactOf(i int) string {
	if a := ps.Steps[i].Artifact; a != "" {
		return a
	}
	return fmt.Sprintf("out/%02d_%s.md", i+1, refName(ps.Steps[i].Piece))
}

func (ps ProcessSpec) gateOf(i int) string {
	if g := ps.Steps[i].Gate; g != "" {
		return g
	}
	return "file_exists('" + ps.artifactOf(i) + "')"
}

// nodes monta a cadeia sequencial: cada passo depende do anterior e consome o artefato dele.
func (ps ProcessSpec) nodes() []dag.Node {
	nodes := make([]dag.Node, len(ps.Steps))
	for i, s := range ps.Steps {
		n := dag.Node{
			ID:         fmt.Sprintf("n%02d_%s", i+1, idPart(s.Piece)),
			Type:       dag.TypeSkill,
			Ref:        "skills/" + refName(s.Piece),
			TimeoutSec: 600,
			Routing:    dag.RoutingConfig{Mode: dag.ModeSuggested, Model: ps.model(), Effort: dag.EffortMedium, Fallback: "code-fast"},
			Output:     dag.HandoffContract{Artifact: ps.artifactOf(i), Gate: ps.gateOf(i)},
		}
		if i > 0 {
			prev := nodes[i-1]
			n.DependsOn = []string{prev.ID}
			n.Inputs = map[string]string{prev.ID: "file('" + prev.Output.Artifact + "')"}
		}
		nodes[i] = n
	}
	return nodes
}

func (ps ProcessSpec) dag() dag.DAG {
	return dag.DAG{
		Schema: dagSchemaURL, ID: ps.Name, Name: ps.Name, Description: ps.Description, Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual}, Nodes: ps.nodes(),
	}
}

// EmitDag gera o <slug>.dag.json executável por `legoloop run <caminho>`. Autoverifica com dag.Parse (estrito).
func EmitDag(ps ProcessSpec) ([]byte, error) {
	if err := ps.validate(); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(ps.dag(), "", "  ")
	if err != nil {
		return nil, err
	}
	if _, err := dag.Parse(data); err != nil {
		return nil, fmt.Errorf("fabricar: o DAG gerado é inválido: %w", err)
	}
	return append(data, '\n'), nil
}
