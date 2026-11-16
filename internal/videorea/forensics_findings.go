package videorea

import (
	"regexp"
	"strings"
)

var (
	conceptPattern = regexp.MustCompile(`(?i)conceito:\s*([^.\n]+)`)
	processPattern = regexp.MustCompile(`(?i)processo:\s*([^.\n]+)`)
	repoPattern    = regexp.MustCompile(`(?i)github\.com/([A-Za-z0-9_-]+/[A-Za-z0-9_-]+(?:\.[A-Za-z0-9_-]+)*)`)
)

// findLabeled devolve os valores únicos que seguem o rótulo do padrão.
func findLabeled(re *regexp.Regexp, text string) []string {
	var out []string
	for _, m := range re.FindAllStringSubmatch(text, -1) {
		out = appendUnique(out, strings.TrimSpace(m[1]))
	}
	return out
}

// findRepositories devolve os repositórios GitHub citados no formato github.com/dono/nome.
func findRepositories(text string) []string {
	var out []string
	for _, m := range repoPattern.FindAllStringSubmatch(text, -1) {
		out = appendUnique(out, "github.com/"+strings.TrimSuffix(m[1], ".git"))
	}
	return out
}

func appendUnique(list []string, item string) []string {
	if item == "" {
		return list
	}
	for _, have := range list {
		if have == item {
			return list
		}
	}
	return append(list, item)
}

// applyBrainFindings lê as linhas CONCEITO, PROCESSO e REPOSITORIO da resposta do cérebro.
func applyBrainFindings(d *ForensicsDossier, clean string) {
	switch {
	case strings.HasPrefix(clean, "CONCEITO:"):
		d.Concepts = appendUnique(d.Concepts, strings.TrimSpace(strings.TrimPrefix(clean, "CONCEITO:")))
	case strings.HasPrefix(clean, "PROCESSO:"):
		d.Processes = appendUnique(d.Processes, strings.TrimSpace(strings.TrimPrefix(clean, "PROCESSO:")))
	case strings.HasPrefix(clean, "REPOSITORIO:"):
		d.Repositories = appendUnique(d.Repositories, strings.TrimSpace(strings.TrimPrefix(clean, "REPOSITORIO:")))
	}
}
