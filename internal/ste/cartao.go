package ste

import (
	"fmt"
	"strings"
)

// Parametro é algo que "personalizar" abre: o usuário pode mudar o valor antes de seguir.
type Parametro struct {
	Nome  string `json:"nome"`
	Valor string `json:"valor"`
	Ajuda string `json:"ajuda"`
}

// Cartao explica uma ação com efeito antes de ela acontecer.
type Cartao struct {
	Acao         string      `json:"acao"`
	Aviso        string      `json:"aviso,omitempty"` // risco; sai ANTES dos passos
	Passos       []string    `json:"passos"`
	Efeitos      []string    `json:"efeitos"`
	SeDerErrado  string      `json:"se_der_errado"`
	Irreversivel bool        `json:"irreversivel"` // exige aviso e pausa extra
	Parametros   []Parametro `json:"parametros,omitempty"`
}

// Decisao é a escolha do usuário diante do cartão.
type Decisao int

const (
	Cancelar Decisao = iota
	Seguir
	Personalizar
)

// ParseDecisao lê a resposta do usuário. Enter vazio cancela (o padrão seguro).
func ParseDecisao(in string) (Decisao, bool) {
	switch strings.ToLower(strings.TrimSpace(in)) {
	case "1", "s", "sim", "y", "yes":
		return Seguir, true
	case "2", "p", "personalizar":
		return Personalizar, true
	case "0", "n", "nao", "não", "no", "":
		return Cancelar, true
	}
	return Cancelar, false
}

// Render monta o texto do cartão. O aviso de risco vem antes dos passos.
func (c Cartao) Render() string {
	var b strings.Builder
	fmt.Fprintf(&b, "AÇÃO: %s\n", c.Acao)
	if c.Aviso != "" {
		fmt.Fprintf(&b, "AVISO: %s\n", c.Aviso)
	}
	b.WriteString("O QUE VAI ACONTECER:\n")
	for i, p := range c.Passos {
		fmt.Fprintf(&b, "  %d. %s\n", i+1, p)
	}
	for _, e := range c.Efeitos {
		fmt.Fprintf(&b, "  Efeito: %s\n", e)
	}
	fmt.Fprintf(&b, "SE DER ERRADO: %s\n", c.SeDerErrado)
	b.WriteString("Escolha: [1] seguir  [2] personalizar  [0] cancelar\n")
	return b.String()
}

// LintCartao aplica o perfil ao texto de cada campo e confere a estrutura do cartão.
func LintCartao(c Cartao) []Violation {
	var out []Violation
	textos := append([]string{c.Acao, c.Aviso, c.SeDerErrado}, c.Passos...)
	for _, t := range append(textos, c.Efeitos...) {
		out = append(out, Lint(t)...)
	}
	if c.Irreversivel && strings.TrimSpace(c.Aviso) == "" {
		out = append(out, Violation{Msg: "ação irreversível precisa de aviso de risco antes dos passos", Regra: "sem-aviso"})
	}
	if strings.TrimSpace(c.SeDerErrado) == "" {
		out = append(out, Violation{Msg: "o cartão precisa dizer como desfazer a ação", Regra: "sem-desfazer"})
	}
	return out
}
