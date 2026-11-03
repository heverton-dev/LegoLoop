package provider

import "testing"

// A credencial do 9Router só pode ir para a base configurada como confiável; qualquer outra URL não a recebe.
func TestTrustedRouterKey(t *testing.T) {
	t.Setenv("NINEROUTER_URL", "http://127.0.0.1:34567")
	t.Setenv("NINEROUTER_API_KEY", "segredo")
	t.Setenv("NINEROUTER_KEY", "")
	cases := []struct{ name, url, want string }{
		{"a base confiável recebe a chave", "http://127.0.0.1:34567", "segredo"},
		{"outra base não recebe", "http://127.0.0.1:99999", ""},
		{"base remota alheia não recebe", "https://evil.example", ""},
		{"URL vazia não recebe", "", ""},
	}
	for _, c := range cases {
		if got := TrustedRouterKey(c.url); got != c.want {
			t.Errorf("%s: TrustedRouterKey(%q) = %q, esperado %q", c.name, c.url, got, c.want)
		}
	}
}
