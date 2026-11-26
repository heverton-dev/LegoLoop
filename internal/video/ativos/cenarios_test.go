package ativos

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

type clienteFlowFake struct {
	falhasAntesSucesso int
	chamadas           int
	ultimoReq          RequisicaoFlow
}

func (c *clienteFlowFake) GerarCenario(_ context.Context, req RequisicaoFlow) ([]byte, error) {
	c.chamadas++
	c.ultimoReq = req
	if c.chamadas <= c.falhasAntesSucesso {
		return nil, errors.New("erro transitorio flow")
	}
	return []byte("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 1280 720\"><rect width=\"1280\" height=\"720\" fill=\"#1e293b\"/></svg>"), nil
}

func TestCenariosV2UsaModoProducao(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	converterProjetoV2(t, dir, "codigo", "molde-de-pilula")
	contratos, err := CarregarContratosEntrada(dir)
	if err != nil {
		t.Fatal(err)
	}
	if ModoProducaoProjeto(contratos.Projeto) != "codigo" {
		t.Fatalf("modo_producao v2 nao chegou ao projeto: %+v", contratos.Projeto)
	}
	flowFake := &clienteFlowFake{}
	if _, err := ExecutarEtapa4Cenarios(context.Background(), dir, contratos.Projeto, contratos.Script, DepsCenarios{Flow: flowFake}); err != nil {
		t.Fatalf("Etapa 4 com contrato v2 falhou: %v", err)
	}
	if flowFake.chamadas != 0 {
		t.Fatalf("modo codigo v2 nao deve chamar Flow, chamadas=%d", flowFake.chamadas)
	}
}

func TestCenariosModoCodigoSVG100Autoral(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	contratos, err := CarregarContratosEntrada(dir)
	if err != nil {
		t.Fatal(err)
	}
	flowFake := &clienteFlowFake{}
	cenarios, err := ExecutarEtapa4Cenarios(context.Background(), dir, contratos.Projeto, contratos.Script, DepsCenarios{
		Flow: flowFake,
	})
	if err != nil {
		t.Fatalf("Etapa 4 modo codigo falhou: %v", err)
	}
	if flowFake.chamadas != 0 {
		t.Fatalf("modo codigo nao deve chamar ClienteFlow, chamadas=%d", flowFake.chamadas)
	}
	if len(cenarios) < 2 {
		t.Fatalf("4 linhas exigem >= 2 cenarios distintos, obtido %d", len(cenarios))
	}
	raw, err := os.ReadFile(filepath.Join(dir, cenarios[0].Arquivo))
	if err != nil || !strings.Contains(string(raw), "<svg") {
		t.Fatalf("SVG autoral invalido: err=%v", err)
	}
}

func TestCenariosModoFlowBackoffRefTravaEIdempotencia(t *testing.T) {
	dir := criarFixtureCiclo1(t, "flow", "11111111-2222-3333-4444-555555555555")
	contratos := carregarContratosTeste(t, dir)
	var esperas []time.Duration
	flowFake := &clienteFlowFake{falhasAntesSucesso: 2}
	deps := DepsCenarios{
		Flow:        flowFake,
		Referencias: []string{"flow/narrador_turnaround.png"},
		Dormir:      func(d time.Duration) { esperas = append(esperas, d) },
		MaxRetries:  3,
	}
	cenarios, err := ExecutarEtapa4Cenarios(context.Background(), dir, contratos.Projeto, contratos.Script, deps)
	if err != nil {
		t.Fatalf("Etapa 4 modo flow falhou: %v", err)
	}
	if len(esperas) != 2 {
		t.Fatalf("esperado 2 backoffs antes do sucesso, obtido %d", len(esperas))
	}
	if flowFake.ultimoReq.FlowProjeto != "11111111-2222-3333-4444-555555555555" || len(flowFake.ultimoReq.Referencias) != 1 {
		t.Fatalf("trava flow_projeto ou referencia --ref nao enviada: %+v", flowFake.ultimoReq)
	}
	chamadasAntes := flowFake.chamadas

	// Segunda execucao com mesma entrada deve ser 100% idempotente por hash.
	cenarios2, err := ExecutarEtapa4Cenarios(context.Background(), dir, contratos.Projeto, contratos.Script, deps)
	if err != nil {
		t.Fatalf("segunda execucao falhou: %v", err)
	}
	if flowFake.chamadas != chamadasAntes {
		t.Fatalf("idempotencia violada: chamadas antes=%d depois=%d", chamadasAntes, flowFake.chamadas)
	}
	if len(cenarios2) != len(cenarios) {
		t.Fatalf("quantidade de cenarios divergiu")
	}
}

