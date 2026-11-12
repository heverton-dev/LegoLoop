package artefato

import (
	"reflect"
	"strings"
	"testing"
)

func TestFrontMatter(t *testing.T) {
	completo := "---\nid: SPEC-I2\nstatus: APPROVED\ndate: 2026-11-12\nschema_version: 2\nupstream: [INTENT-I2]\ncovers: [R1, R2]\n---\n\n# Titulo\n"
	fm, err := ParseFrontMatter(completo)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if fm.ID != "SPEC-I2" || fm.Status != "APPROVED" || fm.Date != "2026-11-12" || fm.SchemaVersion != "2" {
		t.Errorf("campos escalares incorretos: %+v", fm)
	}
	if !reflect.DeepEqual(fm.Upstream, []string{"INTENT-I2"}) {
		t.Errorf("upstream incorreto: %v", fm.Upstream)
	}
	if !reflect.DeepEqual(fm.Covers, []string{"R1", "R2"}) {
		t.Errorf("covers incorreto: %v", fm.Covers)
	}

	semID := "---\nstatus: APPROVED\ndate: 2026-11-12\n---\n"
	_, err = ParseFrontMatter(semID)
	if err == nil || !strings.Contains(err.Error(), "id") {
		t.Errorf("esperava erro citando id, obtido: %v", err)
	}

	if _, err = ParseFrontMatter("# sem front matter\n"); err == nil {
		t.Error("esperava erro sem delimitador inicial")
	}
	if _, err = ParseFrontMatter("---\nid: X\n"); err == nil {
		t.Error("esperava erro sem delimitador final")
	}
}
