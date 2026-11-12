package mcp

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func escreverSkill(t *testing.T, root, nome, desc, categoria string) {
	t.Helper()
	dir := filepath.Join(root, nome)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	conteudo := "---\nname: " + nome + "\ndescription: " + desc + "\ntier: Ouro\nnatureza: Processo\ncategoria: " + categoria + "\n---\n# corpo " + nome + "\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

type respostaMCP struct {
	Result struct {
		IsError bool `json:"isError"`
		Tools   []struct {
			Name string `json:"name"`
		} `json:"tools"`
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"result"`
}

func chamarRouter(t *testing.T, r *Router, req string) respostaMCP {
	t.Helper()
	var out bytes.Buffer
	if err := r.Serve(strings.NewReader(req+"\n"), &out); err != nil {
		t.Fatal(err)
	}
	var resp respostaMCP
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("resposta invalida: %v: %s", err, out.String())
	}
	return resp
}

func TestSkillSearchECallSkill(t *testing.T) {
	root := t.TempDir()
	escreverSkill(t, root, "aidd-tdd", "Ciclo red green refactor", "build")
	escreverSkill(t, root, "aidd-deploy", "Publicacao em producao", "ops")
	escreverSkill(t, root, "aidd-docs", "Escrita de documentacao", "docs")
	engine := NewSkillEngine()
	if err := engine.IndexDirectory(root); err != nil {
		t.Fatal(err)
	}
	r := NewRouter(nil, nil)
	r.SetSkills(engine, root)

	list := chamarRouter(t, r, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	if len(list.Result.Tools) != 4 {
		t.Fatalf("esperava 4 ferramentas, obteve %d", len(list.Result.Tools))
	}

	busca := chamarRouter(t, r, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"skill_search","arguments":{"query":"red green refactor","top_k":3}}}`)
	var matches []SkillMatch
	if err := json.Unmarshal([]byte(busca.Result.Content[0].Text), &matches); err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 || matches[0].ID != "aidd-tdd" || matches[0].Tier != "Ouro" || matches[0].Categoria != "build" || matches[0].Natureza != "Processo" {
		t.Fatalf("busca inesperada: %+v", matches)
	}

	call := chamarRouter(t, r, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"call_skill","arguments":{"skill_id":"aidd-tdd"}}}`)
	if call.Result.IsError || !strings.Contains(call.Result.Content[0].Text, "# corpo aidd-tdd") {
		t.Fatalf("call_skill inesperado: %+v", call.Result)
	}

	fora := chamarRouter(t, r, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"call_skill","arguments":{"skill_id":"../../etc/passwd"}}}`)
	if !fora.Result.IsError {
		t.Fatal("skill_id fora do cofre deveria falhar")
	}
}

func TestCallSkillRecusaCorpoForaDoCofre(t *testing.T) {
	root := t.TempDir()
	outro := t.TempDir()
	escreverSkill(t, outro, "externa", "Skill externa", "build")
	engine := NewSkillEngine()
	if err := engine.IndexDirectory(outro); err != nil {
		t.Fatal(err)
	}
	r := NewRouter(nil, nil)
	r.SetSkills(engine, root)
	resp := chamarRouter(t, r, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"call_skill","arguments":{"skill_id":"externa"}}}`)
	if !resp.Result.IsError {
		t.Fatal("corpo fora do cofre deveria falhar")
	}
}
