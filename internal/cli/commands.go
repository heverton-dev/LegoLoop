package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/canon"
	"legoloop/internal/dag"
	"legoloop/internal/graph"
	"legoloop/internal/orchestrix"
	"legoloop/internal/provider"
	"legoloop/internal/runner"
	"legoloop/internal/studio"
)

// cmdCanonize implementa TICK-104: --dry-run (só plano) ou --apply (links reais).
func cmdCanonize(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("canonize", d)
	var roots listFlag
	fs.Var(&roots, "root", "diretório a varrer (repetível)")
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	workers := fs.Int("workers", canon.DefaultWorkers, "threads do scanner")
	dry := fs.Bool("dry-run", false, "apenas planeja")
	apply := fs.Bool("apply", false, "aplica os links")
	asJSON := fs.Bool("json", false, "saída JSON")
	includeGit := fs.Bool("include-git", false, "inclui working trees git (links podem ser commitados)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if len(roots) == 0 || *dry == *apply {
		fmt.Fprintln(d.Stderr, "erro: informe --root e exatamente um de --dry-run ou --apply")
		return 2
	}
	v, err := resolveVault(*vault)
	if err != nil {
		return fail(d, err)
	}
	if !*includeGit {
		warnGitRoots(d, roots)
	}
	files, err := canon.HashTree(ctx, canon.ScanOptions{Roots: roots, Exclude: []string{v}, Workers: *workers, IncludeGit: *includeGit})
	if err != nil {
		return fail(d, err)
	}
	plan := canon.BuildPlan(files, v)
	if *dry {
		return reportPlan(d, plan, len(files), *asJSON)
	}
	return reportApply(d, canon.Apply(ctx, plan, canon.DefaultLinker()), *asJSON)
}

// warnGitRoots avisa (stderr) cada raiz que será pulada por estar num repositório git.
func warnGitRoots(d Deps, roots []string) {
	home, _ := os.UserHomeDir()
	for _, r := range roots {
		abs, err := filepath.Abs(r)
		if err != nil {
			continue
		}
		if repo, ok := canon.GitRoot(abs, home); ok {
			fmt.Fprintf(d.Stderr, "aviso: raiz %s ignorada (repositório git em %s); use --include-git para incluir\n", abs, repo)
		}
	}
}

func resolveVault(v string) (string, error) {
	if v == "" {
		return canon.DefaultVault()
	}
	return filepath.Abs(v)
}

func reportPlan(d Deps, p canon.Plan, scanned int, asJSON bool) int {
	if asJSON {
		return printJSON(d, p)
	}
	fmt.Fprintf(d.Stdout, "DRY-RUN — %d ativos varridos, %d grupos de duplicatas exatas\n\n", scanned, len(p.Groups))
	w := table(d.Stdout)
	fmt.Fprintln(w, "CANÔNICO\tCÓPIAS\tECONOMIA")
	for _, g := range p.Groups {
		fmt.Fprintf(w, "%s\t%d\t%s\n", g.Canonical.Path, len(g.Links), humanBytes(g.SavedBytes))
	}
	if err := w.Flush(); err != nil {
		return fail(d, err)
	}
	fmt.Fprintf(d.Stdout, "\nCofre: %s\neconomia estimada: %s\n", p.Vault, humanBytes(p.SavedBytes))
	return 0
}

func reportApply(d Deps, r canon.ApplyReport, asJSON bool) int {
	code := 0
	if len(r.Failed) > 0 {
		code = 1
	}
	if asJSON {
		if printJSON(d, r) != 0 {
			return 1
		}
		return code
	}
	for _, l := range r.Linked {
		fmt.Fprintf(d.Stdout, "%-8s %s\n", l.Method, l.Path)
	}
	for _, f := range r.Failed {
		fmt.Fprintf(d.Stdout, "PRESERVADO %s (%s)\n", f.Path, f.Error)
	}
	fmt.Fprintf(d.Stdout, "\n%d links, %d preservados, economia real: %s\n", len(r.Linked), len(r.Failed), humanBytes(r.SavedBytes))
	return code
}

func loadDAG(path string) (*dag.DAG, error) {
	data, err := readFileArg(path)
	if err != nil {
		return nil, err
	}
	return dag.Parse(data)
}

func cmdValidate(_ context.Context, args []string, d Deps) int {
	path, _ := splitPositional(args)
	p, err := loadDAG(path)
	if err != nil {
		return fail(d, err)
	}
	fmt.Fprintf(d.Stdout, "✓ %s v%s — %d nós válidos, sem ciclos\n", p.Name, p.Version, len(p.Nodes))
	for i, layer := range p.Layers() {
		fmt.Fprintf(d.Stdout, "  camada %d: %s\n", i+1, strings.Join(layer, ", "))
	}
	return 0
}

type runFlags struct {
	workdir, skills string
	parallel        int
	timeout         time.Duration
	vars            listFlag
	asJSON          bool
	allowScripts    bool
}

func parseRunFlags(args []string, d Deps) (runFlags, error) {
	var f runFlags
	fs := newFlags("run", d)
	fs.StringVar(&f.workdir, "workdir", ".", "diretório do projeto (handoffs)")
	fs.StringVar(&f.skills, "skills", "", "cofre de skills (padrão ~/.aidd/canonical)")
	fs.IntVar(&f.parallel, "parallel", 0, "nós simultâneos por camada (0 = NumCPU, 1 = sequencial)")
	fs.DurationVar(&f.timeout, "timeout", 30*time.Minute, "timeout padrão por nó")
	fs.Var(&f.vars, "var", "variável chave=valor (repetível)")
	fs.BoolVar(&f.asJSON, "json", false, "saída JSON")
	fs.BoolVar(&f.allowScripts, "allow-scripts", false, "permite nós script (executam comandos locais)")
	return f, fs.Parse(args)
}

