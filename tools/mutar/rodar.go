package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ExecutorTestes abstrai a execução dos testes de um pacote.
type ExecutorTestes interface {
	Rodar(ctx context.Context, dirPacote string) (passou bool, saida string, err error)
}

// ExecutorMutante permite aplicar a mutação de forma isolada antes de executar os testes.
type ExecutorMutante interface {
	RodarMutante(ctx context.Context, m Mutante) (passou bool, saida string, err error)
}

// ExecutorGoTest executa `go test -count=1 -failfast` com isolamento via `-overlay`.
type ExecutorGoTest struct{}

// Rodar invoca `go test -count=1 -failfast .` no diretório do pacote.
func (ExecutorGoTest) Rodar(ctx context.Context, dirPacote string) (bool, string, error) {
	return rodarGoTest(ctx, dirPacote, "")
}

// RodarMutante aplica o código mutado em um arquivo temporário via `-overlay`
// sem modificar o arquivo original em disco.
func (ExecutorGoTest) RodarMutante(ctx context.Context, m Mutante) (bool, string, error) {
	if len(m.Codigo) == 0 {
		return rodarGoTest(ctx, m.Pacote, "")
	}
	tmpDir, err := os.MkdirTemp("", "mutar-overlay-*")
	if err != nil {
		return false, "", fmt.Errorf("criar diretorio temporario: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	overlayPath, err := gravarOverlay(tmpDir, m)
	if err != nil {
		return false, "", err
	}
	return rodarGoTest(ctx, m.Pacote, overlayPath)
}

type arquivoOverlay struct {
	Replace map[string]string `json:"Replace"`
}

func gravarOverlay(tmpDir string, m Mutante) (string, error) {
	origAbs, err := filepath.Abs(filepath.Join(m.Pacote, m.Arquivo))
	if err != nil {
		return "", fmt.Errorf("resolver caminho absoluto de %s: %w", m.Arquivo, err)
	}
	mutPath := filepath.Join(tmpDir, filepath.Base(m.Arquivo))
	if err := os.WriteFile(mutPath, m.Codigo, 0o600); err != nil {
		return "", fmt.Errorf("gravar arquivo mutado: %w", err)
	}
	dados, err := json.Marshal(arquivoOverlay{
		Replace: map[string]string{origAbs: mutPath},
	})
	if err != nil {
		return "", fmt.Errorf("serializar overlay: %w", err)
	}
	overlayPath := filepath.Join(tmpDir, "overlay.json")
	if err := os.WriteFile(overlayPath, dados, 0o600); err != nil {
		return "", fmt.Errorf("gravar overlay.json: %w", err)
	}
	return overlayPath, nil
}

func rodarGoTest(ctx context.Context, dirPacote, overlayPath string) (bool, string, error) {
	args := []string{"test", "-count=1", "-failfast"}
	if overlayPath != "" {
		args = append(args, "-overlay", overlayPath)
	}
	args = append(args, ".")
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dirPacote
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return false, string(out), ctx.Err()
	}
	if err != nil {
		return false, string(out), nil
	}
	return true, string(out), nil
}

func executarCLI(args []string, stdout, stderr io.Writer) int {
	return executar(args, stdout, stderr, ".", ExecutorGoTest{}, OperadoresPadrao())
}

func executar(args []string, stdout, stderr io.Writer, dir string, execTestes ExecutorTestes, operadores []Operador) int {
	cfg, err := parseArgs(args, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "erro de argumentos: %v\n", err)
		return 2
	}
	justificados, err := resolverJustificados(dir, cfg)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao carregar justificados: %v\n", err)
		return 2
	}
	mutantes, err := coletarMutantes(dir, cfg.Pacotes, operadores)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao gerar mutantes: %v\n", err)
		return 2
	}
	rel := Rodar(context.Background(), mutantes, justificados, execTestes, cfg.Orcamento)
	imprimirRelatorio(stdout, rel)
	if rel.Sobreviventes > 0 || rel.Timeouts > 0 {
		return 2
	}
	return 0
}

