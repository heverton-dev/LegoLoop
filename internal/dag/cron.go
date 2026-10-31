package dag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// maxCronLen limita a expressão cron (proteção contra entrada abusiva).
const maxCronLen = 128

// cronBounds são os limites [mín, máx] de minuto, hora, dia do mês, mês e dia da semana
// (0 e 7 = domingo).
var cronBounds = [5][2]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 7}}

var cronFieldNames = [5]string{"minuto", "hora", "dia do mês", "mês", "dia da semana"}

// validateCron aceita cron de 5 campos numéricos com *, listas (a,b), faixas (a-b) e passos (/n).
func validateCron(expr string) error {
	if len(expr) > maxCronLen {
		return fmt.Errorf("trigger cron: expressão excede %d caracteres", maxCronLen)
	}
	fields := strings.Fields(expr)
	if len(fields) != len(cronBounds) {
		return fmt.Errorf("trigger cron: exige 5 campos, recebeu %d", len(fields))
	}
	for i, f := range fields {
		if err := validateCronField(f, cronBounds[i]); err != nil {
			return fmt.Errorf("trigger cron: campo %s: %w", cronFieldNames[i], err)
		}
	}
	return nil
}

func validateCronField(field string, bound [2]int) error {
	for part := range strings.SplitSeq(field, ",") {
		if err := validateCronPart(part, bound); err != nil {
			return err
		}
	}
	return nil
}

func validateCronPart(part string, bound [2]int) error {
	rng, step, hasStep := strings.Cut(part, "/")
	if hasStep {
		n, err := cronNumber(step, 1, bound[1])
		if err != nil || n < 1 {
			return fmt.Errorf("passo inválido %q", step)
		}
	}
	if rng == "*" {
		return nil
	}
	lo, hi, isRange := strings.Cut(rng, "-")
	from, err := cronNumber(lo, bound[0], bound[1])
	if err != nil || !isRange {
		return err
	}
	to, err := cronNumber(hi, bound[0], bound[1])
	if err != nil {
		return err
	}
	if to < from {
		return fmt.Errorf("faixa invertida %q", rng)
	}
	return nil
}

// cronNumber aceita só dígitos (Atoi aceitaria "+1") dentro de [min, max].
func cronNumber(s string, min, max int) (int, error) {
	if s == "" || strings.Trim(s, "0123456789") != "" {
		return 0, fmt.Errorf("valor inválido %q", s)
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("valor fora do intervalo")
	}
	if n < min || n > max {
		return 0, fmt.Errorf("valor %d fora de %d-%d", n, min, max)
	}
	return n, nil
}
