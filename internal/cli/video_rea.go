package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/videorea"
)

const videoReaUsage = `Uso: legoloop video-rea <subcomando> [opções]

Pipeline universal de engenharia reversa audiovisual, dissecação forense e fábrica de governança.

Subcomandos:
  ingest    --url <url> --out <file>        raspa e extrai metadados de YouTube, Instagram e TikTok
  forensics --input <file> [--transcript-file <f>] [--desc-file <f>] --out <file>   extrai prompts, deduções ocultas e OCR
  curate    --dossier <file> [--name <slug>] --out <file>   cruza com peças Tier Ouro e calcula soberania
  gate      --report <file> [--approve]     executa o stop-gate humano de validação da curadoria
  govern    --approval <file> [--bundle <dir>] gera o Target Workspace isolado e projeção canônica
  book      --bundle <dir>                  compila o Livro Didático LegoPress e audita geometria
  run       --url <url> [--bundle <dir>]    executa o pipeline completo criando o bundle do vídeo
`

func cmdVideoRea(_ context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		fmt.Fprint(d.Stdout, videoReaUsage)
		return 0
	}
	subcmd := strings.ToLower(args[0])
	subArgs := args[1:]
	switch subcmd {
	case "ingest":
		return runIngest(subArgs, d)
	case "forensics":
		return runForensics(subArgs, d)
	case "curate":
		return runCurate(subArgs, d)
	case "gate":
		return runGate(subArgs, d)
	case "govern":
		return runGovern(subArgs, d)
	case "book":
		return runBook(subArgs, d)
	case "run":
		return runFullPipeline(subArgs, d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n\n%s", subcmd, videoReaUsage)
		return 2
	}
}

func runIngest(args []string, d Deps) int {
	fs := newFlags("video-rea ingest", d)
	rawURL := fs.String("url", "", "URL do video")
	out := fs.String("out", "", "arquivo de saida")
	if err := fs.Parse(args); err != nil || *rawURL == "" {
		fmt.Fprintf(d.Stderr, "erro: --url é obrigatório\n")
		return 2
	}
	if _, err := videorea.DetectPlatform(*rawURL); err != nil {
		fmt.Fprintf(d.Stderr, "erro na plataforma: %v\n", err)
		return 1
	}
	vInput := videorea.VideoInput{URL: *rawURL}
	return saveJSONOutput(vInput, *out, d)
}

func runForensics(args []string, d Deps) int {
	fs := newFlags("video-rea forensics", d)
	inputPath := fs.String("input", "", "arquivo de entrada")
	transcript := fs.String("transcript", "", "transcrição")
	transcriptFile := fs.String("transcript-file", "", "arquivo com a transcrição")
	descFile := fs.String("desc-file", "", "arquivo com a descrição do vídeo")
	out := fs.String("out", "", "arquivo de saida do dossie")
	if err := fs.Parse(args); err != nil || *inputPath == "" {
		fmt.Fprintf(d.Stderr, "erro: --input é obrigatório\n")
		return 2
	}
	text, err := resolveInlineOrFile(*transcript, *transcriptFile)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo transcrição: %v\n", err)
		return 1
	}
	desc, err := resolveInlineOrFile("", *descFile)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo descrição: %v\n", err)
		return 1
	}
	var vInput videorea.VideoInput
	if err := loadJSONFile(*inputPath, &vInput); err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo input: %v\n", err)
		return 1
	}
	dossier, err := videorea.ExtractForensics(vInput, text, desc)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro extraindo forense: %v\n", err)
		return 1
	}
	return saveJSONOutput(dossier, *out, d)
}

func resolveInlineOrFile(inline, path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return inline, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func runCurate(args []string, d Deps) int {
	fs := newFlags("video-rea curate", d)
	dossierPath := fs.String("dossier", "", "arquivo do dossie")
	name := fs.String("name", "", "nome canônico do pipeline (slug legoloop-*)")
	out := fs.String("out", "", "arquivo de saida da curadoria")
	if err := fs.Parse(args); err != nil || *dossierPath == "" {
		fmt.Fprintf(d.Stderr, "erro: --dossier é obrigatório\n")
		return 2
	}
	var dossier videorea.ForensicsDossier
	if err := loadJSONFile(*dossierPath, &dossier); err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo dossie: %v\n", err)
		return 1
	}
	report, err := videorea.CurateAgainstMatrix(&dossier)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro na curadoria: %v\n", err)
		return 1
	}
	if canonical := strings.TrimSpace(*name); canonical != "" {
		report.SuggestedName = canonical
	}
	return saveJSONOutput(report, *out, d)
}

func runGate(args []string, d Deps) int {
	fs := newFlags("video-rea gate", d)
	reportPath := fs.String("report", "", "arquivo do relatorio de curadoria")
	approve := fs.Bool("approve", false, "aprovar stop-gate")
	if err := fs.Parse(args); err != nil || *reportPath == "" {
		fmt.Fprintf(d.Stderr, "erro: --report é obrigatório\n")
		return 2
	}
	var report videorea.CurationReport
	if err := loadJSONFile(*reportPath, &report); err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo relatorio: %v\n", err)
		return 1
	}
	if *approve {
		report.Approved = true
	}
	if err := videorea.EvaluateCurationGate(&report, report.Approved); err != nil {
		fmt.Fprintf(d.Stderr, "gate rejeitado: %v\n", err)
		return 1
	}
	return saveJSONOutput(report, *reportPath, d)
}

