package refinar_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"legoloop/internal/paths"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/processo"
	"legoloop/internal/refinar"
	"legoloop/internal/ste"
)

var (
	inicio = time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	nome   = "legoloop-x"
)

func etapa(id, tipo, peca string, depende ...string) processo.Etapa {
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

func processoDeTeste() processo.Processo {
	return processo.Processo{
		Nome: nome, Versao: "1.0.0", Natureza: "software", Objetivo: "Extrair o áudio de um vídeo.", CriadoEm: inicio,
		Etapas: []processo.Etapa{
			etapa("e01", "clarificar", "aidd-grill"),
			etapa("e02", "construir", "aidd-tdd", "e01"),
			etapa("e03", "verificar", "aidd-qa", "e02"),
		},
	}
}

type cofre struct{ vault, docs string }

// slugDe é o nome do processo nos documentos: legoloop-extrair-audio vira EXTRAIR_AUDIO.
func slugDe(nome string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimPrefix(nome, "legoloop-"), "-", "_"))
}

// semear grava o processo como o fabricar grava: skill, documento do Studio, DAG e processo.json no cofre; INTENT, SPEC,
// PLAN e ROUTE em docs/. comIntent=false deixa o INTENT de fora.
func semear(t *testing.T, p processo.Processo, data time.Time, comIntent bool) cofre {
	t.Helper()
	root := t.TempDir()
	c := cofre{vault: filepath.Join(root, "vault"), docs: filepath.Join(root, "proj")}
	pr, err := processo.Projetar(p, data)
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
		{Path: filepath.Join(c.vault, paths.DirAntigo, p.Nome+paths.SufixoAntigo), Data: pr.Loop},
		{Path: filepath.Join(c.vault, paths.DirAntigo, p.Nome+".dag.json"), Data: pr.DAG},
		{Path: processo.CaminhoNoCofre(c.vault, p.Nome), Data: pj},
		{Path: doc("specs", pr.SpecNome), Data: []byte(pr.Spec)},
		{Path: doc("plans", pr.PlanNome), Data: []byte(pr.Plan)},
		{Path: doc("routes", pr.RouteNome), Data: pr.Route},
	}
	if comIntent {
		intent := data.Format("2006-01-02") + "_INTENT_" + slugDe(p.Nome) + ".md"
		files = append(files, processo.File{Path: doc("intents", intent), Data: []byte("# INTENT: x\n\nTexto original da entrevista.\n")})
	}
	st := processo.Store{Roots: processo.Roots{Vault: c.vault, Docs: c.docs}}
	if err := st.Write(context.Background(), files, processo.WriteOpts{}); err != nil {
		t.Fatal(err)
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

func depsDe(c cofre, fsys processo.FS, confirm refinar.ConfirmFunc, clock time.Time) refinar.Deps {
	return refinar.Deps{
		Store:   processo.Store{Roots: processo.Roots{Vault: c.vault, Docs: c.docs}, FS: fsys},
		Confirm: confirm,
		Clock:   func() time.Time { return clock },
	}
}

func seguir(refinar.Plano, func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
	return true, nil
}
func cancelar(refinar.Plano, func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
	return false, nil
}

var mudanca = refinar.Mudanca{Etapa: "e02", Peca: "aidd-qa", Casos: []string{"verificacao-de-qualidade"}}

func lerProcesso(t *testing.T, c cofre) processo.Processo {
	t.Helper()
	p, err := processo.Ler(processo.CaminhoNoCofre(c.vault, nome))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func ler(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func hashDe(s string) string { sum := sha256.Sum256([]byte(s)); return hex.EncodeToString(sum[:]) }

func TestRefinar(t *testing.T) {
	t.Run("seguir: diff dos sete, versão 1.1.0, backup e histórico", testSeguir)
	t.Run("cancelar deixa a árvore idêntica", testCancelar)
	t.Run("falha do FS volta a 1.0.0 em qualquer ponto", testFalhaFS)
	t.Run("pedido inválido não grava nada", testInvalidos)
	t.Run("docs de outro dia: reaproveita os nomes", testDataAntiga)
	t.Run("duas rodadas: 1.2.0 e histórico em ordem", testDuasRodadas)
	t.Run("sem INTENT: cria o INTENT", testSemIntent)
	t.Run("personalizar: vale o último plano", testPersonalizar)
	t.Run("peça fora da matriz vira aviso", testPecaForaDaMatriz)
	t.Run("cartão cumpre o STE-PT", testCartao)
	t.Run("contexto cancelado não grava", testContexto)
	t.Run("processo da v1 não se refina", testV1)
	t.Run("arquivo guardado com outro nome", testArquivoTrocado)
	t.Run("casos sem espaço e sem repetição", testCasosLimpos)
	t.Run("só os casos: PLAN e ROUTE ficam iguais", testSoCasos)
	t.Run("duas SPEC: vale a mais nova", testDuasSpecs)
	t.Run("nome com hífen", testNomeComHifen)
	t.Run("raízes vazias", testRaizesVazias)
}

func testRaizesVazias(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	for rotulo, r := range map[string]processo.Roots{"sem docs": {Vault: c.vault}, "sem cofre": {Docs: c.docs}} {
		chamou := false
		confirm := func(refinar.Plano, func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
			chamou = true
			return true, nil
		}
		d := refinar.Deps{Store: processo.Store{Roots: r}, Confirm: confirm, Clock: func() time.Time { return inicio }}
		if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, d); !errors.Is(err, processo.ErrPlano) || chamou {
			t.Errorf("%s: erro %v (confirmou=%v); esperava ErrPlano antes de pedir a confirmação", rotulo, err, chamou)
		}
	}
}

func testV1(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	if err := os.Remove(processo.CaminhoNoCofre(c.vault, nome)); err != nil {
		t.Fatal(err)
	}
	antes := arvore(t, c)
	_, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, seguir, inicio))
	if !errors.Is(err, refinar.ErrNaoExiste) || !strings.Contains(err.Error(), "da v1") {
		t.Fatalf("erro %v; esperava ErrNaoExiste com a explicação da v1", err)
	}
	igual(t, antes, arvore(t, c))
}

