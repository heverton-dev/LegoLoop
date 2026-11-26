package inventario

import (
	"errors"
	"reflect"
	"testing"
)

func estiloValidado(t *testing.T, s Store, id, herda string) {
	t.Helper()
	e := Estilo{SchemaVersion: 1, ID: id, Status: "validado", Herda: herda, Renderizador: "molde-pilula",
		PromptMestre: "p", Paleta: []string{"#fff"}, Referencias: []Referencia{{Tipo: "imagem", Origem: "a.png"}}}
	if err := gravarJSON(s.caminhoEstilo(id), e); err != nil {
		t.Fatal(err)
	}
}

func gravar(t *testing.T, s Store, it Item) {
	t.Helper()
	if err := s.GravarItem(it); err != nil {
		t.Fatalf("GravarItem %s: %v", it.ID, err)
	}
}

func ids(itens []Item) []string {
	out := []string{}
	for _, it := range itens {
		out = append(out, it.ID)
	}
	return out
}

// AC-E07.4: filtra estilo, tipo e todas as tags; só aprovados e não obsoletos; ordena por id.
func TestBuscaFiltraEOrdena(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	estiloValidado(t, s, "pilula", "")
	estiloValidado(t, s, "palito", "")
	gravar(t, s, Item{ID: "c-traje", Tipo: "traje", Estilo: "pilula", Tags: []string{"medieval", "redc"}, Aprovado: true})
	gravar(t, s, Item{ID: "a-traje", Tipo: "traje", Estilo: "pilula", Tags: []string{"medieval", "redc", "x"}, Aprovado: true})
	gravar(t, s, Item{ID: "b-traje", Tipo: "traje", Estilo: "pilula", Tags: []string{"medieval"}, Aprovado: true})
	gravar(t, s, Item{ID: "d-traje", Tipo: "traje", Estilo: "pilula", Tags: []string{"medieval", "redc"}, Aprovado: true, Obsoleto: true})
	gravar(t, s, Item{ID: "e-pele", Tipo: "pele", Estilo: "pilula", Tags: []string{"medieval", "redc"}, Aprovado: true})
	gravar(t, s, Item{ID: "f-traje", Tipo: "traje", Estilo: "pilula", Tags: []string{"medieval", "redc"}})
	gravar(t, s, Item{ID: "g-traje", Tipo: "traje", Estilo: "palito", Tags: []string{"medieval", "redc"}, Aprovado: true})
	got, err := s.Buscar("pilula", "traje", []string{"redc", "medieval"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a-traje", "c-traje"}; !reflect.DeepEqual(ids(got), want) {
		t.Fatalf("Buscar = %v, quer %v", ids(got), want)
	}
}

// AC-E07.4: sem resultado devolve lista vazia e nil.
func TestBuscaSemResultado(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	estiloValidado(t, s, "pilula", "")
	got, err := s.Buscar("pilula", "traje", []string{"nada"})
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("Buscar = %v, %v; quer lista vazia e nil", got, err)
	}
}

func TestBuscaEstiloInexistente(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if _, err := s.Buscar("fantasma", "traje", nil); !errors.Is(err, ErrEstiloInexistente) {
		t.Fatalf("err = %v, quer ErrEstiloInexistente", err)
	}
}

// AC-E07.7: B herda A e recebe só os itens herdáveis de A.
func TestHerancaIncluiHerdaveis(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	estiloValidado(t, s, "estilo-a", "")
	estiloValidado(t, s, "estilo-b", "estilo-a")
	gravar(t, s, Item{ID: "a-herdavel", Tipo: "pele", Estilo: "estilo-a", Aprovado: true, Herdavel: true})
	gravar(t, s, Item{ID: "a-privado", Tipo: "pele", Estilo: "estilo-a", Aprovado: true})
	gravar(t, s, Item{ID: "b-proprio", Tipo: "pele", Estilo: "estilo-b", Aprovado: true})
	got, err := s.Buscar("estilo-b", "pele", nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a-herdavel", "b-proprio"}; !reflect.DeepEqual(ids(got), want) {
		t.Fatalf("Buscar = %v, quer %v", ids(got), want)
	}
}

// AC-E07.7: herança em ciclo retorna ErrHerancaCiclica.
func TestHerancaCiclica(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	estiloValidado(t, s, "estilo-a", "estilo-b")
	estiloValidado(t, s, "estilo-b", "estilo-a")
	if _, err := s.Buscar("estilo-b", "pele", nil); !errors.Is(err, ErrHerancaCiclica) {
		t.Fatalf("err = %v, quer ErrHerancaCiclica", err)
	}
}
