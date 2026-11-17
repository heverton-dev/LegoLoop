package triagem

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestMemoriaSessaoSerializa prova AC-E01.1: sessao_id e pid serializam e validam.
func TestMemoriaSessaoSerializa(t *testing.T) {
	mem := MemoriaExecucao{
		Processo: "proc-teste",
		SessaoID: "abc",
		PID:      100,
	}
	if err := mem.Validate(); err != nil {
		t.Fatalf("Validate() válido retornou erro: %v", err)
	}
	dados, err := json.Marshal(mem)
	if err != nil {
		t.Fatalf("Marshal falhou: %v", err)
	}
	texto := string(dados)
	if !strings.Contains(texto, `"sessao_id":"abc"`) {
		t.Errorf("JSON sem sessao_id: %s", texto)
	}
	if !strings.Contains(texto, `"pid":100`) {
		t.Errorf("JSON sem pid: %s", texto)
	}
	var volta MemoriaExecucao
	if err := json.Unmarshal(dados, &volta); err != nil {
		t.Fatalf("Unmarshal falhou: %v", err)
	}
	if volta.SessaoID != "abc" || volta.PID != 100 {
		t.Errorf("round-trip perdeu sessão: %+v", volta)
	}
}

// TestMemoriaSessaoInvalida prova AC-E01.2: sessão sem identidade é rejeitada
// e nada é gravado em disco.
func TestMemoriaSessaoInvalida(t *testing.T) {
	casos := []struct {
		nome   string
		sessao MemoriaExecucao
	}{
		{"id_vazio", MemoriaExecucao{Processo: "p", SessaoID: "", PID: 100}},
		{"pid_zero", MemoriaExecucao{Processo: "p", SessaoID: "abc", PID: 0}},
		{"pid_negativo", MemoriaExecucao{Processo: "p", SessaoID: "abc", PID: -1}},
	}
	for _, tc := range casos {
		t.Run(tc.nome, func(t *testing.T) {
			err := tc.sessao.Validate()
			if !errors.Is(err, ErrSessaoInvalida) {
				t.Fatalf("Validate() = %v; quero ErrSessaoInvalida", err)
			}
			dir := t.TempDir()
			if err := tc.sessao.Validate(); err == nil {
				salvarSeValido(t, dir, tc.sessao)
			}
			entradas, err := os.ReadDir(dir)
			if err != nil {
				t.Fatalf("ReadDir falhou: %v", err)
			}
			if len(entradas) != 0 {
				t.Errorf("disco não está vazio após rejeição: %v", entradas)
			}
		})
	}
}

// salvarSeValido grava só quando a sessão passa na validação.
func salvarSeValido(t *testing.T, dir string, mem MemoriaExecucao) {
	t.Helper()
	if err := mem.Validate(); err != nil {
		return
	}
	if err := SalvarMemoria(dir, mem); err != nil {
		t.Fatalf("SalvarMemoria falhou: %v", err)
	}
}

// TestMemoriaLockConcorrente prova AC-E01.3: duas sessões simultâneas gravam
// sem corrupção de estado e cada leitura posterior devolve a sessão íntegra.
func TestMemoriaLockConcorrente(t *testing.T) {
	dir := t.TempDir()
	bases := []MemoriaExecucao{
		{Processo: "proc-teste", SessaoID: "sessao-a", PID: 11},
		{Processo: "proc-teste", SessaoID: "sessao-b", PID: 22},
	}
	inicio := make(chan struct{})
	var wg sync.WaitGroup
	for _, base := range bases {
		wg.Add(1)
		go func(base MemoriaExecucao) {
			defer wg.Done()
			<-inicio
			gravarLoteSessao(t, dir, base)
		}(base)
	}
	close(inicio)
	wg.Wait()
	for _, base := range bases {
		verificarSessaoFinal(t, dir, base)
	}
	dados, err := os.ReadFile(filepath.Join(dir, ArquivoMemoria))
	if err != nil {
		t.Fatalf("ReadFile da memória falhou: %v", err)
	}
	if !json.Valid(dados) {
		t.Errorf("JSON do arquivo inválido: %s", dados)
	}
}

