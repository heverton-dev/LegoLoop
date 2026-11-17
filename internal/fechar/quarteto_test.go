package fechar

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/dag"
	"legoloop/internal/runner"
	"legoloop/internal/triagem"
)

// executorVazio executa o nó sem efeitos colaterais além do contrato.
type executorVazio struct{}

func (executorVazio) Execute(_ context.Context, _ dag.Node, _ string) (runner.Result, error) {
	return runner.Result{ExitCode: 0}, nil
}

// TestAuxiliarFecharSai só existe para gerar um PID morto nos testes.
func TestAuxiliarFecharSai(t *testing.T) {
	if os.Getenv("LEGOLOOP_AUX_FECHAR") != "1" {
		return
	}
	os.Exit(0)
}

// pidMortoFechar executa um auxiliar que termina e devolve o PID já morto.
func pidMortoFechar(t *testing.T) int {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestAuxiliarFecharSai$")
	cmd.Env = append(os.Environ(), "LEGOLOOP_AUX_FECHAR=1")
	if err := cmd.Run(); err != nil {
		t.Fatalf("auxiliar de teste falhou: %v", err)
	}
	return cmd.Process.Pid
}

// dagCicloMonta um DAG de um nó com token_budget para a telemetria gravar.
func dagCicloMonta(t *testing.T) []byte {
	t.Helper()
	no := dag.Node{
		ID:          "n1",
		Type:        "script",
		Command:     []string{"x"},
		Inputs:      map[string]string{"prompt": "curto"},
		Routing:     dag.RoutingConfig{Mode: "manual", Model: "none"},
		Output:      dag.HandoffContract{Artifact: "n1.md", Gate: "exit_code == 0"},
		TokenBudget: map[string]any{"max_prompt_tokens": 1000},
	}
	raw, err := json.Marshal(map[string]any{"name": "ciclo", "nodes": []dag.Node{no}})
	if err != nil {
		t.Fatalf("marshal do DAG falhou: %v", err)
	}
	return raw
}

// rodarCicloCenario executa o runner na sessão dona e planta a ativa e a órfã.
func rodarCicloCenario(t *testing.T, dir string) (donaID string) {
	t.Helper()
	donaID = "sessao-dona"
	raw := dagCicloMonta(t)
	d, err := dag.Parse(raw)
	if err != nil {
		t.Fatalf("parse do DAG falhou: %v", err)
	}
	if _, err := runner.Run(context.Background(), d, runner.Options{
		Workdir:  dir,
		Executor: executorVazio{},
		SessaoID: donaID,
	}); err != nil {
		t.Fatalf("runner.Run falhou: %v", err)
	}
	if _, err := triagem.LerMemoria(dir, donaID); err != nil {
		t.Fatalf("runner não gravou a sessão dona: %v", err)
	}
	ativa := triagem.MemoriaExecucao{Processo: "p", SessaoID: "sessao-ativa", PID: os.Getpid(), EtapaAtual: "meio"}
	if err := triagem.SalvarMemoria(dir, ativa); err != nil {
		t.Fatalf("SalvarMemoria da ativa falhou: %v", err)
	}
	morto := pidMortoFechar(t)
	orfa := triagem.MemoriaExecucao{Processo: "p", SessaoID: "sessao-orfa", PID: morto, EtapaAtual: "parada"}
	if err := triagem.SalvarMemoria(dir, orfa); err != nil {
		t.Fatalf("SalvarMemoria da órfã falhou: %v", err)
	}
	lockOrfao := filepath.Join(dir, triagem.ArquivoLockMemoria)
	if err := os.WriteFile(lockOrfao, []byte(fmt.Sprintf("%d sessao-orfa", morto)), 0o644); err != nil {
		t.Fatalf("escrita do lock órfão falhou: %v", err)
	}
	return donaID
}

// verificarPosFechar confirma que só a ativa sobrevive e não sobra órfão.
func verificarPosFechar(t *testing.T, dir string) {
	t.Helper()
	dados, err := os.ReadFile(filepath.Join(dir, triagem.ArquivoMemoria))
	if err != nil {
		t.Fatalf("ReadFile da memória falhou: %v", err)
	}
	var est triagem.EstadoMemoria
	if err := json.Unmarshal(dados, &est); err != nil {
		t.Fatalf("envelope inválido após fechar: %v", err)
	}
	if len(est.Sessoes) != 1 {
		t.Fatalf("esperava só a sessão ativa, obtido: %v", est.Sessoes)
	}
	ativa, ok := est.Sessoes["sessao-ativa"]
	if !ok || ativa.EtapaAtual != "meio" {
		t.Fatalf("sessão ativa perdida ou corrompida: %+v", est.Sessoes)
	}
	if _, err := os.Stat(filepath.Join(dir, triagem.ArquivoLockMemoria)); !os.IsNotExist(err) {
		t.Errorf("lock órfão deveria ser removido: %v", err)
	}
	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir falhou: %v", err)
	}
	for _, e := range entradas {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("arquivo temporário órfão no disco: %s", e.Name())
		}
	}
}

