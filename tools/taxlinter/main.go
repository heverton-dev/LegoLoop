package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"legoloop/internal/taxonomy"
)

type TaxonomyFileViolation struct {
	FilePath  string
	Violation taxonomy.TaxonomyViolation
}

var ignoredPrefixes = []string{
	".git",
	"bin",
	"vendor",
	"docs/sessions",
	".codebase-memory",
	"docs/vocab",
	"internal/taxonomy",
}

var jsonKeyRegex = regexp.MustCompile(`"([^"]+)"\s*:`)

func shouldIgnore(relPath string) bool {
	normalized := filepath.ToSlash(relPath)
	for _, p := range ignoredPrefixes {
		if normalized == p || strings.HasPrefix(normalized, p+"/") {
			return true
		}
	}
	return false
}

func isScannableExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".go" || ext == ".json" || ext == ".md"
}

func isTechnicalLiteral(val string) bool {
	if strings.Contains(val, "/") || strings.Contains(val, "\\") || strings.Contains(val, "-") {
		return true
	}
	if strings.Contains(val, ".") || strings.HasPrefix(val, ".") || strings.HasSuffix(val, ".json") || strings.HasSuffix(val, ".go") {
		return true
	}
	return false
}

func scanGoFile(filePath string) ([]TaxonomyFileViolation, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var violations []TaxonomyFileViolation
	checkText := func(text string, pos token.Pos) {
		vs, err := taxonomy.ValidateText(text)
		if err != nil {
			return
		}
		position := fset.Position(pos)
		for _, v := range vs {
			v.Line = position.Line + v.Line - 1
			violations = append(violations, TaxonomyFileViolation{FilePath: filePath, Violation: v})
		}
	}
	for _, cg := range node.Comments {
		for _, c := range cg.List {
			checkText(c.Text, c.Pos())
		}
	}
	ast.Inspect(node, func(n ast.Node) bool {
		if _, ok := n.(*ast.ImportSpec); ok {
			return false
		}
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			val := strings.Trim(lit.Value, "`\"")
			if !isTechnicalLiteral(val) {
				checkText(val, lit.Pos())
			}
		}
		return true
	})
	return violations, nil
}

func scanJSONKeys(filePath string) ([]TaxonomyFileViolation, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var violations []TaxonomyFileViolation
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		matches := jsonKeyRegex.FindAllStringSubmatch(scanner.Text(), -1)
		for _, m := range matches {
			if len(m) > 1 {
				key := m[1]
				if canon, ok := taxonomy.CheckForbidden(key); ok {
					violations = append(violations, TaxonomyFileViolation{
						FilePath: filePath,
						Violation: taxonomy.TaxonomyViolation{
							ForbiddenTerm: key,
							ApprovedTerm:  canon,
							Line:          lineNum,
							Column:        1,
						},
					})
				}
			}
		}
	}
	return violations, scanner.Err()
}

func scanMarkdownHeadings(filePath string) ([]TaxonomyFileViolation, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var violations []TaxonomyFileViolation
	scanner := bufio.NewScanner(f)
	lineNum := 0
	inCodeBlock := false
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}
		if strings.HasPrefix(line, "#") {
			heading := strings.TrimLeft(line, "# \t")
			vs, _ := taxonomy.ValidateText(heading)
			for _, v := range vs {
				v.Line = lineNum
				violations = append(violations, TaxonomyFileViolation{FilePath: filePath, Violation: v})
			}
		}
	}
	return violations, scanner.Err()
}

func scanFile(filePath string) ([]TaxonomyFileViolation, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return scanGoFile(filePath)
	case ".json":
		return scanJSONKeys(filePath)
	case ".md":
		return scanMarkdownHeadings(filePath)
	default:
		return nil, nil
	}
}

func scanDirectory(root string) ([]TaxonomyFileViolation, error) {
	var results []TaxonomyFileViolation
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(".", p)
		if err != nil {
			rel = p
		}
		if d.IsDir() {
			if shouldIgnore(rel) {
				return filepath.SkipDir
			}
			return nil
		}
		if shouldIgnore(rel) || !isScannableExt(rel) {
			return nil
		}
		fileViolations, err := scanFile(p)
		if err != nil {
			return err
		}
		results = append(results, fileViolations...)
		return nil
	})
	return results, err
}

func runLinter(targets []string) ([]TaxonomyFileViolation, error) {
	var total []TaxonomyFileViolation
	for _, target := range targets {
		info, err := os.Stat(target)
		if err != nil {
			return nil, fmt.Errorf("alvo invalido %q: %w", target, err)
		}
		if info.IsDir() {
			res, err := scanDirectory(target)
			if err != nil {
				return nil, err
			}
			total = append(total, res...)
		} else {
			res, err := scanFile(target)
			if err != nil {
				return nil, err
			}
			total = append(total, res...)
		}
	}
	return total, nil
}

func printViolations(violations []TaxonomyFileViolation) {
	fmt.Fprintf(os.Stderr, "=== VIOLACOES TAXONOMICAS ASD-STE100 ENCONTRADAS (%d) ===\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "%s:%d:%d: termo proibido %q -> use o termo canonico %q (dominio: %s)\n",
			v.FilePath,
			v.Violation.Line,
			v.Violation.Column,
			v.Violation.ForbiddenTerm,
			v.Violation.ApprovedTerm,
			v.Violation.ContextDomain,
		)
	}
}

func main() {
	targets := os.Args[1:]
	if len(targets) == 0 {
		targets = []string{"."}
	}
	violations, err := runLinter(targets)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERRO NO TAXLINTER: %v\n", err)
		os.Exit(2)
	}
	if len(violations) > 0 {
		printViolations(violations)
		os.Exit(1)
	}
	fmt.Println("Linter taxonomico ASD-STE100: zero violacoes.")
	os.Exit(0)
}
