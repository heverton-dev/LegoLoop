package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func indexarRank(t *testing.T, skills map[string]string) *SkillEngine {
	t.Helper()
	root := t.TempDir()
	for nome, conteudo := range skills {
		dir := filepath.Join(root, nome)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	e := NewSkillEngine()
	if err := e.IndexDirectory(root); err != nil {
		t.Fatal(err)
	}
	return e
}

func fmSkill(nome, desc string) string {
	return "---\nname: " + nome + "\ndescription: " + desc + "\n---\ncorpo\n"
}

// TestRankingNomePesaMaisQueDescricao reproduz o achado de 20/11: "aidd bigbrain orquestra tudo"
// devolvia aidd-retro primeiro e não trazia aidd-bigbrain no topo.
func TestRankingNomePesaMaisQueDescricao(t *testing.T) {
	skills := map[string]string{
		"aidd-bigbrain": fmSkill("aidd-bigbrain", "Meta-orchestrator across all maestros."),
	}
	for i := 0; i < 15; i++ {
		n := fmt.Sprintf("aidd-decoy%02d", i)
		skills[n] = fmSkill(n, "orquestra tudo para o aidd e para o time inteiro")
	}
	e := indexarRank(t, skills)
	got := e.Search("aidd bigbrain orquestra tudo", 3)
	if len(got) == 0 || got[0].ID != "aidd-bigbrain" {
		t.Fatalf("topo = %+v, quer aidd-bigbrain primeiro (nome vale mais que descrição)", got)
	}
}

// TestRankingNomeExatoVenceSempre prova que a query igual ao nome da skill a põe em primeiro.
func TestRankingNomeExatoVenceSempre(t *testing.T) {
	skills := map[string]string{
		"legoloop-matriz": fmSkill("legoloop-matriz", "Gera a matriz de peças."),
	}
	for i := 0; i < 10; i++ {
		n := fmt.Sprintf("legoloop-outra%02d", i)
		skills[n] = fmSkill(n, "usa legoloop matriz de peças em muitos passos")
	}
	e := indexarRank(t, skills)
	got := e.Search("legoloop-matriz", 5)
	if len(got) == 0 || got[0].ID != "legoloop-matriz" {
		t.Fatalf("topo = %+v, quer legoloop-matriz", got)
	}
}

// TestRankingDeterministicoEOrdenado prova a ordem por score decrescente, desempate por id, e topK.
func TestRankingDeterministicoEOrdenado(t *testing.T) {
	skills := map[string]string{}
	for i := 0; i < 12; i++ {
		n := fmt.Sprintf("skill-%02d", i)
		skills[n] = fmSkill(n, "grafo semantico")
	}
	skills["grafo-forte"] = fmSkill("grafo-forte", "grafo semantico completo")
	e := indexarRank(t, skills)
	primeiro := e.Search("grafo semantico", 5)
	if len(primeiro) != 5 {
		t.Fatalf("topK = %d, quer 5", len(primeiro))
	}
	for rodada := 0; rodada < 20; rodada++ {
		outro := e.Search("grafo semantico", 5)
		for i := range primeiro {
			if outro[i].ID != primeiro[i].ID {
				t.Fatalf("rodada %d: posição %d mudou de %s para %s (ordem não determinística)", rodada, i, primeiro[i].ID, outro[i].ID)
			}
		}
	}
	for i := 1; i < len(primeiro); i++ {
		if primeiro[i].Score > primeiro[i-1].Score {
			t.Fatalf("scores fora de ordem: %+v", primeiro)
		}
		if primeiro[i].Score == primeiro[i-1].Score && primeiro[i].ID < primeiro[i-1].ID {
			t.Fatalf("empate sem desempate por id: %+v", primeiro)
		}
	}
}

// TestDescricaoMultilinhaYAML prova que `description: |` não vira o texto "|".
func TestDescricaoMultilinhaYAML(t *testing.T) {
	conteudo := "---\nname: gcp-multilinha\ndescription: |\n  Provisiona recursos de infraestrutura\n  com verificação de segurança.\n---\ncorpo\n"
	e := indexarRank(t, map[string]string{"gcp-multilinha": conteudo})
	got := e.Search("provisiona infraestrutura", 3)
	if len(got) == 0 || got[0].ID != "gcp-multilinha" {
		t.Fatalf("busca pela descrição multilinha = %+v, quer gcp-multilinha", got)
	}
	if got[0].Resumo == "|" || got[0].Resumo == "" {
		t.Fatalf("resumo = %q, quer o texto multilinha", got[0].Resumo)
	}
}

// TestRuidoDeMetadadosNaoPontua prova que "Processo Canônico" e "general" (padrões) não entram no placar.
func TestRuidoDeMetadadosNaoPontua(t *testing.T) {
	e := indexarRank(t, map[string]string{"alfa": fmSkill("alfa", "faz coisa específica")})
	if got := e.Search("processo canônico general", 5); len(got) != 0 {
		t.Fatalf("metadados padrão pontuaram: %+v", got)
	}
}

// TestPalavraVaziaNaoCasaDentroDeOutraPalavra reproduz "revisar código com TDD":
// "com" casava dentro de composer e components e empurrava aidd-tdd para o 2º lugar.
func TestPalavraVaziaNaoCasaDentroDeOutraPalavra(t *testing.T) {
	e := indexarRank(t, map[string]string{
		"aidd-tdd":        fmSkill("aidd-tdd", "Ciclo Red-Green-Refactor estrito."),
		"aidd-components": fmSkill("aidd-components", "Biblioteca de telas visuais."),
		"gcp-composer":    fmSkill("gcp-composer", "Orquestra DAGs no Composer."),
	})
	got := e.Search("revisar código com TDD", 3)
	if len(got) == 0 || got[0].ID != "aidd-tdd" {
		t.Fatalf("topo = %+v, quer aidd-tdd (a palavra vazia \"com\" não pontua)", got)
	}
}

// TestRaizComumAchaFormaVizinha prova que "homologar" acha "homologue" (mesma raiz).
func TestRaizComumAchaFormaVizinha(t *testing.T) {
	e := indexarRank(t, map[string]string{
		"aidd-homologue": fmSkill("aidd-homologue", "Deep full-stack E2E acceptance."),
		"aidd-outra":     fmSkill("aidd-outra", "Faz outra coisa."),
	})
	got := e.Search("homologar aplicação", 3)
	if len(got) == 0 || got[0].ID != "aidd-homologue" {
		t.Fatalf("busca por homologar = %+v, quer aidd-homologue", got)
	}
}
