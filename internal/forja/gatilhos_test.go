package forja_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/cli"
	"github.com/heverton-dev/LegoLoop/internal/forja"
	"github.com/heverton-dev/LegoLoop/internal/instalar"
	"github.com/heverton-dev/LegoLoop/internal/mcp"
)

type executorForjaMCP struct{}

func (executorForjaMCP) Run(srv, tool string, args map[string]any) (string, error) {
	if srv != "legoloop" || tool != "forjar" {
		return "", fmt.Errorf("ferramenta nao suportada: %s:%s", srv, tool)
	}
	home, _ := args["home"].(string)
	perfil, _ := args["perfil"].(string)
	aplicar, _ := args["aplicar"].(bool)
	rel, err := forja.Executar(context.Background(), forja.Opcao{
		Home:    home,
		Perfil:  perfil,
		Aplicar: aplicar,
	})
	if err != nil {
		return "", err
	}
	bruto, err := forja.RenderizarRelatorioJSON(rel, nil)
	if err != nil {
		return "", err
	}
	return string(bruto), nil
}

func TestQuatroGatilhosMesmoRelatorio(t *testing.T) {
	home := t.TempDir()
	cofre := criarCofreSkillForjar(t)
	ctx := context.Background()

	execucoes := []struct {
		nome string
		fn   func() []byte
	}{
		{"instalar", func() []byte { return executarGatilhoInstalar(t, ctx, home) }},
		{"cli", func() []byte { return executarGatilhoCLI(t, ctx, home) }},
		{"mcp", func() []byte { return executarGatilhoMCP(t, home) }},
		{"skill", func() []byte { return executarGatilhoSkill(t, ctx, home, cofre) }},
	}

	var referencia []byte
	for i, ex := range execucoes {
		resetarHomeGatilhos(t, home)
		bruto := ex.fn()
		norm := normalizarRelatorioJSON(t, ex.nome, bruto)
		if i == 0 {
			validarRelatorioBase(t, bruto)
			referencia = norm
			continue
		}
		if !bytes.Equal(referencia, norm) {
			t.Fatalf("diff entre instalar e %s:\n--- instalar ---\n%s\n--- %s ---\n%s",
				ex.nome, string(referencia), ex.nome, string(norm))
		}
	}
}

func TestQuatroGatilhosDetectaDivergencia(t *testing.T) {
	home := t.TempDir()
	ctx := context.Background()
	resetarHomeGatilhos(t, home)
	padrao := normalizarRelatorioJSON(t, "padrao", executarForjaDireto(t, ctx, home, "padrao"))
	resetarHomeGatilhos(t, home)
	minimo := normalizarRelatorioJSON(t, "minimo", executarForjaDireto(t, ctx, home, "minimo"))
	if bytes.Equal(padrao, minimo) {
		t.Fatal("esperava diff nao vazio entre perfis diferentes")
	}
}

func executarForjaDireto(t *testing.T, ctx context.Context, home, perfil string) []byte {
	t.Helper()
	rel, err := forja.Executar(ctx, forja.Opcao{Home: home, Perfil: perfil, Aplicar: true})
	if err != nil {
		t.Fatalf("forja.Executar(%s) falhou: %v", perfil, err)
	}
	bruto, err := forja.RenderizarRelatorioJSON(rel, nil)
	if err != nil {
		t.Fatalf("RenderizarRelatorioJSON falhou: %v", err)
	}
	return bruto
}

