package processo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/loop"
	"legoloop/internal/plangen"
)

const (
	dagSchemaURL   = "https://legoloop.dev/schemas/dag.v1.json"
	modeloPadrao   = "sonnet"
	maxDescSkill   = 300 // limite de internal/assets
	produzPadrao   = "Artifact:Spec"
	maxTimeoutNode = 600
)

// Projecoes são as seis visões do mesmo processo. Todas nascem das mesmas etapas, e Verificar prova que
// listam os mesmos ids na mesma ordem.
type Projecoes struct {
	Spec, Plan, Skill string
	Route, DAG, Loop  []byte
	Avisos            []string // avisos do plangen (por exemplo, critério sem ticket)
}

// Projetar gera SPEC, PLAN, ROUTE, DAG, loop e SKILL a partir do processo e confere a coerência entre eles.
// PLAN e ROUTE passam pelo plangen numa pasta temporária; nada é escrito no projeto nem no cofre.
func Projetar(p Processo, now time.Time) (*Projecoes, error) {
	ord, err := validarProjecao(p)
	if err != nil {
		return nil, err
	}
	d := dagDe(p, ord)
	dagJSON, err := marshalDAG(d)
	if err != nil {
		return nil, err
	}
	loopJSON, err := marshalLoop(p, d, now)
	if err != nil {
		return nil, err
	}
	spec := projetarSpec(p, ord)
	plan, route, avisos, err := projetarPlano(p, ord, spec, now)
	if err != nil {
		return nil, err
	}
	pr := &Projecoes{Spec: spec, Plan: plan, Route: route, DAG: dagJSON, Loop: loopJSON, Skill: projetarSkill(p, ord), Avisos: avisos}
	if err := pr.Verificar(p); err != nil {
		return nil, err
	}
	return pr, nil
}

// limpo deixa o texto livre numa linha só e sem crase: assim um título não abre seção, não forja etapa e
// não fecha o trecho de código de uma linha.
func limpo(s string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(s), " "), "`", "'")
}

func limite(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return strings.TrimSpace(string(r[:n]))
	}
	return s
}

// dagDe monta o grafo executável: um nó por etapa, com o id da etapa, ligado às etapas de que ela depende.
func dagDe(p Processo, ord []Etapa) dag.DAG {
	porID := make(map[string]Etapa, len(ord))
	nodes := make([]dag.Node, len(ord))
	for i, e := range ord {
		porID[e.ID] = e
		nodes[i] = noDe(e, porID)
	}
	return dag.DAG{Schema: dagSchemaURL, ID: p.Nome, Name: p.Nome, Description: limpo(p.Objetivo), Version: p.Versao,
		Trigger: &dag.Trigger{Type: dag.TriggerManual}, Nodes: nodes}
}

func noDe(e Etapa, porID map[string]Etapa) dag.Node {
	n := dag.Node{
		ID: e.ID, TimeoutSec: maxTimeoutNode,
		Routing: dag.RoutingConfig{Mode: dag.ModeSuggested, Model: modeloPadrao, Effort: dag.EffortMedium, Fallback: "code-fast"},
		Output:  dag.HandoffContract{Artifact: e.Produz.Caminho, Gate: e.Gate},
	}
	if e.Compoe != "" {
		n.Type, n.Ref = dag.TypeSubloop, nomeComposto(e.Compoe)
	} else {
		n.Type, n.Ref = dag.TypeSkill, "skills/"+e.Peca
	}
	for _, d := range e.Depende {
		n.DependsOn = append(n.DependsOn, d)
		if n.Inputs == nil {
			n.Inputs = map[string]string{}
		}
		n.Inputs[d] = "file('" + porID[d].Produz.Caminho + "')"
	}
	return n
}

func marshalDAG(d dag.DAG) ([]byte, error) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}
	if _, err := dag.Parse(data); err != nil {
		return nil, fmt.Errorf("%w: o DAG gerado é inválido: %v", ErrInvalido, err)
	}
	return append(data, '\n'), nil
}

func marshalLoop(p Processo, d dag.DAG, now time.Time) ([]byte, error) {
	lp := loop.Loop{ID: p.Nome, Name: p.Nome, Description: d.Description, Version: p.Versao,
		Trigger: d.Trigger, Nodes: d.Nodes, CreatedAt: now.UTC(), UpdatedAt: now.UTC()}
	data, err := json.MarshalIndent(lp, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// quemFaz descreve, em uma linha de texto, o que executa a etapa: a peça ou o processo composto.
func quemFaz(e Etapa) string {
	if e.Compoe != "" {
		return "processo `" + e.Compoe + "`"
	}
	return "peça `" + e.Peca + "`"
}

// ticketsDe cria um ticket por etapa, na ordem das etapas. O título leva o id da etapa, o título dela e quem a
// executa; a dependência entre tickets é a dependência entre etapas.
func ticketsDe(ord []Etapa, validateCmd string) []plangen.Ticket {
	idDe := make(map[string]string, len(ord))
	for i, e := range ord {
		idDe[e.ID] = fmt.Sprintf("TICK-01.%d", i+1)
	}
	tickets := make([]plangen.Ticket, len(ord))
	for i, e := range ord {
		t := plangen.Ticket{
			ID: idDe[e.ID], EpicID: "E01", Nature: "micro", ExecutionMode: "sequential",
			Title:         fmt.Sprintf("`%s` — %s — %s", e.ID, limpo(e.Titulo), quemFaz(e)),
			TargetFiles:   []string{e.Produz.Caminho},
			SpecCriteria:  []string{fmt.Sprintf("AC-E01.%d", i+1)},
			ValidationCmd: validateCmd,
			Gate:          e.Gate,
		}
		for _, d := range e.Depende {
			t.DependsOn = append(t.DependsOn, idDe[d])
		}
		tickets[i] = t
	}
	return tickets
}

// projetarPlano gera PLAN e ROUTE com o plangen, numa pasta temporária que some ao fim.
func projetarPlano(p Processo, ord []Etapa, spec string, now time.Time) (plan string, route []byte, avisos []string, err error) {
	stage, err := os.MkdirTemp("", "legoloop-processo-*")
	if err != nil {
		return "", nil, nil, err
	}
	defer os.RemoveAll(stage)
	dir := filepath.Join(stage, "docs", "specs")
	if err = os.MkdirAll(dir, 0o755); err != nil {
		return "", nil, nil, err
	}
	specPath := filepath.Join(dir, fmt.Sprintf("%s_SPEC_%s.md", now.Format("2006-01-02"), strings.ToUpper(strings.ReplaceAll(strings.TrimPrefix(p.Nome, "legoloop-"), "-", "_"))))
	if err = os.WriteFile(specPath, []byte(spec), 0o644); err != nil {
		return "", nil, nil, err
	}
	res, err := plangen.Render(stage, specPath, ticketsDe(ord, "legoloop validate ~/.aidd/canonical/loops/"+p.Nome+".dag.json"), plangen.SupervisionEpic, now)
	if err != nil {
		return "", nil, nil, err
	}
	planB, err := os.ReadFile(res.PlanMarkdownPath)
	if err != nil {
		return "", nil, nil, err
	}
	routeB, err := os.ReadFile(res.RouteJSONPath)
	return string(planB), routeB, res.Warnings, err
}
