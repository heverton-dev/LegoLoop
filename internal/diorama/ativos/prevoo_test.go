package ativos

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
)

func criarFixtureCiclo1(t *testing.T, estilo, flowUUID string) string {
	t.Helper()
	dir := t.TempDir()
	marca := roteiro.Marca{
		SchemaVersion: 1,
		Nome:          "Canal Teste",
		Formatos:      []string{"16:9"},
		Voz:           roteiro.VozMarca{Motor: "edge-tts", Voz: "pt-BR-ThalitaMultilingualNeural", Velocidade: "+0%"},
	}
	gravarJSONTeste(t, filepath.Join(dir, "projeto.json"), projetoFixture(estilo, flowUUID))
	gravarJSONTeste(t, filepath.Join(dir, "marca.json"), marca)
	gravarJSONTeste(t, filepath.Join(dir, "script.json"), scriptFixture())
	if err := os.WriteFile(filepath.Join(dir, "plano.md"), []byte("# Plano\n- Narrador\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func projetoFixture(estilo, flowUUID string) roteiro.Projeto {
	return roteiro.Projeto{
		SchemaVersion:      1,
		Slug:               "video-teste",
		Canal:              "canal-teste",
		Nicho:              "tecnologia",
		PerfisFonte:        []string{"geral"},
		EstruturaNarrativa: "causal",
		Formatos:           []string{"16:9"},
		DuracaoAlvoS:       60,
		Voz:                "pt-BR-ThalitaMultilingualNeural",
		EstiloVisual:       estilo,
		Animacao:           "diorama",
		PoliticaGates:      roteiro.PoliticaGates{Modo: "mecanico"},
		FlowProjeto:        flowUUID,
	}
}

func scriptFixture() roteiro.Script {
	return roteiro.Script{
		SchemaVersion: 1,
		Title:         "Titulo de Teste",
		Idioma:        "pt-BR",
		Estrutura:     "causal",
		TipoPauta:     "perene",
		Sources:       []roteiro.Fonte{{ID: "f1", URL: "https://ibge.gov.br/dados", Nivel: "A", Entidade: "IBGE", Data: "2026-11-24"}},
		Lines: []roteiro.LinhaRoteiro{
			{Text: "Linha 1 sobre dados oficiais", Natureza: "fato_consolidado", Fontes: []string{"f1"}},
			{Text: "Linha 2 mas o impacto cresce", Natureza: "narrativa_transicao"},
			{Text: "Linha 3 portanto a tendencia muda", Natureza: "narrativa_transicao"},
			{Text: "Linha 4 conclusao clara", Natureza: "narrativa_transicao"},
		},
	}
}

func gravarJSONTeste(t *testing.T, caminho string, v any) {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminho, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPreVooAprovadoModoCodigo(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	rel, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if rel.Bloqueado {
		t.Fatalf("pre-voo nao deveria bloquear no modo codigo: %+v", rel)
	}
}

func TestPreVooAprovadoModoFlowComTrava(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/json/version" {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"Browser":"Chrome"}`)); err != nil {
				t.Error(err)
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	dir := criarFixtureCiclo1(t, "flow", "11111111-2222-3333-4444-555555555555")
	rel, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{
		FlowCDPURL:   srv.URL,
		HTTPClient:   srv.Client(),
		ChecarModulo: func(string) bool { return true },
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if rel.Bloqueado {
		t.Fatalf("pre-voo nao deveria bloquear com CDP 9333 e modulos ativos: %+v", rel)
	}
}

func TestPreVooBloqueiaFlowSemChrome9333OuSemRembg(t *testing.T) {
	srvOff := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srvOff.Close()

	dir := criarFixtureCiclo1(t, "hibrido", "11111111-2222-3333-4444-555555555555")
	rel, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{
		FlowCDPURL:   srvOff.URL,
		HTTPClient:   srvOff.Client(),
		ChecarModulo: func(string) bool { return true },
	})
	if !errors.Is(err, ErrPreVooBloqueado) || !rel.Bloqueado {
		t.Fatalf("esperado ErrPreVooBloqueado com CDP fechado, obtido rel=%+v err=%v", rel, err)
	}

	srvOn := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srvOn.Close()

	rel2, err2 := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{
		FlowCDPURL:   srvOn.URL,
		HTTPClient:   srvOn.Client(),
		ChecarModulo: func(nome string) bool { return nome != "rembg" },
	})
	if !errors.Is(err2, ErrPreVooBloqueado) || !rel2.Bloqueado {
		t.Fatalf("esperado ErrPreVooBloqueado sem rembg no modo hibrido, obtido rel=%+v err=%v", rel2, err2)
	}
}

func TestPreVooBloqueiaSemFlowProjeto(t *testing.T) {
	for _, modo := range []string{"flow", "hibrido"} {
		t.Run(modo, func(t *testing.T) {
			dir := criarFixtureCiclo1(t, modo, "")
			rel, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{})
			if !errors.Is(err, ErrFlowTravaAusente) {
				t.Fatalf("esperado ErrFlowTravaAusente, obtido: %v", err)
			}
			if !rel.Bloqueado {
				t.Fatalf("esperado Bloqueado=true")
			}
		})
	}
}

func TestPreVooBloqueiaContratoAusenteOuInvalido(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	if err := os.Remove(filepath.Join(dir, "script.json")); err != nil {
		t.Fatal(err)
	}
	rel, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{})
	if !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("esperado ErrPreVooBloqueado sem script.json, obtido: %v", err)
	}
	if !rel.Bloqueado {
		t.Fatalf("esperado Bloqueado=true")
	}

	dirModo := criarFixtureCiclo1(t, "desconhecido", "")
	rel2, err2 := ExecutarPreVoo(context.Background(), dirModo, DepsPreVoo{})
	if !errors.Is(err2, ErrPreVooBloqueado) || !rel2.Bloqueado {
		t.Fatalf("esperado bloqueio para estilo_visual invalido, obtido rel=%+v err=%v", rel2, err2)
	}

	dirSemPlano := criarFixtureCiclo1(t, "codigo", "")
	if err := os.Remove(filepath.Join(dirSemPlano, "plano.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := ExecutarPreVoo(context.Background(), dirSemPlano, DepsPreVoo{}); !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("esperado bloqueio sem plano.md, obtido: %v", err)
	}
}

