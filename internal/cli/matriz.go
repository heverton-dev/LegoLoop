package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/brain"
	"legoloop/internal/interview"
	"legoloop/internal/matriz"
	"legoloop/internal/ste"
)

func init() {
	for id, msg := range map[string]string{
		"matriz.uso":            "Uso: legoloop matriz gerar | status | vocab (listar | pendentes | aprovar <termo>) | perfil (confirmar <peça> | definir <peça>) | ver.",
		"matriz.gerada":         "Matriz gerada: %d peças, %d recatalogadas, %d reaproveitadas, %d cópias.",
		"matriz.cerebro":        "Cérebro da catalogação: %s.",
		"matriz.sem_cerebro":    "Nenhum cérebro está em uso. O perfil vem só da declaração da peça e de palavras-chave.",
		"matriz.pendentes":      "%d termos novos esperam a sua aprovação. Veja com legoloop matriz vocab pendentes.",
		"matriz.erros":          "%d peças ficaram sem catalogar. Rode legoloop matriz gerar de novo.",
		"matriz.canonize_apply": "Matriz de peças: %d recatalogadas, %d reaproveitadas, %d cópias.",
		"matriz.canonize_dry":   "Matriz de peças: %d peças seriam recatalogadas.",
		"matriz.canonize_erro":  "A matriz de peças não foi gerada: %v.",
		"matriz.em_dia":         "A matriz de peças está em dia.",
		"matriz.vocab_aprovado": "Termo %s aprovado. O vocabulário está na versão %d.",
		"matriz.vocab_vazio":    "Nenhum termo espera aprovação.",
	} {
		ste.Default.Register(id, msg)
	}
}

// MatrizDeps injeta as bordas do comando matriz; o valor zero usa o ambiente real (usado em testes).
type MatrizDeps struct {
	Scan      func(ctx context.Context, home string) ([]assets.Item, error) // nil = assets.Scanner no HOME
	Brain     interview.LLMClient                                           // nil = brain.Pick no inventário
	BrainName string
	Now       func() time.Time
}

func (m *MatrizDeps) scan(ctx context.Context, home string) ([]assets.Item, error) {
	if m != nil && m.Scan != nil {
		return m.Scan(ctx, home)
	}
	return assets.NewScanner(assets.Config{Home: home}).Scan(ctx).Items, nil
}

func (d Deps) matrizHome() string {
	if d.Home != "" {
		return d.Home
	}
	h, _ := os.UserHomeDir()
	return h
}

func (d Deps) matrizDeps(vault string, cat interview.LLMClient) (matriz.Deps, error) {
	v, err := matriz.LoadVocab(vault)
	if err != nil {
		return matriz.Deps{}, err
	}
	md := matriz.Deps{Catalogar: matriz.Catalogador{Vocab: v, Brain: cat}, Home: d.matrizHome(),
		Scan: func(ctx context.Context) ([]assets.Item, error) { return d.Matriz.scan(ctx, d.matrizHome()) }}
	if d.Matriz != nil {
		md.Now = d.Matriz.Now
	}
	return md, nil
}

// cmdMatriz despacha gerar, status e vocab.
func cmdMatriz(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		return matrizUso(d)
	}
	switch args[0] {
	case "gerar":
		return cmdMatrizGerar(ctx, args[1:], d)
	case "status":
		return cmdMatrizStatus(ctx, args[1:], d)
	case "vocab":
		return cmdMatrizVocab(args[1:], d)
	case "perfil":
		return cmdMatrizPerfil(args[1:], d)
	case "ver":
		return cmdMatrizVer(args[1:], d)
	}
	return matrizUso(d)
}

func matrizUso(d Deps) int {
	fmt.Fprintln(d.Stderr, ste.Msg("matriz.uso"))
	return 2
}

func resolveVaultFlag(v string) (string, error) { return resolveVault(v) }

// escolherCerebro decide quem cataloga: nada (--sem-cerebro), o cérebro injetado ou o que o inventário oferece.
// O retorno code != 0 encerra o comando (modelo pedido fora do inventário = exit 2).
func escolherCerebro(ctx context.Context, d Deps, model string, sem bool) (c interview.LLMClient, nome string, code int) {
	if sem {
		return nil, "", 0
	}
	if d.Matriz != nil && d.Matriz.Brain != nil {
		return d.Matriz.Brain, d.Matriz.BrainName, 0
	}
	b, _, err := brain.Pick(d.discover(ctx), model)
	switch {
	case err == nil:
		return b, b.Name(), 0
	case errors.Is(err, brain.ErrModeloDesconhecido):
		fmt.Fprintln(d.Stderr, "erro:", err)
		return nil, "", 2
	}
	return nil, "", 0 // sem cérebro disponível: segue só com declaração e palavras-chave
}

