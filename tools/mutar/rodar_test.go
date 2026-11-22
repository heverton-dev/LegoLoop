package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const fixtureCalcCodigo = `package calc

func Soma(a, b int) int {
	return a + b
}

func EhMaior(a, b int) bool {
	if a > b {
		return true
	}
	return false
}
`

const fixtureCalcTesteCompleto = `package calc

import "testing"

func TestSoma(t *testing.T) {
	if got := Soma(3, 2); got != 5 {
		t.Fatalf("Soma(3,2) = %d, esperado 5", got)
	}
}

func TestEhMaior(t *testing.T) {
	if !EhMaior(5, 3) {
		t.Fatal("esperava EhMaior(5,3) == true")
	}
	if EhMaior(3, 3) {
		t.Fatal("esperava EhMaior(3,3) == false")
	}
	if EhMaior(2, 3) {
		t.Fatal("esperava EhMaior(2,3) == false")
	}
}
`

const fixtureCalcTesteFraco = `package calc

import "testing"

func TestSoma(t *testing.T) {
	if got := Soma(3, 2); got != 5 {
		t.Fatalf("Soma(3,2) = %d, esperado 5", got)
	}
}

func TestEhMaiorFraco(t *testing.T) {
	if !EhMaior(5, 3) {
		t.Fatal("esperava EhMaior(5,3) == true")
	}
	if EhMaior(2, 3) {
		t.Fatal("esperava EhMaior(2,3) == false")
	}
}
`

func TestSobreviventeReprova(t *testing.T) {
	binario := compilarBinarioMutar(t)

	t.Run("fixture_A_todos_mortos_exit_0", func(t *testing.T) {
		dir := criarFixtureCalc(t, fixtureCalcTesteCompleto)
		codigo, saida := rodarBinarioMutar(t, binario, dir, ".")
		if codigo != 0 {
			t.Fatalf("fixture A: esperava exit 0, obtido %d\nsaida:\n%s", codigo, saida)
		}
		if !strings.Contains(saida, "sobreviventes=0") {
			t.Fatalf("fixture A: esperava sobreviventes=0 na saida:\n%s", saida)
		}
		verificarFonteIntacto(t, filepath.Join(dir, "calc.go"), fixtureCalcCodigo)
	})

	t.Run("fixture_B_sobrevivente_reprova_exit_2", func(t *testing.T) {
		dir := criarFixtureCalc(t, fixtureCalcTesteFraco)
		codigo, saida := rodarBinarioMutar(t, binario, dir, ".")
		if codigo != 2 {
			t.Fatalf("fixture B: esperava exit 2, obtido %d\nsaida:\n%s", codigo, saida)
		}
		for _, esperado := range []string{
			"sobrevivente:",
			"id=",
			"arquivo=calc.go",
			"linha=8",
			"operador=gt_geq",
		} {
			if !strings.Contains(saida, esperado) {
				t.Errorf("fixture B: saida nao contem %q:\n%s", esperado, saida)
			}
		}
		verificarFonteIntacto(t, filepath.Join(dir, "calc.go"), fixtureCalcCodigo)
	})

	t.Run("fixture_B_com_justificados_exit_0", func(t *testing.T) {
		dir := criarFixtureCalc(t, fixtureCalcTesteFraco)
		muts, err := Mutantes(dir, OperadoresPadrao())
		if err != nil {
			t.Fatalf("Mutantes falhou: %v", err)
		}
		var alvo *Mutante
		for i := range muts {
			if muts[i].Operador == "gt_geq" && muts[i].Linha == 8 {
				alvo = &muts[i]
				break
			}
		}
		if alvo == nil {
			t.Fatal("mutante gt_geq na linha 8 nao encontrado na fixture B")
		}
		arqJust := filepath.Join(dir, "justificados.json")
		lista := []Justificado{{
			ID:            ".:calc.go:8:gt_geq:1",
			Pacote:        ".",
			Arquivo:       alvo.Arquivo,
			Linha:         alvo.Linha,
			Operador:      alvo.Operador,
			Anticorpo:     "ANT-CALC-001",
			Justificativa: "fronteira de igualdade aceita por especificacao",
		}}
		dados, err := json.MarshalIndent(lista, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(arqJust, dados, 0o600); err != nil {
			t.Fatal(err)
		}

		codigo, saida := rodarBinarioMutar(t, binario, dir, ".", "--justificados", arqJust)
		if codigo != 0 {
			t.Fatalf("fixture B com justificados: esperava exit 0, obtido %d\nsaida:\n%s", codigo, saida)
		}
		if !strings.Contains(saida, "justificados=1") || !strings.Contains(saida, "sobreviventes=0") {
			t.Fatalf("fixture B com justificados: resumo inesperado:\n%s", saida)
		}
	})
}

func compilarBinarioMutar(t *testing.T) string {
	t.Helper()
	nome := "mutar"
	if runtime.GOOS == "windows" {
		nome += ".exe"
	}
	binPath := filepath.Join(t.TempDir(), nome)
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("falha ao compilar tools/mutar: %v\n%s", err, string(out))
	}
	return binPath
}

