package fechar_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/fechar"
)

// semearDocs grava o manual, o livro e a Skill reais mínimos de um processo no repositório de teste.
func semearDocs(t *testing.T, repo, slug string) {
	t.Helper()
	canon := "legoloop-" + slug
	manual := filepath.Join(repo, "docs", "manuals", "2026-11-15_MANUAL_"+strings.ToUpper(strings.ReplaceAll(slug, "-", "_"))+".md")
	livro := filepath.Join(repo, "docs", "books", canon, "livro.json")
	skill := filepath.Join(repo, "skills", canon, "SKILL.md")
	corpoSkill := "---\nname: " + canon + "\ndescription: Skill de teste\n---\n# " + canon + "\n" +
		"Nós do DAG citados: 01-alfa, 02-beta, e01.\n" +
		strings.Repeat("Passo operacional com resultado verificável.\n", 12)
	arquivos := map[string]string{
		manual: "# Manual\n" + strings.Repeat("Instrução operacional.\n", 12),
		livro:  `{"titulo":"Livro de teste","slug":"` + slug + `","capitulos":[{"id":"e01"}]}`,
		skill:  corpoSkill,
	}
	for caminho, conteudo := range arquivos {
		escrever(t, caminho, conteudo)
	}
}

// semearDAG grava o DAG de dois nós do processo na pasta canônica de Processos.
func semearDAG(t *testing.T, repo, slug string) string {
	t.Helper()
	canon := "legoloop-" + slug
	dag := `{"name":"` + canon + `","version":"1.0.0","nodes":[{"id":"01-alfa"},{"id":"02-beta"}]}`
	caminho := filepath.Join(repo, dirDialetoAntigo, canon, caminhos.ArquivoDAG)
	escrever(t, caminho, dag)
	return caminho
}

func escrever(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func depsFechar(repo string) fechar.Deps {
	return fechar.Deps{RepoDir: repo, VaultDir: repo, Home: repo, SkipBuild: true,
		Indexer: func(string) error { return nil }, Portoes: func(string) error { return nil }}
}

func fecharSlug(t *testing.T, repo, slug string) *fechar.Result {
	t.Helper()
	res, err := fechar.Executar(context.Background(), fechar.Request{Slug: slug}, depsFechar(repo))
	if err != nil {
		t.Fatalf("Executar devolveu erro: %v", err)
	}
	return res
}

func exigirEtapaComErro(t *testing.T, res *fechar.Result, etapa int, trecho string) {
	t.Helper()
	if res.Sucesso {
		t.Fatalf("esperava Sucesso=false, obteve true: %+v", res.Etapas)
	}
	e := res.Etapas[etapa-1]
	if e.Status != "error" || !strings.Contains(e.Detalhe, trecho) {
		t.Errorf("etapa %d deveria falhar com %q: %+v", etapa, trecho, e)
	}
}

func TestFecharSemArtefatosFalhaENaoGravaStub(t *testing.T) {
	repo := t.TempDir()
	res := fecharSlug(t, repo, "ausente")
	for etapa := 1; etapa <= 4; etapa++ {
		exigirEtapaComErro(t, res, etapa, "artefato obrigatório ausente")
	}
	for _, rel := range []string{"skills", dirDialetoAntigo, caminhos.DirProcessos, "docs/manuals", "docs/books"} {
		if _, err := os.Stat(filepath.Join(repo, rel)); err == nil {
			t.Errorf("fechar gravou stub em %s", rel)
		}
	}
}

func TestFecharRecusaSkillStub(t *testing.T) {
	repo := t.TempDir()
	semearDocs(t, repo, "curta")
	semearDAG(t, repo, "curta")
	escrever(t, filepath.Join(repo, "skills", "legoloop-curta", "SKILL.md"),
		"---\nname: legoloop-curta\ndescription: Skill universal de curta\n---\n# legoloop-curta\n")
	exigirEtapaComErro(t, fecharSlug(t, repo, "curta"), 3, "stub")
}

func TestFecharRecusaLivroStub(t *testing.T) {
	repo := t.TempDir()
	semearDocs(t, repo, "magro")
	semearDAG(t, repo, "magro")
	escrever(t, filepath.Join(repo, "docs", "books", "legoloop-magro", "livro.json"),
		`{"nome_base": "magro", "titulo": "legoloop-magro"}`)
	exigirEtapaComErro(t, fecharSlug(t, repo, "magro"), 2, "stub")
}

func TestFecharRecusaDAGSemNos(t *testing.T) {
	repo := t.TempDir()
	semearDocs(t, repo, "vazio")
	escrever(t, filepath.Join(repo, dirDialetoAntigo, "legoloop-vazio", caminhos.ArquivoDAG),
		`{"name": "legoloop-vazio", "version": "1.0.0", "nodes": []}`)
	exigirEtapaComErro(t, fecharSlug(t, repo, "vazio"), 4, "stub")
}

func TestFecharManualDeOutroProcessoNaoConta(t *testing.T) {
	repo := t.TempDir()
	semearDocs(t, repo, "video-rea-transcrever")
	semearDAG(t, repo, "video-rea")
	semearDocs(t, repo, "video-rea")
	if err := os.Remove(filepath.Join(repo, "docs", "manuals", "2026-11-15_MANUAL_VIDEO_REA.md")); err != nil {
		t.Fatal(err)
	}
	exigirEtapaComErro(t, fecharSlug(t, repo, "video-rea"), 1, "artefato obrigatório ausente")
}

func TestFecharComArtefatosReaisTemSucesso(t *testing.T) {
	repo := t.TempDir()
	semearDocs(t, repo, "pleno")
	semearDAG(t, repo, "pleno")
	res := fecharSlug(t, repo, "pleno")
	if !res.Sucesso {
		t.Fatalf("esperava sucesso com artefatos reais: %+v", res.Etapas)
	}
}

func TestFecharDryRunTambemRecusaArtefatoAusente(t *testing.T) {
	repo := t.TempDir()
	res, err := fechar.Executar(context.Background(), fechar.Request{Slug: "seco", DryRun: true}, depsFechar(repo))
	if err != nil {
		t.Fatal(err)
	}
	exigirEtapaComErro(t, res, 3, "artefato obrigatório ausente")
}

func TestFecharAuditoriaDoLivroReprovadaFalhaEtapa2(t *testing.T) {
	casos := []struct {
		nome    string
		auditor func(string) (bool, error)
		falha   bool
	}{
		{"reprova sem erro", func(string) (bool, error) { return false, nil }, true},
		{"erro do auditor", func(string) (bool, error) { return true, os.ErrInvalid }, true},
		{"aprova", func(string) (bool, error) { return true, nil }, false},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			repo := t.TempDir()
			semearDocs(t, repo, "auditado")
			semearDAG(t, repo, "auditado")
			deps := depsFechar(repo)
			deps.Auditor = c.auditor
			res, err := fechar.Executar(context.Background(), fechar.Request{Slug: "auditado"}, deps)
			if err != nil {
				t.Fatal(err)
			}
			if c.falha {
				exigirEtapaComErro(t, res, 2, "auditoria do livro")
			} else if !res.Sucesso {
				t.Fatalf("auditoria aprovada deveria fechar: %+v", res.Etapas)
			}
		})
	}
}
