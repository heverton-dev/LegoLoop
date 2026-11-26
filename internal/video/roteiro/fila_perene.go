package roteiro

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// FilaPerene guarda pautas perenes coletadas para futuras producoes.
// Pautas quentes (dependentes de timing) nunca entram nesta fila.
type FilaPerene struct {
	SchemaVersion int              `json:"schema_version"`
	Pautas        []PautaCandidata `json:"pautas"`
}

// CaminhoFilaPerene devolve o caminho padrao de acervo/fila_perene.json a partir de acervo/temas.json.
func CaminhoFilaPerene(acervoTemasPath string) string {
	return filepath.Join(filepath.Dir(acervoTemasPath), "fila_perene.json")
}

// CarregarFilaPerene le acervo/fila_perene.json, retornando fila vazia se nao existir.
func CarregarFilaPerene(caminho string) (FilaPerene, error) {
	dados, err := os.ReadFile(caminho)
	if os.IsNotExist(err) {
		return FilaPerene{SchemaVersion: 1}, nil
	}
	if err != nil {
		return FilaPerene{}, err
	}
	var f FilaPerene
	if err := json.Unmarshal(dados, &f); err != nil {
		return FilaPerene{}, err
	}
	if f.SchemaVersion == 0 {
		f.SchemaVersion = 1
	}
	return f, nil
}

// SalvarFilaPerene grava acervo/fila_perene.json de forma atomica (temp + rename).
func SalvarFilaPerene(caminho string, f FilaPerene) error {
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		return err
	}
	if f.SchemaVersion == 0 {
		f.SchemaVersion = 1
	}
	dados, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	tmp := caminho + ".tmp"
	if err := os.WriteFile(tmp, append(dados, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, caminho)
}

// MesclarFilaPerene adiciona apenas pautas com tipo_pauta == "perene" nao produzidas e sem duplicatas.
func MesclarFilaPerene(atual FilaPerene, coletadas []PautaCandidata, acervo AcervoTemas) FilaPerene {
	vistos := map[string]bool{}
	var saida []PautaCandidata
	for _, it := range atual.Pautas {
		if ok, norm := validarItemFila(it, acervo, vistos); ok {
			vistos[norm.Slug] = true
			saida = append(saida, norm)
		}
	}
	for _, it := range coletadas {
		if ok, norm := validarItemFila(it, acervo, vistos); ok {
			vistos[norm.Slug] = true
			saida = append(saida, norm)
		}
	}
	return FilaPerene{SchemaVersion: 1, Pautas: saida}
}

func validarItemFila(it PautaCandidata, acervo AcervoTemas, vistos map[string]bool) (bool, PautaCandidata) {
	if !strings.EqualFold(strings.TrimSpace(it.TipoPauta), "perene") {
		return false, PautaCandidata{}
	}
	slug := NormalizarSlug(it.Slug)
	if slug == "" {
		slug = NormalizarSlug(it.Titulo)
	}
	if slug == "" || vistos[slug] || acervo.Contem(it.Titulo) || acervo.Contem(slug) {
		return false, PautaCandidata{}
	}
	it.Slug = slug
	it.TipoPauta = "perene"
	return true, it
}

// RemoverDaFilaPerene remove da fila o tema aprovado pelo titulo ou slug.
func RemoverDaFilaPerene(atual FilaPerene, tituloOuSlug string) FilaPerene {
	alvo := NormalizarSlug(tituloOuSlug)
	var resto []PautaCandidata
	for _, it := range atual.Pautas {
		if alvo != "" && (NormalizarSlug(it.Slug) == alvo || NormalizarSlug(it.Titulo) == alvo) {
			continue
		}
		resto = append(resto, it)
	}
	return FilaPerene{SchemaVersion: 1, Pautas: resto}
}

// CombinarPautasComFila une pautas da rodada atual e estoque perene sem duplicatas nem produzidos.
func CombinarPautasComFila(coletadas []PautaCandidata, fila FilaPerene, acervo AcervoTemas) []PautaCandidata {
	vistos := map[string]bool{}
	var saida []PautaCandidata
	saida = anexarCandidatasUnicas(saida, coletadas, acervo, vistos)
	saida = anexarCandidatasUnicas(saida, fila.Pautas, acervo, vistos)
	return saida
}

func anexarCandidatasUnicas(acc, itens []PautaCandidata, acervo AcervoTemas, vistos map[string]bool) []PautaCandidata {
	for _, it := range itens {
		slug := NormalizarSlug(it.Slug)
		if slug == "" {
			slug = NormalizarSlug(it.Titulo)
		}
		if slug == "" || vistos[slug] || acervo.Contem(it.Titulo) || acervo.Contem(slug) {
			continue
		}
		it.Slug = slug
		vistos[slug] = true
		acc = append(acc, it)
	}
	return acc
}
