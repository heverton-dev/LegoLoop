package ativos

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/browser"
)

//go:embed scripts/flow_driver.py
var scriptFlowDriver []byte

//go:embed scripts/sprite_despill.py
var scriptSpriteDespill []byte

// ExecFunc abstrai chamadas de processo externo para manter testes 100% hermeticos.
type ExecFunc func(ctx context.Context, name string, args ...string) ([]byte, error)

// ClienteFlowCDP conecta ao Google Flow via Chrome CDP (porta 9333) usando Playwright.
type ClienteFlowCDP struct {
	CDPURL     string
	ProjetoDir string
	ExecCmd    ExecFunc
}

// NovoClienteFlowCDP cria o driver real de Google Flow conectado ao Chrome CDP (9222/9333).
func NovoClienteFlowCDP(cdpURL, projDir string) *ClienteFlowCDP {
	return &ClienteFlowCDP{CDPURL: strings.TrimSpace(cdpURL), ProjetoDir: projDir}
}

// GerarCenario aciona o driver Playwright via Chrome CDP respeitando a trava flow_projeto.
func (c *ClienteFlowCDP) GerarCenario(ctx context.Context, req RequisicaoFlow) ([]byte, error) {
	if strings.TrimSpace(req.FlowProjeto) == "" {
		return nil, fmt.Errorf("%w: flow_projeto obrigatorio no driver CDP", ErrFlowTravaAusente)
	}
	cdpURL := c.CDPURL
	if cdpURL == "" {
		if resolvida, err := browser.ResolverCDPURL(ctx, nil, nil); err == nil {
			cdpURL = resolvida
		} else {
			cdpURL = "http://127.0.0.1:9222"
		}
	}
	scriptPath, err := materializarScript(c.ProjetoDir, ".legoloop_flow_driver.py", scriptFlowDriver)
	if err != nil {
		return nil, err
	}
	saidaArq, err := ResolverCaminhoSeguro(c.ProjetoDir, filepath.Join("flow", req.CenarioID+".png"))
	if err != nil {
		return nil, err
	}
	args := []string{scriptPath, "--acao", "gerar", "--baixar-todas", "--cdp", cdpURL, "--flow-projeto", req.FlowProjeto, "--prompt", req.Prompt, "--saida-arquivo", saidaArq}
	for _, ref := range req.Referencias {
		args = append(args, "--ref", ref)
	}
	execFn := resolverExecCmd(c.ExecCmd)
	if out, err := execFn(ctx, "python", args...); err != nil {
		return nil, fmt.Errorf("diorama-ativos: flow_driver falhou: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return os.ReadFile(saidaArq)
}

// RecortadorRembgDespill executa rembg + despill verde via script Python embutido.
type RecortadorRembgDespill struct {
	ExecCmd ExecFunc
}

// NovoRecortadorRembgDespill devolve o recortador real com rembg + despill verde.
func NovoRecortadorRembgDespill() *RecortadorRembgDespill {
	return &RecortadorRembgDespill{}
}

// RecortarEFatiar aciona sprite_despill.py e decodifica os descritores gerados.
func (r *RecortadorRembgDespill) RecortarEFatiar(ctx context.Context, projDir string, ficha FichaPersonagem) ([]DescritorSprite, error) {
	scriptPath, err := materializarScript(projDir, ".legoloop_sprite_despill.py", scriptSpriteDespill)
	if err != nil {
		return nil, err
	}
	entrada, err := resolverEntradaSprite(projDir, ficha)
	if err != nil {
		return nil, err
	}
	execFn := resolverExecCmd(r.ExecCmd)
	out, err := execFn(ctx, "python", scriptPath, "--projeto", projDir, "--personagem-id", ficha.ID, "--entrada", entrada)
	if err != nil {
		return nil, fmt.Errorf("diorama-ativos: sprite_despill falhou: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	var sprites []DescritorSprite
	if err := json.Unmarshal(out, &sprites); err != nil {
		return nil, fmt.Errorf("diorama-ativos: saida invalida de sprite_despill: %w", err)
	}
	return sprites, nil
}

// RecortadorDoEstilo devolve o recortador rembg + despill para renderizador de recorte (AC-E09.3).
func RecortadorDoEstilo(renderizador string) (RecortadorSprite, error) {
	if !RenderizadorDeRecorte(renderizador) {
		return nil, fmt.Errorf("%w: renderizador %q", ErrRenderizadorNaoRecorte, renderizador)
	}
	return NovoRecortadorRembgDespill(), nil
}

func resolverEntradaSprite(projDir string, ficha FichaPersonagem) (string, error) {
	caminho := ficha.Fichas["poses_6x"]
	if _, err := os.Stat(filepath.Join(projDir, caminho)); os.IsNotExist(err) {
		for _, cand := range []string{
			ficha.Fichas["turnaround"],
			ficha.Fichas["expressoes_3x2"],
			filepath.Join("flow", "fichas", ficha.ID+"_sheet.png"),
			filepath.Join("flow", "ficha_"+ficha.ID+".png"),
			filepath.Join("flow", ficha.ID+".png"),
		} {
			if strings.TrimSpace(cand) != "" && arquivoExisteEm(projDir, cand) {
				caminho = cand
				break
			}
		}
	}
	return ResolverCaminhoSeguro(projDir, caminho)
}

func arquivoExisteEm(dir, rel string) bool {
	_, err := os.Stat(filepath.Join(dir, rel))
	return err == nil
}

// CriarProjetoFlow cria um projeto novo no Google Flow via CDP e grava flow_projeto em projeto.json.
func CriarProjetoFlow(ctx context.Context, projDir, cdpURL string, execFn ExecFunc) (string, error) {
	cdp := strings.TrimSpace(cdpURL)
	if cdp == "" {
		resolvida, err := browser.ResolverCDPURL(ctx, nil, nil)
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrPreVooBloqueado, err)
		}
		cdp = resolvida
	}
	scriptPath, err := materializarScript(projDir, ".legoloop_flow_driver.py", scriptFlowDriver)
	if err != nil {
		return "", err
	}
	out, err := resolverExecCmd(execFn)(ctx, "python", scriptPath, "--acao", "criar", "--cdp", cdp)
	if err != nil {
		return "", fmt.Errorf("diorama-ativos: criar projeto no Flow falhou: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	var resp struct {
		FlowProjeto string `json:"flow_projeto"`
	}
	if err := json.Unmarshal(out, &resp); err != nil || strings.TrimSpace(resp.FlowProjeto) == "" {
		return "", fmt.Errorf("%w: resposta invalida ao criar projeto no Flow (%s)", ErrFlowTravaAusente, strings.TrimSpace(string(out)))
	}
	if err := atualizarFlowProjetoEmDisco(projDir, resp.FlowProjeto); err != nil {
		return "", err
	}
	return resp.FlowProjeto, nil
}

func atualizarFlowProjetoEmDisco(projDir, flowUUID string) error {
	contratos, err := CarregarContratosEntrada(projDir)
	if err != nil {
		return err
	}
	contratos.Projeto.FlowProjeto = flowUUID
	return gravarJSONSeguro(projDir, "projeto.json", contratos.Projeto)
}

func materializarScript(projDir, nome string, dados []byte) (string, error) {
	destino, err := ResolverCaminhoSeguro(projDir, filepath.Join("flow", nome))
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(destino, dados, 0o644); err != nil {
		return "", err
	}
	return destino, nil
}

func resolverExecCmd(fn ExecFunc) ExecFunc {
	if fn != nil {
		return fn
	}
	return func(ctx context.Context, name string, args ...string) ([]byte, error) {
		out, err := exec.CommandContext(ctx, name, args...).Output()
		if err != nil {
			if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
				return ee.Stderr, err
			}
		}
		return out, err
	}
}
