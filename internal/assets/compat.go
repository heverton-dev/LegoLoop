package assets

import (
	"fmt"
	"strings"
)

// CompatibleItem descreve uma peça compatível e o motivo da conexão ontológica.
type CompatibleItem struct {
	Item   Item   `json:"item"`
	Reason string `json:"reason"`
	Role   string `json:"role"` // "upstream", "downstream", "tool", "hook", "slot"
}

// FindCompatible varre a lista de ativos e retorna todas as peças com afinidade ontológica.
func FindCompatible(target Item, all []Item) []CompatibleItem {
	var out []CompatibleItem
	tConsumes := target.Meta["consumes"]
	tProduces := target.Meta["produces"]

	for _, it := range all {
		if it.ID == target.ID && it.Kind == target.Kind && it.Name == target.Name {
			continue
		}
		if comp, ok := matchCompatibility(target, it, tConsumes, tProduces); ok {
			out = append(out, comp)
		}
	}
	return out
}

func matchCompatibility(target, it Item, tConsumes, tProduces string) (CompatibleItem, bool) {
	itConsumes := it.Meta["consumes"]
	itProduces := it.Meta["produces"]

	if tProduces != "" && itConsumes != "" && strings.EqualFold(tProduces, itConsumes) {
		return CompatibleItem{
			Item:   it,
			Role:   "downstream",
			Reason: fmt.Sprintf("Consome a saída '%s' produzida por esta peça", tProduces),
		}, true
	}
	if tConsumes != "" && itProduces != "" && strings.EqualFold(tConsumes, itProduces) {
		return CompatibleItem{
			Item:   it,
			Role:   "upstream",
			Reason: fmt.Sprintf("Produz a entrada '%s' consumida por esta peça", tConsumes),
		}, true
	}
	return matchAuxiliaryCompatibility(target, it)
}

func matchAuxiliaryCompatibility(target, it Item) (CompatibleItem, bool) {
	if it.Kind == KindMCP && (target.Kind == KindSkill || target.Kind == KindLoop) {
		return CompatibleItem{
			Item:   it,
			Role:   "tool",
			Reason: fmt.Sprintf("Servidor MCP '%s' fornece ferramentas para a execução", it.Name),
		}, true
	}
	if it.Kind == KindHook && (target.Kind == KindSkill || target.Kind == KindScript) {
		return CompatibleItem{
			Item:   it,
			Role:   "hook",
			Reason: fmt.Sprintf("Hook '%s' pode ser acoplado nos slots de ciclo de vida", it.Name),
		}, true
	}
	if it.Kind == KindGate && target.Kind == KindSkill {
		return CompatibleItem{
			Item:   it,
			Role:   "slot",
			Reason: fmt.Sprintf("Gate '%s' pode auditar os artefatos de saída", it.Name),
		}, true
	}
	return CompatibleItem{}, false
}
