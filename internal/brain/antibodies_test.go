package brain

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/orchestrix"
	"legoloop/internal/provider"
)

// TestAnticorposPersistencia prova AC-E02.1: persistência em anticorpos.json do Processo.
func TestAnticorposPersistencia(t *testing.T) {
	dir := t.TempDir()

	db, err := NewAntibodiesDBPersistente(dir)
	if err != nil {
		t.Fatalf("NewAntibodiesDBPersistente: %v", err)
	}
	if err := db.Record("pkgA", "E1", "causa raiz", "diff solução"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(dir, "anticorpos.json"))
	if err != nil {
		t.Fatalf("arquivo anticorpos.json não gravado: %v", err)
	}
	if !json.Valid(raw) {
		t.Fatalf("anticorpos.json inválido: %s", raw)
	}

	db2, err := NewAntibodiesDBPersistente(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	ab, ok := db2.Find("pkgA", "E1")
	if !ok {
		t.Fatal("Find não localizou anticorpo persistido")
	}
	if ab.RootCause != "causa raiz" || ab.DiffSolution != "diff solução" {
		t.Errorf("tripla alterada: %+v", ab)
	}
}

// TestAnticorposCorrompidos prova AC-E02.2: JSON inválido devolve ErrAnticorposCorrompidos sem sobrescrever.
func TestAnticorposCorrompidos(t *testing.T) {
	dir := t.TempDir()
	caminho := filepath.Join(dir, "anticorpos.json")
	lixo := []byte("{json inválido !!")
	if err := os.WriteFile(caminho, lixo, 0o644); err != nil {
		t.Fatalf("preparação: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic em carregamento corrompido: %v", r)
		}
	}()

	_, err := NewAntibodiesDBPersistente(dir)
	if !errors.Is(err, ErrAnticorposCorrompidos) {
		t.Fatalf("erro esperado ErrAnticorposCorrompidos, obtido: %v", err)
	}

	depois, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("arquivo sumiu após erro: %v", err)
	}
	if string(depois) != string(lixo) {
		t.Errorf("arquivo foi sobrescrito: %s", depois)
	}
}

// TestVacinaReincidencia prova AC-E02.3: Record deduplica por errSig; reincidência vacina com antídoto.
func TestVacinaReincidencia(t *testing.T) {
	dir := t.TempDir()

	db, err := NewAntibodiesDBPersistente(dir)
	if err != nil {
		t.Fatalf("NewAntibodiesDBPersistente: %v", err)
	}

	// Executor simulado: falha com "E1" enquanto o antídoto não estiver aplicado.
	rodar := func() error {
		if ant, ok := db.AntidotoAplicado("pkgA", "E1"); ok && ant != "" {
			return nil // antídoto aplicado: o erro não se repete
		}
		return errors.New("E1: falha conhecida")
	}

	// 1ª ocorrência: o erro acontece e é registrado.
	if err := rodar(); err == nil {
		t.Fatal("erro esperado na 1ª execução")
	}
	if err := db.Record("pkgA", "E1", "causa raiz", "diff solução"); err != nil {
		t.Fatalf("Record 1: %v", err)
	}

	// Antes da vacinação o erro ainda se repete.
	if err := rodar(); err == nil {
		t.Fatal("erro resolvido sem vacinação: antídoto aplicado cedo demais")
	}

	// 2ª ocorrência: Record detecta reincidência e vacina com antídoto.
	if err := rodar(); err == nil {
		t.Fatal("erro esperado na 2ª execução")
	}
	if err := db.Record("pkgA", "E1", "causa raiz", "diff solução"); err != nil {
		t.Fatalf("Record 2: %v", err)
	}

	ab, ok := db.Find("pkgA", "E1")
	if !ok {
		t.Fatal("Find não localizou anticorpo")
	}
	if ab.Recorrencias != 2 {
		t.Errorf("Recorrencias = %d, esperado 2", ab.Recorrencias)
	}

	// Prova de vacinação: o mesmo erro não se repete após o antídoto.
	if err := rodar(); err != nil {
		t.Fatalf("erro repetiu após vacinação: %v", err)
	}
}

