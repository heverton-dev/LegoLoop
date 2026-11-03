package fabricar

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/interview"
	"legoloop/internal/plangen"
)

type namedFile struct {
	Name string
	Data []byte
}

// triadFiles são INTENT, SPEC, PLAN e ROUTE já gerados, ainda fora das raízes finais.
type triadFiles struct{ Intent, Spec, Plan, Route namedFile }

// buildTriad é a Camada 4. Reaproveita os motores existentes: a INTENT vem do renderizador do interview,
// a SPEC do specgen (LLM) e PLAN e ROUTE do plangen.Render (sem LLM). Tudo é gerado num diretório temporário
// e só os bytes seguem adiante, para que a gravação final seja uma única operação atômica.
func buildTriad(ctx context.Context, llm interview.LLMClient, clar Clarification, ps ProcessSpec, dagPath string, now time.Time) (triadFiles, error) {
	stage, err := os.MkdirTemp("", "legoloop-fabricar-*")
	if err != nil {
		return triadFiles{}, err
	}
	defer os.RemoveAll(stage)
	date := now.Format("2006-01-02")
	intentPath, err := writeIntent(stage, ps.Name, clar, date)
	if err != nil {
		return triadFiles{}, err
	}
	spec, err := generateSpec(ctx, llm, stage, intentPath, now)
	if err != nil {
		return triadFiles{}, err
	}
	plan, err := plangen.Render(stage, spec.FilePath, ticketsFor(ps, spec.Criteria, dagPath), plangen.SupervisionEpic, now)
	if err != nil {
		return triadFiles{}, err
	}
	return readTriad(intentPath, spec.FilePath, plan)
}

func upperSnake(slug string) string { return strings.ToUpper(strings.ReplaceAll(slug, "-", "_")) }

func writeIntent(stage, name string, clar Clarification, date string) (string, error) {
	file := fmt.Sprintf("%s_INTENT_%s.md", date, upperSnake(strings.TrimPrefix(name, namePrefix)))
	doc := interview.RenderIntent(interview.IntentView{
		Theme: clar.Theme, Date: date, SessionID: "fabricar-" + strings.TrimPrefix(name, namePrefix),
		Confidence: clar.Confidence, Mirror: &clar.Mirror, SpecRel: interview.SpecRelFor(file, clar.Theme, date),
	})
	dir := filepath.Join(stage, "docs", "intents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, file)
	return path, os.WriteFile(path, []byte(doc), 0o644)
}

// ticketsFor cria um ticket por passo do processo, em cadeia, e reparte os critérios da SPEC entre eles
// (todo AC fica coberto). A verificação do ticket é validar o DAG gerado.
func ticketsFor(ps ProcessSpec, criteria []string, dagPath string) []plangen.Ticket {
	tickets := make([]plangen.Ticket, len(ps.Steps))
	for i, s := range ps.Steps {
		tickets[i] = plangen.Ticket{
			ID: fmt.Sprintf("TICK-01.%d", i+1), Title: "Aplicar a skill " + refName(s.Piece), EpicID: "E01",
			Nature: "micro", ExecutionMode: "sequential", TargetFiles: []string{ps.artifactOf(i)},
			ValidationCmd: "legoloop validate " + filepath.ToSlash(dagPath),
		}
		if i > 0 {
			tickets[i].DependsOn = []string{tickets[i-1].ID}
		}
	}
	for k, ac := range criteria {
		t := &tickets[k%len(tickets)]
		t.SpecCriteria = append(t.SpecCriteria, ac)
	}
	return tickets
}

func readTriad(intentPath, specPath string, plan *plangen.PlanResult) (triadFiles, error) {
	var t triadFiles
	for _, f := range []struct {
		dst  *namedFile
		path string
	}{{&t.Intent, intentPath}, {&t.Spec, specPath}, {&t.Plan, plan.PlanMarkdownPath}, {&t.Route, plan.RouteJSONPath}} {
		data, err := os.ReadFile(f.path)
		if err != nil {
			return triadFiles{}, err
		}
		*f.dst = namedFile{Name: filepath.Base(f.path), Data: data}
	}
	return t, nil
}
