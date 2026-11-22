package processo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/gate"
)

// ErrContratoInvalido: o manifesto contrato.json da peça não satisfaz o contrato de I/O.
var ErrContratoInvalido = errors.New("processo: contrato inválido")

// tiposDeSaida são os tipos de artefato que uma saída declara; url e estado só fazem sentido na entrada.
var tiposDeSaida = map[string]bool{"marcador": true, "texto": true, "dados": true, "codigo": true, "json": true}

// Entrada declara o que a peça consome.
type Entrada struct {
	Tipo     string `json:"tipo"` // marcador | texto | dados | codigo | json | url | estado
	Opcional bool   `json:"opcional"`
}

// Saida declara o que a peça produz.
type Saida struct {
	Tipo          string `json:"tipo"`                     // marcador | texto | dados | codigo | json
	CaminhoPadrao string `json:"caminho_padrao,omitempty"` // relativo; vazio = o molde do motor decide
	Promete       bool   `json:"promete"`
	GateSugerido  string `json:"gate_sugerido,omitempty"` // DSL de internal/gate; vazio = GateParaArtefato
}

// Contrato é o manifesto contrato.json da peça: o I/O declarado que os verbos derivam.
type Contrato struct {
	Nome      string    `json:"nome"`
	Tipo      string    `json:"tipo"` // skill | script | hook | gate | mcp | processo
	Naturezas []string  `json:"naturezas"`
	Entra     []Entrada `json:"entra"`
	Sai       []Saida   `json:"sai"`
	Efeitos   []string  `json:"efeitos"` // arquivo | rede | terminal | estado
}

func contratoInvalido(format string, a ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{ErrContratoInvalido}, a...)...)
}

// Validate checa o contrato antes de qualquer E/S: nome, ao menos uma saída, path jail do caminho padrão
// e gate sugerido que compila.
func (c Contrato) Validate() error {
	if strings.TrimSpace(c.Nome) == "" {
		return contratoInvalido("nome vazio")
	}
	if len(c.Sai) == 0 {
		return contratoInvalido("sem saída")
	}
	for i, s := range c.Sai {
		if err := s.validate(i); err != nil {
			return err
		}
	}
	return nil
}

// validate de uma saída: tipo permitido, caminho padrão dentro da raiz e gate compilável.
func (s Saida) validate(i int) error {
	if !tiposDeSaida[s.Tipo] {
		return contratoInvalido("saída %d: tipo %q desconhecido (use marcador, texto, dados, codigo ou json)", i, s.Tipo)
	}
	if s.CaminhoPadrao != "" && !caminhoLocal(s.CaminhoPadrao) {
		return contratoInvalido("saída %d: o caminho_padrao %q precisa de caminho relativo dentro do projeto", i, s.CaminhoPadrao)
	}
	if s.GateSugerido != "" {
		if _, err := gate.Compile(s.GateSugerido); err != nil {
			return contratoInvalido("saída %d: o gate_sugerido não compila: %v", i, err)
		}
	}
	return nil
}

// ContratoDe carrega e valida o manifesto contrato.json. JSON quebrado ou contrato violado viram
// ErrContratoInvalido; falha de E/S volta crua. Nada é criado no disco.
func ContratoDe(caminho string) (Contrato, error) {
	b, err := os.ReadFile(caminho)
	if err != nil {
		return Contrato{}, err
	}
	var c Contrato
	if err := json.Unmarshal(b, &c); err != nil {
		return Contrato{}, contratoInvalido("%s: %v", filepath.Base(caminho), err)
	}
	if err := c.Validate(); err != nil {
		return Contrato{}, err
	}
	return c, nil
}
