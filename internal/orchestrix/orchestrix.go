// Package orchestrix implementa a Orchestrix Cirúrgica: cada nó nasce num
// processo/requisição efêmero, recebe um prompt único montado apenas a partir dos
// artefatos de entrada, grava seu handoff no disco e morre. Nenhum histórico é herdado.
package orchestrix

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
	"legoloop/internal/roteamento"
	"legoloop/internal/runner"
	"legoloop/internal/studiodoc"
)

// Orchestrix é o despachante. Implementa runner.Executor.
type Orchestrix struct {
	Inventory provider.Inventory
	HTTP      *http.Client
	OllamaURL string
	RouterKey string
	// RouterKeyURL é a única base do 9Router que recebe RouterKey.
	RouterKeyURL string
	SkillsDir    string
	Vars         map[string]string
	// ExtraEnv entra no ambiente estéril do subprocesso (além da allowlist).
	ExtraEnv []string
	// AllowScripts é o consentimento explícito (--allow-scripts) para nós script.
	AllowScripts bool
	// AllowUnitTasks é o consentimento explícito (--allow-unit-tasks) para nós unit_task.
	AllowUnitTasks bool
	// LoopStore é o acervo de sub-loops (pode ser nil caso o acervo não esteja carregado).
	LoopStore *studiodoc.Store
	// EstadoHarness devolve o estado de homologação do par harness e via (nil = nenhum harness homologado).
	EstadoHarness func(harness, via string) string
	// Elegiveis lista, por classe, os candidatos homologados do perfil, na ordem de preferência.
	// Nano e micro usam a lista para trocar de candidato quando o declarado falha.
	Elegiveis map[string][]roteamento.CandidatoVia
	// MemoriaLivre devolve a memória física livre em bytes (injetada; nil = sem limite por memória).
	MemoriaLivre func() (uint64, error)
	// PesoHarnessMB devolve o peso_mb homologado do par harness e via (0 = sem medição, vale 750 MB).
	PesoHarnessMB func(harness, via string) int
	// Aviso recebe os avisos de tela, como a Substituicao (nil = sem aviso).
	Aviso func(string)
	// MaxOutputBytes limita stdout de subprocessos (0 = maxResponseBytes).
	MaxOutputBytes int
}

var _ runner.Executor = (*Orchestrix)(nil)

// Execute despacha o nó conforme seu tipo.
func (o *Orchestrix) Execute(ctx context.Context, n dag.Node, workdir string) (runner.Result, error) {
	switch n.Type {
	case dag.TypeScript:
		return o.runScript(ctx, n, workdir)
	case dag.TypeSkill:
		return o.runSkill(ctx, n, workdir)
	case dag.TypeSubloop:
		return o.runSubloop(ctx, n, workdir)
	case dag.TypeUnitTask:
		return o.runUnitTask(ctx, n, workdir)
	case dag.TypeStopGate:
		// Barreira humana: o Execute não faz nada. O runner avalia o prompt do gate depois.
		return runner.Result{Model: "stop_gate/human"}, nil
	}
	return runner.Result{}, fmt.Errorf("orchestrix: tipo de nó não suportado %q", n.Type)
}

func (o *Orchestrix) runSkill(ctx context.Context, n dag.Node, workdir string) (runner.Result, error) {
	if n.Routing.Harness == harnessDireto {
		return o.runDireto(ctx, n, workdir)
	}
	root, err := os.OpenRoot(workdir)
	if err != nil {
		return runner.Result{}, err
	}
	defer root.Close()
	routes, err := o.routes(n.Routing)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	prompt, err := o.buildPrompt(n, root)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	var errs []error
	for _, rt := range routes {
		text, tokens, err := o.dispatch(ctx, rt.src, rt.model, prompt, n.Routing.Effort)
		if err == nil {
			notes := errStrings(errs)
			if n.Routing.Effort != "" {
				notes = append(notes, fmt.Sprintf("effort %s aplicado", n.Routing.Effort))
			}
			res := runner.Result{Tokens: tokens, Model: string(rt.src) + "/" + rt.model, Notes: notes}
			return res, writeArtifact(root, n.Output.Artifact, []byte(text))
		}
		errs = append(errs, fmt.Errorf("via %s/%s: %w", rt.src, rt.model, err))
		if ctx.Err() != nil {
			break // timeout/cancelamento não gasta o fallback
		}
	}
	return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, errors.Join(errs...))
}

// runDireto chama a API do 9Router sem harness e grava só o artefato de saída.
func (o *Orchestrix) runDireto(ctx context.Context, n dag.Node, workdir string) (runner.Result, error) {
	if n.Routing.Via != viaRouter {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: harness direto exige via %s", n.ID, viaRouter)
	}
	root, err := os.OpenRoot(workdir)
	if err != nil {
		return runner.Result{}, err
	}
	defer root.Close()
	prompt, err := o.buildPrompt(n, root)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	text, tokens, err := o.callRouter(ctx, n.Routing.Model, prompt)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, o.semChave(err))
	}
	res := runner.Result{Tokens: tokens, Model: harnessDireto + "/" + viaRouter + "/" + n.Routing.Model}
	return res, writeArtifact(root, n.Output.Artifact, []byte(text))
}

