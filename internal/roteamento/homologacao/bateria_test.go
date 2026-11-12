package homologacao

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type falso struct {
	viola string
}

func (f falso) Executar(ctx context.Context, dir, prompt string, teto time.Duration) (Saida, error) {
	s := Saida{ExitCode: 0, CustoUSD: 0.01, Tokens: 10, PicoMB: 420, Duracao: time.Millisecond}
	if f.viola == "C3" {
		s.CustoUSD, s.Tokens = 0, 0
	}
	switch {
	case strings.Contains(prompt, "ERRO_FORCADO"):
		if f.viola != "C5" {
			s.ExitCode = 1
		}
	case strings.Contains(prompt, "CHAMADA_LONGA"):
		if f.viola == "C4" {
			s.Duracao = time.Hour
		}
	default:
		if err := os.WriteFile(filepath.Join(dir, "permitido.txt"), []byte("conteudo-conhecido"), 0o644); err != nil {
			return s, err
		}
		if f.viola == "C2" {
			if err := os.MkdirAll(filepath.Join(dir, "fora"), 0o755); err != nil {
				return s, err
			}
			if err := os.WriteFile(filepath.Join(dir, "fora", "proibido.txt"), []byte("x"), 0o644); err != nil {
				return s, err
			}
		}
		if f.viola == "C1" {
			return s, context.DeadlineExceeded
		}
	}
	return s, nil
}

func TestBateriaAprovaEReprova(t *testing.T) {
	base := t.TempDir()
	r := Avaliar(context.Background(), falso{}, base, time.Second)
	if r.Estado != EstadoHomologado || r.PesoMB != 420 || r.Motivo != "" {
		t.Fatalf("esperava homologado com peso 420, veio %+v", r)
	}
	for _, k := range []string{"C1", "C2", "C3", "C4", "C5"} {
		r := Avaliar(context.Background(), falso{viola: k}, base, time.Second)
		if r.Estado != EstadoReprovado || r.Criterios[k] || !strings.HasPrefix(r.Motivo, k) {
			t.Errorf("%s: esperava reprovado por %s, veio %+v", k, k, r)
		}
		if r.PesoMB != 420 {
			t.Errorf("%s: peso não gravado: %d", k, r.PesoMB)
		}
	}
	sobras, _ := os.ReadDir(base)
	if len(sobras) != 0 {
		t.Errorf("worktree descartável não removida: %d itens", len(sobras))
	}
}

func TestBateriaC2NaoPodeSerIgnorado(t *testing.T) {
	r := Avaliar(context.Background(), falso{viola: "C2"}, t.TempDir(), time.Second)
	if r.Estado != EstadoReprovado || r.Criterios["C2"] || !strings.HasPrefix(r.Motivo, "C2") {
		t.Fatalf("esperava reprovado por C2, veio %+v", r)
	}
}

func TestForaDoLimiteEReverter(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "permitido.txt"), []byte("a"), 0o644)
	_ = os.MkdirAll(filepath.Join(dir, "fora"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, "fora", "proibido.txt"), []byte("b"), 0o644)
	lista := ForaDoLimite(dir, []string{"permitido.txt"})
	if len(lista) != 1 || lista[0] != "fora/proibido.txt" {
		t.Fatalf("lista inesperada: %v", lista)
	}
	if err := Reverter(dir, lista); err != nil {
		t.Fatal(err)
	}
	if len(ForaDoLimite(dir, []string{"permitido.txt"})) != 0 {
		t.Fatal("escrita fora do limite não foi revertida")
	}
}
