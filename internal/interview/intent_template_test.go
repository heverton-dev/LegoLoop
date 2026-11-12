package interview_test

import (
	"os"
	"strings"
	"testing"

	"legoloop/docs/templates"
	"legoloop/internal/interview"
)

const mirrorResp = `HIPÓTESE: Processo Go automatizado.
CONFIANÇA: 98%
ESPELHO CANÔNICO:
- Resultado Esperado: Processo automatizado de validação Go
- Quem se Beneficia: Desenvolvedores do backend
- Por que Agora: Reduzir regressões em produção
- Critério de Sucesso: Testes rodam em menos de 2 minutos
- Restrições Inegociáveis: Sem subagentes, Segredos só por variável de ambiente
- Fora de Escopo: Deploy em produção, Notificações por e-mail`

func materialize(t *testing.T, tweak func(*interview.Session)) string {
	t.Helper()
	engine := interview.NewEngine(&mockLLM{responses: []string{mirrorResp}})
	sess, _, err := engine.Start(t.Context(), "Processo de testes")
	if err != nil {
		t.Fatal(err)
	}
	if tweak != nil {
		tweak(sess)
	}
	path, err := engine.MaterializeIntent(sess.ID, t.TempDir(), "2026-11-02_INTENT_PROCESSO_TESTES.md")
	if err != nil {
		t.Fatalf("MaterializeIntent: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestIntentTemplate(t *testing.T) {
	const secao = "## 3. Mapeamento Preliminar de Peças"
	campos := []string{"id: INTENT-", "status: ", "date: ", "schema_version: 2", "upstream: [", "covers: ["}
	docs := map[string]string{"template": templates.Intent(), "gerado": materialize(t, nil)}
	for nome, doc := range docs {
		if !strings.HasPrefix(doc, "---\n") {
			t.Errorf("%s: o documento deve abrir com front matter", nome)
		}
		fim := strings.Index(doc[4:], "\n---\n")
		if fim < 0 {
			t.Fatalf("%s: front matter não fechado", nome)
		}
		for _, c := range campos {
			if !strings.Contains(doc[:fim+4], "\n"+c) {
				t.Errorf("%s: front matter sem o campo %q", nome, c)
			}
		}
		if !strings.Contains(doc, secao) {
			t.Errorf("%s: falta a seção %q", nome, secao)
		}
	}
}

func TestMaterializeIntent_FollowsTemplateStructure(t *testing.T) {
	doc := materialize(t, nil)
	pos := 0
	for _, h := range templates.Headings(templates.Intent()) {
		i := strings.Index(doc[pos:], "## "+h)
		if i < 0 {
			t.Fatalf("seção do template ausente ou fora de ordem: %q\n%s", h, doc)
		}
		pos += i + len(h)
	}
	if strings.Contains(doc, "{{") {
		t.Errorf("a INTENT gerada deixou placeholders sem preencher:\n%s", doc)
	}
	for _, want := range []string{
		"# 🎯 INTENT: Processo de testes",
		"**Template:** TEMPLATE_INTENT.md v" + templates.Version,
		"**Confiança da Entrevista:** 98%",
		"STOP-GATE 01",
		"docs/specs/2026-11-02_SPEC_PROCESSO_TESTES.md",
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("a INTENT gerada não contém %q", want)
		}
	}
}

func TestMaterializeIntent_CarriesMirrorIntoTemplateSections(t *testing.T) {
	doc := materialize(t, nil)
	for _, want := range []string{
		"Processo automatizado de validação Go", "Desenvolvedores do backend", "Reduzir regressões em produção",
		"Testes rodam em menos de 2 minutos", "Sem subagentes", "Segredos só por variável de ambiente",
		"Deploy em produção", "Notificações por e-mail",
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("o espelho socrático não chegou ao documento: falta %q", want)
		}
	}
}

func TestMaterializeIntent_NeverSelfApproves(t *testing.T) {
	doc := materialize(t, nil)
	if !strings.Contains(doc, "**Status:** UNDER_REVIEW") || strings.Contains(doc, "APPROVED") {
		t.Errorf("a INTENT deve nascer UNDER_REVIEW: a aprovação é humana\n%s", doc)
	}
	low := materialize(t, func(s *interview.Session) { s.Confidence = 80 })
	if !strings.Contains(low, "**Status:** DRAFT") {
		t.Errorf("confiança abaixo de 95%% deve gerar DRAFT\n%s", low)
	}
}

func TestMaterializeIntent_EmptyMirrorListsStayHonest(t *testing.T) {
	doc := materialize(t, func(s *interview.Session) { s.Mirror.Constraints, s.Mirror.OutOfScope = nil, nil })
	if strings.Contains(doc, "{{") || !strings.Contains(doc, "Nenhuma restrição registrada na entrevista") {
		t.Errorf("listas vazias devem ser declaradas explicitamente, sem placeholder\n%s", doc)
	}
}
