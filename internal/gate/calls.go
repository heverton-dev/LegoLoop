package gate

// contentCalls são os builtins que provam algo sobre o conteúdo (ou a aprovação de uma pessoa), e não só
// que o arquivo existe.
var contentCalls = map[string]bool{
	"word_count": true, "line_count": true, "has_pattern": true,
	"has_binary_criteria": true, "has_acceptance_criteria": true, "schema_valid": true,
	"interactive": true, "prompt": true,
}

// ChecksContent diz se o gate só pode passar quando o conteúdo do artefato (ou a aprovação de uma pessoa)
// foi conferido. Existência de arquivo e código de saída, sozinhos, não provam conteúdo. Em "&&" basta um
// lado provar; em "||" os dois lados precisam provar, senão o lado fraco libera o gate.
func (e *Expr) ChecksContent() bool { return provesContent(e.root) }

func provesContent(n node) bool {
	switch v := n.(type) {
	case callNode:
		return contentCalls[v.name]
	case notNode:
		return provesContent(v.inner)
	case cmpNode:
		return provesContent(v.l) || provesContent(v.r)
	case logicNode:
		if v.op == "&&" {
			return provesContent(v.l) || provesContent(v.r)
		}
		return provesContent(v.l) && provesContent(v.r)
	}
	return false
}
