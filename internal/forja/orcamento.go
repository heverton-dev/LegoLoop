package forja

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// OrcamentoTokens define o teto maximo de tokens permitido para cada harness.
type OrcamentoTokens struct {
	Versao    string         `json:"versao,omitempty"`
	Origem    string         `json:"origem,omitempty"`
	Unidade   string         `json:"unidade,omitempty"`
	Descricao string         `json:"descricao,omitempty"`
	Tetos     map[string]int `json:"tetos"`
	Harnesses map[string]int `json:"harnesses,omitempty"`
}

// ExcessoOrcamento registra um harness cujo custo medido superou o teto.
type ExcessoOrcamento struct {
	Harness string `json:"harness"`
	Medido  int    `json:"medido"`
	Teto    int    `json:"teto"`
	Excesso int    `json:"excesso"`
}

// ResultadoOrcamento resume a verificacao do orcamento de tokens.
type ResultadoOrcamento struct {
	Aprovado bool               `json:"aprovado"`
	Excessos []ExcessoOrcamento `json:"excessos,omitempty"`
}

// CarregarOrcamento le e valida um arquivo JSON de orcamento de tokens.
func CarregarOrcamento(caminho string) (OrcamentoTokens, error) {
	if strings.TrimSpace(caminho) == "" {
		return OrcamentoTokens{}, errors.New("forja: caminho do orcamento vazio")
	}
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return OrcamentoTokens{}, fmt.Errorf("forja: ler orcamento %s: %w", caminho, err)
	}
	var orc OrcamentoTokens
	if err := json.Unmarshal(dados, &orc); err != nil {
		return OrcamentoTokens{}, fmt.Errorf("forja: decodificar orcamento %s: %w", caminho, err)
	}
	if len(orc.Tetos) == 0 && len(orc.Harnesses) > 0 {
		orc.Tetos = orc.Harnesses
	}
	if err := validarOrcamento(orc); err != nil {
		return OrcamentoTokens{}, err
	}
	return orc, nil
}

func validarOrcamento(orc OrcamentoTokens) error {
	if len(orc.Tetos) == 0 {
		return errors.New("forja: orcamento sem tetos definidos")
	}
	for harness, teto := range orc.Tetos {
		if strings.TrimSpace(harness) == "" {
			return errors.New("forja: harness vazio no orcamento")
		}
		if teto < 0 {
			return fmt.Errorf("forja: teto negativo para harness %s (%d)", harness, teto)
		}
	}
	return nil
}

// VerificarOrcamento compara o custo medido por harness contra o teto definido.
func VerificarOrcamento(orc OrcamentoTokens, medidos map[string]int) ResultadoOrcamento {
	nomes := make([]string, 0, len(orc.Tetos))
	for h := range orc.Tetos {
		nomes = append(nomes, h)
	}
	slices.Sort(nomes)
	var excessos []ExcessoOrcamento
	for _, h := range nomes {
		teto := orc.Tetos[h]
		medido := medidos[h]
		if medido > teto {
			excessos = append(excessos, ExcessoOrcamento{
				Harness: h,
				Medido:  medido,
				Teto:    teto,
				Excesso: medido - teto,
			})
		}
	}
	return ResultadoOrcamento{
		Aprovado: len(excessos) == 0,
		Excessos: excessos,
	}
}

// FormatarExcessos gera texto legivel listando harness, medido, teto e excesso.
func FormatarExcessos(res ResultadoOrcamento) string {
	if res.Aprovado || len(res.Excessos) == 0 {
		return "orcamento aprovado: nenhum harness excedeu o teto"
	}
	linhas := make([]string, 0, len(res.Excessos))
	for _, ex := range res.Excessos {
		linhas = append(linhas, fmt.Sprintf(
			"harness=%s medido=%d teto=%d excesso=%d",
			ex.Harness, ex.Medido, ex.Teto, ex.Excesso,
		))
	}
	return strings.Join(linhas, "\n")
}
