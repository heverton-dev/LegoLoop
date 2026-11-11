package orchestrix

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"legoloop/internal/taxonomy"
)

var (
	brandFilter = regexp.MustCompile(`(?i)legoloop`)
	wordExtract = regexp.MustCompile(`[\p{L}]+`)
)

// cleanSlopComments remove linhas de comentários que contêm sinônimos proibidos.
func cleanSlopComments(content string) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	changed := false
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "//go:") {
			commentText := strings.TrimPrefix(trimmed, "//")
			if hasForbiddenComment(commentText) {
				changed = true
				continue
			}
		}
		lines = append(lines, line)
	}
	if !changed {
		return content, false
	}
	return strings.Join(lines, "\n") + "\n", true
}

// hasForbiddenComment verifica se o texto do comentário contém sinônimo proibido.
func hasForbiddenComment(text string) bool {
	clean := brandFilter.ReplaceAllString(text, " ")
	words := wordExtract.FindAllString(clean, -1)
	for _, w := range words {
		lw := strings.ToLower(w)
		if _, ok := taxonomy.CheckForbidden(lw); ok {
			return true
		}
		if len(lw) > 2 && strings.HasSuffix(lw, "s") {
			if _, ok := taxonomy.CheckForbidden(lw[:len(lw)-1]); ok {
				return true
			}
		}
	}
	return false
}

// resolveGofmt acha o executável gofmt no PATH ou no GOROOT.
func resolveGofmt() string {
	if _, err := exec.LookPath("gofmt"); err == nil {
		return "gofmt"
	}
	goroot := strings.TrimSpace(os.Getenv("GOROOT"))
	if goroot == "" {
		if out, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
			goroot = strings.TrimSpace(string(out))
		}
	}
	if goroot == "" {
		return "gofmt"
	}
	exe := "gofmt"
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	return filepath.Join(goroot, "bin", exe)
}

// formatGoFile executa gofmt -w no caminho do arquivo.
func formatGoFile(ctx context.Context, fullPath string) error {
	cmd := exec.CommandContext(ctx, resolveGofmt(), "-w", fullPath)
	return cmd.Run()
}

// hygieneFile processa um arquivo alterado: unslop em comentários e gofmt.
func (o *Orchestrix) hygieneFile(ctx context.Context, dir, rel string) error {
	full := filepath.Join(dir, filepath.FromSlash(rel))
	if !strings.EqualFold(filepath.Ext(full), ".go") {
		return nil
	}
	data, err := os.ReadFile(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	cleaned, changed := cleanSlopComments(string(data))
	if changed {
		if err := os.WriteFile(full, []byte(cleaned), 0o644); err != nil {
			return err
		}
	}
	return formatGoFile(ctx, full)
}

// applyHygiene executa a rotina de higiene (unslop + no-comments + gofmt) antes do gate.
func (o *Orchestrix) applyHygiene(ctx context.Context, dir string, changed []string) error {
	for _, rel := range changed {
		if err := o.hygieneFile(ctx, dir, rel); err != nil {
			return err
		}
	}
	return nil
}
