package fechar

import (
	"context"
	"errors"
	"fmt"
	"legoloop/internal/caminhos"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/brain"
	"legoloop/internal/processo"
)

var ErrDAGDivergente = errors.New("fechar: hash SHA-256 divergente entre docs/dags e a pasta de Processos")

// Deps encapsula as bordas e dependências externas do fechamento.
type Deps struct {
	RepoDir   string
	VaultDir  string
	Home      string
	SkipBuild bool
	Brain     brain.Brain
	Auditor   func(pdfPath string) (bool, error)
	Rebuilder func(dir string) error
	Indexer   func(repo string) error
	Portoes   func(repo string) error
	Clock     func() time.Time
}

func (d *Deps) initDefaults() {
	if d.RepoDir == "" {
		d.RepoDir, _ = os.Getwd()
	}
	if d.VaultDir == "" {
		d.VaultDir = filepath.Join(d.Home, ".aidd", "canonical")
	}
	if d.Clock == nil {
		d.Clock = time.Now
	}
}

// Executar executa as 9 etapas determinísticas de canonização e fechamento.
func Executar(ctx context.Context, req Request, deps Deps) (*Result, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	deps.initDefaults()
	slug, canon, _ := NormalizeSlug(req.Slug)
	start := deps.Clock()
	res := &Result{
		Slug:      slug,
		CanonName: canon,
		Sucesso:   true,
		Timestamp: start,
		Etapas:    make([]EtapaStatus, 0, 9),
	}
	if req.Telemetria != nil {
		res.Telemetria = *req.Telemetria
	}
	executarEtapas(ctx, &deps, slug, canon, req.DryRun, res)
	dur := deps.Clock().Sub(start)
	res.Duracao = dur.Milliseconds()
	if res.Duracao <= 0 {
		res.Duracao = 1
	}
	res.Telemetria.EnsureDefaults(dur)
	return res, nil
}

func executarEtapas(ctx context.Context, d *Deps, slug, canon string, dry bool, res *Result) {
	runEtapa(d, res, 1, "Manual Operacional", func() (string, error) {
		return fecharManual(d, slug)
	})
	runEtapa(d, res, 2, "Livro Didático LegoPress", func() (string, error) {
		return fecharLivro(d, slug, canon)
	})
	runEtapa(d, res, 3, "Skill Canônica Universal", func() (string, error) {
		return fecharSkill(d, canon)
	})
	runEtapa(d, res, 4, "Processo & Grafo DAG", func() (string, error) {
		return fecharDAG(d, slug, canon, dry)
	})
	runEtapaRebuild(d, res, dry, naturezaDoProcesso(d.VaultDir, canon, slug))
	runEtapa(d, res, 6, "Governança Viva (AGENTS e MEMORY)", func() (string, error) {
		return fecharGovernanca(d, slug, canon, dry)
	})
	runEtapa(d, res, 7, "Reindexação do Grafo Semântico", func() (string, error) {
		return fecharGrafo(d, dry)
	})
	runEtapa(d, res, 8, "Portões Mecânicos & Commit", func() (string, error) {
		return fecharPortoes(d, dry)
	})
	runEtapa(d, res, 9, "Relatório de Conclusão & Telemetria", func() (string, error) {
		return fecharReporte(d, slug, canon, res, dry)
	})
}

func runEtapa(d *Deps, res *Result, num int, nome string, fn func() (string, error)) {
	t0 := d.Clock()
	art, err := fn()
	st := "ok"
	det := ""
	if err != nil {
		st = "error"
		det = err.Error()
		res.Sucesso = false
	}
	res.Etapas = append(res.Etapas, EtapaStatus{
		Etapa:    num,
		Nome:     nome,
		Status:   st,
		Artefato: art,
		Duracao:  d.Clock().Sub(t0).Milliseconds(),
		Detalhe:  det,
	})
}

func fecharManual(d *Deps, slug string) (string, error) {
	return exigirManual(d.RepoDir, slug)
}

func fecharLivro(d *Deps, slug, canon string) (string, error) {
	dir := resolverDirLivro(d.RepoDir, slug, canon)
	if _, err := exigirLivro(d.RepoDir, dir); err != nil {
		return dir, err
	}
	if d.Auditor == nil {
		return dir, nil
	}
	pdf := filepath.Join(dir, fmt.Sprintf("%s.pdf", filepath.Base(dir)))
	if ok, err := d.Auditor(pdf); err != nil || !ok {
		return dir, fmt.Errorf("%w: auditoria do livro %s (ok=%v, err=%v)", ErrPortaoFalhou, relativo(d.RepoDir, pdf), ok, err)
	}
	return dir, nil
}

