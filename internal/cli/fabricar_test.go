package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/matriz"
	"legoloop/internal/provider"
	"legoloop/internal/reputation"
)

var fabCLINow = time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)

// specFake é uma SPEC mínima, porém conforme o gabarito, devolvida pelo LLM falso.
const specFake = `# 📋 SPEC: Processo gerado

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

type fakeLLM struct{ calls int }

func (f *fakeLLM) Chat(_ context.Context, msgs []interview.Message) (string, error) {
	f.calls++
	if len(msgs) > 0 && strings.Contains(msgs[0].Content, "vocabulário fechado do LegoLoop") {
		return `{"natureza": "software", "etapas": ["construir", "verificar"], "casos": ["tdd", "verificacao-de-qualidade"]}`, nil
	}
	return specFake, nil
}

// matrizCLI é a matriz em memória dos testes da CLI: alfa constrói (tdd) e beta verifica (qualidade), de software.
func matrizCLI() fabricar.MatrizFunc {
	apta := func(nome string, score float64, etapa, caso string) matriz.Entrada {
		ref := filepath.Join(os.TempDir(), "s", nome, "SKILL.md")
		return matriz.Entrada{Ref: ref, Nome: nome, Tipo: "skill", Hash: "h-" + nome,
			Perfil: &matriz.PerfilDeUso{Etapas: []string{etapa}, Naturezas: []string{"software"}, CasosDeUso: []string{caso},
				Origem: matriz.OrigemDeclarado},
			Rating: reputation.PieceRating{Ref: ref, Name: nome, TotalScore: score, Tier: reputation.TierOuro}}
	}
	return func() (fabricar.FonteMatriz, error) {
		return fabricar.FonteMatriz{Entradas: []matriz.Entrada{apta("alfa", 0.9, "construir", "tdd"), apta("beta", 0.8, "verificar", "verificacao-de-qualidade")},
			Vocab: matriz.DefaultVocab(), GeradaEm: fabCLINow.Add(-48 * time.Hour)}, nil
	}
}

// fabRig monta Deps hermético: inventário com o modelo "sonnet", LLM falso, scanner fixo e relógio fixo.
type fabRig struct {
	d          Deps
	out, errb  *bytes.Buffer
	vault, doc string
	llm        *fakeLLM
}

func newFabRig(t *testing.T) *fabRig {
	t.Helper()
	d, out, errb := deps(t)
	home := t.TempDir()
	llm := &fakeLLM{}
	d.Home = home
	d.Stdin = strings.NewReader("")
	d.Discover = func(context.Context) provider.Inventory {
		return provider.Inventory{CLIs: []provider.CLI{{Name: "claude"}}}
	}
	d.Fabricar = &FabricarDeps{LLM: llm, Matriz: matrizCLI(), Clock: func() time.Time { return fabCLINow }}
	return &fabRig{d: d, out: out, errb: errb, llm: llm, vault: filepath.Join(home, "vault"), doc: filepath.Join(home, "proj")}
}

// args monta a linha de comando com o cofre e o projeto isolados, e o disco C:\ fixo.
func (r *fabRig) args(name string, extra ...string) []string {
	base := []string{"fabricar", name, "--vault", r.vault, "--docs", r.doc, "--model", "sonnet"}
	return append(base, extra...)
}

func (r *fabRig) files(t *testing.T) []string {
	t.Helper()
	var out []string
	for _, root := range []string{r.vault, r.doc} {
		_ = filepath.WalkDir(root, func(p string, e os.DirEntry, err error) error {
			if err == nil && !e.IsDir() {
				out = append(out, p)
			}
			return nil
		})
	}
	return out
}

func TestRunFabricar(t *testing.T) {
	t.Run("autonomo imprime as 5 camadas e os caminhos", func(t *testing.T) {
		r := newFabRig(t)
		if code := run(t, r.d, r.args("extrator de audio", "--yes", "--sem-entrevista", "--goal", "extrair áudio")...); code != 0 {
			t.Fatalf("exit %d\nstdout:%s\nstderr:%s", code, r.out, r.errb)
		}
		for i := 1; i <= 5; i++ {
			if want := "[" + string(rune('0'+i)) + "/5]"; !strings.Contains(r.out.String(), want) {
				t.Errorf("sem progresso %s:\n%s", want, r.out)
			}
		}
		for _, want := range []string{"legoloop-extrator-de-audio", "SKILL.md", ".studio.json", ".dag.json", "INTENT", "SPEC", "PLAN", "ROUTE", "TEMPLATE"} {
			if !strings.Contains(r.out.String(), want) {
				t.Errorf("saída sem %q:\n%s", want, r.out)
			}
		}
		if len(r.files(t)) != 10 {
			t.Errorf("devia gravar 10 arquivos: %v", r.files(t))
		}
	})

	t.Run("o dag gerado passa no validate", func(t *testing.T) {
		r := newFabRig(t)
		run(t, r.d, r.args("extrator de audio", "--yes", "--sem-entrevista")...)
		dagPath := filepath.Join(r.vault, caminhos.DirLoops, "legoloop-extrator-de-audio.dag.json")
		d2, _, errb := deps(t)
		if code := run(t, d2, "validate", dagPath); code != 0 {
			t.Fatalf("validate exit %d: %s", code, errb)
		}
	})

	t.Run("json imprime o Result e deixa o stdout puro", func(t *testing.T) {
		r := newFabRig(t)
		if code := run(t, r.d, r.args("extrator de audio", "--yes", "--sem-entrevista", "--json")...); code != 0 {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
		var res fabricar.Result
		if err := json.Unmarshal(r.out.Bytes(), &res); err != nil {
			t.Fatalf("stdout não é JSON puro: %v\n%s", err, r.out)
		}
		if res.CanonicalName != "legoloop-extrator-de-audio" || !res.CreatedAt.Equal(fabCLINow) {
			t.Errorf("Result inesperado: %+v", res)
		}
		if !strings.Contains(r.errb.String(), "[5/5]") {
			t.Errorf("o progresso deve ir ao stderr com --json:\n%s", r.errb)
		}
	})

	t.Run("nome invalido sai com 2 e nao grava", func(t *testing.T) {
		r := newFabRig(t)
		if code := run(t, r.d, r.args("!!!", "--yes", "--sem-entrevista")...); code != 2 {
			t.Fatalf("exit %d, esperado 2", code)
		}
		if f := r.files(t); len(f) != 0 {
			t.Errorf("gravou com nome inválido: %v", f)
		}
	})

	t.Run("nome com traversal e sanitizado e fica dentro das raizes", func(t *testing.T) {
		r := newFabRig(t)
		if code := run(t, r.d, r.args("../../etc/passwd", "--yes", "--sem-entrevista")...); code != 0 {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
		if !strings.Contains(r.out.String(), "legoloop-etc-passwd") {
			t.Errorf("o nome devia ser sanitizado:\n%s", r.out)
		}
		for _, f := range r.files(t) {
			if !strings.HasPrefix(f, r.vault) && !strings.HasPrefix(f, r.doc) {
				t.Errorf("arquivo fora das raízes: %s", f)
			}
		}
	})

	t.Run("categoria invalida sai com 2", func(t *testing.T) {
		r := newFabRig(t)
		if code := run(t, r.d, r.args("x", "--yes", "--sem-entrevista", "--category", "explodir")...); code != 2 {
			t.Fatalf("exit %d, esperado 2", code)
		}
	})

	t.Run("modelo fora do inventario sai com 2 sem gravar", func(t *testing.T) {
		r := newFabRig(t)
		args := append(r.args("extrator", "--yes", "--sem-entrevista"), "--model", "gpt-inexistente")
		if code := run(t, r.d, args...); code != 2 {
			t.Fatalf("exit %d, esperado 2 (uso incorreto)", code)
		}
		if !strings.Contains(r.errb.String(), "gpt-inexistente") {
			t.Errorf("o erro deve citar o modelo:\n%s", r.errb)
		}
		if f := r.files(t); len(f) != 0 || r.llm.calls != 0 {
			t.Errorf("não devia gastar LLM nem gravar: %v calls=%d", f, r.llm.calls)
		}
	})

	t.Run("processo existente sai com 1 e force sobrescreve com backup", func(t *testing.T) {
		r := newFabRig(t)
		run(t, r.d, r.args("extrator", "--yes", "--sem-entrevista")...)
		before := r.files(t)
		r.out.Reset()
		r.errb.Reset()
		if code := run(t, r.d, r.args("extrator", "--yes", "--sem-entrevista")...); code != 1 {
			t.Fatalf("exit %d, esperado 1 (ErrExists)", code)
		}
		if len(r.files(t)) != len(before) {
			t.Errorf("ErrExists não pode alterar o disco")
		}
		if code := run(t, r.d, r.args("extrator", "--yes", "--sem-entrevista", "--force")...); code != 0 {
			t.Fatalf("--force exit %d: %s", code, r.errb)
		}
	})

	t.Run("dry-run nao grava nada", func(t *testing.T) {
		r := newFabRig(t)
		if code := run(t, r.d, r.args("extrator", "--yes", "--sem-entrevista", "--dry-run")...); code != 0 {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
		if f := r.files(t); len(f) != 0 {
			t.Errorf("dry-run gravou: %v", f)
		}
		if !strings.Contains(r.out.String(), "dry-run") {
			t.Errorf("a saída deve avisar que é dry-run:\n%s", r.out)
		}
	})

	t.Run("sem nenhum modelo disponivel falha com 1 e mensagem clara", func(t *testing.T) {
		r := newFabRig(t)
		r.d.Discover = func(context.Context) provider.Inventory { return provider.Inventory{} }
		r.d.Fabricar = &FabricarDeps{Matriz: matrizCLI(), Clock: func() time.Time { return fabCLINow }}
		if code := run(t, r.d, r.args("extrator", "--yes", "--sem-entrevista")...); code != 1 {
			t.Fatalf("exit %d, esperado 1", code)
		}
		if !strings.Contains(r.errb.String(), "9Router") {
			t.Errorf("a mensagem deve orientar sobre o 9Router:\n%s", r.errb)
		}
	})
}
