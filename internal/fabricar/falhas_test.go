package fabricar_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/cli"
	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/provider"
)

// slowLLM simula um cérebro que expira por timeout na Camada 4 (SPEC)
type slowLLM struct {
	delay time.Duration
}

func (s *slowLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	if strings.Contains(msgs[0].Content, "Arquiteto de Especificações") {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(s.delay):
			return "", context.DeadlineExceeded
		}
	}
	return "HIPÓTESE: extrator\nCONFIANÇA: 98%\nESPELHO CANÔNICO:\n- Resultado Esperado: áudio", nil
}

func TestFalhas(t *testing.T) {
	t.Run("processo ja existe: recusa antes de gravar e oferece refinar", func(t *testing.T) {
		r := newRig(t, goodPieces())
		// primeira vez: sucesso
		r.run(t)

		// segunda vez sem force: ErrExists
		_, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
		if !errors.Is(err, fabricar.ErrExists) {
			t.Fatalf("esperava ErrExists: %v", err)
		}
	})

	t.Run("cerebro expira na Camada 4: cancela sem gravar arquivos", func(t *testing.T) {
		r := newRig(t, goodPieces())
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		r.deps.LLM = &slowLLM{delay: 100 * time.Millisecond}
		_, err := fabricar.Fabricar(ctx, r.req, r.deps)
		if err == nil {
			t.Fatal("esperava erro de context/timeout")
		}
		if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
			t.Fatalf("esperava DeadlineExceeded ou Canceled, veio %v", err)
		}

		// Garante zero arquivos gravados
		for _, root := range []string{r.deps.Roots.Vault, r.deps.Roots.Docs} {
			if _, statErr := os.Stat(root); statErr == nil {
				t.Errorf("não deve gravar nada quando cérebro expira: %s", root)
			}
		}
	})

	t.Run("usuario responde cancelar no cartao: zero arquivos gravados e ErrAborted", func(t *testing.T) {
		r := newRig(t, goodPieces())
		r.deps.Aprovar = func(p fabricar.Plano, replan func([]string) fabricar.Plano) (bool, error) {
			return false, nil // cancelou
		}
		_, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
		if !errors.Is(err, fabricar.ErrAborted) {
			t.Fatalf("esperava ErrAborted, veio %v", err)
		}

		for _, root := range []string{r.deps.Roots.Vault, r.deps.Roots.Docs} {
			if _, statErr := os.Stat(root); statErr == nil {
				t.Errorf("não deve gravar nada quando cancelado no cartão: %s", root)
			}
		}
	})

	t.Run("cli falhas de processo existente oferta refinar e timeout", func(t *testing.T) {
		home := t.TempDir()
		proj := t.TempDir()
		vault := filepath.Join(home, ".aidd", "canonical")

		// Cria processo existente
		r := newRig(t, goodPieces())
		r.deps.Roots.Vault = vault
		r.deps.Roots.Docs = proj
		r.run(t)

		// Executa CLI que detecta existente
		var out, errb strings.Builder
		d := cli.Deps{
			Home:   home,
			Stdout: &out,
			Stderr: &errb,
			Stdin:  strings.NewReader(""),
			Discover: func(context.Context) provider.Inventory {
				return provider.Inventory{CLIs: []provider.CLI{{Name: "claude"}}}
			},
			Fabricar: &cli.FabricarDeps{
				LLM:        r.deps.LLM,
				Scanner:    r.deps.Scanner,
				Clock:      r.deps.Clock,
			},
		}
		code := cli.Run(context.Background(), []string{
			"fabricar", "Extrator de Áudio", "--yes", "--sem-entrevista",
			"--vault", vault, "--docs", proj, "--model", "sonnet",
		}, d)

		if code != 1 {
			t.Errorf("esperava exit 1, obteve %d. Stderr:\n%s", code, errb.String())
		}
		if !strings.Contains(errb.String(), "refinar") {
			t.Errorf("saída de erro deve conter oferta de 'refinar':\n%s", errb.String())
		}
	})
}
