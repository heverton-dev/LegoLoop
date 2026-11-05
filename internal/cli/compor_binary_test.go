package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"legoloop/internal/processo"
)

// compRouter é um 9Router em loopback que lista sonnet e code-fast. Assim os nós de skill (modelo sonnet)
// rodam no router falso, sem rede e sem o claude CLI. Cada resposta vira o artefato do nó.
func compRouter(t *testing.T, chamadas *atomic.Int32) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/models":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]string{{"id": "sonnet"}, {"id": "code-fast"}}})
		case "/v1/chat/completions":
			chamadas.Add(1)
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]string{"content": "# Artefato\n\nfeito\n"}}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func compEtapa(proc, id string, depende ...string) processo.Etapa {
	art := "out/" + strings.TrimPrefix(proc, "legoloop-") + "/" + id + ".md"
	e := processo.Etapa{ID: id, Titulo: "Etapa " + id, Tipo: "construir", Peca: "aidd-tdd", Depende: depende,
		Produz: processo.Artefato{Caminho: art, Tipo: "markdown"}, Gate: "file_exists('" + art + "')"}
	for _, d := range depende {
		e.Consome = append(e.Consome, "out/"+strings.TrimPrefix(proc, "legoloop-")+"/"+d+".md")
	}
	return e
}

