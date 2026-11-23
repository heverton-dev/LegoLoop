package studio

import (
	"net/http"
	"strings"
	"testing"
)

func TestA2UIFabricarHierarquiaAninhada(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	rec := do(t, h, http.MethodGet, "/fabricar", "")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, esperado %d", rec.Code, http.StatusSeeOther)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "/sessao/") {
		t.Fatalf("Location = %q, esperado prefixo /sessao/", loc)
	}
	sessID := strings.TrimPrefix(loc, "/sessao/")

	tree, err := s.A2UIHub().Tree(sessID)
	if err != nil {
		t.Fatalf("erro ao obter árvore: %v", err)
	}
	if len(tree) != 1 {
		t.Fatalf("len(tree) = %d, esperado 1 raiz", len(tree))
	}
	root := tree[0]
	if root.SurfaceID != "proc:fabricar" {
		t.Errorf("raiz = %q, esperado proc:fabricar", root.SurfaceID)
	}
	if len(root.Children) != 2 {
		t.Fatalf("len(children) = %d, esperado 2 filhos", len(root.Children))
	}

	crumb, err := s.A2UIHub().Breadcrumb(sessID, "proc:fabricar:concepcao:stop-gate-01")
	if err != nil {
		t.Fatalf("erro no breadcrumb: %v", err)
	}
	if len(crumb) != 3 {
		t.Fatalf("len(crumb) = %d, esperado 3 níveis", len(crumb))
	}

	if att := s.A2UIHub().PendingAttention(sessID); att != "proc:fabricar:concepcao:stop-gate-01" {
		t.Errorf("PendingAttention = %q, esperado o stop gate", att)
	}
}
