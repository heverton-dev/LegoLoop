package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/forja"
)

func criarSkillTeste(t *testing.T, dirSkill, nome, descricao string) {
	t.Helper()
	if err := os.MkdirAll(dirSkill, 0o755); err != nil {
		t.Fatal(err)
	}
	conteudo := "---\nname: " + nome + "\ndescription: " + descricao + "\n---\n# Corpo\n"
	if err := os.WriteFile(filepath.Join(dirSkill, "SKILL.md"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCmdForjarVerificarEstourado(t *testing.T) {
	home := t.TempDir()
	criarSkillTeste(t, filepath.Join(home, ".claude", "skills", "skill-pesada"), "skill-pesada", strings.Repeat("descricao longa ", 20))

	arqOrc := filepath.Join(t.TempDir(), "orcamento.json")
	if err := os.WriteFile(arqOrc, []byte(`{"versao":"1.0.0","tetos":{"claude":10}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf, Home: home}
	code := Run(context.Background(), []string{"forjar", "--verificar", "--orcamento", arqOrc, "--home", home}, d)
	if code != 1 {
		t.Fatalf("esperava exit 1 para orcamento estourado, obteve %d (stdout=%q stderr=%q)", code, out.String(), errBuf.String())
	}
	combinado := out.String() + errBuf.String()
	if !strings.Contains(combinado, "claude") || !strings.Contains(combinado, "excesso") {
		t.Fatalf("saida deve citar harness e excesso, obteve: %s", combinado)
	}
}

func TestCmdForjarVerificarRespeitado(t *testing.T) {
	home := t.TempDir()
	criarSkillTeste(t, filepath.Join(home, ".claude", "skills", "legoloop-skills"), "legoloop-skills", "indice lazy")

	arqOrc := filepath.Join(t.TempDir(), "orcamento.json")
	if err := os.WriteFile(arqOrc, []byte(`{"versao":"1.0.0","tetos":{"claude":500}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf, Home: home}
	code := Run(context.Background(), []string{"forjar", "--verificar", "--orcamento", arqOrc, "--home", home, "--json"}, d)
	if code != 0 {
		t.Fatalf("esperava exit 0 para orcamento respeitado, obteve %d (stderr=%q)", code, errBuf.String())
	}
	var res forja.ResultadoOrcamento
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatalf("JSON invalido: %v\nsaida: %s", err, out.String())
	}
	if !res.Aprovado || len(res.Excessos) != 0 {
		t.Fatalf("esperava aprovado sem excessos, obteve %+v", res)
	}
}

func TestCmdForjarDryRunApplyEPerfilInvalido(t *testing.T) {
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf, Home: home}

	code := Run(context.Background(), []string{"forjar", "--dry-run", "--home", home, "--sem", "opencode", "--json"}, d)
	if code != 0 {
		t.Fatalf("forjar --dry-run falhou (code=%d): %s", code, errBuf.String())
	}
	var rel forja.Relatorio
	if err := json.Unmarshal(out.Bytes(), &rel); err != nil {
		t.Fatalf("JSON invalido: %v", err)
	}
	if rel.Veredito != forja.VereditoOK {
		t.Fatalf("veredito inesperado: %s", rel.Veredito)
	}

	out.Reset()
	errBuf.Reset()
	code = Run(context.Background(), []string{"forjar", "--perfil", "xyz", "--home", home}, d)
	if code != 2 {
		t.Fatalf("esperava exit 2 para perfil invalido, obteve %d (%s)", code, errBuf.String())
	}
}
