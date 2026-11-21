package mcp

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SkillItem representa uma skill indexada pelo motor.
type SkillItem struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Tier      string `json:"tier"`
	Natureza  string `json:"natureza"`
	Categoria string `json:"categoria"`
	Resumo    string `json:"resumo"`
	CorpoPath string `json:"corpo_path"`
}

// SkillMatch retorno enriquecido com score de busca.
type SkillMatch struct {
	ID        string  `json:"skill_id"`
	Nome      string  `json:"name"`
	Score     float64 `json:"score"`
	Natureza  string  `json:"natureza"`
	Tier      string  `json:"tier"`
	Categoria string  `json:"categoria"`
	Resumo    string  `json:"resumo"`
}

// SkillEngine gerencia o índice e busca de skills.
type SkillEngine struct {
	skills map[string]SkillItem
}

// NewSkillEngine instancia um novo motor de skills.
func NewSkillEngine() *SkillEngine {
	return &SkillEngine{skills: make(map[string]SkillItem)}
}

// IndexDirectory varre recursivamente e indexa skills, ignorando pastas com ponto.
func (e *SkillEngine) IndexDirectory(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.EqualFold(d.Name(), "SKILL.md") {
			item, parseErr := parseSkillFile(path)
			if parseErr == nil && item.ID != "" {
				e.skills[item.ID] = item
			}
		}
		return nil
	})
}

// List retorna todas as skills indexadas.
func (e *SkillEngine) List() []SkillItem {
	res := make([]SkillItem, 0, len(e.skills))
	for _, item := range e.skills {
		res = append(res, item)
	}
	return res
}

// Search busca skills relevantes com base na query.
func (e *SkillEngine) Search(query string, topK int) []SkillMatch {
	terms := tokenize(query)
	var matches []SkillMatch
	for _, s := range e.skills {
		score := calculateScore(terms, s)
		if score > 0 {
			matches = append(matches, SkillMatch{
				ID:        s.ID,
				Nome:      s.Nome,
				Score:     score,
				Natureza:  s.Natureza,
				Tier:      s.Tier,
				Categoria: s.Categoria,
				Resumo:    s.Resumo,
			})
		}
	}
	return topMatches(matches, topK)
}

func parseSkillFile(path string) (SkillItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return SkillItem{}, err
	}
	defer f.Close()

	item := SkillItem{
		ID:        filepath.Base(filepath.Dir(path)),
		Nome:      filepath.Base(filepath.Dir(path)),
		Tier:      "Prata",
		Natureza:  "Processo Canônico",
		Categoria: "general",
		CorpoPath: path,
	}
	leitor := &leitorFrontmatter{item: &item}
	scanner := bufio.NewScanner(f)
	dentro := false
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "---" {
			if dentro {
				break
			}
			dentro = true
			continue
		}
		if dentro {
			leitor.linha(scanner.Text())
		}
	}
	leitor.fechar()
	return item, scanner.Err()
}

// leitorFrontmatter lê o frontmatter linha a linha e junta valores em bloco YAML (|, >).
type leitorFrontmatter struct {
	item  *SkillItem
	chave string
	bloco []string
}

func (l *leitorFrontmatter) fechar() {
	if l.chave != "" {
		aplicarMeta(l.chave, strings.Join(l.bloco, " "), l.item)
	}
	l.chave, l.bloco = "", nil
}

func (l *leitorFrontmatter) linha(bruta string) {
	limpa := strings.TrimSpace(bruta)
	indentada := bruta == "" || bruta[0] == ' ' || bruta[0] == '	'
	if l.chave != "" && indentada {
		if limpa != "" {
			l.bloco = append(l.bloco, limpa)
		}
		return
	}
	l.fechar()
	l.chave = extractMetaLine(limpa, l.item)
}

// extractMetaLine lê uma linha "chave: valor" do frontmatter. Se o valor abre um bloco YAML
// (|, >, |-, >-), devolve a chave para o chamador juntar as linhas indentadas seguintes.
func extractMetaLine(line string, item *SkillItem) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return ""
	}
	key := strings.ToLower(strings.TrimSpace(parts[0]))
	val := strings.TrimSpace(parts[1])
	switch val {
	case "|", ">", "|-", ">-", "|+", ">+":
		return key
	}
	aplicarMeta(key, val, item)
	return ""
}

