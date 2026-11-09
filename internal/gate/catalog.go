package gate

import "sort"

// Doc descreve um gate embutido para o inventário e as sugestões do Studio.
type Doc struct {
	Name        string `json:"name"`
	Signature   string `json:"signature"`
	Description string `json:"description"`
}

// docs é indexado pelo nome do builtin; TestCatalogCoversEveryBuiltin garante a cobertura.
var docs = map[string]Doc{
	"file_exists":             {Signature: "file_exists(caminho)", Description: "Verdadeiro se o arquivo existe dentro do projeto."},
	"word_count":              {Signature: "word_count(caminho)", Description: "Número de palavras do arquivo, para comparar com um mínimo."},
	"line_count":              {Signature: "line_count(caminho)", Description: "Número de linhas do arquivo."},
	"has_pattern":             {Signature: "has_pattern(caminho, regex)", Description: "Verdadeiro se o conteúdo do arquivo casa com a expressão regular."},
	"has_binary_criteria":     {Signature: "has_binary_criteria(caminho)", Description: "Verdadeiro se o arquivo tem checklist de critérios binários (- [ ] item)."},
	"has_acceptance_criteria": {Signature: "has_acceptance_criteria(caminho)", Description: "Verdadeiro se o arquivo tem uma seção de critérios de aceitação."},
	"schema_valid":            {Signature: "schema_valid(schema, alvo)", Description: "Verdadeiro se o arquivo JSON no caminho alvo cumpre o JSON Schema informado."},
	"interactive":             {Signature: "interactive(pergunta)", Description: "Pausa a execução e solicita confirmação humana (Stop-Gate interativo)."},
	"prompt":                  {Signature: "prompt(pergunta)", Description: "Alias de interactive(pergunta) para Stop-Gates humanos."},
	"no_stubs":                {Signature: "no_stubs(caminho)", Description: "Verdadeiro se o arquivo existe e não contém stubs vazios, TODOs ou mocks simulados."},
}

// Catalog lista os gates embutidos em ordem alfabética.
func Catalog() []Doc {
	out := make([]Doc, 0, len(docs))
	for name, d := range docs {
		d.Name = name
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
