package processo_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/processo"
)

// homeCom cria um HOME descartável com as pastas dos harnesses pedidos.
func homeCom(t *testing.T, nomes ...string) string {
	t.Helper()
	home := t.TempDir()
	for _, n := range nomes {
		if err := os.MkdirAll(filepath.Join(home, "."+n), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return home
}

func skillNoCofre(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "skills", "legoloop-x")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: legoloop-x\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestVinculo(t *testing.T) {
	t.Run("detecta so os harnesses que existem", func(t *testing.T) {
		home := homeCom(t, "claude", "gemini", "codex")
		hs := processo.DetectarHarnesses(processo.OSFS(), home)
		if len(hs) != 3 {
			t.Fatalf("esperava 3 harnesses, veio %v", hs)
		}
		if hs[0].Nome != "claude" || hs[0].Skills != filepath.Join(home, ".claude", "skills") {
			t.Errorf("harness incorreto: %+v", hs[0])
		}
		if got := processo.DetectarHarnesses(processo.OSFS(), t.TempDir()); len(got) != 0 {
			t.Errorf("HOME vazio não tem harness: %v", got)
		}
	})

	t.Run("prévia lista os 3 e não toca o disco", func(t *testing.T) {
		home := homeCom(t, "claude", "gemini", "codex")
		hs := processo.DetectarHarnesses(processo.OSFS(), home)
		origem := skillNoCofre(t)
		vs := processo.PlanejarVinculos(hs, "legoloop-x", origem)
		if len(vs) != 3 {
			t.Fatalf("a prévia deve listar 3 vínculos: %v", vs)
		}
		for _, v := range vs {
			if v.Origem != origem || filepath.Base(v.Destino) != "legoloop-x" {
				t.Errorf("vínculo incorreto: %+v", v)
			}
			if _, err := os.Lstat(v.Destino); err == nil {
				t.Errorf("planejar não pode criar nada: %s", v.Destino)
			}
		}
	})

	t.Run("cria 3 links que apontam ao cofre", func(t *testing.T) {
		home := homeCom(t, "claude", "gemini", "codex")
		origem := skillNoCofre(t)
		vs := processo.PlanejarVinculos(processo.DetectarHarnesses(processo.OSFS(), home), "legoloop-x", origem)
		feitos, err := processo.Vincular(vs, nil)
		if err != nil || len(feitos) != 3 {
			t.Fatalf("Vincular: %v (%v)", err, feitos)
		}
		for _, d := range feitos {
			if _, err := os.ReadFile(filepath.Join(d, "SKILL.md")); err != nil {
				t.Errorf("o vínculo deve expor o SKILL.md do cofre: %v", err)
			}
		}
	})

	t.Run("falha no 2o link desfaz o 1o", func(t *testing.T) {
		home := homeCom(t, "claude", "gemini", "codex")
		origem := skillNoCofre(t)
		vs := processo.PlanejarVinculos(processo.DetectarHarnesses(processo.OSFS(), home), "legoloop-x", origem)
		chamadas := 0
		quebra := errors.New("disco cheio")
		link := func(o, d string) error {
			chamadas++
			if chamadas == 2 {
				return quebra
			}
			return processo.SymlinkOuCopia(o, d)
		}
		feitos, err := processo.Vincular(vs, link)
		if !errors.Is(err, quebra) || len(feitos) != 0 {
			t.Fatalf("esperava a falha e nenhum vínculo: %v %v", err, feitos)
		}
		for _, v := range vs {
			if _, err := os.Lstat(v.Destino); err == nil {
				t.Errorf("o vínculo %s devia ter sido desfeito", v.Destino)
			}
		}
		if _, err := os.Stat(filepath.Join(origem, "SKILL.md")); err != nil {
			t.Error("desfazer o link não pode apagar a skill do cofre")
		}
	})

	t.Run("link já correto é reaproveitado e sobrevive à falha seguinte", func(t *testing.T) {
		home := homeCom(t, "claude", "gemini")
		origem := skillNoCofre(t)
		vs := processo.PlanejarVinculos(processo.DetectarHarnesses(processo.OSFS(), home), "legoloop-x", origem)
		if _, err := processo.Vincular(vs[:1], nil); err != nil {
			t.Fatal(err)
		}
		falha := errors.New("x")
		_, err := processo.Vincular(vs, func(o, d string) error { return falha })
		if !errors.Is(err, falha) {
			t.Fatalf("esperava a falha: %v", err)
		}
		if _, err := os.Lstat(vs[0].Destino); err != nil {
			t.Error("o link que já existia antes não é nosso para desfazer")
		}
	})

	t.Run("destino ocupado por outra coisa é recusado", func(t *testing.T) {
		home := homeCom(t, "claude")
		origem := skillNoCofre(t)
		vs := processo.PlanejarVinculos(processo.DetectarHarnesses(processo.OSFS(), home), "legoloop-x", origem)
		if err := os.MkdirAll(vs[0].Destino, 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := processo.Vincular(vs, nil); err == nil {
			t.Fatal("não pode sobrescrever o que já existe")
		}
	})
}
