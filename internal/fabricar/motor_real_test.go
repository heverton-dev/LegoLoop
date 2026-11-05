package fabricar_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
)

// TestMotorRealSemMatrizNaoGastaOCerebro: sem matriz o fabricar para antes da entrevista, sem nenhuma chamada.
func TestMotorRealSemMatrizNaoGastaOCerebro(t *testing.T) {
	r := newRig(t, goodPieces())
	r.deps.Matriz = nil // padrão: lê o cofre, que está vazio
	r.req.Interactive = true
	r.deps.Ask = func(string, string) (string, error) { t.Fatal("não deve perguntar nada sem matriz"); return "", nil }
	_, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
	if !errors.Is(err, fabricar.ErrSemMatriz) || !strings.Contains(err.Error(), "legoloop matriz gerar") {
		t.Fatalf("esperava ErrSemMatriz com a instrução de gerar a matriz, veio %v", err)
	}
	if r.llm.total() != 0 {
		t.Errorf("o cérebro foi chamado %d vezes sem matriz", r.llm.total())
	}
}

// TestMotorRealLeMatrizEVocabularioDoCofre: LerMatriz usa generated_at e o vocab.json do cofre.
func TestMotorRealLeMatrizEVocabularioDoCofre(t *testing.T) {
	vault := t.TempDir()
	reg := `{"generated_at":"2026-11-01T10:00:00Z","matrix_version":"2.0.0","total_pieces":1,"ratings":[],` +
		`"entradas":[{"ref":"x","nome":"alfa","tipo":"skill","hash":"h"}]}`
	voc := `{"versao":7,"etapas":[{"nome":"construir"}],"naturezas":[{"nome":"software"}],"casos_de_uso":[]}`
	for name, data := range map[string]string{"registry_matrix.json": reg, "vocab.json": voc} {
		if err := os.WriteFile(filepath.Join(vault, name), []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	f, err := fabricar.LerMatriz(vault)
	if err != nil {
		t.Fatal(err)
	}
	if !f.GeradaEm.Equal(time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)) || len(f.Entradas) != 1 || f.Vocab.Versao != 7 {
		t.Errorf("matriz lida errada: %+v", f)
	}
	if _, err := fabricar.LerMatriz(t.TempDir()); !errors.Is(err, fabricar.ErrSemMatriz) {
		t.Errorf("cofre sem matriz deve dar ErrSemMatriz, veio %v", err)
	}
	quebrado := t.TempDir()
	if err := os.WriteFile(filepath.Join(quebrado, "registry_matrix.json"), []byte("{quebrado"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := fabricar.LerMatriz(quebrado); err == nil || errors.Is(err, fabricar.ErrSemMatriz) {
		t.Errorf("matriz corrompida deve dar erro próprio, não ErrSemMatriz: %v", err)
	}
}

// TestMotorRealEscolhePelaAptidaoEExplica: a decisão vem da matriz (aptidão), o Score não decide, e o
// relatório traz as necessidades do cérebro e nenhuma decisão inventada.
func TestMotorRealEscolhePelaAptidaoEExplica(t *testing.T) {
	alta := skillApta("site", filepath.Join(os.TempDir(), "s", "site", "SKILL.md"), 0.99, "publicar")
	alta.Perfil.Naturezas = []string{"midia"}
	r := newRig(t, append(goodPieces(), alta))
	res := r.run(t)
	f := res.Report.Fusao
	if f.Necessidades.Natureza != "software" || !slices.Equal(f.Necessidades.Etapas, []string{"construir", "verificar"}) {
		t.Errorf("as necessidades devem vir do cérebro: %+v", f.Necessidades)
	}
	var usadas []string
	for _, p := range res.CuratedPieces {
		usadas = append(usadas, p.Name)
	}
	if !slices.Equal(usadas, []string{"alfa", "beta"}) {
		t.Errorf("a escolha deve seguir a aptidão, não o Score 0.99 de site: %v", usadas)
	}
	for _, d := range f.Decisoes {
		if d.Motivo == "" || d.Motivo == "selecionada na fusão" {
			t.Errorf("decisão sem motivo real: %+v", d)
		}
	}
	if !strings.Contains(r.llm.needsPrompt, r.req.Goal) {
		t.Errorf("o objetivo deve chegar ao cérebro que dá as necessidades:\n%s", r.llm.needsPrompt)
	}
	if !slices.Contains(res.Report.Varredura.Cortes, "site") {
		t.Errorf("a peça excluída deve aparecer nos cortes da Camada 2: %v", res.Report.Varredura.Cortes)
	}
	if r.llm.needsCalls != 1 || r.llm.explainCalls != 1 {
		t.Errorf("o cérebro deve dar as necessidades e explicar uma vez cada: %d %d", r.llm.needsCalls, r.llm.explainCalls)
	}
}

// TestMotorRealIgnoraDrivesComAviso: o pedido antigo com target_drives continua aceito, mas não varre nada.
func TestMotorRealIgnoraDrivesComAviso(t *testing.T) {
	r := newRig(t, goodPieces())
	r.req.TargetDrives = []string{`C:\`}
	res := r.run(t)
	if !slices.ContainsFunc(res.Warnings, func(w string) bool { return strings.Contains(w, "não varre drives") }) {
		t.Errorf("os drives pedidos devem virar aviso: %v", res.Warnings)
	}
}

// TestMotorRealObjetivoChegaAoCerebroMesmoComEspelhoDiferente: com a entrevista feita no Chat, o espelho
// pode não repetir o objetivo; mesmo assim o objetivo do pedido chega ao cérebro que dá as necessidades.
func TestMotorRealObjetivoChegaAoCerebroMesmoComEspelhoDiferente(t *testing.T) {
	r := newRig(t, goodPieces())
	r.deps.Interview = &fabricar.Clarification{Mirror: interview.Mirror{Outcome: "um resultado descrito de outro jeito"},
		Confidence: 96, Reached: true, Interviewed: true}
	r.run(t)
	if !strings.Contains(r.llm.needsPrompt, r.req.Goal) {
		t.Errorf("o objetivo do pedido deve chegar ao cérebro:\n%s", r.llm.needsPrompt)
	}
}

// TestMotorRealDataDaMatrizSoComODia: a matriz real grava generated_at como "2026-11-04" (só o dia). A idade
// vem dessa data, nunca da data do arquivo (que muda numa cópia).
func TestMotorRealDataDaMatrizSoComODia(t *testing.T) {
	vault := t.TempDir()
	reg := `{"generated_at":"2026-11-04","matrix_version":"2.0.0","total_pieces":0,"ratings":[]}`
	if err := os.WriteFile(filepath.Join(vault, "registry_matrix.json"), []byte(reg), 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := fabricar.LerMatriz(vault)
	if err != nil {
		t.Fatal(err)
	}
	if !f.GeradaEm.Equal(time.Date(2026, 11, 4, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("generated_at só com o dia deve valer como data da matriz: %v", f.GeradaEm)
	}
}
