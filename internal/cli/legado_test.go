package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestLoopsLegadoExitDois prova o caminho de falha do rename: comando que
// aponta para o caminho legado `loops/` termina com exit 2 e o guia de
// migração mostra o caminho novo `processos/` (AC-E06.3).
func TestLoopsLegadoExitDois(t *testing.T) {
	casos := []string{
		"loops/legoloop-fechar/pipeline.dag.json",
		"./loops/legoloop-fechar/pipeline.dag.json",
		`vault\loops\legoloop-fechar\pipeline.dag.json`,
	}
	for _, alvo := range casos {
		var saida bytes.Buffer
		code := Run(context.Background(), []string{"validate", alvo}, Deps{Stdout: &saida, Stderr: &saida})
		if code != 2 {
			t.Fatalf("alvo %q: exit %d, esperado 2", alvo, code)
		}
		txt := saida.String()
		if !strings.Contains(txt, "guia de migração") {
			t.Fatalf("alvo %q: sem guia de migração na saída: %q", alvo, txt)
		}
		if !strings.Contains(txt, "processos") {
			t.Fatalf("alvo %q: guia não cita o caminho novo: %q", alvo, txt)
		}
	}
}
