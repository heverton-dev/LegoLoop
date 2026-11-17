package inventario

import (
	"io/fs"
	"os"
	"path/filepath"
)

// ResultadoMigracao resume a movimentação para o arquivo morto.
type ResultadoMigracao struct {
	Movidos    []string // destinos relativos a docs/
	Protegidos []string // relatórios globais mantidos em docs/reports/
}

// Migrar move a documentação de construção de docs/ para docs/archive/.
// Move com rename. Nunca deleta arquivo (contrato AC-E01.2).
func Migrar(raiz string) (ResultadoMigracao, error) {
	var res ResultadoMigracao
	docs := filepath.Join(raiz, "docs")
	if err := os.MkdirAll(filepath.Join(docs, "archive"), 0o755); err != nil {
		return res, err
	}
	for _, dir := range []string{"videorea", "sessions"} {
		origem := filepath.Join(docs, dir)
		destino := filepath.Join(docs, "archive", dir)
		if err := moverDiretorio(origem, destino, docs, &res); err != nil {
			return res, err
		}
	}
	globais, err := RelatoriosGlobais(raiz)
	if err != nil {
		return res, err
	}
	return res, moverRelatorios(docs, globais, &res)
}

// moverDiretorio move todos os arquivos de origem preservando a estrutura.
func moverDiretorio(origem, destino, raizDocs string, res *ResultadoMigracao) error {
	if _, err := os.Stat(origem); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(origem, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, errRel := filepath.Rel(origem, p)
		if errRel != nil {
			return errRel
		}
		alvo := filepath.Join(destino, rel)
		if err := moverArquivo(p, alvo); err != nil {
			return err
		}
		registrarMovido(res, raizDocs, alvo)
		return nil
	})
}

// moverRelatorios move relatórios históricos e protege os relatórios globais.
func moverRelatorios(docs string, globais map[string]bool, res *ResultadoMigracao) error {
	origem := filepath.Join(docs, "reports")
	if _, err := os.Stat(origem); os.IsNotExist(err) {
		return nil
	}
	destino := filepath.Join(docs, "archive", "reports")
	return filepath.WalkDir(origem, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, errRel := filepath.Rel(origem, p)
		if errRel != nil {
			return errRel
		}
		if globais[filepath.Base(p)] {
			res.Protegidos = append(res.Protegidos, filepath.ToSlash(filepath.Join("reports", rel)))
			return nil
		}
		alvo := filepath.Join(destino, rel)
		if err := moverArquivo(p, alvo); err != nil {
			return err
		}
		registrarMovido(res, docs, alvo)
		return nil
	})
}

// moverArquivo move um arquivo com rename. Nunca deleta.
func moverArquivo(origem, destino string) error {
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	return os.Rename(origem, destino)
}

// registrarMovido registra o destino relativo a docs/.
func registrarMovido(res *ResultadoMigracao, raizDocs, destino string) {
	rel, err := filepath.Rel(raizDocs, destino)
	if err != nil {
		rel = destino
	}
	res.Movidos = append(res.Movidos, filepath.ToSlash(rel))
}
