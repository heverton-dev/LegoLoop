package orchestrix

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"legoloop/internal/dag"
)

// allowedWrites devolve os caminhos que um unit_task pode alterar:
// output.artifacts (context_out) e o arquivo de teste do test_seam.
func allowedWrites(n dag.Node) []string {
	out := append([]string{}, n.ContextOut...)
	if f := n.TestSeam["file"]; f != "" {
		out = append(out, f)
	}
	return out
}

// normRel normaliza um caminho relativo com barras e recusa absoluto e "..".
func normRel(p string) (string, error) {
	slash := strings.ReplaceAll(p, `\`, "/")
	if slash == "" || path.IsAbs(slash) || filepath.IsAbs(p) || strings.Contains(slash, ":") {
		return "", fmt.Errorf("caminho %q não é relativo ao projeto", p)
	}
	clean := path.Clean(slash)
	if clean == ".." || strings.HasPrefix(clean, "../") || clean == "." {
		return "", fmt.Errorf("caminho %q sai do projeto", p)
	}
	return clean, nil
}

// checkBoundary aplica o write_boundary strict: todo caminho alterado precisa estar
// na lista de permitidos. Devolve todos os desvios juntos.
// Um caminho permitido terminando com barra ("dir/") ou sem barra representa diretório
// se os arquivos alterados estiverem sob ele.
func checkBoundary(changed, allowed []string) error {
	exact := map[string]bool{}
	var prefixes []string
	for _, a := range allowed {
		n, err := normRel(a)
		if err != nil {
			return fmt.Errorf("saída permitida inválida: %w", err)
		}
		exact[n] = true
		prefixes = append(prefixes, n+"/")
	}
	var errs []error
	for _, c := range changed {
		n, err := normRel(c)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !exact[n] && !hasDirPrefix(n, prefixes) {
			errs = append(errs, fmt.Errorf("write_boundary: %q fora das saídas permitidas", n))
		}
	}
	return errors.Join(errs...)
}

func hasDirPrefix(p string, prefixes []string) bool {
	for _, pfx := range prefixes {
		if strings.HasPrefix(p, pfx) {
			return true
		}
	}
	return false
}

// buildUnitPrompt monta o prompt estéril de um unit_task: só os artefatos do ticket.
func buildUnitPrompt(n dag.Node, work *os.Root) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "# Orchestrix isolada — unit_task %s\n\n", n.ID)
	b.WriteString("Você executa um único ticket, sem histórico. Edite só os arquivos permitidos abaixo. ")
	b.WriteString("Não use sub-agentes. Não faça commit.\n\n")
	if n.SpecRef != "" {
		fmt.Fprintf(&b, "Critério de aceitação: %s\n", n.SpecRef)
	}
	if sig := n.TestSeam["signature"]; sig != "" {
		fmt.Fprintf(&b, "Teste-alvo (TDD Red-Green): %s em %s\n", sig, n.TestSeam["file"])
	}
	fmt.Fprintf(&b, "Condição de saída (gate): `%s`\n", n.Output.Gate)
	if len(n.Command) > 0 {
		fmt.Fprintf(&b, "Verificação: `%s`\n", strings.Join(n.Command, " "))
		b.WriteString("Você trabalha sem shell: não rode a verificação. A Orchestrix roda a verificação depois da sua edição. ")
		b.WriteString("Edite os arquivos permitidos até a verificação passar. Não encerre sem editar.\n")
	}
	b.WriteString("\n## Higiene de Código (lauren-unslop-ship)\n\n")
	b.WriteString("- Remova comentários narrativos supérfluos, especialmente aqueles que contenham termos proibidos pela taxonomia. A regra canônica é: remova o comentário narrativo em vez de renomear o termo ou mantê-lo.\n")
	b.WriteString("- Corrija todas as violações listadas, inclusive em comentários e em arquivos _test.go, e releia a lista no fim.\n")
	b.WriteString("- Não adicione comentários novos desnecessários (unslop, no-comments).\n")
	b.WriteString("- Mantenha o código estritamente formatado com gofmt.\n")
	b.WriteString("\n## Arquivos permitidos para escrita\n\n")
	allowed := allowedWrites(n)
	sort.Strings(allowed)
	for _, a := range allowed {
		fmt.Fprintf(&b, "- %s\n", a)
	}
	return appendContext(&b, n, work)
}

// contextIsDir informa se o context_in é uma pasta. Recusa ".git" e link ou junction no caminho.
func contextIsDir(work *os.Root, rel string) (bool, error) {
	if first, _, _ := strings.Cut(rel, "/"); strings.EqualFold(first, ".git") {
		return false, errors.New("caminho dentro de .git")
	}
	if err := checkNoLinks(work, rel); err != nil {
		return false, err
	}
	info, err := work.Lstat(filepath.FromSlash(rel))
	if err != nil {
		return false, nil
	}
	return info.IsDir(), nil
}

// appendContext injeta cada arquivo de context_in. Arquivo ausente (ainda a criar) é marcado.
func appendContext(b *strings.Builder, n dag.Node, work *os.Root) (string, error) {
	for _, f := range n.ContextIn {
		rel, err := normRel(f)
		if err != nil {
			return "", fmt.Errorf("context_in: %w", err)
		}
		if dir, err := contextIsDir(work, rel); err != nil {
			return "", fmt.Errorf("context_in %q: %w", rel, err)
		} else if dir {
			fmt.Fprintf(b, "\n## Contexto: pasta\n\nPasta de contexto: %s (explore com Glob, Grep e Read)\n", rel)
			continue
		}
		data, err := readLimited(work, filepath.FromSlash(rel))
		switch {
		case errors.Is(err, os.ErrNotExist):
			fmt.Fprintf(b, "\n## Contexto: %s\n\n(arquivo ainda não existe)\n", rel)
		case err != nil:
			return "", fmt.Errorf("context_in %q: %w", rel, err)
		default:
			fmt.Fprintf(b, "\n## Contexto: %s\n\n%s\n", rel, data)
		}
	}
	return b.String(), nil
}
