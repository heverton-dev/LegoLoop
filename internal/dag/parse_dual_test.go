package dag

import (
	"errors"
	"reflect"
	"testing"
)

// dualNos monta um DAG com um nó com ponto no id e um subprocesso/subloop dependente.
func dualNos(trigger, tipo, ref string) string {
	a := `{"id":"TICK-01.1","type":"script","command":["go"],"routing":{"mode":"manual","model":"m"},"output":{"artifact":"a.json","gate":"file_exists('a.json')"}}`
	b := `{"id":"TICK-01.2","type":"` + tipo + `","ref":"` + ref + `","depends_on":["TICK-01.1"],"routing":{"mode":"manual","model":"m"},"output":{"artifact":"b.json","gate":"file_exists('b.json')"}}`
	return v2(trigger, a, b)
}

func parseOuFalha(t *testing.T, src string) *DAG {
	t.Helper()
	d, err := Parse([]byte(src))
	if err != nil {
		t.Fatalf("parse falhou: %v", err)
	}
	return d
}

func TestParseDoisDialetos(t *testing.T) {
	antigo := parseOuFalha(t, dualNos(`{"type":"loop_end","loop":"alvo"}`, "subloop", "outro"))
	novo := parseOuFalha(t, dualNos(`{"type":"processo_fim","processo":"alvo"}`, "subprocesso", "outro"))
	if !reflect.DeepEqual(antigo.Nodes, novo.Nodes) {
		t.Errorf("nós divergem:\nantigo=%+v\nnovo=%+v", antigo.Nodes, novo.Nodes)
	}
	if !reflect.DeepEqual(antigo.Trigger, novo.Trigger) {
		t.Errorf("trigger diverge: antigo=%+v novo=%+v", antigo.Trigger, novo.Trigger)
	}
	if !reflect.DeepEqual(antigo.Layers(), novo.Layers()) {
		t.Errorf("camadas divergem: antigo=%v novo=%v", antigo.Layers(), novo.Layers())
	}
	if antigo.Node("TICK-01.1") == nil {
		t.Error("nó com ponto no id deveria ser aceito")
	}
	misto := dualNos(`{"type":"processo_fim","loop":"a","processo":"b"}`, "subprocesso", "outro")
	if _, err := Parse([]byte(misto)); !errors.Is(err, ErrDialetoMisto) {
		t.Errorf("mistura inválida deveria retornar ErrDialetoMisto, recebeu %v", err)
	}
	for _, tipo := range []string{"subprocesso", "subloop"} {
		if _, err := Parse([]byte(dualNos(`{"type":"manual"}`, tipo, "a.b"))); err == nil {
			t.Errorf("ref %s com ponto deveria ser recusado", tipo)
		}
	}
}
