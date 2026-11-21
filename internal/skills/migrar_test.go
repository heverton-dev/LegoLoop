package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func montarHarnessEager(t *testing.T) (home, vault string) {
	t.Helper()
	home = t.TempDir()
	vault = filepath.Join(home, ".aidd", "canonical")
	pasta := filepath.Join(home, ".cline", "skills")
	gravaSkill(t, pasta, "alfa", fm("alfa"))                          // já no Cofre
	gravaSkill(t, pasta, "beta", fm("beta"))                          // falta no Cofre
	gravaSkill(t, pasta, "quebrada", "sem frontmatter")               // inválida: só backup
	gravaSkill(t, filepath.Join(pasta, "pacote"), "gama", fm("gama")) // pacote sem SKILL.md no topo
	gravaSkill(t, pasta, "velha.disabled", fm("velha"))               // disabled: só backup
	gravaSkill(t, pasta, ".oculta", fm("oculta"))                     // oculta: fica
	gravaSkill(t, filepath.Join(vault, "skills"), "alfa", fm("alfa"))
	return home, vault
}

// TestMigrarSimulaSemTocarNoDisco prova a contagem do dry-run e que nada muda.
func TestMigrarSimulaSemTocarNoDisco(t *testing.T) {
	home, vault := montarHarnessEager(t)
	rel, err := Migrar(MigrarOpcao{Home: home, Vault: vault, Data: "t"})
	if err != nil {
		t.Fatal(err)
	}
	var h *MigrarHarness
	for i := range rel.Harnesses {
		if rel.Harnesses[i].Nome == "cline" {
			h = &rel.Harnesses[i]
		}
	}
	if h == nil || h.Mover != 5 || h.Copiar != 2 || h.Invalidas != 1 {
		t.Fatalf("cline = %+v, quer mover=5 (alfa, beta, quebrada, pacote, velha.disabled), copiar=2 (beta, gama), invalidas=1", h)
	}
	if _, err := os.Stat(filepath.Join(home, ".cline", "skills", "beta", "SKILL.md")); err != nil {
		t.Fatal("dry-run não pode mover nada")
	}
}

// TestMigrarAplicaCopiaMoveEDeixaIndice prova o efeito: Cofre completo, backup com manifesto, pasta só com o índice.
func TestMigrarAplicaCopiaMoveEDeixaIndice(t *testing.T) {
	home, vault := montarHarnessEager(t)
	rel, err := Migrar(MigrarOpcao{Home: home, Vault: vault, Data: "t", Aplicar: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, nome := range []string{"beta", "gama"} {
		if _, err := os.Stat(filepath.Join(vault, "skills", nome, "SKILL.md")); err != nil {
			t.Errorf("skill %q deveria ter ido ao Cofre: %v", nome, err)
		}
	}
	if _, err := os.Stat(filepath.Join(vault, "skills", "quebrada")); err == nil {
		t.Error("skill inválida não deve entrar no Cofre")
	}
	pasta := filepath.Join(home, ".cline", "skills")
	restantes := map[string]bool{}
	entradas, _ := os.ReadDir(pasta)
	for _, e := range entradas {
		restantes[e.Name()] = true
	}
	if len(restantes) != 2 || !restantes["legoloop-skills"] || !restantes[".oculta"] {
		t.Fatalf("pasta do harness = %v, quer só o índice e a oculta", restantes)
	}
	man, err := os.ReadFile(filepath.Join(rel.Backup, "cline", "MANIFESTO.tsv"))
	if err != nil || strings.Count(string(man), "\n") != 5 {
		t.Fatalf("manifesto = %q (%v), quer 5 linhas", man, err)
	}
	if _, err := os.Stat(filepath.Join(rel.Backup, "cline", "quebrada", "SKILL.md")); err != nil {
		t.Error("a skill inválida deve estar no backup")
	}
}

// TestMigrarEIdempotente prova que repetir a migração não mexe em nada.
func TestMigrarEIdempotente(t *testing.T) {
	home, vault := montarHarnessEager(t)
	if _, err := Migrar(MigrarOpcao{Home: home, Vault: vault, Data: "t", Aplicar: true}); err != nil {
		t.Fatal(err)
	}
	rel, err := Migrar(MigrarOpcao{Home: home, Vault: vault, Data: "t", Aplicar: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range rel.Harnesses {
		if h.Mover != 0 || h.Copiar != 0 {
			t.Errorf("segunda rodada em %s: %+v, quer nada a fazer", h.Nome, h)
		}
	}
}

// TestMigrarNaoMoveRepositorioDeFerramenta prova a proteção: pasta sem SKILL.md no topo que é um
// repositório de ferramenta (.git ou manifesto de build) com skills aninhadas fica onde está.
func TestMigrarNaoMoveRepositorioDeFerramenta(t *testing.T) {
	home := t.TempDir()
	vault := filepath.Join(home, ".aidd", "canonical")
	pasta := filepath.Join(home, ".cline", "skills")
	gravaSkill(t, filepath.Join(pasta, "ferramenta", "skills"), "util", fm("util"))
	if err := os.MkdirAll(filepath.Join(pasta, "ferramenta", ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	gravaSkill(t, filepath.Join(pasta, "outra", "skills"), "util2", fm("util2"))
	if err := os.WriteFile(filepath.Join(pasta, "outra", "Cargo.toml"), []byte("[package]"), 0o644); err != nil {
		t.Fatal(err)
	}
	gravaSkill(t, filepath.Join(pasta, "pacote-puro"), "pura", fm("pura"))
	if _, err := Migrar(MigrarOpcao{Home: home, Vault: vault, Data: "t", Aplicar: true}); err != nil {
		t.Fatal(err)
	}
	for _, repo := range []string{"ferramenta", "outra"} {
		if _, err := os.Stat(filepath.Join(pasta, repo)); err != nil {
			t.Errorf("repositório %q não pode sair da pasta: %v", repo, err)
		}
	}
	if _, err := os.Stat(filepath.Join(pasta, "pacote-puro")); err == nil {
		t.Error("pacote só de skills deveria ter sido migrado")
	}
}
