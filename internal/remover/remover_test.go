package remover_test

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/processo"
	"legoloop/internal/remover"
	"legoloop/internal/ste"
)

func etapa(id string) processo.Etapa {
	return processo.Etapa{ID: id, Titulo: "Etapa " + id, Tipo: "construir", Peca: "tdd-workflow",
		Produz: processo.Artefato{Caminho: "out/" + id + ".md", Tipo: "markdown"}, Gate: "file_exists('out/" + id + ".md')"}
}

func escrever(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

// gravar põe no cofre o que o fabricar põe: processo, loop, dag e a skill.
func gravar(t *testing.T, vault string, p processo.Processo) {
	t.Helper()
	p.CriadoEm = time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)
	data, err := p.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	escrever(t, processo.CaminhoNoCofre(vault, p.Nome), string(data))
	escrever(t, filepath.Join(vault, "loops", p.Nome+".loop.json"), `{"id":"`+p.Nome+`"}`)
	escrever(t, filepath.Join(vault, "loops", p.Nome+".dag.json"), `{}`)
	escrever(t, filepath.Join(vault, "skills", p.Nome, "SKILL.md"), "# "+p.Nome+"\n")
}

// cofre monta o AC-E03.3: legoloop-a vinculado a claude e codex, legoloop-b que compõe a. O HOME é descartável.
func cofre(t *testing.T) (remover.Deps, string) {
	t.Helper()
	root := t.TempDir()
	vault, home := filepath.Join(root, "vault"), filepath.Join(root, "home")
	gravar(t, vault, processo.Processo{Nome: "legoloop-a", Versao: "1.2.0", Natureza: "dados", Objetivo: "A.", Etapas: []processo.Etapa{etapa("E01")}})
	sub := processo.Etapa{ID: "E02", Titulo: "Sub", Tipo: "construir", Compoe: "legoloop-a@1.2.0", Depende: []string{"E01"},
		Produz: processo.Artefato{Caminho: "out/E02.md", Tipo: "markdown"}, Gate: "file_exists('out/E02.md')"}
	gravar(t, vault, processo.Processo{Nome: "legoloop-b", Versao: "1.0.0", Natureza: "software", Objetivo: "B.", Etapas: []processo.Etapa{etapa("E01"), sub}})
	for _, h := range []string{"claude", "codex"} {
		escrever(t, filepath.Join(home, "."+h, "skills", "legoloop-a", "SKILL.md"), "# legoloop-a\n")
		escrever(t, filepath.Join(home, "."+h, "skills", "outra-skill", "SKILL.md"), "# outra\n")
	}
	escrever(t, filepath.Join(vault, "loops", "esteira-alheia.loop.json"), `{}`)
	return remover.Deps{Store: processo.Store{Roots: processo.Roots{Vault: vault}}, Home: home}, root
}

