package orchestrix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
	"legoloop/internal/roteamento"
)

// orchestrixDegrada declara o harness omp via oficial como reprovado e o omp via 9router como homologado.
func orchestrixDegrada(avisos *[]string) *Orchestrix {
	o := unitOrchestrix("")
	o.Inventory.CLIs = append(o.Inventory.CLIs, provider.CLI{Name: "omp", Path: os.Args[0]})
	o.ExtraEnv = []string{"LEGOLOOP_HELPER=omp"}
	o.EstadoHarness = func(harness, via string) string {
		if harness == "omp" && via == "9router" {
			return "homologado"
		}
		return "reprovado"
	}
	o.Elegiveis = map[string][]roteamento.CandidatoVia{}
	for _, classe := range []string{roteamento.ClasseNano, roteamento.ClasseMicro, roteamento.ClasseGrande, roteamento.ClasseExtreme} {
		o.Elegiveis[classe] = []roteamento.CandidatoVia{
			{Harness: "omp", Via: roteamento.ViaOficial, Modelo: "omp:x"},
			{Harness: "omp", Via: roteamento.Via9Router, Modelo: "omp:y"},
		}
	}
	o.Aviso = func(m string) { *avisos = append(*avisos, m) }
	return o
}

func noDaClasse(id, classe string) dag.Node {
	n := unitNode(id, "out/"+id+".json", "exit_code == 0", nil)
	n.ContextOut = nil
	n.Nature = classe
	n.Routing = dag.RoutingConfig{Mode: dag.ModeManual, Harness: "omp", Via: "oficial", Model: "omp:x"}
	return n
}

func TestDegradaNanoMicroFalhaGrandeExtreme(t *testing.T) {
	for _, classe := range []string{roteamento.ClasseNano, roteamento.ClasseMicro} {
		var avisos []string
		repo := newRepoT(t)
		res, err := orchestrixDegrada(&avisos).Execute(context.Background(), noDaClasse("n-"+classe, classe), repo)
		if err != nil {
			t.Fatalf("%s: %v", classe, err)
		}
		notas := strings.Join(res.Notes, "\n")
		if !strings.Contains(notas, "substituicao: n-"+classe+" usou omp/9router/omp:y no lugar de omp/oficial/omp:x") {
			t.Fatalf("%s: sem Substituicao nas notas: %q", classe, notas)
		}
		if len(avisos) != 1 || !strings.Contains(avisos[0], "omp:y") {
			t.Fatalf("%s: avisos = %v", classe, avisos)
		}
		if !strings.HasSuffix(res.Model, "omp:y") {
			t.Fatalf("%s: modelo realmente usado = %q", classe, res.Model)
		}
		var rep unitReport
		bruto, rerr := os.ReadFile(filepath.Join(repo, "out", "n-"+classe+".json"))
		if rerr != nil || json.Unmarshal(bruto, &rep) != nil || rep.Model != "omp:y" {
			t.Fatalf("%s: relatório do nó = %q (%v)", classe, bruto, rerr)
		}
	}
	for _, classe := range []string{roteamento.ClasseGrande, roteamento.ClasseExtreme} {
		var avisos []string
		_, err := orchestrixDegrada(&avisos).Execute(context.Background(), noDaClasse("g-"+classe, classe), newRepoT(t))
		if err == nil || !strings.Contains(err.Error(), "não troca de harness") {
			t.Fatalf("%s: err = %v", classe, err)
		}
		if len(avisos) != 0 {
			t.Fatalf("%s: não pode haver Substituicao: %v", classe, avisos)
		}
	}
}

func TestDegradaSemCandidatoElegivelFalha(t *testing.T) {
	var avisos []string
	o := orchestrixDegrada(&avisos)
	o.Elegiveis[roteamento.ClasseNano] = o.Elegiveis[roteamento.ClasseNano][:1]
	_, err := o.Execute(context.Background(), noDaClasse("n1", roteamento.ClasseNano), newRepoT(t))
	if err == nil || !strings.Contains(err.Error(), "nenhum candidato") {
		t.Fatalf("err = %v", err)
	}
}

func TestNoDiretoChamaApiDoRouter(t *testing.T) {
	var corpo map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" || r.Header.Get("Authorization") != "Bearer "+chaveFalsa {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewDecoder(r.Body).Decode(&corpo)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"saida direta"}}],"usage":{"prompt_tokens":2100,"completion_tokens":5}}`))
	}))
	defer srv.Close()
	o := &Orchestrix{Inventory: provider.Inventory{RouterURL: srv.URL}, HTTP: srv.Client(),
		RouterKey: chaveFalsa, RouterKeyURL: srv.URL, SkillsDir: skillsDir(t)}
	n := skillNode(dag.RoutingConfig{Mode: dag.ModeManual, Harness: "direto", Via: "9router", Model: "code-pro"}, nil)
	work := t.TempDir()
	res, err := o.Execute(context.Background(), n, work)
	if err != nil {
		t.Fatal(err)
	}
	if corpo["model"] != "code-pro" || readArtifact(t, work) != "saida direta" || res.Model != "direto/9router/code-pro" {
		t.Fatalf("corpo = %v, modelo = %q", corpo, res.Model)
	}
	entradas, _ := os.ReadDir(work)
	if len(entradas) != 1 {
		t.Fatalf("o nó direto grava só o artefato de saída: %v", entradas)
	}
	n.Routing.Via = "oficial"
	if _, err := o.Execute(context.Background(), n, t.TempDir()); err == nil || !strings.Contains(err.Error(), "exige via") {
		t.Fatalf("direto com via oficial deveria falhar: %v", err)
	}
}