// gravarLoteSessao grava 50 atualizações sequenciais da mesma sessão.
func gravarLoteSessao(t *testing.T, dir string, base MemoriaExecucao) {
	t.Helper()
	for j := 1; j <= 50; j++ {
		mem := base
		mem.EtapaAtual = fmt.Sprintf("etapa-%d", j)
		mem.ProximoPasso = fmt.Sprintf("passo-%d", j)
		if err := SalvarMemoria(dir, mem); err != nil {
			t.Errorf("SalvarMemoria(%s) falhou: %v", base.SessaoID, err)
			return
		}
	}
}

// verificarSessaoFinal confere que a sessão lida mantém a última atualização íntegra.
func verificarSessaoFinal(t *testing.T, dir string, base MemoriaExecucao) {
	t.Helper()
	lida, err := LerMemoria(dir, base.SessaoID)
	if err != nil {
		t.Fatalf("LerMemoria(%s) falhou: %v", base.SessaoID, err)
	}
	if lida.SessaoID != base.SessaoID || lida.PID != base.PID || lida.Processo != base.Processo {
		t.Fatalf("sessão %s divergiu na leitura: %+v", base.SessaoID, lida)
	}
	if lida.EtapaAtual != "etapa-50" || lida.ProximoPasso != "passo-50" {
		t.Errorf("sessão %s perdeu atualização: etapa=%q proximo=%q",
			base.SessaoID, lida.EtapaAtual, lida.ProximoPasso)
	}
}

// TestMemoriaLockOcupado prova AC-E01.4: lock retido além do tempo limite
// devolve ErrLockOcupado e nada é gravado em disco.
func TestMemoriaLockOcupado(t *testing.T) {
	dir := t.TempDir()
	antigo := timeoutLockMemoria
	timeoutLockMemoria = 50 * time.Millisecond
	defer func() { timeoutLockMemoria = antigo }()
	lk, err := adquirirLockMemoria(dir, 99, "sessao-a")
	if err != nil {
		t.Fatalf("adquirirLockMemoria para reter falhou: %v", err)
	}
	defer lk.liberar()
	memB := MemoriaExecucao{Processo: "proc-teste", SessaoID: "sessao-b", PID: 22, EtapaAtual: "etapa-b"}
	if err := SalvarMemoria(dir, memB); !errors.Is(err, ErrLockOcupado) {
		t.Fatalf("SalvarMemoria com lock retido = %v; quero ErrLockOcupado", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ArquivoMemoria)); !os.IsNotExist(err) {
		t.Errorf("lock retido gravou estado parcial: %v", err)
	}
}

// TestAuxiliarSaiSemEfeito só existe para gerar um PID morto nos testes.
func TestAuxiliarSaiSemEfeito(t *testing.T) {
	if os.Getenv("LEGOLOOP_AUX") != "1" {
		return
	}
	os.Exit(0)
}

// pidMortoParaTeste executa um auxiliar que termina e devolve o PID já morto.
func pidMortoParaTeste(t *testing.T) int {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestAuxiliarSaiSemEfeito$")
	cmd.Env = append(os.Environ(), "LEGOLOOP_AUX=1")
	if err := cmd.Run(); err != nil {
		t.Fatalf("auxiliar de teste falhou: %v", err)
	}
	return cmd.Process.Pid
}

