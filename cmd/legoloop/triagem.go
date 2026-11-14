package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"legoloop/internal/triagem"
)

func runTriagem(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "finalizar" {
		return runTriagemFinalizar(args[1:], stdout, stderr)
	}
	fs := flag.NewFlagSet("triagem", flag.ContinueOnError)
	fs.SetOutput(stderr)
	workspace := fs.String("workspace", "", "diretório de trabalho para verificar artefatos ativos")
	asJSON := fs.Bool("json", false, "emite decisão em formato JSON")
	loopsDir := fs.String("loops-dir", "loops", "diretório raiz dos processos")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	entrada := ""
	if fs.NArg() > 0 {
		entrada = fs.Arg(0)
	}
	if entrada == "" && *workspace == "" {
		fmt.Fprintln(stderr, "uso: legoloop triagem \"<entrada>\" [--workspace DIR] [--json]")
		return 2
	}
	catalogo, err := carregarCatalogoTriagem(*loopsDir)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao carregar catálogo de triagem: %v\n", err)
		return 1
	}
	dec, err := decidirComFallback(ctx, entrada, *workspace, catalogo)
	if err != nil {
		fmt.Fprintf(stderr, "erro na triagem: %v\n", err)
		return 1
	}
	return imprimirDecisaoTriagem(stdout, dec, *asJSON)
}

func decidirComFallback(ctx context.Context, texto, workspace string, catalogo []triagem.TriagemSpec) (triagem.Decisao, error) {
	dec, err := triagem.Decidir(ctx, triagem.Entrada{Texto: texto, Workspace: workspace}, catalogo, nil)
	if errors.Is(err, triagem.ErrResolvedorAusente) {
		return triagem.Decisao{
			PrecisaConfirmacao: true,
			Camada:             "ambiguidade",
			Justificativa:      err.Error(),
		}, nil
	}
	return dec, err
}

func imprimirDecisaoTriagem(stdout io.Writer, dec triagem.Decisao, asJSON bool) int {
	if asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(dec)
		return 0
	}
	fmt.Fprintf(stdout, "Processo: %s\nEtapa inicial: %s\nCamada: %s\nConfiança: %.2f\nJustificativa: %s\n",
		dec.Processo, dec.EtapaInicial, dec.Camada, dec.Confianca, dec.Justificativa)
	return 0
}

func carregarCatalogoTriagem(loopsDir string) ([]triagem.TriagemSpec, error) {
	entries, err := os.ReadDir(loopsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var catalogo []triagem.TriagemSpec
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		spec, ok, err := lerTriagemSpec(filepath.Join(loopsDir, e.Name(), "triagem.json"))
		if err != nil {
			return nil, err
		}
		if ok {
			catalogo = append(catalogo, spec)
		}
	}
	return catalogo, nil
}

func lerTriagemSpec(caminho string) (triagem.TriagemSpec, bool, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		if os.IsNotExist(err) {
			return triagem.TriagemSpec{}, false, nil
		}
		return triagem.TriagemSpec{}, false, err
	}
	var spec triagem.TriagemSpec
	if err := json.Unmarshal(dados, &spec); err != nil {
		return triagem.TriagemSpec{}, false, err
	}
	if err := spec.Validar(); err != nil {
		return triagem.TriagemSpec{}, false, err
	}
	return spec, true, nil
}

func runTriagemFinalizar(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("triagem finalizar", flag.ContinueOnError)
	fs.SetOutput(stderr)
	loopsDir := fs.String("loops-dir", "loops", "diretório raiz dos processos")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(stderr, "uso: legoloop triagem finalizar <slug> [--loops-dir DIR]")
		return 2
	}
	slug := fs.Arg(0)
	if err := triagem.LimparMemoria(filepath.Join(*loopsDir, slug)); err != nil {
		fmt.Fprintf(stderr, "erro ao limpar memória de %s: %v\n", slug, err)
		return 1
	}
	fmt.Fprintf(stdout, "Memória de execução finalizada para %s\n", slug)
	return 0
}
