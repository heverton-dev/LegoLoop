package gate

import (
	"strings"
	"testing"
)

// Caracterização de ValidateJSONSchema: fixa o comportamento atual de cada ramo antes de dividir matchSchema
// (gate funclen do TICK-02.9). Se um caso falhar depois da refatoração, ela mudou o comportamento.
func TestGateTrivialCaracterizaSchema(t *testing.T) {
	for _, c := range []struct {
		nome, schema, doc string
		want              bool
	}{
		{"sem restrições", `{}`, `{"a":1}`, true},
		{"tipo desconhecido não restringe", `{"type":"banana"}`, `1`, true},
		{"object ok", `{"type":"object"}`, `{}`, true},
		{"object recusa lista", `{"type":"object"}`, `[]`, false},
		{"array ok", `{"type":"array"}`, `[]`, true},
		{"array recusa objeto", `{"type":"array"}`, `{}`, false},
		{"string ok", `{"type":"string"}`, `"x"`, true},
		{"string recusa número", `{"type":"string"}`, `1`, false},
		{"number ok", `{"type":"number"}`, `1.5`, true},
		{"number recusa texto", `{"type":"number"}`, `"1"`, false},
		{"integer ok", `{"type":"integer"}`, `3`, true},
		{"integer aceita 3.0", `{"type":"integer"}`, `3.0`, true},
		{"integer recusa fração", `{"type":"integer"}`, `3.5`, false},
		{"integer recusa texto", `{"type":"integer"}`, `"3"`, false},
		{"boolean ok", `{"type":"boolean"}`, `true`, true},
		{"boolean recusa número", `{"type":"boolean"}`, `1`, false},
		{"null ok", `{"type":"null"}`, `null`, true},
		{"null recusa valor", `{"type":"null"}`, `0`, false},
		{"type que não é texto é ignorado", `{"type":7}`, `1`, true},
		{"required ok", `{"required":["a","b"]}`, `{"a":1,"b":null}`, true},
		{"required falta chave", `{"required":["a","b"]}`, `{"a":1}`, false},
		{"required em não-objeto", `{"required":["a"]}`, `[1]`, false},
		{"required com item que não é texto", `{"required":[1]}`, `{"a":1}`, false},
		{"required vazio", `{"required":[]}`, `{}`, true},
		{"properties ok", `{"properties":{"a":{"type":"number"}}}`, `{"a":1}`, true},
		{"properties tipo errado", `{"properties":{"a":{"type":"number"}}}`, `{"a":"x"}`, false},
		{"properties ausente é opcional", `{"properties":{"a":{"type":"number"}}}`, `{}`, true},
		{"properties em não-objeto", `{"properties":{"a":{}}}`, `[1]`, false},
		{"properties com definição que não é objeto", `{"properties":{"a":5}}`, `{"a":"x"}`, true},
		{"properties aninhadas", `{"properties":{"a":{"properties":{"b":{"type":"string"}}}}}`, `{"a":{"b":1}}`, false},
		{"items ok", `{"items":{"type":"number"}}`, `[1,2]`, true},
		{"items item errado", `{"items":{"type":"number"}}`, `[1,"x"]`, false},
		{"items em não-lista", `{"items":{"type":"number"}}`, `{}`, false},
		{"items lista vazia", `{"items":{"type":"number"}}`, `[]`, true},
		{"items definição que não é objeto", `{"items":5}`, `[1]`, true},
		{"objeto completo", `{"type":"object","required":["itens"],"properties":{"itens":{"type":"array","items":{"type":"object","required":["preco"],"properties":{"preco":{"type":"number"}}}}}}`,
			`{"itens":[{"preco":1},{"preco":2.5}]}`, true},
		{"objeto completo com item ruim", `{"type":"object","required":["itens"],"properties":{"itens":{"type":"array","items":{"type":"object","required":["preco"]}}}}`,
			`{"itens":[{"preco":1},{}]}`, false},
	} {
		if got := ValidateJSONSchema([]byte(c.schema), []byte(c.doc)); got != c.want {
			t.Errorf("%s: ValidateJSONSchema(%s, %s) = %v, queria %v", c.nome, c.schema, c.doc, got, c.want)
		}
	}
}

func TestGateTrivialCaracterizaSchemaEntradasInvalidas(t *testing.T) {
	if ValidateJSONSchema([]byte(`{quebrado`), []byte(`{}`)) {
		t.Error("esquema que não é JSON deve falhar")
	}
	if ValidateJSONSchema([]byte(`{}`), []byte(`{quebrado`)) {
		t.Error("alvo que não é JSON deve falhar")
	}
	if ValidateJSONSchema([]byte(``), []byte(`{}`)) {
		t.Error("esquema vazio deve falhar")
	}
}

func TestGateTrivialCaracterizaSchemaProfundidade(t *testing.T) {
	aninha := func(n int) (schema, doc string) {
		schema, doc = `{"type":"number"}`, `1`
		for i := 0; i < n; i++ {
			schema = `{"properties":{"a":` + schema + `}}`
			doc = `{"a":` + doc + `}`
		}
		return schema, doc
	}
	s, d := aninha(MaxDepth)
	if !ValidateJSONSchema([]byte(s), []byte(d)) {
		t.Errorf("aninhamento de %d níveis deve passar", MaxDepth)
	}
	s, d = aninha(MaxDepth + 2)
	if ValidateJSONSchema([]byte(s), []byte(d)) {
		t.Errorf("aninhamento de %d níveis deve falhar", MaxDepth+2)
	}
	if strings.Count(s, "properties") != MaxDepth+2 {
		t.Fatal("o teste montou o aninhamento errado")
	}
}
