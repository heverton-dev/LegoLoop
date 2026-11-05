package dag

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"legoloop/internal/gate"
)

// MaxDAGBytes limita o tamanho do contrato aceito (proteção contra payload abusivo).
const MaxDAGBytes = 1 << 20

// MaxTimeoutSec é o teto de timeout_sec por nó (24h).
const MaxTimeoutSec = 24 * 60 * 60

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)
	// shortNameRe detecta nomes curtos 8.3 do NTFS (GIT~1 ⇒ .git).
	shortNameRe = regexp.MustCompile(`~[0-9]`)
)

// sensitiveDirs nunca podem receber artefatos: gravar ali vira hook/config executável.
func sensitiveDirs() []string {
	return []string{".git", ".claude", ".vscode", ".github", ".idea", ".husky", ".gemini", ".cursor", ".agents"}
}

// CycleError indica deadlock: os nós listados formam (ou dependem de) um ciclo.
type CycleError struct{ Nodes []string }

func (e *CycleError) Error() string {
	return "dag: ciclo detectado entre " + strings.Join(e.Nodes, ", ")
}

// Parse decodifica de forma estrita (campos desconhecidos rejeitados) e valida o DAG.
func Parse(data []byte) (*DAG, error) {
	if len(data) > MaxDAGBytes {
		return nil, fmt.Errorf("dag: contrato excede %d bytes", MaxDAGBytes)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var d DAG
	if err := dec.Decode(&d); err != nil {
		return nil, fmt.Errorf("dag: json inválido: %w", err)
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return nil, errors.New("dag: conteúdo extra após o objeto JSON")
	}
	if err := d.validate(); err != nil {
		return nil, err
	}
	return &d, nil
}

// Marshal serializa o DAG indentado (formato gravado em pipeline.dag.json).
func (d *DAG) Marshal() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

func (d *DAG) validate() error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("dag: campo name obrigatório")
	}
	if len(d.Nodes) == 0 {
		return errors.New("dag: ao menos um nó é obrigatório")
	}
	if err := d.indexNodes(); err != nil {
		return err
	}
	for _, n := range d.Nodes {
		if err := d.checkDeps(n); err != nil {
			return err
		}
	}
	layers, err := d.topoLayers()
	if err != nil {
		return err
	}
	d.layers = layers
	if err := d.checkTrigger(); err != nil {
		return err
	}
	return d.checkExecution()
}

func (d *DAG) indexNodes() error {
	d.byID = make(map[string]int, len(d.Nodes))
	d.gates = make(map[string]*gate.Expr, len(d.Nodes))
	artifacts := make(map[string]string, len(d.Nodes))
	for i, n := range d.Nodes {
		if err := d.indexNode(i); err != nil {
			return err
		}
		key := normPath(n.Output.Artifact)
		if other, dup := artifacts[key]; dup {
			return fmt.Errorf("dag: nós %q e %q gravam o mesmo artefato %q", other, n.ID, n.Output.Artifact)
		}
		artifacts[key] = n.ID
	}
	return nil
}

func (d *DAG) indexNode(i int) error {
	n := d.Nodes[i]
	if !idRe.MatchString(n.ID) {
		return fmt.Errorf("dag: id inválido %q", n.ID)
	}
	if _, dup := d.byID[n.ID]; dup {
		return fmt.Errorf("dag: id duplicado %q", n.ID)
	}
	d.byID[n.ID] = i
	if err := validateNode(n); err != nil {
		return fmt.Errorf("dag: nó %q: %w", n.ID, err)
	}
	g, err := gate.Compile(n.Output.Gate)
	if err != nil {
		return fmt.Errorf("dag: nó %q: %w", n.ID, err)
	}
	d.gates[n.ID] = g
	return nil
}

func validateNode(n Node) error {
	switch n.Type {
	case TypeSkill:
		if strings.TrimSpace(n.Ref) == "" {
			return errors.New("skill exige ref")
		}
	case TypeScript:
		if err := validateCommand(n.Command); err != nil {
			return err
		}
	case TypeSubloop:
		if !idRe.MatchString(n.Ref) {
			return fmt.Errorf("subloop exige ref com o nome de um loop válido, recebeu %q", n.Ref)
		}
	case "mcp":
		return errors.New("nó MCP ainda não é executável: não há cliente MCP (fase 2)")
	default:
		return fmt.Errorf("tipo desconhecido %q", n.Type)
	}
	if n.TimeoutSec < 0 || n.TimeoutSec > MaxTimeoutSec {
		return fmt.Errorf("timeout_sec deve estar entre 0 e %d", MaxTimeoutSec)
	}
	if err := validateOutput(n); err != nil {
		return err
	}
	if err := validateExecution(n.Execution); err != nil {
		return err
	}
	if err := validateSlots(n.Slots); err != nil {
		return err
	}
	return validateRouting(n.Routing)
}

// ValidateCommand aplica ao argv a mesma regra dos nós script e dos slots (SEC-01): o executável
// é um nome no PATH ou um caminho absoluto.
func ValidateCommand(cmd []string) error { return validateCommand(cmd) }

