package fabricar_test

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/assets"
	"github.com/heverton-dev/LegoLoop/internal/fabricar"
	"github.com/heverton-dev/LegoLoop/internal/matriz"
)

func criarSkillDeclarada(t *testing.T, projDir, nome, etapa, caso string) string {
	t.Helper()
	skillDir := filepath.Join(projDir, ".claude", "skills", nome)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projDir, "go.mod"), []byte("module "+nome+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	conteudo := "---\nname: " + nome + "\ndescription: skill " + nome + " para desenvolvimento de software\netapas: [" + etapa + "]\nnaturezas: [software]\ncasos_de_uso: [" + caso + "]\n---\n" +
		"# " + nome + "\n\n" +
		"## I/O Contracts\n- Consumes: spec.md\n- Produces: artifact.json\n\n" +
		"## Negative Guardrails\n- NEVER skip tests or leave stubs.\n- DO NOT modify unrelated files.\n\n" +
		"## Verification\n- Run `go test ./...` and assert exit_code == 0.\n"
	path := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(path, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func gerarMatrizSemCerebroTeste(t *testing.T, home, vault, raiz string) {
	t.Helper()
	voc, err := matriz.LoadVocab(vault)
	if err != nil {
		t.Fatal(err)
	}
	scanner := assets.NewScanner(assets.Config{
		Home:           home,
		AutoDiscover:   true,
		DiscoveryRoots: []string{raiz},
		LookPath:       func(string) (string, error) { return "", os.ErrNotExist },
		Getenv:         func(string) string { return "" },
	})
	md := matriz.Deps{
		Home:      home,
		Catalogar: matriz.Catalogador{Vocab: voc, Brain: nil},
		Scan:      func(ctx context.Context) ([]assets.Item, error) { return scanner.Scan(ctx).Items, nil },
		Now:       func() time.Time { return fabNow },
	}
	if _, err := matriz.Generate(context.Background(), matriz.Options{Vault: vault}, md); err != nil {
		t.Fatalf("matriz.Generate inicial: %v", err)
	}
}

func TestFabricarAtualizaMatrizDefasada(t *testing.T) {
	t.Run("chama Atualizar uma vez e usa a peça nova na Fusão", func(t *testing.T) {
		home := t.TempDir()
		vault := filepath.Join(home, ".aidd", "canonical")
		raiz := filepath.Join(home, "projetos")
		docs := filepath.Join(home, "docs")

		criarSkillDeclarada(t, filepath.Join(raiz, "proj-alfa"), "alfa", "construir", "tdd")
		gerarMatrizSemCerebroTeste(t, home, vault, raiz)

		// Peça nova colocada na raiz de teste sem rodar `matriz gerar` antes.
		criarSkillDeclarada(t, filepath.Join(raiz, "proj-beta"), "beta", "verificar", "verificacao-de-qualidade")

		chamadas := 0
		llm := &routerLLM{interviewReplies: []string{lowConfidence}}
		deps := fabricar.Deps{
			LLM:            llm,
			Clock:          func() time.Time { return fabNow },
			Roots:          fabricar.Roots{Vault: vault, Docs: docs},
			FS:             fabricar.OSFS(),
			Home:           home,
			DiscoveryRoots: []string{raiz},
			Atualizar: func(ctx context.Context) error {
				chamadas++
				gerarMatrizSemCerebroTeste(t, home, vault, raiz)
				return nil
			},
		}
		req := fabricar.Request{
			RawName:     "Extrator de Áudio",
			Goal:        "extrair áudio de vídeos",
			Category:    fabricar.CategoryBuild,
			ModelChoice: "sonnet",
			DryRun:      true,
		}

		res, err := fabricar.Fabricar(context.Background(), req, deps)
		if err != nil {
			t.Fatalf("Fabricar falhou: %v", err)
		}
		if chamadas != 1 {
			t.Fatalf("Atualizar deveria ser chamado 1 vez, foi chamado %d vezes", chamadas)
		}
		var usadas []string
		for _, p := range res.CuratedPieces {
			usadas = append(usadas, p.Name)
		}
		if !slices.Contains(usadas, "beta") {
			t.Errorf("peça nova beta deveria aparecer na Matriz usada pela Fusão; usadas=%v vistas=%v decisoes=%+v", usadas, res.Report.Varredura.PecasVistas, res.Report.Fusao.Decisoes)
		}
		if res.Report == nil || !res.Report.Varredura.Defasagem.Defasada() {
			t.Errorf("relatório deveria registrar a defasagem detectada: %+v", res.Report)
		}
	})

	t.Run("matriz em dia não chama Atualizar", func(t *testing.T) {
		home := t.TempDir()
		vault := filepath.Join(home, ".aidd", "canonical")
		raiz := filepath.Join(home, "projetos")
		docs := filepath.Join(home, "docs")

		criarSkillDeclarada(t, filepath.Join(raiz, "proj-alfa"), "alfa", "construir", "tdd")
		criarSkillDeclarada(t, filepath.Join(raiz, "proj-beta"), "beta", "verificar", "verificacao-de-qualidade")
		gerarMatrizSemCerebroTeste(t, home, vault, raiz)

		chamadas := 0
		llm := &routerLLM{interviewReplies: []string{lowConfidence}}
		deps := fabricar.Deps{
			LLM:            llm,
			Clock:          func() time.Time { return fabNow },
			Roots:          fabricar.Roots{Vault: vault, Docs: docs},
			FS:             fabricar.OSFS(),
			Home:           home,
			DiscoveryRoots: []string{raiz},
			Atualizar: func(ctx context.Context) error {
				chamadas++
				return nil
			},
		}
		req := fabricar.Request{
			RawName:     "Extrator de Áudio",
			Goal:        "extrair áudio de vídeos",
			Category:    fabricar.CategoryBuild,
			ModelChoice: "sonnet",
			DryRun:      true,
		}

		if _, err := fabricar.Fabricar(context.Background(), req, deps); err != nil {
			t.Fatalf("Fabricar falhou: %v", err)
		}
		if chamadas != 0 {
			t.Errorf("Atualizar não deveria ser chamado com Matriz em dia; chamadas=%d", chamadas)
		}
	})

	t.Run("Atualizar padrão sem cérebro lista peças pendentes no relatório", func(t *testing.T) {
		home := t.TempDir()
		vault := filepath.Join(home, ".aidd", "canonical")
		raiz := filepath.Join(home, "projetos")
		docs := filepath.Join(home, "docs")

		criarSkillDeclarada(t, filepath.Join(raiz, "proj-alfa"), "alfa", "construir", "tdd")
		criarSkillDeclarada(t, filepath.Join(raiz, "proj-beta"), "beta", "verificar", "verificacao-de-qualidade")
		gerarMatrizSemCerebroTeste(t, home, vault, raiz)

		// Peça sem casos_de_uso fica pendente (Retentar=true) no modo -sem-cerebro.
		pendDir := filepath.Join(raiz, "proj-pend", ".claude", "skills", "sem-perfil")
		if err := os.MkdirAll(pendDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(raiz, "proj-pend", "go.mod"), []byte("module pend\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(pendDir, "SKILL.md"), []byte("---\nname: sem-perfil\ndescription: xyz123\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		llm := &routerLLM{interviewReplies: []string{lowConfidence}}
		deps := fabricar.Deps{
			LLM:            llm,
			Clock:          func() time.Time { return fabNow },
			Roots:          fabricar.Roots{Vault: vault, Docs: docs},
			FS:             fabricar.OSFS(),
			Home:           home,
			DiscoveryRoots: []string{raiz},
		}
		req := fabricar.Request{
			RawName:     "Extrator de Áudio",
			Goal:        "extrair áudio de vídeos",
			Category:    fabricar.CategoryBuild,
			ModelChoice: "sonnet",
			DryRun:      true,
		}

		res, err := fabricar.Fabricar(context.Background(), req, deps)
		if err != nil {
			t.Fatalf("Fabricar falhou: %v", err)
		}
		if res.Report == nil || !slices.Contains(res.Report.Varredura.Pendentes, "sem-perfil") {
			t.Errorf("peça pendente sem-perfil deveria constar em Report.Varredura.Pendentes; veio %+v", res.Report)
		}
	})
}

func TestFabricarNaoAtualizaMatrizEmDia(t *testing.T) {
	home := t.TempDir()
	vault := filepath.Join(home, ".aidd", "canonical")
	raiz := filepath.Join(home, "projetos")
	docs := filepath.Join(home, "docs")
	criarSkillDeclarada(t, filepath.Join(raiz, "proj-alfa"), "alfa", "construir", "tdd")
	criarSkillDeclarada(t, filepath.Join(raiz, "proj-beta"), "beta", "verificar", "verificacao-de-qualidade")
	gerarMatrizSemCerebroTeste(t, home, vault, raiz)

	chamadas := 0
	deps := fabricar.Deps{
		LLM:            &routerLLM{interviewReplies: []string{lowConfidence}},
		Clock:          func() time.Time { return fabNow },
		Roots:          fabricar.Roots{Vault: vault, Docs: docs},
		FS:             fabricar.OSFS(),
		Home:           home,
		DiscoveryRoots: []string{raiz},
		Atualizar:      func(context.Context) error { chamadas++; return nil },
	}
	req := fabricar.Request{
		RawName: "Extrator de Áudio", Goal: "extrair áudio de vídeos",
		Category: fabricar.CategoryBuild, ModelChoice: "sonnet", DryRun: true,
	}
	res, err := fabricar.Fabricar(context.Background(), req, deps)
	if err != nil {
		t.Fatalf("Fabricar falhou: %v", err)
	}
	if chamadas != 0 {
		t.Fatalf("Atualizar deveria ter 0 chamadas com Matriz em dia; chamadas=%d", chamadas)
	}
	if res == nil || len(res.CuratedPieces) == 0 || (res.Report != nil && res.Report.Varredura.Defasagem.Defasada()) {
		t.Fatalf("esperava conclusão normal usando Matriz em dia; res=%+v", res)
	}
}

func scanRaizTeste(t *testing.T, home, raiz string) []assets.Item {
	t.Helper()
	scanner := assets.NewScanner(assets.Config{
		Home:           home,
		AutoDiscover:   true,
		DiscoveryRoots: []string{raiz},
		LookPath:       func(string) (string, error) { return "", os.ErrNotExist },
		Getenv:         func(string) string { return "" },
	})
	return scanner.Scan(context.Background()).Items
}

func rodarFabricarSemCerebroTeste(t *testing.T, home, vault, raiz, docs string) *fabricar.Result {
	t.Helper()
	deps := fabricar.Deps{
		LLM:            &routerLLM{interviewReplies: []string{lowConfidence}},
		Clock:          func() time.Time { return fabNow },
		Roots:          fabricar.Roots{Vault: vault, Docs: docs},
		FS:             fabricar.OSFS(),
		Home:           home,
		DiscoveryRoots: []string{raiz},
	}
	req := fabricar.Request{
		RawName: "Extrator de Áudio", Goal: "extrair áudio de vídeos",
		Category: fabricar.CategoryBuild, ModelChoice: "sonnet", DryRun: true,
	}
	res, err := fabricar.Fabricar(context.Background(), req, deps)
	if err != nil {
		t.Fatalf("Fabricar falhou: %v", err)
	}
	return res
}

func checarStalenessDisco(t *testing.T, home, vault, raiz string) matriz.Stale {
	t.Helper()
	reg, err := matriz.LerRegistro(filepath.Join(vault, "registry_matrix.json"))
	if err != nil {
		t.Fatalf("LerRegistro após fabricar: %v", err)
	}
	items := scanRaizTeste(t, home, raiz)
	stale, err := matriz.Staleness(context.Background(), reg, items, matriz.Deps{Home: home})
	if err != nil {
		t.Fatalf("Staleness após fabricar: %v", err)
	}
	return stale
}

func criarSkillSemPerfil(t *testing.T, raiz, nome string) {
	t.Helper()
	proj := filepath.Join(raiz, "proj-"+nome)
	dir := filepath.Join(proj, ".claude", "skills", nome)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module "+nome+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	conteudo := "---\nname: " + nome + "\ndescription: xyz123\n---\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestMatrizStatusZeroDepoisDoFabricar(t *testing.T) {
	t.Run("peca nova declarada zera defasagem em disco", func(t *testing.T) {
		home := t.TempDir()
		vault, raiz, docs := filepath.Join(home, ".aidd", "canonical"), filepath.Join(home, "projetos"), filepath.Join(home, "docs")
		criarSkillDeclarada(t, filepath.Join(raiz, "proj-alfa"), "alfa", "construir", "tdd")
		gerarMatrizSemCerebroTeste(t, home, vault, raiz)
		criarSkillDeclarada(t, filepath.Join(raiz, "proj-beta"), "beta", "verificar", "verificacao-de-qualidade")

		res := rodarFabricarSemCerebroTeste(t, home, vault, raiz, docs)
		stale := checarStalenessDisco(t, home, vault, raiz)
		if stale.Defasada() {
			t.Fatalf("esperava matriz em dia (Defasada()==false) após fabricar; veio %+v", stale)
		}
		if res.Report == nil || len(res.Report.Varredura.Pendentes) != 0 {
			t.Errorf("esperava 0 pendentes no relatório; veio %+v", res.Report)
		}
	})

	t.Run("peca pendente fica só em Retentar e aparece no relatório", func(t *testing.T) {
		home := t.TempDir()
		vault, raiz, docs := filepath.Join(home, ".aidd", "canonical"), filepath.Join(home, "projetos"), filepath.Join(home, "docs")
		criarSkillDeclarada(t, filepath.Join(raiz, "proj-alfa"), "alfa", "construir", "tdd")
		gerarMatrizSemCerebroTeste(t, home, vault, raiz)
		criarSkillDeclarada(t, filepath.Join(raiz, "proj-beta"), "beta", "verificar", "verificacao-de-qualidade")
		criarSkillSemPerfil(t, raiz, "sem-perfil")

		res := rodarFabricarSemCerebroTeste(t, home, vault, raiz, docs)
		stale := checarStalenessDisco(t, home, vault, raiz)
		if len(stale.Novas)+len(stale.Alteradas)+len(stale.Removidas) != 0 {
			t.Fatalf("esperava 0 novas/alteradas/removidas após fabricar; veio %+v", stale)
		}
		if len(stale.Retentar) != 1 || res.Report == nil || !slices.Contains(res.Report.Varredura.Pendentes, "sem-perfil") {
			t.Fatalf("esperava sem-perfil apenas em Retentar e Report.Varredura.Pendentes; stale=%+v report=%+v", stale, res.Report)
		}
	})
}
