package board

import (
	"os"
	"os/exec"
	"testing"
)

// deadPID devolve o PID de um processo real que ja terminou.
func deadPID(t *testing.T) int {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	return cmd.Process.Pid
}

func TestLiveViewMarcaPIDMortoSemMutarEntrada(t *testing.T) {
	in := []Worker{
		{Name: "w1", Status: WorkerWorking, PID: deadPID(t)},
		{Name: "w2", Status: WorkerKilled},
		{Name: "w3", Status: WorkerWorking, PID: os.Getpid()},
	}
	out := LiveView(in)
	if out[0].Status != WorkerStopped || out[0].PID != 0 {
		t.Errorf("w1 deveria virar stopped com PID 0: %+v", out[0])
	}
	if out[1].Status != WorkerKilled || out[2].Status != WorkerWorking || out[2].PID != os.Getpid() {
		t.Errorf("w2 e w3 nao deviam mudar: %+v %+v", out[1], out[2])
	}
	if in[0].Status != WorkerWorking || in[0].PID == 0 {
		t.Errorf("a entrada nao pode ser mutada: %+v", in[0])
	}
	if got := LiveView(nil); got == nil || len(got) != 0 {
		t.Errorf("LiveView(nil) deve ser lista vazia nao nula: %#v", got)
	}
}
