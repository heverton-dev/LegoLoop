package roteiro

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCoberturaCatalogoEAcervo(t *testing.T) {
	dir := t.TempDir()

	if _, err := CarregarCatalogo(filepath.Join(dir, "nao_existe.json")); !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("esperava ErrPreVooBloqueado para arquivo ausente, obteve %v", err)
	}
	invalido := filepath.Join(dir, "invalido.json")
	_ = os.WriteFile(invalido, []byte("{json_quebrado"), 0o644)
	if _, err := CarregarCatalogo(invalido); !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("esperava ErrPreVooBloqueado para JSON invalido, obteve %v", err)
	}
	vazio := filepath.Join(dir, "vazio.json")
	_ = os.WriteFile(vazio, []byte(`{"perfis":[]}`), 0o644)
	if _, err := CarregarCatalogo(vazio); !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("esperava ErrPreVooBloqueado para catalogo vazio, obteve %v", err)
	}
	valido := filepath.Join(dir, "valido.json")
	b, _ := json.Marshal(CatalogoPadrao())
	_ = os.WriteFile(valido, b, 0o644)
	cat, err := CarregarCatalogo(valido)
	if err != nil || len(cat.Perfis) == 0 {
		t.Fatalf("esperava catalogo valido, obteve err=%v", err)
	}
	if ClassificarFonte("://url-invalida", projetoValidoTeste(), cat).Nivel != "C" {
		t.Fatalf("URL invalida deve cair em fail-closed C")
	}
	if ClassificarFonte("https://portal.mec.gov.br/doc", projetoValidoTeste(), cat).Nivel != "A" {
		t.Fatalf("sufixo .gov.br fora da lista explicita deve ser nivel A")
	}
	if _, err := ColetarPautas(context.Background(), projetoValidoTeste(), CatalogoFontes{}, AcervoTemas{}, nil); err != nil {
		t.Fatalf("client nil sem conectores nao deve falhar: %v", err)
	}
	arqBloqueio := filepath.Join(dir, "arquivo_nao_dir")
	_ = os.WriteFile(arqBloqueio, []byte("x"), 0o644)
	if err := SalvarAcervoTemas(filepath.Join(arqBloqueio, "temas.json"), AcervoTemas{}); err == nil {
		t.Fatalf("esperava erro ao salvar acervo sob arquivo regular")
	}
	if !casaListaDominios("sub.reuters.com", []string{"", "reuters.com"}) {
		t.Fatalf("subdominio deveria casar e entrada vazia ser ignorada")
	}
	p := projetoValidoTeste()
	p.FontesExcluidas = []string{"hacker_news"}
	for _, c := range cat.ConectoresParaProjeto(p) {
		if c.ID == "hacker_news" {
			t.Fatalf("conector excluido nao deveria aparecer")
		}
	}

	ac := AcervoTemas{}
	if ac.Contem("   ") {
		t.Fatalf("titulo vazio nao deve constar no acervo")
	}
	acPath := filepath.Join(dir, "acervo", "temas.json")
	if err := SalvarAcervoTemas(acPath, AcervoTemas{Temas: []RegistroTema{{Slug: "s1", Titulo: "T1"}}}); err != nil {
		t.Fatalf("erro ao salvar acervo: %v", err)
	}
	_ = os.WriteFile(acPath, []byte(`{"temas":[{"slug":"s1","titulo":"T1"}]}`), 0o644)
	lido, err := CarregarAcervoTemas(acPath)
	if err != nil || lido.SchemaVersion != 1 || !lido.Contem("T1") {
		t.Fatalf("esperava schema_version 1 e T1 presente, obteve err=%v ac=%+v", err, lido)
	}
	_ = os.WriteFile(acPath, []byte("{quebrado"), 0o644)
	if _, err := CarregarAcervoTemas(acPath); err == nil {
		t.Fatalf("esperava erro em JSON corrompido do acervo")
	}
}

