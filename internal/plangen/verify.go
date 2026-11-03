package plangen

import (
	"path"
	"path/filepath"
	"strings"

	"legoloop/internal/dag"
	"legoloop/internal/gate"
)

// shellMeta são os caracteres que só têm sentido para um shell. Os comandos do ROUTE são argv, sem shell.
const shellMeta = "|;<>$`\"'(){}*?&\\"

// splitVerify converte "cmd a b && cmd2 c" em argv. Devolve ok=false quando o comando depende de shell
// (pipe, redirecionamento, aspas, glob) ou não é um executável aceito pelo contrato do DAG.
func splitVerify(cmd string) (argv [][]string, ok bool) {
	for part := range strings.SplitSeq(cmd, "&&") {
		fields := strings.Fields(strings.Trim(strings.TrimSpace(part), "`"))
		if len(fields) == 0 || strings.ContainsAny(strings.Join(fields, ""), shellMeta) || dag.ValidateCommand(fields) != nil {
			return nil, false
		}
		argv = append(argv, fields)
	}
	return argv, len(argv) > 0
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
