package cli

import (
	"context"
	"fmt"

	"legoloop/internal/caminhos"
)

const execUsage = `uso: legoloop exec <slug> [--processos-dir DIR] [--workdir DIR] [--json]

Executa o processo <slug> pela via CLI: resolve processos/<slug>, roda o DAG
com as opções de run e grava os artefatos de saída no disco.
`

// cmdExec é a via CLI da tríplice via de invocação: executa o processo pelo
// slug, sem passo manual, com o mesmo contrato do comando run.
func cmdExec(ctx context.Context, args []string, d Deps) int {
	slug, resto := splitPositional(args)
	f, err := parseRunFlags("exec", resto, d)
	if err != nil {
		return 2
	}
	if slug == "" {
		fmt.Fprintf(d.Stderr, "erro: informe o slug do processo.\n\n%s", execUsage)
		return 2
	}
	dir := resolverDirProcesso(slug, f.processosDir)
	if dir == "" {
		fmt.Fprintf(d.Stderr, "exec: processo %q não encontrado\n", slug)
		return 1
	}
	p, err := loadDAG(caminhos.CaminhoDAGLeitura(dir))
	if err != nil {
		return fail(d, err)
	}
	return rodarDAG(ctx, d, p, f)
}