func cmdMatrizGerar(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("matriz gerar", d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	model := fs.String("model", "", "modelo do cérebro da catalogação (padrão: o primeiro disponível)")
	sem := fs.Bool("sem-cerebro", false, "cataloga só por declaração e palavras-chave")
	forcar := fs.Bool("forcar", false, "recataloga todas as peças, mesmo sem mudança")
	asJSON := fs.Bool("json", false, "saída JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	v, err := resolveVaultFlag(*vault)
	if err != nil {
		return fail(d, err)
	}
	cat, nome, code := escolherCerebro(ctx, d, *model, *sem)
	if code != 0 {
		return code
	}
	info := d.Stdout
	if *asJSON {
		info = d.Stderr // o stdout do --json fica só com o JSON
	}
	if cat != nil {
		fmt.Fprintln(info, ste.Msg("matriz.cerebro", nome))
	} else {
		fmt.Fprintln(info, ste.Msg("matriz.sem_cerebro"))
	}
	sum, err := gerarMatriz(ctx, d, v, cat, *forcar)
	if err != nil {
		return fail(d, err)
	}
	return reportarMatriz(d, sum, *asJSON)
}

// gerarMatriz roda o motor e é usado por `matriz gerar` e pelo fim do `canonize --apply`.
func gerarMatriz(ctx context.Context, d Deps, vault string, cat interview.LLMClient, forcar bool) (*matriz.Summary, error) {
	md, err := d.matrizDeps(vault, cat)
	if err != nil {
		return nil, err
	}
	return matriz.Generate(ctx, matriz.Options{Vault: vault, Forcar: forcar}, md)
}

func reportarMatriz(d Deps, s *matriz.Summary, asJSON bool) int {
	if asJSON {
		return printJSON(d, s)
	}
	fmt.Fprintln(d.Stdout, ste.Msg("matriz.gerada", s.Total, s.Reprocessadas, s.Reaproveitadas, s.Copias))
	if s.Pendentes > 0 {
		fmt.Fprintln(d.Stdout, ste.Msg("matriz.pendentes", s.Pendentes))
	}
	if len(s.Erros) > 0 {
		fmt.Fprintln(d.Stdout, ste.Msg("matriz.erros", len(s.Erros)))
		for _, e := range s.Erros {
			fmt.Fprintln(d.Stdout, "  "+e)
		}
	}
	for _, a := range slices.Concat(s.Alertas, s.Invalidadas) {
		fmt.Fprintln(d.Stderr, "aviso:", a)
	}
	return 0
}

func lerMatrizOuNil(vault string) (*matriz.Registro, error) {
	reg, err := matriz.LerRegistro(filepath.Join(vault, "registry_matrix.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return reg, err
}

func cmdMatrizStatus(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("matriz status", d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	asJSON := fs.Bool("json", false, "saída JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	v, err := resolveVaultFlag(*vault)
	if err != nil {
		return fail(d, err)
	}
	reg, err := lerMatrizOuNil(v)
	if err != nil {
		return fail(d, err)
	}
	items, err := d.Matriz.scan(ctx, d.matrizHome())
	if err != nil {
		return fail(d, fmt.Errorf("a varredura falhou: %w", err))
	}
	md, err := d.matrizDeps(v, nil)
	if err != nil {
		return fail(d, err)
	}
	s, err := matriz.Staleness(ctx, reg, items, md)
	if err != nil {
		return fail(d, err)
	}
	return reportarStatus(d, s, *asJSON)
}

func reportarStatus(d Deps, s matriz.Stale, asJSON bool) int {
	code := 0
	if s.Defasada() {
		code = 1
	}
	if asJSON {
		if printJSON(d, s) != 0 {
			return 1
		}
		return code
	}
	if m := s.Mensagem(); m != "" {
		fmt.Fprintln(d.Stdout, m)
	} else {
		fmt.Fprintln(d.Stdout, ste.Msg("matriz.em_dia"))
	}
	return code
}

// matrizAposCanonize gera a matriz ao fim do `canonize --apply`, sem cérebro (o canonize não gasta modelo).
// Com --json o aviso vai ao stderr, para o stdout continuar sendo um JSON só.
func matrizAposCanonize(ctx context.Context, d Deps, vault string, asJSON bool) {
	out := d.Stdout
	if asJSON {
		out = d.Stderr
	}
	sum, err := gerarMatriz(ctx, d, vault, nil, false)
	if err != nil {
		fmt.Fprintln(d.Stderr, "aviso:", ste.Msg("matriz.canonize_erro", err))
		return
	}
	fmt.Fprintln(out, ste.Msg("matriz.canonize_apply", sum.Reprocessadas, sum.Reaproveitadas, sum.Copias))
	for _, a := range sum.Alertas {
		fmt.Fprintln(d.Stderr, "aviso:", a)
	}
}

// matrizPreviaCanonize diz, sem gravar nada, quantas peças o --apply recatalogaria.
func matrizPreviaCanonize(ctx context.Context, d Deps, vault string, asJSON bool) {
	out := d.Stdout
	if asJSON {
		out = d.Stderr
	}
	reg, err := lerMatrizOuNil(vault)
	if err == nil {
		var items []assets.Item
		if items, err = d.Matriz.scan(ctx, d.matrizHome()); err == nil {
			var md matriz.Deps
			if md, err = d.matrizDeps(vault, nil); err == nil {
				var n int
				if n, err = matriz.Previa(ctx, reg, items, md); err == nil {
					fmt.Fprintln(out, ste.Msg("matriz.canonize_dry", n))
					return
				}
			}
		}
	}
	fmt.Fprintln(d.Stderr, "aviso:", ste.Msg("matriz.canonize_erro", err))
}

func cmdMatrizVocab(args []string, d Deps) int {
	if len(args) == 0 {
		return matrizUso(d)
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "listar", "pendentes", "aprovar":
		return vocabSub(sub, rest, d)
	}
	return matrizUso(d)
}

func vocabSub(sub string, args []string, d Deps) int {
	termo, rest := splitPositional(args)
	fs := newFlags("matriz vocab "+sub, d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if termo == "" {
		termo = fs.Arg(0)
	}
	v, err := resolveVaultFlag(*vault)
	if err != nil {
		return fail(d, err)
	}
	switch sub {
	case "listar":
		return vocabListar(d, v)
	case "pendentes":
		return vocabPendentes(d, v)
	}
	if termo == "" {
		return matrizUso(d)
	}
	return vocabAprovar(d, v, termo)
}

func vocabListar(d Deps, vault string) int {
	voc, err := matriz.LoadVocab(vault)
	if err != nil {
		return fail(d, err)
	}
	fmt.Fprintf(d.Stdout, "Vocabulário, versão %d\n", voc.Versao)
	for _, c := range []matriz.Categoria{matriz.CatEtapa, matriz.CatNatureza, matriz.CatCasoDeUso} {
		nomes := []string{}
		for _, t := range voc.Termos(c) {
			nomes = append(nomes, t.Nome)
		}
		fmt.Fprintf(d.Stdout, "%s: %s\n", c, strings.Join(nomes, ", "))
	}
	return 0
}

func vocabPendentes(d Deps, vault string) int {
	fila, err := matriz.LerPendentes(vault)
	if err != nil {
		return fail(d, err)
	}
	if len(fila) == 0 {
		fmt.Fprintln(d.Stdout, ste.Msg("matriz.vocab_vazio"))
		return 0
	}
	w := table(d.Stdout)
	fmt.Fprintln(w, "TERMO\tCATEGORIA\tPEÇAS")
	for _, p := range fila {
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.Termo, p.Categoria, strings.Join(p.Pecas, ", "))
	}
	if err := w.Flush(); err != nil {
		return fail(d, err)
	}
	return 0
}

func vocabAprovar(d Deps, vault, termo string) int {
	if err := matriz.Aprovar(vault, termo); err != nil {
		return fail(d, err)
	}
	voc, err := matriz.LoadVocab(vault)
	if err != nil {
		return fail(d, err)
	}
	fmt.Fprintln(d.Stdout, ste.Msg("matriz.vocab_aprovado", termo, voc.Versao))
	return 0
}
