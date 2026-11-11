//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/cli"
	"legoloop/internal/dag"
	"legoloop/internal/fabricar"
	"legoloop/internal/gate"
	"legoloop/internal/interview"
	"legoloop/internal/matriz"
	"legoloop/internal/processo"
	"legoloop/internal/provider"
	"legoloop/internal/reputation"
	"legoloop/internal/studio"
)

var e2eNow = time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)

type cenarioRef struct {
	nome      string
	objetivo  string
	categoria fabricar.Category
	porta     string // "cli", "chat", "studio"
	cerebro   string // "code-pro", "agy:gemini-3.8-flash"
}

func matrizCenarios() fabricar.MatrizFunc {
	apta := func(nome, tipo, etapa, natureza, caso string, score float64) matriz.Entrada {
		ref := filepath.Join(os.TempDir(), "s", nome, "SKILL.md")
		return matriz.Entrada{
			Ref:  ref,
			Nome: nome,
			Tipo: tipo,
			Hash: "h-" + nome,
			Perfil: &matriz.PerfilDeUso{
				Etapas:     []string{etapa},
				Naturezas:  []string{natureza},
				CasosDeUso: []string{caso},
				Origem:     matriz.OrigemDeclarado,
			},
			Rating: reputation.PieceRating{
				Ref:        ref,
				Name:       nome,
				TotalScore: score,
				Tier:       reputation.TierOuro,
			},
		}
	}
	entradas := []matriz.Entrada{
		// Peças para Software / build-engine
		apta("compilador-tdd", "skill", "construir", "software", "tdd", 0.95),
		apta("linter-qualidade", "skill", "verificar", "software", "verificacao-de-qualidade", 0.92),
		// Peças para Dados / coleta-de-precos
		apta("web-scraper-precos", "skill", "coletar", "dados", "coleta-de-dados", 0.94),
		apta("conversor-tabelas", "skill", "transformar", "dados", "analise-de-dados", 0.90),
		// Peças para Documentação / manual-de-processo
		apta("redator-manual", "skill", "documentar", "documentacao", "documentacao-tecnica", 0.93),
		apta("revisor-governanca", "skill", "verificar", "documentacao", "governanca-de-processo", 0.91),
	}
	return func() (fabricar.FonteMatriz, error) {
		return fabricar.FonteMatriz{
			Entradas: entradas,
			Vocab:    matriz.DefaultVocab(),
			GeradaEm: e2eNow.Add(-24 * time.Hour),
		}, nil
	}
}

type fakeE2EBrain struct {
	nome string
}

const e2eSpecMD = `# 📋 SPEC: Processo de Referência

## 1. Premissas Técnicas Assumidas (Addy Osmani)
1. Conformidade estrita com o contrato LegoLoop.
2. Isolamento de ambiente e execução hermética.

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
* Sempre: TDD com testes determinísticos.
* Nunca: bypass de gates de qualidade.

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
Costura primária do processo validada em cada camada.

## 4. Contratos de Dados & Schemas Formais
Schemas e payloads definidos no padrão canônico.

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)
* [ ] AC-E01.1: Execução do primeiro passo
  * **Dado que:** ambiente isolado e entradas fornecidas
  * **Quando:** a primeira etapa processa os dados
  * **Então:** exit_code == 0 && file_exists('step1_output.txt')
* [ ] AC-E01.2: Verificação do segundo passo
  * **Dado que:** saída da etapa anterior disponível
  * **Quando:** a segunda etapa valida o resultado
  * **Então:** exit_code == 0 && file_exists('final_output.txt')

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
* Test: go test ./...

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**
`

const e2eMirrorReply = `HIPÓTESE: processo de referência
CONFIANÇA: 98%
ESPELHO CANÔNICO:
- Resultado Esperado: execução automatizada ponta a ponta
- Quem se Beneficia: time de engenharia e operações
- Por que Agora: ganho de escala e governança determinística
- Critério de Sucesso: execução completa aprovada com zero falhas
- Restrições Inegociáveis: hermetismo absoluto e zero stubs
- Fora de Escopo: alterações fora das peças autorizadas`