// testArquivoTrocado: o processo.json de legoloop-x guardado como legoloop-y não vale (gravaria nos nomes errados).
func testArquivoTrocado(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	copia := ler(t, processo.CaminhoNoCofre(c.vault, nome))
	if err := os.WriteFile(processo.CaminhoNoCofre(c.vault, "legoloop-y"), []byte(copia), 0o644); err != nil {
		t.Fatal(err)
	}
	antes := arvore(t, c)
	_, err := refinar.Run(context.Background(), refinar.Request{Nome: "legoloop-y", Mudanca: mudanca}, depsDe(c, nil, seguir, inicio))
	if !errors.Is(err, processo.ErrInvalido) {
		t.Fatalf("erro %v; esperava ErrInvalido", err)
	}
	igual(t, antes, arvore(t, c))
}

func testCasosLimpos(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	m := refinar.Mudanca{Etapa: "e02", Casos: []string{" tdd ", "tdd", "", "qualidade", "tdd"}}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: m}, depsDe(c, nil, seguir, inicio)); err != nil {
		t.Fatal(err)
	}
	if got := lerProcesso(t, c).Etapas[1].Casos; !slices.Equal(got, []string{"tdd", "qualidade"}) {
		t.Errorf("casos %q; esperava [tdd qualidade]", got)
	}
}

// testSoCasos: mudar só os casos sobe a versão e troca SPEC, skill, documento do Studio, DAG e processo; PLAN e ROUTE não
// mudam, ficam fora da gravação e não ganham .legoloop-bak.
func testSoCasos(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	m := refinar.Mudanca{Etapa: "e02", Casos: []string{"tdd"}}
	res, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: m}, depsDe(c, nil, seguir, inicio))
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range res.Arquivos {
		quer := "troca"
		if a.Tipo == "PLAN" {
			quer = "igual"
		}
		if a.Estado != quer {
			t.Errorf("%s: estado %q; esperava %q", a.Tipo, a.Estado, quer)
		}
		if _, err := os.Stat(a.Caminho + ".legoloop-bak"); (err == nil) != (quer == "troca") {
			t.Errorf("%s: backup presente=%v com estado %q", a.Tipo, err == nil, a.Estado)
		}
	}
}

