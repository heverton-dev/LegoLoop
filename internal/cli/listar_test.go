package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/processo"
)

// cofreListar monta o cofre do AC-E03.1: legoloop-a vinculado a claude e codex, legoloop-b que compõe a, e
// legoloop-velho da v1 (só o .loop.json). Devolve o HOME (o cofre fica em HOME/vault).
func cofreListar(t *testing.T) (home, vault string) {
	t.Helper()
	home = t.TempDir()
	vault = filepath.Join(home, "vault")
	e01 := processo.Etapa{ID: "E01", Titulo: "Construir", Tipo: "construir", Peca: "tdd-workflow",
		Produz: processo.Artefato{Caminho: "out/E01.md", Tipo: "markdown"}, Gate: "file_exists('out/E01.md')"}
	sub := processo.Etapa{ID: "E02", Titulo: "Sub", Tipo: "construir", Compoe: "legoloop-a@1.2.0", Depende: []string{"E01"},
		Produz: processo.Artefato{Caminho: "out/E02.md", Tipo: "markdown"}, Gate: "file_exists('out/E02.md')"}
	for _, p := range []processo.Processo{
		{Nome: "legoloop-a", Versao: "1.2.0", Natureza: "dados", Objetivo: "A.", Etapas: []processo.Etapa{e01}},
		{Nome: "legoloop-b", Versao: "1.0.0", Natureza: "software", Objetivo: "B.", Etapas: []processo.Etapa{e01, sub}},
	} {
		p.CriadoEm = time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)
		data, err := p.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		gravarArq(t, processo.CaminhoNoCofre(vault, p.Nome), string(data))
	}
	gravarArq(t, filepath.Join(vault, "loops", "legoloop-velho.loop.json"), `{"id":"legoloop-velho","version":"0.3.0"}`)
	for _, h := range []string{"claude", "codex"} {
		gravarArq(t, filepath.Join(home, "."+h, "skills", "legoloop-a", "SKILL.md"), "# a\n")
	}
	return home, vault
}

func gravarArq(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

type itemListar struct {
	Nome      string   `json:"nome"`
	Versao    string   `json:"versao"`
	Natureza  string   `json:"natureza"`
	Harnesses []string `json:"harnesses"`
	UsadoPor  []string `json:"usado_por"`
	Origem    string   `json:"origem"`
}

// conferirTextoListar confere a tabela: cabeçalho, uma linha por processo e as colunas alinhadas.
func conferirTextoListar(t *testing.T, out string) {
	t.Helper()
	for _, quer := range []string{"NOME", "VERSÃO", "NATUREZA", "HARNESSES", "USADO POR",
		"legoloop-a", "1.2.0", "dados", "claude, codex", "legoloop-b", "legoloop-velho", "0.3.0", "desconhecida"} {
		if !strings.Contains(out, quer) {
			t.Errorf("falta %q na saída:\n%s", quer, out)
		}
	}
	var linhas []string
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(l, "NOME") || strings.HasPrefix(l, "legoloop-") {
			linhas = append(linhas, l)
		}
	}
	if len(linhas) != 4 {
		t.Fatalf("esperava cabeçalho e 3 linhas:\n%s", out)
	}
	for _, col := range []string{"VERSÃO", "NATUREZA", "HARNESSES", "USADO POR"} {
		pos := len([]rune(linhas[0][:strings.Index(linhas[0], col)]))
		for _, l := range linhas[1:] {
			if r := []rune(l); pos >= len(r) || r[pos-1] != ' ' || r[pos] == ' ' {
				t.Errorf("coluna %s desalinhada (posição %d) na linha %q", col, pos, l)
			}
		}
	}
	if !strings.Contains(linhaDe(linhas, "legoloop-a"), "legoloop-b") {
		t.Errorf("legoloop-a é usado por legoloop-b:\n%s", out)
	}
}

func linhaDe(ls []string, nome string) string {
	for _, l := range ls {
		if strings.HasPrefix(l, nome+" ") {
			return l
		}
	}
	return ""
}

