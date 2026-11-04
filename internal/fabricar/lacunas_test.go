package fabricar_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/matriz"
	"legoloop/internal/ste"
)

var precisaConstruirEVerificar = matriz.Needs{Natureza: "software", Etapas: []string{"construir", "verificar"}}

func semPerfil(ref string, score float64) matriz.Entrada {
	e := skill(ref, score, nil, nil)
	e.Perfil = nil
	return e
}

func decisaoDe(t *testing.T, m fabricar.Merged, etapa string) fabricar.Decisao {
	t.Helper()
	for _, d := range m.Decisoes {
		if d.Etapa == etapa {
			return d
		}
	}
	t.Fatalf("sem decisão para a etapa %q: %+v", etapa, m.Decisoes)
	return fabricar.Decisao{}
}

func TestLacunasEtapaSemPecaAptaViraLacuna(t *testing.T) {
	es := []matriz.Entrada{skill("tdd", 0.58, []string{"construir"}, sw)}
	m := fabricar.MergeByFit(precisaConstruirEVerificar, es, fabricar.MergeOptions{})
	if !slices.Equal(m.Lacunas, []string{"verificar"}) {
		t.Fatalf("Lacunas = %v", m.Lacunas)
	}
	d := decisaoDe(t, m, "verificar")
	if d.Peca != "" || d.Entrou || !strings.Contains(d.Motivo, "lacuna") {
		t.Errorf("a etapa sem peça é decisão com Peca vazia: %+v", d)
	}
	if got := nomes(m); !slices.Equal(got, []string{"tdd"}) {
		t.Errorf("a lacuna não vira passo: %v", got)
	}
	if !slices.ContainsFunc(m.Warnings, func(w string) bool { return strings.Contains(w, "Lacuna") && strings.Contains(w, "verificar") }) {
		t.Errorf("a lacuna deve ser destacada nos avisos: %v", m.Warnings)
	}
}

func TestLacunasPecaSemPerfilNaoEscolhidaEApareceComoDesconhecida(t *testing.T) {
	es := []matriz.Entrada{skill("tdd", 0.58, []string{"construir", "verificar"}, sw), semPerfil("sem-perfil", 0.95)}
	m := fabricar.MergeByFit(precisaConstruirEVerificar, es, fabricar.MergeOptions{})
	if got := nomes(m); slices.Contains(got, "sem-perfil") {
		t.Errorf("peça sem perfil nunca é escolhida: %v", got)
	}
	if len(m.Unknown) != 1 || m.Unknown[0].Name != "sem-perfil" || m.Unknown[0].Score != 0.95 {
		t.Errorf("Unknown = %+v", m.Unknown)
	}
	i := slices.IndexFunc(m.Decisoes, func(d fabricar.Decisao) bool { return d.Peca == "sem-perfil" })
	if i < 0 {
		t.Fatalf("a peça sem perfil deve constar nas decisões: %+v", m.Decisoes)
	}
	if d := m.Decisoes[i]; d.Entrou || d.Score != 0.95 || !strings.Contains(d.Motivo, "aptidão desconhecida") {
		t.Errorf("decisão da peça sem perfil: %+v", d)
	}
	if !slices.ContainsFunc(m.Warnings, func(w string) bool {
		return strings.Contains(w, "Aptidão desconhecida") && strings.Contains(w, "sem-perfil")
	}) {
		t.Errorf("o aviso deve citar a peça: %v", m.Warnings)
	}
}

func TestLacunasAvisosSeguemOPerfilSTE(t *testing.T) {
	es := []matriz.Entrada{skill("tdd", 0.58, []string{"construir"}, sw), semPerfil("sem-perfil", 0.95)}
	m := fabricar.MergeByFit(precisaConstruirEVerificar, es, fabricar.MergeOptions{})
	n := 0
	for _, w := range m.Warnings {
		if strings.Contains(w, "Lacuna") || strings.Contains(w, "Aptidão desconhecida") {
			n++
			if v := ste.Lint(w); len(v) != 0 {
				t.Errorf("aviso fora do perfil STE-PT: %q → %v", w, v)
			}
		}
	}
	if n != 2 {
		t.Errorf("esperava 2 avisos (lacuna e peça sem perfil), houve %d: %v", n, m.Warnings)
	}
}

