package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/processo"
	"github.com/heverton-dev/LegoLoop/internal/studiodoc"
)

// dirDeProcesso cria processos/<slug> num cofre temporário e devolve a pasta.
func dirDeProcesso(t *testing.T, slug string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), caminhos.DirProcessos, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// TestValidateDerrubaEmDriftStudio prova o gate de paridade do Studio: o documento do
// Studio é derivado do DAG canônico e nós diferentes entre os dois derrubam o validate.
func TestValidateDerrubaEmDriftStudio(t *testing.T) {
	dir := dirDeProcesso(t, "legoloop-drift")
	dagPath := caminhos.CaminhoDAG(dir)
	if err := os.WriteFile(dagPath, []byte(processoScript), 0o644); err != nil {
		t.Fatal(err)
	}
	derivado, err := processo.DerivarStudio([]byte(processoScript), time.Now())
	if err != nil {
		t.Fatalf("DerivarStudio: %v", err)
	}
	var s studiodoc.Processo
	if err := json.Unmarshal(derivado, &s); err != nil {
		t.Fatal(err)
	}
	s.Nodes[0].Output.Artifact = "INVASOR.txt"
	drift, _ := json.Marshal(s)
	studioPath := filepath.Join(dir, "legoloop-drift"+caminhos.SufixoStudio)
	if err := os.WriteFile(studioPath, drift, 0o644); err != nil {
		t.Fatal(err)
	}
	d, _, errb := deps(t)
	if code := run(t, d, "validate", dagPath); code == 0 || !strings.Contains(errb.String(), "diverge") {
		t.Fatalf("drift não derrubou o gate: code=%d stderr=%s", code, errb.String())
	}
}

// TestValidateRegeneraStudioDerivado prova que validate --regenerar-studio regrava o
// documento do Studio a partir do DAG, e que o regenerado passa no gate de paridade.
func TestValidateRegeneraStudioDerivado(t *testing.T) {
	dir := dirDeProcesso(t, "legoloop-regen")
	dagPath := caminhos.CaminhoDAG(dir)
	if err := os.WriteFile(dagPath, []byte(processoScript), 0o644); err != nil {
		t.Fatal(err)
	}
	studioPath := filepath.Join(dir, "legoloop-regen"+caminhos.SufixoStudio)

	d, _, _ := deps(t)
	if code := run(t, d, "validate", dagPath, "--regenerar-studio"); code != 0 {
		t.Fatalf("regenerar: code=%d", code)
	}
	novo, err := os.ReadFile(studioPath)
	if err != nil {
		t.Fatalf("studio.json não regenerado: %v", err)
	}
	if err := processo.ParidadeStudio([]byte(processoScript), novo); err != nil {
		t.Errorf("regenerado não é derivado do DAG: %v", err)
	}

	var s studiodoc.Processo
	_ = json.Unmarshal(novo, &s)
	s.Nodes[0].Name = "drift"
	drift, _ := json.Marshal(s)
	if err := os.WriteFile(studioPath, drift, 0o644); err != nil {
		t.Fatal(err)
	}
	d2, _, _ := deps(t)
	if code := run(t, d2, "validate", dagPath, "--regenerar-studio"); code != 0 {
		t.Fatalf("regenerar sobre drift: code=%d", code)
	}
	depois, _ := os.ReadFile(studioPath)
	if err := processo.ParidadeStudio([]byte(processoScript), depois); err != nil {
		t.Errorf("regeneração não corrigiu o drift: %v", err)
	}
}

// TestValidateSemStudioNaoEscreve prova que validate, sem o flag, não grava nada:
// processo sem documento do Studio passa sem criar arquivo (a escrita é do fabricar).
func TestValidateSemStudioNaoEscreve(t *testing.T) {
	dir := dirDeProcesso(t, "legoloop-semstudio")
	dagPath := caminhos.CaminhoDAG(dir)
	if err := os.WriteFile(dagPath, []byte(processoScript), 0o644); err != nil {
		t.Fatal(err)
	}
	d, _, _ := deps(t)
	if code := run(t, d, "validate", dagPath); code != 0 {
		t.Fatalf("validate sem studio: code=%d", code)
	}
	if _, err := os.Stat(filepath.Join(dir, "legoloop-semstudio"+caminhos.SufixoStudio)); err == nil {
		t.Error("validate escreveu studio.json sem o flag --regenerar-studio")
	}
}
