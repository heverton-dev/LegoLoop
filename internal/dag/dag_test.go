package dag

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

// specExample é o contrato literal da seção 3 do SPEC.md.
const specExample = `{
  "$schema": "https://legoloop.dev/schemas/dag.v1.json",
  "name": "greenfield-clean-craft",
  "version": "1.0.0",
  "nodes": [
    {
      "id": "node_01_interview",
      "type": "skill",
      "ref": "skills/interview",
      "routing": {"mode": "suggested", "model": "gemini-3.5-flash-lite", "fallback": "haiku"},
      "inputs": {"prompt": "var.idea"},
      "output": {"artifact": "BRIEFING.md", "gate": "file_exists('BRIEFING.md') && word_count('BRIEFING.md') > 100"}
    },
    {
      "id": "node_02_spec",
      "type": "skill",
      "ref": "skills/aidd-spec",
      "depends_on": ["node_01_interview"],
      "routing": {"mode": "auto", "combo": "code-pro"},
      "inputs": {"briefing": "file('BRIEFING.md')"},
      "output": {"artifact": "SPEC.md", "gate": "has_binary_criteria('SPEC.md')"}
    }
  ]
}`

func TestParseSpecExample(t *testing.T) {
	d, err := Parse([]byte(specExample))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if d.Name != "greenfield-clean-craft" || len(d.Nodes) != 2 {
		t.Fatalf("unexpected DAG: %+v", d)
	}
	n := d.Nodes[0]
	if n.Routing.Mode != ModeSuggested || n.Routing.Fallback != "haiku" || n.Output.Artifact != "BRIEFING.md" {
		t.Fatalf("node 0 mal decodificado: %+v", n)
	}
	if d.Nodes[1].Routing.Combo != "code-pro" || d.Nodes[1].DependsOn[0] != "node_01_interview" {
		t.Fatalf("node 1 mal decodificado: %+v", d.Nodes[1])
	}
	layers := d.Layers()
	want := [][]string{{"node_01_interview"}, {"node_02_spec"}}
	if !reflect.DeepEqual(layers, want) {
		t.Fatalf("Layers = %v, want %v", layers, want)
	}
	if d.Gate("node_02_spec") == nil || d.Node("node_02_spec") == nil || d.Node("x") != nil {
		t.Fatal("lookup de gate/nó falhou")
	}
}

func TestLayersParallelDiamond(t *testing.T) {
	src := dagJSON(
		node("a", ""),
		node("b", `"a"`),
		node("c", `"a"`),
		node("d", `"b","c"`),
	)
	d, err := Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"a"}, {"b", "c"}, {"d"}}
	if got := d.Layers(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Layers = %v, want %v", got, want)
	}
}

func TestCycleDetected(t *testing.T) {
	src := dagJSON(node("a", `"c"`), node("b", `"a"`), node("c", `"b"`), node("z", ""))
	_, err := Parse([]byte(src))
	var ce *CycleError
	if !errors.As(err, &ce) {
		t.Fatalf("esperado CycleError, recebido %v", err)
	}
	if !reflect.DeepEqual(ce.Nodes, []string{"a", "b", "c"}) {
		t.Fatalf("nós do ciclo = %v", ce.Nodes)
	}
}

