// Package assets inventaria o que o usuário tem na máquina e que pode virar peça, slot ou
// ferramenta de um Processo: skills, MCPs, scripts, hooks, gates, harnesses e binários.
//
// Regras: só biblioteca padrão; leitura apenas (nada é executado nem alterado); segredos
// (env, args, headers, query de URL) nunca entram no inventário.
package assets

import (
	"os"
	"os/exec"
	"time"
)

// Kind é a categoria de um ativo.
type Kind string

const (
	KindSkill   Kind = "skill"
	KindMCP     Kind = "mcp"
	KindScript  Kind = "script"
	KindHook    Kind = "hook"
	KindGate    Kind = "gate"
	KindHarness Kind = "harness"
	KindBinary  Kind = "binary"
	KindLoop    Kind = "loop"
)

// Kinds lista as categorias na ordem de exibição.
func Kinds() []Kind {
	return []Kind{KindSkill, KindMCP, KindScript, KindHook, KindGate, KindHarness, KindBinary, KindLoop}
}

// Item é um ativo inventariado. Scope: global, project, local, builtin ou system.
type Item struct {
	ID          string            `json:"id"`
	Kind        Kind              `json:"kind"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Path        string            `json:"path,omitempty"`
	Scope       string            `json:"scope,omitempty"`
	Origin      string            `json:"origin,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`

	// Command é o comando bruto de um hook. Fica só no servidor (json:"-"): a descrição pública é redigida.
	Command string `json:"-"`
}

// Snapshot é o resultado de uma varredura. Version só muda quando o conteúdo muda.
type Snapshot struct {
	Items     []Item       `json:"-"`
	Counts    map[Kind]int `json:"counts"`
	ScannedAt time.Time    `json:"scanned_at"`
	Version   uint64       `json:"version"`
}

// Of devolve os itens de uma categoria, já ordenados.
func (s Snapshot) Of(k Kind) []Item {
	var out []Item
	for _, it := range s.Items {
		if it.Kind == k {
			out = append(out, it)
		}
	}
	return out
}

// Config injeta as dependências do scanner (testável, sem estado global).
type Config struct {
	Home           string
	Projects       []string
	AutoDiscover   bool
	DiscoveryRoots []string
	LoopsDir       string
	LookPath       func(string) (string, error)
	Getenv         func(string) string
}

func (c Config) withDefaults() Config {
	if c.Home == "" {
		c.Home, _ = os.UserHomeDir()
	}
	if c.LookPath == nil {
		c.LookPath = exec.LookPath
	}
	if c.Getenv == nil {
		c.Getenv = os.Getenv
	}
	return c
}

func cloneItems(in []Item) []Item {
	out := make([]Item, len(in))
	for i, it := range in {
		out[i] = it
		if it.Meta != nil {
			out[i].Meta = make(map[string]string, len(it.Meta))
			for k, v := range it.Meta {
				out[i].Meta[k] = v
			}
		}
	}
	return out
}
