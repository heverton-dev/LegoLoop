package ste

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	cabecalhoRespostas = "### Respostas ao Usuário\n"
	bomUTF8            = "\xef\xbb\xbf"
)

func marcadoresSecao(nome string) (string, string) {
	if nome == "stept" {
		return "<!-- stept:inicio (gerado por tools/steptsync/gemini_md.py; não edite) -->", "<!-- stept:fim -->"
	}
	ini := fmt.Sprintf("<!-- %s:inicio (gerado por legoloop; não edite) -->", nome)
	fim := fmt.Sprintf("<!-- %s:fim -->", nome)
	return ini, fim
}

func removerSecaoGerenciada(texto, nome string) string {
	esc := regexp.QuoteMeta(nome)
	padrao := fmt.Sprintf(`(?s)<!-- %s:inicio[^>]*-->.*?<!-- %s:fim -->\n?`, esc, esc)
	re := regexp.MustCompile(padrao)
	return re.ReplaceAllString(texto, "")
}

func removerSecaoRespostas(texto string) string {
	for {
		inicio := strings.Index(texto, cabecalhoRespostas)
		if inicio < 0 {
			return texto
		}
		resto := texto[inicio+len(cabecalhoRespostas):]
		prox := strings.Index(resto, "\n### ")
		if prox < 0 {
			return texto[:inicio]
		}
		texto = texto[:inicio] + resto[prox:]
	}
}

func montarConteudoGerenciado(atual, nome, conteudo string) string {
	texto := strings.TrimPrefix(atual, bomUTF8)
	corpo := strings.TrimSpace(conteudo)
	texto = removerSecaoGerenciada(texto, nome)
	if nome == "stept" {
		texto = removerSecaoRespostas(texto)
	}
	ini, fim := marcadoresSecao(nome)
	secao := ini + "\n" + corpo + "\n" + fim + "\n"
	base := strings.TrimRight(texto, "\n")
	if base == "" {
		return secao
	}
	return base + "\n\n" + secao
}

// AplicarBloco aplica de forma idempotente uma seção gerenciada em arquivo de texto.
func AplicarBloco(arquivo, nome, conteudo string) error {
	if strings.TrimSpace(arquivo) == "" || strings.TrimSpace(nome) == "" {
		return ErrInvalido
	}
	dados, err := os.ReadFile(arquivo)
	existe := true
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		existe = false
	}
	novo := montarConteudoGerenciado(string(dados), nome, conteudo)
	if existe && string(dados) == novo {
		return nil
	}
	return gravarRegraAtomica(arquivo, novo)
}
