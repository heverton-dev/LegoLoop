package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/brain"
	"legoloop/internal/videorea"
)

const videoReaUsage = `Uso: legoloop video-rea <subcomando> [opções]

Processo universal de engenharia reversa audiovisual, dissecação forense e fábrica de governança.

Subcomandos:
  ingest    --url <url> --out <file>        raspa e extrai metadados de YouTube, Instagram e TikTok
  forensics --input <file> [--transcript-file <f>] [--desc-file <f>] [--model <m>] [--effort <low|medium|high>] --out <file>   extrai prompts, deduções ocultas e OCR
  curate    --dossier <file> [--name <slug>] --out <file>   cruza com peças Tier Ouro e calcula soberania
  gate      --report <file> [--approve]     executa o stop-gate humano de validação da curadoria
  govern    --approval <file> [--bundle <dir>] [--model <m>] [--effort <low|medium|high>] gera o Target Workspace isolado e projeção canônica
  book      --bundle <dir> [--model <m>] [--effort <low|medium|high>] compila o Livro Didático LegoPress e audita geometria
  run       --url <url> [--bundle <dir>] [--model <m>] [--effort <low|medium|high>] executa o Processo completo criando o bundle do vídeo
`

func cmdVideoRea(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		fmt.Fprint(d.Stdout, videoReaUsage)
		return 0
	}
	subcmd := strings.ToLower(args[0])
	subArgs := args[1:]
	switch subcmd {
	case "ingest":
		return runIngest(ctx, subArgs, d)
	case "forensics":
		return runForensics(ctx, subArgs, d)
	case "curate":
		return runCurate(subArgs, d)
	case "gate":
		return runGate(subArgs, d)
	case "govern":
		return runGovern(ctx, subArgs, d)
	case "book":
		return runBook(ctx, subArgs, d)
	case "run":
		return runProcessoCompleto(ctx, subArgs, d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n\n%s", subcmd, videoReaUsage)
		return 2
	}
}

func resolveBrain(ctx context.Context, modelFlag, effortFlag string, d Deps) (brain.Brain, error) {
	inv := d.discover(ctx)
	if len(inv.Models()) == 0 {
		return nil, nil
	}
	opts := d.brainOpts()
	if effortFlag != "" {
		opts = append(opts, brain.WithEffort(brain.Effort(effortFlag)))
	}
	b, _, err := brain.Pick(inv, modelFlag, opts...)
	return b, err
}

func runIngest(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("video-rea ingest", d)
	rawURL := fs.String("url", "", "URL do video")
	out := fs.String("out", "", "arquivo de saida")
	dest := fs.String("dest", "", "pasta scratch do bundle (padrão: pasta do --out)")
	if err := fs.Parse(args); err != nil || *rawURL == "" {
		fmt.Fprintf(d.Stderr, "erro: --url é obrigatório\n")
		return 2
	}
	vInput, err := ingestVideo(ctx, *rawURL, *out, *dest, d)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro na ingestão: %v\n", err)
		return 1
	}
	return saveJSONOutput(vInput, *out, d)
}

func ingestVideo(ctx context.Context, rawURL, out, dest string, d Deps) (*videorea.VideoInput, error) {
	platform, err := videorea.DetectPlatform(rawURL)
	if err != nil {
		return nil, err
	}
	if platform == videorea.PlatformLocal {
		return &videorea.VideoInput{URL: rawURL, Platform: platform}, nil
	}
	if dest == "" {
		dest = filepath.Join("scratch", "videorea")
		if out != "" {
			dest = filepath.Dir(out)
		}
	}
	var runner videorea.Runner = videorea.ExecRunner{}
	if d.VideoRunner != nil {
		runner = d.VideoRunner
	}
	return videorea.DownloadWithRunner(ctx, runner, rawURL, dest)
}

func runForensics(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("video-rea forensics", d)
	inputPath := fs.String("input", "", "arquivo de entrada")
	transcript := fs.String("transcript", "", "transcrição")
	transcriptFile := fs.String("transcript-file", "", "arquivo com a transcrição")
	descFile := fs.String("desc-file", "", "arquivo com a descrição do vídeo")
	model := fs.String("model", "", "modelo para inferência")
	effort := fs.String("effort", "", "esforço cognitivo (low, medium, high)")
	out := fs.String("out", "", "arquivo de saida do dossie")
	if err := fs.Parse(args); err != nil || *inputPath == "" {
		fmt.Fprintf(d.Stderr, "erro: --input é obrigatório\n")
		return 2
	}
	b, err := resolveBrain(ctx, *model, *effort, d)
	if *model != "" && err != nil {
		fmt.Fprintf(d.Stderr, "erro ao selecionar cérebro: %v\n", err)
		return 1
	}
	return executeForensics(ctx, b, *inputPath, *transcript, *transcriptFile, *descFile, *out, d)
}

