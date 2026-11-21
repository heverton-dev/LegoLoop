package espinha

import (
	"os"
	"path/filepath"
	"testing"
)

func localizarRaizRepo(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("obter diretorio atual: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			t.Fatalf("go.mod nao encontrado acima de %s", dir)
		}
		dir = pai
	}
}

func TestSemReferenciaAoNomeAntigo(t *testing.T) {
	raiz := localizarRaizRepo(t)
	caminho := filepath.Join(raiz, "espinha", "renomeados.json")
	renomeados, err := CarregarRenomeados(caminho)
	if err != nil {
		t.Fatalf("carregar renomeados: %v", err)
	}
	if len(renomeados) != 10 {
		t.Fatalf("esperado 10 nomes antigos, obtido %d", len(renomeados))
	}

	t.Run("processo-sem-prefixo", func(t *testing.T) {
		dirProc := filepath.Join(raiz, "processos")
		violacoes, err := VerificarPrefixoProcessos(dirProc)
		if err != nil {
			t.Fatalf("verificar prefixo: %v", err)
		}
		if len(violacoes) > 0 {
			t.Fatalf("processos sem prefixo legoloop-: %v", violacoes)
		}
	})

	t.Run("referencias-no-repositorio", func(t *testing.T) {
		achados, err := BuscarReferenciasAntigas(raiz, renomeados)
		if err != nil {
			t.Fatalf("varrer repositorio: %v", err)
		}
		if len(achados) > 0 {
			for _, a := range achados {
				t.Errorf("referencia ao nome antigo: %s", a)
			}
		}
	})

	t.Run("mutacao-pasta-sem-prefixo", func(t *testing.T) {
		tmp := t.TempDir()
		if err := os.Mkdir(filepath.Join(tmp, "sem-prefixo"), 0o755); err != nil {
			t.Fatalf("criar pasta mutante: %v", err)
		}
		violacoes, err := VerificarPrefixoProcessos(tmp)
		if err != nil {
			t.Fatalf("verificar prefixo mutante: %v", err)
		}
		if len(violacoes) != 1 || violacoes[0] != "sem-prefixo" {
			t.Fatalf("esperado detectar sem-prefixo, obtido %v", violacoes)
		}
	})

	t.Run("mutacao-referencia-antiga", func(t *testing.T) {
		tmp := t.TempDir()
		alvo := filepath.Join(tmp, "README.md")
		conteudo := "Exemplo: processos/addy-osmani/dag.json\n"
		if err := os.WriteFile(alvo, []byte(conteudo), 0o644); err != nil {
			t.Fatalf("escrever mutante: %v", err)
		}
		achados, err := BuscarReferenciasAntigas(tmp, renomeados)
		if err != nil {
			t.Fatalf("varrer mutante: %v", err)
		}
		if len(achados) == 0 {
			t.Fatalf("esperado detectar referencia mutante em README.md")
		}
	})

	t.Run("verbos-cli-nao-dao-falso-positivo", func(t *testing.T) {
		tmp := t.TempDir()
		alvo := filepath.Join(tmp, "README.md")
		conteudo := "Use `legoloop fechar` e `legoloop listar` ou processos/legoloop-fechar/dag.json\n"
		if err := os.WriteFile(alvo, []byte(conteudo), 0o644); err != nil {
			t.Fatalf("escrever fixture: %v", err)
		}
		achados, err := BuscarReferenciasAntigas(tmp, renomeados)
		if err != nil {
			t.Fatalf("varrer fixture: %v", err)
		}
		if len(achados) > 0 {
			t.Fatalf("falso positivo para verbos CLI: %v", achados)
		}
	})
}
