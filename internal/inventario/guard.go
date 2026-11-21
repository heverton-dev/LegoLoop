package inventario

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrDelecaoRecusada sinaliza deleção recusada em diretório de lei.
var ErrDelecaoRecusada = errors.New("deleção recusada em diretório de lei (Lei 6/15)")

// diretoriosDeLei lista os diretórios protegidos pela Lei 6/15.
var diretoriosDeLei = []string{"templates", "vocab"}

// normalizarSeparadores converte barras invertidas e separadores nativos para '/'.
func normalizarSeparadores(alvo string) string {
	return strings.ReplaceAll(filepath.ToSlash(alvo), `\`, `/`)
}

// Protegido informa se o caminho está sob um diretório de lei.
func Protegido(alvo string) bool {
	normalizado := normalizarSeparadores(alvo)
	for _, dir := range diretoriosDeLei {
		if normalizado == dir || normalizado == "docs/"+dir ||
			strings.HasPrefix(normalizado, dir+"/") || strings.HasPrefix(normalizado, "docs/"+dir+"/") {
			return true
		}
	}
	return false
}

// RecusaDelecao devolve erro quando a deleção atinge diretório de lei.
func RecusaDelecao(alvo string) error {
	if Protegido(alvo) {
		return fmt.Errorf("%w: %s", ErrDelecaoRecusada, alvo)
	}
	return nil
}

// ExitCode devolve o código de saída da operação: 0 aceita, 2 recusada.
func ExitCode(err error) int {
	if err != nil {
		return 2
	}
	return 0
}

// RegistrarTentativa anexa a tentativa de deleção no relatório do guard.
func RegistrarTentativa(raizRepo, alvo string, recusa error) error {
	status, fundamento := "PERMITIDA", "-"
	if recusa != nil {
		status, fundamento = "RECUSADA", "Lei 6/15"
	}
	registro := filepath.Join(raizRepo, "docs", "reports", "2026-11-16_GUARD_LOG.md")
	if err := os.MkdirAll(filepath.Dir(registro), 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(registro); os.IsNotExist(err) {
		cabecalho := "# Registro do guard de lei — 2026-11-16\n\n| Alvo | Status | Fundamento |\n| :--- | :--- | :--- |\n"
		if err := os.WriteFile(registro, []byte(cabecalho), 0o644); err != nil {
			return err
		}
	}
	linha := fmt.Sprintf("| `%s` | %s | %s |\n", filepath.ToSlash(alvo), status, fundamento)
	arquivo, err := os.OpenFile(registro, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer arquivo.Close()
	_, err = arquivo.WriteString(linha)
	return err
}