func TestPreVooBloqueiaPathJailEscape(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	if _, err := ResolverCaminhoSeguro(dir, "../fora.json"); err == nil {
		t.Fatalf("esperado erro de Path Jail para travessia")
	}
	if _, err := ResolverCaminhoSeguro(dir, "flow/elenco.json"); err != nil {
		t.Fatalf("caminho interno deveria ser permitido: %v", err)
	}
}

// converterProjetoV2 reescreve projeto.json no contrato v2 (modo_producao + estilo).
func converterProjetoV2(t *testing.T, dir, modo, estilo string) {
	t.Helper()
	caminho := filepath.Join(dir, "projeto.json")
	raw, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	delete(m, "estilo_visual")
	m["modo_producao"] = modo
	if estilo != "" {
		m["estilo"] = estilo
	}
	gravarJSONTeste(t, caminho, m)
}

func criarInventarioEstilo(t *testing.T, id, status string) string {
	t.Helper()
	inv := t.TempDir()
	pasta := filepath.Join(inv, "estilos", id)
	if err := os.MkdirAll(pasta, 0o755); err != nil {
		t.Fatal(err)
	}
	gravarJSONTeste(t, filepath.Join(pasta, "estilo.json"), map[string]any{"id": id, "status": status, "versao": "1.0.0"})
	return inv
}

func itemPreVoo(rel RelatorioPreVoo, nome string) (ItemPreVoo, bool) {
	for _, it := range rel.Itens {
		if it.Nome == nome {
			return it, true
		}
	}
	return ItemPreVoo{}, false
}

