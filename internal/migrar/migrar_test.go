package migrar

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/triagem"
)

func TestMigrarLoops_DryRunNaoAlteraDisco(t *testing.T) {
	dir := t.TempDir()
	criarDir(t, filepath.Join(dir, "vazio"))
	criarDAG(t, filepath.Join(dir, "valido"), `{"nodes":[{"id":"n1"},{"id":"n2"}]}`)

	rel, err := MigrarLoops(dir, false)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(rel.Removidos) != 1 || rel.Removidos[0] != "vazio" {
		t.Errorf("esperado removido=vazio, obtido %v", rel.Removidos)
	}
	if len(rel.TriagensGeradas) != 1 || rel.TriagensGeradas[0] != "valido" {
		t.Errorf("esperado triagem=valido, obtido %v", rel.TriagensGeradas)
	}
	if _, err := os.Stat(filepath.Join(dir, "vazio")); err != nil {
		t.Errorf("dry-run não deveria remover diretório vazio")
	}
	if _, err := os.Stat(filepath.Join(dir, "valido", "triagem.json")); err == nil {
		t.Errorf("dry-run não deveria criar triagem.json")
	}
}

func TestMigrarLoops_ApplyRemoveQuebradosEGeraTriagem(t *testing.T) {
	dir := t.TempDir()
	criarDir(t, filepath.Join(dir, "quebrado"))
	_ = os.WriteFile(filepath.Join(dir, "quebrado", "pipeline.dag.json"), []byte("{invalido"), 0644)
	criarDAG(t, filepath.Join(dir, "proc-ok"), `{"nodes":[{"id":"etapa1"}]}`)

	rel, err := MigrarLoops(dir, true)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(rel.Removidos) != 1 || rel.Removidos[0] != "quebrado" {
		t.Errorf("esperado removido=quebrado, obtido %v", rel.Removidos)
	}
	if _, err := os.Stat(filepath.Join(dir, "quebrado")); !os.IsNotExist(err) {
		t.Errorf("apply deveria remover diretório quebrado")
	}
	data, err := os.ReadFile(filepath.Join(dir, "proc-ok", "triagem.json"))
	if err != nil {
		t.Fatalf("triagem.json não criado: %v", err)
	}
	var spec triagem.TriagemSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("triagem.json inválido: %v", err)
	}
	if err := spec.Validar(); err != nil {
		t.Errorf("spec inválida: %v", err)
	}
	if len(spec.CriterioAmbiguidade.Etapas) != 1 || spec.CriterioAmbiguidade.Etapas[0] != "etapa1" {
		t.Errorf("etapas inesperadas: %v", spec.CriterioAmbiguidade.Etapas)
	}
}

func criarDir(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0755); err != nil {
		t.Fatal(err)
	}
}

func criarDAG(t *testing.T, p, conteudo string) {
	t.Helper()
	criarDir(t, p)
	if err := os.WriteFile(filepath.Join(p, "pipeline.dag.json"), []byte(conteudo), 0644); err != nil {
		t.Fatal(err)
	}
}