func TestValidationErrors(t *testing.T) {
	cases := map[string]string{
		"json inválido":      `{`,
		"campo desconhecido": `{"name":"x","nodes":[],"evil":1}`,
		"sem nome":           dagJSONNamed("", node("a", "")),
		"sem nós":            `{"name":"x","nodes":[]}`,
		"id duplicado":       dagJSON(node("a", ""), node("a", "")),
		"id inválido":        dagJSON(node("a b", "")),
		"dep inexistente":    dagJSON(node("a", `"x"`)),
		"auto-dependência":   dagJSON(node("a", `"a"`)),
		"tipo inválido":      strings.Replace(dagJSON(node("a", "")), `"script"`, `"virus"`, 1),
		"modo inválido":      strings.Replace(dagJSON(node("a", "")), `"manual"`, `"yolo"`, 1),
		"artifact absoluto":  strings.Replace(dagJSON(node("a", "")), `"OUT_a.md"`, `"/etc/passwd"`, 1),
		"artifact traversal": strings.Replace(dagJSON(node("a", "")), `"OUT_a.md"`, `"../x.md"`, 1),
		"gate inválido":      strings.Replace(dagJSON(node("a", "")), `exit_code == 0`, `exit_code ==`, 1),
		"timeout negativo":   strings.Replace(dagJSON(node("a", "")), `"timeout_sec": 5`, `"timeout_sec": -1`, 1),
		"script sem comando": strings.Replace(dagJSON(node("a", "")), `"command": ["go","version"]`, `"command": []`, 1),
		"skill sem ref":      `{"name":"x","nodes":[{"id":"a","type":"skill","routing":{"mode":"manual","model":"m"},"output":{"artifact":"A.md","gate":"true"}}]}`,
		"manual sem modelo":  `{"name":"x","nodes":[{"id":"a","type":"skill","ref":"r","routing":{"mode":"manual"},"output":{"artifact":"A.md","gate":"true"}}]}`,
		"auto sem combo":     `{"name":"x","nodes":[{"id":"a","type":"skill","ref":"r","routing":{"mode":"auto"},"output":{"artifact":"A.md","gate":"true"}}]}`,
		"dois json":          dagJSON(node("a", "")) + `{}`,
	}
	for name, src := range cases {
		if _, err := Parse([]byte(src)); err == nil {
			t.Errorf("%s: esperado erro", name)
		}
	}
}

func TestParseRejectsOversizedInput(t *testing.T) {
	big := make([]byte, MaxDAGBytes+1)
	if _, err := Parse(big); err == nil {
		t.Fatal("esperado erro de tamanho")
	}
}

func TestMarshalRoundTrip(t *testing.T) {
	d, err := Parse([]byte(specExample))
	if err != nil {
		t.Fatal(err)
	}
	out, err := d.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	d2, err := Parse(out)
	if err != nil {
		t.Fatalf("round-trip: %v", err)
	}
	if !reflect.DeepEqual(d.Nodes, d2.Nodes) {
		t.Fatal("round-trip alterou os nós")
	}
}

func node(id, deps string) string {
	return `{"id":"` + id + `","type":"script","depends_on":[` + deps + `],
	"command": ["go","version"], "timeout_sec": 5,
	"routing":{"mode":"manual","model":"none"},
	"output":{"artifact":"OUT_` + id + `.md","gate":"exit_code == 0"}}`
}

func dagJSON(nodes ...string) string { return dagJSONNamed("t", nodes...) }

func dagJSONNamed(name string, nodes ...string) string {
	return `{"name":"` + name + `","version":"1.0.0","nodes":[` + strings.Join(nodes, ",") + `]}`
}

// SEC-01/03/07/11: validações de segurança do contrato.
func TestSecurityValidation(t *testing.T) {
	base := dagJSON(node("a", ""))
	cases := map[string]string{
		"hook git":            strings.Replace(base, `"OUT_a.md"`, `".git/hooks/pre-commit"`, 1),
		"settings claude":     strings.Replace(base, `"OUT_a.md"`, `".claude/settings.json"`, 1),
		"tasks vscode":        strings.Replace(base, `"OUT_a.md"`, `".VSCode/tasks.json"`, 1),
		"workflow github":     strings.Replace(base, `"OUT_a.md"`, `"sub/.github/workflows/x.yml"`, 1),
		"comando relativo":    strings.Replace(base, `["go","version"]`, `["tools/x.exe"]`, 1),
		"comando relativo \\": strings.Replace(base, `["go","version"]`, `["tools\\x.exe"]`, 1),
		"timeout absurdo":     strings.Replace(base, `"timeout_sec": 5`, `"timeout_sec": 9999999999`, 1),
		"nome curto 8.3":      strings.Replace(base, `"OUT_a.md"`, `"GIT~1/hooks/pre-commit"`, 1),
		"comando drive-rel":   strings.Replace(base, `["go","version"]`, `["C:x.exe"]`, 1),
		"artefato duplicado":  strings.Replace(dagJSON(node("a", ""), node("b", "")), `"OUT_b.md"`, `"OUT_a.md"`, 1),
	}
	for name, src := range cases {
		if _, err := Parse([]byte(src)); err == nil {
			t.Errorf("%s: deveria ser rejeitado", name)
		}
	}
	exe, _ := json.Marshal(os.Args[0]) // caminho absoluto válido no SO corrente
	abs := strings.Replace(base, `["go","version"]`, `[`+string(exe)+`,"version"]`, 1)
	if _, err := Parse([]byte(abs)); err != nil {
		t.Fatalf("comando absoluto deveria ser aceito: %v", err)
	}
}

