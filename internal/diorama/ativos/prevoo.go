package ativos

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/browser"
	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
	"github.com/heverton-dev/LegoLoop/internal/jail"
)

// AvisoEstiloVisualLegado nomeia o aviso do pre-voo quando projeto.json ainda usa estilo_visual.
const AvisoEstiloVisualLegado = "AVISO_ESTILO_VISUAL_LEGADO"

var (
	ErrEstiloInexistente = errors.New("diorama-ativos: estilo inexistente no inventario")
	ErrEstiloNaoValidado = errors.New("diorama-ativos: estilo nao validado")
)

// ContratoProducao guarda os campos v2 de projeto.json (modo_producao e estilo) e o legado estilo_visual.
type ContratoProducao struct {
	ModoProducao string `json:"modo_producao,omitempty"`
	Estilo       string `json:"estilo,omitempty"`
	EstiloVisual string `json:"estilo_visual,omitempty"`
}

// Legado indica projeto.json sem modo_producao, lido pelo campo antigo estilo_visual.
func (c ContratoProducao) Legado() bool {
	return strings.TrimSpace(c.ModoProducao) == ""
}

// ModoEfetivo devolve modo_producao ou, no contrato legado, estilo_visual.
func (c ContratoProducao) ModoEfetivo() string {
	if c.Legado() {
		return strings.TrimSpace(c.EstiloVisual)
	}
	return strings.TrimSpace(c.ModoProducao)
}

type ItemPreVoo struct {
	Nome    string `json:"nome"`
	Nivel   string `json:"nivel"`
	Detalhe string `json:"detalhe,omitempty"`
}

type RelatorioPreVoo struct {
	Bloqueado bool         `json:"bloqueado"`
	Itens     []ItemPreVoo `json:"itens"`
}

type DepsPreVoo struct {
	FlowCDPURL     string
	HTTPClient     *http.Client
	ChecarModulo   func(nome string) bool
	ChecarOpcional func(nome string) bool
	InventarioDir  string
}

type ContratosEntrada struct {
	Projeto  roteiro.Projeto
	Producao ContratoProducao
	Marca    roteiro.Marca
	Script   roteiro.Script
	PlanoMD  string
}

// ResolverCaminhoSeguro garante que o caminho permanece dentro da pasta do projeto.
func ResolverCaminhoSeguro(projDir, relativo string) (string, error) {
	return jail.DentroDaRaiz(projDir, relativo)
}

// CarregarContratosEntrada le e valida os artefatos imutaveis produzidos pelo Ciclo 1.
func CarregarContratosEntrada(projDir string) (ContratosEntrada, error) {
	var c ContratosEntrada
	if err := lerProjetoV2(projDir, &c); err != nil {
		return c, fmt.Errorf("%w: projeto.json: %v", ErrPreVooBloqueado, err)
	}
	if err := lerJSONSeguro(projDir, "marca.json", &c.Marca); err != nil {
		return c, fmt.Errorf("%w: marca.json: %v", ErrPreVooBloqueado, err)
	}
	if err := lerJSONSeguro(projDir, "script.json", &c.Script); err != nil {
		return c, fmt.Errorf("%w: script.json: %v", ErrPreVooBloqueado, err)
	}
	if err := carregarPlanoEValidarLinhas(projDir, &c); err != nil {
		return c, err
	}
	return c, nil
}

// lerProjetoV2 le projeto.json v1 ou v2 e normaliza o modo efetivo em Projeto.EstiloVisual.
func lerProjetoV2(projDir string, c *ContratosEntrada) error {
	if err := lerJSONSeguro(projDir, "projeto.json", &c.Projeto); err != nil {
		return err
	}
	if err := lerJSONSeguro(projDir, "projeto.json", &c.Producao); err != nil {
		return err
	}
	c.Projeto.EstiloVisual = c.Producao.ModoEfetivo()
	return nil
}

func carregarPlanoEValidarLinhas(projDir string, c *ContratosEntrada) error {
	planoPath, err := ResolverCaminhoSeguro(projDir, "plano.md")
	if err != nil {
		return fmt.Errorf("%w: plano.md: %v", ErrPreVooBloqueado, err)
	}
	raw, err := os.ReadFile(planoPath)
	if err != nil || strings.TrimSpace(string(raw)) == "" {
		return fmt.Errorf("%w: plano.md ausente ou vazio", ErrPreVooBloqueado)
	}
	c.PlanoMD = string(raw)
	if len(c.Script.Lines) == 0 {
		return fmt.Errorf("%w: script.json sem linhas", ErrPreVooBloqueado)
	}
	return nil
}

