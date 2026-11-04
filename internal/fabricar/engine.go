package fabricar

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/interview"
	"legoloop/internal/processo"
)

// Deps reúne tudo que o motor recebe de fora, para ser hermético e determinístico.
type Deps struct {
	LLM          interview.LLMClient  // entrevista e SPEC; obrigatório
	Scanner      Scanner              // nil usa o DiskScanner real
	Clock        func() time.Time     // nil usa time.Now
	Roots        Roots                // obrigatório
	FS           FS                   // nil usa o disco real
	Harnesses    []processo.Harness   // harnesses detectados; a skill é ligada a todos, salvo Request.NoLink
	Link         processo.LinkFunc    // nil = symlink ou cópia; os testes trocam para provar o desfazer
	Ask          AskFunc              // modo interativo: entrega as perguntas ao usuário
	In           io.Reader            // modo interativo sem Ask: lê as respostas daqui
	Out          io.Writer            // idem: escreve as perguntas aqui
	Progress     ProgressFunc         // opcional: avisa o início de cada uma das Layers camadas
	OnTurn       func(confidence int) // opcional: recebe a confiança antes de cada pergunta da entrevista
	Interview    *Clarification       // opcional: entrevista já feita (Chat); o motor a usa em vez de entrevistar
	ConfirmDraft ConfirmDraftFunc     // opcional: confirma gerar a SPEC quando a entrevista ficou abaixo da meta
	Context      interview.Context    // opcional: o que a entrevista sabe do projeto (objetivo, AGENTS.md, peças aptas)
	Aprovar      ApproveFunc          // opcional: o plano único (arquivos, links, pausas extras) e a última palavra do usuário
}

// ErrAborted: o usuário recusou a gravação; nada foi escrito.
var ErrAborted = errors.New("fabricar: cancelado pelo usuário")

// Layers é o número de camadas que Fabricar percorre: Clarificação, Varredura, Fusão, Tríade e Template.
const Layers = 5

// ProgressFunc recebe o número da camada (1..Layers) e o seu nome, antes de ela começar.
type ProgressFunc func(layer int, name string)

func (d Deps) progress(layer int, name string) {
	if d.Progress != nil {
		d.Progress(layer, name)
	}
}

// Result representa a entrega canônica da fábrica de processos. Sucesso é err == nil.
type Result struct {
	CanonicalName string             `json:"canonical_name"` // Ex: "legoloop-processo-do-joaozinho"
	Slug          string             `json:"slug"`           // "processo-do-joaozinho"
	CreatedAt     time.Time          `json:"created_at"`     // UTC, vindo de Deps.Clock
	Artifacts     CanonicalArtifacts `json:"artifacts"`
	CuratedPieces []DiscoveredPiece  `json:"curated_pieces"` // Ordem estável: Score desc, Path asc
	SkippedRoots  []SkippedRoot      `json:"skipped_roots"`
	Warnings      []string           `json:"warnings"`
	Report        *Report            `json:"report,omitempty"`
}

// prepare aplica os padrões e recusa o que não tem padrão seguro.
func (d Deps) prepare(req Request) (Deps, error) {
	if d.LLM == nil {
		return d, errors.New("fabricar: a Tríade exige um LLM (Deps.LLM)")
	}
	if d.Roots.Vault == "" || d.Roots.Docs == "" {
		return d, errors.New("fabricar: Deps.Roots.Vault e Deps.Roots.Docs são obrigatórios")
	}
	if d.Scanner == nil {
		d.Scanner = NewDiskScanner()
	}
	if d.Clock == nil {
		d.Clock = time.Now
	}
	if d.FS == nil {
		d.FS = OSFS()
	}
	if req.Interactive && d.Ask == nil {
		if d.In == nil || d.Out == nil {
			return d, errors.New("fabricar: o modo interativo exige Deps.Ask ou Deps.In e Deps.Out")
		}
		d.Ask = lineAsker(d.In, d.Out)
	}
	return d, nil
}

// lineAsker pergunta em Out e lê uma linha de In. Entrada encerrada aborta a entrevista.
func lineAsker(in io.Reader, out io.Writer) AskFunc {
	r := bufio.NewReader(in)
	return func(question, guess string) (string, error) {
		fmt.Fprintf(out, "%s\n[palpite: %s]\n> ", question, guess)
		line, err := r.ReadString('\n')
		if strings.TrimSpace(line) == "" && err != nil {
			return "", errors.New("fabricar: entrada encerrada durante a entrevista")
		}
		return strings.TrimSpace(line), nil
	}
}

// Fabricar conduz o pedido pelas camadas e materializa o processo `legoloop-<slug>`: SKILL.md (Chat),
// .loop.json (Studio) e .dag.json (CLI), mais INTENT, SPEC, PLAN, ROUTE e TEMPLATE em docs/.
// Nada é gravado antes de a Tríade estar pronta, e a gravação é atômica: falha de LLM não deixa artefato parcial.
func Fabricar(ctx context.Context, req Request, deps Deps) (*Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	deps, err := deps.prepare(req)
	if err != nil {
		return nil, err
	}
	name, _ := NormalizeProcessName(req.RawName)
	if err := preflight(deps, name, req.Force); err != nil {
		return nil, err
	}
	deps.progress(1, "Clarificação")
	clar, err := clarify(ctx, req, deps)
	if err != nil {
		return nil, err
	}
	if err := deps.confirmDraft(clar); err != nil {
		return nil, err
	}
	scan, merged, err := curate(ctx, req, deps)
	if err != nil {
		return nil, err
	}
	return assemble(ctx, req, deps, name, clar, scan, merged)
}

