package dag

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
)

// Validador mínimo de JSON Schema Draft-07. Implementa só as palavras-chave
// que dag.v1.json e dag.v2.json usam. Palavra-chave desconhecida falha o
// teste: o validador não aprova em silêncio o que não entende.
var palavrasAnotacao = map[string]bool{
	"$schema": true, "title": true, "description": true,
}

func validarSchema(schema map[string]any, v any, caminho string) []string {
	var erros []string
	for k, sv := range schema {
		switch k {
		case "type":
			if !tipoConfere(sv, v) {
				erros = append(erros, fmt.Sprintf("%s: tipo %v esperado", caminho, sv))
				return erros
			}
		case "enum":
			ok := false
			for _, e := range sv.([]any) {
				if e == v {
					ok = true
				}
			}
			if !ok {
				erros = append(erros, fmt.Sprintf("%s: valor %v fora do enum", caminho, v))
			}
		case "pattern":
			if s, ok := v.(string); ok && !regexp.MustCompile(sv.(string)).MatchString(s) {
				erros = append(erros, fmt.Sprintf("%s: %q não casa com %s", caminho, s, sv))
			}
		case "minimum":
			if n, ok := v.(float64); ok && n < sv.(float64) {
				erros = append(erros, fmt.Sprintf("%s: %v abaixo do mínimo", caminho, n))
			}
		case "maximum":
			if n, ok := v.(float64); ok && n > sv.(float64) {
				erros = append(erros, fmt.Sprintf("%s: %v acima do máximo", caminho, n))
			}
		case "required":
			if obj, ok := v.(map[string]any); ok {
				for _, r := range sv.([]any) {
					if _, tem := obj[r.(string)]; !tem {
						erros = append(erros, fmt.Sprintf("%s: campo %q ausente", caminho, r))
					}
				}
			}
		case "properties":
			if obj, ok := v.(map[string]any); ok {
				for nome, sub := range sv.(map[string]any) {
					if val, tem := obj[nome]; tem {
						erros = append(erros, validarSchema(sub.(map[string]any), val, caminho+"."+nome)...)
					}
				}
			}
		case "items":
			if arr, ok := v.([]any); ok {
				for i, el := range arr {
					erros = append(erros, validarSchema(sv.(map[string]any), el, fmt.Sprintf("%s[%d]", caminho, i))...)
				}
			}
		default:
			if !palavrasAnotacao[k] {
				erros = append(erros, fmt.Sprintf("%s: palavra-chave %q não suportada pelo validador", caminho, k))
			}
		}
	}
	return erros
}

func tipoConfere(t any, v any) bool {
	if lista, ok := t.([]any); ok {
		for _, x := range lista {
			if tipoConfere(x, v) {
				return true
			}
		}
		return false
	}
	switch t.(string) {
	case "object":
		_, ok := v.(map[string]any)
		return ok
	case "array":
		_, ok := v.([]any)
		return ok
	case "string":
		_, ok := v.(string)
		return ok
	case "boolean":
		_, ok := v.(bool)
		return ok
	case "null":
		return v == nil
	case "integer":
		n, ok := v.(float64)
		return ok && n == float64(int64(n))
	case "number":
		_, ok := v.(float64)
		return ok
	}
	return false
}

func carregarSchema(t *testing.T, nome string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "schemas", nome))
	if err != nil {
		t.Fatalf("lendo schema %s: %v", nome, err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("schema %s inválido: %v", nome, err)
	}
	return m
}

func validarContra(t *testing.T, schema map[string]any, raw []byte) []string {
	t.Helper()
	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	return validarSchema(schema, doc, "$")
}

// arquivosProcessos lista os JSON das pastas de Processos do repositório (os dois dialetos).
func arquivosProcessos(t *testing.T) []string {
	t.Helper()
	var out []string
	pastas := caminhos.DirsLeitura(filepath.Join("..", ".."))
	if len(pastas) == 0 {
		t.Fatal("nenhuma pasta de Processos na raiz do repositório")
	}
	for _, raiz := range pastas {
		err := filepath.WalkDir(raiz, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(p, ".json") {
				out = append(out, p)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return out
}

// Gate do TICK-01.3: todo Processo do repositório valida contra v1 e contra v2.
func TestSchemasV1eV2ValidamProcessos(t *testing.T) {
	v1, v2 := carregarSchema(t, "dag.v1.json"), carregarSchema(t, "dag.v2.json")
	validados := 0
	for _, f := range arquivosProcessos(t) {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		// Arquivo que o parser recusa não é Processo válido (ex.: o documento
		// legado do Studio ai-motion-graphics, fora do teste dourado): fica de fora.
		if _, err := Parse(raw); err != nil {
			t.Logf("fora do gate (Parse recusa): %s: %v", filepath.Base(f), err)
			continue
		}
		validados++
		for nome, sc := range map[string]map[string]any{"v1": v1, "v2": v2} {
			if erros := validarContra(t, sc, raw); len(erros) > 0 {
				t.Errorf("%s não valida contra %s: %v", filepath.Base(f), nome, erros)
			}
		}
	}
	// Pelo menos 23 arquivos passam no parser (inclui os 15 do teste dourado)
	// e 1 legado antigo é recusado pelo parser.
	if validados < 23 {
		t.Fatalf("validei %d Processos do repositório, esperado pelo menos 23", validados)
	}
}

const dialetoNovo = `{
  "name": "novo",
  "trigger": {"type": "processo_fim", "processo": "outro"},
  "nodes": [
    {"id": "TICK-01.3", "type": "subprocesso", "ref": "filho",
     "routing": {"mode": "manual", "model": "sonnet"}, "output": {"artifact": "a.json", "gate": "file_exists('a.json')"}}
  ]
}`

// O dialeto novo passa em v2, falha em v1 e o parser o aceita.
func TestSchemaV2AceitaDialetoNovo(t *testing.T) {
	v1, v2 := carregarSchema(t, "dag.v1.json"), carregarSchema(t, "dag.v2.json")
	if erros := validarContra(t, v2, []byte(dialetoNovo)); len(erros) > 0 {
		t.Errorf("v2 deveria aceitar o dialeto novo: %v", erros)
	}
	if erros := validarContra(t, v1, []byte(dialetoNovo)); len(erros) == 0 {
		t.Error("v1 deveria recusar o dialeto novo")
	}
	if _, err := Parse([]byte(dialetoNovo)); err != nil {
		t.Errorf("Parse deveria aceitar o dialeto novo: %v", err)
	}
}

// Mutação: o validador recusa Processo quebrado em v1 e v2.
func TestSchemaRecusaProcessoQuebrado(t *testing.T) {
	quebrados := map[string]string{
		"sem nodes":      `{"name": "x"}`,
		"tipo ruim":      strings.Replace(dialetoNovo, `"subprocesso"`, `"inexistente"`, 1),
		"id com @":       strings.Replace(dialetoNovo, `TICK-01.3`, `TICK@1`, 1),
		"id com . final": strings.Replace(dialetoNovo, `TICK-01.3`, `TICK-01.`, 1),
	}
	for _, nome := range []string{"dag.v1.json", "dag.v2.json"} {
		sc := carregarSchema(t, nome)
		for caso, raw := range quebrados {
			if erros := validarContra(t, sc, []byte(raw)); len(erros) == 0 {
				t.Errorf("%s deveria recusar %q", nome, caso)
			}
		}
	}
}
