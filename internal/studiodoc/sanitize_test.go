package studiodoc

import (
	"testing"

	"legoloop/internal/dag"
)

func TestToKebabCase(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"Meu Novo Loop 2026", "meu-novo-loop-2026"},
		{"Criar_Código & Testes! V1", "criar-codigo-testes-v1"},
		{"  --loop---teste--  ", "loop-teste"},
		{"Kebab-Ja-Pronto", "kebab-ja-pronto"},
		{"AÇÃO/REAÇÃO", "acao-reacao"},
		{"", "unnamed-loop"},
	}

	for _, tc := range cases {
		got := ToKebabCase(tc.input)
		if got != tc.expected {
			t.Errorf("ToKebabCase(%q) = %q, esperava %q", tc.input, got, tc.expected)
		}
	}
}

func TestStoreSaveNormalizesToKebab(t *testing.T) {
	tmp := t.TempDir()
	st, err := NewStore(tmp)
	if err != nil {
		t.Fatal(err)
	}

	rawName := "Meu Fluxo Especial V1"
	lp := &Processo{
		ID:   ToKebabCase(rawName),
		Name: rawName,
		Nodes: []dag.Node{
			{
				ID:   "n1",
				Type: dag.TypeSkill,
				Ref:  "interview",
				Routing: dag.RoutingConfig{
					Mode:  dag.ModeManual,
					Model: "local",
				},
				Output: dag.HandoffContract{
					Artifact: "out.txt",
					Gate:     "file_exists('out.txt')",
				},
			},
		},
	}

	if err := st.Save(lp); err != nil {
		t.Fatal(err)
	}

	if lp.ID != "meu-fluxo-especial-v1" {
		t.Errorf("lp.ID não foi normalizado para kebab-case: %q", lp.ID)
	}

	loaded, err := st.Get("meu-fluxo-especial-v1")
	if err != nil {
		t.Fatalf("não conseguiu carregar loop pelo id kebab: %v", err)
	}
	if loaded.ID != "meu-fluxo-especial-v1" {
		t.Errorf("loaded.ID = %q", loaded.ID)
	}
}
