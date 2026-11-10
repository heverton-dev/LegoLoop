package cli

import (
	"encoding/json"
	"legoloop/internal/caminhos"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"testing"

	"legoloop/internal/processo"
)

// artefatoReal é a resposta do router falso aos nós de skill: passa no gate markdown da biblioteca (30 palavras).
var artefatoReal = "# Artefato\n\n" + strings.TrimSpace(strings.Repeat("palavra do artefato entregue pelo executor ", 10)) + "\n"

// routerReal é o 9Router falso dos processos reais: responde ao fabricar como o cerebroPorPapel e, aos nós de
// skill do legoloop run, com um artefato que passa no gate. Lista sonnet e code-fast: nenhum nó sai do router.
type routerReal struct {
	cer   cerebroPorPapel
	execs atomic.Int32
}

func (r *routerReal) servidor(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/models":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]string{{"id": "sonnet"}, {"id": "code-fast"}}})
		case "/v1/chat/completions":
			var corpo struct {
				Messages []struct{ Role, Content string } `json:"messages"`
			}
			_ = json.NewDecoder(req.Body).Decode(&corpo)
			primeira := ""
			if len(corpo.Messages) > 0 {
				primeira = corpo.Messages[0].Content
			}
			reply := artefatoReal
			if strings.Contains(primeira, "executor de tarefa única") {
				r.execs.Add(1)
			} else {
				reply = r.cer.responder(primeira)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]string{"content": reply}}}})
		default:
			http.NotFound(w, req)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

// fabricarDois cria legoloop-cp-a e legoloop-cp-b pelo legoloop fabricar de verdade, no mesmo cofre e projeto.
func fabricarDois(t *testing.T, bin string, router *httptest.Server, root, vault, proj string) {
	t.Helper()
	for _, nome := range []string{"cp a", "cp b"} {
		r := runBinaryEm(t, bin, router, root, proj, "fabricar", nome, "--goal", "construir e testar um extrator", "--vault", vault,
			"--docs", proj, "--yes", "--sem-entrevista", "--no-link")
		if r.code != 0 {
			t.Fatalf("fabricar %s: código %d\nstdout:\n%s\nstderr:\n%s", nome, r.code, r.out, r.errb)
		}
	}
	for _, nome := range []string{"legoloop-cp-a", "legoloop-cp-b"} {
		p, err := processo.Ler(processo.CaminhoNoCofre(vault, nome))
		if err != nil || len(p.Etapas) != 2 || p.Etapas[0].Tipo != "construir" || p.Etapas[1].Tipo != "verificar" ||
			p.Etapas[1].Produz.Caminho != "out/02_verificar.md" {
			t.Fatalf("%s: o fabricar deveria gravar construir e verificar em out/02_verificar.md (%v): %+v", nome, err, p.Etapas)
		}
	}
}

// TestComporProcessosReais prova AC-E03.4 com dois processos do fabricar real, que gravam os mesmos caminhos
// (out/01_construir.md e out/02_verificar.md): o compor liga b em a, o DAG valida, o run executa o subloop e o
// artefato de b chega em out/legoloop-cp-b/ sem tocar o de a. O ciclo, o processo inexistente e a versão
// errada saem com 1 e a árvore fica idêntica.
func TestComporProcessosReais(t *testing.T) {
	bin := buildBinary(t)
	rr := &routerReal{}
	router := rr.servidor(t)
	root := t.TempDir()
	vault, proj := filepath.Join(root, "vault"), filepath.Join(root, "proj")
	matrizDeTeste(t, vault)
	fabricarDois(t, bin, router, root, vault, proj)
	r := runBinaryEm(t, bin, router, root, proj, "compor", "legoloop-cp-a", "legoloop-cp-b@1.0.0", "--vault", vault, "--docs", proj, "--yes")
	if r.code != 0 || !strings.Contains(r.out, "Processo legoloop-cp-b@1.0.0 ligado a legoloop-cp-a na etapa e03") {
		t.Fatalf("compor: código %d\nstdout:\n%s\nstderr:\n%s", r.code, r.out, r.errb)
	}
	dagA := filepath.Join(vault, caminhos.DirLoops, "legoloop-cp-a.dag.json")
	if v := runBinaryEm(t, bin, router, root, proj, "validate", dagA); v.code != 0 {
		t.Fatalf("validate: código %d\n%s%s", v.code, v.out, v.errb)
	}
	realRodar(t, bin, router, root, proj, vault, dagA, rr)
	realListar(t, bin, router, root, proj, vault)
	antes := refinArvore(t, root)
	for _, f := range []struct {
		args  []string
		trech string
	}{
		{[]string{"legoloop-cp-b", "legoloop-cp-a@1.1.0"}, "ciclo"},
		{[]string{"legoloop-cp-a", "legoloop-cp-zzz@1.0.0"}, "não existe"},
		{[]string{"legoloop-cp-a", "legoloop-cp-b@9.0.0"}, "versão"},
	} {
		args := append(append([]string{"compor"}, f.args...), "--vault", vault, "--docs", proj, "--yes")
		if r := runBinaryEm(t, bin, router, root, proj, args...); r.code != 1 || !strings.Contains(r.errb, f.trech) {
			t.Errorf("compor %v: código %d (esperava 1) e stderr sem %q:\n%s", f.args, r.code, f.trech, r.errb)
		}
	}
	refinIgual(t, antes, refinArvore(t, root))
}

