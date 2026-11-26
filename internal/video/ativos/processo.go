package ativos

import (
	"context"
	"path/filepath"
	"slices"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

const (
	GateEtapa3Elenco       = "etapa3_elenco"
	GateEtapa4QuadrosTeste = "etapa4_quadros_teste"
	GateEtapa5Storyboard   = "etapa5_storyboard"
)

type DepsProcesso struct {
	PreVoo     DepsPreVoo
	Recortador RecortadorSprite
	Cenarios   DepsCenarios
}

type OpcoesExecucao struct {
	ProjetoDir     string
	Elenco         []FichaPersonagem
	GatesAprovados []string
}

type ResultadoExecucao struct {
	PreVoo         RelatorioPreVoo   `json:"prevoo"`
	Fichas         []FichaAtivo      `json:"fichas"`
	Sprites        []DescritorSprite `json:"sprites"`
	Cenarios       []Cenario         `json:"cenarios"`
	Storyboard     Storyboard        `json:"storyboard"`
	AguardandoGate string            `json:"aguardando_gate,omitempty"`
}

type Processo struct {
	deps DepsProcesso
}

// NovoProcesso cria o orquestrador das Etapas 3, 4 e 5 do Ciclo 2.
func NovoProcesso(deps DepsProcesso) *Processo {
	return &Processo{deps: deps}
}

// Executar percorre Pre-voo e Etapas 3, 4 e 5 respeitando politica_gates de projeto.json.
func (p *Processo) Executar(ctx context.Context, opts OpcoesExecucao) (*ResultadoExecucao, error) {
	rel, err := ExecutarPreVoo(ctx, opts.ProjetoDir, p.deps.PreVoo)
	if err != nil {
		return &ResultadoExecucao{PreVoo: rel}, err
	}
	contratos, err := CarregarContratosEntrada(opts.ProjetoDir)
	if err != nil {
		return &ResultadoExecucao{PreVoo: rel}, err
	}
	res := &ResultadoExecucao{PreVoo: rel}
	fichas4F := ExtrairFichasDeEntrada(contratos.Script, contratos.PlanoMD)
	fichasGeradas, err := ExecutarGeracaoFichas4Faixas(ctx, opts.ProjetoDir, contratos.Projeto, fichas4F, p.deps.Cenarios)
	if err != nil {
		return res, err
	}
	res.Fichas = fichasGeradas
	elenco := resolverElencoPadrao(opts.Elenco, contratos.Projeto, fichasGeradas)
	sprites, err := ExecutarEtapa3Elenco(ctx, opts.ProjetoDir, elenco, p.deps.Recortador)
	if err != nil {
		return res, err
	}
	res.Sprites = sprites
	if exigeParadaGate(contratos.Projeto.PoliticaGates, GateEtapa3Elenco, opts.GatesAprovados) {
		res.AguardandoGate = GateEtapa3Elenco
		return res, nil
	}
	return p.executarEtapas4E5(ctx, opts, contratos, res)
}

func (p *Processo) executarEtapas4E5(ctx context.Context, opts OpcoesExecucao, contratos ContratosEntrada, res *ResultadoExecucao) (*ResultadoExecucao, error) {
	depsCen := p.deps.Cenarios
	if len(depsCen.Dispositivos) == 0 {
		depsCen.Dispositivos = ExtrairDispositivosPlanoMD(contratos.PlanoMD)
	}
	cenarios, err := ExecutarEtapa4Cenarios(ctx, opts.ProjetoDir, contratos.Projeto, contratos.Script, depsCen)
	if err != nil {
		return res, err
	}
	res.Cenarios = cenarios
	if exigeParadaGate(contratos.Projeto.PoliticaGates, GateEtapa4QuadrosTeste, opts.GatesAprovados) {
		res.AguardandoGate = GateEtapa4QuadrosTeste
		return res, nil
	}
	sb, err := ExecutarEtapa5Storyboard(opts.ProjetoDir, contratos.Projeto, contratos.Script, res.Sprites, cenarios)
	if err != nil {
		return res, err
	}
	res.Storyboard = sb
	if exigeParadaGate(contratos.Projeto.PoliticaGates, GateEtapa5Storyboard, opts.GatesAprovados) {
		res.AguardandoGate = GateEtapa5Storyboard
		return res, nil
	}
	return res, nil
}

func exigeParadaGate(pg roteiro.PoliticaGates, gateID string, aprovados []string) bool {
	if !pg.ExigeAprovacaoHumana(gateID) {
		return false
	}
	return !slices.Contains(aprovados, gateID)
}

func resolverElencoPadrao(elenco []FichaPersonagem, proj roteiro.Projeto, fichasGeradas []FichaAtivo) []FichaPersonagem {
	if len(elenco) > 0 {
		return elenco
	}
	var res []FichaPersonagem
	for _, f := range fichasGeradas {
		if f.Categoria == "personagem" {
			res = append(res, converterFichaEmPersonagem(f))
		}
	}
	if len(res) > 0 {
		return res
	}
	return elencoNarradorPadrao(proj)
}

func converterFichaEmPersonagem(f FichaAtivo) FichaPersonagem {
	folha := f.FolhaMestrePNG
	if folha == "" {
		folha = filepath.ToSlash(filepath.Join("flow", "fichas", f.ID+"_sheet.png"))
	}
	return FichaPersonagem{
		ID:                f.ID,
		Nome:              f.Nome,
		DescricaoImutavel: f.DescricaoImutavel,
		Fichas: map[string]string{
			"turnaround":     folha,
			"expressoes_3x2": folha,
			"poses_6x":       folha,
		},
	}
}

func elencoNarradorPadrao(proj roteiro.Projeto) []FichaPersonagem {
	folha := filepath.ToSlash(filepath.Join("flow", "fichas", "narrador_sheet.png"))
	return []FichaPersonagem{
		{
			ID:                "narrador",
			Nome:              "Narrador",
			DescricaoImutavel: "Personagem padrao do projeto " + proj.Slug + " em estilo " + proj.EstiloVisual,
			Fichas: map[string]string{
				"turnaround":     folha,
				"expressoes_3x2": folha,
				"poses_6x":       folha,
			},
		},
	}
}