// TestFecharSessaoDona prova AC-E01.5: o ciclo executar, gravar por sessão,
// fechar zera só a dona e limpa a órfã preserva a sessão ativa.
func TestFecharSessaoDona(t *testing.T) {
	dir := t.TempDir()
	raw := dagCicloMonta(t)
	dagPath := filepath.Join(dir, caminhos.ArquivoDAG)
	if err := os.WriteFile(dagPath, raw, 0o644); err != nil {
		t.Fatalf("escrita do DAG falhou: %v", err)
	}
	donaID := rodarCicloCenario(t, dir)
	if err := garantirQuarteto(dir, dagPath, donaID); err != nil {
		t.Fatalf("garantirQuarteto falhou: %v", err)
	}
	verificarPosFechar(t, dir)
}

// TestFecharSessaoInexistenteNaoAlteraDisco prova AC-E01.6: sessão inexistente
// devolve erro citando o ID e nenhum arquivo de memória é alterado.
func TestFecharSessaoInexistenteNaoAlteraDisco(t *testing.T) {
	dir := t.TempDir()
	mem := triagem.MemoriaExecucao{Processo: "p", SessaoID: "sessao-viva", PID: os.Getpid(), EtapaAtual: "meio"}
	if err := triagem.SalvarMemoria(dir, mem); err != nil {
		t.Fatalf("SalvarMemoria falhou: %v", err)
	}
	caminho := filepath.Join(dir, triagem.ArquivoMemoria)
	antes, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ReadFile falhou: %v", err)
	}
	err = zerarMemoriaSessao(dir, "sessao-fantasma")
	if err == nil || !strings.Contains(err.Error(), "sessao-fantasma") {
		t.Fatalf("erro %v deveria citar o ID inexistente", err)
	}
	depois, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ReadFile após falha: %v", err)
	}
	if string(antes) != string(depois) {
		t.Errorf("memória alterada apesar do erro:\nantes=%s\ndepois=%s", antes, depois)
	}
}

// TestZeraMemoria prova AC-E03.3 do mega plano: ao fim do fechar, a memória de
// execução está zerada — envelope válido no schema e sem histórico inflado.
func TestZeraMemoria(t *testing.T) {
	dir := t.TempDir()
	raw := dagCicloMonta(t)
	dagPath := filepath.Join(dir, caminhos.ArquivoDAG)
	if err := os.WriteFile(dagPath, raw, 0o644); err != nil {
		t.Fatalf("escrita do DAG falhou: %v", err)
	}
	donaID := rodarCicloCenario(t, dir)
	if err := garantirQuarteto(dir, dagPath, donaID); err != nil {
		t.Fatalf("garantirQuarteto falhou: %v", err)
	}
	// Sessão ativa plantada também sai no fechar total: zera tudo e valida o schema.
	if err := zerarMemoriaSessao(dir, "sessao-ativa"); err != nil {
		t.Fatalf("zerarMemoriaSessao da ativa falhou: %v", err)
	}
	dados, err := os.ReadFile(filepath.Join(dir, triagem.ArquivoMemoria))
	if err != nil {
		t.Fatalf("memória deveria existir zerada: %v", err)
	}
	var est triagem.EstadoMemoria
	if err := json.Unmarshal(dados, &est); err != nil {
		t.Fatalf("envelope zerado inválido: %v", err)
	}
	if est.SchemaVersion != triagem.SchemaMemoria {
		t.Errorf("schema_version = %q, esperado %q", est.SchemaVersion, triagem.SchemaMemoria)
	}
	if len(est.Sessoes) != 0 {
		t.Errorf("memória deveria estar zerada sem histórico, obtido: %v", est.Sessoes)
	}
}

// TestGravaEmProcessos prova AC-E06.2: após o rename, o dialeto de gravação do
// runtime é processos/ — o quarteto do Processo nasce em processos/<slug>/ e o
// caminho legado nunca é recriado.
func TestGravaEmProcessos(t *testing.T) {
	repo := t.TempDir()
	dir := resolverDirProcesso(repo, "legoloop-x", "")
	esp := filepath.Join(repo, caminhos.DirProcessos, "legoloop-x")
	if dir != esp {
		t.Fatalf("resolverDirProcesso devolveu %q, esperado %q", dir, esp)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("criação de %s: %v", dir, err)
	}
	raw := dagCicloMonta(t)
	dagPath := filepath.Join(dir, caminhos.ArquivoDAG)
	if err := os.WriteFile(dagPath, raw, 0o644); err != nil {
		t.Fatalf("escrita do DAG falhou: %v", err)
	}
	donaID := rodarCicloCenario(t, dir)
	if err := garantirQuarteto(dir, dagPath, donaID); err != nil {
		t.Fatalf("garantirQuarteto falhou: %v", err)
	}
	for _, nome := range []string{caminhos.ArquivoDAG, nomeEspelho("legoloop-x"), triagem.ArquivoMemoria} {
		if _, err := os.Stat(filepath.Join(dir, nome)); err != nil {
			t.Errorf("artefato %q ausente em processos/: %v", nome, err)
		}
	}
	if _, err := os.Stat(filepath.Join(repo, caminhos.DirLoops)); !os.IsNotExist(err) {
		t.Fatalf("runtime recriou o diretório legado: %v", err)
	}
}
