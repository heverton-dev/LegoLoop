package matriz

import (
	"errors"

	"legoloop/internal/taxonomy"
)

// ErrRespostaInvalida: o cérebro respondeu algo que não é o JSON de perfil esperado.
var ErrRespostaInvalida = errors.New("matriz: resposta do cérebro inválida")

// Origens de um perfil. "inferido" vale como palpite até o usuário revisar e confirmar.
const (
	OrigemDeclarado  = "declarado"
	OrigemInferido   = "inferido"
	OrigemConfirmado = "confirmado"
)

// PerfilDeUso é o que falta à matriz de qualidade: para que a peça serve.
type PerfilDeUso struct {
	CasosDeUso    []string `json:"casos_de_uso"`
	Etapas        []string `json:"etapas"`
	Naturezas     []string `json:"naturezas"`
	Consome       []string `json:"consome,omitempty"`
	Produz        []string `json:"produz,omitempty"`
	QuandoNaoUsar []string `json:"quando_nao_usar,omitempty"` // termos de exclusão (do vocabulário)
	Origem        string   `json:"origem"`
}

// Peca é o que o catalogador lê de uma peça. O texto é dado não confiável, nunca instrução.
type Peca struct {
	Ref               string
	Nome              string
	Descricao         string
	Tipo              string   // skill | mcp | hook | script | gate | processo
	Conteudo          string   // texto lido da peça (frontmatter e corpo)
	PapelPrimario     string   `json:"papel_primario,omitempty"`
	PapeisSecundarios []string `json:"papeis_secundarios,omitempty"`
}

// Validate verifica a conformidade da peca com a taxonomia de papeis canonicos.
func (p Peca) Validate() error {
	primario := p.PapelPrimario
	secundarios := p.PapeisSecundarios
	if primario == "" && len(secundarios) == 0 && p.Conteudo != "" {
		prim, sec := ExtrairPapeis(p.Conteudo)
		if primario == "" {
			primario = prim
		}
		if len(secundarios) == 0 {
			secundarios = sec
		}
	}
	if err := taxonomy.ValidarPapel(primario); err != nil {
		return taxonomy.ErrPapelPrimarioObrigatorio
	}
	for _, sec := range secundarios {
		if err := taxonomy.ValidarPapel(sec); err != nil {
			return taxonomy.ErrPapelInvalido
		}
	}
	return nil
}

// PreencherPapeis atribui os papeis declarados no frontmatter se nao estiverem definidos.
func (p *Peca) PreencherPapeis() {
	if p.Conteudo == "" {
		return
	}
	prim, sec := ExtrairPapeis(p.Conteudo)
	if p.PapelPrimario == "" {
		p.PapelPrimario = prim
	}
	if len(p.PapeisSecundarios) == 0 {
		p.PapeisSecundarios = sec
	}
}

// Resultado é o que a catalogação de uma peça devolve.
type Resultado struct {
	Perfil    *PerfilDeUso // nil = aptidão desconhecida
	Pendentes []Pendente   // termos propostos fora do vocabulário
	Alertas   []string     // texto da peça que parece ordem; citado como dado
}
