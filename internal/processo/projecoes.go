package processo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/plangen"
	"legoloop/internal/studiodoc"
	"legoloop/internal/triagem"
)

const (
	dagSchemaURL   = "https://legoloop.dev/schemas/dag.v1.json"
	modeloPadrao   = "sonnet"
	maxDescSkill   = 300 // limite de internal/assets
	produzPadrao   = "Artifact:Spec"
	maxTimeoutNode = 600
)

// FichaTecnica representa a documentação técnica de uma etapa/nó do processo.
type FichaTecnica struct {
	EtapaID  string
	Conteudo string
}

// Projecoes são as visões do mesmo processo no Quinteto ISO. Todas nascem das mesmas etapas, e Verificar prova que
// listam os mesmos ids na mesma ordem.
type Projecoes struct {
	Spec, Plan, Skill string
	Route, DAG, Loop  []byte
	Triagem, Livro    []byte
	Manual            string
	Fichas            []FichaTecnica
	Avisos            []string // avisos do plangen (por exemplo, critério sem ticket)
	// Nomes dos arquivos em docs/specs, docs/plans, docs/routes, docs/dags e docs/manuals.
	SpecNome, PlanNome, RouteNome string
	DAGNome, ManualNome           string
}

// CheckParity compara o hash SHA-256 de dois arquivos no disco.
func CheckParity(pathA, pathB string) (bool, error) {
	dataA, err := os.ReadFile(pathA)
	if err != nil {
		return false, err
	}
	dataB, err := os.ReadFile(pathB)
	if err != nil {
		return false, err
	}
	hashA := sha256.Sum256(dataA)
	hashB := sha256.Sum256(dataB)
	return hex.EncodeToString(hashA[:]) == hex.EncodeToString(hashB[:]), nil
}

// Projetar gera SPEC, PLAN, ROUTE, DAG, documento do Studio e SKILL a partir do processo e confere a coerência entre eles.
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
	studioJSON, err := marshalStudio(p, d, now)
	if err != nil {
		return nil, err
	}
	spec := projetarSpec(p, ord)
	pl, err := projetarPlano(p, ord, spec, now)
	if err != nil {
		return nil, err
	}
	dagNome := strings.Replace(strings.Replace(pl.specNome, "SPEC", "DAG", 1), ".md", ".json", 1)
	triagemJSON, err := projetarTriagem(p, ord)
	if err != nil {
		return nil, err
	}
	livroJSON, err := projetarLivro(p, ord)
	if err != nil {
		return nil, err
	}
	manualNome, manualMD := projetarManual(p, ord, now)
	pr := &Projecoes{
		Spec: spec, Plan: pl.plan, Route: pl.route, DAG: dagJSON, Loop: studioJSON, Skill: projetarSkill(p, ord),
		Triagem: triagemJSON, Livro: livroJSON, Manual: manualMD, Fichas: projetarFichas(p, ord),
		Avisos: pl.avisos, SpecNome: pl.specNome, PlanNome: pl.planNome, RouteNome: pl.routeNome,
		DAGNome: dagNome, ManualNome: manualNome,
	}
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
		n.Output.SubArtifact = e.Origem
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
	adicionarMetadadosCognitivos(&n, e, porID)
	return n
}