// realRodar roda o DAG de a: e01 e e02 de a pelo router, e03 é o subloop de b (mais duas chamadas). O artefato
// de b chega em out/legoloop-cp-b/02_verificar.md e o out/02_verificar.md de a continua o de a.
func realRodar(t *testing.T, bin string, router *httptest.Server, root, proj, vault, dagA string, rr *routerReal) {
	t.Helper()
	pecaNoCofre(t, vault)
	work := t.TempDir()
	r := runBinaryEm(t, bin, router, root, proj, "run", dagA, "--workdir", work, "--skills", vault, "--yes", "--json")
	var rep struct {
		OK    bool `json:"ok"`
		Nodes []struct{ ID, Status, Model string }
	}
	if err := json.Unmarshal([]byte(r.out), &rep); err != nil || r.code != 0 || !rep.OK {
		t.Fatalf("run: código %d, json %v\nstdout:\n%s\nstderr:\n%s", r.code, err, r.out, r.errb)
	}
	i := slices.IndexFunc(rep.Nodes, func(n struct{ ID, Status, Model string }) bool { return n.ID == "e03" })
	if i < 0 || rep.Nodes[i].Status != "passed" || rep.Nodes[i].Model != "subloop/legoloop-cp-b" || len(rep.Nodes) != 3 {
		t.Fatalf("o nó e03 não rodou o subloop de b: %+v", rep.Nodes)
	}
	deB, errB := os.ReadFile(filepath.Join(work, "out", "legoloop-cp-b", "02_verificar.md"))
	deA, errA := os.ReadFile(filepath.Join(work, "out", "02_verificar.md"))
	if errB != nil || errA != nil || string(deB) != artefatoReal || string(deA) != artefatoReal || rr.execs.Load() != 4 {
		t.Errorf("artefatos: b %v, a %v; chamadas de nó %d (esperava 2 de a e 2 de b)", errB, errA, rr.execs.Load())
	}
	if _, err := os.Stat(filepath.Join(work, "out", "01_construir.md")); err != nil {
		t.Errorf("o artefato de e01 de a deveria estar no pai: %v", err)
	}
}

// pecaNoCofre põe a peça aidd-tdd no cofre de skills do run (o fabricar só referencia a peça da matriz).
func pecaNoCofre(t *testing.T, vault string) {
	t.Helper()
	p := filepath.Join(vault, "skills", "aidd-tdd", "SKILL.md")
	if _, err := os.Stat(p); err == nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("# aidd-tdd\n\nEscreva o artefato.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func realListar(t *testing.T, bin string, router *httptest.Server, root, proj, vault string) {
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
	if a.Nome != "legoloop-cp-a" || a.Versao != "1.1.0" || b.Versao != "1.0.0" || !slices.Equal(b.UsadoPor, []string{"legoloop-cp-a"}) {
		t.Errorf("listar: %+v; esperava a 1.1.0 e b usado por a", res.Processos)
	}
	txt := runBinaryEm(t, bin, router, root, proj, "listar", "--vault", vault)
	if linha := linhaCom(txt.out, "legoloop-cp-b"); !strings.Contains(linha, "legoloop-cp-a") {
		t.Errorf("o texto do listar não mostra b usado por a:\n%s", txt.out)
	}
}
