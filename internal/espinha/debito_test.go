package espinha

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func raizComGo(t *testing.T, conteudo string) string {
	t.Helper()
	raiz := t.TempDir()
	dir := filepath.Join(raiz, `internal`, `x`)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, `x.go`), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	return raiz
}

func TestDebitoMarcadoReprovaComentarioQueAdiaCorrecao(t *testing.T) {
	casos := []string{
		`// TODO: ligar o Studio`,
		`// FIXME(ana): corrigir depois`,
		`// LACUNA: falta o teste`,
		`/* HACK: gambiarra */`,
		`// DÉBITO: refatorar`,
	}
	for _, c := range casos {
		vs := verificarDebitoMarcado(raizComGo(t, "package x\n\n"+c+"\nfunc F() {}\n"))
		if len(vs) != 1 || vs[0].Regra != `debito-marcado` || !strings.Contains(vs[0].Alvo, `x.go:3`) {
			t.Errorf(`comentário %q: violações = %+v`, c, vs)
		}
	}
}

func TestDebitoMarcadoIgnoraPortuguesEStrings(t *testing.T) {
	fonte := `package x

// todo mundo precisa ler isto
// Todo: palavra comum no início da frase não é marcador
// O TODO da string abaixo não conta
const s = "TODO: texto de prompt"
`
	if vs := verificarDebitoMarcado(raizComGo(t, fonte)); len(vs) != 0 {
		t.Fatalf(`falso positivo: %+v`, vs)
	}
}

func TestDebitoMarcadoPulaTestdata(t *testing.T) {
	raiz := raizComGo(t, "package x\n")
	dir := filepath.Join(raiz, `internal`, `x`, `testdata`)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, `quebrado.go`), []byte(`nao compila // TODO: x`), 0o644); err != nil {
		t.Fatal(err)
	}
	if vs := verificarDebitoMarcado(raiz); len(vs) != 0 {
		t.Fatalf(`testdata devia ser ignorado: %+v`, vs)
	}
}
