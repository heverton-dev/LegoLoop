package skills_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/cli"
	"legoloop/internal/processo"
)

const dagUmNo = `{"name":"p","nodes":[
 {"id":"a","type":"script","command":["go","version"],"timeout_sec":60,
  "routing":{"mode":"manual","model":"none"},
  "output":{"artifact":"GO_VERSION.txt","gate":"exit_code == 0"}}]}`

// montarCenario cria um projeto com um processo, um HOME com harnesses
// simulados e devolve os caminhos usados pelo teste.
func montarCenario(t *testing.T, harnesses []string) (base, dirProcessos, home, vault string) {
	t.Helper()
	base = t.TempDir()
	dirProcessos = filepath.Join(base, caminhos.DirProcessos)
	procDir := filepath.Join(dirProcessos, "meu-proc")
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminhos.CaminhoDAG(procDir), []byte(dagUmNo), 0o644); err != nil {
		t.Fatal(err)
	}
	home = filepath.Join(base, "home")
	for _, h := range harnesses {
		if err := os.MkdirAll(filepath.Join(home, "."+h), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return base, dirProcessos, home, filepath.Join(base, "vault")
}

func rodarInstalar(t *testing.T, home, vault, dirProcessos string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	d := cli.Deps{Stdout: &out, Stderr: &errb, Home: home}
	code := cli.Run(context.Background(), []string{
		"instalar-skills", "--vault", vault, "--processos-dir", dirProcessos,
	}, d)
	return code, out.String(), errb.String()
}

// TestDistribuiTodosHarness (AC-E07.2): dado processo criado, a skill
// skills/<slug>/SKILL.md e o slash command /legoloop-<slug> existem e o
// legoloop instalar-skills os distribui para todos os harnesses, sem passo manual.
func TestDistribuiTodosHarness(t *testing.T) {
	base, dirProcessos, home, vault := montarCenario(t, []string{"claude", "codex", "gemini"})
	code, out, errb := rodarInstalar(t, home, vault, dirProcessos)
	if code != 0 {
		t.Fatalf("instalar-skills code=%d out=%s err=%s", code, out, errb)
	}

	// 1. A skill do processo existe no projeto, com front matter e contrato de execução.
	skProj := filepath.Join(base, "skills", "meu-proc", "SKILL.md")
	dados, err := os.ReadFile(skProj)
	if err != nil {
		t.Fatalf("skills/meu-proc/SKILL.md não existe: %v", err)
	}
	if !strings.Contains(string(dados), "name: meu-proc") || !strings.Contains(string(dados), "description:") {
		t.Errorf("front matter incompleto em %s", skProj)
	}
	if !strings.Contains(string(dados), "legoloop exec meu-proc") {
		t.Errorf("skill não delega a legoloop exec: %s", skProj)
	}

	// 2. O slash command /legoloop-<slug> existe no projeto.
	cmdProj := filepath.Join(base, "skills", "meu-proc", "comandos", "legoloop-meu-proc.md")
	dadosCmd, err := os.ReadFile(cmdProj)
	if err != nil {
		t.Fatalf("slash command não existe: %v", err)
	}
	if !strings.Contains(string(dadosCmd), "legoloop exec meu-proc") {
		t.Errorf("slash command não delega a legoloop exec: %s", cmdProj)
	}

	// 3. Tudo distribuído para TODOS os harnesses detectados.
	hs := processo.DetectarHarnesses(processo.OSFS(), home)
	if len(hs) != 3 {
		t.Fatalf("esperava 3 harnesses, obteve %d", len(hs))
	}
	if vinculados := processo.HarnessesVinculados(processo.OSFS(), hs, "meu-proc"); len(vinculados) != 3 {
		t.Errorf("skill distribuída em %v, esperava os 3 harnesses", vinculados)
	}
	for _, h := range hs {
		dest := filepath.Join(filepath.Dir(h.Skills), "commands", "legoloop-meu-proc.md")
		if _, err := os.Stat(dest); err != nil {
			t.Errorf("slash command ausente no harness %s: %v", h.Nome, err)
		}
	}

	// 4. Regressão: as skills nativas continuam distribuídas.
	if vinculados := processo.HarnessesVinculados(processo.OSFS(), hs, "legoloop-fabricar"); len(vinculados) != 3 {
		t.Errorf("skill nativa distribuída em %v, esperava os 3 harnesses", vinculados)
	}
}

// TestSkillInvalidoExit (AC-E07.4): processo com SKILL.md inválido (front
// matter quebrado) faz instalar-skills sair com exit ≠ 0 apontando o arquivo,
// sem corromper nenhum harness.
func TestSkillInvalidoExit(t *testing.T) {
	base, dirProcessos, home, vault := montarCenario(t, []string{"claude", "gemini"})

	// Quebra o front matter da skill do processo e cria um segundo processo válido.
	skProj := filepath.Join(base, "skills", "meu-proc", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(skProj), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skProj, []byte("---\nname: meu-proc\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	procBom := filepath.Join(dirProcessos, "bom")
	if err := os.MkdirAll(procBom, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminhos.CaminhoDAG(procBom), []byte(dagUmNo), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, errb := rodarInstalar(t, home, vault, dirProcessos)
	if code == 0 {
		t.Fatal("esperava exit ≠ 0 com SKILL.md inválido")
	}
	if !strings.Contains(errb, skProj) {
		t.Errorf("erro não aponta o arquivo inválido: %s", errb)
	}

	// Nenhum harness corrompido: nada do processo chegou a nenhum harness.
	hs := processo.DetectarHarnesses(processo.OSFS(), home)
	for _, h := range hs {
		for _, nome := range []string{"meu-proc", "bom"} {
			if _, err := os.Stat(filepath.Join(h.Skills, nome)); err == nil {
				t.Errorf("harness %s recebeu a skill %s apesar do erro", h.Nome, nome)
			}
		}
		if _, err := os.Stat(filepath.Join(filepath.Dir(h.Skills), "commands", "legoloop-meu-proc.md")); err == nil {
			t.Errorf("harness %s recebeu o slash command apesar do erro", h.Nome)
		}
	}
}
