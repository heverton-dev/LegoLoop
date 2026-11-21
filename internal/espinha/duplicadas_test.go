package espinha

import (
	"path/filepath"
	"testing"
)

func TestFuncaoDuplicadaEntrePacotes(t *testing.T) {
	t.Run("dispara para mesmo nome e mesma assinatura em pacotes distintos", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverArq(t, filepath.Join(raiz, "internal", "alfa", "util.go"), `package alfa

func NormalizarNome(entrada string) string {
	return entrada + "_norm"
}
`)
		escreverArq(t, filepath.Join(raiz, "internal", "beta", "util.go"), `package beta

func NormalizarNome(texto string) string {
	return texto + "_norm"
}
`)
		vs := Verificar(raiz)
		if !temRegra(vs, "funcao-duplicada-entre-pacotes") {
			t.Fatalf("esperada regra funcao-duplicada-entre-pacotes; obtidas: %v", vs)
		}
	})

	t.Run("nao dispara quando assinaturas diferem", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverArq(t, filepath.Join(raiz, "internal", "alfa", "util.go"), `package alfa

func NormalizarNome(entrada string) string {
	return entrada
}
`)
		escreverArq(t, filepath.Join(raiz, "internal", "beta", "util.go"), `package beta

func NormalizarNome(entrada string, limite int) string {
	return entrada
}
`)
		vs := Verificar(raiz)
		if temRegra(vs, "funcao-duplicada-entre-pacotes") {
			t.Fatalf("nao esperada regra funcao-duplicada-entre-pacotes; obtidas: %v", vs)
		}
	})

	t.Run("nao dispara quando nomes diferem", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverArq(t, filepath.Join(raiz, "internal", "alfa", "util.go"), `package alfa

func NormalizarA(entrada string) string {
	return entrada
}
`)
		escreverArq(t, filepath.Join(raiz, "internal", "beta", "util.go"), `package beta

func NormalizarB(entrada string) string {
	return entrada
}
`)
		vs := Verificar(raiz)
		if temRegra(vs, "funcao-duplicada-entre-pacotes") {
			t.Fatalf("nao esperada regra funcao-duplicada-entre-pacotes; obtidas: %v", vs)
		}
	})

	t.Run("ignora testes metodos init e main", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverArq(t, filepath.Join(raiz, "internal", "alfa", "util.go"), `package alfa

type S struct{}

func (S) Executar(x string) string { return x }
func init() {}
func main() {}
`)
		escreverArq(t, filepath.Join(raiz, "internal", "beta", "util.go"), `package beta

type S struct{}

func (S) Executar(x string) string { return x }
func init() {}
func main() {}
`)
		escreverArq(t, filepath.Join(raiz, "internal", "alfa", "util_test.go"), `package alfa

func AuxiliarTeste(x string) string { return x }
`)
		escreverArq(t, filepath.Join(raiz, "internal", "beta", "util_test.go"), `package beta

func AuxiliarTeste(x string) string { return x }
`)
		vs := Verificar(raiz)
		if temRegra(vs, "funcao-duplicada-entre-pacotes") {
			t.Fatalf("nao esperada regra funcao-duplicada-entre-pacotes; obtidas: %v", vs)
		}
	})

	t.Run("nao dispara quando tipos locais pertencem a pacotes distintos", func(t *testing.T) {
		raiz := criarFixtureEspinha(t)
		escreverArq(t, filepath.Join(raiz, "internal", "alfa", "req.go"), `package alfa

type Request struct{}

func Planejar(r Request) error { return nil }
`)
		escreverArq(t, filepath.Join(raiz, "internal", "beta", "req.go"), `package beta

type Request struct{}

func Planejar(r Request) error { return nil }
`)
		vs := Verificar(raiz)
		if temRegra(vs, "funcao-duplicada-entre-pacotes") {
			t.Fatalf("nao esperada regra para tipos locais de pacotes distintos; obtidas: %v", vs)
		}
	})
}
