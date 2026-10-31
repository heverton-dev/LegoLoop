package provider

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// candidateEnvPaths retorna os locais mais prováveis de conter .env no SO.
func candidateEnvPaths(home string) []string {
	var candidates []string
	if home != "" {
		candidates = append(candidates,
			filepath.Join(home, "Desktop", "9router", ".env"),
			filepath.Join(home, ".aidd", ".env"),
			filepath.Join(home, ".env"),
		)
	}
	candidates = append(candidates, ".env")
	return candidates
}

// parseEnvFile lê um arquivo .env simples ignorando comentários e aspas.
func parseEnvFile(path string) map[string]string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	out := make(map[string]string)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(k) == "" {
			continue
		}
		key := strings.TrimSpace(k)
		val := strings.Trim(strings.TrimSpace(v), `"'`)
		if val != "" {
			out[key] = val
		}
	}
	return out
}

// findRouterCreds busca NINEROUTER_URL e chave no ambiente e depois nos arquivos .env.
func findRouterCreds(home string, getenv func(string) string) (string, string) {
	url := getenv("NINEROUTER_URL")
	key := routerKeyFromEnv(getenv)
	if url != "" && key != "" {
		return url, key
	}
	for _, p := range candidateEnvPaths(home) {
		m := parseEnvFile(p)
		if url == "" && m["NINEROUTER_URL"] != "" {
			url = m["NINEROUTER_URL"]
		}
		if key == "" {
			if k := m["NINEROUTER_API_KEY"]; k != "" {
				key = k
			} else if k := m["NINEROUTER_KEY"]; k != "" {
				key = k
			}
		}
		if url != "" && key != "" {
			break
		}
	}
	return url, key
}
