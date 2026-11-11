package fabricar_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/studiodoc"
)

var canonicalRe = regexp.MustCompile(`^legoloop-[a-z0-9]+(-[a-z0-9]+)*$`)

func TestNormalizeProcessName_Table(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Processo do Joãozinho", "legoloop-processo-do-joaozinho"},
		{"legoloop_teste", "legoloop-teste"},
		{"MEU PROCESSO", "legoloop-meu-processo"},
		{"legoloop-ja-tem", "legoloop-ja-tem"},
		{"  Ação   Rápida!! ", "legoloop-acao-rapida"},
		{"ÀÉÎÕÜ ç ñ ß", "legoloop-aeiou-c-n-ss"},
		{"../../etc/passwd", "legoloop-etc-passwd"},
		{"CON", "legoloop-con"},
		{strings.Repeat("a", 200), "legoloop-" + strings.Repeat("a", 55)},
	}
	for _, c := range cases {
		got, err := fabricar.NormalizeProcessName(c.in)
		if err != nil || got != c.want {
			t.Errorf("NormalizeProcessName(%q) = %q, %v; quer %q", c.in, got, err, c.want)
		}
	}
}

func TestNormalizeProcessName_Invalid(t *testing.T) {
	for _, in := range []string{"", "   ", "---", "legoloop", "legoloop-", "日本語", "!!!"} {
		got, err := fabricar.NormalizeProcessName(in)
		if !errors.Is(err, fabricar.ErrInvalidName) || got != "" {
			t.Errorf("NormalizeProcessName(%q) = %q, %v; quer ErrInvalidName", in, got, err)
		}
	}
}

func TestNormalizeProcessName_TruncationLeavesNoTrailingHyphen(t *testing.T) {
	in := strings.Repeat("a", 54) + " " + strings.Repeat("b", 20)
	got, err := fabricar.NormalizeProcessName(in)
	if err != nil || len(got) > 64 || strings.HasSuffix(got, "-") {
		t.Fatalf("truncagem incorreta: %q (%d) %v", got, len(got), err)
	}
}

func assertInvariants(t *testing.T, in string) {
	t.Helper()
	got, err := fabricar.NormalizeProcessName(in)
	if err != nil {
		if !errors.Is(err, fabricar.ErrInvalidName) {
			t.Fatalf("erro inesperado para %q: %v", in, err)
		}
		return
	}
	if !canonicalRe.MatchString(got) || len(got) > 64 {
		t.Fatalf("fora da forma canônica: %q -> %q", in, got)
	}
	if err := studiodoc.ValidateID(got); err != nil {
		t.Fatalf("studiodoc.ValidateID rejeitou %q: %v", got, err)
	}
	again, err := fabricar.NormalizeProcessName(got)
	if err != nil || again != got {
		t.Fatalf("não idempotente: %q -> %q -> %q (%v)", in, got, again, err)
	}
}

func TestNormalizeProcessName_Idempotent(t *testing.T) {
	for _, in := range []string{"Meu Processo", "legoloop-x", "ÀÉ", "a--b__c", "  x  ", strings.Repeat("é", 100), "legoloop legoloop-x"} {
		assertInvariants(t, in)
	}
}

func FuzzNormalizeProcessName(f *testing.F) {
	for _, s := range []string{"", "a", "Processo do Joãozinho", "legoloop-", "../x", "日本語", "ß", strings.Repeat("z", 80)} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, in string) { assertInvariants(t, in) })
}
