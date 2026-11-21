package espinha_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var padraoLinhaViolacao = regexp.MustCompile(`^regra: ([a-z0-9-]+); alvo: (.+); motivo: (.+)$`)

func TestCadaRegraSaiComExitDois(t *testing.T) {
	bin := compilarBinarioEspinha(t)
	casos := []struct {
		regra    string
		preparar func(t *testing.T, raiz string)
	}{
		{"verbo-sem-registro", quebrarVerboSemRegistro},
		{"verbo-orfao", quebrarVerboOrfao},
		{"lei-sem-verificador", quebrarLeiSemVerificador},
		{"funcao-duplicada-entre-pacotes", quebrarFuncaoDuplicada},
		{"processo-incompleto", quebrarProcessoIncompleto},
		{"processo-sem-prefixo", quebrarProcessoSemPrefixo},
		{"processo-duplicado", quebrarProcessoDuplicado},
		{"executor-incoerente", quebrarExecutorIncoerente},
		{"tipo-de-no-fora-da-tabela", quebrarTipoDeNo},
		{"erro-ignorado-acima-do-limite", quebrarErroIgnorado},
		{"global-mutavel-nao-justificada", quebrarGlobalMutavel},
		{"comando-exclusivo-de-so", quebrarComandoExclusivoSO},
	}
	if len(casos) != 12 {
		t.Fatalf("esperado 12 regras, obtido %d", len(casos))
	}
	for _, tc := range casos {
		t.Run(tc.regra, func(t *testing.T) {
			raiz := criarFixtureBase12(t)
			tc.preparar(t, raiz)
			executarEExigirExitDois(t, bin, raiz, tc.regra)
		})
	}
}

func compilarBinarioEspinha(t *testing.T) string {
	t.Helper()
	nome := "legoloop"
	if runtime.GOOS == "windows" {
		nome += ".exe"
	}
	bin := filepath.Join(t.TempDir(), nome)
	cmd := exec.Command("go", "build", "-o", bin, "../../cmd/legoloop")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("falha ao compilar legoloop: %v\n%s", err, out)
	}
	return bin
}

func executarEExigirExitDois(t *testing.T, bin, raiz, regraEsperada string) {
	t.Helper()
	cmd := exec.Command(bin, "espinha", "verificar", "--raiz", raiz)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 2 {
		t.Fatalf("regra %s: esperado exit 2, obtido err=%v\nstdout:\n%s\nstderr:\n%s",
			regraEsperada, err, stdout.String(), stderr.String())
	}
	encontrou := false
	for _, linha := range strings.Split(strings.TrimSpace(stdout.String()), "\n") {
		linha = strings.TrimSpace(linha)
		if linha == "" {
			continue
		}
		m := padraoLinhaViolacao.FindStringSubmatch(linha)
		if m == nil {
			t.Errorf("linha fora do formato canonico: %q", linha)
			continue
		}
		if m[1] == regraEsperada {
			encontrou = true
		}
	}
	if !encontrou {
		t.Fatalf("regra %s nao encontrada na saida:\n%s", regraEsperada, stdout.String())
	}
}

func criarFixtureBase12(t *testing.T) string {
	t.Helper()
	raiz := t.TempDir()
	gravarArq(t, filepath.Join(raiz, "AGENTS.md"), "## 1. Leis\n\n1. Lei um.\n")
	gravarArq(t, filepath.Join(raiz, "internal", "espinha", "base.go"), "package espinha\n")
	gravarArq(t, filepath.Join(raiz, "internal", "cli", "cli.go"), `package cli
func mapaComandos() map[string]func() {
	return map[string]func(){"espinha": nil}
}
`)
	gravarArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "grafo.dag.json"), dagPortavelJSON("legoloop-alpha"))
	gravarArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "ANTICORPOS.md"), "# Anticorpos\n")
	gravarArq(t, filepath.Join(raiz, "espinha", "registro.json"), registroBaseJSON(`{"slug": "legoloop-alpha"}`))
	return raiz
}

func registroBaseJSON(processos string) string {
	return `{
  "schema": "espinha.v1",
  "verbos": [
    {"nome": "espinha", "camada": "espinha", "motor": "internal/espinha", "contrato": "verifica espinha"}
  ],
  "motores": [
    {"pacote": "internal/espinha", "interface": "Verificar", "deterministico": true}
  ],
  "processos": [` + processos + `],
  "requisitos": [
    {"lei": "Lei 1", "verificador": "internal/espinha", "mecanico": true}
  ]
}`
}

