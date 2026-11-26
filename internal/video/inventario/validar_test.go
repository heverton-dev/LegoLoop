package inventario

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func estiloRascunho(id string) Estilo {
	return Estilo{
		SchemaVersion: 1, ID: id, Nome: "Teste", Status: "rascunho", Versao: "0.1.0",
		Renderizador: "molde-pilula", PromptMestre: "diorama", Paleta: []string{"#ffffff"},
		Referencias: []Referencia{{Tipo: "imagem", Origem: "ref.png"}},
	}
}

func gravarAprovacao(t *testing.T, dir string, a Aprovacao) string {
	t.Helper()
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	caminho := filepath.Join(dir, "aprovacao.json")
	if err := os.WriteFile(caminho, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return caminho
}

func novoStoreComEstilo(t *testing.T, id string) Store {
	t.Helper()
	s := Store{Raiz: t.TempDir()}
	if err := s.NovoEstilo(estiloRascunho(id)); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestValidarEstiloComAprovacaoMudaStatus(t *testing.T) {
	s := novoStoreComEstilo(t, "molde-de-pilula")
	arq := gravarAprovacao(t, t.TempDir(), Aprovacao{
		Estilo: "molde-de-pilula", AprovadoPor: "usuario", PersonagensTeste: []string{"rei"},
	})
	agora := time.Date(2026, 11, 26, 10, 0, 0, 0, time.UTC)
	e, err := s.ValidarEstilo("molde-de-pilula", arq, agora)
	if err != nil {
		t.Fatalf("ValidarEstilo: %v", err)
	}
	if e.Status != "validado" {
		t.Fatalf("status = %q, esperado validado", e.Status)
	}
	lido, err := s.LerEstilo("molde-de-pilula")
	if err != nil || lido.Status != "validado" {
		t.Fatalf("estilo gravado: status=%q err=%v", lido.Status, err)
	}
	v, err := s.LerValidacao("molde-de-pilula")
	if err != nil {
		t.Fatalf("LerValidacao: %v", err)
	}
	if v.ValidadoEm != "2026-11-26T10:00:00Z" || v.ArquivoAprovacao != arq {
		t.Fatalf("validacao gravada errada: %+v", v)
	}
}

func TestValidarEstiloSemAprovacaoContinuaRascunho(t *testing.T) {
	s := novoStoreComEstilo(t, "molde-de-pilula")
	ausente := filepath.Join(t.TempDir(), "nao-existe.json")
	_, err := s.ValidarEstilo("molde-de-pilula", ausente, time.Now())
	if !errors.Is(err, ErrAprovacaoAusente) {
		t.Fatalf("erro = %v, esperado ErrAprovacaoAusente", err)
	}
	_, err = s.ValidarEstilo("molde-de-pilula", "", time.Now())
	if !errors.Is(err, ErrAprovacaoAusente) {
		t.Fatalf("caminho vazio: erro = %v, esperado ErrAprovacaoAusente", err)
	}
	lido, err := s.LerEstilo("molde-de-pilula")
	if err != nil || lido.Status != "rascunho" {
		t.Fatalf("status = %q err=%v, esperado rascunho", lido.Status, err)
	}
	if _, err := s.LerValidacao("molde-de-pilula"); err == nil {
		t.Fatal("validacao.json nao deveria existir")
	}
}

func TestValidarEstiloAprovacaoInvalida(t *testing.T) {
	casos := map[string]Aprovacao{
		"sem personagens": {Estilo: "molde-de-pilula", AprovadoPor: "usuario"},
		"outro estilo":    {Estilo: "outro", AprovadoPor: "usuario", PersonagensTeste: []string{"rei"}},
		"sem aprovador":   {Estilo: "molde-de-pilula", PersonagensTeste: []string{"rei"}},
	}
	for nome, a := range casos {
		t.Run(nome, func(t *testing.T) {
			s := novoStoreComEstilo(t, "molde-de-pilula")
			arq := gravarAprovacao(t, t.TempDir(), a)
			_, err := s.ValidarEstilo("molde-de-pilula", arq, time.Now())
			if !errors.Is(err, ErrAprovacaoInvalida) {
				t.Fatalf("erro = %v, esperado ErrAprovacaoInvalida", err)
			}
			lido, err := s.LerEstilo("molde-de-pilula")
			if err != nil || lido.Status != "rascunho" {
				t.Fatalf("status = %q err=%v", lido.Status, err)
			}
		})
	}
}

func TestValidarEstiloInexistente(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	arq := gravarAprovacao(t, t.TempDir(), Aprovacao{Estilo: "nada", AprovadoPor: "u", PersonagensTeste: []string{"x"}})
	_, err := s.ValidarEstilo("nada", arq, time.Now())
	if !errors.Is(err, ErrEstiloInexistente) {
		t.Fatalf("erro = %v, esperado ErrEstiloInexistente", err)
	}
}

func TestValidarEstiloPermiteItemAprovadoDepois(t *testing.T) {
	s := novoStoreComEstilo(t, "molde-de-pilula")
	arq := gravarAprovacao(t, t.TempDir(), Aprovacao{
		Estilo: "molde-de-pilula", AprovadoPor: "usuario", PersonagensTeste: []string{"rei"},
	})
	if _, err := s.ValidarEstilo("molde-de-pilula", arq, time.Now()); err != nil {
		t.Fatal(err)
	}
	it := Item{SchemaVersion: 1, ID: "rei", Tipo: "personagem", Estilo: "molde-de-pilula", Aprovado: true}
	if err := s.GravarItem(it); err != nil {
		t.Fatalf("item aprovado apos validacao: %v", err)
	}
}
