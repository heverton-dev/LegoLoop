package plangen

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/gate"
)

// shellMeta são os caracteres que só têm sentido para um shell. Os comandos do ROUTE são argv, sem shell.
const shellMeta = "|;<>$`\"'(){}*?&\\"

// splitVerify converte "cmd a b && cmd2 c" em argv. Devolve ok=false quando o comando depende de shell
// (pipe, redirecionamento, aspas simples, glob) ou não é um executável aceito pelo contrato do DAG.
// Argumento entre aspas duplas com regex simples (go test -run "A|B") é um argv único.
func splitVerify(cmd string) (argv [][]string, ok bool) {
	for part := range strings.SplitSeq(cmd, "&&") {
		fields, okF := fieldsComAspas(strings.Trim(strings.TrimSpace(part), "`"))
		if !okF || len(fields) == 0 || dag.ValidateCommand(fields) != nil {
			return nil, false
		}
		argv = append(argv, fields)
	}
	return argv, len(argv) > 0
}

// fieldsComAspas separa por espaço e agrupa aspas duplas. Fora das aspas vale shellMeta;
// dentro, só letras, dígitos e _ | . ^ / -. Aspas sem fechar reprovam.
func fieldsComAspas(s string) ([]string, bool) {
	var out []string
	for s = strings.TrimSpace(s); s != ""; s = strings.TrimSpace(s) {
		if s[0] == '"' {
			fim := strings.IndexByte(s[1:], '"')
			if fim < 0 || !regexSeguro(s[1:1+fim]) || (len(s) > fim+2 && s[fim+2] != 32) {
				return nil, false
			}
			out, s = append(out, s[1:1+fim]), s[fim+2:]
			continue
		}
		tok, resto, _ := strings.Cut(s, " ")
		if strings.ContainsAny(tok, shellMeta) {
			return nil, false
		}
		out, s = append(out, tok), resto
	}
	return out, true
}

func regexSeguro(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune("_|.^/-", r) && (r < '0' || r > '9') && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return false
		}
	}
	return s != ""
}

// gateSafe diz se o caminho pode ir como literal em file_exists(...).
func gateSafe(p string) bool {
	p = filepath.ToSlash(p)
	return p != "" && !path.IsAbs(p) && filepath.IsLocal(filepath.FromSlash(p)) && !strings.ContainsAny(p, "':\\")
}

// artifactGate prova mecanicamente: exit code 0 e todos os artefatos presentes.
// Se a expressão composta não compilar, cai para a forma mínima, que sempre compila.
func artifactGate(artifacts []string) string {
	var sb strings.Builder
	sb.WriteString("exit_code == 0")
	for _, a := range artifacts {
		if gateSafe(a) {
			sb.WriteString(" && file_exists('" + filepath.ToSlash(a) + "')")
		}
	}
	expr := sb.String()
	if _, err := gate.Compile(expr); err != nil {
		return "exit_code == 0"
	}
	return expr
}

// manualGate é o fallback honesto para comandos que não são argv: pede confirmação humana explícita.
func manualGate(cmd string) string {
	clean := strings.NewReplacer("'", "\"", "\n", " ", "\r", " ").Replace(strings.TrimSpace(cmd))
	return "prompt('Execute manualmente e confirme exit code 0: " + clean + "')"
}

// verifyAndGate devolve os argv de verificação e o gate DSL de um ticket.
func verifyAndGate(cmd string, artifacts []string) ([][]string, string) {
	if argv, ok := splitVerify(cmd); ok {
		return argv, artifactGate(artifacts)
	}
	return [][]string{}, manualGate(cmd)
}