func adicionarMetadadosCognitivos(n *dag.Node, e Etapa, porID map[string]Etapa) {
	n.Name = limpo(e.Titulo)
	n.Nature = "micro"
	n.CognitiveComplexity = "simples"
	n.Isolation = map[string]any{"workspace": "local", "executor": "orchestrix_hermetic", "subagents_allowed": false}
	n.TokenBudget = map[string]any{"max_prompt_tokens": 10000, "max_completion_tokens": 3000}
	n.Tooling = map[string]any{"runtime": "powershell"}
	n.Status = "PENDING"
	if e.Produz.Caminho != "" {
		n.ContextOut = []string{e.Produz.Caminho}
	}
	for _, d := range e.Depende {
		if prod := porID[d].Produz.Caminho; prod != "" {
			n.ContextIn = append(n.ContextIn, prod)
		}
	}
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

func marshalStudio(p Processo, d dag.DAG, now time.Time) ([]byte, error) {
	lp := studiodoc.Processo{ID: p.Nome, Name: p.Nome, Description: d.Description, Version: p.Versao,
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

// planoProjetado é o que o plangen devolve: PLAN e ROUTE, os avisos e os nomes dos três documentos.
type planoProjetado struct {
	plan                          string
	route                         []byte
	avisos                        []string
	specNome, planNome, routeNome string
}

// projetarPlano gera PLAN e ROUTE com o plangen, numa pasta temporária que some ao fim.
func projetarPlano(p Processo, ord []Etapa, spec string, now time.Time) (planoProjetado, error) {
	stage, err := os.MkdirTemp("", "legoloop-processo-*")
	if err != nil {
		return planoProjetado{}, err
	}
	defer os.RemoveAll(stage)
	specPath, err := gravarSpec(stage, p.Nome, spec, now)
	if err != nil {
		return planoProjetado{}, err
	}
	res, err := plangen.Render(stage, specPath, ticketsDe(ord, "legoloop validate "+dagNoCofre(p.Nome)), plangen.SupervisionEpic, now)
	if err != nil {
		return planoProjetado{}, err
	}
	planB, err := os.ReadFile(res.PlanMarkdownPath)
	if err != nil {
		return planoProjetado{}, err
	}
	routeB, err := os.ReadFile(res.RouteJSONPath)
	return planoProjetado{plan: string(planB), route: routeB, avisos: res.Warnings, specNome: filepath.Base(specPath),
		planNome: filepath.Base(res.PlanMarkdownPath), routeNome: filepath.Base(res.RouteJSONPath)}, err
}

// gravarSpec grava a SPEC em <stage>/docs/specs com o nome <data>_SPEC_<NOME>.md, que o plangen espera.
func gravarSpec(stage, nome, spec string, now time.Time) (string, error) {
	dir := filepath.Join(stage, "docs", "specs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	slug := strings.ToUpper(strings.ReplaceAll(strings.TrimPrefix(nome, "legoloop-"), "-", "_"))
	specPath := filepath.Join(dir, fmt.Sprintf("%s_SPEC_%s.md", now.Format("2006-01-02"), slug))
	return specPath, os.WriteFile(specPath, []byte(spec), 0o644)
}

func projetarTriagem(p Processo, ord []Etapa) ([]byte, error) {
	ids := make([]string, len(ord))
	for i, e := range ord {
		ids[i] = e.ID
	}
	spec := triagem.GerarTriagemPadrao(p.Nome, ids)
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func projetarLivro(p Processo, ord []Etapa) ([]byte, error) {
	caps := make([]map[string]string, len(ord))
	for i, e := range ord {
		caps[i] = map[string]string{"id": e.ID, "titulo": limpo(e.Titulo), "peca": e.Peca}
	}
	livro := map[string]any{
		"slug":      p.Nome,
		"titulo":    "Manual Didático: " + p.Nome,
		"descricao": limpo(p.Objetivo),
		"capitulos": caps,
	}
	data, err := json.MarshalIndent(livro, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func projetarManual(p Processo, ord []Etapa, now time.Time) (string, string) {
	slug := strings.ToUpper(strings.ReplaceAll(strings.TrimPrefix(p.Nome, "legoloop-"), "-", "_"))
	nome := fmt.Sprintf("%s_MANUAL_%s.md", now.Format("2006-01-02"), slug)
	var b strings.Builder
	fmt.Fprintf(&b, "# Manual Operacional: %s\n\n%s\n\n## Etapas\n\n", p.Nome, limpo(p.Objetivo))
	for _, e := range ord {
		fmt.Fprintf(&b, "- `%s`: %s\n", e.ID, limpo(e.Titulo))
	}
	return nome, b.String()
}

func projetarFichas(p Processo, ord []Etapa) []FichaTecnica {
	fichas := make([]FichaTecnica, len(ord))
	for i, e := range ord {
		var b strings.Builder
		fmt.Fprintf(&b, "# Ficha Técnica: %s\n\n- **Processo:** %s\n- **Título:** %s\n- **Peça:** %s\n- **Saída:** %s\n",
			e.ID, p.Nome, limpo(e.Titulo), e.Peca, e.Produz.Caminho)
		fichas[i] = FichaTecnica{EtapaID: e.ID, Conteudo: b.String()}
	}
	return fichas
}