func conferirJSONListar(t *testing.T, out string) {
	t.Helper()
	var doc struct {
		Vault     string       `json:"vault"`
		Processos []itemListar `json:"processos"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("stdout deve ser JSON puro: %v\n%s", err, out)
	}
	quer := []itemListar{
		{"legoloop-a", "1.2.0", "dados", []string{"claude", "codex"}, []string{"legoloop-b"}, "v2"},
		{"legoloop-b", "1.0.0", "software", []string{}, []string{}, "v2"},
		{"legoloop-velho", "0.3.0", "desconhecida", []string{}, []string{}, "v1"},
	}
	if len(doc.Processos) != len(quer) {
		t.Fatalf("processos: %+v", doc.Processos)
	}
	for i, q := range quer {
		g := doc.Processos[i]
		if g.Nome != q.Nome || g.Versao != q.Versao || g.Natureza != q.Natureza || g.Origem != q.Origem ||
			!slices.Equal(g.Harnesses, q.Harnesses) || !slices.Equal(g.UsadoPor, q.UsadoPor) {
			t.Errorf("processo %d:\n got %+v\nquer %+v", i, g, q)
		}
	}
	if !strings.Contains(out, `"harnesses": []`) || !strings.Contains(out, `"usado_por": []`) {
		t.Errorf("listas vazias saem como [] e não null:\n%s", out)
	}
}

func rodarListar(home string, args ...string) (int, string, string) {
	var so, se bytes.Buffer
	code := Run(context.Background(), append([]string{"listar"}, args...), Deps{Stdout: &so, Stderr: &se, Home: home})
	return code, so.String(), se.String()
}

func TestListarCLI(t *testing.T) {
	home, vault := cofreListar(t)
	t.Run("texto", func(t *testing.T) {
		code, out, errb := rodarListar(home, "--vault", vault)
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, errb)
		}
		conferirTextoListar(t, out)
	})
	t.Run("json", func(t *testing.T) {
		code, out, errb := rodarListar(home, "--vault", vault, "--json")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, errb)
		}
		conferirJSONListar(t, out)
	})
	t.Run("cofre padrão no HOME", func(t *testing.T) {
		h := t.TempDir()
		gravarArq(t, filepath.Join(h, ".aidd", "canonical", "loops", "legoloop-velho.loop.json"), `{"version":"0.1.0"}`)
		code, out, _ := rodarListar(h)
		if code != 0 || !strings.Contains(out, "legoloop-velho") {
			t.Fatalf("exit %d; o cofre padrão é HOME/.aidd/canonical:\n%s", code, out)
		}
	})
	t.Run("vazio", func(t *testing.T) {
		v := t.TempDir()
		code, out, _ := rodarListar(t.TempDir(), "--vault", v)
		if code != 0 || !strings.Contains(out, "O cofre "+v+" não tem processos") || strings.Contains(out, "NOME") {
			t.Fatalf("exit %d; cofre vazio imprime a mensagem e sai com 0:\n%s", code, out)
		}
		code, out, errb := rodarListar(t.TempDir(), "--vault", v, "--json")
		if code != 0 || !strings.Contains(out, `"processos": []`) || !strings.Contains(errb, "não tem processos") {
			t.Fatalf("exit %d; --json vazio dá processos [] e a mensagem no stderr:\n%s\n%s", code, out, errb)
		}
	})
	t.Run("cofre inexistente", func(t *testing.T) {
		code, _, errb := rodarListar(t.TempDir(), "--vault", filepath.Join(t.TempDir(), "x"))
		if code != 1 || !strings.Contains(errb, "Informe o cofre com --vault") {
			t.Fatalf("exit %d; esperava 1 com orientação:\n%s", code, errb)
		}
	})
	t.Run("argumento a mais", func(t *testing.T) {
		if code, _, _ := rodarListar(home, "--vault", vault, "sobra"); code != 2 {
			t.Fatalf("exit %d; argumento posicional é uso incorreto (2)", code)
		}
	})
}

// TestListarBinario prova o AC-E03.1 pelo binário: texto, --json e cofre vazio, com HOME descartável.
func TestListarBinario(t *testing.T) {
	bin := buildBinary(t)
	home, vault := cofreListar(t)
	srv := binRouter(t, 90)
	r := runBinaryEm(t, bin, srv, home, filepath.Join(home, "proj"), "listar", "--vault", vault)
	if r.code != 0 {
		t.Fatalf("exit %d\n%s", r.code, r.errb)
	}
	conferirTextoListar(t, r.out)
	r = runBinaryEm(t, bin, srv, home, filepath.Join(home, "proj"), "listar", "--vault", vault, "--json")
	if r.code != 0 {
		t.Fatalf("exit %d\n%s", r.code, r.errb)
	}
	conferirJSONListar(t, r.out)
	vazio := t.TempDir()
	r = runBinaryEm(t, bin, srv, vazio, filepath.Join(vazio, "proj"), "listar")
	if r.code != 0 || !strings.Contains(r.out, "não tem processos") {
		t.Fatalf("exit %d; cofre vazio (HOME/.aidd/canonical) imprime a mensagem e sai com 0:\n%s\n%s", r.code, r.out, r.errb)
	}
}
