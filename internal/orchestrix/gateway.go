package orchestrix

import (
	"errors"
	"fmt"
	"strings"

	"legoloop/internal/dag"
)

// viaRouter é o valor de routing.via que manda o filho do claude pelo 9Router.
const viaRouter = "9router"

// noViaRouter informa se o nó roda o claude CLI pelo gateway do 9Router.
func noViaRouter(n dag.Node) bool {
	return n.Routing.Via == viaRouter && !unitOutroHarness(n)
}

// envGateway monta as variáveis do gateway só para o filho de um nó via 9router.
// Nó via oficial devolve nil. O ambiente do processo pai e a configuração global do gw não mudam.
// A chave só sai para a base confiável (a mesma regra do callRouter).
func (o *Orchestrix) envGateway(n dag.Node) ([]string, error) {
	if !noViaRouter(n) {
		return nil, nil
	}
	key := o.routerCredential()
	if key == "" || o.Inventory.RouterURL == "" {
		return nil, fmt.Errorf("nó %q pede via 9router, mas não há base confiável com chave", n.ID)
	}
	return []string{"ANTHROPIC_BASE_URL=" + o.Inventory.RouterURL, "ANTHROPIC_AUTH_TOKEN=" + key}, nil
}

// semChave remove a chave do 9Router do texto de um erro antes de ele chegar a log ou a artefato.
func (o *Orchestrix) semChave(err error) error {
	if err == nil || o.RouterKey == "" || !strings.Contains(err.Error(), o.RouterKey) {
		return err
	}
	return errors.New(o.semChaveTexto(err.Error()))
}

// semChaveTexto remove a chave do 9Router de um texto.
func (o *Orchestrix) semChaveTexto(t string) string {
	if o.RouterKey == "" {
		return t
	}
	return strings.ReplaceAll(t, o.RouterKey, "[chave omitida]")
}
