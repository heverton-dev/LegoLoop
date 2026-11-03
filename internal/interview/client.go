package interview

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RouterLLMClient conecta o motor socrático ao 9Router ou a qualquer endpoint compatível com OpenAI.
type RouterLLMClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

// NewRouterClient cria um cliente configurado para o router local ou remoto.
func NewRouterClient(baseURL, apiKey, model string) *RouterLLMClient {
	if model == "" {
		model = "code-pro"
	}
	return &RouterLLMClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
		HTTPClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

func (c *RouterLLMClient) Chat(ctx context.Context, messages []Message) (string, error) {
	if c.BaseURL == "" {
		return "", errors.New("router_url não configurada")
	}

	url := c.BaseURL + "/v1/chat/completions"
	reqBody := map[string]any{
		"model":       c.Model,
		"messages":    messages,
		"temperature": 0.3,
		"stream":      false, // o 9Router transmite em SSE por padrão; este cliente só decodifica JSON
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("falha na requisição ao 9Router: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("9Router retornou status %d: %s", resp.StatusCode, string(body))
	}

	var completion struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return "", fmt.Errorf("falha ao decodificar JSON do 9Router: %w", err)
	}
	if len(completion.Choices) == 0 {
		return "", errors.New("9Router retornou choices vazias")
	}

	return completion.Choices[0].Message.Content, nil
}