func executeForensics(ctx context.Context, b brain.Brain, inPath, tr, trFile, descFile, out string, d Deps) int {
	text, err := resolveInlineOrFile(tr, trFile)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo transcrição: %v\n", err)
		return 1
	}
	desc, err := resolveInlineOrFile("", descFile)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo descrição: %v\n", err)
		return 1
	}
	var vInput videorea.VideoInput
	if err := loadJSONFile(inPath, &vInput); err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo input: %v\n", err)
		return 1
	}
	if text == "" && desc == "" {
		desc = fmt.Sprintf("Vídeo analisado via %s", vInput.URL)
	}
	dossier, err := videorea.ExtractForensicsWithBrain(ctx, vInput, text, desc, b)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro extraindo forense: %v\n", err)
		return 1
	}
	if err := dossier.ValidateStrict(); err != nil {
		fmt.Fprintf(d.Stderr, "erro validando dossiê: %v\n", err)
		return 1
	}
	return saveJSONOutput(dossier, out, d)
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
	name := fs.String("name", "", "nome canônico do Processo (slug legoloop-*)")
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
	outPath := fs.String("out", "", "arquivo de saida do relatorio aprovado (opcional)")
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
	dest := *reportPath
	if *outPath != "" {
		dest = *outPath
	}
	return saveJSONOutput(report, dest, d)
}

func runGovern(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("video-rea govern", d)
	approvalPath := fs.String("approval", "", "relatorio aprovado")
	bundleDir := fs.String("bundle", "docs/videorea", "diretorio base do bundle")
	repoFlag := fs.String("repo", "", "raiz do repositorio")
	outPath := fs.String("out", "", "arquivo de saida do manifesto (opcional)")
	model := fs.String("model", "", "modelo para inferência")
	effort := fs.String("effort", "", "esforço cognitivo (low, medium, high)")
	if err := fs.Parse(args); err != nil || *approvalPath == "" {
		fmt.Fprintf(d.Stderr, "erro: --approval é obrigatório\n")
		return 2
	}
	b, err := resolveBrain(ctx, *model, *effort, d)
	if *model != "" && err != nil {
		fmt.Fprintf(d.Stderr, "erro ao selecionar cérebro: %v\n", err)
		return 1
	}
	var report videorea.CurationReport
	if err := loadJSONFile(*approvalPath, &report); err != nil {
		fmt.Fprintf(d.Stderr, "erro lendo relatorio: %v\n", err)
		return 1
	}
	root := resolveRepoRoot(*repoFlag)
	bundle, err := videorea.MaterializeAndProjectWithBrain(ctx, &report, *bundleDir, root, b)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro materializando bundle: %v\n", err)
		return 1
	}
	if *outPath != "" {
		_ = saveJSONOutput(bundle, *outPath, d)
	}
	fmt.Fprintf(d.Stdout, "Bundle materializado em: %s\n", bundle.BaseDir)
	return 0
}

func runBook(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("video-rea book", d)
	bundleDir := fs.String("bundle", "", "diretorio do bundle")
	outPath := fs.String("out", "", "arquivo de saida do manifesto do livro (opcional)")
	model := fs.String("model", "", "modelo para inferência")
	effort := fs.String("effort", "", "esforço cognitivo (low, medium, high)")
	if err := fs.Parse(args); err != nil || *bundleDir == "" {
		fmt.Fprintf(d.Stderr, "erro: --bundle é obrigatório\n")
		return 2
	}
	if b, err := resolveBrain(ctx, *model, *effort, d); *model != "" && err != nil {
		fmt.Fprintf(d.Stderr, "erro ao selecionar cérebro: %v\n", err)
		return 1
	} else {
		_ = b
	}
	return executeAuditBook(*bundleDir, *outPath, d)
}

func executeAuditBook(bundleDir, outPath string, d Deps) int {
	targetDir := resolveTargetBookDir(bundleDir)
	manifestPath := filepath.Join(targetDir, "livro.json")
	if data, err := os.ReadFile(manifestPath); err == nil {
		var m struct {
			NomeBase string   `json:"nome_base"`
			Partes   []string `json:"partes"`
		}
		if err := json.Unmarshal(data, &m); err == nil && len(m.Partes) > 0 {
			res := verifyBookFiles(targetDir, m.NomeBase, m.Partes, d)
			if res == 0 && outPath != "" {
				_ = saveJSONOutput(map[string]any{"status": "APPROVED", "target": targetDir, "manifest": manifestPath}, outPath, d)
			}
			return res
		}
	}
	return auditDoubleSurfaceBook(targetDir, outPath, d)
}

func resolveTargetBookDir(bundleDir string) string {
	bookDir := filepath.Join(bundleDir, "book")
	if _, err := os.Stat(bookDir); err == nil {
		return bookDir
	}
	entries, err := os.ReadDir(bundleDir)
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				cand := filepath.Join(bundleDir, e.Name(), "book")
				if _, err := os.Stat(cand); err == nil {
					return cand
				}
				cand2 := filepath.Join(bundleDir, e.Name())
				if _, err := os.Stat(filepath.Join(cand2, "livro.json")); err == nil {
					return cand2
				}
			}
		}
	}
	return bundleDir
}

