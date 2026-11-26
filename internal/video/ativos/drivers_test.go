package ativos

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDriversScriptsEmbutidosDisponiveis(t *testing.T) {
	if len(scriptFlowDriver) == 0 || !strings.Contains(string(scriptFlowDriver), "connect_over_cdp") {
		t.Fatalf("script flow_driver.py embutido invalido")
	}
	if len(scriptSpriteDespill) == 0 || !strings.Contains(string(scriptSpriteDespill), "rembg") {
		t.Fatalf("script sprite_despill.py embutido invalido")
	}
	if !strings.Contains(string(scriptFlowDriver), "--baixar-todas") || !strings.Contains(string(scriptSpriteDespill), "--layout") {
		t.Fatalf("scripts devem suportar --baixar-todas (xN) e --layout (3x2/6x)")
	}
}

func TestDriversClienteFlowCDPComExecutorInjetado(t *testing.T) {
	dir := t.TempDir()
	var cmdCapturado string
	var argsCapturados []string

	cli := &ClienteFlowCDP{
		CDPURL:     "http://127.0.0.1:9333",
		ProjetoDir: dir,
		ExecCmd: func(_ context.Context, name string, args ...string) ([]byte, error) {
			cmdCapturado = name
			argsCapturados = args
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "--saida-arquivo" {
					if err := os.WriteFile(args[i+1], []byte("PNG_REAL_SIMULADO"), 0o644); err != nil {
						return nil, err
					}
				}
			}
			return []byte(`{"ok":true}`), nil
		},
	}

	out, err := cli.GerarCenario(context.Background(), RequisicaoFlow{
		CenarioID:   "cen_01",
		Prompt:      "Escritorio isometrico com grafico",
		FlowProjeto: "11111111-2222-3333-4444-555555555555",
		Referencias: []string{"narrador_turnaround"},
	})
	if err != nil {
		t.Fatalf("GerarCenario falhou: %v", err)
	}
	if string(out) != "PNG_REAL_SIMULADO" || cmdCapturado != "python" || len(argsCapturados) == 0 {
		t.Fatalf("execucao inesperada: cmd=%s out=%q", cmdCapturado, string(out))
	}
	if _, err := cli.GerarCenario(context.Background(), RequisicaoFlow{CenarioID: "cen_01"}); !errors.Is(err, ErrFlowTravaAusente) {
		t.Fatalf("esperado ErrFlowTravaAusente sem flow_projeto, obtido: %v", err)
	}
}

func TestDriversRecortadorRembgDespillComExecutorInjetado(t *testing.T) {
	dir := t.TempDir()
	rec := &RecortadorRembgDespill{
		ExecCmd: func(_ context.Context, _ string, args ...string) ([]byte, error) {
			sprites := []DescritorSprite{
				{
					ID:           "p1_pose_1",
					PersonagemID: "p1",
					Tipo:         "pose",
					ArquivoPNG:   "sprites/p1_pose_1.png",
					Largura:      256,
					Altura:       256,
					AlfaValido:   true,
					DespillVerde: true,
					BoundingBox:  BoundingBox{X: 0, Y: 0, W: 128, H: 128},
				},
			}
			return json.Marshal(sprites)
		},
	}
	ficha := FichaPersonagem{
		ID:                "p1",
		Nome:              "P1",
		DescricaoImutavel: "Palito",
		Fichas: map[string]string{
			"turnaround":     "flow/p1_t.png",
			"expressoes_3x2": "flow/p1_e.png",
			"poses_6x":       "flow/p1_p.png",
		},
	}
	res, err := rec.RecortarEFatiar(context.Background(), dir, ficha)
	if err != nil || len(res) != 1 || res[0].ID != "p1_pose_1" {
		t.Fatalf("RecortarEFatiar falhou: res=%+v err=%v", res, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "flow", ".legoloop_sprite_despill.py")); err != nil {
		t.Fatalf("script sprite_despill.py nao materializado no projDir: %v", err)
	}
}

