package fechar_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/fechar"
)

func TestFecharSemStubsETelemetriaReal(t *testing.T) {
	t.Run("telemetria sem valores ficticios 80000/65000 e duracao em ms reais", func(t *testing.T) {
		tmp := t.TempDir()
		semearDocs(t, tmp, "video-motion")
		semearDAG(t, tmp, "video-motion")
		tick := 0
		clock := func() time.Time {
			tick++
			return time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC).Add(time.Duration(tick*25) * time.Millisecond)
		}
		deps := fechar.Deps{
			RepoDir:   tmp,
			VaultDir:  tmp,
			Home:      tmp,
			SkipBuild: true,
			Clock:     clock,
			Indexer:   func(string) error { return nil },
			Portoes:   func(string) error { return nil },
		}

		res, err := fechar.Executar(context.Background(), fechar.Request{Slug: "video-motion"}, deps)
		if err != nil {
			t.Fatalf("Executar falhou: %v", err)
		}
		if !res.Sucesso {
			t.Fatalf("esperava sucesso, etapas=%+v", res.Etapas)
		}
		if res.Telemetria.TokensInEstimados == 80000 || res.Telemetria.TokensInGastos == 65000 ||
			res.Telemetria.TokensOutEstimados == 20000 || res.Telemetria.TokensOutGastos == 14500 {
			t.Fatalf("telemetria ainda contem valores ficticios: %+v", res.Telemetria)
		}
		if res.Telemetria.TokensInEstimados != 0 || res.Telemetria.TokensInGastos != 0 {
			t.Errorf("tokens nao informados devem permanecer 0 (nao medido): %+v", res.Telemetria)
		}
		if res.Telemetria.StatusMedicao != "nao medido" {
			t.Errorf("esperava status_medicao='nao medido' quando tokens nao medidos, obteve %q", res.Telemetria.StatusMedicao)
		}

		raw, err := json.Marshal(res)
		if err != nil {
			t.Fatalf("marshal result: %v", err)
		}
		var doc struct {
			DuracaoTotalMS int64 `json:"duracao_total_ms"`
			Etapas         []struct {
				Etapa     int    `json:"etapa"`
				Artefato  string `json:"artefato"`
				DuracaoMS int64  `json:"duracao_ms"`
			} `json:"etapas"`
		}
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatalf("unmarshal json: %v", err)
		}
		if doc.DuracaoTotalMS <= 0 || doc.DuracaoTotalMS > 10000 {
			t.Errorf("duracao_total_ms deve estar em milissegundos reais (nao nanossegundos), obteve %d", doc.DuracaoTotalMS)
		}
		for _, e := range doc.Etapas {
			if e.DuracaoMS < 0 || e.DuracaoMS > 10000 {
				t.Errorf("etapa %d duracao_ms fora da escala de ms: %d", e.Etapa, e.DuracaoMS)
			}
			for _, stub := range []string{"build_skipped", "rebuilt_via_custom", "graph_skipped_dry_run", "reindexed_custom", "portoes_verificados"} {
				if e.Artefato == stub {
					t.Errorf("etapa %d contem string stub proibida %q", e.Etapa, e.Artefato)
				}
			}
		}

		reportBytes, err := os.ReadFile(res.ReportPath)
		if err != nil {
			t.Fatalf("ler report: %v", err)
		}
		if strings.Contains(string(reportBytes), "80000") || strings.Contains(string(reportBytes), "65000") {
			t.Errorf("relatorio MD contem tokens inventados 80000/65000:\n%s", string(reportBytes))
		}
		if !strings.Contains(string(reportBytes), "nao medido") {
			t.Errorf("relatorio MD deve declarar 'nao medido' quando sem telemetria real:\n%s", string(reportBytes))
		}
	})

	t.Run("execucao real de Rebuilder, Governanca, Indexer e Portoes via seams", func(t *testing.T) {
		tmp := t.TempDir()
		semearDocs(t, tmp, "video-rea")
		semearDAG(t, tmp, "video-rea")
		var rebuildCalled, indexCalled, portoesCalled bool
		deps := fechar.Deps{
			RepoDir:   tmp,
			VaultDir:  tmp,
			Home:      tmp,
			SkipBuild: false,
			Rebuilder: func(dir string) error {
				rebuildCalled = true
				binDir := filepath.Join(dir, "bin")
				_ = os.MkdirAll(binDir, 0o755)
				return os.WriteFile(filepath.Join(binDir, "legoloop.exe"), []byte("bin"), 0o755)
			},
			Indexer: func(repo string) error {
				indexCalled = true
				gDir := filepath.Join(repo, ".codebase-memory")
				_ = os.MkdirAll(gDir, 0o755)
				return os.WriteFile(filepath.Join(gDir, "graph.db.zst"), []byte("zst"), 0o644)
			},
			Portoes: func(repo string) error {
				portoesCalled = true
				return nil
			},
		}

		res, err := fechar.Executar(context.Background(), fechar.Request{Slug: "video-rea", DryRun: false}, deps)
		if err != nil || !res.Sucesso {
			t.Fatalf("esperava sucesso: err=%v etapas=%+v", err, res.Etapas)
		}
		if !rebuildCalled || !indexCalled || !portoesCalled {
			t.Errorf("seams nao foram acionados: rebuild=%v index=%v portoes=%v", rebuildCalled, indexCalled, portoesCalled)
		}
		registroData, err := os.ReadFile(filepath.Join(tmp, "docs", "PROCESSOS_CANONIZADOS.md"))
		if err != nil || !strings.Contains(string(registroData), "legoloop-video-rea") {
			t.Errorf("registro de Processos canonizados nao foi atualizado na etapa 6: err=%v data=%q", err, string(registroData))
		}
		agentsData, err := os.ReadFile(filepath.Join(tmp, "AGENTS.md"))
		if err != nil || !strings.Contains(string(agentsData), "docs/PROCESSOS_CANONIZADOS.md") {
			t.Errorf("AGENTS.md nao referencia o registro na etapa 6: err=%v data=%q", err, string(agentsData))
		}
		memData, err := os.ReadFile(filepath.Join(tmp, "MEMORY.md"))
		if err != nil || !strings.Contains(string(memData), "legoloop-video-rea") {
			t.Errorf("MEMORY.md nao foi atualizado na etapa 6: err=%v data=%q", err, string(memData))
		}
	})

	t.Run("erro explicito quando codebase-memory-mcp ausente no DefaultIndexer", func(t *testing.T) {
		tmp := t.TempDir()
		t.Setenv("USERPROFILE", tmp)
		t.Setenv("HOME", tmp)
		err := fechar.DefaultIndexerWithLookPath(tmp, func(file string) (string, error) {
			return "", errors.New("executable file not found")
		}, nil)
		if err == nil || !strings.Contains(err.Error(), "codebase-memory-mcp") {
			t.Fatalf("esperava erro explicito de codebase-memory-mcp ausente, obteve: %v", err)
		}
		_ = bytes.Buffer{}
	})

	t.Run("fallback encontra codebase-memory-mcp.exe quando ausente no PATH", func(t *testing.T) {
		tmp := t.TempDir()
		t.Setenv("USERPROFILE", tmp)
		t.Setenv("HOME", tmp)
		esperado := filepath.Join(tmp, "tools", "codebase-memory-mcp", "codebase-memory-mcp.exe")
		if err := os.MkdirAll(filepath.Dir(esperado), 0o755); err != nil {
			t.Fatalf("mkdir fallback: %v", err)
		}
		if err := os.WriteFile(esperado, []byte("stub"), 0o755); err != nil {
			t.Fatalf("write fallback: %v", err)
		}
		var chamadoBin string
		var chamadoArgs []string
		err := fechar.DefaultIndexerWithLookPath(tmp, func(file string) (string, error) {
			return "", errors.New("executable file not found")
		}, func(bin string, args ...string) error {
			chamadoBin = bin
			chamadoArgs = args
			return nil
		})
		if err != nil {
			t.Fatalf("esperava sucesso via fallback, obteve erro: %v", err)
		}
		if chamadoBin != esperado {
			t.Fatalf("binario esperado=%q obteve=%q", esperado, chamadoBin)
		}
		if len(chamadoArgs) != 3 || chamadoArgs[0] != "cli" || chamadoArgs[1] != "index_repository" {
			t.Fatalf("args inesperados: %+v", chamadoArgs)
		}
	})
}
