package fechar

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDestilarELimpar_SemStopGateNaoLimpa(t *testing.T) {
	dir := t.TempDir()
	slug := "meu-processo"
	epPath := filepath.Join(dir, ".legoloop", "episodic", slug+".json")
	if err := os.MkdirAll(filepath.Dir(epPath), 0o755); err != nil {
		t.Fatal(err)
	}
	eventos := []ErroEpisodico{
		{Assinatura: "ERR_LINT_FMT", Mensagem: "gofmt falhou"},
		{Assinatura: "ERR_LINT_FMT", Mensagem: "gofmt falhou de novo"},
	}
	data, _ := json.Marshal(eventos)
	if err := os.WriteFile(epPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := DestilarELimpar(dir, slug, false, 2, false)
	if err == nil {
		t.Fatalf("esperava erro quando stopGateAprovado=false, obteve nil (res=%+v)", res)
	}
	if _, statErr := os.Stat(epPath); os.IsNotExist(statErr) {
		t.Fatalf("memoria episodica nao deveria ser removida sem Stop-Gate aprovado")
	}
}

func TestDestilarELimpar_DestilaRecorrentesELimpaEpisodica(t *testing.T) {
	dir := t.TempDir()
	slug := "meu-processo"
	epPath := filepath.Join(dir, ".legoloop", "episodic", slug+".json")
	if err := os.MkdirAll(filepath.Dir(epPath), 0o755); err != nil {
		t.Fatal(err)
	}
	eventos := []ErroEpisodico{
		{Assinatura: "ERR_LINT_FMT", Mensagem: "gofmt falhou"},
		{Assinatura: "ERR_LINT_FMT", Mensagem: "gofmt falhou novamente"},
		{Assinatura: "ERR_UNICO", Mensagem: "erro pontual"},
	}
	data, _ := json.Marshal(eventos)
	if err := os.WriteFile(epPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := DestilarELimpar(dir, slug, true, 2, false)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(res.GatesDestilados) != 1 {
		t.Fatalf("esperava 1 gate destilado, obteve %d", len(res.GatesDestilados))
	}
	if res.GatesDestilados[0].ErroAssinatura != "ERR_LINT_FMT" {
		t.Fatalf("esperava assinatura ERR_LINT_FMT, obteve %s", res.GatesDestilados[0].ErroAssinatura)
	}
	if res.GatesDestilados[0].Ocorrencias != 2 {
		t.Fatalf("esperava 2 ocorrencias, obteve %d", res.GatesDestilados[0].Ocorrencias)
	}
	if !res.MemoriaLimpa {
		t.Fatalf("esperava MemoriaLimpa=true")
	}
	if _, statErr := os.Stat(epPath); !os.IsNotExist(statErr) {
		t.Fatalf("arquivo de memoria episodica deveria ter sido removido")
	}
	gatePath := filepath.Join(dir, ".legoloop", "gates", slug+"_destilados.json")
	if _, statErr := os.Stat(gatePath); os.IsNotExist(statErr) {
		t.Fatalf("arquivo de gates destilados deveria existir em %s", gatePath)
	}
}

func TestDestilarELimpar_DryRunNaoAlteraDisco(t *testing.T) {
	dir := t.TempDir()
	slug := "meu-processo"
	epPath := filepath.Join(dir, ".legoloop", "episodic", slug+".json")
	if err := os.MkdirAll(filepath.Dir(epPath), 0o755); err != nil {
		t.Fatal(err)
	}
	eventos := []ErroEpisodico{
		{Assinatura: "ERR_VET", Mensagem: "go vet falhou"},
		{Assinatura: "ERR_VET", Mensagem: "go vet falhou"},
	}
	data, _ := json.Marshal(eventos)
	if err := os.WriteFile(epPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := DestilarELimpar(dir, slug, true, 2, true)
	if err != nil {
		t.Fatalf("erro inesperado em dry-run: %v", err)
	}
	if len(res.GatesDestilados) != 1 {
		t.Fatalf("esperava 1 gate destilado em dry-run, obteve %d", len(res.GatesDestilados))
	}
	if res.MemoriaLimpa {
		t.Fatalf("em dry-run MemoriaLimpa deve ser false")
	}
	if _, statErr := os.Stat(epPath); os.IsNotExist(statErr) {
		t.Fatalf("em dry-run arquivo episodico nao deve ser removido")
	}
}
