package processo_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/processo"
)

func TestContratoValidate(t *testing.T) {
	valido := processo.Contrato{
		Nome: "coletor-csv",
		Tipo: "script",
		Sai:  []processo.Saida{{Tipo: "dados", CaminhoPadrao: "out/coleta.csv", Promete: true, GateSugerido: "line_count('out/coleta.csv') >= 2"}},
	}
	casos := []struct {
		nome    string
		muda    func(*processo.Contrato)
		querErr bool
	}{
		{"contrato valido passa", func(*processo.Contrato) {}, false},
		{"sem nome e recusado", func(c *processo.Contrato) { c.Nome = " " }, true},
		{"sem saida e recusado", func(c *processo.Contrato) { c.Sai = nil }, true},
		{"saida sem tipo e recusada", func(c *processo.Contrato) { c.Sai[0].Tipo = "" }, true},
		{"caminho_padrao com .. e recusado", func(c *processo.Contrato) { c.Sai[0].CaminhoPadrao = "../coleta.csv" }, true},
		{"caminho_padrao absoluto e recusado", func(c *processo.Contrato) { c.Sai[0].CaminhoPadrao = "/tmp/coleta.csv" }, true},
		{"gate_sugerido que nao compila e recusado", func(c *processo.Contrato) { c.Sai[0].GateSugerido = "line_count(" }, true},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			contrato := valido
			contrato.Sai = append([]processo.Saida(nil), valido.Sai...)
			c.muda(&contrato)
			err := contrato.Validate()
			if c.querErr && !errors.Is(err, processo.ErrContratoInvalido) {
				t.Fatalf("Validate() = %v; queria ErrContratoInvalido", err)
			}
			if !c.querErr && err != nil {
				t.Fatalf("Validate() = %v; queria nil", err)
			}
		})
	}
}

func TestContratoDe(t *testing.T) {
	dir := t.TempDir()
	gravar := func(t *testing.T, conteudo string) string {
		t.Helper()
		caminho := filepath.Join(dir, "contrato.json")
		if err := os.WriteFile(caminho, []byte(conteudo), 0o600); err != nil {
			t.Fatal(err)
		}
		return caminho
	}
	t.Run("carrega manifesto valido", func(t *testing.T) {
		c, err := processo.ContratoDe(gravar(t, `{"nome":"coletor-csv","tipo":"script","sai":[{"tipo":"dados","caminho_padrao":"out/coleta.csv","promete":true}]}`))
		if err != nil || c.Nome != "coletor-csv" || c.Sai[0].CaminhoPadrao != "out/coleta.csv" {
			t.Fatalf("ContratoDe = %+v, %v", c, err)
		}
	})
	t.Run("sai vazio e recusado sem criar arquivo", func(t *testing.T) {
		caminho := gravar(t, `{"nome":"vazio","tipo":"script","sai":[]}`)
		if _, err := processo.ContratoDe(caminho); !errors.Is(err, processo.ErrContratoInvalido) {
			t.Fatalf("ContratoDe = %v; queria ErrContratoInvalido", err)
		}
		entradas, err := os.ReadDir(dir)
		if err != nil || len(entradas) != 1 {
			t.Fatalf("dir = %v, %v; queria só o contrato.json", entradas, err)
		}
	})
	t.Run("caminho_padrao com .. e recusado", func(t *testing.T) {
		caminho := gravar(t, `{"nome":"hostil","tipo":"script","sai":[{"tipo":"dados","caminho_padrao":"../fuga.csv"}]}`)
		if _, err := processo.ContratoDe(caminho); !errors.Is(err, processo.ErrContratoInvalido) {
			t.Fatalf("ContratoDe = %v; queria ErrContratoInvalido", err)
		}
	})
	t.Run("JSON quebrado e ErrContratoInvalido", func(t *testing.T) {
		if _, err := processo.ContratoDe(gravar(t, `{"nome":`)); !errors.Is(err, processo.ErrContratoInvalido) {
			t.Fatalf("ContratoDe = %v; queria ErrContratoInvalido", err)
		}
	})
	t.Run("arquivo inexistente nao e ErrContratoInvalido", func(t *testing.T) {
		_, err := processo.ContratoDe(filepath.Join(dir, "nao-existe.json"))
		if err == nil || errors.Is(err, processo.ErrContratoInvalido) {
			t.Fatalf("ContratoDe = %v; queria erro de E/S sem ErrContratoInvalido", err)
		}
	})
}
