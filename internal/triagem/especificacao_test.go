package triagem

import (
	"context"
	"path/filepath"
	"testing"
)

type resolvedorQueFalha struct {
	t *testing.T
}

func (r *resolvedorQueFalha) Resolver(_ context.Context, _ string, _ []TriagemSpec) (Decisao, error) {
	r.t.Fatal("LLM chamado")
	return Decisao{}, nil
}

func TestTriagemRoteiaEspecificacao(t *testing.T) {
	arquivos := []string{
		filepath.Join("..", "..", "processos", "legoloop-especificacao", "triagem.json"),
		filepath.Join("..", "..", "processos", "legoloop-concepcao", "triagem.json"),
		filepath.Join("..", "..", "processos", "legoloop-planejamento", "triagem.json"),
	}

	var catalogo []TriagemSpec
	for _, arq := range arquivos {
		spec, ok, err := LerSpec(arq)
		if err != nil {
			t.Fatalf("erro ao ler %s: %v", arq, err)
		}
		if !ok {
			t.Fatalf("arquivo %s não encontrado", arq)
		}
		catalogo = append(catalogo, spec)
	}

	resolvedor := &resolvedorQueFalha{t: t}
	entrada := Entrada{
		Texto: "gerar a SPEC da intent aprovada",
	}

	dec, err := Decidir(context.Background(), entrada, catalogo, resolvedor)
	if err != nil {
		t.Fatalf("erro inesperado na triagem: %v", err)
	}

	if dec.Processo != "legoloop-especificacao" {
		t.Errorf("esperava processo %q, obteve %q", "legoloop-especificacao", dec.Processo)
	}
	if dec.Camada != "deterministica" {
		t.Errorf("esperava camada %q, obteve %q", "deterministica", dec.Camada)
	}
}
