package triagem

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExigida(t *testing.T) {
	if Exigida(0) || Exigida(1) {
		t.Fatal("processo com 0 ou 1 nó não deve exigir triagem")
	}
	if !Exigida(2) || !Exigida(5) {
		t.Fatal("processo com 2+ nós deve exigir triagem")
	}
}

func TestCompletar(t *testing.T) {
	vazia := Completar(TriagemSpec{}, "meu-proc", []string{"01-inicio", "02-fim"})
	if vazia.SchemaVersion != SchemaVersionV1 || len(vazia.Etapas) != 2 {
		t.Fatalf("esperava triagem padrão com 2 etapas, obtive %+v", vazia)
	}
	if vazia.Etapas[0].SlashCommand != "/meu-proc-inicio" {
		t.Fatalf("esperava remoção do prefixo numérico, obtive %q", vazia.Etapas[0].SlashCommand)
	}

	parcial := TriagemSpec{
		SchemaVersion: SchemaVersionV1,
		Processo:      "meu-proc",
		Etapas: []RegraEtapa{
			{ID: "01-inicio", PalavrasChave: []string{"custom"}},
		},
		CriterioAmbiguidade: CriterioAmbiguidade{
			Descricao: "desc",
			Etapas:    []string{"01-inicio"},
		},
	}
	comp := Completar(parcial, "meu-proc", []string{"01-inicio", "02-meio_etapa"})
	if len(comp.Etapas) != 2 {
		t.Fatalf("esperava 2 etapas, obtive %d", len(comp.Etapas))
	}
	if comp.Etapas[0].PalavrasChave[0] != "custom" {
		t.Fatalf("edição existente não foi preservada: %+v", comp.Etapas[0])
	}
	if len(comp.CriterioAmbiguidade.Etapas) != 2 || comp.CriterioAmbiguidade.Etapas[1] != "02-meio_etapa" {
		t.Fatalf("etapas de ambiguidade incompletas: %+v", comp.CriterioAmbiguidade.Etapas)
	}
}

func TestConferir(t *testing.T) {
	dir := t.TempDir()
	if err := Conferir(dir, []string{"a", "b"}); err == nil {
		t.Fatal("esperava erro quando triagem.json não existe")
	}

	arq := filepath.Join(dir, ArquivoTriagem)
	if err := os.WriteFile(arq, []byte("{invalido"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Conferir(dir, []string{"a", "b"}); err == nil || !strings.Contains(err.Error(), "ilegível") {
		t.Fatalf("esperava erro de JSON ilegível, obtive %v", err)
	}

	invalida, _ := json.Marshal(TriagemSpec{SchemaVersion: "errada"})
	if err := os.WriteFile(arq, invalida, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Conferir(dir, []string{"a", "b"}); err == nil || !strings.Contains(err.Error(), "inválida") {
		t.Fatalf("esperava erro de triagem inválida, obtive %v", err)
	}

	incompleta, _ := json.Marshal(GerarTriagemPadrao("proc", []string{"a"}))
	if err := os.WriteFile(arq, incompleta, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Conferir(dir, []string{"a", "b"}); err == nil || !strings.Contains(err.Error(), "não cobre") {
		t.Fatalf("esperava erro de nó faltando, obtive %v", err)
	}

	completa, _ := json.Marshal(GerarTriagemPadrao("proc", []string{"a", "b"}))
	if err := os.WriteFile(arq, completa, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Conferir(dir, []string{"a", "b"}); err != nil {
		t.Fatalf("esperava sucesso, obtive %v", err)
	}
}

type resolvedorErro struct{}

func (resolvedorErro) Resolver(context.Context, string, []TriagemSpec) (Decisao, error) {
	return Decisao{}, errors.New("falha llm")
}

func TestMotorCasosBorda(t *testing.T) {
	spec := GerarTriagemPadrao("demo", nil)
	if primeiraEtapa(spec) != "" {
		t.Fatalf("esperava etapa vazia, obtive %q", primeiraEtapa(spec))
	}
	if _, ok := matchEtapa("", spec); ok {
		t.Fatal("texto vazio não deve casar etapa")
	}
	if casaTexto("", spec.RegrasDeterministicas[0]) {
		t.Fatal("texto vazio não deve casar regra")
	}
	if casaWorkspace("", []string{"a.txt"}) {
		t.Fatal("workspace vazio não deve casar")
	}
	if casaWorkspace(t.TempDir(), []string{"", "nao_existe.txt"}) {
		t.Fatal("artefato vazio/ausente não deve casar")
	}
	if _, err := Decidir(context.Background(), Entrada{Texto: "ambiguo"}, []TriagemSpec{spec}, resolvedorErro{}); err == nil {
		t.Fatal("esperava erro propagado do resolvedor")
	}
}