func testDuasSpecs(t *testing.T) {
	c := semear(t, processoDeTeste(), time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC), true)
	velha := achar(t, c.docs, "specs")
	nova := filepath.Join(filepath.Dir(velha), "2026-10-15_SPEC_X.md")
	if err := os.WriteFile(nova, []byte(ler(t, velha)), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, seguir, inicio))
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range res.Arquivos {
		if a.Tipo == "SPEC" && a.Caminho != nova {
			t.Errorf("a SPEC trocada é %s; esperava a mais nova, %s", a.Caminho, nova)
		}
	}
}

// testNomeComHifen: legoloop-extrair-audio guarda os documentos como EXTRAIR_AUDIO; o refinar acha o INTENT e a SPEC.
func testNomeComHifen(t *testing.T) {
	p := processoDeTeste()
	p.Nome = "legoloop-extrair-audio"
	c := semear(t, p, inicio, true)
	res, err := refinar.Run(context.Background(), refinar.Request{Nome: p.Nome, Mudanca: mudanca}, depsDe(c, nil, seguir, inicio))
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range res.Arquivos {
		if a.Estado != "troca" {
			t.Errorf("%s: estado %q; esperava troca", a.Tipo, a.Estado)
		}
	}
	if txt := ler(t, achar(t, c.docs, "intents")); !strings.Contains(txt, "Texto original da entrevista.") {
		t.Errorf("o INTENT perdeu o texto original:\n%s", txt)
	}
}

func testSeguir(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	antes := arvore(t, c)
	var diff string
	confirm := func(pl refinar.Plano, r func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
		diff = pl.Diff()
		return seguir(pl, r)
	}
	res, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, confirm, inicio.Add(24*time.Hour)))
	if err != nil {
		t.Fatal(err)
	}
	for _, tipo := range []string{"INTENT", "SPEC", "PLAN", "ROUTE", "DAG", "loop", "skill"} {
		if !strings.Contains(diff, tipo) {
			t.Errorf("o diff não cita %s:\n%s", tipo, diff)
		}
	}
	if res.De != "1.0.0" || res.Para != "1.1.0" {
		t.Errorf("versão %s -> %s; esperava 1.0.0 -> 1.1.0", res.De, res.Para)
	}
	p := lerProcesso(t, c)
	if p.Versao != "1.1.0" || !slices.Equal(p.Historico, []string{"1.0.0"}) {
		t.Errorf("versão %q e histórico %v", p.Versao, p.Historico)
	}
	if e := p.Etapas[1]; e.Peca != "aidd-qa" || !slices.Equal(e.Casos, []string{"verificacao-de-qualidade"}) {
		t.Errorf("a etapa e02 não ficou com a peça e os casos pedidos: %+v", e)
	}
	if p.Etapas[0].Peca != "aidd-grill" || p.Etapas[2].Peca != "aidd-qa" {
		t.Errorf("as outras etapas mudaram: %+v", p.Etapas)
	}
	conferirBackups(t, antes, res)
	conferirProjecoes(t, c)
	conferirDiff(t, diff, res)
}

// conferirDiff confere a contagem de linhas de cada arquivo, a amostra limitada e as linhas cortadas.
func conferirDiff(t *testing.T, diff string, res *refinar.Resultado) {
	t.Helper()
	quer := map[string][2]int{"INTENT": {7, 0}, "SPEC": {3, 3}, "PLAN": {1, 1}, "ROUTE": {2, 2}, "DAG": {2, 2}, "skill": {2, 2}}
	for _, a := range res.Arquivos {
		if q, ok := quer[a.Tipo]; ok && (a.Mais != q[0] || a.Menos != q[1]) {
			t.Errorf("%s: +%d -%d; esperava +%d -%d", a.Tipo, a.Mais, a.Menos, q[0], q[1])
		}
		if a.Tipo == "INTENT" && (len(a.Linhas) != 5 || !strings.Contains(a.Linhas[4], "mais 3 linhas")) {
			t.Errorf("a amostra do INTENT deve ter 4 linhas e o aviso de mais 3: %q", a.Linhas)
		}
	}
	for _, l := range strings.Split(diff, "\n") {
		if (strings.HasPrefix(l, "  + ") || strings.HasPrefix(l, "  - ")) && len([]rune(l)) > 104 {
			t.Errorf("linha do diff sem corte (%d): %s", len([]rune(l)), l)
		}
	}
	if !strings.Contains(diff, "ARQUIVO  ESTADO MUDANÇA   CAMINHO\n") || !strings.Contains(diff, "INTENT   troca  +7 -0     ") {
		t.Errorf("as colunas do diff não estão alinhadas:\n%s", diff)
	}
}

