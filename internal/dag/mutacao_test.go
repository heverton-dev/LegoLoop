package dag

import (
	"errors"
	"strings"
	"testing"
)

func TestMutacaoCronLimitesEFaixas(t *testing.T) {
	exato := "0" + strings.Repeat(" ", maxCronLen-8) + "0 1 1 0"
	if len(exato) != maxCronLen {
		t.Fatalf("tamanho esperado %d, obtido %d", maxCronLen, len(exato))
	}
	for _, expr := range []string{exato, "*/1 * * * *", "5-5 * * * *"} {
		if err := validateCron(expr); err != nil {
			t.Errorf("validateCron(%q) falhou: %v", expr, err)
		}
	}
	for _, expr := range []string{"0-abc * * * *", "0-99 * * * *"} {
		if err := validateCron(expr); err == nil {
			t.Errorf("validateCron(%q) deveria falhar", expr)
		}
	}
}

func TestMutacaoExecutorDeclaradoECoerencia(t *testing.T) {
	if executorDeclarado(map[string]any{"executor": nil}) {
		t.Error("executor nil não deveria contar como declarado")
	}
	if !executorDeclarado(map[string]any{"executor": 123}) {
		t.Error("executor não-string não-nil deveria contar como declarado")
	}
	if executorCoerente(NodeType("desconhecido"), ExecutorOrchestrixHermetic) {
		t.Error("tipo desconhecido não deveria ser coerente")
	}
	if got := executoresAceitos(TypeStopGate); got != ExecutorLocalProcess {
		t.Errorf("executoresAceitos(stop_gate) = %q, esperado %q", got, ExecutorLocalProcess)
	}
	if got := executoresAceitos(TypeScript); !strings.Contains(got, ExecutorOrchestrixHermetic) {
		t.Errorf("executoresAceitos(script) = %q, faltou %q", got, ExecutorOrchestrixHermetic)
	}
}

func TestMutacaoParseTamanhoLimiteEExecutor(t *testing.T) {
	base := `{"name":"lim","nodes":[{"id":"a","type":"script","command":["echo","ok"],"output":{"artifact":"out/a.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}]}`
	pad := make([]byte, MaxDAGBytes)
	copy(pad, base)
	for i := len(base); i < len(pad); i++ {
		pad[i] = ' '
	}
	if _, err := Parse(pad); err != nil {
		t.Fatalf("Parse com exatamente MaxDAGBytes falhou: %v", err)
	}
	invalido := `{"name":"ex","nodes":[{"id":"g","type":"stop_gate","isolation":{"executor":"orchestrix_hermetic"},"output":{"artifact":"out/g.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}]}`
	if _, err := Parse([]byte(invalido)); !errors.Is(err, ErrExecutorIncoerente) {
		t.Fatalf("Parse com executor incoerente deveria retornar ErrExecutorIncoerente, obtido: %v", err)
	}
	nulo := `{"name":"nl","nodes":[{"id":"a","type":"script","command":["echo","ok"],"isolation":{"executor":null},"output":{"artifact":"out/a.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}]}`
	if _, err := Parse([]byte(nulo)); err != nil {
		t.Fatalf("Parse com executor null deveria derivar padrão, erro: %v", err)
	}
	num := `{"name":"nm","nodes":[{"id":"a","type":"script","command":["echo","ok"],"isolation":{"executor":123},"output":{"artifact":"out/a.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}]}`
	if _, err := Parse([]byte(num)); !errors.Is(err, ErrExecutorIncoerente) {
		t.Fatalf("Parse com executor numérico deveria retornar ErrExecutorIncoerente, obtido: %v", err)
	}
}

func TestMutacaoParseCheckDepsNaoViraCiclo(t *testing.T) {
	auto := `{"name":"auto","nodes":[{"id":"a","type":"script","command":["echo"],"depends_on":["a"],"output":{"artifact":"out/a.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}]}`
	_, err := Parse([]byte(auto))
	if err == nil || !strings.Contains(err.Error(), "depende de si mesmo") {
		t.Fatalf("esperado erro de auto-dependência, obtido: %v", err)
	}
	inex := `{"name":"inex","nodes":[{"id":"a","type":"script","command":["echo"],"depends_on":["z"],"output":{"artifact":"out/a.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}]}`
	_, err = Parse([]byte(inex))
	if err == nil || !strings.Contains(err.Error(), "inexistente") {
		t.Fatalf("esperado erro de dependência inexistente, obtido: %v", err)
	}
}

