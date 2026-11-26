package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/browser"
	"github.com/heverton-dev/LegoLoop/internal/diorama/ativos"
)

const usoDioramaAtivos = `Uso: legoloop diorama-ativos <subcomando> [opções]

Subcomandos:
  prevoo              --projeto DIR [--json]
  executar            --projeto DIR [--aprovar gate1,gate2] [--json]
  validar-storyboard  --projeto DIR [--json]
  abrir-flow          [--porta 9222] [--url URL] [--json]
  criar-projeto-flow  --projeto DIR [--cdp URL] [--json]
`

func cmdDioramaAtivos(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		fmt.Fprint(d.Stdout, usoDioramaAtivos)
		return 0
	}
	switch args[0] {
	case "prevoo":
		return cmdDioramaAtivosPreVoo(ctx, args[1:], d)
	case "executar":
		return cmdDioramaAtivosExecutar(ctx, args[1:], d)
	case "validar-storyboard":
		return cmdDioramaAtivosValidarStoryboard(args[1:], d)
	case "abrir-flow":
		return cmdDioramaAtivosAbrirFlow(ctx, args[1:], d)
	case "criar-projeto-flow":
		return cmdDioramaAtivosCriarProjetoFlow(ctx, args[1:], d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n%s", args[0], usoDioramaAtivos)
		return 2
	}
}

func cmdDioramaAtivosPreVoo(ctx context.Context, args []string, d Deps) int {
	projDir, _, emJSON, ok := parseDioramaAtivosFlags("prevoo", args, d.Stderr)
	if !ok {
		return 2
	}
	rel, err := ativos.ExecutarPreVoo(ctx, projDir, ativos.DepsPreVoo{})
	if err != nil || rel.Bloqueado {
		return responderDiorama(d, emJSON, map[string]any{"bloqueado": true, "relatorio": rel, "erro": errStr(err)}, 1)
	}
	return responderDiorama(d, emJSON, rel, 0)
}

func cmdDioramaAtivosExecutar(ctx context.Context, args []string, d Deps) int {
	projDir, aprovar, emJSON, ok := parseDioramaAtivosFlags("executar", args, d.Stderr)
	if !ok {
		return 2
	}
	contratos, err := ativos.CarregarContratosEntrada(projDir)
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	var deps ativos.DepsProcesso
	if contratos.Projeto.EstiloVisual == "flow" || contratos.Projeto.EstiloVisual == "hibrido" {
		deps.Recortador = ativos.NovoRecortadorRembgDespill()
		deps.Cenarios = ativos.DepsCenarios{
			Flow: ativos.NovoClienteFlowCDP("", projDir),
		}
	}
	proc := ativos.NovoProcesso(deps)
	res, err := proc.Executar(ctx, ativos.OpcoesExecucao{
		ProjetoDir:     projDir,
		GatesAprovados: aprovar,
	})
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"ok": false, "erro": err.Error(), "resultado": res}, 1)
	}
	return responderDiorama(d, emJSON, map[string]any{"ok": true, "resultado": res}, 0)
}

func cmdDioramaAtivosValidarStoryboard(args []string, d Deps) int {
	projDir, _, emJSON, ok := parseDioramaAtivosFlags("validar-storyboard", args, d.Stderr)
	if !ok {
		return 2
	}
	contratos, err := ativos.CarregarContratosEntrada(projDir)
	if err != nil {
		return responderDiorama(d, emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	var cenarios []ativos.Cenario
	if err := lerJSONCLI(filepath.Join(projDir, "flow", "cenarios.json"), &cenarios); err != nil {
		return responderDiorama(d, emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	var sb ativos.Storyboard
	if err := lerJSONCLI(filepath.Join(projDir, "cenas.json"), &sb); err != nil {
		return responderDiorama(d, emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	if err := ativos.ValidarStoryboard(sb, contratos.Script, cenarios); err != nil {
		return responderDiorama(d, emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	return responderDiorama(d, emJSON, map[string]any{"ok": true, "total_cenas": len(sb.Cenas)}, 0)
}

func cmdDioramaAtivosAbrirFlow(ctx context.Context, args []string, d Deps) int {
	fs := flag.NewFlagSet("diorama-ativos abrir-flow", flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	porta := fs.Int("porta", browser.PortaPadrao, "porta CDP do navegador persistente")
	urlInicial := fs.String("url", "https://flow.google.com", "URL inicial do navegador")
	emJSON := fs.Bool("json", false, "saida em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if err := browser.AbrirPersistente(ctx, browser.OpcoesAbertura{
		Porta:      *porta,
		URLInicial: *urlInicial,
	}); err != nil {
		return responderDiorama(d, *emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	return responderDiorama(d, *emJSON, map[string]any{"ok": true, "porta": *porta, "perfil": browser.PerfilPadrao("")}, 0)
}

func cmdDioramaAtivosCriarProjetoFlow(ctx context.Context, args []string, d Deps) int {
	fs := flag.NewFlagSet("diorama-ativos criar-projeto-flow", flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	projDir := fs.String("projeto", ".", "diretorio do projeto")
	cdpURL := fs.String("cdp", "", "URL CDP opcional (padrao auto 9222/9333)")
	emJSON := fs.Bool("json", false, "saida em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	flowUUID, err := ativos.CriarProjetoFlow(ctx, *projDir, *cdpURL, nil)
	if err != nil {
		return responderDiorama(d, *emJSON, map[string]any{"ok": false, "erro": err.Error()}, 1)
	}
	return responderDiorama(d, *emJSON, map[string]any{"ok": true, "flow_projeto": flowUUID}, 0)
}

func parseDioramaAtivosFlags(nome string, args []string, stderr io.Writer) (string, []string, bool, bool) {
	fs := flag.NewFlagSet("diorama-ativos "+nome, flag.ContinueOnError)
	fs.SetOutput(stderr)
	projDir := fs.String("projeto", ".", "diretorio do projeto")
	aprovarCSV := fs.String("aprovar", "", "gates humanos aprovados separados por virgula")
	emJSON := fs.Bool("json", false, "saida em JSON")
	if err := fs.Parse(args); err != nil {
		return "", nil, false, false
	}
	var gates []string
	for _, g := range strings.Split(*aprovarCSV, ",") {
		if t := strings.TrimSpace(g); t != "" {
			gates = append(gates, t)
		}
	}
	return *projDir, gates, *emJSON, true
}

func lerJSONCLI(caminho string, dst any) error {
	b, err := os.ReadFile(caminho)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

func errStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
