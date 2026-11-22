package progresso

import (
	"fmt"
	"path/filepath"
	"strings"
)

// CaminhoDoDAG valida o campo `progresso` de um DAG e o resolve a partir de base.
// O DAG pode vir de terceiros: só vale caminho local .json, sem `..` e sem raiz absoluta (Path Jail).
// Campo vazio devolve vazio e desliga o progresso vivo.
func CaminhoDoDAG(base, campo string) (string, error) {
	campo = strings.TrimSpace(campo)
	if campo == `` {
		return ``, nil
	}
	limpo := filepath.Clean(filepath.FromSlash(campo))
	if !filepath.IsLocal(limpo) || !strings.EqualFold(filepath.Ext(limpo), `.json`) {
		return ``, fmt.Errorf(`campo progresso inválido %q: use caminho local .json dentro do projeto`, campo)
	}
	return filepath.Join(base, limpo), nil
}