// TestZerarSessaoRemoveSoADona prova AC-E01.5: a dona sai zerada e a ativa fica.
func TestZerarSessaoRemoveSoADona(t *testing.T) {
	dir := t.TempDir()
	dona := MemoriaExecucao{Processo: "p", SessaoID: "sessao-dona", PID: os.Getpid(), EtapaAtual: "fim"}
	ativa := MemoriaExecucao{Processo: "p", SessaoID: "sessao-ativa", PID: os.Getpid(), EtapaAtual: "meio"}
	salvarSeValido(t, dir, dona)
	salvarSeValido(t, dir, ativa)
	if err := ZerarSessao(dir, "sessao-dona"); err != nil {
		t.Fatalf("ZerarSessao falhou: %v", err)
	}
	est, err := lerEstadoDisco(dir)
	if err != nil {
		t.Fatalf("lerEstadoDisco falhou: %v", err)
	}
	if _, ok := est.Sessoes["sessao-dona"]; ok {
		t.Errorf("sessão dona deveria sair do envelope: %v", est.Sessoes)
	}
	got, ok := est.Sessoes["sessao-ativa"]
	if !ok || got.EtapaAtual != "meio" {
		t.Errorf("sessão ativa corrompida: %+v", est.Sessoes)
	}
}

// TestZerarSessaoInexistenteCitaID prova AC-E01.6: erro cita o ID e nada muda.
func TestZerarSessaoInexistenteCitaID(t *testing.T) {
	dir := t.TempDir()
	salvarSeValido(t, dir, MemoriaExecucao{Processo: "p", SessaoID: "sessao-viva", PID: os.Getpid()})
	caminho := filepath.Join(dir, ArquivoMemoria)
	antes, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ReadFile falhou: %v", err)
	}
	err = ZerarSessao(dir, "sessao-fantasma")
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

// TestLimparOrfaosLimpaSessaoELock prova AC-E01.5: PID morto limpa sessão e lock.
func TestLimparOrfaosLimpaSessaoELock(t *testing.T) {
	dir := t.TempDir()
	morto := pidMortoParaTeste(t)
	salvarSeValido(t, dir, MemoriaExecucao{Processo: "p", SessaoID: "sessao-orfa", PID: morto, EtapaAtual: "parada"})
	salvarSeValido(t, dir, MemoriaExecucao{Processo: "p", SessaoID: "sessao-viva", PID: os.Getpid(), EtapaAtual: "meio"})
	lockOrfao := filepath.Join(dir, ArquivoLockMemoria)
	if err := os.WriteFile(lockOrfao, []byte(fmt.Sprintf("%d sessao-orfa", morto)), 0o644); err != nil {
		t.Fatalf("escrita do lock órfão falhou: %v", err)
	}
	n, err := LimparOrfaos(dir)
	if err != nil {
		t.Fatalf("LimparOrfaos falhou: %v", err)
	}
	if n != 1 {
		t.Fatalf("esperava 1 órfã removida, obtido %d", n)
	}
	est, err := lerEstadoDisco(dir)
	if err != nil {
		t.Fatalf("lerEstadoDisco falhou: %v", err)
	}
	if _, ok := est.Sessoes["sessao-orfa"]; ok {
		t.Errorf("sessão órfã deveria sair do envelope: %v", est.Sessoes)
	}
	if got, ok := est.Sessoes["sessao-viva"]; !ok || got.EtapaAtual != "meio" {
		t.Errorf("sessão viva corrompida: %+v", est.Sessoes)
	}
	if _, err := os.Stat(lockOrfao); !os.IsNotExist(err) {
		t.Errorf("lock órfão deveria ser removido: %v", err)
	}
}

