package fabricar

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"legoloop/internal/paths"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/interview"
	"legoloop/internal/matriz"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

// Deps reúne tudo que o motor recebe de fora, para ser hermético e determinístico.
type Deps struct {
	LLM          interview.LLMClient  // entrevista e SPEC; obrigatório
	Matriz       MatrizFunc           // nil lê registry_matrix.json e vocab.json do cofre (Roots.Vault)
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
	if d.Matriz == nil {
		vault := d.Roots.Vault
		d.Matriz = func() (FonteMatriz, error) { return LerMatriz(vault) }
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
// documento do Studio e .dag.json (CLI), mais INTENT, SPEC, PLAN, ROUTE e TEMPLATE em docs/.
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
	fonte, err := deps.Matriz() // antes da entrevista: sem matriz não vale a pena perguntar nada
	if err != nil {
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
	c, err := consultar(ctx, req, deps, clar, fonte)
	if err != nil {
		return nil, err
	}
	return assemble(ctx, req, deps, name, clar, c)
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
		{Path: filepath.Join(d.Roots.Vault, paths.DirAntigo, name+paths.SufixoStudio)},
		{Path: filepath.Join(d.Roots.Vault, paths.DirAntigo, name+paths.SufixoAntigo)},
		{Path: filepath.Join(d.Roots.Vault, paths.DirAntigo, name+".dag.json")},
		{Path: processo.CaminhoNoCofre(d.Roots.Vault, name)},
	}
	if found := processo.Existentes(d.FS, probe); len(found) > 0 {
		return fmt.Errorf("%w: %s", ErrExists, strings.Join(found, ", "))
	}
	return nil
}

// consulta é o que as Camadas 2 e 3 entregam à Tríade.
type consulta struct {
	fonte  FonteMatriz
	needs  matriz.Needs
	m      Merged
	avisos []string
}

// consultar é a Camada 2 (matriz) e a Camada 3 (fusão por aptidão). O fabricar só lê a matriz: nunca varre
// drives nem a regenera. O cérebro dá as necessidades e explica as decisões; quem decide é o ranking.
func consultar(ctx context.Context, req Request, d Deps, clar Clarification, fonte FonteMatriz) (consulta, error) {
	d.progress(2, "Matriz")
	d.progress(3, "Fusão")
	needs, avisos, err := Necessidades(ctx, d.LLM, objetivoDe(req), clar.Mirror, fonte.Vocab)
	if err != nil {
		return consulta{}, err
	}
	if nat := strings.TrimSpace(clar.Natureza); nat != "" && fonte.Vocab.Has(matriz.CatNatureza, nat) {
		needs.Natureza = nat // a natureza confirmada na entrevista vale sobre a proposta do cérebro
	}
	aptas, cortes := confiaveis(fonte.Entradas, req.TrustRemovable)
	m := MergeByFit(needs, aptas, MergeOptions{Vocab: fonte.Vocab})
	m.Explain(ctx, d.LLM, needs)
	if len(req.TargetDrives) > 0 {
		avisos = append(avisos, ste.Msg("fabricar.drives_ignorados"))
	}
	return consulta{fonte: fonte, needs: needs, m: m, avisos: append(avisos, cortes...)}, nil
}

// objetivoDe é o objetivo do pedido; sem ele, o próprio nome.
func objetivoDe(req Request) string {
	if strings.TrimSpace(req.Goal) != "" {
		return req.Goal
	}
	return req.RawName
}

// assemble é a Camada 4 (o Processo e as seis projeções), a Camada 5 (Template) e a gravação.
func assemble(ctx context.Context, req Request, d Deps, name string, clar Clarification, c consulta) (*Result, error) {
	m := c.m
	now := d.Clock().UTC()
	d.progress(4, "Tríade")
	p, pr, err := projetar(name, objetivoDoProcesso(req, clar), c, now)
	if req.DryRun && errors.Is(err, ErrLacuna) {
		return relatorioDaLacuna(d, name, clar, c, p, now), err // a lacuna não para antes do relatório
	}
	if err != nil {
		return nil, err
	}
	d.progress(5, "Template e gravação")
	files, err := plannedFiles(d.Roots, p, pr, intentDe(name, clar, p.Natureza, now))
	if err != nil {
		return nil, err
	}
	vinculos, lwarns, err := aprovar(d, req, files, m, name)
	if err != nil {
		return nil, err
	}
	plano := novoPlano(d, req, files, m, vinculos)
	rep := buildReport(d, clar, c, p, files, plano)
	art, err := Materialize(ctx, MaterializePlan{Name: name, Roots: d.Roots, Files: files, Force: req.Force, DryRun: req.DryRun,
		Vinculos: vinculos, Link: d.Link}, d.FS)
	if err != nil {
		return nil, err
	}
	warns := append(append(append(clar.Warnings(), c.avisos...), m.Warnings...), pr.Avisos...)
	return &Result{CanonicalName: name, Slug: strings.TrimPrefix(name, namePrefix), CreatedAt: now, Artifacts: art,
		CuratedPieces: m.Used, Warnings: append(warns, lwarns...), Report: rep}, nil
}

func filePaths(files []PlannedFile) []string {
	lista := make([]string, len(files))
	for i, f := range files {
		lista[i] = f.Path
	}
	return lista
}
