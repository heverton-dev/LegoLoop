package mcp

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

// gwServe envia requisições JSON-RPC ao Serve e devolve as respostas já decodificadas.
func gwServe(t *testing.T, r *Router, reqs ...string) []map[string]any {
	t.Helper()
	var out bytes.Buffer
	if err := r.Serve(strings.NewReader(strings.Join(reqs, "\n")+"\n"), &out); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	var resps []map[string]any
	for _, l := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		var m map[string]any
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			t.Fatalf("resposta inválida %q: %v", l, err)
		}
		resps = append(resps, m)
	}
	return resps
}

func gatewayRouter(t *testing.T) *Router {
	t.Helper()
	home := t.TempDir()
	r := NewConfiguredRouter(home, filepath.Join(home, ".aidd", "canonical"), t.TempDir())
	r.EnableGateway()
	t.Cleanup(r.Close)
	return r
}

func chamadaTool(id int, server, tool string, args map[string]any) string {
	p, _ := json.Marshal(map[string]any{"name": "call_mcp_tool", "arguments": map[string]any{"server": server, "tool": tool, "arguments": args}})
	return `{"jsonrpc":"2.0","id":` + string(rune('0'+id)) + `,"method":"tools/call","params":` + string(p) + `}`
}

func textoDe(t *testing.T, resp map[string]any) (string, bool) {
	t.Helper()
	res, _ := resp["result"].(map[string]any)
	erro, _ := res["isError"].(bool)
	cont, _ := res["content"].([]any)
	if len(cont) == 0 {
		t.Fatalf("resposta sem content: %v", resp)
	}
	return cont[0].(map[string]any)["text"].(string), erro
}

func nomesDeTools(t *testing.T, resp map[string]any) []string {
	t.Helper()
	res, _ := resp["result"].(map[string]any)
	var nomes []string
	for _, x := range res["tools"].([]any) {
		nomes = append(nomes, x.(map[string]any)["name"].(string))
	}
	sort.Strings(nomes)
	return nomes
}

// TestGatewayExpoeUmaFerramenta prova o piso de contexto: em modo gateway, tools/list
// devolve só call_mcp_tool; o modo padrão continua expondo swarm_*.
func TestGatewayExpoeUmaFerramenta(t *testing.T) {
	r := gatewayRouter(t)
	got := nomesDeTools(t, gwServe(t, r, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)[0])
	want := []string{"call_mcp_tool"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("tools/list = %v, quer %v", got, want)
	}
	home := t.TempDir()
	padrao := NewConfiguredRouter(home, filepath.Join(home, ".aidd", "canonical"), t.TempDir())
	defer padrao.Close()
	if n := nomesDeTools(t, gwServe(t, padrao, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)[0]); len(n) <= 4 {
		t.Fatalf("modo padrão deve manter swarm_*, veio %v", n)
	}
}

// TestGatewayWsBash prova que o shell sai pelo servidor embutido ws, sem ferramenta nativa.
func TestGatewayWsBash(t *testing.T) {
	r := gatewayRouter(t)
	txt, erro := textoDe(t, gwServe(t, r, chamadaTool(1, "ws", "bash", map[string]any{"command": "echo ola-gateway"}))[0])
	if erro || !strings.Contains(txt, "ola-gateway") {
		t.Fatalf("bash via ws = %q (erro=%v)", txt, erro)
	}
	txt, erro = textoDe(t, gwServe(t, r, chamadaTool(2, "ws", "bash", map[string]any{"command": "exit 3"}))[0])
	if !erro || !strings.Contains(txt, "3") {
		t.Fatalf("falha do comando deve virar erro com o código, veio %q (erro=%v)", txt, erro)
	}
}

// TestGatewayWsArquivos prova leitura, escrita e edição de arquivos via ws, e a recusa de edição ambígua.
func TestGatewayWsArquivos(t *testing.T) {
	r := gatewayRouter(t)
	arq := filepath.Join(t.TempDir(), "a", "nota.txt")
	if _, erro := textoDe(t, gwServe(t, r, chamadaTool(1, "ws", "write_file", map[string]any{"path": arq, "content": "um dois dois"}))[0]); erro {
		t.Fatal("write_file falhou")
	}
	if txt, _ := textoDe(t, gwServe(t, r, chamadaTool(2, "ws", "read_file", map[string]any{"path": arq}))[0]); txt != "um dois dois" {
		t.Fatalf("read_file = %q", txt)
	}
	if _, erro := textoDe(t, gwServe(t, r, chamadaTool(3, "ws", "edit_file", map[string]any{"path": arq, "old": "um", "new": "UM"}))[0]); erro {
		t.Fatal("edit_file único falhou")
	}
	if txt, _ := textoDe(t, gwServe(t, r, chamadaTool(4, "ws", "read_file", map[string]any{"path": arq}))[0]); txt != "UM dois dois" {
		t.Fatalf("após edição = %q", txt)
	}
	if _, erro := textoDe(t, gwServe(t, r, chamadaTool(5, "ws", "edit_file", map[string]any{"path": arq, "old": "dois", "new": "x"}))[0]); !erro {
		t.Fatal("edit_file ambíguo (2 ocorrências) deve falhar")
	}
	if _, erro := textoDe(t, gwServe(t, r, chamadaTool(6, "ws", "edit_file", map[string]any{"path": arq, "old": "ausente", "new": "x"}))[0]); !erro {
		t.Fatal("edit_file sem ocorrência deve falhar")
	}
}

