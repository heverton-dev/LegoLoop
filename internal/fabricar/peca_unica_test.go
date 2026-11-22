package fabricar_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/fabricar"
	"github.com/heverton-dev/LegoLoop/internal/processo"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

func TestRequestValidate_PecaTipos(t *testing.T) {
	for _, tipo := range []string{"", "script", "skill", "gate", "hook", "mcp", "processo"} {
		r := validRequest()
		r.Peca = tipo
		if err := r.Validate(); err != nil {
			t.Errorf("tipo %q deve ser válido: %v", tipo, err)
		}
	}
	r := validRequest()
	r.Peca = "banana"
	if err := r.Validate(); !errors.Is(err, fabricar.ErrInvalidRequest) {
		t.Errorf("tipo desconhecido: erro = %v; quer ErrInvalidRequest", err)
	}
}

func TestFabricar_PecaUnicaGeraUmNoSemTriagem(t *testing.T) {
	r := newRig(t, nil) // matriz vazia: a peça única não depende de peça apta
	r.req.Peca = "script"
	res := r.run(t)
	a := res.Artifacts
	for nome, p := range map[string]string{"skill": a.SkillPath, "processo": a.ProcessoPath, "dag": a.DagPath,
		"intent": a.IntentPath, "spec": a.SpecPath, "plan": a.PlanPath, "route": a.RoutePath, "dag_iso": a.DagIsoPath, "swarm": a.SwarmPlanPath} {
		if p == "" {
			t.Fatalf("caminho de %s vazio: %+v", nome, a)
		}
		if _, err := os.Stat(p); err != nil {
			t.Errorf("artefato ausente (%s): %v", nome, err)
		}
	}
	dir := filepath.Dir(a.DagPath)
	if _, err := os.Stat(filepath.Join(dir, "triagem.json")); err == nil {
		t.Error("peça única não gera TRIAGEM")
	}
	if _, err := os.Stat(filepath.Join(dir, triagem.ArquivoMemoria)); err != nil {
		t.Errorf("a memória deve existir: %v", err)
	}
	if r.llm.needsCalls != 0 {
		t.Errorf("peça única não consulta a Fusão: %d chamadas", r.llm.needsCalls)
	}
}

func TestFabricar_PecaUnicaTemUmaEtapaComAPropriaPeca(t *testing.T) {
	r := newRig(t, nil)
	r.req.Peca = "gate"
	res := r.run(t)
	var p processo.Processo
	if err := json.Unmarshal([]byte(read(t, res.Artifacts.ProcessoPath)), &p); err != nil {
		t.Fatal(err)
	}
	if len(p.Etapas) != 1 || p.Etapas[0].Peca != res.CanonicalName || p.Etapas[0].Tipo != "gate" {
		t.Fatalf("esperava 1 etapa do tipo gate com a própria peça: %+v", p.Etapas)
	}
}

func TestFabricar_PecaUnicaSwarmPlanTemContratoEstrito(t *testing.T) {
	r := newRig(t, nil)
	r.req.Peca = "skill"
	res := r.run(t)
	var sp struct {
		Peca          string   `json:"peca"`
		Tipo          string   `json:"tipo"`
		Workers       int      `json:"workers"`
		WriteBoundary string   `json:"write_boundary"`
		Ferramentas   []string `json:"ferramentas"`
		Ticket        string   `json:"ticket"`
	}
	if err := json.Unmarshal([]byte(read(t, res.Artifacts.SwarmPlanPath)), &sp); err != nil {
		t.Fatal(err)
	}
	if sp.Peca != res.CanonicalName || sp.Tipo != "skill" || sp.Workers != 1 || sp.WriteBoundary != "strict" ||
		len(sp.Ferramentas) == 0 || sp.Ticket == "" {
		t.Fatalf("swarm.plan.json fora do contrato: %+v", sp)
	}
}

func TestFabricar_PecaUnicaNaoExigeAMatriz(t *testing.T) {
	r := newRig(t, nil)
	r.deps.Matriz = func() (fabricar.FonteMatriz, error) { return fabricar.FonteMatriz{}, errors.New("matriz ausente") }
	r.req.Peca = "script"
	r.run(t)
	r.req.Peca, r.req.Force = "", true
	if _, err := fabricar.Fabricar(context.Background(), r.req, r.deps); err == nil {
		t.Error("o processo completo continua exigindo a matriz")
	}
}
