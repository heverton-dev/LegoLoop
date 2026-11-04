package ste

import (
	"fmt"
	"sort"
	"sync"
)

// Catalog guarda as mensagens do produto por identificador. Toda mensagem mostrada ao usuário passa por
// aqui, para que um teste possa aplicar o perfil STE-PT a todas elas de uma vez.
type Catalog struct {
	mu   sync.RWMutex
	msgs map[string]string
}

// Default é o catálogo do produto. Cada pacote registra as suas mensagens em init.
var Default = NewCatalog()

func NewCatalog() *Catalog { return &Catalog{msgs: map[string]string{}} }

// Register guarda a mensagem. Registrar o mesmo id duas vezes é erro de programação.
func (c *Catalog) Register(id, msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, dup := c.msgs[id]; dup {
		panic(fmt.Sprintf("ste: mensagem %q registrada duas vezes", id))
	}
	c.msgs[id] = msg
}

// Msg devolve a mensagem formatada. Um id desconhecido aparece na saída em vez de sumir.
func (c *Catalog) Msg(id string, args ...any) string {
	c.mu.RLock()
	m, ok := c.msgs[id]
	c.mu.RUnlock()
	if !ok {
		return "[mensagem desconhecida: " + id + "]"
	}
	if len(args) == 0 {
		return m
	}
	return fmt.Sprintf(m, args...)
}

// IDs lista os identificadores registrados em ordem alfabética.
func (c *Catalog) IDs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]string, 0, len(c.msgs))
	for id := range c.msgs {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// Violations aplica o perfil a todas as mensagens. Cada violação traz o id da mensagem.
func (c *Catalog) Violations() []Violation {
	var out []Violation
	for _, id := range c.IDs() {
		c.mu.RLock()
		m := c.msgs[id]
		c.mu.RUnlock()
		for _, v := range Lint(m) {
			out = append(out, Violation{Msg: id + ": " + v.Msg, Regra: v.Regra})
		}
	}
	return out
}

// Msg é atalho para o catálogo do produto.
func Msg(id string, args ...any) string { return Default.Msg(id, args...) }