// TestRegistroInvalido prova AC-E02.4: errSig ou pacote vazio devolve ErrAnticorpoInvalido sem alterar o banco.
func TestRegistroInvalido(t *testing.T) {
	dir := t.TempDir()

	db, err := NewAntibodiesDBPersistente(dir)
	if err != nil {
		t.Fatalf("NewAntibodiesDBPersistente: %v", err)
	}
	if err := db.Record("pkgA", "E1", "causa raiz", "diff solução"); err != nil {
		t.Fatalf("Record válido: %v", err)
	}
	antes, err := os.ReadFile(filepath.Join(dir, "anticorpos.json"))
	if err != nil {
		t.Fatalf("lendo banco antes: %v", err)
	}

	casos := []struct{ pkg, errSig string }{
		{"pkgA", ""},
		{"", "E1"},
		{"pkgA", "   "},
	}
	for _, c := range casos {
		if err := db.Record(c.pkg, c.errSig, "causa", "solução"); !errors.Is(err, ErrAnticorpoInvalido) {
			t.Errorf("Record(%q, %q) = %v, esperado ErrAnticorpoInvalido", c.pkg, c.errSig, err)
		}
	}

	depois, err := os.ReadFile(filepath.Join(dir, "anticorpos.json"))
	if err != nil {
		t.Fatalf("lendo banco depois: %v", err)
	}
	if string(depois) != string(antes) {
		t.Errorf("banco persistido foi alterado:\nantes: %s\ndepois: %s", antes, depois)
	}
}

// TestInjecaoPreventiva prova AC-E02.5/AC-E02.6: o contexto preventivo vira passo do processo e
// entra antes do Nó 1; sem anticorpos o contexto do Nó 1 segue inalterado.
func TestInjecaoPreventiva(t *testing.T) {
	t.Run("ComAnticorpos_InjetaAntesDoNo1", func(t *testing.T) {
		db := NewAntibodiesDB()
		if err := db.Record("pkgA", "E1", "causa raiz", "diff solução"); err != nil {
			t.Fatalf("Record: %v", err)
		}

		contextos := []string{"CONTEXTO DO NÓ 1", "CONTEXTO DO NÓ 2"}
		resultado := db.InjecaoPreventiva("pkgA", contextos)
		if len(resultado) != 2 {
			t.Fatalf("esperava 2 contextos de nó, obteve %d", len(resultado))
		}

		no1 := resultado[0]
		idxAdv := strings.Index(no1, "ANTICORPO PREVENTIVO")
		idxNo1 := strings.Index(no1, "CONTEXTO DO NÓ 1")
		if idxAdv < 0 {
			t.Fatalf("advertência preventiva ausente do contexto do Nó 1: %s", no1)
		}
		if idxNo1 < 0 {
			t.Fatalf("contexto do Nó 1 ausente: %s", no1)
		}
		if idxAdv > idxNo1 {
			t.Errorf("advertência deve vir antes do contexto do Nó 1: adv em %d, nó em %d", idxAdv, idxNo1)
		}
		if resultado[1] != "CONTEXTO DO NÓ 2" {
			t.Errorf("contexto do Nó 2 não deve ser alterado: %s", resultado[1])
		}
	})

	t.Run("SemAnticorpos_ContextoDoNo1Inalterado", func(t *testing.T) {
		db := NewAntibodiesDB()
		contextos := []string{"CONTEXTO DO NÓ 1", "CONTEXTO DO NÓ 2"}
		resultado := db.InjecaoPreventiva("pkgSem", contextos)
		if resultado[0] != "CONTEXTO DO NÓ 1" {
			t.Errorf("sem anticorpos o contexto do Nó 1 deve seguir inalterado: %s", resultado[0])
		}
		if resultado[1] != "CONTEXTO DO NÓ 2" {
			t.Errorf("contexto do Nó 2 não deve ser alterado: %s", resultado[1])
		}
	})
}

