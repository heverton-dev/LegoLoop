package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"legoloop/internal/caminhos"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/canon"
	"legoloop/internal/dag"
	"legoloop/internal/gate"
	"legoloop/internal/graph"
	"legoloop/internal/loop"
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
	return finalizarCanonize(ctx, d, v, canon.BuildPlan(files, v), len(files), *dry, *asJSON)
}

// finalizarCanonize mostra o plano (--dry-run) ou aplica os links e, nos dois casos, trata a matriz de peças:
// a prévia só informa quantas peças seriam recatalogadas; o --apply gera a matriz ao final.
func finalizarCanonize(ctx context.Context, d Deps, vault string, plan canon.Plan, scanned int, dry, asJSON bool) int {
	if dry {
		code := reportPlan(d, plan, scanned, asJSON)
		matrizPreviaCanonize(ctx, d, vault, asJSON)
		return code
	}
	code := reportApply(d, canon.Apply(ctx, plan, canon.DefaultLinker()), asJSON)
	matrizAposCanonize(ctx, d, vault, asJSON)
	return code
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
	allowUnitTasks  bool
	autoApprove     bool
	nonInteractive  bool
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
	fs.BoolVar(&f.allowUnitTasks, "allow-unit-tasks", false, "permite nós unit_task (claude com ferramentas em worktree git)")
	fs.BoolVar(&f.autoApprove, "yes", false, "aprova automaticamente Stop-Gates interativos")
	fs.BoolVar(&f.autoApprove, "y", false, "alias de --yes")
	fs.BoolVar(&f.nonInteractive, "non-interactive", false, "rejeita Stop-Gates interativos se requerer interação")
	return f, fs.Parse(args)
}

type terminalPrompter struct {
	reader *bufio.Reader
	out    io.Writer
}

func (tp *terminalPrompter) Prompt(msg string) (bool, error) {
	fmt.Fprintf(tp.out, "\n🛑 [STOP-GATE INTERATIVO] %s\nConfirmar execução? [S/n]: ", msg)
	ans, err := tp.reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}
	ans = strings.TrimSpace(strings.ToLower(ans))
	if ans == "" || ans == "s" || ans == "sim" || ans == "y" || ans == "yes" || ans == "a" {
		return true, nil
	}
	return false, nil
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
	ox.AllowUnitTasks = f.allowUnitTasks
	rep, err := runner.Run(ctx, p, f.runnerOptions(ox, f.prompter(d)))
	if err != nil {
		return fail(d, err)
	}
	return reportRun(d, rep, f.asJSON)
}

// prompter devolve o prompt de aprovação no terminal; nil quando a aprovação é automática ou não interativa.
func (f runFlags) prompter(d Deps) gate.Prompter {
	if f.autoApprove || f.nonInteractive {
		return nil
	}
	in := d.Stdin
	if in == nil {
		in = os.Stdin
	}
	return &terminalPrompter{reader: bufio.NewReader(in), out: d.Stdout}
}

func (f runFlags) runnerOptions(ox *orchestrix.Orchestrix, prompter gate.Prompter) runner.Options {
	return runner.Options{
		Workdir:        f.workdir,
		Parallelism:    f.parallel,
		DefaultTimeout: f.timeout,
		Executor:       ox,
		AllowScripts:   f.allowScripts,
		AutoApprove:    f.autoApprove,
		NonInteractive: f.nonInteractive,
		Prompter:       prompter,
		Vars:           ox.Vars,
	}
}

func newOrchestrix(ctx context.Context, d Deps, skills string, vars map[string]string) (*orchestrix.Orchestrix, error) {
	dir, err := resolveVault(skills)
	if err != nil {
		return nil, err
	}
	cfg := provider.DefaultConfig()
	ox := &orchestrix.Orchestrix{
		Inventory: d.discover(ctx), HTTP: &http.Client{}, OllamaURL: cfg.OllamaURL,
		RouterKey: cfg.RouterKey, RouterKeyURL: cfg.TrustedRouterURL(), SkillsDir: dir, Vars: vars,
	}
	// O acervo de loops do cofre executa os nós subloop (etapa que compõe outro processo). Só liga se a pasta
	// existe: o run não cria pasta no cofre.
	if loops := caminhos.DirLoopsEm(dir); dirExists(loops) {
		if ox.LoopStore, err = loop.NewStore(loops); err != nil {
			return nil, err
		}
	}
	return ox, nil
}

func dirExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
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
	project := fs.String("project", ".", "projeto onde "+caminhos.ArquivoDAG+" é gravado")
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

func cmdForge(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("forge", d)
	addr := fs.String("addr", "127.0.0.1:7878", "endereço de escuta (somente loopback)")
	project := fs.String("project", ".", "projeto onde "+caminhos.ArquivoDAG+" é gravado")
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
	// TICK-1203: Dispara a primeira varredura antes de servir
	w.Refresh(ctx)
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

// cmdImport importa uma skill de repositório GitHub para a pasta global e calcula SHA-256.
func cmdImport(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("import", d)
	name := fs.String("name", "", "nome customizado para a skill")
	target := fs.String("target", "", "pasta destino (padrão: ~/.gemini/config ou cofre canônico)")
	asJSON := fs.Bool("json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Fprintln(d.Stderr, "erro: informe a URL do repositório/skill do GitHub")
		return 2
	}
	rawURL := remaining[0]

	destDir := importDest(*target, d.Home)

	res, err := canon.ImportSkillFromURL(ctx, rawURL, *name, destDir)
	if err != nil {
		return fail(d, err)
	}

	return reportImport(d, res, *asJSON)
}

// importDest resolve a pasta destino da importação: --target, senão ~/.gemini/config (se existir), senão o cofre canônico.
func importDest(target, home string) string {
	if target != "" {
		return target
	}
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	if home == "" {
		return ""
	}
	geminiConfig := filepath.Join(home, ".gemini", "config")
	if fi, err := os.Stat(geminiConfig); err == nil && fi.IsDir() {
		return geminiConfig
	}
	return filepath.Join(home, ".aidd", "canonical")
}

func reportImport(d Deps, res canon.ImportResult, asJSON bool) int {
	if asJSON {
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Fprintln(d.Stdout, string(data))
		return 0
	}
	fmt.Fprintf(d.Stdout, "Skill %q importada e canonizada com sucesso!\n", res.Name)
	fmt.Fprintf(d.Stdout, "Destino: %s\n", res.TargetDir)
	fmt.Fprintf(d.Stdout, "SHA-256: %s\n", res.SHA256)
	fmt.Fprintf(d.Stdout, "Arquivos: %d\n", res.FilesCount)
	return 0
}
