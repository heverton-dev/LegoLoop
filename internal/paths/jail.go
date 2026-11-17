package paths

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrEscapeRaiz sinaliza entrada hostil que tenta escapar da raiz sanitizada.
var ErrEscapeRaiz = errors.New("paths: escape da raiz sanitizada")

// ErrPIDInvalido sinaliza PID falso (não positivo) vindo de entrada hostil.
var ErrPIDInvalido = errors.New("paths: PID inválido")

// DentroDaRaiz resolve alvo contra raiz e devolve o caminho absoluto somente se ele
// permanecer dentro da raiz sanitizada. Recusa vazio, "..", absolutos externos e UNC.
func DentroDaRaiz(raiz, alvo string) (string, error) {
	if strings.TrimSpace(alvo) == "" {
		return "", fmt.Errorf("%w: caminho vazio", ErrEscapeRaiz)
	}
	base, err := filepath.Abs(raiz)
	if err != nil {
		return "", err
	}
	alvoAbs := alvo
	if !filepath.IsAbs(alvo) {
		alvoAbs = filepath.Join(base, alvo)
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
