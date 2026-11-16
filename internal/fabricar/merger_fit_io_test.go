package fabricar_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/matriz"
)

// TestMergeByFitPorIO cobre a AC-E05.1 no caminho real: MergeByFit carrega o contrato.json ao lado
// do SKILL.md e a peça que produz o que a etapa consome é escolhida mesmo sem caso de uso.
func TestMergeByFitPorIO(t *testing.T) {
	vault := t.TempDir()
	dir := filepath.Join(vault, "skills", "coletor-io")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifesto := `{"nome":"coletor-io","tipo":"script","naturezas":["dados"],
		"entra":[{"tipo":"texto"}],"sai":[{"tipo":"dados","promete":true}]}`
	if err := os.WriteFile(filepath.Join(dir, "contrato.json"), []byte(manifesto), 0o644); err != nil {
		t.Fatal(err)
	}
	ref := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(ref, []byte("# coletor-io"), 0o644); err != nil {
		t.Fatal(err)
	}
	e := skill(ref, 0.5, []string{"coletar"}, []string{"dados"})
	e.Nome = "coletor-io"
	n := matriz.Needs{Natureza: "dados", Etapas: []string{"coletar"},
		Casos: []string{"coleta-de-dados"}, Consome: []string{"dados"}}

	t.Run("peca com contrato em disco é escolhida para a etapa", func(t *testing.T) {
		m := fabricar.MergeByFit(n, []matriz.Entrada{e}, fabricar.MergeOptions{})
		if !slices.Contains(nomes(m), "coletor-io") {
			t.Fatalf("a peça que produz dados devia compor a etapa; passos: %v", nomes(m))
		}
	})

	t.Run("sem contrato em disco vale o comportamento atual", func(t *testing.T) {
		avulsa := skill(filepath.Join(vault, "skills", "avulsa", "SKILL.md"), 0.5,
			[]string{"coletar"}, []string{"dados"})
		avulsa.Nome = "avulsa"
		m := fabricar.MergeByFit(n, []matriz.Entrada{avulsa}, fabricar.MergeOptions{})
		if slices.Contains(nomes(m), "avulsa") {
			t.Fatalf("sem contrato a peça sem caso de uso não pode compor: %v", nomes(m))
		}
	})
}
