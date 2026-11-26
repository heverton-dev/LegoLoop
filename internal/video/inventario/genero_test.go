package inventario

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// AC-E07.9: vocabulário de gênero aceita masculino, feminino e neutro.
func TestGeneroVocabulario(t *testing.T) {
	for _, g := range []string{"masculino", "feminino", "neutro"} {
		if err := ValidarGenero(g); err != nil {
			t.Fatalf("ValidarGenero(%q) = %v, quer nil", g, err)
		}
	}
	for _, g := range []string{"", "Masculino", "outro"} {
		if err := ValidarGenero(g); !errors.Is(err, ErrGeneroInvalido) {
			t.Fatalf("ValidarGenero(%q) = %v, quer ErrGeneroInvalido", g, err)
		}
	}
}

// AC-E07.9: vocabulário de faixa etária aceita crianca, adulto e idoso.
func TestGeneroFaixaVocabulario(t *testing.T) {
	for _, f := range []string{"crianca", "adulto", "idoso"} {
		if err := ValidarFaixa(f); err != nil {
			t.Fatalf("ValidarFaixa(%q) = %v, quer nil", f, err)
		}
	}
	for _, f := range []string{"", "jovem", "Adulto", "x"} {
		if err := ValidarFaixa(f); !errors.Is(err, ErrFaixaInvalida) {
			t.Fatalf("ValidarFaixa(%q) = %v, quer ErrFaixaInvalida", f, err)
		}
	}
}

// AC-E07.9: personagem exige os dois campos; outro tipo não exige.
func TestGeneroPersonagemExigeCampos(t *testing.T) {
	casos := []struct {
		it   Item
		quer error
	}{
		{Item{Tipo: "personagem", Genero: "feminino", FaixaEtaria: "adulto"}, nil},
		{Item{Tipo: "personagem", FaixaEtaria: "adulto"}, ErrGeneroInvalido},
		{Item{Tipo: "personagem", Genero: "feminino"}, ErrFaixaInvalida},
		{Item{Tipo: "personagem", Genero: "x", FaixaEtaria: "adulto"}, ErrGeneroInvalido},
		{Item{Tipo: "traje"}, nil},
	}
	for i, c := range casos {
		err := ValidarPersonagem(c.it)
		if (c.quer == nil && err != nil) || (c.quer != nil && !errors.Is(err, c.quer)) {
			t.Fatalf("caso %d: ValidarPersonagem = %v, quer %v", i, err, c.quer)
		}
	}
}

// AC-E07.9: item antigo sem os campos continua serializando sem eles.
func TestGeneroOmitempty(t *testing.T) {
	b, err := json.Marshal(Item{ID: "a", Tipo: "traje"})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"genero", "faixa_etaria"} {
		if _, ok := m[k]; ok {
			t.Fatalf("item antigo serializou %q: %s", k, b)
		}
	}
}

func elenco(t *testing.T) Store {
	t.Helper()
	s := Store{Raiz: t.TempDir()}
	estiloValidado(t, s, "pilula", "")
	p := func(id, g, f string) Item {
		return Item{ID: id, Tipo: "personagem", Estilo: "pilula", Genero: g, FaixaEtaria: f, Aprovado: true}
	}
	gravar(t, s, p("ana", "feminino", "adulto"))
	gravar(t, s, p("beto", "masculino", "adulto"))
	gravar(t, s, p("caio", "masculino", "crianca"))
	gravar(t, s, p("dito", "masculino", "idoso"))
	gravar(t, s, Item{ID: "eva-traje", Tipo: "traje", Estilo: "pilula", Aprovado: true})
	return s
}