// clarify é a Camada 1. Com Deps.Interview a entrevista já foi feita (no Chat) e o motor não pergunta de novo.
func clarify(ctx context.Context, req Request, deps Deps) (Clarification, error) {
	if deps.Interview != nil {
		return *deps.Interview, nil
	}
	c := NewClarifier(deps.LLM)
	c.OnTurn = deps.OnTurn
	c.Context = deps.Context
	return c.Clarify(ctx, req, deps.Ask)
}

// preflight antecipa o conflito antes de gastar LLM e varredura.
func preflight(d Deps, name string, force bool) error {
	if force {
		return nil
	}
	probe := []PlannedFile{
		{Path: filepath.Join(d.Roots.Vault, "skills", name, "SKILL.md")},
		{Path: filepath.Join(d.Roots.Vault, "loops", name+".loop.json")},
		{Path: filepath.Join(d.Roots.Vault, "loops", name+".dag.json")},
	}
	if found := processo.Existentes(d.FS, probe); len(found) > 0 {
		return fmt.Errorf("%w: %s", ErrExists, strings.Join(found, ", "))
	}
	return nil
}

// curate é a Camada 2 (varredura) e a Camada 3 (fusão).
func curate(ctx context.Context, req Request, d Deps) (ScanResult, Merged, error) {
	d.progress(2, "Varredura")
	roots, skipped := req.TargetDrives, []SkippedRoot(nil)
	if len(roots) == 0 {
		roots, skipped = SystemRoots(ctx)
	}
	scan, err := d.Scanner.Scan(ctx, roots, ScanOptions{})
	if err != nil {
		return scan, Merged{}, err
	}
	scan.Skipped = append(skipped, scan.Skipped...)
	d.progress(3, "Fusão")
	return scan, ArchitectMerge(scan.Pieces, MergeOptions{TrustRemovable: req.TrustRemovable}), nil
}

func producesFor(c Category) string {
	switch c {
	case CategoryReverse:
		return "Artifact:ReverseStructure"
	case CategoryWeb:
		return "Artifact:RenderedPage"
	}
	return "Artifact:Code"
}

// assemble é a Camada 4 (Tríade), a Camada 5 (Template) e a gravação.
func assemble(ctx context.Context, req Request, d Deps, name string, clar Clarification, scan ScanResult, m Merged) (*Result, error) {
	now := d.Clock().UTC()
	slug := strings.TrimPrefix(name, namePrefix)
	ps := ProcessSpec{Name: name, Description: clar.Mirror.Outcome, Goal: clar.Mirror.Outcome, Model: req.ModelChoice, Produces: producesFor(req.Category), Steps: m.Steps}
	dagPath := filepath.Join(d.Roots.Vault, "loops", name+".dag.json")
	d.progress(4, "Tríade")
	tri, err := buildTriad(ctx, d.LLM, clar, ps, dagPath, now)
	if err != nil {
		return nil, err
	}
	d.progress(5, "Template e gravação")
	files, err := plannedFiles(d.Roots, ps, tri, now)
	if err != nil {
		return nil, err
	}
	vinculos, lwarns, err := aprovar(d, req, files, m, name)
	if err != nil {
		return nil, err
	}
	plano := novoPlano(d, req, files, m, vinculos)
	rep := buildReport(d, clar, scan, m, tri, files, plano)
	art, err := Materialize(ctx, MaterializePlan{Name: name, Roots: d.Roots, Files: files, Force: req.Force, DryRun: req.DryRun,
		Vinculos: vinculos, Link: d.Link}, d.FS)
	if err != nil {
		return nil, err
	}
	warns := append(append(clar.Warnings(), scan.Warnings...), m.Warnings...)
	return &Result{CanonicalName: name, Slug: slug, CreatedAt: now, Artifacts: art, CuratedPieces: m.Used,
		SkippedRoots: scan.Skipped, Warnings: append(warns, lwarns...), Report: rep}, nil
}

func filePaths(files []PlannedFile) []string {
	paths := make([]string, len(files))
	for i, f := range files {
		paths[i] = f.Path
	}
	return paths
}

// plannedFiles monta a lista final de arquivos: os três meios no cofre e a governança em docs/.
func plannedFiles(roots Roots, ps ProcessSpec, tri triadFiles, now time.Time) ([]PlannedFile, error) {
	skill, err := EmitSkill(ps)
	if err != nil {
		return nil, err
	}
	loopJSON, err := EmitLoop(ps, now)
	if err != nil {
		return nil, err
	}
	dagJSON, err := EmitDag(ps)
	if err != nil {
		return nil, err
	}
	docs := func(sub, file string) string { return filepath.Join(roots.Docs, "docs", sub, file) }
	vault := func(parts ...string) string { return filepath.Join(append([]string{roots.Vault}, parts...)...) }
	return []PlannedFile{
		{Kind: "skill", Path: vault("skills", ps.Name, "SKILL.md"), Data: []byte(skill)},
		{Kind: "loop", Path: vault("loops", ps.Name+".loop.json"), Data: loopJSON},
		{Kind: "dag", Path: vault("loops", ps.Name+".dag.json"), Data: dagJSON},
		{Kind: "intent", Path: docs("intents", tri.Intent.Name), Data: tri.Intent.Data},
		{Kind: "spec", Path: docs("specs", tri.Spec.Name), Data: tri.Spec.Data},
		{Kind: "plan", Path: docs("plans", tri.Plan.Name), Data: tri.Plan.Data},
		{Kind: "route", Path: docs("routes", tri.Route.Name), Data: tri.Route.Data},
		{Kind: "template", Path: docs("templates", templateFileName(ps.Name)), Data: []byte(EmitTemplate(ps))},
	}, nil
}
