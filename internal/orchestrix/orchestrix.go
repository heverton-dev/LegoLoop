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
	"time"

	"legoloop/internal/combo"
	"legoloop/internal/dag"
	"legoloop/internal/provider"
	"legoloop/internal/roteamento"
	"legoloop/internal/runner"
	"legoloop/internal/studiodoc"
)

// AntibodyStore persiste e consulta lições de falhas anteriores sem acoplamento circular.
type AntibodyStore interface {
	Record(pkg, errSig, rootCause, diffSol string) error
	InjectPreventiveContext(pkg string) string
}

// SpeculativeAttempt executa uma tentativa especulativa isolada.
type SpeculativeAttempt func(ctx context.Context, prompt string) (string, runner.TokenUsage, error)

// SpeculativeCandidate descreve um candidato rápido ou pesado para hot-swap especulativo.
type SpeculativeCandidate struct {
	Name    string
	Timeout time.Duration
	Run     SpeculativeAttempt
}

// SpeculativeRunner encapsula ExecuteSpeculative para failover com contexto negativo.
type SpeculativeRunner func(ctx context.Context, fast, heavy SpeculativeCandidate, prompt string) (string, string, runner.TokenUsage, error)

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
	// Reprovados guarda os pares reprovados em execução (nil = sem registro).
	Reprovados *RegistroReprovados
	// AoReprovar persiste a reprovação de um par (nil = nada a persistir).
	AoReprovar func(harness, via, motivo string)
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
	// DefaultTimeout define o timeout herdado por sub-processos e corrida especulativa.
	DefaultTimeout time.Duration
	// AutoApprove e NonInteractive propagam flags para sub-processos.
	AutoApprove    bool
	NonInteractive bool
	// Antibodies armazena e injeta anticorpos preventivos no prompt dos nós.
	Antibodies AntibodyStore
	// Speculative executa hot-swap especulativo em ModeAuto/ModeSuggested.
	Speculative SpeculativeRunner
}

var _ runner.Executor = (*Orchestrix)(nil)

// SetAcervo configura o acervo de sub-processos sem expor termos legados.
func (o *Orchestrix) SetAcervo(s *studiodoc.Store) {
	o.LoopStore = s
}

// PermitirUnitNos ativa consentimento explícito para nós unit_no.
func (o *Orchestrix) PermitirUnitNos() {
	o.AllowUnitTasks = true
}

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
	case dag.TypeMCP:
		return o.runMCP(ctx, n, workdir)
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
	prompt, err := o.prepareSkillPrompt(n, root)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	if err := validateNodeCombo(n, prompt); err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	text, model, tokens, errs, execErr := o.executeSkillRoutes(ctx, n, routes, prompt)
	if execErr != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, execErr)
	}
	return buildSkillResult(root, n, model, text, tokens, errs)
}

func (o *Orchestrix) prepareSkillPrompt(n dag.Node, root *os.Root) (string, error) {
	prompt, err := o.buildPrompt(n, root)
	if err != nil {
		return "", err
	}
	if o.Antibodies != nil {
		if advice := o.Antibodies.InjectPreventiveContext(antibodyScope(n)); advice != "" {
			prompt = advice + prompt
		}
	}
	return prompt, nil
}

func buildSkillResult(root *os.Root, n dag.Node, model, text string, tokens runner.TokenUsage, errs []error) (runner.Result, error) {
	notes := errStrings(errs)
	if n.Routing.Effort != "" {
		notes = append(notes, fmt.Sprintf("effort %s aplicado", n.Routing.Effort))
	}
	res := runner.Result{Tokens: tokens, Model: model, Notes: notes}
	return res, writeArtifact(root, n.Output.Artifact, []byte(text))
}

func (o *Orchestrix) executeSkillRoutes(ctx context.Context, n dag.Node, routes []route, prompt string) (string, string, runner.TokenUsage, []error, error) {
	if len(routes) >= 2 && (n.Routing.Mode == dag.ModeAuto || n.Routing.Mode == dag.ModeSuggested) {
		return o.executeSpeculativeRoutes(ctx, n, routes, prompt)
	}
	return o.executeRoutesSequential(ctx, n, routes, prompt)
}

func (o *Orchestrix) executeRoutesSequential(ctx context.Context, n dag.Node, routes []route, prompt string) (string, string, runner.TokenUsage, []error, error) {
	var errs []error
	curPrompt := prompt
	for _, rt := range routes {
		text, tokens, err := o.dispatch(ctx, rt.src, rt.model, curPrompt, n.Routing.Effort)
		if err == nil {
			return text, string(rt.src) + "/" + rt.model, tokens, errs, nil
		}
		o.recordAntibody(n, rt.model, err)
		errs = append(errs, fmt.Errorf("via %s/%s: %w", rt.src, rt.model, err))
		if ctx.Err() != nil {
			break
		}
		curPrompt = enrichNegativeContext(prompt, err)
	}
	return "", "", runner.TokenUsage{}, errs, errors.Join(errs...)
}

