package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/forja"
	"github.com/heverton-dev/LegoLoop/internal/instalar"
)

func cmdInstalar(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("instalar", d)
	dryRun := fs.Bool("dry-run", false, "simula a instalacao sem gravar no disco")
	apply := fs.Bool("apply", false, "aplica a instalacao e grava os relatorios")
	perfil := fs.String("perfil", "padrao", "perfil de contexto: minimo, padrao ou gateway")
	homeFlag := fs.String("home", "", "diretorio HOME alvo (padrao: HOME do usuario)")
	asJSON := fs.Bool("json", false, "emite o resultado em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *dryRun && *apply {
		fmt.Fprintln(d.Stderr, "erro: use apenas --dry-run ou --apply")
		return 2
	}
	home, err := resolverHomeInstalar(d, *homeFlag)
	if err != nil {
		return fail(d, err)
	}
	res, err := instalar.Executar(ctx, instalar.Opcao{
		Home:    home,
		Perfil:  *perfil,
		Aplicar: *apply,
	})
	if err != nil {
		return tratarErroInstalar(d, err)
	}
	return exibirResultadoInstalar(d, res, *asJSON)
}

func resolverHomeInstalar(d Deps, flagHome string) (string, error) {
	if strings.TrimSpace(flagHome) != "" {
		return flagHome, nil
	}
	if strings.TrimSpace(d.Home) != "" {
		return d.Home, nil
	}
	return os.UserHomeDir()
}

func tratarErroInstalar(d Deps, err error) int {
	if errors.Is(err, forja.ErrPerfilDesconhecido) || errors.Is(err, forja.ErrSemHome) {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	return fail(d, err)
}

func exibirResultadoInstalar(d Deps, res instalar.Resultado, asJSON bool) int {
	if asJSON {
		return printJSON(d, res)
	}
	modo := "DRY-RUN"
	if res.Aplicado {
		modo = "APLICADO"
	}
	fmt.Fprintf(d.Stdout, "Instalar (%s): veredito=%s perfil=%s etapas=%d\n",
		modo, res.Relatorio.Veredito, res.Relatorio.Perfil, len(res.Relatorio.Etapas))
	if res.CaminhoJSON != "" {
		fmt.Fprintf(d.Stdout, "Relatorio JSON: %s\n", res.CaminhoJSON)
	}
	if res.CaminhoMD != "" {
		fmt.Fprintf(d.Stdout, "Relatorio MD: %s\n", res.CaminhoMD)
	}
	if res.AtalhoGrafo != "" {
		fmt.Fprintf(d.Stdout, "Atalho do grafo: %s\n", res.AtalhoGrafo)
	}
	return 0
}
