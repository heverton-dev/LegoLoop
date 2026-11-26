// Package motor guarda o pré-voo e a integração do motor local de render (AC-E10).
package motor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// ErrMotorAusente indica que o motor ou uma dependência de render não está instalada (AC-E10.2).
var ErrMotorAusente = errors.New("motor: dependência de render ausente")

// scriptChromium sai com código 0 só se o Chromium do Playwright existir em disco.
const scriptChromium = "import os,sys\n" +
	"from playwright.sync_api import sync_playwright\n" +
	"with sync_playwright() as p: sys.exit(0 if os.path.exists(p.chromium.executable_path) else 1)"

// Dependencias injeta o acesso ao sistema para o pré-voo ser testável.
type Dependencias struct {
	MotorDir string
	Procurar func(nome string) (string, error)
	Rodar    func(nome string, args ...string) error
}

// checagem descreve um item do pré-voo e a instrução de instalação.
type checagem struct {
	nome      string
	instrucao string
	verificar func(d Dependencias) error
}

// comPadroes preenche os campos vazios com o sistema real.
func (d Dependencias) comPadroes() Dependencias {
	if d.MotorDir == "" {
		d.MotorDir = "motor-local"
	}
	if d.Procurar == nil {
		d.Procurar = exec.LookPath
	}
	if d.Rodar == nil {
		d.Rodar = func(nome string, args ...string) error {
			return exec.Command(nome, args...).Run()
		}
	}
	return d
}

// checagens lista os itens na ordem de dependência: motor, python, Playwright, Chromium, ffmpeg.
func checagens() []checagem {
	return []checagem{
		{"motor-local", "clone o motor em motor-local/ com LICENCA.json", func(d Dependencias) error {
			_, err := os.Stat(d.MotorDir)
			return err
		}},
		{"python", "Python 3.10+ no PATH (https://www.python.org/downloads/)", func(d Dependencias) error {
			_, err := d.Procurar("python")
			return err
		}},
		{"playwright", "pip install playwright", func(d Dependencias) error {
			return d.Rodar("python", "-c", "import playwright")
		}},
		{"chromium", "python -m playwright install chromium", func(d Dependencias) error {
			return d.Rodar("python", "-c", scriptChromium)
		}},
		{"ffmpeg", "ffmpeg no PATH (https://ffmpeg.org/download.html)", func(d Dependencias) error {
			_, err := d.Procurar("ffmpeg")
			return err
		}},
	}
}

// PreVoo aprova o render só com motor, python, Playwright, Chromium e ffmpeg presentes.
func PreVoo(d Dependencias) error {
	d = d.comPadroes()
	for _, c := range checagens() {
		if err := c.verificar(d); err != nil {
			return fmt.Errorf("%w: %s (%v); instale: %s", ErrMotorAusente, c.nome, err, c.instrucao)
		}
	}
	return nil
}
