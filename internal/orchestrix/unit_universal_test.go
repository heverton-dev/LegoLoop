package orchestrix

import (
	"context"
	"os"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
	"legoloop/internal/roteamento"
)

// semClaude monta um motor sem o claude CLI no inventário e com o omp homologado como candidato da classe.
func semClaude(avisos *[]string) *Orchestrix {
	o := &Orchestrix{
		Inventory:      provider.Inventory{CLIs: []provider.CLI{{Name: "omp", Path: os.Args[0]}}},
		ExtraEnv:       []string{"LEGOLOOP_HELPER=omp"},
		AllowUnitTasks: true,
		EstadoHarness: func(harness, via string) string {
			if harness == "omp" && via == "oficial" {
				return "homologado"
			}
			return ""
		},
		Elegiveis: map[string][]roteamento.CandidatoVia{
			roteamento.ClasseNano:  {{Harness: "omp", Via: roteamento.ViaOficial, Modelo: "omp:x"}},
			roteamento.ClasseMicro: {{Harness: "omp", Via: roteamento.ViaOficial, Modelo: "omp:x"}},
		},
	}
	o.Aviso = func(m string) { *avisos = append(*avisos, m) }
	return o
}

func noClaude(id, classe string) dag.Node {
	n := unitNode(id, "out/"+id+".json", "exit_code == 0", nil)
	n.ContextOut = nil
	n.Nature = classe
	n.Routing = dag.RoutingConfig{Mode: dag.ModeSuggested, Harness: "claude", Via: "oficial", Model: "sonnet"}
	return n
}

// Sem o claude instalado, um nó nano ou micro não quebra: cai no candidato homologado da classe.
func TestNoDeCodigoSemClaudeDegradaNanoMicro(t *testing.T) {
	for _, classe := range []string{roteamento.ClasseNano, roteamento.ClasseMicro} {
		var avisos []string
		repo := newRepoT(t)
		res, err := semClaude(&avisos).Execute(context.Background(), noClaude("n-"+classe, classe), repo)
		if err != nil {
			t.Fatalf("%s: o nó não pode quebrar sem o claude: %v", classe, err)
		}
		notas := strings.Join(res.Notes, "\n")
		if !strings.Contains(notas, "usou omp/oficial/omp:x no lugar de claude/oficial/sonnet") {
			t.Errorf("%s: sem a Substituicao nas notas: %q", classe, notas)
		}
		if len(avisos) != 1 {
			t.Errorf("%s: esperava 1 aviso, obteve %v", classe, avisos)
		}
	}
}

// Classe acima de micro não troca de harness: falha com mensagem clara, sem worktree vazada.
func TestNoDeCodigoSemClaudeFalhaForaDeNanoMicro(t *testing.T) {
	var avisos []string
	repo := newRepoT(t)
	_, err := semClaude(&avisos).Execute(context.Background(), noClaude("n-g", "epic"), repo)
	if err == nil || !strings.Contains(err.Error(), "claude") {
		t.Fatalf("classe epic sem claude deve falhar citando o claude: %v", err)
	}
	if c := worktreeCount(t, repo); c != 1 {
		t.Fatalf("worktree vazou: %d", c)
	}
}
