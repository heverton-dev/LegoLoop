package inventario

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Classes de classificação do inventário de docs/ (contrato AC-E01.1).
const (
	ClasseImplementado = "IMPLEMENTADO"
	ClasseUltrapassado = "ULTRAPASSADO"
	ClasseSuperado     = "SUPERADO"
	ClassePendente     = "PENDENTE"
	ClasseRever        = "REVER"
)

// todasClasses lista as classes na ordem do relatório.
var todasClasses = []string{
	ClasseImplementado, ClasseUltrapassado, ClasseSuperado, ClassePendente, ClasseRever,
}

// Entrada é a classificação de um arquivo de docs/.
type Entrada struct {
	Caminho string // relativo a docs/, com barra simples
	Classe  string
	Dono    string // processo ou épico dono do documento
	Motivo  string
}

var (
	padraoISO       = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})_(.+)$`)
	padraoData      = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})_`)
	padraoRelatorio = regexp.MustCompile(`docs/reports/([A-Za-z0-9._-]+\.md)`)
)

// tiposDocumento lista os prefixos de tipo após a data no nome ISO.
var tiposDocumento = []string{
	"MANUAL_FUNCAO_", "INTENT_", "SPEC_", "PLAN_", "ROUTE_",
	"DAG_", "REPORT_", "MANUAL_", "INVENTARIO_",
}

// arquivosGovernanca lista os documentos de governança viva que citam relatórios globais.
var arquivosGovernanca = []string{"AGENTS.md", "README.md", "MEMORY.md"}

// Classificar decide classe, dono e motivo de um caminho relativo a docs/.
func Classificar(rel string) Entrada {
	rel = path.Clean(filepath.ToSlash(rel))
	dir, nome := path.Split(rel)
	categoria := strings.TrimSuffix(dir, "/")
	if i := strings.Index(categoria, "/"); i >= 0 {
		categoria = categoria[:i]
	}
	switch categoria {
	case "templates", "vocab":
		return Entrada{Caminho: rel, Classe: ClasseImplementado, Dono: "lei-6-15", Motivo: "lei Lei 6/15"}
	case "intents", "specs", "plans", "routes", "dags":
		return classificarQuinteto(rel, nome)
	case "manuals":
		return classificarISO(rel, nome, ClasseImplementado, "manual de função")
	case "fichas":
		return Entrada{Caminho: rel, Classe: ClasseImplementado, Dono: strings.TrimSuffix(nome, path.Ext(nome)), Motivo: "ficha de governança"}
	case "books":
		return Entrada{Caminho: rel, Classe: ClasseImplementado, Dono: donoLivro(dir), Motivo: "livro didático"}
	case "reports":
		return classificarISO(rel, nome, ClasseUltrapassado, "documentação de construção")
	case "videorea":
		return Entrada{Caminho: rel, Classe: ClasseUltrapassado, Dono: "video-rea", Motivo: "documentação de construção"}
	case "sessions":
		return Entrada{Caminho: rel, Classe: ClasseUltrapassado, Dono: "sessao", Motivo: "documentação de construção"}
	case "archive":
		return Entrada{Caminho: rel, Classe: ClasseUltrapassado, Dono: "arquivo-morto", Motivo: "já em arquivo morto"}
	case "qa-web":
		return Entrada{Caminho: rel, Classe: ClasseImplementado, Dono: "web-qa", Motivo: "relatório de homologação"}
	}
	return Entrada{Caminho: rel, Classe: ClasseRever, Motivo: "diretório fora do previsto"}
}

// classificarQuinteto classifica documento do Quinteto de Governança ISO.
func classificarQuinteto(rel, nome string) Entrada {
	semExt := strings.TrimSuffix(nome, path.Ext(nome))
	if !padraoISO.MatchString(semExt) {
		return Entrada{Caminho: rel, Classe: ClassePendente, Motivo: "demanda sem padrão ISO"}
	}
	return Entrada{Caminho: rel, Classe: ClasseImplementado, Dono: slugDe(semExt), Motivo: "quinteto de governança"}
}

// classificarISO classifica documento com nome no padrão ISO da governança.
func classificarISO(rel, nome, classe, motivo string) Entrada {
	semExt := strings.TrimSuffix(nome, path.Ext(nome))
	if !padraoISO.MatchString(semExt) {
		return Entrada{Caminho: rel, Classe: ClasseRever, Motivo: "nome fora do padrão ISO"}
	}
	return Entrada{Caminho: rel, Classe: classe, Dono: slugDe(semExt), Motivo: motivo}
}

// donoLivro devolve o slug do livro em books/<slug>/....
func donoLivro(dir string) string {
	partes := strings.Split(strings.Trim(dir, "/"), "/")
	if len(partes) >= 2 && partes[1] != "" {
		return partes[1]
	}
	return "books"
}

// slugDe extrai o slug da demanda de um nome de arquivo no padrão ISO.
func slugDe(semExt string) string {
	m := padraoISO.FindStringSubmatch(semExt)
	if m == nil {
		return semExt
	}
	resto := m[2]
	for _, tipo := range tiposDocumento {
		if strings.HasPrefix(resto, tipo) {
			return strings.TrimPrefix(resto, tipo)
		}
	}
	return resto
}

// classeValida informa se a classe pertence ao contrato AC-E01.1.
func classeValida(c string) bool {
	switch c {
	case ClasseImplementado, ClasseUltrapassado, ClasseSuperado, ClassePendente, ClasseRever:
		return true
	}
	return false
}