func TestLacunasMuitasPecasSemPerfilViramResumo(t *testing.T) {
	es := []matriz.Entrada{skill("tdd", 0.58, []string{"construir", "verificar"}, sw)}
	for i := 0; i < 8; i++ {
		es = append(es, semPerfil(fmt.Sprintf("p%02d", i), 0.5))
	}
	m := fabricar.MergeByFit(precisaConstruirEVerificar, es, fabricar.MergeOptions{})
	if len(m.Unknown) != 8 {
		t.Fatalf("Unknown guarda todas: %d", len(m.Unknown))
	}
	var linhas []string
	for _, w := range m.Warnings {
		if strings.Contains(w, "Aptidão desconhecida") {
			linhas = append(linhas, w)
		}
	}
	if len(linhas) != 6 || !strings.Contains(linhas[5], "mais 3") {
		t.Errorf("5 peças citadas e um resumo com as 3 restantes: %v", linhas)
	}
}

func TestLacunasSemLacunaNemDesconhecidaNaoAvisa(t *testing.T) {
	es := []matriz.Entrada{skill("tdd", 0.58, []string{"construir", "verificar"}, sw)}
	m := fabricar.MergeByFit(precisaConstruirEVerificar, es, fabricar.MergeOptions{})
	if len(m.Lacunas) != 0 || len(m.Unknown) != 0 {
		t.Errorf("nada a declarar: %v %v", m.Lacunas, m.Unknown)
	}
	for _, w := range m.Warnings {
		if strings.Contains(w, "Lacuna") || strings.Contains(w, "Aptidão desconhecida") {
			t.Errorf("aviso indevido: %q", w)
		}
	}
	for _, e := range []string{"construir", "verificar"} {
		if d := decisaoDe(t, m, e); d.Peca != "tdd" || !d.Entrou || d.Aptidao != 1 || d.Motivo == "" || d.Score != 0.58 {
			t.Errorf("decisão da etapa %s: %+v", e, d)
		}
	}
}

func TestLacunasTodasAsEtapasSemPecaCaemNosPassosPadraoMasContinuamLacunas(t *testing.T) {
	m := fabricar.MergeByFit(precisaConstruirEVerificar, []matriz.Entrada{skill("pub", 0.9, []string{"publicar"}, sw)}, fabricar.MergeOptions{})
	if !slices.Equal(m.Lacunas, []string{"construir", "verificar"}) {
		t.Errorf("as lacunas seguem a ordem do pedido: %v", m.Lacunas)
	}
	if got := nomes(m); !slices.Equal(got, nomes(fabricar.Merged{Steps: fabricar.DefaultSteps()})) {
		t.Errorf("passos padrão: %v", got)
	}
}

func TestLacunasDesconhecidasIgnoraNaoSkillECopias(t *testing.T) {
	mcp := semPerfil("um-mcp", 0.9)
	mcp.Tipo = "mcp"
	a, b := semPerfil("a/x", 0.5), semPerfil("b/x", 0.5)
	a.Hash, b.Hash = "igual", "igual"
	m := fabricar.MergeByFit(precisaConstruirEVerificar, []matriz.Entrada{mcp, a, b}, fabricar.MergeOptions{})
	if len(m.Unknown) != 1 || m.Unknown[0].Path != "a/x" {
		t.Errorf("só skills, e cópias contam uma vez: %+v", m.Unknown)
	}
}

func TestLacunasSemEtapasPedidasNaoGeraLacuna(t *testing.T) {
	m := fabricar.MergeByFit(matriz.Needs{Natureza: "software"}, nil, fabricar.MergeOptions{})
	if len(m.Lacunas) != 0 {
		t.Errorf("sem etapas pedidas não há etapa para faltar: %v", m.Lacunas)
	}
}

func TestLacunasDesconhecidasVemDoMaiorScoreParaOMenorEDesempatamPorRef(t *testing.T) {
	es := []matriz.Entrada{semPerfil("b", 0.60), semPerfil("a", 0.60), semPerfil("z", 0.90), semPerfil("m", 0.30)}
	m := fabricar.MergeByFit(precisaConstruirEVerificar, es, fabricar.MergeOptions{})
	var got []string
	for _, p := range m.Unknown {
		got = append(got, p.Name)
	}
	if !slices.Equal(got, []string{"z", "a", "b", "m"}) {
		t.Errorf("Score desc e depois Ref asc: %v", got)
	}
}
