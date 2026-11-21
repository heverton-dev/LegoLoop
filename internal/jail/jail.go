package jail

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrEscapeRaiz sinaliza entrada hostil que tenta escapar da raiz sanitizada.
var ErrEscapeRaiz = errors.New("jail: escape da raiz sanitizada")

// ErrPIDInvalido sinaliza PID falso (não positivo) vindo de entrada hostil.
var ErrPIDInvalido = errors.New("jail: PID inválido")

func temDriveLetter(s string) bool {
	return len(s) >= 2 && ((s[0] >= 'A' && s[0] <= 'Z') || (s[0] >= 'a' && s[0] <= 'z')) && s[1] == ':'
}

func normalizarAlvo(alvo string) (string, error) {
	if strings.TrimSpace(alvo) == "" {
		return "", fmt.Errorf("%w: caminho vazio", ErrEscapeRaiz)
	}
	norm := strings.ReplaceAll(alvo, `\`, "/")
	if strings.HasPrefix(norm, "//") {
		return "", fmt.Errorf("%w: %q", ErrEscapeRaiz, alvo)
	}
	alvoNorm := filepath.FromSlash(norm)
	if (temDriveLetter(norm) || strings.HasPrefix(norm, "/")) && !filepath.IsAbs(alvoNorm) {
		return "", fmt.Errorf("%w: %q", ErrEscapeRaiz, alvo)
	}
	return alvoNorm, nil
}

// DentroDaRaiz resolve alvo contra raiz e devolve o caminho absoluto somente se ele
// permanecer dentro da raiz sanitizada. Recusa vazio, "..", absolutos externos e UNC.
func DentroDaRaiz(raiz, alvo string) (string, error) {
	alvoNorm, err := normalizarAlvo(alvo)
	if err != nil {
		return "", err
	}
	base, err := filepath.Abs(raiz)
	if err != nil {
		return "", err
	}
	alvoAbs := alvoNorm
	if !filepath.IsAbs(alvoNorm) {
		alvoAbs = filepath.Join(base, alvoNorm)
	}
	alvoAbs, err = filepath.Abs(alvoAbs)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(base, alvoAbs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("%w: %q", ErrEscapeRaiz, alvo)
	}
	return alvoAbs, nil
}

// CampoDeJSON valida campo de caminho vindo de JSON hostil: recusa byte nulo e
// delega a DentroDaRaiz o bloqueio de escape.
func CampoDeJSON(raiz, campo string) (string, error) {
	if strings.ContainsRune(campo, 0) {
		return "", fmt.Errorf("%w: byte nulo em campo de JSON", ErrEscapeRaiz)
	}
	return DentroDaRaiz(raiz, campo)
}

// ValidarPID recusa PID falso (não positivo) vindo de entrada hostil.
func ValidarPID(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("%w: %d", ErrPIDInvalido, pid)
	}
	return nil
}

// ValidarNomeProcesso recusa nomes ou slugs de Processo com travessia,
// separadores de caminho (portável entre Windows e Linux), drive ou byte nulo.
func ValidarNomeProcesso(nome string) error {
	limpo := strings.TrimSpace(nome)
	if limpo == "" || limpo == "." || strings.ContainsRune(nome, 0) ||
		strings.Contains(nome, "..") || strings.ContainsAny(nome, `/\:`) {
		return fmt.Errorf("%w: nome de processo inválido %q", ErrEscapeRaiz, nome)
	}
	if _, err := DentroDaRaiz("processos", limpo); err != nil {
		return err
	}
	return nil
}

// Renomear move de/para dentro da raiz sanitizada; escape em qualquer ponta é
// bloqueado sem efeito em disco.
func Renomear(raiz, de, para string) error {
	deAbs, err := DentroDaRaiz(raiz, de)
	if err != nil {
		return err
	}
	paraAbs, err := DentroDaRaiz(raiz, para)
	if err != nil {
		return err
	}
	return os.Rename(deAbs, paraAbs)
}
