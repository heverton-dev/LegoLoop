package triagem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// ArquivoTriagem é o nome do arquivo da Peça de Triagem dentro da pasta do processo.
const ArquivoTriagem = "triagem.json"

// Exigida diz se um processo com esse número de nós precisa de triagem. A triagem escolhe qual nó o pedido
// chama; com um nó só não há o que escolher.
func Exigida(nos int) bool { return nos > 1 }

// Completar devolve a triagem com uma regra para cada etapa. Mantém as regras que já existem (e suas edições)
// e acrescenta as que faltam, na ordem de etapas.
func Completar(s TriagemSpec, processo string, etapas []string) TriagemSpec {
	base := GerarTriagemPadrao(processo, etapas)
	if s.SchemaVersion == "" {
		return base
	}
	tem := make(map[string]bool, len(s.Etapas))
	for _, e := range s.Etapas {
		tem[e.ID] = true
	}
	s.Etapas = slices.Clone(s.Etapas)
	for _, e := range base.Etapas {
		if !tem[e.ID] {
			s.Etapas = append(s.Etapas, e)
		}
	}
	s.CriterioAmbiguidade.Etapas = slices.Clone(s.CriterioAmbiguidade.Etapas)
	for _, id := range etapas {
		if !slices.Contains(s.CriterioAmbiguidade.Etapas, id) {
			s.CriterioAmbiguidade.Etapas = append(s.CriterioAmbiguidade.Etapas, id)
		}
	}
	return s
}

// Conferir lê a triagem da pasta do processo e prova que ela é válida e cobre todos os nós.
func Conferir(dir string, nos []string) error {
	b, err := os.ReadFile(filepath.Join(dir, ArquivoTriagem))
	if err != nil {
		return fmt.Errorf("o processo tem %d nós e precisa de %s: %w", len(nos), ArquivoTriagem, err)
	}
	var s TriagemSpec
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("%s ilegível: %w", ArquivoTriagem, err)
	}
	if err := s.Validar(); err != nil {
		return fmt.Errorf("%s inválida: %w", ArquivoTriagem, err)
	}
	var falta []string
	for _, id := range nos {
		if !slices.ContainsFunc(s.Etapas, func(e RegraEtapa) bool { return e.ID == id }) {
			falta = append(falta, id)
		}
	}
	if len(falta) > 0 {
		return fmt.Errorf("%s não cobre os nós %v", ArquivoTriagem, falta)
	}
	return nil
}