func TestDriversRecortadorDoEstilo(t *testing.T) {
	rec, err := RecortadorDoEstilo("recorte-rembg")
	if err != nil {
		t.Fatalf("estilo de recorte deve usar rembg: %v", err)
	}
	if _, ok := rec.(*RecortadorRembgDespill); !ok {
		t.Fatalf("esperado *RecortadorRembgDespill, obtido %T", rec)
	}
	if _, err := RecortadorDoEstilo("molde-pilula"); !errors.Is(err, ErrRenderizadorNaoRecorte) {
		t.Fatalf("esperado ErrRenderizadorNaoRecorte para vetorial, obtido: %v", err)
	}
}

func TestDriversErrosEBordas(t *testing.T) {
	dir := t.TempDir()

	t.Run("falhas de execucao e Path Jail no ClienteFlowCDP", func(t *testing.T) {
		cli := &ClienteFlowCDP{
			CDPURL:     "http://127.0.0.1:9333",
			ProjetoDir: dir,
			ExecCmd: func(_ context.Context, _ string, _ ...string) ([]byte, error) {
				return []byte("boom"), errors.New("exit 1")
			},
		}
		if _, err := cli.GerarCenario(context.Background(), RequisicaoFlow{
			CenarioID:   "cen_01",
			FlowProjeto: "uuid-1",
		}); err == nil {
			t.Fatalf("esperado erro quando ExecCmd falha")
		}
		if _, err := cli.GerarCenario(context.Background(), RequisicaoFlow{
			CenarioID:   "../../fora",
			FlowProjeto: "uuid-1",
		}); err == nil {
			t.Fatalf("esperado erro de Path Jail no CenarioID")
		}
		if _, err := materializarScript(dir, "../../x.py", []byte("ok")); err == nil {
			t.Fatalf("esperado erro de Path Jail em materializarScript")
		}
	})

	t.Run("falhas de execucao e JSON invalido no RecortadorRembgDespill", func(t *testing.T) {
		ficha := FichaPersonagem{
			ID:                "p1",
			Nome:              "P1",
			DescricaoImutavel: "Palito",
			Fichas:            map[string]string{"turnaround": "a", "expressoes_3x2": "b", "poses_6x": "flow/p1_p.png"},
		}
		recErr := &RecortadorRembgDespill{
			ExecCmd: func(_ context.Context, _ string, _ ...string) ([]byte, error) {
				return []byte("erro python"), errors.New("exit 1")
			},
		}
		if _, err := recErr.RecortarEFatiar(context.Background(), dir, ficha); err == nil {
			t.Fatalf("esperado erro quando python falha")
		}

		recJSONRuim := &RecortadorRembgDespill{
			ExecCmd: func(_ context.Context, _ string, _ ...string) ([]byte, error) {
				return []byte("{json_invalido"), nil
			},
		}
		if _, err := recJSONRuim.RecortarEFatiar(context.Background(), dir, ficha); err == nil {
			t.Fatalf("esperado erro para JSON invalido")
		}

		fichaEscape := ficha
		fichaEscape.Fichas = map[string]string{"poses_6x": "../../escape.png"}
		if _, err := recJSONRuim.RecortarEFatiar(context.Background(), dir, fichaEscape); err == nil {
			t.Fatalf("esperado erro de Path Jail para entrada")
		}
	})

	t.Run("resolverExecCmd padrao e extensao png em flow", func(t *testing.T) {
		fn := resolverExecCmd(nil)
		if out, err := fn(context.Background(), "python", "--version"); err != nil || len(out) == 0 {
			t.Fatalf("execucao padrao python --version falhou: err=%v out=%q", err, string(out))
		}
		if err := gerarSpritePNGLocal(dir, "../fora.png", "pose"); err == nil {
			t.Fatalf("gerarSpritePNGLocal deveria bloquear Path Jail")
		}
		if _, err := CarregarContratosEntrada("\x00invalido"); err == nil {
			t.Fatalf("esperado erro para projDir invalido")
		}
		var relTmp RelatorioPreVoo
		if err := checarEscritaSubpastas(&relTmp, "\x00invalido"); err == nil {
			t.Fatalf("esperado erro para projDir invalido em checarEscritaSubpastas")
		}
		if _, err := (&ClienteFlowCDP{ProjetoDir: dir}).GerarCenario(context.Background(), RequisicaoFlow{CenarioID: "c1", FlowProjeto: "u1"}); err == nil {
			t.Fatalf("esperado erro quando script python real nao tem Chrome 9333")
		}
		dirArqFlow := t.TempDir()
		if err := os.WriteFile(filepath.Join(dirArqFlow, "flow"), []byte("arquivo"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := materializarScript(dirArqFlow, "x.py", []byte("ok")); err == nil {
			t.Fatalf("esperado erro quando flow nao e diretorio")
		}
		if _, err := (&RecortadorRembgDespill{}).RecortarEFatiar(context.Background(), dirArqFlow, FichaPersonagem{ID: "p1", Fichas: map[string]string{"poses_6x": "a.png"}}); err == nil {
			t.Fatalf("esperado erro quando materializarScript falha no RecortadorRembgDespill")
		}
		var ct ContratosEntrada
		if err := carregarPlanoEValidarLinhas("\x00invalido", &ct); err == nil {
			t.Fatalf("esperado erro para projDir invalido em carregarPlanoEValidarLinhas")
		}
		if _, err := ExecutarPreVoo(context.Background(), criarFixtureCiclo1(t, "codigo", ""), DepsPreVoo{
			ChecarOpcional: func(string) bool { return true },
		}); err != nil {
			t.Fatalf("prevoo codigo com opcionais ok falhou: %v", err)
		}
		if _, err := ExecutarPreVoo(context.Background(), criarFixtureCiclo1(t, "flow", "uuid-1"), DepsPreVoo{
			FlowCDPURL: "http://[::1]:namedport",
		}); err == nil {
			t.Fatalf("esperado erro de URL invalida no CDP")
		}
		dirFlow := criarFixtureCiclo1(t, "flow", "uuid-png-1")
		contratos, err := CarregarContratosEntrada(dirFlow)
		if err != nil {
			t.Fatal(err)
		}
		cliPNG := &clienteFlowPNGFake{}
		cenarios, err := ExecutarEtapa4Cenarios(context.Background(), dirFlow, contratos.Projeto, contratos.Script, DepsCenarios{
			Flow: cliPNG,
		})
		if err != nil || !strings.HasSuffix(cenarios[0].Arquivo, ".png") {
			t.Fatalf("esperado arquivo .png no modo flow, obtido cenarios=%+v err=%v", cenarios, err)
		}
	})

	t.Run("CriarProjetoFlow grava flow_projeto em projeto.json", func(t *testing.T) {
		dirProj := criarFixtureCiclo1(t, "hibrido", "")
		uuidReal := "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
		pid, err := CriarProjetoFlow(context.Background(), dirProj, "http://127.0.0.1:9222", func(_ context.Context, _ string, _ ...string) ([]byte, error) {
			return []byte(`{"ok":true,"flow_projeto":"` + uuidReal + `"}`), nil
		})
		if err != nil || pid != uuidReal {
			t.Fatalf("CriarProjetoFlow falhou: pid=%q err=%v", pid, err)
		}
		contratos, err := CarregarContratosEntrada(dirProj)
		if err != nil || contratos.Projeto.FlowProjeto != uuidReal {
			t.Fatalf("flow_projeto nao gravado em projeto.json: %+v err=%v", contratos.Projeto, err)
		}
		if _, err := CriarProjetoFlow(context.Background(), dirProj, "http://127.0.0.1:9222", func(_ context.Context, _ string, _ ...string) ([]byte, error) {
			return []byte(`{"ok":true,"flow_projeto":""}`), nil
		}); err == nil {
			t.Fatalf("esperado erro quando flow_projeto retornado e vazio")
		}
		if _, err := CriarProjetoFlow(context.Background(), dirProj, "http://127.0.0.1:9222", func(_ context.Context, _ string, _ ...string) ([]byte, error) {
			return []byte("boom"), errors.New("exit 1")
		}); err == nil {
			t.Fatalf("esperado erro quando script criar falha")
		}
		if _, err := CriarProjetoFlow(context.Background(), t.TempDir(), "http://127.0.0.1:9222", func(_ context.Context, _ string, _ ...string) ([]byte, error) {
			return []byte(`{"ok":true,"flow_projeto":"` + uuidReal + `"}`), nil
		}); err == nil {
			t.Fatalf("esperado erro quando projeto.json nao existe")
		}
		if _, err := CriarProjetoFlow(context.Background(), "\x00invalido", "http://127.0.0.1:9222", nil); err == nil {
			t.Fatalf("esperado erro quando projDir invalido")
		}
		if _, err := CriarProjetoFlow(context.Background(), dirProj, "", func(_ context.Context, _ string, _ ...string) ([]byte, error) {
			return []byte("sem cdp"), errors.New("exit 1")
		}); err == nil {
			t.Fatalf("esperado erro sem porta CDP ativa ou execucao falhando")
		}
		if err := atualizarFlowProjetoEmDisco("\x00invalido", "u1"); err == nil {
			t.Fatalf("esperado erro em atualizarFlowProjetoEmDisco com path invalido")
		}
	})

	t.Run("RecortadorRembgDespill faz fallback para folha mestre quando poses_6x ausente", func(t *testing.T) {
		dir := t.TempDir()
		folha := filepath.Join(dir, "flow", "fichas", "heroi_sheet.png")
		if err := os.MkdirAll(filepath.Dir(folha), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(folha, []byte("SIMULADO"), 0o644); err != nil {
			t.Fatal(err)
		}

		var entradaCapturada string
		rec := &RecortadorRembgDespill{
			ExecCmd: func(_ context.Context, _ string, args ...string) ([]byte, error) {
				for i := 0; i < len(args)-1; i++ {
					if args[i] == "--entrada" {
						entradaCapturada = args[i+1]
					}
				}
				return []byte(`[{"id":"heroi_pose_1","personagem_id":"heroi","tipo":"pose","arquivo_png":"sprites/heroi_pose_1.png","largura":128,"altura":128,"alfa_valido":true,"despill_verde":true,"bounding_box":{"x":0,"y":0,"w":128,"h":128}}]`), nil
			},
		}

		ficha := FichaPersonagem{
			ID:                "heroi",
			Nome:              "Herói",
			DescricaoImutavel: "Miniatura",
			Fichas: map[string]string{
				"turnaround":     "flow/heroi_turnaround.png",
				"expressoes_3x2": "flow/heroi_expressoes.png",
				"poses_6x":       "flow/heroi_poses.png",
			},
		}
		sprites, err := rec.RecortarEFatiar(context.Background(), dir, ficha)
		if err != nil {
			t.Fatalf("RecortarEFatiar falhou: %v", err)
		}
		if len(sprites) != 1 || !strings.HasSuffix(entradaCapturada, "heroi_sheet.png") {
			t.Fatalf("esperado fallback para heroi_sheet.png, obtido: %q sprites=%+v", entradaCapturada, sprites)
		}
	})
}

type clienteFlowPNGFake struct{}

func (clienteFlowPNGFake) GerarCenario(_ context.Context, _ RequisicaoFlow) ([]byte, error) {
	return []byte("\x89PNG\r\n\x1a\n"), nil
}
