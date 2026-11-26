package inventario

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestStore(t *testing.T) {
	t.Run("AC-E07.1 NovoEstilo grava estilo.json", testNovoEstiloGrava)
	t.Run("AC-E07.1 id repetido retorna ErrEstiloDuplicado sem alterar", testNovoEstiloDuplicado)
	t.Run("AC-E07.1 estilo invalido nao grava", testNovoEstiloInvalido)
	t.Run("AC-E07.3 rascunho recusa item aprovado", testRascunhoRecusaAprovado)
	t.Run("AC-E07.3 rascunho grava item nao aprovado", testRascunhoGravaNaoAprovado)
	t.Run("AC-E07.3 validado grava item aprovado", testValidadoGravaAprovado)
	t.Run("GravarItem recusa estilo inexistente e tipo invalido", testGravarItemRecusas)
}

func estiloBase(id string) Estilo {
	return Estilo{
		SchemaVersion: 1, ID: id, Nome: "Molde", Status: "rascunho", Versao: "1.0.0",
		Renderizador: "molde-pilula", PromptMestre: "pilula 3d", Paleta: []string{"#fff"},
		Referencias: []Referencia{{Tipo: "imagem", Origem: "ref.png"}},
	}
}

func itemBase(estilo string, aprovado bool) Item {
	return Item{SchemaVersion: 1, ID: "coque", Tipo: "penteado", Estilo: estilo, EstiloVersao: "1.0.0", Origem: "p/2026-11-26", Aprovado: aprovado}
}

func lerArquivo(t *testing.T, caminho string) []byte {
	t.Helper()
	b, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func semTemporarios(t *testing.T, dir string) {
	t.Helper()
	ents, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range ents {
		if filepath.Ext(e.Name()) != ".json" {
			t.Fatalf("arquivo temporario restante: %s", e.Name())
		}
	}
}

func testNovoEstiloGrava(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if err := s.NovoEstilo(estiloBase("molde-de-pilula")); err != nil {
		t.Fatal(err)
	}
	caminho := filepath.Join(s.Raiz, "inventario", "video", "estilos", "molde-de-pilula", "estilo.json")
	var e Estilo
	if err := json.Unmarshal(lerArquivo(t, caminho), &e); err != nil {
		t.Fatal(err)
	}
	if e.ID != "molde-de-pilula" || e.Status != "rascunho" {
		t.Fatalf("estilo gravado inesperado: %+v", e)
	}
	semTemporarios(t, filepath.Dir(caminho))
}

func testNovoEstiloDuplicado(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if err := s.NovoEstilo(estiloBase("molde")); err != nil {
		t.Fatal(err)
	}
	caminho := s.caminhoEstilo("molde")
	antes := lerArquivo(t, caminho)
	outro := estiloBase("molde")
	outro.Nome = "Outro"
	if err := s.NovoEstilo(outro); !errors.Is(err, ErrEstiloDuplicado) {
		t.Fatalf("esperado ErrEstiloDuplicado, veio %v", err)
	}
	if !bytes.Equal(antes, lerArquivo(t, caminho)) {
		t.Fatal("arquivo do estilo foi alterado")
	}
}

func testNovoEstiloInvalido(t *testing.T) {
	casos := map[string]func(*Estilo){
		"id nao kebab":       func(e *Estilo) { e.ID = "Molde_Pilula" },
		"sem renderizador":   func(e *Estilo) { e.Renderizador = "" },
		"status validado":    func(e *Estilo) { e.Status = "validado" },
		"sem prompt":         func(e *Estilo) { e.PromptMestre = "" },
		"sem paleta":         func(e *Estilo) { e.Paleta = nil },
		"sem referencia":     func(e *Estilo) { e.Referencias = nil },
		"referencia vazia":   func(e *Estilo) { e.Referencias = []Referencia{{Tipo: "imagem"}} },
		"id com traco duplo": func(e *Estilo) { e.ID = "molde--pilula" },
	}
	for nome, muda := range casos {
		s := Store{Raiz: t.TempDir()}
		e := estiloBase("molde")
		muda(&e)
		if err := s.NovoEstilo(e); !errors.Is(err, ErrEstiloInvalido) {
			t.Fatalf("%s: esperado ErrEstiloInvalido, veio %v", nome, err)
		}
		if _, err := os.Stat(filepath.Join(s.Raiz, "inventario")); !os.IsNotExist(err) {
			t.Fatalf("%s: nada deveria ser gravado (%v)", nome, err)
		}
	}
}

func testRascunhoRecusaAprovado(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if err := s.NovoEstilo(estiloBase("molde")); err != nil {
		t.Fatal(err)
	}
	if err := s.GravarItem(itemBase("molde", true)); !errors.Is(err, ErrEstiloNaoValidado) {
		t.Fatalf("esperado ErrEstiloNaoValidado, veio %v", err)
	}
	if _, err := os.Stat(s.caminhoItem(itemBase("molde", true))); !os.IsNotExist(err) {
		t.Fatalf("item nao deveria existir (%v)", err)
	}
}

func testRascunhoGravaNaoAprovado(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if err := s.NovoEstilo(estiloBase("molde")); err != nil {
		t.Fatal(err)
	}
	it := itemBase("molde", false)
	if err := s.GravarItem(it); err != nil {
		t.Fatal(err)
	}
	var lido Item
	if err := json.Unmarshal(lerArquivo(t, s.caminhoItem(it)), &lido); err != nil {
		t.Fatal(err)
	}
	if lido.ID != "coque" || lido.Aprovado {
		t.Fatalf("item gravado inesperado: %+v", lido)
	}
	semTemporarios(t, filepath.Dir(s.caminhoItem(it)))
}

func testValidadoGravaAprovado(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if err := s.NovoEstilo(estiloBase("molde")); err != nil {
		t.Fatal(err)
	}
	e := estiloBase("molde")
	e.Status = "validado"
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.caminhoEstilo("molde"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.GravarItem(itemBase("molde", true)); err != nil {
		t.Fatalf("estilo validado deveria aceitar item aprovado: %v", err)
	}
}

func testGravarItemRecusas(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if err := s.GravarItem(itemBase("fantasma", false)); !errors.Is(err, ErrEstiloInexistente) {
		t.Fatalf("esperado ErrEstiloInexistente, veio %v", err)
	}
	if err := s.NovoEstilo(estiloBase("molde")); err != nil {
		t.Fatal(err)
	}
	it := itemBase("molde", false)
	it.Tipo = "chapeu"
	if err := s.GravarItem(it); !errors.Is(err, ErrTipoInvalido) {
		t.Fatalf("esperado ErrTipoInvalido, veio %v", err)
	}
	if _, err := os.Stat(filepath.Join(s.Raiz, "inventario", "video", "estilos", "molde", "itens")); !os.IsNotExist(err) {
		t.Fatalf("nada deveria ser gravado (%v)", err)
	}
}
