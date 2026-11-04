package matriz

import (
	"encoding/json"
	"fmt"
	"os"

	"legoloop/internal/reputation"
)

// VersaoMatriz é a versão do formato do registro. A 1.0.0 tinha só as notas de qualidade (`ratings`).
const VersaoMatriz = "2.0.0"

// Entrada é uma peça na matriz: identidade, hash do conteúdo, nota de qualidade e perfil de uso.
type Entrada struct {
	Ref      string                 `json:"ref"`
	Nome     string                 `json:"nome"`
	Tipo     string                 `json:"tipo"` // skill | mcp | script | hook | gate | loop
	Hash     string                 `json:"hash"` // SHA-256 do conteúdo (ou da identidade, em MCP e hook)
	Perfil   *PerfilDeUso           `json:"perfil,omitempty"`
	Tipado   map[string]any         `json:"tipado,omitempty"`   // campos da matriz por tipo; só nos arquivos por tipo
	Retentar bool                   `json:"retentar,omitempty"` // o cérebro falhou nesta peça: tentar de novo na próxima geração
	Alertas  []string               `json:"alertas,omitempty"`  // texto da peça que parece ordem (citado como dado)
	Rating   reputation.PieceRating `json:"-"`                  // preenchido ao ler, a partir de Ratings (mesma Ref)
}

// Registro é o arquivo registry_matrix.json. O array `ratings` mantém o formato da v1, para os leitores
// antigos (Studio, `legoloop matrix`) seguirem funcionando; `entradas` traz o perfil de uso.
type Registro struct {
	GeneratedAt string                   `json:"generated_at"`
	Versao      string                   `json:"matrix_version"`
	VocabVersao int                      `json:"vocab_version,omitempty"`
	TotalPecas  int                      `json:"total_pieces"`
	Ratings     []reputation.PieceRating `json:"ratings"`
	Entradas    []Entrada                `json:"entradas,omitempty"`
}

// LerRegistro lê a matriz, v1 ou v2. Na v1 não há `entradas`: as notas são lidas inteiras.
func LerRegistro(path string) (*Registro, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var r Registro
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, fmt.Errorf("matriz: %s inválido: %w", path, err)
	}
	porRef := make(map[string]reputation.PieceRating, len(r.Ratings))
	for _, n := range r.Ratings {
		porRef[n.Ref] = n
	}
	for i := range r.Entradas {
		r.Entradas[i].Rating = porRef[r.Entradas[i].Ref]
	}
	return &r, nil
}

// arquivoPorTipo é o conteúdo de matrix/<tipo>.json.
type arquivoPorTipo struct {
	Versao   string    `json:"matrix_version"`
	Tipo     string    `json:"tipo"`
	Total    int       `json:"total"`
	Entradas []Entrada `json:"entradas"`
}
