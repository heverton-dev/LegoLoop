package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFalhaDeTelemetriaEhRelatada(t *testing.T) {
	t.Run("destino eh diretorio existente", func(t *testing.T) {
		workdir := t.TempDir()
		badPath := filepath.Join(workdir, "telemetria_bloqueada")
		if err := os.MkdirAll(badPath, 0o755); err != nil {
			t.Fatalf("criar diretorio bloqueio: %v", err)
		}

		d := diamond(t, 5)
		rep, err := Run(context.Background(), d, Options{
			Workdir:          workdir,
			Executor:         &fakeExec{},
			FlightRecordPath: badPath,
		})
		if err != nil {
			t.Fatalf("Run retornou erro inesperado: %v", err)
		}
		if !rep.OK {
			t.Fatalf("esperava execucao OK mesmo com aviso de telemetria, obteve: %+v", rep)
		}
		if len(rep.Warnings) == 0 {
			t.Fatalf("esperava falha de telemetria relatada em Report.Warnings, obteve vazio")
		}
		junto := strings.Join(rep.Warnings, "; ")
		if !strings.Contains(junto, "telemetria") && !strings.Contains(junto, "flight_record") {
			t.Fatalf("aviso nao menciona telemetria/flight_record: %q", junto)
		}
	})

	t.Run("diretorio pai eh arquivo", func(t *testing.T) {
		workdir := t.TempDir()
		blocker := filepath.Join(workdir, "arquivo_bloqueio")
		if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
			t.Fatalf("criar arquivo bloqueio: %v", err)
		}
		badPath := filepath.Join(blocker, "flight_record.json")

		d := diamond(t, 5)
		rep, err := Run(context.Background(), d, Options{
			Workdir:          workdir,
			Executor:         &fakeExec{},
			FlightRecordPath: badPath,
		})
		if err != nil {
			t.Fatalf("Run retornou erro inesperado: %v", err)
		}
		if len(rep.Warnings) == 0 {
			t.Fatalf("esperava falha de mkdir da telemetria relatada em Report.Warnings, obteve vazio")
		}
		junto := strings.Join(rep.Warnings, "; ")
		if !strings.Contains(junto, "telemetria") {
			t.Fatalf("aviso nao menciona telemetria: %q", junto)
		}
	})
}
