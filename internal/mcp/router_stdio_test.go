package mcp

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type rpcEnvelope struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  struct {
		IsError bool `json:"isError"`
		Tools   []struct {
			Name string `json:"name"`
		} `json:"tools"`
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"result"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestMCPRouterCofreSwarmEErro32700(t *testing.T) {
	t.Run("CofreRealSkillsESwarmNoToolsList", func(t *testing.T) {
		home := t.TempDir()
		vault := filepath.Join(home, ".aidd", "canonical")
		proj := t.TempDir()
		escreverSkill(t, filepath.Join(vault, "skills"), "aidd-tdd", "Ciclo estrito Red Green Refactor", "build")

		r := NewConfiguredRouter(home, vault, proj)

		reqs := strings.Join([]string{
			`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05"}}`,
			`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
			`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
			`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"mcp_search","arguments":{"query":"tdd"}}}`,
			`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"call_mcp_tool","arguments":{"server":"skills","tool":"aidd-tdd","arguments":{}}}}`,
		}, "\n") + "\n"

		var out bytes.Buffer
		if err := r.Serve(strings.NewReader(reqs), &out); err != nil {
			t.Fatalf("Serve falhou: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if len(lines) != 4 {
			t.Fatalf("esperava 4 respostas (notificacao nao gera resposta), obteve %d: %s", len(lines), out.String())
		}

		var listEnv rpcEnvelope
		if err := json.Unmarshal([]byte(lines[1]), &listEnv); err != nil {
			t.Fatalf("unmarshal tools/list: %v", err)
		}
		toolNames := make(map[string]bool)
		for _, tl := range listEnv.Result.Tools {
			toolNames[tl.Name] = true
		}
		for _, want := range []string{
			"mcp_search", "call_mcp_tool",
			"skill_search", "call_skill",
			"swarm_add", "swarm_spawn", "swarm_post", "swarm_read", "swarm_wait", "swarm_status", "swarm_kill", "swarm_monitor",
		} {
			if !toolNames[want] {
				t.Errorf("ferramenta esperada %q ausente em tools/list: %v", want, toolNames)
			}
		}

		var searchEnv rpcEnvelope
		if err := json.Unmarshal([]byte(lines[2]), &searchEnv); err != nil {
			t.Fatalf("unmarshal mcp_search: %v", err)
		}
		if len(searchEnv.Result.Content) == 0 || !strings.Contains(searchEnv.Result.Content[0].Text, "aidd-tdd") {
			t.Fatalf("mcp_search deveria encontrar aidd-tdd no cofre real, obteve: %+v", searchEnv.Result)
		}

		var callEnv rpcEnvelope
		if err := json.Unmarshal([]byte(lines[3]), &callEnv); err != nil {
			t.Fatalf("unmarshal call_mcp_tool: %v", err)
		}
		if callEnv.Result.IsError || len(callEnv.Result.Content) == 0 || !strings.Contains(callEnv.Result.Content[0].Text, "# corpo aidd-tdd") {
			t.Fatalf("call_mcp_tool deveria executar a skill do cofre, obteve: %+v", callEnv.Result)
		}
	})

	t.Run("JSONInvalidoRetornaErro32700ELinhaMaiorQue64KiB", func(t *testing.T) {
		home := t.TempDir()
		vault := filepath.Join(home, ".aidd", "canonical")
		r := NewConfiguredRouter(home, vault, t.TempDir())

		bigQuery := strings.Repeat("a", 80*1024)
		bigReq := `{"jsonrpc":"2.0","id":99,"method":"tools/call","params":{"name":"mcp_search","arguments":{"query":"` + bigQuery + `"}}}`
		input := "{json_malformado_sem_fechar\n" + bigReq + "\n"

		var out bytes.Buffer
		if err := r.Serve(strings.NewReader(input), &out); err != nil {
			t.Fatalf("Serve falhou com payload > 64KiB: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if len(lines) != 2 {
			t.Fatalf("esperava 2 respostas (erro -32700 + resposta >64KiB), obteve %d: %s", len(lines), out.String())
		}

		var errEnv rpcEnvelope
		if err := json.Unmarshal([]byte(lines[0]), &errEnv); err != nil {
			t.Fatalf("unmarshal erro -32700: %v", err)
		}
		if errEnv.Error == nil || errEnv.Error.Code != -32700 {
			t.Fatalf("esperava JSON-RPC error code -32700, obteve: %+v", errEnv.Error)
		}

		var okEnv rpcEnvelope
		if err := json.Unmarshal([]byte(lines[1]), &okEnv); err != nil {
			t.Fatalf("unmarshal resposta >64KiB: %v", err)
		}
		if okEnv.ID != float64(99) || okEnv.Error != nil {
			t.Fatalf("esperava resposta id=99 sem erro para linha >64KiB, obteve: %+v", okEnv)
		}
	})

	t.Run("SwarmStatusEstritamenteSomenteLeitura", func(t *testing.T) {
		emptyDir := t.TempDir()
		r := NewConfiguredRouter(t.TempDir(), t.TempDir(), emptyDir)

		statusReq := `{"jsonrpc":"2.0","id":10,"method":"tools/call","params":{"name":"swarm_status","arguments":{"dir":` + quoteJSON(emptyDir) + `}}}`
		resp := chamarRouter(t, r, statusReq)
		if resp.Result.IsError {
			t.Fatalf("swarm_status falhou em diretorio vazio: %+v", resp.Result)
		}
		if _, err := os.Stat(filepath.Join(emptyDir, ".legoloop-swarm")); !os.IsNotExist(err) {
			t.Fatalf("swarm_status NAO pode criar .legoloop-swarm no disco (err=%v)", err)
		}

		activeDir := t.TempDir()
		addReq := `{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"swarm_add","arguments":{"dir":` + quoteJSON(activeDir) + `,"name":"w-tdd","role":"engenheiro-tdd"}}}`
		postReq := `{"jsonrpc":"2.0","id":12,"method":"tools/call","params":{"name":"swarm_post","arguments":{"dir":` + quoteJSON(activeDir) + `,"from":"coordenador","to":"@w-tdd","body":"oi"}}}`
		stReq2 := `{"jsonrpc":"2.0","id":13,"method":"tools/call","params":{"name":"swarm_status","arguments":{"dir":` + quoteJSON(activeDir) + `}}}`

		_ = chamarRouter(t, r, addReq)
		_ = chamarRouter(t, r, postReq)
		stResp := chamarRouter(t, r, stReq2)
		if stResp.Result.IsError || !strings.Contains(stResp.Result.Content[0].Text, "w-tdd") || !strings.Contains(stResp.Result.Content[0].Text, `"message_count":1`) {
			t.Fatalf("swarm_status deveria ler .legoloop-swarm/swarm.json unificado, obteve: %+v", stResp.Result)
		}
	})
}

func quoteJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
