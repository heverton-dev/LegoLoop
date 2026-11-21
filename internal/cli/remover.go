package cli

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"legoloop/internal/jail"
	"legoloop/internal/processo"
	"legoloop/internal/remover"
	"legoloop/internal/ste"
)

func init() {
	ste.Default.Register("remover.uso", "Uso: legoloop remover <nome> [--vault DIR] [--dry-run] [--force] [--yes] [--json].")
	ste.Default.Register("remover.pergunta", "Remover o processo? [s/N]: ")
	ste.Default.Register("remover.sem_terminal", "A confirmação exige um terminal. Use --yes para confirmar sem perguntar.")
}

type removerOpcoes struct {
	vault     string
	workspace bool
	yes       bool
	dryRun    bool
	force     bool
	asJSON    bool
}

type removerResultadoJSON struct {
	Nome      string   `json:"nome"`
	Versao    string   `json:"versao"`
	Harnesses []string `json:"harnesses"`
	Links     []string `json:"links"`
	Arquivos  []string `json:"arquivos"`
	DryRun    bool     `json:"dry_run"`
	Removido  bool     `json:"removido"`
}

// cmdRemover apaga um processo do cofre e os vínculos dele. Mostra o cartão e pede confirmação; --yes/--force confirma.
func cmdRemover(ctx context.Context, args []string, d Deps) int {
	nome, op, code := parseRemoverFlags(args, d)
	if code != 0 {
		return code
	}
	if err := jail.ValidarNomeProcesso(nome); err != nil {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	rd, err := montarRemoverDeps(op, d)
	if err != nil {
		return fail(d, err)
	}
	plano, err := remover.Planejar(ctx, rd, nome)
	if err != nil {
		if op.force && errors.Is(err, remover.ErrEmUso) {
			plano, err = planejarRemoverForcado(rd, nome)
		}
		if err != nil {
			return failRemover(d, err)
		}
	}
	return executarPlanoRemover(ctx, d, plano, op)
}

// failRemover mapeia erro de uso para exit 2 e o resto para exit 1.
func failRemover(d Deps, err error) int {
	if errors.Is(err, remover.ErrNomeInvalido) {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	return fail(d, err)
}

func parseRemoverFlags(args []string, d Deps) (string, removerOpcoes, int) {
	nome, rest := splitPositional(args)
	fs := newFlags("remover", d)
	var op removerOpcoes
	fs.StringVar(&op.vault, "vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	fs.BoolVar(&op.workspace, "workspace", false, "apaga também os artefatos locais do workspace")
	fs.BoolVar(&op.yes, "yes", false, "confirma sem perguntar")
	fs.BoolVar(&op.dryRun, "dry-run", false, "simula a remoção sem apagar arquivos")
	fs.BoolVar(&op.force, "force", false, "força a remoção mesmo em uso")
	fs.BoolVar(&op.asJSON, "json", false, "saída em JSON")
	if err := fs.Parse(rest); err != nil {
		return "", op, 2
	}
	if nome == "" && fs.NArg() == 1 {
		nome = fs.Arg(0)
	} else if nome == "" || fs.NArg() > 0 {
		fmt.Fprintln(d.Stderr, "erro:", ste.Msg("remover.uso"))
		return "", op, 2
	}
	return nome, op, 0
}

func montarRemoverDeps(op removerOpcoes, d Deps) (remover.Deps, error) {
	home := d.Home
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	v, _, err := vaultListar(op.vault, d.Home)
	if err != nil {
		return remover.Deps{}, err
	}
	var wsDir string
	if op.workspace {
		wsDir, _ = os.Getwd()
	}
	return remover.Deps{Store: processo.Store{Roots: processo.Roots{Vault: v}}, Home: home, Workspace: wsDir}, nil
}

func planejarRemoverForcado(rd remover.Deps, nome string) (remover.Plano, error) {
	vault := rd.Store.Roots.Vault
	procs, err := processo.ProcessosDoCofre(vault)
	if err != nil {
		return remover.Plano{}, err
	}
	var alvo *processo.ProcessoCofre
	for i := range procs {
		if procs[i].Nome == nome {
			alvo = &procs[i]
			break
		}
	}
	if alvo == nil {
		return remover.Plano{}, remover.ErrNaoEncontrado
	}
	return construirPlanoForcado(rd, *alvo), nil
}

func construirPlanoForcado(rd remover.Deps, it processo.ProcessoCofre) remover.Plano {
	vault := rd.Store.Roots.Vault
	fsys := processo.OSFS()
	hs := processo.DetectarHarnesses(fsys, rd.Home)
	vinculados := processo.HarnessesVinculados(fsys, hs, it.Nome)
	p := remover.Plano{Nome: it.Nome, Versao: it.Versao, Harnesses: vinculados, Links: []string{}, Arquivos: []string{}}
	for _, h := range hs {
		if slices.Contains(vinculados, h.Nome) {
			p.Links = append(p.Links, filepath.Join(h.Skills, it.Nome))
		}
	}
	procJSON := processo.CaminhoNoCofre(vault, it.Nome)
	_ = filepath.WalkDir(vault, func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == vault || path == procJSON {
			return nil
		}
		base := filepath.Base(path)
		if base == it.Nome || strings.HasPrefix(base, it.Nome+".") {
			p.Arquivos = append(p.Arquivos, path)
			if d.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	if _, err := os.Lstat(procJSON); err == nil {
		p.Arquivos = append(p.Arquivos, procJSON)
	}
	return p
}

func executarPlanoRemover(ctx context.Context, d Deps, plano remover.Plano, op removerOpcoes) int {
	if op.dryRun {
		if op.asJSON {
			return printJSON(d, removerResultadoJSON{
				Nome: plano.Nome, Versao: plano.Versao, Harnesses: plano.Harnesses,
				Links: plano.Links, Arquivos: plano.Arquivos, DryRun: true, Removido: false,
			})
		}
		fmt.Fprintf(d.Stdout, "\n%s[dry-run] Nenhum arquivo foi apagado.\n", plano.Cartao.Texto())
		return 0
	}
	if !op.asJSON {
		fmt.Fprintf(d.Stdout, "\n%s", plano.Cartao.Texto())
	}
	if !op.yes && !op.force {
		if code := pedirConfirmacaoRemover(d); code != 0 {
			return code
		}
	}
	if err := remover.Executar(ctx, plano); err != nil {
		return fail(d, err)
	}
	if op.asJSON {
		return printJSON(d, removerResultadoJSON{
			Nome: plano.Nome, Versao: plano.Versao, Harnesses: plano.Harnesses,
			Links: plano.Links, Arquivos: plano.Arquivos, DryRun: false, Removido: true,
		})
	}
	fmt.Fprintln(d.Stdout, remover.MensagemFeito(plano))
	return 0
}

func pedirConfirmacaoRemover(d Deps) int {
	if d.Stdin == nil {
		return fail(d, errors.New(ste.Msg("remover.sem_terminal")))
	}
	ok, err := newPrompter(d.Stdin, d.Stdout).yesNo(ste.Msg("remover.pergunta"))
	if err != nil {
		return fail(d, err)
	}
	if !ok {
		fmt.Fprintln(d.Stderr, remover.MensagemCancelado())
		return 1
	}
	return 0
}
