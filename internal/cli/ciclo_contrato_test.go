package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
)

type skillContrato struct {
	nome     string
	consumes string
	produces string
}

func encontrarRaizRepo(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("obter wd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			t.Fatal("go.mod não encontrado")
		}
		dir = pai
	}
}

func extrairContratoSkill(t *testing.T, caminho string) skillContrato {
	t.Helper()
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler skill %s: %v", caminho, err)
	}
	sc := bufio.NewScanner(bytes.NewReader(dados))
	dentro := false
	var res skillContrato
	for sc.Scan() {
		linha := strings.TrimSpace(sc.Text())
		if linha == "---" {
			if dentro {
				break
			}
			dentro = true
			continue
		}
		if dentro {
			preencherCampoContrato(&res, linha)
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scanner erro em %s: %v", caminho, err)
	}
	return res
}

func preencherCampoContrato(res *skillContrato, linha string) {
	chave, valor, ok := strings.Cut(linha, ":")
	if !ok {
		return
	}
	v := strings.Trim(strings.TrimSpace(valor), "[]\"'")
	switch strings.TrimSpace(chave) {
	case "name":
		res.nome = v
	case "consumes":
		res.consumes = v
	case "produces":
		res.produces = v
	}
}

type dagProcesso struct {
	Nodes []struct {
		ID         string   `json:"id"`
		Ref        string   `json:"ref"`
		DependsOn  []string `json:"depends_on"`
		ContextIn  []string `json:"context_in"`
		ContextOut []string `json:"context_out"`
		Output     struct {
			Artifact string `json:"artifact"`
			Gate     string `json:"gate"`
		} `json:"output"`
	} `json:"nodes"`
}

func validarOrdemNosDAG(t *testing.T, raiz string) {
	t.Helper()
	dagPath := filepath.Join(raiz, "processos", "legoloop-ciclo", caminhos.ArquivoDAG)
	dados, err := os.ReadFile(dagPath)
	if err != nil {
		t.Fatalf("ler dag legoloop-ciclo: %v", err)
	}
	var dp dagProcesso
	if err := json.Unmarshal(dados, &dp); err != nil {
		t.Fatalf("unmarshal dag legoloop-ciclo: %v", err)
	}
	if len(dp.Nodes) != 3 {
		t.Fatalf("esperava 3 nós no dag legoloop-ciclo, obteve %d", len(dp.Nodes))
	}
	verificarNosEncadeados(t, dp)
}

func verificarNosEncadeados(t *testing.T, dp dagProcesso) {
	t.Helper()
	ids := []string{dp.Nodes[0].ID, dp.Nodes[1].ID, dp.Nodes[2].ID}
	esperados := []string{"concepcao", "especificacao", "planejamento"}
	if !slices.Equal(ids, esperados) {
		t.Fatalf("ordem dos nós incorreta: got %v, want %v", ids, esperados)
	}
	if len(dp.Nodes[0].DependsOn) != 0 {
		t.Errorf("concepcao não deve ter depends_on: %v", dp.Nodes[0].DependsOn)
	}
	if !slices.Contains(dp.Nodes[1].DependsOn, "concepcao") {
		t.Errorf("especificacao deve depender de concepcao: %v", dp.Nodes[1].DependsOn)
	}
	if !slices.Contains(dp.Nodes[2].DependsOn, "especificacao") {
		t.Errorf("planejamento deve depender de especificacao: %v", dp.Nodes[2].DependsOn)
	}
	verificarArtefatosEncadeados(t, dp)
}

func verificarArtefatosEncadeados(t *testing.T, dp dagProcesso) {
	t.Helper()
	artIntent := "docs/intents/INTENT.md"
	artSpec := "docs/specs/SPEC.md"
	artPlan := "docs/plans/PLAN.md"

	if dp.Nodes[0].Output.Artifact != artIntent {
		t.Errorf("concepcao output %q != %q", dp.Nodes[0].Output.Artifact, artIntent)
	}
	if !slices.Contains(dp.Nodes[1].ContextIn, artIntent) {
		t.Errorf("especificacao context_in não contém %q", artIntent)
	}
	if dp.Nodes[1].Output.Artifact != artSpec {
		t.Errorf("especificacao output %q != %q", dp.Nodes[1].Output.Artifact, artSpec)
	}
	if !slices.Contains(dp.Nodes[2].ContextIn, artSpec) {
		t.Errorf("planejamento context_in não contém %q", artSpec)
	}
	if dp.Nodes[2].Output.Artifact != artPlan {
		t.Errorf("planejamento output %q != %q", dp.Nodes[2].Output.Artifact, artPlan)
	}
}

func testarValidateCiclo(t *testing.T, raiz string) {
	t.Helper()
	dagPath := filepath.Join(raiz, "processos", "legoloop-ciclo", caminhos.ArquivoDAG)
	var stdout, stderr bytes.Buffer
	d := Deps{Stdout: &stdout, Stderr: &stderr}
	code := cmdValidate(context.Background(), []string{dagPath}, d)
	if code != 0 {
		t.Fatalf("cmdValidate falhou com exit %d: %s", code, stderr.String())
	}
	saida := stdout.String()
	camadas := []string{
		"camada 1: concepcao",
		"camada 2: especificacao",
		"camada 3: planejamento",
	}
	for _, cam := range camadas {
		if !strings.Contains(saida, cam) {
			t.Errorf("saída do validate não contém %q: %s", cam, saida)
		}
	}
}

func TestCicloContratoEncadeia(t *testing.T) {
	raiz := encontrarRaizRepo(t)
	cConcepcao := extrairContratoSkill(t, filepath.Join(raiz, "skills", "legoloop-concepcao", "SKILL.md"))
	cEspec := extrairContratoSkill(t, filepath.Join(raiz, "skills", "legoloop-especificacao", "SKILL.md"))
	cPlano := extrairContratoSkill(t, filepath.Join(raiz, "skills", "legoloop-planejamento", "SKILL.md"))

	if cConcepcao.produces != "Artifact:Intent" {
		t.Errorf("concepcao produces inválido: got %q, want Artifact:Intent", cConcepcao.produces)
	}
	if cEspec.consumes != cConcepcao.produces {
		t.Errorf("quebra de contrato concepcao -> especificacao: produces %q != consumes %q", cConcepcao.produces, cEspec.consumes)
	}
	if cEspec.produces != "Artifact:Spec" {
		t.Errorf("especificacao produces inválido: got %q, want Artifact:Spec", cEspec.produces)
	}
	if cPlano.consumes != cEspec.produces {
		t.Errorf("quebra de contrato especificacao -> planejamento: produces %q != consumes %q", cEspec.produces, cPlano.consumes)
	}
	if cPlano.produces != "Artifact:Plan" {
		t.Errorf("planejamento produces inválido: got %q, want Artifact:Plan", cPlano.produces)
	}

	validarOrdemNosDAG(t, raiz)
	testarValidateCiclo(t, raiz)
}
