package render

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// preparoRender reune o que o render precisa depois de validar todas as entradas.
type preparoRender struct {
	deps       RenderDeps
	cfg        ConfigRender
	tl         Timeline
	resumo     resumoEntradas
	imagens    []string
	soundtrack string
	rendersDir string
}

func preparar(projDir string, c ContratosRender, deps RenderDeps) (preparoRender, error) {
	p := preparoRender{deps: deps, cfg: c.Config}
	var err error
	if p.tl, err = lerTimelineSalva(projDir); err != nil {
		return p, err
	}
	if p.resumo, err = lerResumo(projDir); err != nil {
		return p, err
	}
	if p.soundtrack, err = existeArquivo(projDir, "audio/soundtrack.wav"); err != nil {
		return p, err
	}
	plano, err := montarFiltroRender(p.tl, resolucaoDe(ResolucaoPadrao), "")
	if err != nil {
		return p, err
	}
	if p.imagens, err = imagensAbsolutas(projDir, plano.Entradas); err != nil {
		return p, err
	}
	p.rendersDir, err = resolverSeguro(projDir, "renders")
	return p, err
}

// ExecutarRender e a Etapa 8b (master 16:9): uma saida por resolucao pedida, a primeira como vN.mp4 e as
// demais como vN_<res>.mp4. Nada e gravado em definitivo se alguma saida reprovar no gate.
func ExecutarRender(ctx context.Context, projDir string, deps RenderDeps) (RelatorioRender, error) {
	if err := ctx.Err(); err != nil {
		return RelatorioRender{}, err
	}
	c, err := CarregarContratos(projDir)
	if err != nil {
		return RelatorioRender{}, err
	}
	p, err := preparar(projDir, c, normalizarRenderDeps(deps))
	if err != nil {
		return RelatorioRender{}, err
	}
	return p.executar(ctx)
}

type saidaRender struct{ parcial, final, res string }

func (p preparoRender) renderUma(ctx context.Context, s saidaRender, ass, enc string) (string, error) {
	plano, err := montarFiltroRender(p.tl, resolucaoDe(s.res), ass)
	if err != nil {
		return enc, err
	}
	usado, err := ExecutarComFallback(ctx, p.deps.Exec, enc, func(e string) []string {
		return argsRender(p.imagens, p.soundtrack, plano.Filtro, s.parcial, e)
	})
	if err != nil {
		return usado, err
	}
	return usado, ConferirRender(ctx, p.deps.Med, s.parcial)
}

func (p preparoRender) planejar(n int) []saidaRender {
	var saidas []saidaRender
	for i, res := range p.cfg.Resolucoes {
		saidas = append(saidas, saidaRender{
			parcial: filepath.Join(p.rendersDir, fmt.Sprintf(".v%d_%s.parcial.mp4", n, res)),
			final:   filepath.Join(p.rendersDir, nomeDaSaida(n, res, i == 0)),
			res:     res,
		})
	}
	return saidas
}

func removerTodos(caminhos ...string) error {
	var falhas []error
	for _, c := range caminhos {
		falhas = append(falhas, removerTemp(c))
	}
	return errors.Join(falhas...)
}

func (p preparoRender) executar(ctx context.Context) (RelatorioRender, error) {
	rel := RelatorioRender{ProvedorVoz: p.resumo.provedor, LUFS: p.resumo.lufs, TruePeakDBTP: p.resumo.tp, Renders: map[string]string{}}
	if err := os.MkdirAll(p.rendersDir, 0o755); err != nil {
		return rel, err
	}
	n, err := proximaVersao(p.rendersDir)
	if err != nil {
		return rel, err
	}
	assParcial := filepath.Join(p.rendersDir, ".legenda.parcial.ass")
	if err := os.WriteFile(assParcial, []byte(GerarASS(p.tl)), 0o644); err != nil {
		return rel, err
	}
	saidas := p.planejar(n)
	pendentes := []string{assParcial}
	for _, s := range saidas {
		pendentes = append(pendentes, s.parcial)
	}
	enc := EscolherEncoder(ctx, p.deps.Sonda)
	for _, s := range saidas {
		if enc, err = p.renderUma(ctx, s, assParcial, enc); err != nil {
			return rel, errors.Join(err, removerTodos(pendentes...))
		}
		rel.Renders[s.res] = "renders/" + filepath.Base(s.final)
	}
	rel.Encoder = enc
	return rel, p.publicar(saidas, assParcial, pendentes)
}

// publicar renomeia os parciais aprovados para os nomes finais.
func (p preparoRender) publicar(saidas []saidaRender, assParcial string, pendentes []string) error {
	for _, s := range saidas {
		if err := os.Rename(s.parcial, s.final); err != nil {
			return errors.Join(err, removerTodos(pendentes...))
		}
	}
	return os.Rename(assParcial, filepath.Join(p.rendersDir, "legenda.ass"))
}
