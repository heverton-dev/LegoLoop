package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"legoloop/internal/brain"
	"legoloop/internal/caminhos"
)

const anticorposUso = `uso: legoloop anticorpos <sync|status> [slug] [--raiz DIR] [--json]

Comandos:
  sync    Sincroniza e projeta o arquivo ANTICORPOS.md a partir de anticorpos.json.
  status  Exibe o resumo executivo de imunidade e anticorpos cadastrados.
`

func cmdAnticorpos(_ context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		fmt.Fprint(d.Stderr, anticorposUso)
		return 2
	}
	subcomando := args[0]
	resto := args[1:]
	fs := flag.NewFlagSet("anticorpos "+subcomando, flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	var raiz string
	var asJSON bool
	fs.StringVar(&raiz, "raiz", ".", "diretório raiz para busca ou sincronização")
	fs.BoolVar(&asJSON, "json", false, "emite a resposta em formato JSON")
	if err := fs.Parse(resto); err != nil {
		return 2
	}
	switch subcomando {
	case "sync":
		return executarSync(fs.Args(), raiz, d)
	case "status":
		return executarStatus(fs.Args(), raiz, d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido: %s\n\n%s", subcomando, anticorposUso)
		return 2
	}
}

func executarSync(posicionais []string, raiz string, d Deps) int {
	if len(posicionais) > 0 {
		slug := posicionais[0]
		dir := resolverDirAlvo(slug, raiz)
		if err := brain.SincronizarProcesso(dir); err != nil {
			fmt.Fprintf(d.Stderr, "erro sincronizando %s: %v\n", slug, err)
			return 1
		}
		fmt.Fprintf(d.Stdout, "✓ Processo %s sincronizado com sucesso: ANTICORPOS.md gerado.\n", slug)
		return 0
	}
	total, err := brain.SincronizarTodos(raiz)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro na sincronização em lote: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "✓ Sincronização concluída: %d processos atualizados com ANTICORPOS.md.\n", total)
	return 0
}

func executarStatus(posicionais []string, raiz string, d Deps) int {
	if len(posicionais) == 0 {
		fmt.Fprintln(d.Stderr, "erro: informe o slug do processo para ver o status.")
		return 2
	}
	slug := posicionais[0]
	dir := resolverDirAlvo(slug, raiz)
	caminhoMD := filepath.Join(dir, brain.AntibodiesMarkdownFile)
	dados, err := os.ReadFile(caminhoMD)
	if err != nil {
		fmt.Fprintf(d.Stderr, "anticorpos não encontrados para %s: %v\n", slug, err)
		return 1
	}
	fmt.Fprintln(d.Stdout, string(dados))
	return 0
}

func resolverDirAlvo(slug, raiz string) string {
	candidatos := []string{
		filepath.Join(raiz, slug),
		filepath.Join(raiz, "processos", slug),
		filepath.Join(raiz, caminhos.DirProcessos, slug),
	}
	for _, c := range candidatos {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return filepath.Join(raiz, slug)
}