// conferirBackups: cada arquivo trocado tem o .legoloop-bak com o conteúdo exato da versão anterior.
func conferirBackups(t *testing.T, antes map[string]string, res *refinar.Resultado) {
	t.Helper()
	trocados := 0
	for _, a := range res.Arquivos {
		if a.Estado != "troca" {
			continue
		}
		trocados++
		if hashDe(ler(t, a.Caminho+".legoloop-bak")) != antes[a.Caminho] {
			t.Errorf("o backup de %s não guarda a versão anterior", a.Caminho)
		}
		if a.Backup != a.Caminho+".legoloop-bak" {
			t.Errorf("Backup %q de %s", a.Backup, a.Caminho)
		}
	}
	if trocados != 8 {
		t.Errorf("trocou %d arquivos; esperava 8 (7 projeções e o processo.json)", trocados)
	}
}

func conferirProjecoes(t *testing.T, c cofre) {
	t.Helper()
	skill := ler(t, filepath.Join(c.vault, "skills", nome, "SKILL.md"))
	studioJSON := ler(t, filepath.Join(c.vault, paths.DirAntigo, nome+paths.SufixoAntigo))
	dagJSON := ler(t, filepath.Join(c.vault, paths.DirAntigo, nome+".dag.json"))
	for rotulo, quer := range map[string]struct{ texto, termo string }{
		"skill traz a peça nova":   {skill, "aidd-qa"},
		"skill traz a versão nova": {skill, "1.1.0"},
		"studio traz a versão":     {studioJSON, `"version": "1.1.0"`},
		"dag liga a peça nova":     {dagJSON, "skills/aidd-qa"},
	} {
		if !strings.Contains(quer.texto, quer.termo) {
			t.Errorf("%s: falta %q", rotulo, quer.termo)
		}
	}
	if !strings.Contains(studioJSON, `"created_at": "2026-11-05T10:00:00Z"`) || !strings.Contains(studioJSON, `"updated_at": "2026-11-06T10:00:00Z"`) {
		t.Errorf("o documento deve manter created_at e levar a data do refinar em updated_at:\n%s", studioJSON)
	}
	if strings.Contains(dagJSON, "skills/aidd-tdd") {
		t.Error("o DAG ainda cita a peça antiga")
	}
	intent := ler(t, achar(t, c.docs, "intents"))
	if !strings.Contains(intent, "Texto original da entrevista.") || !strings.Contains(intent, "Refinamento 1.1.0") {
		t.Errorf("o INTENT perdeu o texto original ou não registra o refinamento:\n%s", intent)
	}
	for _, f := range arvoreNomes(t, c) {
		if strings.HasSuffix(f, ".legoloop-tmp") {
			t.Errorf("sobrou temporário: %s", f)
		}
	}
}

func arvoreNomes(t *testing.T, c cofre) []string {
	t.Helper()
	var out []string
	for p := range arvore(t, c) {
		out = append(out, p)
	}
	slices.Sort(out)
	return out
}

// achar devolve o único arquivo atual (sem .legoloop-bak) de docs/<sub>.
func achar(t *testing.T, docs, sub string) string {
	t.Helper()
	es, err := os.ReadDir(filepath.Join(docs, "docs", sub))
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range es {
		if !strings.Contains(e.Name(), ".legoloop-") {
			out = append(out, filepath.Join(docs, "docs", sub, e.Name()))
		}
	}
	if len(out) != 1 {
		t.Fatalf("docs/%s tem %d arquivos atuais: %v", sub, len(out), out)
	}
	return out[0]
}