func criarFixtureCalc(t *testing.T, codigoTeste string) string {
	t.Helper()
	dir := t.TempDir()
	goMod := "module exemplo.local/calc\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "calc.go"), []byte(fixtureCalcCodigo), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "calc_test.go"), []byte(codigoTeste), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

func rodarBinarioMutar(t *testing.T, binario, dir string, args ...string) (int, string) {
	t.Helper()
	cmd := exec.Command(binario, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err == nil {
		return 0, string(out)
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), string(out)
	}
	t.Fatalf("erro inesperado ao executar %s: %v\n%s", binario, err, string(out))
	return -1, string(out)
}

func verificarFonteIntacto(t *testing.T, caminho, esperado string) {
	t.Helper()
	atual, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("erro ao ler fonte %s: %v", caminho, err)
	}
	if string(atual) != esperado {
		t.Fatalf("arquivo fonte %s foi alterado apos mutacao", caminho)
	}
}

func TestAmbienteSubteste_IncluiGoMemLimitEGoMaxProcs(t *testing.T) {
	env := ambienteSubteste([]string{"PATH=/usr/bin"})
	if !contemEnv(env, "GOMAXPROCS=2") {
		t.Fatalf("esperava GOMAXPROCS=2 no ambiente: %v", env)
	}
	if !contemEnv(env, "GOMEMLIMIT=512MiB") {
		t.Fatalf("esperava GOMEMLIMIT=512MiB no ambiente: %v", env)
	}

	envCustom := ambienteSubteste([]string{"GOMAXPROCS=4", "GOMEMLIMIT=256MiB"})
	if !contemEnv(envCustom, "GOMAXPROCS=4") || contemEnv(envCustom, "GOMAXPROCS=2") {
		t.Fatalf("deveria respeitar GOMAXPROCS existente: %v", envCustom)
	}
	if !contemEnv(envCustom, "GOMEMLIMIT=256MiB") || contemEnv(envCustom, "GOMEMLIMIT=512MiB") {
		t.Fatalf("deveria respeitar GOMEMLIMIT existente: %v", envCustom)
	}
}

func TestArgsGoTest_IncluiTimeoutPorMutante(t *testing.T) {
	args := argsGoTest("/tmp/overlay.json", 15*time.Second)
	junto := strings.Join(args, " ")
	for _, esperado := range []string{"-p 1", "-count=1", "-failfast", "-timeout 15s", "-overlay /tmp/overlay.json"} {
		if !strings.Contains(junto, esperado) {
			t.Fatalf("argsGoTest nao contem %q: %v", esperado, args)
		}
	}
}

func TestResolverTimeoutMutante(t *testing.T) {
	t.Setenv("MUTAR_TIMEOUT_MUTANTE", "")
	if got := resolverTimeoutMutante(); got != timeoutMutantePadrao {
		t.Fatalf("padrao esperado %v, obtido %v", timeoutMutantePadrao, got)
	}
	t.Setenv("MUTAR_TIMEOUT_MUTANTE", "7s")
	if got := resolverTimeoutMutante(); got != 7*time.Second {
		t.Fatalf("esperado 7s, obtido %v", got)
	}
	t.Setenv("MUTAR_TIMEOUT_MUTANTE", "invalido")
	if got := resolverTimeoutMutante(); got != timeoutMutantePadrao {
		t.Fatalf("esperado fallback %v para valor invalido, obtido %v", timeoutMutantePadrao, got)
	}
}

func TestRodarGoTest_MataMutanteEmLoopInfinito(t *testing.T) {
	t.Setenv("MUTAR_TIMEOUT_MUTANTE", "2s")
	dir := t.TempDir()
	goMod := "module exemplo.local/loop\n\ngo 1.22\n"
	loopTest := `package loop

import "testing"

func TestLoopInfinito(t *testing.T) {
	var lixo [][]byte
	for {
		lixo = append(lixo, make([]byte, 1<<20))
		if len(lixo) > 200 {
			lixo = lixo[:0]
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "loop_test.go"), []byte(loopTest), 0o600); err != nil {
		t.Fatal(err)
	}
	inicio := time.Now()
	passou, _, err := rodarGoTest(context.Background(), dir, "")
	duracao := time.Since(inicio)
	if err != nil {
		t.Fatalf("erro inesperado (esperava mutante morto por timeout): %v", err)
	}
	if passou {
		t.Fatal("esperava passou=false para teste em loop infinito")
	}
	if duracao > 15*time.Second {
		t.Fatalf("rodarGoTest demorou demais (%v) para encerrar loop infinito", duracao)
	}
}

func contemEnv(env []string, alvo string) bool {
	for _, kv := range env {
		if kv == alvo {
			return true
		}
	}
	return false
}
