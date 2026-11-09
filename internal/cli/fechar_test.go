package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"legoloop/internal/cli"
)

func TestCmdFechar(t *testing.T) {
	t.Run("help do comando fechar", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}
		code := cli.Run(context.Background(), []string{"fechar", "--help"}, deps)
		if code != 0 {
			t.Fatalf("esperava exit 0 para help, obteve %d", code)
		}
		if !strings.Contains(stdout.String(), "legoloop fechar") {
			t.Errorf("esperava manual de uso na saída: %s", stdout.String())
		}
	})

	t.Run("execucao dry-run com json", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}
		code := cli.Run(context.Background(), []string{"fechar", "video-motion", "--dry-run", "--json"}, deps)
		if code != 0 {
			t.Fatalf("esperava exit 0 para dry-run json, obteve %d. Stderr: %s", code, stderr.String())
		}

		var res struct {
			Slug       string `json:"slug"`
			CanonName  string `json:"canonical_name"`
			Sucesso    bool   `json:"sucesso"`
			Telemetria struct {
				TempoEstimadoMin  float64 `json:"tempo_estimado_min"`
				TokensInEstimados int     `json:"tokens_in_estimados"`
				CustoEstimadoUSD  float64 `json:"custo_estimado_usd"`
			} `json:"telemetria"`
			Etapas []struct {
				Etapa int    `json:"etapa"`
				Nome  string `json:"nome"`
			} `json:"etapas"`
		}
		if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
			t.Fatalf("esperava JSON válido na saída, obteve erro: %v (raw: %s)", err, stdout.String())
		}
		if res.Slug != "video-motion" || res.CanonName != "legoloop-video-motion" {
			t.Errorf("slug inesperado no retorno: %+v", res)
		}
		if !res.Sucesso || len(res.Etapas) != 9 {
			t.Errorf("esperava 9 etapas com sucesso, obteve: %+v", res)
		}
	})

	t.Run("execucao com telemetria personalizada", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}
		code := cli.Run(context.Background(), []string{
			"fechar", "video-motion", "--dry-run", "--json",
			"--tempo-min", "30", "--tokens-in", "90000", "--custo-usd", "0.55",
		}, deps)
		if code != 0 {
			t.Fatalf("esperava exit 0, obteve %d. Stderr: %s", code, stderr.String())
		}
		var res struct {
			Telemetria struct {
				TempoEstimadoMin  float64 `json:"tempo_estimado_min"`
				TokensInEstimados int     `json:"tokens_in_estimados"`
				CustoEstimadoUSD  float64 `json:"custo_estimado_usd"`
			} `json:"telemetria"`
		}
		if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if res.Telemetria.TempoEstimadoMin != 30.0 || res.Telemetria.TokensInEstimados != 90000 {
			t.Errorf("telemetria customizada não propagada: %+v", res.Telemetria)
		}
	})

	t.Run("rejeita sem slug", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}
		code := cli.Run(context.Background(), []string{"fechar"}, deps)
		if code == 0 {
			t.Fatalf("esperava erro ao chamar fechar sem slug, obteve 0")
		}
	})

	t.Run("suporte a model e effort", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}
		code := cli.Run(context.Background(), []string{"fechar", "video-motion", "--dry-run", "--model", "modelo-inexistente-abc"}, deps)
		if code == 0 {
			t.Fatalf("esperava erro ao chamar fechar com modelo inexistente")
		}
		if !strings.Contains(stderr.String(), "cérebro") && !strings.Contains(stderr.String(), "modelo") {
			t.Errorf("esperava erro de modelo no stderr: %s", stderr.String())
		}
	})
}