// foto guarda caminho e hash de tudo sob root. Duas fotos iguais provam que nada mudou.
func foto(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() {
			out[p] = "dir"
			return nil
		}
		b, err := os.ReadFile(p)
		out[p] = fmt.Sprintf("%x", sha256.Sum256(b))
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func iguais(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func existe(p string) bool { _, err := os.Lstat(p); return err == nil }

func arquivosDe(d remover.Deps, nome string) []string {
	v := d.Store.Roots.Vault
	return []string{processo.CaminhoNoCofre(v, nome), filepath.Join(v, "loops", nome+".loop.json"),
		filepath.Join(v, "loops", nome+".dag.json"), filepath.Join(v, "skills", nome)}
}

func linksDe(d remover.Deps, nome string) []string {
	return []string{filepath.Join(d.Home, ".claude", "skills", nome), filepath.Join(d.Home, ".codex", "skills", nome)}
}

// TestRemover prova o AC-E03.3: remover a falha com ErrEmUso e nada muda; remover b funciona só com
// confirmação; depois remover a apaga o processo e seus 2 links; sem confirmação nada muda.
func TestRemover(t *testing.T) {
	ctx := context.Background()
	d, root := cofre(t)
	antes := foto(t, root)

	t.Run("a está em uso: ErrEmUso e nada muda", func(t *testing.T) {
		_, err := remover.Planejar(ctx, d, "legoloop-a")
		if !errors.Is(err, remover.ErrEmUso) {
			t.Fatalf("esperava ErrEmUso, veio %v", err)
		}
		if !strings.Contains(err.Error(), "legoloop-b") {
			t.Errorf("o erro deve dizer quem usa: %v", err)
		}
		if !iguais(antes, foto(t, root)) {
			t.Fatal("o cofre ou os links mudaram")
		}
	})

	t.Run("planejar b não muda nada e mostra o cartão", func(t *testing.T) {
		pl, err := remover.Planejar(ctx, d, "legoloop-b")
		if err != nil {
			t.Fatal(err)
		}
		if !pl.Cartao.Irreversivel || pl.Cartao.Aviso == "" {
			t.Errorf("o cartão deve avisar que a ação não tem volta: %+v", pl.Cartao)
		}
		if vs := ste.LintCartao(pl.Cartao); len(vs) != 0 {
			t.Errorf("cartão fora do STE-PT: %v", vs)
		}
		if len(pl.Links) != 0 || len(pl.Arquivos) != 4 {
			t.Errorf("b não tem link e tem 4 itens no cofre: links=%v arquivos=%v", pl.Links, pl.Arquivos)
		}
		if ult := pl.Arquivos[len(pl.Arquivos)-1]; ult != processo.CaminhoNoCofre(d.Store.Roots.Vault, "legoloop-b") {
			t.Errorf("o processo.json deve sair por último, e o último é %s", ult)
		}
		if !iguais(antes, foto(t, root)) {
			t.Fatal("planejar sem confirmar mudou o disco")
		}
	})

	t.Run("b com confirmação funciona e a fica", func(t *testing.T) {
		pl, err := remover.Planejar(ctx, d, "legoloop-b")
		if err != nil {
			t.Fatal(err)
		}
		if err := remover.Executar(ctx, pl); err != nil {
			t.Fatal(err)
		}
		for _, p := range arquivosDe(d, "legoloop-b") {
			if existe(p) {
				t.Errorf("sobrou %s", p)
			}
		}
		for _, p := range append(arquivosDe(d, "legoloop-a"), linksDe(d, "legoloop-a")...) {
			if !existe(p) {
				t.Errorf("a perdeu %s", p)
			}
		}
		procs, err := processo.ProcessosDoCofre(d.Store.Roots.Vault)
		usadoPor := processo.MapearUsadoPor(procs)
		if err != nil || len(procs) != 1 || len(usadoPor["legoloop-a"]) != 0 {
			t.Errorf("depois de remover b, a não é usado por ninguém: %+v %v", procs, err)
		}
	})

	t.Run("a sem confirmação não muda nada", func(t *testing.T) {
		meio := foto(t, root)
		pl, err := remover.Planejar(ctx, d, "legoloop-a")
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(pl.Harnesses, []string{"claude", "codex"}) || len(pl.Links) != 2 {
			t.Errorf("a tem 2 links: %+v", pl)
		}
		if !strings.Contains(strings.Join(pl.Cartao.Passos, "\n"), "2 vínculos nos harnesses: claude, codex") {
			t.Errorf("o cartão de a deve dizer os 2 vínculos: %v", pl.Cartao.Passos)
		}
		if !iguais(meio, foto(t, root)) {
			t.Fatal("planejar sem confirmar mudou o disco")
		}
	})

	t.Run("a confirmado apaga o processo e os 2 links", func(t *testing.T) {
		pl, err := remover.Planejar(ctx, d, "legoloop-a")
		if err != nil {
			t.Fatal(err)
		}
		if err := remover.Executar(ctx, pl); err != nil {
			t.Fatal(err)
		}
		for _, p := range append(arquivosDe(d, "legoloop-a"), linksDe(d, "legoloop-a")...) {
			if existe(p) {
				t.Errorf("sobrou %s", p)
			}
		}
		for _, p := range []string{filepath.Join(d.Store.Roots.Vault, "loops", "esteira-alheia.loop.json"),
			filepath.Join(d.Home, ".claude", "skills", "outra-skill", "SKILL.md"), filepath.Join(d.Home, ".codex", "skills", "outra-skill", "SKILL.md")} {
			if !existe(p) {
				t.Errorf("remover apagou o que não era do processo: %s", p)
			}
		}
	})
}

// TestRemoverErros: nome fora do padrão, processo que não existe e contexto cancelado não tocam o disco.
func TestRemoverErros(t *testing.T) {
	d, root := cofre(t)
	antes := foto(t, root)
	ctx := context.Background()
	for _, nome := range []string{"", "../x", "legoloop-a/..", "legoloop-a/b", "outra-skill", "legoloop-", `legoloop-a\..\x`} {
		if _, err := remover.Planejar(ctx, d, nome); !errors.Is(err, remover.ErrNomeInvalido) {
			t.Errorf("nome %q: esperava ErrNomeInvalido, veio %v", nome, err)
		}
	}
	if _, err := remover.Planejar(ctx, d, "legoloop-zzz"); !errors.Is(err, remover.ErrNaoEncontrado) {
		t.Errorf("esperava ErrNaoEncontrado, veio %v", err)
	}
	arqVault := filepath.Join(t.TempDir(), "arquivo-vault")
	if err := os.WriteFile(arqVault, []byte("não é pasta"), 0o644); err != nil {
		t.Fatal(err)
	}
	dVault := remover.Deps{Store: processo.Store{Roots: processo.Roots{Vault: arqVault}}, Home: d.Home}
	if _, err := remover.Planejar(ctx, dVault, "legoloop-b"); !errors.Is(err, remover.ErrNaoEncontrado) {
		t.Errorf("esperava ErrNaoEncontrado para vault arquivo, veio %v", err)
	}
	conteudo, err := os.ReadFile(arqVault)
	if err != nil || string(conteudo) != "não é pasta" {
		t.Fatalf("o arquivo de vault foi alterado: %s %v", conteudo, err)
	}
	cancelado, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := remover.Planejar(cancelado, d, "legoloop-b"); !errors.Is(err, context.Canceled) {
		t.Errorf("planejar cancelado: esperava context.Canceled, veio %v", err)
	}
	pl, err := remover.Planejar(ctx, d, "legoloop-b")
	if err != nil {
		t.Fatal(err)
	}
	if err := remover.Executar(cancelado, pl); !errors.Is(err, context.Canceled) {
		t.Errorf("esperava context.Canceled, veio %v", err)
	}
	if !iguais(antes, foto(t, root)) {
		t.Fatal("erro ou cancelamento mudou o disco")
	}
}

// TestRemoverSymlink: o vínculo por symlink some sem apagar a skill do cofre antes da hora, e a pasta de
// outra skill no mesmo harness fica.
func TestRemoverSymlink(t *testing.T) {
	d, _ := cofre(t)
	link := filepath.Join(d.Home, ".claude", "skills", "legoloop-b")
	alvo := filepath.Join(d.Store.Roots.Vault, "skills", "legoloop-b")
	if err := ligar(alvo, link); err != nil {
		t.Skipf("sem symlink nem junction neste sistema: %v", err)
	}
	pl, err := remover.Planejar(context.Background(), d, "legoloop-b")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(pl.Links, []string{link}) {
		t.Fatalf("links = %v", pl.Links)
	}
	if err := remover.Executar(context.Background(), pl); err != nil {
		t.Fatal(err)
	}
	if existe(link) || existe(filepath.Join(d.Store.Roots.Vault, "skills", "legoloop-b")) {
		t.Error("link ou skill sobrou")
	}
	if !existe(filepath.Join(d.Home, ".claude", "skills", "outra-skill", "SKILL.md")) {
		t.Error("apagou outra skill")
	}
}

// ligar cria um symlink; no Windows sem privilégio, cai para junction (também um link para a pasta do cofre).
func ligar(alvo, link string) error {
	err := os.Symlink(alvo, link)
	if err != nil && runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "mklink", "/J", link, alvo).Run()
	}
	return err
}

// TestRemoverSoOQueExiste: um processo sem DAG nem skill no cofre (por exemplo, um que a v1 gravou só com o
// loop) tem no plano só o que existe, e remover funciona.
func TestRemoverSoOQueExiste(t *testing.T) {
	d, _ := cofre(t)
	v := d.Store.Roots.Vault
	for _, p := range []string{filepath.Join(v, "loops", "legoloop-b.dag.json"), filepath.Join(v, "skills", "legoloop-b")} {
		if err := os.RemoveAll(p); err != nil {
			t.Fatal(err)
		}
	}
	pl, err := remover.Planejar(context.Background(), d, "legoloop-b")
	if err != nil {
		t.Fatal(err)
	}
	if len(pl.Arquivos) != 2 || !strings.Contains(strings.Join(pl.Cartao.Passos, ""), "Apagar 2 itens") {
		t.Fatalf("esperava só o loop e o processo: %v", pl.Arquivos)
	}
	if err := remover.Executar(context.Background(), pl); err != nil {
		t.Fatal(err)
	}
	if existe(processo.CaminhoNoCofre(v, "legoloop-b")) {
		t.Error("o processo sobrou")
	}
}
