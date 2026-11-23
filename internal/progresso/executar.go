package progresso

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Executor executa o comando de validação de um Ticket e devolve exit code e saída resumida.
type Executor func(ctx context.Context, dir, comando string) (int, string, error)

// Iniciar marca um Ticket como IN_PROGRESS registrando responsável e ponto de retomada.
func Iniciar(p *Progresso, ticketID, responsavel, nota, agora string) error {
	idx := buscarIndice(p.Tickets, ticketID)
	if idx < 0 {
		return fmt.Errorf("ticket nao encontrado: %s", ticketID)
	}
	it := &p.Tickets[idx]
	it.Estado = EstadoEmAndamento
	it.Responsavel = strings.TrimSpace(responsavel)
	it.Prova = resumirSaida(nota)
	it.AtualizadoEm = agora
	p.AtualizadoEm = agora
	p.RecalcularResumo()
	return nil
}

// ValidarTicket executa o Gate do Ticket e promove para PASSED apenas se exit == 0.
func ValidarTicket(ctx context.Context, p *Progresso, ticketID, dir, agora string, execFn Executor) (bool, error) {
	idx := buscarIndice(p.Tickets, ticketID)
	if idx < 0 {
		return false, fmt.Errorf("ticket nao encontrado: %s", ticketID)
	}
	it := &p.Tickets[idx]
	if strings.TrimSpace(it.Comando) == "" {
		return false, fmt.Errorf("ticket %s sem comando de verificacao", ticketID)
	}
	if execFn == nil {
		execFn = ExecutorPadrao
	}
	code, saida, err := execFn(ctx, dir, it.Comando)
	if err != nil {
		it.Estado = EstadoReprovado
		it.ExitCode = code
		it.Prova = resumirSaida(err.Error())
		it.AtualizadoEm = agora
		p.AtualizadoEm = agora
		p.RecalcularResumo()
		return false, err
	}
	return aplicarResultadoValidacao(p, it, code, saida, agora), nil
}

func aplicarResultadoValidacao(p *Progresso, it *ItemTicket, code int, saida, agora string) bool {
	it.ExitCode = code
	it.Prova = resumirSaida(saida)
	it.AtualizadoEm = agora
	p.AtualizadoEm = agora
	if code == 0 {
		it.Estado = EstadoAprovado
	} else {
		it.Estado = EstadoReprovado
	}
	p.RecalcularResumo()
	return code == 0
}

// ExecutorPadrao executa o comando de verificação no diretório informado.
func ExecutorPadrao(ctx context.Context, dir, comando string) (int, string, error) {
	cmd := construirComando(ctx, comando)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	saida := resumirSaida(buf.String())
	if err == nil {
		return 0, saida, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), saida, nil
	}
	return 1, saida, err
}

func construirComando(ctx context.Context, comando string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		if sh, err := exec.LookPath("sh"); err == nil {
			return exec.CommandContext(ctx, sh, "-c", comando)
		}
		return exec.CommandContext(ctx, "cmd", "/c", comando)
	}
	return exec.CommandContext(ctx, "sh", "-c", comando)
}

func buscarIndice(itens []ItemTicket, id string) int {
	for i := range itens {
		if itens[i].ID == id {
			return i
		}
	}
	return -1
}

func resumirSaida(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", " | "), "\n", " | "))
	runas := []rune(s)
	if len(runas) > 220 {
		return string(runas[:217]) + "..."
	}
	return s
}
