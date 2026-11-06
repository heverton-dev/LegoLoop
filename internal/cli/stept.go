package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"legoloop/internal/ste"
)

func resolveHome(d Deps, flagHome string) string {
	if flagHome != "" {
		return flagHome
	}
	if d.Home != "" {
		return d.Home
	}
	return ste.HomeDirPadrao()
}

func resolveConfigPath(home string, customCfg string) string {
	if customCfg != "" {
		return customCfg
	}
	return filepath.Join(home, ".stept", "config.json")
}

func renderTabelaStatus(w io.Writer, statusList []ste.HarnessStatus, cfgPath string) {
	fmt.Fprint(w, ste.RenderTabelaBox(statusList, cfgPath))
}

func renderTabelaScan(w io.Writer, harnesses []ste.HarnessInfo) {
	tw := table(w)
	fmt.Fprintln(tw, "HARNESS\tINSTALADO\tATIVO\tPASTA BASE\tCAMINHO REGRAS")
	for _, h := range harnesses {
		inst := "não"
		if h.Instalado {
			inst = "sim"
		}
		atv := "não"
		if h.Ativo {
			atv = "sim"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", h.Nome, inst, atv, h.PastaBase, h.CaminhoRegras)
	}
	tw.Flush()
}

func executarSteptScan(home string, asJSON bool, d Deps) int {
	lista := ste.ScanHarnesses(home)
	if asJSON {
		return printJSON(d, lista)
	}
	renderTabelaScan(d.Stdout, lista)
	return 0
}

func aplicarPerfilCLI(cfg *ste.Config, harness string, modo ste.Modo, pausa ste.Pausa) error {
	if harness == "" {
		p := cfg.Padrao
		if modo != "" {
			p.Modo = modo
		}
		if pausa != "" {
			p.Pausa = pausa
		}
		return cfg.DefinirPadrao(p)
	}
	hID := ste.HarnessID(strings.ToLower(harness))
	p, _ := cfg.ResolverPerfil(hID)
	if modo != "" {
		p.Modo = modo
	}
	if pausa != "" {
		p.Pausa = pausa
	}
	return cfg.DefinirHarness(hID, p)
}

func salvarESincronizar(cfg *ste.Config, cfgPath, home string, asJSON bool, d Deps) int {
	if err := ste.SalvarConfig(cfgPath, cfg); err != nil {
		fmt.Fprintf(d.Stderr, "erro ao salvar config: %v\n", err)
		return 1
	}
	statusList, err := ste.Sync(cfg, home, false)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro ao sincronizar harnesses: %v\n", err)
		return 1
	}
	if asJSON {
		return printJSON(d, statusList)
	}
	renderTabelaStatus(d.Stdout, statusList, cfgPath)
	return 0
}

func despacharSubcomandoStept(sub string, fs *flag.FlagSet, args []string, d Deps, home, cfgPath string, asJSON bool) int {
	if sub == "manual" || sub == "help" || sub == "-h" || sub == "--help" {
		fmt.Fprint(d.Stdout, ste.TextoManual())
		return 0
	}
	cfg, err := ste.CarregarConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro ao carregar config: %v\n", err)
		return 1
	}
	switch sub {
	case "status", "":
		return salvarESincronizar(cfg, cfgPath, home, asJSON, d)
	case "scan":
		return executarSteptScan(home, asJSON, d)
	case "set", "on", "off":
		return processarComandoMutacao(sub, fs, args, d, cfg, cfgPath, home, asJSON)
	default:
		fmt.Fprintf(d.Stderr, "subcomando stept desconhecido: %q\n", sub)
		return 2
	}
}

func processarComandoMutacao(sub string, fs *flag.FlagSet, args []string, d Deps, cfg *ste.Config, cfgPath, home string, asJSON bool) int {
	harness := fs.String("harness", "", "identificador do harness")
	modoStr := fs.String("modo", "", "modo de governança")
	pausaStr := fs.String("pausa", "", "política de pausa")
	_ = fs.Parse(args)

	m := ste.Modo(strings.ToLower(*modoStr))
	p := ste.Pausa(strings.ToLower(*pausaStr))
	if sub == "on" && m == "" {
		m = ste.ModoHibrido
	}
	if sub == "off" {
		m = ste.ModoOff
	}
	if err := aplicarPerfilCLI(cfg, *harness, m, p); err != nil {
		fmt.Fprintf(d.Stderr, "erro na configuração: %v\n", err)
		return 1
	}
	return salvarESincronizar(cfg, cfgPath, home, asJSON, d)
}

func cmdStept(_ context.Context, args []string, d Deps) int {
	sub := ""
	restArgs := args
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		sub = args[0]
		restArgs = args[1:]
	}
	fs := flag.NewFlagSet("stept "+sub, flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	homeFlag := fs.String("home", "", "diretório base")
	cfgFlag := fs.String("config", "", "caminho do arquivo config.json")
	asJSON := fs.Bool("json", false, "saída JSON")

	home := resolveHome(d, *homeFlag)
	cfgPath := resolveConfigPath(home, *cfgFlag)
	return despacharSubcomandoStept(sub, fs, restArgs, d, home, cfgPath, *asJSON)
}
