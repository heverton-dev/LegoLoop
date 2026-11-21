package mcp

import (
	"fmt"
	"strings"
	"testing"
)

// vaultDeBusca monta um cofre com o grafo semântico e servidores de ruído, sem subir processo.
func vaultDeBusca() *localVault {
	specs := map[string]ServerSpec{"codebase-memory-mcp": {}}
	catalog := map[string][]CatalogEntry{
		"codebase-memory-mcp": {
			{Name: "trace_path", Description: "Trace callers/callees, data flow, or cross-service paths."},
			{Name: "search_graph", Description: "Search the code graph by name or label."},
			{Name: "index_repository", Description: "Index a repository into the graph."},
		},
	}
	for i := 0; i < 8; i++ {
		srv := fmt.Sprintf("ruido-%d", i)
		specs[srv] = ServerSpec{}
		catalog[srv] = []CatalogEntry{{Name: "ping", Description: "Replies pong for health checks."}}
	}
	return &localVault{specs: specs, catalog: catalog}
}

func nomesDosMatches(ms []ToolMatch) string {
	var out []string
	for _, m := range ms {
		out = append(out, m.Server+":"+m.Name)
	}
	return strings.Join(out, ",")
}

func TestBuscaComVariasPalavrasAchaAFerramenta(t *testing.T) {
	got := vaultDeBusca().Search("trace path graph")
	if len(got) == 0 || got[0].Name != "trace_path" {
		t.Fatalf("Search(\"trace path graph\") = %s, quer trace_path primeiro", nomesDosMatches(got))
	}
}

func TestBuscaEmPortuguesAchaOServidorPeloNome(t *testing.T) {
	got := nomesDosMatches(vaultDeBusca().Search("codebase memory grafo semântico"))
	for _, quer := range []string{"codebase-memory-mcp:trace_path", "codebase-memory-mcp:search_graph"} {
		if !strings.Contains(got, quer) {
			t.Fatalf("faltou %s em %s", quer, got)
		}
	}
	if strings.Contains(got, "ruido-") {
		t.Fatalf("servidor de ruído entrou na busca: %s", got)
	}
}

func TestBuscaSemTermoEmComumNaoDevolveNada(t *testing.T) {
	if got := vaultDeBusca().Search("banana laranja"); len(got) != 0 {
		t.Fatalf("termos sem relação devolveram %s", nomesDosMatches(got))
	}
}

func TestBuscaDeUmTermoContinuaFuncionando(t *testing.T) {
	got := vaultDeBusca().Search("trace_path")
	if len(got) != 1 || got[0].Name != "trace_path" {
		t.Fatalf("Search(\"trace_path\") = %s", nomesDosMatches(got))
	}
}

func TestBuscaTemOrdemDeterministica(t *testing.T) {
	v := vaultDeBusca()
	primeiro := nomesDosMatches(v.Search("health pong ping"))
	if primeiro == "" {
		t.Fatal("busca por health pong ping não achou os servidores de ruído")
	}
	for i := 0; i < 30; i++ {
		if got := nomesDosMatches(v.Search("health pong ping")); got != primeiro {
			t.Fatalf("ordem mudou entre chamadas:\n%s\n%s", primeiro, got)
		}
	}
}
