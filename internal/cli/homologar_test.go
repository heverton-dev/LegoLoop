package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/roteamento"
	"legoloop/internal/roteamento/homologacao"
)

type harnessContador struct {
	chamadas *int
	saida    *bytes.Buffer
	antes    *string
}

func (h harnessContador) Executar(_ context.Context, _, prompt string, _ time.Duration) (homologacao.Saida, error) {
	if *h.chamadas == 0 {
		*h.antes = h.saida.String()
	}
	*h.chamadas++
	s := homologacao.Saida{CustoUSD: 0.01, Tokens: 10, PicoMB: 420}
	if prompt == homologacao.PromptB {
		s.ExitCode = 1
	}
	return s, nil
}

func ambienteHomologar(t *testing.T, livre uint64) (Deps, *bytes.Buffer, *int, *string) {
	t.Helper()
	out := &bytes.Buffer{}
	chamadas, antes := new(int), new(string)
	inv := []roteamento.ItemInv{{
		Par:     roteamento.Par{Harness: "claude", Via: roteamento.ViaOficial},
		Binario: "/bin/claude", Modelos: []string{"sonnet"}, Estado: roteamento.EstadoPendente,
	}}
	d := Deps{
		Stdin: strings.NewReader(""), Stdout: out, Stderr: &bytes.Buffer{}, Home: t.TempDir(),
		Rota: &RotaDeps{Inventario: func(context.Context) ([]roteamento.ItemInv, error) { return inv, nil }},
		Homologar: &HomologarDeps{
			Harness: func(_, _, _ string) homologacao.Harness {
				return harnessContador{chamadas: chamadas, saida: out, antes: antes}
			},
			MemoriaLivre: func() (uint64, error) { return livre, nil },
		},
	}
	return d, out, chamadas, antes
}

func TestHomologarConfirmacaoEPendente(t *testing.T) {
	const chave = "claude|oficial"
	for _, args := range [][]string{{}, {"--non-interactive"}} {
		d, out, chamadas, _ := ambienteHomologar(t, 4<<30)
		if code := Run(context.Background(), append([]string{"homologar"}, args...), d); code != 0 {
			t.Fatalf("exit %d: %s", code, out)
		}
		if *chamadas != 0 {
			t.Fatalf("sem confirmação não pode chamar: %d", *chamadas)
		}
		if !strings.Contains(out.String(), "Teto total estimado: 3 chamadas") {
			t.Fatalf("falta o teto total: %s", out)
		}
		st, err := homologacao.Ler(homologacao.Caminho(d.Home))
		if err != nil || st[chave].Estado != homologacao.EstadoPendente {
			t.Fatalf("esperava pendente: %v %v", st, err)
		}
	}

	d, out, chamadas, antes := ambienteHomologar(t, 4<<30)
	if code := Run(context.Background(), []string{"homologar", "--yes", "--teto", "1s"}, d); code != 0 {
		t.Fatalf("exit %d: %s", code, out)
	}
	if *chamadas != 3 || !strings.Contains(*antes, "Teto total estimado") {
		t.Fatalf("chamadas=%d, teto antes da primeira chamada: %q", *chamadas, *antes)
	}
	st, err := homologacao.Ler(homologacao.Caminho(d.Home))
	if err != nil || st[chave].Estado != homologacao.EstadoHomologado || st[chave].PesoMB != 420 {
		t.Fatalf("esperava homologado com peso: %v %v", st, err)
	}
	lido := rotaLerHomologacao(filepath.Join(d.Home, ".aidd", "homologacao.json"))
	if lido[chave].Estado != roteamento.EstadoHomologado || lido[chave].PesoMB != 420 {
		t.Fatalf("o inventário não lê o resultado: %v", lido)
	}

	d, out, chamadas, _ = ambienteHomologar(t, 100<<20)
	if code := Run(context.Background(), []string{"homologar", "--yes"}, d); code != 1 {
		t.Fatalf("memória baixa deve dar exit 1, veio %d", code)
	}
	st, _ = homologacao.Ler(homologacao.Caminho(d.Home))
	if *chamadas != 0 || len(st) != 0 {
		t.Fatalf("memória baixa não pode chamar nem gravar: %d %v", *chamadas, st)
	}
}
