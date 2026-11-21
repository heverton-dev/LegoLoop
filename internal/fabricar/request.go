package fabricar

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidRequest: categoria ou modelo ausentes ou desconhecidos.
var ErrInvalidRequest = errors.New("fabricar: pedido inválido")

// Category restringe a categoria a valores conhecidos.
type Category string

const (
	CategoryBuild   Category = "build"
	CategoryReverse Category = "reverse"
	CategoryWeb     Category = "web"
	CategoryMobile  Category = "mobile"
)

// Request representa o payload de entrada para a fábrica de processos.
type Request struct {
	RawName        string   `json:"raw_name"`        // Ex: "Processo do Joãozinho"
	Goal           string   `json:"goal"`            // Objetivo em linguagem natural
	Category       Category `json:"category"`        // Validado contra as constantes acima
	TargetDrives   []string `json:"target_drives"`   // Vazio = todos os drives do SO
	ModelChoice    string   `json:"model_choice"`    // Conferido contra o inventário pela CLI e pelo Studio
	Interactive    bool     `json:"interactive"`     // true = solicita prompts; false = autônomo (--yes)
	Force          bool     `json:"force"`           // true = sobrescreve após backup .legoloop-bak
	TrustRemovable bool     `json:"trust_removable"` // true = aceita peças Untrusted na fusão
	LinkHarnesses  []string `json:"link_harnesses"`  // Ex: ["claude"]; vazio = todos os harnesses detectados
	NoLink         bool     `json:"no_link"`         // true = só o cofre, nenhum vínculo a harness
	Eager          bool     `json:"eager"`           // true = liga a skill em toda pasta de harness (padrão lazy: só no Cofre; o skill_search a acha)
	DryRun         bool     `json:"dry_run"`         // true = calcula tudo e não grava nada
	Peca           string   `json:"peca,omitempty"`  // tipo da peça única (script, skill, gate, hook, mcp, processo); vazio = processo completo
}

// Validate normaliza e checa o pedido antes de qualquer E/S.
func (r Request) Validate() error {
	if _, err := NormalizeProcessName(r.RawName); err != nil {
		return err
	}
	switch r.Category {
	case CategoryBuild, CategoryReverse, CategoryWeb, CategoryMobile:
	default:
		return fmt.Errorf("%w: categoria %q", ErrInvalidRequest, r.Category)
	}
	if r.Peca != "" && !TipoPecaValido(r.Peca) {
		return errTipoPeca(r.Peca)
	}
	if strings.TrimSpace(r.ModelChoice) == "" {
		return fmt.Errorf("%w: model_choice é obrigatório", ErrInvalidRequest)
	}
	return nil
}
