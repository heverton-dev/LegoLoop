package fechar_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/fechar"
)

// TestNosNaoCitados_CitaPorIdOuNome prova o critério do Gate de paridade:
// o Nó do DAG vale como citado quando o id ou o nome aparece na Skill.
func TestNosNaoCitados_CitaPorIdOuNome(t *testing.T) {
	nos := []fechar.NoDAG{
		{ID: "e01_dissecao", Nome: "Dissecar o alvo"},
		{ID: "e02_lacunas", Nome: "Listar lacunas"},
	}
	casos := []struct {
		nome     string
		skill    string
		faltando []string
	}{
		{
			nome:     "sem citacao derruba o Gate",
			skill:    "Texto sem nenhuma referencia aos nos.",
			faltando: []string{"e01_dissecao (Dissecar o alvo)", "e02_lacunas (Listar lacunas)"},
		},
		{
			nome:  "citacao por id basta",
			skill: "Executa e01_dissecao e depois e02_lacunas.",
		},
		{
			nome:  "citacao por nome basta",
			skill: "Dissecar o alvo e Listar lacunas.",
		},
		{
			nome:     "no faltando e listado",
			skill:    "So falo de e01_dissecao.",
			faltando: []string{"e02_lacunas (Listar lacunas)"},
		},
		{
			nome:  "citacao ignora maiusculas e minusculas",
			skill: "DISSECAR O ALVO e LISTAR LACUNAS.",
		},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got := fechar.NosNaoCitados(c.skill, nos)
			if len(got) != len(c.faltando) {
				t.Fatalf("faltando = %v, want %v", got, c.faltando)
			}
			for i := range got {
				if got[i] != c.faltando[i] {
					t.Errorf("faltando[%d] = %q, want %q", i, got[i], c.faltando[i])
				}
			}
		})
	}
}

// TestFecharFalhaQuandoSkillNaoCitaNoDoDAG roda o Gate no fechamento: a Skill
// sem citacao do No derruba a etapa 4 com falha binaria.
func TestFecharFalhaQuandoSkillNaoCitaNoDoDAG(t *testing.T) {
	repo := t.TempDir()
	semearDocs(t, repo, "mudo")
	semearDAG(t, repo, "mudo")
	escrever(t, filepath.Join(repo, "skills", "legoloop-mudo", "SKILL.md"),
		"---\nname: legoloop-mudo\ndescription: Skill de teste\n---\n# legoloop-mudo\n"+
			strings.Repeat("Passo operacional com resultado verificavel.\n", 12))
	res := fecharSlug(t, repo, "mudo")
	exigirEtapaComErro(t, res, 4, "não cita os nós")
}

// TestFecharPassaQuandoSkillCitaTodosOsNos prova o lado verde do Gate.
func TestFecharPassaQuandoSkillCitaTodosOsNos(t *testing.T) {
	repo := t.TempDir()
	semearDocs(t, repo, "citado")
	semearDAG(t, repo, "citado")
	if res := fecharSlug(t, repo, "citado"); !res.Sucesso {
		t.Fatalf("Skill com citacao de todos os nos deve fechar: %+v", res.Etapas)
	}
}

// TestRelatorioParidadeSkillDAG_ListaFaltas prova o modo relatorio.
func TestRelatorioParidadeSkillDAG_ListaFaltas(t *testing.T) {
	repo := t.TempDir()
	semearDocs(t, repo, "inteiro")
	semearDAG(t, repo, "inteiro")
	semearDocs(t, repo, "falho")
	semearDAG(t, repo, "falho")
	escrever(t, filepath.Join(repo, "skills", "legoloop-falho", "SKILL.md"),
		"---\nname: legoloop-falho\ndescription: Skill de teste\n---\n# legoloop-falho\n"+
			strings.Repeat("Passo operacional com resultado verificavel.\n", 12))
	semearDAG(t, repo, "semskill")
	relatos := fechar.RelatorioParidadeSkillDAG(repo)
	porNome := map[string]fechar.ParidadeRelato{}
	for _, r := range relatos {
		porNome[r.Processo] = r
	}
	if r := porNome["legoloop-inteiro"]; !r.OK {
		t.Errorf("legoloop-inteiro deveria passar: %+v", r)
	}
	if r := porNome["legoloop-falho"]; r.OK || len(r.Faltando) != 2 {
		t.Errorf("legoloop-falho deveria falhar com 2 nos: %+v", r)
	}
	if r := porNome["legoloop-semskill"]; r.OK || r.Motivo != "SKILL.md ausente" {
		t.Errorf("legoloop-semskill deveria reportar SKILL.md ausente: %+v", r)
	}
}