// TestGatewaySwarmViaCallMcpTool prova que o swarm continua acessível sem estar em tools/list.
func TestGatewaySwarmViaCallMcpTool(t *testing.T) {
	r := gatewayRouter(t)
	dir := t.TempDir()
	txt, erro := textoDe(t, gwServe(t, r, chamadaTool(1, "legoloop", "swarm_status", map[string]any{"dir": dir}))[0])
	if erro || !strings.Contains(txt, "workers") {
		t.Fatalf("swarm_status via call_mcp_tool = %q (erro=%v)", txt, erro)
	}
}

// TestGatewayDescricaoAnunciaServidoresEmbutidos prova que o modelo descobre ws e legoloop sem busca.
func TestGatewayDescricaoAnunciaServidoresEmbutidos(t *testing.T) {
	r := gatewayRouter(t)
	res, _ := gwServe(t, r, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)[0]["result"].(map[string]any)
	for _, x := range res["tools"].([]any) {
		m := x.(map[string]any)
		if m["name"] == "call_mcp_tool" {
			d := m["description"].(string)
			for _, trecho := range []string{"ws", "bash", "swarm_", "skill_search", "call_skill", "mcp_search"} {
				if !strings.Contains(d, trecho) {
					t.Fatalf("descrição de call_mcp_tool não anuncia %q: %q", trecho, d)
				}
			}
			if len(m) > 0 && len(d) > 200 {
				t.Fatalf("descrição com %d bytes; o gateway precisa ser compacto (<= 200)", len(d))
			}
			return
		}
	}
	t.Fatal("call_mcp_tool ausente")
}

// TestGatewayEscondeToolsDeProcessoMasChamaPorLegoloop prova que tools de Processo não furam o piso
// de contexto e continuam chamáveis por call_mcp_tool no servidor legoloop.
func TestGatewayEscondeToolsDeProcessoMasChamaPorLegoloop(t *testing.T) {
	base := t.TempDir()
	prepararProcessoComTools(t, base, "alpha", manifestoEco("eco_alpha"))
	r := gatewayRouter(t)
	if err := r.CarregarToolsDeProcessos(base); err != nil {
		t.Fatalf("CarregarToolsDeProcessos: %v", err)
	}
	got := nomesDeTools(t, gwServe(t, r, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)[0])
	if len(got) != 1 {
		t.Fatalf("tools/list com Processo carregado = %v, quer só call_mcp_tool", got)
	}
	txt, erro := textoDe(t, gwServe(t, r, chamadaTool(2, "legoloop", "eco_alpha", map[string]any{"texto": "abc"}))[0])
	if erro || !strings.Contains(txt, "abc") {
		t.Fatalf("eco_alpha via legoloop = %q (erro=%v)", txt, erro)
	}
}

// TestGatewayBuscaESkillsViaLegoloop prova que skill_search, call_skill e mcp_search saem pelo servidor embutido.
func TestGatewayBuscaESkillsViaLegoloop(t *testing.T) {
	home := t.TempDir()
	vault := filepath.Join(home, ".aidd", "canonical")
	dir := filepath.Join(vault, "skills", "skill-gw")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	conteudo := "---" + "\n" + "name: skill-gw" + "\n" + "description: Prova de gateway para buscar skills." + "\n" + "---" + "\n" + "corpo da skill gw" + "\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	r := NewConfiguredRouter(home, vault, t.TempDir())
	r.EnableGateway()
	defer r.Close()
	busca, erro := textoDe(t, gwServe(t, r, chamadaTool(1, "legoloop", "skill_search", map[string]any{"query": "skill-gw gateway"}))[0])
	if erro || !strings.Contains(busca, "skill-gw") {
		t.Fatalf("skill_search via legoloop = %q (erro=%v)", busca, erro)
	}
	corpo, erro := textoDe(t, gwServe(t, r, chamadaTool(2, "legoloop", "call_skill", map[string]any{"skill_id": "skill-gw"}))[0])
	if erro || !strings.Contains(corpo, "corpo da skill gw") {
		t.Fatalf("call_skill via legoloop = %q (erro=%v)", corpo, erro)
	}
	if _, erro := textoDe(t, gwServe(t, r, chamadaTool(3, "legoloop", "mcp_search", map[string]any{"query": "qualquer"}))[0]); erro {
		t.Fatal("mcp_search via legoloop falhou")
	}
}

// TestGatewaySessaoCriarViaLegoloop prova que sessao_* (A2UI) é despachado via call_mcp_tool no modo gateway.
func TestGatewaySessaoCriarViaLegoloop(t *testing.T) {
	r := gatewayRouter(t)
	hub := a2ui.NewHub(a2ui.DefaultCatalog())
	r.SetSessaoHub(hub, "http://127.0.0.1:7878")
	r.SetSessaoAbridor(func(string) error { return nil })

	txt, erro := textoDe(t, gwServe(t, r, chamadaTool(1, "legoloop", "sessao_criar", map[string]any{"processo": "legoloop-diorama-render", "abrir_navegador": false}))[0])
	if erro {
		t.Fatalf("sessao_criar via legoloop retornou erro: %s", txt)
	}
	if !strings.Contains(txt, "id") || !strings.Contains(txt, "http://127.0.0.1:7878/sessao/") {
		t.Fatalf("resposta inesperada de sessao_criar: %s", txt)
	}
}
