package escrita

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestGravarAtomicoCriaESobrescreve(t *testing.T) {
	dir := t.TempDir()
	alvo := filepath.Join(dir, "dados.json")

	inicial := []byte(`{"versao":1}`)
	if err := GravarAtomico(alvo, inicial); err != nil {
		t.Fatalf("GravarAtomico inicial falhou: %v", err)
	}
	lido, err := os.ReadFile(alvo)
	if err != nil {
		t.Fatalf("ReadFile inicial falhou: %v", err)
	}
	if !bytes.Equal(lido, inicial) {
		t.Fatalf("conteudo inicial divergente: obtido %q, esperado %q", lido, inicial)
	}

	atualizado := []byte(`{"versao":2}`)
	if err := GravarAtomicoComPerm(alvo, atualizado, 0o600); err != nil {
		t.Fatalf("GravarAtomicoComPerm sobrescrita falhou: %v", err)
	}
	lido2, err := os.ReadFile(alvo)
	if err != nil {
		t.Fatalf("ReadFile atualizado falhou: %v", err)
	}
	if !bytes.Equal(lido2, atualizado) {
		t.Fatalf("conteudo atualizado divergente: obtido %q, esperado %q", lido2, atualizado)
	}
	garantirSemTemporarios(t, dir, "dados.json")
}

func TestGravarAtomicoFalhaAntesDoRenamePreservaOriginalELimpaTemporario(t *testing.T) {
	dir := t.TempDir()
	alvo := filepath.Join(dir, "critico.json")
	original := []byte("conteudo-original-intacto-byte-a-byte")
	if err := os.WriteFile(alvo, original, 0o644); err != nil {
		t.Fatalf("setup original falhou: %v", err)
	}

	errInjetado := errors.New("falha injetada antes do rename")
	novo := []byte("novo-conteudo-que-nao-pode-corromper-original")
	var tmpObservado string

	err := gravarAtomicoComGancho(alvo, novo, 0o644, func(tmpPath string) error {
		tmpObservado = tmpPath
		// Confirma que o temporario existe e contem os novos bytes antes do rename.
		b, readErr := os.ReadFile(tmpPath)
		if readErr != nil {
			t.Errorf("temporario ausente no gancho: %v", readErr)
		}
		if !bytes.Equal(b, novo) {
			t.Errorf("temporario com conteudo divergente: obtido %q, esperado %q", b, novo)
		}
		return errInjetado
	})
	if !errors.Is(err, errInjetado) {
		t.Fatalf("erro esperado %v, obtido %v", errInjetado, err)
	}
	if tmpObservado == "" {
		t.Fatal("gancho antesDeRenomear nao foi executado")
	}

	atual, err := os.ReadFile(alvo)
	if err != nil {
		t.Fatalf("ReadFile apos falha: %v", err)
	}
	if !bytes.Equal(atual, original) {
		t.Fatalf("arquivo original foi alterado: obtido %q, esperado %q", atual, original)
	}
	if _, statErr := os.Stat(tmpObservado); !os.IsNotExist(statErr) {
		t.Fatalf("temporario orfao permaneceu em disco: %s", tmpObservado)
	}
	garantirSemTemporarios(t, dir, "critico.json")
}

func TestGravarAtomicoDiretorioInexistenteRetornaErro(t *testing.T) {
	alvo := filepath.Join(t.TempDir(), "nao-existe", "arquivo.json")
	if err := GravarAtomico(alvo, []byte("abc")); err == nil {
		t.Fatal("esperava erro para diretorio pai inexistente")
	}
}

func garantirSemTemporarios(t *testing.T, dir, esperado string) {
	t.Helper()
	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir falhou: %v", err)
	}
	if len(entradas) != 1 || entradas[0].Name() != esperado {
		nomes := make([]string, 0, len(entradas))
		for _, e := range entradas {
			nomes = append(nomes, e.Name())
		}
		t.Fatalf("arquivos inesperados no diretorio: %v (esperado apenas %q)", nomes, esperado)
	}
}
