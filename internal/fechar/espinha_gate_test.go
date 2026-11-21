package fechar_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/fechar"
)

func TestFecharRodaEspinhaVerificar(t *testing.T) {
	t.Run("reprova quando a espinha tem violacao", func(t *testing.T) {
		repo := semearRepoComEspinha(t, "alfa")
		if err := os.MkdirAll(filepath.Join(repo, "processos", "sem-prefixo"), 0o755); err != nil {
			t.Fatal(err)
		}
		res, err := fechar.Executar(context.Background(), fechar.Request{Slug: "alfa"}, depsComPortoesReais(repo))
		if err != nil {
			t.Fatalf("Executar falhou: %v", err)
		}
		exigirEtapaComErro(t, res, 8, "espinha")
		if !strings.Contains(res.Etapas[7].Detalhe, "processo-sem-prefixo") {
			t.Errorf("detalhe da etapa 8 deveria citar regra processo-sem-prefixo: %q", res.Etapas[7].Detalhe)
		}
	})

	t.Run("aprova quando a espinha esta integra", func(t *testing.T) {
		repo := semearRepoComEspinha(t, "alfa")
		res, err := fechar.Executar(context.Background(), fechar.Request{Slug: "alfa"}, depsComPortoesReais(repo))
		if err != nil {
			t.Fatalf("Executar falhou: %v", err)
		}
		if !res.Sucesso {
			t.Fatalf("esperava sucesso com espinha integra: %+v", res.Etapas)
		}
		if e8 := etapaDe(t, res, 8); e8.Status != "ok" {
			t.Errorf("etapa 8 deveria passar: %+v", e8)
		}
	})
}

func depsComPortoesReais(repo string) fechar.Deps {
	return fechar.Deps{
		RepoDir:   repo,
		VaultDir:  repo,
		Home:      repo,
		SkipBuild: false,
		Rebuilder: func(string) error { return nil },
		Indexer:   func(string) error { return nil },
	}
}

func semearRepoComEspinha(t *testing.T, slug string) string {
	t.Helper()
	repo := t.TempDir()
	canon := "legoloop-" + slug
	semearDocs(t, repo, slug)
	escrever(t, filepath.Join(repo, "AGENTS.md"), "## 1. Leis\n\n1. Lei um.\n")
	escrever(t, filepath.Join(repo, "internal", "espinha", "stub.go"), "package espinha\n")
	escrever(t, filepath.Join(repo, "internal", "cli", "cli.go"),
		"package cli\nfunc mapaComandos() map[string]func() { return map[string]func(){\"espinha\": nil} }\n")
	escrever(t, filepath.Join(repo, caminhos.DirProcessos, canon, caminhos.ArquivoProcessoDAG), dagEspinhaValido(canon))
	escrever(t, filepath.Join(repo, caminhos.DirProcessos, canon, "ANTICORPOS.md"), "# Anticorpos\n")
	escrever(t, filepath.Join(repo, "espinha", "registro.json"), registroEspinhaValido(canon))
	return repo
}

func dagEspinhaValido(nome string) string {
	return `{
  "name": "` + nome + `",
  "version": "1.0.0",
  "trigger": {"type": "manual"},
  "nodes": [
    {
      "id": "01-alfa",
      "type": "script",
      "command": ["echo", "ok"],
      "isolation": {"executor": "orchestrix_hermetic"},
      "routing": {"mode": "auto", "combo": "padrao"},
      "output": {"artifact": "out/01.txt", "gate": "file_exists('out/01.txt')"}
    }
  ]
}`
}

func registroEspinhaValido(canon string) string {
	return `{
  "schema": "espinha.v1",
  "verbos": [
    {"nome": "espinha", "camada": "espinha", "motor": "internal/espinha", "contrato": "verifica espinha"}
  ],
  "motores": [
    {"pacote": "internal/espinha", "interface": "Verificar", "deterministico": true}
  ],
  "processos": [
    {"slug": "` + canon + `"}
  ],
  "requisitos": [
    {"lei": "Lei 1", "verificador": "internal/espinha", "mecanico": true}
  ]
}`
}
