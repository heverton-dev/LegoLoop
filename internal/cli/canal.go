package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/redc"
)

func cmdCanal(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		fmt.Fprintln(d.Stderr, "uso: legoloop canal redc|gerar [opções]")
		return 2
	}
	switch args[0] {
	case "redc":
		return runCanalREDC(args[1:], d)
	case "gerar":
		return runCanalGerar(args[1:], d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando de canal desconhecido: %s\n", args[0])
		return 2
	}
}

func parseArgsAndFlags(rawArgs []string) (string, []string) {
	var pos string
	var flags []string
	for _, a := range rawArgs {
		if !strings.HasPrefix(a, "-") && pos == "" {
			pos = a
		} else {
			flags = append(flags, a)
		}
	}
	return pos, flags
}

func runCanalREDC(rawArgs []string, d Deps) int {
	origem, flagsToParse := parseArgsAndFlags(rawArgs)
	fs := newFlags("canal redc", d)
	canal := fs.String("canal", "", "slug do canal")
	root := fs.String("root", ".", "diretório raiz do projeto")
	asJSON := fs.Bool("json", false, "saída JSON")
	if err := fs.Parse(flagsToParse); err != nil || (origem == "" && fs.NArg() == 0) {
		fmt.Fprintln(d.Stderr, "uso: legoloop canal redc <origem-dir> [--canal <nome>] [--root DIR] [--json]")
		return 2
	}
	if origem == "" {
		origem = fs.Arg(0)
	}
	canalSlug := *canal
	if canalSlug == "" {
		canalSlug = filepath.Base(origem)
	}
	return executeREDC(origem, canalSlug, *root, *asJSON, d)
}

func executeREDC(origem, canalSlug, root string, asJSON bool, d Deps) int {
	videos, clusters, dna, err := loadAndAnalyzeChannel(origem, canalSlug)
	if err != nil {
		fmt.Fprintln(d.Stderr, err.Error())
		return 1
	}
	createdREPC, err := redc.WriteREPC(root, dna)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro ao gravar repc: %v\n", err)
		return 1
	}
	createdPipes, err := redc.GenerateMethodPipelines(root, canalSlug, clusters)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro ao gerar pipelines: %v\n", err)
		return 1
	}
	report := redc.ChannelReport{
		CanalSlug:       canalSlug,
		TotalVideos:     len(videos),
		AmostraDeepREDC: len(dna.TopGanchos),
		MetodosUnicos:   len(clusters),
		DiretorioSaida:  filepath.Join(root, "docs", "channels", canalSlug),
		ArquivosCriados: append(createdREPC, createdPipes...),
	}
	return printCanalReport(report, asJSON, d)
}

func loadAndAnalyzeChannel(origem, canalSlug string) ([]redc.VideoEntry, []redc.MethodCluster, redc.ChannelDNA, error) {
	metricsPath := filepath.Join(origem, "metricas.json")
	raw, err := os.ReadFile(metricsPath)
	if err != nil {
		return nil, nil, redc.ChannelDNA{}, fmt.Errorf("erro ao ler %s: %w", metricsPath, err)
	}
	videos, err := redc.CatalogVideos(raw)
	if err != nil {
		return nil, nil, redc.ChannelDNA{}, fmt.Errorf("erro ao catalogar vídeos: %w", err)
	}
	clusters := redc.ClusterMethods(videos)
	dna, err := redc.ExtractDNA(canalSlug, videos)
	if err != nil {
		return nil, nil, redc.ChannelDNA{}, fmt.Errorf("erro ao extrair DNA: %w", err)
	}
	return videos, clusters, dna, nil
}

func printCanalReport(report redc.ChannelReport, asJSON bool, d Deps) int {
	if asJSON {
		b, _ := json.MarshalIndent(report, "", "  ")
		fmt.Fprintln(d.Stdout, string(b))
		return 0
	}
	fmt.Fprintf(d.Stdout, "Canal %s analisado com sucesso!\n", report.CanalSlug)
	fmt.Fprintf(d.Stdout, "Total de vídeos: %d | Métodos únicos: %d\n", report.TotalVideos, report.MetodosUnicos)
	fmt.Fprintf(d.Stdout, "Diretório de saída: %s\n", report.DiretorioSaida)
	return 0
}

func runCanalGerar(rawArgs []string, d Deps) int {
	canalSlug, flagsToParse := parseArgsAndFlags(rawArgs)
	fs := newFlags("canal gerar", d)
	tema := fs.String("tema", "", "tema do novo vídeo")
	marca := fs.String("nome-marca", "Marca Autoral", "nome da marca")
	corPrim := fs.String("cor-primaria", "#3B82F6", "cor primária")
	root := fs.String("root", ".", "diretório raiz do projeto")
	asJSON := fs.Bool("json", false, "saída JSON")
	if err := fs.Parse(flagsToParse); err != nil || (canalSlug == "" && fs.NArg() == 0) {
		fmt.Fprintln(d.Stderr, "uso: legoloop canal gerar <canal-slug> --tema <tema> [opções]")
		return 2
	}
	if canalSlug == "" {
		canalSlug = fs.Arg(0)
	}
	return executeGerar(canalSlug, *tema, *marca, *corPrim, *root, *asJSON, d)
}

func executeGerar(canalSlug, tema, marca, corPrim, root string, asJSON bool, d Deps) int {
	dna := loadChannelDNAOrFallback(root, canalSlug)
	req := redc.GenerationRequest{
		CanalSlug: canalSlug,
		Tema:      tema,
		Marca: redc.BrandTokens{
			NomeMarca:   marca,
			CorPrimaria: corPrim,
		},
	}
	outPath, err := redc.GenerateStructuralVideo(root, dna, req)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro ao gerar vídeo estrutural: %v\n", err)
		return 1
	}
	if asJSON {
		fmt.Fprintf(d.Stdout, "{\"status\":\"ok\",\"roteiro\":\"%s\"}\n", filepath.ToSlash(outPath))
		return 0
	}
	fmt.Fprintf(d.Stdout, "Vídeo estrutural gerado com sucesso em: %s\n", outPath)
	return 0
}

func loadChannelDNAOrFallback(root, canalSlug string) redc.ChannelDNA {
	dnaPath := filepath.Join(root, "docs", "channels", canalSlug, "repc", "01_dossie_dna.json")
	var dna redc.ChannelDNA
	rawDNA, err := os.ReadFile(dnaPath)
	if err == nil {
		_ = json.Unmarshal(rawDNA, &dna)
		return dna
	}
	return redc.ChannelDNA{
		CanalSlug:   canalSlug,
		Arquetipo:   "Mestre",
		RitmoCortes: "Cortes a cada 2.5s",
	}
}
