package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"legoloop/internal/dissecar"
)

type DissecarCommand struct {
	fs     *flag.FlagSet
	output io.Writer
}

func NewDissecarCommand() *DissecarCommand {
	return &DissecarCommand{
		output: os.Stdout,
	}
}

func (c *DissecarCommand) Run(args []string) error {
	fs := flag.NewFlagSet("dissecar", flag.ContinueOnError)
	fs.SetOutput(c.output)

	var target string
	var outputJSON bool

	fs.StringVar(&target, "target", "", "Caminho do codigo-fonte a dissecar")
	fs.BoolVar(&outputJSON, "json", false, "Exibir saida em formato JSON headless")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if target == "" {
		return fmt.Errorf("parametro --target e obrigatorio")
	}

	return c.executarDisseccao(target, outputJSON)
}

func (c *DissecarCommand) executarDisseccao(target string, outputJSON bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	safePath, err := dissecar.ValidatePath(cwd, target)
	if err != nil {
		return fmt.Errorf("falha no path jail: %w", err)
	}

	astNodes, err := dissecar.ExtractAST(safePath)
	if err != nil {
		return fmt.Errorf("erro ao extrair AST: %w", err)
	}

	syscalls, err := dissecar.ExtractSyscalls(safePath)
	if err != nil {
		return fmt.Errorf("erro ao extrair syscalls: %w", err)
	}

	result := dissecar.Result{
		TargetID:   safePath,
		AstNodes:   astNodes,
		Syscalls:   syscalls,
		ReverseDAG: dissecar.BuildReverseDAG(astNodes),
	}

	return c.renderizarResultado(result, outputJSON)
}

func (c *DissecarCommand) renderizarResultado(result dissecar.Result, outputJSON bool) error {
	if outputJSON {
		enc := json.NewEncoder(c.output)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	_, _ = fmt.Fprintln(c.output, "=== RELATORIO DE DISSECCAO ===")
	_, _ = fmt.Fprintf(c.output, "Alvo: %s\n", result.TargetID)
	_, _ = fmt.Fprintf(c.output, "Nos AST: %d\n", len(result.AstNodes))
	_, _ = fmt.Fprintf(c.output, "Syscalls: %d\n", len(result.Syscalls))
	return nil
}
