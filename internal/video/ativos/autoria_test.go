package ativos

import (
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

func TestEnriquecerStoryboardGeraTextosFocoCamadasECartoesEAlternaHeroi(t *testing.T) {
	linhas := []roteiro.LinhaRoteiro{
		{Text: "Relatos revelam semanas de 90 horas nas equipes de software.", Natureza: "alegacao_terceiro", Atribuicao: "BBC"},
		{Text: "Dados oficiais mostram alta de 42% na pressao por entregas.", Natureza: "fato_consolidado"},
		{Text: "Estudos apontam queda no descanso de engenheiros.", Natureza: "evidencia_preliminar"},
		{Text: "A esteira continua de codigo dobrou as metas mensais.", Natureza: "narrativa_transicao"},
		{Text: "Por que tanta automacao gerou mais trabalho?", Natureza: "narrativa_transicao"},
	}
	cenarios := []Cenario{{ID: "servidor"}, {ID: "escritorio-noite"}, {ID: "sala-reuniao"}}
	cenas := CompilarCenasAutorais(linhas, cenarios, "narrador_pose_1", []string{"dev", "analista"})
	if len(cenas) != 5 {
		t.Fatalf("esperava 5 cenas, obteve %d", len(cenas))
	}
	if cenas[0].TextoTela == "" || len(cenas[0].Camadas) == 0 {
		t.Fatalf("linha 0 deveria ter texto foco (90 horas) e camadas (stamp/stopwatch): %+v", cenas[0])
	}
	if cenas[1].TextoTela == "" || len(cenas[1].Camadas) == 0 {
		t.Fatalf("linha 1 deveria destacar 42%% e camada doc/bars: %+v", cenas[1])
	}
	// Regra AUTORIA.md: no maximo 3 cenas seguidas com personagem em quadro; a 4a (indice 3) nao traz personagem_id.
	temHeroi3 := false
	for _, el := range cenas[3].ElementosMoveis {
		if el.PersonagemID != "" {
			temHeroi3 = true
		}
	}
	if temHeroi3 {
		t.Fatalf("4a cena seguida nao deveria ter personagem para respeitar teto de 3 seguidas: %+v", cenas[3])
	}
	if !cenas[4].Cartao {
		t.Fatalf("pergunta curta na linha 4 deveria virar cartao digitado: %+v", cenas[4])
	}
}
