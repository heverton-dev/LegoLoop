package fechar

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func depsGovernanca(t *testing.T) *Deps {
	t.Helper()
	return &Deps{
		RepoDir: t.TempDir(),
		Clock:   func() time.Time { return time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC) },
	}
}

func lerArquivo(t *testing.T, caminho string) string {
	t.Helper()
	b, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("leitura de %s: %v", caminho, err)
	}
	return string(b)
}

// AC: o fechamento grava o Processo canonizado no registro à parte e não no AGENTS.md.
func TestFecharGovernancaGravaRegistroAparte(t *testing.T) {
	d := depsGovernanca(t)
	if _, err := fecharGovernanca(d, "web-qa", "legoloop-web-qa", false); err != nil {
		t.Fatalf("fecharGovernanca: %v", err)
	}
	registro := lerArquivo(t, filepath.Join(d.RepoDir, "docs", "PROCESSOS_CANONIZADOS.md"))
	if !strings.Contains(registro, "`legoloop-web-qa` (2026-11-18)") {
		t.Errorf("registro sem o Processo canonizado: %q", registro)
	}
	agents := lerArquivo(t, filepath.Join(d.RepoDir, "AGENTS.md"))
	if strings.Contains(agents, "legoloop-web-qa") {
		t.Errorf("AGENTS.md não pode receber o Processo canonizado: %q", agents)
	}
	if !strings.Contains(agents, "docs/PROCESSOS_CANONIZADOS.md") {
		t.Errorf("AGENTS.md deve referenciar o registro: %q", agents)
	}
}

// AC: fechar repetidamente não faz o AGENTS.md crescer (limite de 24000 bytes do agy).
func TestFecharGovernancaNaoFazAgentsCrescer(t *testing.T) {
	d := depsGovernanca(t)
	for _, canon := range []string{"legoloop-a", "legoloop-b", "legoloop-c", "legoloop-a"} {
		if _, err := fecharGovernanca(d, canon, canon, false); err != nil {
			t.Fatalf("fecharGovernanca %s: %v", canon, err)
		}
	}
	agentsPath := filepath.Join(d.RepoDir, "AGENTS.md")
	primeiro := lerArquivo(t, agentsPath)
	if _, err := fecharGovernanca(d, "legoloop-d", "legoloop-d", false); err != nil {
		t.Fatal(err)
	}
	if segundo := lerArquivo(t, agentsPath); segundo != primeiro {
		t.Errorf("AGENTS.md cresceu entre fechamentos:\nantes=%q\ndepois=%q", primeiro, segundo)
	}
	registro := lerArquivo(t, filepath.Join(d.RepoDir, "docs", "PROCESSOS_CANONIZADOS.md"))
	if n := strings.Count(registro, "`legoloop-a`"); n != 1 {
		t.Errorf("registro deve ser idempotente: legoloop-a aparece %d vezes", n)
	}
}

// AC: nome que é prefixo de outro Processo canonizado ainda é registrado.
func TestFecharGovernancaDistingueNomePrefixo(t *testing.T) {
	d := depsGovernanca(t)
	for _, canon := range []string{"legoloop-video-rea-relatorio", "legoloop-video-rea"} {
		if _, err := fecharGovernanca(d, canon, canon, false); err != nil {
			t.Fatal(err)
		}
	}
	registro := lerArquivo(t, filepath.Join(d.RepoDir, "docs", "PROCESSOS_CANONIZADOS.md"))
	if !strings.Contains(registro, "`legoloop-video-rea` (") {
		t.Errorf("legoloop-video-rea não foi registrado: %q", registro)
	}
}

// AC: dry-run não escreve nada.
func TestFecharGovernancaDryRunNaoEscreve(t *testing.T) {
	d := depsGovernanca(t)
	if _, err := fecharGovernanca(d, "x", "legoloop-x", true); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"AGENTS.md", "MEMORY.md", filepath.Join("docs", "PROCESSOS_CANONIZADOS.md")} {
		if _, err := os.Stat(filepath.Join(d.RepoDir, rel)); err == nil {
			t.Errorf("dry-run criou %s", rel)
		}
	}
}
