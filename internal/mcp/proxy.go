package mcp

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Proxy sobe servidores MCP externos so na primeira chamada e os encerra por ociosidade.
type Proxy struct {
	mu      sync.Mutex
	specs   map[string]ServerSpec
	clients map[string]*stdioClient
	timers  map[string]*time.Timer
	idle    time.Duration
	// recarregar relê o registro antes de subir um servidor; sem ele vale o mapa inicial.
	recarregar func() (map[string]ServerSpec, error)
}

// spec devolve a configuração atual do servidor; registro ilegível mantém a última boa.
func (p *Proxy) spec(name string) (ServerSpec, bool) {
	if p.recarregar != nil {
		if atual, err := p.recarregar(); err == nil {
			p.specs = atual
		}
	}
	spec, ok := p.specs[name]
	return spec, ok
}

// NewProxy cria o proxy; idle <= 0 desliga o encerramento por ociosidade.
func NewProxy(specs map[string]ServerSpec, idle time.Duration) *Proxy {
	return &Proxy{specs: specs, clients: map[string]*stdioClient{}, timers: map[string]*time.Timer{}, idle: idle}
}

func (p *Proxy) ativos() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.clients)
}

// client devolve o cliente do servidor, subindo-o se ainda nao estiver ativo.
func (p *Proxy) client(name string) (*stdioClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if c, ok := p.clients[name]; ok {
		return c, nil
	}
	spec, ok := p.spec(name)
	if !ok {
		return nil, fmt.Errorf("servidor MCP %q nao esta no registro", name)
	}
	c, err := startStdio(spec)
	if err != nil {
		return nil, fmt.Errorf("falha ao subir %q: %w", name, err)
	}
	p.clients[name] = c
	return c, nil
}

func (p *Proxy) touch(name string) {
	if p.idle <= 0 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if t, ok := p.timers[name]; ok {
		t.Stop()
	}
	p.timers[name] = time.AfterFunc(p.idle, func() { p.stop(name) })
}

func (p *Proxy) stop(name string) {
	p.mu.Lock()
	c, ok := p.clients[name]
	delete(p.clients, name)
	if t, has := p.timers[name]; has {
		t.Stop()
		delete(p.timers, name)
	}
	p.mu.Unlock()
	if ok {
		c.close()
	}
}

// Close encerra todos os servidores ativos.
func (p *Proxy) Close() {
	p.mu.Lock()
	names := make([]string, 0, len(p.clients))
	for n := range p.clients {
		names = append(names, n)
	}
	p.mu.Unlock()
	for _, n := range names {
		p.stop(n)
	}
}

// CallTool executa uma ferramenta do servidor, subindo-o sob demanda.
func (p *Proxy) CallTool(server, tool string, args map[string]any) (string, error) {
	c, err := p.client(server)
	if err != nil {
		return "", err
	}
	if args == nil {
		args = map[string]any{}
	}
	raw, err := c.call("tools/call", map[string]any{"name": tool, "arguments": args})
	if err != nil {
		p.stop(server)
		return "", err
	}
	p.touch(server)
	return flattenToolResult(raw)
}

func flattenToolResult(raw json.RawMessage) (string, error) {
	var res struct {
		IsError bool             `json:"isError"`
		Content []map[string]any `json:"content"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return "", err
	}
	parts := make([]string, 0, len(res.Content))
	for _, item := range res.Content {
		if txt, ok := item["text"].(string); ok {
			parts = append(parts, txt)
			continue
		}
		b, _ := json.Marshal(item)
		parts = append(parts, string(b))
	}
	text := strings.Join(parts, "\n")
	if res.IsError {
		return "", fmt.Errorf("%s", text)
	}
	return text, nil
}

// listTools pergunta ao servidor quais ferramentas ele expoe.
func (p *Proxy) listTools(server string) ([]CatalogEntry, error) {
	c, err := p.client(server)
	if err != nil {
		return nil, err
	}
	raw, err := c.call("tools/list", map[string]any{})
	if err != nil {
		return nil, err
	}
	var res struct {
		Tools []CatalogEntry `json:"tools"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, err
	}
	return res.Tools, nil
}

// Refresh sobe cada servidor uma vez, indexa suas ferramentas e o encerra.
func (p *Proxy) Refresh() (map[string][]CatalogEntry, map[string]error) {
	cat := map[string][]CatalogEntry{}
	errs := map[string]error{}
	for name, spec := range p.specs {
		if spec.Command == "" {
			continue
		}
		tools, err := p.listTools(name)
		p.stop(name)
		if err != nil {
			errs[name] = err
			continue
		}
		cat[name] = tools
	}
	return cat, errs
}
