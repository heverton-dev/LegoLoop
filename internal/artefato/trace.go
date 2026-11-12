package artefato

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type No struct{ ID, SpecRef string }

type Cadeia struct {
	Requisitos map[string][]string
	ACs        []string
	Nos        []No
	Perguntas  []string
	Covers     []string
}

var (
	reqRow = regexp.MustCompile(`^\|\s*(R\d+)\s*\|(.*)$`)
	qRow   = regexp.MustCompile(`^\|\s*(Q\d+)\s*\|`)
	acDef  = regexp.MustCompile("^\\s*\\* \\[.\\]\\s+`(AC-[A-Z0-9]+-\\d+\\.\\d+)`")
	acRef  = regexp.MustCompile(`AC-[A-Z0-9]+-\d+\.\d+`)
	idRef  = regexp.MustCompile(`\b(?:Q|R)\d+\b`)
)

func (c Cadeia) Orfaos() []string {
	var out []string
	for r, acs := range c.Requisitos {
		if len(acs) == 0 {
			out = append(out, r)
		}
	}
	refs := map[string]bool{}
	for _, n := range c.Nos {
		if n.SpecRef == "" {
			out = append(out, n.ID)
			continue
		}
		refs[n.SpecRef] = true
	}
	for _, ac := range c.ACs {
		if !refs[ac] {
			out = append(out, ac)
		}
	}
	cobertas := map[string]bool{}
	for _, id := range c.Covers {
		cobertas[id] = true
	}
	for _, q := range c.Perguntas {
		if !cobertas[q] {
			out = append(out, q)
		}
	}
	sort.Strings(out)
	return compact(out)
}

func compact(s []string) []string {
	var out []string
	for i, v := range s {
		if i == 0 || v != s[i-1] {
			out = append(out, v)
		}
	}
	return out
}

func (c *Cadeia) addMarkdown(texto string) {
	var perguntas, covers []string
	for _, linha := range strings.Split(texto, "\n") {
		linha = strings.TrimRight(linha, "\r")
		if m := reqRow.FindStringSubmatch(linha); m != nil {
			c.addRequisito(m[1], acRef.FindAllString(m[2], -1))
		} else if m := qRow.FindStringSubmatch(linha); m != nil {
			perguntas = append(perguntas, m[1])
		} else if m := acDef.FindStringSubmatch(linha); m != nil {
			c.ACs = append(c.ACs, m[1])
		} else if strings.HasPrefix(linha, "covers:") {
			covers = append(covers, idRef.FindAllString(linha, -1)...)
		}
	}
	c.Perguntas = append(c.Perguntas, perguntas...)
	if len(perguntas) > 0 {
		return
	}
	c.Covers = append(c.Covers, covers...)
	for _, id := range covers {
		if strings.HasPrefix(id, "R") {
			c.addRequisito(id, nil)
		}
	}
}

func (c *Cadeia) addRequisito(id string, acs []string) {
	if c.Requisitos == nil {
		c.Requisitos = map[string][]string{}
	}
	c.Requisitos[id] = append(c.Requisitos[id], acs...)
}

func (c *Cadeia) addDAG(dados []byte) error {
	var d struct {
		Nodes []struct {
			ID      string `json:"id"`
			Type    string `json:"type"`
			SpecRef string `json:"spec_ref"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(dados, &d); err != nil {
		return err
	}
	for _, n := range d.Nodes {
		if n.Type != "stop_gate" || n.SpecRef != "" {
			c.Nos = append(c.Nos, No{ID: n.ID, SpecRef: n.SpecRef})
		}
	}
	return nil
}

func CarregarCadeia(arquivos []string) (Cadeia, error) {
	var c Cadeia
	for _, a := range arquivos {
		dados, err := os.ReadFile(a)
		if err != nil {
			return Cadeia{}, err
		}
		if strings.HasSuffix(strings.ToLower(a), ".json") {
			if err := c.addDAG(dados); err != nil {
				return Cadeia{}, fmt.Errorf("%s: %w", a, err)
			}
			continue
		}
		c.addMarkdown(string(dados))
	}
	return c, nil
}
