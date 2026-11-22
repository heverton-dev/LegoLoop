package forja

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCorrigirDeriva(t *testing.T) {
	home := t.TempDir()
	cofre := filepath.Join(t.TempDir(), "cofre-skills")
	criarSkillDerivaTeste(t, filepath.Join(home, ".claude", "skills", "legoloop-skills"), "legoloop-skills")
	pastaEager := filepath.Join(home, ".claude", "skills", "skill-fora")
	criarSkillDerivaTeste(t, pastaEager, "skill-fora")

	res, err := CorrigirDeriva(home, cofre, "2026-11-22")
	if err != nil {
		t.Fatalf("CorrigirDeriva falhou: %v", err)
	}
	if len(res.Corrigidos) != 1 {
		t.Fatalf("esperava 1 item corrigido, obteve %+v", res)
	}
	validarArtefatosCorrecao(t, home, cofre, pastaEager)
}

func validarArtefatosCorrecao(t *testing.T, home, cofre, pastaEager string) {
	t.Helper()
	if _, err := os.Stat(pastaEager); !os.IsNotExist(err) {
		t.Fatalf("skill eager ainda existe na pasta do harness: %s", pastaEager)
	}
	if _, err := os.Stat(filepath.Join(cofre, "skill-fora", "SKILL.md")); err != nil {
		t.Fatalf("skill nao foi copiada ao cofre: %v", err)
	}
	manifesto := filepath.Join(home, ".aidd", "backups", "deriva_2026-11-22", "claude", "MANIFESTO.tsv")
	dados, err := os.ReadFile(manifesto)
	if err != nil {
		t.Fatalf("MANIFESTO.tsv nao criado em %s: %v", manifesto, err)
	}
	if !strings.Contains(string(dados), "skill-fora") {
		t.Fatalf("MANIFESTO.tsv deve citar skill-fora: %s", string(dados))
	}
	if _, err := os.Stat(filepath.Join(home, ".aidd", "backups", "deriva_2026-11-22", "claude", "skill-fora", "SKILL.md")); err != nil {
		t.Fatalf("backup da skill nao encontrado: %v", err)
	}
	restantes, err := DetectarDeriva(home)
	if err != nil {
		t.Fatalf("DetectarDeriva pos-correcao falhou: %v", err)
	}
	if len(restantes) != 0 {
		t.Fatalf("esperava 0 derivas apos correcao, obteve %+v", restantes)
	}
}

func TestDerivaCorrigir(t *testing.T) {
	TestCorrigirDeriva(t)
}

func TestCorrigirDerivaHomeVazio(t *testing.T) {
	if _, err := CorrigirDeriva("", "", "2026-11-22"); !errors.Is(err, ErrSemHome) {
		t.Fatalf("esperava ErrSemHome, obteve %v", err)
	}
}
