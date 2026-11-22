package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/cli"
	"github.com/heverton-dev/LegoLoop/internal/provider"
)

// semearRepoFechar monta um repositório de teste com o manual, o livro, a Skill e o DAG reais mínimos.
func semearRepoFechar(t *testing.T, slug string) string {
	t.Helper()
	repo := t.TempDir()
	canon := "legoloop-" + slug
	manual := "2026-11-15_MANUAL_" + strings.ToUpper(strings.ReplaceAll(slug, "-", "_")) + ".md"
	skill := "---\nname: " + canon + "\n---\nNós do DAG citados: n01.\n" + strings.Repeat("Passo operacional.\n", 12)
	dag := `{"name":"` + canon + `","version":"1.0.0","nodes":[{"id":"n01"}]}`
	arquivos := map[string]string{
		filepath.Join("docs", "manuals", manual):                                 "# Manual\n",
		filepath.Join("docs", "books", canon, "livro.json"):                      `{"titulo":"t","slug":"s","capitulos":[]}`,
		filepath.Join("skills", canon, "SKILL.md"):                               skill,
		filepath.Join(caminhos.DirProcessos, canon, caminhos.ArquivoProcessoDAG): dag,
	}
	for rel, conteudo := range arquivos {
		caminho := filepath.Join(repo, rel)
		if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return repo
}

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
		deps := cli.Deps{Stdout: &stdout, Stderr: &stderr, Discover: noModels, Fechar: &cli.FecharDeps{RepoDir: semearRepoFechar(t, "video-motion")}}
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
		deps := cli.Deps{Stdout: &stdout, Stderr: &stderr, Discover: noModels, Fechar: &cli.FecharDeps{RepoDir: semearRepoFechar(t, "video-motion")}}
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
		deps := cli.Deps{
			Stdout: &stdout,
			Stderr: &stderr,
			Discover: func(context.Context) provider.Inventory {
				return provider.Inventory{
					RouterURL:    "http://127.0.0.1:20128",
					RouterModels: []string{"fast-brain"},
				}
			},
		}
		code := cli.Run(context.Background(), []string{"fechar", "video-motion", "--dry-run", "--model", "modelo-inexistente-abc"}, deps)
		if code == 0 {
			t.Fatalf("esperava erro ao chamar fechar com modelo inexistente")
		}
		if !strings.Contains(stderr.String(), "cérebro") && !strings.Contains(stderr.String(), "modelo") {
			t.Errorf("esperava erro de modelo no stderr: %s", stderr.String())
		}
	})
}

func TestCmdFecharFalhaNaoImprimeSucessoESaiComCodigo1(t *testing.T) {
	repo := t.TempDir()
	ok := func(string) error { return nil }
	deps := func(out, errb *bytes.Buffer) cli.Deps {
		return cli.Deps{Stdout: out, Stderr: errb, Discover: noModels, Fechar: &cli.FecharDeps{RepoDir: repo, Indexer: ok, Portoes: ok, Rebuilder: ok}}
	}
	t.Run("texto", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run(context.Background(), []string{"fechar", "sem-artefatos", "--skip-build"}, deps(&stdout, &stderr))
		if code != 1 {
			t.Fatalf("esperava exit 1 com etapa em erro, obteve %d. Stdout: %s", code, stdout.String())
		}
		if strings.Contains(stdout.String(), "canonizado com sucesso") {
			t.Errorf("imprimiu sucesso falso:\n%s", stdout.String())
		}
		if !strings.Contains(stdout.String(), "falhou") {
			t.Errorf("esperava a palavra falhou na saída:\n%s", stdout.String())
		}
	})
	t.Run("json", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run(context.Background(), []string{"fechar", "sem-artefatos", "--skip-build", "--json"}, deps(&stdout, &stderr))
		if code != 1 {
			t.Fatalf("esperava exit 1 também com --json, obteve %d", code)
		}
		var res struct {
			Sucesso bool `json:"sucesso"`
		}
		if err := json.Unmarshal(stdout.Bytes(), &res); err != nil || res.Sucesso {
			t.Errorf("esperava JSON válido com sucesso=false: err=%v raw=%s", err, stdout.String())
		}
	})
}
