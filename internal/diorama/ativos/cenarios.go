package ativos

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

type RequisicaoFlow struct {
	CenarioID   string
	Prompt      string
	FlowProjeto string
	Referencias []string
}

type ClienteFlow interface {
	GerarCenario(ctx context.Context, req RequisicaoFlow) ([]byte, error)
}

type DepsCenarios struct {
	Flow         ClienteFlow
	Referencias  []string
	Dispositivos []string
	Dormir       func(time.Duration)
	MaxRetries   int
}

// ModoProducaoProjeto devolve o modo de producao efetivo; CarregarContratosEntrada ja normaliza v1 e v2.
func ModoProducaoProjeto(proj roteiro.Projeto) string {
	return strings.TrimSpace(proj.EstiloVisual)
}

// ExecutarEtapa4Cenarios gera cenarios nos modos codigo, flow ou hibrido com backoff, trava e idempotencia.
func ExecutarEtapa4Cenarios(ctx context.Context, projDir string, proj roteiro.Projeto, script roteiro.Script, deps DepsCenarios) ([]Cenario, error) {
	modo := ModoProducaoProjeto(proj)
	if !ModoVisualValido(modo) {
		return nil, fmt.Errorf("%w: modo_producao %q", ErrPreVooBloqueado, modo)
	}
	if (modo == "flow" || modo == "hibrido") && strings.TrimSpace(proj.FlowProjeto) == "" {
		return nil, fmt.Errorf("%w: flow_projeto obrigatorio", ErrFlowTravaAusente)
	}
	deps = normalizarDepsCenarios(deps)
	cache := carregarCacheCenarios(projDir)
	qtd := minimoCenarios(len(script.Lines))
	saida := make([]Cenario, 0, qtd)
	for i := 0; i < qtd; i++ {
		c, err := gerarOuReusarCenario(ctx, projDir, proj, script, i, cache, deps)
		if err != nil {
			return nil, err
		}
		saida = append(saida, c)
	}
	if err := ValidarDensidadeCenarios(saida, len(script.Lines)); err != nil {
		return nil, err
	}
	if err := gravarJSONSeguro(projDir, filepath.Join("flow", "cenarios.json"), saida); err != nil {
		return nil, err
	}
	return saida, nil
}

// ValidarDensidadeCenarios aplica o Gate da Etapa 4: minimo de 1 cenario distinto a cada 2 linhas.
func ValidarDensidadeCenarios(cenarios []Cenario, totalLinhas int) error {
	minimo := minimoCenarios(totalLinhas)
	vistos := make(map[string]struct{}, len(cenarios))
	for _, c := range cenarios {
		if strings.TrimSpace(c.ID) != "" {
			vistos[c.ID] = struct{}{}
		}
	}
	if len(vistos) < minimo {
		return fmt.Errorf("%w: %d cenarios distintos para %d linhas (minimo %d)", ErrDensidadeCenarios, len(vistos), totalLinhas, minimo)
	}
	return nil
}

func minimoCenarios(totalLinhas int) int {
	if totalLinhas <= 0 {
		return 1
	}
	return (totalLinhas + 1) / 2
}

func normalizarDepsCenarios(d DepsCenarios) DepsCenarios {
	if d.Dormir == nil {
		d.Dormir = time.Sleep
	}
	if d.MaxRetries <= 0 {
		d.MaxRetries = 3
	}
	return d
}

func gerarOuReusarCenario(ctx context.Context, projDir string, proj roteiro.Projeto, script roteiro.Script, idx int, cache map[string]Cenario, deps DepsCenarios) (Cenario, error) {
	id := fmt.Sprintf("cen_%02d", idx+1)
	modo := resolverModoCenario(ModoProducaoProjeto(proj), idx)
	texto := textoBaseCenario(script, deps.Dispositivos, idx)
	hash := calcularHashCenario(id, modo, proj.FlowProjeto, texto, deps.Referencias)
	if anterior, ok := cache[id]; ok && anterior.HashEntrada == hash && arquivoExisteSeguro(projDir, anterior.Arquivo) {
		return anterior, nil
	}
	conteudo, err := produzirConteudoCenario(ctx, id, modo, proj.FlowProjeto, texto, deps)
	if err != nil {
		return Cenario{}, err
	}
	ext := ".svg"
	if modo == "flow" && !strings.Contains(string(conteudo), "<svg") {
		ext = ".png"
	}
	relArq := filepath.ToSlash(filepath.Join("flow", id+ext))
	if err := gravarBytesSeguro(projDir, relArq, conteudo); err != nil {
		return Cenario{}, err
	}
	return Cenario{
		ID:          id,
		Nome:        fmt.Sprintf("Cenario %02d", idx+1),
		Modo:        modo,
		Arquivo:     relArq,
		HashEntrada: hash,
		Referencias: append([]string(nil), deps.Referencias...),
		FlowProjeto: proj.FlowProjeto,
	}, nil
}

