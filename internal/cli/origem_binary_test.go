package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const needsConstruirTDD = `{"natureza": "software", "etapas": ["construir"], "casos": ["tdd"]}`

// matrizDeOrigem grava no cofre uma skill por item, todas aptas a construir/tdd, cada uma com score e origem.
func matrizDeOrigem(t *testing.T, vault string, pecas ...[3]string) {
	t.Helper()
	var ratings, entradas []map[string]any
	for _, p := range pecas {
		nome, origem, score := p[0], p[1], p[2]
		ref := "~/.claude/skills/" + nome + "/SKILL.md"
		var s float64
		if err := json.Unmarshal([]byte(score), &s); err != nil {
			t.Fatal(err)
		}
		ratings = append(ratings, map[string]any{"ref": ref, "name": nome, "total_score": s})
		entradas = append(entradas, map[string]any{"ref": ref, "nome": nome, "tipo": "skill", "hash": "h-" + nome,
			"perfil": map[string]any{"etapas": []string{"construir"}, "naturezas": []string{"software"},
				"casos_de_uso": []string{"tdd"}, "origem": origem}})
	}
	reg := map[string]any{"generated_at": "2026-11-01T10:00:00Z", "matrix_version": "2.0.0", "total_pieces": len(pecas),
		"ratings": ratings, "entradas": entradas}
	b, err := json.Marshal(reg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(vault, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vault, "registry_matrix.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func rodarOrigem(t *testing.T, bin string, extra string, pecas ...[3]string) binRun {
	t.Helper()
	cer := &cerebroPorPapel{needsJSON: needsConstruirTDD}
	root := t.TempDir()
	vault, proj := filepath.Join(root, "vault"), filepath.Join(root, "proj")
	matrizDeOrigem(t, vault, pecas...)
	args := []string{"fabricar", "extrair audio", "--goal", "construir um extrator", "--vault", vault, "--docs", proj,
		"--yes", "--sem-entrevista", "--dry-run", "--no-link"}
	if extra != "" {
		args = append(args, extra)
	}
	return runBinaryEm(t, bin, cer.servidor(t), root, proj, args...)
}

// TestOrigemDoPerfilBinario prova o AC-E03.8 pelo binário: a declarada vence a inferida de Score maior, o
// motivo cita a origem, e o texto e o --json mostram a origem. Com só a inferida, ela é escolhida.
func TestOrigemDoPerfilBinario(t *testing.T) {
	bin := buildBinary(t)
	duas := [][3]string{{"peca-inferida", "inferido", "0.95"}, {"peca-declarada", "declarado", "0.60"}}
	t.Run("texto", func(t *testing.T) {
		r := rodarOrigem(t, bin, "", duas...)
		if r.code != 0 {
			t.Fatalf("exit %d\n%s\n%s", r.code, r.out, r.errb)
		}
		linha := linhaDaEtapa(r.out, "construir")
		if !strings.Contains(linha, "peca-declarada") || !strings.Contains(linha, "perfil: declarado") || strings.Contains(linha, "peca-inferida") {
			t.Errorf("a linha de construir deve ter a declarada e a origem: %q\n%s", linha, r.out)
		}
	})
	t.Run("json", func(t *testing.T) {
		r := rodarOrigem(t, bin, "--json", duas...)
		if r.code != 0 {
			t.Fatalf("exit %d\n%s\n%s", r.code, r.out, r.errb)
		}
		d := decisaoConstruir(t, r.out)
		if d.Peca != "peca-declarada" || d.Origem != "declarado" || !strings.Contains(d.Motivo, "perfil: declarado") {
			t.Errorf("decisão = %+v", d)
		}
	})
	t.Run("só a inferida", func(t *testing.T) {
		r := rodarOrigem(t, bin, "", [3]string{"peca-inferida", "inferido", "0.95"})
		if r.code != 0 {
			t.Fatalf("exit %d\n%s\n%s", r.code, r.out, r.errb)
		}
		linha := linhaDaEtapa(r.out, "construir")
		if !strings.Contains(linha, "peca-inferida") || !strings.Contains(linha, "perfil: inferido") {
			t.Errorf("a inferida sozinha é escolhida: %q", linha)
		}
		rj := rodarOrigem(t, bin, "--json", [3]string{"peca-inferida", "inferido", "0.95"})
		if d := decisaoConstruir(t, rj.out); d.Peca != "peca-inferida" || d.Origem != "inferido" {
			t.Errorf("decisão = %+v", d)
		}
	})
}

func linhaDaEtapa(out, etapa string) string {
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), etapa+" ") {
			return l
		}
	}
	return ""
}

func decisaoConstruir(t *testing.T, out string) (d struct{ Etapa, Peca, Motivo, Origem string }) {
	t.Helper()
	var s struct {
		Report struct {
			Fusao struct {
				Decisoes []struct{ Etapa, Peca, Motivo, Origem string } `json:"decisoes"`
			} `json:"fusao"`
		} `json:"report"`
	}
	if err := json.Unmarshal([]byte(out), &s); err != nil {
		t.Fatalf("stdout não é JSON: %v\n%s", err, out)
	}
	for _, x := range s.Report.Fusao.Decisoes {
		if x.Etapa == "construir" {
			return x
		}
	}
	t.Fatalf("sem decisão para construir: %s", out)
	return
}
