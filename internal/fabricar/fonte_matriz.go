package fabricar

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/assets"
	"github.com/heverton-dev/LegoLoop/internal/matriz"
	"github.com/heverton-dev/LegoLoop/internal/processo"
	"github.com/heverton-dev/LegoLoop/internal/ste"
)

func init() {
	ste.Default.Register("fabricar.sem_matriz", "A matriz de peças não existe em %s. Rode legoloop matriz gerar e tente de novo.")
	ste.Default.Register("fabricar.peca_removivel", "A peça %s está em mídia removível ou de rede. Ela ficou fora. Use --trust-removable para aceitá-la.")
	ste.Default.Register("fabricar.peca_nome_inseguro", "Uma peça tem um nome que não é seguro para o grafo: %s. Ela ficou fora.")
}

// ErrSemMatriz: o cofre não tem registry_matrix.json. O fabricar só consulta a matriz; nunca a gera.
var ErrSemMatriz = errors.New("fabricar: matriz de peças ausente")

// FonteMatriz é o que a Camada 2 lê do cofre: as peças com perfil de uso, o vocabulário e a data da matriz.
type FonteMatriz struct {
	Entradas []matriz.Entrada
	Vocab    matriz.Vocab
	GeradaEm time.Time
}

// MatrizFunc entrega a matriz já lida; os testes a trocam por uma matriz em memória.
type MatrizFunc func() (FonteMatriz, error)

// AtualizarMatrizFunc atualiza a Matriz gravada no cofre (em modo -sem-cerebro por padrão).
type AtualizarMatrizFunc func(ctx context.Context) error

func deveVerificarDefasagem(d Deps) bool {
	return len(d.DiscoveryRoots) > 0 || d.Scan != nil || d.Atualizar != nil
}

func homeDe(d Deps) string {
	if d.Home != "" {
		return d.Home
	}
	h, _ := os.UserHomeDir()
	return h
}

func varrerPecas(ctx context.Context, req Request, d Deps) ([]assets.Item, error) {
	if d.Scan != nil {
		return d.Scan(ctx)
	}
	cfg := assets.Config{
		Home:             homeDe(d),
		AutoDiscover:     len(d.DiscoveryRoots) > 0 || len(req.TargetDrives) > 0,
		DiscoveryRoots:   d.DiscoveryRoots,
		TargetDrives:     req.TargetDrives,
		TrustRemovable:   req.TrustRemovable,
		ProvedorUnidades: d.ProvedorUnidades,
		LookPath:         func(string) (string, error) { return "", os.ErrNotExist },
		Getenv:           func(string) string { return "" },
	}
	return assets.NewScanner(cfg).Scan(ctx).Items, nil
}

func atualizarMatrizSemCerebro(ctx context.Context, d Deps, items []assets.Item) error {
	voc, err := matriz.LoadVocab(d.Roots.Vault)
	if err != nil {
		return err
	}
	md := matriz.Deps{
		Home:      homeDe(d),
		Catalogar: matriz.Catalogador{Vocab: voc, Brain: nil},
		Scan:      func(context.Context) ([]assets.Item, error) { return items, nil },
		Now:       d.Clock,
	}
	_, err = matriz.Generate(ctx, matriz.Options{Vault: d.Roots.Vault}, md)
	return err
}

// atualizarSeDefasada consulta matriz.Staleness antes da Fusão; se Defasada() == true, chama Atualizar uma vez
// e recarrega a Matriz para a Fusão usar as peças atualizadas.
func atualizarSeDefasada(ctx context.Context, req Request, d Deps, fonte FonteMatriz) (FonteMatriz, matriz.Stale, error) {
	if !deveVerificarDefasagem(d) {
		return fonte, matriz.Stale{}, nil
	}
	reg, err := matriz.LerRegistro(filepath.Join(d.Roots.Vault, "registry_matrix.json"))
	if errors.Is(err, os.ErrNotExist) {
		reg, err = nil, nil
	}
	if err != nil {
		return fonte, matriz.Stale{}, err
	}
	items, err := varrerPecas(ctx, req, d)
	if err != nil {
		return fonte, matriz.Stale{}, err
	}
	stale, err := matriz.Staleness(ctx, reg, items, matriz.Deps{Home: homeDe(d)})
	if err != nil || !stale.Defasada() {
		return fonte, stale, err
	}
	if d.Atualizar != nil {
		err = d.Atualizar(ctx)
	} else {
		err = atualizarMatrizSemCerebro(ctx, d, items)
	}
	if err != nil {
		return fonte, stale, err
	}
	recarregada, err := d.Matriz()
	return recarregada, stale, err
}

