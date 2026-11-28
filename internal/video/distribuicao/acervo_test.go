package distribuicao

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var t0 = time.Date(2026, 11, 27, 10, 0, 0, 0, time.UTC)

func itemBase(id, tema, tipo string) Item {
	it := Item{ID: id, Tema: tema, Tipo: tipo, Formatos: []string{"16:9"}, Canais: []string{"youtube"}, Status: StatusPronto, Visibilidade: VisPrivado, Autorizado: true}
	if tipo == TipoQuente {
		it.Validade = t0.Add(30 * time.Hour)
	}
	return it
}

func TestAcervoAusenteEhVazioEGravaAtomico(t *testing.T) {
	dir := t.TempDir()
	a, err := CarregarAcervo(dir)
	if err != nil || len(a.Itens) != 0 || a.SchemaVersion != 1 {
		t.Fatalf("acervo vazio esperado: %+v %v", a, err)
	}
	if err := a.Adicionar(itemBase("v1", "juros compostos", TipoPerene)); err != nil {
		t.Fatal(err)
	}
	if err := a.Gravar(dir); err != nil {
		t.Fatal(err)
	}
	b, err := CarregarAcervo(dir)
	if err != nil || len(b.Itens) != 1 || b.Itens[0].ID != "v1" {
		t.Fatalf("releitura: %+v %v", b, err)
	}
	if _, err := os.Stat(CaminhoAcervo(dir) + ".tmp"); err == nil {
		t.Fatal("arquivo temporario ficou para tras")
	}
}

func TestAcervoRecusaArquivoInvalido(t *testing.T) {
	dir := t.TempDir()
	p := CaminhoAcervo(dir)
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	_ = os.WriteFile(p, []byte("{"), 0o644)
	if _, err := CarregarAcervo(dir); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("quer ErrAcervoInvalido, obteve %v", err)
	}
	// diretorio no lugar do arquivo: erro de leitura que nao e ausencia
	d2 := t.TempDir()
	_ = os.MkdirAll(CaminhoAcervo(d2), 0o755)
	if _, err := CarregarAcervo(d2); err == nil {
		t.Fatal("quer erro de leitura")
	}
}

func TestItemValidar(t *testing.T) {
	casos := map[string]func(*Item){
		"sem id":              func(i *Item) { i.ID = " " },
		"tipo invalido":       func(i *Item) { i.Tipo = "morno" },
		"status invalido":     func(i *Item) { i.Status = "x" },
		"visibilidade":        func(i *Item) { i.Visibilidade = "" },
		"sem canais":          func(i *Item) { i.Canais = nil },
		"quente sem validade": func(i *Item) { i.Tipo = TipoQuente; i.Validade = time.Time{} },
	}
	for nome, mudar := range casos {
		t.Run(nome, func(t *testing.T) {
			it := itemBase("a", "tema", TipoPerene)
			mudar(&it)
			if err := it.Validar(); !errors.Is(err, ErrAcervoInvalido) {
				t.Fatalf("quer ErrAcervoInvalido, obteve %v", err)
			}
		})
	}
}

func TestAdicionarAplicaOriginalidade(t *testing.T) {
	a := Acervo{SchemaVersion: 1}
	if err := a.Adicionar(itemBase("a", "A historia do dolar no Brasil", TipoPerene)); err != nil {
		t.Fatal(err)
	}
	if err := a.Adicionar(itemBase("a", "outro tema", TipoPerene)); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("id repetido: %v", err)
	}
	if err := a.Adicionar(itemBase("b", "a historia do dolar no brasil!", TipoPerene)); !errors.Is(err, ErrTemaRepetido) {
		t.Fatalf("tema repetido: %v", err)
	}
	if err := a.Adicionar(itemBase("c", "Como funciona um cofre", TipoPerene)); err != nil {
		t.Fatalf("tema novo deveria entrar: %v", err)
	}
	derivado := itemBase("a-reel", "A historia do dolar no Brasil", TipoPerene)
	derivado.Origem = "a"
	if err := a.Adicionar(derivado); err != nil {
		t.Fatalf("derivado com origem nao conta como repeticao: %v", err)
	}
	if err := a.Adicionar(Item{}); err == nil {
		t.Fatal("item vazio deveria falhar")
	}
}

func TestAcervoValidarIdsRepetidosEGravarInvalido(t *testing.T) {
	a := Acervo{Itens: []Item{itemBase("a", "x", TipoPerene), itemBase("a", "y", TipoPerene)}}
	if err := a.Validar(); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("ids repetidos: %v", err)
	}
	if err := a.Gravar(t.TempDir()); err == nil {
		t.Fatal("gravar acervo invalido deveria falhar")
	}
	bad := Acervo{Itens: []Item{{}}}
	if err := bad.Validar(); err == nil {
		t.Fatal("item invalido deveria falhar")
	}
}

func TestAutorizarERevogar(t *testing.T) {
	a := Acervo{Itens: []Item{itemBase("a", "x", TipoPerene)}}
	if err := a.Autorizar("a", false); err != nil || a.Itens[0].Autorizado {
		t.Fatalf("revogar: %v %+v", err, a.Itens[0])
	}
	if err := a.Autorizar("zz", true); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("inexistente: %v", err)
	}
	a.Itens[0].Status = StatusPublic
	if err := a.Autorizar("a", true); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("publicado nao muda: %v", err)
	}
}

func TestTemasSemelhantes(t *testing.T) {
	casos := []struct {
		a, b string
		quer bool
	}{
		{"", "x", false}, {"x", "", false},
		{"o primeiro banco do brasil", "O primeiro banco do Brasil", true},
		{"o primeiro banco do brasil", "a queda do banco", false},
		{"a b c d e", "a b c d f", false},
	}
	for _, c := range casos {
		if got := TemasSemelhantes(c.a, c.b); got != c.quer {
			t.Errorf("TemasSemelhantes(%q,%q)=%v, quer %v", c.a, c.b, got, c.quer)
		}
	}
}
