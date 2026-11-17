package triagem

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestCinquentaWorkers prova AC-NF.2: 50 workers sob carga gravam memória de processos
// distintos. O lock por arquivo não serializa processos distintos (a barreira fecha com
// os 50 locks retidos ao mesmo tempo) e a memória por processo permanece íntegra.
func TestCinquentaWorkers(t *testing.T) {
	const workers = 50
	raiz := t.TempDir()
	dirDe := func(i int) string { return filepath.Join(raiz, fmt.Sprintf("proc-%02d", i)) }

	t.Run("LockPorArquivoNaoSerializaProcessosDistintos", func(t *testing.T) {
		// Barreira determinística: cada worker só libera o lock quando os 50 estão
		// retidos ao mesmo tempo. Lock global reprovaria por timeout sem abrir a barreira.
		entraram := make(chan struct{}, workers)
		barreira := make(chan struct{})
		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				lk, err := adquirirLockMemoria(dirDe(i), os.Getpid(), fmt.Sprintf("sessao-%d", i))
				if err != nil {
					t.Errorf("worker %d: adquirirLockMemoria: %v", i, err)
					return
				}
				defer lk.liberar()
				entraram <- struct{}{}
				<-barreira
			}(i)
		}
		prazo := time.Now().Add(30 * time.Second)
		for retidos := 0; retidos < workers; retidos++ {
			select {
			case <-entraram:
			case <-time.After(time.Until(prazo)):
				t.Fatalf("só %d/%d locks retidos ao mesmo tempo: o lock serializa processos distintos", retidos, workers)
			}
		}
		close(barreira)
		wg.Wait()
	})

	t.Run("MemoriaPorProcessoIntegraSobCarga", func(t *testing.T) {
		const rodadas = 3
		for r := 0; r < rodadas; r++ {
			var wg sync.WaitGroup
			for i := 0; i < workers; i++ {
				wg.Add(1)
				go func(i, r int) {
					defer wg.Done()
					mem := MemoriaExecucao{
						Processo: fmt.Sprintf("proc-%02d", i), SessaoID: fmt.Sprintf("sessao-%d", i),
						PID: os.Getpid(), EtapaAtual: fmt.Sprintf("rodada-%d-worker-%d", r, i),
					}
					if err := SalvarMemoria(dirDe(i), mem); err != nil {
						t.Errorf("worker %d rodada %d: SalvarMemoria: %v", i, r, err)
					}
				}(i, r)
			}
			wg.Wait()
		}
		for i := 0; i < workers; i++ {
			lida, err := LerMemoria(dirDe(i), fmt.Sprintf("sessao-%d", i))
			if err != nil {
				t.Fatalf("worker %d: LerMemoria: %v", i, err)
			}
			esperado := fmt.Sprintf("rodada-%d-worker-%d", rodadas-1, i)
			if lida.Processo != fmt.Sprintf("proc-%02d", i) || lida.EtapaAtual != esperado {
				t.Errorf("worker %d: memória corrompida sob carga: %+v (quero etapa %q)", i, lida, esperado)
			}
		}
	})
}
