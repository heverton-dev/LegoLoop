package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"legoloop/internal/matriz"
	"legoloop/internal/reputation"
	"legoloop/internal/ste"
	"legoloop/internal/taxonomy"
)

func init() {
	ste.Default.Register("matriz.matrix_descontinuado", "O comando matrix foi descontinuado. Use legoloop matriz ver.")
}

var (
	tiposVer   = []string{"skill", "mcp", "script", "hook", "gate", "processo"}
	tiersVer   = []string{"ouro", "prata", "bronze", "toxica"}
	origensVer = []string{"declarado", "inferido", "confirmado", "nenhuma"}
)

type verOpts struct {
	vault, registro, tipo, tier, natureza, etapa, caso, origem string
	papel, guilda                                              string
	min                                                        float64
	copias, json                                               bool
}

type papeisInfo struct {
	PapelPrimario     string   `json:"papel_primario"`
	PapeisSecundarios []string `json:"papeis_secundarios"`
}

// linhaVer é uma linha da tabela: a peça e as cópias idênticas que ela representa.
type linhaVer struct {
	Entrada matriz.Entrada
	Refs    []string
}

// entradaVer é a forma da linha no --json.
type entradaVer struct {
	Ref               string              `json:"ref"`
	Nome              string              `json:"nome"`
	Tipo              string              `json:"tipo"`
	Tier              string              `json:"tier"`
	Score             float64             `json:"total_score"`
	Perfil            *matriz.PerfilDeUso `json:"perfil,omitempty"`
	Origem            string              `json:"origem"`
	Copias            int                 `json:"copias"`
	Refs              []string            `json:"refs"`
	PapelPrimario     string              `json:"papel_primario,omitempty"`
	PapeisSecundarios []string            `json:"papeis_secundarios,omitempty"`
}

// cmdMatrizVer mostra a matriz de peças: nota, tipo e perfil de uso, com as cópias idênticas numa linha só.
func cmdMatrizVer(args []string, d Deps) int {
	o, code, ok := parseVer(args, d)
	if !ok {
		return code
	}
	reg, err := carregarParaVer(o.registro)
	if err != nil {
		return fail(d, err)
	}
	papeis := carregarMapaPapeis(o.registro)
	var passam []matriz.Entrada
	for _, e := range reg.Entradas {
		if passaVer(e, o, papeis) {
			passam = append(passam, e)
		}
	}
	return imprimirVer(d, agruparVer(passam, o.copias), o, papeis)
}

