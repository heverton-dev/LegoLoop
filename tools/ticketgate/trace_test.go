package main

import (
	"strings"
	"testing"
)

func TestTraceStepListaOrfaos(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "spec.md", "* [ ] `AC-X-1.1`: a\n* [ ] `AC-X-1.2`: b\n")
	write(t, dir, "dag.json", `{"nodes":[{"id":"N1","spec_ref":"AC-X-1.1"},{"id":"N2"}]}`)
	out, err := traceStep(config{Root: dir, Trace: []string{"."}}).Run()
	if err == nil {
		t.Fatal("esperava falha com órfãos")
	}
	for _, id := range []string{"AC-X-1.2", "N2"} {
		if !strings.Contains(out, id) {
			t.Fatalf("saída %q não cita %s", out, id)
		}
	}
}

func TestTraceStepCadeiaIntegra(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "spec.md", "* [ ] `AC-X-1.1`: a\n")
	write(t, dir, "dag.json", `{"nodes":[{"id":"N1","spec_ref":"AC-X-1.1"}]}`)
	if out, err := traceStep(config{Root: dir, Trace: []string{"."}}).Run(); err != nil {
		t.Fatalf("cadeia íntegra falhou: %v %s", err, out)
	}
}

func TestTraceStepCaminhoInexistente(t *testing.T) {
	if _, err := traceStep(config{Root: t.TempDir(), Trace: []string{"nada"}}).Run(); err == nil {
		t.Fatal("esperava erro para caminho inexistente")
	}
}

func TestParseArgsTrace(t *testing.T) {
	cfg, err := parseArgs([]string{"-trace", "a,b"}, &strings.Builder{})
	if err != nil || len(cfg.Trace) != 2 {
		t.Fatalf("cfg=%+v err=%v", cfg, err)
	}
}
