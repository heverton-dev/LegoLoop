package cli

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"
)

type saidaCLI struct {
	code   int
	stdout string
	stderr string
}

func rodarCLI(t *testing.T, args ...string) saidaCLI {
	t.Helper()
	var out, errb bytes.Buffer
	code := Run(context.Background(), args, Deps{Stdout: &out, Stderr: &errb})
	return saidaCLI{code: code, stdout: out.String(), stderr: errb.String()}
}

func TestVideoAliasComandosNovosRegistrados(t *testing.T) {
	verbos := Verbos()
	for _, nome := range []string{"video-roteiro", "video-ativos", "video-render", "diorama-roteiro", "diorama-ativos"} {
		if !slices.Contains(verbos, nome) {
			t.Errorf("verbo %q ausente do mapa de comandos", nome)
		}
		if !strings.Contains(usage, "  "+nome+" ") {
			t.Errorf("verbo %q ausente do texto de uso", nome)
		}
	}
}

func TestVideoAliasMesmoCodigoComAviso(t *testing.T) {
	pares := map[string]string{"diorama-roteiro": "video-roteiro", "diorama-ativos": "video-ativos"}
	for antigo, novo := range pares {
		for _, args := range [][]string{{"subcomando-inexistente"}, {}} {
			n := rodarCLI(t, append([]string{novo}, args...)...)
			a := rodarCLI(t, append([]string{antigo}, args...)...)
			if a.code != n.code {
				t.Errorf("%s %v: exit %d, esperado %d (igual a %s)", antigo, args, a.code, n.code, novo)
			}
			if a.stdout != n.stdout {
				t.Errorf("%s %v: stdout diverge de %s", antigo, args, novo)
			}
			aviso := avisoDepreciacao(antigo, novo)
			if !strings.HasPrefix(a.stderr, aviso) || a.stderr[len(aviso):] != n.stderr {
				t.Errorf("%s %v: stderr %q, esperado aviso %q seguido de %q", antigo, args, a.stderr, aviso, n.stderr)
			}
			if strings.Contains(n.stderr, "depreciado") {
				t.Errorf("%s não deve imprimir aviso de depreciação: %q", novo, n.stderr)
			}
		}
	}
}

func TestVideoAliasRenderSubcomandos(t *testing.T) {
	ajuda := rodarCLI(t, "video-render")
	if ajuda.code != 0 || !strings.Contains(ajuda.stdout, "timeline") || !strings.Contains(ajuda.stdout, "executar") {
		t.Fatalf("video-render sem argumentos: exit %d, stdout %q", ajuda.code, ajuda.stdout)
	}
	if r := rodarCLI(t, "video-render", "subcomando-inexistente"); r.code != 2 {
		t.Errorf("subcomando desconhecido: exit %d, esperado 2", r.code)
	}
	for _, sub := range []string{"timeline", "executar"} {
		if r := rodarCLI(t, "video-render", sub); r.code != 2 || !strings.Contains(r.stderr, "--projeto") {
			t.Errorf("%s sem --projeto: exit %d, stderr %q", sub, r.code, r.stderr)
		}
		dir := t.TempDir()
		if r := rodarCLI(t, "video-render", sub, "--projeto", dir); r.code != 1 {
			t.Errorf("%s em projeto vazio: exit %d, esperado 1; stderr %q", sub, r.code, r.stderr)
		}
	}
}
