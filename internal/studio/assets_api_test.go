package studio

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/gate"
)

// newAssetsServer liga o Studio a um inventário isolado (HOME, projeto e cofre temporários).
func newAssetsServer(t *testing.T) (*Server, string, string) {
	t.Helper()
	s, _ := newServer(t)
	home := t.TempDir()
	cfg := assets.Config{
		Home: home, Projects: []string{s.ProjectDir},
		LookPath: func(string) (string, error) { return "", errors.New("fora do PATH") },
		Getenv:   func(string) string { return "" },
	}
	s.Assets = assets.NewWatcher(assets.NewScanner(cfg), time.Hour)
	s.Home = home
	s.Vault = t.TempDir()
	return s, home, s.ProjectDir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// leaks diz se o corpo contém o caminho absoluto (cru ou escapado como JSON no Windows).
func leaks(body, path string) bool {
	return strings.Contains(body, path) || strings.Contains(body, strings.ReplaceAll(path, `\`, `\\`))
}

type summaryBody struct {
	Version uint64         `json:"version"`
	Counts  map[string]int `json:"counts"`
}

type listBody struct {
	Kind  string `json:"kind"`
	Items []struct {
		Kind, Name, Path, Scope, Origin string
	} `json:"items"`
}

func TestAssetsSummaryCountsEveryTab(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	writeFile(t, filepath.Join(home, ".claude", "skills", "alpha", "SKILL.md"), "---\nname: alpha\n---\n")
	rec := do(t, s.Handler(), "GET", "/api/assets", "")
	if rec.Code != 200 {
		t.Fatalf("status %d: %s", rec.Code, rec.Body)
	}
	var b summaryBody
	decode(t, rec, &b)
	for _, k := range []string{"skill", "mcp", "script", "hook", "gate", "harness", "binary", "model", "loop"} {
		if _, ok := b.Counts[k]; !ok {
			t.Errorf("contagem de %q ausente: %v", k, b.Counts)
		}
	}
	if b.Counts["skill"] != 1 || b.Counts["gate"] != len(gate.Catalog()) || b.Counts["model"] != 5 || b.Counts["loop"] != 0 || b.Version == 0 {
		t.Errorf("contagens inesperadas: %+v", b)
	}
}

func TestAssetsListByKindHidesAbsolutePaths(t *testing.T) {
	s, home, project := newAssetsServer(t)
	writeFile(t, filepath.Join(home, ".claude", "skills", "alpha", "SKILL.md"), "---\nname: alpha\n---\n")
	writeFile(t, filepath.Join(project, "scripts", "build.sh"), "echo")
	for _, kind := range []string{"skill", "script"} {
		rec := do(t, s.Handler(), "GET", "/api/assets/"+kind, "")
		body := rec.Body.String()
		if rec.Code != 200 {
			t.Fatalf("%s: %d %s", kind, rec.Code, body)
		}
		if leaks(body, home) || leaks(body, project) {
			t.Errorf("%s: o corpo vaza caminho absoluto: %s", kind, body)
		}
		var b listBody
		decode(t, rec, &b)
		if b.Kind != kind || len(b.Items) != 1 {
			t.Fatalf("%s: %+v", kind, b)
		}
	}
	rec := do(t, s.Handler(), "GET", "/api/assets/skill", "")
	var b listBody
	decode(t, rec, &b)
	if p := filepath.ToSlash(b.Items[0].Path); !strings.HasPrefix(p, "~/.claude/skills/alpha") {
		t.Errorf("caminho de exibição = %q", p)
	}
}

func TestAssetsModelsAndLoopsTabs(t *testing.T) {
	s, _, _ := newAssetsServer(t)
	var models, loops listBody
	decode(t, do(t, s.Handler(), "GET", "/api/assets/model", ""), &models)
	if len(models.Items) != 5 || models.Items[0].Name != "code-pro" || models.Items[0].Origin != "9router" {
		t.Errorf("modelos = %+v", models.Items)
	}
	rec := do(t, s.Handler(), "GET", "/api/assets/loop", "")
	decode(t, rec, &loops)
	if rec.Code != 200 || loops.Items == nil || len(loops.Items) != 0 {
		t.Errorf("a aba Loops lista vazia até o acervo existir: %d %+v", rec.Code, loops)
	}
	if strings.Contains(rec.Body.String(), `"items":null`) {
		t.Error("items deve ser [] e não null")
	}
}

func TestAssetsUnknownKindIs404(t *testing.T) {
	s, _, _ := newAssetsServer(t)
	if rec := do(t, s.Handler(), "GET", "/api/assets/segredos", ""); rec.Code != 404 {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAssetsRescanPicksUpNewSkill(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	var before, after summaryBody
	decode(t, do(t, s.Handler(), "GET", "/api/assets", ""), &before)
	writeFile(t, filepath.Join(home, ".claude", "skills", "nova", "SKILL.md"), "---\nname: nova\n---\n")
	rec := do(t, s.Handler(), "POST", "/api/assets/rescan", "{}")
	if rec.Code != 200 {
		t.Fatalf("rescan: %d %s", rec.Code, rec.Body)
	}
	decode(t, rec, &after)
	if after.Counts["skill"] != before.Counts["skill"]+1 || after.Version <= before.Version {
		t.Errorf("antes %+v, depois %+v", before, after)
	}
}

func TestAssetsUnavailableWithoutWatcher(t *testing.T) {
	s, _ := newServer(t)
	for _, path := range []string{"/api/assets", "/api/assets/skill"} {
		if rec := do(t, s.Handler(), "GET", path, ""); rec.Code != 503 {
			t.Errorf("%s: status %d, quer 503", path, rec.Code)
		}
	}
	if rec := do(t, s.Handler(), "POST", "/api/canonize/plan", "{}"); rec.Code != 503 {
		t.Errorf("canonize sem inventário: status %d, quer 503", rec.Code)
	}
}

type planBody struct {
	FilesScanned int   `json:"files_scanned"`
	SavedBytes   int64 `json:"saved_bytes"`
	Groups       []struct {
		SHA256    string   `json:"sha256"`
		Canonical string   `json:"canonical"`
		Links     []string `json:"links"`
	} `json:"groups"`
}

type applyBody struct {
	Linked     int                            `json:"linked"`
	Failed     []struct{ Path, Error string } `json:"failed"`
	SavedBytes int64                          `json:"saved_bytes"`
}

const dupContent = "---\nname: a\ndescription: duplicada\n---\n# corpo igual nas duas pastas\n"

func twoCopies(t *testing.T, home string) (string, string) {
	t.Helper()
	a := filepath.Join(home, ".claude", "skills", "a", "SKILL.md")
	b := filepath.Join(home, ".agents", "skills", "a", "SKILL.md")
	writeFile(t, a, dupContent)
	writeFile(t, b, dupContent)
	return a, b
}

func TestCanonizePlanIsReadOnly(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	a, b := twoCopies(t, home)
	rec := do(t, s.Handler(), "POST", "/api/canonize/plan", "{}")
	if rec.Code != 200 {
		t.Fatalf("plan: %d %s", rec.Code, rec.Body)
	}
	var p planBody
	decode(t, rec, &p)
	if len(p.Groups) != 1 || len(p.Groups[0].Links) != 2 || p.SavedBytes <= 0 || p.FilesScanned < 2 {
		t.Errorf("plano = %+v", p)
	}
	if leaks(rec.Body.String(), home) {
		t.Errorf("o plano vaza o caminho absoluto: %s", rec.Body)
	}
	for _, f := range []string{a, b} {
		if st, err := os.Lstat(f); err != nil || st.Mode()&os.ModeSymlink != 0 {
			t.Errorf("simular não pode alterar %s (%v)", f, err)
		}
	}
	if entries, _ := os.ReadDir(s.Vault); len(entries) != 0 {
		t.Errorf("simular não pode escrever no cofre: %v", entries)
	}
}

func TestCanonizeApplyNeedsExplicitFlagAndIsIdempotent(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	a, b := twoCopies(t, home)
	for name, body := range map[string]string{"vazio": `{}`, "falso": `{"apply":false}`} {
		if rec := do(t, s.Handler(), "POST", "/api/canonize/apply", body); rec.Code != 400 {
			t.Errorf("%s: status %d, quer 400", name, rec.Code)
		}
	}
	rec := do(t, s.Handler(), "POST", "/api/canonize/apply", `{"apply":true}`)
	if rec.Code != 200 {
		t.Fatalf("apply: %d %s", rec.Code, rec.Body)
	}
	var r applyBody
	decode(t, rec, &r)
	if r.Linked != 2 || len(r.Failed) != 0 {
		t.Errorf("relatório = %+v", r)
	}
	for _, f := range []string{a, b} {
		got, err := os.ReadFile(f)
		if err != nil || string(got) != dupContent {
			t.Errorf("%s perdeu o conteúdo após canonizar: %v", f, err)
		}
	}
	var p planBody
	decode(t, do(t, s.Handler(), "POST", "/api/canonize/plan", "{}"), &p)
	if len(p.Groups) != 0 {
		t.Errorf("depois de aplicar não há mais o que canonizar: %+v", p)
	}
}

// As raízes da canonização vêm só da configuração do servidor: o cliente nunca escolhe caminhos.
func TestCanonizeRejectsRootsFromRequest(t *testing.T) {
	s, _, _ := newAssetsServer(t)
	for _, path := range []string{"/api/canonize/plan", "/api/canonize/apply"} {
		rec := do(t, s.Handler(), "POST", path, `{"apply":true,"roots":["C:\\"]}`)
		if rec.Code != 400 {
			t.Errorf("%s com roots: status %d, quer 400", path, rec.Code)
		}
	}
}

func TestInventoryPageIsServed(t *testing.T) {
	s, _ := newServer(t)
	rec := do(t, s.Handler(), "GET", "/inventario", "")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Inventário") {
		t.Fatalf("/inventario: %d", rec.Code)
	}
	if !strings.HasPrefix(rec.Header().Get("Content-Type"), "text/html") {
		t.Errorf("content-type = %q", rec.Header().Get("Content-Type"))
	}
	if r := do(t, s.Handler(), "GET", "/inventory.js", ""); r.Code != 200 {
		t.Errorf("/inventory.js: %d", r.Code)
	}
}

// Fora do HOME e do projeto nunca se mostra caminho absoluto: só o nome do arquivo.
func TestDisplayPathOutsideKnownBasesShowsOnlyBaseName(t *testing.T) {
	s, home, project := newAssetsServer(t)
	outside := filepath.Join(t.TempDir(), "segredo", "arquivo.md")
	cases := map[string]string{
		filepath.Join(home, ".claude", "skills", "a"): "~/.claude/skills/a",
		filepath.Join(project, "scripts", "x.sh"):     "<projeto>/scripts/x.sh",
		home:    "~",
		project: "<projeto>",
		outside: "arquivo.md",
		"":      "",
	}
	for in, want := range cases {
		if got := s.displayPath(in); got != want {
			t.Errorf("displayPath(%q) = %q, quer %q", in, got, want)
		}
	}
}

// O cofre padrão (~/.aidd/canonical) contém uma raiz de skills: a varredura não pode incluí-lo,
// senão o próprio cofre viraria "duplicata" de si mesmo.
func TestCanonizeNeverScansTheVault(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	s.Vault = filepath.Join(home, ".aidd", "canonical")
	writeFile(t, filepath.Join(s.Vault, "skills", "a", "SKILL.md"), dupContent)
	writeFile(t, filepath.Join(home, ".claude", "skills", "a", "SKILL.md"), dupContent)
	var p planBody
	rec := do(t, s.Handler(), "POST", "/api/canonize/plan", "{}")
	decode(t, rec, &p)
	if len(p.Groups) != 0 || p.FilesScanned != 1 {
		t.Errorf("o arquivo do cofre não pode entrar na varredura: %+v", p)
	}
}

func TestCanonizeApplyRefusesConcurrentRuns(t *testing.T) {
	s, _, _ := newAssetsServer(t)
	s.canonMu.Lock()
	defer s.canonMu.Unlock()
	if rec := do(t, s.Handler(), "POST", "/api/canonize/apply", `{"apply":true}`); rec.Code != 409 {
		t.Errorf("com uma canonização em andamento: status %d, quer 409", rec.Code)
	}
}

// O comando de um hook costuma citar caminhos do HOME ou do projeto: a descrição também é mascarada.
func TestAssetsMaskHomeAndProjectInsideDescriptions(t *testing.T) {
	s, home, project := newAssetsServer(t)
	cmd := filepath.ToSlash(filepath.Join(home, ".orca", "hook.exe")) + " --dir " + filepath.ToSlash(project)
	writeFile(t, filepath.Join(home, ".claude", "settings.json"),
		`{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"`+cmd+`"}]}]}}`)
	rec := do(t, s.Handler(), "GET", "/api/assets/hook", "")
	body := rec.Body.String()
	for _, p := range []string{home, project, filepath.ToSlash(home), filepath.ToSlash(project)} {
		if leaks(body, p) {
			t.Errorf("a descrição vaza %q: %s", p, body)
		}
	}
	var b struct {
		Items []struct{ Description string } `json:"items"`
	}
	decode(t, rec, &b)
	if want := "~/.orca/hook.exe --dir <projeto>"; len(b.Items) != 1 || b.Items[0].Description != want {
		t.Errorf("descrição = %+v, quer %q", b.Items, want)
	}
}

// A exibição do comando de hook é curta e termina em reticências.
func TestHookDescriptionIsTruncatedForDisplay(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	writeFile(t, filepath.Join(home, ".claude", "settings.json"),
		`{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"`+strings.Repeat("a", 900)+`"}]}]}}`)
	var b struct {
		Items []struct{ Description string } `json:"items"`
	}
	decode(t, do(t, s.Handler(), "GET", "/api/assets/hook", ""), &b)
	if n := len([]rune(b.Items[0].Description)); n != maxHookDisplay+1 || !strings.HasSuffix(b.Items[0].Description, "…") {
		t.Errorf("descrição exibida tem %d runas: %q", n, b.Items[0].Description)
	}
}

func TestAssetImportRejectsInvalidURL(t *testing.T) {
	s, _, _ := newAssetsServer(t)

	// URL não github
	rec := do(t, s.Handler(), "POST", "/api/assets/import", `{"url":"https://evil.example.com/repo"}`)
	if rec.Code != 400 {
		t.Errorf("esperava status 400 para domínio não-github, obteve %d", rec.Code)
	}

	// URL vazia
	recEmpty := do(t, s.Handler(), "POST", "/api/assets/import", `{"url":""}`)
	if recEmpty.Code != 400 {
		t.Errorf("esperava status 400 para url vazia, obteve %d", recEmpty.Code)
	}
}

