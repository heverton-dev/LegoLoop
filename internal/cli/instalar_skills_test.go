package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/processo"
	"legoloop/skills"
)

func TestInstalarSkills(t *testing.T) {
	// 1. Manuais das 7 funções devem existir na raiz do repositório
	manuaisEsperados := []string{
		"2026-11-03_MANUAL_FUNCAO_FABRICAR_LEGOLOOP.md",
		"2026-11-06_MANUAL_FUNCAO_MATRIZ_LEGOLOOP.md",
		"2026-11-06_MANUAL_FUNCAO_LISTAR_LEGOLOOP.md",
		"2026-11-06_MANUAL_FUNCAO_REFINAR_LEGOLOOP.md",
		"2026-11-06_MANUAL_FUNCAO_REMOVER_LEGOLOOP.md",
		"2026-11-06_MANUAL_FUNCAO_COMPOR_LEGOLOOP.md",
		"2026-11-06_MANUAL_FUNCAO_LEGOPRESS_LEGOLOOP.md",
	}
	manualDir := filepath.Join("..", "..", "docs", "manuals")
	for _, m := range manuaisEsperados {
		path := filepath.Join(manualDir, m)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("manual esperado %q não encontrado: %v", m, err)
		}
	}

	// 2. Criação do ambiente hermético de teste
	tmp := t.TempDir()
	vault := filepath.Join(tmp, "vault")
	home := filepath.Join(tmp, "home")
	_ = os.MkdirAll(vault, 0o755)

	// Cria 3 harnesses simulados no home
	simulados := []string{"claude", "gemini", "codex"}
	for _, h := range simulados {
		if err := os.MkdirAll(filepath.Join(home, "."+h), 0o755); err != nil {
			t.Fatalf("erro ao criar harness simulado: %v", err)
		}
	}

	// 3. Testa Dry-Run (não deve gravar nada)
	var stdoutDry, stderrDry bytes.Buffer
	dDry := Deps{Stdout: &stdoutDry, Stderr: &stderrDry, Home: home}
	codeDry := Run(context.Background(), []string{"instalar-skills", "--vault", vault, "--dry-run"}, dDry)
	if codeDry != 0 {
		t.Fatalf("esperava exit 0 no dry-run, obteve %d. Stderr: %s", codeDry, stderrDry.String())
	}
	entries, _ := os.ReadDir(filepath.Join(vault, "skills"))
	if len(entries) > 0 {
		t.Fatalf("dry-run gravou %d itens no cofre, esperava 0", len(entries))
	}

	// 4. Executa Instalação Real
	var stdout, stderr bytes.Buffer
	d := Deps{Stdout: &stdout, Stderr: &stderr, Home: home}
	code := Run(context.Background(), []string{"instalar-skills", "--vault", vault}, d)
	if code != 0 {
		t.Fatalf("esperava exit 0 na instalação, obteve %d. Stderr: %s", code, stderr.String())
	}

	// 5. Verifica as 6 funções no cofre e nos 3 harnesses
	fsys := processo.OSFS()
	hs := processo.DetectarHarnesses(fsys, home)
	if len(hs) != 3 {
		t.Fatalf("esperava 3 harnesses detectados, obteve %d", len(hs))
	}

	for _, fn := range skills.FuncoesNativas {
		skillCofre := filepath.Join(vault, "skills", fn, "SKILL.md")
		data, err := os.ReadFile(skillCofre)
		if err != nil {
			t.Errorf("skill %q não foi gravada no cofre: %v", fn, err)
			continue
		}
		conteudo := string(data)
		// Verifica frontmatter canônico
		if !strings.Contains(conteudo, "name: "+fn) {
			t.Errorf("skill %q não contém frontmatter name", fn)
		}
		if !strings.Contains(conteudo, "description:") {
			t.Errorf("skill %q não contém frontmatter description", fn)
		}
		if !strings.Contains(conteudo, "REGRAS NEGATIVAS") {
			t.Errorf("skill %q não contém seção de regras negativas", fn)
		}
		if !strings.Contains(conteudo, "legoloop") {
			t.Errorf("skill %q não delega ao binário legoloop", fn)
		}

		// Verifica vínculos nos 3 harnesses
		vinculados := processo.HarnessesVinculados(fsys, hs, fn)
		if len(vinculados) != 3 {
			t.Errorf("skill %q esperava 3 harnesses vinculados, obteve %v", fn, vinculados)
		}
	}

	// 6. Testa formato JSON
	var stdoutJSON, stderrJSON bytes.Buffer
	dJSON := Deps{Stdout: &stdoutJSON, Stderr: &stderrJSON, Home: home}
	codeJSON := Run(context.Background(), []string{"instalar-skills", "--vault", vault, "--json"}, dJSON)
	if codeJSON != 0 {
		t.Fatalf("esperava exit 0 com --json, obteve %d", codeJSON)
	}
	var res map[string]any
	if err := json.Unmarshal(stdoutJSON.Bytes(), &res); err != nil {
		t.Fatalf("saída com --json inválida: %v", err)
	}
}