func dagPortavelJSON(nome string) string {
	return `{
  "name": "` + nome + `",
  "trigger": {"type": "manual"},
  "nodes": [
    {
      "id": "etapa1",
      "type": "script",
      "command": ["legoloop", "validate", "--file", "grafo.dag.json"],
      "isolation": {"executor": "orchestrix_hermetic"},
      "routing": {"mode": "auto", "combo": "padrao"},
      "output": {"artifact": "out/etapa1.txt", "gate": "file_exists('out/etapa1.txt')"}
    }
  ]
}`
}

func gravarArq(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func quebrarVerboSemRegistro(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "internal", "cli", "cli.go"), `package cli
func mapaComandos() map[string]func() {
	return map[string]func(){"espinha": nil, "extra": nil}
}
`)
}

func quebrarVerboOrfao(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "internal", "cli", "cli.go"), `package cli
func mapaComandos() map[string]func() {
	return map[string]func(){}
}
`)
}

func quebrarLeiSemVerificador(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "AGENTS.md"), "## 1. Leis\n\n1. Lei um.\n2. Lei dois sem registro.\n")
}

func quebrarFuncaoDuplicada(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "internal", "alfa", "a.go"), "package alfa\nfunc Somar(a, b int) int { return a + b }\n")
	gravarArq(t, filepath.Join(raiz, "internal", "beta", "b.go"), "package beta\nfunc Somar(x, y int) int { return x + y }\n")
}

func quebrarProcessoIncompleto(t *testing.T, raiz string) {
	if err := os.Remove(filepath.Join(raiz, "processos", "legoloop-alpha", "ANTICORPOS.md")); err != nil {
		t.Fatal(err)
	}
}

func quebrarProcessoSemPrefixo(t *testing.T, raiz string) {
	if err := os.MkdirAll(filepath.Join(raiz, "processos", "sem-prefixo"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func quebrarProcessoDuplicado(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "processos", "legoloop-beta", "grafo.dag.json"), dagPortavelJSON("legoloop-alpha"))
	gravarArq(t, filepath.Join(raiz, "processos", "legoloop-beta", "ANTICORPOS.md"), "# Anticorpos\n")
	gravarArq(t, filepath.Join(raiz, "espinha", "registro.json"),
		registroBaseJSON(`{"slug": "legoloop-alpha"}, {"slug": "legoloop-beta"}`))
}

func quebrarExecutorIncoerente(t *testing.T, raiz string) {
	dagRuim := `{
  "name": "legoloop-alpha",
  "trigger": {"type": "manual"},
  "nodes": [
    {
      "id": "sg1",
      "type": "stop_gate",
      "isolation": {"executor": "orchestrix_hermetic"},
      "routing": {"mode": "auto", "combo": "padrao"},
      "output": {"artifact": "out/sg.md", "gate": "file_exists('out/sg.md')"}
    }
  ]
}`
	gravarArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "grafo.dag.json"), dagRuim)
}

func quebrarTipoDeNo(t *testing.T, raiz string) {
	dagTipo := `{
  "name": "legoloop-alpha",
  "trigger": {"type": "manual"},
  "nodes": [
    {
      "id": "n1",
      "type": "tipo_inventado",
      "routing": {"mode": "auto", "combo": "padrao"},
      "output": {"artifact": "out/n1.md", "gate": "file_exists('out/n1.md')"}
    }
  ]
}`
	gravarArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "grafo.dag.json"), dagTipo)
}

func quebrarErroIgnorado(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "internal", "espinha", "ignora.go"), `package espinha
import "os"
func f() {
	_ = os.Remove("x")
}
`)
}

func quebrarGlobalMutavel(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "internal", "espinha", "global.go"), `package espinha
var cacheMutavel = map[string]int{}
`)
}

func quebrarComandoExclusivoSO(t *testing.T, raiz string) {
	dagSO := `{
  "name": "legoloop-alpha",
  "trigger": {"type": "manual"},
  "nodes": [
    {
      "id": "win1",
      "type": "script",
      "command": ["powershell", "-Command", "Write-Output ok"],
      "isolation": {"executor": "orchestrix_hermetic"},
      "routing": {"mode": "auto", "combo": "padrao"},
      "output": {"artifact": "out/win1.txt", "gate": "file_exists('out/win1.txt')"}
    }
  ]
}`
	gravarArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "grafo.dag.json"), dagSO)
}
