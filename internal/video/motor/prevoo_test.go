package motor

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// depsFalsas devolve dependências com tudo instalado, exceto o item em ausente (AC-E10.2).
func depsFalsas(t *testing.T, ausente string) Dependencias {
	t.Helper()
	motorDir := t.TempDir()
	if ausente == "motor-local" {
		motorDir = motorDir + "-inexistente"
	}
	return Dependencias{
		MotorDir: motorDir,
		Procurar: func(nome string) (string, error) {
			if nome == ausente {
				return "", os.ErrNotExist
			}
			return "/bin/" + nome, nil
		},
		Rodar: func(nome string, args ...string) error {
			if ausente != "" && strings.Contains(strings.Join(args, " "), ausente) {
				return os.ErrNotExist
			}
			return nil
		},
	}
}

func TestPreVooAprovaComTudoInstalado(t *testing.T) {
	if err := PreVoo(depsFalsas(t, "")); err != nil {
		t.Fatalf("pré-voo com tudo instalado deveria aprovar: %v", err)
	}
}

func TestPreVooAusenciaRetornaErrMotorAusente(t *testing.T) {
	casos := []struct {
		ausente   string
		instrucao string
	}{
		{"motor-local", "motor-local"},
		{"python", "python"},
		{"playwright", "pip install playwright"},
		{"chromium", "playwright install chromium"},
		{"ffmpeg", "ffmpeg"},
	}
	for _, c := range casos {
		err := PreVoo(depsFalsas(t, c.ausente))
		if !errors.Is(err, ErrMotorAusente) {
			t.Errorf("%s ausente: esperado ErrMotorAusente, veio %v", c.ausente, err)
			continue
		}
		if !strings.Contains(err.Error(), "instale") || !strings.Contains(err.Error(), c.instrucao) {
			t.Errorf("%s ausente: erro sem instrução de instalação %q: %v", c.ausente, c.instrucao, err)
		}
	}
}

func TestPreVooDependenciasPadraoPreenchidas(t *testing.T) {
	d := Dependencias{}.comPadroes()
	if d.Procurar == nil || d.Rodar == nil || d.MotorDir == "" {
		t.Fatalf("dependências padrão incompletas: %+v", d)
	}
}
