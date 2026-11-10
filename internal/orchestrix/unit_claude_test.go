package orchestrix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
)

// helperClaudeTools é o claude falso do UT-03: devolve argv, cwd e stdin como resultado.
func helperClaudeTools() {
	prompt, _ := io.ReadAll(os.Stdin)
	cwd, _ := os.Getwd()
	out, _ := json.Marshal(map[string]any{
		"type":   "result",
		"result": fmt.Sprintf("argv=%s\ncwd=%s\nprompt=%s", strings.Join(os.Args[1:], "|"), cwd, prompt),
		"usage":  map[string]int{"input_tokens": 7, "output_tokens": 3},
	})
	os.Stdout.Write(out)
	os.Exit(0)
}

func toolsOrchestrix() *Orchestrix {
	return &Orchestrix{
		Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: "claude", Path: os.Args[0]}}},
		ExtraEnv:  []string{"LEGOLOOP_HELPER=claude-tools"},
	}
}

func TestCallClaudeToolsExactArgs(t *testing.T) {
	wt := t.TempDir()
	text, usage, err := toolsOrchestrix().callClaudeTools(context.Background(), "sonnet", "faça X", wt, nil)
	if err != nil {
		t.Fatal(err)
	}
	wantArgs := "-p|--model|sonnet|--output-format|json|--no-session-persistence|--strict-mcp-config|--safe-mode|" +
		"--setting-sources|user|--tools|Read,Edit,Write,Glob,Grep|--allowed-tools|Read,Edit,Write,Glob,Grep|" +
		"--disallowed-tools|Agent,Task|--disable-slash-commands|--effort|medium|--max-budget-usd|2.50"
	if !strings.Contains(text, "argv="+wantArgs+"\n") {
		t.Fatalf("argumentos inexatos:\n%s\nesperado: %s", text, wantArgs)
	}
	want, _ := filepath.EvalSymlinks(wt)
	if !strings.Contains(text, "cwd="+wt+"\n") && !strings.Contains(text, "cwd="+want+"\n") {
		t.Fatalf("cwd não é a worktree %s:\n%s", wt, text)
	}
	if !strings.Contains(text, "prompt=faça X") {
		t.Fatalf("prompt não chegou via stdin:\n%s", text)
	}
	if usage.Input != 7 || usage.Output != 3 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestCallClaudeToolsNeverGrantsAgentsOrBash(t *testing.T) {
	args := strings.Join(unitClaudeArgs("sonnet", unitTools, unitLimitsFor(dag.Node{})), " ")
	for _, bad := range []string{"Bash", "WebFetch", "WebSearch"} {
		if strings.Contains(args, bad) {
			t.Errorf("argumentos concedem %s: %s", bad, args)
		}
	}
	// Agent e Task só aparecem em --disallowed-tools.
	if strings.Count(args, "Agent") != 1 || !strings.HasSuffix(strings.Split(args, "--disallowed-tools ")[1], "Agent,Task --disable-slash-commands --effort medium --max-budget-usd 2.50") {
		t.Errorf("Agent fora de --disallowed-tools: %s", args)
	}
}

func TestCallClaudeToolsRejectsToolOutsideList(t *testing.T) {
	for _, tools := range [][]string{{"Bash"}, {"Read", "Agent"}, {"Task"}, {"read"}, {"Read,Bash"}, {}} {
		_, _, err := toolsOrchestrix().callClaudeTools(context.Background(), "sonnet", "x", t.TempDir(), tools)
		if err == nil {
			t.Errorf("tools %v: esperava erro", tools)
		}
	}
}

func TestCallClaudeToolsRejectsEmptyWorktree(t *testing.T) {
	if _, _, err := toolsOrchestrix().callClaudeTools(context.Background(), "sonnet", "x", "", nil); err == nil {
		t.Fatal("worktree vazia deveria falhar")
	}
}

func TestCallClaudeToolsSubsetIsHonored(t *testing.T) {
	text, _, err := toolsOrchestrix().callClaudeTools(context.Background(), "sonnet", "x", t.TempDir(), []string{"Read", "Grep"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "--tools|Read,Grep|--allowed-tools|Read,Grep|") {
		t.Fatalf("subconjunto ignorado:\n%s", text)
	}
}

func TestCallClaudeToolsPropagatesFailure(t *testing.T) {
	o := toolsOrchestrix()
	o.ExtraEnv = []string{"LEGOLOOP_HELPER=claude-api-error"}
	if _, _, err := o.callClaudeTools(context.Background(), "sonnet", "x", t.TempDir(), nil); err == nil ||
		!strings.Contains(err.Error(), "Credit balance") {
		t.Fatalf("err = %v", err)
	}
}

func TestUnitLimitsFollowNature(t *testing.T) {
	cases := []struct {
		nature string
		effort string
		budget string
	}{{"micro", "low", "0.80"}, {"nano", "low", "0.80"}, {"macro", "medium", "2.50"}, {"", "medium", "2.50"}}
	for _, c := range cases {
		lim := unitLimitsFor(dag.Node{Nature: c.nature})
		args := strings.Join(unitClaudeArgs("sonnet", unitTools, lim), " ")
		if !strings.Contains(args, "--effort "+c.effort) || !strings.Contains(args, "--max-budget-usd "+c.budget) {
			t.Errorf("nature %q: args = %s", c.nature, args)
		}
	}
}

func TestParseClaudeOutputReadsCost(t *testing.T) {
	r := parseClaudeOutput([]byte(`{"result":"ok","total_cost_usd":0.4321,"usage":{"input_tokens":3,"output_tokens":9}}`))
	if r.usage.CostUSD != 0.4321 || r.usage.Output != 9 {
		t.Fatalf("usage = %+v", r.usage)
	}
}
