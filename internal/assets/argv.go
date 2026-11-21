package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

const (
	maxCommandLen = 4096
	maxArgs       = 64
	// shellMeta são caracteres que um shell interpretaria. Não há shell aqui (o argv vai direto
	// ao processo), então um comando que dependa deles não pode ser representado com fidelidade.
	shellMeta = "|&;<>$`()%"
)

// SplitCommand converte uma linha de comando em argv, sem shell: separa por espaço, respeita
// aspas simples e duplas e trata a barra invertida como caractere comum (caminhos do Windows).
// Recusa o que um shell interpretaria (pipes, redirecionamentos, $VAR, %VAR%, crases...).
func SplitCommand(s string) ([]string, error) {
	if len(s) > maxCommandLen {
		return nil, fmt.Errorf("comando excede %d caracteres", maxCommandLen)
	}
	var sp splitter
	for _, r := range s {
		if err := sp.feed(r); err != nil {
			return nil, err
		}
	}
	if sp.quote != 0 {
		return nil, errors.New("aspas não fechadas")
	}
	sp.flush()
	if len(sp.args) == 0 {
		return nil, errors.New("comando vazio")
	}
	if len(sp.args) > maxArgs {
		return nil, fmt.Errorf("comando excede %d argumentos", maxArgs)
	}
	return sp.args, nil
}

// splitter é o autômato do SplitCommand.
type splitter struct {
	args    []string
	cur     strings.Builder
	inToken bool
	quote   rune
}

func (sp *splitter) flush() {
	if sp.inToken {
		sp.args = append(sp.args, sp.cur.String())
		sp.cur.Reset()
		sp.inToken = false
	}
}

func (sp *splitter) feed(r rune) error {
	switch {
	case r == '\n' || r == '\r' || r == 0:
		return errors.New("comando com quebra de linha ou byte nulo")
	case sp.quote != 0:
		return sp.feedQuoted(r)
	case r == '"' || r == '\'':
		sp.quote, sp.inToken = r, true
	case unicode.IsSpace(r):
		sp.flush()
	case strings.ContainsRune(shellMeta, r):
		return fmt.Errorf("o caractere %q exige shell; o runner executa sem shell", r)
	default:
		sp.cur.WriteRune(r)
		sp.inToken = true
	}
	return nil
}

func (sp *splitter) feedQuoted(r rune) error {
	switch {
	case r == sp.quote:
		sp.quote = 0
	case r == '$' || r == '`':
		return fmt.Errorf("%q dentro de aspas depende de shell", r)
	default:
		sp.cur.WriteRune(r)
	}
	return nil
}

// ArgvFor devolve o argv que executa um script ou hook do inventário. Roda no servidor: o
// navegador nunca vê nem envia caminhos absolutos nem o comando bruto de um hook.
func ArgvFor(it Item) ([]string, error) {
	return ArgvForOS(it, runtime.GOOS)
}

// ArgvForOS devolve o argv que executa um script ou hook resolvendo binários para o SO informado.
func ArgvForOS(it Item, goos string) ([]string, error) {
	switch it.Kind {
	case KindScript:
		return scriptArgvForOS(it, goos)
	case KindHook:
		if it.Command == "" {
			return nil, errors.New("hook sem comando utilizável")
		}
		return SplitCommand(it.Command)
	}
	return nil, fmt.Errorf("só scripts e hooks viram comando, recebeu %q", it.Kind)
}

// PythonForOS devolve o executável Python padrão do SO.
func PythonForOS(goos string) string {
	if goos == "windows" {
		return "python"
	}
	return "python3"
}

// PowerShellForOS devolve o executável PowerShell padrão do SO.
func PowerShellForOS(goos string) string {
	if goos == "windows" {
		return "powershell"
	}
	return "pwsh"
}

// ShellForOS devolve o shell e flag de comando padrão do SO.
func ShellForOS(goos string) []string {
	if goos == "windows" {
		return []string{"cmd", "/c"}
	}
	return []string{"sh", "-c"}
}

func scriptArgvForOS(it Item, goos string) ([]string, error) {
	if it.Path == "" {
		return nil, errors.New("script sem caminho")
	}
	switch strings.ToLower(filepath.Ext(it.Path)) {
	case ".ps1":
		return []string{PowerShellForOS(goos), "-NoProfile", "-File", it.Path}, nil
	case ".sh", ".bash":
		return []string{"bash", it.Path}, nil
	case ".py":
		return []string{PythonForOS(goos), it.Path}, nil
	case ".js", ".mjs", ".cjs":
		return []string{"node", it.Path}, nil
	case ".bat", ".cmd":
		if goos == "windows" {
			return []string{"cmd", "/c", it.Path}, nil
		}
		return []string{"sh", it.Path}, nil
	}
	return nil, fmt.Errorf("extensão %q sem interpretador conhecido", filepath.Ext(it.Path))
}

// itemID é um identificador opaco e estável do item (não revela o caminho).
func itemID(it Item) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{string(it.Kind), it.Name, it.Scope, it.Origin, it.Path, it.Description}, "\x00")))
	return hex.EncodeToString(sum[:])[:12]
}

// ByID procura um item pelo id opaco.
func (s Snapshot) ByID(id string) (Item, bool) {
	for _, it := range s.Items {
		if it.ID == id {
			return it, true
		}
	}
	return Item{}, false
}
