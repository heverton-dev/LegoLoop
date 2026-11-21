package espinha

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestViolacaoString(t *testing.T) {
	v := Violacao{
		Regra:  "verbo-sem-registro",
		Alvo:   "demo",
		Motivo: "ausente no registro",
	}
	esperado := "regra: verbo-sem-registro; alvo: demo; motivo: ausente no registro"
	if got := v.String(); got != esperado {
		t.Fatalf("String() = %q, esperado %q", got, esperado)
	}
}

func TestVerificarRepositorioReal(t *testing.T) {
	raiz := filepath.Join("..", "..")
	v1 := Verificar(raiz)
	for _, v := range v1 {
		t.Errorf("violacao inesperada: %s", v.String())
	}
	v2 := Verificar(raiz)
	if !slices.Equal(v1, v2) {
		t.Fatalf("Verificar nao deterministico: %v != %v", v1, v2)
	}
}

func TestVerificarFixtureValida(t *testing.T) {
	raiz := criarFixtureEspinha(t)
	if vs := Verificar(raiz); len(vs) != 0 {
		t.Fatalf("esperado 0 violacoes na fixture valida, obtido: %v", vs)
	}
}

func TestVerificarRegras(t *testing.T) {
	t.Run("verbo-sem-registro", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverArq(t, filepath.Join(raiz, "internal", "cli", "cli.go"), `package cli
func mapaComandos() map[string]func() {
	return map[string]func(){"espinha": nil, "extra": nil}
}
`)
		exigirRegra(t, Verificar(raiz), "verbo-sem-registro", "extra")
	})

	t.Run("verbo-orfao", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverArq(t, filepath.Join(raiz, "internal", "cli", "cli.go"), `package cli
func mapaComandos() map[string]func() {
	return map[string]func(){}
}
`)
		exigirRegra(t, Verificar(raiz), "verbo-orfao", "espinha")
	})

	t.Run("lei-sem-verificador", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverArq(t, filepath.Join(raiz, "AGENTS.md"), "## 1. Leis\n\n1. Lei um.\n2. Lei dois sem registro.\n")
		exigirRegra(t, Verificar(raiz), "lei-sem-verificador", "Lei 2")
	})

	t.Run("processo-sem-prefixo", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		if err := os.MkdirAll(filepath.Join(raiz, "processos", "sem-prefixo"), 0o755); err != nil {
			t.Fatal(err)
		}
		exigirRegra(t, Verificar(raiz), "processo-sem-prefixo", "sem-prefixo")
	})

	t.Run("processo-duplicado por nome de DAG", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverArq(t, filepath.Join(raiz, "processos", "legoloop-beta", "grafo.dag.json"), dagValidoJSON("legoloop-alpha"))
		escreverArq(t, filepath.Join(raiz, "processos", "legoloop-beta", "ANTICORPOS.md"), "# Anticorpos\n")
		escreverRegistroComProcessos(t, raiz, []string{"legoloop-alpha", "legoloop-beta"})
		exigirRegra(t, Verificar(raiz), "processo-duplicado", "processos/legoloop-beta/grafo.dag.json")
	})

	t.Run("processo-duplicado ausente no disco", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverRegistroComProcessos(t, raiz, []string{"legoloop-alpha", "legoloop-fantasma"})
		exigirRegra(t, Verificar(raiz), "processo-duplicado", "legoloop-fantasma")
	})

	t.Run("executor-incoerente", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
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
		escreverArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "grafo.dag.json"), dagRuim)
		exigirRegra(t, Verificar(raiz), "executor-incoerente", "processos/legoloop-alpha/grafo.dag.json")
	})

	t.Run("tipo-de-no-fora-da-tabela", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
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
		escreverArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "grafo.dag.json"), dagTipo)
		exigirRegra(t, Verificar(raiz), "tipo-de-no-fora-da-tabela", "processos/legoloop-alpha/grafo.dag.json:n1")
	})

	t.Run("motor inexistente e sem deterministico", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		regSemDet := `{
  "schema": "espinha.v1",
  "verbos": [
    {"nome": "espinha", "camada": "espinha", "motor": "internal/espinha", "contrato": "verifica espinha"}
  ],
  "motores": [
    {"pacote": "internal/espinha", "interface": "Verificar"},
    {"pacote": "internal/fantasma", "interface": "X", "deterministico": true}
  ],
  "processos": [{"slug": "legoloop-alpha"}],
  "requisitos": [{"lei": "Lei 1", "verificador": "internal/espinha", "mecanico": true}]
}`
		escreverArq(t, filepath.Join(raiz, "espinha", "registro.json"), regSemDet)
		vs := Verificar(raiz)
		exigirRegra(t, vs, "motor-sem-deterministico", "internal/espinha")
		exigirRegra(t, vs, "motor-inexistente", "internal/fantasma")
	})
}

func criarFixtureEspinha(t *testing.T) string {
	t.Helper()
	raiz := t.TempDir()
	escreverArq(t, filepath.Join(raiz, "AGENTS.md"), "## 1. Leis\n\n1. Lei um.\n")
	escreverArq(t, filepath.Join(raiz, "internal", "espinha", "stub.go"), "package espinha\n")
	escreverArq(t, filepath.Join(raiz, "internal", "cli", "cli.go"), `package cli
func mapaComandos() map[string]func() {
	return map[string]func(){"espinha": nil}
}
`)
	escreverArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "grafo.dag.json"), dagValidoJSON("legoloop-alpha"))
	escreverArq(t, filepath.Join(raiz, "processos", "legoloop-alpha", "ANTICORPOS.md"), "# Anticorpos\n")
	escreverRegistroComProcessos(t, raiz, []string{"legoloop-alpha"})
	return raiz
}

func escreverRegistroComProcessos(t *testing.T, raiz string, slugs []string) {
	t.Helper()
	var itens []string
	for _, s := range slugs {
		itens = append(itens, `{"slug": "`+s+`"}`)
	}
	conteudo := `{
  "schema": "espinha.v1",
  "verbos": [
    {"nome": "espinha", "camada": "espinha", "motor": "internal/espinha", "contrato": "verifica espinha"}
  ],
  "motores": [
    {"pacote": "internal/espinha", "interface": "Verificar", "deterministico": true}
  ],
  "processos": [` + strings.Join(itens, ", ") + `],
  "requisitos": [
    {"lei": "Lei 1", "verificador": "internal/espinha", "mecanico": true}
  ]
}`
	escreverArq(t, filepath.Join(raiz, "espinha", "registro.json"), conteudo)
}

func dagValidoJSON(nome string) string {
	return `{
  "name": "` + nome + `",
  "trigger": {"type": "manual"},
  "nodes": [
    {
      "id": "etapa1",
      "type": "script",
      "command": ["echo", "ok"],
      "isolation": {"executor": "orchestrix_hermetic"},
      "routing": {"mode": "auto", "combo": "padrao"},
      "output": {"artifact": "out/etapa1.txt", "gate": "file_exists('out/etapa1.txt')"}
    }
  ]
}`
}

func escreverArq(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func exigirRegra(t *testing.T, vs []Violacao, regra, alvo string) {
	t.Helper()
	for _, v := range vs {
		if v.Regra == regra && v.Alvo == alvo {
			return
		}
	}
	t.Fatalf("esperada violacao regra=%q alvo=%q; obtidas: %v", regra, alvo, vs)
}