// ExecutarPreVoo verifica contratos do Ciclo 1, modo visual, trava flow_projeto, dependencias reais e escrita.
func ExecutarPreVoo(ctx context.Context, projDir string, deps DepsPreVoo) (RelatorioPreVoo, error) {
	var rel RelatorioPreVoo
	deps = normalizarDepsPreVoo(deps)
	contratos, err := CarregarContratosEntrada(projDir)
	if err != nil {
		adicionarPreVoo(&rel, ItemPreVoo{Nome: "contratos_ciclo1", Nivel: "FALTA", Detalhe: err.Error()})
		return rel, err
	}
	adicionarPreVoo(&rel, ItemPreVoo{Nome: "contratos_ciclo1", Nivel: "OK"})
	if err := validarModoETrava(&rel, contratos.Projeto); err != nil {
		return rel, err
	}
	if err := validarEstiloPreVoo(&rel, projDir, contratos.Producao, deps.InventarioDir); err != nil {
		return rel, err
	}
	if err := checarEscritaSubpastas(&rel, projDir); err != nil {
		return rel, err
	}
	if modo := ModoProducaoProjeto(contratos.Projeto); modo == "flow" || modo == "hibrido" {
		if err := checarDependenciasFlow(ctx, &rel, deps); err != nil {
			return rel, err
		}
		return rel, nil
	}
	registrarAvisosOpcionais(&rel, deps.ChecarOpcional)
	return rel, nil
}

func normalizarDepsPreVoo(d DepsPreVoo) DepsPreVoo {
	if d.HTTPClient == nil {
		d.HTTPClient = &http.Client{Timeout: 2 * time.Second}
	}
	if d.ChecarModulo == nil {
		d.ChecarModulo = func(nome string) bool {
			return exec.Command("python", "-c", "import "+nome).Run() == nil
		}
	}
	if d.ChecarOpcional == nil {
		d.ChecarOpcional = func(nome string) bool {
			_, err := exec.LookPath(nome)
			return err == nil
		}
	}
	return d
}

func checarDependenciasFlow(ctx context.Context, rel *RelatorioPreVoo, deps DepsPreVoo) error {
	var candidatas []string
	if strings.TrimSpace(deps.FlowCDPURL) != "" {
		candidatas = []string{deps.FlowCDPURL}
	}
	cdpURL, err := browser.ResolverCDPURL(ctx, deps.HTTPClient, candidatas)
	if err != nil {
		cdpErr := fmt.Errorf("%w: nenhuma porta CDP ativa (9222/9333); rode 'legoloop diorama-ativos abrir-flow' (%v)", ErrPreVooBloqueado, err)
		adicionarPreVoo(rel, ItemPreVoo{Nome: "chrome_flow_cdp", Nivel: "FALTA", Detalhe: cdpErr.Error()})
		return cdpErr
	}
	adicionarPreVoo(rel, ItemPreVoo{Nome: "chrome_flow_cdp", Nivel: "OK", Detalhe: cdpURL})
	for _, mod := range []string{"playwright", "PIL", "numpy", "rembg"} {
		if !deps.ChecarModulo(mod) {
			err := fmt.Errorf("%w: modulo Python %q ausente (instale antes de rodar flow/hibrido)", ErrPreVooBloqueado, mod)
			adicionarPreVoo(rel, ItemPreVoo{Nome: "python_" + mod, Nivel: "FALTA", Detalhe: err.Error()})
			return err
		}
		adicionarPreVoo(rel, ItemPreVoo{Nome: "python_" + mod, Nivel: "OK"})
	}
	return nil
}

func validarModoETrava(rel *RelatorioPreVoo, p roteiro.Projeto) error {
	modo := ModoProducaoProjeto(p)
	if !ModoVisualValido(modo) {
		err := fmt.Errorf("%w: modo_producao invalido %q", ErrPreVooBloqueado, modo)
		adicionarPreVoo(rel, ItemPreVoo{Nome: "modo_producao", Nivel: "FALTA", Detalhe: err.Error()})
		return err
	}
	adicionarPreVoo(rel, ItemPreVoo{Nome: "modo_producao", Nivel: "OK", Detalhe: modo})
	if (modo == "flow" || modo == "hibrido") && strings.TrimSpace(p.FlowProjeto) == "" {
		err := fmt.Errorf("%w: modo %q exige flow_projeto", ErrFlowTravaAusente, modo)
		adicionarPreVoo(rel, ItemPreVoo{Nome: "trava_flow_projeto", Nivel: "FALTA", Detalhe: err.Error()})
		return err
	}
	adicionarPreVoo(rel, ItemPreVoo{Nome: "trava_flow_projeto", Nivel: "OK"})
	return nil
}

