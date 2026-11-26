package ativos

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
)

func TestProcessoExecutarMecanicoCompleto(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	proc := NovoProcesso(DepsProcesso{})
	res, err := proc.Executar(context.Background(), OpcoesExecucao{
		ProjetoDir: dir,
	})
	if err != nil {
		t.Fatalf("Executar falhou: %v", err)
	}
	if res.AguardandoGate != "" {
		t.Fatalf("modo mecanico nao deveria pausar em gate, pausou em %q", res.AguardandoGate)
	}
	if len(res.Sprites) == 0 || len(res.Cenarios) < 2 || len(res.Storyboard.Cenas) != 4 {
		t.Fatalf("resultado incompleto: sprites=%d cenarios=%d cenas=%d", len(res.Sprites), len(res.Cenarios), len(res.Storyboard.Cenas))
	}
}

func TestProcessoPoliticaGatesHumanoEHibrido(t *testing.T) {
	t.Run("modo humano para no primeiro stop-gate e avanca com aprovacao", func(t *testing.T) {
		dir := criarFixtureCiclo1(t, "codigo", "")
		atualizarPoliticaGatesTeste(t, dir, roteiro.PoliticaGates{Modo: "humano"})
		proc := NovoProcesso(DepsProcesso{})

		res1, err := proc.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if res1.AguardandoGate != GateEtapa3Elenco {
			t.Fatalf("esperado pausa em %q, obtido %q", GateEtapa3Elenco, res1.AguardandoGate)
		}

		res2, err := proc.Executar(context.Background(), OpcoesExecucao{
			ProjetoDir:     dir,
			GatesAprovados: []string{GateEtapa3Elenco},
		})
		if err != nil || res2.AguardandoGate != GateEtapa4QuadrosTeste {
			t.Fatalf("esperado pausa em %q, obtido %q err=%v", GateEtapa4QuadrosTeste, res2.AguardandoGate, err)
		}

		res3, err := proc.Executar(context.Background(), OpcoesExecucao{
			ProjetoDir:     dir,
			GatesAprovados: []string{GateEtapa3Elenco, GateEtapa4QuadrosTeste},
		})
		if err != nil || res3.AguardandoGate != GateEtapa5Storyboard {
			t.Fatalf("esperado pausa em %q, obtido %q err=%v", GateEtapa5Storyboard, res3.AguardandoGate, err)
		}

		res4, err := proc.Executar(context.Background(), OpcoesExecucao{
			ProjetoDir:     dir,
			GatesAprovados: []string{GateEtapa3Elenco, GateEtapa4QuadrosTeste, GateEtapa5Storyboard},
		})
		if err != nil || res4.AguardandoGate != "" {
			t.Fatalf("esperado conclusao sem pausa, obtido gate=%q err=%v", res4.AguardandoGate, err)
		}
	})

	t.Run("modo hibrido para apenas nos gates configurados", func(t *testing.T) {
		dir := criarFixtureCiclo1(t, "codigo", "")
		atualizarPoliticaGatesTeste(t, dir, roteiro.PoliticaGates{
			Modo:         "hibrido",
			GatesHumanos: []string{GateEtapa5Storyboard},
		})
		proc := NovoProcesso(DepsProcesso{})
		res, err := proc.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir})
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if res.AguardandoGate != GateEtapa5Storyboard {
			t.Fatalf("esperado pausa apenas em %q, obtido %q", GateEtapa5Storyboard, res.AguardandoGate)
		}
	})
}