// Rodar avalia cada mutante usando o executor fornecido, respeitando orçamento e justificativas.
func Rodar(ctx context.Context, mutantes []Mutante, justificados []Justificado, execTestes ExecutorTestes, orcamento time.Duration) Relatorio {
	if orcamento > 0 {
		var cancelar context.CancelFunc
		ctx, cancelar = context.WithTimeout(ctx, orcamento)
		defer cancelar()
	}
	rel := Relatorio{Total: len(mutantes), Resultados: make([]Resultado, 0, len(mutantes))}
	if _, paralelo := execTestes.(ExecutorGoTest); paralelo && len(mutantes) > 1 {
		resultados := rodarEmParalelo(ctx, mutantes, justificados, execTestes, 6)
		for _, res := range resultados {
			acumularResultado(&rel, res)
		}
		return rel
	}
	for _, m := range mutantes {
		res := avaliarMutante(ctx, m, justificados, execTestes)
		acumularResultado(&rel, res)
	}
	return rel
}

func rodarEmParalelo(ctx context.Context, mutantes []Mutante, justificados []Justificado, execTestes ExecutorTestes, workers int) []Resultado {
	if workers > len(mutantes) {
		workers = len(mutantes)
	}
	resultados := make([]Resultado, len(mutantes))
	jobs := make(chan int, len(mutantes))
	for i := range mutantes {
		jobs <- i
	}
	close(jobs)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				resultados[idx] = avaliarMutante(ctx, mutantes[idx], justificados, execTestes)
			}
		}()
	}
	wg.Wait()
	return resultados
}

func avaliarMutante(ctx context.Context, m Mutante, justificados []Justificado, execTestes ExecutorTestes) Resultado {
	if j, ok := buscarJustificado(m, justificados); ok {
		return Resultado{Mutante: m, Status: StatusJustificado, Justificativa: j}
	}
	if ctx.Err() != nil {
		return Resultado{Mutante: m, Status: StatusTimeout, Saida: ctx.Err().Error()}
	}
	passou, saida, err := executarTesteMutante(ctx, m, execTestes)
	if err != nil && (errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
		return Resultado{Mutante: m, Status: StatusTimeout, Saida: saida}
	}
	if err != nil {
		return Resultado{Mutante: m, Status: StatusSobrevivente, Saida: err.Error()}
	}
	if passou {
		return Resultado{Mutante: m, Status: StatusSobrevivente, Saida: saida}
	}
	return Resultado{Mutante: m, Status: StatusMorto, Saida: saida}
}

func executarTesteMutante(ctx context.Context, m Mutante, execTestes ExecutorTestes) (bool, string, error) {
	if em, ok := execTestes.(ExecutorMutante); ok {
		return em.RodarMutante(ctx, m)
	}
	return execTestes.Rodar(ctx, m.Pacote)
}

func acumularResultado(rel *Relatorio, res Resultado) {
	rel.Resultados = append(rel.Resultados, res)
	switch res.Status {
	case StatusMorto:
		rel.Mortos++
	case StatusSobrevivente:
		rel.Sobreviventes++
	case StatusJustificado:
		rel.Justificados++
	case StatusTimeout:
		rel.Timeouts++
	}
}

func buscarJustificado(m Mutante, lista []Justificado) (*Justificado, bool) {
	for i := range lista {
		j := &lista[i]
		if strings.TrimSpace(j.Anticorpo) == "" {
			continue
		}
		if j.ID == m.ID {
			return j, true
		}
		if filepath.Clean(j.Pacote) == filepath.Clean(m.Pacote) &&
			j.Arquivo == m.Arquivo && j.Linha == m.Linha && j.Operador == m.Operador {
			return j, true
		}
	}
	return nil, false
}

func imprimirRelatorio(w io.Writer, rel Relatorio) {
	fmt.Fprintf(w, "mutacao: total=%d mortos=%d sobreviventes=%d justificados=%d timeouts=%d\n",
		rel.Total, rel.Mortos, rel.Sobreviventes, rel.Justificados, rel.Timeouts)
	for _, r := range rel.Resultados {
		if r.Status == StatusSobrevivente || r.Status == StatusTimeout {
			fmt.Fprintf(w, "%s: id=%s pacote=%s arquivo=%s linha=%d operador=%s\n",
				r.Status, r.Mutante.ID, r.Mutante.Pacote, r.Mutante.Arquivo, r.Mutante.Linha, r.Mutante.Operador)
		}
	}
}
