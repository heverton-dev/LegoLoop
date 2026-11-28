package distribuicao

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRascunharEDecidir(t *testing.T) {
	r := Respostas{SchemaVersion: 1}
	c := Comentario{ID: "c1", Rede: "youtube", Autor: "ana", Texto: "otimo video"}
	if err := r.Rascunhar(c, "  Obrigado, Ana!  "); err != nil {
		t.Fatal(err)
	}
	if r.Itens[0].Rascunho != "Obrigado, Ana!" || r.Itens[0].Status != StatusRascunho {
		t.Fatalf("rascunho: %+v", r.Itens[0])
	}
	for nome, f := range map[string]func() error{
		"duplicado":  func() error { return r.Rascunhar(c, "outro") },
		"sem id":     func() error { return r.Rascunhar(Comentario{}, "x") },
		"sem texto":  func() error { return r.Rascunhar(Comentario{ID: "c9"}, " ") },
		"respondido": func() error { return r.Rascunhar(Comentario{ID: "c8", Respondido: true}, "x") },
	} {
		if err := f(); !errors.Is(err, ErrPedidoInvalido) {
			t.Errorf("%s: quer ErrPedidoInvalido, obteve %v", nome, err)
		}
	}
	if len(r.Aprovadas()) != 0 {
		t.Fatal("rascunho nao e aprovado: nada sai sem aprovacao")
	}
	if err := r.Decidir("c1", true, "Obrigado pelo carinho, Ana!"); err != nil {
		t.Fatal(err)
	}
	if ap := r.Aprovadas(); len(ap) != 1 || ap[0].Rascunho != "Obrigado pelo carinho, Ana!" {
		t.Fatalf("aprovada com edicao: %+v", ap)
	}
	if err := r.Decidir("c1", false, ""); !errors.Is(err, ErrPedidoInvalido) {
		t.Fatalf("decisao tomada nao muda: %v", err)
	}
	if err := r.Decidir("zz", true, ""); !errors.Is(err, ErrPedidoInvalido) {
		t.Fatalf("sem rascunho: %v", err)
	}
	_ = r.Rascunhar(Comentario{ID: "c2"}, "resposta")
	if err := r.Decidir("c2", false, ""); err != nil || r.Itens[1].Status != StatusIgnorada {
		t.Fatalf("ignorar: %v %+v", err, r.Itens[1])
	}
	_ = r.Rascunhar(Comentario{ID: "c3"}, "mantida")
	if err := r.Decidir("c3", true, "  "); err != nil || r.Itens[2].Rascunho != "mantida" {
		t.Fatalf("aprovar sem edicao mantem texto: %v %+v", err, r.Itens[2])
	}
}

func TestRespostasGravarECarregar(t *testing.T) {
	dir := t.TempDir()
	r, err := CarregarRespostas(dir)
	if err != nil || len(r.Itens) != 0 {
		t.Fatalf("ausente: %+v %v", r, err)
	}
	_ = r.Rascunhar(Comentario{ID: "c1", Rede: "instagram"}, "oi")
	if err := r.Gravar(dir); err != nil {
		t.Fatal(err)
	}
	b, err := CarregarRespostas(dir)
	if err != nil || len(b.Itens) != 1 || b.Itens[0].Rede != "instagram" {
		t.Fatalf("releitura: %+v %v", b, err)
	}
	ruim := t.TempDir()
	_ = os.MkdirAll(filepath.Dir(CaminhoRespostas(ruim)), 0o755)
	_ = os.WriteFile(CaminhoRespostas(ruim), []byte("{"), 0o644)
	if _, err := CarregarRespostas(ruim); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("json quebrado: %v", err)
	}
	dirComoArquivo := t.TempDir()
	_ = os.MkdirAll(CaminhoRespostas(dirComoArquivo), 0o755)
	if _, err := CarregarRespostas(dirComoArquivo); err == nil {
		t.Fatal("erro de leitura deve voltar")
	}
}

