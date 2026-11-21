package espinha

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const prefixoProcesso = "legoloop-"

// EntradaRenomeado representa um item de espinha/renomeados.json.
type EntradaRenomeado struct {
	Antigo string `json:"antigo"`
	Novo   string `json:"novo"`
	Data   string `json:"data"`
}

// CarregarRenomeados le o registro de Processos renomeados.
func CarregarRenomeados(caminho string) ([]EntradaRenomeado, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil, fmt.Errorf("ler %s: %w", caminho, err)
	}
	var lista []EntradaRenomeado
	if err := json.Unmarshal(dados, &lista); err != nil {
		return nil, fmt.Errorf("decodificar %s: %w", caminho, err)
	}
	return lista, nil
}

// VerificarPrefixoProcessos aplica a regra processo-sem-prefixo.
func VerificarPrefixoProcessos(dirProcessos string) ([]string, error) {
	entradas, err := os.ReadDir(dirProcessos)
	if err != nil {
		return nil, fmt.Errorf("ler %s: %w", dirProcessos, err)
	}
	var violacoes []string
	for _, e := range entradas {
		if !e.IsDir() {
			continue
		}
		if !strings.HasPrefix(e.Name(), prefixoProcesso) {
			violacoes = append(violacoes, e.Name())
		}
	}
	return violacoes, nil
}

// BuscarReferenciasAntigas varre alvos ativos por nomes antigos de Processo.
func BuscarReferenciasAntigas(raiz string, lista []EntradaRenomeado) ([]string, error) {
	padroes := compilarPadroes(lista)
	var achados []string
	for _, alvo := range alvosVarredura(raiz) {
		info, err := os.Stat(alvo)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			achados = append(achados, varrerArquivo(raiz, alvo, padroes)...)
			continue
		}
		err = filepath.WalkDir(alvo, func(p string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if d.IsDir() {
				if ignorarDir(d.Name()) {
					return filepath.SkipDir
				}
				return nil
			}
			if ignorarArquivo(raiz, p) {
				return nil
			}
			achados = append(achados, varrerArquivo(raiz, p, padroes)...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return achados, nil
}

type padraoNome struct {
	antigo string
	re     *regexp.Regexp
}

func compilarPadroes(lista []EntradaRenomeado) []padraoNome {
	out := make([]padraoNome, 0, len(lista))
	for _, item := range lista {
		q := regexp.QuoteMeta(item.Antigo)
		expr := fmt.Sprintf(
			`(?i)(?:processos[/\\]%s(?:[^a-z0-9_-]|$)|"(?:name|slug|processo)"\s*:\s*"%s"|legoloop\s+(?:exec|run|rodar)\s+%s(?:[^a-z0-9_-]|$))`,
			q, q, q,
		)
		out = append(out, padraoNome{antigo: item.Antigo, re: regexp.MustCompile(expr)})
	}
	return out
}

func alvosVarredura(raiz string) []string {
	relativos := []string{
		"processos",
		"internal",
		"cmd",
		"tools",
		"skills",
		filepath.Join("docs", "PROCESSOS_CANONIZADOS.md"),
		filepath.Join("docs", "manuals"),
		"README.md",
		"AGENTS.md",
	}
	out := make([]string, 0, len(relativos))
	for _, r := range relativos {
		out = append(out, filepath.Join(raiz, r))
	}
	return out
}

func ignorarDir(nome string) bool {
	switch nome {
	case ".git", ".backups", "bin", "node_modules", ".venv", "venv", "__pycache__", "dist", "vendor":
		return true
	}
	return false
}

func ignorarArquivo(raiz, caminho string) bool {
	rel, err := filepath.Rel(raiz, caminho)
	if err != nil {
		return false
	}
	norm := filepath.ToSlash(rel)
	if norm == "espinha/renomeados.json" {
		return true
	}
	if strings.HasPrefix(norm, "internal/espinha/") {
		return true
	}
	ext := strings.ToLower(filepath.Ext(norm))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".ico", ".exe", ".dll", ".zip", ".pdf", ".mp4", ".wav":
		return true
	}
	return false
}

func varrerArquivo(raiz, caminho string, padroes []padraoNome) []string {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil
	}
	rel, err := filepath.Rel(raiz, caminho)
	if err != nil {
		rel = caminho
	}
	linhas := strings.Split(string(dados), "\n")
	var achados []string
	for idx, linha := range linhas {
		for _, p := range padroes {
			if p.re.MatchString(linha) {
				achados = append(achados, fmt.Sprintf("%s:%d (%s): %s",
					filepath.ToSlash(rel), idx+1, p.antigo, strings.TrimSpace(linha)))
			}
		}
	}
	return achados
}
