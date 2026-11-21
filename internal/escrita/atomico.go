package escrita

import (
	"fmt"
	"os"
	"path/filepath"
)

const permPadrao os.FileMode = 0o644

// GravarAtomico grava dados em caminho usando arquivo temporário no mesmo diretório,
// Sync, Close e os.Rename com permissão padrão 0o644.
func GravarAtomico(caminho string, dados []byte) error {
	return GravarAtomicoComPerm(caminho, dados, permPadrao)
}

// GravarAtomicoComPerm grava dados de forma atômica em caminho aplicando a permissão informada.
func GravarAtomicoComPerm(caminho string, dados []byte, perm os.FileMode) error {
	return gravarAtomicoComGancho(caminho, dados, perm, nil)
}

// gravarAtomicoComGancho executa a escrita atômica e permite injetar gancho de teste antes do rename.
func gravarAtomicoComGancho(caminho string, dados []byte, perm os.FileMode, antesDeRenomear func(string) error) error {
	tmpPath, err := escreverTemporario(filepath.Dir(caminho), dados, perm)
	if err != nil {
		return err
	}
	if antesDeRenomear != nil {
		if err := antesDeRenomear(tmpPath); err != nil {
			return limparTemporario(tmpPath, err)
		}
	}
	if err := os.Rename(tmpPath, caminho); err != nil {
		return limparTemporario(tmpPath, fmt.Errorf("escrita: renomeando temporario para %s: %w", caminho, err))
	}
	return nil
}

func escreverTemporario(dir string, dados []byte, perm os.FileMode) (string, error) {
	tmp, err := os.CreateTemp(dir, ".legoloop-tmp-*")
	if err != nil {
		return "", fmt.Errorf("escrita: criando temporario em %s: %w", dir, err)
	}
	tmpPath := tmp.Name()
	if err := preencherESincronizar(tmp, dados, perm); err != nil {
		return "", limparTemporario(tmpPath, err)
	}
	return tmpPath, nil
}

func preencherESincronizar(tmp *os.File, dados []byte, perm os.FileMode) error {
	if _, err := tmp.Write(dados); err != nil {
		return fecharComErro(tmp, fmt.Errorf("escrita: gravando temporario: %w", err))
	}
	if err := tmp.Chmod(perm); err != nil {
		return fecharComErro(tmp, fmt.Errorf("escrita: ajustando permissao do temporario: %w", err))
	}
	if err := tmp.Sync(); err != nil {
		return fecharComErro(tmp, fmt.Errorf("escrita: sincronizando temporario: %w", err))
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("escrita: fechando temporario: %w", err)
	}
	return nil
}

func fecharComErro(f *os.File, causa error) error {
	if err := f.Close(); err != nil {
		return fmt.Errorf("%w (fechamento: %v)", causa, err)
	}
	return causa
}

func limparTemporario(tmpPath string, causa error) error {
	if err := os.Remove(tmpPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("%w (limpeza %s: %v)", causa, tmpPath, err)
	}
	return causa
}
