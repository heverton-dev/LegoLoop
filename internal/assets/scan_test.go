package assets

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestScanAggregatesSortsAndCounts(t *testing.T) {
	cfg, home, project := testCfg(t, "git")
	write(t, filepath.Join(home, ".claude", "skills", "zeta", "SKILL.md"), "---\nname: zeta\n---\n")
	write(t, filepath.Join(home, ".claude", "skills", "alpha", "SKILL.md"), "---\nname: alpha\n---\n")
	write(t, filepath.Join(project, "scripts", "b.sh"), "1")
	snap := NewScanner(cfg).Scan(t.Context())

	skills := snap.Of(KindSkill)
	if len(skills) != 2 || skills[0].Name != "alpha" || skills[1].Name != "zeta" {
		t.Fatalf("skills fora de ordem: %+v", skills)
	}
	if snap.Counts[KindSkill] != 2 || snap.Counts[KindScript] != 1 || snap.Counts[KindGate] != 6 || snap.Counts[KindBinary] != 1 {
		t.Errorf("counts = %v", snap.Counts)
	}
	for _, k := range Kinds() {
		if _, ok := snap.Counts[k]; !ok {
			t.Errorf("contagem de %s ausente (deve existir mesmo zerada)", k)
		}
	}
	if snap.ScannedAt.IsZero() {
		t.Error("ScannedAt não preenchido")
	}
}

// R-INV-2: a re-varredura só relê o que mudou (caminho + mtime + tamanho).
func TestCacheAvoidsReparsingUnchangedFiles(t *testing.T) {
	cfg, home, _ := testCfg(t)
	skill := filepath.Join(home, ".claude", "skills", "alpha", "SKILL.md")
	write(t, skill, "---\nname: alpha\ndescription: v1\n---\n")
	s := NewScanner(cfg)
	var parses atomic.Int32
	s.onParse = func(string) { parses.Add(1) }

	s.Scan(t.Context())
	first := parses.Load()
	if first == 0 {
		t.Fatal("primeira varredura deveria parsear")
	}
	s.Scan(t.Context())
	if parses.Load() != first {
		t.Errorf("arquivo intacto foi relido (%d → %d)", first, parses.Load())
	}

	write(t, skill, "---\nname: alpha\ndescription: versão dois, maior\n---\n")
	later := time.Now().Add(2 * time.Second)
	if err := os.Chtimes(skill, later, later); err != nil {
		t.Fatal(err)
	}
	snap := s.Scan(t.Context())
	if parses.Load() == first {
		t.Error("arquivo alterado deveria ser relido")
	}
	if d := mustFind(t, snap.Of(KindSkill), "alpha").Description; d != "versão dois, maior" {
		t.Errorf("descrição = %q", d)
	}
}

func TestCacheDropsDeletedFiles(t *testing.T) {
	cfg, home, _ := testCfg(t)
	skill := filepath.Join(home, ".claude", "skills", "alpha", "SKILL.md")
	write(t, skill, "---\nname: alpha\n---\n")
	s := NewScanner(cfg)
	if len(s.Scan(t.Context()).Of(KindSkill)) != 1 {
		t.Fatal("esperava 1 skill")
	}
	if err := os.Remove(skill); err != nil {
		t.Fatal(err)
	}
	if n := len(s.Scan(t.Context()).Of(KindSkill)); n != 0 {
		t.Errorf("skill removida ainda aparece (%d)", n)
	}
	if len(s.cache) != 0 {
		t.Errorf("cache retém %d entradas de arquivos que não existem mais", len(s.cache))
	}
}

func TestRootsListsOnlyExistingDirs(t *testing.T) {
	cfg, home, project := testCfg(t)
	write(t, filepath.Join(home, ".claude", "skills", "a", "SKILL.md"), "x")
	write(t, filepath.Join(project, "scripts", "a.sh"), "x")
	roots := NewScanner(cfg).Roots()
	want := map[string]bool{
		filepath.Join(home, ".claude", "skills"): true,
		filepath.Join(project, "scripts"):        true,
	}
	for _, r := range roots {
		if !want[r] {
			t.Errorf("raiz inesperada (ou inexistente): %s", r)
		}
		delete(want, r)
	}
	if len(want) != 0 {
		t.Errorf("raízes ausentes: %v", want)
	}
}

func TestScanHonorsCancelledContext(t *testing.T) {
	cfg, home, _ := testCfg(t)
	write(t, filepath.Join(home, ".claude", "skills", "a", "SKILL.md"), "x")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	done := make(chan struct{})
	go func() { NewScanner(cfg).Scan(ctx); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Scan não respeitou o contexto cancelado")
	}
}

// R-INV-3: um ativo novo aparece sozinho, sem ação do usuário.
func TestWatcherPicksUpNewAssets(t *testing.T) {
	cfg, home, _ := testCfg(t)
	w := NewWatcher(NewScanner(cfg), 20*time.Millisecond)
	ctx, cancel := context.WithCancel(t.Context())
	w.Start(ctx)

	first := w.Refresh(ctx)
	if first.Version == 0 {
		t.Fatal("Refresh deveria publicar a versão 1")
	}
	write(t, filepath.Join(home, ".claude", "skills", "nova", "SKILL.md"), "---\nname: nova\n---\n")
	deadline := time.Now().Add(5 * time.Second)
	for {
		snap := w.Snapshot()
		if _, ok := find(snap.Of(KindSkill), "nova"); ok {
			if snap.Version <= first.Version {
				t.Errorf("a versão deveria subir (%d → %d)", first.Version, snap.Version)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("o vigia não detectou a skill nova")
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	select {
	case <-w.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("o vigia não parou ao cancelar o contexto")
	}
}

func TestWatcherKeepsVersionWhenNothingChanged(t *testing.T) {
	cfg, home, _ := testCfg(t)
	write(t, filepath.Join(home, ".claude", "skills", "a", "SKILL.md"), "---\nname: a\n---\n")
	w := NewWatcher(NewScanner(cfg), time.Hour)
	v1 := w.Refresh(t.Context()).Version
	v2 := w.Refresh(t.Context()).Version
	if v1 != v2 {
		t.Errorf("sem mudança a versão não pode subir (%d → %d)", v1, v2)
	}
}

func TestSnapshotBeforeAnyScanIsEmptyVersionZero(t *testing.T) {
	cfg, _, _ := testCfg(t)
	w := NewWatcher(NewScanner(cfg), time.Hour)
	if s := w.Snapshot(); s.Version != 0 || len(s.Items) != 0 {
		t.Errorf("snapshot inicial = %+v", s)
	}
}

// O mesmo arquivo alimenta dois parsers (settings.json tem hooks E mcpServers). O cache precisa
// distinguir quem o leu, senão uma re-varredura devolve os itens de um parser ao outro.
func TestCacheDoesNotMixParsersReadingTheSameFile(t *testing.T) {
	cfg, home, _ := testCfg(t)
	write(t, filepath.Join(home, ".claude", "settings.json"), `{
	  "hooks": {"Stop": [{"hooks": [{"type": "command", "command": "notifier done"}]}],
	            "PreToolUse": [{"matcher": "Bash", "hooks": [{"type": "command", "command": "rtk run"}]}]},
	  "mcpServers": {"rea": {"command": "npx"}}
	}`)
	s := NewScanner(cfg)
	for i := range 4 {
		snap := s.Scan(t.Context())
		if h, m := len(snap.Of(KindHook)), len(snap.Of(KindMCP)); h != 2 || m != 1 {
			t.Fatalf("varredura %d: hooks=%d (quer 2), mcps=%d (quer 1)", i+1, h, m)
		}
	}
}
