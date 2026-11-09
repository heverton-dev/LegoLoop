package compor_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/compor"
	"legoloop/internal/dag"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

var inicio = time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)

const (
	nomeA = "legoloop-a"
	nomeB = "legoloop-b"
	nomeC = "legoloop-c"
)

// etapa grava em out/<processo>/<id>.md: cada processo tem os seus artefatos.
func etapa(proc, id, tipo, peca string, depende ...string) processo.Etapa {
	art := "out/" + strings.TrimPrefix(proc, "legoloop-") + "/" + id + ".md"
	e := processo.Etapa{
		ID: id, Titulo: "Etapa " + id, Tipo: tipo, Peca: peca, Depende: depende,
		Produz: processo.Artefato{Caminho: art, Tipo: "markdown"},
		Gate:   "file_exists('" + art + "')",
	}
	for _, d := range depende {
		e.Consome = append(e.Consome, "out/"+strings.TrimPrefix(proc, "legoloop-")+"/"+d+".md")
	}
	return e
}

func processoDe(nome string, n int) processo.Processo {
	p := processo.Processo{Nome: nome, Versao: "1.0.0", Natureza: "software", Objetivo: "Objetivo de " + nome + ".", CriadoEm: inicio}
	tipos := []string{"clarificar", "construir", "verificar"}
	for i := range n {
		id := "e0" + string(rune('1'+i))
		var dep []string
		if i > 0 {
			dep = []string{p.Etapas[i-1].ID}
		}
		p.Etapas = append(p.Etapas, etapa(nome, id, tipos[i%3], "aidd-tdd", dep...))
	}
	return p
}

type cofre struct{ vault, docs string }

func slugDe(nome string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimPrefix(nome, "legoloop-"), "-", "_"))
}

