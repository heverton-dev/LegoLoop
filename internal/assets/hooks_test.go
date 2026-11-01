package assets

import (
	"path/filepath"
	"strings"
	"testing"
)

const settingsWithHooks = `{
  "hooks": {
    "PreToolUse": [{"matcher": "Bash", "hooks": [{"type": "command", "command": "rtk hook --token=abc123 run"}]}],
    "Stop": [{"hooks": [{"type": "command", "command": "echo done"}]}],
    "Quebrado": "texto",
    "SemHooks": [{"matcher": "x"}]
  }
}`

func TestScanHooksAllScopes(t *testing.T) {
	cfg, home, project := testCfg(t)
	write(t, filepath.Join(home, ".claude", "settings.json"), settingsWithHooks)
	write(t, filepath.Join(project, ".claude", "settings.json"), `{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"echo projeto"}]}]}}`)
	write(t, filepath.Join(project, ".claude", "settings.local.json"), `{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"echo local"}]}]}}`)
	items := scanKind(t, cfg, KindHook)

	pre := mustFind(t, items, "PreToolUse:Bash")
	if pre.Scope != "global" || pre.Origin != "claude" || pre.Meta["event"] != "PreToolUse" || pre.Meta["matcher"] != "Bash" {
		t.Errorf("pre = %+v", pre)
	}
	if !strings.Contains(pre.Description, "token=***") {
		t.Errorf("comando deveria vir redigido: %q", pre.Description)
	}
	mustNotLeak(t, items, "abc123")

	scopes := map[string]bool{}
	for _, it := range items {
		if it.Name == "Stop" {
			scopes[it.Scope] = true
		}
	}
	for _, s := range []string{"global", "project", "local"} {
		if !scopes[s] {
			t.Errorf("hook Stop do escopo %s ausente (%v)", s, scopes)
		}
	}
	if _, bad := find(items, "Quebrado"); bad {
		t.Error("evento malformado deve ser ignorado")
	}
}

func TestHookCommandIsTruncated(t *testing.T) {
	cfg, home, _ := testCfg(t)
	long := strings.Repeat("a", 5000)
	write(t, filepath.Join(home, ".claude", "settings.json"),
		`{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"`+long+`"}]}]}}`)
	d := mustFind(t, scanKind(t, cfg, KindHook), "Stop").Description
	if n := len([]rune(d)); n > 1001 || n < 900 {
		t.Errorf("descrição deveria ser truncada (%d runas)", len([]rune(d)))
	}
}

func TestRedact(t *testing.T) {
	cases := map[string]string{
		"run --token=abc":                    "run --token=***",
		"API_KEY=zzz node x.js":              "API_KEY=*** node x.js",
		"PASSWORD: s3nha":                    "PASSWORD: ***",
		"echo ok":                            "echo ok",
		"curl -H Authorization: Bearer q123": "curl -H Authorization: *** ***",
		"x --header bearer abc9":             "x --header bearer ***",
	}
	for in, want := range cases {
		if got := redact(in); got != want {
			t.Errorf("redact(%q) = %q, quer %q", in, got, want)
		}
	}
}

func TestGuardHooksRecognition(t *testing.T) {
	cfg, home, project := testCfg(t)
	guardSettings := `{
  "hooks": {
    "PreToolUse": [
      {"matcher": "Bash", "hooks": [{"type": "command", "command": "python3 .claude/hooks/guard-delete-outside.py"}]},
      {"matcher": "Bash", "hooks": [{"type": "command", "command": "python3 .claude/hooks/guard-worktree-path.py"}]}
    ],
    "PreCommit": [
      {"hooks": [{"type": "command", "command": "python3 .claude/hooks/guard-secrets.py --key=sk-123456"}]},
      {"hooks": [{"type": "command", "command": "python3 .claude/hooks/guard-protected-push.py"}]}
    ]
  }
}`
	write(t, filepath.Join(project, ".claude", "settings.json"), guardSettings)
	_ = home
	items := scanKind(t, cfg, KindHook)

	preBash := mustFind(t, items, "PreToolUse:Bash")
	if !strings.Contains(preBash.Command, "guard-delete-outside.py") {
		t.Errorf("esperava comando guard-delete-outside.py preservado em Command, obteve %q", preBash.Command)
	}

	preCommit := mustFind(t, items, "PreCommit")
	if strings.Contains(preCommit.Description, "sk-123456") {
		t.Errorf("segredo não deveria vazar na descrição do hook: %q", preCommit.Description)
	}
}

