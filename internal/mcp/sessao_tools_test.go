package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

func chamarToolJSONRPC(t *testing.T, r *Router, nome string, args map[string]any) (string, bool) {
	t.Helper()
	params, err := json.Marshal(map[string]any{"name": nome, "arguments": args})
	if err != nil {
		t.Fatalf("falha ao serializar params: %v", err)
	}
	req, err := json.Marshal(JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params:  params,
	})
	if err != nil {
		t.Fatalf("falha ao serializar req: %v", err)
	}
	var out bytes.Buffer
	if err := r.Serve(bytes.NewReader(append(req, '\n')), &out); err != nil {
		t.Fatalf("Serve falhou: %v", err)
	}
	return extrairTextoResposta(t, out.Bytes())
}

func extrairTextoResposta(t *testing.T, raw []byte) (string, bool) {
	t.Helper()
	var resp struct {
		Result struct {
			IsError bool `json:"isError"`
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("resposta JSON-RPC invalida: %v (%s)", err, string(raw))
	}
	if resp.Error != nil {
		t.Fatalf("erro JSON-RPC inesperado: %v", resp.Error)
	}
	if len(resp.Result.Content) == 0 {
		t.Fatalf("resposta sem content: %s", string(raw))
	}
	return resp.Result.Content[0].Text, resp.Result.IsError
}

func novoRouterSessaoTeste() (*Router, *a2ui.Hub) {
	hub := a2ui.NewHub(a2ui.DefaultCatalog())
	r := NewRouter(nil, nil)
	r.SetSessaoHub(hub, "http://127.0.0.1:4242")
	return r, hub
}

func criarSessaoTeste(t *testing.T, r *Router) (string, string) {
	t.Helper()
	txt, isErr := chamarToolJSONRPC(t, r, "sessao_criar", map[string]any{"processo": "legoloop-concepcao"})
	if isErr {
		t.Fatalf("sessao_criar retornou erro: %s", txt)
	}
	var out struct {
		SessionID string `json:"session_id"`
		URL       string `json:"url"`
	}
	if err := json.Unmarshal([]byte(txt), &out); err != nil {
		t.Fatalf("JSON invalido de sessao_criar: %v (%s)", err, txt)
	}
	return out.SessionID, out.URL
}

func TestSessaoTools_CriarListaEPublicar(t *testing.T) {
	r, hub := novoRouterSessaoTeste()
	sessID, url := criarSessaoTeste(t, r)
	if sessID == "" || url != "http://127.0.0.1:4242/sessao/"+sessID {
		t.Fatalf("sessao_criar inesperado: id=%q url=%q", sessID, url)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := hub.Subscribe(ctx, sessID)
	if err != nil {
		t.Fatalf("Subscribe falhou: %v", err)
	}

	txt, isErr := chamarToolJSONRPC(t, r, "sessao_publicar", map[string]any{
		"session_id":       sessID,
		"surface_id":       "proc:concepcao",
		"protocol_version": "v0.9",
		"message_type":     "createSurface",
		"components":       []map[string]any{{"id": "m1", "kind": "LegoProgressMeter"}},
	})
	if isErr {
		t.Fatalf("sessao_publicar valido falhou: %s", txt)
	}
	select {
	case env := <-ch:
		if env.SurfaceID != "proc:concepcao" {
			t.Fatalf("surface_id inesperado: %q", env.SurfaceID)
		}
	case <-time.After(time.Second):
		t.Fatal("envelope publicado nao chegou no hub")
	}
}

func TestSessaoTools_PublicarRejeitaComponenteProibido(t *testing.T) {
	r, _ := novoRouterSessaoTeste()
	sessID, _ := criarSessaoTeste(t, r)
	txt, isErr := chamarToolJSONRPC(t, r, "sessao_publicar", map[string]any{
		"session_id":       sessID,
		"surface_id":       "proc:concepcao",
		"protocol_version": "v0.9",
		"message_type":     "createSurface",
		"components":       []map[string]any{{"id": "x1", "kind": "RawHTML"}},
	})
	if !isErr {
		t.Fatalf("esperava erro para componente proibido, obteve: %s", txt)
	}
}

func TestSessaoTools_AguardarRecebeActionParalela(t *testing.T) {
	r, hub := novoRouterSessaoTeste()
	sessID, _ := criarSessaoTeste(t, r)

	go func() {
		time.Sleep(25 * time.Millisecond)
		_ = hub.SubmitAction(a2ui.Action{
			SessionID: sessID,
			SurfaceID: "proc:concepcao",
			Origin:    "ui",
			Name:      "submit_answer",
			Payload:   map[string]any{"resposta": "confirmado"},
		})
	}()

	txt, isErr := chamarToolJSONRPC(t, r, "sessao_aguardar", map[string]any{
		"session_id": sessID,
		"timeout_ms": 1000,
	})
	if isErr {
		t.Fatalf("sessao_aguardar retornou erro: %s", txt)
	}
	var act a2ui.Action
	if err := json.Unmarshal([]byte(txt), &act); err != nil {
		t.Fatalf("JSON invalido de sessao_aguardar: %v (%s)", err, txt)
	}
	if act.Name != "submit_answer" || act.Origin != "ui" || act.Payload["resposta"] != "confirmado" {
		t.Fatalf("Action recebida diverge: %+v", act)
	}
}

func TestSessaoTools_AguardarTimeoutFallbackTextual(t *testing.T) {
	r, _ := novoRouterSessaoTeste()
	sessID, _ := criarSessaoTeste(t, r)

	txt, isErr := chamarToolJSONRPC(t, r, "sessao_aguardar", map[string]any{
		"session_id": sessID,
		"timeout_ms": 40,
	})
	if isErr {
		t.Fatalf("timeout nao deve retornar erro fatal: %s", txt)
	}
	var resp struct {
		FallbackTextual bool   `json:"fallback_textual"`
		SessionID       string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(txt), &resp); err != nil {
		t.Fatalf("JSON de fallback invalido: %v (%s)", err, txt)
	}
	if !resp.FallbackTextual || resp.SessionID != sessID {
		t.Fatalf("esperava fallback_textual=true para %s, obteve: %s", sessID, txt)
	}
}

func TestSessaoTools_ToolsListIncluiFerramentasSessao(t *testing.T) {
	r, _ := novoRouterSessaoTeste()
	req := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}` + "\n"
	var out bytes.Buffer
	if err := r.Serve(bytes.NewBufferString(req), &out); err != nil {
		t.Fatalf("Serve tools/list falhou: %v", err)
	}
	var resp struct {
		Result struct {
			Tools []struct {
				Name string `json:"name"`
			} `json:"tools"`
		} `json:"result"`
	}
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("JSON invalido: %v", err)
	}
	encontrados := map[string]bool{}
	for _, tool := range resp.Result.Tools {
		encontrados[tool.Name] = true
	}
	for _, esperado := range []string{"sessao_criar", "sessao_publicar", "sessao_aguardar"} {
		if !encontrados[esperado] {
			t.Fatalf("tool %q ausente em tools/list", esperado)
		}
	}
}
