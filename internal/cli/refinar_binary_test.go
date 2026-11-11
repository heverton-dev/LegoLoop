package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/paths"
	"legoloop/internal/processo"
)

const refinNome = "legoloop-x"

func refinEtapa(id, tipo, peca string, depende ...string) processo.Etapa {
	e := processo.Etapa{
		ID: id, Titulo: "Etapa " + id, Tipo: tipo, Peca: peca, Depende: depende,
		Produz: processo.Artefato{Caminho: "out/" + id + ".md", Tipo: "markdown"},
		Gate:   "file_exists('out/" + id + ".md')",
	}
	if len(depende) > 0 {
		e.Consome = []string{"out/" + depende[0] + ".md"}
	}
	return e
}

// refinSemear grava o processo legoloop-x 1.0.0 como o fabricar grava (cofre e docs/) e devolve as raízes.
func refinSemear(t *testing.T) (root, vault, proj string) {
	t.Helper()
	root = t.TempDir()
	vault, proj = filepath.Join(root, "vault"), filepath.Join(root, "proj")
	matrizDeTeste(t, vault)
	data := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	p := processo.Processo{
		Nome: refinNome, Versao: "1.0.0", Natureza: "software", Objetivo: "Extrair o áudio de um vídeo.", CriadoEm: data,
		Etapas: []processo.Etapa{refinEtapa("e01", "clarificar", "sem-perfil"), refinEtapa("e02", "construir", "sem-perfil", "e01"),
			refinEtapa("e03", "verificar", "aidd-tdd", "e02")},
	}
	pr, err := processo.Projetar(p, data)
	if err != nil {
		t.Fatal(err)
	}
	pj, err := p.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	doc := func(sub, file string) string { return filepath.Join(proj, "docs", sub, file) }
	files := []processo.File{
		{Path: filepath.Join(vault, "skills", refinNome, "SKILL.md"), Data: []byte(pr.Skill)},
		{Path: filepath.Join(vault, paths.DirAntigo, refinNome+paths.SufixoAntigo), Data: pr.Loop},
		{Path: filepath.Join(vault, paths.DirAntigo, refinNome+".dag.json"), Data: pr.DAG},
		{Path: processo.CaminhoNoCofre(vault, refinNome), Data: pj},
		{Path: doc("intents", "2026-11-05_INTENT_X.md"), Data: []byte("# INTENT: x\n\nTexto original da entrevista.\n")},
		{Path: doc("specs", pr.SpecNome), Data: []byte(pr.Spec)},
		{Path: doc("plans", pr.PlanNome), Data: []byte(pr.Plan)},
		{Path: doc("routes", pr.RouteNome), Data: pr.Route},
	}
	st := processo.Store{Roots: processo.Roots{Vault: vault, Docs: proj}}
	if err := st.Write(context.Background(), files, processo.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
	return root, vault, proj
}

// refinArvore devolve o hash de cada arquivo sob root: é a prova de "árvore idêntica".
func refinArvore(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, f := range filesUnder(root) {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(b)
		out[f] = hex.EncodeToString(sum[:])
	}
	return out
}

func refinIgual(t *testing.T, antes, depois map[string]string) {
	t.Helper()
	for p, h := range antes {
		if depois[p] != h {
			t.Errorf("o arquivo mudou ou sumiu: %s", p)
		}
	}
	for p := range depois {
		if _, ok := antes[p]; !ok {
			t.Errorf("apareceu arquivo novo: %s", p)
		}
	}
}

func refinLer(t *testing.T, vault string) processo.Processo {
	t.Helper()
	p, err := processo.Ler(processo.CaminhoNoCofre(vault, refinNome))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// refinRodar executa o binário: legoloop refinar legoloop-x --etapa e02 --peca aidd-tdd --caso tdd, mais extras.
func refinRodar(t *testing.T, bin, root, vault, proj, stdin string, extra ...string) binRun {
	t.Helper()
	args := append([]string{"refinar", refinNome, "--etapa", "e02", "--peca", "aidd-tdd", "--caso", "tdd", "--vault", vault, "--docs", proj}, extra...)
	return runBinaryStdin(t, bin, binRouter(t, 40), root, proj, stdin, args...)
}

// TestRefinarBinario prova o AC-E03.2 pelo executável: seguir (com --yes e no terminal), cancelar, personalizar
// e os pedidos inválidos.
func TestRefinarBinario(t *testing.T) {
	bin := buildBinary(t)
	t.Run("--yes: o diff cita os sete, versão 1.1.0, backup e histórico", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		antes := refinArvore(t, root)
		r := refinRodar(t, bin, root, vault, proj, "", "--yes")
		if r.code != 0 {
			t.Fatalf("exit %d\n%s\n%s", r.code, r.out, r.errb)
		}
		for _, quer := range []string{"INTENT ", "SPEC ", "PLAN ", "ROUTE ", "DAG ", "skill ", "versão 1.0.0 para 1.1.0",
			"AÇÃO: Refinar o processo legoloop-x", "AVISO:", ".legoloop-bak", "Segue sem perguntar", "Processo legoloop-x refinado",
			"ARQUIVO   ESTADO  CAMINHO", "legoloop-x.processo.json.legoloop-bak"} {
			if !strings.Contains(r.out, quer) {
				t.Errorf("falta %q na saída:\n%s", quer, r.out)
			}
		}
		p := refinLer(t, vault)
		if p.Versao != "1.1.0" || !slices.Equal(p.Historico, []string{"1.0.0"}) || p.Etapas[1].Peca != "aidd-tdd" || !slices.Equal(p.Etapas[1].Casos, []string{"tdd"}) {
			t.Errorf("processo depois do refinar: %+v", p)
		}
		bak := filepath.Join(vault, paths.DirAntigo, refinNome+".processo.json.legoloop-bak")
		if b, err := os.ReadFile(bak); err != nil || !strings.Contains(string(b), `"versao": "1.0.0"`) {
			t.Errorf("o backup não guarda a versão 1.0.0: %v\n%s", err, b)
		}
		for f, h := range antes {
			if strings.HasSuffix(f, ".processo.json") && refinArvore(t, root)[f] == h {
				t.Errorf("o processo.json não mudou: %s", f)
			}
		}
	})
	t.Run("terminal: seguir com 1", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		r := refinRodar(t, bin, root, vault, proj, "1\n")
		if r.code != 0 || refinLer(t, vault).Versao != "1.1.0" {
			t.Fatalf("exit %d, versão %s\n%s\n%s", r.code, refinLer(t, vault).Versao, r.out, r.errb)
		}
		if !strings.Contains(r.out, "Escolha: [1] seguir  [2] personalizar  [0] cancelar") {
			t.Errorf("falta a linha de escolha:\n%s", r.out)
		}
	})
	t.Run("terminal: cancelar com 0 deixa a árvore idêntica", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		antes := refinArvore(t, root)
		r := refinRodar(t, bin, root, vault, proj, "0\n")
		if r.code != 1 || !strings.Contains(r.errb, "Nenhum arquivo mudou") {
			t.Fatalf("exit %d; esperava 1 e a mensagem de cancelamento\n%s\n%s", r.code, r.out, r.errb)
		}
		refinIgual(t, antes, refinArvore(t, root))
	})
	t.Run("terminal: Enter vazio e fim da entrada cancelam", func(t *testing.T) {
		for _, stdin := range []string{"\n", ""} {
			root, vault, proj := refinSemear(t)
			antes := refinArvore(t, root)
			if r := refinRodar(t, bin, root, vault, proj, stdin); r.code != 1 {
				t.Errorf("stdin %q: exit %d; esperava 1\n%s\n%s", stdin, r.code, r.out, r.errb)
			}
			refinIgual(t, antes, refinArvore(t, root))
		}
	})
	t.Run("terminal: personalizar os casos e seguir", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		r := refinRodar(t, bin, root, vault, proj, "2\n\noutro-caso\n1\n")
		if r.code != 0 {
			t.Fatalf("exit %d\n%s\n%s", r.code, r.out, r.errb)
		}
		if got := refinLer(t, vault).Etapas[1].Casos; !slices.Equal(got, []string{"outro-caso"}) {
			t.Errorf("casos %v; esperava o valor do personalizar", got)
		}
	})
	t.Run("terminal: personalizar a peça e seguir", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		r := refinRodar(t, bin, root, vault, proj, "2\naidd-site-dominator\n\n1\n")
		if got := refinLer(t, vault).Etapas[1]; r.code != 0 || got.Peca != "aidd-site-dominator" || !slices.Equal(got.Casos, []string{"tdd"}) {
			t.Fatalf("exit %d, etapa %+v\n%s\n%s", r.code, got, r.out, r.errb)
		}
	})
	t.Run("terminal: personalizar inválido mostra o motivo e mantém o plano", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		r := refinRodar(t, bin, root, vault, proj, "2\n\nnenhum\n1\n")
		if got := refinLer(t, vault).Etapas[1]; r.code != 0 || !slices.Equal(got.Casos, []string{"tdd"}) ||
			!strings.Contains(r.out, "erro: refinar: mudança inválida: A etapa e02 precisa de ao menos um caso de uso") {
			t.Fatalf("exit %d, etapa %+v\n%s\n%s", r.code, got, r.out, r.errb)
		}
	})
	t.Run("terminal: três entradas inválidas cancelam", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		antes := refinArvore(t, root)
		r := refinRodar(t, bin, root, vault, proj, "x\nx\nx\n")
		if r.code != 1 || strings.Count(r.out, "Entrada inválida: responda 1, 2 ou 0.") != 3 || strings.Count(r.out, "AÇÃO: Refinar") != 3 {
			t.Fatalf("exit %d\n%s\n%s", r.code, r.out, r.errb)
		}
		refinIgual(t, antes, refinArvore(t, root))
	})
	t.Run("sem --vault e --docs: cofre em HOME/.aidd/canonical e projeto no diretório atual", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		padrao := filepath.Join(root, ".aidd", "canonical")
		if err := os.MkdirAll(filepath.Dir(padrao), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(vault, padrao); err != nil {
			t.Fatal(err)
		}
		args := []string{"refinar", refinNome, "--etapa", "e02", "--peca", "aidd-tdd", "--caso", "tdd", "--yes"}
		r := runBinaryStdin(t, bin, binRouter(t, 40), root, proj, "", args...)
		if r.code != 0 || refinLer(t, padrao).Versao != "1.1.0" {
			t.Fatalf("exit %d\n%s\n%s", r.code, r.out, r.errb)
		}
		if b, err := os.ReadFile(filepath.Join(proj, "docs", "intents", "2026-11-05_INTENT_X.md")); err != nil || !strings.Contains(string(b), "Refinamento 1.1.0") {
			t.Errorf("o INTENT do diretório atual não recebeu o refinamento: %v", err)
		}
		if _, err := os.Stat(filepath.Join(proj, "docs", "docs")); err == nil {
			t.Error("o refinar gravou documentos fora de docs/ do diretório atual")
		}
	})
	t.Run("--json: stdout é JSON puro", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		r := refinRodar(t, bin, root, vault, proj, "", "--yes", "--json")
		var res struct {
			De, Para string
			Arquivos []struct{ Tipo, Estado string }
		}
		if err := json.Unmarshal([]byte(r.out), &res); err != nil || r.code != 0 {
			t.Fatalf("exit %d; stdout deve ser JSON: %v\n%s\n%s", r.code, err, r.out, r.errb)
		}
		if res.De != "1.0.0" || res.Para != "1.1.0" || len(res.Arquivos) != 8 {
			t.Errorf("resultado: %+v", res)
		}
	})
	t.Run("pedidos inválidos não gravam", func(t *testing.T) {
		casos := []struct {
			rotulo string
			args   []string
			code   int
			msg    string
		}{
			{"etapa inexistente", []string{"refinar", refinNome, "--etapa", "e99", "--peca", "aidd-tdd", "--caso", "tdd", "--yes"}, 2, "não tem a etapa e99"},
			{"sem caso", []string{"refinar", refinNome, "--etapa", "e02", "--peca", "aidd-tdd", "--yes"}, 2, "precisa de ao menos um caso"},
			{"sem etapa", []string{"refinar", refinNome, "--yes"}, 2, "Uso: legoloop refinar"},
			{"nome inválido", []string{"refinar", "x", "--etapa", "e02", "--caso", "tdd", "--yes"}, 2, "não é de processo"},
			{"processo ausente", []string{"refinar", "legoloop-nada", "--etapa", "e02", "--caso", "tdd", "--yes"}, 1, "não existe no cofre"},
		}
		for _, tc := range casos {
			root, vault, proj := refinSemear(t)
			antes := refinArvore(t, root)
			args := append(tc.args, "--vault", vault, "--docs", proj)
			r := runBinaryStdin(t, bin, binRouter(t, 40), root, proj, "", args...)
			if r.code != tc.code || !strings.Contains(r.errb, tc.msg) {
				t.Errorf("%s: exit %d, stderr %q; esperava %d e %q", tc.rotulo, r.code, r.errb, tc.code, tc.msg)
			}
			refinIgual(t, antes, refinArvore(t, root))
		}
	})
	t.Run("peça fora da matriz: aviso, sem barrar", func(t *testing.T) {
		root, vault, proj := refinSemear(t)
		r := refinRodar(t, bin, root, vault, proj, "", "--yes", "--peca", "peca-que-nao-existe")
		if r.code != 0 || !strings.Contains(r.out, "A peça peca-que-nao-existe não está na matriz") {
			t.Fatalf("exit %d\n%s\n%s", r.code, r.out, r.errb)
		}
	})
}

// TestRefinarSemTerminal: sem --yes e sem stdin o comando não pergunta nem grava; o help lista o verbo.
func TestRefinarSemTerminal(t *testing.T) {
	root, vault, proj := refinSemear(t)
	antes := refinArvore(t, root)
	var out, errb bytes.Buffer
	args := []string{"refinar", refinNome, "--etapa", "e02", "--peca", "aidd-tdd", "--caso", "tdd", "--vault", vault, "--docs", proj}
	code := Run(context.Background(), args, Deps{Stdout: &out, Stderr: &errb})
	if code != 1 || !strings.Contains(errb.String(), "Não há terminal para confirmar. Use --yes") {
		t.Fatalf("exit %d\n%s\n%s", code, out.String(), errb.String())
	}
	refinIgual(t, antes, refinArvore(t, root))
	out.Reset()
	if code := Run(context.Background(), []string{"help"}, Deps{Stdout: &out, Stderr: &errb}); code != 0 || !strings.Contains(out.String(), "refinar    <nome> --etapa ID") {
		t.Errorf("o help não lista o refinar (exit %d):\n%s", code, out.String())
	}
}