func resetarHomeGatilhos(t *testing.T, home string) {
	t.Helper()
	entradas, err := os.ReadDir(home)
	if err != nil {
		t.Fatalf("falha ao ler home: %v", err)
	}
	for _, e := range entradas {
		if err := os.RemoveAll(filepath.Join(home, e.Name())); err != nil {
			t.Fatalf("falha ao limpar %s: %v", e.Name(), err)
		}
	}
	for _, h := range []string{"claude", "gemini", "codex"} {
		dir := filepath.Join(home, "."+h)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		arq := filepath.Join(dir, "usuario.txt")
		if err := os.WriteFile(arq, []byte("conteudo-"+h+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func executarGatilhoInstalar(t *testing.T, ctx context.Context, home string) []byte {
	t.Helper()
	res, err := instalar.Executar(ctx, instalar.Opcao{
		Home:     home,
		Perfil:   "padrao",
		Aplicar:  true,
		Ambiente: []string{},
	})
	if err != nil {
		t.Fatalf("gatilho instalar falhou: %v", err)
	}
	bruto, err := os.ReadFile(res.CaminhoJSON)
	if err != nil {
		t.Fatalf("falha ao ler RELATORIO.json de instalar: %v", err)
	}
	return bruto
}

func executarGatilhoCLI(t *testing.T, ctx context.Context, home string) []byte {
	t.Helper()
	var out, errBuf bytes.Buffer
	d := cli.Deps{Stdout: &out, Stderr: &errBuf, Home: home}
	args := []string{"forjar", "--apply", "--perfil", "padrao", "--home", home, "--json"}
	code := cli.Run(ctx, args, d)
	if code != 0 {
		t.Fatalf("gatilho CLI falhou (code=%d): %s", code, errBuf.String())
	}
	return out.Bytes()
}

func executarGatilhoMCP(t *testing.T, home string) []byte {
	t.Helper()
	r := mcp.NewRouter(nil, executorForjaMCP{})
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]any{
			"name": "call_mcp_tool",
			"arguments": map[string]any{
				"server": "legoloop",
				"tool":   "forjar",
				"arguments": map[string]any{
					"home":    home,
					"perfil":  "padrao",
					"aplicar": true,
				},
			},
		},
	}
	return []byte(chamarRouterJSONRPC(t, r, req))
}

func executarGatilhoSkill(t *testing.T, ctx context.Context, home, cofre string) []byte {
	t.Helper()
	eng := mcp.NewSkillEngine()
	if err := eng.IndexDirectory(cofre); err != nil {
		t.Fatalf("falha ao indexar cofre: %v", err)
	}
	r := mcp.NewRouter(nil, nil)
	r.SetSkills(eng, cofre)
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "call_skill",
			"arguments": map[string]any{"skill_id": "legoloop-forjar"},
		},
	}
	corpo := chamarRouterJSONRPC(t, r, req)
	args := extrairArgsSkillForjar(t, corpo, home)
	var out, errBuf bytes.Buffer
	d := cli.Deps{Stdout: &out, Stderr: &errBuf, Home: home}
	if code := cli.Run(ctx, args, d); code != 0 {
		t.Fatalf("gatilho skill falhou (code=%d): %s", code, errBuf.String())
	}
	return out.Bytes()
}

func criarCofreSkillForjar(t *testing.T) string {
	t.Helper()
	cofre := t.TempDir()
	dir := filepath.Join(cofre, "legoloop-forjar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	conteudo := "---\nname: legoloop-forjar\ndescription: forjar ambiente com a Forja\n---\n" +
		"# Skill legoloop-forjar\n\nExecute:\nlegoloop forjar --apply --perfil padrao --json\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	return cofre
}

func extrairArgsSkillForjar(t *testing.T, corpo, home string) []string {
	t.Helper()
	for _, linha := range strings.Split(corpo, "\n") {
		limpa := strings.TrimSpace(linha)
		if strings.HasPrefix(limpa, "legoloop forjar") {
			partes := strings.Fields(limpa)
			return append(partes[1:], "--home", home)
		}
	}
	t.Fatalf("comando legoloop forjar ausente na skill:\n%s", corpo)
	return nil
}

func chamarRouterJSONRPC(t *testing.T, r *mcp.Router, req map[string]any) string {
	t.Helper()
	linha, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("falha ao serializar requisicao: %v", err)
	}
	var out bytes.Buffer
	if err := r.Serve(bytes.NewReader(append(linha, '\n')), &out); err != nil {
		t.Fatalf("router.Serve falhou: %v", err)
	}
	var resp struct {
		Result struct {
			IsError bool `json:"isError"`
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &resp); err != nil {
		t.Fatalf("resposta JSON-RPC invalida: %v (%s)", err, out.String())
	}
	if resp.Error != nil || resp.Result.IsError || len(resp.Result.Content) == 0 {
		t.Fatalf("chamada JSON-RPC falhou: %s", out.String())
	}
	return resp.Result.Content[0].Text
}

func validarRelatorioBase(t *testing.T, bruto []byte) {
	t.Helper()
	var rel forja.Relatorio
	if err := json.Unmarshal(bruto, &rel); err != nil {
		t.Fatalf("RELATORIO.json invalido: %v", err)
	}
	if rel.Veredito != forja.VereditoOK {
		t.Fatalf("veredito=%q, esperava %q", rel.Veredito, forja.VereditoOK)
	}
	if len(rel.Etapas) < 3 {
		t.Fatalf("esperava ao menos 3 etapas, obteve %d", len(rel.Etapas))
	}
}

func normalizarRelatorioJSON(t *testing.T, gatilho string, bruto []byte) []byte {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(bruto, &doc); err != nil {
		t.Fatalf("gatilho %s gerou JSON invalido: %v\n%s", gatilho, err, string(bruto))
	}
	for _, campo := range []string{"data", "timestamp", "gerado_em", "criado_em"} {
		delete(doc, campo)
	}
	norm, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("falha ao normalizar JSON de %s: %v", gatilho, err)
	}
	return norm
}
