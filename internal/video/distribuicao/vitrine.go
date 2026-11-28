package distribuicao

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Limites do pacote de vitrine.
const (
	NumTitulos         = 5
	NumCapasPorFormato = 3
	MaxTitulo          = 100
	TituloIdeal        = 60
)

var reCapitulo = regexp.MustCompile(`^\d{1,2}:\d{2}$`)

// Titulo e uma opcao de titulo com nota de 0 a 10.
type Titulo struct {
	Texto   string   `json:"texto"`
	Nota    int      `json:"nota"`
	Motivos []string `json:"motivos,omitempty"`
}

// Capitulo marca um trecho do video na descricao.
type Capitulo struct {
	Inicio string `json:"inicio"`
	Titulo string `json:"titulo"`
}

// Vitrine e o pacote de publicacao de um video: titulos, capas, descricao, tags e rotulo de IA.
type Vitrine struct {
	Video     string              `json:"video"`
	Titulos   []Titulo            `json:"titulos"`
	Capas     map[string][]string `json:"capas"`
	Resumo    string              `json:"resumo"`
	Capitulos []Capitulo          `json:"capitulos"`
	Fontes    []string            `json:"fontes"`
	Tags      []string            `json:"tags"`
	RotuloIA  bool                `json:"rotulo_ia"`
	SRT       string              `json:"srt"`
}

// NotaTitulo pontua um titulo: palavra principal no inicio, ate 60 caracteres, numero, sem gritaria.
func NotaTitulo(texto, palavra string) Titulo {
	t := Titulo{Texto: strings.TrimSpace(texto)}
	n := len([]rune(t.Texto))
	if n == 0 || n > MaxTitulo {
		t.Motivos = []string{fmt.Sprintf("titulo fora de 1 a %d caracteres", MaxTitulo)}
		return t
	}
	regras := []struct {
		pontos int
		ok     bool
		motivo string
	}{
		{3, n <= TituloIdeal, "ate 60 caracteres"},
		{3, comecaCom(t.Texto, palavra), "palavra principal no inicio"},
		{2, strings.IndexFunc(t.Texto, unicode.IsDigit) >= 0, "tem numero"},
		{1, !gritado(t.Texto), "sem caixa alta"},
		{1, !strings.Contains(t.Texto, "!!"), "sem exclamacao repetida"},
	}
	for _, r := range regras {
		if r.ok {
			t.Nota += r.pontos
			t.Motivos = append(t.Motivos, r.motivo)
		}
	}
	return t
}

func comecaCom(texto, palavra string) bool {
	palavra = strings.ToLower(strings.TrimSpace(palavra))
	if palavra == "" {
		return false
	}
	r := []rune(strings.ToLower(texto))
	if len(r) > 20 {
		r = r[:20]
	}
	return strings.Contains(string(r), palavra)
}

func gritado(texto string) bool {
	letras, altas := 0, 0
	for _, r := range texto {
		if unicode.IsLetter(r) {
			letras++
			if unicode.IsUpper(r) {
				altas++
			}
		}
	}
	return letras > 0 && float64(altas)/float64(letras) >= 0.6
}

// PontuarTitulos devolve os titulos em ordem decrescente de nota, mantendo a ordem original nos empates.
func PontuarTitulos(textos []string, palavra string) []Titulo {
	out := make([]Titulo, 0, len(textos))
	for _, t := range textos {
		out = append(out, NotaTitulo(t, palavra))
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Nota > out[j].Nota })
	return out
}

// Validar confere o pacote contra o disco; nada e aprovado com arquivo ausente.
func (v Vitrine) Validar(projDir string) error {
	if len(v.Titulos) != NumTitulos {
		return fmt.Errorf("%w: %d titulos, esperado %d", ErrVitrineInvalida, len(v.Titulos), NumTitulos)
	}
	for _, t := range v.Titulos {
		if n := len([]rune(t.Texto)); n == 0 || n > MaxTitulo {
			return fmt.Errorf("%w: titulo fora de 1 a %d caracteres", ErrVitrineInvalida, MaxTitulo)
		}
	}
	if len(v.Fontes) == 0 || len(v.Tags) == 0 || strings.TrimSpace(v.Resumo) == "" {
		return fmt.Errorf("%w: resumo, fontes e tags sao obrigatorios", ErrVitrineInvalida)
	}
	if !v.RotuloIA {
		return ErrSemRotuloIA
	}
	for _, c := range v.Capitulos {
		if !reCapitulo.MatchString(c.Inicio) || strings.TrimSpace(c.Titulo) == "" {
			return fmt.Errorf("%w: capitulo %q invalido", ErrVitrineInvalida, c.Inicio)
		}
	}
	return v.validarArquivos(projDir)
}

