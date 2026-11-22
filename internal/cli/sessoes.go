package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/jail"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

const sessoesUsage = `uso: legoloop sessoes <slug> [--json] [--processos-dir DIR]

Lista as sessões ativas (PID vivo) do processo, com etapa atual e próximo
passo da memória de execução (memoria_execucao.json). Saída nunca expõe
caminho absoluto.
`

// SessaoAtiva é a visão enxuta de uma sessão viva para a saída da CLI.
type SessaoAtiva struct {
	SessaoID     string `json:"sessao_id"`
	PID          int    `json:"pid"`
	EtapaAtual   string `json:"etapa_atual"`
	ProximoPasso string `json:"proximo_passo"`
}

// cmdSessoes lista as sessões ativas de um processo (controle de progresso Lei 11).
func cmdSessoes(_ context.Context, args []string, d Deps) int {
	slug, resto := splitPositional(args)
	fs := newFlags("sessoes", d)
	asJSON := fs.Bool("json", false, "saída em json")
	dirBase := fs.String("processos-dir", "", "pasta que contém <slug>")
	if err := fs.Parse(resto); err != nil {
		return 2
	}
	if slug == "" || fs.NArg() != 0 {
		fmt.Fprintf(d.Stderr, "erro: informe o slug do processo.\n\n%s", sessoesUsage)
		return 2
	}
	if err := jail.ValidarNomeProcesso(slug); err != nil {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	dir := resolverDirProcesso(slug, *dirBase)
	if dir == "" {
		fmt.Fprintln(d.Stderr, "sessoes: processo não encontrado")
		return 1
	}
	return imprimirSessoes(triagem.ListarSessoesAtivas(dir), *asJSON, d)
}

// resolverDirProcesso acha a pasta do processo na base informada ou nos dialetos
// padrão sob o diretório corrente (novo e legado).
func resolverDirProcesso(slug, dirBase string) string {
	bases := []string{dirBase}
	if dirBase == "" {
		bases = caminhos.DirsLeitura(".")
	}
	for _, base := range bases {
		dir := filepath.Join(base, slug)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	return ""
}

// imprimirSessoes grava a lista de sessões ativas em JSON ou texto.
func imprimirSessoes(ativas []triagem.MemoriaExecucao, asJSON bool, d Deps) int {
	saidas := make([]SessaoAtiva, 0, len(ativas))
	for _, m := range ativas {
		saidas = append(saidas, SessaoAtiva{
			SessaoID:     m.SessaoID,
			PID:          m.PID,
			EtapaAtual:   m.EtapaAtual,
			ProximoPasso: m.ProximoPasso,
		})
	}
	if asJSON {
		return printJSON(d, saidas)
	}
	for _, s := range saidas {
		fmt.Fprintf(d.Stdout, "%s pid=%d etapa=%s\n", s.SessaoID, s.PID, s.EtapaAtual)
	}
	return 0
}
