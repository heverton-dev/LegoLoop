package orchestrix

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"legoloop/internal/dag"
)

const maxGateOutputBytes = 20000

var (
	rxViolationLine = regexp.MustCompile(`(^|[\\/])[a-zA-Z0-9_\-\.]+\.[a-zA-Z0-9]+:\d+(:\d+)?:`)
	rxExplicitCount = regexp.MustCompile(`(?i)(?:ENCONTRAD[AO]S?|DIAGN[OÓ]STICOS?|TOTAL)\s*\(\s*(\d+)\s*\)|limite\s*-max\s*\d+:\s*(\d+)`)
	rxGeneralCount  = regexp.MustCompile(`(?i)(\d+)\s+(?:viola[çc][õo]es|violac|violations?|errors?|failures?)`)
	rxFailedSteps   = regexp.MustCompile(`(?i)(\d+)\s+passo\(s\)\s+falharam`)
)

func isViolationLine(line string) bool {
	return rxViolationLine.MatchString(line)
}

func diagnosticCommand(n dag.Node) []string {
	if n.Verification == nil {
		return nil
	}
	if v, ok := n.Verification["diagnostic_command"].([]string); ok && len(v) > 0 {
		return v
	}
	if v, ok := n.Verification["diagnostic_command"].([]any); ok && len(v) > 0 {
		var out []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return diagnosticArgsFromVerification(n)
}

func diagnosticArgsFromVerification(n dag.Node) []string {
	if v, ok := n.Verification["diagnostic_args"].([]string); ok && len(n.Command) > 0 {
		return append(append([]string{}, n.Command...), v...)
	}
	if b, ok := n.Verification["diagnostic"].(bool); ok && b && len(n.Command) > 0 {
		return append(append([]string{}, n.Command...), "-diag")
	}
	if s, ok := n.Verification["diagnostic"].(string); ok && s != "" {
		return strings.Fields(s)
	}
	if b, ok := n.Verification["diagnostic_mode"].(bool); ok && b && len(n.Command) > 0 {
		return append(append([]string{}, n.Command...), "-diag")
	}
	return nil
}

func formatGateOutput(out []byte, maxBytes int) string {
	clean := strings.TrimSpace(string(out))
	if len(clean) <= maxBytes {
		return clean
	}
	scanner := bufio.NewScanner(strings.NewReader(clean))
	var headers, violations []string
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if isViolationLine(line) {
			violations = append(violations, line)
		} else if isHeaderLine(trimmed) {
			headers = append(headers, line)
		}
	}
	return assemblePrioritizedOutput(headers, violations, clean, maxBytes)
}

func isHeaderLine(trimmed string) bool {
	return strings.HasPrefix(trimmed, "[FALHA]") ||
		strings.HasPrefix(trimmed, "=== VIOLACOES") ||
		strings.HasPrefix(trimmed, "ticketgate:") ||
		strings.HasPrefix(trimmed, "--- FAIL:")
}

func assemblePrioritizedOutput(headers, violations []string, full string, maxBytes int) string {
	if len(violations) == 0 {
		if len(full) > maxBytes {
			return "..." + full[len(full)-maxBytes:]
		}
		return full
	}
	var b strings.Builder
	for _, h := range headers {
		b.WriteString(h + "\n")
	}
	omitted := 0
	for _, v := range violations {
		line := v + "\n"
		if b.Len()+len(line) > maxBytes-60 {
			omitted++
			continue
		}
		b.WriteString(line)
	}
	if omitted > 0 {
		fmt.Fprintf(&b, "... (%d violações adicionais omitidas)\n", omitted)
	}
	return strings.TrimSpace(b.String())
}

type gateScore struct {
	violations    int
	hasViolations bool
	failedSteps   int
	exitCode      int
}

func extractGateScore(out []byte, code int) gateScore {
	text := string(out)
	score := gateScore{exitCode: code}
	if m := rxExplicitCount.FindStringSubmatch(text); len(m) > 1 {
		val := m[1]
		if val == "" && len(m) > 2 {
			val = m[2]
		}
		if n, err := strconv.Atoi(val); err == nil {
			score.violations, score.hasViolations = n, true
		}
	} else if m := rxGeneralCount.FindStringSubmatch(text); len(m) > 1 {
		if n, err := strconv.Atoi(m[1]); err == nil {
			score.violations, score.hasViolations = n, true
		}
	}
	if !score.hasViolations {
		count := countViolationLines(text)
		if count > 0 {
			score.violations, score.hasViolations = count, true
		}
	}
	score.failedSteps = extractFailedSteps(text, code)
	return score
}

func countViolationLines(text string) int {
	scanner := bufio.NewScanner(strings.NewReader(text))
	count := 0
	for scanner.Scan() {
		if isViolationLine(scanner.Text()) {
			count++
		}
	}
	return count
}

func extractFailedSteps(text string, code int) int {
	if m := rxFailedSteps.FindStringSubmatch(text); len(m) > 1 {
		if n, err := strconv.Atoi(m[1]); err == nil {
			return n
		}
	}
	scanner := bufio.NewScanner(strings.NewReader(text))
	failed := 0
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(trimmed, "[FALHA]") || strings.HasPrefix(trimmed, "--- FAIL:") {
			failed++
		}
	}
	if failed == 0 && code != 0 {
		return 1
	}
	return failed
}

func isGateImprovement(initial, final gateScore) bool {
	if initial.hasViolations && final.hasViolations {
		return final.violations < initial.violations
	}
	if initial.failedSteps > 0 && final.failedSteps >= 0 {
		return final.failedSteps < initial.failedSteps
	}
	return false
}

func checkBuild(ctx context.Context, dir string) error {
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); os.IsNotExist(err) {
		return nil
	}
	cmd := exec.CommandContext(ctx, "go", "build", "./...")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &buf, &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build falhou: %w: %s", err, strings.TrimSpace(buf.String()))
	}
	return nil
}
