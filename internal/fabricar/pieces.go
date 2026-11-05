package fabricar

import "legoloop/internal/reputation"

// PieceKind é o tipo da peça descoberta.
type PieceKind string

const (
	KindSkill  PieceKind = "skill"
	KindScript PieceKind = "script"
	KindHook   PieceKind = "hook"
	KindGate   PieceKind = "gate"
	KindLoop   PieceKind = "loop"
	KindMCP    PieceKind = "mcp"
)

// DiscoveredPiece representa uma peça localizada em qualquer drive montado.
type DiscoveredPiece struct {
	Name      string          `json:"name"`
	Kind      PieceKind       `json:"kind"`
	Drive     string          `json:"drive"` // Ex: "C:\", "D:\", "/Volumes/Pendrive"
	Path      string          `json:"path"`
	SHA256    string          `json:"sha256"`    // Proveniência e desduplicação
	Score     float64         `json:"score"`     // 0.00 a 1.00, calculado por internal/reputation
	Tier      reputation.Tier `json:"tier"`      // ouro, prata, bronze, toxica
	Untrusted bool            `json:"untrusted"` // true = mídia removível ou rede
}
