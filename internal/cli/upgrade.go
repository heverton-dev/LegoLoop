package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"legoloop/internal/reputation"
)

// RunUpgrade comanda a diagnose e a validação de elevação de uma peça para Padrão Ouro.
func RunUpgrade(stdout, stderr io.Writer, args []string) int {
	target, rest := splitPositional(args)

	fs := flag.NewFlagSet("legoloop upgrade", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		diagnose = fs.Bool("diagnose", false, "Diagnostica gaps e score atual da peça")
		apply    = fs.Bool("apply", false, "Valida se a peça reescrita atingiu o Padrão Ouro")
		asJSON   = fs.Bool("json", false, "Saída estruturada em JSON")
	)

	if err := fs.Parse(rest); err != nil {
		return 2
	}

	targetPath := target
	if targetPath == "" {
		if fs.NArg() > 0 {
			targetPath = fs.Arg(0)
		} else {
			fmt.Fprintln(stderr, "erro: informe o caminho do arquivo da peça (ex: legoloop upgrade path/to/SKILL.md)")
			return 2
		}
	}
	rating, err := reputation.ScoreFile(targetPath)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao ler arquivo alvo %s: %v\n", targetPath, err)
		return 1
	}

	// Gaps diagnosticados
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

	if *asJSON {
		out := map[string]any{
			"target":     targetPath,
			"rating":     rating,
			"gaps":       gaps,
			"is_gold":    rating.TotalScore >= 0.85,
			"actionable": len(gaps) > 0,
		}
		data, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(stdout, string(data))
		return 0
	}

	// Saída em modo Texto/Tabela
	if *diagnose || (!*diagnose && !*apply) {
		fmt.Fprintf(stdout, "\n🔍 DIAGNÓSTICO DE ELEVAÇÃO: %s\n", targetPath)
		fmt.Fprintf(stdout, "• Score Atual: %.2f | Tier: %s (Badge: %s)\n", rating.TotalScore, strings.ToUpper(string(rating.Tier)), rating.BadgeColor)
		fmt.Fprintf(stdout, "• Autor/Origem: %s\n", rating.Author)
		fmt.Fprintf(stdout, "• Dimensões: Rigor: %.2f | Contratos: %.2f | Proveniência: %.2f | Histórico: %.2f\n\n",
			rating.Dimensions.Rigor, rating.Dimensions.Contracts, rating.Dimensions.Provenance, rating.Dimensions.History)

		if len(gaps) > 0 {
			fmt.Fprintln(stdout, "⚠️ GAPS PARA ATINGIR O PADRÃO OURO (Score >= 0.85):")
			for _, g := range gaps {
				fmt.Fprintf(stdout, "  - %s\n", g)
			}
			fmt.Fprintf(stdout, "\n💡 Recomendação: Aplique %s para elevar a peça.\n", templateHint(targetPath))
		} else {
			fmt.Fprintln(stdout, "🏆 PARABÉNS: Esta peça já atinge o Padrão Ouro de engenharia!")
		}
		return 0
	}

	if *apply {
		if rating.TotalScore >= 0.85 {
			fmt.Fprintf(stdout, "✅ SUCESSO: A peça %s foi elevada ao PADRÃO OURO com Score %.2f!\n", filepath.Base(targetPath), rating.TotalScore)
			return 0
		}
		fmt.Fprintf(stderr, "❌ FALHA: A peça ainda não atingiu o Padrão Ouro (Score atual: %.2f < 0.85). Gaps pendentes: %s\n",
			rating.TotalScore, strings.Join(gaps, "; "))
		return 1
	}

	return 0
}

// templateHint aponta o template canônico (docs/templates) que corresponde ao tipo da peça.
func templateHint(path string) string {
	const dir = "docs/templates/"
	name := strings.ToLower(filepath.Base(path))
	switch ext := filepath.Ext(name); {
	case strings.HasSuffix(name, ".loop.json"):
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