func TestProcessoCoberturaBordas(t *testing.T) {
	t.Run("pre-voo falho interrompe antes da etapa 3", func(t *testing.T) {
		dir := criarFixtureCiclo1(t, "flow", "")
		proc := NovoProcesso(DepsProcesso{})
		if _, err := proc.Executar(context.Background(), OpcoesExecucao{ProjetoDir: dir}); !errors.Is(err, ErrFlowTravaAusente) {
			t.Fatalf("esperado ErrFlowTravaAusente, obtido: %v", err)
		}
	})

	t.Run("validacoes extras de contratos e bordas", func(t *testing.T) {
		if (BoundingBox{W: 0, H: 10}).Sobrepoe(BoundingBox{W: 10, H: 10}) {
			t.Fatalf("caixa vazia nao sobrepoe")
		}
		if err := ValidarSpritesSemSobreposicao(nil); !errors.Is(err, ErrSpriteInvalido) {
			t.Fatalf("lista vazia deveria falhar")
		}
		if minimoCenarios(0) != 1 {
			t.Fatalf("minimoCenarios(0) deve ser 1")
		}
		if primeiroSpriteID(nil) != "narrador_pose_1" {
			t.Fatalf("fallback de primeiroSpriteID inesperado")
		}
		for _, nat := range []string{"fato_consolidado", "evidencia_preliminar", "alegacao_terceiro", "narrativa_transicao"} {
			if derivarMetaforaVisual(roteiro.LinhaRoteiro{Natureza: nat}) == "" {
				t.Fatalf("metafora vazia para %s", nat)
			}
		}
		if err := ValidarStoryboard(Storyboard{
			Cenas: []Cena{{IndiceLinha: 2, CenarioID: "c1", MetaforaVisual: "m", ElementosMoveis: make([]ElementoMovel, 4)}},
		}, roteiro.Script{Lines: []roteiro.LinhaRoteiro{{Text: "a"}}}, []Cenario{{ID: "c1"}}); !errors.Is(err, ErrStoryboardInvalido) {
			t.Fatalf("esperado erro para indice fora de ordem")
		}
	})

	t.Run("erro JSON invalido em projeto, marca e script", func(t *testing.T) {
		dir := t.TempDir()
		_ = os.WriteFile(filepath.Join(dir, "projeto.json"), []byte("{invalido"), 0o644)
		if _, err := CarregarContratosEntrada(dir); !errors.Is(err, ErrPreVooBloqueado) {
			t.Fatalf("esperado ErrPreVooBloqueado para projeto.json invalido")
		}
		dir2 := criarFixtureCiclo1(t, "codigo", "")
		_ = os.WriteFile(filepath.Join(dir2, "marca.json"), []byte("{invalido"), 0o644)
		if _, err := CarregarContratosEntrada(dir2); !errors.Is(err, ErrPreVooBloqueado) {
			t.Fatalf("esperado ErrPreVooBloqueado para marca.json invalido")
		}
		dir3 := criarFixtureCiclo1(t, "codigo", "")
		gravarJSONTeste(t, filepath.Join(dir3, "script.json"), roteiro.Script{SchemaVersion: 1, Title: "Vazio"})
		if _, err := CarregarContratosEntrada(dir3); !errors.Is(err, ErrPreVooBloqueado) {
			t.Fatalf("esperado ErrPreVooBloqueado para script sem linhas")
		}
	})

	t.Run("flow sem cliente injetado falha", func(t *testing.T) {
		dir := criarFixtureCiclo1(t, "flow", "uuid-ok")
		contratos, _ := CarregarContratosEntrada(dir)
		if _, err := ExecutarEtapa4Cenarios(context.Background(), dir, contratos.Projeto, contratos.Script, DepsCenarios{}); !errors.Is(err, ErrPreVooBloqueado) {
			t.Fatalf("esperado erro sem cliente Flow: %v", err)
		}
		if _, err := ExecutarEtapa4Cenarios(context.Background(), dir, roteiro.Projeto{EstiloVisual: "invalido"}, contratos.Script, DepsCenarios{}); !errors.Is(err, ErrPreVooBloqueado) {
			t.Fatalf("esperado erro com modo invalido: %v", err)
		}
	})
}

func atualizarPoliticaGatesTeste(t *testing.T, dir string, pg roteiro.PoliticaGates) {
	t.Helper()
	contratos, err := CarregarContratosEntrada(dir)
	if err != nil {
		t.Fatal(err)
	}
	contratos.Projeto.PoliticaGates = pg
	gravarJSONTeste(t, filepath.Join(dir, "projeto.json"), contratos.Projeto)
}