func (f *fakeE2EBrain) Chat(_ context.Context, msgs []interview.Message) (string, error) {
	if len(msgs) == 0 {
		return "", nil
	}
	last := msgs[len(msgs)-1].Content
	first := msgs[0].Content

	switch {
	case strings.Contains(first, "Arquiteto de Especificações"):
		return e2eSpecMD, nil
	case strings.Contains(first, "vocabulário fechado do LegoLoop"):
		if strings.Contains(last, "coleta-de-precos") || strings.Contains(last, "preços") || strings.Contains(last, "dados") {
			return `{"natureza": "dados", "etapas": ["coletar", "transformar"], "casos": ["coleta-de-dados", "analise-de-dados"]}`, nil
		}
		if strings.Contains(last, "manual-de-processo") || strings.Contains(last, "documentação") {
			return `{"natureza": "documentacao", "etapas": ["documentar", "verificar"], "casos": ["documentacao-tecnica", "governanca-de-processo"]}`, nil
		}
		return `{"natureza": "software", "etapas": ["construir", "verificar"], "casos": ["tdd", "verificacao-de-qualidade"]}`, nil
	case strings.Contains(first, "explicador da Fusão"):
		return `{"explicacoes": {}}`, nil
	}
	return e2eMirrorReply, nil
}

func setupHermeticEnv(t *testing.T) (string, string, string) {
	t.Helper()
	tmp := t.TempDir()
	vault := filepath.Join(tmp, "vault")
	docs := filepath.Join(tmp, "docs")
	home := filepath.Join(tmp, "home")
	for _, d := range []string{vault, docs, home} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return vault, docs, home
}

func fakeClarification() *fabricar.Clarification {
	return &fabricar.Clarification{
		Reached:    true,
		Confidence: 98,
		Turns:      1,
		Mirror: interview.Mirror{
			Outcome:     "execução automatizada ponta a ponta",
			User:        "time de engenharia e operações",
			WhyNow:      "ganho de escala e governança determinística",
			Success:     "execução completa aprovada com zero falhas",
			Constraints: []string{"hermetismo absoluto e zero stubs"},
			OutOfScope:  []string{"alterações fora das peças autorizadas"},
		},
	}
}

