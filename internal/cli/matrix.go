package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"legoloop/internal/reputation"
)

// RunMatrix lista e filtra a matriz algorítmica de reputação das peças.
func RunMatrix(stdout, stderr io.Writer, args []string) int {
	fs := flag.NewFlagSet("legoloop matrix", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		kindPath = fs.String("path", "", "Caminho alternativo para o registry_matrix.json")
		kind     = fs.String("kind", "", "Filtra por categoria: skill, script, hook, mcp, gate, loop")
		tier     = fs.String("tier", "", "Filtra por tier: ouro, prata, bronze, toxica")
		asJSON   = fs.Bool("json", false, "Saída estruturada em JSON")
		minScore = fs.Float64("min", 0.0, "Score mínimo (0.0 a 1.0)")
	)

	if err := fs.Parse(args); err != nil {
		return 2
	}

	p := *kindPath
	if p == "" {
		p = reputation.DefaultRegistryPath("")
	}

	reg, err := reputation.LoadRegistry(p)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao ler registro matricial em %s: %v\n", p, err)
		return 1
	}

	var filtered []reputation.PieceRating
	filterKind := strings.ToLower(strings.TrimSpace(*kind))
	filterTier := strings.ToLower(strings.TrimSpace(*tier))

	for _, r := range reg.Ratings {
		if *minScore > 0 && r.TotalScore < *minScore {
			continue
		}
		if filterTier != "" && strings.ToLower(string(r.Tier)) != filterTier {
			continue
		}
		if filterKind != "" {
			refLower := strings.ToLower(r.Ref)
			matched := false
			switch filterKind {
			case "skill":
				matched = strings.Contains(refLower, "skill")
			case "script":
				matched = strings.HasSuffix(refLower, ".ps1") || strings.HasSuffix(refLower, ".sh") || strings.HasSuffix(refLower, ".bat") || strings.HasSuffix(refLower, ".cmd") || strings.Contains(refLower, "scripts")
			case "hook":
				matched = strings.Contains(refLower, "hook")
			case "mcp":
				matched = strings.Contains(refLower, "mcp")
			case "gate":
				matched = strings.Contains(refLower, "gate")
			case "loop":
				matched = strings.Contains(refLower, "loop") || strings.HasSuffix(refLower, ".dag.json")
			default:
				matched = strings.Contains(refLower, filterKind)
			}
			if !matched {
				continue
			}
		}
		filtered = append(filtered, r)
	}

	if *asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]any{
			"total_found": len(filtered),
			"ratings":     filtered,
		})
		return 0
	}

	fmt.Fprintf(stdout, "\n=== MATRIZ ALGORÍTMICA DE REPUTAÇÃO (%d PEÇAS ENCONTRADAS) ===\n\n", len(filtered))
	w := tabwriter.NewWriter(stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "TIER\tSCORE\tNOME\tAUTOR\tRIGOR\tCONTRATOS\tPROVENIÊNCIA\tHISTÓRICO")
	fmt.Fprintln(w, "----\t-----\t----\t-----\t-----\t---------\t-----------\t---------")

	for _, r := range filtered {
		tierBadge := strings.ToUpper(string(r.Tier))
		name := filepath.Base(r.Name)
		if name == "" || name == "." {
			name = filepath.Base(r.Ref)
		}
		author := r.Author
		if author == "" {
			author = "Desconhecido"
		}
		fmt.Fprintf(w, "%s\t%.2f\t%s\t%s\t%.2f\t%.2f\t%.2f\t%.2f\n",
			tierBadge, r.TotalScore, name, author,
			r.Dimensions.Rigor, r.Dimensions.Contracts, r.Dimensions.Provenance, r.Dimensions.History,
		)
	}
	w.Flush()
	fmt.Fprintln(stdout, "")
	return 0
}
