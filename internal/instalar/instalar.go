// Package instalar orquestra a instalação do LegoLoop encadeando a Forja,
// gravando os relatórios de comprovação e configurando o atalho do grafo.
package instalar

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/escrita"
	"github.com/heverton-dev/LegoLoop/internal/forja"
	"github.com/heverton-dev/LegoLoop/internal/jail"
)

// Opcao configura a execução da instalação do LegoLoop.
type Opcao struct {
	Home         string                       `json:"home"`
	Perfil       string                       `json:"perfil"`
	Sem          []string                     `json:"sem,omitempty"`
	Aplicar      bool                         `json:"aplicar"`
	Data         string                       `json:"-"`
	SO           string                       `json:"-"`
	Versao       string                       `json:"-"`
	Harnesses    []string                     `json:"-"`
	Geradores    map[string]forja.GeradorFunc `json:"-"`
	DirBinarios  string                       `json:"-"`
	BinarioGrafo string                       `json:"-"`
	Ambiente     []string                     `json:"-"`
}

// Resultado descreve o desfecho da instalação e os caminhos produzidos.
type Resultado struct {
	Relatorio   forja.Relatorio `json:"relatorio"`
	CaminhoJSON string          `json:"caminho_json,omitempty"`
	CaminhoMD   string          `json:"caminho_md,omitempty"`
	AtalhoGrafo string          `json:"atalho_grafo,omitempty"`
	Aplicado    bool            `json:"aplicado"`
}

// Executar encadeia a Forja e, no modo aplicar, grava os relatórios atomicamente
// e cria o atalho do grafo quando o binário fonte existe.
func Executar(ctx context.Context, op Opcao) (Resultado, error) {
	op = normalizarOpcao(op)
	rel, err := forja.Executar(ctx, montarOpcaoForja(op))
	if err != nil {
		return Resultado{Relatorio: rel}, err
	}
	res := Resultado{Relatorio: rel, Aplicado: op.Aplicar}
	dirBin := resolverDirBinarios(op)
	binGrafo := resolverBinarioGrafo(op)
	atalho, err := CriarAtalhoGrafo(dirBin, binGrafo, op.SO, op.Aplicar)
	if err != nil {
		return res, err
	}
	res.AtalhoGrafo = atalho
	if !op.Aplicar {
		return res, nil
	}
	arqJSON, arqMD, err := gravarRelatorios(op, rel)
	if err != nil {
		return res, err
	}
	res.CaminhoJSON = arqJSON
	res.CaminhoMD = arqMD
	return res, nil
}

func normalizarOpcao(op Opcao) Opcao {
	if strings.TrimSpace(op.Perfil) == "" {
		op.Perfil = "padrao"
	}
	if strings.TrimSpace(op.SO) == "" {
		op.SO = runtime.GOOS
	}
	return op
}

func montarOpcaoForja(op Opcao) forja.Opcao {
	return forja.Opcao{
		Home:      op.Home,
		Perfil:    op.Perfil,
		Sem:       op.Sem,
		Aplicar:   op.Aplicar,
		Data:      op.Data,
		SO:        op.SO,
		Versao:    op.Versao,
		Harnesses: op.Harnesses,
		Geradores: op.Geradores,
	}
}

func gravarRelatorios(op Opcao, rel forja.Relatorio) (string, string, error) {
	dirRel, err := jail.DentroDaRaiz(op.Home, filepath.Join(".legoloop", "forja", rel.Data))
	if err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(dirRel, 0o755); err != nil {
		return "", "", err
	}
	segredos := coletarSegredos(op.Ambiente)
	dadosJSON, err := forja.RenderizarRelatorioJSON(rel, segredos)
	if err != nil {
		return "", "", err
	}
	dadosMD := []byte(forja.RenderizarRelatorioMD(rel, segredos))
	arqJSON := filepath.Join(dirRel, "RELATORIO.json")
	arqMD := filepath.Join(dirRel, "RELATORIO.md")
	if err := escrita.GravarAtomico(arqJSON, dadosJSON); err != nil {
		return "", "", err
	}
	if err := escrita.GravarAtomico(arqMD, dadosMD); err != nil {
		return "", "", err
	}
	return arqJSON, arqMD, nil
}

func coletarSegredos(env []string) []string {
	if env == nil {
		return forja.ColetarSegredosAmbiente(os.Environ())
	}
	return forja.ColetarSegredosAmbiente(env)
}
