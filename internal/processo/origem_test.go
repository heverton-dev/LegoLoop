package processo_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/processo"
)

// composta devolve quatro() com a etapa e03 executando outro processo, que grava out/e03.md dentro dele.
func composta(origem string) processo.Processo {
	p := quatro()
	e := &p.Etapas[2]
	e.Peca, e.Compoe, e.Origem = "", "legoloop-outro@1.2.0", origem
	e.Produz.Caminho, e.Gate = "out/legoloop-outro/e03.md", "file_exists('out/legoloop-outro/e03.md')"
	return p
}

func TestOrigemViraSubArtifactDoNo(t *testing.T) {
	pr := projetar(t, composta("out/e03.md"))
	d, err := dag.Parse(pr.DAG)
	if err != nil {
		t.Fatal(err)
	}
	n := d.Nodes[2]
	if n.Type != dag.TypeSubloop || n.Output.Artifact != "out/legoloop-outro/e03.md" || n.Output.SubArtifact != "out/e03.md" {
		t.Errorf("nó composto: %+v", n.Output)
	}
	if got := d.Nodes[3].Inputs["e03"]; got != "file('out/legoloop-outro/e03.md')" {
		t.Errorf("a etapa seguinte lê o caminho do pai: %q", got)
	}
	if !strings.Contains(string(pr.Loop), `"sub_artifact": "out/e03.md"`) {
		t.Errorf("o documento do Studio leva o sub_artifact:\n%s", pr.Loop)
	}
}

func TestOrigemVaziaMantemOContratoAntigo(t *testing.T) {
	p := composta("")
	pr := projetar(t, p)
	pj, err := p.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(pr.DAG), "sub_artifact") || strings.Contains(string(pr.Loop), "sub_artifact") ||
		strings.Contains(string(pj), "origem") {
		t.Errorf("sem origem, nada novo nos arquivos:\n%s", pr.DAG)
	}
	var volta processo.Processo
	q := composta("out/e03.md")
	b, err := q.Marshal()
	if err != nil || json.Unmarshal(b, &volta) != nil || volta.Etapas[2].Origem != "out/e03.md" {
		t.Errorf("a origem vai e volta no processo.json (%v): %s", err, b)
	}
}

// TestOrigemRecusada: a origem só vale na etapa composta e com caminho relativo dentro do processo composto.
func TestOrigemRecusada(t *testing.T) {
	for nome, c := range map[string]struct {
		p    processo.Processo
		quer string
	}{
		"etapa com peça": {func() processo.Processo { p := quatro(); p.Etapas[1].Origem = "out/x.md"; return p }(), "origem"},
		"absoluta":       {composta("/etc/passwd"), "origem"},
		"fora":           {composta("../x.md"), "origem"},
		"com espaço":     {composta("out/a b.md"), "origem"},
	} {
		if _, err := processo.Projetar(c.p, agora); !errors.Is(err, processo.ErrInvalido) || !strings.Contains(err.Error(), c.quer) {
			t.Errorf("%s: esperava ErrInvalido citando %q, veio %v", nome, c.quer, err)
		}
		if err := c.p.Validate(); nome != "com espaço" && err == nil {
			t.Errorf("%s: Validate deveria recusar", nome)
		}
	}
}
