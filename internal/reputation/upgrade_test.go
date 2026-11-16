package reputation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpgradePieceFile_MarkdownToGold(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	orig := "# Minha Skill\n\nInstrucao simples sem contratos.\n"
	if err := os.WriteFile(path, []byte(orig), 0o600); err != nil {
		t.Fatal(err)
	}

	before, err := ScoreFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if before.Tier != TierToxica {
		t.Fatalf("esperava tier inicial toxica, obteve %s (%.2f)", before.Tier, before.TotalScore)
	}

	rep, err := UpgradePieceFile(path)
	if err != nil {
		t.Fatalf("UpgradePieceFile falhou: %v", err)
	}
	if !rep.Changed || rep.After.Tier != TierOuro || rep.After.TotalScore < 0.85 {
		t.Fatalf("esperava elevacao para ouro >= 0.85, obteve %+v", rep)
	}
	if _, err := os.Stat(rep.BackupPath); err != nil {
		t.Fatalf("esperava backup criado em %s: %v", rep.BackupPath, err)
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(updated), "Instrucao simples sem contratos.") {
		t.Fatalf("conteudo original nao preservado:\n%s", string(updated))
	}

	// Segunda chamada deve ser idempotente (ja e ouro).
	rep2, err := UpgradePieceFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if rep2.Changed {
		t.Fatalf("peca ja ouro nao deveria ser reescrita: %+v", rep2)
	}
}

func TestUpgradePieceFile_ScriptPreservesHashComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "check.py")
	orig := "print('ok')\n"
	if err := os.WriteFile(path, []byte(orig), 0o600); err != nil {
		t.Fatal(err)
	}

	rep, err := UpgradePieceFile(path)
	if err != nil {
		t.Fatalf("UpgradePieceFile falhou: %v", err)
	}
	if rep.After.Tier != TierOuro || rep.After.TotalScore < 0.85 {
		t.Fatalf("script deveria atingir ouro, obteve %+v", rep.After)
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), "print('ok')") {
		t.Fatalf("codigo python original perdido:\n%s", string(raw))
	}
}