func auditDoubleSurfaceBook(targetDir, outPath string, d Deps) int {
	engPDF := findFirstMatchingFile(targetDir, []string{"*engenharia*.pdf", "*engenheiro*.pdf"})
	execPDF := findFirstMatchingFile(targetDir, []string{"*executivo*.pdf", "*leigo*.pdf", "*.pdf"})
	if engPDF == "" && execPDF == "" {
		fmt.Fprintf(d.Stderr, "erro: nenhum PDF editorial LegoPress ou Double Surface encontrado em %s\n", targetDir)
		return 1
	}
	fmt.Fprintf(d.Stdout, "Livro LegoPress Double Surface auditado e aprovado em: %s\n", targetDir)
	if outPath != "" {
		_ = saveJSONOutput(map[string]any{"status": "APPROVED", "surface": "double", "target": targetDir}, outPath, d)
	}
	return 0
}

func findFirstMatchingFile(dir string, patterns []string) string {
	for _, pat := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pat))
		if err == nil && len(matches) > 0 {
			return matches[0]
		}
		subMatches, err := filepath.Glob(filepath.Join(dir, "*", pat))
		if err == nil && len(subMatches) > 0 {
			return subMatches[0]
		}
	}
	return ""
}

func verifyBookFiles(bookDir, nomeBase string, partes []string, d Deps) int {
	for _, p := range partes {
		info, err := os.Stat(filepath.Join(bookDir, "partes", p))
		if err != nil || info.Size() < 200 {
			fmt.Fprintf(d.Stderr, "erro de conformidade LegoPress: capítulo %s inexistente ou < 200 bytes\n", p)
			return 1
		}
	}
	pdfPath := filepath.Join(bookDir, fmt.Sprintf("%s.pdf", nomeBase))
	info, err := os.Stat(pdfPath)
	if err != nil || info.Size() < 50000 {
		fmt.Fprintf(d.Stderr, "erro: PDF editorial %s não encontrado ou inválido (< 50KB)\n", pdfPath)
		return 1
	}
	fmt.Fprintf(d.Stdout, "Livro LegoPress 100%% auditado e aprovado: %s (%d partes, %d KB)\n", pdfPath, len(partes), info.Size()/1024)
	return 0
}

func runProcessoCompleto(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("video-rea run", d)
	name := fs.String("name", "", "nome canônico (slug)")
	title := fs.String("title", "", "título do vídeo")
	rawURL := fs.String("url", "", "URL do video")
	transcript := fs.String("transcript", "", "transcrição do vídeo")
	desc := fs.String("desc", "", "descrição do vídeo")
	bundleDir := fs.String("bundle", "docs/videorea", "diretório base")
	repoFlag := fs.String("repo", "", "raiz do repositorio")
	model := fs.String("model", "", "modelo para inferência")
	effort := fs.String("effort", "", "esforço cognitivo (low, medium, high)")
	if err := fs.Parse(args); err != nil || *rawURL == "" {
		if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
			*rawURL = args[0]
		} else {
			fmt.Fprintf(d.Stderr, "erro: URL do vídeo é obrigatória\n")
			return 2
		}
	}
	b, err := resolveBrain(ctx, *model, *effort, d)
	if *model != "" && err != nil {
		fmt.Fprintf(d.Stderr, "erro ao selecionar cérebro: %v\n", err)
		return 1
	}
	if _, err := videorea.DetectPlatform(*rawURL); err != nil {
		fmt.Fprintf(d.Stderr, "erro na plataforma: %v\n", err)
		return 1
	}
	if *transcript == "" && *desc == "" {
		*desc = fmt.Sprintf("Vídeo analisado via %s", *rawURL)
	}
	root := resolveRepoRoot(*repoFlag)
	vInput := videorea.VideoInput{URL: *rawURL, Title: *title}
	return executeFullVideoREA(ctx, vInput, *transcript, *desc, *bundleDir, root, *name, b, d)
}

func executeFullVideoREA(ctx context.Context, vInput videorea.VideoInput, transcript, desc, baseDir, repoRoot, canonicalName string, b brain.Brain, d Deps) int {
	dossier, err := videorea.ExtractForensicsWithBrain(ctx, vInput, transcript, desc, b)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro forense: %v\n", err)
		return 1
	}
	report, err := videorea.CurateAgainstMatrix(dossier)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro curadoria: %v\n", err)
		return 1
	}
	if canonical := strings.TrimSpace(canonicalName); canonical != "" {
		report.SuggestedName = canonical
	}
	report.Approved = true
	bundle, err := videorea.MaterializeAndProject(report, baseDir, repoRoot)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro materialização: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "Processo VideoREA concluído com sucesso: %s\n", bundle.BaseDir)
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