func fecharSkill(d *Deps, canon string) (string, error) {
	return exigirSkill(d.RepoDir, filepath.Join(d.RepoDir, "skills", canon, "SKILL.md"))
}

func fecharDAG(d *Deps, slug, canon string, dry bool) (string, error) {
	dir := resolverDirProcesso(d.RepoDir, slug, canon)
	dagPath := caminhos.CaminhoDAGLeitura(dir)
	if err := exigirNosDoDAG(d.RepoDir, dagPath); err != nil {
		return dagPath, err
	}
	if err := auditDAGParity(d.RepoDir, slug, dagPath); err != nil {
		return dagPath, err
	}
	if dry {
		return dagPath, nil
	}
	return dagPath, garantirQuarteto(dir, dagPath)
}

// resolverDirProcesso devolve a pasta existente do Processo nos dois dialetos,
// com preferência para o novo e para o nome canônico sobre o slug. Sem pasta
// existente, devolve a do dialeto antigo, que é o dialeto de gravação.
func resolverDirProcesso(repoDir, slug, canon string) string {
	nomes := []string{slug}
	if canon != "" {
		nomes = []string{canon, slug}
	}
	for _, base := range caminhos.DirsLeitura(repoDir) {
		for _, nome := range nomes {
			if dir := filepath.Join(base, nome); ehDiretorio(dir) {
				return dir
			}
		}
	}
	return filepath.Join(caminhos.DirLoopsEm(repoDir), nomes[0])
}

// nomeEspelho devolve o nome do arquivo espelho do DAG na pasta do Processo.
func nomeEspelho(nome string) string {
	return nome + caminhos.SufixoLoop
}

func auditDAGParity(repoDir, slug, procDAGPath string) error {
	dagsDir := filepath.Join(repoDir, "docs", "dags")
	pattern := filepath.Join(dagsDir, fmt.Sprintf("*_DAG_%s.json", upperSlug(slug)))
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		return nil
	}
	par, err := processo.CheckParity(matches[len(matches)-1], procDAGPath)
	if err != nil {
		return err
	}
	if !par {
		return ErrDAGDivergente
	}
	return nil
}

func upperSlug(slug string) string {
	s := strings.TrimPrefix(slug, "legoloop-")
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}

func fecharRebuild(d *Deps, dry bool) (string, error) {
	binPath := filepath.Join(d.RepoDir, "bin", "legoloop.exe")
	if d.SkipBuild || dry {
		return binPath, nil
	}
	rebuild := d.Rebuilder
	if rebuild == nil {
		rebuild = DefaultRebuilder
	}
	return binPath, rebuild(d.RepoDir)
}

