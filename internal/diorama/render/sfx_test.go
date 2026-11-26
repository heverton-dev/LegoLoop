package render

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/diorama/ativos"
)

// batidaObj e uma batida com camada de aderecos (objeto em destaque) quando comObjeto.
func batidaObj(linha int, ini, fim float64, cenario, camera string, comObjeto bool) Batida {
	b := batidaUnica(linha, ini, fim, camera)
	b.CenarioID = cenario
	if !comObjeto {
		return b
	}
	spr, leg := b.Camadas[1], b.Camadas[2]
	spr.Profundidade, leg.Profundidade = 2, 3
	b.Camadas = []CamadaItem{b.Camadas[0], {ID: "o", Camada: "aderecos", Profundidade: 1, Caixa: ativos.BoundingBox{W: 4, H: 4}}, spr, leg}
	return b
}

func catalogoEmDisco(t *testing.T, dir string) CatalogoSFX {
	t.Helper()
	cat := CatalogoPadrao()
	for _, rel := range cat {
		caminho := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(caminho, wavTeste(1, 100), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return cat
}

func timelineParaSFX() Timeline {
	return timelineDe(
		batidaObj(0, 0, 3.25, "cen_01", "pan", true),
		batidaObj(0, 3.25, 6.5, "cen_01", "zoom", true),
		batidaObj(1, 6.5, 9.75, "cen_01", "shake", true),
		batidaObj(2, 9.75, 13, "cen_02", "pan", true),
		batidaObj(3, 13, 16.25, "cen_02", "pan", false),
	)
}

func ev(tipo string, t float64) EventoSFX {
	return EventoSFX{Tipo: tipo, TempoS: t, Arquivo: CatalogoPadrao()[tipo]}
}

func TestSFXEventosPorRegra(t *testing.T) {
	dir := t.TempDir()
	saida, err := GerarEventosSFX(dir, timelineParaSFX(), catalogoEmDisco(t, dir))
	if err != nil {
		t.Fatal(err)
	}
	esperado := [][]EventoSFX{
		{ev("troca_cena", 0), ev("movimento_objeto", 0.25)},
		nil,
		{ev("destaque", 6.5), ev("movimento_objeto", 6.75)},
		{ev("troca_cena", 9.75), ev("movimento_objeto", 10.0)},
		nil,
	}
	for i, b := range saida.Batidas {
		if !reflect.DeepEqual(b.SFX, esperado[i]) {
			t.Errorf("batida %d: obteve %+v, quer %+v", i, b.SFX, esperado[i])
		}
	}
	if err := saida.Validar(); err != nil {
		t.Fatalf("timeline com SFX deve validar: %v", err)
	}
}

func TestSFXSubstituiEventosAntigosENaoAlteraEntrada(t *testing.T) {
	dir := t.TempDir()
	cat := catalogoEmDisco(t, dir)
	entrada := timelineParaSFX()
	entrada.Batidas[1].SFX = []EventoSFX{{Tipo: "destaque", TempoS: 4, Arquivo: "velho.wav"}}
	a, err := GerarEventosSFX(dir, entrada, cat)
	if err != nil || len(a.Batidas[1].SFX) != 0 {
		t.Fatalf("eventos antigos devem ser substituidos, sem duplicar: %v %+v", err, a.Batidas[1].SFX)
	}
	b, _ := GerarEventosSFX(dir, a, cat)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("gerar de novo sobre o resultado deve ser identico (idempotente)")
	}
	if len(entrada.Batidas[1].SFX) != 1 || entrada.Batidas[1].SFX[0].Arquivo != "velho.wav" {
		t.Fatal("a entrada nao pode ser alterada")
	}
}

func TestSFXObjetoNaoPassaDoFimDaBatida(t *testing.T) {
	dir := t.TempDir()
	saida, err := GerarEventosSFX(dir, timelineDe(batidaObj(0, 0, 0.1, "cen_01", "pan", true)), catalogoEmDisco(t, dir))
	if err != nil {
		t.Fatal(err)
	}
	if got := saida.Batidas[0].SFX[1].TempoS; got != 0.1 {
		t.Fatalf("atraso do objeto deve parar no fim da batida: %v", got)
	}
}

func TestSFXSemArquivoNoCatalogoBloqueia(t *testing.T) {
	dir := t.TempDir()
	completo := catalogoEmDisco(t, dir)
	semEntrada := CatalogoSFX{"movimento_objeto": completo["movimento_objeto"], "destaque": completo["destaque"]}
	semArquivo := CatalogoSFX{"troca_cena": "audio/sfx/nao_existe.wav", "movimento_objeto": completo["movimento_objeto"], "destaque": completo["destaque"]}
	foraDaRaiz := CatalogoSFX{"troca_cena": "../fora.wav", "movimento_objeto": completo["movimento_objeto"], "destaque": completo["destaque"]}
	for nome, cat := range map[string]CatalogoSFX{"sem entrada": semEntrada, "arquivo ausente": semArquivo, "fora da raiz": foraDaRaiz} {
		t.Run(nome, func(t *testing.T) {
			if _, err := GerarEventosSFX(dir, timelineParaSFX(), cat); !errors.Is(err, ErrSFXSemArquivo) {
				t.Fatalf("esperava ErrSFXSemArquivo, obteve %v", err)
			}
		})
	}
}

func TestSFXCatalogoDeRenderJSON(t *testing.T) {
	dir := t.TempDir()
	if cat, err := CarregarCatalogo(dir, ConfigPadrao()); err != nil || !reflect.DeepEqual(cat, CatalogoPadrao()) {
		t.Fatalf("sem render.json usa o catalogo padrao: %v %+v", err, cat)
	}
	escrever := func(rel, conteudo string) {
		if err := os.WriteFile(filepath.Join(dir, rel), []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	escrever("cat.json", `{"troca_cena":"meus/a.wav"}`)
	cfg := ConfigRender{SchemaVersion: 1, SFXCatalogo: "cat.json"}
	if cat, err := CarregarCatalogo(dir, cfg); err != nil || cat["troca_cena"] != "meus/a.wav" {
		t.Fatalf("catalogo proprio: %v %+v", err, cat)
	}
	escrever("cat.json", `{"explosao":"x.wav"}`)
	if _, err := CarregarCatalogo(dir, cfg); !errors.Is(err, ErrSchemaInvalido) {
		t.Fatalf("tipo desconhecido deve falhar: %v", err)
	}
	escrever("cat.json", `{`)
	if _, err := CarregarCatalogo(dir, cfg); !errors.Is(err, ErrSFXSemArquivo) {
		t.Fatalf("json invalido deve falhar com erro tipado: %v", err)
	}
	cfg.SFXCatalogo = "nao_existe.json"
	if _, err := CarregarCatalogo(dir, cfg); !errors.Is(err, ErrSFXSemArquivo) {
		t.Fatalf("catalogo ausente deve falhar: %v", err)
	}
}
