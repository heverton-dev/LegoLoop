package fabricar_test

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/processo"
)

// comHarnesses cria as pastas dos harnesses no HOME do rig e os entrega ao motor como detectados.
func comHarnesses(t *testing.T, r *rig, nomes ...string) {
	t.Helper()
	for _, n := range nomes {
		if err := os.MkdirAll(filepath.Join(r.home, "."+n), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	r.deps.Harnesses = processo.DetectarHarnesses(processo.OSFS(), r.home)
}

func arquivosSob(t *testing.T, dirs ...string) []string {
	t.Helper()
	var out []string
	for _, d := range dirs {
		_ = filepath.WalkDir(d, func(p string, e fs.DirEntry, err error) error {
			if err == nil && !e.IsDir() {
				out = append(out, p)
			}
			return nil
		})
	}
	return out
}

func TestVinculoTodosOsHarnessesDetectados(t *testing.T) {
	r := newRig(t, goodPieces())
	comHarnesses(t, r, "claude", "gemini", "codex")
	var previa []string
	r.deps.Confirm = func(paths []string) (bool, error) { previa = paths; return true, nil }
	res := r.run(t)
	if len(res.Artifacts.SkillLinks) != 3 {
		t.Fatalf("esperava 3 vínculos: %v", res.Artifacts.SkillLinks)
	}
	for _, l := range res.Artifacts.SkillLinks {
		if !strings.Contains(read(t, filepath.Join(l, "SKILL.md")), "name: legoloop-extrator-de-audio") {
			t.Errorf("o vínculo %s deve expor o SKILL.md do cofre", l)
		}
		if !contem(previa, l) {
			t.Errorf("a prévia deve listar %s: %v", l, previa)
		}
	}
}

func contem(xs []string, x string) bool {
	for _, s := range xs {
		if s == x {
			return true
		}
	}
	return false
}

func TestVinculoNoLinkNaoCriaNenhum(t *testing.T) {
	r := newRig(t, goodPieces())
	comHarnesses(t, r, "claude", "gemini", "codex")
	r.req.NoLink = true
	res := r.run(t)
	if len(res.Artifacts.SkillLinks) != 0 {
		t.Fatalf("--no-link não cria vínculo: %v", res.Artifacts.SkillLinks)
	}
	if _, err := os.Stat(filepath.Join(r.home, ".claude", "skills")); err == nil {
		t.Error("com --no-link a pasta do harness não pode ser tocada")
	}
}

func TestVinculoFalhaNoSegundoDesfazOsLinksEOCofre(t *testing.T) {
	r := newRig(t, goodPieces())
	comHarnesses(t, r, "claude", "gemini", "codex")
	chamadas, quebra := 0, errors.New("disco cheio")
	r.deps.Link = func(o, d string) error {
		chamadas++
		if chamadas == 2 {
			return quebra
		}
		return processo.SymlinkOuCopia(o, d)
	}
	_, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
	if !errors.Is(err, quebra) {
		t.Fatalf("esperava a falha do 2º link: %v", err)
	}
	if sobras := arquivosSob(t, r.deps.Roots.Vault, r.deps.Roots.Docs); len(sobras) != 0 {
		t.Errorf("a falha deve desfazer o cofre e os docs: %v", sobras)
	}
	for _, n := range []string{"claude", "gemini", "codex"} {
		if _, err := os.Lstat(filepath.Join(r.home, "."+n, "skills", "legoloop-extrator-de-audio")); err == nil {
			t.Errorf("o link de %s devia ter sido desfeito", n)
		}
	}
}

func TestVinculoDryRunListaMasNaoCria(t *testing.T) {
	r := newRig(t, goodPieces())
	comHarnesses(t, r, "claude", "gemini")
	r.req.DryRun = true
	res := r.run(t)
	if len(res.Artifacts.SkillLinks) != 2 {
		t.Fatalf("o dry-run deve listar os 2 vínculos que criaria: %v", res.Artifacts.SkillLinks)
	}
	for _, l := range res.Artifacts.SkillLinks {
		if _, err := os.Lstat(l); err == nil {
			t.Errorf("dry-run não cria: %s", l)
		}
	}
}

func TestVinculoLinkHarnessRestringeEAvisaOQueNaoExiste(t *testing.T) {
	r := newRig(t, goodPieces())
	comHarnesses(t, r, "claude", "gemini")
	r.req.LinkHarnesses = []string{"claude", "kiro"}
	res := r.run(t)
	if len(res.Artifacts.SkillLinks) != 1 || !strings.Contains(res.Artifacts.SkillLinks[0], ".claude") {
		t.Fatalf("só o claude deve ser ligado: %v", res.Artifacts.SkillLinks)
	}
	if _, err := os.Stat(filepath.Join(r.home, ".gemini", "skills")); err == nil {
		t.Error("harness fora do pedido não pode ser tocado")
	}
	avisou := false
	for _, w := range res.Warnings {
		avisou = avisou || strings.Contains(w, "kiro")
	}
	if !avisou {
		t.Errorf("harness pedido e não detectado deve gerar aviso: %v", res.Warnings)
	}
}

func TestVinculoSemHarnessDetectadoNaoLigaNada(t *testing.T) {
	r := newRig(t, goodPieces())
	if res := r.run(t); len(res.Artifacts.SkillLinks) != 0 {
		t.Fatalf("sem harness detectado não há vínculo: %v", res.Artifacts.SkillLinks)
	}
}
