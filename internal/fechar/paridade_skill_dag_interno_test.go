package fechar

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestAuditParidade_DAGComJSONInvalidoNaoEhStub prova que o JSON inválido do
// DAG vira erro próprio do Gate de paridade, e nunca ErrArtefatoStub.
func TestAuditParidade_DAGComJSONInvalidoNaoEhStub(t *testing.T) {
	dir := t.TempDir()
	dagPath := filepath.Join(dir, "dag.json")
	if err := os.WriteFile(dagPath, []byte("{não é JSON"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := auditParidadeSkillDAG(dir, "mudo", "legoloop-mudo", dagPath)
	if err == nil {
		t.Fatal("DAG com JSON inválido deve falhar o Gate")
	}
	if errors.Is(err, ErrArtefatoStub) {
		t.Errorf("erro não pode envolver ErrArtefatoStub: %v", err)
	}
	if !errors.Is(err, errParidadeDAGInvalido) {
		t.Errorf("erro deve envolver errParidadeDAGInvalido: %v", err)
	}
}
