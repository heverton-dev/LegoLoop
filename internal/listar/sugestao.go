package listar

import (
	"cmp"
	"path/filepath"
	"slices"
	"strings"

	"legoloop/internal/matriz"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

// MaxComandosTexto limita os comandos por sugestão na saída em texto; o --json traz todos.
const MaxComandosTexto = 5

func init() {
	for id, msg := range map[string]string{
		"listar.lacuna":          "A etapa %s do processo %s não tem peça.",
		"listar.lacuna_sugestao": "Estas peças fazem a etapa %s, mas não declaram o caso %s. Para dar o caso a uma peça, rode:",
		"listar.lacuna_sem_peca": "Nenhuma peça da matriz faz a etapa %s e pode receber o caso. Veja as peças com legoloop matriz ver --etapa %s.",
		"listar.lacuna_mais":     "Mais %d peças fazem a etapa %s. Veja todas com legoloop listar --json.",
		"listar.lacuna_mais_uma": "Mais 1 peça faz a etapa %s. Veja todas com legoloop listar --json.",
	} {
		ste.Default.Register(id, msg)
	}
}

// Sugestao diz, para uma etapa em lacuna e um caso pedido, quais peças da matriz fazem a etapa mas não declaram
// o caso, e o comando que dá o caso a cada uma. O comando repete o perfil atual da peça, porque `matriz perfil
// definir` substitui etapas, naturezas e casos: sem repetir, a peça perderia a etapa.
type Sugestao struct {
	Etapa    string   `json:"etapa"`
	Caso     string   `json:"caso"`
	Pecas    []string `json:"pecas"`
	Comandos []string `json:"comandos"`
}

// fonte é a matriz e o vocabulário do cofre. Sem matriz, não há candidatas (a lacuna aparece sem sugestão).
type fonte struct {
	entradas []matriz.Entrada
	voc      matriz.Vocab
}

func lerFonte(vault string) fonte {
	f := fonte{voc: matriz.DefaultVocab()}
	if v, err := matriz.LoadVocab(vault); err == nil {
		f.voc = v
	}
	if reg, err := matriz.LerRegistro(filepath.Join(vault, "registry_matrix.json")); err == nil {
		f.entradas = reg.Entradas
	}
	return f
}

// lacunas devolve as etapas sem peça e sem composição, sem repetir, na ordem do processo.
func lacunas(etapas []processo.Etapa) []processo.Etapa {
	var out []processo.Etapa
	for _, e := range etapas {
		if e.Peca == "" && e.Compoe == "" && !slices.ContainsFunc(out, func(o processo.Etapa) bool { return o.Tipo == e.Tipo }) {
			out = append(out, e)
		}
	}
	return out
}

// sugerir preenche Lacunas e Sugestoes de cada processo. A matriz só é lida se algum processo tem lacuna.
func sugerir(vault string, its []Item, etapasPor map[string][]processo.Etapa) {
	var f *fonte
	for i := range its {
		for _, e := range lacunas(etapasPor[its[i].Nome]) {
			if f == nil {
				lida := lerFonte(vault)
				f = &lida
			}
			its[i].Lacunas = append(its[i].Lacunas, e.Tipo)
			for _, caso := range f.casosDe(e) {
				if s, ok := f.sugestao(e.Tipo, caso); ok {
					its[i].Sugestoes = append(its[i].Sugestoes, s)
				}
			}
		}
	}
}

// casosDe usa os casos gravados na etapa; sem eles (arquivo antigo), os casos que o vocabulário liga à etapa.
func (f fonte) casosDe(e processo.Etapa) []string {
	if len(e.Casos) > 0 {
		return e.Casos
	}
	var todos []string
	for _, t := range f.voc.CasosDeUso {
		todos = append(todos, t.Nome)
	}
	return f.voc.CasosDaEtapa(todos, e.Tipo)
}

// candidata é uma peça (cópias idênticas contam uma vez) que faz a etapa e não declara o caso.
type candidata struct {
	alvo   string
	score  float64
	perfil *matriz.PerfilDeUso
}

func (f fonte) sugestao(etapa, caso string) (Sugestao, bool) {
	cs := f.candidatas(etapa, caso)
	if len(cs) == 0 {
		return Sugestao{}, false
	}
	s := Sugestao{Etapa: etapa, Caso: caso}
	for _, c := range cs {
		s.Pecas = append(s.Pecas, c.alvo)
		s.Comandos = append(s.Comandos, comando(c, caso))
	}
	return s, true
}

func (f fonte) candidatas(etapa, caso string) []candidata {
	porHash, hashesDoNome := map[string]matriz.Entrada{}, map[string]map[string]bool{}
	for _, e := range f.entradas {
		p := e.Perfil
		if p == nil || !slices.Contains(p.Etapas, etapa) || slices.Contains(p.CasosDeUso, caso) {
			continue
		}
		chave := e.Tipo + "|" + e.Hash
		if _, visto := porHash[chave]; !visto {
			porHash[chave] = e
		}
		if hashesDoNome[e.Nome] == nil {
			hashesDoNome[e.Nome] = map[string]bool{}
		}
		hashesDoNome[e.Nome][chave] = true
	}
	out := make([]candidata, 0, len(porHash))
	for _, e := range porHash {
		alvo := e.Nome
		if len(hashesDoNome[e.Nome]) > 1 {
			alvo = e.Ref // nomes iguais com conteúdos diferentes: só o caminho identifica a peça
		}
		out = append(out, candidata{alvo: alvo, score: e.Rating.TotalScore, perfil: e.Perfil})
	}
	slices.SortFunc(out, func(a, b candidata) int {
		return cmp.Or(cmp.Compare(b.score, a.score), strings.Compare(a.alvo, b.alvo))
	})
	return out
}

// comando monta o `perfil definir` com o caso novo primeiro e o perfil atual repetido.
func comando(c candidata, caso string) string {
	partes := []string{"legoloop matriz perfil definir", aspas(c.alvo), "--caso", caso}
	for _, g := range []struct {
		flag   string
		termos []string
	}{{"--caso", c.perfil.CasosDeUso}, {"--etapa", c.perfil.Etapas}, {"--natureza", c.perfil.Naturezas}} {
		for _, t := range g.termos {
			partes = append(partes, g.flag, aspas(t))
		}
	}
	return strings.Join(partes, " ")
}

// aspas protege o argumento que o shell quebraria (espaço ou caractere especial).
func aspas(s string) string {
	if strings.ContainsAny(s, " \t\"'&|<>();$`") {
		return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return s
}
