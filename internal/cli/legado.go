package cli

// Este arquivo implementa o caminho de falha do rename do acervo de processos
// (AC-E06.3): comando que aponta para o caminho legado `loops` recebe exit 2 e
// o guia de migração com o caminho novo `processos`. A menção ao caminho antigo
// aqui é o próprio guia: é a única razão do literal permanecer no código.

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

// padraoCaminhoLegado casa o segmento de pasta antigo num caminho informado.
var padraoCaminhoLegado = regexp.MustCompile(`(^|[/\\])loops([/\\]|$)`)

// caminhoLegado devolve o primeiro argumento de caminho que aponta para o
// segmento antigo. Flag sem valor embutido não é caminho e é ignorada.
func caminhoLegado(args []string) string {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			i := strings.Index(a, "=")
			if i < 0 {
				continue
			}
			a = a[i+1:]
		}
		if a != "" && padraoCaminhoLegado.MatchString(a) {
			return a
		}
	}
	return ""
}

// guiaMigracaoEscrever grava o erro claro com o caminho novo e devolve exit 2.
func guiaMigracaoEscrever(w io.Writer, alvo string) int {
	novo := padraoCaminhoLegado.ReplaceAllString(alvo, "${1}processos${2}")
	fmt.Fprintf(w, "caminho legado %q foi renomeado\n", alvo)
	fmt.Fprintf(w, "guia de migração: use o caminho novo %q\n", novo)
	fmt.Fprintf(w, "o acervo de processos agora vive em processos/; o caminho antigo não é mais aceito\n")
	return 2
}
