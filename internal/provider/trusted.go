package provider

// TrustedRouterKey devolve a chave do 9Router só quando routerURL é a base confiável (a primeira configurada).
// Qualquer outra URL recebe "": a credencial nunca trafega para uma base que o operador não designou.
func TrustedRouterKey(routerURL string) string {
	cfg := DefaultConfig()
	if trusted := cfg.TrustedRouterURL(); trusted != "" && routerURL == trusted {
		return cfg.RouterKey
	}
	return ""
}
