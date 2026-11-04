// Package processo é o núcleo compartilhado das funções do LegoLoop que lidam com processos fabricados
// (fabricar, listar, refinar, remover, compor): o modelo do processo, a validação e a gravação atômica no
// cofre. Não tem regra de negócio de nenhuma função, e as funções nunca se importam entre si.
package processo

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"legoloop/internal/gate"
)

var (
	// ErrInvalido: o contrato do processo não é satisfeito.
	ErrInvalido = errors.New("processo: contrato inválido")
	// ErrCiclo: as etapas (ou, no futuro, os processos compostos) formam um ciclo.
	ErrCiclo = errors.New("processo: composição com ciclo")
)

var (
	reNome   = regexp.MustCompile(`^legoloop-[a-z0-9]+(-[a-z0-9]+)*$`)
	reSemver = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
)

// Artefato é o que uma etapa produz.
type Artefato struct {
	Caminho string `json:"caminho"`
	Tipo    string `json:"tipo"`    // markdown | json | codigo | dados | marcador
	Promete bool   `json:"promete"` // true = promete conteúdo; gate só de existência não basta
}

// Etapa é a unidade única do processo. SPEC, PLAN, ROUTE, DAG, loop e skill são projeções dela.
type Etapa struct {
	ID      string   `json:"id"`
	Titulo  string   `json:"titulo"`
	Tipo    string   `json:"tipo"`             // etapa do vocabulário (ex.: "construir")
	Peca    string   `json:"peca"`             // "" = lacuna declarada
	Compoe  string   `json:"compoe,omitempty"` // "nome@versao" se a etapa é outro processo
	Consome []string `json:"consome"`
	Produz  Artefato `json:"produz"`
	Gate    string   `json:"gate"` // expressão da DSL de internal/gate
	Depende []string `json:"depende"`
}

// Processo é o que o cofre guarda: durável e reutilizável. Cada execução dele é efêmera.
type Processo struct {
	Nome      string    `json:"nome"`             // legoloop-<kebab>
	Versao    string    `json:"versao"`           // semver
	Natureza  string    `json:"natureza"`         // vocabulário; "desconhecida" nos processos da v1
	Rotulo    string    `json:"rotulo,omitempty"` // antiga categoria, opcional
	Objetivo  string    `json:"objetivo"`
	Etapas    []Etapa   `json:"etapas"`
	CriadoEm  time.Time `json:"criado_em"`
	Historico []string  `json:"historico,omitempty"` // versões anteriores guardadas
}

func invalido(format string, a ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{ErrInvalido}, a...)...)
}

// Validate confere o contrato antes de qualquer E/S.
func (p Processo) Validate() error {
	switch {
	case !reNome.MatchString(p.Nome):
		return invalido("nome %q (use legoloop-<kebab-case>)", p.Nome)
	case !reSemver.MatchString(p.Versao):
		return invalido("versão %q (use maior.menor.correção)", p.Versao)
	case strings.TrimSpace(p.Natureza) == "":
		return invalido("natureza vazia")
	case strings.TrimSpace(p.Objetivo) == "":
		return invalido("objetivo vazio")
	case len(p.Etapas) == 0:
		return invalido("o processo precisa de ao menos uma etapa")
	}
	ids, err := indexar(p.Etapas)
	if err != nil {
		return err
	}
	for _, e := range p.Etapas {
		if err := e.validate(ids); err != nil {
			return err
		}
	}
	return grafo(p.Etapas)
}

// indexar garante IDs não vazios e únicos.
func indexar(es []Etapa) (map[string]bool, error) {
	ids := make(map[string]bool, len(es))
	for _, e := range es {
		if strings.TrimSpace(e.ID) == "" {
			return nil, invalido("etapa sem id")
		}
		if ids[e.ID] {
			return nil, invalido("id de etapa duplicado %q", e.ID)
		}
		ids[e.ID] = true
	}
	return ids, nil
}

func (e Etapa) validate(ids map[string]bool) error {
	for _, d := range e.Depende {
		if !ids[d] {
			return invalido("etapa %q depende de %q, que não existe", e.ID, d)
		}
	}
	if strings.TrimSpace(e.Gate) == "" {
		return invalido("etapa %q sem gate", e.ID)
	}
	if _, err := gate.Compile(e.Gate); err != nil {
		return invalido("etapa %q: o gate não compila: %v", e.ID, err)
	}
	c := e.Produz.Caminho
	if c == "" || strings.HasPrefix(c, "/") || !filepath.IsLocal(filepath.FromSlash(c)) {
		return invalido("etapa %q: o artefato %q precisa de caminho relativo dentro do projeto", e.ID, c)
	}
	return nil
}

// grafo recusa ciclos e etapas sem conexão (quando há mais de uma etapa, toda etapa precisa ligar-se a outra).
func grafo(es []Etapa) error {
	deps, usada := map[string][]string{}, map[string]bool{}
	for _, e := range es {
		deps[e.ID] = e.Depende
		for _, d := range e.Depende {
			usada[d] = true
		}
	}
	if len(es) > 1 {
		for _, e := range es {
			if len(e.Depende) == 0 && !usada[e.ID] {
				return invalido("etapa %q sem conexão com as outras", e.ID)
			}
		}
	}
	return semCiclo(deps)
}

func semCiclo(deps map[string][]string) error {
	estado := map[string]int{} // 1 = visitando, 2 = pronto
	var visita func(string) error
	visita = func(n string) error {
		switch estado[n] {
		case 1:
			return fmt.Errorf("%w: %w: a etapa %q volta a si mesma", ErrInvalido, ErrCiclo, n)
		case 2:
			return nil
		}
		estado[n] = 1
		for _, d := range deps[n] {
			if err := visita(d); err != nil {
				return err
			}
		}
		estado[n] = 2
		return nil
	}
	for n := range deps {
		if err := visita(n); err != nil {
			return err
		}
	}
	return nil
}
