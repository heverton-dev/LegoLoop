package fabricar

import (
	"errors"
	"fmt"
	"strings"

	"legoloop/internal/interview"
	"legoloop/internal/ste"
)

// ConfidenceTarget é a confiança da entrevista a partir da qual o INTENT sai de DRAFT.
const ConfidenceTarget = 95

// ErrLowConfidence: a entrevista parou abaixo da meta e o usuário não confirmou que quer seguir assim.
var ErrLowConfidence = errors.New("fabricar: confiança abaixo da meta")

// ConfirmDraftFunc recebe a entrevista que ficou abaixo da meta e devolve true se o usuário quer gerar a SPEC
// mesmo assim.
type ConfirmDraftFunc func(c Clarification) (bool, error)

func init() {
	ste.Default.Register("fabricar.confianca_baixa", "A entrevista parou em %d%%. A meta é %d%%. O INTENT sai como rascunho.")
	ste.Default.Register("fabricar.lacuna", "Lacuna: %s.")
	ste.Default.Register("fabricar.sem_confirmacao",
		"A confiança ficou abaixo da meta. A SPEC não foi gerada. Confirme que quer seguir assim ou responda mais perguntas.")
}

// mirrorGaps lista o que o espelho ainda não respondeu, na ordem do espelho canônico.
func mirrorGaps(m interview.Mirror) []string {
	var gaps []string
	add := func(empty bool, label string) {
		if empty {
			gaps = append(gaps, label)
		}
	}
	add(isBlank(m.User), "quem se beneficia")
	add(isBlank(m.WhyNow), "por que agora")
	add(isBlank(m.Success), "critério de sucesso")
	add(allBlank(m.Constraints), "restrições inegociáveis")
	add(allBlank(m.OutOfScope), "fora de escopo")
	return gaps
}

func isBlank(s string) bool { return strings.TrimSpace(s) == "" }

func allBlank(items []string) bool {
	for _, it := range items {
		if !isBlank(it) {
			return false
		}
	}
	return true
}

// Warnings devolve os avisos de uma entrevista que ficou abaixo da meta: a confiança medida e cada lacuna.
// Entrevista concluída, ou que não aconteceu (modo autônomo), não gera aviso.
func (c Clarification) Warnings() []string {
	if !c.Interviewed || c.Reached {
		return nil
	}
	out := []string{ste.Msg("fabricar.confianca_baixa", c.Confidence, ConfidenceTarget)}
	for _, g := range c.Gaps {
		out = append(out, ste.Msg("fabricar.lacuna", g))
	}
	return out
}

// confirmDraft é o portão antes da SPEC: entrevista abaixo da meta só segue se o usuário confirmar.
func (d Deps) confirmDraft(c Clarification) error {
	if !c.Interviewed || c.Reached {
		return nil
	}
	if d.ConfirmDraft == nil {
		return fmt.Errorf("%w: %s", ErrLowConfidence, ste.Msg("fabricar.sem_confirmacao"))
	}
	ok, err := d.ConfirmDraft(c)
	if err != nil {
		return err
	}
	if !ok {
		return ErrAborted
	}
	return nil
}
