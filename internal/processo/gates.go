package processo

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"legoloop/internal/gate"
	"legoloop/internal/ste"
)

// PalavrasMinimas é o mínimo de palavras que o gate da biblioteca exige de um artefato Markdown.
const PalavrasMinimas = 30

// LinhasMinimasDados é o mínimo de linhas do gate de dados: o cabeçalho e ao menos um registro.
const LinhasMinimasDados = 2

// ErrGateTrivial: a etapa promete conteúdo, mas o gate só confere que o arquivo existe (ou nem isso).
var ErrGateTrivial = errors.New("processo: gate trivial")

func init() {
	ste.Default.Register("processo.gate_trivial", "A etapa %s promete conteúdo, mas o gate só confere que o arquivo existe.")
	ste.Default.Register("processo.gate_sugerido", "Use este gate: %s")
}

// tiposComConteudo são os artefatos em que a promessa de conteúdo faz sentido; marcador só precisa existir.
const tipoMarcador = "marcador"

// GateParaArtefato gera o gate da biblioteca para o artefato: o mínimo que prova que o conteúdo existe.
//   - markdown: existe, tem ao menos PalavrasMinimas palavras e não deixou {{placeholder}};
//   - json: existe e cumpre o esquema (se houver), senão começa como objeto ou lista JSON;
//   - codigo: existe e não está vazio; dados: existe e tem o cabeçalho e ao menos um registro;
//   - marcador: existe.
func GateParaArtefato(a Artefato) (string, error) {
	p, err := refGate("caminho", a.Caminho)
	if err != nil {
		return "", err
	}
	switch a.Tipo {
	case "markdown":
		return fmt.Sprintf(`file_exists(%[1]s) && word_count(%[1]s) >= %[2]d && !has_pattern(%[1]s, '\{\{')`, p, PalavrasMinimas), nil
	case "json":
		return gateJSON(p, a.Esquema)
	case "codigo", "dados":
		return gateCodigoOuDados(p, a.Tipo), nil
	case tipoMarcador:
		return fmt.Sprintf(`file_exists(%s)`, p), nil
	}
	return "", invalido("tipo de artefato %q desconhecido (use markdown, json, codigo, dados ou marcador)", a.Tipo)
}

func gateJSON(p, esquema string) (string, error) {
	if esquema == "" {
		return fmt.Sprintf(`file_exists(%[1]s) && has_pattern(%[1]s, '^\s*[\{\[]')`, p), nil
	}
	e, err := refGate("esquema", esquema)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`file_exists(%[1]s) && schema_valid(%[2]s, %[1]s)`, p, e), nil
}

// refGate valida um caminho que vai virar literal da expressão do gate e o devolve entre aspas simples.
// Aspa, barra invertida, crase e espaço não entram: o caminho nunca pode fechar o literal nem virar outro termo.
func refGate(campo, caminho string) (string, error) {
	if caminho == "" || strings.ContainsAny(caminho, "'\\`\r\n\t ") || !filepath.IsLocal(filepath.FromSlash(caminho)) || strings.HasPrefix(caminho, "/") {
		return "", invalido("%s %q não serve para o gate (use caminho relativo, com /, sem espaço nem aspa)", campo, caminho)
	}
	return "'" + caminho + "'", nil
}

// gateNaoTrivial recusa a etapa que promete conteúdo e usa um gate que não o confere.
func gateNaoTrivial(e Etapa, g *gate.Expr) error {
	if !e.Produz.Promete || e.Produz.Tipo == tipoMarcador || g.ChecksContent() {
		return nil
	}
	msg := ste.Msg("processo.gate_trivial", e.ID)
	if sugerido, err := GateParaArtefato(e.Produz); err == nil {
		msg += " " + ste.Msg("processo.gate_sugerido", sugerido)
	}
	return fmt.Errorf("%w: %w: %s", ErrInvalido, ErrGateTrivial, msg)
}

// gateCodigoOuDados: código não pode estar vazio; dados precisam de ao menos duas linhas (cabeçalho e um registro).
func gateCodigoOuDados(p, tipo string) string {
	if tipo == "dados" {
		return fmt.Sprintf(`file_exists(%[1]s) && line_count(%[1]s) >= %[2]d`, p, LinhasMinimasDados)
	}
	return fmt.Sprintf(`file_exists(%[1]s) && word_count(%[1]s) >= 1`, p)
}