func resolverModoCenario(estilo string, _ int) string {
	if estilo == "hibrido" {
		return "flow"
	}
	return estilo
}

func produzirConteudoCenario(ctx context.Context, id, modo, flowProjeto, texto string, deps DepsCenarios) ([]byte, error) {
	if modo == "codigo" {
		return gerarSVGAutoral(id, texto), nil
	}
	if deps.Flow == nil {
		return nil, fmt.Errorf("%w: cliente Flow ausente para %s", ErrPreVooBloqueado, id)
	}
	req := RequisicaoFlow{CenarioID: id, Prompt: texto, FlowProjeto: flowProjeto, Referencias: deps.Referencias}
	return executarComBackoff(ctx, req, deps)
}

func executarComBackoff(ctx context.Context, req RequisicaoFlow, deps DepsCenarios) ([]byte, error) {
	var ultimoErr error
	espera := 50 * time.Millisecond
	for tentativa := 1; tentativa <= deps.MaxRetries; tentativa++ {
		dados, err := deps.Flow.GerarCenario(ctx, req)
		if err == nil {
			return dados, nil
		}
		ultimoErr = err
		if tentativa < deps.MaxRetries {
			deps.Dormir(espera)
			espera *= 2
		}
	}
	return nil, fmt.Errorf("diorama-ativos: fila Flow esgotada para %s: %w", req.CenarioID, ultimoErr)
}

