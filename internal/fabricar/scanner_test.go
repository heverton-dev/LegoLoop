package fabricar_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/fabricar"
)

const goodSkill = "---\nname: %s\ndescription: skill de teste\n---\n# Skill\n\n## Negative Guardrails\nNUNCA faça X. Zero Stubs.\n\n- go test ./...\n- exit code == 0\n\nConsumes: a.md\nProduces: b.md\n\n## Stopping Checklist\n"

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func skillAt(t *testing.T, dir, name string) string {
	t.Helper()
	p := filepath.Join(dir, name, "SKILL.md")
	writeFile(t, p, strings.Replace(goodSkill, "%s", name, 1))
	return p
}

func scanner(untrusted ...string) *fabricar.DiskScanner {
	s := fabricar.NewDiskScanner()
	s.Untrusted = func(root string) bool { return slices.Contains(untrusted, root) }
	return s
}

func TestScanStorageLocations_FindsPiecesInEveryRootAndFlagsUntrusted(t *testing.T) {
	trusted, removable := t.TempDir(), t.TempDir()
	skillAt(t, filepath.Join(trusted, ".claude", "skills"), "alfa")
	skillAt(t, filepath.Join(removable, "skills"), "beta")
	skillAt(t, filepath.Join(removable, "Users", "ana", ".claude", "skills"), "gama")

	res, err := scanner(removable).Scan(context.Background(), []string{trusted, removable}, fabricar.ScanOptions{})
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]fabricar.DiscoveredPiece{}
	for _, p := range res.Pieces {
		byName[p.Name] = p
	}
	for _, n := range []string{"alfa", "beta", "gama"} {
		if _, ok := byName[n]; !ok {
			t.Fatalf("peça %q não encontrada: %+v", n, res.Pieces)
		}
	}
	if byName["alfa"].Untrusted || !byName["beta"].Untrusted || !byName["gama"].Untrusted {
		t.Errorf("Untrusted incorreto: %+v", byName)
	}
	if byName["beta"].Drive != removable || byName["beta"].Kind != fabricar.KindSkill {
		t.Errorf("drive/kind incorretos: %+v", byName["beta"])
	}
}

func TestScanStorageLocations_HashesAndScoresWithReputation(t *testing.T) {
	root := t.TempDir()
	path := skillAt(t, filepath.Join(root, "skills"), "alfa")
	res, err := scanner().Scan(context.Background(), []string{root}, fabricar.ScanOptions{})
	if err != nil || len(res.Pieces) != 1 {
		t.Fatalf("esperava 1 peça: %v %+v", err, res.Pieces)
	}
	raw, _ := os.ReadFile(path)
	sum := sha256.Sum256(raw)
	p := res.Pieces[0]
	if p.SHA256 != hex.EncodeToString(sum[:]) {
		t.Errorf("SHA256 incorreto: %s", p.SHA256)
	}
	if p.Score <= 0 || p.Score > 1 || p.Tier == "" {
		t.Errorf("score/tier devem vir de internal/reputation: %+v", p)
	}
}

