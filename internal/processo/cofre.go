package processo

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// NaturezaDesconhecida é a natureza dos processos gravados pela v1 (sem .processo.json).
const NaturezaDesconhecida = "desconhecida"

const (
	sufixoLoopV1    = ".loop.json"
	prefixoProcesso = "legoloop-"
)

// ProcessoCofre representa um processo descoberto na pasta loops/ do cofre.
type ProcessoCofre struct {
	Nome     string
	Versao   string
	Natureza string
	Origem   string // "v2" ou "v1"
	Erro     string
	Etapas   []Etapa
}

// ProcessosDoCofre lê loops/ do cofre e devolve os processos ordenados por nome.
func ProcessosDoCofre(vault string) ([]ProcessoCofre, error) {
	entradas, err := os.ReadDir(filepath.Join(vault, "loops"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []ProcessoCofre{}, nil
		}
		return nil, err
	}
	nomes := nomesProcessos(entradas)
	out := make([]ProcessoCofre, 0, len(nomes))
	for _, nome := range nomes {
		out = append(out, lerProcessoCofre(vault, nome))
	}
	return out, nil
}

func nomesProcessos(entradas []os.DirEntry) []string {
	vistos := map[string]bool{}
	for _, e := range entradas {
		n := e.Name()
		if e.IsDir() || !strings.HasPrefix(n, prefixoProcesso) {
			continue
		}
		for _, suf := range []string{sufixoProcesso, sufixoLoopV1} {
			if base, ok := strings.CutSuffix(n, suf); ok && base != strings.TrimSuffix(prefixoProcesso, "-") {
				vistos[base] = true
			}
		}
	}
	out := make([]string, 0, len(vistos))
	for n := range vistos {
		out = append(out, n)
	}
	slices.Sort(out)
	return out
}

func lerProcessoCofre(vault, nome string) ProcessoCofre {
	proc := ProcessoCofre{Nome: nome, Origem: "v2"}
	caminho := CaminhoNoCofre(vault, nome)
	if _, err := os.Stat(caminho); err == nil {
		p, err := Ler(caminho)
		if err != nil {
			proc.Erro = err.Error()
			return proc
		}
		proc.Versao, proc.Natureza, proc.Etapas = p.Versao, p.Natureza, p.Etapas
		return proc
	}
	proc.Origem, proc.Natureza = "v1", NaturezaDesconhecida
	proc.Versao = versaoLoopV1(filepath.Join(vault, "loops", nome+sufixoLoopV1))
	return proc
}

func versaoLoopV1(path string) string {
	var lp struct {
		Version string `json:"version"`
	}
	if b, err := os.ReadFile(path); err == nil && json.Unmarshal(b, &lp) == nil && lp.Version != "" {
		return lp.Version
	}
	return "?"
}

// MapearUsadoPor devolve um mapa ordenado de quem usa cada processo do cofre.
func MapearUsadoPor(procs []ProcessoCofre) map[string][]string {
	compostos := make(map[string][]Etapa, len(procs))
	for _, p := range procs {
		compostos[p.Nome] = p.Etapas
	}
	return UsadoPorEtapas(compostos)
}

// UsadoPorEtapas calcula quem compõe quem a partir de um mapa de etapas por processo.
func UsadoPorEtapas(compostos map[string][]Etapa) map[string][]string {
	usadoPor := map[string][]string{}
	for nome, etapas := range compostos {
		for _, e := range etapas {
			alvo, _, ok := strings.Cut(e.Compoe, "@")
			if ok && !slices.Contains(usadoPor[alvo], nome) {
				usadoPor[alvo] = append(usadoPor[alvo], nome)
			}
		}
	}
	for _, quem := range usadoPor {
		slices.Sort(quem)
	}
	return usadoPor
}
