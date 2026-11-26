package ativos

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
)

func TestFichas4FaixasValidacaoEExtracao(t *testing.T) {
	ficha := FichaAtivo{
		ID:                "narrador",
		Categoria:         "personagem",
		Nome:              "Narrador",
		DescricaoImutavel: "Miniatura terracota com mascara preta e olhos brancos",
		Turnaround:        []string{"Front", "Three-quarter", "Side", "Back three-quarter", "Back"},
		Expressoes:        []string{"Neutral", "Wide", "Narrowed", "Slanted in", "Squeezed"},
		Poses:             []string{"Explaining", "Pointing", "Thinking", "Surprise"},
		BuildNotes: BuildNotes{
			CorpoMaterial:     "Terracotta matte block",
			OlhosDetalhes:     "Two white glowing slot eyes",
			AcessoriosEfeitos: "Black cotton cape",
			Regras:            "No colour change; white eyes in every shot",
		},
	}
	if err := ficha.Validar(); err != nil {
		t.Fatalf("ficha 4 faixas valida falhou: %v", err)
	}

	invalida := ficha
	invalida.BuildNotes.Regras = ""
	if !errors.Is(invalida.Validar(), ErrFichaInvalida) {
		t.Fatalf("esperado ErrFichaInvalida sem build_notes.regras, obtido: %v", invalida.Validar())
	}

	invalidaCat := ficha
	invalidaCat.Categoria = "invalida"
	if !errors.Is(invalidaCat.Validar(), ErrFichaInvalida) {
		t.Fatalf("esperado ErrFichaInvalida para categoria invalida")
	}

	prompt := CompilarPromptFicha4Faixas(ficha)
	for _, want := range []string{"TURNAROUND", "EXPRESSIONS", "POSES", "BUILD NOTES", "No colour change"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt 4 faixas sem secao %q: %s", want, prompt)
		}
	}

	sc := roteiro.Script{
		Lines: []roteiro.LinhaRoteiro{{Text: "L1"}},
		Ativos: []roteiro.AtivoPrevisto{
			{ID: "moedas", Categoria: "objeto", Nome: "Moedas", DescricaoImutavel: "Moedas douradas", BuildNotes: []string{"Brilho metalico", "Sem texto"}},
		},
	}
	plano := "| 1 | 00-05s | transicao | L1 | Diorama do Banco Central com torre de barras | — |"
	extraidas := ExtrairFichasDeEntrada(sc, plano)
	if len(extraidas) < 3 {
		t.Fatalf("esperado personagem padrao + objeto do script + dispositivo do plano, obtido %d", len(extraidas))
	}
}

type clienteFlowContadorFake struct {
	chamadas int
}

func (c *clienteFlowContadorFake) GerarCenario(_ context.Context, _ RequisicaoFlow) ([]byte, error) {
	c.chamadas++
	return []byte("\x89PNG\r\n\x1a\n"), nil
}

func TestFichas4FaixasGeracaoEIdempotencia(t *testing.T) {
	dir := criarFixtureCiclo1(t, "flow", "uuid-fichas-1")
	contratos, _ := CarregarContratosEntrada(dir)
	fichas := ExtrairFichasDeEntrada(contratos.Script, contratos.PlanoMD)
	flowFake := &clienteFlowContadorFake{}
	deps := DepsCenarios{Flow: flowFake, Dormir: func(time.Duration) {}}

	geradas, err := ExecutarGeracaoFichas4Faixas(context.Background(), dir, contratos.Projeto, fichas, deps)
	if err != nil {
		t.Fatalf("ExecutarGeracaoFichas4Faixas falhou: %v", err)
	}
	if len(geradas) == 0 || flowFake.chamadas != len(geradas) {
		t.Fatalf("esperado %d chamadas ao Flow, obtido %d", len(geradas), flowFake.chamadas)
	}
	if _, err := os.Stat(filepath.Join(dir, "flow", "fichas.json")); err != nil {
		t.Fatalf("flow/fichas.json nao gravado: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(geradas[0].FolhaMestrePNG))); err != nil {
		t.Fatalf("folha mestre PNG nao gravada: %v", err)
	}

	chamadasAntes := flowFake.chamadas
	if _, err := ExecutarGeracaoFichas4Faixas(context.Background(), dir, contratos.Projeto, fichas, deps); err != nil {
		t.Fatalf("segunda execucao falhou: %v", err)
	}
	if flowFake.chamadas != chamadasAntes {
		t.Fatalf("idempotencia das fichas violada: antes=%d depois=%d", chamadasAntes, flowFake.chamadas)
	}

	dirCod := criarFixtureCiclo1(t, "codigo", "")
	contratosCod, _ := CarregarContratosEntrada(dirCod)
	geradasCod, err := ExecutarGeracaoFichas4Faixas(context.Background(), dirCod, contratosCod.Projeto, fichas, DepsCenarios{})
	if err != nil || len(geradasCod) == 0 {
		t.Fatalf("modo codigo para fichas 4 faixas falhou: err=%v", err)
	}

	if _, err := ExecutarGeracaoFichas4Faixas(context.Background(), dir, contratos.Projeto, []FichaAtivo{{ID: "x"}}, deps); !errors.Is(err, ErrFichaInvalida) {
		t.Fatalf("esperado ErrFichaInvalida, obtido: %v", err)
	}
	dirSemFlow := criarFixtureCiclo1(t, "flow", "uuid-sem-flow")
	if _, err := ExecutarGeracaoFichas4Faixas(context.Background(), dirSemFlow, contratos.Projeto, fichas, DepsCenarios{}); !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("esperado ErrPreVooBloqueado sem Flow, obtido: %v", err)
	}
	if _, err := ExecutarGeracaoFichas4Faixas(context.Background(), dirSemFlow, contratos.Projeto, fichas, DepsCenarios{
		Flow: &clienteFlowFake{falhasAntesSucesso: 5}, Dormir: func(time.Duration) {}, MaxRetries: 1,
	}); err == nil {
		t.Fatalf("esperado erro ao esgotar retries em ficha")
	}
	if _, err := ExecutarGeracaoFichas4Faixas(context.Background(), "\x00invalido", contratosCod.Projeto, fichas, DepsCenarios{}); err == nil {
		t.Fatalf("esperado erro para projDir invalido em modo codigo")
	}
	conv := converterAtivoPrevisto(roteiro.AtivoPrevisto{ID: "m1", Categoria: "desconhecida", Nome: "Efeito", DescricaoImutavel: "Zoom"})
	if conv.Categoria != "objeto" || conv.BuildNotes.Regras == "" {
		t.Fatalf("converterAtivoPrevisto fallback inesperado: %+v", conv)
	}
	if err := (FichaAtivo{ID: "a", Categoria: "personagem", Nome: "A", DescricaoImutavel: "D"}).Validar(); !errors.Is(err, ErrFichaInvalida) {
		t.Fatalf("esperado ErrFichaInvalida para faixas < 3")
	}
}