func TestScanStorageLocations_PrunesNodeModulesAndGit(t *testing.T) {
	root := t.TempDir()
	skills := filepath.Join(root, "skills")
	skillAt(t, skills, "legitima")
	skillAt(t, filepath.Join(skills, "node_modules"), "veneno")
	skillAt(t, filepath.Join(skills, ".git"), "oculta")
	res, err := scanner().Scan(context.Background(), []string{root}, fabricar.ScanOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range res.Pieces {
		if p.Name == "veneno" || p.Name == "oculta" {
			t.Errorf("a poda falhou: %s", p.Path)
		}
	}
	if len(res.Pieces) != 1 {
		t.Errorf("esperava só a peça legítima: %+v", res.Pieces)
	}
}

func TestScanStorageLocations_SymlinkLoopTerminates(t *testing.T) {
	root := t.TempDir()
	skills := filepath.Join(root, "skills")
	skillAt(t, skills, "alfa")
	makeDirLink(t, skills, filepath.Join(skills, "alfa", "laco"))
	done := make(chan struct{})
	go func() {
		_, _ = scanner().Scan(context.Background(), []string{root}, fabricar.ScanOptions{})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("o laço de symlink travou a varredura")
	}
}

func TestScanStorageLocations_SlowRootIsInterruptedAndRecorded(t *testing.T) {
	s := scanner()
	s.ScanRoot = func(ctx context.Context, root string) ([]assets.Item, error) {
		if root == "lenta" {
			<-ctx.Done()
			return nil, ctx.Err()
		}
		return []assets.Item{{Kind: assets.KindSkill, Name: "rapida", Path: filepath.Join(root, "x")}}, nil
	}
	start := time.Now()
	res, err := s.Scan(context.Background(), []string{"lenta", "rapida"}, fabricar.ScanOptions{PerRootTimeout: 150 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if time.Since(start) > 3*time.Second {
		t.Fatal("a raiz lenta atrasou a varredura")
	}
	if len(res.Skipped) != 1 || res.Skipped[0].Root != "lenta" || res.Skipped[0].Reason != "timeout" {
		t.Errorf("a raiz lenta deve ficar em Skipped com timeout: %+v", res.Skipped)
	}
	if len(res.Pieces) != 1 || res.Pieces[0].Name != "rapida" {
		t.Errorf("as demais raízes devem voltar: %+v", res.Pieces)
	}
}

func TestScanStorageLocations_FailingRootIsRecordedNotFatal(t *testing.T) {
	s := scanner()
	s.ScanRoot = func(context.Context, string) ([]assets.Item, error) { return nil, errors.New("negado") }
	res, err := s.Scan(context.Background(), []string{"X:\\"}, fabricar.ScanOptions{})
	if err != nil || len(res.Skipped) != 1 || res.Skipped[0].Reason != "inacessivel" {
		t.Fatalf("raiz com erro deve ser registrada: %v %+v", err, res.Skipped)
	}
}

func TestScanStorageLocations_OrderIsScoreDescThenPathAsc(t *testing.T) {
	root := t.TempDir()
	skills := filepath.Join(root, "skills")
	skillAt(t, skills, "b-boa")
	skillAt(t, skills, "a-boa")
	writeFile(t, filepath.Join(skills, "c-fraca", "SKILL.md"), "---\nname: c-fraca\ndescription: x\n---\ntexto solto\n")
	res, err := scanner().Scan(context.Background(), []string{root}, fabricar.ScanOptions{})
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, p := range res.Pieces {
		names = append(names, p.Name)
	}
	if !slices.Equal(names, []string{"a-boa", "b-boa", "c-fraca"}) {
		t.Errorf("ordem = %v; quer Score desc, Path asc", names)
	}
}

func TestScanStorageLocations_LimitsBecomeWarningsNotErrors(t *testing.T) {
	root := t.TempDir()
	for _, n := range []string{"a", "b", "c", "d", "e"} {
		skillAt(t, filepath.Join(root, "skills"), n)
	}
	res, err := scanner().Scan(context.Background(), []string{root}, fabricar.ScanOptions{MaxPieces: 2})
	if err != nil || len(res.Pieces) != 2 {
		t.Fatalf("MaxPieces deve limitar sem erro: %v %d", err, len(res.Pieces))
	}
	if len(res.Warnings) == 0 || !strings.Contains(res.Warnings[0], "limite") {
		t.Errorf("o limite deve gerar aviso: %v", res.Warnings)
	}
	res, err = scanner().Scan(context.Background(), []string{root}, fabricar.ScanOptions{MaxFiles: 3})
	if err != nil || len(res.Warnings) == 0 {
		t.Errorf("MaxFiles deve gerar aviso sem erro: %v %v", err, res.Warnings)
	}
}

func TestScanStorageLocations_PackageLevelEntryPoint(t *testing.T) {
	root := t.TempDir()
	skillAt(t, filepath.Join(root, "skills"), "alfa")
	res, err := fabricar.ScanStorageLocations(context.Background(), []string{root}, fabricar.ScanOptions{})
	if err != nil || len(res.Pieces) != 1 {
		t.Fatalf("ScanStorageLocations: %v %+v", err, res.Pieces)
	}
}

// makeDirLink cria um laço de diretório: symlink e, no Windows sem privilégio, junction.
// Só pula o teste se nenhuma das duas formas for possível.
func makeDirLink(t *testing.T, target, link string) {
	t.Helper()
	err := os.Symlink(target, link)
	if err == nil {
		return
	}
	if runtime.GOOS == "windows" {
		if out, jerr := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); jerr == nil {
			return
		} else {
			t.Skipf("nem symlink nem junction: %v / %v %s", err, jerr, out)
		}
	}
	t.Skipf("symlink não permitido neste ambiente: %v", err)
}
