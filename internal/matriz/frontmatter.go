package matriz

import (
	"regexp"
	"strings"
)

// frontmatter devolve os campos de lista da seção "---" no topo do texto. Entende "[a, b]", "a, b" e
// listas "- item" nas linhas seguintes. Só nos interessam os campos do perfil de uso.
func frontmatter(texto string) map[string][]string {
	linhas := strings.Split(strings.ReplaceAll(texto, "\r\n", "\n"), "\n")
	if len(linhas) == 0 || strings.TrimSpace(linhas[0]) != "---" {
		return nil
	}
	out := map[string][]string{}
	chave := ""
	for _, l := range linhas[1:] {
		if strings.TrimSpace(l) == "---" {
			break
		}
		if item, ok := strings.CutPrefix(strings.TrimSpace(l), "- "); ok && chave != "" {
			out[chave] = append(out[chave], limpa(item))
			continue
		}
		k, v, ok := strings.Cut(l, ":")
		if !ok || strings.HasPrefix(l, " ") {
			continue
		}
		chave = strings.TrimSpace(k)
		if lista := itens(v); len(lista) > 0 {
			out[chave] = lista
		}
	}
	return out
}

func itens(v string) []string {
	v = strings.TrimSpace(v)
	v = strings.TrimSuffix(strings.TrimPrefix(v, "["), "]")
	var out []string
	for _, p := range strings.Split(v, ",") {
		if p = limpa(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func limpa(s string) string { return strings.Trim(strings.TrimSpace(s), `"'`) }

var reQuandoNaoUsar = regexp.MustCompile(`(?i)quando\s+n[ãa]o\s+usar[:*\s]*([^\n]+)`)

// ExtrairPapeis extrai papel_primario e papeis_secundarios do frontmatter da peca.
// Reconhece as chaves papel_primario / papel e papeis_secundarios / papeis.
func ExtrairPapeis(texto string) (primario string, secundarios []string) {
	fm := frontmatter(texto)
	if fm == nil {
		return "", nil
	}
	if v, ok := fm["papel_primario"]; ok && len(v) > 0 {
		primario = v[0]
	} else if v, ok := fm["papel"]; ok && len(v) > 0 {
		primario = v[0]
	}
	if v, ok := fm["papeis_secundarios"]; ok && len(v) > 0 {
		secundarios = v
	} else if v, ok := fm["papeis"]; ok && len(v) > 0 {
		secundarios = v
	}
	return primario, secundarios
}

// trechoQuandoNaoUsar devolve a frase "Quando NÃO usar: ..." do texto da peça, se existir.
func trechoQuandoNaoUsar(texto string) string {
	if m := reQuandoNaoUsar.FindStringSubmatch(texto); m != nil {
		return m[1]
	}
	return ""
}
