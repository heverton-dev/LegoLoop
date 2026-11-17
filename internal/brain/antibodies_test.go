package brain

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