func parseVars(kvs []string) (map[string]string, error) {
	out := map[string]string{}
	for _, kv := range kvs {
		k, v, ok := strings.Cut(kv, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("--var %q deve ser chave=valor", kv)
		}
		out[k] = v
	}
	return out, nil
}

func cmdRun(ctx context.Context, args []string, d Deps) int {
	path, rest := splitPositional(args)
	f, err := parseRunFlags(rest, d)
	if err != nil {
		return 2
	}
	p, err := loadDAG(path)
	if err != nil {
		return fail(d, err)
	}
	vars, err := parseVars(f.vars)
	if err != nil {
		return fail(d, err)
	}
	ox, err := newOrchestrix(ctx, d, f.skills, vars)
	if err != nil {
		return fail(d, err)
	}
	ox.AllowScripts = f.allowScripts
	rep, err := runner.Run(ctx, p, runner.Options{Workdir: f.workdir, Parallelism: f.parallel, DefaultTimeout: f.timeout, Executor: ox})
	if err != nil {
		return fail(d, err)
	}
	return reportRun(d, rep, f.asJSON)
}

func newOrchestrix(ctx context.Context, d Deps, skills string, vars map[string]string) (*orchestrix.Orchestrix, error) {
	dir, err := resolveVault(skills)
	if err != nil {
		return nil, err
	}
	cfg := provider.DefaultConfig()
	return &orchestrix.Orchestrix{
		Inventory: d.discover(ctx), HTTP: &http.Client{}, OllamaURL: cfg.OllamaURL,
		RouterKey: cfg.RouterKey, RouterKeyURL: cfg.TrustedRouterURL(), SkillsDir: dir, Vars: vars,
	}, nil
}

func reportRun(d Deps, rep runner.Report, asJSON bool) int {
	code := 0
	if !rep.OK {
		code = 1
	}
	if asJSON {
		if printJSON(d, rep) != 0 {
			return 1
		}
		return code
	}
	w := table(d.Stdout)
	fmt.Fprintln(w, "NÓ\tSTATUS\tMODELO\tDURAÇÃO\tTOKENS IN/OUT\tDETALHE")
	for _, n := range rep.Nodes {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d/%d\t%s\n", n.ID, n.Status, n.Model, n.Duration.Round(time.Millisecond), n.Tokens.Input, n.Tokens.Output, n.Error)
	}
	if err := w.Flush(); err != nil {
		return fail(d, err)
	}
	for _, n := range rep.Nodes {
		for _, note := range n.Notes {
			fmt.Fprintf(d.Stdout, "  ↳ %s: %s\n", n.ID, note)
		}
	}
	fmt.Fprintf(d.Stdout, "\nTotal de tokens: %d in / %d out\n", rep.Tokens.Input, rep.Tokens.Output)
	return code
}

func cmdRoute(_ context.Context, args []string, d Deps) int {
	fs := newFlags("route", d)
	from := fs.String("from", string(graph.Idea), "artefato de entrada")
	to := fs.String("to", string(graph.Code), "artefato objetivo")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	idx, err := graph.NewIndex(graph.FullCatalog())
	if err != nil {
		return fail(d, err)
	}
	path, cost, err := idx.Route(graph.ArtifactType(*from), graph.ArtifactType(*to))
	if err != nil {
		return fail(d, err)
	}
	ids := make([]string, len(path))
	for i, b := range path {
		ids[i] = b.ID
	}
	fmt.Fprintf(d.Stdout, "%s → %s\n  %s\n  custo estimado: %dk tokens\n", *from, *to, strings.Join(ids, " → "), cost)
	return 0
}

// assetScanInterval é a cadência da re-varredura do inventário (barata: só relê o que mudou).
const assetScanInterval = 30 * time.Second

func cmdStudio(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("studio", d)
	addr := fs.String("addr", "127.0.0.1:7878", "endereço de escuta (somente loopback)")
	project := fs.String("project", ".", "projeto onde pipeline.dag.json é gravado")
	var canonRoots listFlag
	fs.Var(&canonRoots, "canon-root", "raiz extra da canonização na interface (repetível)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	dir, err := filepath.Abs(*project)
	if err != nil {
		return fail(d, err)
	}
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return fail(d, errors.New("--project deve ser um diretório existente"))
	}
	for _, r := range canonRoots {
		if st, err := os.Stat(r); err != nil || !st.IsDir() {
			return fail(d, fmt.Errorf("--canon-root %q deve ser um diretório existente", r))
		}
	}
	idx, err := graph.NewIndex(graph.FullCatalog())
	if err != nil {
		return fail(d, err)
	}
	w := assets.NewWatcher(assets.NewScanner(assets.Config{Home: d.Home, Projects: []string{dir}, AutoDiscover: true}), assetScanInterval)
	w.Start(ctx)
	srv := &studio.Server{Index: idx, ProjectDir: dir, Discover: d.discover, Assets: w, Home: d.Home, CanonRoots: canonRoots}
	return serveStudio(ctx, srv, *addr, d)
}

func serveStudio(ctx context.Context, srv *studio.Server, addr string, d Deps) int {
	ready := make(chan string, 1)
	go func() {
		if a, ok := <-ready; ok {
			fmt.Fprintf(d.Stdout, "Looping Studio em http://%s (Ctrl+C para sair)\n", a)
			if d.Ready != nil {
				d.Ready <- a
			}
		}
	}()
	err := srv.ListenAndServe(ctx, addr, ready)
	close(ready)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fail(d, err)
	}
	return 0
}
