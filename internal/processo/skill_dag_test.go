package processo_test

import (
	"strings"
	"testing"

	"legoloop/internal/caminhos"
)

// A skill gerada manda validar o DAG espelho em docs/dags, que existe e não usa o termo proibido da Lei #15.
// O SPEC aponta o DAG do cofre no layout que o fabricar grava: uma pasta por processo.
func TestProjetar_SkillApontaParaODagQueExiste(t *testing.T) {
	pr := projetar(t, quatro())
	if !strings.Contains(pr.Skill, "legoloop validate docs/dags/"+pr.DAGNome) {
		t.Errorf("a skill deve validar docs/dags/%s:\n%s", pr.DAGNome, pr.Skill)
	}
	if strings.Contains(pr.Skill, ".aidd/canonical") || strings.Contains(pr.Skill, caminhos.DirAntigo+"/") || strings.Contains(pr.Skill, caminhos.DirProcessos+"/") {
		t.Error("a skill não aponta o cofre nem o caminho de processo")
	}
	quer := "~/.aidd/canonical/" + caminhos.DirProcessos + "/legoloop-build-engine/" + caminhos.ArquivoDAG
	if !strings.Contains(pr.Spec, quer) {
		t.Errorf("o SPEC deve apontar a pasta do processo (%s):\n%s", quer, pr.Spec)
	}
}
