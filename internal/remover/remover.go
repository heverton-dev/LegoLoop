// Package remover apaga um processo do cofre e os vínculos dele nos harnesses. Um processo que outro processo
// compõe não sai: a função falha com ErrEmUso e nada muda. É uma fatia independente do Épico 03: depende apenas do
// núcleo (processo, ste).
package remover

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

const prefixo = "legoloop-"

func init() {
	for id, msg := range map[string]string{
		"remover.em_uso":        "O processo %s não sai: o processo %s usa ele. Remova primeiro quem usa.",
		"remover.nao_achou":     "O cofre %s não tem o processo %s. Veja os processos com legoloop listar.",
		"remover.nome_invalido": "O nome %q não é de um processo. O nome começa com legoloop- e não tem barra.",
		"remover.cancelado":     "Remoção cancelada. Nada mudou.",
		"remover.feito":         "Processo %s removido. Itens apagados: %d. Vínculos desfeitos: %d.",
	} {
		ste.Default.Register(id, msg)
	}
}

var (
	// ErrEmUso: outro processo compõe o processo pedido. Nada é removido.
	ErrEmUso = errors.New("remover: processo em uso")
	// ErrNaoEncontrado: o cofre não tem o processo.
	ErrNaoEncontrado = errors.New("remover: processo não encontrado")
	// ErrNomeInvalido: o nome não é de um processo do cofre (prefixo ausente ou caminho no nome).
	ErrNomeInvalido = errors.New("remover: nome inválido")
)

// Deps são as bordas do remover. Store.Roots.Vault é o cofre; Home é onde ficam os harnesses.
type Deps struct {
	Store processo.Store
	Home  string
}

// Plano é o que Executar vai apagar. Planejar não toca o disco.
type Plano struct {
	Nome      string     `json:"nome"`
	Versao    string     `json:"versao"`
	Harnesses []string   `json:"harnesses"`
	Links     []string   `json:"links"`    // pastas em ~/.<harness>/skills
	Arquivos  []string   `json:"arquivos"` // itens do cofre, com o processo.json por último
	Cartao    ste.Cartao `json:"cartao"`
}

// MensagemCancelado é o texto de quando o usuário não confirma.
func MensagemCancelado() string { return ste.Msg("remover.cancelado") }

// MensagemFeito é o texto de quando a remoção termina.
func MensagemFeito(p Plano) string {
	return ste.Msg("remover.feito", p.Nome, len(p.Arquivos), len(p.Links))
}

// Planejar confere o nome, acha o processo e recusa com ErrEmUso se outro processo o compõe.
func Planejar(ctx context.Context, d Deps, nome string) (Plano, error) {
	if !nomeValido(nome) {
		return Plano{}, fmt.Errorf("%w: %s", ErrNomeInvalido, ste.Msg("remover.nome_invalido", nome))
	}
	if err := ctx.Err(); err != nil {
		return Plano{}, err
	}
	vault := d.Store.Roots.Vault
	if info, err := os.Stat(vault); err != nil || !info.IsDir() {
		return Plano{}, fmt.Errorf("%w: %s", ErrNaoEncontrado, ste.Msg("remover.nao_achou", vault, nome))
	}
	procs, err := processo.ProcessosDoCofre(vault)
	if err != nil {
		return Plano{}, err
	}
	alvo := acharProcesso(procs, nome)
	if alvo == nil {
		return Plano{}, fmt.Errorf("%w: %s", ErrNaoEncontrado, ste.Msg("remover.nao_achou", vault, nome))
	}
	usadoPor := processo.MapearUsadoPor(procs)
	if quem := usadoPor[nome]; len(quem) > 0 {
		return Plano{}, fmt.Errorf("%w: %s", ErrEmUso, ste.Msg("remover.em_uso", nome, strings.Join(quem, ", ")))
	}
	return montar(d, *alvo), nil
}

func acharProcesso(procs []processo.ProcessoCofre, nome string) *processo.ProcessoCofre {
	for i := range procs {
		if procs[i].Nome == nome {
			return &procs[i]
		}
	}
	return nil
}

// nomeValido aceita só legoloop-<algo> sem separador de caminho: o nome vira caminho no disco.
func nomeValido(nome string) bool {
	return strings.HasPrefix(nome, prefixo) && len(nome) > len(prefixo) &&
		!strings.ContainsAny(nome, `/\`)
}

func montar(d Deps, it processo.ProcessoCofre) Plano {
	vault := d.Store.Roots.Vault
	hs := processo.DetectarHarnesses(fsDe(d), d.Home)
	vinculados := processo.HarnessesVinculados(fsDe(d), hs, it.Nome)
	p := Plano{Nome: it.Nome, Versao: it.Versao, Harnesses: vinculados, Links: []string{}, Arquivos: []string{}}
	for _, h := range hs {
		if slices.Contains(vinculados, h.Nome) {
			p.Links = append(p.Links, filepath.Join(h.Skills, it.Nome))
		}
	}
	for _, c := range caminhosCofre(vault, it.Nome) {
		if _, err := os.Lstat(c); err == nil {
			p.Arquivos = append(p.Arquivos, c)
		}
	}
	p.Cartao = cartao(p)
	return p
}

func caminhosCofre(vault, nome string) []string {
	return []string{
		filepath.Join(vault, "skills", nome),
		filepath.Join(vault, "loops", nome+".loop.json"),
		filepath.Join(vault, "loops", nome+".dag.json"),
		processo.CaminhoNoCofre(vault, nome),
	}
}

func fsDe(d Deps) processo.FS {
	if d.Store.FS != nil {
		return d.Store.FS
	}
	return processo.OSFS()
}

// Executar apaga os vínculos e, depois, os itens do cofre. O processo.json sai por último: se algo falhar no
// meio, o processo ainda existe e o usuário pode rodar remover de novo. RemoveAll em symlink apaga só o link.
func Executar(ctx context.Context, p Plano) error {
	for _, caminho := range append(append([]string{}, p.Links...), p.Arquivos...) {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := os.RemoveAll(caminho); err != nil {
			return fmt.Errorf("remover %s: %w", caminho, err)
		}
	}
	return nil
}