// validarEstiloPreVoo exige estilo validado no inventario; o contrato legado so recebe aviso.
func validarEstiloPreVoo(rel *RelatorioPreVoo, projDir string, p ContratoProducao, invDir string) error {
	if p.Legado() {
		adicionarPreVoo(rel, ItemPreVoo{Nome: AvisoEstiloVisualLegado, Nivel: "AVISO",
			Detalhe: "estilo_visual lido como modo_producao; migre para modo_producao e estilo"})
		if strings.TrimSpace(p.Estilo) == "" {
			return nil
		}
	}
	status, err := statusEstilo(localizarInventario(projDir, invDir), p.Estilo)
	if err == nil && status != "validado" {
		err = fmt.Errorf("%w: %w: %q esta %q", ErrPreVooBloqueado, ErrEstiloNaoValidado, p.Estilo, status)
	}
	if err != nil {
		adicionarPreVoo(rel, ItemPreVoo{Nome: "estilo", Nivel: "FALTA", Detalhe: err.Error()})
		return err
	}
	adicionarPreVoo(rel, ItemPreVoo{Nome: "estilo", Nivel: "OK", Detalhe: p.Estilo})
	return nil
}

// statusEstilo le inventario/video/estilos/<id>/estilo.json e devolve o campo status.
func statusEstilo(invDir, id string) (string, error) {
	if !idEstiloValido(id) {
		return "", fmt.Errorf("%w: %w: id %q", ErrPreVooBloqueado, ErrEstiloInexistente, id)
	}
	raw, err := os.ReadFile(filepath.Join(invDir, "estilos", id, "estilo.json"))
	if err != nil {
		return "", fmt.Errorf("%w: %w: %q em %s", ErrPreVooBloqueado, ErrEstiloInexistente, id, invDir)
	}
	var e struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF}), &e); err != nil {
		return "", fmt.Errorf("%w: estilo.json de %q invalido: %v", ErrPreVooBloqueado, id, err)
	}
	return e.Status, nil
}

// idEstiloValido aceita so kebab-case minusculo, o que impede travessia de pasta.
func idEstiloValido(id string) bool {
	if id == "" || id[0] == '-' || id[len(id)-1] == '-' {
		return false
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}

// localizarInventario usa o caminho explicito, LEGOLOOP_INVENTARIO_DIR ou sobe pastas ate achar inventario/video.
func localizarInventario(projDir, explicito string) string {
	if strings.TrimSpace(explicito) != "" {
		return explicito
	}
	if env := strings.TrimSpace(os.Getenv("LEGOLOOP_INVENTARIO_DIR")); env != "" {
		return env
	}
	dir, err := filepath.Abs(projDir)
	if err != nil {
		dir = projDir
	}
	for {
		alvo := filepath.Join(dir, "inventario", "video")
		if info, err := os.Stat(filepath.Join(alvo, "estilos")); err == nil && info.IsDir() {
			return alvo
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			return filepath.Join("inventario", "video")
		}
		dir = pai
	}
}

func checarEscritaSubpastas(rel *RelatorioPreVoo, projDir string) error {
	for _, sub := range []string{"flow", "sprites"} {
		alvo, err := ResolverCaminhoSeguro(projDir, sub)
		if err != nil {
			adicionarPreVoo(rel, ItemPreVoo{Nome: "pasta_" + sub, Nivel: "FALTA", Detalhe: err.Error()})
			return fmt.Errorf("%w: %v", ErrPreVooBloqueado, err)
		}
		if err := testarPermissaoEscrita(alvo); err != nil {
			adicionarPreVoo(rel, ItemPreVoo{Nome: "pasta_" + sub, Nivel: "FALTA", Detalhe: err.Error()})
			return fmt.Errorf("%w: pasta %s sem escrita: %v", ErrPreVooBloqueado, sub, err)
		}
		adicionarPreVoo(rel, ItemPreVoo{Nome: "pasta_" + sub, Nivel: "OK"})
	}
	return nil
}

func testarPermissaoEscrita(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	probe := filepath.Join(dir, ".prevoo_probe")
	if err := os.WriteFile(probe, []byte("ok"), 0o644); err != nil {
		return err
	}
	return os.Remove(probe)
}

func registrarAvisosOpcionais(rel *RelatorioPreVoo, checar func(string) bool) {
	for _, opt := range []string{"rembg", "chrome_flow_9333"} {
		if checar(opt) {
			adicionarPreVoo(rel, ItemPreVoo{Nome: opt, Nivel: "OK"})
		} else {
			adicionarPreVoo(rel, ItemPreVoo{Nome: opt, Nivel: "AVISO", Detalhe: "ausente no PATH local"})
		}
	}
}

func adicionarPreVoo(rel *RelatorioPreVoo, it ItemPreVoo) {
	rel.Itens = append(rel.Itens, it)
	if it.Nivel == "FALTA" {
		rel.Bloqueado = true
	}
}

func lerJSONSeguro(projDir, nome string, dst any) error {
	caminho, err := ResolverCaminhoSeguro(projDir, nome)
	if err != nil {
		return err
	}
	b, err := os.ReadFile(caminho)
	if err != nil {
		return err
	}
	b = bytes.TrimPrefix(b, []byte{0xEF, 0xBB, 0xBF})
	return json.Unmarshal(b, dst)
}
