// Package matriz é o motor das matrizes de peças: perfil de uso (para que a peça serve), catalogador,
// geração incremental e detecção de defasagem. O vocabulário de termos é fechado e versionado; cresce só
// com aprovação do usuário. A seleção de peças pelo fabricar consulta este motor e nunca o regenera.
package matriz

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"legoloop/internal/processo"
)

//go:embed vocab.default.json
var vocabDefaultJSON []byte

// ErrVocabulario: termo fora do vocabulário, mal formado ou sem proposta pendente.
var ErrVocabulario = errors.New("matriz: termo fora do vocabulário")

// Categoria separa os três vocabulários do perfil.
type Categoria string

const (
	CatCasoDeUso Categoria = "caso_de_uso"
	CatEtapa     Categoria = "etapa"
	CatNatureza  Categoria = "natureza"
)

const (
	vocabFile    = "vocab.json"
	pendenteFile = "vocab.pendente.json"
)

var reTermo = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// TermoValido diz se o nome pode ser um termo do vocabulário (kebab-case ASCII).
func TermoValido(nome string) bool { return reTermo.MatchString(nome) }

// Termo é um item do vocabulário; Palavras são as pistas usadas na inferência determinística.
type Termo struct {
	Nome     string   `json:"nome"`
	Palavras []string `json:"palavras,omitempty"`
}

// Vocab é o vocabulário fechado de termos do perfil de uso.
type Vocab struct {
	Versao     int     `json:"versao"`
	Etapas     []Termo `json:"etapas"`
	Naturezas  []Termo `json:"naturezas"`
	CasosDeUso []Termo `json:"casos_de_uso"`
}

// DefaultVocab devolve o vocabulário inicial, embutido no binário.
func DefaultVocab() Vocab {
	var v Vocab
	if err := json.Unmarshal(vocabDefaultJSON, &v); err != nil || v.Versao == 0 {
		panic(fmt.Sprintf("matriz: vocab.default.json inválido: %v", err))
	}
	return v
}

// LoadVocab lê o vocabulário do cofre; sem arquivo, usa o padrão (e não grava nada).
func LoadVocab(vault string) (Vocab, error) {
	b, err := os.ReadFile(filepath.Join(vault, vocabFile))
	if errors.Is(err, os.ErrNotExist) {
		return DefaultVocab(), nil
	}
	if err != nil {
		return Vocab{}, err
	}
	var v Vocab
	if err := json.Unmarshal(b, &v); err != nil || v.Versao == 0 {
		return Vocab{}, fmt.Errorf("matriz: %s inválido: %v", vocabFile, err)
	}
	return v, nil
}

// Termos devolve a lista da categoria.
func (v Vocab) Termos(c Categoria) []Termo {
	switch c {
	case CatCasoDeUso:
		return v.CasosDeUso
	case CatEtapa:
		return v.Etapas
	case CatNatureza:
		return v.Naturezas
	}
	return nil
}

// Has diz se o termo existe na categoria.
func (v Vocab) Has(c Categoria, nome string) bool {
	return slices.ContainsFunc(v.Termos(c), func(t Termo) bool { return t.Nome == nome })
}

// HasAny diz se o termo existe em alguma categoria.
func (v Vocab) HasAny(nome string) bool {
	return v.Has(CatCasoDeUso, nome) || v.Has(CatEtapa, nome) || v.Has(CatNatureza, nome)
}

func (v *Vocab) add(c Categoria, nome string) {
	t := Termo{Nome: nome}
	switch c {
	case CatCasoDeUso:
		v.CasosDeUso = append(v.CasosDeUso, t)
	case CatEtapa:
		v.Etapas = append(v.Etapas, t)
	case CatNatureza:
		v.Naturezas = append(v.Naturezas, t)
	}
}

// Pendente é uma proposta de termo novo, à espera da aprovação do usuário.
type Pendente struct {
	Termo     string    `json:"termo"`
	Categoria Categoria `json:"categoria"`
	Peca      string    `json:"peca"`            // a peça que pediu o termo
	Pecas     []string  `json:"pecas,omitempty"` // todas as peças que pediram (preenchido ao gravar)
}

// LerPendentes devolve a fila de propostas; sem arquivo, a fila é vazia.
func LerPendentes(vault string) ([]Pendente, error) {
	b, err := os.ReadFile(filepath.Join(vault, pendenteFile))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var ps []Pendente
	if err := json.Unmarshal(b, &ps); err != nil {
		return nil, fmt.Errorf("matriz: %s inválido: %w", pendenteFile, err)
	}
	return ps, nil
}

// AdicionarPendentes junta novas propostas à fila, sem duplicar (mesmo termo e categoria) e lembrando as peças.
func AdicionarPendentes(vault string, novas []Pendente) error {
	atual, err := LerPendentes(vault)
	if err != nil {
		return err
	}
	for _, n := range novas {
		if !TermoValido(n.Termo) {
			return fmt.Errorf("%w: %q não é um termo válido (use kebab-case)", ErrVocabulario, n.Termo)
		}
		atual = mesclar(atual, n)
	}
	return gravar(vault, nil, atual)
}

func mesclar(fila []Pendente, n Pendente) []Pendente {
	i := slices.IndexFunc(fila, func(p Pendente) bool { return p.Termo == n.Termo && p.Categoria == n.Categoria })
	if i < 0 {
		n.Pecas = unicas(nil, n.Peca, n.Pecas)
		return append(fila, n)
	}
	fila[i].Pecas = unicas(fila[i].Pecas, n.Peca, n.Pecas)
	return fila
}

func unicas(base []string, um string, varias []string) []string {
	for _, p := range append([]string{um}, varias...) {
		if p != "" && !slices.Contains(base, p) {
			base = append(base, p)
		}
	}
	return base
}

// Aprovar move o termo da fila para o vocabulário e sobe a versão. Só o usuário chama isto.
func Aprovar(vault, termo string) error {
	if !TermoValido(termo) {
		return fmt.Errorf("%w: %q não é um termo válido (use kebab-case)", ErrVocabulario, termo)
	}
	fila, err := LerPendentes(vault)
	if err != nil {
		return err
	}
	v, err := LoadVocab(vault)
	if err != nil {
		return err
	}
	var resto []Pendente
	achou := false
	for _, p := range fila {
		if p.Termo != termo {
			resto = append(resto, p)
			continue
		}
		achou = true
		if !v.Has(p.Categoria, termo) {
			v.add(p.Categoria, termo)
		}
	}
	if !achou {
		return fmt.Errorf("%w: %q não está na fila de propostas", ErrVocabulario, termo)
	}
	v.Versao++
	return gravar(vault, &v, resto)
}

// gravar escreve o vocabulário (se houver) e a fila numa só operação atômica.
func gravar(vault string, v *Vocab, fila []Pendente) error {
	files := []processo.File{}
	if v != nil {
		b, _ := json.MarshalIndent(v, "", "  ")
		files = append(files, processo.File{Kind: "vocab", Path: filepath.Join(vault, vocabFile), Data: append(b, '\n')})
	}
	if fila == nil {
		fila = []Pendente{}
	}
	b, _ := json.MarshalIndent(fila, "", "  ")
	files = append(files, processo.File{Kind: "pendente", Path: filepath.Join(vault, pendenteFile), Data: append(b, '\n')})
	st := processo.Store{Roots: processo.Roots{Vault: vault}}
	return st.Write(context.Background(), files, processo.WriteOpts{Force: true})
}