func aplicarMeta(key, val string, item *SkillItem) {
	switch key {
	case "name":
		item.Nome = val
		item.ID = val
	case "description":
		item.Resumo = val
	case "tier":
		item.Tier = val
	case "natureza":
		item.Natureza = val
	case "categoria":
		item.Categoria = val
	}
}

// palavrasVazias não carregam sentido na busca e casavam dentro de outras palavras ("com" em composer).
var palavrasVazias = map[string]bool{
	"com": true, "de": true, "da": true, "do": true, "das": true, "dos": true, "para": true, "por": true,
	"um": true, "uma": true, "uns": true, "umas": true, "em": true, "no": true, "na": true, "nos": true, "nas": true,
	"ao": true, "aos": true, "que": true, "ou": true, "se": true, "como": true, "the": true, "and": true,
	"for": true, "with": true, "of": true, "to": true, "in": true, "on": true, "an": true,
}

func tokenize(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	var res []string
	for _, w := range words {
		w = strings.Trim(w, ",.-_!?;:()[]")
		if len(w) > 1 && !palavrasVazias[w] {
			res = append(res, w)
		}
	}
	return res
}

const (
	pesoNomeExato  = 3.0 // termo igual a um pedaço do nome (aidd-bigbrain -> aidd, bigbrain)
	pesoNomeParte  = 2.5 // termo contido no nome
	pesoNomeRaiz   = 1.5 // mesma raiz (6 letras) de um pedaço do nome: homologar ~ homologue
	pesoResumo     = 1.0
	pesoResumoRaiz = 0.5
	tamanhoRaiz    = 6
	bonusNomeExato = 0.5
)

// calculateScore pontua pelo nome e pela descrição. Categoria e natureza ficam fora:
// os padrões "general" e "Processo Canônico" casavam com quase toda busca.
func calculateScore(terms []string, item SkillItem) float64 {
	if len(terms) == 0 {
		return 0
	}
	nome := strings.ToLower(item.ID + " " + item.Nome)
	pedacosNome := partesDoTexto(nome)
	resumo := strings.ToLower(item.Resumo)
	pedacosResumo := partesDoTexto(resumo)
	var soma float64
	for _, t := range terms {
		soma += pesoDoTermo(t, nome, pedacosNome, resumo, pedacosResumo)
	}
	score := soma / (pesoNomeExato * float64(len(terms)))
	if score > 0 && strings.Join(terms, "-") == strings.ToLower(item.ID) {
		score += bonusNomeExato
	}
	return score
}

// pesoDoTermo devolve o maior peso com que o termo casa no nome ou na descrição.
func pesoDoTermo(t, nome string, pedacosNome []string, resumo string, pedacosResumo []string) float64 {
	switch {
	case contem(pedacosNome, t):
		return pesoNomeExato
	case len(t) >= 4 && strings.Contains(nome, t):
		return pesoNomeParte
	case temRaiz(pedacosNome, t):
		return pesoNomeRaiz
	case contem(pedacosResumo, t) || (len(t) >= 4 && strings.Contains(resumo, t)):
		return pesoResumo
	case temRaiz(pedacosResumo, t):
		return pesoResumoRaiz
	}
	return 0
}

func partesDoTexto(texto string) []string {
	return strings.FieldsFunc(texto, func(r rune) bool {
		return !(r == '+' || r == '#' || (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || r >= 0x80)
	})
}

func contem(partes []string, t string) bool {
	for _, p := range partes {
		if p == t {
			return true
		}
	}
	return false
}

// temRaiz diz se o termo e algum pedaço do texto começam com as mesmas tamanhoRaiz letras.
func temRaiz(partes []string, t string) bool {
	if len(t) < tamanhoRaiz {
		return false
	}
	raiz := t[:tamanhoRaiz]
	for _, p := range partes {
		if len(p) >= tamanhoRaiz && p[:tamanhoRaiz] == raiz {
			return true
		}
	}
	return false
}

// topMatches ordena por score decrescente (empate por id) e devolve os topK primeiros.
func topMatches(matches []SkillMatch, topK int) []SkillMatch {
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		return matches[i].ID < matches[j].ID
	})
	if topK > 0 && len(matches) > topK {
		return matches[:topK]
	}
	return matches
}