// inventarioCenarioFixture cria estilo validado com um cenario aprovado de tags banco-central e noite.
func inventarioCenarioFixture(t *testing.T) (inventario.Store, inventario.Item) {
	t.Helper()
	s := inventario.Store{Raiz: t.TempDir()}
	base := filepath.Join(s.Raiz, "inventario", "video", "estilos", "molde-pilula")
	if err := os.MkdirAll(filepath.Join(base, "arquivos"), 0o755); err != nil {
		t.Fatal(err)
	}
	gravarJSONTeste(t, filepath.Join(base, "estilo.json"), inventario.Estilo{ID: "molde-pilula", Status: "validado", Versao: "1.0.0", PromptMestre: "estilo pilula"})
	rel := filepath.ToSlash(filepath.Join("inventario", "video", "estilos", "molde-pilula", "arquivos", "banco-noite.svg"))
	if err := os.WriteFile(filepath.Join(s.Raiz, rel), []byte("<svg>banco</svg>"), 0o644); err != nil {
		t.Fatal(err)
	}
	h, err := inventario.HashArquivo(filepath.Join(s.Raiz, rel))
	if err != nil {
		t.Fatal(err)
	}
	it := inventario.Item{SchemaVersion: 1, ID: "banco-noite", Tipo: "cenario", Estilo: "molde-pilula", Tags: []string{"banco-central", "noite"}, Arquivo: rel, Hash: h, Aprovado: true}
	if err := s.GravarItem(it); err != nil {
		t.Fatal(err)
	}
	return s, it
}

func TestCenariosInventarioReutilizaAprovado(t *testing.T) {
	s, it := inventarioCenarioFixture(t)
	proj := t.TempDir()
	flowFake := &clienteFlowFake{}
	ped := PedidoFundo{ID: "fundo-banco", Tags: []string{"banco-central", "noite"}, Prompt: "fachada do banco central"}
	res, err := ObterFundoInventario(context.Background(), proj, s, "molde-pilula", ped, DepsCenarios{Flow: flowFake}, "proj-x 2026-11-26")
	if err != nil {
		t.Fatalf("reuso do inventario falhou: %v", err)
	}
	if flowFake.chamadas != 0 || !res.Reutilizado || res.Item.ID != it.ID {
		t.Fatalf("esperado reuso sem Flow: chamadas=%d res=%+v", flowFake.chamadas, res)
	}
	h, err := inventario.HashArquivo(filepath.Join(proj, res.Arquivo))
	if err != nil || h != it.Hash {
		t.Fatalf("copia no projeto com hash divergente: err=%v hash=%s", err, h)
	}
}

func TestCenariosInventarioHashDivergenteBloqueia(t *testing.T) {
	s, it := inventarioCenarioFixture(t)
	if err := os.WriteFile(filepath.Join(s.Raiz, it.Arquivo), []byte("adulterado"), 0o644); err != nil {
		t.Fatal(err)
	}
	ped := PedidoFundo{ID: "fundo-banco", Tags: []string{"banco-central", "noite"}}
	_, err := ObterFundoInventario(context.Background(), t.TempDir(), s, "molde-pilula", ped, DepsCenarios{Flow: &clienteFlowFake{}}, "proj-x")
	if !errors.Is(err, inventario.ErrHashDivergente) {
		t.Fatalf("esperado ErrHashDivergente, obtido %v", err)
	}
}

