package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"legoloop/internal/paths"
	"legoloop/internal/reputation"
)

// goldScore é a nota mínima do Padrão Ouro.
const goldScore = 0.85

type upgradeFlags struct{ diagnose, apply, asJSON bool }

// RunUpgrade comanda a diagnose e a validação de elevação de uma peça para Padrão Ouro.
func RunUpgrade(stdout, stderr io.Writer, args []string) int {
	target, rest := splitPositional(args)
	fs := flag.NewFlagSet("legoloop upgrade", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var f upgradeFlags
	fs.BoolVar(&f.diagnose, "diagnose", false, "Diagnostica gaps e score atual da peça")
	fs.BoolVar(&f.apply, "apply", false, "Valida se a peça reescrita atingiu o Padrão Ouro")
	fs.BoolVar(&f.asJSON, "json", false, "Saída estruturada em JSON")
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	targetPath := target
	if targetPath == "" {
		if fs.NArg() == 0 {
			fmt.Fprintln(stderr, "erro: informe o caminho do arquivo da peça (ex: legoloop upgrade path/to/SKILL.md)")
			return 2
		}
		targetPath = fs.Arg(0)
	}
	rating, err := reputation.ScoreFile(targetPath)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao ler arquivo alvo %s: %v\n", targetPath, err)
		return 1
	}
	return reportUpgrade(stdout, stderr, f, targetPath, rating, upgradeGaps(rating))
}

// upgradeGaps lista o que falta para a peça chegar ao Padrão Ouro, dimensão por dimensão.
func upgradeGaps(rating reputation.PieceRating) []string {
	var gaps []string
	if rating.Dimensions.Rigor < 0.85 {
		gaps = append(gaps, "Falta rigor: adicione guardrails negativos explícitos (NEVER/Zero Stubs) e comandos de teste de terminal")
	}
	if rating.Dimensions.Contracts < 0.85 {
		gaps = append(gaps, "Faltam contratos: declare formalmente I/O (Consumes/Produces) e schemas tipados")
	}
	if rating.Dimensions.History < 0.80 {
		gaps = append(gaps, "Histórico neutro ou com falhas: execute testes com sucesso para consolidar a confiabilidade local")
	}
	return gaps
}

// reportUpgrade escolhe a saída: JSON, diagnóstico em texto (padrão) ou validação do --apply.
func reportUpgrade(stdout, stderr io.Writer, f upgradeFlags, path string, rating reputation.PieceRating, gaps []string) int {
	switch {
	case f.asJSON:
		data, _ := json.MarshalIndent(map[string]any{
			"target":     path,
			"rating":     rating,
			"gaps":       gaps,
			"is_gold":    rating.TotalScore >= goldScore,
			"actionable": len(gaps) > 0,
		}, "", "  ")
		fmt.Fprintln(stdout, string(data))
	case f.diagnose || !f.apply:
		printDiagnosis(stdout, path, rating, gaps)
	default:
		return reportUpgradeApply(stdout, stderr, path, rating, gaps)
	}
	return 0
}

func printDiagnosis(stdout io.Writer, path string, rating reputation.PieceRating, gaps []string) {
	fmt.Fprintf(stdout, "\n🔍 DIAGNÓSTICO DE ELEVAÇÃO: %s\n", path)
	fmt.Fprintf(stdout, "• Score Atual: %.2f | Tier: %s (Badge: %s)\n", rating.TotalScore, strings.ToUpper(string(rating.Tier)), rating.BadgeColor)
	fmt.Fprintf(stdout, "• Autor/Origem: %s\n", rating.Author)
	fmt.Fprintf(stdout, "• Dimensões: Rigor: %.2f | Contratos: %.2f | Proveniência: %.2f | Histórico: %.2f\n\n",
		rating.Dimensions.Rigor, rating.Dimensions.Contracts, rating.Dimensions.Provenance, rating.Dimensions.History)
	if len(gaps) == 0 {
		fmt.Fprintln(stdout, "🏆 PARABÉNS: Esta peça já atinge o Padrão Ouro de engenharia!")
		return
	}
	fmt.Fprintln(stdout, "⚠️ GAPS PARA ATINGIR O PADRÃO OURO (Score >= 0.85):")
	for _, g := range gaps {
		fmt.Fprintf(stdout, "  - %s\n", g)
	}
	fmt.Fprintf(stdout, "\n💡 Recomendação: Aplique %s para elevar a peça.\n", templateHint(path))
}

func reportUpgradeApply(stdout, stderr io.Writer, path string, rating reputation.PieceRating, gaps []string) int {
	if rating.TotalScore >= goldScore {
		fmt.Fprintf(stdout, "✅ SUCESSO: A peça %s foi elevada ao PADRÃO OURO com Score %.2f!\n", filepath.Base(path), rating.TotalScore)
		return 0
	}
	fmt.Fprintf(stderr, "❌ FALHA: A peça ainda não atingiu o Padrão Ouro (Score atual: %.2f < 0.85). Gaps pendentes: %s\n",
		rating.TotalScore, strings.Join(gaps, "; "))
	return 1
}

// templateHint aponta o template canônico (docs/templates) que corresponde ao tipo da peça.
func templateHint(path string) string {
	const dir = "docs/templates/"
	name := strings.ToLower(filepath.Base(path))
	switch ext := filepath.Ext(name); {
	case strings.HasSuffix(name, paths.SufixoAntigo):
		return dir + "TEMPLATE_PIECE_LOOP.json"
	case ext == ".md":
		return dir + "TEMPLATE_PIECE_SKILL.md"
	case slices.Contains([]string{".ps1", ".sh", ".bash", ".py", ".js", ".mjs", ".cjs", ".bat", ".cmd"}, ext):
		return dir + "TEMPLATE_PIECE_SCRIPT.ps1"
	case ext == ".json":
		return jsonTemplateHint(name, dir)
	}
	return dir + "README.md (escolha o template do tipo da peça)"
}

func jsonTemplateHint(name, dir string) string {
	switch {
	case strings.Contains(name, "mcp"):
		return dir + "TEMPLATE_PIECE_MCP.json"
	case strings.Contains(name, "hook") || strings.Contains(name, "settings"):
		return dir + "TEMPLATE_PIECE_HOOK.json"
	case strings.Contains(name, "gate"):
		return dir + "TEMPLATE_PIECE_GATE.json"
	}
	return dir + "TEMPLATE_PIECE_{GATE,HOOK,LOOP,MCP}.json (conforme o tipo)"
}
