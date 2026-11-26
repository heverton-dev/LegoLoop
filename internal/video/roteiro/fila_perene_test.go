package roteiro

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFilaPerene_EstocaApenasPerenesEDescartaQuentes(t *testing.T) {
	acervo := AcervoTemas{
		SchemaVersion: 1,
		Temas:         []RegistroTema{{Slug: "historia-ja-feita", Titulo: "Historia ja feita"}},
	}
	coletadas := []PautaCandidata{
		{Titulo: "Noticia urgente de hoje", TipoPauta: "quente", Nivel: "A", Score: 10},
		{Titulo: "Historia do Cafe no Brasil", TipoPauta: "perene", Nivel: "A", Score: 10},
		{Titulo: "Historia ja feita", TipoPauta: "perene", Nivel: "A", Score: 10},
		{Titulo: "Historia do Cafe no Brasil", TipoPauta: "perene", Nivel: "B", Score: 7},
	}

	fila := MesclarFilaPerene(FilaPerene{SchemaVersion: 1}, coletadas, acervo)
	if len(fila.Pautas) != 1 {
		t.Fatalf("esperava 1 pauta perene estocada, obteve %d: %+v", len(fila.Pautas), fila.Pautas)
	}
	if fila.Pautas[0].Slug != "historia-do-cafe-no-brasil" {
		t.Fatalf("slug inesperado: %q", fila.Pautas[0].Slug)
	}
}

func TestFilaPerene_SalvarCarregarERemoverAprovado(t *testing.T) {
	dir := t.TempDir()
	acervoPath := filepath.Join(dir, "acervo", "temas.json")
	filaPath := CaminhoFilaPerene(acervoPath)

	vazia, err := CarregarFilaPerene(filaPath)
	if err != nil || vazia.SchemaVersion != 1 || len(vazia.Pautas) != 0 {
		t.Fatalf("fila ausente deveria retornar vazia: %+v err=%v", vazia, err)
	}

	inicial := FilaPerene{
		Pautas: []PautaCandidata{
			{Titulo: "Curiosidade A", Slug: "curiosidade-a", TipoPauta: "perene", Score: 7},
			{Titulo: "Curiosidade B", Slug: "curiosidade-b", TipoPauta: "perene", Score: 10},
		},
	}
	if err := SalvarFilaPerene(filaPath, inicial); err != nil {
		t.Fatalf("erro ao salvar fila: %v", err)
	}
	lida, err := CarregarFilaPerene(filaPath)
	if err != nil {
		t.Fatalf("erro ao carregar fila: %v", err)
	}
	if len(lida.Pautas) != 2 {
		t.Fatalf("esperava 2 itens na fila, obteve %d", len(lida.Pautas))
	}

	restante := RemoverDaFilaPerene(lida, "Curiosidade A")
	if len(restante.Pautas) != 1 || restante.Pautas[0].Slug != "curiosidade-b" {
		t.Fatalf("remocao incorreta: %+v", restante.Pautas)
	}

	// Erro de JSON invalido
	invalido := filepath.Join(dir, "invalido.json")
	_ = os.WriteFile(invalido, []byte("{invalido"), 0o644)
	if _, err := CarregarFilaPerene(invalido); err == nil {
		t.Fatalf("esperava erro ao carregar JSON invalido")
	}
	// Mesclar preservando itens ja existentes na fila e ignorando vazios
	mesclada := MesclarFilaPerene(lida, []PautaCandidata{
		{Titulo: "   ", TipoPauta: "perene"},
		{Titulo: "Nova Perene C", TipoPauta: "perene"},
	}, AcervoTemas{SchemaVersion: 1})
	if len(mesclada.Pautas) != 3 {
		t.Fatalf("esperava 3 itens apos mesclar com atual, obteve %d", len(mesclada.Pautas))
	}
	// CarregarFilaPerene com schema_version 0 deve normalizar para 1
	semVersao := filepath.Join(dir, "sem_versao.json")
	_ = os.WriteFile(semVersao, []byte(`{"pautas":[]}`), 0o644)
	norm, err := CarregarFilaPerene(semVersao)
	if err != nil || norm.SchemaVersion != 1 {
		t.Fatalf("esperava schema_version 1, obteve %+v err=%v", norm, err)
	}
	// Erro de leitura em diretorio e escrita sobre arquivo bloqueado como pai
	if _, err := CarregarFilaPerene(dir); err == nil {
		t.Fatalf("esperava erro ao ler diretorio como arquivo")
	}
	if err := SalvarFilaPerene(filepath.Join(invalido, "fila.json"), inicial); err == nil {
		t.Fatalf("esperava erro ao criar subpasta sob arquivo")
	}
	if err := baixarTemaDaFilaPerene(filepath.Join(dir, "inexistente", "temas.json"), "x"); err != nil {
		t.Fatalf("baixar tema com fila ausente deveria retornar nil, obteve: %v", err)
	}
}

