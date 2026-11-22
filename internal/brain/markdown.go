package brain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/escrita"
)

// AntibodiesMarkdownFile é o nome do espelho executivo para leitura rápida de agentes.
const AntibodiesMarkdownFile = "ANTICORPOS.md"

// SincronizarProcesso lê o anticorpos.json de dir e projeta o respectivo ANTICORPOS.md.
func SincronizarProcesso(dir string) error {
	caminhoJSON := filepath.Join(dir, AntibodiesDBFile)
	dados, err := os.ReadFile(caminhoJSON)
	if err != nil {
		return fmt.Errorf("brain: lendo %s: %w", AntibodiesDBFile, err)
	}
	var banco bancoJSON
	if err := json.Unmarshal(dados, &banco); err != nil {
		return fmt.Errorf("brain: decodificando %s: %w", AntibodiesDBFile, err)
	}
	caminhoMD := filepath.Join(dir, AntibodiesMarkdownFile)
	md := RenderizarMarkdown(banco.Processo, banco.Globais, banco.Dominio)
	if err := escrita.GravarAtomico(caminhoMD, md); err != nil {
		return fmt.Errorf("brain: gravando %s: %w", AntibodiesMarkdownFile, err)
	}
	return nil
}

// SincronizarTodos varre recursivamente a raiz e sincroniza todos os anticorpos.json encontrados.
func SincronizarTodos(raiz string) (int, error) {
	total := 0
	err := filepath.Walk(raiz, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || info.Name() != AntibodiesDBFile {
			return nil
		}
		if err := SincronizarProcesso(filepath.Dir(path)); err == nil {
			total++
		}
		return nil
	})
	return total, err
}

// RenderizarMarkdown gera o espelho executivo em Markdown com cabeçalho YAML ultra-compacto.
func RenderizarMarkdown(processo string, globais, dominio []vacinaJSON) []byte {
	var sb strings.Builder
	escreverCabecalhoYAML(&sb, processo, globais, dominio)
	escreverCorpoMarkdown(&sb, processo, globais, dominio)
	return []byte(sb.String())
}

// escreverCabecalhoYAML gera o frontmatter YAML de leitura telegráfica para consumo de agentes.
func escreverCabecalhoYAML(sb *strings.Builder, processo string, globais, dominio []vacinaJSON) {
	total := len(globais) + len(dominio)
	sb.WriteString("---\n")
	sb.WriteString("schema: anticorpos.v1\n")
	fmt.Fprintf(sb, "processo: %s\n", processo)
	fmt.Fprintf(sb, "total_anticorpos: %d\n", total)
	fmt.Fprintf(sb, "globais: %d\n", len(globais))
	fmt.Fprintf(sb, "dominio: %d\n", len(dominio))
	sb.WriteString("regras_inviolaveis:\n")
	for _, v := range dominio {
		fmt.Fprintf(sb, "  - pacote: \"%s\"\n    antidoto: \"%s\"\n", escaparYAML(v.Pacote), escaparYAML(resumirAntidoto(v.Antidoto)))
	}
	for _, v := range globais {
		fmt.Fprintf(sb, "  - pacote: \"%s\"\n    antidoto: \"%s\"\n", escaparYAML(v.Pacote), escaparYAML(resumirAntidoto(v.Antidoto)))
	}
	sb.WriteString("---\n\n")
}

// escreverCorpoMarkdown gera o detalhamento em seções para consulta sob demanda.
func escreverCorpoMarkdown(sb *strings.Builder, processo string, globais, dominio []vacinaJSON) {
	fmt.Fprintf(sb, "# 🛡️ Banco de Anticorpos — %s\n\n", processo)
	sb.WriteString("> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.\n\n")
	escreverSecaoVacinas(sb, "Regras de Domínio", dominio)
	escreverSecaoVacinas(sb, "Regras Globais", globais)
}

// escreverSecaoVacinas renderiza uma lista de vacinas em tópicos estruturados.
func escreverSecaoVacinas(sb *strings.Builder, titulo string, vacinas []vacinaJSON) {
	fmt.Fprintf(sb, "## %s (%d)\n\n", titulo, len(vacinas))
	if len(vacinas) == 0 {
		sb.WriteString("_Nenhuma regra cadastrada nesta seção._\n\n")
		return
	}
	for _, v := range vacinas {
		fmt.Fprintf(sb, "### `%s`\n", v.Pacote)
		fmt.Fprintf(sb, "- **Antídoto:** %s\n", v.Antidoto)
		fmt.Fprintf(sb, "- **Erro Evitado:** %s\n", v.Erro)
		if v.CausaRaiz != "" {
			fmt.Fprintf(sb, "- **Causa Raiz:** %s\n", v.CausaRaiz)
		}
		fmt.Fprintf(sb, "- **Assinatura:** `%s`\n\n", v.Assinatura)
	}
}

// resumirAntidoto encurta o antídoto para caber no cabeçalho YAML em uma linha concisa.
func resumirAntidoto(antidoto string) string {
	limpo := strings.ReplaceAll(antidoto, "\n", " ")
	limpo = strings.TrimSpace(limpo)
	runas := []rune(limpo)
	if len(runas) > 120 {
		return string(runas[:117]) + "..."
	}
	return limpo
}

// escaparYAML escapa aspas duplas para inserção segura em string YAML.
func escaparYAML(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
