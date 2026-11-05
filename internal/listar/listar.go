// Package listar mostra os processos que existem no cofre: nome, versão, natureza, harnesses vinculados e quem
// os usa. É uma fatia independente: depende só do núcleo (processo, ste) e nunca de outra função.
package listar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

func init() {
	ste.Default.Register("listar.vazio", "O cofre %s não tem processos. Para criar um processo, rode legoloop fabricar.")
	ste.Default.Register("listar.sem_cofre", "O cofre %s não existe. Informe o cofre com --vault.")
}

// NaturezaDesconhecida é a natureza dos processos gravados pela v1, que não têm o arquivo .processo.json.
const NaturezaDesconhecida = "desconhecida"

const (
	sufixoProcesso = ".processo.json"
	sufixoLoopV1   = ".loop.json"
	prefixo        = "legoloop-"
)

// ErrSemCofre: a pasta do cofre não existe.
var ErrSemCofre = errors.New("listar: cofre não encontrado")

// Deps são as bordas do listar. Store.Roots.Vault é o cofre; Store.FS nil usa o disco real.
type Deps struct {
	Store processo.Store
	Home  string // HOME onde os harnesses ficam (vazio = sem harnesses)
}

// Item é um processo do cofre.
type Item struct {
	Nome      string   `json:"nome"`
	Versao    string   `json:"versao"`
	Natureza  string   `json:"natureza"`
	Harnesses []string `json:"harnesses"`
	UsadoPor  []string `json:"usado_por"`
	Origem    string   `json:"origem"`         // "v2" (tem .processo.json) ou "v1" (só o .loop.json)
	Erro      string   `json:"erro,omitempty"` // o .processo.json não pôde ser lido
}

// Resultado é a lista ordenada por nome.
type Resultado struct {
	Vault     string `json:"vault"`
	Processos []Item `json:"processos"`
}

// MensagemVazio é o texto para o cofre sem processos.
func MensagemVazio(vault string) string { return ste.Msg("listar.vazio", vault) }

// Run lê loops/ do cofre e monta a lista. Um .processo.json inválido aparece com Erro e não para a lista.
func Run(ctx context.Context, d Deps) (Resultado, error) {
	vault := d.Store.Roots.Vault
	res := Resultado{Vault: vault, Processos: []Item{}}
	if err := ctx.Err(); err != nil {
		return res, err
	}
	if info, err := os.Stat(vault); err != nil || !info.IsDir() {
		return res, fmt.Errorf("%w: %s", ErrSemCofre, ste.Msg("listar.sem_cofre", vault))
	}
	entradas, err := os.ReadDir(filepath.Join(vault, "loops"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return res, err
	}
	compostos := map[string][]processo.Etapa{}
	for _, nome := range nomes(entradas) {
		if err := ctx.Err(); err != nil {
			return res, err
		}
		it, etapas := ler(vault, nome)
		compostos[nome] = etapas
		res.Processos = append(res.Processos, it)
	}
	completar(res.Processos, compostos, harnesses(d))
	return res, nil
}

// nomes devolve os processos de loops/, por nome e sem repetição: legoloop-*.processo.json (v2) ou, sem ele,
// legoloop-*.loop.json (v1). Outras esteiras não são processos fabricados.
func nomes(entradas []os.DirEntry) []string {
	vistos := map[string]bool{}
	for _, e := range entradas {
		n := e.Name()
		if e.IsDir() || !strings.HasPrefix(n, prefixo) {
			continue
		}
		for _, suf := range []string{sufixoProcesso, sufixoLoopV1} {
			if base, ok := strings.CutSuffix(n, suf); ok && base != strings.TrimSuffix(prefixo, "-") {
				vistos[base] = true
			}
		}
	}
	out := make([]string, 0, len(vistos))
	for n := range vistos {
		out = append(out, n)
	}
	slices.Sort(out)
	return out
}

// ler lê o processo da v2; sem o .processo.json, lê o .loop.json da v1 com natureza desconhecida.
func ler(vault, nome string) (Item, []processo.Etapa) {
	it := Item{Nome: nome, Harnesses: []string{}, UsadoPor: []string{}, Origem: "v2"}
	caminho := processo.CaminhoNoCofre(vault, nome)
	if _, err := os.Stat(caminho); err == nil {
		p, err := processo.Ler(caminho)
		if err != nil {
			it.Erro = err.Error()
			return it, nil
		}
		it.Versao, it.Natureza = p.Versao, p.Natureza
		return it, p.Etapas
	}
	it.Origem, it.Natureza = "v1", NaturezaDesconhecida
	it.Versao = versaoV1(filepath.Join(vault, "loops", nome+sufixoLoopV1))
	return it, nil
}

// versaoV1 lê o campo version do .loop.json da v1; sem ele, devolve "?".
func versaoV1(path string) string {
	var lp struct {
		Version string `json:"version"`
	}
	if b, err := os.ReadFile(path); err == nil && json.Unmarshal(b, &lp) == nil && lp.Version != "" {
		return lp.Version
	}
	return "?"
}

// harnesses devolve, por processo, os harnesses que têm a pasta skills/<nome> (vínculo ou cópia).
func harnesses(d Deps) func(nome string) []string {
	fsys := d.Store.FS
	if fsys == nil {
		fsys = processo.OSFS()
	}
	hs := processo.DetectarHarnesses(fsys, d.Home)
	return func(nome string) []string {
		out := []string{}
		for _, h := range hs {
			if _, err := fsys.Stat(filepath.Join(h.Skills, nome)); err == nil {
				out = append(out, h.Nome)
			}
		}
		return out
	}
}

// completar preenche harnesses e "usado por": quem tem uma etapa que compõe o processo.
func completar(its []Item, compostos map[string][]processo.Etapa, vinculados func(string) []string) {
	usadoPor := map[string][]string{}
	for nome, etapas := range compostos {
		for _, e := range etapas {
			if alvo, _, ok := strings.Cut(e.Compoe, "@"); ok && !slices.Contains(usadoPor[alvo], nome) {
				usadoPor[alvo] = append(usadoPor[alvo], nome)
			}
		}
	}
	for i := range its {
		its[i].Harnesses = vinculados(its[i].Nome)
		if u := usadoPor[its[i].Nome]; len(u) > 0 {
			slices.Sort(u)
			its[i].UsadoPor = u
		}
	}
}
