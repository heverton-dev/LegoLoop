// Package upstream aplica a trava da cadeia de artefatos: nenhuma SPEC, PLAN ou DAG nasce sem uma
// INTENT aprovada, que por sua vez só existe depois da concepção (entrevista socrática).
package upstream

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// ErrSemIntencao: a SPEC não declara a Referência da Intenção.
	ErrSemIntencao = errors.New("upstream: a SPEC não referencia uma INTENT; rode a concepção (legoloop-concepcao) antes")
	// ErrIntencaoAusente: a INTENT referenciada não existe no disco.
	ErrIntencaoAusente = errors.New("upstream: a INTENT referenciada não existe")
	// ErrIntencaoNaoAprovada: a INTENT existe, mas o Stop-Gate 01 não a aprovou (status diferente de APPROVED).
	ErrIntencaoNaoAprovada = errors.New("upstream: a INTENT não está aprovada (status APPROVED no cabeçalho)")
)

// Recusa diz se o erro é uma recusa da trava (e não uma falha de infraestrutura).
// Quem expõe a trava por HTTP devolve 422 nesse caso, e não 500.
func Recusa(err error) bool {
	return errors.Is(err, ErrSemIntencao) || errors.Is(err, ErrIntencaoAusente) || errors.Is(err, ErrIntencaoNaoAprovada)
}

var (
	refIntencaoRe = regexp.MustCompile(`Referência da Intenção:\*\*\s*(?:\[([^\]]+)\]|(\S+))`)
	statusRe      = regexp.MustCompile(`(?im)^status:\s*(\S+)\s*$`)
)

// StatusDoDocumento devolve o status do cabeçalho (front matter) de um documento do sexteto.
// Só o cabeçalho vale: um "status:" no corpo nunca aprova nada.
func StatusDoDocumento(b []byte) string {
	texto := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.HasPrefix(texto, "---\n") {
		return ""
	}
	fim := strings.Index(texto[4:], "\n---")
	if fim < 0 {
		return ""
	}
	if m := statusRe.FindStringSubmatch(texto[4 : 4+fim]); m != nil {
		return strings.ToUpper(m[1])
	}
	return ""
}

// IntencaoAprovada confere que o arquivo existe e que o cabeçalho diz status APPROVED.
func IntencaoAprovada(caminho string) error {
	b, err := os.ReadFile(caminho)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrIntencaoAusente, filepath.ToSlash(caminho))
	}
	if StatusDoDocumento(b) != "APPROVED" {
		return fmt.Errorf("%w: %s", ErrIntencaoNaoAprovada, filepath.ToSlash(caminho))
	}
	return nil
}

// ReferenciaDaIntencao extrai da SPEC o caminho da INTENT, ou "" quando a SPEC não a declara.
func ReferenciaDaIntencao(specTexto string) string {
	m := refIntencaoRe.FindStringSubmatch(specTexto)
	if m == nil {
		return ""
	}
	v := strings.TrimSpace(m[1] + m[2])
	return strings.ReplaceAll(v, `\`, "/")
}

// VerificarSpec confere a trava de uma SPEC: ela referencia uma INTENT que existe e está aprovada.
// O caminho da INTENT vale a partir da raiz do projeto ou, na falta, a partir da pasta da SPEC.
func VerificarSpec(baseDir, specPath string) error {
	b, err := os.ReadFile(specPath)
	if err != nil {
		return err
	}
	ref := ReferenciaDaIntencao(string(b))
	if ref == "" {
		return ErrSemIntencao
	}
	candidatos := []string{filepath.Join(baseDir, filepath.FromSlash(ref)), filepath.Join(filepath.Dir(specPath), filepath.FromSlash(ref))}
	for _, c := range candidatos {
		if _, err := os.Stat(c); err == nil {
			return IntencaoAprovada(c)
		}
	}
	return fmt.Errorf("%w: %s", ErrIntencaoAusente, ref)
}
