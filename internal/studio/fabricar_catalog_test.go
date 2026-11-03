package studio

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"legoloop/internal/fabricar"
)

// catalogServer: HOME, projeto e cofre isolados, com o cofre dentro do HOME para o assets.Scanner enxergá-lo.
func catalogServer(t *testing.T, llm *specOnlyLLM) (*Server, string, string) {
	t.Helper()
	s, home, proj := newAssetsServer(t)
	s.Vault = filepath.Join(home, ".aidd", "canonical")
	s.Fabricar = realFabricar(s, llm)
	return s, home, proj
}

func listedSkill(t *testing.T, s *Server, name string) bool {
	t.Helper()
	rec := do(t, s.Handler(), "GET", "/api/assets/skill", "")
	return rec.Code == 200 && strings.Contains(rec.Body.String(), `"`+name+`"`)
}

func listedLoop(t *testing.T, s *Server, id string) bool {
	t.Helper()
	rec := do(t, s.Handler(), "GET", "/api/loops", "")
	return rec.Code == 200 && strings.Contains(rec.Body.String(), `"`+id+`"`)
}

func TestFabricarRunCatalog(t *testing.T) {
	s, _, _ := catalogServer(t, &specOnlyLLM{})
	const name = "legoloop-extrator-de-audio"
	if listedSkill(t, s, name) || listedLoop(t, s, name) {
		t.Fatal("o processo não pode existir antes de fabricar")
	}
	if code, body := postFabricar(t, s, okBody); code != 200 {
		t.Fatalf("fabricar = %d: %s", code, body)
	}
	if !listedLoop(t, s, name) {
		t.Error("GET /api/loops deve listar o loop novo, sem reiniciar o servidor")
	}
	if !listedSkill(t, s, name) {
		t.Error("GET /api/assets/skill deve listar a skill nova: o handler precisa acionar Assets.Refresh")
	}
}

func TestFabricarRunMutualExclusion(t *testing.T) {
	s, _, _ := catalogServer(t, &specOnlyLLM{})
	started, release := make(chan struct{}), make(chan struct{})
	s.Fabricar = func(ctx context.Context, req fabricar.Request) (*fabricar.Result, error) {
		close(started)
		<-release
		return &fabricar.Result{CanonicalName: "legoloop-x"}, nil
	}
	var first int
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); first, _ = postFabricar(t, s, okBody) }()
	<-started
	code, body := postFabricar(t, s, okBody)
	if code != http.StatusConflict || !strings.Contains(body, `"code":"busy"`) {
		t.Errorf("2ª execução simultânea = %d %s, esperado 409 com code=busy (não pode oferecer sobrescrita)", code, body)
	}
	close(release)
	wg.Wait()
	if first != 200 {
		t.Errorf("a 1ª execução deve concluir normalmente, veio %d", first)
	}
	s.Fabricar = func(context.Context, fabricar.Request) (*fabricar.Result, error) { return &fabricar.Result{}, nil }
	if code, _ := postFabricar(t, s, okBody); code != 200 {
		t.Errorf("após terminar, o lock deve ser liberado: %d", code)
	}
}

func TestFabricarRunSurvivesTheClientAndExtendsTheDeadline(t *testing.T) {
	s, _, _ := catalogServer(t, &specOnlyLLM{})
	var ctxErr error
	var deadline time.Time
	s.Fabricar = func(ctx context.Context, req fabricar.Request) (*fabricar.Result, error) {
		ctxErr = ctx.Err()
		deadline, _ = ctx.Deadline()
		return &fabricar.Result{}, nil
	}
	reqCtx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("POST", "http://"+host+"/api/fabricar/run", strings.NewReader(okBody)).WithContext(reqCtx)
	req.Host = host
	req.Header.Set("Content-Type", "application/json")
	cancel() // a aba fechou
	w := &deadlineRecorder{ResponseRecorder: httptest.NewRecorder()}
	s.Handler().ServeHTTP(w, req)
	if ctxErr != nil {
		t.Errorf("a fabricação não pode morrer com o cliente (WithoutCancel): %v", ctxErr)
	}
	if time.Until(deadline) < time.Minute {
		t.Errorf("a fabricação precisa de um teto longo, veio %s", time.Until(deadline))
	}
	if time.Until(w.deadline) < 5*time.Minute {
		t.Errorf("o WriteTimeout de 30 s do servidor cortaria a resposta; o prazo de escrita deveria ser estendido, veio %s", time.Until(w.deadline))
	}
}

type deadlineRecorder struct {
	*httptest.ResponseRecorder
	deadline time.Time
}

func (d *deadlineRecorder) SetWriteDeadline(t time.Time) error { d.deadline = t; return nil }

func TestFabricarRunFailsCleanly(t *testing.T) {
	cases := []struct {
		name string
		llm  *specOnlyLLM
		want int
	}{
		{"falha de transporte do LLM", &specOnlyLLM{fail: errors.New(`9Router caiu em C:\Users\segredo\token.txt`)}, http.StatusInternalServerError},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, home, proj := catalogServer(t, c.llm)
			code, body := postFabricar(t, s, okBody)
			if code != c.want {
				t.Fatalf("= %d, esperado %d: %s", code, c.want, body)
			}
			if strings.Contains(body, "segredo") || strings.Contains(body, "token.txt") || strings.Contains(body, home) {
				t.Errorf("o erro bruto vazou para o navegador: %s", body)
			}
			var files []string
			for _, root := range []string{filepath.Join(home, ".aidd"), filepath.Join(proj, "docs")} {
				_ = filepath.WalkDir(root, func(p string, e os.DirEntry, _ error) error {
					if e != nil && !e.IsDir() {
						files = append(files, p)
					}
					return nil
				})
			}
			if len(files) != 0 {
				t.Errorf("falha de LLM não pode deixar artefatos parciais: %v", files)
			}
			if listedLoop(t, s, "legoloop-extrator-de-audio") {
				t.Error("o catálogo não pode listar um processo que falhou")
			}
		})
	}
}

func TestFabricarRunContractFailureIs502(t *testing.T) {
	s, _, _ := catalogServer(t, &specOnlyLLM{})
	s.Fabricar = func(ctx context.Context, req fabricar.Request) (*fabricar.Result, error) {
		return fabricar.Fabricar(ctx, req, fabricar.Deps{LLM: badLLM{}, Scanner: studioScan{}, Roots: fabricar.Roots{Vault: s.Vault, Docs: s.ProjectDir}})
	}
	if code, body := postFabricar(t, s, okBody); code != http.StatusBadGateway {
		t.Errorf("SPEC fora do contrato = %d, esperado 502: %s", code, body)
	}
}
