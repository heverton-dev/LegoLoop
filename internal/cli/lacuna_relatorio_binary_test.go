package cli

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// needsComLacuna pede construir e verificar, mas só o caso tdd (que pertence a construir): verificar fica sem
// caso pedido e vira lacuna.
const needsComLacuna = `{"natureza": "software", "etapas": ["construir", "verificar"], "casos": ["tdd"]}`

// TestLacunaRelatorioBinario prova o AC-E02.7 pelo binário nos dois modos. Com --dry-run: exit 1, o relatório
// e o cartão saem com a lacuna destacada (em texto e em --json). Sem --dry-run: exit 1 e nada gravado.
func TestLacunaRelatorioBinario(t *testing.T) {
	bin := buildBinary(t)
	t.Run("dry-run texto", func(t *testing.T) {
		r, root := rodarLacunaBin(t, bin, "--dry-run")
		saida := r.out + r.errb
		if r.code != 0 {
			t.Fatalf("exit %d; esperava 0\n%s", r.code, saida)
		}
		for _, quer := range []string{"verificar  LACUNA |", "aptidão desconhecida: 1", "aviso: Lacuna: a etapa verificar", "sugestão externa (GitHub/Fontes): curadoria-externa-verificar"} {
			if !strings.Contains(saida, quer) {
				t.Errorf("falta %q na saída:\n%s", quer, saida)
			}
		}
		semArquivos(t, root)
	})
	t.Run("dry-run json", func(t *testing.T) {
		r, root := rodarLacunaBin(t, bin, "--dry-run", "--json")
		if r.code != 0 {
			t.Fatalf("exit %d; esperava 0 com cascata de 3 fontes\n%s\n%s", r.code, r.out, r.errb)
		}
		var s struct {
			CuratedPieces []any `json:"curated_pieces"`
		}
		if err := json.Unmarshal([]byte(r.out), &s); err != nil {
			t.Fatalf("stdout deve ser JSON puro: %v\n%s", err, r.out)
		}
		if len(s.CuratedPieces) == 0 {
			t.Errorf("esperava sugestões externas na cascata de 3 fontes: %+v", s)
		}
		semArquivos(t, root)
	})
	t.Run("sem dry-run", func(t *testing.T) {
		r, _ := rodarLacunaBin(t, bin)
		if r.code != 0 {
			t.Fatalf("exit %d; esperava 0 com cascata de 3 fontes\n%s\n%s", r.code, r.out, r.errb)
		}
	})
}

func rodarLacunaBin(t *testing.T, bin string, extra ...string) (binRun, string) {
	t.Helper()
	cer := &cerebroPorPapel{needsJSON: needsComLacuna}
	root := t.TempDir()
	vault, proj := filepath.Join(root, "vault"), filepath.Join(root, "proj")
	matrizDeTeste(t, vault)
	args := append([]string{"fabricar", "extrair audio", "--goal", "construir e testar um extrator", "--vault", vault,
		"--docs", proj, "--yes", "--sem-entrevista", "--no-link"}, extra...)
	return runBinaryEm(t, bin, cer.servidor(t), root, proj, args...), root
}

// semArquivos confere que o cofre só tem a matriz e que nenhum documento foi gravado.
func semArquivos(t *testing.T, root string) {
	t.Helper()
	for _, f := range filesUnder(root) {
		if !strings.HasSuffix(f, "registry_matrix.json") {
			t.Errorf("a lacuna não pode gravar nada: %s", f)
		}
	}
}
