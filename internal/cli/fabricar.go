package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"legoloop/internal/brain"
	"legoloop/internal/canon"
	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
)

// FabricarDeps injeta o que o `fabricar` não deve buscar sozinho; o valor zero usa o ambiente real.
// É assim que os testes rodam sem rede, sem disco real e com data fixa.
type FabricarDeps struct {
	LLM         interview.LLMClient            // nil = cliente do 9Router descoberto no inventário
	Scanner     fabricar.Scanner               // nil = varredura real dos discos
	Clock       func() time.Time               // nil = time.Now
	HarnessDirs []string                       // nil = ~/.<harness>/skills para cada --link-harness
	DriveRoots  func(context.Context) []string // nil = drives do sistema; só o modo interativo os lista
	BrainOpts   []brain.Option                 // ajustes do cérebro (URL do Ollama, prazo, cliente HTTP); usado em testes
}

// fabricarOpts reúne as flags do comando.
type fabricarOpts struct {
	yes, asJSON, force, dry, trust bool
	goal, model, category          string
	vault, docs                    string
	link                           listFlag
}

func newFabricarFlags(o *fabricarOpts) *flag.FlagSet {
	fs := flag.NewFlagSet("fabricar", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&o.yes, "yes", false, "modo autônomo: não pergunta nada")
	fs.StringVar(&o.goal, "goal", "", "objetivo em linguagem natural (padrão: o próprio nome)")
	fs.StringVar(&o.model, "model", "", "modelo que conduz a fabricação, executor 9router (padrão: o primeiro listado por legoloop doctor)")
	fs.StringVar(&o.category, "category", "build", "categoria: build, reverse, web ou mobile")
	fs.StringVar(&o.vault, "vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	fs.StringVar(&o.docs, "docs", "", "raiz do projeto que recebe docs/ (padrão: diretório atual)")
	fs.BoolVar(&o.force, "force", false, "sobrescreve um processo existente (guarda .legoloop-bak)")
	fs.BoolVar(&o.dry, "dry-run", false, "calcula tudo e não grava nada")
	fs.BoolVar(&o.trust, "trust-removable", false, "aceita peças de mídia removível na fusão")
	fs.BoolVar(&o.asJSON, "json", false, "imprime o Result como JSON")
	fs.Var(&o.link, "link-harness", "liga a skill à pasta do harness, ex.: claude (repetível)")
	return fs
}

const fabricarUsage = `Uso: legoloop fabricar "<nome do processo>" [opções]

Cria o processo legoloop-<nome>: skill (Chat), loop (Studio) e DAG (CLI), mais os documentos
de governança (INTENT, SPEC, PLAN, ROUTE e TEMPLATE). Não confundir com canonize, que desduplica.

Opções:
`

func printFabricarHelp(d Deps) {
	var o fabricarOpts
	fs := newFabricarFlags(&o)
	fs.SetOutput(d.Stdout)
	fmt.Fprint(d.Stdout, fabricarUsage)
	fs.PrintDefaults()
}

func cmdFabricar(ctx context.Context, args []string, d Deps) int {
	var o fabricarOpts
	fs := newFabricarFlags(&o)
	name, rest := splitPositional(args)
	if err := fs.Parse(rest); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printFabricarHelp(d)
			return 0
		}
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	if name == "" {
		name = fs.Arg(0)
	}
	if name == "" {
		fmt.Fprint(d.Stderr, fabricarUsage)
		return 2
	}
	return runFabricar(ctx, name, o, d)
}

