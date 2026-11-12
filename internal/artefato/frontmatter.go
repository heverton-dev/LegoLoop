package artefato

import (
	"fmt"
	"strings"
)

type FrontMatter struct {
	ID             string
	Status         string
	Date           string
	SchemaVersion  string
	Upstream       []string
	Covers         []string
	Interpretative []string
}

var camposObrigatorios = []string{"id", "status", "date"}

func ParseFrontMatter(conteudo string) (FrontMatter, error) {
	campos, err := extrairCampos(conteudo)
	if err != nil {
		return FrontMatter{}, err
	}
	for _, nome := range camposObrigatorios {
		if campos[nome] == "" {
			return FrontMatter{}, fmt.Errorf("front matter: campo ausente: %s", nome)
		}
	}
	return FrontMatter{
		ID:             campos["id"],
		Status:         campos["status"],
		Date:           campos["date"],
		SchemaVersion:  campos["schema_version"],
		Upstream:       lista(campos["upstream"]),
		Covers:         lista(campos["covers"]),
		Interpretative: lista(campos["interpretative"]),
	}, nil
}

func extrairCampos(conteudo string) (map[string]string, error) {
	linhas := strings.Split(strings.ReplaceAll(conteudo, "\r\n", "\n"), "\n")
	if len(linhas) == 0 || strings.TrimSpace(linhas[0]) != "---" {
		return nil, fmt.Errorf("front matter: delimitador inicial ausente")
	}
	campos := map[string]string{}
	for _, linha := range linhas[1:] {
		if strings.TrimSpace(linha) == "---" {
			return campos, nil
		}
		chave, valor, ok := strings.Cut(linha, ":")
		if ok {
			campos[strings.TrimSpace(chave)] = strings.TrimSpace(valor)
		}
	}
	return nil, fmt.Errorf("front matter: delimitador final ausente")
}

func lista(valor string) []string {
	valor = strings.TrimSpace(strings.Trim(valor, "[]"))
	if valor == "" {
		return nil
	}
	var itens []string
	for _, item := range strings.Split(valor, ",") {
		itens = append(itens, strings.TrimSpace(item))
	}
	return itens
}
