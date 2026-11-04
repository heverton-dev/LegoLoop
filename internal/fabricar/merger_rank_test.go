package fabricar_test

import (
	"slices"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/matriz"
	"legoloop/internal/reputation"
)

func skill(ref string, score float64, etapas, naturezas []string) matriz.Entrada {
	return matriz.Entrada{Ref: ref, Nome: ref, Tipo: "skill", Hash: "h-" + ref,
		Perfil: &matriz.PerfilDeUso{Etapas: etapas, Naturezas: naturezas, Origem: matriz.OrigemDeclarado},
		Rating: reputation.PieceRating{Ref: ref, TotalScore: score, Tier: reputation.TierPrata}}
}

func nomes(m fabricar.Merged) []string {
	var out []string
	for _, s := range m.Steps {
		out = append(out, s.Piece)
	}
	return out
}

var sw = []string{"software"}

func TestRankMergeUmaPecaPorEtapa(t *testing.T) {
	es := []matriz.Entrada{
		skill("tdd", 0.58, []string{"construir"}, sw),
		skill("revisor", 0.90, []string{"verificar"}, sw),
		skill("publicador", 0.99, []string{"publicar"}, sw),
	}
	m := fabricar.MergeByFit(matriz.Needs{Natureza: "software", Etapas: []string{"construir", "verificar"}}, es, fabricar.MergeOptions{})
	if got := nomes(m); !slices.Equal(got, []string{"tdd", "revisor"}) {
		t.Errorf("um passo por etapa, pela aptidão e não pelo Score: %v", got)
	}
	if len(m.Used) != 2 || m.Used[0].Name != "tdd" || m.Used[0].Path != "tdd" || m.Used[0].Score != 0.58 || m.Used[0].Tier != reputation.TierPrata {
		t.Errorf("Used deve descrever as peças escolhidas: %+v", m.Used)
	}
}

func TestRankMergePecaQueCobreDuasEtapasViraUmPasso(t *testing.T) {
	es := []matriz.Entrada{skill("tdd", 0.58, []string{"construir", "verificar"}, sw)}
	m := fabricar.MergeByFit(matriz.Needs{Natureza: "software", Etapas: []string{"construir", "verificar"}}, es, fabricar.MergeOptions{})
	if got := nomes(m); !slices.Equal(got, []string{"tdd"}) {
		t.Errorf("a mesma peça não repete o passo: %v", got)
	}
}

func TestRankMergeIgnoraOQueNaoESkillEJuntaCopias(t *testing.T) {
	mcp := skill("um-mcp", 0.99, []string{"construir"}, sw)
	mcp.Tipo = "mcp"
	copiaB := skill("b/tdd", 0.90, []string{"construir"}, sw) // melhor Score, mas Ref maior
	copiaA := skill("a/tdd", 0.50, []string{"construir"}, sw)
	copiaA.Hash, copiaB.Hash = "mesmo", "mesmo" // cópias idênticas, em caminhos diferentes
	m := fabricar.MergeByFit(matriz.Needs{Natureza: "software", Etapas: []string{"construir"}}, []matriz.Entrada{mcp, copiaB, copiaA}, fabricar.MergeOptions{})
	if len(m.Used) != 1 || m.Used[0].Path != "a/tdd" {
		t.Errorf("MCP nunca é herdado e a cópia de menor Ref representa o grupo: %+v", m.Used)
	}
}

func TestRankMergeAvisaExclusaoPorUsoUmaVez(t *testing.T) {
	site := skill("aidd-site-dominator", 0.90, []string{"publicar"}, []string{"midia"})
	site.Perfil.QuandoNaoUsar = []string{"software"}
	soEtapa := skill("so-publica", 0.95, []string{"publicar"}, sw) // sai só por etapa: não é caso de aviso
	es := []matriz.Entrada{skill("tdd", 0.58, []string{"construir", "verificar"}, sw), site, soEtapa}
	m := fabricar.MergeByFit(matriz.Needs{Natureza: "software", Etapas: []string{"construir", "verificar"}}, es, fabricar.MergeOptions{})
	n := 0
	for _, w := range m.Warnings {
		if strings.Contains(w, "aidd-site-dominator") {
			n++
			if !strings.Contains(w, "quando_nao_usar") {
				t.Errorf("o aviso deve trazer o motivo: %q", w)
			}
		}
	}
	if n != 1 {
		t.Errorf("a exclusão deve ser avisada uma vez, não %d: %v", n, m.Warnings)
	}
	i := slices.IndexFunc(m.Decisoes, func(d fabricar.Decisao) bool { return d.Peca == "aidd-site-dominator" })
	if i < 0 || m.Decisoes[i].Entrou || !strings.Contains(m.Decisoes[i].Motivo, "quando_nao_usar") || m.Decisoes[i].Score != 0.90 {
		t.Errorf("a exclusão por uso deve constar nas decisões, com motivo e Score: %+v", m.Decisoes)
	}
	if slices.ContainsFunc(m.Warnings, func(w string) bool { return strings.Contains(w, "so-publica") }) {
		t.Errorf("sair só por etapa não é exclusão a avisar: %v", m.Warnings)
	}
	if len(m.Excluded) != 1 || m.Excluded[0].Name != "aidd-site-dominator" {
		t.Errorf("Excluded deve listar a peça: %+v", m.Excluded)
	}
}

func TestRankMergeSemPecaAptaUsaOsPassosPadrao(t *testing.T) {
	es := []matriz.Entrada{skill("publicador", 0.99, []string{"publicar"}, sw)}
	m := fabricar.MergeByFit(matriz.Needs{Natureza: "software", Etapas: []string{"construir"}}, es, fabricar.MergeOptions{})
	if got := nomes(m); !slices.Equal(got, nomes(fabricar.Merged{Steps: fabricar.DefaultSteps()})) {
		t.Errorf("sem peça apta valem os passos padrão: %v", got)
	}
	if len(m.Warnings) == 0 {
		t.Error("o uso dos passos padrão deve ser avisado")
	}
}

func TestRankMergeRespeitaOLimiteDePassos(t *testing.T) {
	var es []matriz.Entrada
	for _, e := range []string{"clarificar", "planejar", "construir"} {
		es = append(es, skill("s-"+e, 0.5, []string{e}, sw))
	}
	m := fabricar.MergeByFit(matriz.Needs{Natureza: "software", Etapas: []string{"clarificar", "planejar", "construir"}}, es, fabricar.MergeOptions{MaxSteps: 2})
	if len(m.Steps) != 2 {
		t.Errorf("MaxSteps=2: %v", nomes(m))
	}
	if !slices.ContainsFunc(m.Warnings, func(w string) bool { return strings.Contains(w, "limite de 2 passos") }) {
		t.Errorf("o corte deve ser avisado: %v", m.Warnings)
	}
}

func TestRankMergeSemEtapasPedidasUsaARankingGeral(t *testing.T) {
	es := []matriz.Entrada{skill("a", 0.9, []string{"construir"}, sw), skill("b", 0.5, []string{"verificar"}, sw)}
	m := fabricar.MergeByFit(matriz.Needs{Natureza: "software"}, es, fabricar.MergeOptions{})
	if got := nomes(m); !slices.Equal(got, []string{"a", "b"}) {
		t.Errorf("sem etapas pedidas valem as melhores peças, na ordem do ranking: %v", got)
	}
}
