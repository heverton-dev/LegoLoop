package kraften

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// CloudflareClient fala com a API REST v4 da Cloudflare de forma idempotente.
type CloudflareClient struct {
	base   string
	zoneID string
	token  string
	hc     *http.Client
}

// NovoCloudflareClient cria o cliente da API v4 com o HTTP client injetado.
func NovoCloudflareClient(base, zoneID, token string, hc *http.Client) *CloudflareClient {
	if hc == nil {
		hc = http.DefaultClient
	}
	return &CloudflareClient{base: strings.TrimRight(base, "/"), zoneID: zoneID, token: token, hc: hc}
}

// EnsureCNAME cria o CNAME apenas se ainda não existir (idempotente).
func (c *CloudflareClient) EnsureCNAME(nome, alvo string) error {
	if strings.TrimSpace(nome) == "" || strings.TrimSpace(alvo) == "" {
		return ErrInvalidDeployInput
	}
	quantos, err := c.contarCNAME(nome)
	if err != nil {
		return err
	}
	if quantos > 0 {
		return nil
	}
	return c.criarCNAME(nome, alvo)
}

// contarCNAME consulta GET /zones/{zoneId}/dns_records?name={fqdn}.
func (c *CloudflareClient) contarCNAME(nome string) (int, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records?name=%s", c.base, c.zoneID, nome)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	c.autenticar(req)
	resp, err := c.hc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var corpo struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&corpo); err != nil {
		return 0, err
	}
	return len(corpo.Result), nil
}

// criarCNAME emite POST /zones/{zoneId}/dns_records com CNAME proxied=false.
func (c *CloudflareClient) criarCNAME(nome, alvo string) error {
	url := fmt.Sprintf("%s/zones/%s/dns_records", c.base, c.zoneID)
	carga := map[string]any{"type": "CNAME", "name": nome, "content": alvo, "proxied": false, "ttl": 1}
	buf, err := json.Marshal(carga)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(buf)))
	if err != nil {
		return err
	}
	c.autenticar(req)
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("kraften: criação de CNAME recusada (%d)", resp.StatusCode)
	}
	return nil
}

func (c *CloudflareClient) autenticar(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
}