func TestCenariosInventarioRegistraFundoNovo(t *testing.T) {
	s, _ := inventarioCenarioFixture(t)
	proj := t.TempDir()
	flowFake := &clienteFlowFake{}
	ped := PedidoFundo{ID: "praca-dia", Tags: []string{"praca", "dia"}, Prompt: "praca central ao meio-dia"}
	res, err := ObterFundoInventario(context.Background(), proj, s, "molde-pilula", ped, DepsCenarios{Flow: flowFake, Dormir: func(time.Duration) {}}, "proj-x 2026-11-26")
	if err != nil {
		t.Fatalf("geracao de fundo novo falhou: %v", err)
	}
	if flowFake.chamadas != 1 || res.Reutilizado {
		t.Fatalf("fundo ausente deve gerar uma vez: chamadas=%d res=%+v", flowFake.chamadas, res)
	}
	if !strings.Contains(flowFake.ultimoReq.Prompt, "estilo pilula") || !strings.Contains(flowFake.ultimoReq.Prompt, ped.Prompt) {
		t.Fatalf("prompt sem prompt_mestre do estilo: %q", flowFake.ultimoReq.Prompt)
	}
	var it inventario.Item
	lerJSONTesteCenario(t, filepath.Join(s.Raiz, "inventario", "video", "estilos", "molde-pilula", "itens", "cenario", "praca-dia.json"), &it)
	if it.Parametros["prompt"] != flowFake.ultimoReq.Prompt || strings.Join(it.Tags, ",") != "praca,dia" || it.Aprovado {
		t.Fatalf("item novo sem prompt, tags ou com aprovacao indevida: %+v", it)
	}
	if err := s.Verificar(it); err != nil || it.Hash == "" {
		t.Fatalf("hash do item novo invalido: %v", err)
	}
	if !arquivoExisteSeguro(proj, res.Arquivo) {
		t.Fatalf("fundo novo nao copiado ao projeto: %q", res.Arquivo)
	}
}

func lerJSONTesteCenario(t *testing.T, caminho string, v any) {
	t.Helper()
	b, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatal(err)
	}
}

func carregarContratosTeste(t *testing.T, dir string) ContratosEntrada {
	t.Helper()
	contratos, err := CarregarContratosEntrada(dir)
	if err != nil {
		t.Fatal(err)
	}
	return contratos
}

func TestCenariosModoHibrido(t *testing.T) {
	t.Run("hibrido gera todos os cenarios no flow", func(t *testing.T) {
		dir := criarFixtureCiclo1(t, "hibrido", "uuid-hibrido-1")
		contratos := carregarContratosTeste(t, dir)
		flowFake := &clienteFlowFake{}
		cenarios, err := ExecutarEtapa4Cenarios(context.Background(), dir, contratos.Projeto, contratos.Script, DepsCenarios{
			Flow:   flowFake,
			Dormir: func(time.Duration) {},
		})
		if err != nil {
			t.Fatalf("hibrido falhou: %v", err)
		}
		if len(cenarios) < 2 || flowFake.chamadas != len(cenarios) {
			t.Fatalf("hibrido deveria gerar todos os cenarios via flow: cenarios=%d chamadas=%d", len(cenarios), flowFake.chamadas)
		}
	})
}

func TestCenariosBloqueios(t *testing.T) {
	t.Run("flow sem trava bloqueia", func(t *testing.T) {
		dir := criarFixtureCiclo1(t, "flow", "")
		contratos := carregarContratosTeste(t, dir)
		_, err := ExecutarEtapa4Cenarios(context.Background(), dir, contratos.Projeto, contratos.Script, DepsCenarios{Flow: &clienteFlowFake{}})
		if !errors.Is(err, ErrFlowTravaAusente) {
			t.Fatalf("esperado ErrFlowTravaAusente, obtido: %v", err)
		}
	})

	t.Run("densidade insuficiente bloqueia", func(t *testing.T) {
		err := ValidarDensidadeCenarios([]Cenario{{ID: "c1"}}, 4)
		if !errors.Is(err, ErrDensidadeCenarios) {
			t.Fatalf("esperado ErrDensidadeCenarios, obtido: %v", err)
		}
	})

	t.Run("falha apos esgotar retries", func(t *testing.T) {
		dir := criarFixtureCiclo1(t, "flow", "uuid-1")
		contratos := carregarContratosTeste(t, dir)
		flowFake := &clienteFlowFake{falhasAntesSucesso: 10}
		_, err := ExecutarEtapa4Cenarios(context.Background(), dir, contratos.Projeto, contratos.Script, DepsCenarios{
			Flow:       flowFake,
			Dormir:     func(time.Duration) {},
			MaxRetries: 2,
		})
		if err == nil {
			t.Fatalf("esperado erro ao esgotar retries")
		}
	})
}