// TestMemoriaSegueSchemaV1 prova a conformidade do envelope com o schema memoria_execucao.v1
// (fonte da persistência, TICK-INT-02): schema_version, processo e sessoes sempre gravados.
func TestMemoriaSegueSchemaV1(t *testing.T) {
	dir := t.TempDir()
	if err := SalvarMemoria(dir, MemoriaExecucao{Processo: "proc-x", SessaoID: "s1", PID: os.Getpid(), EtapaAtual: "TRI-001"}); err != nil {
		t.Fatalf("SalvarMemoria falhou: %v", err)
	}
	dados, err := os.ReadFile(filepath.Join(dir, ArquivoMemoria))
	if err != nil {
		t.Fatalf("lendo envelope: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(dados, &raw); err != nil {
		t.Fatalf("envelope não é JSON: %v", err)
	}
	if raw["schema_version"] != SchemaMemoria {
		t.Errorf("schema_version = %v, esperado %q", raw["schema_version"], SchemaMemoria)
	}
	if raw["processo"] != "proc-x" {
		t.Errorf("processo = %v, esperado proc-x", raw["processo"])
	}
	sessoes, ok := raw["sessoes"].(map[string]any)
	if !ok {
		t.Fatalf("sessoes ausente ou não é objeto: %v", raw["sessoes"])
	}
	if _, ok := sessoes["s1"]; !ok {
		t.Errorf("sessão s1 ausente do envelope: %v", sessoes)
	}
}

// TestLockOrfaoTomadoNaGravacao prova AC-E03.2 do mega plano: lock retido por PID
// morto é tomado pela segunda sessão e o atualizado avança (ISO 8601).
func TestLockOrfaoTomadoNaGravacao(t *testing.T) {
	dir := t.TempDir()
	antigo := time.Now().UTC().Add(-2 * time.Hour)
	if err := SalvarMemoria(dir, MemoriaExecucao{Processo: "p", SessaoID: "s2", PID: os.Getpid(), AtualizadoEm: antigo}); err != nil {
		t.Fatalf("primeira gravação: %v", err)
	}
	morto := pidMortoParaTeste(t)
	lockOrfao := filepath.Join(dir, ArquivoLockMemoria)
	if err := os.WriteFile(lockOrfao, []byte(fmt.Sprintf("%d sessao-morto", morto)), 0o644); err != nil {
		t.Fatalf("preparação do lock órfão: %v", err)
	}
	if err := SalvarMemoria(dir, MemoriaExecucao{Processo: "p", SessaoID: "s2", PID: os.Getpid(), EtapaAtual: "retomada"}); err != nil {
		t.Fatalf("gravação com lock órfão devia tomar o lock: %v", err)
	}
	mem, err := LerMemoria(dir, "s2")
	if err != nil {
		t.Fatalf("LerMemoria: %v", err)
	}
	if !mem.AtualizadoEm.After(antigo) {
		t.Errorf("atualizado não avançou: %s", mem.AtualizadoEm)
	}
	if mem.EtapaAtual != "retomada" {
		t.Errorf("etapa = %q, esperado retomada", mem.EtapaAtual)
	}
}

// TestMemoriaCorrompidaSentinelaBackup prova AC-E03.4 do mega plano: JSON inválido
// devolve erro sentinela, cria backup .legoloop-bak e não permite escrita cega por cima.
func TestMemoriaCorrompidaSentinelaBackup(t *testing.T) {
	dir := t.TempDir()
	caminho := filepath.Join(dir, ArquivoMemoria)
	lixo := []byte("{nao é json !!")
	if err := os.WriteFile(caminho, lixo, 0o644); err != nil {
		t.Fatalf("preparação: %v", err)
	}
	_, err := LerMemoria(dir, "s1")
	if !errors.Is(err, ErrMemoriaCorrompida) {
		t.Fatalf("esperado ErrMemoriaCorrompida, obtido: %v", err)
	}
	backup, err := os.ReadFile(caminho + SufixoBackup)
	if err != nil {
		t.Fatalf("backup .legoloop-bak não criado: %v", err)
	}
	if string(backup) != string(lixo) {
		t.Errorf("backup difere do original: %s", backup)
	}
	depois, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("original sumiu: %v", err)
	}
	if string(depois) != string(lixo) {
		t.Errorf("original foi sobrescrito: %s", depois)
	}
	// Nenhuma escrita cega por cima: gravar sobre corrompido sinaliza e preserva.
	if err := SalvarMemoria(dir, MemoriaExecucao{Processo: "p", SessaoID: "s9", PID: os.Getpid()}); !errors.Is(err, ErrMemoriaCorrompida) {
		t.Errorf("esperado ErrMemoriaCorrompida na gravação, obtido: %v", err)
	}
}
