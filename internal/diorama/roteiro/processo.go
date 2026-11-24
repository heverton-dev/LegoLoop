package roteiro

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/jail"
)

const maxCorrecoesRoteiro = 2

type FuncaoGerarRoteiro func(ctx context.Context, pauta PautaCandidata, anterior *RelatorioValidacao) (Script, string, error)
type FuncaoEscolherPauta func(pautas []PautaCandidata) (PautaCandidata, error)
type FuncaoAprovarHumano func(gateID string, resumo any) (bool, error)

type OpcoesExecucao struct {
	ProjetoDir string
	AcervoPath string
	Projeto    Projeto
	Marca      Marca
	Catalogo   CatalogoFontes
}

type ResultadoExecucao struct {
	PreVoo          RelatorioPreVoo    `json:"pre_voo"`
	PautaEscolhida  PautaCandidata     `json:"pauta_escolhida"`
	Validacao       RelatorioValidacao `json:"validacao"`
	CorrecoesUsadas int                `json:"correcoes_usadas"`
	ScriptPath      string             `json:"script_path"`
	PlanoPath       string             `json:"plano_path"`
	TemaRegistrado  bool               `json:"tema_registrado"`
}

type Processo struct {
	DepsPreVoo    DepsPreVoo
	GerarRoteiro  FuncaoGerarRoteiro
	EscolherPauta FuncaoEscolherPauta
	AprovarHumano FuncaoAprovarHumano
	Agora         func() time.Time
}

// Executar orquestra validacao de config, pre-voo, descoberta, geracao, ate 2 correcoes e gravacao pos-Gate.
func (p *Processo) Executar(ctx context.Context, opts OpcoesExecucao) (*ResultadoExecucao, error) {
	if err := ValidarConfiguracao(opts.Projeto, opts.Marca, opts.Catalogo); err != nil {
		return nil, err
	}
	pre := ExecutarPreVoo(ctx, opts.Projeto, opts.Catalogo, opts.ProjetoDir, filepath.Dir(opts.AcervoPath), p.DepsPreVoo)
	if pre.Bloqueado {
		return &ResultadoExecucao{PreVoo: pre}, ErrPreVooBloqueado
	}
	acervo, err := CarregarAcervoTemas(opts.AcervoPath)
	if err != nil {
		return nil, err
	}
	pautas, err := ColetarPautas(ctx, opts.Projeto, opts.Catalogo, acervo, p.DepsPreVoo.HTTPClient)
	if err != nil {
		return nil, err
	}
	pauta, err := p.selecionarPauta(opts.Projeto, pautas)
	if err != nil {
		return nil, err
	}
	return p.gerarEValidar(ctx, opts, pre, acervo, pauta)
}

func (p *Processo) selecionarPauta(proj Projeto, pautas []PautaCandidata) (PautaCandidata, error) {
	if len(pautas) == 0 {
		return PautaCandidata{}, fmt.Errorf("%w: nenhuma pauta nova disponivel", ErrTemaRepetido)
	}
	if proj.PoliticaGates.ExigeAprovacaoHumana("etapa1_pauta") && p.EscolherPauta != nil {
		return p.EscolherPauta(pautas)
	}
	melhor := pautas[0]
	for _, cand := range pautas[1:] {
		if cand.Score > melhor.Score {
			melhor = cand
		}
	}
	return melhor, nil
}

func (p *Processo) gerarEValidar(ctx context.Context, opts OpcoesExecucao, pre RelatorioPreVoo, acervo AcervoTemas, pauta PautaCandidata) (*ResultadoExecucao, error) {
	if p.GerarRoteiro == nil {
		return nil, fmt.Errorf("%w: gerador de roteiro nao configurado", ErrGateRoteiroReprovado)
	}
	var ultimoRel *RelatorioValidacao
	for tentativa := 0; tentativa <= maxCorrecoesRoteiro; tentativa++ {
		sc, plano, err := p.GerarRoteiro(ctx, pauta, ultimoRel)
		if err != nil {
			return nil, err
		}
		rel := ValidarRoteiro(sc, opts.Projeto, opts.Catalogo)
		if rel.Aprovado {
			return p.persistirAprovado(opts, pre, acervo, pauta, sc, plano, rel, tentativa)
		}
		ultimoRel = &rel
	}
	return &ResultadoExecucao{PreVoo: pre, PautaEscolhida: pauta, Validacao: *ultimoRel, CorrecoesUsadas: maxCorrecoesRoteiro},
		fmt.Errorf("%w: %d violacoes apos %d correcoes", ErrGateRoteiroReprovado, len(ultimoRel.Violacoes), maxCorrecoesRoteiro)
}

func (p *Processo) persistirAprovado(opts OpcoesExecucao, pre RelatorioPreVoo, acervo AcervoTemas, pauta PautaCandidata, sc Script, plano string, rel RelatorioValidacao, correcoes int) (*ResultadoExecucao, error) {
	if opts.Projeto.PoliticaGates.ExigeAprovacaoHumana("etapa2_roteiro") && p.AprovarHumano != nil {
		ok, err := p.AprovarHumano("etapa2_roteiro", sc)
		if err != nil || !ok {
			return nil, fmt.Errorf("%w: recusado no Stop-Gate humano", ErrGateRoteiroReprovado)
		}
	}
	scriptPath, err := jail.DentroDaRaiz(opts.ProjetoDir, "script.json")
	if err != nil {
		return nil, err
	}
	planoPath, err := jail.DentroDaRaiz(opts.ProjetoDir, "plano.md")
	if err != nil {
		return nil, err
	}
	if err := gravarSaidas(scriptPath, planoPath, sc, plano); err != nil {
		return nil, err
	}
	data := p.dataAtual()
	acervo.Temas = append(acervo.Temas, RegistroTema{Slug: pauta.Slug, Titulo: pauta.Titulo, Data: data, TipoPauta: pauta.TipoPauta})
	if err := SalvarAcervoTemas(opts.AcervoPath, acervo); err != nil {
		return nil, err
	}
	return &ResultadoExecucao{
		PreVoo: pre, PautaEscolhida: pauta, Validacao: rel, CorrecoesUsadas: correcoes,
		ScriptPath: scriptPath, PlanoPath: planoPath, TemaRegistrado: true,
	}, nil
}

func gravarSaidas(scriptPath, planoPath string, sc Script, plano string) error {
	dados, err := json.MarshalIndent(sc, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(scriptPath, append(dados, '\n'), 0o644); err != nil {
		return err
	}
	return os.WriteFile(planoPath, []byte(plano), 0o644)
}

func (p *Processo) dataAtual() string {
	if p.Agora != nil {
		return p.Agora().Format("2006-01-02")
	}
	return time.Now().Format("2006-01-02")
}