func testCancelar(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	antes := arvore(t, c)
	_, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, cancelar, inicio))
	if !errors.Is(err, refinar.ErrCancelado) {
		t.Fatalf("erro %v; esperava ErrCancelado", err)
	}
	igual(t, antes, arvore(t, c))
	if p := lerProcesso(t, c); p.Versao != "1.0.0" || len(p.Historico) != 0 {
		t.Errorf("o cancelamento mexeu no processo: %q %v", p.Versao, p.Historico)
	}
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
	for _, op := range []string{"MkdirAll", "WriteFile", "Rename", "RemoveAll"} {
		cobertas := 0
		for n := 1; ; n++ {
			c := semear(t, processoDeTeste(), inicio, true)
			antes := arvore(t, c)
			f := &falhaFS{FS: processo.OSFS(), op: op, n: n, chamadas: map[string]int{}}
			_, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, f, seguir, inicio))
			if !f.disparou {
				if err != nil {
					t.Fatalf("%s #%d: sem falha injetada, o refinar não deveria falhar: %v", op, n, err)
				}
				break
			}
			cobertas++
			if !errors.Is(err, errInjetado) {
				t.Errorf("%s #%d: erro %v; esperava a falha do FS", op, n, err)
			}
			igual(t, antes, arvore(t, c))
			if p := lerProcesso(t, c); p.Versao != "1.0.0" || len(p.Historico) != 0 {
				t.Errorf("%s #%d: o processo não voltou a 1.0.0: %q %v", op, n, p.Versao, p.Historico)
			}
		}
		if cobertas < 8 {
			t.Errorf("%s: só %d pontos de falha provados; esperava ao menos 8", op, cobertas)
		}
	}
}

func testInvalidos(t *testing.T) {
	composto := processoDeTeste()
	composto.Etapas[2].Peca, composto.Etapas[2].Compoe = "", "legoloop-a@1.0.0"
	comCasos := processoDeTeste()
	comCasos.Etapas[1].Casos = []string{"tdd"}
	casos := []struct {
		rotulo string
		p      processo.Processo
		nome   string
		m      refinar.Mudanca
		quer   error
		msg    string // trecho que só a regra testada produz
	}{
		{"processo ausente", processoDeTeste(), "legoloop-nao-existe", mudanca, refinar.ErrNaoExiste, "não existe no cofre"},
		{"nome que foge do cofre", processoDeTeste(), "legoloop-../x", mudanca, refinar.ErrNome, "não é de processo"},
		{"nome sem prefixo", processoDeTeste(), "x", mudanca, refinar.ErrNome, "não é de processo"},
		{"etapa inexistente", processoDeTeste(), nome, refinar.Mudanca{Etapa: "e99", Peca: "aidd-qa", Casos: []string{"a"}}, refinar.ErrMudanca, "não tem a etapa e99. As etapas são: e01, e02, e03"},
		{"sem mudança", comCasos, nome, refinar.Mudanca{Etapa: "e02"}, refinar.ErrMudanca, "Diga o que mudar"},
		{"sem caso", processoDeTeste(), nome, refinar.Mudanca{Etapa: "e02", Peca: "aidd-qa"}, refinar.ErrMudanca, "precisa de ao menos um caso"},
		{"etapa composta", composto, nome, refinar.Mudanca{Etapa: "e03", Titulo: "Outro", Casos: []string{"a"}}, refinar.ErrMudanca, "compõe outro processo"},
		{"peça com nome inseguro", processoDeTeste(), nome, refinar.Mudanca{Etapa: "e02", Peca: "~/x/SKILL.md", Casos: []string{"a"}}, processo.ErrInvalido, "não é seguro"},
	}
	for _, tc := range casos {
		t.Run(tc.rotulo, func(t *testing.T) {
			c := semear(t, tc.p, inicio, true)
			antes := arvore(t, c)
			chamou := false
			confirm := func(refinar.Plano, func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
				chamou = true
				return true, nil
			}
			_, err := refinar.Run(context.Background(), refinar.Request{Nome: tc.nome, Mudanca: tc.m}, depsDe(c, nil, confirm, inicio))
			if !errors.Is(err, tc.quer) || !strings.Contains(err.Error(), tc.msg) {
				t.Fatalf("erro %v; esperava %v com %q", err, tc.quer, tc.msg)
			}
			if chamou {
				t.Error("pediu confirmação de um pedido inválido")
			}
			igual(t, antes, arvore(t, c))
		})
	}
}

