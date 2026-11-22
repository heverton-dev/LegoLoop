package forja

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/jail"
)

// ResolverPastaHarness resolve o diretório do harness dentro de home,
// bloqueando caminhos vazios, absolutos, travessia com ".." e fugas do home.
func ResolverPastaHarness(home, nomeHarness string) (string, error) {
	if strings.TrimSpace(home) == "" {
		return "", fmt.Errorf("%w: home vazio", jail.ErrEscapeRaiz)
	}
	if err := validarNomeHarness(nomeHarness); err != nil {
		return "", err
	}
	resolvido, err := jail.DentroDaRaiz(home, nomeHarness)
	if err != nil {
		return "", err
	}
	homeAbs, err := filepath.Abs(home)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(homeAbs, resolvido)
	if err != nil || rel == "." || rel == "" || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("%w: %q", jail.ErrEscapeRaiz, nomeHarness)
	}
	return resolvido, nil
}

func validarNomeHarness(nome string) error {
	limpo := strings.TrimSpace(nome)
	if limpo == "" || strings.ContainsRune(nome, 0) {
		return fmt.Errorf("%w: nome de harness inválido %q", jail.ErrEscapeRaiz, nome)
	}
	norm := strings.ReplaceAll(limpo, `\`, "/")
	if filepath.IsAbs(limpo) || strings.HasPrefix(norm, "/") || temDrive(norm) {
		return fmt.Errorf("%w: caminho absoluto proibido %q", jail.ErrEscapeRaiz, nome)
	}
	for _, parte := range strings.Split(norm, "/") {
		if parte == ".." || parte == "." {
			return fmt.Errorf("%w: travessia proibida %q", jail.ErrEscapeRaiz, nome)
		}
	}
	return nil
}

func temDrive(s string) bool {
	if len(s) < 2 || s[1] != ':' {
		return false
	}
	c := s[0]
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}