// DefaultRebuilder recompila e instala o binário Go via exec sem shell.
func DefaultRebuilder(repoDir string) error {
	binPath := filepath.Join(repoDir, "bin", "legoloop.exe")
	if err := os.MkdirAll(filepath.Dir(binPath), 0o755); err != nil {
		return err
	}
	cmdBuild := exec.Command("go", "build", "-o", binPath, "./cmd/legoloop")
	cmdBuild.Dir = repoDir
	if out, err := cmdBuild.CombinedOutput(); err != nil {
		return fmt.Errorf("fechar: go build falhou: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	cmdInstall := exec.Command("go", "install", "./cmd/legoloop")
	cmdInstall.Dir = repoDir
	if out, err := cmdInstall.CombinedOutput(); err != nil {
		return fmt.Errorf("fechar: go install falhou: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func fecharGovernanca(d *Deps, slug, canon string, dry bool) (string, error) {
	agents := filepath.Join(d.RepoDir, "AGENTS.md")
	memory := filepath.Join(d.RepoDir, "MEMORY.md")
	if !dry {
		dataStr := d.Clock().Format("2006-01-02")
		if err := registrarGovernancaArquivo(agents, slug, canon, dataStr); err != nil {
			return "", err
		}
		if err := registrarGovernancaArquivo(memory, slug, canon, dataStr); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s, %s", filepath.Base(agents), filepath.Base(memory)), nil
}

func registrarGovernancaArquivo(caminho, slug, canon, dataStr string) error {
	atual, err := os.ReadFile(caminho)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	texto := string(atual)
	if strings.Contains(texto, canon) {
		return nil
	}
	linha := fmt.Sprintf("\n* **Processo Canonizado `%s` (%s):** `skills/%s/SKILL.md` · `docs/reports/`\n", canon, dataStr, canon)
	return os.WriteFile(caminho, []byte(texto+linha), 0o644)
}

func fecharGrafo(d *Deps, dry bool) (string, error) {
	graphPath := filepath.Join(d.RepoDir, ".codebase-memory", "graph.db.zst")
	if dry {
		return graphPath, nil
	}
	if d.Indexer != nil {
		return graphPath, d.Indexer(d.RepoDir)
	}
	if d.SkipBuild {
		return graphPath, nil
	}
	return graphPath, DefaultIndexer(d.RepoDir)
}

// DefaultIndexer reindexa o grafo semântico via codebase-memory-mcp cli index_repository.
func DefaultIndexer(repoDir string) error {
	return DefaultIndexerWithLookPath(repoDir, exec.LookPath, nil)
}

// DefaultIndexerWithLookPath executa a indexação do grafo com LookPath e executor injetáveis.
func DefaultIndexerWithLookPath(repoDir string, lookPath func(string) (string, error), run func(string, ...string) error) error {
	binPath, err := resolveCodebaseMemoryMCP(lookPath)
	if err != nil {
		return err
	}
	payload := fmt.Sprintf(`{"repo_path":%q}`, repoDir)
	if run != nil {
		return run(binPath, "cli", "index_repository", payload)
	}
	cmd := exec.Command(binPath, "cli", "index_repository", payload)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("fechar: index_repository falhou: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func resolveCodebaseMemoryMCP(lookPath func(string) (string, error)) (string, error) {
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	binPath, err := lookPath("codebase-memory-mcp")
	if err == nil {
		return binPath, nil
	}
	if home, homeErr := os.UserHomeDir(); homeErr == nil && home != "" {
		candidatos := []string{
			filepath.Join(home, "tools", "codebase-memory-mcp", "codebase-memory-mcp.exe"),
			filepath.Join(home, "tools", "codebase-memory-mcp", "codebase-memory-mcp"),
			filepath.Join(home, ".local", "bin", "codebase-memory-mcp.exe"),
			filepath.Join(home, ".local", "bin", "codebase-memory-mcp"),
		}
		for _, cand := range candidatos {
			if info, statErr := os.Stat(cand); statErr == nil && !info.IsDir() {
				return cand, nil
			}
		}
	}
	return "", fmt.Errorf("fechar: codebase-memory-mcp ausente no PATH: %w", err)
}

func fecharPortoes(d *Deps, dry bool) (string, error) {
	alvo := filepath.Join(d.RepoDir, "go.mod")
	if dry {
		return alvo, nil
	}
	if d.Portoes != nil {
		return alvo, d.Portoes(d.RepoDir)
	}
	if d.SkipBuild {
		return alvo, nil
	}
	return alvo, DefaultPortoes(d.RepoDir)
}

// DefaultPortoes executa go vet sobre o repositório via exec sem shell.
func DefaultPortoes(repoDir string) error {
	if _, err := os.Stat(filepath.Join(repoDir, "go.mod")); os.IsNotExist(err) {
		return nil
	}
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: go vet: %v (%s)", ErrPortaoFalhou, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func fecharReporte(d *Deps, slug, canon string, res *Result, dry bool) (string, error) {
	dir := filepath.Join(d.RepoDir, "docs", "reports")
	if !dry {
		_ = os.MkdirAll(dir, 0o755)
	}
	name := fmt.Sprintf("%s_REPORT_%s.md", d.Clock().Format("2006-01-02"), strings.ToUpper(strings.ReplaceAll(slug, "-", "_")))
	target := filepath.Join(dir, name)
	res.ReportPath = target
	if !dry {
		res.Telemetria.EnsureDefaults(time.Duration(res.Duracao) * time.Millisecond)
		md := GerarReportMD(slug, canon, res, d.Clock())
		if err := os.WriteFile(target, []byte(md), 0o644); err != nil {
			return target, err
		}
	}
	return target, nil
}