func (v Vitrine) validarArquivos(projDir string) error {
	if len(v.Capas) == 0 {
		return fmt.Errorf("%w: sem capas", ErrVitrineInvalida)
	}
	for formato, capas := range v.Capas {
		if len(capas) != NumCapasPorFormato {
			return fmt.Errorf("%w: formato %s com %d capas, esperado %d", ErrVitrineInvalida, formato, len(capas), NumCapasPorFormato)
		}
		for _, c := range capas {
			if err := existe(projDir, c); err != nil {
				return err
			}
		}
	}
	for _, rel := range []string{v.Video, v.SRT} {
		if err := existe(projDir, rel); err != nil {
			return err
		}
	}
	return nil
}

func existe(projDir, rel string) error {
	if rel == "" || filepath.IsAbs(rel) || strings.Contains(filepath.ToSlash(filepath.Clean(rel)), "..") {
		return fmt.Errorf("%w: caminho %q fora do projeto", ErrVitrineInvalida, rel)
	}
	if _, err := os.Stat(filepath.Join(projDir, filepath.FromSlash(rel))); err != nil {
		return fmt.Errorf("%w: arquivo %q ausente", ErrVitrineInvalida, rel)
	}
	return nil
}

// Descricao monta o texto de publicacao: resumo, capitulos, fontes e aviso de IA.
func (v Vitrine) Descricao() string {
	var b strings.Builder
	b.WriteString(strings.TrimSpace(v.Resumo) + "\n")
	if len(v.Capitulos) > 0 {
		b.WriteString("\nCapitulos:\n")
		for _, c := range v.Capitulos {
			b.WriteString(c.Inicio + " " + c.Titulo + "\n")
		}
	}
	b.WriteString("\nFontes:\n")
	for _, f := range v.Fontes {
		b.WriteString("- " + f + "\n")
	}
	if v.RotuloIA {
		b.WriteString("\nConteudo com apoio de IA (imagens e narracao).\n")
	}
	return b.String()
}

// GravarPublicacao valida o pacote e grava publicacao.md com o titulo escolhido (indice base 0).
func GravarPublicacao(projDir string, v Vitrine, escolhido int) error {
	if err := v.Validar(projDir); err != nil {
		return err
	}
	if err := ExigirLiberada(projDir, v.Video); err != nil {
		return err
	}
	if escolhido < 0 || escolhido >= len(v.Titulos) {
		return fmt.Errorf("%w: titulo %d fora de 0 a %d", ErrVitrineInvalida, escolhido, len(v.Titulos)-1)
	}
	var b strings.Builder
	b.WriteString("# Publicacao\n\n- Video: " + v.Video + "\n- Titulo: " + v.Titulos[escolhido].Texto + "\n- Tags: " + strings.Join(v.Tags, ", ") + "\n- Rotulo de IA: sim\n- Legenda: " + v.SRT + "\n\n## Descricao\n\n" + v.Descricao())
	return gravarAtomico(filepath.Join(projDir, "publicacao.md"), []byte(b.String()))
}

// ValidarTextoCapa limita o texto da capa: ate 4 palavras (3 no curto).
func ValidarTextoCapa(linha1, linha2 string, curto bool) error {
	max := 4
	if curto {
		max = 3
	}
	n := len(strings.Fields(linha1)) + len(strings.Fields(linha2))
	if strings.TrimSpace(linha1) == "" || n > max {
		return fmt.Errorf("%w: texto da capa com %d palavras (1 a %d)", ErrVitrineInvalida, n, max)
	}
	return nil
}
