package routedag_test

import (
	"github.com/heverton-dev/LegoLoop/internal/routedag"
	"testing"
)

func TestRoutedagBasico(t *testing.T) {
	_, _, err := routedag.Projetar("arquivo_inexistente.json")
	if err == nil {
		t.Fatal("esperava erro ao abrir arquivo inexistente")
	}
}
