package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// execCmd roda um comando em dir e devolve a saída combinada. Os testes o substituem.
var execCmd = func(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func defaultBin() string {
	if runtime.GOOS == "windows" {
		return filepath.Join("bin", "legoloop.exe")
	}
	return filepath.Join("bin", "legoloop")
}

func cmdStep(name, dir, bin string, args ...string) step {
	return step{Name: name, Run: func() (string, error) { return execCmd(dir, bin, args...) }}
}

func testStep(c config) step {
	args := append([]string{"test", "-count=1"}, c.Pkgs...)
	return cmdStep("go test "+strings.Join(c.Pkgs, " "), c.Root, "go", args...)
}

func vetStep(c config) step { return cmdStep("go vet ./...", c.Root, "go", "vet", "./...") }

func funclenStep(c config) step {
	return cmdStep("funclen internal", c.Root, "go", "run", "./tools/funclen", "internal")
}

func modStep(c config) step {
	return cmdStep("go.mod e go.sum intactos", c.Root, "git", "diff", "--exit-code", "--", "go.mod", "go.sum")
}

func goldStep(c config) step {
	return cmdStep("teste dourado", c.Root, "go", "test", "-count=1", "./internal/golden/...")
}

func buildStep(c config) step { return cmdStep("go build ./...", c.Root, "go", "build", "./...") }

func taxStep(c config) step {
	args := append([]string{"run", "./tools/taxlinter"}, c.Tax...)
	return cmdStep("taxlinter "+strings.Join(c.Tax, " "), c.Root, "go", args...)
}

// gofmtPaths devolve a raiz mais as pastas dos pacotes pedidos, sem o sufixo "/...".
func gofmtPaths(c config) []string {
	paths := []string{"."}
	for _, p := range c.Pkgs {
		p = strings.TrimSuffix(strings.TrimSuffix(p, "..."), "/")
		if p != "" && p != "." {
			paths = append(paths, p)
		}
	}
	return paths
}

// gofmtStep falha quando gofmt -l imprime algum arquivo: gofmt sai com 0 nesse caso.
func gofmtStep(c config) step {
	return step{Name: "gofmt -l", Run: func() (string, error) {
		out, err := execCmd(c.Root, "gofmt", append([]string{"-l"}, gofmtPaths(c)...)...)
		if err != nil {
			return out, err
		}
		if strings.TrimSpace(out) != "" {
			return out, errors.New("arquivos sem formatação")
		}
		return "", nil
	}}
}

// listPDFs devolve os PDFs de uma pasta. Pasta inexistente ou sem PDF é erro.
func listPDFs(root, dir string) ([]string, error) {
	var pdfs []string
	base := filepath.Join(root, dir)
	err := filepath.WalkDir(base, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.EqualFold(filepath.Ext(p), ".pdf") {
			pdfs = append(pdfs, p)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("pasta de livros %q: %w", dir, err)
	}
	if len(pdfs) == 0 {
		return nil, fmt.Errorf("nenhum PDF em %q", dir)
	}
	return pdfs, nil
}

// livrosStep audita cada PDF com "legoloop legopress auditar <pdf>".
func livrosStep(c config) step {
	return step{Name: "legopress auditar " + strings.Join(c.Livros, " "), Run: func() (string, error) {
		var report []string
		for _, dir := range c.Livros {
			pdfs, err := listPDFs(c.Root, dir)
			if err != nil {
				return "", err
			}
			for _, pdf := range pdfs {
				if out, err := execCmd(c.Root, c.Bin, "legopress", "auditar", pdf); err != nil {
					report = append(report, pdf+": "+strings.TrimSpace(out))
				}
			}
		}
		if len(report) > 0 {
			return strings.Join(report, "\n"), errors.New("PDF reprovado na auditoria")
		}
		return "", nil
	}}
}
