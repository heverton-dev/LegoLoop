package board

import "testing"

// TestAddWorkerReusaNomeDeWorkerEncerrado cobre o achado A18: rerodar um no do DAG falhava porque o
// worker encerrado (killed) continuava ocupando o nome. Worker vivo continua recusando duplicata.
func TestAddWorkerReusaNomeDeWorkerEncerrado(t *testing.T) {
	sb, err := OpenSwarmBoard(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	req := AddWorkerRequest{Name: "w1", Role: "engenheiro-tdd", Model: "m1"}
	if _, err := sb.AddWorker(req); err != nil {
		t.Fatalf("primeiro AddWorker: %v", err)
	}
	if _, err := sb.AddWorker(req); err == nil {
		t.Fatal("worker vivo deveria recusar nome duplicado")
	}
	if err := sb.KillWorker("w1", "."); err != nil {
		t.Fatalf("KillWorker: %v", err)
	}
	if _, err := sb.AddWorker(req); err != nil {
		t.Fatalf("nome de worker encerrado deveria liberar: %v", err)
	}
	st, err := sb.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Workers) != 1 || st.Workers[0].Status != WorkerWaiting {
		t.Fatalf("esperava 1 worker waiting, veio %+v", st.Workers)
	}
}
