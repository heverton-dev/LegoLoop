package distribuicao

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed scripts/capa.py
var scriptCapa []byte

// LimiteCapaBytes e o aviso de peso da capa; o limite oficial da rede ainda precisa de conferencia.
const LimiteCapaBytes = 2 << 20

// Executor roda um programa externo sem shell; injetavel nos testes.
type Executor interface {
	Executar(ctx context.Context, nome string, args ...string) error
}

// ExecSistema e o Executor real.
type ExecSistema struct{}

// Executar roda o programa e devolve a saida no erro.
func (ExecSistema) Executar(ctx context.Context, nome string, args ...string) error {
	saida, err := exec.CommandContext(ctx, nome, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w: %s", nome, err, strings.TrimSpace(string(saida)))
	}
	return nil
}

// PedidoCapa descreve uma capa: fundo sem texto (Flow ou codigo) e ate 4 palavras.
type PedidoCapa struct {
	Fundo        string `json:"fundo"`
	Linha1       string `json:"linha1"`
	Linha2       string `json:"linha2"`
	Saida        string `json:"saida"`
	Logo         string `json:"logo,omitempty"`
	Fonte        string `json:"fonte,omitempty"`
	Ouro         string `json:"ouro,omitempty"`
	Texto        string `json:"texto,omitempty"`
	TextoDireita bool   `json:"texto_direita,omitempty"`
	Curto        bool   `json:"curto,omitempty"`
}

// ResultadoCapa informa o tamanho gravado e se passou do peso de aviso.
type ResultadoCapa struct {
	Saida string `json:"saida"`
	Bytes int64  `json:"bytes"`
	Aviso string `json:"aviso,omitempty"`
}

// GerarCapa valida o pedido, roda o script embutido e confere o arquivo gravado.
func GerarCapa(ctx context.Context, ex Executor, projDir string, p PedidoCapa) (res ResultadoCapa, err error) {
	if err := ValidarTextoCapa(p.Linha1, p.Linha2, p.Curto); err != nil {
		return res, err
	}
	fundo, saida, err := caminhosCapa(projDir, p)
	if err != nil {
		return res, err
	}
	script, limpar, err := gravarScriptCapa()
	if err != nil {
		return res, err
	}
	defer func() { err = errors.Join(err, limpar()) }()
	if err := ex.Executar(ctx, "python", argsCapa(script, fundo, saida, p)...); err != nil {
		return res, err
	}
	info, err := os.Stat(saida)
	if err != nil {
		return res, fmt.Errorf("%w: o script nao gravou %s", ErrPedidoInvalido, p.Saida)
	}
	res = ResultadoCapa{Saida: p.Saida, Bytes: info.Size()}
	if info.Size() > LimiteCapaBytes {
		res.Aviso = "capa acima de 2 MB; confira o limite oficial da rede"
	}
	return res, nil
}

func caminhosCapa(projDir string, p PedidoCapa) (string, string, error) {
	for _, rel := range []string{p.Fundo, p.Saida} {
		if rel == "" || filepath.IsAbs(rel) || strings.Contains(filepath.ToSlash(filepath.Clean(rel)), "..") {
			return "", "", fmt.Errorf("%w: caminho %q fora do projeto", ErrPedidoInvalido, rel)
		}
	}
	fundo := filepath.Join(projDir, filepath.FromSlash(p.Fundo))
	if _, err := os.Stat(fundo); err != nil {
		return "", "", fmt.Errorf("%w: fundo %q ausente", ErrPedidoInvalido, p.Fundo)
	}
	return fundo, filepath.Join(projDir, filepath.FromSlash(p.Saida)), nil
}

func gravarScriptCapa() (string, func() error, error) {
	dir, err := os.MkdirTemp("", "capa-*")
	if err != nil {
		return "", nil, err
	}
	caminho := filepath.Join(dir, "capa.py")
	if err := os.WriteFile(caminho, scriptCapa, 0o600); err != nil {
		return "", nil, errors.Join(err, os.RemoveAll(dir))
	}
	return caminho, func() error { return os.RemoveAll(dir) }, nil
}

func argsCapa(script, fundo, saida string, p PedidoCapa) []string {
	args := []string{script, "--fundo", fundo, "--linha1", p.Linha1, "--saida", saida}
	opcoes := []struct{ flag, valor string }{{"--linha2", p.Linha2}, {"--logo", p.Logo}, {"--fonte", p.Fonte}, {"--ouro", p.Ouro}, {"--texto", p.Texto}}
	for _, o := range opcoes {
		if o.valor != "" {
			args = append(args, o.flag, o.valor)
		}
	}
	if p.TextoDireita {
		args = append(args, "--texto-direita")
	}
	return args
}