func runCenarioCLI(t *testing.T, c cenarioRef, vault, docs, home string) (*fabricar.Result, error) {
	var stdout, stderr bytes.Buffer
	llm := &fakeE2EBrain{nome: c.cerebro}
	d := cli.Deps{
		Stdout: &stdout,
		Stderr: &stderr,
		Home:   home,
		Discover: func(context.Context) provider.Inventory {
			return provider.Inventory{
				CLIs: []provider.CLI{{Name: "claude"}, {Name: "agy"}},
			}
		},
		Fabricar: &cli.FabricarDeps{
			LLM:    llm,
			Matriz: matrizCenarios(),
			Clock:  func() time.Time { return e2eNow },
			Harnesses: []processo.Harness{
				{Nome: "chat", Skills: filepath.Join(home, "skills")},
			},
		},
	}

	// 1. Inicia entrevista CLI
	stdout.Reset()
	stderr.Reset()
	if code := cli.Run(context.Background(), []string{"fabricar", "entrevista", "iniciar", c.nome, "--vault", vault, "--docs", docs, "--model", "code-pro"}, d); code != 0 {
		return nil, fmt.Errorf("entrevista iniciar falhou (exit %d): %s", code, stderr.String())
	}
	var p struct {
		Sessao string `json:"sessao"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &p); err != nil || p.Sessao == "" {
		return nil, fmt.Errorf("sessao inválida na entrevista: %w (out: %s)", err, stdout.String())
	}

	// 2. Se não estiver finalizada, responde entrevista CLI
	if p.Status != "done" {
		stdout.Reset()
		stderr.Reset()
		_ = cli.Run(context.Background(), []string{"fabricar", "entrevista", "responder", p.Sessao, "tdd com compilador", "--vault", vault, "--docs", docs, "--model", "code-pro"}, d)
	}

	// 3. Executa fabricar com a entrevista
	stdout.Reset()
	stderr.Reset()
	args := []string{
		"fabricar", c.nome,
		"--yes",
		"--entrevista", p.Sessao,
		"--category", string(c.categoria),
		"--vault", vault,
		"--docs", docs,
		"--json",
	}
	code := cli.Run(context.Background(), args, d)
	if code != 0 {
		return nil, fmt.Errorf("cli fabricar falhou (exit %d): %s", code, stderr.String())
	}
	var res fabricar.Result
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		return nil, fmt.Errorf("json inválido na saída da CLI: %w\n%s", err, stdout.String())
	}
	return &res, nil
}

func runCenarioChat(t *testing.T, c cenarioRef, vault, docs, home string) (*fabricar.Result, error) {
	deps := fabricar.Deps{
		Roots:     fabricar.Roots{Vault: vault, Docs: docs},
		Clock:     func() time.Time { return e2eNow },
		LLM:       &fakeE2EBrain{nome: c.cerebro},
		Matriz:    matrizCenarios(),
		FS:        fabricar.OSFS(),
		Interview: fakeClarification(),
		Harnesses: []processo.Harness{
			{Nome: "chat", Skills: filepath.Join(home, "skills")},
		},
	}
	req := fabricar.Request{
		RawName:     c.nome,
		Goal:        c.objetivo,
		Category:    c.categoria,
		ModelChoice: c.cerebro,
		Force:       true,
	}
	return fabricar.Fabricar(context.Background(), req, deps)
}

func runCenarioStudio(t *testing.T, c cenarioRef, vault, docs, home string) (*fabricar.Result, error) {
	srv := &studio.Server{
		Host:       "127.0.0.1:7878",
		ProjectDir: docs,
		Fabricar: func(ctx context.Context, req fabricar.Request) (*fabricar.Result, error) {
			deps := fabricar.Deps{
				Roots:     fabricar.Roots{Vault: vault, Docs: docs},
				Clock:     func() time.Time { return e2eNow },
				LLM:       &fakeE2EBrain{nome: c.cerebro},
				Matriz:    matrizCenarios(),
				FS:        fabricar.OSFS(),
				Interview: fakeClarification(),
				Harnesses: []processo.Harness{
					{Nome: "chat", Skills: filepath.Join(home, "skills")},
				},
			}
			return fabricar.Fabricar(ctx, req, deps)
		},
	}

	body, _ := json.Marshal(map[string]any{
		"name":     c.nome,
		"goal":     c.objetivo,
		"category": string(c.categoria),
		"model":    c.cerebro,
		"force":    true,
	})
	req := httptest.NewRequest("POST", "http://127.0.0.1:7878/api/fabricar/run", bytes.NewReader(body))
	req.Host = "127.0.0.1:7878"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://127.0.0.1:7878")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		return nil, fmt.Errorf("studio /api/fabricar/run retornou %d: %s", rec.Code, rec.Body.String())
	}
	var res fabricar.Result
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		return nil, fmt.Errorf("json inválido do studio: %w", err)
	}
	return &res, nil
}

func validarChecklistQualidade(t *testing.T, res *fabricar.Result, c cenarioRef, vault, docs, home string) {
	t.Helper()
	if res.CanonicalName == "" {
		t.Fatalf("[%s] CanonicalName vazio", c.nome)
	}
	realHome, _ := os.UserHomeDir()
	a := res.Artifacts
	caminhos := []string{a.IntentPath, a.SpecPath, a.PlanPath, a.RoutePath, a.DagPath, a.LoopPath, a.SkillPath}
	for _, p := range caminhos {
		if p == "" {
			t.Errorf("[%s] Caminho de artefato vazio", c.nome)
			continue
		}
		// Desmascara se for resultado do Studio (<projeto> -> docs, ~ -> realHome)
		pReal := strings.ReplaceAll(p, "<projeto>", docs)
		pReal = strings.ReplaceAll(pReal, "~", realHome)
		if _, err := os.Stat(pReal); err != nil {
			t.Errorf("[%s] Artefato não encontrado no disco: %s (real: %s)", c.nome, p, pReal)
		}
	}

	// 1. INTENT sem 'Não informado'
	intentPathReal := strings.ReplaceAll(a.IntentPath, "<projeto>", docs)
	intentPathReal = strings.ReplaceAll(intentPathReal, "~", realHome)
	intentBytes, err := os.ReadFile(intentPathReal)
	if err == nil {
		if strings.Contains(string(intentBytes), "Não informado") {
			t.Errorf("[%s] INTENT contém 'Não informado':\n%s", c.nome, string(intentBytes))
		}
	}

	// 2. Gates não triviais
	routePathReal := strings.ReplaceAll(a.RoutePath, "<projeto>", docs)
	routePathReal = strings.ReplaceAll(routePathReal, "~", realHome)
	routeBytes, err := os.ReadFile(routePathReal)
	if err == nil {
		var rt map[string]any
		if err := json.Unmarshal(routeBytes, &rt); err == nil {
			if nodes, ok := rt["nodes"].([]any); ok {
				for _, n := range nodes {
					if nm, ok := n.(map[string]any); ok {
						if out, ok := nm["output"].(map[string]any); ok {
							gStr, _ := out["gate"].(string)
							if gStr == "" || gStr == "true" {
								t.Errorf("[%s] Nó %v com gate trivial ou vazio: %q", c.nome, nm["id"], gStr)
							}
							if _, err := gate.Compile(gStr); err != nil {
								t.Errorf("[%s] Gate %q não compila: %v", c.nome, gStr, err)
							}
						}
					}
				}
			}
		}
	}

	// 3. DAG parseável e válido
	dagPathReal := strings.ReplaceAll(strings.ReplaceAll(a.DagPath, "<projeto>", docs), "~", realHome)
	dagBytes, err := os.ReadFile(dagPathReal)
	if err == nil {
		if _, err := dag.Parse(dagBytes); err != nil {
			t.Errorf("[%s] DAG inválido: %v", c.nome, err)
		}
	}
}

func TestE2ECenarios(t *testing.T) {
	cenarios := []cenarioRef{
		{
			nome:      "build-engine",
			objetivo:  "Compilar e testar núcleo central de software com rigor de compilação",
			categoria: fabricar.CategoryBuild,
			porta:     "cli",
			cerebro:   "code-pro",
		},
		{
			nome:      "coleta-de-precos",
			objetivo:  "coletar dados estruturados de preços via web",
			categoria: fabricar.CategoryWeb,
			porta:     "chat",
			cerebro:   "agy:gemini-3.8-flash",
		},
		{
			nome:      "manual-de-processo",
			objetivo:  "documentar tecnicamente processos com governança de processo",
			categoria: fabricar.CategoryBuild,
			porta:     "studio",
			cerebro:   "code-pro",
		},
	}

	for _, c := range cenarios {
		c := c
		t.Run(fmt.Sprintf("%s_%s_%s", c.nome, c.porta, strings.ReplaceAll(c.cerebro, ":", "_")), func(t *testing.T) {
			vault, docs, home := setupHermeticEnv(t)
			var res *fabricar.Result
			var err error

			switch c.porta {
			case "cli":
				res, err = runCenarioCLI(t, c, vault, docs, home)
			case "chat":
				res, err = runCenarioChat(t, c, vault, docs, home)
			case "studio":
				res, err = runCenarioStudio(t, c, vault, docs, home)
			default:
				t.Fatalf("porta desconhecida %s", c.porta)
			}

			if err != nil {
				t.Fatalf("[%s via %s] erro na fabricação: %v", c.nome, c.porta, err)
			}
			validarChecklistQualidade(t, res, c, vault, docs, home)
		})
	}
}