func pecasPendentes(entradas []matriz.Entrada) []string {
	var out []string
	for _, e := range entradas {
		if e.Retentar {
			out = append(out, e.Nome)
		}
	}
	slices.Sort(out)
	return out
}

// LerMatriz lê registry_matrix.json e vocab.json do cofre. A data é a de generated_at; sem ela, a do arquivo.
func LerMatriz(vault string) (FonteMatriz, error) {
	path := filepath.Join(vault, "registry_matrix.json")
	reg, err := matriz.LerRegistro(path)
	if errors.Is(err, os.ErrNotExist) {
		return FonteMatriz{}, fmt.Errorf("%w: %s", ErrSemMatriz, ste.Msg("fabricar.sem_matriz", path))
	}
	if err != nil {
		return FonteMatriz{}, err
	}
	voc, err := matriz.LoadVocab(vault)
	if err != nil {
		return FonteMatriz{}, err
	}
	return FonteMatriz{Entradas: reg.Entradas, Vocab: voc, GeradaEm: dataDaMatriz(reg.GeneratedAt, path)}, nil
}

// dataDaMatriz lê generated_at com hora (RFC 3339) ou só com o dia, como a matriz real grava ("2026-11-04").
func dataDaMatriz(generatedAt, path string) time.Time {
	for _, layout := range []string{time.RFC3339, time.DateOnly} {
		if t, err := time.Parse(layout, generatedAt); err == nil {
			return t
		}
	}
	if fi, err := os.Stat(path); err == nil {
		return fi.ModTime()
	}
	return time.Time{}
}

// confiaveis tira as skills que a Fusão não pode usar, cada uma com um aviso: nome que não serve de peça
// (sempre) e skill em mídia removível ou de rede (salvo trust). Só skills são julgadas: são as únicas peças
// que a Fusão usa (MCPs e hooks nunca viram passo).
func confiaveis(es []matriz.Entrada, trust bool) (kept []matriz.Entrada, cortes []string) {
	for _, e := range es {
		switch {
		case e.Tipo != "skill":
		case !processo.NomeDePecaValido(e.Nome):
			cortes = append(cortes, ste.Msg("fabricar.peca_nome_inseguro", strconv.Quote(limiteRunas(e.Nome, 40))))
			continue
		case !trust && naoConfiavel(e.Ref):
			cortes = append(cortes, ste.Msg("fabricar.peca_removivel", e.Nome))
			continue
		}
		kept = append(kept, e)
	}
	return kept, cortes
}

// limiteRunas corta o texto em n runas: um nome hostil não vira um aviso enorme.
func limiteRunas(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n]) + "..."
	}
	return s
}

// naoConfiavel aplica a regra das raízes ao volume do arquivo da peça.
func naoConfiavel(ref string) bool {
	p := caminhoDoRef(ref)
	vol := filepath.VolumeName(p)
	if vol == "" {
		return isUntrustedRoot(p)
	}
	return isUntrustedRoot(vol + string(filepath.Separator))
}

// caminhoDoRef traduz o ref da matriz para um caminho de arquivo. A matriz guarda "~/caminho" (relativo ao
// HOME) e, em MCP e hook, "arquivo#nome": o "~" vira o HOME e o "#nome" sai. Sem HOME o ref fica como está
// e cai na regra de não confiável, que é o lado seguro.
func caminhoDoRef(ref string) string {
	if i := strings.Index(ref, "#"); i >= 0 {
		ref = ref[:i]
	}
	if ref == "~" || strings.HasPrefix(ref, "~/") || strings.HasPrefix(ref, `~\`) {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			return filepath.Join(home, filepath.FromSlash(ref[1:]))
		}
	}
	return filepath.FromSlash(ref)
}
