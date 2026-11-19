package brain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAnticorposMarkdownPersistencia prova que a persistência emite o espelho executivo ANTICORPOS.md.
func TestAnticorposMarkdownPersistencia(t *testing.T) {
	dir := t.TempDir()

	db, err := NewAntibodiesDBPersistente(dir)
	if err != nil {
		t.Fatalf("NewAntibodiesDBPersistente: %v", err)
	}
	if err := db.Record("pkgA", "E1", "causa raiz", "diff solucao"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	caminhoMD := filepath.Join(dir, AntibodiesMarkdownFile)
	raw, err := os.ReadFile(caminhoMD)
	if err != nil {
		t.Fatalf("arquivo %s não gerado: %v", AntibodiesMarkdownFile, err)
	}

	conteudo := string(raw)
	if !strings.HasPrefix(conteudo, "---\n") {
		t.Fatalf("esperava cabeçalho YAML com delimitador '---', obteve:\n%s", conteudo)
	}
	if !strings.Contains(conteudo, "schema: anticorpos.v1") {
		t.Errorf("metadado schema ausente no YAML:\n%s", conteudo)
	}
	if !strings.Contains(conteudo, "total_anticorpos: 1") {
		t.Errorf("total_anticorpos incorreto no YAML:\n%s", conteudo)
	}
	if !strings.Contains(conteudo, "pkgA") {
		t.Errorf("pacote pkgA ausente no markdown:\n%s", conteudo)
	}
	if !strings.Contains(conteudo, "diff solucao") {
		t.Errorf("antídoto ausente no markdown:\n%s", conteudo)
	}
}

// TestSincronizarProcesso prova a emissão do espelho a partir de um JSON existente.
func TestSincronizarProcesso(t *testing.T) {
	dir := t.TempDir()
	jsonDummy := `{
  "schema_version": "anticorpos.v1",
  "processo": "meu-processo",
  "globais": [],
  "dominio": [
    {
      "assinatura": "sha256:abc",
      "erro": "falha na validacao",
      "causa_raiz": "parametro nulo",
      "antidoto": "checar ponteiro antes",
      "recorrencias": 1,
      "pacote": "validador/ponteiro",
      "timestamp": "2026-11-19T10:00:00Z",
      "antidoto_aplicado": true
    }
  ]
}`
	if err := os.WriteFile(filepath.Join(dir, AntibodiesDBFile), []byte(jsonDummy), 0o644); err != nil {
		t.Fatalf("gravando dummy: %v", err)
	}

	if err := SincronizarProcesso(dir); err != nil {
		t.Fatalf("SincronizarProcesso: %v", err)
	}

	caminhoMD := filepath.Join(dir, AntibodiesMarkdownFile)
	raw, err := os.ReadFile(caminhoMD)
	if err != nil {
		t.Fatalf("arquivo %s não foi criado pelo sincronizador: %v", AntibodiesMarkdownFile, err)
	}

	conteudo := string(raw)
	if !strings.Contains(conteudo, "processo: meu-processo") {
		t.Errorf("processo ausente no YAML:\n%s", conteudo)
	}
	if !strings.Contains(conteudo, "validador/ponteiro") {
		t.Errorf("pacote ausente no markdown:\n%s", conteudo)
	}
}
