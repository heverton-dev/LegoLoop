package assets

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSplitCommand(t *testing.T) {
	ok := map[string][]string{
		`echo ok`:         {"echo", "ok"},
		"  echo   a\tb  ": {"echo", "a", "b"},
		`"C:\Program Files\x\a.exe" --flag "v w"`: {`C:\Program Files\x\a.exe`, "--flag", "v w"},
		`node 'um arquivo.js' --x=1`:              {"node", "um arquivo.js", "--x=1"},
		`C:\tools\a.exe sub\caminho`:              {`C:\tools\a.exe`, `sub\caminho`},
		`"C:\Program Files (x86)\a.exe" run`:      {`C:\Program Files (x86)\a.exe`, "run"},
		`prog "a|b" 'c;d'`:                        {"prog", "a|b", "c;d"},
		`prog --empty ""`:                         {"prog", "--empty", ""},
	}
	for in, want := range ok {
		got, err := SplitCommand(in)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("SplitCommand(%q) = %q, %v; quer %q", in, got, err, want)
		}
	}
	bad := []string{
		``, `   `, `"aberta`, `'aberta`, `a | b`, `a && b`, `a ; b`, `a > f`, `a < f`, `echo $HOME`,
		`echo %PATH%`, `(a)`, "echo `x`", "a\nb", `echo "$X"`, "echo \"`x`\"", `a & b`,
		strings.Repeat("a ", maxArgs+1), strings.Repeat("x", maxCommandLen+1),
	}
	for _, in := range bad {
		if got, err := SplitCommand(in); err == nil {
			t.Errorf("SplitCommand(%q) deveria falhar, veio %q", in, got)
		}
	}
}

func TestArgvForScriptsByExtension(t *testing.T) {
	cases := map[string][]string{
		"/x/a.ps1":  {"powershell", "-NoProfile", "-File", "/x/a.ps1"},
		"/x/a.sh":   {"bash", "/x/a.sh"},
		"/x/a.bash": {"bash", "/x/a.bash"},
		"/x/a.py":   {"python", "/x/a.py"},
		"/x/a.js":   {"node", "/x/a.js"},
		"/x/a.mjs":  {"node", "/x/a.mjs"},
		"/x/a.cjs":  {"node", "/x/a.cjs"},
		"/x/A.BAT":  {"cmd", "/c", "/x/A.BAT"},
		"/x/a.cmd":  {"cmd", "/c", "/x/a.cmd"},
	}
	for path, want := range cases {
		got, err := ArgvFor(Item{Kind: KindScript, Name: filepath.Base(path), Path: path})
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("ArgvFor(%s) = %q, %v; quer %q", path, got, err, want)
		}
	}
	if _, err := ArgvFor(Item{Kind: KindScript, Path: "/x/a.txt"}); err == nil {
		t.Error("extensão desconhecida deveria falhar")
	}
	if _, err := ArgvFor(Item{Kind: KindScript}); err == nil {
		t.Error("script sem caminho deveria falhar")
	}
	if _, err := ArgvFor(Item{Kind: KindSkill, Path: "/x/a.sh"}); err == nil {
		t.Error("só scripts e hooks viram argv")
	}
}

// O comando completo do hook (com o segredo) fica só no servidor: a API mostra a descrição redigida.
func TestHookKeepsRawCommandServerSideOnly(t *testing.T) {
	cfg, home, _ := testCfg(t)
	write(t, filepath.Join(home, ".claude", "settings.json"),
		`{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"rtk run --token=abc123 \"a b\""}]}]}}`)
	hook := mustFind(t, scanKind(t, cfg, KindHook), "Stop")
	if hook.Command != `rtk run --token=abc123 "a b"` {
		t.Fatalf("comando bruto = %q", hook.Command)
	}
	b, _ := json.Marshal(hook)
	if strings.Contains(string(b), "abc123") || strings.Contains(string(b), `"Command":`) || strings.Contains(string(b), `"command":`) {
		t.Errorf("o JSON expõe o comando bruto: %s", b)
	}
	argv, err := ArgvFor(hook)
	if err != nil || !reflect.DeepEqual(argv, []string{"rtk", "run", "--token=abc123", "a b"}) {
		t.Errorf("argv = %q, %v", argv, err)
	}
	if _, err := ArgvFor(Item{Kind: KindHook}); err == nil {
		t.Error("hook sem comando deveria falhar")
	}
}

func TestItemIDsAreStableUniqueAndPublic(t *testing.T) {
	cfg, home, _ := testCfg(t)
	write(t, filepath.Join(home, ".claude", "skills", "a", "SKILL.md"), "---\nname: a\n---\n")
	write(t, filepath.Join(home, ".claude", "skills", "b", "SKILL.md"), "---\nname: b\n---\n")
	s := NewScanner(cfg)
	first, second := s.Scan(t.Context()), s.Scan(t.Context())
	seen := map[string]string{}
	for i, it := range first.Items {
		if len(it.ID) != 12 || it.ID != second.Items[i].ID {
			t.Fatalf("ID instável ou malformado: %q vs %q", it.ID, second.Items[i].ID)
		}
		if other, dup := seen[it.ID]; dup {
			t.Errorf("ID repetido entre %s e %s/%s", other, it.Kind, it.Name)
		}
		seen[it.ID] = string(it.Kind) + "/" + it.Name
	}
	b, _ := json.Marshal(mustFind(t, first.Of(KindSkill), "a"))
	if !strings.Contains(string(b), `"id":"`) {
		t.Errorf("o JSON deve expor o id: %s", b)
	}
	if got, ok := first.ByID(mustFind(t, first.Of(KindSkill), "b").ID); !ok || got.Name != "b" {
		t.Errorf("ByID falhou: %+v %v", got, ok)
	}
	if _, ok := first.ByID("nao-existe"); ok {
		t.Error("ByID de id inexistente deve falhar")
	}
}
