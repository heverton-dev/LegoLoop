package templates_test

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/dag"
	"legoloop/internal/gate"
	"legoloop/internal/roteamento"
	studio "legoloop/internal/studiodoc"
)

const (
	arquivoUnit = "TEMPLATE_PIECE_UNIT.json"
	apoioUnit   = "TEMPLATE_PIECE_UNIT.md"
	idUnit      = "exemplo-peca-unit-no"
)

func carregarUnit(t *testing.T) (*dag.DAG, dag.Node) {
	t.Helper()
	d, err := dag.Parse(read(t, arquivoUnit))
	if err != nil {
		t.Fatalf("%s não passa em dag.Parse: %v", arquivoUnit, err)
	}
	if len(d.Nodes) != 1 || d.Nodes[0].Type != dag.NodeType("unit_no") {
		t.Fatalf("o template deve ter um único nó unit_no: %+v", d.Nodes)
	}
	return d, d.Nodes[0]
}

func TestUnitTemplateEhDAGValidoComUmNoUnit(t *testing.T) {
	d, n := carregarUnit(t)
	if d.SchemaVersion != 2 || d.EffectiveTrigger().Type != dag.TriggerManual {
		t.Errorf("esperava schema_version 2 e trigger manual: %+v", d.Trigger)
	}
	if strings.Contains(string(read(t, arquivoUnit)), "{{") {
		t.Error("o template deve ser exemplo válido, sem placeholders")
	}
	if _, err := dag.Parse(read(t, arquivoUnit)); err != nil {
		t.Error(err)
	}
	if n.ID == "" || n.Name == "" {
		t.Errorf("nó sem id ou nome: %+v", n)
	}
}

func TestUnitTemplateContratoDoExecutor(t *testing.T) {
	_, n := carregarUnit(t)
	if n.RetryPolicy == nil || n.RetryPolicy.MaxRetries < 0 || n.RetryPolicy.MaxRetries > dag.MaxUnitRetries {
		t.Errorf("retry_policy.max_retries deve estar entre 0 e %d: %+v", dag.MaxUnitRetries, n.RetryPolicy)
	}
	if n.Isolation["write_boundary"] != "strict" || n.Isolation["subagents_allowed"] != false {
		t.Errorf("isolamento deve ser strict e sem subagentes: %+v", n.Isolation)
	}
	if _, err := gate.Compile(n.Output.Gate); err != nil {
		t.Errorf("o gate não compila: %v", err)
	}
	if !strings.Contains(n.Output.Gate, "exit_code") {
		t.Errorf("o gate do unit_no avalia o exit_code do command: %q", n.Output.Gate)
	}
	if len(n.Command) == 0 {
		t.Fatal("o unit_no do exemplo declara o command de verificação")
	}
	if err := dag.ValidateCommand(n.Command); err != nil {
		t.Errorf("command inválido: %v", err)
	}
	if n.TestSeam["file"] == "" || n.TestSeam["signature"] == "" {
		t.Errorf("test_seam exige file e signature: %+v", n.TestSeam)
	}
}

// Espelha allowedWrites: context_out mais test_seam.file são as únicas escritas permitidas, relativas ao projeto.
func TestUnitTemplateEscritasPermitidasSaoRelativas(t *testing.T) {
	_, n := carregarUnit(t)
	permitidas := append(append([]string{}, n.ContextOut...), n.TestSeam["file"])
	if len(n.ContextOut) == 0 {
		t.Error("context_out não pode ser vazio: o unit_no precisa de ao menos uma escrita permitida")
	}
	for _, p := range permitidas {
		if p == "" || path.IsAbs(p) || filepath.IsAbs(p) || strings.Contains(p, "..") || strings.Contains(p, ":") {
			t.Errorf("escrita permitida inválida %q", p)
		}
		if p == n.Output.Artifact {
			t.Errorf("o artefato do relatório %q não pode estar entre as escritas permitidas", p)
		}
	}
}

// Universalidade: o roteamento é preferência editável, e a classe degrada para o próximo candidato homologado.
func TestUnitTemplateRoteamentoUniversal(t *testing.T) {
	_, n := carregarUnit(t)
	if n.Routing.Mode != dag.ModeSuggested {
		t.Errorf("o roteamento do template é suggested, não %q", n.Routing.Mode)
	}
	if roteamento.PoliticaDeFalha(n.Nature) != roteamento.PoliticaDegrada {
		t.Errorf("a natureza %q do template não degrada; o nó quebraria sem o harness declarado", n.Nature)
	}
	if n.Routing.Fallback == "" {
		t.Error("o template declara um fallback")
	}
	texto := string(read(t, arquivoUnit))
	for _, shell := range []string{`"powershell"`, `"pwsh"`, `"cmd"`, `"bash"`, `"sh"`} {
		if strings.Contains(strings.ToLower(texto), shell) {
			t.Errorf("o template fixa o shell %s", shell)
		}
	}
}

func TestUnitTemplateDeclaraPecaAgnostica(t *testing.T) {
	texto := string(read(t, arquivoUnit))
	for _, c := range []string{"execution_modes", "os_support", "harness_support", `"load": "lazy"`} {
		if !strings.Contains(texto, c) {
			t.Errorf("%s sem %s", arquivoUnit, c)
		}
	}
}

func TestUnitTemplateCarregaNoStudio(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, idUnit+caminhos.SufixoAntigo), read(t, arquivoUnit), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := studio.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if doc, err := st.Get(idUnit); err != nil || len(doc.Nodes) != 1 {
		t.Fatalf("Store.Get: %v %+v", err, doc)
	}
}

func TestUnitTemplateTemDocumentoDeApoioERegistroNoReadme(t *testing.T) {
	apoio := string(read(t, apoioUnit))
	for _, termo := range []string{"consentimento", "--allow-unit-no", "retry_policy", "write_boundary", "context_out", "test_seam", "worktree", "degrada", "legoloop rota"} {
		if !strings.Contains(apoio, termo) {
			t.Errorf("%s não explica %q", apoioUnit, termo)
		}
	}
	if !strings.Contains(string(read(t, "README.md")), arquivoUnit) {
		t.Errorf("o README deve citar %s", arquivoUnit)
	}
}
