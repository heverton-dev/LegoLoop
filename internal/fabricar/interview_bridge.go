package fabricar

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"legoloop/internal/caminhos"
	"legoloop/internal/dag"
	"legoloop/internal/interview"
)

// FabricarFromInterview consome a sessão e a INTENT (UNDER_REVIEW, confiança >= 95%)
// produzidas por internal/interview e as conecta como context_in da Curadoria e Composição.
func FabricarFromInterview(ctx context.Context, req Request, deps Deps, sess *interview.Session, intentPath string) (*Result, error) {
	if sess == nil || sess.Confidence < 95 {
		return nil, ErrLowConfidence
	}
	if strings.TrimSpace(intentPath) != "" {
		data, err := os.ReadFile(intentPath)
		if err != nil {
			return nil, fmt.Errorf("fabricar: ler INTENT de context_in: %w", err)
		}
		if len(strings.TrimSpace(string(data))) == 0 {
			return nil, fmt.Errorf("fabricar: INTENT vazia em %s", intentPath)
		}
	}
	clar := ClarificationFromSession(sess)
	clar.ContextIn = intentPath
	deps.Interview = &clar
	return Fabricar(ctx, req, deps)
}

// FabricarFromInterviewWithStopGate executa o fluxo pós-entrevista exigindo parada
// obrigatória no Stop-Gate após a cascata de 3 camadas (Cofre -> SO -> GitHub)
// antes da materialização quando req.Yes é falso.
func FabricarFromInterviewWithStopGate(ctx context.Context, req Request, deps Deps, sess *interview.Session, intentPath string, sugestoesGitHub []string) (*Result, CuradoriaStopGateSummary, error) {
	if sess == nil || sess.Confidence < 95 {
		return nil, CuradoriaStopGateSummary{}, ErrLowConfidence
	}
	resumo := fmt.Sprintf("processo=%s tema=%s", req.RawName, sess.Theme)
	summary, err := AvaliarStopGateCuradoria(req, resumo, sugestoesGitHub)
	if err != nil {
		return nil, summary, err
	}
	res, err := FabricarFromInterview(ctx, req, deps, sess, intentPath)
	return res, summary, err
}

func injectIntentFiles(files []PlannedFile, intentPath string) error {
	if strings.TrimSpace(intentPath) == "" {
		return nil
	}
	for i := range files {
		if strings.HasSuffix(files[i].Path, caminhos.ArquivoDAG) {
			updated, err := injectIntentContextIn(files[i].Data, intentPath)
			if err != nil {
				return err
			}
			files[i].Data = updated
		}
	}
	return nil
}

func injectIntentContextIn(dagBytes []byte, intentPath string) ([]byte, error) {
	if strings.TrimSpace(intentPath) == "" {
		return dagBytes, nil
	}
	d, err := dag.Parse(dagBytes)
	if err != nil {
		return nil, err
	}
	for i := range d.Nodes {
		if len(d.Nodes[i].DependsOn) == 0 {
			d.Nodes[i].ContextIn = append([]string{intentPath}, d.Nodes[i].ContextIn...)
		}
	}
	out, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}
	out = append(out, '\n')
	if _, err := dag.Parse(out); err != nil {
		return nil, err
	}
	return out, nil
}
