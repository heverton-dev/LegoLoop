package cli

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/mcp"
)

const usoMcpRegistry = "Uso: legoloop mcp-registry (list|refresh)"

// cmdMcpRegistry gerencia o registro de MCPs externos usados pelo router lazy.
func cmdMcpRegistry(_ context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		fmt.Fprintln(d.Stderr, usoMcpRegistry)
		return 2
	}
	home := resolveHome(d, "")
	switch args[0] {
	case "list":
		return registryList(home, d)
	case "refresh":
		return registryRefresh(home, d)
	}
	fmt.Fprintln(d.Stderr, usoMcpRegistry)
	return 2
}

// registryList mostra os servidores sem nunca imprimir valores de env (podem ser segredos).
func registryList(home string, d Deps) int {
	specs, err := mcp.LoadRegistry(mcp.RegistryPath(home))
	if err != nil {
		fmt.Fprintf(d.Stderr, "registro invalido: %v\n", err)
		return 1
	}
	cat, _ := mcp.LoadCatalog(mcp.CatalogPath(home))
	for _, name := range sortedSpecNames(specs) {
		s := specs[name]
		fmt.Fprintf(d.Stdout, "%s  %s %s  (%d args, %d vars de ambiente, %d ferramentas indexadas)\n",
			name, s.Command, strings.Join(s.Args[:min(len(s.Args), 2)], " "), len(s.Args), len(s.Env), len(cat[name]))
	}
	return 0
}

func sortedSpecNames(specs map[string]mcp.ServerSpec) []string {
	names := make([]string, 0, len(specs))
	for n := range specs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// registryRefresh sobe cada servidor uma vez, indexa as ferramentas e o encerra.
func registryRefresh(home string, d Deps) int {
	specs, err := mcp.LoadRegistry(mcp.RegistryPath(home))
	if err != nil {
		fmt.Fprintf(d.Stderr, "registro invalido: %v\n", err)
		return 1
	}
	proxy := mcp.NewProxy(specs, time.Minute)
	defer proxy.Close()
	fresh, errs := proxy.Refresh()
	old, _ := mcp.LoadCatalog(mcp.CatalogPath(home))
	for name, tools := range fresh {
		old[name] = tools
	}
	if err := mcp.SaveCatalog(mcp.CatalogPath(home), old); err != nil {
		fmt.Fprintf(d.Stderr, "falha ao gravar catalogo: %v\n", err)
		return 1
	}
	return printRefresh(d, specs, fresh, errs)
}

func printRefresh(d Deps, specs map[string]mcp.ServerSpec, fresh map[string][]mcp.CatalogEntry, errs map[string]error) int {
	for _, name := range sortedSpecNames(specs) {
		if err, bad := errs[name]; bad {
			fmt.Fprintf(d.Stdout, "%s: falha: %v\n", name, err)
		} else if tools, ok := fresh[name]; ok {
			fmt.Fprintf(d.Stdout, "%s: ok, %d ferramentas indexadas\n", name, len(tools))
		}
	}
	if len(errs) > 0 {
		return 1
	}
	return 0
}