func gerarSVGAutoral(id, titulo string) []byte {
	esc := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(titulo)
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1280 720" data-cenario="%s">
  <defs>
    <linearGradient id="bg_%s" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%%" stop-color="#0b132b"/>
      <stop offset="100%%" stop-color="#1c2541"/>
    </linearGradient>
  </defs>
  <rect width="1280" height="720" fill="url(#bg_%s)"/>
  <polygon points="180,500 640,340 1100,500 640,660" fill="#2b3a67" stroke="#5bc0be" stroke-width="3"/>
  <polygon points="180,500 640,660 640,695 180,535" fill="#1c2541"/>
  <polygon points="1100,500 640,660 640,695 1100,535" fill="#141c33"/>
  <rect x="480" y="250" width="320" height="190" rx="14" fill="#3a506b" stroke="#6fffe9" stroke-width="2"/>
  <circle cx="240" cy="170" r="72" fill="#5bc0be" fill-opacity="0.22"/>
  <text x="96" y="640" fill="#f8fafc" font-size="24" font-family="sans-serif">%s</text>
</svg>
`, id, id, id, esc)
	return []byte(svg)
}

func textoBaseCenario(script roteiro.Script, dispositivos []string, idx int) string {
	pos := idx * 2
	if pos < len(dispositivos) && strings.TrimSpace(dispositivos[pos]) != "" {
		return dispositivos[pos]
	}
	if pos >= len(script.Lines) {
		pos = len(script.Lines) - 1
	}
	return script.Lines[pos].Text
}

func calcularHashCenario(id, modo, flowProjeto, texto string, refs []string) string {
	h := sha256.New()
	h.Write([]byte(id + "|" + modo + "|" + flowProjeto + "|" + texto + "|" + strings.Join(refs, ",")))
	return hex.EncodeToString(h.Sum(nil))
}

func carregarCacheCenarios(projDir string) map[string]Cenario {
	out := make(map[string]Cenario)
	var lista []Cenario
	if err := lerJSONSeguro(projDir, filepath.Join("flow", "cenarios.json"), &lista); err != nil {
		return out
	}
	for _, c := range lista {
		out[c.ID] = c
	}
	return out
}

// PedidoFundo descreve o fundo que a Etapa 4 precisa: id kebab-case, tags de busca e prompt da cena.
type PedidoFundo struct {
	ID     string
	Tags   []string
	Prompt string
}

// ResultadoFundo informa o item usado, a copia no projeto e se veio do inventario.
type ResultadoFundo struct {
	Item        inventario.Item
	Arquivo     string
	Reutilizado bool
}

// ObterFundoInventario consulta o inventario antes de gerar: item aprovado com as tags e reutilizado
// com hash conferido; ausente, gera no Flow e registra o item com prompt, tags e hash.
func ObterFundoInventario(ctx context.Context, projDir string, s inventario.Store, estilo string, p PedidoFundo, deps DepsCenarios, origem string) (ResultadoFundo, error) {
	itens, err := s.Buscar(estilo, "cenario", p.Tags)
	if err != nil {
		return ResultadoFundo{}, err
	}
	if len(itens) > 0 {
		return reutilizarFundo(projDir, s, itens[0])
	}
	it, err := gerarFundoNovo(ctx, s, estilo, p, normalizarDepsCenarios(deps), origem)
	if err != nil {
		return ResultadoFundo{}, err
	}
	res, err := reutilizarFundo(projDir, s, it)
	res.Reutilizado = false
	return res, err
}

// reutilizarFundo confere o hash na origem, copia ao projeto e confere o hash da copia.
func reutilizarFundo(projDir string, s inventario.Store, it inventario.Item) (ResultadoFundo, error) {
	if err := s.Verificar(it); err != nil {
		return ResultadoFundo{}, err
	}
	dados, err := os.ReadFile(filepath.Join(s.Raiz, it.Arquivo))
	if err != nil {
		return ResultadoFundo{}, err
	}
	rel := filepath.ToSlash(filepath.Join("flow", it.ID+filepath.Ext(it.Arquivo)))
	if err := gravarBytesSeguro(projDir, rel, dados); err != nil {
		return ResultadoFundo{}, err
	}
	copia, err := ResolverCaminhoSeguro(projDir, rel)
	if err != nil {
		return ResultadoFundo{}, err
	}
	h, err := inventario.HashArquivo(copia)
	if err != nil {
		return ResultadoFundo{}, err
	}
	if h != it.Hash {
		return ResultadoFundo{}, fmt.Errorf("%w: copia de %q: esperado %s, obtido %s", inventario.ErrHashDivergente, it.ID, it.Hash, h)
	}
	return ResultadoFundo{Item: it, Arquivo: rel, Reutilizado: true}, nil
}

// gerarFundoNovo gera no Flow com o prompt_mestre do estilo e grava arquivo e item sem aprovacao.
func gerarFundoNovo(ctx context.Context, s inventario.Store, estilo string, p PedidoFundo, deps DepsCenarios, origem string) (inventario.Item, error) {
	e, err := s.LerEstilo(estilo)
	if err != nil {
		return inventario.Item{}, err
	}
	if deps.Flow == nil {
		return inventario.Item{}, fmt.Errorf("%w: cliente Flow ausente para %s", ErrPreVooBloqueado, p.ID)
	}
	prompt := strings.TrimSpace(e.PromptMestre + " " + p.Prompt)
	dados, err := executarComBackoff(ctx, RequisicaoFlow{CenarioID: p.ID, Prompt: prompt}, deps)
	if err != nil {
		return inventario.Item{}, err
	}
	ext := ".png"
	if strings.Contains(string(dados), "<svg") {
		ext = ".svg"
	}
	rel := filepath.ToSlash(filepath.Join("inventario", "video", "estilos", e.ID, "arquivos", "cenario", p.ID+ext))
	if err := gravarBytesSeguro(s.Raiz, rel, dados); err != nil {
		return inventario.Item{}, err
	}
	return registrarFundoNovo(s, e, p, prompt, rel, origem)
}

// registrarFundoNovo calcula o hash do arquivo e grava o item com prompt e tags.
func registrarFundoNovo(s inventario.Store, e inventario.Estilo, p PedidoFundo, prompt, rel, origem string) (inventario.Item, error) {
	h, err := inventario.HashArquivo(filepath.Join(s.Raiz, rel))
	if err != nil {
		return inventario.Item{}, err
	}
	it := inventario.Item{
		SchemaVersion: 1,
		ID:            p.ID,
		Tipo:          "cenario",
		Estilo:        e.ID,
		EstiloVersao:  e.Versao,
		Tags:          append([]string(nil), p.Tags...),
		Parametros:    map[string]any{"prompt": prompt},
		Arquivo:       rel,
		Hash:          h,
		Origem:        origem,
	}
	return it, s.GravarItem(it)
}

func arquivoExisteSeguro(projDir, rel string) bool {
	p, err := ResolverCaminhoSeguro(projDir, rel)
	if err != nil {
		return false
	}
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func gravarBytesSeguro(projDir, rel string, dados []byte) error {
	destino, err := ResolverCaminhoSeguro(projDir, rel)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destino, dados, 0o644)
}