func TestPreVooV2EstiloValidadoAprova(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	converterProjetoV2(t, dir, "codigo", "molde-de-pilula")
	inv := criarInventarioEstilo(t, "molde-de-pilula", "validado")
	rel, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{InventarioDir: inv})
	if err != nil || rel.Bloqueado {
		t.Fatalf("v2 com estilo validado deveria aprovar: rel=%+v err=%v", rel, err)
	}
	if it, ok := itemPreVoo(rel, "estilo"); !ok || it.Nivel != "OK" {
		t.Fatalf("item estilo OK ausente: %+v", rel.Itens)
	}
	if _, ok := itemPreVoo(rel, AvisoEstiloVisualLegado); ok {
		t.Fatalf("v2 nao deve emitir aviso legado: %+v", rel.Itens)
	}
	contratos, err := CarregarContratosEntrada(dir)
	if err != nil || contratos.Producao.ModoProducao != "codigo" || contratos.Producao.Estilo != "molde-de-pilula" {
		t.Fatalf("contrato v2 mal lido: %+v err=%v", contratos.Producao, err)
	}
}

func TestPreVooV2EstiloRascunhoBloqueia(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	converterProjetoV2(t, dir, "codigo", "palito")
	inv := criarInventarioEstilo(t, "palito", "rascunho")
	rel, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{InventarioDir: inv})
	if !errors.Is(err, ErrEstiloNaoValidado) || !errors.Is(err, ErrPreVooBloqueado) || !rel.Bloqueado {
		t.Fatalf("esperado ErrEstiloNaoValidado, obtido rel=%+v err=%v", rel, err)
	}
}

func TestPreVooV2EstiloInexistenteBloqueia(t *testing.T) {
	inv := criarInventarioEstilo(t, "molde-de-pilula", "validado")
	for _, estilo := range []string{"nao-existe", "", "../fuga"} {
		dir := criarFixtureCiclo1(t, "codigo", "")
		converterProjetoV2(t, dir, "codigo", estilo)
		rel, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{InventarioDir: inv})
		if !errors.Is(err, ErrEstiloInexistente) || !rel.Bloqueado {
			t.Fatalf("estilo %q: esperado ErrEstiloInexistente, obtido rel=%+v err=%v", estilo, rel, err)
		}
	}
}

func TestPreVooV2ModoProducaoInvalidoBloqueia(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	converterProjetoV2(t, dir, "cinema", "molde-de-pilula")
	inv := criarInventarioEstilo(t, "molde-de-pilula", "validado")
	if _, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{InventarioDir: inv}); !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("modo_producao invalido deveria bloquear, obtido: %v", err)
	}
}

func TestPreVooLegadoEstiloVisualAprovaComAviso(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	rel, err := ExecutarPreVoo(context.Background(), dir, DepsPreVoo{InventarioDir: t.TempDir()})
	if err != nil || rel.Bloqueado {
		t.Fatalf("estilo_visual legado deveria aprovar: rel=%+v err=%v", rel, err)
	}
	it, ok := itemPreVoo(rel, AvisoEstiloVisualLegado)
	if !ok || it.Nivel != "AVISO" {
		t.Fatalf("esperado AVISO_ESTILO_VISUAL_LEGADO, itens=%+v", rel.Itens)
	}
}

func TestPreVooLocalizaInventarioSubindoPastas(t *testing.T) {
	raiz := t.TempDir()
	pasta := filepath.Join(raiz, "inventario", "video", "estilos", "molde-de-pilula")
	if err := os.MkdirAll(pasta, 0o755); err != nil {
		t.Fatal(err)
	}
	gravarJSONTeste(t, filepath.Join(pasta, "estilo.json"), map[string]any{"id": "molde-de-pilula", "status": "validado"})
	proj := filepath.Join(raiz, "projetos", "canal", "video")
	if got := localizarInventario(proj, ""); got != filepath.Join(raiz, "inventario", "video") {
		t.Fatalf("inventario nao localizado subindo pastas: %q", got)
	}
	if got := localizarInventario(proj, "explicito"); got != "explicito" {
		t.Fatalf("InventarioDir explicito deve prevalecer: %q", got)
	}
}
