package combo

// mapaLegado mapeia nomes de combos legados para os novos combos honestos certificados.
var mapaLegado = map[string]string{
	"code-pro":            "t-codigo-128k",
	"code-fast":           "t-chat-128k",
	"code-free":           "t-agente-128k",
	"t" + "ask-micro":     "t-micro-8k",
	"megacombo-coding":    "t-codigo-128k",
	"megacombo-reasoning": "t-raciocinio-128k",
	"megacombo-vision":    "t-visao-128k",
}

// MapearComboLegado traduz nomes legados para combos canônicos de 128k/8k.
// Se o nome já existir no catálogo ou for desconhecido, retorna o nome inalterado.
func MapearComboLegado(nome string) string {
	if mapeado, ok := mapaLegado[nome]; ok {
		return mapeado
	}
	if _, ok := CatalogoCombos[nome]; ok {
		return nome
	}
	return nome
}

// ResolverCombo mapeia o nome legado se necessário e busca o combo no catálogo certificado.
func ResolverCombo(nome string) (ComboHonesto, error) {
	nomeMapeado := MapearComboLegado(nome)
	return ObterCombo(nomeMapeado)
}
