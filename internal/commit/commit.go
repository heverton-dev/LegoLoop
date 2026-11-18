package commit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Opcoes encapsula as opções de execução do commit automático.
type Opcoes struct {
	Mensagem string
	DryRun   bool
	NoPush   bool
	Verbose  bool
}

// Executar roda o ciclo de commit determinístico.
func Executar(ctx context.Context, stdout, stderr io.Writer, opt Opcoes) error {
	if err := checarRepoGit(ctx); err != nil {
		return err
	}
	if err := prepararStage(ctx); err != nil {
		return err
	}
	staged, err := obterArquivosStaged(ctx)
	if err != nil {
		return err
	}
	if len(staged) == 0 {
		fmt.Fprintln(stdout, "Nenhuma alteração pendente para commit.")
		return nil
	}
	msg := opt.Mensagem
	if msg == "" {
		msg = gerarMensagemFallback(staged)
	}
	if opt.DryRun {
		fmt.Fprintf(stdout, "[dry-run] Commit previsto (%d arquivos): %s\n", len(staged), msg)
		return nil
	}
	if err := rodarCommit(ctx, msg); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Commit realizado com sucesso: %s\n", msg)
	if opt.NoPush {
		return nil
	}
	return rodarPush(ctx, stdout)
}

func checarRepoGit(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--is-inside-work-tree")
	if out, err := cmd.CombinedOutput(); err != nil || strings.TrimSpace(string(out)) != "true" {
		return errors.New("diretório atual não é um repositório git válido")
	}
	return nil
}

func prepararStage(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "add", "-A")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("falha ao adicionar arquivos ao stage: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func obterArquivosStaged(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "diff", "--cached", "--name-status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("falha ao verificar arquivos em stage: %s", strings.TrimSpace(string(out)))
	}
	linhas := strings.Split(strings.TrimSpace(string(out)), "\n")
	var arquivos []string
	for _, l := range linhas {
		if strings.TrimSpace(l) != "" {
			arquivos = append(arquivos, l)
		}
	}
	return arquivos, nil
}

func gerarMensagemFallback(staged []string) string {
	if len(staged) == 1 {
		partes := strings.Fields(staged[0])
		if len(partes) >= 2 {
			return fmt.Sprintf("chore: atualizar %s", partes[1])
		}
	}
	return fmt.Sprintf("chore: sincronizar %d arquivos do repositório", len(staged))
}

func rodarCommit(ctx context.Context, msg string) error {
	cmd := exec.CommandContext(ctx, "git", "commit", "-m", msg)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("falha no git commit: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func rodarPush(ctx context.Context, stdout io.Writer) error {
	cmdBranch := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	outBranch, err := cmdBranch.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao detectar branch atual: %s", strings.TrimSpace(string(outBranch)))
	}
	branch := strings.TrimSpace(string(outBranch))
	cmdPush := exec.CommandContext(ctx, "git", "push", "origin", branch)
	if outPush, err := cmdPush.CombinedOutput(); err != nil {
		cmdSetUpstream := exec.CommandContext(ctx, "git", "push", "--set-upstream", "origin", branch)
		if outUp, errUp := cmdSetUpstream.CombinedOutput(); errUp != nil {
			return fmt.Errorf("falha no git push: %s", strings.TrimSpace(string(outPush))+"\n"+strings.TrimSpace(string(outUp)))
		}
	}
	fmt.Fprintf(stdout, "Push concluído para origin/%s\n", branch)
	return nil
}
