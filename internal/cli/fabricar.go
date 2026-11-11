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
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

// FabricarDeps injeta o que o `fabricar` não deve buscar sozinho; o valor zero usa o ambiente real.
// É assim que os testes rodam sem rede, sem disco real e com data fixa.
type FabricarDeps struct {
	LLM       interview.LLMClient // nil = cliente do 9Router descoberto no inventário
	Matriz    fabricar.MatrizFunc // nil = registry_matrix.json do cofre
	Clock     func() time.Time    // nil = time.Now
	Harnesses []processo.Harness  // nil = harnesses detectados no HOME
	Link      processo.LinkFunc   // nil = symlink ou cópia; testes trocam para provar o desfazer
	BrainOpts []brain.Option      // ajustes do cérebro (URL do Ollama, prazo, cliente HTTP); usado em testes
}

// fabricarOpts reúne as flags do comando.
type fabricarOpts struct {
	yes, asJSON, force, dry, trust bool
	noLink                         bool
	semEntrevista, aceitar         bool
	entrevista                     string
	goal, model, category, effort  string
	vault, docs                    string
	link                           listFlag
}

func newFabricarFlags(o *fabricarOpts) *flag.FlagSet {
	fs := flag.NewFlagSet("fabricar", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&o.yes, "yes", false, "modo autônomo: não pergunta nada")
	fs.StringVar(&o.entrevista, "entrevista", "", "usa a entrevista já feita no Chat (identificador da sessão)")
	fs.BoolVar(&o.semEntrevista, "sem-entrevista", false, "com --yes: segue sem entrevista; o INTENT sai sem ela")
	fs.BoolVar(&o.aceitar, "aceitar-rascunho", false, "com --entrevista: gera a SPEC mesmo com a confiança abaixo da meta")
	fs.StringVar(&o.goal, "goal", "", "objetivo em linguagem natural (padrão: o próprio nome)")
	fs.StringVar(&o.model, "model", "", "modelo que conduz a fabricação, executor 9router (padrão: o primeiro listado por legoloop doctor)")
	fs.StringVar(&o.effort, "effort", "", "esforço cognitivo (low, medium, high)")
	fs.StringVar(&o.category, "category", "build", "categoria: build, reverse, web ou mobile")
	fs.StringVar(&o.vault, "vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	fs.StringVar(&o.docs, "docs", "", "raiz do projeto que recebe docs/ (padrão: diretório atual)")
	fs.BoolVar(&o.force, "force", false, "sobrescreve um processo existente (guarda .legoloop-bak)")
	fs.BoolVar(&o.dry, "dry-run", false, "calcula tudo e não grava nada")
	fs.BoolVar(&o.trust, "trust-removable", false, "aceita peças de mídia removível na fusão")
	fs.BoolVar(&o.asJSON, "json", false, "imprime o Result como JSON")
	fs.Var(&o.link, "link-harness", "liga só aos harnesses citados, ex.: claude (repetível); padrão: todos os detectados")
	fs.BoolVar(&o.noLink, "no-link", false, "grava só no cofre, sem ligar a nenhum harness")
	return fs
}

const fabricarUsage = `Uso: legoloop fabricar "<nome do processo>" [opções]

Cria o processo legoloop-<nome>: skill (Chat), documento do Studio e DAG (CLI), mais os documentos
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
	if len(args) > 0 && args[0] == "entrevista" {
		return cmdEntrevista(ctx, args[1:], d)
	}
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
	if code, stop := checkEntrevista(o, d); stop {
		return code
	}
	req := fabricar.Request{RawName: name, Goal: o.goal, Category: fabricar.Category(o.category), ModelChoice: o.model,
		Interactive: !o.yes, Force: o.force, DryRun: o.dry, TrustRemovable: o.trust, LinkHarnesses: o.link, NoLink: o.noLink}
	probe := req
	if probe.ModelChoice == "" {
		probe.ModelChoice = "auto" // só para validar nome e categoria antes de consultar o inventário
	}
	if err := probe.Validate(); err != nil {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	b, why, err := resolveFabricarBrain(ctx, o, d)
	if err != nil {
		return failBrain(d, err)
	}
	req.ModelChoice = b.Model()
	announceBrain(d, why, o.asJSON)
	fd, err := buildFabricarDeps(o, req, d, b)
	if err == nil {
		err = prepareInterview(ctx, &req, &fd, o, d)
	}
	if err != nil {
		return failFabricar(d, err)
	}
	res, err := fabricar.Fabricar(ctx, req, fd)
	if err != nil {
		if res != nil && errors.Is(err, fabricar.ErrLacuna) {
			reportLacuna(d, res, o) // dry-run: o relatório e o cartão destacam a lacuna antes do erro
		}
		return failFabricar(d, err)
	}
	return reportFabricar(d, res, o)
}

func resolveFabricarBrain(ctx context.Context, o fabricarOpts, d Deps) (brain.Brain, string, error) {
	opts := d.brainOpts()
	if o.effort != "" {
		switch o.effort {
		case "low", "medium", "high":
			opts = append(opts, brain.WithEffort(brain.Effort(o.effort)))
		default:
			return nil, "", fmt.Errorf("esforço cognitivo inválido %q (use low, medium ou high)", o.effort)
		}
	}
	inv := d.discover(ctx)
	return brain.Pick(inv, o.model, opts...)
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
	if errors.Is(err, fabricar.ErrExists) {
		fmt.Fprintln(d.Stderr, "erro:", err)
		fmt.Fprintln(d.Stderr, ste.Msg("fabricar.existe_oferta_refinar"))
		return 1
	}
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Fprintln(d.Stderr, "erro:", ste.Msg("fabricar.cerebro_timeout"))
		return 1
	}
	if errors.Is(err, fabricar.ErrAborted) {
		fmt.Fprintln(d.Stderr, "cancelado: zero arquivos gravados")
		return 1
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
	if x.Harnesses == nil && !req.NoLink {
		x.Harnesses = detectHarnesses(d.Home)
	}
	progress := d.Stdout
	if o.asJSON {
		progress = d.Stderr // mantém o stdout como JSON puro
	}
	return fabricar.Deps{LLM: x.LLM, Matriz: x.Matriz, Clock: x.Clock, Roots: roots, Harnesses: x.Harnesses, Link: x.Link,
		Progress: printProgress(progress), Aprovar: aprovarSemPerguntar(progress)}, nil
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

// detectHarnesses procura os harnesses no HOME informado (vazio = o HOME real).
func detectHarnesses(home string) []processo.Harness {
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	return processo.DetectarHarnesses(processo.OSFS(), home)
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
	rows := [][2]string{{"skill (Chat)", a.SkillPath}, {"documento (Studio)", a.LoopPath}, {"dag (CLI)", a.DagPath},
		{"INTENT", a.IntentPath}, {"SPEC", a.SpecPath}, {"PLAN", a.PlanPath}, {"ROUTE", a.RoutePath}, {"TEMPLATE", a.TemplatePath}}
	for _, r := range rows {
		fmt.Fprintf(d.Stdout, "  %-14s %s\n", r[0], r[1])
	}
	for _, l := range a.SkillLinks {
		fmt.Fprintf(d.Stdout, "  %-14s %s\n", "link harness", l)
	}
	fmt.Fprintf(d.Stdout, "\nPeças curadas: %d\nExecute com: legoloop run %s\n", len(res.CuratedPieces), a.DagPath)
	PrintReportTexto(d.Stdout, res.Report)
	return 0
}

func init() {
	for id, msg := range map[string]string{
		"fabricar.sem_entrevista_exige":  "INTENT sem entrevista não vale com --yes. Faça a entrevista com legoloop fabricar entrevista iniciar. Ou use --sem-entrevista para seguir sem ela.",
		"fabricar.aviso_sem_entrevista":  "Aviso: INTENT sem entrevista. O objetivo vem só do que você escreveu.",
		"fabricar.entrevista_dupla":      "Use --entrevista ou --sem-entrevista, não os dois.",
		"fabricar.rascunho_exige_aceite": "A entrevista ficou em %d%%. A meta é %d%%. Para gerar a SPEC mesmo assim, use --aceitar-rascunho.",
		"fabricar.existe_oferta_refinar": "O processo já existe. Use legoloop refinar para alterá-lo ou --force para sobrescrever.",
		"fabricar.cerebro_timeout":       "O cérebro excedeu o tempo limite. Aumente o prazo com a opção de configuração ou tente outro modelo.",
		"fabricar.cancelado":             "Operação cancelada pelo usuário. Zero arquivos gravados.",
	} {
		ste.Default.Register(id, msg)
	}
}

// checkEntrevista aplica a regra do modo autônomo: --yes sem entrevista só segue com --sem-entrevista.
// Devolve stop=true quando o comando deve terminar com o código dado.
func checkEntrevista(o fabricarOpts, d Deps) (code int, stop bool) {
	switch {
	case o.semEntrevista && o.entrevista != "":
		fmt.Fprintln(d.Stderr, "erro:", ste.Msg("fabricar.entrevista_dupla"))
		return 2, true
	case o.yes && o.entrevista == "" && !o.semEntrevista:
		fmt.Fprintln(d.Stderr, "erro:", ste.Msg("fabricar.sem_entrevista_exige"))
		return 2, true
	case o.semEntrevista:
		fmt.Fprintln(d.Stderr, ste.Msg("fabricar.aviso_sem_entrevista"))
	}
	return 0, false
}

// prepareInterview liga o motor à entrevista: o terminal (modo interativo) e/ou a sessão já feita no Chat.
func prepareInterview(ctx context.Context, req *fabricar.Request, fd *fabricar.Deps, o fabricarOpts, d Deps) error {
	if req.Interactive {
		if err := setupInteractive(req, fd, d); err != nil {
			return err
		}
	}
	if o.entrevista == "" {
		return nil
	}
	store := interview.Store{Dir: filepath.Join(fd.Roots.Vault, "sessions")}
	sess, err := store.Load(o.entrevista)
	if err != nil {
		return err
	}
	clar := fabricar.ClarificationFromSession(sess)
	fd.Interview = &clar
	if !req.Interactive {
		fd.ConfirmDraft = draftByFlag(o.aceitar)
	}
	return nil
}

// draftByFlag é a confirmação do modo sem terminal: só o --aceitar-rascunho autoriza gerar a SPEC com confiança baixa.
func draftByFlag(aceitar bool) fabricar.ConfirmDraftFunc {
	return func(c fabricar.Clarification) (bool, error) {
		if aceitar {
			return true, nil
		}
		return false, fmt.Errorf("%w: %s", fabricar.ErrLowConfidence,
			ste.Msg("fabricar.rascunho_exige_aceite", c.Confidence, fabricar.ConfidenceTarget))
	}
}