func (o *Orchestrix) executeSpeculativeRoutes(ctx context.Context, n dag.Node, routes []route, prompt string) (string, string, runner.TokenUsage, []error, error) {
	if o.Speculative == nil {
		return o.executeRoutesSequential(ctx, n, routes, prompt)
	}
	var errs []error
	fastRt, heavyRt := routes[0], routes[1]
	fast := o.makeCandidate(n, fastRt, &errs)
	heavy := o.makeCandidate(n, heavyRt, &errs)
	text, usedName, tokens, err := o.Speculative(ctx, fast, heavy, prompt)
	if err == nil {
		return text, usedName, tokens, errs, nil
	}
	if len(routes) > 2 && ctx.Err() == nil {
		t2, m2, tok2, remErrs, seqErr := o.executeRoutesSequential(ctx, n, routes[2:], enrichNegativeContext(prompt, err))
		errs = append(errs, remErrs...)
		if seqErr == nil {
			return t2, m2, tok2, errs, nil
		}
	}
	return "", "", runner.TokenUsage{}, errs, errors.Join(errs...)
}

func (o *Orchestrix) makeCandidate(n dag.Node, rt route, errs *[]error) SpeculativeCandidate {
	fullName := string(rt.src) + "/" + rt.model
	return SpeculativeCandidate{
		Name:    fullName,
		Timeout: speculativeFastTimeout(o.DefaultTimeout),
		Run: func(ctx context.Context, p string) (string, runner.TokenUsage, error) {
			txt, tok, err := o.dispatch(ctx, rt.src, rt.model, p, n.Routing.Effort)
			if err != nil {
				o.recordAntibody(n, rt.model, err)
				*errs = append(*errs, fmt.Errorf("via %s: %w", fullName, err))
			}
			return txt, tok, err
		},
	}
}

func speculativeFastTimeout(d time.Duration) time.Duration {
	if d > 0 {
		return d
	}
	return 15 * time.Second
}

func enrichNegativeContext(prompt string, err error) string {
	return prompt + "\n\n[CONTEXTO NEGATIVO - EVITAR ABORDAGEM ANTERIOR QUE FALHOU]\n" + err.Error()
}

func (o *Orchestrix) recordAntibody(n dag.Node, model string, err error) {
	if o.Antibodies == nil || err == nil {
		return
	}
	cause := fmt.Sprintf("Falha ao executar %s: %s", model, trimStderr([]byte(err.Error())))
	_ = o.Antibodies.Record(antibodyScope(n), err.Error(), cause, "Aplicar hot-swap e contexto preventivo")
}

func antibodyScope(n dag.Node) string {
	if strings.TrimSpace(n.Ref) != "" {
		return n.Ref
	}
	return n.ID
}

func validateNodeCombo(n dag.Node, prompt string) error {
	raw := strings.TrimSpace(n.Routing.Combo)
	if raw == "" && n.Routing.Mode != dag.ModeAuto {
		raw = strings.TrimSpace(n.Routing.Model)
	}
	if raw == "" {
		return nil
	}
	mapped := combo.MapearComboLegado(raw)
	if _, err := combo.ObterCombo(mapped); err != nil {
		return nil
	}
	estTokens := len(prompt) / 4
	return combo.ValidarRequisicao(mapped, estTokens, nodeRequiresTools(n))
}

func nodeRequiresTools(n dag.Node) bool {
	if n.Tooling == nil {
		return false
	}
	if v, ok := n.Tooling["require_tools"].(bool); ok && v {
		return true
	}
	if v, ok := n.Tooling["tools_required"].(bool); ok && v {
		return true
	}
	return false
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
	prompt, err := o.prepareSkillPrompt(n, root)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	if err := validateNodeCombo(n, prompt); err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	text, tokens, err := o.callRouter(ctx, n.Routing.Model, prompt)
	if err != nil {
		o.recordAntibody(n, n.Routing.Model, err)
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, o.semChave(err))
	}
	res := runner.Result{Tokens: tokens, Model: harnessDireto + "/" + viaRouter + "/" + n.Routing.Model, Harness: harnessDireto, Via: viaRouter}
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
		src, ok := o.Inventory.Source(r.Combo)
		if !ok || src != provider.SourceRouter {
			return nil, fmt.Errorf("combo %q não está no inventário (9Router offline?)", r.Combo)
		}
		out := []route{{src, r.Combo}}
		for _, fb := range []string{r.Fallback, r.Model, "t-raciocinio-128k", "t-chat-128k"} {
			if fb != "" && fb != r.Combo {
				if fbSrc, fbOK := o.Inventory.Source(fb); fbOK {
					out = append(out, route{fbSrc, fb})
					break
				}
			}
		}
		return out, nil
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
