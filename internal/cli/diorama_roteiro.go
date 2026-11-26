package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
)

const usoDioramaRoteiro = `Uso: legoloop diorama-roteiro <subcomando> [opções]

Subcomandos:
  prevoo           --projeto DIR [--acervo DIR] [--catalogo ARQ] [--json]
  validar-config   --projeto DIR [--catalogo ARQ] [--json]
  coletar          --projeto DIR [--acervo ARQ] [--catalogo ARQ] [--json]
  validar-roteiro  --projeto DIR [--catalogo ARQ] [--json]
`

func cmdDioramaRoteiro(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		fmt.Fprint(d.Stdout, usoDioramaRoteiro)
		return 0
	}
	switch args[0] {
	case "prevoo":
		return cmdDioramaPreVoo(ctx, args[1:], d)
	case "validar-config":
		return cmdDioramaValidarConfig(args[1:], d)
	case "coletar":
		return cmdDioramaColetar(ctx, args[1:], d)
	case "validar-roteiro":
		return cmdDioramaValidarRoteiro(args[1:], d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n%s", args[0], usoDioramaRoteiro)
		return 2
	}
}

func cmdDioramaValidarConfig(args []string, d Deps) int {
	projDir, _, catPath, emJSON, ok := parseDioramaFlags("validar-config", args, d.Stderr)
	if !ok {
		return 2
	}
	proj, marca, cat, err := carregarConfigDiorama(projDir, catPath)
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	if err := roteiro.ValidarConfiguracao(proj, marca, cat); err != nil {
		return responderDiorama(d, emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	return responderDiorama(d, emJSON, map[string]any{"ok": true, "nicho": proj.Nicho}, 0)
}

func cmdDioramaPreVoo(ctx context.Context, args []string, d Deps) int {
	projDir, acervoArg, catPath, emJSON, ok := parseDioramaFlags("prevoo", args, d.Stderr)
	if !ok {
		return 2
	}
	proj, _, cat, err := carregarConfigDiorama(projDir, catPath)
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"bloqueado": true, "erro": err.Error()}, 1)
	}
	acervoDir := acervoArg
	if acervoDir == "" {
		acervoDir = filepath.Join(projDir, "acervo")
	}
	rel := roteiro.ExecutarPreVoo(ctx, proj, cat, projDir, acervoDir, roteiro.DepsPreVoo{})
	code := 0
	if rel.Bloqueado {
		code = 1
	}
	return responderDiorama(d, emJSON, rel, code)
}

func cmdDioramaColetar(ctx context.Context, args []string, d Deps) int {
	projDir, acervoArg, catPath, emJSON, ok := parseDioramaFlags("coletar", args, d.Stderr)
	if !ok {
		return 2
	}
	proj, _, cat, err := carregarConfigDiorama(projDir, catPath)
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"erro": err.Error()}, 1)
	}
	acervoPath := acervoArg
	if acervoPath == "" {
		acervoPath = filepath.Join(projDir, "acervo", "temas.json")
	}
	acervo, err := roteiro.CarregarAcervoTemas(acervoPath)
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"erro": err.Error()}, 1)
	}
	pautas, err := roteiro.ColetarPautas(ctx, proj, cat, acervo, nil)
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"erro": err.Error()}, 1)
	}
	fila, filaPath, err := sincronizarFilaPerene(acervoPath, pautas, acervo)
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"erro": err.Error()}, 1)
	}
	combinadas := roteiro.CombinarPautasComFila(pautas, fila, acervo)
	return responderDiorama(d, emJSON, map[string]any{
		"total":             len(combinadas),
		"pautas":            combinadas,
		"fila_perene_total": len(fila.Pautas),
		"fila_perene_path":  filaPath,
	}, 0)
}

func sincronizarFilaPerene(acervoPath string, pautas []roteiro.PautaCandidata, acervo roteiro.AcervoTemas) (roteiro.FilaPerene, string, error) {
	filaPath := roteiro.CaminhoFilaPerene(acervoPath)
	fila, err := roteiro.CarregarFilaPerene(filaPath)
	if err != nil {
		return fila, filaPath, err
	}
	fila = roteiro.MesclarFilaPerene(fila, pautas, acervo)
	if len(fila.Pautas) > 0 {
		if err := roteiro.SalvarFilaPerene(filaPath, fila); err != nil {
			return fila, filaPath, err
		}
	}
	return fila, filaPath, nil
}

func cmdDioramaValidarRoteiro(args []string, d Deps) int {
	projDir, _, catPath, emJSON, ok := parseDioramaFlags("validar-roteiro", args, d.Stderr)
	if !ok {
		return 2
	}
	proj, cat, err := carregarProjetoECatalogo(projDir, catPath)
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"aprovado": false, "erro": err.Error()}, 1)
	}
	var sc roteiro.Script
	if err := lerJSONArquivo(filepath.Join(projDir, "script.json"), &sc); err != nil {
		return responderDiorama(d, emJSON, map[string]any{"aprovado": false, "erro": err.Error()}, 1)
	}
	rel := roteiro.ValidarRoteiro(sc, proj, cat)
	code := 0
	if !rel.Aprovado {
		code = 1
	}
	return responderDiorama(d, emJSON, rel, code)
}

func parseDioramaFlags(nome string, args []string, errw io.Writer) (string, string, string, bool, bool) {
	fs := flag.NewFlagSet("diorama-roteiro "+nome, flag.ContinueOnError)
	fs.SetOutput(errw)
	proj := fs.String("projeto", ".", "pasta do projeto")
	acervo := fs.String("acervo", "", "pasta ou arquivo do acervo")
	cat := fs.String("catalogo", "", "arquivo JSON de catalogo de fontes")
	emJSON := fs.Bool("json", false, "saida em JSON")
	if err := fs.Parse(args); err != nil {
		return "", "", "", false, false
	}
	return *proj, *acervo, *cat, *emJSON, true
}

func carregarConfigDiorama(projDir, catPath string) (roteiro.Projeto, roteiro.Marca, roteiro.CatalogoFontes, error) {
	proj, cat, err := carregarProjetoECatalogo(projDir, catPath)
	if err != nil {
		return roteiro.Projeto{}, roteiro.Marca{}, roteiro.CatalogoFontes{}, err
	}
	var marca roteiro.Marca
	if err := lerJSONArquivo(filepath.Join(projDir, "marca.json"), &marca); err != nil {
		return roteiro.Projeto{}, roteiro.Marca{}, roteiro.CatalogoFontes{}, err
	}
	return proj, marca, cat, nil
}

func carregarProjetoECatalogo(projDir, catPath string) (roteiro.Projeto, roteiro.CatalogoFontes, error) {
	var proj roteiro.Projeto
	if err := lerJSONArquivo(filepath.Join(projDir, "projeto.json"), &proj); err != nil {
		return roteiro.Projeto{}, roteiro.CatalogoFontes{}, err
	}
	if catPath == "" {
		return proj, roteiro.CatalogoPadrao(), nil
	}
	cat, err := roteiro.CarregarCatalogo(catPath)
	return proj, cat, err
}

func lerJSONArquivo(caminho string, alvo any) error {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return err
	}
	return json.Unmarshal(dados, alvo)
}

func responderDiorama(d Deps, emJSON bool, payload any, code int) int {
	if emJSON {
		b, _ := json.MarshalIndent(payload, "", "  ")
		fmt.Fprintln(d.Stdout, string(b))
		return code
	}
	fmt.Fprintf(d.Stdout, "%v\n", payload)
	return code
}
