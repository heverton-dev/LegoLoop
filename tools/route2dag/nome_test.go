package main

import (
	"encoding/json"
	"os"
	"testing"
)

// Achado da I02: a ROUTE gerada pelo plangen traz name "<data>_ROUTE_<TEMA>", e o route2dag
// repetia data e a palavra ROUTE no nome do arquivo e no nome do Processo do DAG.
func TestNomeDoDagNaoRepeteDataNemRoute(t *testing.T) {
	raw, err := os.ReadFile("../../docs/routes/2026-11-21_ROUTE_DISTRIBUICAO_FORJA_LEGOLOOP.json")
	if err != nil {
		t.Fatal(err)
	}
	out, nome, err := Converter(raw)
	if err != nil {
		t.Fatal(err)
	}
	if nome != "2026-11-21_DAG_DISTRIBUICAO_FORJA_LEGOLOOP.json" {
		t.Errorf("nome do arquivo = %q", nome)
	}
	var d struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(out, &d); err != nil {
		t.Fatal(err)
	}
	if d.Name != "legoloop-distribuicao-forja-legoloop" {
		t.Errorf("name do DAG = %q", d.Name)
	}
}

func TestNomeLimpoMantemNomeSemPrefixo(t *testing.T) {
	if got := nomeLimpo("NORMALIZACAO_TAXONOMICA_LEI15"); got != "NORMALIZACAO_TAXONOMICA_LEI15" {
		t.Errorf("nome sem prefixo mudou: %q", got)
	}
	if got := nomeLimpo("2026-11-21_ROUTE_X_Y"); got != "X_Y" {
		t.Errorf("nome com prefixo = %q", got)
	}
}