// compSemear grava legoloop-a (3 etapas) e legoloop-b (2 etapas) como o fabricar grava, e a peça aidd-tdd.
func compSemear(t *testing.T) (root, vault, proj string) {
	t.Helper()
	root = t.TempDir()
	vault, proj = filepath.Join(root, "vault"), filepath.Join(root, "proj")
	data := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	var files []processo.File
	for nome, n := range map[string]int{"legoloop-a": 3, "legoloop-b": 2} {
		p := processo.Processo{Nome: nome, Versao: "1.0.0", Natureza: "software", Objetivo: "Objetivo de " + nome + ".", CriadoEm: data}
		for i := range n {
			id := "e0" + string(rune('1'+i))
			if i == 0 {
				p.Etapas = append(p.Etapas, compEtapa(nome, id))
			} else {
				p.Etapas = append(p.Etapas, compEtapa(nome, id, p.Etapas[i-1].ID))
			}
		}
		files = append(files, compArquivos(t, p, vault, proj, data)...)
	}
	files = append(files, processo.File{Path: filepath.Join(vault, "skills", "aidd-tdd", "SKILL.md"), Data: []byte("# aidd-tdd\n\nEscreva o artefato.\n")})
	st := processo.Store{Roots: processo.Roots{Vault: vault, Docs: proj}}
	if err := st.Write(context.Background(), files, processo.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
	return root, vault, proj
}

func compArquivos(t *testing.T, p processo.Processo, vault, proj string, data time.Time) []processo.File {
	t.Helper()
	pr, err := processo.Projetar(p, data)
	if err != nil {
		t.Fatal(err)
	}
	pj, err := p.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	doc := func(sub, file string) string { return filepath.Join(proj, "docs", sub, file) }
	return []processo.File{
		{Path: filepath.Join(vault, "skills", p.Nome, "SKILL.md"), Data: []byte(pr.Skill)},
		{Path: filepath.Join(vault, "loops", p.Nome+".loop.json"), Data: pr.Loop},
		{Path: filepath.Join(vault, "loops", p.Nome+".dag.json"), Data: pr.DAG},
		{Path: processo.CaminhoNoCofre(vault, p.Nome), Data: pj},
		{Path: doc("specs", pr.SpecNome), Data: []byte(pr.Spec)},
		{Path: doc("plans", pr.PlanNome), Data: []byte(pr.Plan)},
		{Path: doc("routes", pr.RouteNome), Data: pr.Route},
	}
}

// TestComporBinario prova AC-E03.4 pelo executável: liga b@1.0.0 em a, valida e roda o DAG (o subloop executa
// b), recusa o ciclo e o processo inexistente sem mudar a árvore, e o listar mostra b usado por a.
func TestComporBinario(t *testing.T) {
	bin := buildBinary(t)
	var chamadas atomic.Int32
	router := compRouter(t, &chamadas)
	root, vault, proj := compSemear(t)
	r := runBinaryEm(t, bin, router, root, proj, "compor", "legoloop-a", "legoloop-b@1.0.0", "--vault", vault, "--docs", proj, "--yes", "--json")
	var res struct {
		Para, Composto, Etapa string
		Arquivos              []struct{ Backup string }
	}
	if err := json.Unmarshal([]byte(r.out), &res); err != nil || r.code != 0 || res.Etapa != "e04" || res.Para != "1.1.0" ||
		res.Composto != "legoloop-b@1.0.0" || len(res.Arquivos) != 8 || !strings.HasSuffix(res.Arquivos[7].Backup, ".processo.json.legoloop-bak") {
		t.Fatalf("compor --json: código %d, json %v %+v\nstdout:\n%s\nstderr:\n%s", r.code, err, res, r.out, r.errb)
	}
	if !strings.Contains(r.errb, "COMPOR legoloop-a") || strings.Contains(r.out, "COMPOR") {
		t.Errorf("com --json, o plano e o cartão vão para o stderr:\n%s", r.errb)
	}
	dagA := filepath.Join(vault, "loops", "legoloop-a.dag.json")
	if v := runBinaryEm(t, bin, router, root, proj, "validate", dagA); v.code != 0 || !strings.Contains(v.out, "e04") {
		t.Fatalf("validate: código %d\n%s%s", v.code, v.out, v.errb)
	}
	compRodar(t, bin, router, root, proj, vault, dagA, &chamadas)
	compListar(t, bin, router, root, proj, vault)
	antes := refinArvore(t, root)
	falhas := []struct {
		args  []string
		trech string
	}{
		{[]string{"legoloop-b", "legoloop-a@1.1.0"}, "ciclo"},
		{[]string{"legoloop-a", "legoloop-zzz@1.0.0"}, "não existe"},
		{[]string{"legoloop-a", "legoloop-b@9.0.0"}, "versão"},
	}
	for _, f := range falhas {
		args := append(append([]string{"compor"}, f.args...), "--vault", vault, "--docs", proj, "--yes")
		r := runBinaryEm(t, bin, router, root, proj, args...)
		if r.code != 1 || !strings.Contains(r.errb, f.trech) {
			t.Errorf("compor %v: código %d (esperava 1) e stderr sem %q:\n%s", f.args, r.code, f.trech, r.errb)
		}
	}
	refinIgual(t, antes, refinArvore(t, root))
}

// compRodar roda o DAG de a com o router falso: o nó e04 é o subloop de b e grava o artefato final de b.
func compRodar(t *testing.T, bin string, router *httptest.Server, root, proj, vault, dagA string, chamadas *atomic.Int32) {
	t.Helper()
	work := t.TempDir()
	r := runBinaryEm(t, bin, router, root, proj, "run", dagA, "--workdir", work, "--skills", vault, "--yes", "--json")
	var rep struct {
		OK    bool `json:"ok"`
		Nodes []struct {
			ID, Status, Model string
		} `json:"nodes"`
	}
	if err := json.Unmarshal([]byte(r.out), &rep); err != nil || r.code != 0 || !rep.OK {
		t.Fatalf("run: código %d, json %v\nstdout:\n%s\nstderr:\n%s", r.code, err, r.out, r.errb)
	}
	i := slices.IndexFunc(rep.Nodes, func(n struct{ ID, Status, Model string }) bool { return n.ID == "e04" })
	if i < 0 || rep.Nodes[i].Status != "passed" || rep.Nodes[i].Model != "subloop/legoloop-b" {
		t.Fatalf("o nó e04 não rodou o subloop de b: %+v", rep.Nodes)
	}
	for _, n := range rep.Nodes {
		if n.ID != "e04" && !strings.HasPrefix(n.Model, "9router/") {
			t.Errorf("o nó %s rodou fora do router falso: %s", n.ID, n.Model)
		}
	}
	// TICK-03.4b: o artefato de b chega ao pai em out/<b>/<arquivo>, não no caminho que ele tem dentro de b
	if _, err := os.Stat(filepath.Join(work, "out", "legoloop-b", "e02.md")); err != nil || chamadas.Load() != 5 {
		t.Errorf("artefato do subloop: %v; chamadas ao router: %d (esperava 3 de a e 2 de b)", err, chamadas.Load())
	}
	if _, err := os.Stat(filepath.Join(work, "out", "b", "e02.md")); err == nil {
		t.Error("o caminho de dentro de b não pode aparecer no pai")
	}
}

func compListar(t *testing.T, bin string, router *httptest.Server, root, proj, vault string) {
	t.Helper()
	r := runBinaryEm(t, bin, router, root, proj, "listar", "--vault", vault, "--json")
	var res struct {
		Processos []struct {
			Nome     string   `json:"nome"`
			Versao   string   `json:"versao"`
			UsadoPor []string `json:"usado_por"`
		} `json:"processos"`
	}
	if err := json.Unmarshal([]byte(r.out), &res); err != nil || r.code != 0 || len(res.Processos) != 2 {
		t.Fatalf("listar: código %d, json %v\n%s%s", r.code, err, r.out, r.errb)
	}
	a, b := res.Processos[0], res.Processos[1]
	if a.Nome != "legoloop-a" || a.Versao != "1.1.0" || len(a.UsadoPor) != 0 || !slices.Equal(b.UsadoPor, []string{"legoloop-a"}) {
		t.Errorf("listar: %+v; esperava a 1.1.0 e b usado por a", res.Processos)
	}
	txt := runBinaryEm(t, bin, router, root, proj, "listar", "--vault", vault)
	if linha := linhaCom(txt.out, "legoloop-b"); !strings.Contains(linha, "legoloop-a") {
		t.Errorf("o texto do listar não mostra b usado por a:\n%s", txt.out)
	}
}

func linhaCom(texto, prefixo string) string {
	for _, l := range strings.Split(texto, "\n") {
		if strings.HasPrefix(l, prefixo) {
			return l
		}
	}
	return ""
}

// TestComporBinarioInterativo: sem --yes, 0 cancela sem mudar a árvore; 2 personaliza o título e o depende.
func TestComporBinarioInterativo(t *testing.T) {
	bin := buildBinary(t)
	var chamadas atomic.Int32
	router := compRouter(t, &chamadas)
	root, vault, proj := compSemear(t)
	args := []string{"compor", "legoloop-a", "legoloop-b@1.0.0", "--vault", vault, "--docs", proj}
	antes := refinArvore(t, root)
	if r := runBinaryStdin(t, bin, router, root, proj, "0\n", args...); r.code != 1 || !strings.Contains(r.errb, "cancelada") {
		t.Errorf("cancelar: código %d\n%s%s", r.code, r.out, r.errb)
	}
	refinIgual(t, antes, refinArvore(t, root))
	if r := runBinaryEm(t, bin, router, root, proj, args...); r.code != 1 {
		t.Errorf("sem stdin: código %d; esperava 1", r.code)
	}
	r := runBinaryStdin(t, bin, router, root, proj, "2\nRodar b\ne99\n2\nRodar b\ne02\n1\n", args...)
	if r.code != 0 || !strings.Contains(r.out, "e99") {
		t.Fatalf("personalizar: código %d\nstdout:\n%s\nstderr:\n%s", r.code, r.out, r.errb)
	}
	p, err := processo.Ler(processo.CaminhoNoCofre(vault, "legoloop-a"))
	if err != nil {
		t.Fatal(err)
	}
	if e := p.Etapas[3]; e.Titulo != "Rodar b" || !slices.Equal(e.Depende, []string{"e02"}) {
		t.Errorf("etapa e04 %+v; esperava o título e o depende do personalizar", e)
	}
}

// TestComporBinarioUso: argumentos que faltam ou pedido malformado saem com 2.
func TestComporBinarioUso(t *testing.T) {
	bin := buildBinary(t)
	var chamadas atomic.Int32
	router := compRouter(t, &chamadas)
	root, vault, proj := compSemear(t)
	for _, args := range [][]string{
		{"compor"}, {"compor", "legoloop-a"},
		{"compor", "legoloop-a", "legoloop-b", "--vault", vault, "--docs", proj, "--yes"},
		{"compor", "legoloop-a", "legoloop-b@1.0.0", "--bogus"},
	} {
		r := runBinaryEm(t, bin, router, root, proj, args...)
		if r.code != 2 || len(args) < 3 && !strings.Contains(r.errb, "Uso: legoloop compor") {
			t.Errorf("%v: código %d; esperava 2 com o uso\n%s", args, r.code, r.errb)
		}
	}
	// flags antes dos nomes; --titulo e --depende chegam à etapa nova
	r := runBinaryEm(t, bin, router, root, proj, "compor", "--vault", vault, "--docs", proj, "--yes", "--titulo", "Rodar b",
		"--depende", "e01", "legoloop-a", "legoloop-b@1.0.0")
	if r.code != 0 || !strings.Contains(r.out, "Processo legoloop-b@1.0.0 ligado a legoloop-a na etapa e04") {
		t.Fatalf("flags antes dos nomes: código %d\n%s%s", r.code, r.out, r.errb)
	}
	p, err := processo.Ler(processo.CaminhoNoCofre(vault, "legoloop-a"))
	if err != nil {
		t.Fatal(err)
	}
	if e := p.Etapas[3]; e.Titulo != "Rodar b" || !slices.Equal(e.Depende, []string{"e01"}) {
		t.Errorf("etapa e04 %+v; esperava o título e o depende das flags", e)
	}
}