func TestFilaPerene_CombinarSemDuplicarEIgnorandoProduzidos(t *testing.T) {
	acervo := AcervoTemas{
		SchemaVersion: 1,
		Temas:         []RegistroTema{{Slug: "ja-produzido", Titulo: "Ja Produzido"}},
	}
	fila := FilaPerene{
		SchemaVersion: 1,
		Pautas: []PautaCandidata{
			{Titulo: "Perene Estocada", Slug: "perene-estocada", TipoPauta: "perene", Score: 8},
			{Titulo: "Ja Produzido", Slug: "ja-produzido", TipoPauta: "perene", Score: 8},
		},
	}
	coletadas := []PautaCandidata{
		{Titulo: "Quente Nova", TipoPauta: "quente", Score: 10},
		{Titulo: "Perene Estocada", TipoPauta: "perene", Score: 8},
	}
	comb := CombinarPautasComFila(coletadas, fila, acervo)
	if len(comb) != 2 {
		t.Fatalf("esperava 2 candidatas combinadas, obteve %d: %+v", len(comb), comb)
	}
}

func TestProcesso_EstocaPerenesERemoveAprovadaDaFila(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/json" {
			_, _ = w.Write([]byte(`[
				{"title":"Historia do Ouro no Brasil","url":"https://github.com/golang/go/releases"},
				{"title":"Origem da Burocracia Imperial","url":"https://github.com/golang/go/releases"}
			]`))
			return
		}
		_, _ = w.Write([]byte(`<rss><channel><item><title>Noticia Urgente Hoje</title><link>https://reuters.com/urgente</link></item></channel></rss>`))
	}))
	defer srv.Close()

	dir := t.TempDir()
	acervoPath := filepath.Join(dir, "acervo", "temas.json")
	cat := CatalogoPadrao()
	cat.Conectores = []ConectorFonte{
		{ID: "perene-json", Nome: "Perene", Tipo: "json", URL: srv.URL + "/json", Perfis: []string{"tecnologia"}},
		{ID: "quente-rss", Nome: "Quente", Tipo: "rss", URL: srv.URL + "/rss", Perfis: []string{"tecnologia"}},
	}
	proj := projetoValidoTeste()
	proj.PoliticaGates = PoliticaGates{Modo: "mecanico"}

	proc := &Processo{
		DepsPreVoo: DepsPreVoo{
			ChecarLegoloop: func(context.Context) error { return nil },
			ChecarCofreRea: func() bool { return true },
			ChecarOpcional: func(string) bool { return true },
			HTTPClient:     srv.Client(),
		},
		GerarRoteiro: func(_ context.Context, p PautaCandidata, _ *RelatorioValidacao) (Script, string, error) {
			return Script{
				SchemaVersion: 1,
				Title:         p.Titulo,
				Sources:       []Fonte{{ID: "FA", URL: "https://github.com/golang/go/releases"}},
				Lines:         []LinhaRoteiro{{Text: "Fato comprovado.", Natureza: "fato_consolidado", Fontes: []string{"FA"}}},
			}, "# Plano\n", nil
		},
	}
	res, err := proc.Executar(context.Background(), OpcoesExecucao{
		ProjetoDir: dir, AcervoPath: acervoPath, Projeto: proj, Marca: marcaValidaTeste(), Catalogo: cat,
	})
	if err != nil {
		t.Fatalf("execucao falhou: %v", err)
	}
	if !res.TemaRegistrado {
		t.Fatalf("esperava tema registrado em acervo/temas.json")
	}
	fila, err := CarregarFilaPerene(CaminhoFilaPerene(acervoPath))
	if err != nil {
		t.Fatalf("erro ao ler fila_perene.json: %v", err)
	}
	// Das 2 perenes coletadas, 1 foi escolhida e aprovada (saiu da fila), sobrando 1; a quente nao entrou na fila.
	if len(fila.Pautas) != 1 {
		t.Fatalf("esperava 1 pauta perene restante na fila, obteve %d: %+v", len(fila.Pautas), fila.Pautas)
	}
	if fila.Pautas[0].TipoPauta != "perene" {
		t.Fatalf("item na fila deveria ser perene, obteve %q", fila.Pautas[0].TipoPauta)
	}
	if fila.Pautas[0].Slug == res.PautaEscolhida.Slug {
		t.Fatalf("pauta aprovada %q nao deveria continuar na fila", res.PautaEscolhida.Slug)
	}
}