// validateCommand aceita só nome puro (resolvido no PATH) ou caminho absoluto:
// um binário relativo ao repositório seria escolhido por quem escreveu o DAG.
func validateCommand(cmd []string) error {
	if len(cmd) == 0 || cmd[0] == "" {
		return errors.New("script exige command")
	}
	bin := cmd[0]
	if !filepath.IsAbs(bin) && strings.ContainsAny(bin, `/\:`) {
		return fmt.Errorf("command %q deve ser nome no PATH ou caminho absoluto", bin)
	}
	return nil
}

func validateArtifact(artifact string) error { return validatePath("artifact", artifact) }

// validateOutput confere os caminhos do contrato de saída: artifact, sub_artifact e schema.
func validateOutput(n Node) error {
	if err := validateArtifact(n.Output.Artifact); err != nil {
		return err
	}
	if err := validateSubArtifact(n); err != nil {
		return err
	}
	if n.Output.Schema != "" {
		return validatePath("schema", n.Output.Schema)
	}
	return nil
}

// validateSubArtifact: sub_artifact só existe no nó subloop e segue a mesma regra de caminho do artifact.
func validateSubArtifact(n Node) error {
	if n.Output.SubArtifact == "" {
		return nil
	}
	if n.Type != TypeSubloop {
		return fmt.Errorf("sub_artifact só vale no nó subloop (tipo %q)", n.Type)
	}
	return validatePath("sub_artifact", n.Output.SubArtifact)
}

// ValidateArtifact aplica ao nome de um artefato a mesma regra do contrato: caminho relativo ao
// projeto, sem ":", sem nome curto 8.3 e fora dos diretórios sensíveis.
func ValidateArtifact(artifact string) error { return validateArtifact(artifact) }

// validatePath exige caminho relativo ao projeto, sem ":" (letra de unidade e fluxos NTFS),
// sem nomes curtos 8.3 e fora dos diretórios sensíveis.
func validatePath(label, p string) error {
	native := filepath.FromSlash(p)
	if strings.Contains(p, ":") || !filepath.IsLocal(native) {
		return fmt.Errorf("%s %q deve ser relativo ao projeto", label, p)
	}
	for part := range strings.SplitSeq(filepath.ToSlash(native), "/") {
		if shortNameRe.MatchString(part) {
			return fmt.Errorf("%s %q usa nome curto 8.3", label, p)
		}
		for _, bad := range sensitiveDirs() {
			if strings.EqualFold(part, bad) {
				return fmt.Errorf("%s %q não pode ficar em %s", label, p, bad)
			}
		}
	}
	return nil
}

func validateRouting(r RoutingConfig) error {
	switch r.Mode {
	case ModeSuggested, ModeManual:
		if r.Model == "" {
			return fmt.Errorf("roteamento %s exige model", r.Mode)
		}
	case ModeAuto:
		if r.Combo == "" {
			return errors.New("roteamento auto exige combo")
		}
	default:
		return fmt.Errorf("modo de roteamento desconhecido %q", r.Mode)
	}
	switch r.Effort {
	case "", EffortLow, EffortMedium, EffortHigh:
	default:
		return fmt.Errorf("effort inválido %q (permitidos: low, medium, high)", r.Effort)
	}
	return nil
}

func (d *DAG) checkDeps(n Node) error {
	for _, dep := range n.DependsOn {
		if dep == n.ID {
			return fmt.Errorf("dag: nó %q depende de si mesmo", n.ID)
		}
		if _, ok := d.byID[dep]; !ok {
			return fmt.Errorf("dag: nó %q depende de %q inexistente", n.ID, dep)
		}
	}
	return nil
}

// topoLayers aplica Kahn por camadas; sobras com grau > 0 caracterizam ciclo.
func (d *DAG) topoLayers() ([][]string, error) {
	indeg := make(map[string]int, len(d.Nodes))
	children := make(map[string][]string)
	for _, n := range d.Nodes {
		indeg[n.ID] += 0
		for _, dep := range n.DependsOn {
			indeg[n.ID]++
			children[dep] = append(children[dep], n.ID)
		}
	}
	var layers [][]string
	for current := zeroDegree(indeg); len(current) > 0; current = nextLayer(current, indeg, children) {
		layers = append(layers, current)
	}
	if len(indeg) > 0 {
		return nil, &CycleError{Nodes: sortedKeys(indeg)}
	}
	return layers, nil
}

func zeroDegree(indeg map[string]int) []string {
	var out []string
	for id, deg := range indeg {
		if deg == 0 {
			out = append(out, id)
		}
	}
	for _, id := range out {
		delete(indeg, id)
	}
	sort.Strings(out)
	return out
}

func nextLayer(current []string, indeg map[string]int, children map[string][]string) []string {
	for _, id := range current {
		for _, c := range children[id] {
			indeg[c]--
		}
	}
	return zeroDegree(indeg)
}

func sortedKeys(m map[string]int) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
