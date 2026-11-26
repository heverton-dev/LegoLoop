package render

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// removerTemp apaga o temporario; arquivo que ja nao existe conta como sucesso.
func removerTemp(caminho string) error {
	if err := os.Remove(caminho); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remover temporario %s: %w", caminho, err)
	}
	return nil
}

// saidaArquivo e um arquivo a gravar, com caminho relativo a pasta do projeto.
type saidaArquivo struct {
	rel   string
	dados []byte
}

// gravarEmLote grava todos os arquivos em temporarios e so entao os renomeia. Se qualquer
// gravacao falhar, nenhum arquivo final muda e nenhum temporario sobra.
func gravarEmLote(projDir string, saidas []saidaArquivo) error {
	tmps := make([]string, len(saidas))
	finais := make([]string, len(saidas))
	limpar := func() error {
		var falhas []error
		for _, t := range tmps {
			if t != "" {
				falhas = append(falhas, removerTemp(t))
			}
		}
		return errors.Join(falhas...)
	}
	for i, s := range saidas {
		final, err := resolverSeguro(projDir, s.rel)
		if err != nil {
			return errors.Join(err, limpar())
		}
		if err := os.MkdirAll(filepath.Dir(final), 0o755); err != nil {
			return errors.Join(err, limpar())
		}
		tmp := final + ".tmp"
		if err := os.WriteFile(tmp, s.dados, 0o644); err != nil {
			return errors.Join(err, limpar())
		}
		tmps[i], finais[i] = tmp, final
	}
	for i := range tmps {
		if err := os.Rename(tmps[i], finais[i]); err != nil {
			return err
		}
	}
	return nil
}
