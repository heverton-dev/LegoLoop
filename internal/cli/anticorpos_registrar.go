package cli

import (
	"flag"
	"fmt"

	"legoloop/internal/brain"
)

// dadosAnticorpo reúne as opções do subcomando registrar.
type dadosAnticorpo struct{ pacote, erro, causa, antidoto, escopo string }

// parseIntercalado lê as opções em qualquer posição: o pacote flag para na primeira palavra solta,
// e o uso do comando manda pôr o slug antes das opções.
func parseIntercalado(fs *flag.FlagSet, args []string) ([]string, error) {
	var soltos []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		if fs.NArg() == 0 {
			return soltos, nil
		}
		soltos = append(soltos, fs.Arg(0))
		args = fs.Args()[1:]
	}
}

// executarRegistrar grava o anticorpo no anticorpos.json do Processo e regenera o ANTICORPOS.md.
func executarRegistrar(posicionais []string, raiz string, r dadosAnticorpo, d Deps) int {
	if len(posicionais) == 0 || r.pacote == "" || r.erro == "" {
		fmt.Fprintln(d.Stderr, "erro: informe o slug, --pacote e --erro para registrar um anticorpo.")
		return 2
	}
	escopo := brain.EscopoDominio
	if r.escopo == "global" {
		escopo = brain.EscopoGlobal
	}
	slug := posicionais[0]
	dir := resolverDirAlvo(slug, raiz)
	db, err := brain.NewAntibodiesDBPersistente(dir)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro abrindo anticorpos de %s: %v\n", slug, err)
		return 1
	}
	if err := db.RecordComEscopo(escopo, r.pacote, r.erro, r.causa, r.antidoto); err != nil {
		fmt.Fprintf(d.Stderr, "erro registrando anticorpo: %v\n", err)
		return 1
	}
	if err := brain.SincronizarProcesso(dir); err != nil {
		fmt.Fprintf(d.Stderr, "erro sincronizando %s: %v\n", slug, err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "✓ Anticorpo registrado em %s (escopo %s).\n", slug, r.escopo)
	return 0
}
