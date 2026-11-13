package dissecar

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrPathEscape = errors.New("caminho fora da sandbox (Path Jail)")
)

// Target representa o alvo a ser dissecado.
type Target struct {
	Path string
	Type string // "code", "binary", "web"
}

// Result representa o resultado da dissecação.
type Result struct {
	TargetID   string            `json:"target_id"`
	AstNodes   []string          `json:"ast_nodes,omitempty"`
	Syscalls   []string          `json:"syscalls,omitempty"`
	ReverseDAG map[string][]string `json:"reverse_dag,omitempty"`
}

// ValidatePath valida se o caminho respeita o Path Jail (restrito a baseDir).
func ValidatePath(baseDir, targetPath string) (string, error) {
	cleanBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}

	fullTarget := targetPath
	if !filepath.IsAbs(fullTarget) {
		fullTarget = filepath.Join(cleanBase, targetPath)
	}

	cleanTarget, err := filepath.Abs(fullTarget)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(cleanBase, cleanTarget)
	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." {
		return "", ErrPathEscape
	}

	return cleanTarget, nil
}