func testDataAntiga(t *testing.T) {
	antiga := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	c := semear(t, processoDeTeste(), antiga, true)
	nomes := func() []string {
		var out []string
		for _, sub := range []string{"intents", "specs", "plans", "routes"} {
			es, _ := os.ReadDir(filepath.Join(c.docs, "docs", sub))
			for _, e := range es {
				if !strings.Contains(e.Name(), ".legoloop-") {
					out = append(out, sub+"/"+e.Name())
				}
			}
		}
		return out
	}
	antes := nomes()
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, seguir, inicio.Add(24*time.Hour))); err != nil {
		t.Fatal(err)
	}
	if depois := nomes(); !slices.Equal(antes, depois) {
		t.Errorf("os documentos mudaram de nome:\nantes  %v\ndepois %v", antes, depois)
	}
	route := ler(t, achar(t, c.docs, "routes"))
	if !strings.Contains(route, "2026-10-01_SPEC_X.md") || strings.Contains(route, "2026-11-06_SPEC_X.md") {
		t.Error("o ROUTE tem de citar a SPEC que existe, com o nome antigo")
	}
	if studioJSON := ler(t, filepath.Join(c.vault, paths.DirAntigo, nome+paths.SufixoAntigo)); !strings.Contains(studioJSON, "2026-11-06") {
		t.Error("o documento deve levar a data do refinar em updated_at")
	}
}

func testDuasRodadas(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	for i := range 2 {
		m := refinar.Mudanca{Etapa: "e03", Titulo: fmt.Sprintf("Verificar rodada %d", i), Casos: []string{"verificacao-de-qualidade"}}
		if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: m}, depsDe(c, nil, seguir, inicio)); err != nil {
			t.Fatal(err)
		}
	}
	p := lerProcesso(t, c)
	if p.Versao != "1.2.0" || !slices.Equal(p.Historico, []string{"1.0.0", "1.1.0"}) {
		t.Errorf("versão %q e histórico %v", p.Versao, p.Historico)
	}
	if p.Etapas[2].Titulo != "Verificar rodada 1" || p.Etapas[2].Peca != "aidd-qa" {
		t.Errorf("a etapa e03: %+v", p.Etapas[2])
	}
}

func testSemIntent(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, false)
	res, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, seguir, inicio))
	if err != nil {
		t.Fatal(err)
	}
	var intent *refinar.Arquivo
	for i, a := range res.Arquivos {
		if a.Tipo == "INTENT" {
			intent = &res.Arquivos[i]
		}
	}
	if intent == nil || intent.Estado != "novo" || filepath.Base(intent.Caminho) != "2026-11-05_INTENT_X.md" {
		t.Fatalf("o INTENT deveria ser novo: %+v", intent)
	}
	if txt := ler(t, intent.Caminho); !strings.Contains(txt, "Refinamento 1.1.0") || !strings.Contains(txt, "Extrair o áudio") {
		t.Errorf("INTENT novo sem o refinamento ou o objetivo:\n%s", txt)
	}
}

func testPersonalizar(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	confirm := func(pl refinar.Plano, replan func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
		m := pl.Mudanca
		m.Casos = []string{"outro-caso"}
		novo, err := replan(m)
		if err != nil || !slices.Equal(novo.Mudanca.Casos, []string{"outro-caso"}) {
			t.Errorf("replan: %v %+v", err, novo.Mudanca)
		}
		return true, nil
	}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, confirm, inicio)); err != nil {
		t.Fatal(err)
	}
	if e := lerProcesso(t, c).Etapas[1]; !slices.Equal(e.Casos, []string{"outro-caso"}) {
		t.Errorf("valeu o plano antigo: %v", e.Casos)
	}
}