func TestMutacaoParseDependenciasLegado(t *testing.T) {
	raw := `{"name":"dep","nodes":[
		{"id":"a","type":"script","command":["echo"],"dependencies":[],"output":{"artifact":"out/a.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}},
		{"id":"b","type":"script","command":["echo"],"dependencies":["a"],"output":{"artifact":"out/b.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}},
		{"id":"c","type":"script","command":["echo"],"depends_on":["a"],"dependencies":["b"],"output":{"artifact":"out/c.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}
	]}`
	d, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("Parse falhou: %v", err)
	}
	if d.Nodes[0].DependsOn != nil {
		t.Errorf("nó a deveria manter DependsOn nil, obtido %#v", d.Nodes[0].DependsOn)
	}
	if len(d.Nodes[1].DependsOn) != 1 || d.Nodes[1].DependsOn[0] != "a" {
		t.Errorf("nó b deveria copiar dependencies, obtido %#v", d.Nodes[1].DependsOn)
	}
	if len(d.Nodes[2].DependsOn) != 1 || d.Nodes[2].DependsOn[0] != "a" {
		t.Errorf("nó c deveria preservar depends_on, obtido %#v", d.Nodes[2].DependsOn)
	}
}

func TestMutacaoParseUnitNoComandoEPolitica(t *testing.T) {
	polInvalida := `{"name":"u1","nodes":[{"id":"u","type":"unit_no","command":["echo","ok"],"isolation":{"write_boundary":"open"},"output":{"artifact":"out/u.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}]}`
	if _, err := Parse([]byte(polInvalida)); err == nil {
		t.Error("unit_no com comando válido e write_boundary inválido deveria falhar")
	}
	cmdInvalido := `{"name":"u2","nodes":[{"id":"u","type":"unit_no","command":["./relativo.sh"],"output":{"artifact":"out/u.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}]}`
	if _, err := Parse([]byte(cmdInvalido)); err == nil {
		t.Error("unit_no com comando relativo deveria falhar")
	}
}

func TestMutacaoParseRolePapel(t *testing.T) {
	no := func(role, papel string) Node {
		return Node{Role: role, Papel: papel}
	}
	for _, tc := range []struct {
		nome string
		n    Node
		erro bool
	}{
		{"role_invalido", no("invalido_xyz", ""), true},
		{"papel_invalido", no("", "invalido_xyz"), true},
		{"so_role", no("arquiteto", ""), false},
		{"so_papel", no("", "arquiteto"), false},
		{"iguais", no("arquiteto", "arquiteto"), false},
		{"divergentes", no("arquiteto", "planejador"), true},
	} {
		err := validateRolePapel(tc.n)
		if (err != nil) != tc.erro {
			t.Errorf("%s: erro=%v, esperado=%v", tc.nome, err, tc.erro)
		}
	}
}

func TestMutacaoProvenanceTransitivaEContextOut(t *testing.T) {
	raw := `{"name":"prov","nodes":[
		{"id":"avo","type":"script","command":["echo"],"context_out":["gerado/ctx.txt"],"output":{"artifact":"out/avo.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}},
		{"id":"pai","type":"script","command":["echo"],"depends_on":["avo"],"context_in":["gerado/ctx.txt"],"output":{"artifact":"out/pai.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}},
		{"id":"filho","type":"script","command":["echo"],"depends_on":["pai"],"context_in":["out/avo.txt"],"output":{"artifact":"out/filho.txt","gate":"exit_code == 0"},"routing":{"mode":"manual","model":"m"}}
	]}`
	d, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("Parse falhou: %v", err)
	}
	if err := d.ValidateProvenance(t.TempDir()); err != nil {
		t.Fatalf("ValidateProvenance falhou: %v", err)
	}
	visited := make(map[string]bool)
	d.searchAncestors(&d.Nodes[2], "inexistente.txt", visited)
	if !visited["pai"] || !visited["avo"] {
		t.Errorf("searchAncestors deveria marcar ancestrais visitados, obtido %#v", visited)
	}
}

func TestMutacaoHasSlotsCadaCampo(t *testing.T) {
	casos := []struct {
		nome string
		s    *Slots
		want bool
	}{
		{"nil", nil, false},
		{"vazio", &Slots{}, false},
		{"so_gates", &Slots{Gates: []string{"exit_code == 0"}}, true},
		{"so_scripts", &Slots{Scripts: [][]string{{"echo"}}}, true},
		{"so_pre", &Slots{Hooks: Hooks{Pre: [][]string{{"echo"}}}}, true},
		{"so_post", &Slots{Hooks: Hooks{Post: [][]string{{"echo"}}}}, true},
	}
	for _, tc := range casos {
		if got := (Node{Slots: tc.s}).HasSlots(); got != tc.want {
			t.Errorf("%s: HasSlots() = %v, esperado %v", tc.nome, got, tc.want)
		}
	}
}

func TestMutacaoV2VersaoESlotsNoLimite(t *testing.T) {
	for _, v := range []string{"0", "3", "abc"} {
		if _, err := majorVersion(v); err == nil {
			t.Errorf("majorVersion(%q) deveria falhar", v)
		}
	}
	gates := make([]string, MaxSlotEntries)
	cmds := make([][]string, MaxSlotEntries)
	for i := range MaxSlotEntries {
		gates[i] = "exit_code == 0"
		cmds[i] = []string{"echo"}
	}
	if err := validateSlots(&Slots{Gates: gates, Scripts: cmds}); err != nil {
		t.Errorf("validateSlots no limite exato falhou: %v", err)
	}
}
