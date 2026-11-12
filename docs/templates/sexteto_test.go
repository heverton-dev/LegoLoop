package templates_test

import (
	"encoding/json"
	"strings"
	"testing"

	"legoloop/docs/templates"
)

func TestInterviewEHandoffEmbutidos(t *testing.T) {
	casos := []struct {
		arquivo string
		prefixo string
		obtido  string
	}{
		{"TEMPLATE_INTERVIEW.md", "---\nid: INTERVIEW-", templates.Interview()},
		{"TEMPLATE_HANDOFF.md", "---\nid: HANDOFF-", templates.Handoff()},
	}
	for _, c := range casos {
		if !strings.HasPrefix(c.obtido, c.prefixo) {
			t.Errorf("%s embutido deve começar com %q", c.arquivo, c.prefixo)
		}
		if c.obtido != string(read(t, c.arquivo)) {
			t.Errorf("%s embutido difere do arquivo em disco", c.arquivo)
		}
	}
}

func TestTemplatesDoSexteto(t *testing.T) {
	markdowns := []struct {
		arquivo  string
		prefixo  string
		secoes   []string
		campos   []string
		minTotal int
	}{
		{"TEMPLATE_INTERVIEW.md", "---\nid: INTERVIEW-", []string{"## 1. Respostas", "## 2. Confiança por rodada", "## 3. Fatos medidos durante a entrevista", "## 4. Pendências para a SPEC"}, []string{"\nschema_version: 2\n", "\nupstream: [", "\ncovers: ["}, 4},
		{"TEMPLATE_INTENT.md", "---\nid: INTENT-", []string{"## 1. Contexto e Motivação", "## 2. Escopo Inegociável", "## 3. Mapeamento Preliminar de Peças", "## 6. Handoff de Governança"}, []string{"\nschema_version: 2\n", "\nupstream: [", "\ncovers: ["}, 6},
		{"TEMPLATE_SPEC.md", "---\nid: SPEC-", []string{"## 1. Premissas Técnicas", "## 4. Contratos de Dados", "## 5. Critérios Binários de Aceite", "## 6. Comandos Operacionais"}, []string{"\nschema_version: 2\n", "\nupstream: [", "\ncovers: [", "| R | Requisito | AC |"}, 6},
		{"TEMPLATE_PLAN.md", "---\nid: PLAN-", []string{"## 1. Estratégia de Execução", "## 2. Mapa Geral de Épicos", "## 3. Decomposição Estruturada de Tickets", "## 4. Homologação Final"}, []string{"\nschema_version: 2\n", "\nupstream: [", "\ncovers: ["}, 4},
		{"TEMPLATE_REPORT.md", "", []string{"## 1. Resumo Executivo", "## 2. Telemetria Estrita", "## 3. Entregas por Épico", "## 4. Matriz de Qualidade", "## 5. Inventário de Artefatos", "## 6. Conclusão, Recomendações"}, nil, 6},
	}
	for _, m := range markdowns {
		texto := string(read(t, m.arquivo))
		if m.prefixo != "" && !strings.HasPrefix(texto, m.prefixo) {
			t.Errorf("%s deve começar com %q", m.arquivo, m.prefixo)
		}
		for _, s := range m.secoes {
			if !strings.Contains(texto, "\n"+s) {
				t.Errorf("%s sem a seção %q", m.arquivo, s)
			}
		}
		for _, c := range m.campos {
			if !strings.Contains(texto, c) {
				t.Errorf("%s sem o campo %q", m.arquivo, c)
			}
		}
		if n := strings.Count(texto, "\n## "); n < m.minTotal {
			t.Errorf("%s: esperava ao menos %d seções, obteve %d", m.arquivo, m.minTotal, n)
		}
	}
	for _, arquivo := range []string{"TEMPLATE_ROUTE.json", "TEMPLATE_DAG.json"} {
		var doc map[string]any
		if err := json.Unmarshal(read(t, arquivo), &doc); err != nil {
			t.Fatalf("%s: %v", arquivo, err)
		}
		if nos, _ := doc["nodes"].([]any); len(nos) == 0 {
			t.Errorf("%s sem nós", arquivo)
		}
	}
	var rota map[string]any
	if err := json.Unmarshal(read(t, "TEMPLATE_ROUTE.json"), &rota); err != nil {
		t.Fatal(err)
	}
	for _, c := range []string{"date", "execution_strategy", "upstream_signatures", "nodes"} {
		if _, ok := rota[c]; !ok {
			t.Errorf("TEMPLATE_ROUTE.json sem o campo %q", c)
		}
	}
}

func TestReadmeCitaSexteto(t *testing.T) {
	texto := string(read(t, "README.md"))
	for _, s := range []string{"TEMPLATE_INTERVIEW.md", "TEMPLATE_HANDOFF.md", "TEMPLATE_REPORT.md", "os_support", "execution_modes", "harness_support"} {
		if !strings.Contains(texto, s) {
			t.Errorf("README.md sem a citação %q", s)
		}
	}
}

func TestHandoffTemplate(t *testing.T) {
	texto := string(read(t, "TEMPLATE_HANDOFF.md"))
	for _, s := range []string{"---\nid: HANDOFF-", "\nschema_version: 2\n", "\nupstream: [REPORT-"} {
		if !strings.Contains(texto, s) {
			t.Errorf("TEMPLATE_HANDOFF.md sem %q", s)
		}
	}
}

func TestPecasAgnosticas(t *testing.T) {
	campos := []string{"execution_modes", "os_support", "harness_support"}
	pecas := []string{"TEMPLATE_PIECE_SKILL.md", "TEMPLATE_PIECE_MCP.json", "TEMPLATE_PIECE_HOOK.json", "TEMPLATE_PIECE_GATE.json", "TEMPLATE_PIECE_PROCESSO.json", "TEMPLATE_PIECE_SCRIPT.sh", "TEMPLATE_PIECE_SCRIPT.ps1"}
	for _, p := range pecas {
		texto := string(read(t, p))
		for _, c := range campos {
			if !strings.Contains(texto, c) {
				t.Errorf("%s sem o campo %s", p, c)
			}
		}
	}
	if strings.Contains(string(read(t, "TEMPLATE_PIECE_SKILL.md")), "suggested_model") {
		t.Error("a peça skill declara classe, não modelo")
	}
	if strings.Contains(string(read(t, "TEMPLATE_PIECE_MCP.json")), "autoApprove") {
		t.Error("a peça MCP não fixa autoApprove")
	}
	sh := string(read(t, "TEMPLATE_PIECE_SCRIPT.sh"))
	ps := string(read(t, "TEMPLATE_PIECE_SCRIPT.ps1"))
	for _, s := range []string{"--target-file", "--project-root", "--dry-run", "--json"} {
		if !strings.Contains(sh, s) {
			t.Errorf("script sh sem a opção %s", s)
		}
	}
	for _, s := range []string{"TargetFile", "ProjectRoot", "DryRun", "Json"} {
		if !strings.Contains(ps, s) {
			t.Errorf("script ps1 sem o parâmetro %s", s)
		}
	}
}
