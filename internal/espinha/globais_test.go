package espinha

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGlobaisDaEspinhaEstaoJustificados(t *testing.T) {
	t.Run("repositorio real tem globais e inits justificados", func(t *testing.T) {
		raiz := filepath.Join("..", "..")
		reg, err := CarregarRegistro(filepath.Join(raiz, "espinha", "registro.json"))
		if err != nil {
			t.Fatalf("carregar registro: %v", err)
		}
		for _, v := range verificarGlobaisMutaveis(raiz, reg.GlobaisJustificadas) {
			t.Errorf("global mutavel sem justificativa: %s", v.String())
		}
		for _, v := range verificarInits(raiz, reg.InitsJustificados) {
			t.Errorf("init sem justificativa: %s", v.String())
		}
	})

	t.Run("mutacao global mutavel sem justificativa falha", func(t *testing.T) {
		raiz := criarFixtureGlobais(t, `package demo
var estadoMutavel = map[string]int{"a": 1}
var fatiaMutavel = []string{"x"}
var ponteiroMutavel = new(int)
`)
		vs := verificarGlobaisMutaveis(raiz, nil)
		if len(vs) != 3 {
			t.Fatalf("esperado 3 violacoes de global mutavel, obtido %d: %v", len(vs), vs)
		}
		for _, v := range vs {
			if v.Regra != "global-mutavel-nao-justificada" {
				t.Errorf("regra inesperada: %s", v.Regra)
			}
		}
	})

	t.Run("mutacao init sem justificativa falha", func(t *testing.T) {
		raiz := criarFixtureGlobais(t, `package demo
func init() {}
`)
		vs := verificarInits(raiz, nil)
		if len(vs) != 1 || vs[0].Regra != "init-nao-justificado" || vs[0].Alvo != "internal/demo/demo.go" {
			t.Fatalf("esperado 1 violacao init-nao-justificado em internal/demo/demo.go, obtido: %v", vs)
		}
		just := map[string]string{"internal/demo/demo.go": "registro estatico de teste"}
		if vsOK := verificarInits(raiz, just); len(vsOK) != 0 {
			t.Fatalf("esperado 0 violacoes quando init esta justificado, obtido: %v", vsOK)
		}
	})

	t.Run("constantes e imutaveis por construcao passam", func(t *testing.T) {
		raiz := criarFixtureGlobais(t, `package demo
import "errors"
const limite = 10
var errPadrao = errors.New("erro imutavel")
`)
		if vs := verificarGlobaisMutaveis(raiz, nil); len(vs) != 0 {
			t.Fatalf("esperado 0 violacoes para imutaveis, obtido: %v", vs)
		}
		if vs := verificarInits(raiz, nil); len(vs) != 0 {
			t.Fatalf("esperado 0 violacoes sem init, obtido: %v", vs)
		}
	})
}

func criarFixtureGlobais(t *testing.T, codigo string) string {
	t.Helper()
	raiz := t.TempDir()
	dir := filepath.Join(raiz, "internal", "demo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("criar diretorio fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "demo.go"), []byte(codigo), 0o644); err != nil {
		t.Fatalf("gravar arquivo fixture: %v", err)
	}
	return raiz
}
