package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type payloadInstalarTeste struct {
	RelatorioJSON string `json:"relatorio_json"`
	Veredito      string `json:"veredito"`
	DryRun        bool   `json:"dry_run"`
	Aplicado      bool   `json:"aplicado"`
}

func TestInstalarViaMCP(t *testing.T) {
	t.Run("call_mcp_tool dry_run nao grava no disco", testarInstalarMCPDryRun)
	t.Run("chamada direta e listagem no modo padrao", testarInstalarMCPDireto)
	t.Run("apply grava RELATORIO.json no caminho retornado", testarInstalarMCPApply)
	t.Run("perfil invalido devolve erro", testarInstalarMCPPerfilInvalido)
	t.Run("mcp_search encontra legoloop:instalar", testarInstalarMCPSearch)
}

func testarInstalarMCPDryRun(t *testing.T) {
	home := t.TempDir()
	r := gatewayRouter(t)
	req := chamadaTool(1, "legoloop", "instalar", map[string]any{
		"dry_run": true,
		"home":    home,
		"perfil":  "padrao",
	})
	txt, isErr := textoDe(t, gwServe(t, r, req)[0])
	if isErr {
		t.Fatalf("esperava sucesso no dry_run, obteve erro: %s", txt)
	}
	validarDryRunInstalar(t, home, decodificarPayloadInstalar(t, txt))
}

func testarInstalarMCPDireto(t *testing.T) {
	home := t.TempDir()
	r := NewConfiguredRouter(home, filepath.Join(home, ".aidd", "canonical"), t.TempDir())
	defer r.Close()
	nomes := nomesDeTools(t, gwServe(t, r, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)[0])
	if !contemString(nomes, "instalar") {
		t.Fatalf("tools/list sem instalar: %v", nomes)
	}
	req := chamadaDiretaInstalar(t, 2, map[string]any{
		"dry_run": true,
		"home":    home,
		"perfil":  "padrao",
	})
	txt, isErr := textoDe(t, gwServe(t, r, req)[0])
	if isErr {
		t.Fatalf("esperava sucesso na chamada direta, obteve erro: %s", txt)
	}
	validarDryRunInstalar(t, home, decodificarPayloadInstalar(t, txt))
}

func testarInstalarMCPApply(t *testing.T) {
	home := t.TempDir()
	r := gatewayRouter(t)
	req := chamadaTool(3, "legoloop", "instalar", map[string]any{
		"dry_run": false,
		"home":    home,
		"perfil":  "padrao",
	})
	txt, isErr := textoDe(t, gwServe(t, r, req)[0])
	if isErr {
		t.Fatalf("esperava sucesso no apply, obteve erro: %s", txt)
	}
	payload := decodificarPayloadInstalar(t, txt)
	if payload.Veredito != "OK" || payload.DryRun || !payload.Aplicado {
		t.Fatalf("payload inesperado no apply: %+v", payload)
	}
	if _, err := os.Stat(payload.RelatorioJSON); err != nil {
		t.Fatalf("RELATORIO.json ausente em %q: %v", payload.RelatorioJSON, err)
	}
}

func testarInstalarMCPPerfilInvalido(t *testing.T) {
	r := gatewayRouter(t)
	req := chamadaTool(4, "legoloop", "instalar", map[string]any{
		"dry_run": true,
		"home":    t.TempDir(),
		"perfil":  "invalido",
	})
	txt, isErr := textoDe(t, gwServe(t, r, req)[0])
	if !isErr || !strings.Contains(txt, "perfil") {
		t.Fatalf("esperava erro de perfil invalido, obteve isErr=%v txt=%q", isErr, txt)
	}
}

func testarInstalarMCPSearch(t *testing.T) {
	r := gatewayRouter(t)
	req := chamadaTool(5, "legoloop", "mcp_search", map[string]any{"query": "instalar"})
	txt, isErr := textoDe(t, gwServe(t, r, req)[0])
	if isErr || !strings.Contains(txt, "instalar") {
		t.Fatalf("mcp_search nao encontrou instalar: isErr=%v txt=%q", isErr, txt)
	}
}

func decodificarPayloadInstalar(t *testing.T, txt string) payloadInstalarTeste {
	t.Helper()
	var p payloadInstalarTeste
	if err := json.Unmarshal([]byte(txt), &p); err != nil {
		t.Fatalf("resposta nao e JSON valido (%v): %s", err, txt)
	}
	return p
}

func validarDryRunInstalar(t *testing.T, home string, p payloadInstalarTeste) {
	t.Helper()
	if p.Veredito != "OK" {
		t.Fatalf("veredito=%q, esperava OK", p.Veredito)
	}
	if !p.DryRun || p.Aplicado {
		t.Fatalf("esperava dry_run=true e aplicado=false, obteve %+v", p)
	}
	if !strings.HasSuffix(p.RelatorioJSON, "RELATORIO.json") {
		t.Fatalf("relatorio_json=%q nao termina em RELATORIO.json", p.RelatorioJSON)
	}
	if !strings.HasPrefix(p.RelatorioJSON, home) {
		t.Fatalf("relatorio_json=%q fora do home=%q", p.RelatorioJSON, home)
	}
	entradas, err := os.ReadDir(home)
	if err != nil {
		t.Fatalf("ReadDir(home): %v", err)
	}
	if len(entradas) != 0 {
		t.Fatalf("dry_run gravou arquivos no home: %v", entradas)
	}
}

func chamadaDiretaInstalar(t *testing.T, id int, args map[string]any) string {
	t.Helper()
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "instalar",
			"arguments": args,
		},
	}
	dados, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal chamadaDiretaInstalar: %v", err)
	}
	return string(dados)
}

func contemString(lista []string, alvo string) bool {
	for _, item := range lista {
		if item == alvo {
			return true
		}
	}
	return false
}
