package fabricar

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSubprocessoCuradoriaFronteira_SeisEtapas(t *testing.T) {
	etapas := SubprocessoCuradoriaFronteira()
	if len(etapas) != 6 {
		t.Fatalf("esperado 6 etapas, obtido %d", len(etapas))
	}

	esperados := []struct {
		id     string
		trecho string
	}{
		{"entrevista_escopo", "0 a 10"},
		{"inventario_local", "Matriz de Peças"},
		{"curadoria_externa", "legoloop import"},
		{"elevacao_canonizacao", "legoloop upgrade --diagnose"},
		{"composicao_dag", "context_in"},
		{"projecao_padrao", "legoloop validate"},
	}

	for i, esp := range esperados {
		if etapas[i].Ordem != i+1 {
			t.Errorf("etapa %d: Ordem = %d, esperado %d", i, etapas[i].Ordem, i+1)
		}
		if etapas[i].ID != esp.id {
			t.Errorf("etapa %d: ID = %q, esperado %q", i, etapas[i].ID, esp.id)
		}
		if !strings.Contains(etapas[i].Descricao, esp.trecho) {
			t.Errorf("etapa %s: Descricao %q não contém %q", esp.id, etapas[i].Descricao, esp.trecho)
		}
	}
}

func TestFabricarArtifacts_ContemCuradoriaFronteira(t *testing.T) {
	root := filepath.Join("..", "..")
	skillBytes, err := os.ReadFile(filepath.Join(root, "skills", "legoloop-fabricar", "SKILL.md"))
	if err != nil {
		t.Fatalf("ler SKILL.md: %v", err)
	}
	skill := string(skillBytes)
	for _, termo := range []string{"Curadoria e Composição de Fronteira", "legoloop import", "legoloop upgrade --diagnose", "legoloop validate"} {
		if !strings.Contains(skill, termo) {
			t.Errorf("SKILL.md não contém %q", termo)
		}
	}
}

func TestAvaliarStopGateCuradoria_ParaSemYesEProssegueComYes(t *testing.T) {
	sugestoes := []string{"https://github.com/exemplo/repo-a"}
	summary, err := AvaliarStopGateCuradoria(Request{RawName: "teste", Interactive: true}, "resumo do processo", sugestoes)
	if err == nil || !strings.Contains(err.Error(), "Stop-Gate") {
		t.Fatalf("esperado ErrStopGatePending sem --yes, obtido: %v", err)
	}
	if !summary.AwaitingApproval || summary.ProcessoResumo != "resumo do processo" || len(summary.SugestoesGitHub) != 1 {
		t.Fatalf("resumo do Stop-Gate inválido: %+v", summary)
	}
	summaryYes, err := AvaliarStopGateCuradoria(Request{RawName: "teste", Interactive: false}, "resumo do processo", sugestoes)
	if err != nil {
		t.Fatalf("com --yes não deveria parar no Stop-Gate: %v", err)
	}
	if summaryYes.AwaitingApproval {
		t.Fatalf("com --yes AwaitingApproval deve ser false")
	}
}
