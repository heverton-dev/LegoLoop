package ativos

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
