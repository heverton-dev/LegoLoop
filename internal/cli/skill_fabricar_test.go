package cli

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	"legoloop/internal/assets"
	"legoloop/internal/gate"
)

const fabricarSkillName = "legoloop-fabricar"

// repoRoot sobe de internal/cli até a raiz do repositório.
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func readFabricarSkill(t *testing.T) (front map[string]string, body string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(repoRoot(t), "skills", fabricarSkillName, "SKILL.md"))
	if err != nil {
		t.Fatalf("a skill precisa existir: %v", err)
	}
	text := strings.ReplaceAll(string(b), "\r\n", "\n")
	rest, ok := strings.CutPrefix(text, "---\n")
	if !ok {
		t.Fatal("SKILL.md sem frontmatter")
	}
	head, body, ok := strings.Cut(rest, "\n---\n")
	if !ok {
		t.Fatal("frontmatter não fechado")
	}
	front = map[string]string{}
	for _, l := range strings.Split(head, "\n") {
		if k, v, ok := strings.Cut(l, ": "); ok {
			front[k] = strings.TrimSpace(v)
		}
	}
	return front, body
}

func TestFabricarSkill(t *testing.T) {
	t.Run("o assets.Scanner lista a skill do repositório", func(t *testing.T) {
		snap := assets.NewScanner(assets.Config{
			Home: t.TempDir(), Projects: []string{repoRoot(t)},
			LookPath: func(string) (string, error) { return "", os.ErrNotExist }, Getenv: func(string) string { return "" },
		}).Scan(context.Background())
		for _, it := range snap.Of(assets.KindSkill) {
			if it.Name == fabricarSkillName {
				return
			}
		}
		t.Errorf("assets.Scanner não listou %s", fabricarSkillName)
	})

	t.Run("frontmatter valido com name igual ao diretorio", func(t *testing.T) {
		front, _ := readFabricarSkill(t)
		if front["name"] != fabricarSkillName {
			t.Errorf("name = %q, esperado %q (igual ao diretório)", front["name"], fabricarSkillName)
		}
		desc := front["description"]
		if n := utf8.RuneCountInString(desc); n == 0 || n > 300 {
			t.Errorf("description com %d runas; o limite é 300", n)
		}
		if !strings.HasPrefix(desc, "Use quando") || !strings.Contains(desc, "NÃO use") {
			t.Errorf("description deve começar pelo gatilho e dizer quando não usar: %q", desc)
		}
		for _, k := range []string{"consumes", "produces", "artifact", "gate", "suggested_model"} {
			if front[k] == "" {
				t.Errorf("frontmatter sem %q", k)
			}
		}
		if _, err := gate.Compile(front["gate"]); err != nil {
			t.Errorf("gate do frontmatter não compila: %v", err)
		}
	})

	t.Run("o corpo conduz entrevista, para no stop-gate e delega ao comando", func(t *testing.T) {
		_, body := readFabricarSkill(t)
		low := strings.ToLower(body)
		for _, want := range []string{"entrevista", "stop-gate", "legoloop fabricar", "--yes"} {
			if !strings.Contains(low, want) {
				t.Errorf("o corpo não menciona %q", want)
			}
		}
		assertGateBeforeWrite(t, body)
		if regexp.MustCompile("(?i)legoloop canonize").MatchString(body) {
			t.Error("a skill não pode citar `legoloop canonize`: é outra função (desduplicação)")
		}
		if strings.Contains(body, "{{") {
			t.Error("o corpo tem placeholder {{...}} sem preencher")
		}
		if lines := strings.Count(body, "\n"); lines > 120 {
			t.Errorf("a skill duplica lógica em prompt: %d linhas (máx. 120); a lógica fica no comando", lines)
		}
	})
}

// assertGateBeforeWrite exige que a primeira invocação que GRAVA (sem --dry-run) venha depois do stop-gate do
// procedimento. A prévia com --dry-run pode vir antes, porque não escreve nada.
func assertGateBeforeWrite(t *testing.T, body string) {
	t.Helper()
	proc := strings.Index(body, "### Procedimento")
	if proc < 0 {
		t.Fatal("a skill precisa ter a seção '### Procedimento'")
	}
	section := body[proc:]
	gateAt := strings.Index(strings.ToLower(section), "stop-gate humano")
	var writeAt = -1
	for _, m := range regexp.MustCompile("legoloop fabricar \"<nome>\"[^\n`]*").FindAllStringIndex(section, -1) {
		if !strings.Contains(section[m[0]:m[1]], "--dry-run") {
			writeAt = m[0]
			break
		}
	}
	if gateAt < 0 || writeAt < 0 || gateAt > writeAt {
		t.Errorf("o stop-gate humano deve vir antes da invocação que grava (gate@%d, grava@%d)", gateAt, writeAt)
	}
	if !strings.Contains(section, "--dry-run") {
		t.Error("a skill deve fazer uma prévia com --dry-run antes do gate")
	}
}
