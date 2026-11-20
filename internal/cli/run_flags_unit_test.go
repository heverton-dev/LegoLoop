package cli

import (
	"bytes"
	"strings"
	"testing"
)

// A flag de consentimento do unit_no tem o nome canônico e mantém o alias legado, sem quebrar scripts existentes.
func TestRunFlagConsentimentoUnitNoENomeLegado(t *testing.T) {
	for _, nome := range []string{"allow-unit-no", "allow-unit-tasks"} {
		f, err := parseRunFlags("run", []string{"--" + nome, "x.json"}, Deps{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}})
		if err != nil {
			t.Fatalf("--%s: %v", nome, err)
		}
		if !f.allowUnitNo {
			t.Errorf("--%s deve ligar o consentimento do unit_no", nome)
		}
	}
	f, err := parseRunFlags("run", []string{"x.json"}, Deps{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}})
	if err != nil || f.allowUnitNo {
		t.Errorf("sem a flag o consentimento fica desligado: %v %+v", err, f)
	}
}

func TestRunFlagAjudaCitaUnitNoEMarcaOAliasLegado(t *testing.T) {
	var out bytes.Buffer
	_, _ = parseRunFlags("run", []string{"-h"}, Deps{Stdout: &out, Stderr: &out})
	ajuda := out.String()
	if !strings.Contains(ajuda, "allow-unit-no") || !strings.Contains(ajuda, "permite nós unit_no") {
		t.Errorf("a ajuda deve citar allow-unit-no e unit_no: %q", ajuda)
	}
	if !strings.Contains(ajuda, "alias legado") {
		t.Errorf("a ajuda deve marcar o alias legado: %q", ajuda)
	}
}
