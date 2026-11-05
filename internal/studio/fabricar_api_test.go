package studio

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/matriz"
	"legoloop/internal/reputation"
)

var studioNow = time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)

const studioSpec = `# 📋 SPEC: Processo gerado

## 1. Premissas Técnicas Assumidas (Addy Osmani)
1. Go conforme o go.mod.

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
* Sempre: TDD.

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
Costura primária: Fabricar.

## 4. Contratos de Dados & Schemas Formais
Structs.

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)
* [ ] AC-E01.1: Primeiro passo
  * **Dado que:** x
  * **Quando:** y
  * **Então:** exit_code == 0
* [ ] AC-E01.2: Segundo passo
  * **Dado que:** x
  * **Quando:** y
  * **Então:** exit_code == 0

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
* Test: go test ./...

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**
`

type specOnlyLLM struct {
	fail  error
	calls atomic.Int32
}

func (l *specOnlyLLM) Chat(_ context.Context, msgs []interview.Message) (string, error) {
	l.calls.Add(1)
	if l.fail != nil {
		return "", l.fail
	}
	if strings.Contains(msgs[0].Content, "vocabulário fechado do LegoLoop") {
		return `{"natureza": "software", "etapas": ["construir", "verificar"], "casos": ["tdd", "verificacao-de-qualidade"]}`, nil
	}
	return studioSpec, nil
}

// studioMatriz é a matriz em memória dos testes do Studio: alfa constrói (tdd) e beta verifica (qualidade), de software.
func studioMatriz() fabricar.MatrizFunc {
	apta := func(nome string, score float64, etapa, caso string) matriz.Entrada {
		ref := filepath.Join(os.TempDir(), "s", nome, "SKILL.md")
		return matriz.Entrada{Ref: ref, Nome: nome, Tipo: "skill", Hash: "h-" + nome,
			Perfil: &matriz.PerfilDeUso{Etapas: []string{etapa}, Naturezas: []string{"software"}, CasosDeUso: []string{caso},
				Origem: matriz.OrigemDeclarado},
			Rating: reputation.PieceRating{Ref: ref, Name: nome, TotalScore: score, Tier: reputation.TierOuro}}
	}
	return func() (fabricar.FonteMatriz, error) {
		return fabricar.FonteMatriz{Entradas: []matriz.Entrada{apta("alfa", 0.9, "construir", "tdd"), apta("beta", 0.8, "verificar", "verificacao-de-qualidade")},
			Vocab: matriz.DefaultVocab(), GeradaEm: studioNow.Add(-48 * time.Hour)}, nil
	}
}

// realFabricar liga ao Studio a fábrica de verdade (LLM e scanner falsos, raízes no servidor), do jeito que o
// padrão faria: as raízes vêm do servidor e nunca do cliente.
func realFabricar(s *Server, llm *specOnlyLLM) FabricarFunc {
	return func(ctx context.Context, req fabricar.Request) (*fabricar.Result, error) {
		vault, err := s.vault()
		if err != nil {
			return nil, err
		}
		return fabricar.Fabricar(ctx, req, fabricar.Deps{LLM: llm, Matriz: studioMatriz(), Clock: func() time.Time { return studioNow },
			Roots: fabricar.Roots{Vault: vault, Docs: s.ProjectDir}})
	}
}

const okBody = `{"name":"Extrator de Áudio","goal":"extrair áudio de vídeos","category":"build","model":"code-pro"}`

func postFabricar(t *testing.T, s *Server, body string, mut ...func(*http.Request)) (int, string) {
	t.Helper()
	rec := do(t, s.Handler(), "POST", "/api/fabricar/run", body, mut...)
	return rec.Code, rec.Body.String()
}

