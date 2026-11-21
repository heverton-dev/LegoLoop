// Comando cobertura verifica se os pacotes listados atingem a cobertura mínima de testes.
//
// Uso: go run ./tools/cobertura [--limites espinha/cobertura.json]
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var reCover = regexp.MustCompile(`^(?:ok|\?)\s+(\S+)\s+.*coverage:\s+([0-9]+(?:\.[0-9]+)?)%\s+of\s+statements`)
var reSemTestes = regexp.MustCompile(`^\?\s+(\S+)\s+\[no test files\]`)

type resultadoPacote struct {
	Pacote string
	Obtida float64
	Minima float64
	Abaixo bool
}

func main() {
	os.Exit(executar(os.Args[1:], os.Stdout, os.Stderr, "."))
}

func executar(args []string, stdout, stderr io.Writer, dir string) int {
	fs := flag.NewFlagSet("cobertura", flag.ContinueOnError)
	fs.SetOutput(stderr)
	arqLimites := fs.String("limites", filepath.Join("espinha", "cobertura.json"), "arquivo JSON com limites mínimos por pacote")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	caminho := *arqLimites
	if !filepath.IsAbs(caminho) {
		caminho = filepath.Join(dir, caminho)
	}
	limites, err := lerLimites(caminho)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao ler limites: %v\n", err)
		return 2
	}
	pacotes := chavesOrdenadas(limites)
	saida, err := rodarGoTest(dir, pacotes)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao executar go test: %v\n%s\n", err, saida)
		return 2
	}
	modulo := lerNomeModulo(dir)
	obtidas := extrairCoberturas(saida)
	resultados, falhou := avaliar(pacotes, limites, obtidas, modulo)
	imprimirResultados(stdout, resultados)
	if falhou {
		return 2
	}
	return 0
}

func lerLimites(caminho string) (map[string]float64, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil, err
	}
	var limites map[string]float64
	if err := json.Unmarshal(dados, &limites); err != nil {
		return nil, err
	}
	if len(limites) == 0 {
		return nil, fmt.Errorf("arquivo de limites vazio: %s", caminho)
	}
	return limites, nil
}

func chavesOrdenadas(m map[string]float64) []string {
	chaves := make([]string, 0, len(m))
	for k := range m {
		chaves = append(chaves, k)
	}
	sort.Strings(chaves)
	return chaves
}

func rodarGoTest(dir string, pacotes []string) (string, error) {
	cmdArgs := append([]string{"test", "-count=1", "-cover"}, pacotes...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func lerNomeModulo(dir string) string {
	atual, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}
	for {
		f, err := os.Open(filepath.Join(atual, "go.mod"))
		if err == nil {
			defer f.Close()
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				linha := strings.TrimSpace(sc.Text())
				if strings.HasPrefix(linha, "module ") {
					return strings.TrimSpace(strings.TrimPrefix(linha, "module "))
				}
			}
			return ""
		}
		pai := filepath.Dir(atual)
		if pai == atual {
			return ""
		}
		atual = pai
	}
}

func extrairCoberturas(saida string) map[string]float64 {
	res := make(map[string]float64)
	sc := bufio.NewScanner(strings.NewReader(saida))
	for sc.Scan() {
		linha := strings.TrimSpace(sc.Text())
		if m := reCover.FindStringSubmatch(linha); m != nil {
			if pct, err := strconv.ParseFloat(m[2], 64); err == nil {
				res[m[1]] = pct
			}
			continue
		}
		if m := reSemTestes.FindStringSubmatch(linha); m != nil {
			res[m[1]] = 0
		}
	}
	return res
}

func resolverCobertura(pkg, modulo string, obtidas map[string]float64) (float64, bool) {
	if v, ok := obtidas[pkg]; ok {
		return v, true
	}
	rel := strings.TrimPrefix(strings.ReplaceAll(pkg, "\\", "/"), "./")
	if modulo != "" {
		chave := modulo + "/" + rel
		if rel == "." || rel == "" {
			chave = modulo
		}
		if v, ok := obtidas[chave]; ok {
			return v, true
		}
	}
	for k, v := range obtidas {
		if k == rel || strings.HasSuffix(k, "/"+rel) {
			return v, true
		}
	}
	return 0, false
}

func avaliar(pacotes []string, limites, obtidas map[string]float64, modulo string) ([]resultadoPacote, bool) {
	resultados := make([]resultadoPacote, 0, len(pacotes))
	falhou := false
	for _, pkg := range pacotes {
		minima := limites[pkg]
		obtida, ok := resolverCobertura(pkg, modulo, obtidas)
		abaixo := !ok || obtida < minima
		if abaixo {
			falhou = true
		}
		resultados = append(resultados, resultadoPacote{
			Pacote: pkg,
			Obtida: obtida,
			Minima: minima,
			Abaixo: abaixo,
		})
	}
	return resultados, falhou
}

func imprimirResultados(w io.Writer, resultados []resultadoPacote) {
	for _, r := range resultados {
		status := "OK"
		if r.Abaixo {
			status = "ABAIXO"
		}
		fmt.Fprintf(w, "%-6s %-25s obtida=%.1f%% minima=%.1f%%\n", status, r.Pacote, r.Obtida, r.Minima)
	}
}
