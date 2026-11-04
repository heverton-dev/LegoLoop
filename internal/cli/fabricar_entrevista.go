package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"legoloop/internal/brain"
	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/ste"
)

func init() {
	for id, msg := range map[string]string{
		"entrevista.uso":           "Uso: legoloop fabricar entrevista iniciar \"<ideia>\" ou legoloop fabricar entrevista responder <sessao> \"<texto>\".",
		"entrevista.opcao_palpite": "Aceitar o palpite: %s",
		"entrevista.opcao_livre":   "Responder com as suas palavras",
	} {
		ste.Default.Register(id, msg)
	}
}

// espelhoJSON é o espelho canônico com chaves em português, para o Chat ler sem tradução.
type espelhoJSON struct {
	Resultado    string   `json:"resultado_esperado"`
	Quem         string   `json:"quem_se_beneficia"`
	PorQueAgora  string   `json:"por_que_agora"`
	Criterio     string   `json:"criterio_de_sucesso"`
	Restricoes   []string `json:"restricoes_inegociaveis"`
	ForaDeEscopo []string `json:"fora_de_escopo"`
}

// passoJSON é a saída de cada passo da entrevista. O Chat lê este JSON e faz a pergunta ao usuário.
type passoJSON struct {
	Sessao    string       `json:"sessao"`
	Status    string       `json:"status"`
	Pergunta  string       `json:"pergunta"`
	Palpite   string       `json:"palpite"`
	Opcoes    []string     `json:"opcoes"`
	Confianca int          `json:"confianca"`
	Meta      int          `json:"meta"`
	Turnos    int          `json:"turnos"`
	Espelho   *espelhoJSON `json:"espelho"`
	Lacunas   []string     `json:"lacunas,omitempty"`
}

type entrevistaOpts struct{ vault, docs, model string }

func newEntrevistaFlags(o *entrevistaOpts) *flag.FlagSet {
	fs := flag.NewFlagSet("fabricar entrevista", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&o.vault, "vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	fs.StringVar(&o.docs, "docs", "", "raiz do projeto, de onde vem o contexto (padrão: diretório atual)")
	fs.StringVar(&o.model, "model", "", "modelo que conduz a entrevista (padrão: o primeiro listado por legoloop doctor)")
	return fs
}

// parseInterspersed lê flags e argumentos em qualquer ordem e devolve os argumentos soltos.
func parseInterspersed(fs *flag.FlagSet, args []string) ([]string, error) {
	var positional []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		if fs.NArg() == 0 {
			return positional, nil
		}
		positional = append(positional, fs.Arg(0))
		args = fs.Args()[1:]
	}
}

// cmdEntrevista conduz a entrevista passo a passo: cada chamada faz um passo e guarda a sessão em disco.
func cmdEntrevista(ctx context.Context, args []string, d Deps) int {
	var o entrevistaOpts
	pos, err := parseInterspersed(newEntrevistaFlags(&o), args)
	if err != nil || len(pos) == 0 {
		return usoEntrevista(d, err)
	}
	switch verb, rest := pos[0], pos[1:]; {
	case verb == "iniciar" && len(rest) == 1:
		return runIniciar(ctx, rest[0], o, d)
	case verb == "responder" && len(rest) == 2:
		return runResponder(ctx, rest[0], rest[1], o, d)
	}
	return usoEntrevista(d, nil)
}

func usoEntrevista(d Deps, err error) int {
	if err != nil {
		fmt.Fprintln(d.Stderr, "erro:", err)
	}
	fmt.Fprintln(d.Stderr, ste.Msg("entrevista.uso"))
	return 2
}

// entrevistaEnv resolve o cérebro e a pasta das sessões. O LLM injetado em testes vale mais que o inventário.
func entrevistaEnv(ctx context.Context, o entrevistaOpts, d Deps) (interview.LLMClient, interview.Store, error) {
	roots, err := fabricarRoots(fabricarOpts{vault: o.vault, docs: o.docs}, d)
	if err != nil {
		return nil, interview.Store{}, err
	}
	store := interview.Store{Dir: filepath.Join(roots.Vault, "sessions")}
	if d.Fabricar != nil && d.Fabricar.LLM != nil {
		return d.Fabricar.LLM, store, nil
	}
	b, why, err := brain.Pick(d.discover(ctx), o.model, d.brainOpts()...)
	if err != nil {
		return nil, store, err
	}
	fmt.Fprintln(d.Stderr, why)
	return b, store, nil
}

func runIniciar(ctx context.Context, idea string, o entrevistaOpts, d Deps) int {
	llm, store, err := entrevistaEnv(ctx, o, d)
	if err != nil {
		return failBrain(d, err)
	}
	docs := o.docs
	if docs == "" {
		docs = "."
	}
	eng := interview.NewEngine(llm).WithContext(interview.Context{Goal: idea, Docs: interview.LoadProjectContext(docs)})
	sess, turn, err := eng.Start(ctx, idea)
	if err != nil {
		return fail(d, err)
	}
	return gravaEMostra(d, store, sess, turn)
}

func runResponder(ctx context.Context, id, text string, o entrevistaOpts, d Deps) int {
	llm, store, err := entrevistaEnv(ctx, o, d)
	if err != nil {
		return failBrain(d, err)
	}
	sess, err := store.Load(id)
	if err != nil {
		return failSessao(d, err)
	}
	eng := interview.NewEngine(llm)
	eng.Adopt(sess)
	turn, err := eng.Step(ctx, sess.ID, text)
	if err != nil {
		return fail(d, err)
	}
	return gravaEMostra(d, store, sess, turn)
}

// failSessao: identificador malformado é uso incorreto (2); sessão ausente ou ilegível é falha (1).
func failSessao(d Deps, err error) int {
	if errors.Is(err, interview.ErrSessaoInvalida) {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	return fail(d, err)
}

func gravaEMostra(d Deps, store interview.Store, sess *interview.Session, turn *interview.TurnResult) int {
	if err := store.Save(sess); err != nil {
		return fail(d, err)
	}
	return printJSON(d, montaPasso(sess, turn))
}

func montaPasso(sess *interview.Session, turn *interview.TurnResult) passoJSON {
	clar := fabricar.ClarificationFromSession(sess)
	p := passoJSON{Sessao: sess.ID, Status: string(sess.Status), Pergunta: turn.Question, Palpite: turn.Guess,
		Opcoes: opcoes(turn.Guess), Confianca: sess.Confidence, Meta: fabricar.ConfidenceTarget, Turnos: sess.Turns, Lacunas: clar.Gaps}
	if turn.Mirror != nil {
		m := turn.Mirror
		p.Espelho = &espelhoJSON{m.Outcome, m.User, m.WhyNow, m.Success, m.Constraints, m.OutOfScope}
	}
	return p
}

// opcoes oferece aceitar o palpite (se houver) ou responder com palavras próprias.
func opcoes(guess string) []string {
	var out []string
	if guess != "" {
		out = append(out, ste.Msg("entrevista.opcao_palpite", guess))
	}
	return append(out, ste.Msg("entrevista.opcao_livre"))
}
