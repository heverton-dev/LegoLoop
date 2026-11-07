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
			Slug      string `json:"slug"`
			CanonName string `json:"canonical_name"`
			Sucesso   bool   `json:"sucesso"`
			Etapas    []struct {
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
		if !res.Sucesso || len(res.Etapas) != 8 {
			t.Errorf("esperava 8 etapas com sucesso, obteve: %+v", res)
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
}