func testPecaForaDaMatriz(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	reg := `{"matrix_version":"2.0.0","entradas":[{"ref":"~/.claude/skills/aidd-qa/SKILL.md","nome":"aidd-qa","tipo":"skill","hash":"h1"}]}`
	if err := os.WriteFile(filepath.Join(c.vault, "registry_matrix.json"), []byte(reg), 0o644); err != nil {
		t.Fatal(err)
	}
	var avisos map[string]int
	confirm := func(pl refinar.Plano, r func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
		avisos = map[string]int{"qa": len(pl.Avisos)}
		other := pl.Mudanca
		other.Peca = "peca-que-nao-existe"
		novo, err := r(other)
		avisos["fora"] = len(novo.Avisos)
		if aviso := "A peça peca-que-nao-existe não está na matriz"; !strings.Contains(novo.Diff(), "AVISO: "+aviso) || !strings.Contains(novo.Cartao().Aviso, aviso) {
			t.Errorf("o aviso da peça deve sair no diff e no cartão:\n%s\n%s", novo.Diff(), novo.Cartao().Aviso)
		}
		return err == nil, err
	}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, confirm, inicio)); err != nil {
		t.Fatal(err)
	}
	if avisos["qa"] != 0 || avisos["fora"] != 1 {
		t.Errorf("avisos da peça na matriz: %v", avisos)
	}
	// só o título muda (a peça fica a mesma): sem aviso, mesmo que a peça atual não conste na matriz; matriz
	// sem entradas (v1) nunca avisa.
	casos := []struct {
		rotulo string
		m      refinar.Mudanca
	}{
		{"só título", refinar.Mudanca{Etapa: "e01", Titulo: "Outro título", Casos: []string{"tdd"}}},
		{"matriz v1", refinar.Mudanca{Etapa: "e02", Peca: "qualquer-peca", Casos: []string{"tdd"}}},
	}
	for _, tc := range casos {
		rotulo, m := tc.rotulo, tc.m
		if rotulo == "matriz v1" {
			if err := os.WriteFile(filepath.Join(c.vault, "registry_matrix.json"), []byte(`{"matrix_version":"1.0.0","ratings":[]}`), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		n := -1
		confirm := func(pl refinar.Plano, _ func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
			n = len(pl.Avisos)
			return false, nil
		}
		_, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: m}, depsDe(c, nil, confirm, inicio))
		if !errors.Is(err, refinar.ErrCancelado) || n != 0 {
			t.Errorf("%s: erro %v e %d avisos; esperava o plano pronto, cancelado, com 0 avisos", rotulo, err, n)
		}
	}
}

func testCartao(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	var texto string
	confirm := func(pl refinar.Plano, _ func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
		ficha := pl.Cartao()
		if vs := ste.LintCartao(ficha); len(vs) != 0 {
			t.Errorf("o cartão viola o STE-PT: %v", vs)
		}
		texto = ficha.Render()
		if len(ficha.Parametros) != 2 || ficha.Parametros[0].Nome != "peca" || ficha.Parametros[0].Valor != "aidd-qa" ||
			ficha.Parametros[1].Nome != "casos" || ficha.Parametros[1].Valor != "verificacao-de-qualidade" {
			t.Errorf("parâmetros do personalizar: %+v", ficha.Parametros)
		}
		return false, nil
	}
	_, _ = refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, confirm, inicio))
	for _, quer := range []string{".legoloop-bak", "1.0.0", "1.1.0", "volta à versão 1.0.0", "[1] seguir", "troca 8 arquivos", "Efeito: troca "} {
		if !strings.Contains(texto, quer) {
			t.Errorf("o cartão não diz %q:\n%s", quer, texto)
		}
	}
	if i, j := strings.Index(texto, "AVISO:"), strings.Index(texto, "O QUE VAI ACONTECER"); i < 0 || i > j {
		t.Errorf("o aviso de risco deve vir antes dos passos:\n%s", texto)
	}
	if vs := ste.Default.Violations(); len(vs) != 0 {
		t.Errorf("mensagens do catálogo fora do perfil: %v", vs)
	}
}

func testContexto(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	antes := arvore(t, c)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := refinar.Run(ctx, refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, seguir, inicio)); !errors.Is(err, context.Canceled) {
		t.Fatalf("erro %v; esperava context.Canceled", err)
	}
	igual(t, antes, arvore(t, c))
}