func TestHandleFabricarRun(t *testing.T) {
	t.Run("GET devolve 405", func(t *testing.T) {
		s, _, _ := newAssetsServer(t)
		s.Fabricar = realFabricar(s, &specOnlyLLM{})
		if rec := do(t, s.Handler(), "GET", "/api/fabricar/run", ""); rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("GET = %d, esperado 405", rec.Code)
		}
	})

	t.Run("origem cruzada, host invalido e content-type errado sao rejeitados", func(t *testing.T) {
		s, _, _ := newAssetsServer(t)
		llm := &specOnlyLLM{}
		s.Fabricar = realFabricar(s, llm)
		if code, _ := postFabricar(t, s, okBody, func(r *http.Request) { r.Header.Set("Origin", "http://evil.example") }); code != http.StatusForbidden {
			t.Errorf("origem cruzada = %d, esperado 403", code)
		}
		if code, _ := postFabricar(t, s, okBody, func(r *http.Request) { r.Host = "evil.example" }); code != http.StatusMisdirectedRequest {
			t.Errorf("host inválido = %d, esperado 421", code)
		}
		if code, _ := postFabricar(t, s, okBody, func(r *http.Request) { r.Header.Set("Content-Type", "text/plain") }); code != http.StatusUnsupportedMediaType {
			t.Errorf("content-type = %d, esperado 415", code)
		}
		if llm.calls.Load() != 0 {
			t.Error("requisição rejeitada não pode gastar LLM")
		}
	})

	t.Run("nome invalido devolve 400 sem chamar a fabrica", func(t *testing.T) {
		s, _, _ := newAssetsServer(t)
		var called atomic.Bool
		s.Fabricar = func(context.Context, fabricar.Request) (*fabricar.Result, error) { called.Store(true); return nil, nil }
		code, body := postFabricar(t, s, `{"name":"!!!","model":"code-pro"}`)
		if code != http.StatusBadRequest || called.Load() {
			t.Errorf("nome inválido = %d (fábrica chamada=%v): %s", code, called.Load(), body)
		}
	})

	t.Run("campos desconhecidos e caminhos do cliente sao recusados", func(t *testing.T) {
		s, _, _ := newAssetsServer(t)
		s.Fabricar = realFabricar(s, &specOnlyLLM{})
		for _, extra := range []string{`"vault":"C:\x"`, `"docs":"/tmp"`, `"trust_removable":true`, `"link_harnesses":["claude"]`, `"target_drives":["C:\\"]`} {
			body := `{"name":"x","model":"code-pro",` + extra + `}`
			if code, _ := postFabricar(t, s, body); code != http.StatusBadRequest {
				t.Errorf("%s = %d, esperado 400: o cliente não define caminhos nem confiança", extra, code)
			}
		}
	})

	t.Run("modelo fora do inventario devolve 400", func(t *testing.T) {
		s, _, _ := newAssetsServer(t)
		var called atomic.Bool
		s.Fabricar = func(context.Context, fabricar.Request) (*fabricar.Result, error) { called.Store(true); return nil, nil }
		if code, _ := postFabricar(t, s, `{"name":"x","model":"gpt-inexistente"}`); code != http.StatusBadRequest || called.Load() {
			t.Errorf("modelo fora do inventário = %d (chamada=%v)", code, called.Load())
		}
	})

	t.Run("processo existente devolve 409", func(t *testing.T) {
		s, _, _ := newAssetsServer(t)
		s.Fabricar = realFabricar(s, &specOnlyLLM{})
		if code, body := postFabricar(t, s, okBody); code != 200 {
			t.Fatalf("1ª fabricação = %d: %s", code, body)
		}
		code, body := postFabricar(t, s, okBody)
		if code != http.StatusConflict {
			t.Errorf("duplicata = %d, esperado 409", code)
		}
		if !strings.Contains(body, `"code":"exists"`) {
			t.Errorf("o 409 de duplicata deve trazer code=exists, para a UI oferecer a sobrescrita: %s", body)
		}
	})

	t.Run("nome com ../ nao grava fora das raizes do servidor", func(t *testing.T) {
		s, home, proj := newAssetsServer(t)
		s.Vault = filepath.Join(home, ".aidd", "canonical")
		s.Fabricar = realFabricar(s, &specOnlyLLM{})
		code, body := postFabricar(t, s, `{"name":"../../etc/passwd","model":"code-pro"}`)
		if code != 200 || !strings.Contains(body, "legoloop-etc-passwd") {
			t.Fatalf("= %d: %s", code, body)
		}
		for _, root := range []string{home, proj} {
			_ = filepath.WalkDir(root, func(p string, e os.DirEntry, _ error) error {
				if !e.IsDir() && !strings.HasPrefix(p, s.Vault) && !strings.HasPrefix(p, filepath.Join(proj, "docs")) {
					t.Errorf("arquivo fora das raízes do servidor: %s", p)
				}
				return nil
			})
		}
	})

	t.Run("sucesso devolve o Result em JSON sem caminhos absolutos", func(t *testing.T) {
		s, home, proj := newAssetsServer(t)
		s.Vault = filepath.Join(home, ".aidd", "canonical")
		s.Fabricar = realFabricar(s, &specOnlyLLM{})
		rec := do(t, s.Handler(), "POST", "/api/fabricar/run", okBody)
		if rec.Code != 200 || !strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
			t.Fatalf("= %d, Content-Type %q: %s", rec.Code, rec.Header().Get("Content-Type"), rec.Body)
		}
		var res fabricar.Result
		decode(t, rec, &res)
		if res.CanonicalName != "legoloop-extrator-de-audio" || len(res.CuratedPieces) == 0 {
			t.Errorf("Result inesperado: %+v", res)
		}
		for _, abs := range []string{home, proj} {
			for _, v := range []string{abs, filepath.ToSlash(abs), strings.ReplaceAll(abs, "\\", "\\\\")} {
				if strings.Contains(rec.Body.String(), v) {
					t.Errorf("a resposta expõe o caminho local %q:\n%s", v, rec.Body)
				}
			}
		}
	})
}

// badLLM sempre responde fora do contrato da SPEC.
type badLLM struct{}

func (badLLM) Chat(context.Context, []interview.Message) (string, error) {
	return "# 📋 SPEC: Ruim\n\nsem as seções do contrato", nil
}