func TestCoberturaColetorJSONEErros(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/json":
			_, _ = w.Write([]byte(`[
				{"title":"Release Oficial","url":"https://github.com/golang/go/releases"},
				{"title":"Noticia Reuters","url":"https://www.reuters.com/tech/1"},
				{"title":"Blog Qualquer","url":"https://blog.exemplo.com/1"}
			]`))
		case "/json_ruim":
			_, _ = w.Write([]byte(`{nao_e_lista}`))
		case "/xml_ruim":
			_, _ = w.Write([]byte(`<rss><quebrado`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	cat := CatalogoPadrao()
	cat.Conectores = []ConectorFonte{
		{ID: "c1", URL: srv.URL + "/json", Tipo: "json", Perfis: []string{"tecnologia"}},
		{ID: "c2", URL: srv.URL + "/json_ruim", Tipo: "json", Perfis: []string{"tecnologia"}},
		{ID: "c3", URL: srv.URL + "/xml_ruim", Tipo: "rss", Perfis: []string{"tecnologia"}},
		{ID: "c4", URL: srv.URL + "/erro", Tipo: "rss", Perfis: []string{"tecnologia"}},
		{ID: "c5", URL: "://url-invalida", Tipo: "rss", Perfis: []string{"tecnologia"}},
	}
	p := projetoValidoTeste()
	p.PerfisFonte = []string{"tecnologia"}
	pautas, err := ColetarPautas(context.Background(), p, cat, AcervoTemas{}, srv.Client())
	if err != nil || len(pautas) != 3 {
		t.Fatalf("esperava 3 pautas JSON (A, B, C), obteve err=%v len=%d", err, len(pautas))
	}
	if pautas[0].Score != 10 || pautas[1].Score != 7 || pautas[2].Score != 2 {
		t.Fatalf("scores inesperados: %+v", pautas)
	}
}

func TestCoberturaPreVooSchemasEValidador(t *testing.T) {
	p := projetoValidoTeste()
	p.DuracaoAlvoS = 0
	if err := ValidarConfiguracao(p, marcaValidaTeste(), CatalogoPadrao()); !errors.Is(err, ErrSchemaInvalido) {
		t.Fatalf("esperava ErrSchemaInvalido para duracao 0")
	}
	m := marcaValidaTeste()
	m.Nome = ""
	if err := ValidarConfiguracao(projetoValidoTeste(), m, CatalogoPadrao()); !errors.Is(err, ErrSchemaInvalido) {
		t.Fatalf("esperava ErrSchemaInvalido para marca.nome vazio")
	}
	m = marcaValidaTeste()
	m.Formatos = nil
	if err := ValidarConfiguracao(projetoValidoTeste(), m, CatalogoPadrao()); !errors.Is(err, ErrSchemaInvalido) {
		t.Fatalf("esperava ErrSchemaInvalido para marca.formatos vazio")
	}
	if !(PoliticaGates{Modo: "humano"}).ExigeAprovacaoHumana("qualquer") {
		t.Fatalf("modo humano deve exigir aprovacao")
	}
	if (PoliticaGates{Modo: "mecanico"}).ExigeAprovacaoHumana("qualquer") {
		t.Fatalf("modo mecanico nao deve exigir aprovacao")
	}
	if !(PoliticaGates{Modo: "hibrido", GatesHumanos: []string{"etapa1_pauta"}}).ExigeAprovacaoHumana("etapa1_pauta") {
		t.Fatalf("modo hibrido deve exigir gate listado")
	}

	// Pre-voo bordas: catalogo vazio, diretorio vazio, sem conectores, URL invalida, defaults
	rel := ExecutarPreVoo(context.Background(), projetoValidoTeste(), CatalogoFontes{}, "", "", DepsPreVoo{
		ChecarLegoloop: func(context.Context) error { return nil },
		ChecarOpcional: func(string) bool { return true },
	})
	if !rel.Bloqueado {
		t.Fatalf("esperava bloqueio com catalogo vazio e pastas vazias")
	}
	if checarCofreVideoRea(func() bool { return true }).Nivel != "OK" {
		t.Fatalf("esperava OK quando legoloop-video-rea existe no Cofre")
	}
	if sondaHTTP200(context.Background(), http.DefaultClient, "://invalida") {
		t.Fatalf("URL invalida nao pode retornar 200")
	}
	depsDef := preencherDepsPadrao(DepsPreVoo{})
	if depsDef.HTTPClient == nil || depsDef.ChecarLegoloop == nil || depsDef.ChecarCofreRea == nil || depsDef.ChecarOpcional == nil {
		t.Fatalf("deps padrao devem ser preenchidas")
	}
	_ = depsDef.ChecarCofreRea()
	_ = depsDef.ChecarOpcional("binario_inexistente_xyz_123")

	// Validador bordas
	cat := CatalogoPadrao()
	if ValidarRoteiro(Script{}, projetoValidoTeste(), cat).Aprovado {
		t.Fatalf("roteiro sem linhas deve reprovar")
	}
	scBordas := Script{
		Sources: []Fonte{{ID: "F1", URL: "https://arxiv.org/abs/2601.1"}},
		Lines: []LinhaRoteiro{
			{Text: "   ", Natureza: "narrativa_transicao"},
			{Text: "Tipo desconhecido", Natureza: "desconhecida"},
			{Text: "Sem fonte", Natureza: "evidencia_preliminar"},
			{Text: "Sem escopo", Natureza: "evidencia_preliminar", Fontes: []string{"F1"}},
			{Text: "Alegacao sem fonte", Natureza: "alegacao_terceiro", Atribuicao: "autor"},
		},
	}
	relVal := ValidarRoteiro(scBordas, projetoValidoTeste(), cat)
	if relVal.Aprovado || len(relVal.Violacoes) != 5 {
		t.Fatalf("esperava 5 violacoes de borda, obteve %+v", relVal.Violacoes)
	}
}

func TestCoberturaProcessoRamosAdicionais(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[
			{"title":"Pauta Fraca","url":"https://blog.x.com/1"},
			{"title":"Pauta Forte","url":"https://github.com/golang/go/releases"}
		]`))
	}))
	defer srv.Close()

	cat := CatalogoPadrao()
	cat.Conectores = []ConectorFonte{{ID: "j", URL: srv.URL, Tipo: "json", Perfis: []string{"tecnologia"}}}
	depsOK := DepsPreVoo{
		ChecarLegoloop: func(context.Context) error { return nil },
		ChecarCofreRea: func() bool { return true },
		ChecarOpcional: func(string) bool { return true },
		HTTPClient:     srv.Client(),
	}

	t.Run("config invalida e pre-voo bloqueado", func(t *testing.T) {
		dir := t.TempDir()
		p := projetoValidoTeste()
		p.Nicho = ""
		pipe := Processo{DepsPreVoo: depsOK}
		if _, err := pipe.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir, AcervoPath: filepath.Join(dir, "a.json"), Projeto: p, Marca: marcaValidaTeste(), Catalogo: cat}); !errors.Is(err, ErrSchemaInvalido) {
			t.Fatalf("esperava ErrSchemaInvalido, obteve %v", err)
		}
		pipeFail := Processo{DepsPreVoo: DepsPreVoo{
			ChecarLegoloop: func(context.Context) error { return ErrPreVooBloqueado },
			ChecarCofreRea: func() bool { return true },
			ChecarOpcional: func(string) bool { return true },
			HTTPClient:     srv.Client(),
		}}
		if _, err := pipeFail.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir, AcervoPath: filepath.Join(dir, "a.json"), Projeto: projetoValidoTeste(), Marca: marcaValidaTeste(), Catalogo: cat}); !errors.Is(err, ErrPreVooBloqueado) {
			t.Fatalf("esperava ErrPreVooBloqueado, obteve %v", err)
		}
		acCorrompido := filepath.Join(dir, "corrompido.json")
		_ = os.WriteFile(acCorrompido, []byte("{invalido"), 0o644)
		if _, err := pipe.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir, AcervoPath: acCorrompido, Projeto: projetoValidoTeste(), Marca: marcaValidaTeste(), Catalogo: cat}); err == nil {
			t.Fatalf("esperava erro ao carregar acervo corrompido em Executar")
		}
	})

	t.Run("sem pautas novas, sem gerador, erro no gerador e recusa humana", func(t *testing.T) {
		dir := t.TempDir()
		acPath := filepath.Join(dir, "acervo", "temas.json")
		_ = SalvarAcervoTemas(acPath, AcervoTemas{Temas: []RegistroTema{
			{Slug: "pauta-fraca", Titulo: "Pauta Fraca"},
			{Slug: "pauta-forte", Titulo: "Pauta Forte"},
		}})
		pipe := Processo{DepsPreVoo: depsOK}
		if _, err := pipe.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir, AcervoPath: acPath, Projeto: projetoValidoTeste(), Marca: marcaValidaTeste(), Catalogo: cat}); !errors.Is(err, ErrTemaRepetido) {
			t.Fatalf("esperava ErrTemaRepetido quando todas as pautas ja existem, obteve %v", err)
		}

		acNovo := filepath.Join(dir, "acervo2", "temas.json")
		if _, err := pipe.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir, AcervoPath: acNovo, Projeto: projetoValidoTeste(), Marca: marcaValidaTeste(), Catalogo: cat}); !errors.Is(err, ErrGateRoteiroReprovado) {
			t.Fatalf("esperava ErrGateRoteiroReprovado sem gerador, obteve %v", err)
		}

		errGen := errors.New("falha gerador")
		pipeErr := Processo{
			DepsPreVoo: depsOK,
			GerarRoteiro: func(context.Context, PautaCandidata, *RelatorioValidacao) (Script, string, error) {
				return Script{}, "", errGen
			},
		}
		if _, err := pipeErr.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir, AcervoPath: acNovo, Projeto: projetoValidoTeste(), Marca: marcaValidaTeste(), Catalogo: cat}); !errors.Is(err, errGen) {
			t.Fatalf("esperava propagacao do erro do gerador, obteve %v", err)
		}

		pHum := projetoValidoTeste()
		pHum.PoliticaGates = PoliticaGates{Modo: "humano"}
		pipeHum := Processo{
			DepsPreVoo:    depsOK,
			EscolherPauta: func(ps []PautaCandidata) (PautaCandidata, error) { return ps[1], nil },
			AprovarHumano: func(string, any) (bool, error) { return false, nil },
			Agora:         func() time.Time { return time.Date(2026, 11, 24, 12, 0, 0, 0, time.UTC) },
			GerarRoteiro: func(ctx context.Context, pauta PautaCandidata, _ *RelatorioValidacao) (Script, string, error) {
				return Script{
					SchemaVersion: 1, Title: pauta.Titulo,
					Sources: []Fonte{{ID: "F1", URL: "https://github.com/golang/go/releases"}},
					Lines:   []LinhaRoteiro{{Text: "Linha valida.", Natureza: "fato_consolidado", Fontes: []string{"F1"}}},
				}, "# plano", nil
			},
		}
		if _, err := pipeHum.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir, AcervoPath: acNovo, Projeto: pHum, Marca: marcaValidaTeste(), Catalogo: cat}); !errors.Is(err, ErrGateRoteiroReprovado) {
			t.Fatalf("esperava ErrGateRoteiroReprovado quando humano recusa, obteve %v", err)
		}
		if pipeHum.dataAtual() != "2026-11-24" {
			t.Fatalf("esperava data fixa 2026-11-24")
		}
	})
}
