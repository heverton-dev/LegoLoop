package migrar

import (
	"errors"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"testing"
)

func escrever(t *testing.T, raiz, rel, conteudo string) {
	t.Helper()
	p := filepath.Join(raiz, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(conteudo), 0o600); err != nil {
		t.Fatal(err)
	}
}

func ler(t *testing.T, raiz, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(raiz, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("ler %s: %v", rel, err)
	}
	return string(data)
}

func ausente(t *testing.T, raiz, rel string) {
	t.Helper()
	if _, err := os.Lstat(filepath.Join(raiz, filepath.FromSlash(rel))); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("%s deveria não existir (err=%v)", rel, err)
	}
}

// fotografar registra todo o conteúdo sob a raiz para comparar antes e depois.
func fotografar(t *testing.T, raiz string) map[string]string {
	t.Helper()
	foto := map[string]string{}
	err := filepath.WalkDir(raiz, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(raiz, p)
		if d.IsDir() {
			foto[filepath.ToSlash(rel)] = "<pasta>"
			return nil
		}
		data, err := os.ReadFile(p)
		foto[filepath.ToSlash(rel)] = string(data)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return foto
}

func processoV1(t *testing.T) string {
	t.Helper()
	raiz := t.TempDir()
	escrever(t, raiz, "pesquisa.loop.json", `{"id":"pesquisa"}`)
	escrever(t, raiz, "pesquisa.dag.json", `{"name":"pesquisa-dag"}`)
	escrever(t, raiz, "pesquisa.processo.json", `{"nome":"pesquisa"}`)
	escrever(t, raiz, "extrator/pipeline.dag.json", `{"name":"extrator"}`)
	return raiz
}

func TestMigrarIdaEVolta(t *testing.T) {
	raiz := processoV1(t)
	rel, err := Migrar(raiz)
	if err != nil {
		t.Fatalf("primeira execução: %v", err)
	}
	if len(rel.Movimentos) != 2 {
		t.Fatalf("esperava 2 movimentos, veio %+v", rel.Movimentos)
	}
	conferirMigrado(t, raiz)

	antes := fotografar(t, raiz)
	rel2, err := Migrar(raiz)
	if err != nil {
		t.Fatalf("segunda execução: %v", err)
	}
	if len(rel2.Movimentos) != 0 {
		t.Fatalf("segunda execução moveu %+v", rel2.Movimentos)
	}
	if depois := fotografar(t, raiz); !maps.Equal(antes, depois) {
		t.Fatalf("segunda execução alterou o disco:\nantes=%v\ndepois=%v", antes, depois)
	}
}

// conferirMigrado verifica os nomes novos, os arquivos intactos e o backup.
func conferirMigrado(t *testing.T, raiz string) {
	t.Helper()
	if got := ler(t, raiz, "pesquisa.studio.json"); got != `{"id":"pesquisa"}` {
		t.Fatalf("studio.json com conteúdo %q", got)
	}
	ausente(t, raiz, "pesquisa.loop.json")
	ausente(t, raiz, "extrator/pipeline.dag.json")
	if got := ler(t, raiz, "extrator/processo.dag.json"); got != `{"name":"extrator"}` {
		t.Fatalf("processo.dag.json com conteúdo %q", got)
	}
	if ler(t, raiz, "pesquisa.dag.json") != `{"name":"pesquisa-dag"}` || ler(t, raiz, "pesquisa.processo.json") != `{"nome":"pesquisa"}` {
		t.Fatal(".dag.json ou .processo.json foram alterados")
	}
	if got := ler(t, raiz, DirBackup+"/pesquisa.loop.json"); got != `{"id":"pesquisa"}` {
		t.Fatalf("backup do .loop.json com conteúdo %q", got)
	}
	if got := ler(t, raiz, DirBackup+"/extrator/pipeline.dag.json"); got != `{"name":"extrator"}` {
		t.Fatalf("backup do pipeline.dag.json com conteúdo %q", got)
	}
}

func TestMigrarDestinoExistenteNaoAlteraDisco(t *testing.T) {
	raiz := processoV1(t)
	escrever(t, raiz, "extrator/processo.dag.json", `{"name":"ja-existe"}`)
	antes := fotografar(t, raiz)
	if _, err := Migrar(raiz); !errors.Is(err, ErrDestinoExiste) {
		t.Fatalf("esperava ErrDestinoExiste, veio %v", err)
	}
	if depois := fotografar(t, raiz); !maps.Equal(antes, depois) {
		t.Fatalf("conflito alterou o disco:\nantes=%v\ndepois=%v", antes, depois)
	}
}

func TestMigrarBackupExistenteNaoAlteraDisco(t *testing.T) {
	raiz := processoV1(t)
	escrever(t, raiz, DirBackup+"/pesquisa.loop.json", `velho`)
	antes := fotografar(t, raiz)
	if _, err := Migrar(raiz); !errors.Is(err, ErrDestinoExiste) {
		t.Fatalf("esperava ErrDestinoExiste, veio %v", err)
	}
	if depois := fotografar(t, raiz); !maps.Equal(antes, depois) {
		t.Fatal("conflito no backup alterou o disco")
	}
}

func TestMigrarRecusaRaizInvalida(t *testing.T) {
	raiz := t.TempDir()
	for _, r := range []string{"", raiz + "/../" + filepath.Base(raiz), filepath.Join(raiz, "nao-existe")} {
		if _, err := Migrar(r); !errors.Is(err, ErrRaizInvalida) {
			t.Fatalf("raiz %q: esperava ErrRaizInvalida, veio %v", r, err)
		}
	}
}

func TestMigrarRecusaSymlink(t *testing.T) {
	raiz := t.TempDir()
	fora := t.TempDir()
	escrever(t, fora, "alheio.loop.json", `{}`)
	if err := os.Symlink(filepath.Join(fora, "alheio.loop.json"), filepath.Join(raiz, "alheio.loop.json")); err != nil {
		t.Skipf("symlink indisponível: %v", err)
	}
	if _, err := Migrar(raiz); !errors.Is(err, ErrSymlink) {
		t.Fatalf("esperava ErrSymlink, veio %v", err)
	}
	if ler(t, fora, "alheio.loop.json") != `{}` {
		t.Fatal("arquivo externo alterado")
	}
}