// Varrer classifica todos os arquivos de docs/ sem exceção.
func Varrer(raizDocs string) ([]Entrada, error) {
	var entradas []Entrada
	err := filepath.WalkDir(raizDocs, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, errRel := filepath.Rel(raizDocs, p)
		if errRel != nil {
			return errRel
		}
		entradas = append(entradas, Classificar(filepath.ToSlash(rel)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	globais, err := RelatoriosGlobais(filepath.Dir(raizDocs))
	if err != nil {
		return nil, err
	}
	sort.Slice(entradas, func(i, j int) bool { return entradas[i].Caminho < entradas[j].Caminho })
	return refinar(entradas, globais), nil
}

// RelatoriosGlobais extrai nomes de relatórios citados na governança viva.
func RelatoriosGlobais(raizRepo string) (map[string]bool, error) {
	globais := make(map[string]bool)
	for _, nome := range arquivosGovernanca {
		conteudo, err := os.ReadFile(filepath.Join(raizRepo, nome))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, m := range padraoRelatorio.FindAllStringSubmatch(string(conteudo), -1) {
			globais[m[1]] = true
		}
	}
	return globais, nil
}

// refinar aplica as decisões de classe que exigem visão do conjunto.
func refinar(entradas []Entrada, globais map[string]bool) []Entrada {
	fechados := fechamentos(entradas)
	maisNovo := datasMaisRecentes(entradas)
	for i := range entradas {
		e := &entradas[i]
		base := path.Base(e.Caminho)
		switch {
		case strings.HasPrefix(e.Caminho, "reports/") && globais[base]:
			e.Classe, e.Dono = ClasseImplementado, "governanca-viva"
			e.Motivo = "relatório global da governança viva"
		case strings.HasPrefix(e.Caminho, "reports/") && e.Classe == ClasseUltrapassado && dataDe(base) < maisNovo[e.Dono]:
			e.Classe, e.Motivo = ClasseSuperado, "substituído por relatório mais recente"
		case ehQuinteto(e.Caminho) && e.Classe == ClasseImplementado && e.Dono != "" && !fechados[e.Dono]:
			e.Classe, e.Motivo = ClassePendente, "demanda sem relatório de fechamento"
		}
	}
	return entradas
}

// ehQuinteto informa se o caminho pertence ao Quinteto de Governança ISO.
func ehQuinteto(rel string) bool {
	switch strings.SplitN(rel, "/", 2)[0] {
	case "intents", "specs", "plans", "routes", "dags":
		return true
	}
	return false
}

// fechamentos devolve os slugs de demandas com relatório de fechamento.
func fechamentos(entradas []Entrada) map[string]bool {
	fechados := make(map[string]bool)
	for _, e := range entradas {
		if strings.HasPrefix(e.Caminho, "reports/") && e.Dono != "" {
			fechados[e.Dono] = true
		}
	}
	return fechados
}

// datasMaisRecentes devolve a data ISO mais recente por slug de relatório.
func datasMaisRecentes(entradas []Entrada) map[string]string {
	maisNovo := make(map[string]string)
	for _, e := range entradas {
		if !strings.HasPrefix(e.Caminho, "reports/") || e.Dono == "" {
			continue
		}
		if data := dataDe(path.Base(e.Caminho)); data > maisNovo[e.Dono] {
			maisNovo[e.Dono] = data
		}
	}
	return maisNovo
}

// dataDe extrai a data ISO do nome de um arquivo.
func dataDe(nome string) string {
	m := padraoData.FindStringSubmatch(nome)
	if m == nil {
		return ""
	}
	return m[1]
}

// GerarRelatorio escreve o inventário completo em markdown.
func GerarRelatorio(entradas []Entrada, saida string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# Inventário do docs/ — 2026-11-16\n\n")
	fmt.Fprintf(&b, "Total de arquivos classificados: %d. Zero arquivo sem classificação.\n\n", len(entradas))
	fmt.Fprintf(&b, "Classes do contrato AC-E01.1: %s.\n\n", strings.Join(todasClasses, ", "))
	escreverResumo(&b, entradas)
	for _, classe := range todasClasses {
		escreverTabela(&b, classe, entradas)
	}
	return os.WriteFile(saida, []byte(b.String()), 0o644)
}

// escreverResumo escreve as contagens por classe e por dono.
func escreverResumo(b *strings.Builder, entradas []Entrada) {
	porClasse := make(map[string]int)
	porDono := make(map[string]int)
	for _, e := range entradas {
		porClasse[e.Classe]++
		porDono[e.Dono]++
	}
	b.WriteString("## Resumo por classe\n\n| Classe | Quantidade |\n| :--- | ---: |\n")
	for _, classe := range todasClasses {
		fmt.Fprintf(b, "| %s | %d |\n", classe, porClasse[classe])
	}
	b.WriteString("\n## Resumo por dono\n\n| Dono | Quantidade |\n| :--- | ---: |\n")
	for _, dono := range chavesOrdenadas(porDono) {
		fmt.Fprintf(b, "| %s | %d |\n", dono, porDono[dono])
	}
	b.WriteString("\n")
}

// escreverTabela escreve a tabela de uma classe com caminho, dono e motivo.
func escreverTabela(b *strings.Builder, classe string, entradas []Entrada) {
	fmt.Fprintf(b, "## %s\n\n| Caminho | Dono | Motivo |\n| :--- | :--- | :--- |\n", classe)
	for _, e := range entradas {
		if e.Classe == classe {
			fmt.Fprintf(b, "| `%s` | %s | %s |\n", e.Caminho, e.Dono, e.Motivo)
		}
	}
	b.WriteString("\n")
}

// chavesOrdenadas devolve as chaves de um mapa em ordem alfabética.
func chavesOrdenadas(m map[string]int) []string {
	chaves := make([]string, 0, len(m))
	for k := range m {
		chaves = append(chaves, k)
	}
	sort.Strings(chaves)
	return chaves
}
