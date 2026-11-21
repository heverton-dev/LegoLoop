package espinha

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var padraoLeiNumerada = regexp.MustCompile(`^(\d+)\.\s+`)

// ExtrairLeisAgents extrai os numeros das leis de nivel superior da secao 1 do AGENTS.md.
func ExtrairLeisAgents(conteudo string) []int {
	linhas := strings.Split(strings.ReplaceAll(conteudo, "\r\n", "\n"), "\n")
	emSecao1 := false
	var leis []int
	for _, linha := range linhas {
		if strings.HasPrefix(linha, "## 1.") {
			emSecao1 = true
			continue
		}
		if emSecao1 && strings.HasPrefix(linha, "## ") {
			break
		}
		if !emSecao1 {
			continue
		}
		m := padraoLeiNumerada.FindStringSubmatch(linha)
		if m == nil {
			continue
		}
		if num, err := strconv.Atoi(m[1]); err == nil {
			leis = append(leis, num)
		}
	}
	return leis
}

// ConferirLeis garante que toda lei numerada do AGENTS.md consta no Registro
// com verificador (quando mecanica) ou justificativa (quando nao mecanica).
func ConferirLeis(leis []int, reg Registro) error {
	if len(leis) == 0 {
		return errors.New("nenhuma lei encontrada no AGENTS.md")
	}
	indice, err := indexarRequisitos(reg.Requisitos)
	if err != nil {
		return err
	}
	for _, num := range leis {
		chave := fmt.Sprintf("Lei %d", num)
		req, ok := indice[chave]
		if !ok {
			return fmt.Errorf("%s ausente em espinha/registro.json", chave)
		}
		if err := validarRequisitoLei(chave, req); err != nil {
			return err
		}
	}
	return nil
}

func indexarRequisitos(reqs []Requisito) (map[string]Requisito, error) {
	indice := make(map[string]Requisito, len(reqs))
	for _, r := range reqs {
		chave := strings.TrimSpace(r.Lei)
		if chave == "" {
			return nil, errors.New("requisito com lei vazia")
		}
		if _, dup := indice[chave]; dup {
			return nil, fmt.Errorf("requisito repetido: %q", chave)
		}
		indice[chave] = r
	}
	return indice, nil
}

func validarRequisitoLei(chave string, req Requisito) error {
	if req.Mecanico && strings.TrimSpace(req.Verificador) == "" {
		return fmt.Errorf("%s mecanica sem verificador", chave)
	}
	if !req.Mecanico && strings.TrimSpace(req.Justificativa) == "" {
		return fmt.Errorf("%s nao mecanica sem justificativa", chave)
	}
	return nil
}