func TestHandoffSchemaContract(t *testing.T) {
	valid := `{
		"name": "schema-test",
		"version": "1.0.0",
		"nodes": [{
			"id": "node_1",
			"type": "script",
			"command": ["go", "version"],
			"routing": {"mode": "suggested", "model": "go"},
			"output": {
				"artifact": "out.json",
				"schema": "schemas/output.schema.json",
				"gate": "exit_code == 0"
			}
		}]
	}`
	d, err := Parse([]byte(valid))
	if err != nil {
		t.Fatalf("Parse com output.schema válido falhou: %v", err)
	}
	if d.Nodes[0].Output.Schema != "schemas/output.schema.json" {
		t.Fatalf("schema não decodificado corretamente: %q", d.Nodes[0].Output.Schema)
	}

	invalidPaths := []string{
		"../outside.json",
		"C:/schema.json",
		".git/schema.json",
		"GIT~1/schema.json",
	}
	for _, p := range invalidPaths {
		bad := strings.Replace(valid, "schemas/output.schema.json", p, 1)
		if _, err := Parse([]byte(bad)); err == nil {
			t.Errorf("deveria rejeitar schema inválido %q", p)
		}
	}
}

func TestRoutingEffortValidation(t *testing.T) {
	template := `{
		"name": "effort-test",
		"version": "1.0.0",
		"nodes": [{
			"id": "node_1",
			"type": "skill",
			"ref": "skills/interview",
			"routing": {"mode": "suggested", "model": "flash", "effort": "EFFORT_VALUE"},
			"output": {"artifact": "out.md", "gate": "exit_code == 0"}
		}]
	}`

	validEfforts := []EffortLevel{EffortLow, EffortMedium, EffortHigh}
	for _, eff := range validEfforts {
		src := strings.Replace(template, "EFFORT_VALUE", string(eff), 1)
		d, err := Parse([]byte(src))
		if err != nil {
			t.Fatalf("Parse com effort %q falhou: %v", eff, err)
		}
		if d.Nodes[0].Routing.Effort != eff {
			t.Fatalf("effort não decodificado corretamente: want %q, got %q", eff, d.Nodes[0].Routing.Effort)
		}
	}

	invalidEfforts := []string{"extreme", "zero", "low_high", "123"}
	for _, bad := range invalidEfforts {
		src := strings.Replace(template, "EFFORT_VALUE", bad, 1)
		if _, err := Parse([]byte(src)); err == nil {
			t.Errorf("deveria rejeitar effort inválido %q", bad)
		}
	}
}

func TestParseStopGateNode(t *testing.T) {
	src := `{
		"name": "stopgate-test",
		"version": "1.0.0",
		"nodes": [
			{
				"id": "step_01",
				"type": "script",
				"command": ["go", "test", "./..."],
				"routing": {"mode": "suggested", "model": "code-fast"},
				"output": {"artifact": "res.txt", "gate": "exit_code == 0"}
			},
			{
				"id": "sg_approval",
				"type": "stop_gate",
				"depends_on": ["step_01"],
				"routing": {"mode": "suggested", "model": "code-pro"},
				"output": {"artifact": "res_approved.txt", "gate": "file_exists('res.txt') && prompt('Aprovar etapa?')"}
			}
		]
	}`
	d, err := Parse([]byte(src))
	if err != nil {
		t.Fatalf("Parse com stop_gate falhou: %v", err)
	}
	if len(d.Nodes) != 2 || d.Nodes[1].Type != TypeStopGate {
		t.Fatalf("nó stop_gate não decodificado: %+v", d.Nodes[1])
	}
}