// runFabricar valida, resolve o ambiente e executa a fábrica. Pedido inválido sai com 2; falha de execução, com 1.
func runFabricar(ctx context.Context, name string, o fabricarOpts, d Deps) int {
	req := fabricar.Request{RawName: name, Goal: o.goal, Category: fabricar.Category(o.category), ModelChoice: o.model,
		Interactive: !o.yes, Force: o.force, DryRun: o.dry, TrustRemovable: o.trust, LinkHarnesses: o.link}
	probe := req
	if probe.ModelChoice == "" {
		probe.ModelChoice = "auto" // só para validar nome e categoria antes de consultar o inventário
	}
	if err := probe.Validate(); err != nil {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	inv := d.discover(ctx)
	b, why, err := brain.Pick(inv, o.model, d.brainOpts()...)
	if err != nil {
		return failBrain(d, err)
	}
	req.ModelChoice = b.Model()
	announceBrain(d, why, o.asJSON)
	fd, err := buildFabricarDeps(o, req, d, b)
	if err == nil && req.Interactive {
		err = setupInteractive(ctx, &req, &fd, d)
	}
	if err != nil {
		return failFabricar(d, err)
	}
	res, err := fabricar.Fabricar(ctx, req, fd)
	if err != nil {
		return failFabricar(d, err)
	}
	return reportFabricar(d, res, o)
}

// failBrain traduz o erro da escolha do cérebro: modelo fora do inventário é uso incorreto (2); sem nenhum
// cérebro disponível é falha de ambiente (1).
func failBrain(d Deps, err error) int {
	if errors.Is(err, brain.ErrModeloDesconhecido) {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	return fail(d, err)
}

// announceBrain diz ao usuário qual cérebro conduz a fabricação. Com --json o aviso vai ao stderr.
func announceBrain(d Deps, why string, asJSON bool) {
	w := d.Stdout
	if asJSON {
		w = d.Stderr
	}
	fmt.Fprintln(w, why)
}

func (d Deps) brainOpts() []brain.Option {
	if d.Fabricar != nil {
		return d.Fabricar.BrainOpts
	}
	return nil
}

func failFabricar(d Deps, err error) int {
	if errors.Is(err, fabricar.ErrInvalidName) || errors.Is(err, fabricar.ErrInvalidRequest) {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	return fail(d, err)
}

// buildFabricarDeps resolve as raízes e as dependências do motor; o que não vier injetado usa o ambiente real.
func buildFabricarDeps(o fabricarOpts, req fabricar.Request, d Deps, b brain.Brain) (fabricar.Deps, error) {
	x := FabricarDeps{}
	if d.Fabricar != nil {
		x = *d.Fabricar
	}
	if x.LLM == nil {
		x.LLM = b
	}
	roots, err := fabricarRoots(o, d)
	if err != nil {
		return fabricar.Deps{}, err
	}
	if x.HarnessDirs == nil {
		x.HarnessDirs = harnessSkillDirs(req.LinkHarnesses, d.Home)
	}
	progress := d.Stdout
	if o.asJSON {
		progress = d.Stderr // mantém o stdout como JSON puro
	}
	return fabricar.Deps{LLM: x.LLM, Scanner: x.Scanner, Clock: x.Clock, Roots: roots, HarnessDirs: x.HarnessDirs,
		Progress: printProgress(progress)}, nil
}

func printProgress(w io.Writer) fabricar.ProgressFunc {
	return func(layer int, name string) { fmt.Fprintf(w, "[%d/%d] %s...\n", layer, fabricar.Layers, name) }
}

func fabricarRoots(o fabricarOpts, d Deps) (fabricar.Roots, error) {
	vault := o.vault
	switch {
	case vault != "":
	case d.Home != "":
		vault = filepath.Join(d.Home, ".aidd", "canonical")
	default:
		v, err := canon.DefaultVault()
		if err != nil {
			return fabricar.Roots{}, err
		}
		vault = v
	}
	docs := o.docs
	if docs == "" {
		docs = "."
	}
	v, err := filepath.Abs(vault)
	if err != nil {
		return fabricar.Roots{}, err
	}
	p, err := filepath.Abs(docs)
	return fabricar.Roots{Vault: v, Docs: p}, err
}

// harnessSkillDirs devolve ~/.<harness>/skills para cada harness pedido.
func harnessSkillDirs(harnesses []string, home string) []string {
	if len(harnesses) == 0 {
		return nil
	}
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	var dirs []string
	for _, h := range harnesses {
		dirs = append(dirs, filepath.Join(home, "."+h, "skills"))
	}
	return dirs
}

func reportFabricar(d Deps, res *fabricar.Result, o fabricarOpts) int {
	for _, w := range res.Warnings {
		fmt.Fprintln(d.Stderr, "aviso:", w)
	}
	if o.asJSON {
		return printJSON(d, res)
	}
	if o.dry {
		fmt.Fprintln(d.Stdout, "\n(dry-run) nada foi gravado; o processo seria criado em:")
	} else {
		fmt.Fprintf(d.Stdout, "\n✓ Processo %s criado\n", res.CanonicalName)
	}
	a := res.Artifacts
	rows := [][2]string{{"skill (Chat)", a.SkillPath}, {"loop (Studio)", a.LoopPath}, {"dag (CLI)", a.DagPath},
		{"INTENT", a.IntentPath}, {"SPEC", a.SpecPath}, {"PLAN", a.PlanPath}, {"ROUTE", a.RoutePath}, {"TEMPLATE", a.TemplatePath}}
	for _, r := range rows {
		fmt.Fprintf(d.Stdout, "  %-14s %s\n", r[0], r[1])
	}
	for _, l := range a.SkillLinks {
		fmt.Fprintf(d.Stdout, "  %-14s %s\n", "link harness", l)
	}
	fmt.Fprintf(d.Stdout, "\nPeças curadas: %d\nExecute com: legoloop run %s\n", len(res.CuratedPieces), a.DagPath)
	return 0
}
