// Comando skillgate emite o modo relatório do Gate de paridade Skill×DAG
// (Lei 11): percorre os Processos do repositório e lista quem não cita todos
// os nós do DAG no SKILL.md correspondente. Somente leitura.
//
// Uso: go run ./tools/skillgate [repoDir]
package main

import (
	"fmt"
	"os"
	"strings"

	"legoloop/internal/fechar"
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	for _, r := range fechar.RelatorioParidadeSkillDAG(dir) {
		fmt.Println(linhaRelato(r))
	}
}

func linhaRelato(r fechar.ParidadeRelato) string {
	if r.OK {
		return fmt.Sprintf("OK    %s | %s", r.Processo, r.Skill)
	}
	if r.Motivo != "" {
		return fmt.Sprintf("FALHA %s | %s", r.Processo, r.Motivo)
	}
	return fmt.Sprintf("FALHA %s | %s não cita os nós: %s", r.Processo, r.Skill, strings.Join(r.Faltando, "; "))
}
