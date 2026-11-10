package suggest

import (
	"fmt"
	"strings"
	"testing"

	"legoloop/internal/assets"
	"legoloop/internal/gate"
	"legoloop/internal/graph"
)

func index(t *testing.T) *graph.Index {
	t.Helper()
	idx, err := graph.NewIndex(graph.DefaultCatalog())
	if err != nil {
		t.Fatal(err)
	}
	return idx
}

func blockCtx(t *testing.T, idx *graph.Index, id string) Context {
	t.Helper()
	b, ok := idx.Block(id)
	if !ok {
		t.Fatalf("bloco %q não existe", id)
	}
	return Context{Block: &b, Artifact: b.Artifact}
}

func exprs(gs []Gate) []string {
	var out []string
	for _, g := range gs {
		out = append(out, g.Expr)
	}
	return out
}

func has(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

func snap(items ...assets.Item) assets.Snapshot { return assets.Snapshot{Items: items} }

func script(name, scope, path string) assets.Item {
	return assets.Item{ID: "s-" + name + scope, Kind: assets.KindScript, Name: name, Scope: scope, Path: path}
}

func hook(event, cmd string) assets.Item {
	return assets.Item{ID: "h-" + event + cmd, Kind: assets.KindHook, Name: event, Scope: "global", Origin: "claude",
		Meta: map[string]string{"event": event}, Command: cmd}
}

func TestGatesForSpecBlockExcludeOwnGateAndAlwaysCompile(t *testing.T) {
	idx := index(t)
	r := Suggest(blockCtx(t, idx, "aidd-spec"), idx, snap())
	got := exprs(r.Gates)
	for _, want := range []string{"file_exists('SPEC.md')", "word_count('SPEC.md') > 100", "line_count('SPEC.md') > 10",
		"has_acceptance_criteria('SPEC.md')"} {
		if !has(got, want) {
			t.Errorf("faltou a sugestão %q em %v", want, got)
		}
	}
	if has(got, "has_binary_criteria('SPEC.md')") {
		t.Error("o gate próprio do bloco não deve ser sugerido de novo")
	}
	for _, g := range r.Gates {
		if _, err := gate.Compile(g.Expr); err != nil {
			t.Errorf("sugestão inválida %q: %v", g.Expr, err)
		}
		if strings.TrimSpace(g.Label) == "" {
			t.Errorf("sugestão sem rótulo: %+v", g)
		}
	}
}

func TestGatesDependOnArtifactTypeAndExtension(t *testing.T) {
	idx := index(t)
	// JSON: só existência (o gate próprio do bloco rea já é esse ⇒ nada sobra: estado vazio explícito)
	if r := Suggest(blockCtx(t, idx, "rea"), idx, snap()); len(r.Gates) != 0 {
		t.Errorf("rea (json, gate próprio file_exists) não deveria sugerir nada: %v", exprs(r.Gates))
	}
	// peça genérica .md: oferece existência, tamanho e critérios
	g := exprs(Suggest(Context{Kind: "skill", Artifact: "SAIDA.md"}, idx, snap()).Gates)
	for _, want := range []string{"file_exists('SAIDA.md')", "word_count('SAIDA.md') > 100", "has_binary_criteria('SAIDA.md')"} {
		if !has(g, want) {
			t.Errorf("genérica .md sem %q: %v", want, g)
		}
	}
	// peça genérica .json: só existência
	if g := exprs(Suggest(Context{Kind: "skill", Artifact: "dados.json"}, idx, snap()).Gates); len(g) != 1 || g[0] != "file_exists('dados.json')" {
		t.Errorf("genérica .json = %v", g)
	}
}

func TestGatesRefuseArtifactNamesThatCouldBreakTheExpression(t *testing.T) {
	idx := index(t)
	for _, a := range []string{"", "a'b.md", `a\b.md`, "a\nb.md"} {
		if g := Suggest(Context{Kind: "skill", Artifact: a}, idx, snap()).Gates; len(g) != 0 {
			t.Errorf("artefato %q não deveria gerar gates: %v", a, exprs(g))
		}
	}
}

func TestScriptsRankProjectFirstAndSkipUnrunnable(t *testing.T) {
	idx := index(t)
	items := []assets.Item{
		script("zeta.sh", "global", "/g/zeta.sh"),
		script("alfa.sh", "global", "/g/alfa.sh"),
		script("build.ps1", "project", "/p/build.ps1"),
		script("notas.txt", "project", "/p/notas.txt"), // sem interpretador ⇒ fora
	}
	r := Suggest(Context{Kind: "skill", Artifact: "A.md"}, idx, snap(items...))
	var names []string
	for _, s := range r.Scripts {
		names = append(names, s.Name)
	}
	if strings.Join(names, ",") != "build.ps1,alfa.sh,zeta.sh" {
		t.Errorf("scripts = %v", names)
	}
}

func TestSuggestionListsAreCapped(t *testing.T) {
	idx := index(t)
	var items []assets.Item
	for i := range MaxItems + 10 {
		items = append(items, script(fmt.Sprintf("s%03d.sh", i), "global", fmt.Sprintf("/g/s%03d.sh", i)))
	}
	if n := len(Suggest(Context{Kind: "skill", Artifact: "A.md"}, idx, snap(items...)).Scripts); n != MaxItems {
		t.Errorf("scripts = %d, quer %d", n, MaxItems)
	}
}

func TestHooksSplitByEventAndSkipWhatNeedsAShell(t *testing.T) {
	idx := index(t)
	items := []assets.Item{
		hook("PreToolUse", `rtk hook run`),
		hook("Stop", `notifier done`),
		hook("Elicitation", `rtk elicit`),          // evento desconhecido ⇒ pre e post
		hook("PostToolUse", `echo a | tee b`),      // exige shell ⇒ fora
		hook("SessionStart", `tools/relativo.exe`), // caminho relativo ⇒ fora
		hook("SessionEnd", ``),                     // sem comando ⇒ fora
	}
	r := Suggest(Context{Kind: "skill", Artifact: "A.md"}, idx, snap(items...))
	names := func(l []assets.Item) string {
		var n []string
		for _, h := range l {
			n = append(n, h.Name)
		}
		return strings.Join(n, ",")
	}
	if got := names(r.HooksPre); got != "Elicitation,PreToolUse" {
		t.Errorf("pre = %q", got)
	}
	if got := names(r.HooksPost); got != "Elicitation,Stop" {
		t.Errorf("post = %q", got)
	}
	if r.SkippedHooks != 3 {
		t.Errorf("SkippedHooks = %d, quer 3", r.SkippedHooks)
	}
}

func TestNeighborsForTypedBlockListInstalledFirstThenMCPs(t *testing.T) {
	idx := index(t)
	mcps := []assets.Item{{ID: "m1", Kind: assets.KindMCP, Name: "playwright"}, {ID: "m2", Kind: assets.KindMCP, Name: "rea"}}
	// interview produz Briefing; consumidores: aidd-spec (instalada) — só ele consome Briefing
	skills := []assets.Item{{ID: "k1", Kind: assets.KindSkill, Name: "AIDD-Spec"}}
	r := Suggest(blockCtx(t, idx, "interview"), idx, snap(append(skills, mcps...)...))
	if len(r.Neighbors) < 3 {
		t.Fatalf("vizinhos = %+v", r.Neighbors)
	}
	first := r.Neighbors[0]
	if first.Kind != "block" || first.ID != "aidd-spec" || !first.Installed {
		t.Errorf("primeiro vizinho = %+v", first)
	}
	var kinds []string
	for _, n := range r.Neighbors {
		kinds = append(kinds, n.Kind)
	}
	if !has(kinds, "mcp") {
		t.Errorf("MCPs do inventário devem aparecer como vizinhos: %v", kinds)
	}
}

func TestNeighborsMarkMissingSkillsAndSortThemLast(t *testing.T) {
	idx := index(t)
	// peça genérica aceita qualquer bloco; só 'aidd-tdd' está instalada
	skills := []assets.Item{{ID: "k1", Kind: assets.KindSkill, Name: "aidd-tdd"}}
	r := Suggest(Context{Kind: "skill", Artifact: "A.md"}, idx, snap(skills...))
	var blocks []Neighbor
	for _, n := range r.Neighbors {
		if n.Kind == "block" {
			blocks = append(blocks, n)
		}
	}
	if len(blocks) == 0 || blocks[0].ID != "aidd-tdd" || !blocks[0].Installed {
		t.Fatalf("instalada deve vir primeiro: %+v", blocks)
	}
	for _, b := range blocks[1:] {
		if b.Installed {
			t.Errorf("%s não está instalada", b.ID)
		}
	}
	if len(blocks) > MaxNeighbors {
		t.Errorf("vizinhos acima do limite: %d", len(blocks))
	}
}

func TestSuggestWithEmptyInventoryNeverPanicsAndReturnsEmptyListsNotNil(t *testing.T) {
	idx := index(t)
	r := Suggest(Context{Kind: "skill", Artifact: "A.md"}, idx, snap())
	if r.Scripts == nil || r.HooksPre == nil || r.HooksPost == nil || r.Neighbors == nil || r.Gates == nil {
		t.Errorf("listas devem ser vazias, não nil: %+v", r)
	}
}

// Com um catálogo maior que o limite, só MaxNeighbors peças tipadas são sugeridas.
func TestNeighborsAreCappedWithALargeCatalog(t *testing.T) {
	var catalog []graph.Block
	for i := range MaxNeighbors + 5 {
		catalog = append(catalog, graph.Block{
			ID: fmt.Sprintf("b%02d", i), Label: fmt.Sprintf("Bloco %02d", i), Ref: fmt.Sprintf("skills/b%02d", i),
			Consumes: []graph.ArtifactType{graph.Idea}, Produces: graph.Briefing,
			Artifact: fmt.Sprintf("B%02d.md", i), Gate: fmt.Sprintf("file_exists('B%02d.md')", i), Cost: 1,
		})
	}
	idx, err := graph.NewIndex(catalog)
	if err != nil {
		t.Fatal(err)
	}
	r := Suggest(Context{Kind: "skill", Artifact: "A.md"}, idx, snap())
	if len(r.Neighbors) != MaxNeighbors {
		t.Errorf("vizinhos = %d, quer %d", len(r.Neighbors), MaxNeighbors)
	}
}

// MCPs com o mesmo nome (configurados em vários harnesses) aparecem uma vez; o total tem teto.
func TestMCPNeighborsAreDeduplicatedAndCapped(t *testing.T) {
	idx := index(t)
	var items []assets.Item
	for i := range MaxMCPNeighbors + 4 {
		items = append(items, assets.Item{ID: fmt.Sprintf("m%02d", i), Kind: assets.KindMCP, Name: fmt.Sprintf("mcp%02d", i), Origin: "claude"})
	}
	items = append(items, assets.Item{ID: "dup", Kind: assets.KindMCP, Name: "mcp00", Origin: "cursor"})
	var mcps []Neighbor
	for _, n := range Suggest(Context{Kind: "skill", Artifact: "A.md"}, idx, snap(items...)).Neighbors {
		if n.Kind == "mcp" {
			mcps = append(mcps, n)
		}
	}
	if len(mcps) != MaxMCPNeighbors {
		t.Fatalf("MCPs sugeridos = %d, quer %d", len(mcps), MaxMCPNeighbors)
	}
	seen := map[string]bool{}
	for _, m := range mcps {
		if seen[m.Name] {
			t.Errorf("MCP %s repetido", m.Name)
		}
		seen[m.Name] = true
	}
}

// Um comando com credencial nunca é sugerido: ele iria literalmente para o arquivo DAG.
func TestHooksWithCredentialsAreNeverSuggested(t *testing.T) {
	idx := index(t)
	items := []assets.Item{
		hook("Stop", `notifier done`),
		hook("Stop", `notifier --token=abc123`),
		hook("Stop", `API_KEY=zzz node x.js`),
		hook("Stop", `curl -H "Authorization: Bearer q123"`),
	}
	r := Suggest(Context{Kind: "skill", Artifact: "A.md"}, idx, snap(items...))
	if len(r.HooksPost) != 1 || r.HooksPost[0].Command != "notifier done" {
		t.Errorf("hooks pós = %+v", r.HooksPost)
	}
	if r.SkippedHooks != 3 {
		t.Errorf("SkippedHooks = %d, quer 3", r.SkippedHooks)
	}
}

func TestContainsSecretMatchesTheRedactionRules(t *testing.T) {
	for cmd, want := range map[string]bool{
		"notifier done": false, "run --token=abc": true, "API_KEY=zzz x": true,
		"curl -H Authorization: Bearer q": true, "PASSWORD: s3": true, "echo ok": false,
	} {
		if got := assets.ContainsSecret(cmd); got != want {
			t.Errorf("ContainsSecret(%q) = %v, quer %v", cmd, got, want)
		}
	}
}

func TestLoopNeighborsAreSuggested(t *testing.T) {
	idx := index(t)
	loopItem := assets.Item{
		ID:          "loop-pesquisa",
		Kind:        assets.KindLoop,
		Name:        "Pesquisa de Mercado",
		Description: "Loop salvo no acervo",
	}
	r := Suggest(Context{Kind: "skill", Artifact: "A.md"}, idx, snap(loopItem))
	var found bool
	for _, n := range r.Neighbors {
		if n.Kind == "loop" && n.ID == "loop-pesquisa" && n.Name == "Pesquisa de Mercado" && n.Installed {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("sub-loop não encontrado nos vizinhos sugeridos: %+v", r.Neighbors)
	}
}