// TestInjecaoPreventivaNo1CaminhoReal prova AC-E02.5/AC-E02.6 no caminho real de execução
// (Orchestrix.Execute → prepareSkillPrompt → despacho HTTP): o contexto preventivo entra
// antes do conteúdo do Nó 1; sem anticorpos o prompt do Nó 1 segue inalterado.
func TestInjecaoPreventivaNo1CaminhoReal(t *testing.T) {
	t.Run("ComAnticorpos_InjetaAntesDoConteudoDoNo1", func(t *testing.T) {
		prompt := executarNo1Real(t, true)
		idxAdv := strings.Index(prompt, "ANTICORPO PREVENTIVO")
		idxNo1 := strings.Index(prompt, "Orchestrix isolada")
		if idxAdv < 0 {
			t.Fatalf("advertência preventiva ausente do prompt real do Nó 1: %s", prompt)
		}
		if idxNo1 < 0 {
			t.Fatalf("conteúdo do Nó 1 ausente do prompt real: %s", prompt)
		}
		if idxAdv > idxNo1 {
			t.Errorf("advertência deve vir antes do conteúdo do Nó 1: adv em %d, nó em %d", idxAdv, idxNo1)
		}
	})

	t.Run("SemAnticorpos_PromptDoNo1Inalterado", func(t *testing.T) {
		prompt := executarNo1Real(t, false)
		if strings.Contains(prompt, "ANTICORPO PREVENTIVO") {
			t.Errorf("sem anticorpos o prompt do Nó 1 não deve conter advertência: %s", prompt)
		}
		if !strings.Contains(prompt, "Orchestrix isolada") {
			t.Fatalf("prompt do Nó 1 sem conteúdo: %s", prompt)
		}
	})
}

// executarNo1Real roda o Nó 1 pelo caminho real (Orchestrix.Execute) e devolve o prompt
// capturado do despacho HTTP ao roteador.
func executarNo1Real(t *testing.T, comAnticorpo bool) string {
	t.Helper()
	abDB := NewAntibodiesDB()
	if comAnticorpo {
		if err := abDB.Record("skills/interview", "E1", "causa raiz", "diff solução"); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	var prompt string
	srv := servidorCapturaPrompt(t, &prompt)
	defer srv.Close()
	ox := &orchestrix.Orchestrix{
		Inventory:    provider.Inventory{RouterURL: srv.URL, RouterModels: []string{"code-pro"}},
		HTTP:         srv.Client(),
		RouterKey:    "k9",
		RouterKeyURL: srv.URL,
		SkillsDir:    skillsDirDeTeste(t),
		Antibodies:   abDB,
	}
	n := dag.Node{ID: "n1", Type: dag.TypeSkill, Ref: "skills/interview",
		Routing: dag.RoutingConfig{Mode: dag.ModeAuto, Combo: "code-pro"},
		Output:  dag.HandoffContract{Artifact: "out/BRIEFING.md", Gate: "true"}}
	if _, err := ox.Execute(context.Background(), n, t.TempDir()); err != nil {
		t.Fatalf("Execute do Nó 1: %v", err)
	}
	return prompt
}

// servidorCapturaPrompt expõe um roteador falso que grava o prompt do despacho.
func servidorCapturaPrompt(t *testing.T, prompt *string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corpo, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("ler body: %v", err)
		}
		*prompt = string(corpo)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"# SPEC"}}],"usage":{"prompt_tokens":10,"completion_tokens":5}}`))
	}))
}

// skillsDirDeTeste cria uma Skill mínima para o Nó 1 do caminho real.
func skillsDirDeTeste(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "skills", "interview")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p, "SKILL.md"), []byte("Instrucao de teste."), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}
