package legopress_test

import (
	"strings"
	"testing"

	"legoloop/internal/legopress"
)

func TestCompilerBuildArgsAndDoctor(t *testing.T) {
	manifest := legopress.Manifest{
		Title:     "Obra de Teste",
		Template:  "ativos/livro.typst",
		OutputPDF: "saida.pdf",
		Parts: []legopress.Part{
			{ID: "P01", Title: "Capitulo 1", File: "partes/01.md", Level: 1},
		},
	}

	compiler := legopress.NewCompiler()

	// 1. Validar montagem dos argumentos de compilação (AC-E02.1)
	args := compiler.BuildArgs("concatenado.md", manifest)
	if len(args) == 0 {
		t.Fatal("esperava argumentos para o pandoc, obteve lista vazia")
	}

	argsStr := strings.Join(args, " ")
	if !strings.Contains(argsStr, "--template ativos/livro.typst") {
		t.Errorf("esperava --template ativos/livro.typst nos argumentos, obteve: %s", argsStr)
	}
	if !strings.Contains(argsStr, "--pdf-engine typst") {
		t.Errorf("esperava --pdf-engine typst nos argumentos, obteve: %s", argsStr)
	}
	if !strings.Contains(argsStr, "-o saida.pdf") {
		t.Errorf("esperava -o saida.pdf nos argumentos, obteve: %s", argsStr)
	}

	// 2. Validar verificacao de dependencias Doctor (AC-E02.2)
	report := compiler.Doctor()
	if report == nil {
		t.Fatal("esperava relatorio de dependencias, obteve nil")
	}
	// No ambiente atual do usuário, pandoc e typst existem e estão funcionais
	if !report.PandocAvailable {
		t.Errorf("esperava pandoc disponivel no sistema: %+v", report)
	}
}

func TestCompilerCompile_InvalidManifest(t *testing.T) {
	compiler := legopress.NewCompiler()
	invalidManifest := legopress.Manifest{Title: ""} // sem partes
	err := compiler.Compile(invalidManifest)
	if err != legopress.ErrManifestoInvalido {
		t.Fatalf("esperava ErrManifestoInvalido, obteve: %v", err)
	}
}