func errStrings(errs []error) []string {
	var out []string
	for _, e := range errs {
		out = append(out, "descartado "+e.Error())
	}
	return out
}

func (o *Orchestrix) dispatch(ctx context.Context, src provider.Source, model, prompt string, effort dag.EffortLevel) (string, runner.TokenUsage, error) {
	switch src {
	case provider.SourceRouter:
		return o.callRouter(ctx, model, prompt)
	case provider.SourceOllama:
		return o.callOllama(ctx, model, prompt)
	case provider.SourceClaudeCLI:
		return o.callClaude(ctx, model, prompt)
	default:
		return o.callHarness(ctx, src, model, prompt, effort)
	}
}

// route é um candidato executável (modelo presente no inventário).
type route struct {
	src   provider.Source
	model string
}

// routes aplica a Regra de Ouro e devolve os candidatos em ordem de tentativa:
// auto ⇒ só o combo do 9Router; manual ⇒ só o modelo travado; suggested ⇒ modelo
// e depois fallback (failover em tempo de execução). Nada fora do inventário roda.
func (o *Orchestrix) routes(r dag.RoutingConfig) ([]route, error) {
	if r.Mode == dag.ModeAuto {
		if src, ok := o.Inventory.Source(r.Combo); ok && src == provider.SourceRouter {
			return []route{{src, r.Combo}}, nil
		}
		return nil, fmt.Errorf("combo %q não está no inventário (9Router offline?)", r.Combo)
	}
	candidates := []string{r.Model}
	if r.Mode == dag.ModeSuggested && r.Fallback != r.Model {
		candidates = append(candidates, r.Fallback)
	}
	var out []route
	for _, m := range candidates {
		if src, ok := o.Inventory.Source(m); ok && m != "" {
			out = append(out, route{src, m})
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("nenhum modelo de %v está no inventário", candidates)
	}
	return out, nil
}

// writeArtifact grava via temporário + rename: troca a entrada de diretório, então um
// artefato que seja hardlink (canonização) nunca propaga a escrita para outros projetos.
// Diretórios-pai que sejam symlink/junction são recusados: um link commitado
// (ex.: out -> .git) desviaria o artefato para hooks/configs executáveis.
func writeArtifact(root *os.Root, name string, data []byte) error {
	native := filepath.FromSlash(name)
	if dir := filepath.Dir(native); dir != "." {
		if err := ensureRealDirs(root, dir); err != nil {
			return err
		}
	}
	tmp, tmpName, err := createExclusive(root, native)
	if err != nil {
		return err
	}
	_, err = tmp.Write(data)
	err = errors.Join(err, tmp.Close())
	if err == nil {
		err = root.Rename(tmpName, native)
	}
	if err != nil {
		return errors.Join(err, root.Remove(tmpName))
	}
	return nil
}

// ensureRealDirs percorre o caminho segmento a segmento: cria o que falta e
// recusa qualquer pai que seja link ANTES de criar algo através dele.
func ensureRealDirs(root *os.Root, dir string) error {
	cur := ""
	for _, part := range strings.Split(dir, string(filepath.Separator)) {
		cur = filepath.Join(cur, part)
		info, err := root.Lstat(cur)
		if errors.Is(err, os.ErrNotExist) {
			if err := root.Mkdir(cur, 0o755); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Mode()&(os.ModeSymlink|os.ModeIrregular) != 0 {
			return fmt.Errorf("diretório %q do artefato é link; gravação recusada", cur)
		}
	}
	return nil
}

// createExclusive cria um temporário com nome aleatório e O_EXCL (nunca segue link pré-existente).
func createExclusive(root *os.Root, native string) (*os.File, string, error) {
	var rnd [8]byte
	if _, err := rand.Read(rnd[:]); err != nil {
		return nil, "", err
	}
	name := native + "." + hex.EncodeToString(rnd[:]) + ".legoloop-tmp"
	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	return f, name, err
}

func artifactExists(root *os.Root, name string) bool {
	_, err := root.Stat(filepath.FromSlash(name))
	return !errors.Is(err, os.ErrNotExist)
}

// Ask despacha um prompt avulso (sem nó de DAG) a um executor do inventário e devolve o texto. O cérebro
// plugável (internal/brain) usa esta porta para falar com Ollama, claude CLI e harnesses.
func (o *Orchestrix) Ask(ctx context.Context, src provider.Source, model, prompt string) (string, error) {
	text, _, err := o.dispatch(ctx, src, model, prompt, "")
	return text, err
}