func TestPendentesECarregarComentarios(t *testing.T) {
	dir := t.TempDir()
	caminho := filepath.Join(dir, "c.json")
	_ = os.WriteFile(caminho, []byte(`[{"id":"1","texto":"a"},{"id":"2","texto":"b","respondido":true},{"id":"3","texto":"c"}]`), 0o644)
	cs, err := CarregarComentarios(caminho)
	if err != nil || len(cs) != 3 {
		t.Fatalf("carregar: %+v %v", cs, err)
	}
	r := Respostas{}
	_ = r.Rascunhar(cs[2], "resp")
	p := Pendentes(cs, r)
	if len(p) != 1 || p[0].ID != "1" {
		t.Fatalf("pendentes: %+v", p)
	}
	_ = os.WriteFile(caminho, []byte("nao e json"), 0o644)
	if _, err := CarregarComentarios(caminho); !errors.Is(err, ErrPedidoInvalido) {
		t.Fatalf("json ruim: %v", err)
	}
	if _, err := CarregarComentarios(filepath.Join(dir, "nao.json")); err == nil {
		t.Fatal("arquivo ausente")
	}
}

func metricasBase() []Metrica {
	return []Metrica{
		{ItemID: "a", Tema: "juros", Formato: "9:16", DuracaoS: 40, Retencao: 0.7, Alcance: 1000, Hora: 14, NovosSeguidores: 5},
		{ItemID: "b", Tema: "cofre", Formato: "16:9", DuracaoS: 600, Retencao: 0.3, Alcance: 500, Hora: 19, NovosSeguidores: 50},
		{ItemID: "c", Tema: "juros", Formato: "9:16", DuracaoS: 20, Retencao: 0.9, Alcance: 3000, Hora: 14, NovosSeguidores: 7},
	}
}

func TestGerarAprendizado(t *testing.T) {
	texto, err := GerarAprendizado(metricasBase(), t0)
	if err != nil {
		t.Fatal(err)
	}
	for _, quer := range []string{`"juros" (80%)`, "9:16 (80%)", "ate 30 s (90%)", "14h (2000)", "b (50)", "Gerado em 2026-11-27 com 3 medidas"} {
		if !strings.Contains(texto, quer) {
			t.Fatalf("relatorio sem %q:\n%s", quer, texto)
		}
	}
	if _, err := GerarAprendizado(metricasBase()[:2], t0); !errors.Is(err, ErrPedidoInvalido) {
		t.Fatalf("poucas medidas: %v", err)
	}
	ruins := map[string]func(*Metrica){
		"retencao": func(m *Metrica) { m.Retencao = 1.5 },
		"hora":     func(m *Metrica) { m.Hora = 24 },
		"duracao":  func(m *Metrica) { m.DuracaoS = 0 },
		"sem item": func(m *Metrica) { m.ItemID = "" },
		"alcance":  func(m *Metrica) { m.Alcance = -1 },
	}
	for nome, mudar := range ruins {
		ms := metricasBase()
		mudar(&ms[1])
		if _, err := GerarAprendizado(ms, t0); !errors.Is(err, ErrPedidoInvalido) {
			t.Errorf("%s: quer ErrPedidoInvalido, obteve %v", nome, err)
		}
	}
}

func TestFaixasDeDuracaoEEmpates(t *testing.T) {
	for s, quer := range map[float64]string{10: "ate 30 s", 45: "30 a 60 s", 100: "1 a 3 min", 900: "mais de 3 min"} {
		if got := faixaDuracao(s); got != quer {
			t.Errorf("faixaDuracao(%v)=%q, quer %q", s, got, quer)
		}
	}
	ms := metricasBase()
	for i := range ms {
		ms[i].Retencao = 0.5
	}
	texto, err := GerarAprendizado(ms, t0)
	if err != nil || !strings.Contains(texto, `"cofre" (50%)`) {
		t.Fatalf("empate resolve em ordem alfabetica: %v\n%s", err, texto)
	}
}

func TestGravarAprendizado(t *testing.T) {
	dir := t.TempDir()
	caminho, err := GravarAprendizado(dir, metricasBase(), t0)
	if err != nil || filepath.Base(caminho) != "aprendizado.md" {
		t.Fatalf("gravar: %q %v", caminho, err)
	}
	if dados, _ := os.ReadFile(caminho); !strings.Contains(string(dados), "# Aprendizado semanal") {
		t.Fatalf("conteudo: %s", dados)
	}
	if _, err := GravarAprendizado(dir, nil, t0); err == nil {
		t.Fatal("sem medidas deve falhar")
	}
}
