package brain

import (
	"context"
	"fmt"
	"time"
)

// ExecuteSpeculative despacha inicialmente para fast; em caso de erro, realiza hot-swap para heavy (Motor IV).
func ExecuteSpeculative(ctx context.Context, fast, heavy Brain, prompt string, fastTimeout time.Duration) (string, string, error) {
	fastCtx, cancel := context.WithTimeout(ctx, fastTimeout)
	defer cancel()

	resp, err := fast.Infer(fastCtx, prompt)
	if err == nil {
		return resp, fast.Name(), nil
	}

	if heavy == nil {
		return "", fast.Name(), fmt.Errorf("speculative cascade: modelo rápido falhou e não há modelo pesado de contingência: %w", err)
	}

	enrichedPrompt := fmt.Sprintf("%s\n\n[CONTEXTO NEGATIVO]: A tentativa prévia com modelo rápido falhou com o erro: %v", prompt, err)
	heavyResp, heavyErr := heavy.Infer(ctx, enrichedPrompt)
	if heavyErr != nil {
		return "", heavy.Name(), fmt.Errorf("speculative cascade: modelo pesado de contingência também falhou: %w", heavyErr)
	}

	return heavyResp, heavy.Name(), nil
}
