package dag

import "testing"

const noScript = `{"id":"a","type":"script","command":["go"],"routing":{"mode":"manual","model":"m"},"output":{"artifact":"a.json","gate":"file_exists('a.json')"}}`

func TestLeitorDoisDialetos(t *testing.T) {
	sub := func(tipo, ref string) string {
		return `{"id":"b","type":"` + tipo + `","ref":"` + ref + `","depends_on":["a"],"routing":{"mode":"manual","model":"m"},"output":{"artifact":"b.json","gate":"file_exists('b.json')"}}`
	}
	casos := []struct {
		nome, src string
		ok        bool
	}{
		{"nó subprocesso novo", v2(`{"type":"manual"}`, noScript, sub("subprocesso", "outro")), true},
		{"subprocesso ref inválida", v2(`{"type":"manual"}`, noScript, sub("subprocesso", "../x")), false},
		{"trigger processo_fim novo", v2(`{"type":"processo_fim","processo":"alvo"}`, noScript), true},
		{"processo_fim com o campo antigo", v2(`{"type":"processo_fim","loop":"alvo"}`, noScript), true},
		{"gatilho antigo com campo processo", v2(`{"type":"loop_end","processo":"alvo"}`, noScript), true},
		{"campos divergentes", v2(`{"type":"processo_fim","loop":"a","processo":"b"}`, noScript), false},
		{"processo_fim sem alvo", v2(`{"type":"processo_fim"}`, noScript), false},
		{"processo inválido", v2(`{"type":"processo_fim","processo":"a b"}`, noScript), false},
		{"processo em trigger manual", v2(`{"type":"manual","processo":"x"}`, noScript), false},
	}
	for _, c := range casos {
		_, err := Parse([]byte(c.src))
		if (err == nil) != c.ok {
			t.Errorf("%s: ok=%v erro=%v", c.nome, c.ok, err)
		}
	}
}

func TestLeitorNormalizaParaRepresentacaoAntiga(t *testing.T) {
	d, err := Parse([]byte(v2(`{"type":"processo_fim","processo":"alvo"}`, noScript)))
	if err != nil {
		t.Fatal(err)
	}
	if d.Trigger.Type != gatilhoFimInterno || *d.Trigger.alvo() != "alvo" || d.Trigger.Processo != "" {
		t.Errorf("trigger não normalizado: %+v", d.Trigger)
	}
}

func TestIDDeNoComPonto(t *testing.T) {
	no := func(id string) string {
		return `{"id":"` + id + `","type":"script","command":["go"],"routing":{"mode":"manual","model":"m"},"output":{"artifact":"o.json","gate":"file_exists('o.json')"}}`
	}
	validos := []string{"TICK-01.2", "a.b", "e04", "x_y-z", "a.b.c"}
	invalidos := []string{".a", "a.", "..", ".", "a b", "a/b", "a\b", ""}
	for _, id := range validos {
		if _, err := Parse([]byte(v2(`{"type":"manual"}`, no(id)))); err != nil {
			t.Errorf("id %q deveria valer: %v", id, err)
		}
	}
	for _, id := range invalidos {
		if _, err := Parse([]byte(v2(`{"type":"manual"}`, no(id)))); err == nil {
			t.Errorf("id %q deveria falhar", id)
		}
	}
}
