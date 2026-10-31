// Package provider faz o Inventário Passivo de provedores locais (SPEC_ORCHESTRIX_PROVIDER_AWARE):
// Ollama, 9Router, CLIs no PATH e presença (nunca o valor) de chaves de API.
package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// ScanTimeout é o orçamento de cada sonda HTTP passiva em loopback.
const ScanTimeout = 200 * time.Millisecond

// RemoteScanTimeout é o orçamento para bases remotas. Medido em 31/10/2026: o 9Router
// na VPS devolve /v1/models (217 modelos, 82 KB) em 1,2–2,7 s; 3 s falhava de forma intermitente.
const RemoteScanTimeout = 8 * time.Second

// scanTimeoutFor escolhe o orçamento pelo host: loopback é local e rápido; o resto é rede.
func scanTimeoutFor(base string) time.Duration {
	u, err := url.Parse(base)
	if err != nil {
		return ScanTimeout
	}
	host := u.Hostname()
	if ip := net.ParseIP(host); host == "localhost" || (ip != nil && ip.IsLoopback()) {
		return ScanTimeout
	}
	return RemoteScanTimeout
}

// maxScanBody limita respostas de sondas (defesa contra respostas gigantes).
const maxScanBody = 1 << 20

// CLI é um harness encontrado no PATH.
type CLI struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// ScanOllama lista os modelos instalados via GET /api/tags com timeout de 200ms.
func ScanOllama(ctx context.Context, c *http.Client, baseURL string) ([]string, error) {
	var body struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := getJSON(ctx, c, baseURL+"/api/tags", "", &body); err != nil {
		return nil, fmt.Errorf("ollama: %w", err)
	}
	out := make([]string, 0, len(body.Models))
	for _, m := range body.Models {
		out = append(out, m.Name)
	}
	sort.Strings(out)
	return out, nil
}

// ScanNineRouter testa as bases em ordem e devolve a primeira que responde a GET /v1/models.
// A credencial só é enviada à primeira base (a confiável); alternativas são sondadas sem ela.
func ScanNineRouter(ctx context.Context, c *http.Client, bases []string, apiKey string) (string, []string, error) {
	errs := []error{errors.New("9router: nenhuma base respondeu")}
	for i, base := range bases {
		key := apiKey
		if i > 0 {
			key = ""
		}
		var body struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if err := getJSON(ctx, c, base+"/v1/models", key, &body); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", base, err))
			continue
		}
		ids := make([]string, 0, len(body.Data))
		for _, m := range body.Data {
			ids = append(ids, m.ID)
		}
		ids = sortRouterModels(ids)
		return base, ids, nil
	}
	return "", nil, errors.Join(errs...)
}

// sortRouterModels põe os combos (IDs sem "/", ex.: code-pro) antes dos modelos
// "provedor/modelo", ambos em ordem alfabética: os combos são as rotas recomendadas.
func sortRouterModels(ids []string) []string {
	sort.Slice(ids, func(i, j int) bool {
		ci, cj := !strings.Contains(ids[i], "/"), !strings.Contains(ids[j], "/")
		if ci != cj {
			return ci
		}
		return ids[i] < ids[j]
	})
	return ids
}

// DetectCLIs devolve os harnesses presentes no PATH, na ordem pedida.
func DetectCLIs(lookPath func(string) (string, error), names []string) []CLI {
	var out []CLI
	for _, n := range names {
		if p, err := lookPath(n); err == nil {
			out = append(out, CLI{Name: n, Path: p})
		}
	}
	return out
}

// apiKeyVars mapeia provedor → variável de ambiente inspecionada.
func apiKeyVars() [][2]string {
	return [][2]string{
		{"anthropic", "ANTHROPIC_API_KEY"},
		{"openai", "OPENAI_API_KEY"},
		{"gemini", "GEMINI_API_KEY"},
		{"groq", "GROQ_API_KEY"},
	}
}

// DetectAPIKeys informa QUAIS provedores têm chave definida; o valor nunca sai desta função.
func DetectAPIKeys(getenv func(string) string) []string {
	var out []string
	for _, kv := range apiKeyVars() {
		if strings.TrimSpace(getenv(kv[1])) != "" {
			out = append(out, kv[0])
		}
	}
	return out
}

func getJSON(ctx context.Context, c *http.Client, endpoint, bearer string, dst any) error {
	ctx, cancel := context.WithTimeout(ctx, scanTimeoutFor(endpoint))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status HTTP %d", resp.StatusCode)
	}
	return json.NewDecoder(io.LimitReader(resp.Body, maxScanBody)).Decode(dst)
}