// AC-E07.9: o filtro devolve só os personagens que casam; filtro vazio vale qualquer valor.
func TestBuscarPersonagensFiltra(t *testing.T) {
	s := elenco(t)
	casos := []struct {
		g, f string
		quer []string
	}{
		{"masculino", "adulto", []string{"beto"}},
		{"masculino", "", []string{"beto", "caio", "dito"}},
		{"", "adulto", []string{"ana", "beto"}},
		{"", "", []string{"ana", "beto", "caio", "dito"}},
		{"neutro", "", []string{}},
	}
	for _, c := range casos {
		got, err := s.BuscarPersonagens("pilula", c.g, c.f)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ids(got), c.quer) {
			t.Fatalf("BuscarPersonagens(%q,%q) = %v, quer %v", c.g, c.f, ids(got), c.quer)
		}
	}
}

// AC-E07.9: filtro fora do vocabulário retorna erro tipado.
func TestBuscarPersonagensFiltroInvalido(t *testing.T) {
	s := elenco(t)
	if _, err := s.BuscarPersonagens("pilula", "outro", ""); !errors.Is(err, ErrGeneroInvalido) {
		t.Fatalf("err = %v, quer ErrGeneroInvalido", err)
	}
	if _, err := s.BuscarPersonagens("pilula", "", "jovem"); !errors.Is(err, ErrFaixaInvalida) {
		t.Fatalf("err = %v, quer ErrFaixaInvalida", err)
	}
}

type fichaPersonagem struct {
	ID          string         `json:"id"`
	Genero      string         `json:"genero"`
	FaixaEtaria string         `json:"faixa_etaria"`
	Aprovado    bool           `json:"aprovado"`
	PoseTeste   []string       `json:"pose_teste"`
	HumorTeste  []string       `json:"humor_teste"`
	Cast        map[string]any `json:"cast"`
}

func lerElencoReal(t *testing.T) map[string]fichaPersonagem {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "inventario", "video", "estilos", "molde-de-pilula", "personagens.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Personagens []fichaPersonagem `json:"personagens"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	out := map[string]fichaPersonagem{}
	for _, p := range doc.Personagens {
		out[p.ID] = p
	}
	return out
}

// AC-E07.9: os seis personagens de teste declaram gênero e faixa válidos e ficam sem aprovação.
func TestGeneroElencoMoldeDePilula(t *testing.T) {
	quer := map[string][2]string{
		"narrador-teste": {"masculino", "adulto"}, "executivo-teste": {"masculino", "adulto"},
		"analista-teste": {"feminino", "adulto"}, "mulher-teste": {"feminino", "adulto"},
		"idoso-teste": {"masculino", "idoso"}, "crianca-teste": {"masculino", "crianca"},
		"menina-teste": {"feminino", "crianca"}, "adolescente-teste": {"masculino", "adolescente"},
		"adolescente-feminina-teste": {"feminino", "adolescente"}, "idosa-teste": {"feminino", "idoso"},
	}
	elenco := lerElencoReal(t)
	for id, gf := range quer {
		p, ok := elenco[id]
		if !ok {
			t.Fatalf("personagem %q ausente", id)
		}
		if p.Genero != gf[0] || p.FaixaEtaria != gf[1] {
			t.Fatalf("%s = %s/%s, quer %s/%s", id, p.Genero, p.FaixaEtaria, gf[0], gf[1])
		}
		if p.Aprovado || len(p.PoseTeste) != 3 || len(p.HumorTeste) != 3 {
			t.Fatalf("%s: aprovado=%t pose=%d humor=%d", id, p.Aprovado, len(p.PoseTeste), len(p.HumorTeste))
		}
	}
}

// AC-E07.9: crianca-teste declara cast.proporcao com cabeca 1.3 e corpo 0.62.
func TestGeneroCriancaProporcao(t *testing.T) {
	p := lerElencoReal(t)["crianca-teste"]
	prop, ok := p.Cast["proporcao"].(map[string]any)
	if !ok || prop["cabeca"] != 1.3 || prop["corpo"] != 0.62 {
		t.Fatalf("crianca-teste cast.proporcao = %v, quer cabeca 1.3 e corpo 0.62", p.Cast["proporcao"])
	}
}
