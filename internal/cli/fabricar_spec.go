package cli

import (
	"context"
	"fmt"

	"github.com/heverton-dev/LegoLoop/internal/brain"
	"github.com/heverton-dev/LegoLoop/internal/interview"
	"github.com/heverton-dev/LegoLoop/internal/specgen"
	"github.com/heverton-dev/LegoLoop/internal/upstream"
)

const fabricarSpecUso = "uso: legoloop fabricar spec --intent INTENT [--raiz DIR] [--json]"

type opcoesFabricarSpec struct {
	raiz, intent string
	asJSON       bool
}

func parseFabricarSpecFlags(args []string, d Deps) (opcoesFabricarSpec, bool) {
	var o opcoesFabricarSpec
	fs := newFlags("fabricar spec", d)
	fs.StringVar(&o.raiz, "raiz", ".", "raiz do projeto (docs/specs fica sob ela)")
	fs.StringVar(&o.intent, "intent", "", "caminho do arquivo INTENT aprovado")
	fs.BoolVar(&o.asJSON, "json", false, "imprime o resultado em formato JSON")
	if err := fs.Parse(args); err != nil || o.intent == "" {
		fmt.Fprintln(d.Stderr, fabricarSpecUso)
		return o, false
	}
	return o, true
}

func obterLLMFabricarSpec(ctx context.Context, d Deps) (interview.LLMClient, error) {
	if d.Fabricar != nil && d.Fabricar.LLM != nil {
		return d.Fabricar.LLM, nil
	}
	b, _, err := brain.Pick(d.discover(ctx), "", d.brainOpts()...)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// cmdFabricarSpec orquestra a geração formal de SPEC a partir de uma INTENT aprovada.
func cmdFabricarSpec(ctx context.Context, args []string, d Deps) int {
	o, ok := parseFabricarSpecFlags(args, d)
	if !ok {
		return 2
	}
	llm, err := obterLLMFabricarSpec(ctx, d)
	if err != nil {
		return fail(d, err)
	}
	gen := specgen.NewGenerator(llm, nil)
	res, err := gen.Generate(ctx, o.raiz, o.intent)
	if err != nil {
		if upstream.Recusa(err) {
			fmt.Fprintln(d.Stderr, "erro:", err)
			return 2
		}
		return fail(d, err)
	}
	if o.asJSON {
		return printJSON(d, res)
	}
	fmt.Fprintf(d.Stdout, "SPEC gravada em %s\n", res.FilePath)
	return 0
}
