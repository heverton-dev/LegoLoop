package motor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// licencaMotor espelha motor-local/LICENCA.json (AC-E10.1).
type licencaMotor struct {
	Status string `json:"status"`
	Motor  string `json:"motor"`
}

// exigeIgnorado decide se motor-local precisa ficar fora do git.
// pendente exige ignorar; liberada ou reescrito permite mover para internal/video/motor/.
func exigeIgnorado(status string) (bool, error) {
	switch status {
	case "pendente":
		return true, nil
	case "liberada", "reescrito":
		return false, nil
	}
	return false, fmt.Errorf("status de licença desconhecido: %q", status)
}

// raizRepo sobe a partir do diretório do teste até achar go.mod.
func raizRepo(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			t.Fatal("go.mod não encontrado acima do pacote")
		}
		dir = pai
	}
}

func lerLicenca(t *testing.T, raiz string) licencaMotor {
	t.Helper()
	dados, err := os.ReadFile(filepath.Join(raiz, "motor-local", "LICENCA.json"))
	if err != nil {
		t.Fatalf("ler LICENCA.json: %v", err)
	}
	var lic licencaMotor
	if err := json.Unmarshal(dados, &lic); err != nil {
		t.Fatalf("LICENCA.json inválido: %v", err)
	}
	return lic
}

func TestLicencaDecisaoPorStatus(t *testing.T) {
	casos := []struct {
		status  string
		ignorar bool
		erro    bool
	}{
		{"pendente", true, false},
		{"liberada", false, false},
		{"reescrito", false, false},
		{"", false, true},
		{"negada", false, true},
	}
	for _, c := range casos {
		ignorar, err := exigeIgnorado(c.status)
		if (err != nil) != c.erro {
			t.Errorf("status %q: erro=%v, esperado erro=%v", c.status, err, c.erro)
		}
		if ignorar != c.ignorar {
			t.Errorf("status %q: ignorar=%v, esperado %v", c.status, ignorar, c.ignorar)
		}
	}
}

func TestLicencaMotorLocalGovernanca(t *testing.T) {
	raiz := raizRepo(t)
	lic := lerLicenca(t, raiz)
	if lic.Motor == "" {
		t.Error("LICENCA.json sem o campo motor")
	}
	ignorar, err := exigeIgnorado(lic.Status)
	if err != nil {
		t.Fatal(err)
	}
	if !ignorar {
		t.Logf("licença %q: motor-local pode migrar para internal/video/motor/", lic.Status)
		return
	}
	cmd := exec.Command("git", "check-ignore", "-q", "motor-local")
	cmd.Dir = raiz
	if err := cmd.Run(); err != nil {
		t.Fatalf("licença pendente exige motor-local ignorado no git: git check-ignore falhou: %v", err)
	}
}
