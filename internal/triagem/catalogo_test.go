package triagem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCarregarSpecs(t *testing.T) {
	tempDir := t.TempDir()
	proc1 := filepath.Join(tempDir, "proc-1")
	if err := os.MkdirAll(proc1, 0o755); err != nil {
		t.Fatalf("erro ao criar pasta proc-1: %v", err)
	}

	spec1 := GerarTriagemPadrao("proc-1", []string{"etapa1"})
	if err := os.WriteFile(filepath.Join(proc1, "triagem.json"), []byte(`{
		"schema_version": "triagem.v1",
		"processo": "proc-1",
		"regras_deterministicas": [{"padrao_slug": "proc-1", "comandos_cli": ["legoloop proc-1"], "prioridade": 100}],
		"criterio_ambiguidade": {"descricao": "desc", "etapas": ["etapa1"]},
		"pontos_entrada": ["cli"]
	}`), 0o644); err != nil {
		t.Fatalf("erro ao salvar triagem.json: %v", err)
	}

	specs, err := CarregarSpecs(tempDir)
	if err != nil {
		t.Fatalf("CarregarSpecs falhou: %v", err)
	}
	if len(specs) != 1 {
		t.Fatalf("esperava 1 spec, obteve %d", len(specs))
	}
	if specs[0].Processo != spec1.Processo {
		t.Errorf("esperava processo %s, obteve %s", spec1.Processo, specs[0].Processo)
	}

	naoExiste, err := CarregarSpecs(filepath.Join(tempDir, "nao-existe"))
	if err != nil {
		t.Fatalf("erro inesperado para diretório inexistente: %v", err)
	}
	if naoExiste != nil {
		t.Errorf("esperava nil para diretório inexistente")
	}
}

func TestLerSpecInexistenteEInvalida(t *testing.T) {
	tempDir := t.TempDir()
	_, ok, err := LerSpec(filepath.Join(tempDir, "ausente.json"))
	if err != nil || ok {
		t.Fatalf("esperava ok=false e err=nil para arquivo inexistente, obteve ok=%v err=%v", ok, err)
	}

	invalido := filepath.Join(tempDir, "invalido.json")
	if err := os.WriteFile(invalido, []byte(`{`), 0o644); err != nil {
		t.Fatalf("erro ao escrever json invalido: %v", err)
	}
	if _, _, err := LerSpec(invalido); err == nil {
		t.Fatal("esperava erro de json inválido")
	}
}
