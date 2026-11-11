// Package listar mostra os processos que existem no cofre: nome, versão, natureza, harnesses vinculados e quem
// os usa. É uma fatia independente: depende só do núcleo (processo, ste) e nunca de outra função.
package listar

import (
	"context"
	"errors"
	"fmt"
	"os"

	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

func init() {
	ste.Default.Register("listar.vazio", "O cofre %s não tem processos. Para criar um processo, rode legoloop fabricar.")
	ste.Default.Register("listar.sem_cofre", "O cofre %s não existe. Informe o cofre com --vault.")
}

// NaturezaDesconhecida é a natureza dos processos gravados pela v1, que não têm o arquivo .processo.json.
const NaturezaDesconhecida = processo.NaturezaDesconhecida

// ErrSemCofre: a pasta do cofre não existe.
var ErrSemCofre = errors.New("listar: cofre não encontrado")

// Deps são as bordas do listar. Store.Roots.Vault é o cofre; Store.FS nil usa o disco real.
type Deps struct {
	Store processo.Store
	Home  string // HOME onde os harnesses ficam (vazio = sem harnesses)
}

// Item é um processo do cofre.
type Item struct {
	Nome      string     `json:"nome"`
	Versao    string     `json:"versao"`
	Natureza  string     `json:"natureza"`
	Harnesses []string   `json:"harnesses"`
	UsadoPor  []string   `json:"usado_por"`
	Origem    string     `json:"origem"`              // "v2" (tem .processo.json) ou "v1" (só o documento do Studio)
	Erro      string     `json:"erro,omitempty"`      // o .processo.json não pôde ser lido
	Lacunas   []string   `json:"lacunas,omitempty"`   // etapas sem peça
	Sugestoes []Sugestao `json:"sugestoes,omitempty"` // peças que podem cobrir a lacuna com perfil definir
}

// Resultado é a lista ordenada por nome.
type Resultado struct {
	Vault     string `json:"vault"`
	Processos []Item `json:"processos"`
}

// MensagemVazio é o texto para o cofre sem processos.
func MensagemVazio(vault string) string { return ste.Msg("listar.vazio", vault) }

// Run lê os processos do cofre e monta a lista. Um .processo.json inválido aparece com Erro e não para a lista.
func Run(ctx context.Context, d Deps) (Resultado, error) {
	vault := d.Store.Roots.Vault
	res := Resultado{Vault: vault, Processos: []Item{}}
	if err := ctx.Err(); err != nil {
		return res, err
	}
	if info, err := os.Stat(vault); err != nil || !info.IsDir() {
		return res, fmt.Errorf("%w: %s", ErrSemCofre, ste.Msg("listar.sem_cofre", vault))
	}
	procs, err := processo.ProcessosDoCofre(vault)
	if err != nil {
		return res, err
	}
	fsys := d.Store.FS
	if fsys == nil {
		fsys = processo.OSFS()
	}
	hs := processo.DetectarHarnesses(fsys, d.Home)
	usadoPor := processo.MapearUsadoPor(procs)
	compostos := make(map[string][]processo.Etapa, len(procs))
	for _, p := range procs {
		res.Processos = append(res.Processos, converterItem(p, fsys, hs, usadoPor[p.Nome]))
		compostos[p.Nome] = p.Etapas
	}
	sugerir(vault, res.Processos, compostos)
	return res, nil
}

func converterItem(p processo.ProcessoCofre, fsys processo.FS, hs []processo.Harness, usadoPor []string) Item {
	harnesses := processo.HarnessesVinculados(fsys, hs, p.Nome)
	if usadoPor == nil {
		usadoPor = []string{}
	}
	return Item{
		Nome:      p.Nome,
		Versao:    p.Versao,
		Natureza:  p.Natureza,
		Harnesses: harnesses,
		UsadoPor:  usadoPor,
		Origem:    p.Origem,
		Erro:      p.Erro,
	}
}
