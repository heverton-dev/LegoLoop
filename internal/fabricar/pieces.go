package fabricar

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/reputation"
)

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

// ScanOptions limita a varredura. Zero usa o padrão.
type ScanOptions struct {
	MaxDepth       int           // diretórios entre o perfil e o arquivo da peça (padrão 4)
	MaxFiles       int           // itens inspecionados por raiz (padrão 20000)
	MaxPieces      int           // peças devolvidas por raiz (padrão 500)
	PerRootTimeout time.Duration // prazo por raiz (padrão 60 s)
}

func (o ScanOptions) withDefaults() ScanOptions {
	if o.MaxDepth <= 0 {
		o.MaxDepth = 4
	}
	if o.MaxFiles <= 0 {
		o.MaxFiles = 20000
	}
	if o.MaxPieces <= 0 {
		o.MaxPieces = 500
	}
	if o.PerRootTimeout <= 0 {
		o.PerRootTimeout = 60 * time.Second
	}
	return o
}

// ScanResult devolve as peças e o que foi pulado, com o motivo.
type ScanResult struct {
	Pieces   []DiscoveredPiece `json:"pieces"`
	Skipped  []SkippedRoot     `json:"skipped"`
	Warnings []string          `json:"warnings"`
}

// Scanner é a borda de descoberta de peças.
type Scanner interface {
	Scan(ctx context.Context, roots []string, opts ScanOptions) (ScanResult, error)
}

// kindOf traduz as categorias do inventário. Gates, harnesses e binários não são peças de unidade.
var kindOf = map[assets.Kind]PieceKind{
	assets.KindSkill: KindSkill, assets.KindScript: KindScript, assets.KindHook: KindHook,
	assets.KindLoop: KindLoop, assets.KindMCP: KindMCP,
}

// maxHashBytes limita a leitura para o SHA-256 (proteção de memória e de tempo).
const maxHashBytes = 16 << 20

// newPiece converte um item do inventário em peça. Só lê: nada é executado. O SHA-256 é preenchido depois do corte por MaxPieces.
func newPiece(it assets.Item, root string, untrusted bool) (DiscoveredPiece, bool) {
	kind, ok := kindOf[it.Kind]
	if !ok || it.Scope == "builtin" || it.Scope == "system" {
		return DiscoveredPiece{}, false
	}
	score, tier := rate(it, kind)
	return DiscoveredPiece{
		Name: it.Name, Kind: kind, Drive: root, Path: it.Path,
		Score: score, Tier: tier, Untrusted: untrusted,
	}, true
}

// rate pontua com internal/reputation: pelo arquivo quando a peça é um arquivo próprio (skill, script, loop)
// e pela descrição quando vem de uma configuração compartilhada (hook, MCP).
func rate(it assets.Item, kind PieceKind) (float64, reputation.Tier) {
	if kind == KindSkill || kind == KindScript || kind == KindLoop {
		if r, err := reputation.ScoreFile(it.Path); err == nil {
			return r.TotalScore, r.Tier
		}
	}
	r := reputation.ScorePiece(it.Path, it.Name, it.Description, 0, 0)
	return r.TotalScore, r.Tier
}

func fileSHA(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, io.LimitReader(f, maxHashBytes)); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
