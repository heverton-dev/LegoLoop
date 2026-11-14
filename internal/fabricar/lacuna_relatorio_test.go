package fabricar_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"legoloop/internal/cli"
	"legoloop/internal/fabricar"
	"legoloop/internal/matriz"
	"legoloop/internal/reputation"
	"legoloop/internal/ste"
)

// matrizComLacuna tem alfa, apta a construir (tdd), e uma skill sem perfil de Score 0.95. Nenhuma peça é apta a
// verificar, então a etapa verificar vira lacuna (AC-E02.7).
func matrizComLacuna() []matriz.Entrada {
	alfa := goodPieces()[0]
	ref := filepath.Join(os.TempDir(), "s", "sem-perfil", "SKILL.md")
	sem := matriz.Entrada{Ref: ref, Nome: "sem-perfil", Tipo: "skill", Hash: "h-sem",
		Rating: reputation.PieceRating{Ref: ref, Name: "sem-perfil", TotalScore: 0.95, Tier: reputation.TierOuro}}
	return []matriz.Entrada{alfa, sem}
}

// TestLacunaRelatorio prova o AC-E02.7 nos dois modos: no --dry-run a lacuna não para antes do relatório (o
// Report e o cartão saem com a etapa e a orientação STE-PT, e o erro continua sendo ErrLacuna); sem --dry-run
// a lacuna para sem gravar nada e sem relatório. A peça sem perfil aparece como aptidão desconhecida.
func TestLacunaRelatorio(t *testing.T) {
	t.Run("dry-run sai com o relatório e o cartão", testLacunaDryRun)
	t.Run("sem dry-run para sem gravar", testLacunaSemDryRun)
	t.Run("texto destaca a lacuna", testLacunaTexto)
}

func rodarLacuna(t *testing.T, dry bool) (*fabricar.Result, error, *rig) {
	t.Helper()
	r := newRig(t, matrizComLacuna())
	r.req.DryRun = dry
	res, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
	return res, err, r
}

func nadaGravado(t *testing.T, r *rig) {
	t.Helper()
	for _, root := range []string{r.deps.Roots.Vault, r.deps.Roots.Docs} {
		if _, err := os.Stat(root); err == nil {
			t.Errorf("a lacuna não pode gravar nada: %s existe", root)
		}
	}
}

func testLacunaDryRun(t *testing.T) {
	res, err, r := rodarLacuna(t, true)
	if err != nil {
		t.Fatalf("cascata de 3 fontes não deve retornar erro no dry-run com lacuna: %v", err)
	}
	if res == nil || res.Report == nil {
		t.Fatal("no dry-run Result e Report são obrigatórios")
	}
	nadaGravado(t, r)
	if !slices.ContainsFunc(res.Warnings, func(w string) bool { return strings.Contains(w, "verificar") }) {
		t.Errorf("o aviso da cascata deve citar verificar: %v", res.Warnings)
	}
	rep := res.Report
	assertDecisoesDaLacuna(t, rep.Fusao.Decisoes)
	if len(rep.Gravacao.Arquivos) == 0 {
		t.Errorf("com cascata de 3 fontes o dry-run planeja os arquivos normalmente: %+v", rep.Gravacao)
	}
	if rep.Triade.ACs != 2 || len(rep.Triade.Gates) != 2 || !strings.Contains(rep.Triade.Gates[1], "02_verificar") {
		t.Errorf("a Tríade mostra as duas etapas: %+v", rep.Triade)
	}
	avisos := strings.Join(res.Warnings, "\n")
	if !strings.Contains(avisos, "verificar") || !strings.Contains(avisos, "sugestão externa") {
		t.Errorf("os avisos mostram a cascata externa para verificar:\n%s", avisos)
	}
	if res.CanonicalName != "legoloop-extrator-de-audio" || res.Artifacts.SkillPath == "" {
		t.Errorf("o resultado nomeia o processo e aponta os artefatos planejados: %q %+v", res.CanonicalName, res.Artifacts)
	}
}

func assertDecisoesDaLacuna(t *testing.T, ds []fabricar.Decisao) {
	t.Helper()
	var lacuna, sem, alfa bool
	for _, d := range ds {
		switch {
		case d.Etapa == "verificar" && d.Peca == "" && !d.Entrou && strings.HasPrefix(d.Motivo, "lacuna"):
			lacuna = true
		case d.Peca == "sem-perfil" && !d.Entrou && strings.Contains(d.Motivo, "aptidão desconhecida"):
			sem = true
		case d.Peca == "alfa" && d.Entrou && d.Etapa == "construir":
			alfa = true
		}
	}
	if !lacuna || !sem || !alfa {
		t.Errorf("decisões: lacuna em verificar=%v, sem-perfil desconhecida=%v, alfa em construir=%v\n%+v", lacuna, sem, alfa, ds)
	}
}

func assertCartaoDaLacuna(t *testing.T, c ste.Cartao) {
	t.Helper()
	if !strings.Contains(c.Acao, "lacuna") {
		t.Errorf("a ação do cartão diz que há lacuna: %q", c.Acao)
	}
	if !strings.Contains(c.Aviso, "A etapa verificar não tem peça apta") || !strings.Contains(c.Aviso, "legoloop matriz perfil definir") {
		t.Errorf("o aviso do cartão destaca a etapa e a orientação: %q", c.Aviso)
	}
	if !slices.ContainsFunc(c.Passos, func(p string) bool { return strings.Contains(p, "Lacuna: a etapa verificar") }) {
		t.Errorf("os passos destacam a etapa da lacuna: %v", c.Passos)
	}
	if len(c.Efeitos) != 0 {
		t.Errorf("com lacuna o cartão não promete gravar nada: %v", c.Efeitos)
	}
	if v := ste.LintCartao(c); len(v) != 0 {
		t.Errorf("o cartão da lacuna segue o STE-PT: %+v", v)
	}
}

func testLacunaSemDryRun(t *testing.T) {
	res, err, _ := rodarLacuna(t, false)
	if err != nil {
		t.Fatalf("cascata de 3 fontes não deve bloquear fabricação sem dry-run: %v", err)
	}
	if res == nil {
		t.Fatal("resultado não deve ser nil")
	}
	if !slices.ContainsFunc(res.Warnings, func(w string) bool { return strings.Contains(w, "verificar") }) {
		t.Errorf("o aviso da cascata deve citar verificar: %v", res.Warnings)
	}
}

func testLacunaTexto(t *testing.T) {
	res, _, _ := rodarLacuna(t, true)
	if res == nil {
		t.Fatal("sem resultado no dry-run")
	}
	var b bytes.Buffer
	cli.PrintReportTexto(&b, res.Report)
	out := b.String()
	if !strings.Contains(out, "verificar  LACUNA |") {
		t.Errorf("o resumo destaca a etapa da lacuna:\n%s", out)
	}
	if !strings.Contains(out, "aptidão desconhecida: 1") {
		t.Errorf("o resumo conta a peça sem perfil como aptidão desconhecida:\n%s", out)
	}
}
