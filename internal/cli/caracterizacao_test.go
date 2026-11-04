package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Testes de caracterização: fixam o comportamento atual de upgrade e import antes de dividir as funções
// longas (gate funclen do TICK-02.4). Se um deles falhar depois da refatoração, a refatoração mudou algo.

func skillAnemica(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "SKILL.md")
	if err := os.WriteFile(p, []byte("# Anemic Skill\nFaça um bom trabalho."), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCaracterizaUpgradeJSON(t *testing.T) {
	d, out, errb := deps(t)
	p := skillAnemica(t)
	if code := run(t, d, "upgrade", p, "--json"); code != 0 {
		t.Fatalf("exit %d: %s", code, errb)
	}
	var got map[string]any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("saída não é JSON: %v\n%s", err, out)
	}
	if got["target"] != p || got["is_gold"] != false || got["actionable"] != true {
		t.Errorf("campos inesperados: %v", got)
	}
	if gaps, _ := got["gaps"].([]any); len(gaps) == 0 {
		t.Error("a skill anêmica deve ter gaps")
	}
	if _, ok := got["rating"].(map[string]any); !ok {
		t.Error("falta o campo rating")
	}
}

func TestCaracterizaUpgradeErros(t *testing.T) {
	d, _, errb := deps(t)
	if code := run(t, d, "upgrade"); code != 2 || !strings.Contains(errb.String(), "informe o caminho") {
		t.Errorf("sem alvo: exit %d\n%s", code, errb)
	}
	d2, _, errb2 := deps(t)
	if code := run(t, d2, "upgrade", filepath.Join(t.TempDir(), "nao-existe.md")); code != 1 || !strings.Contains(errb2.String(), "erro ao ler arquivo alvo") {
		t.Errorf("arquivo ausente: exit %d\n%s", code, errb2)
	}
}

func TestCaracterizaUpgradeSemFlagsFazDiagnostico(t *testing.T) {
	d, out, errb := deps(t)
	if code := run(t, d, "upgrade", skillAnemica(t)); code != 0 {
		t.Fatalf("exit %d: %s", code, errb)
	}
	for _, want := range []string{"DIAGNÓSTICO DE ELEVAÇÃO", "Score Atual", "GAPS PARA ATINGIR O PADRÃO OURO", "Recomendação"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("saída sem %q:\n%s", want, out)
		}
	}
}

func TestCaracterizaImportDestino(t *testing.T) {
	home := t.TempDir()
	canon := filepath.Join(home, ".aidd", "canonical")
	if got := importDest("", home); got != canon {
		t.Errorf("sem ~/.gemini/config o destino é o cofre canônico: %q", got)
	}
	gemini := filepath.Join(home, ".gemini", "config")
	if err := os.MkdirAll(gemini, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := importDest("", home); got != gemini {
		t.Errorf("com ~/.gemini/config o destino é ele: %q", got)
	}
	if got := importDest("/x/y", home); got != "/x/y" {
		t.Errorf("--target vence: %q", got)
	}
	if got := importDest("", ""); got == "" {
		// sem home informado cai em os.UserHomeDir(); só falha se o ambiente não tiver HOME.
		t.Log("ambiente sem HOME: destino vazio")
	}
}
