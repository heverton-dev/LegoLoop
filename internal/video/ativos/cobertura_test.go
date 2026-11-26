package ativos

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

func TestCoberturaErrosDeBordaEPathJail(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")

	t.Run("elenco customizado e falha propagada nas etapas 3 e 4", func(t *testing.T) {
		custom := []FichaPersonagem{
			{
				ID:                "host",
				Nome:              "Host",
				DescricaoImutavel: "Palito vermelho",
				Fichas: map[string]string{
					"turnaround":     "flow/h_t.png",
					"expressoes_3x2": "flow/h_e.png",
					"poses_6x":       "flow/h_p.png",
				},
			},
		}
		procOK := NovoProcesso(DepsProcesso{})
		res, err := procOK.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir, Elenco: custom})
		if err != nil || len(res.Sprites) == 0 {
			t.Fatalf("elenco customizado falhou: err=%v", err)
		}

		errE3 := errors.New("falha etapa 3")
		procE3 := NovoProcesso(DepsProcesso{Recortador: &recortadorFake{falhar: errE3}})
		if _, err := procE3.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir}); !errors.Is(err, errE3) {
			t.Fatalf("esperado erro da etapa 3, obtido: %v", err)
		}

		dirFlow := criarFixtureCiclo1(t, "flow", "uuid-x")
		procE4 := NovoProcesso(DepsProcesso{
			Cenarios: DepsCenarios{Flow: &clienteFlowFake{falhasAntesSucesso: 5}, MaxRetries: 1},
		})
		if _, err := procE4.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dirFlow}); err == nil {
			t.Fatalf("esperado erro da etapa 4 no processo")
		}
	})

	t.Run("subpasta bloqueada por arquivo comum reprova pre-voo", func(t *testing.T) {
		dirBloq := criarFixtureCiclo1(t, "codigo", "")
		if err := os.WriteFile(filepath.Join(dirBloq, "flow"), []byte("nao sou pasta"), 0o644); err != nil {
			t.Fatal(err)
		}
		rel, err := ExecutarPreVoo(context.Background(), dirBloq, DepsPreVoo{})
		if !errors.Is(err, ErrPreVooBloqueado) || !rel.Bloqueado {
			t.Fatalf("esperado bloqueio de escrita em flow, obtido rel=%+v err=%v", rel, err)
		}
	})

	t.Run("helpers seguros rejeitam escape de Path Jail", func(t *testing.T) {
		if err := gravarJSONSeguro(dir, "../escape.json", map[string]int{"a": 1}); err == nil {
			t.Fatalf("gravarJSONSeguro deveria rejeitar escape")
		}
		if err := gravarBytesSeguro(dir, "../escape.svg", []byte("<svg/>")); err == nil {
			t.Fatalf("gravarBytesSeguro deveria rejeitar escape")
		}
		var dst map[string]any
		if err := lerJSONSeguro(dir, "../escape.json", &dst); err == nil {
			t.Fatalf("lerJSONSeguro deveria rejeitar escape")
		}
		if arquivoExisteSeguro(dir, "../escape.svg") {
			t.Fatalf("arquivoExisteSeguro deveria rejeitar escape")
		}
	})

	t.Run("textoBaseCenario com indice acima do limite e dispositivo do plano", func(t *testing.T) {
		sc := roteiro.Script{Lines: []roteiro.LinhaRoteiro{{Text: "unica"}}}
		if got := textoBaseCenario(sc, nil, 5); got != "unica" {
			t.Fatalf("esperado 'unica', obtido %q", got)
		}
		if got := textoBaseCenario(sc, []string{"Diorama banco central"}, 0); got != "Diorama banco central" {
			t.Fatalf("esperado dispositivo do plano, obtido %q", got)
		}
	})

	t.Run("lerJSONSeguro tolera BOM UTF-8 do Windows", func(t *testing.T) {
		dirBOM := criarFixtureCiclo1(t, "codigo", "")
		raw, err := os.ReadFile(filepath.Join(dirBOM, "projeto.json"))
		if err != nil {
			t.Fatal(err)
		}
		comBOM := append([]byte{0xEF, 0xBB, 0xBF}, raw...)
		if err := os.WriteFile(filepath.Join(dirBOM, "projeto.json"), comBOM, 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := CarregarContratosEntrada(dirBOM); err != nil {
			t.Fatalf("deveria aceitar projeto.json com BOM UTF-8: %v", err)
		}
	})

	t.Run("NovoClienteFlowCDP rejeita execucao sem flow_projeto", func(t *testing.T) {
		cli := NovoClienteFlowCDP("", dir)
		if _, err := cli.GerarCenario(context.Background(), RequisicaoFlow{CenarioID: "cen_01"}); !errors.Is(err, ErrFlowTravaAusente) {
			t.Fatalf("esperado ErrFlowTravaAusente: %v", err)
		}
		if NovoRecortadorRembgDespill() == nil {
			t.Fatalf("NovoRecortadorRembgDespill nao deve retornar nil")
		}
	})
}
