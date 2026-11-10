package dag

import (
	"encoding/json"
	"strings"
	"testing"
)

// subDAG é um DAG pai com um nó subprocesso que devolve o artefato do subprocesso num caminho próprio do pai.
const subDAG = `{
  "name": "pai",
  "version": "1.1.0",
  "nodes": [
    {
      "id": "e01", "type": "skill", "ref": "skills/x",
      "routing": {"mode": "manual", "model": "m"},
      "output": {"artifact": "out/02_verificar.md", "gate": "file_exists('out/02_verificar.md')"}
    },
    {
      "id": "e02", "type": "subprocesso", "ref": "legoloop-b", "depends_on": ["e01"],
      "routing": {"mode": "manual", "model": "m"},
      "output": {"artifact": "out/legoloop-b/02_verificar.md", "sub_artifact": "out/02_verificar.md",
        "gate": "file_exists('out/legoloop-b/02_verificar.md')"}
    }
  ]
}`

// TestSubArtifactSeparaOCaminhoDoSubprocessoDoCaminhoNoPai: o nó subprocesso lê o artefato no
// subprocesso por sub_artifact e grava no pai por artifact; os dois podem ter o mesmo nome de
// arquivo sem colidir.
func TestSubArtifactSeparaOCaminhoDoSubprocessoDoCaminhoNoPai(t *testing.T) {
	d, err := Parse([]byte(subDAG))
	if err != nil {
		t.Fatalf("DAG com sub_artifact deveria valer: %v", err)
	}
	n := d.Nodes[1]
	if n.Output.Artifact != "out/legoloop-b/02_verificar.md" || n.Output.SubArtifact != "out/02_verificar.md" {
		t.Errorf("contrato lido: %+v", n.Output)
	}
	b, err := json.Marshal(d)
	if err != nil || !strings.Contains(string(b), `"sub_artifact":"out/02_verificar.md"`) {
		t.Errorf("o sub_artifact deve voltar no JSON (%v): %s", err, b)
	}
}

// TestSubArtifactOmitidoNoDAGAntigo: o DAG sem sub_artifact continua igual na ida e na volta.
func TestSubArtifactOmitidoNoDAGAntigo(t *testing.T) {
	antigo := strings.Replace(subDAG, `"sub_artifact": "out/02_verificar.md",`, "", 1)
	antigo = strings.Replace(antigo, "out/legoloop-b/02_verificar.md", "out/b.md", 2)
	d, err := Parse([]byte(antigo))
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(d)
	if d.Nodes[1].Output.SubArtifact != "" || strings.Contains(string(b), "sub_artifact") {
		t.Errorf("sem sub_artifact, nada muda no JSON: %s", b)
	}
}

// TestSubArtifactRecusado: só nó subprocesso usa sub_artifact, e o caminho segue a regra do artifact.
func TestSubArtifactRecusado(t *testing.T) {
	casos := map[string]string{
		"nó skill":       strings.Replace(subDAG, `"output": {"artifact": "out/02_verificar.md",`, `"output": {"sub_artifact": "x.md", "artifact": "out/02_verificar.md",`, 1),
		"absoluto":       strings.Replace(subDAG, `"sub_artifact": "out/02_verificar.md"`, `"sub_artifact": "/etc/passwd"`, 1),
		"fora do pai":    strings.Replace(subDAG, `"sub_artifact": "out/02_verificar.md"`, `"sub_artifact": "../x.md"`, 1),
		"letra de disco": strings.Replace(subDAG, `"sub_artifact": "out/02_verificar.md"`, `"sub_artifact": "C:/x.md"`, 1),
		"diretório .git": strings.Replace(subDAG, `"sub_artifact": "out/02_verificar.md"`, `"sub_artifact": ".git/config"`, 1),
	}
	for nome, src := range casos {
		if src == subDAG {
			t.Fatalf("%s: a troca não aconteceu", nome)
		}
		_, err := Parse([]byte(src))
		if err == nil || !strings.Contains(err.Error(), "sub_artifact") {
			t.Errorf("%s: esperava erro citando sub_artifact, veio %v", nome, err)
		}
	}
}
