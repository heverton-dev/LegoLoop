package atualizar

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/escrita"
	"github.com/heverton-dev/LegoLoop/internal/forja"
)

const permBinario os.FileMode = 0o755

func substituirERevalidar(op Opcao, dadosNovo []byte, artefato, shaNovo string) (Resultado, error) {
	prevPath, shaAnterior, err := substituirComPrev(op.Binario, dadosNovo)
	if err != nil {
		return Resultado{}, err
	}
	if err := revalidarForja(op.Home); err != nil {
		return Resultado{}, err
	}
	return Resultado{
		Canal:          op.Canal,
		Artefato:       artefato,
		Binario:        op.Binario,
		BinarioPrev:    prevPath,
		SHA256Anterior: shaAnterior,
		SHA256Novo:     shaNovo,
		Revalidado:     true,
	}, nil
}

func substituirComPrev(binario string, dadosNovo []byte) (string, string, error) {
	if err := os.MkdirAll(filepath.Dir(binario), 0o755); err != nil {
		return "", "", fmt.Errorf("atualizar: criar pasta do binario: %w", err)
	}
	prevPath := binario + ".prev"
	shaAnterior, perm, err := guardarPrev(binario, prevPath)
	if err != nil {
		return "", "", err
	}
	if err := escrita.GravarAtomicoComPerm(binario, dadosNovo, perm); err != nil {
		return "", "", fmt.Errorf("atualizar: gravar novo binario: %w", err)
	}
	return prevPath, shaAnterior, nil
}

func guardarPrev(binario, prevPath string) (string, os.FileMode, error) {
	info, err := os.Stat(binario)
	if err != nil {
		if os.IsNotExist(err) {
			return "", permBinario, nil
		}
		return "", permBinario, fmt.Errorf("atualizar: inspecionar binario atual: %w", err)
	}
	dadosAtual, err := os.ReadFile(binario)
	if err != nil {
		return "", permBinario, fmt.Errorf("atualizar: ler binario atual: %w", err)
	}
	perm := info.Mode().Perm()
	if perm == 0 {
		perm = permBinario
	}
	if err := escrita.GravarAtomicoComPerm(prevPath, dadosAtual, perm); err != nil {
		return "", perm, fmt.Errorf("atualizar: guardar binario anterior em .prev: %w", err)
	}
	return calcularSHA256(dadosAtual), perm, nil
}

func revalidarForja(home string) error {
	derivas, err := forja.DetectarDeriva(home)
	if err != nil {
		return fmt.Errorf("atualizar: revalidar Forja: %w", err)
	}
	if len(derivas) > 0 {
		return fmt.Errorf("%w: %s", ErrDerivaDetectada, forja.FormatarDerivas(derivas))
	}
	return nil
}
