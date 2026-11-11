package assets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanSkills(t *testing.T) {
	cfg, home, project := testCfg(t)
	write(t, filepath.Join(home, ".claude", "skills", "alpha", "SKILL.md"),
		"---\nname: alpha\ndescription: Faz a coisa alpha\n---\n# corpo\n")
	write(t, filepath.Join(home, ".aidd", "canonical", "skills", "beta", "SKILL.md"), "# sem frontmatter\n")
	write(t, filepath.Join(project, ".claude", "skills", "gamma", "SKILL.md"), "---\nname: gamma\n---\n")
	write(t, filepath.Join(home, ".claude", "skills", "plugin", "skills", "delta", "SKILL.md"), "---\nname: delta\n---\n")
	write(t, filepath.Join(home, ".claude", "skills", "alpha", "sub", "SKILL.md"), "---\nname: aninhada\n---\n")

	items := scanKind(t, cfg, KindSkill)
	a := mustFind(t, items, "alpha")
	if a.Description != "Faz a coisa alpha" || a.Scope != "global" || a.Origin != "claude" ||
		a.Path != filepath.Join(home, ".claude", "skills", "alpha", "SKILL.md") {
		t.Errorf("alpha = %+v", a)
	}
	if b := mustFind(t, items, "beta"); b.Description != "" || b.Origin != "aidd" {
		t.Errorf("beta (nome pelo diretório) = %+v", b)
	}
	if g := mustFind(t, items, "gamma"); g.Scope != "project" {
		t.Errorf("gamma = %+v", g)
	}
	mustFind(t, items, "delta")
	if _, ok := find(items, "aninhada"); !ok {
		t.Error("esperava encontrar skill aninhada em varredura recursiva profunda")
	}
}

func TestSkillFrontmatterVariants(t *testing.T) {
	cfg, home, _ := testCfg(t)
	root := filepath.Join(home, ".claude", "skills")
	write(t, filepath.Join(root, "folded", "SKILL.md"), "---\r\nname: folded\r\ndescription: >\r\n  texto dobrado\r\n  segunda linha\r\n---\r\n")
	write(t, filepath.Join(root, "quoted", "SKILL.md"), "\xef\xbb\xbf---\nname: \"quoted\"\ndescription: 'com aspas'\n---\n")
	write(t, filepath.Join(root, "huge", "SKILL.md"), "---\nname: huge\n---\n"+strings.Repeat("x", 2<<20))
	items := scanKind(t, cfg, KindSkill)
	if f := mustFind(t, items, "folded"); f.Description != "texto dobrado segunda linha" {
		t.Errorf("folded = %q", f.Description)
	}
	if q := mustFind(t, items, "quoted"); q.Description != "com aspas" {
		t.Errorf("quoted = %q", q.Description)
	}
	mustFind(t, items, "huge")
}

func TestSkillRootsSkipNodeModulesAndDeepTrees(t *testing.T) {
	cfg, home, _ := testCfg(t)
	root := filepath.Join(home, ".claude", "skills")
	write(t, filepath.Join(root, "node_modules", "x", "SKILL.md"), "---\nname: modulo\n---\n")
	write(t, filepath.Join(root, "a", "b", "c", "d", "e", "f", "g", "SKILL.md"), "---\nname: fundo\n---\n")
	write(t, filepath.Join(root, "ok", "SKILL.md"), "---\nname: ok\n---\n")
	items := scanKind(t, cfg, KindSkill)
	for _, n := range []string{"modulo", "fundo"} {
		if _, ok := find(items, n); ok {
			t.Errorf("%s não deveria ser inventariada", n)
		}
	}
	mustFind(t, items, "ok")
}

func TestDuplicateSkillsAreFlagged(t *testing.T) {
	cfg, home, project := testCfg(t)
	same := "---\nname: dup\ndescription: igual\n---\n"
	write(t, filepath.Join(home, ".claude", "skills", "dup", "SKILL.md"), same)
	write(t, filepath.Join(project, ".claude", "skills", "dup", "SKILL.md"), same)
	write(t, filepath.Join(home, ".claude", "skills", "solo", "SKILL.md"), "---\nname: solo\n---\n")
	items := scanKind(t, cfg, KindSkill)
	n := 0
	for _, it := range items {
		if it.Name == "dup" {
			n++
			if it.Meta["duplicates"] != "2" || it.Meta["sha256"] == "" {
				t.Errorf("dup = %+v", it)
			}
		}
	}
	if n != 2 {
		t.Fatalf("esperadas 2 ocorrências de dup, vieram %d", n)
	}
	if s := mustFind(t, items, "solo"); s.Meta["duplicates"] != "" {
		t.Errorf("solo não é duplicada: %+v", s)
	}
}

// Skills viram symlink depois da canonização: o link deve ser inventariado sem descer nele (sem ciclos).
func TestSkillSymlinkIsInventoriedWithoutLooping(t *testing.T) {
	cfg, home, _ := testCfg(t)
	target := filepath.Join(t.TempDir(), "cofre", "ligada")
	write(t, filepath.Join(target, "SKILL.md"), "---\nname: ligada\n---\n")
	root := filepath.Join(home, ".claude", "skills")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(root, "ligada")); err != nil {
		t.Skipf("sem permissão para criar symlink neste ambiente (Windows sem admin): %v", err)
	}
	if err := os.Symlink(root, filepath.Join(target, "ciclo")); err != nil {
		t.Fatal(err)
	}
	it := mustFind(t, scanKind(t, cfg, KindSkill), "ligada")
	if it.Scope != "global" || !strings.HasSuffix(it.Path, "SKILL.md") {
		t.Errorf("ligada = %+v", it)
	}
}

func TestSkillOntologyInference(t *testing.T) {
	cfg, home, _ := testCfg(t)
	root := filepath.Join(home, ".claude", "skills")

	// 1. Skill com frontmatter explícito de contrato ontológico
	write(t, filepath.Join(root, "explicit", "SKILL.md"), `---
name: explicit
description: Skill com contrato explícito
consumes: [Artifact:Briefing]
produces: Artifact:Spec
artifact: SPEC.md
gate: file_exists('SPEC.md')
---
`)

	// 2. Skill canônica sem frontmatter que deve receber inferência heurística
	write(t, filepath.Join(root, "aidd-legoloop-interview", "SKILL.md"), `---
name: aidd-legoloop-interview
description: Entrevista guiada
---
`)

	items := scanKind(t, cfg, KindSkill)

	exp := mustFind(t, items, "explicit")
	if exp.Meta["consumes"] != "Artifact:Briefing" || exp.Meta["produces"] != "Artifact:Spec" ||
		exp.Meta["artifact"] != "SPEC.md" || exp.Meta["gate"] != "file_exists('SPEC.md')" {
		t.Errorf("explicit metadata ontológica incorreta: %+v", exp.Meta)
	}

	inferida := mustFind(t, items, "aidd-legoloop-interview")
	if inferida.Meta["consumes"] == "" || inferida.Meta["produces"] == "" || inferida.Meta["artifact"] == "" || inferida.Meta["gate"] == "" {
		t.Errorf("aidd-legoloop-interview não recebeu inferência ontológica: %+v", inferida.Meta)
	}
}
