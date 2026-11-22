package a2ui

// Catalog mantém a lista branca de componentes permitidos em superfícies A2UI.
type Catalog struct {
	permitidos map[string]bool
}

// DefaultCatalog devolve o catálogo básico A2UI somado aos componentes de domínio do LegoLoop.
func DefaultCatalog() Catalog {
	nomes := []string{
		"Text", "Heading", "Card", "Row", "Column", "List", "Tabs", "Modal",
		"Button", "TextField", "CheckBox", "MultipleChoice", "Divider", "Badge",
		"LegoDagGraph", "LegoProgressMeter", "LegoStopGateCard",
		"LegoArtifactViewer", "LegoLogStream", "LegoSwarmBoard",
	}
	m := make(map[string]bool, len(nomes))
	for _, n := range nomes {
		m[n] = true
	}
	return Catalog{permitidos: m}
}

// Allows informa se o tipo de componente pertence ao catálogo confiável.
func (c Catalog) Allows(kind string) bool {
	return c.permitidos[kind]
}
