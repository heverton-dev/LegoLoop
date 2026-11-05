package fabricar

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"legoloop/internal/matriz"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

func init() {
	ste.Default.Register("fabricar.sem_matriz", "A matriz de peças não existe em %s. Rode legoloop matriz gerar e tente de novo.")
	ste.Default.Register("fabricar.drives_ignorados", "O fabricar não varre drives. As peças vêm da matriz. Os drives pedidos foram ignorados.")
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
// que a Fusão usa (MCPs, hooks e loops nunca viram passo).
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
