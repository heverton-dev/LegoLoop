package runner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"legoloop/internal/dag"
)

func TestCronNextTick(t *testing.T) {
	// 2026-10-31 12:00:00 (Sábado)
	ref := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		expr     string
		expected time.Time
	}{
		{
			expr:     "*/5 * * * *",
			expected: time.Date(2026, 10, 31, 12, 5, 0, 0, time.UTC),
		},
		{
			expr:     "30 14 * * *",
			expected: time.Date(2026, 10, 31, 14, 30, 0, 0, time.UTC),
		},
		{
			expr:     "0 0 * * *",
			expected: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range cases {
		next, err := NextCronTime(tc.expr, ref)
		if err != nil {
			t.Fatalf("erro ao calcular cron %q: %v", tc.expr, err)
		}
		if !next.Equal(tc.expected) {
			t.Errorf("cron %q: esperado %v, obteve %v", tc.expr, tc.expected, next)
		}
	}
}

func TestFileChangeWatcher(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "input.txt")
	if err := os.WriteFile(targetFile, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	triggered := make(chan string, 1)
	watcher := NewFileWatcher(tmpDir, "input.txt", 20*time.Millisecond)

	go func() {
		_ = watcher.Watch(ctx, func(path string) {
			triggered <- path
		})
	}()

	// Aguarda watcher iniciar
	time.Sleep(30 * time.Millisecond)

	// Modifica o arquivo
	if err := os.WriteFile(targetFile, []byte("v2-modificado"), 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case p := <-triggered:
		if p != "input.txt" {
			t.Fatalf("esperado input.txt, obteve %s", p)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout aguardando trigger de file_change")
	}
}

func TestTriggerSchedulerManual(t *testing.T) {
	trg := dag.Trigger{Type: dag.TriggerManual}
	sched := NewTriggerScheduler(t.TempDir())

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	fired := false
	err := sched.Listen(ctx, trg, func() error {
		fired = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !fired {
		t.Fatal("trigger manual deveria disparar imediatamente")
	}
}
