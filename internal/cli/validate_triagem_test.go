package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

const noScript = `{"id":"%s","type":"script","command":["go","version"],"timeout_sec":60,
 "routing":{"mode":"manual","model":"none"},
 "output":{"artifact":"%s.txt","gate":"exit_code == 0"}}`

func dagComNos(ids ...string) string {
	var nos []string
	for _, id := range ids {
		nos = append(nos, strings.NewReplacer("%s", id).Replace(noScript))
	}
	return `{"name":"legoloop-dois","nodes":[` + strings.Join(nos, ",") + `]}`
}

func TestValidateExigeTriagemComMaisDeUmNo(t *testing.T) {
	pasta := filepath.Join(t.TempDir(), "legoloop-dois")
	if err := os.MkdirAll(pasta, 0o755); err != nil {
		t.Fatal(err)
	}
	p := caminhos.CaminhoDAG(pasta)
	if err := os.WriteFile(p, []byte(dagComNos("a", "b")), 0o644); err != nil {
		t.Fatal(err)
	}
	d, _, errb := deps(t)
	if code := run(t, d, "validate", p); code == 0 || !strings.Contains(errb.String(), "triagem.json") {
		t.Fatalf("sem triagem deve falhar: code=%d err=%s", code, errb.String())
	}
	incompleta, _ := json.Marshal(triagem.GerarTriagemPadrao("legoloop-dois", []string{"a"}))
	os.WriteFile(filepath.Join(pasta, triagem.ArquivoTriagem), incompleta, 0o644)
	d, _, errb = deps(t)
	if code := run(t, d, "validate", p); code == 0 || !strings.Contains(errb.String(), "não cobre") {
		t.Fatalf("triagem sem o nó b deve falhar: code=%d err=%s", code, errb.String())
	}
	completa, _ := json.Marshal(triagem.GerarTriagemPadrao("legoloop-dois", []string{"a", "b"}))
	os.WriteFile(filepath.Join(pasta, triagem.ArquivoTriagem), completa, 0o644)
	d, _, errb = deps(t)
	if code := run(t, d, "validate", p); code != 0 {
		t.Fatalf("triagem completa deve passar: code=%d err=%s", code, errb.String())
	}
}

func TestValidateNaoExigeTriagemDeUmNoNemDeOutroArquivo(t *testing.T) {
	pasta := filepath.Join(t.TempDir(), "legoloop-dois")
	os.MkdirAll(pasta, 0o755)
	um := caminhos.CaminhoDAG(pasta)
	os.WriteFile(um, []byte(dagComNos("a")), 0o644)
	solto := filepath.Join(pasta, "espelho.dag.json")
	os.WriteFile(solto, []byte(dagComNos("a", "b")), 0o644)
	for _, p := range []string{um, solto} {
		d, _, errb := deps(t)
		if code := run(t, d, "validate", p); code != 0 {
			t.Errorf("%s: code=%d err=%s", filepath.Base(p), code, errb.String())
		}
	}
}
