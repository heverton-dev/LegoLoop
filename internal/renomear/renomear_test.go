package renomear

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/jail"
)

func dataFixa() time.Time {
	return time.Date(2026, 11, 21, 15, 30, 0, 0, time.UTC)
}

type fixtureProjeto struct {
	raiz   string
	vault  string
	home   string
	docs   string
	claude string
}

func montarFixture(t *testing.T) fixtureProjeto {
	t.Helper()
	base := t.TempDir()
	f := fixtureProjeto{
		raiz:   filepath.Join(base, "proj"),
		vault:  filepath.Join(base, "vault"),
		home:   filepath.Join(base, "home"),
		docs:   filepath.Join(base, "proj", "docs"),
		claude: filepath.Join(base, "home", ".claude", "skills"),
	}
	dirs := []string{
		filepath.Join(f.raiz, caminhos.DirProcessos, "velho"),
		filepath.Join(f.raiz, caminhos.DirProcessos, "legoloop-velho"),
		filepath.Join(f.vault, caminhos.DirProcessos, "velho"),
		filepath.Join(f.vault, "skills", "velho"),
		f.claude,
		filepath.Join(f.docs, "manuals"),
		filepath.Join(f.docs, "archive"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	escreverArquivo(t, filepath.Join(f.raiz, caminhos.DirProcessos, "velho", caminhos.ArquivoDAG),
		`{"name":"velho","version":"1.0.0","nodes":[]}`+"\n")
	escreverArquivo(t, filepath.Join(f.raiz, caminhos.DirProcessos, "velho", "triagem.json"),
		`{"schema_version":"triagem.v1","processo":"velho","regras_deterministicas":[{"padrao_slug":"velho","comandos_cli":["legoloop velho","/legoloop-velho"]}],"criterio_ambiguidade":{"descricao":"Classificação para o processo velho"}}`+"\n")
	escreverArquivo(t, filepath.Join(f.raiz, caminhos.DirProcessos, "velho", "anticorpos.json"),
		`{"schema_version":"anticorpos.v1","processo":"velho","globais":[],"dominio":[]}`+"\n")
	escreverArquivo(t, filepath.Join(f.raiz, caminhos.DirProcessos, "legoloop-velho", caminhos.ArquivoDAG),
		`{"name":"legoloop-velho","version":"1.0.0","nodes":[]}`+"\n")
	escreverArquivo(t, filepath.Join(f.vault, caminhos.DirProcessos, "velho", caminhos.ArquivoDAG),
		`{"name":"velho","version":"1.0.0","nodes":[]}`+"\n")
	escreverArquivo(t, filepath.Join(f.vault, "skills", "velho", "SKILL.md"),
		"---\nname: velho\n---\n")
	// Vínculo do harness (symlink ou pasta copiada, portável no Windows).
	linkAlvo := filepath.Join(f.claude, "velho")
	if err := os.Symlink(filepath.Join(f.vault, "skills", "velho"), linkAlvo); err != nil {
		if err := os.MkdirAll(linkAlvo, 0o755); err != nil {
			t.Fatal(err)
		}
		escreverArquivo(t, filepath.Join(linkAlvo, "SKILL.md"), "---\nname: velho\n---\n")
	}
	escreverArquivo(t, filepath.Join(f.docs, "manuals", "guia.md"),
		"Veja processos/velho/"+caminhos.ArquivoDAG+" e rode `legoloop exec velho`.\nO termo comum velho na frase permanece.\n")
	escreverArquivo(t, filepath.Join(f.docs, "archive", "historico.md"),
		"Registro histórico de processos/velho que não deve mudar.\n")
	return f
}

func escreverArquivo(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func lerArquivo(t *testing.T, caminho string) string {
	t.Helper()
	b, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler %s: %v", caminho, err)
	}
	return string(b)
}

func TestRenomear(t *testing.T) {
	t.Run("dry-run nao altera disco e lista alteracoes", func(t *testing.T) {
		f := montarFixture(t)
		rel, err := Executar(Opcoes{
			Antigo: "velho",
			Novo:   "legoloop-novo",
			Raiz:   f.raiz,
			Vault:  f.vault,
			Home:   f.home,
			Docs:   f.docs,
			DryRun: true,
			Agora:  dataFixa,
		})
		if err != nil {
			t.Fatalf("dry-run falhou: %v", err)
		}
		if !rel.DryRun {
			t.Fatal("esperava DryRun=true no relatório")
		}
		if len(rel.Alteracoes) < 5 {
			t.Fatalf("esperava ao menos 5 alterações previstas, obteve %d: %+v", len(rel.Alteracoes), rel.Alteracoes)
		}
		if _, err := os.Stat(filepath.Join(f.raiz, caminhos.DirProcessos, "velho")); err != nil {
			t.Fatal("dry-run removeu pasta antiga")
		}
		if _, err := os.Stat(filepath.Join(f.raiz, caminhos.DirProcessos, "legoloop-novo")); !os.IsNotExist(err) {
			t.Fatal("dry-run criou pasta nova")
		}
		if _, err := os.Stat(filepath.Join(f.raiz, "espinha", "renomeados.json")); !os.IsNotExist(err) {
			t.Fatal("dry-run criou espinha/renomeados.json")
		}
	})

	t.Run("execucao completa atualiza tudo com backup e manifesto", func(t *testing.T) {
		f := montarFixture(t)
		rel, err := Executar(Opcoes{
			Antigo: "velho",
			Novo:   "legoloop-novo",
			Raiz:   f.raiz,
			Vault:  f.vault,
			Home:   f.home,
			Docs:   f.docs,
			Agora:  dataFixa,
		})
		if err != nil {
			t.Fatalf("execucao falhou: %v", err)
		}
		// Pasta antiga sumiu, nova existe.
		if _, err := os.Stat(filepath.Join(f.raiz, caminhos.DirProcessos, "velho")); !os.IsNotExist(err) {
			t.Fatal("pasta antiga ainda existe")
		}
		novaPasta := filepath.Join(f.raiz, caminhos.DirProcessos, "legoloop-novo")
		// DAG name atualizado.
		dagTexto := lerArquivo(t, filepath.Join(novaPasta, caminhos.ArquivoDAG))
		if !strings.Contains(dagTexto, `"legoloop-novo"`) || strings.Contains(dagTexto, `"velho"`) {
			t.Fatalf("DAG não atualizado: %s", dagTexto)
		}
		// triagem.json atualizado.
		triTexto := lerArquivo(t, filepath.Join(novaPasta, "triagem.json"))
		if !strings.Contains(triTexto, `"processo": "legoloop-novo"`) || !strings.Contains(triTexto, `"padrao_slug": "legoloop-novo"`) {
			t.Fatalf("triagem.json não atualizado: %s", triTexto)
		}
		// anticorpos.json atualizado.
		antiTexto := lerArquivo(t, filepath.Join(novaPasta, "anticorpos.json"))
		if !strings.Contains(antiTexto, `"processo": "legoloop-novo"`) {
			t.Fatalf("anticorpos.json não atualizado: %s", antiTexto)
		}
		// Cofre e vínculo do harness atualizados.
		if _, err := os.Stat(filepath.Join(f.vault, caminhos.DirProcessos, "legoloop-novo")); err != nil {
			t.Fatalf("processo no cofre não renomeado: %v", err)
		}
		if _, err := os.Lstat(filepath.Join(f.claude, "legoloop-novo")); err != nil {
			t.Fatalf("vínculo no harness não renomeado: %v", err)
		}
		if _, err := os.Lstat(filepath.Join(f.claude, "velho")); !os.IsNotExist(err) {
			t.Fatal("vínculo antigo ainda existe no harness")
		}
		// Docs atualizado com cuidado (mantém palavra comum, ignora archive).
		guia := lerArquivo(t, filepath.Join(f.docs, "manuals", "guia.md"))
		if !strings.Contains(guia, "processos/legoloop-novo/"+caminhos.ArquivoDAG) || !strings.Contains(guia, "legoloop exec legoloop-novo") {
			t.Fatalf("doc não atualizado: %s", guia)
		}
		if !strings.Contains(guia, "O termo comum velho na frase permanece.") {
			t.Fatalf("palavra comum foi corrompida: %s", guia)
		}
		hist := lerArquivo(t, filepath.Join(f.docs, "archive", "historico.md"))
		if !strings.Contains(hist, "processos/velho") {
			t.Fatalf("docs/archive não deveria ser alterado: %s", hist)
		}
		// Backup + MANIFESTO.tsv existem.
		manifesto := lerArquivo(t, rel.Manifesto)
		if !strings.Contains(manifesto, "origem\tbackup\tacao") || !strings.Contains(manifesto, "velho") {
			t.Fatalf("MANIFESTO.tsv incompleto: %s", manifesto)
		}
		// espinha/renomeados.json gravado sem duplicatas.
		regBytes := lerArquivo(t, filepath.Join(f.raiz, "espinha", "renomeados.json"))
		var regs []RegistroRenomeado
		if err := json.Unmarshal([]byte(regBytes), &regs); err != nil {
			t.Fatalf("renomeados.json inválido: %v", err)
		}
		if len(regs) != 1 || regs[0].Antigo != "velho" || regs[0].Novo != "legoloop-novo" || regs[0].Data != "2026-11-21" {
			t.Fatalf("registro inesperado: %+v", regs)
		}
		// BuscarRenomeado encontra o registro.
		achou, ok := BuscarRenomeado(f.raiz, "velho")
		if !ok || achou.Novo != "legoloop-novo" {
			t.Fatalf("BuscarRenomeado falhou: ok=%v reg=%+v", ok, achou)
		}
	})

	t.Run("colisao devolve ErrColisao sem tocar em nada", func(t *testing.T) {
		f := montarFixture(t)
		_, err := Executar(Opcoes{
			Antigo: "velho",
			Novo:   "legoloop-velho",
			Raiz:   f.raiz,
			Vault:  f.vault,
			Home:   f.home,
			Docs:   f.docs,
			Agora:  dataFixa,
		})
		if !errors.Is(err, ErrColisao) {
			t.Fatalf("esperava ErrColisao, obteve %v", err)
		}
		if _, err := os.Stat(filepath.Join(f.raiz, caminhos.DirProcessos, "velho")); err != nil {
			t.Fatal("colisão alterou pasta antiga")
		}
		if _, err := os.Stat(filepath.Join(f.raiz, "espinha", "renomeados.json")); !os.IsNotExist(err) {
			t.Fatal("colisão criou renomeados.json")
		}
	})

	t.Run("nome hostil devolve erro de jail sem tocar em nada", func(t *testing.T) {
		f := montarFixture(t)
		for _, hostil := range []string{"../x", "a/b", `..\x`} {
			_, err := Executar(Opcoes{
				Antigo: hostil,
				Novo:   "legoloop-novo",
				Raiz:   f.raiz,
				Agora:  dataFixa,
			})
			if !errors.Is(err, jail.ErrEscapeRaiz) {
				t.Fatalf("antigo %q: esperava ErrEscapeRaiz, obteve %v", hostil, err)
			}
			_, err = Executar(Opcoes{
				Antigo: "velho",
				Novo:   hostil,
				Raiz:   f.raiz,
				Agora:  dataFixa,
			})
			if !errors.Is(err, jail.ErrEscapeRaiz) {
				t.Fatalf("novo %q: esperava ErrEscapeRaiz, obteve %v", hostil, err)
			}
		}
	})

	t.Run("processo antigo inexistente devolve ErrNaoEncontrado", func(t *testing.T) {
		f := montarFixture(t)
		_, err := Executar(Opcoes{
			Antigo: "nao-existe",
			Novo:   "legoloop-novo",
			Raiz:   f.raiz,
			Agora:  dataFixa,
		})
		if !errors.Is(err, ErrNaoEncontrado) {
			t.Fatalf("esperava ErrNaoEncontrado, obteve %v", err)
		}
	})
}
