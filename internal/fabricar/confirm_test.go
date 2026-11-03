package fabricar_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/fabricar"
)

// written conta os arquivos gravados nas raízes do rig.
func (r *rig) written(t *testing.T) int {
	t.Helper()
	n := 0
	for _, root := range []string{r.deps.Roots.Vault, r.deps.Roots.Docs} {
		_ = filepath.WalkDir(root, func(_ string, e os.DirEntry, err error) error {
			if err == nil && !e.IsDir() {
				n++
			}
			return nil
		})
	}
	return n
}

func TestFabricar_ConfirmReceivesAllPathsBeforeWriting(t *testing.T) {
	r := newRig(t, goodPieces())
	var got []string
	var wroteAtConfirm int
	r.deps.Confirm = func(paths []string) (bool, error) {
		got = paths
		wroteAtConfirm = r.written(t)
		return true, nil
	}
	res := r.run(t)
	if len(got) != 8 {
		t.Errorf("Confirm devia receber os 8 caminhos: %v", got)
	}
	if wroteAtConfirm != 0 {
		t.Errorf("nada pode estar gravado quando Confirm é chamado: %d", wroteAtConfirm)
	}
	if _, err := os.Stat(res.Artifacts.SkillPath); err != nil {
		t.Errorf("Confirm=true devia gravar: %v", err)
	}
}

func TestFabricar_ConfirmDeniedAbortsWithoutWriting(t *testing.T) {
	r := newRig(t, goodPieces())
	r.deps.Confirm = func([]string) (bool, error) { return false, nil }
	_, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
	if !errors.Is(err, fabricar.ErrAborted) {
		t.Fatalf("esperava ErrAborted, veio %v", err)
	}
	if n := r.written(t); n != 0 {
		t.Errorf("abortou mas gravou %d arquivos", n)
	}
}

func TestFabricar_ConfirmErrorAbortsWithoutWriting(t *testing.T) {
	r := newRig(t, goodPieces())
	boom := errors.New("stdin fechou")
	r.deps.Confirm = func([]string) (bool, error) { return false, boom }
	if _, err := fabricar.Fabricar(context.Background(), r.req, r.deps); !errors.Is(err, boom) {
		t.Fatalf("o erro do Confirm deve subir: %v", err)
	}
	if n := r.written(t); n != 0 {
		t.Errorf("gravou %d arquivos", n)
	}
}

func TestFabricar_DryRunDoesNotAskForConfirmation(t *testing.T) {
	r := newRig(t, goodPieces())
	r.req.DryRun = true
	r.deps.Confirm = func([]string) (bool, error) {
		t.Error("dry-run não grava, não deve pedir confirmação")
		return false, nil
	}
	r.run(t)
}
