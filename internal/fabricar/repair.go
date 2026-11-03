package fabricar

import (
	"context"
	"errors"
	"slices"
	"time"

	"legoloop/internal/interview"
	"legoloop/internal/specgen"
)

// maxSpecAttempts limita as chamadas ao LLM para gerar a SPEC (1 original + 2 correções).
const maxSpecAttempts = 3

// repairLLM devolve ao modelo a resposta que foi rejeitada, com o motivo, para ele corrigir. Sem isto um único
// desvio de formato de um LLM real derruba a fabricação inteira.
type repairLLM struct {
	inner    interview.LLMClient
	last     string
	feedback string
}

func (r *repairLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	if r.feedback != "" {
		msgs = append(slices.Clone(msgs),
			interview.Message{Role: "assistant", Content: r.last},
			interview.Message{Role: "user", Content: "Sua SPEC foi rejeitada pelo validador: " + r.feedback +
				". Reescreva a SPEC COMPLETA seguindo estritamente o formato do contrato, sem comentários fora dela."})
	}
	out, err := r.inner.Chat(ctx, msgs)
	r.last = out
	return out, err
}

// generateSpec gera a SPEC com até maxSpecAttempts tentativas. Só repete quando o LLM respondeu fora do contrato;
// falha de transporte ou de gravação sobe na hora (repetir só custaria tempo e dinheiro).
func generateSpec(ctx context.Context, llm interview.LLMClient, stage, intentPath string, now time.Time) (*specgen.SpecResult, error) {
	rl := &repairLLM{inner: llm}
	gen := specgen.NewGenerator(rl, nil)
	gen.Now = func() time.Time { return now }
	for attempt := 1; ; attempt++ {
		spec, err := gen.Generate(ctx, stage, intentPath)
		if err == nil || !errors.Is(err, specgen.ErrContract) || attempt >= maxSpecAttempts {
			return spec, err
		}
		rl.feedback = err.Error()
	}
}
