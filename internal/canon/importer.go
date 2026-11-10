package canon

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GitHubInfo armazena os metadados parseados de uma URL do GitHub.
type GitHubInfo struct {
	Owner         string
	Repo          string
	Branch        string
	Subpath       string
	SuggestedName string
}

// ImportResult contém os dados do ativo após cópia e canonização com SHA-256.
type ImportResult struct {
	Name       string `json:"name"`
	TargetDir  string `json:"target_dir"`
	SHA256     string `json:"sha256"`
	FilesCount int    `json:"files_count"`
}

// ParseGitHubURL analisa e valida com segurança uma URL apontando para um repositório ou subpasta do GitHub.
func ParseGitHubURL(raw string) (GitHubInfo, error) {
	if raw == "" {
		return GitHubInfo{}, errors.New("url não pode ser vazia")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return GitHubInfo{}, fmt.Errorf("url inválida: %w", err)
	}
	parts, err := validateGitHubPath(u)
	if err != nil {
		return GitHubInfo{}, err
	}
	return githubInfoFromParts(parts), nil
}

// validateGitHubPath exige host github.com, rejeita traversal e devolve os segmentos do caminho.
func validateGitHubPath(u *url.URL) ([]string, error) {
	host := strings.ToLower(u.Host)
	if host != "github.com" && host != "www.github.com" {
		return nil, fmt.Errorf("apenas URLs do github.com são aceitas, obteve %q", u.Host)
	}
	if strings.Contains(u.Path, "..") {
		return nil, errors.New("path traversal não permitido")
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return nil, errors.New("a URL deve conter pelo menos /owner/repo")
	}
	return parts, nil
}

func githubInfoFromParts(parts []string) GitHubInfo {
	info := GitHubInfo{
		Owner: parts[0],
		Repo:  strings.TrimSuffix(parts[1], ".git"),
	}

	// Trata URL de subpasta: /owner/repo/tree/branch/sub/path...
	if len(parts) >= 4 && parts[2] == "tree" {
		info.Branch = parts[3]
		if len(parts) > 4 {
			info.Subpath = strings.Join(parts[4:], "/")
			info.SuggestedName = parts[len(parts)-1]
		}
	}

	if info.SuggestedName == "" {
		info.SuggestedName = info.Repo
	}

	return info
}

// InstallAndCanonizeSkill copia os arquivos da skill para a pasta global canônica e calcula o SHA-256.
func InstallAndCanonizeSkill(srcDir string, skillName string, destRoot string) (ImportResult, error) {
	if skillName == "" {
		return ImportResult{}, errors.New("skillName obrigatório")
	}

	// Sanitiza o nome para evitar escape de diretório
	cleanName := filepath.Base(skillName)
	if cleanName == "." || cleanName == "/" || cleanName == "\\" {
		return ImportResult{}, errors.New("nome de skill inválido")
	}

	targetDir := filepath.Join(destRoot, "skills", cleanName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return ImportResult{}, fmt.Errorf("falha ao criar pasta destino %s: %w", targetDir, err)
	}

	hasher := sha256.New()
	filesCount, err := copyTreeHashed(srcDir, targetDir, hasher)
	if err != nil {
		return ImportResult{}, fmt.Errorf("falha ao copiar arquivos para canonização: %w", err)
	}

	return ImportResult{
		Name:       cleanName,
		TargetDir:  targetDir,
		SHA256:     hex.EncodeToString(hasher.Sum(nil)),
		FilesCount: filesCount,
	}, nil
}

// copyTreeHashed copia srcDir para targetDir (ignora .git) alimentando o hash com cada byte copiado.
func copyTreeHashed(srcDir, targetDir string, hasher io.Writer) (int, error) {
	filesCount := 0
	err := filepath.Walk(srcDir, func(path string, d os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if err := copyFileHashed(path, filepath.Join(targetDir, rel), hasher); err != nil {
			return err
		}
		filesCount++
		return nil
	})
	return filesCount, err
}

func copyFileHashed(srcPath, destPath string, hasher io.Writer) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(io.MultiWriter(destFile, hasher), srcFile)
	return err
}

// ImportSkillFromURL clona a skill do GitHub e a instala no diretório global canônico.
func ImportSkillFromURL(ctx context.Context, rawURL string, customName string, destRoot string) (ImportResult, error) {
	info, err := ParseGitHubURL(rawURL)
	if err != nil {
		return ImportResult{}, err
	}

	skillName := customName
	if skillName == "" {
		skillName = info.SuggestedName
	}

	tmpDir, err := os.MkdirTemp("", "lego_gh_import_*")
	if err != nil {
		return ImportResult{}, fmt.Errorf("falha ao criar pasta temporária: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := cloneRepo(ctx, info, tmpDir); err != nil {
		return ImportResult{}, err
	}

	sourceSkillDir := tmpDir
	if info.Subpath != "" {
		sourceSkillDir = filepath.Join(tmpDir, filepath.FromSlash(info.Subpath))
	}

	// Verifica se a pasta existe e contém arquivos
	if _, err := os.Stat(sourceSkillDir); err != nil {
		return ImportResult{}, fmt.Errorf("subpasta %q não encontrada no repositório: %w", info.Subpath, err)
	}

	return InstallAndCanonizeSkill(sourceSkillDir, skillName, destRoot)
}

// cloneRepo clona o repositório (raso, na branch pedida) em tmpDir com timeout de 45s.
func cloneRepo(ctx context.Context, info GitHubInfo, tmpDir string) error {
	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", info.Owner, info.Repo)
	cloneCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	args := []string{"clone", "--depth", "1"}
	if info.Branch != "" {
		args = append(args, "--branch", info.Branch)
	}
	args = append(args, repoURL, tmpDir)
	if output, err := exec.CommandContext(cloneCtx, "git", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("falha ao clonar repositório %s: %s (%w)", repoURL, string(output), err)
	}
	return nil
}
