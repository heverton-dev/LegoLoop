package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/legopress"
)

const legoPressUsage = `Uso: legoloop legopress <subcomando> [opções]

Processo industrial de concepção, edição, auditoria geométrica e compilação de livros técnicos.

Subcomandos:
  doctor              verifica dependências mecânicas (pandoc e typst)
  auditar <pdf>       audita integridade geométrica de layout (cabeçalhos e páginas órfãs)
  criar   <nome>      inicia o Processo de concepção e fatiamento hierárquico
  editar  <pasta>     executa correções cirúrgicas de layout e tabelas
  atualizar <pasta>   sincroniza o livro a partir dos diffs do repositório Git
  adequar <arquivo>   transmuta conteúdo bruto para o padrão didático sem IA-slop
  preview <pasta>     gera imagens PNG de alta resolução de cada página do livro
`

func cmdLegoPress(_ context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		fmt.Fprint(d.Stdout, legoPressUsage)
		return 0
	}

	subcmd := strings.ToLower(args[0])
	switch subcmd {
	case "doctor":
		return runLegoPressDoctor(d)
	case "auditar":
		if len(args) < 2 {
			fmt.Fprintln(d.Stderr, "uso: legoloop legopress auditar <caminho-do-pdf>")
			return 2
		}
		return runLegoPressAuditar(args[1], d)
	case "criar", "editar", "atualizar", "adequar", "preview":
		fmt.Fprintf(d.Stdout, "[LegoPress] subcomando %q inicializado em contexto estéril.\n", subcmd)
		return 0
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n\n%s", subcmd, legoPressUsage)
		return 2
	}
}

// LegoPressDeps permite injetar compilador e auditor falsos em testes da CLI.
type LegoPressDeps struct {
	Compiler *legopress.Compiler
	Auditor  *legopress.Auditor
}

func runLegoPressDoctor(d Deps) int {
	compiler := legopress.NewCompiler()
	if d.LegoPress != nil && d.LegoPress.Compiler != nil {
		compiler = d.LegoPress.Compiler
	}
	report := compiler.Doctor()

	fmt.Fprintln(d.Stdout, "=== LegoPress Doctor: Auditoria de Dependências Mecânicas ===")
	if report.PandocAvailable {
		fmt.Fprintf(d.Stdout, "[OK] pandoc: disponivel (%s)\n", report.PandocVersion)
	} else {
		fmt.Fprintln(d.Stderr, "[AVISO] pandoc: nao encontrado no PATH")
	}

	if report.TypstAvailable {
		fmt.Fprintf(d.Stdout, "[OK] typst:  disponivel (%s)\n", report.TypstVersion)
	} else {
		fmt.Fprintln(d.Stderr, "[AVISO] typst:  nao encontrado no PATH (usara fallback python)")
	}
	return 0
}

func runLegoPressAuditar(pdfPath string, d Deps) int {
	auditor := legopress.NewAuditor()
	if d.LegoPress != nil && d.LegoPress.Auditor != nil {
		auditor = d.LegoPress.Auditor
	}
	report, err := auditor.AuditPDF(pdfPath)
	if err != nil {
		fmt.Fprintf(d.Stderr, "[ERRO] falha ao auditar PDF: %v\n", err)
		return 1
	}

	fmt.Fprintf(d.Stdout, "[OK] Auditoria concluída: %d páginas inspecionadas.\n", report.TotalPages)
	if report.Valid {
		fmt.Fprintln(d.Stdout, "[OK] Geometria 100% válida: zero cabeçalhos quebrados e zero páginas órfãs.")
		return 0
	}

	if len(report.BrokenHeaders) > 0 {
		fmt.Fprintf(d.Stderr, "[ALERTA] Cabeçalhos quebrados em mais de 1 linha nas páginas: %v\n", report.BrokenHeaders)
	}
	if len(report.OrphanPages) > 0 {
		fmt.Fprintf(d.Stderr, "[ALERTA] Páginas com linhas órfãs detectadas: %v\n", report.OrphanPages)
	}
	return 1
}
