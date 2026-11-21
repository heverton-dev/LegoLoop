package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const fixtureOperadores = `package alvo

import "errors"

func Igualdade(a, b int) bool {
	if a == b {
		return true
	}
	return a != b
}

func Menor(a, b int) bool {
	if a < b {
		return a <= b
	}
	return false
}

func Maior(a, b int) bool {
	if a > b {
		return a >= b
	}
	return false
}

func Logica(a, b bool) bool {
	if a && b {
		return true
	}
	return a || b
}

func Booleano(flag bool) bool {
	if flag {
		return true
	}
	return false
}

func RetornoErro(falhar bool) (int, error) {
	err := errors.New("falha")
	if falhar {
		if err != nil {
			return 0, err
		}
	}
	return 1, nil
}

func RetornoErroDireto() error {
	err := errors.New("direto")
	return err
}

func AritmeticaEString(a, b int, s1, s2 string) (int, string) {
	soma := a + b
	sub := a - b
	concatLit := "ola" + " mundo"
	concatVar := s1 + s2
	return soma + sub, concatLit + concatVar
}
`

func TestOperadores(t *testing.T) {
	dir := t.TempDir()
	escreverFixturePacote(t, dir)

	ops := OperadoresPadrao()
	esperados := []string{
		"eq_neq",
		"lt_leq",
		"gt_geq",
		"and_or",
		"bool_flip",
		"rm_return_err",
		"add_sub",
	}
	if len(ops) != len(esperados) {
		t.Fatalf("esperava %d operadores padrao, obtido %d", len(esperados), len(ops))
	}

	muts1, err := Mutantes(dir, ops)
	if err != nil {
		t.Fatalf("Mutantes falhou: %v", err)
	}
	muts2, err := Mutantes(dir, ops)
	if err != nil {
		t.Fatalf("segunda chamada Mutantes falhou: %v", err)
	}
	if len(muts1) != len(muts2) {
		t.Fatalf("geracao nao deterministica: %d vs %d", len(muts1), len(muts2))
	}

	porOperador := make(map[string]int, len(esperados))
	ids := make(map[string]bool, len(muts1))
	for i, m := range muts1 {
		if m.ID != muts2[i].ID || string(m.Codigo) != string(muts2[i].Codigo) {
			t.Fatalf("mutante %d nao deterministico: %q vs %q", i, m.ID, muts2[i].ID)
		}
		if ids[m.ID] {
			t.Fatalf("ID duplicado: %s", m.ID)
		}
		ids[m.ID] = true
		validarFormatoID(t, m)
		if len(m.Codigo) == 0 {
			t.Fatalf("mutante %s sem codigo gerado", m.ID)
		}
		porOperador[m.Operador]++
		compilarMutante(t, m)
	}

	for _, nome := range esperados {
		if porOperador[nome] == 0 {
			t.Errorf("operador %s nao gerou nenhum mutante", nome)
		}
	}
	t.Logf("mutantes compilaveis por operador: %v (total=%d)", porOperador, len(muts1))
}

func escreverFixturePacote(t *testing.T, dir string) {
	t.Helper()
	goMod := "module exemplo.local/alvo\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "alvo.go"), []byte(fixtureOperadores), 0o600); err != nil {
		t.Fatal(err)
	}
}

func validarFormatoID(t *testing.T, m Mutante) {
	t.Helper()
	prefixo := m.Pacote + ":" + m.Arquivo + ":"
	if !strings.HasPrefix(m.ID, prefixo) {
		t.Fatalf("ID %q fora do formato pacote:arquivo:linha:operador:indice", m.ID)
	}
	resto := strings.TrimPrefix(m.ID, prefixo)
	partes := strings.Split(resto, ":")
	if len(partes) != 3 || partes[1] != m.Operador || partes[0] == "" || partes[2] == "" {
		t.Fatalf("ID %q fora do formato pacote:arquivo:linha:operador:indice", m.ID)
	}
}

func compilarMutante(t *testing.T, m Mutante) {
	t.Helper()
	tmp := t.TempDir()
	goMod := "module exemplo.local/alvo\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, m.Arquivo), m.Codigo, 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "build", ".")
	cmd.Dir = tmp
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("mutante %s (%s linha %d: %s -> %s) nao compilou: %v\nsaida:\n%s\ncodigo:\n%s",
			m.ID, m.Operador, m.Linha, m.Original, m.Mutado, err, string(out), string(m.Codigo))
	}
}