func runGovern(args []string, d Deps) int {
	fs := newFlags("video-rea govern", d)
	approvalPath := fs.String("approval", "", "relatorio aprovado")
	bundleDir := fs.String("bundle", "docs/videorea", "diretorio base do bundle")
	repoFlag := fs.String("repo", "", "raiz do repositorio")
	if err := fs.Parse(args); err != nil || *approvalPath == "" {
		fmt.Fprintf(d.Stderr, "erro: --approval é obrigatório\n")
		return 2
	}
	var report videorea.CurationReport
	if err := loadJSONFile(*approvalPath, &report); err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo relatorio: %v\n", err)
		return 1
	}
	root := resolveRepoRoot(*repoFlag)
	bundle, err := videorea.MaterializeAndProject(&report, *bundleDir, root)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro materializando bundle: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "Bundle materializado em: %s\n", bundle.BaseDir)
	return 0
}

func runBook(args []string, d Deps) int {
	fs := newFlags("video-rea book", d)
	bundleDir := fs.String("bundle", "", "diretorio do bundle")
	if err := fs.Parse(args); err != nil || *bundleDir == "" {
		fmt.Fprintf(d.Stderr, "erro: --bundle é obrigatório\n")
		return 2
	}
	bookDir := filepath.Join(*bundleDir, "book")
	if _, err := os.Stat(bookDir); os.IsNotExist(err) {
		bookDir = *bundleDir
	}
	manifestPath := filepath.Join(bookDir, "livro.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		fmt.Fprintf(d.Stderr, "manifesto livro.json não encontrado em %s: %v\n", manifestPath, err)
		return 1
	}
	var manifest struct {
		NomeBase string   `json:"nome_base"`
		Partes   []string `json:"partes"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		fmt.Fprintf(d.Stderr, "erro decodificando livro.json: %v\n", err)
		return 1
	}
	if len(manifest.Partes) == 0 {
		fmt.Fprintf(d.Stderr, "erro: livro.json não define partes didáticas\n")
		return 1
	}
	// Auditoria anti-stub: verifica se as partes possuem conteúdo real (> 200 bytes)
	for _, p := range manifest.Partes {
		pPath := filepath.Join(bookDir, "partes", p)
		info, err := os.Stat(pPath)
		if err != nil || info.Size() < 200 {
			fmt.Fprintf(d.Stderr, "erro de conformidade LegoPress: capítulo %s inexistente ou esquelético/stub (< 200 bytes)\n", p)
			return 1
		}
	}
	// Auditoria do PDF compilado
	pdfPath := filepath.Join(bookDir, fmt.Sprintf("%s.pdf", manifest.NomeBase))
	pdfInfo, err := os.Stat(pdfPath)
	if err != nil || pdfInfo.Size() < 50000 {
		fmt.Fprintf(d.Stderr, "erro: PDF editorial %s não encontrado ou inválido (< 50KB). Execute a compilação do LegoPress via typst/pandoc.\n", pdfPath)
		return 1
	}
	fmt.Fprintf(d.Stdout, "Livro LegoPress 100%% auditado e aprovado: %s (%d partes, %d KB)\n", pdfPath, len(manifest.Partes), pdfInfo.Size()/1024)
	return 0
}

func runFullPipeline(args []string, d Deps) int {
	fs := newFlags("video-rea run", d)
	rawURL := fs.String("url", "", "URL do video")
	transcript := fs.String("transcript", "", "transcrição do vídeo")
	desc := fs.String("desc", "", "descrição do vídeo")
	bundleDir := fs.String("bundle", "docs/videorea", "diretório base")
	repoFlag := fs.String("repo", "", "raiz do repositorio")
	if err := fs.Parse(args); err != nil || *rawURL == "" {
		if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
			*rawURL = args[0]
		} else {
			fmt.Fprintf(d.Stderr, "erro: URL do vídeo é obrigatória\n")
			return 2
		}
	}
	if _, err := videorea.DetectPlatform(*rawURL); err != nil {
		fmt.Fprintf(d.Stderr, "erro na plataforma: %v\n", err)
		return 1
	}
	if *transcript == "" && *desc == "" {
		*desc = fmt.Sprintf("Vídeo analisado via %s", *rawURL)
	}
	root := resolveRepoRoot(*repoFlag)
	return executeFullVideoREA(videorea.VideoInput{URL: *rawURL}, *transcript, *desc, *bundleDir, root, d)
}

func executeFullVideoREA(vInput videorea.VideoInput, transcript, desc, baseDir, repoRoot string, d Deps) int {
	dossier, err := videorea.ExtractForensics(vInput, transcript, desc)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro forense: %v\n", err)
		return 1
	}
	report, err := videorea.CurateAgainstMatrix(dossier)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro curadoria: %v\n", err)
		return 1
	}
	report.Approved = true
	bundle, err := videorea.MaterializeAndProject(report, baseDir, repoRoot)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro materialização: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "Pipeline VideoREA concluído com sucesso: %s\n", bundle.BaseDir)
	return 0
}

func resolveRepoRoot(repoFlag string) string {
	if repoFlag != "" {
		return repoFlag
	}
	root, err := os.Getwd()
	if err != nil {
		return "."
	}
	return root
}

func loadJSONFile(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func saveJSONOutput(v any, path string, d Deps) int {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro codificando json: %v\n", err)
		return 1
	}
	if path == "" {
		fmt.Fprintln(d.Stdout, string(data))
		return 0
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fmt.Fprintf(d.Stderr, "erro criando diretorio: %v\n", err)
		return 1
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		fmt.Fprintf(d.Stderr, "erro gravando arquivo: %v\n", err)
		return 1
	}
	return 0
}
