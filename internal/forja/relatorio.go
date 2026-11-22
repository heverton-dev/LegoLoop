package forja

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const marcadorRedigido = "[REDIGIDO]"

// ColetarSegredosAmbiente extrai valores não vazios de variáveis cujas chaves
// contenham SECRET, TOKEN, KEY ou PASSWORD (sem distinção de caixa).
func ColetarSegredosAmbiente(env []string) []string {
	vistos := make(map[string]struct{})
	var lista []string
	for _, entrada := range env {
		chave, valor, ok := strings.Cut(entrada, "=")
		if !ok || valor == "" || !chaveEhSecreta(chave) {
			continue
		}
		if _, existe := vistos[valor]; existe {
			continue
		}
		vistos[valor] = struct{}{}
		lista = append(lista, valor)
	}
	ordenarPorTamanhoDesc(lista)
	return lista
}

// SanitizarTexto substitui cada segredo não vazio em texto por [REDIGIDO].
func SanitizarTexto(texto string, segredos []string) string {
	if texto == "" || len(segredos) == 0 {
		return texto
	}
	ordenados := append([]string(nil), segredos...)
	ordenarPorTamanhoDesc(ordenados)
	resultado := texto
	for _, segredo := range ordenados {
		if segredo == "" {
			continue
		}
		resultado = strings.ReplaceAll(resultado, segredo, marcadorRedigido)
	}
	return resultado
}

// RenderizarRelatorioJSON serializa dados em JSON indentado, garantindo que
// nenhum valor de segredos apareça no resultado.
func RenderizarRelatorioJSON(dados any, segredos []string) ([]byte, error) {
	bruto, err := json.Marshal(dados)
	if err != nil {
		return nil, err
	}
	var generico any
	if err := json.Unmarshal(bruto, &generico); err != nil {
		return nil, err
	}
	limpo := sanitizarValorJSON(generico, segredos)
	return json.MarshalIndent(limpo, "", "  ")
}

// RenderizarRelatorioMD gera o relatório Markdown a partir de dados,
// redigindo todos os segredos antes de retornar.
func RenderizarRelatorioMD(dados any, segredos []string) string {
	if texto, ok := dados.(string); ok {
		return SanitizarTexto(texto, segredos)
	}
	jsonLimpo, err := RenderizarRelatorioJSON(dados, segredos)
	if err != nil {
		return SanitizarTexto(fmt.Sprintf("%v", dados), segredos)
	}
	var sb strings.Builder
	sb.WriteString("# Relatório da Forja\n\n```json\n")
	sb.Write(jsonLimpo)
	sb.WriteString("\n```\n")
	return SanitizarTexto(sb.String(), segredos)
}

func chaveEhSecreta(chave string) bool {
	alta := strings.ToUpper(chave)
	for _, frag := range []string{"SECRET", "TOKEN", "KEY", "PASSWORD"} {
		if strings.Contains(alta, frag) {
			return true
		}
	}
	return false
}

func ordenarPorTamanhoDesc(itens []string) {
	sort.SliceStable(itens, func(i, j int) bool {
		return len(itens[i]) > len(itens[j])
	})
}

func sanitizarValorJSON(v any, segredos []string) any {
	switch atual := v.(type) {
	case string:
		return SanitizarTexto(atual, segredos)
	case []any:
		saida := make([]any, len(atual))
		for i, item := range atual {
			saida[i] = sanitizarValorJSON(item, segredos)
		}
		return saida
	case map[string]any:
		saida := make(map[string]any, len(atual))
		for k, val := range atual {
			saida[SanitizarTexto(k, segredos)] = sanitizarValorJSON(val, segredos)
		}
		return saida
	default:
		return v
	}
}
