package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
)

func criarProjetoAtivosCLI(t *testing.T, estilo string) string {
	t.Helper()
	dir := t.TempDir()
	proj := roteiro.Projeto{
		SchemaVersion: 1, Slug: "cli-teste", Canal: "c", Nicho: "tech",
		PerfisFonte: []string{"geral"}, EstruturaNarrativa: "causal",
		Formatos: []string{"16:9"}, DuracaoAlvoS: 60, Voz: "pt-BR-ThalitaMultilingualNeural",
		EstiloVisual: estilo, Animacao: "diorama", PoliticaGates: roteiro.PoliticaGates{Modo: "mecanico"},
	}
	marca := roteiro.Marca{
		SchemaVersion: 1, Nome: "Canal", Formatos: []string{"16:9"},
		Voz: roteiro.VozMarca{Motor: "edge-tts", Voz: "pt-BR-ThalitaMultilingualNeural", Velocidade: "+0%"},
	}
	script := roteiro.Script{
		SchemaVersion: 1, Title: "Teste CLI", Idioma: "pt-BR", Estrutura: "causal", TipoPauta: "perene",
		Lines: []roteiro.LinhaRoteiro{
			{Text: "L1", Natureza: "narrativa_transicao"},
			{Text: "L2", Natureza: "narrativa_transicao"},
		},
	}
	for nome, val := range map[string]any{"projeto.json": proj, "marca.json": marca, "script.json": script} {
		b, _ := json.MarshalIndent(val, "", "  ")
		if err := os.WriteFile(filepath.Join(dir, nome), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "plano.md"), []byte("# Plano\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestDioramaAtivosCLI(t *testing.T) {
	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf}

	if code := cmdDioramaAtivos(context.Background(), []string{"--help"}, d); code != 0 || !strings.Contains(out.String(), "diorama-ativos") {
		t.Fatalf("help falhou: code=%d out=%s", code, out.String())
	}
	if code := cmdDioramaAtivos(context.Background(), []string{"desconhecido"}, d); code != 2 {
		t.Fatalf("esperado exit 2 para subcomando invalido, obtido %d", code)
	}

	dir := criarProjetoAtivosCLI(t, "codigo")
	out.Reset()
	if code := cmdDioramaAtivos(context.Background(), []string{"prevoo", "--projeto", dir, "--json"}, d); code != 0 {
		t.Fatalf("prevoo falhou: code=%d out=%s", code, out.String())
	}

	out.Reset()
	if code := cmdDioramaAtivos(context.Background(), []string{"executar", "--projeto", dir, "--aprovar", "etapa3_elenco", "--json"}, d); code != 0 {
		t.Fatalf("executar falhou: code=%d out=%s", code, out.String())
	}

	out.Reset()
	if code := cmdDioramaAtivos(context.Background(), []string{"validar-storyboard", "--projeto", dir, "--json"}, d); code != 0 {
		t.Fatalf("validar-storyboard falhou: code=%d out=%s", code, out.String())
	}

	dirFlowSemTrava := criarProjetoAtivosCLI(t, "flow")
	if code := cmdDioramaAtivos(context.Background(), []string{"prevoo", "--projeto", dirFlowSemTrava, "--json"}, d); code != 1 {
		t.Fatalf("esperado exit 1 para flow sem trava, obtido %d", code)
	}
	if code := cmdDioramaAtivos(context.Background(), []string{"executar", "--projeto", dirFlowSemTrava, "--json"}, d); code != 1 {
		t.Fatalf("esperado exit 1 ao executar flow sem trava, obtido %d", code)
	}
	if code := cmdDioramaAtivos(context.Background(), []string{"validar-storyboard", "--projeto", t.TempDir(), "--json"}, d); code != 1 {
		t.Fatalf("esperado exit 1 para validar-storyboard em pasta vazia, obtido %d", code)
	}

	dirHibrido := criarProjetoAtivosCLI(t, "hibrido")
	rawProj, _ := os.ReadFile(filepath.Join(dirHibrido, "projeto.json"))
	var pj roteiro.Projeto
	_ = json.Unmarshal(rawProj, &pj)
	pj.FlowProjeto = "uuid-hibrido-cli"
	b, _ := json.MarshalIndent(pj, "", "  ")
	_ = os.WriteFile(filepath.Join(dirHibrido, "projeto.json"), b, 0o644)
	if code := cmdDioramaAtivos(context.Background(), []string{"executar", "--projeto", dirHibrido, "--json"}, d); code != 1 {
		t.Fatalf("esperado exit 1 para executar modo hibrido sem Chrome CDP 9333 aberto, obtido %d", code)
	}
}