// gravar põe o processo no cofre e no projeto como o fabricar grava.
func gravar(t *testing.T, c cofre, p processo.Processo) {
	t.Helper()
	pr, err := processo.Projetar(p, inicio)
	if err != nil {
		t.Fatal(err)
	}
	pj, err := p.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	doc := func(sub, file string) string { return filepath.Join(c.docs, "docs", sub, file) }
	files := []processo.File{
		{Path: filepath.Join(c.vault, "skills", p.Nome, "SKILL.md"), Data: []byte(pr.Skill)},
		{Path: filepath.Join(c.vault, "loops", p.Nome+".loop.json"), Data: pr.Loop},
		{Path: filepath.Join(c.vault, "loops", p.Nome+".dag.json"), Data: pr.DAG},
		{Path: processo.CaminhoNoCofre(c.vault, p.Nome), Data: pj},
		{Path: doc("specs", pr.SpecNome), Data: []byte(pr.Spec)},
		{Path: doc("plans", pr.PlanNome), Data: []byte(pr.Plan)},
		{Path: doc("routes", pr.RouteNome), Data: pr.Route},
		{Path: doc("dags", pr.DAGNome), Data: pr.DAG},
		{Path: doc("intents", "2026-11-05_INTENT_"+slugDe(p.Nome)+".md"), Data: []byte("# INTENT: " + p.Nome + "\n\nTexto original.\n")},
	}
	st := processo.Store{Roots: processo.Roots{Vault: c.vault, Docs: c.docs}}
	if err := st.Write(context.Background(), files, processo.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
}

// semear cria o cofre com os processos dados (a com 3 etapas, os outros com 2).
func semear(t *testing.T, nomes ...string) cofre {
	t.Helper()
	root := t.TempDir()
	c := cofre{vault: filepath.Join(root, "vault"), docs: filepath.Join(root, "proj")}
	for _, n := range nomes {
		etapas := 2
		if n == nomeA {
			etapas = 3
		}
		gravar(t, c, processoDe(n, etapas))
	}
	return c
}

// arvore devolve o hash de cada arquivo sob as duas raízes: é a prova de "árvore idêntica".
func arvore(t *testing.T, c cofre) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, root := range []string{c.vault, c.docs} {
		err := filepath.WalkDir(root, func(p string, e fs.DirEntry, err error) error {
			if err != nil || e.IsDir() {
				return err
			}
			b, err := os.ReadFile(p)
			sum := sha256.Sum256(b)
			out[p] = hex.EncodeToString(sum[:])
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return out
}

func igual(t *testing.T, antes, depois map[string]string) {
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

func seguir(compor.Plano, func(compor.Ligacao) (compor.Plano, error)) (bool, error) { return true, nil }
func cancelar(compor.Plano, func(compor.Ligacao) (compor.Plano, error)) (bool, error) {
	return false, nil
}

func depsDe(c cofre, fsys processo.FS, confirm compor.ConfirmFunc) compor.Deps {
	return compor.Deps{
		Store:   processo.Store{Roots: processo.Roots{Vault: c.vault, Docs: c.docs}, FS: fsys},
		Confirm: confirm,
		Clock:   func() time.Time { return inicio.Add(time.Hour) },
	}
}

func ligar(c cofre, nome, composto string) (*compor.Resultado, error) {
	req := compor.Request{Nome: nome, Ligacao: compor.Ligacao{Composto: composto}}
	return compor.Run(context.Background(), req, depsDe(c, nil, seguir))
}

func lerProcesso(t *testing.T, c cofre, nome string) processo.Processo {
	t.Helper()
	p, err := processo.Ler(processo.CaminhoNoCofre(c.vault, nome))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCompor(t *testing.T) {
	t.Run("ligar b@1.0.0 em a: etapa composta, DAG com subloop, versão 1.1.0", testLigar)
	t.Run("ciclo direto: a em b depois de b em a", testCicloDireto)
	t.Run("ciclo transitivo: a->b->c->a", testCicloTransitivo)
	t.Run("ligar o processo nele mesmo é ciclo", testAutoCiclo)
	t.Run("processo inexistente: erro claro", testInexistente)
	t.Run("versão que não confere: erro claro", testVersaoErrada)
	t.Run("pedido inválido não grava nada", testInvalidos)
	t.Run("falha do FS volta a 1.0.0 em qualquer ponto", testFalhaFS)
	t.Run("cancelar deixa a árvore idêntica", testCancelar)
	t.Run("contexto cancelado não grava", testContexto)
	t.Run("depende escolhido pelo usuário", testDepende)
	t.Run("personalizar: vale o último plano", testPersonalizar)
	t.Run("artefato repetido: erro claro", testArtefatoRepetido)
	t.Run("mesmo artefato em a e b: caminho próprio no pai", testMesmoArtefato)
	t.Run("cartão cumpre o STE-PT", testCartao)
	t.Run("raízes vazias", testRaizesVazias)
	t.Run("outro dia: mantém os nomes e a criação", testOutroDia)
	t.Run("cadeia com processo ausente ou ilegível", testCadeia)
	t.Run("arquivo guardado com outro nome", testArquivoTrocado)
	t.Run("replan que falha não troca o plano", testReplanFalho)
}

func testLigar(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	antesB := arvoreDe(t, c, nomeB)
	res, err := ligar(c, nomeA, nomeB+"@1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if res.De != "1.0.0" || res.Para != "1.1.0" || res.Etapa != "e04" || res.Composto != nomeB+"@1.0.0" {
		t.Errorf("resultado %+v; esperava e04, 1.0.0 -> 1.1.0, legoloop-b@1.0.0", res)
	}
	a := lerProcesso(t, c, nomeA)
	if a.Versao != "1.1.0" || !slices.Equal(a.Historico, []string{"1.0.0"}) || len(a.Etapas) != 4 {
		t.Fatalf("processo a: versão %q, histórico %v, %d etapas", a.Versao, a.Historico, len(a.Etapas))
	}
	e := a.Etapas[3]
	// o artefato de b chega ao pai em out/<b>/<arquivo>; a origem é o caminho dentro de b (TICK-03.4b)
	quer := processoDe(nomeB, 2).Etapas[1].Produz
	quer.Caminho = "out/legoloop-b/e02.md"
	if e.ID != "e04" || e.Compoe != nomeB+"@1.0.0" || e.Peca != "" || !slices.Equal(e.Depende, []string{"e03"}) ||
		e.Produz != quer || e.Origem != "out/b/e02.md" || e.Gate != "file_exists('out/legoloop-b/e02.md')" ||
		!slices.Equal(e.Consome, []string{"out/a/e03.md"}) ||
		e.Titulo != "Executar o processo legoloop-b" {
		t.Errorf("etapa composta: %+v", e)
	}
	conferirDAG(t, c)
	igual(t, antesB, arvoreDe(t, c, nomeB))
	for _, arq := range res.Arquivos {
		if _, err := os.Stat(arq.Caminho + ".legoloop-bak"); arq.Estado != "troca" || err != nil || arq.Backup != arq.Caminho+".legoloop-bak" {
			t.Errorf("%s: estado %q, backup %q (%v); esperava troca com .legoloop-bak", arq.Tipo, arq.Estado, arq.Backup, err)
		}
	}
	intent, err := os.ReadFile(filepath.Join(c.docs, "docs", "intents", "2026-11-05_INTENT_A.md"))
	if txt := string(intent); err != nil || !strings.HasPrefix(txt, "# INTENT: legoloop-a\n\nTexto original.\n") ||
		!strings.Contains(txt, "## Composição 1.1.0") || !strings.Contains(txt, "* Processo: legoloop-b@1.0.0") {
		t.Errorf("INTENT (%v):\n%s", err, txt)
	}
}

// testOutroDia: o compor feito dias depois mantém os nomes de SPEC, PLAN, ROUTE e INTENT e a criação do loop.
func testOutroDia(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	nomes := func() []string {
		var out []string
		for p := range arvore(t, c) {
			if strings.HasPrefix(p, c.docs) && !strings.HasSuffix(p, ".legoloop-bak") {
				out = append(out, p)
			}
		}
		slices.Sort(out)
		return out
	}
	antes := nomes()
	d := depsDe(c, nil, seguir)
	d.Clock = func() time.Time { return inicio.Add(48 * time.Hour) }
	req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0"}}
	if _, err := compor.Run(context.Background(), req, d); err != nil {
		t.Fatal(err)
	}
	if depois := nomes(); !slices.Equal(antes, depois) {
		t.Errorf("os documentos mudaram de nome:\nantes  %v\ndepois %v", antes, depois)
	}
	var lp struct {
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	b, err := os.ReadFile(filepath.Join(c.vault, "loops", nomeA+".loop.json"))
	if err != nil || json.Unmarshal(b, &lp) != nil || !lp.CreatedAt.Equal(inicio) || !lp.UpdatedAt.Equal(inicio.Add(48*time.Hour)) {
		t.Errorf("loop: criado %v, atualizado %v (%v); esperava a criação mantida", lp.CreatedAt, lp.UpdatedAt, err)
	}
}

// testCadeia: um processo ausente no meio da cadeia não bloqueia; um processo ilegível bloqueia sem gravar.
func testCadeia(t *testing.T) {
	c := semear(t, nomeA, nomeB, nomeC)
	if _, err := ligar(c, nomeB, nomeC+"@1.0.0"); err != nil { // b -> c
		t.Fatal(err)
	}
	if got := lerProcesso(t, c, nomeB).Etapas[2].Compoe; got != nomeC+"@1.0.0" {
		t.Errorf("b compõe %q; esperava legoloop-c@1.0.0", got)
	}
	caminhoC := processo.CaminhoNoCofre(c.vault, nomeC)
	if err := os.WriteFile(caminhoC, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	esperarFalha(t, c, nomeA, nomeB+"@1.1.0", processo.ErrInvalido)
	if err := os.Remove(caminhoC); err != nil {
		t.Fatal(err)
	}
	if _, err := ligar(c, nomeA, nomeB+"@1.1.0"); err != nil {
		t.Fatalf("c fora do cofre não forma ciclo; o compor deveria seguir: %v", err)
	}
}

// testArquivoTrocado: o processo.json de legoloop-a guardado como legoloop-y não vale (gravaria nos nomes errados).
func testArquivoTrocado(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	b, err := os.ReadFile(processo.CaminhoNoCofre(c.vault, nomeA))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(processo.CaminhoNoCofre(c.vault, "legoloop-y"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	esperarFalha(t, c, "legoloop-y", nomeB+"@1.0.0", processo.ErrInvalido, "legoloop-y")
}

// testReplanFalho: um personalizar que não vale não troca o plano; seguir grava o plano anterior.
func testReplanFalho(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	confirm := func(pl compor.Plano, replan func(compor.Ligacao) (compor.Plano, error)) (bool, error) {
		l := pl.Ligacao
		l.Titulo, l.Depende = "Não vale", []string{"e99"}
		if _, err := replan(l); !errors.Is(err, compor.ErrPedido) {
			t.Errorf("replan: erro %v; esperava ErrPedido", err)
		}
		return true, nil
	}
	req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0"}}
	if _, err := compor.Run(context.Background(), req, depsDe(c, nil, confirm)); err != nil {
		t.Fatal(err)
	}
	if e := lerProcesso(t, c, nomeA).Etapas[3]; e.Titulo != "Executar o processo legoloop-b" || !slices.Equal(e.Depende, []string{"e03"}) {
		t.Errorf("etapa %+v; esperava o plano anterior ao replan que falhou", e)
	}
}

// arvoreDe é a árvore só dos arquivos do processo dado (no cofre e em docs/).
func arvoreDe(t *testing.T, c cofre, nome string) map[string]string {
	t.Helper()
	out := map[string]string{}
	sep := string(filepath.Separator)
	for p, h := range arvore(t, c) {
		if strings.Contains(p, sep+nome+".") || strings.Contains(p, sep+nome+sep) || strings.Contains(p, "_"+slugDe(nome)+".") {
			out[p] = h
		}
	}
	return out
}

func conferirDAG(t *testing.T, c cofre) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(c.vault, "loops", nomeA+".dag.json"))
	if err != nil {
		t.Fatal(err)
	}
	d, err := dag.Parse(b)
	if err != nil {
		t.Fatalf("o DAG gravado não é válido: %v", err)
	}
	n := d.Nodes[len(d.Nodes)-1]
	if n.ID != "e04" || n.Type != dag.TypeSubloop || n.Ref != nomeB || d.Version != "1.1.0" {
		t.Errorf("nó composto: id %s tipo %s ref %s versão %s", n.ID, n.Type, n.Ref, d.Version)
	}
	var lp struct {
		Nodes []dag.Node `json:"nodes"`
	}
	lb, err := os.ReadFile(filepath.Join(c.vault, "loops", nomeA+".loop.json"))
	if err != nil || json.Unmarshal(lb, &lp) != nil || len(lp.Nodes) != 4 || lp.Nodes[3].Type != dag.TypeSubloop {
		t.Errorf("o loop de a não tem o subloop: %v", err)
	}
}

func esperarFalha(t *testing.T, c cofre, nome, composto string, alvo error, trechos ...string) {
	t.Helper()
	antes := arvore(t, c)
	_, err := ligar(c, nome, composto)
	if !errors.Is(err, alvo) {
		t.Fatalf("compor %s %s: erro %v; esperava %v", nome, composto, err, alvo)
	}
	for _, tr := range trechos {
		if !strings.Contains(err.Error(), tr) {
			t.Errorf("a mensagem não diz %q: %v", tr, err)
		}
	}
	igual(t, antes, arvore(t, c))
}

func testCicloDireto(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	if _, err := ligar(c, nomeA, nomeB+"@1.0.0"); err != nil {
		t.Fatal(err)
	}
	esperarFalha(t, c, nomeB, nomeA+"@1.1.0", processo.ErrCiclo, nomeA, nomeB)
	esperarFalha(t, c, nomeB, nomeA+"@1.0.0", processo.ErrCiclo) // o ciclo vem antes da versão
	if b := lerProcesso(t, c, nomeB); b.Versao != "1.0.0" || len(b.Etapas) != 2 {
		t.Errorf("o processo b mudou: %q, %d etapas", b.Versao, len(b.Etapas))
	}
}

func testCicloTransitivo(t *testing.T) {
	c := semear(t, nomeA, nomeB, nomeC)
	if _, err := ligar(c, nomeC, nomeA+"@1.0.0"); err != nil { // c -> a
		t.Fatal(err)
	}
	if _, err := ligar(c, nomeB, nomeC+"@1.1.0"); err != nil { // b -> c
		t.Fatal(err)
	}
	if got := lerProcesso(t, c, nomeB).Etapas[2].Compoe; got != nomeC+"@1.1.0" {
		t.Errorf("b compõe %q; esperava a versão pedida, legoloop-c@1.1.0", got)
	}
	esperarFalha(t, c, nomeA, nomeB+"@1.1.0", processo.ErrCiclo, nomeA, nomeB, nomeC)
}

func testAutoCiclo(t *testing.T) {
	c := semear(t, nomeA)
	esperarFalha(t, c, nomeA, nomeA+"@1.0.0", processo.ErrCiclo, nomeA)
}

func testInexistente(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	esperarFalha(t, c, nomeA, "legoloop-zzz@1.0.0", compor.ErrNaoExiste, "legoloop-zzz", "não existe", "legoloop listar")
	esperarFalha(t, c, "legoloop-zzz", nomeB+"@1.0.0", compor.ErrNaoExiste, "legoloop-zzz")
	// processo da v1: só o .loop.json, sem o .processo.json
	if err := os.Remove(processo.CaminhoNoCofre(c.vault, nomeB)); err != nil {
		t.Fatal(err)
	}
	esperarFalha(t, c, nomeA, nomeB+"@1.0.0", compor.ErrNaoExiste, "v1")
}

func testVersaoErrada(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	esperarFalha(t, c, nomeA, nomeB+"@2.0.0", compor.ErrVersao, "2.0.0", "1.0.0")
	esperarFalha(t, c, nomeA, nomeB+"@1.0", compor.ErrVersao, "1.0.0")
}

func testInvalidos(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	esperarFalha(t, c, "Legoloop-A", nomeB+"@1.0.0", compor.ErrNome)
	esperarFalha(t, c, nomeA, nomeB, compor.ErrPedido, "@")
	esperarFalha(t, c, nomeA, nomeB+"@", compor.ErrPedido, "@<versão>")
	esperarFalha(t, c, nomeA, "b@1.0.0", compor.ErrPedido)
	esperarFalha(t, c, nomeA, "", compor.ErrPedido)
	esperarFalha(t, c, nomeA, "../legoloop-b@1.0.0", compor.ErrPedido)
	if _, err := ligar(c, nomeA, nomeB+"@1.0.0"); err != nil {
		t.Fatal(err)
	}
	esperarFalha(t, c, nomeA, nomeB+"@1.0.0", compor.ErrPedido, "já usa", "e04") // já compõe b
	req := compor.Request{Nome: nomeB, Ligacao: compor.Ligacao{Composto: nomeC + "@1.0.0"}}
	gravar(t, c, processoDe(nomeC, 2))
	antes := arvore(t, c)
	req.Ligacao.Depende = []string{"e09"}
	if _, err := compor.Run(context.Background(), req, depsDe(c, nil, seguir)); !errors.Is(err, compor.ErrPedido) || !strings.Contains(err.Error(), "e09") {
		t.Errorf("depende inexistente: erro %v; esperava ErrPedido citando e09", err)
	}
	igual(t, antes, arvore(t, c))
}

var errInjetado = errors.New("falha injetada no FS")

// falhaFS falha na n-ésima chamada da operação op e delega todo o resto ao FS real.
type falhaFS struct {
	processo.FS
	op       string
	n        int
	chamadas map[string]int
	disparou bool
}

func (f *falhaFS) tick(op string) error {
	f.chamadas[op]++
	if op == f.op && f.chamadas[op] == f.n {
		f.disparou = true
		return errInjetado
	}
	return nil
}

func (f *falhaFS) MkdirAll(p string, perm os.FileMode) error {
	if err := f.tick("MkdirAll"); err != nil {
		return err
	}
	return f.FS.MkdirAll(p, perm)
}

func (f *falhaFS) WriteFile(n string, d []byte, perm os.FileMode) error {
	if err := f.tick("WriteFile"); err != nil {
		return err
	}
	return f.FS.WriteFile(n, d, perm)
}

func (f *falhaFS) Rename(a, b string) error {
	if err := f.tick("Rename"); err != nil {
		return err
	}
	return f.FS.Rename(a, b)
}

func (f *falhaFS) RemoveAll(p string) error {
	if err := f.tick("RemoveAll"); err != nil {
		return err
	}
	return f.FS.RemoveAll(p)
}

func testFalhaFS(t *testing.T) {
	req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0"}}
	for _, op := range []string{"MkdirAll", "WriteFile", "Rename", "RemoveAll"} {
		cobertas := 0
		for n := 1; ; n++ {
			c := semear(t, nomeA, nomeB)
			antes := arvore(t, c)
			f := &falhaFS{FS: processo.OSFS(), op: op, n: n, chamadas: map[string]int{}}
			_, err := compor.Run(context.Background(), req, depsDe(c, f, seguir))
			if !f.disparou {
				if err != nil {
					t.Fatalf("%s #%d: sem falha injetada, o compor não deveria falhar: %v", op, n, err)
				}
				break
			}
			cobertas++
			if !errors.Is(err, errInjetado) {
				t.Errorf("%s #%d: erro %v; esperava a falha do FS", op, n, err)
			}
			igual(t, antes, arvore(t, c))
			if a := lerProcesso(t, c, nomeA); a.Versao != "1.0.0" || len(a.Etapas) != 3 {
				t.Errorf("%s #%d: o processo não voltou a 1.0.0: %q, %d etapas", op, n, a.Versao, len(a.Etapas))
			}
		}
		if cobertas < 8 {
			t.Errorf("%s: só %d pontos de falha provados; esperava ao menos 8", op, cobertas)
		}
	}
}

func testCancelar(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	antes := arvore(t, c)
	req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0"}}
	if _, err := compor.Run(context.Background(), req, depsDe(c, nil, cancelar)); !errors.Is(err, compor.ErrCancelado) {
		t.Fatalf("erro %v; esperava ErrCancelado", err)
	}
	igual(t, antes, arvore(t, c))
	errConfirm := errors.New("terminal fechado")
	falha := func(compor.Plano, func(compor.Ligacao) (compor.Plano, error)) (bool, error) { return true, errConfirm }
	if _, err := compor.Run(context.Background(), req, depsDe(c, nil, falha)); !errors.Is(err, errConfirm) {
		t.Fatalf("erro %v; esperava o erro da confirmação", err)
	}
	igual(t, antes, arvore(t, c))
}

func testContexto(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	antes := arvore(t, c)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0"}}
	if _, err := compor.Run(ctx, req, depsDe(c, nil, seguir)); !errors.Is(err, context.Canceled) {
		t.Fatalf("erro %v; esperava context.Canceled", err)
	}
	igual(t, antes, arvore(t, c))
}

func testDepende(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0", Depende: []string{" e01 ", "e01"}, Titulo: "Rodar b"}}
	if _, err := compor.Run(context.Background(), req, depsDe(c, nil, seguir)); err != nil {
		t.Fatal(err)
	}
	e := lerProcesso(t, c, nomeA).Etapas[3]
	if !slices.Equal(e.Depende, []string{"e01"}) || e.Titulo != "Rodar b" || !slices.Equal(e.Consome, []string{"out/a/e01.md"}) {
		t.Errorf("etapa: %+v; esperava depender só de e01, com o título pedido", e)
	}
}

func testPersonalizar(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	var replanErr error
	confirm := func(pl compor.Plano, replan func(compor.Ligacao) (compor.Plano, error)) (bool, error) {
		l := pl.Ligacao
		l.Depende = []string{"e99"}
		if _, replanErr = replan(l); replanErr == nil {
			t.Error("replan com etapa inexistente deveria falhar")
		}
		l.Depende, l.Titulo = []string{"e02"}, "Título novo"
		_, err := replan(l)
		return true, err
	}
	req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0"}}
	if _, err := compor.Run(context.Background(), req, depsDe(c, nil, confirm)); err != nil {
		t.Fatal(err)
	}
	e := lerProcesso(t, c, nomeA).Etapas[3]
	if e.Titulo != "Título novo" || !slices.Equal(e.Depende, []string{"e02"}) {
		t.Errorf("etapa %+v; esperava o último plano (e02, Título novo)", e)
	}
}

// testArtefatoRepetido: a colisão agora é com o caminho do artefato no pai (out/<b>/<arquivo>). O caminho
// dentro de b (out/b/e02.md) não colide mais: isso é o testMesmoArtefato.
func testArtefatoRepetido(t *testing.T) {
	c := semear(t, nomeB)
	a := processoDe(nomeA, 3)
	a.Etapas[0].Produz.Caminho, a.Etapas[0].Gate = "out/legoloop-b/e02.md", "file_exists('out/legoloop-b/e02.md')"
	a.Etapas[1].Consome = []string{"out/legoloop-b/e02.md"}
	gravar(t, c, a)
	esperarFalha(t, c, nomeA, nomeB+"@1.0.0", compor.ErrPedido, "out/legoloop-b/e02.md", "e01", nomeB)
}

// comoFabricar devolve o processo como o fabricar grava: construir e verificar em out/01_construir.md e
// out/02_verificar.md, com o gate da biblioteca. Dois processos assim gravam os mesmos caminhos.
func comoFabricar(t *testing.T, nome string) processo.Processo {
	t.Helper()
	p := processo.Processo{Nome: nome, Versao: "1.0.0", Natureza: "software", Objetivo: "Objetivo de " + nome + ".", CriadoEm: inicio}
	for i, tipo := range []string{"construir", "verificar"} {
		art := processo.Artefato{Caminho: "out/0" + string(rune('1'+i)) + "_" + tipo + ".md", Tipo: "markdown", Promete: true}
		g, err := processo.GateParaArtefato(art)
		if err != nil {
			t.Fatal(err)
		}
		e := processo.Etapa{ID: "e0" + string(rune('1'+i)), Titulo: "Etapa " + tipo, Tipo: tipo, Peca: "aidd-tdd", Produz: art, Gate: g}
		if i > 0 {
			e.Depende, e.Consome = []string{"e01"}, []string{"out/01_construir.md"}
		}
		p.Etapas = append(p.Etapas, e)
	}
	return p
}

// testMesmoArtefato: a e b gravam out/02_verificar.md. O compor liga b em a: o artefato de b chega ao pai em
// out/legoloop-b/02_verificar.md, o gate confere esse caminho e o DAG vale. A etapa e02 de a não muda.
func testMesmoArtefato(t *testing.T) {
	root := t.TempDir()
	c := cofre{vault: filepath.Join(root, "vault"), docs: filepath.Join(root, "proj")}
	gravar(t, c, comoFabricar(t, nomeA))
	pb := comoFabricar(t, nomeB)
	pb.Etapas[0].Produz.Promete = false // o artefato da etapa nova vem da etapa final de b, não da primeira
	gravar(t, c, pb)
	if _, err := ligar(c, nomeA, nomeB+"@1.0.0"); err != nil {
		t.Fatalf("dois processos com os mesmos artefatos deveriam se ligar: %v", err)
	}
	a := lerProcesso(t, c, nomeA)
	e, velha := a.Etapas[2], comoFabricar(t, nomeA).Etapas[1]
	quer, _ := processo.GateParaArtefato(processo.Artefato{Caminho: "out/legoloop-b/02_verificar.md", Tipo: "markdown"})
	if e.Produz.Caminho != "out/legoloop-b/02_verificar.md" || e.Origem != "out/02_verificar.md" || e.Gate != quer ||
		!e.Produz.Promete || !slices.Equal(e.Consome, []string{"out/02_verificar.md"}) {
		t.Errorf("etapa composta: %+v", e)
	}
	if strings.Contains(e.Gate, "'out/02_verificar.md'") || a.Etapas[1].Produz != velha.Produz || a.Etapas[1].Gate != velha.Gate {
		t.Errorf("o gate da etapa nova não pode conferir o arquivo de a, e a e02 de a não muda: %+v", a.Etapas)
	}
	b, err := os.ReadFile(filepath.Join(c.vault, "loops", nomeA+".dag.json"))
	if err != nil {
		t.Fatal(err)
	}
	d, err := dag.Parse(b)
	if err != nil {
		t.Fatalf("o DAG gravado não é válido: %v", err)
	}
	if o := d.Nodes[2].Output; o.Artifact != e.Produz.Caminho || o.SubArtifact != "out/02_verificar.md" || o.Gate != quer {
		t.Errorf("nó do subloop: %+v", o)
	}
}

func testCartao(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	var texto, diff string
	confirm := func(pl compor.Plano, _ func(compor.Ligacao) (compor.Plano, error)) (bool, error) {
		card := pl.Cartao()
		if vs := ste.LintCartao(card); len(vs) != 0 {
			t.Errorf("o cartão viola o STE-PT: %v", vs)
		}
		if len(card.Parametros) != 2 || card.Parametros[0].Nome != "titulo" || card.Parametros[1].Nome != "depende" ||
			card.Parametros[1].Valor != "e03" {
			t.Errorf("parâmetros do personalizar: %+v", card.Parametros)
		}
		texto, diff = card.Render(), pl.Diff()
		return false, nil
	}
	req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0"}}
	_, _ = compor.Run(context.Background(), req, depsDe(c, nil, confirm))
	for _, quer := range []string{".legoloop-bak", "1.0.0", "1.1.0", nomeB + "@1.0.0", "e04", "[1] seguir", "Efeito: troca "} {
		if !strings.Contains(texto, quer) {
			t.Errorf("o cartão não diz %q:\n%s", quer, texto)
		}
	}
	if i, j := strings.Index(texto, "AVISO:"), strings.Index(texto, "O QUE VAI ACONTECER"); i < 0 || i > j {
		t.Errorf("o aviso de risco deve vir antes dos passos:\n%s", texto)
	}
	for _, quer := range []string{"COMPOR", "DAG", "processo", "e03", "out/b/e02.md", "out/legoloop-b/e02.md"} {
		if !strings.Contains(diff, quer) {
			t.Errorf("o diff não diz %q:\n%s", quer, diff)
		}
	}
	if vs := ste.Default.Violations(); len(vs) != 0 {
		t.Errorf("mensagens do catálogo fora do perfil: %v", vs)
	}
}

func testRaizesVazias(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	for rotulo, r := range map[string]processo.Roots{"sem docs": {Vault: c.vault}, "sem cofre": {Docs: c.docs}} {
		d := compor.Deps{Store: processo.Store{Roots: r}, Confirm: seguir}
		req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0"}}
		if _, err := compor.Run(context.Background(), req, d); !errors.Is(err, processo.ErrPlano) {
			t.Errorf("%s: erro %v; esperava ErrPlano", rotulo, err)
		}
	}
}

func TestCompor_RegeraDAGIso(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	req := compor.Request{Nome: nomeA, Ligacao: compor.Ligacao{Composto: nomeB + "@1.0.0"}}
	res, err := compor.Run(context.Background(), req, depsDe(c, nil, seguir))
	if err != nil {
		t.Fatalf("compor.Run: %v", err)
	}
	var achouDAGIso bool
	for _, a := range res.Arquivos {
		if strings.Contains(a.Caminho, filepath.Join("docs", "dags")) {
			achouDAGIso = true
			if b, err := os.ReadFile(a.Caminho); err != nil || len(b) == 0 {
				t.Fatalf("arquivo DAG ISO em %s não pode estar vazio", a.Caminho)
			}
		}
	}
	if !achouDAGIso {
		t.Errorf("compor deve regerar e persistir docs/dags: %+v", res.Arquivos)
	}
}
