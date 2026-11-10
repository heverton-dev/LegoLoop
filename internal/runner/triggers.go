package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"legoloop/internal/dag"
)

// NextCronTime calcula o próximo momento exato (minuto cheio) em que a expressão cron dispara
// a partir do instante `from`.
func NextCronTime(expr string, from time.Time) (time.Time, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return time.Time{}, fmt.Errorf("cron inválido: exige 5 campos, obteve %d", len(fields))
	}

	// Começa no próximo minuto limpo
	t := from.Truncate(time.Minute).Add(time.Minute)

	// Procura até 1 ano à frente (525.600 minutos)
	limit := t.Add(366 * 24 * time.Hour)
	for t.Before(limit) {
		if matchCronField(fields[0], t.Minute(), 0, 59) &&
			matchCronField(fields[1], t.Hour(), 0, 23) &&
			matchCronField(fields[2], t.Day(), 1, 31) &&
			matchCronField(fields[3], int(t.Month()), 1, 12) &&
			matchCronField(fields[4], int(t.Weekday()), 0, 7) {
			return t, nil
		}
		t = t.Add(time.Minute)
	}

	return time.Time{}, fmt.Errorf("nenhum disparo encontrado no próximo ano para cron %q", expr)
}

func matchCronField(pattern string, value, min, max int) bool {
	// Normaliza domingo (0 ou 7) se max for 7
	if max == 7 && value == 0 && strings.Contains(pattern, "7") {
		value = 7
	}

	for _, part := range strings.Split(pattern, ",") {
		if part == "*" {
			return true
		}
		if strings.HasPrefix(part, "*/") {
			step, err := strconv.Atoi(part[2:])
			if err == nil && step > 0 && value%step == 0 {
				return true
			}
			continue
		}
		if strings.Contains(part, "-") {
			var lo, hi int
			if _, err := fmt.Sscanf(part, "%d-%d", &lo, &hi); err == nil {
				if value >= lo && value <= hi {
					return true
				}
			}
			continue
		}
		if num, err := strconv.Atoi(part); err == nil {
			if num == value {
				return true
			}
		}
	}
	return false
}

// FileWatcher monitora alterações em arquivos relativos ao workspace.
type FileWatcher struct {
	workdir      string
	relPath      string
	pollInterval time.Duration
}

func NewFileWatcher(workdir, relPath string, interval time.Duration) *FileWatcher {
	if interval <= 0 {
		interval = 250 * time.Millisecond
	}
	return &FileWatcher{
		workdir:      workdir,
		relPath:      relPath,
		pollInterval: interval,
	}
}

func (w *FileWatcher) Watch(ctx context.Context, onChange func(path string)) error {
	fullPath := filepath.Join(w.workdir, w.relPath)
	var lastMod time.Time

	if info, err := os.Stat(fullPath); err == nil {
		lastMod = info.ModTime()
	}

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			info, err := os.Stat(fullPath)
			if err != nil {
				continue
			}
			if !lastMod.IsZero() && info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				onChange(w.relPath)
			} else if lastMod.IsZero() {
				lastMod = info.ModTime()
			}
		}
	}
}

// TriggerScheduler gerencia os gatilhos formais do contrato de loop.
type TriggerScheduler struct {
	workdir string
}

func NewTriggerScheduler(workdir string) *TriggerScheduler {
	return &TriggerScheduler{workdir: workdir}
}

func (s *TriggerScheduler) Listen(ctx context.Context, trg dag.Trigger, onTrigger func() error) error {
	switch trg.Type {
	case dag.TriggerManual, "":
		return onTrigger()
	case dag.TriggerFileChange:
		w := NewFileWatcher(s.workdir, trg.Path, 100*time.Millisecond)
		return w.Watch(ctx, func(path string) {
			_ = onTrigger()
		})
	case dag.TriggerCron:
		return listenCron(ctx, trg.Cron, onTrigger)
	case dag.TriggerLoopEnd:
		// Em modo standalone, aguarda cancelamento do contexto
		<-ctx.Done()
		return ctx.Err()
	default:
		return fmt.Errorf("tipo de trigger não suportado no runtime: %q", trg.Type)
	}
}

// listenCron dispara onTrigger em cada instante da expressão cron até o contexto terminar.
func listenCron(ctx context.Context, cron string, onTrigger func() error) error {
	for {
		next, err := NextCronTime(cron, time.Now())
		if err != nil {
			return err
		}
		timer := time.NewTimer(time.Until(next))
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			if err := onTrigger(); err != nil {
				return err
			}
		}
	}
}
