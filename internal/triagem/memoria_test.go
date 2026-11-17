package triagem

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