func configurarFlagsVer(o *verOpts, fs *flag.FlagSet) {
	fs.StringVar(&o.vault, "vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	fs.StringVar(&o.registro, "registro", "", "caminho de um registry_matrix.json (padrão: o do cofre)")
	fs.StringVar(&o.tipo, "tipo", "", "skill, mcp, script, hook, gate ou processo")
	fs.StringVar(&o.tier, "tier", "", "ouro, prata, bronze ou toxica")
	fs.Float64Var(&o.min, "min", 0, "nota mínima, de 0 a 1")
	fs.StringVar(&o.natureza, "natureza", "", "natureza do vocabulário")
	fs.StringVar(&o.etapa, "etapa", "", "etapa do vocabulário")
	fs.StringVar(&o.caso, "caso", "", "caso de uso do vocabulário")
	fs.StringVar(&o.origem, "origem", "", "declarado, inferido, confirmado ou nenhuma")
	fs.StringVar(&o.papel, "papel", "", "papel operacional da taxonomia")
	fs.StringVar(&o.guilda, "guilda", "", "guilda operacional da taxonomia")
	fs.BoolVar(&o.copias, "copias", false, "uma linha por caminho, sem agrupar as cópias")
	fs.BoolVar(&o.json, "json", false, "saída JSON")
}

func parseVer(args []string, d Deps) (o verOpts, code int, ok bool) {
	fs := newFlags("matriz ver", d)
	configurarFlagsVer(&o, fs)
	if err := fs.Parse(args); err != nil {
		return o, 2, false
	}
	if fs.NArg() > 0 {
		return o, matrizUso(d), false
	}
	v, err := resolveVaultFlag(o.vault)
	if err != nil {
		return o, fail(d, err), false
	}
	if o.registro == "" {
		o.registro = filepath.Join(v, "registry_matrix.json")
	}
	o.tipo, o.tier, o.origem = strings.ToLower(o.tipo), strings.ToLower(o.tier), strings.ToLower(o.origem)
	o.papel, o.guilda = strings.ToLower(strings.TrimSpace(o.papel)), strings.ToLower(strings.TrimSpace(o.guilda))
	voc, err := matriz.LoadVocab(v)
	if err == nil {
		err = validarVer(o, voc)
	}
	if err != nil {
		fmt.Fprintln(d.Stderr, "erro:", err)
		code := 2
		if errors.Is(err, taxonomy.ErrPapelInvalido) || errors.Is(err, taxonomy.ErrGuildaInvalida) {
			code = 1
		}
		return o, code, false
	}
	return o, 0, true
}

func validarVer(o verOpts, voc matriz.Vocab) error {
	for _, c := range []struct {
		flag, val string
		aceitos   []string
	}{{"tipo", o.tipo, tiposVer}, {"tier", o.tier, tiersVer}, {"origem", o.origem, origensVer}} {
		if c.val != "" && !slices.Contains(c.aceitos, c.val) {
			return fmt.Errorf("--%s %q não existe. Valores aceitos: %s", c.flag, c.val, strings.Join(c.aceitos, ", "))
		}
	}
	for _, c := range []struct {
		flag, val string
		cat       matriz.Categoria
	}{{"natureza", o.natureza, matriz.CatNatureza}, {"etapa", o.etapa, matriz.CatEtapa}, {"caso", o.caso, matriz.CatCasoDeUso}} {
		if c.val != "" && !voc.Has(c.cat, c.val) {
			return fmt.Errorf("--%s %q não está no vocabulário. Veja os termos com legoloop matriz vocab listar", c.flag, c.val)
		}
	}
	if o.min < 0 || o.min > 1 {
		return errors.New("--min precisa estar entre 0 e 1")
	}
	if o.papel != "" {
		if err := taxonomy.ValidarPapel(o.papel); err != nil {
			return fmt.Errorf("--papel %q inválido ou fora do catálogo canônico: %w", o.papel, err)
		}
	}
	if o.guilda != "" {
		if err := taxonomy.ValidarGuilda(taxonomy.Guilda(o.guilda)); err != nil {
			return fmt.Errorf("--guilda %q inválida ou fora do catálogo canônico: %w", o.guilda, err)
		}
	}
	return nil
}

// carregarParaVer lê a matriz. A da v1 não tem `entradas`: as notas viram entradas sem perfil, com o tipo
// deduzido do caminho (como o comando matrix fazia).
func carregarParaVer(path string) (*matriz.Registro, error) {
	reg, err := matriz.LerRegistro(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("a matriz não existe em %s. Rode legoloop matriz gerar", path)
	}
	if err != nil {
		return nil, err
	}
	if len(reg.Entradas) == 0 {
		for _, r := range reg.Ratings {
			reg.Entradas = append(reg.Entradas, matriz.Entrada{Ref: r.Ref, Nome: r.Name, Tipo: tipoPorRef(r.Ref), Rating: r})
		}
	}
	return reg, nil
}

func tipoPorRef(ref string) string {
	r := strings.ToLower(ref)
	for _, t := range []struct {
		tipo string
		cond bool
	}{
		{"script", strings.HasSuffix(r, ".ps1") || strings.HasSuffix(r, ".sh") || strings.HasSuffix(r, ".bat") || strings.HasSuffix(r, ".cmd") || strings.Contains(r, "scripts")},
		{"skill", strings.Contains(r, "skill")}, {"hook", strings.Contains(r, "hook")}, {"mcp", strings.Contains(r, "mcp")},
		{"gate", strings.Contains(r, "gate")}, {"processo", strings.Contains(r, "processo") || strings.HasSuffix(r, ".dag.json")},
	} {
		if t.cond {
			return t.tipo
		}
	}
	return ""
}

func origemDaEntrada(e matriz.Entrada) string {
	if e.Perfil == nil {
		return "nenhuma"
	}
	return e.Perfil.Origem
}

func carregarMapaPapeis(path string) map[string]papeisInfo {
	res := make(map[string]papeisInfo)
	b, err := os.ReadFile(path)
	if err != nil {
		return res
	}
	var doc struct {
		Entradas []struct {
			Ref               string         `json:"ref"`
			Nome              string         `json:"nome"`
			PapelPrimario     string         `json:"papel_primario"`
			PapeisSecundarios []string       `json:"papeis_secundarios"`
			Tipado            map[string]any `json:"tipado"`
		} `json:"entradas"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		return res
	}
	for _, it := range doc.Entradas {
		info := mapearPapeisEntrada(it.PapelPrimario, it.PapeisSecundarios, it.Tipado)
		if it.Ref != "" {
			res[it.Ref] = info
		}
		if it.Nome != "" {
			if _, exists := res[it.Nome]; !exists {
				res[it.Nome] = info
			}
		}
	}
	return res
}

func mapearPapeisEntrada(prim string, secs []string, tipado map[string]any) papeisInfo {
	info := papeisInfo{PapelPrimario: prim, PapeisSecundarios: secs}
	if info.PapelPrimario == "" && tipado != nil {
		if p, ok := tipado["papel_primario"].(string); ok {
			info.PapelPrimario = p
		}
		if sList, ok := tipado["papeis_secundarios"].([]any); ok {
			for _, item := range sList {
				if s, ok := item.(string); ok {
					info.PapeisSecundarios = append(info.PapeisSecundarios, s)
				}
			}
		}
	}
	return info
}

func obterPapeisEntrada(e matriz.Entrada, mapa map[string]papeisInfo) papeisInfo {
	if info, ok := mapa[e.Ref]; ok && (info.PapelPrimario != "" || len(info.PapeisSecundarios) > 0) {
		return info
	}
	if info, ok := mapa[e.Nome]; ok && (info.PapelPrimario != "" || len(info.PapeisSecundarios) > 0) {
		return info
	}
	if e.Tipado != nil {
		if info := mapearPapeisEntrada("", nil, e.Tipado); info.PapelPrimario != "" || len(info.PapeisSecundarios) > 0 {
			return info
		}
	}
	if e.Ref != "" {
		if b, err := os.ReadFile(e.Ref); err == nil {
			prim, sec := matriz.ExtrairPapeis(string(b))
			if prim != "" || len(sec) > 0 {
				return papeisInfo{PapelPrimario: prim, PapeisSecundarios: sec}
			}
		}
	}
	return papeisInfo{}
}

func papelPertenceAGuilda(slug, guilda string) bool {
	if slug == "" {
		return false
	}
	p, err := taxonomy.ObterPapel(slug)
	if err != nil {
		return false
	}
	return string(p.Guilda) == guilda
}

func passaPapeis(e matriz.Entrada, o verOpts, mapa map[string]papeisInfo) bool {
	info := obterPapeisEntrada(e, mapa)
	if o.papel != "" {
		bateu := info.PapelPrimario == o.papel || slices.Contains(info.PapeisSecundarios, o.papel)
		if !bateu {
			return false
		}
	}
	if o.guilda != "" {
		bateuGuilda := papelPertenceAGuilda(info.PapelPrimario, o.guilda)
		if !bateuGuilda {
			for _, sec := range info.PapeisSecundarios {
				if papelPertenceAGuilda(sec, o.guilda) {
					bateuGuilda = true
					break
				}
			}
		}
		if !bateuGuilda {
			return false
		}
	}
	return true
}

func passaVer(e matriz.Entrada, o verOpts, mapa map[string]papeisInfo) bool {
	switch {
	case o.tipo != "" && e.Tipo != o.tipo,
		o.tier != "" && string(e.Rating.Tier) != o.tier,
		o.min > 0 && e.Rating.TotalScore < o.min,
		o.origem != "" && origemDaEntrada(e) != o.origem:
		return false
	}
	if (o.papel != "" || o.guilda != "") && !passaPapeis(e, o, mapa) {
		return false
	}
	if o.natureza == "" && o.etapa == "" && o.caso == "" {
		return true
	}
	p := e.Perfil
	return p != nil && (o.natureza == "" || slices.Contains(p.Naturezas, o.natureza)) &&
		(o.etapa == "" || slices.Contains(p.Etapas, o.etapa)) && (o.caso == "" || slices.Contains(p.CasosDeUso, o.caso))
}

// agruparVer junta as cópias idênticas (mesmo tipo e hash) numa linha, na ordem da matriz. Com copias, não agrupa.
func agruparVer(es []matriz.Entrada, copias bool) []linhaVer {
	var out []linhaVer
	pos := map[string]int{}
	for _, e := range es {
		k := e.Tipo + "|" + e.Hash
		i, achou := pos[k]
		if copias || e.Hash == "" || !achou {
			pos[k] = len(out)
			out = append(out, linhaVer{Entrada: e, Refs: []string{e.Ref}})
			continue
		}
		out[i].Refs = append(out[i].Refs, e.Ref)
	}
	return out
}

func imprimirVer(d Deps, ls []linhaVer, o verOpts, mapa map[string]papeisInfo) int {
	if o.json {
		doc := struct {
			Total    int                      `json:"total_found"`
			Entradas []entradaVer             `json:"entradas"`
			Ratings  []reputation.PieceRating `json:"ratings"`
		}{Total: len(ls), Entradas: []entradaVer{}, Ratings: []reputation.PieceRating{}}
		for _, l := range ls {
			e := l.Entrada
			info := obterPapeisEntrada(e, mapa)
			doc.Entradas = append(doc.Entradas, entradaVer{
				Ref:               e.Ref,
				Nome:              e.Nome,
				Tipo:              e.Tipo,
				Tier:              string(e.Rating.Tier),
				Score:             e.Rating.TotalScore,
				Perfil:            e.Perfil,
				Origem:            origemDaEntrada(e),
				Copias:            len(l.Refs),
				Refs:              l.Refs,
				PapelPrimario:     info.PapelPrimario,
				PapeisSecundarios: info.PapeisSecundarios,
			})
			doc.Ratings = append(doc.Ratings, e.Rating)
		}
		return printJSON(d, doc)
	}
	return tabelaVer(d, ls, o.copias)
}

func tabelaVer(d Deps, ls []linhaVer, porCaminho bool) int {
	ultima := "CÓPIAS"
	if porCaminho {
		ultima = "CAMINHO"
	}
	fmt.Fprintf(d.Stdout, "\n=== MATRIZ DE PEÇAS (%d ENCONTRADAS) ===\n\n", len(ls))
	w := table(d.Stdout)
	fmt.Fprintf(w, "TIER\tSCORE\tNOME\tTIPO\tNATUREZAS\tETAPAS\tCASOS\tORIGEM\t%s\n", ultima)
	fmt.Fprintln(w, "----\t-----\t----\t----\t---------\t------\t-----\t------\t------")
	for _, l := range ls {
		e := l.Entrada
		fim := fmt.Sprint(len(l.Refs))
		if porCaminho {
			fim = e.Ref
		}
		fmt.Fprintf(w, "%s\t%.2f\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", strings.ToUpper(string(e.Rating.Tier)), e.Rating.TotalScore,
			e.Nome, e.Tipo, termosVer(e, func(p *matriz.PerfilDeUso) []string { return p.Naturezas }),
			termosVer(e, func(p *matriz.PerfilDeUso) []string { return p.Etapas }),
			termosVer(e, func(p *matriz.PerfilDeUso) []string { return p.CasosDeUso }), origemDaEntrada(e), fim)
	}
	if err := w.Flush(); err != nil {
		return fail(d, err)
	}
	return 0
}

func termosVer(e matriz.Entrada, campo func(*matriz.PerfilDeUso) []string) string {
	if e.Perfil == nil || len(campo(e.Perfil)) == 0 {
		return "-"
	}
	return strings.Join(campo(e.Perfil), ",")
}
