package espinha_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

type mutanteRegra struct {
	regra   string
	aplicar func(t *testing.T, raiz string)
}

var padraoRegraCodigo = regexp.MustCompile(`Regra:\s*"([a-z0-9-]+)"`)

func TestMutantesDasRegras(t *testing.T) {
	bin := compilarBinarioEspinha(t)
	exigirBaselineLimpa(t, bin, criarFixtureBase12(t))

	mutantes := mutantesDasRegras()
	exigirCoberturaCompletaRegras(t, mutantes)

	mortos := 0
	for _, m := range mutantes {
		t.Run(m.regra, func(t *testing.T) {
			raiz := criarFixtureBase12(t)
			m.aplicar(t, raiz)
			if !executarMutanteRegra(t, bin, raiz, m.regra) {
				t.Fatalf("mutante sobreviveu para regra %s", m.regra)
			}
			mortos++
		})
	}

	taxa := float64(mortos) * 100 / float64(len(mutantes))
	t.Logf("mutantes das regras da espinha: %d/%d mortos (%.1f%%)", mortos, len(mutantes), taxa)
	if mortos != len(mutantes) {
		t.Fatalf("taxa de mutantes mortos abaixo de 100%%: %d/%d (%.1f%%)", mortos, len(mutantes), taxa)
	}
}

func mutantesDasRegras() []mutanteRegra {
	return []mutanteRegra{
		{"comando-exclusivo-de-so", quebrarComandoExclusivoSO},
		{"dag-invalido", quebrarDAGInvalido},
		{"erro-ignorado-acima-do-limite", quebrarErroIgnorado},
		{"executor-incoerente", quebrarExecutorIncoerente},
		{"funcao-duplicada-entre-pacotes", quebrarFuncaoDuplicada},
		{"global-mutavel-nao-justificada", quebrarGlobalMutavel},
		{"init-nao-justificado", quebrarInitNaoJustificado},
		{"lei-sem-verificador", quebrarLeiSemVerificador},
		{"motor-inexistente", quebrarMotorInexistente},
		{"motor-sem-deterministico", quebrarMotorSemDeterministico},
		{"processo-duplicado", quebrarProcessoDuplicado},
		{"processo-incompleto", quebrarProcessoIncompleto},
		{"processo-sem-prefixo", quebrarProcessoSemPrefixo},
		{"registro-invalido", quebrarRegistroInvalido},
		{"tipo-de-no-fora-da-tabela", quebrarTipoDeNo},
		{"verbo-orfao", quebrarVerboOrfao},
		{"verbo-sem-registro", quebrarVerboSemRegistro},
	}
}

func descobrirRegrasDoGate(t *testing.T) []string {
	t.Helper()
	entradas, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("ler diretorio espinha: %v", err)
	}
	set := make(map[string]bool)
	for _, e := range entradas {
		nome := e.Name()
		if e.IsDir() || !strings.HasSuffix(nome, ".go") || strings.HasSuffix(nome, "_test.go") {
			continue
		}
		dados, err := os.ReadFile(nome)
		if err != nil {
			t.Fatalf("ler %s: %v", nome, err)
		}
		for _, m := range padraoRegraCodigo.FindAllSubmatch(dados, -1) {
			set[string(m[1])] = true
		}
	}
	regras := make([]string, 0, len(set))
	for r := range set {
		regras = append(regras, r)
	}
	slices.Sort(regras)
	return regras
}

func exigirCoberturaCompletaRegras(t *testing.T, mutantes []mutanteRegra) {
	t.Helper()
	regrasCodigo := descobrirRegrasDoGate(t)
	cobertas := make(map[string]bool, len(mutantes))
	for _, m := range mutantes {
		if cobertas[m.regra] {
			t.Fatalf("regra duplicada na lista de mutantes: %s", m.regra)
		}
		cobertas[m.regra] = true
	}
	for _, r := range regrasCodigo {
		if !cobertas[r] {
			t.Errorf("regra do gate sem mutante em TestMutantesDasRegras: %s", r)
		}
	}
	if len(regrasCodigo) == 0 {
		t.Fatal("nenhuma regra encontrada no codigo do pacote espinha")
	}
}

func exigirBaselineLimpa(t *testing.T, bin, raiz string) {
	t.Helper()
	cmd := exec.Command(bin, "espinha", "verificar", "--raiz", raiz)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("baseline limpa deveria sair com exit 0, obtido err=%v\n%s", err, out)
	}
}

func executarMutanteRegra(t *testing.T, bin, raiz, regraEsperada string) bool {
	t.Helper()
	cmd := exec.Command(bin, "espinha", "verificar", "--raiz", raiz)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 2 {
		t.Errorf("regra %s: esperado exit 2, obtido err=%v\nstdout:\n%s\nstderr:\n%s",
			regraEsperada, err, stdout.String(), stderr.String())
		return false
	}
	return saidaContemRegra(stdout.String(), regraEsperada)
}

func saidaContemRegra(saida, regraEsperada string) bool {
	for _, linha := range strings.Split(strings.TrimSpace(saida), "\n") {
		m := padraoLinhaViolacao.FindStringSubmatch(strings.TrimSpace(linha))
		if m != nil && m[1] == regraEsperada {
			return true
		}
	}
	return false
}

func quebrarRegistroInvalido(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "espinha", "registro.json"), `{"schema": "invalido"}`)
}

func quebrarDAGInvalido(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "grafo.dag.json"), `{json quebrado`)
}

func quebrarMotorSemDeterministico(t *testing.T, raiz string) {
	conteudo := `{
  "schema": "espinha.v1",
  "verbos": [
    {"nome": "espinha", "camada": "espinha", "motor": "internal/espinha", "contrato": "verifica espinha"}
  ],
  "motores": [
    {"pacote": "internal/espinha", "interface": "Verificar"}
  ],
  "processos": [{"slug": "legoloop-alpha"}],
  "requisitos": [
    {"lei": "Lei 1", "verificador": "internal/espinha", "mecanico": true}
  ]
}`
	gravarArq(t, filepath.Join(raiz, "espinha", "registro.json"), conteudo)
}

func quebrarMotorInexistente(t *testing.T, raiz string) {
	conteudo := `{
  "schema": "espinha.v1",
  "verbos": [
    {"nome": "espinha", "camada": "espinha", "motor": "internal/espinha", "contrato": "verifica espinha"}
  ],
  "motores": [
    {"pacote": "internal/espinha", "interface": "Verificar", "deterministico": true},
    {"pacote": "internal/fantasma", "interface": "X", "deterministico": true}
  ],
  "processos": [{"slug": "legoloop-alpha"}],
  "requisitos": [
    {"lei": "Lei 1", "verificador": "internal/espinha", "mecanico": true}
  ]
}`
	gravarArq(t, filepath.Join(raiz, "espinha", "registro.json"), conteudo)
}

func quebrarInitNaoJustificado(t *testing.T, raiz string) {
	gravarArq(t, filepath.Join(raiz, "internal", "espinha", "init_mut.go"), "package espinha\nfunc init() {}\n")
}
