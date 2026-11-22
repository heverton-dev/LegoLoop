package forja

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLancadores(t *testing.T) {
	homeDev, _ := os.UserHomeDir()

	t.Run("windows", func(t *testing.T) {
		dir := t.TempDir()
		arqs, err := GerarLancadores(dir, "windows")
		if err != nil {
			t.Fatalf("GerarLancadores windows falhou: %v", err)
		}
		if len(arqs) != 4 {
			t.Fatalf("esperava 4 lancadores, obteve %d: %v", len(arqs), arqs)
		}
		for _, arq := range arqs {
			if !strings.HasSuffix(arq, ".cmd") {
				t.Errorf("windows deve gerar .cmd, obteve %s", arq)
			}
			verificarConteudoLancador(t, arq, homeDev, "@echo off")
		}
	})

	for _, so := range []string{"linux", "darwin"} {
		t.Run(so, func(t *testing.T) {
			dir := t.TempDir()
			arqs, err := GerarLancadores(dir, so)
			if err != nil {
				t.Fatalf("GerarLancadores %s falhou: %v", so, err)
			}
			if len(arqs) != 4 {
				t.Fatalf("esperava 4 lancadores, obteve %d: %v", len(arqs), arqs)
			}
			for _, arq := range arqs {
				if strings.HasSuffix(arq, ".cmd") {
					t.Errorf("%s nao deve gerar .cmd: %s", so, arq)
				}
				verificarConteudoLancador(t, arq, homeDev, "#!/usr/bin/env sh")
				if runtime.GOOS != "windows" {
					info, err := os.Stat(arq)
					if err != nil {
						t.Fatalf("stat falhou em %s: %v", arq, err)
					}
					if info.Mode().Perm()&0o111 == 0 {
						t.Errorf("arquivo %s sem permissao executavel: %o", arq, info.Mode().Perm())
					}
				}
			}
		})
	}
}

func TestLancadoresErros(t *testing.T) {
	if _, err := GerarLancadores("", "linux"); err == nil {
		t.Error("esperava erro para dirBin vazio")
	}
	if _, err := GerarLancadores(t.TempDir(), "plan9"); err == nil {
		t.Error("esperava erro para SO desconhecido")
	}
}

func verificarConteudoLancador(t *testing.T, arq, homeDev, cabecalho string) {
	t.Helper()
	dados, err := os.ReadFile(arq)
	if err != nil {
		t.Fatalf("falha ao ler %s: %v", arq, err)
	}
	texto := string(dados)
	if !strings.HasPrefix(texto, cabecalho) {
		t.Errorf("cabecalho inesperado em %s: %q", arq, texto)
	}
	if strings.Contains(strings.ToLower(texto), "set path=") {
		t.Errorf("lancador %s nao pode conter set PATH=", arq)
	}
	if homeDev != "" && (strings.Contains(texto, homeDev) || strings.Contains(texto, filepath.ToSlash(homeDev))) {
		t.Errorf("lancador %s contem HOME absoluto do desenvolvedor (%s)", arq, homeDev)
	}
}
