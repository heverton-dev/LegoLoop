package distribuicao

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/heverton-dev/LegoLoop/internal/redc"
)

// IdentidadeCanal guarda simbolo, paleta e fontes do nucleo do canal.
type IdentidadeCanal struct {
	Simbolo string            `json:"simbolo"`
	Paleta  map[string]string `json:"paleta"`
	Fontes  map[string]string `json:"fontes"`
}

// RedeCanal representa os campos adaptados de uma rede social no Kit de Canal.
type RedeCanal struct {
	NomeExibicao         string   `json:"nome_exibicao,omitempty"`
	Handle               string   `json:"handle,omitempty"`
	Descricao            string   `json:"descricao,omitempty"`
	Bio                  string   `json:"bio,omitempty"`
	Tagline              string   `json:"tagline,omitempty"`
	Sobre                string   `json:"sobre,omitempty"`
	PalavrasChave        []string `json:"palavras_chave,omitempty"`
	Links                []string `json:"links,omitempty"`
	Playlists            []string `json:"playlists,omitempty"`
	Destaques            []string `json:"destaques,omitempty"`
	VideoEmDestaque      string   `json:"video_em_destaque,omitempty"`
	PrimeiroPostFixado   string   `json:"primeiro_post_fixado,omitempty"`
	OrdemDePreenchimento []string `json:"ordem_de_preenchimento,omitempty"`
}

// PautaCanal descreve uma pauta inicial planejada para o canal.
type PautaCanal struct {
	Titulo string `json:"titulo"`
	Status string `json:"status"`
}

// TextosCanal espelha textos.json do Kit de Canal original.
type TextosCanal struct {
	Nome                string               `json:"nome"`
	Handle              string               `json:"handle"`
	Promessa            string               `json:"promessa"`
	Tagline             string               `json:"tagline"`
	Tom                 string               `json:"tom"`
	Pilares             []string             `json:"pilares"`
	IdentidadeVisual    IdentidadeCanal      `json:"identidade_visual"`
	DisponibilidadeNome map[string]string    `json:"disponibilidade_nome"`
	Redes               map[string]RedeCanal `json:"redes"`
	PrimeirasPautas     []PautaCanal         `json:"primeiras_pautas"`
	HandleFinal         string               `json:"handle_final,omitempty"`
}

// RedeSpec define os limites de caracteres e se foram verificados em pagina oficial.
type RedeSpec struct {
	Nome       string
	Limites    map[string]int
	Verificado bool
}

// RelatorioTextosCanal informa o arquivo colar.md gerado e quantos campos excederam limites.
type RelatorioTextosCanal struct {
	Saida     string   `json:"saida"`
	Excedidos int      `json:"excedidos"`
	Campos    []string `json:"campos_excedidos,omitempty"`
}

// CatalogoRedes devolve os limites por rede equivalentes a ferramentas/redes.json.
func CatalogoRedes() map[string]RedeSpec {
	return map[string]RedeSpec{
		"youtube":   {Nome: "YouTube", Limites: map[string]int{"nome_exibicao": 50, "handle": 30, "descricao": 1000}},
		"instagram": {Nome: "Instagram", Limites: map[string]int{"nome_exibicao": 30, "handle": 30, "bio": 150}},
		"tiktok":    {Nome: "TikTok", Limites: map[string]int{"nome_exibicao": 30, "handle": 24, "bio": 80}},
		"linkedin":  {Nome: "LinkedIn (pagina da empresa)", Limites: map[string]int{"nome_exibicao": 100, "tagline": 120, "sobre": 2000}},
	}
}

// GravarTextosCanal valida o nucleo e salva textos.json de forma atomica.
func GravarTextosCanal(canalDir string, tc TextosCanal) error {
	if err := tc.Validar(); err != nil {
		return err
	}
	dados, err := json.MarshalIndent(tc, "", "  ")
	if err != nil {
		return err
	}
	return gravarAtomico(filepath.Join(canalDir, "textos.json"), append(dados, '\n'))
}

// Validar confere os campos essenciais de textos.json.
func (tc TextosCanal) Validar() error {
	if strings.TrimSpace(tc.Nome) == "" || strings.TrimSpace(tc.Handle) == "" || strings.TrimSpace(tc.Promessa) == "" || strings.TrimSpace(tc.Tagline) == "" {
		return fmt.Errorf("%w: nome, handle, promessa e tagline sao obrigatorios", ErrPedidoInvalido)
	}
	if len(tc.Pilares) == 0 || len(tc.Redes) == 0 || strings.TrimSpace(tc.IdentidadeVisual.Simbolo) == "" {
		return fmt.Errorf("%w: pilares, simbolo e redes sao obrigatorios", ErrPedidoInvalido)
	}
	return nil
}

// GerarTextosCanal le textos.json, confere limites de redes e grava colar.md.
func GerarTextosCanal(canalDir string) (RelatorioTextosCanal, error) {
	dados, err := os.ReadFile(filepath.Join(canalDir, "textos.json"))
	if err != nil {
		return RelatorioTextosCanal{}, err
	}
	var tc TextosCanal
	if err := json.Unmarshal(dados, &tc); err != nil {
		return RelatorioTextosCanal{}, fmt.Errorf("%w: textos.json: %v", ErrPedidoInvalido, err)
	}
	if err := tc.Validar(); err != nil {
		return RelatorioTextosCanal{}, err
	}
	md, rel := renderizarColarMD(tc, CatalogoRedes())
	saida := filepath.Join(canalDir, "colar.md")
	if err := gravarAtomico(saida, []byte(md)); err != nil {
		return RelatorioTextosCanal{}, err
	}
	rel.Saida = saida
	if rel.Excedidos > 0 {
		return rel, fmt.Errorf("%w: %d limite(s) excedido(s): %s", ErrPedidoInvalido, rel.Excedidos, strings.Join(rel.Campos, ", "))
	}
	return rel, nil
}

func renderizarColarMD(tc TextosCanal, cat map[string]RedeSpec) (string, RelatorioTextosCanal) {
	var b strings.Builder
	var rel RelatorioTextosCanal
	escreverCabecalhoColar(&b, tc)
	for _, chave := range []string{"youtube", "instagram", "tiktok", "linkedin"} {
		rc, ok := tc.Redes[chave]
		if !ok {
			continue
		}
		escreverSecaoRede(&b, chave, rc, cat[chave], &rel)
	}
	b.WriteString("\n---\n\n## Primeiras pautas\n\n")
	for _, p := range tc.PrimeirasPautas {
		fmt.Fprintf(&b, "- %s — %s\n", p.Titulo, p.Status)
	}
	return b.String(), rel
}

func escreverCabecalhoColar(b *strings.Builder, tc TextosCanal) {
	fmt.Fprintf(b, "# Kit de canal: %s (@%s)\n\n", tc.Nome, tc.Handle)
	b.WriteString("Cole cada campo na rede indicada. Os números entre parênteses são caracteres usados / limite.\n")
	b.WriteString("Limite marcado com ⚠ não foi verificado em página oficial: confira na ajuda da rede.\n\n")
	fmt.Fprintf(b, "## Núcleo da identidade\n\n- Promessa: %s\n- Tagline: %s\n- Tom: %s\n- Símbolo: %s\n- Pilares:\n", tc.Promessa, tc.Tagline, tc.Tom, tc.IdentidadeVisual.Simbolo)
	for _, p := range tc.Pilares {
		fmt.Fprintf(b, "  - %s\n", p)
	}
}

func escreverSecaoRede(b *strings.Builder, chave string, rc RedeCanal, spec RedeSpec, rel *RelatorioTextosCanal) {
	nome := spec.Nome
	if nome == "" {
		nome = chave
	}
	fmt.Fprintf(b, "\n---\n\n## %s\n\n", nome)
	pares := [][2]string{
		{"nome_exibicao", rc.NomeExibicao}, {"handle", rc.Handle}, {"descricao", rc.Descricao},
		{"bio", rc.Bio}, {"tagline", rc.Tagline}, {"sobre", rc.Sobre},
	}
	for _, kv := range pares {
		if kv[1] != "" {
			escreverCampoMedido(b, chave, kv[0], kv[1], spec, rel)
		}
	}
	escreverListasRede(b, rc)
}

func escreverCampoMedido(b *strings.Builder, rede, campo, valor string, spec RedeSpec, rel *RelatorioTextosCanal) {
	n := utf8.RuneCountInString(valor)
	lim, tem := spec.Limites[campo]
	marca := ""
	if tem {
		aviso := ""
		if !spec.Verificado {
			aviso = " ⚠"
		}
		marca = fmt.Sprintf(" (%d/%d)%s", n, lim, aviso)
		if n > lim {
			marca += "  **EXCEDE O LIMITE**"
			rel.Excedidos++
			rel.Campos = append(rel.Campos, rede+"."+campo)
		}
	}
	fmt.Fprintf(b, "**%s**%s\n```\n%s\n```\n\n", campo, marca, valor)
}

func escreverListasRede(b *strings.Builder, rc RedeCanal) {
	listas := []struct {
		nome  string
		itens []string
	}{{"palavras_chave", rc.PalavrasChave}, {"playlists", rc.Playlists}, {"destaques", rc.Destaques}, {"ordem_de_preenchimento", rc.OrdemDePreenchimento}}
	for _, l := range listas {
		if len(l.itens) == 0 {
			continue
		}
		fmt.Fprintf(b, "**%s:**\n", l.nome)
		for _, x := range l.itens {
			fmt.Fprintf(b, "- %s\n", x)
		}
		b.WriteString("\n")
	}
}

// GerarTextosCanalDeDNA monta um TextosCanal completo via Brand Swap, dentro dos limites de cada rede.
func GerarTextosCanalDeDNA(dna redc.ChannelDNA, tokens redc.BrandTokens, handle, nicho string) (TextosCanal, error) {
	nome := strings.TrimSpace(tokens.NomeMarca)
	hd := strings.TrimPrefix(strings.TrimSpace(handle), "@")
	if nome == "" || hd == "" || strings.TrimSpace(nicho) == "" {
		return TextosCanal{}, fmt.Errorf("%w: nome_marca, handle e nicho obrigatorios", ErrPedidoInvalido)
	}
	if len(hd) > 24 {
		hd = hd[:24]
	}
	tom := strings.TrimSpace(tokens.TomPreferido)
	if tom == "" {
		tom = dna.TomDeVoz
	}
	return montarTextosDeDNA(nome, hd, nicho, tom, tokens), nil
}

func montarTextosDeDNA(nome, hd, nicho, tom string, t redc.BrandTokens) TextosCanal {
	pri := corOu(t.CorPrimaria, "#FFD24A")
	tag := cortarRunas("Análises reais sobre "+nicho+".", 40)
	descYT := fmt.Sprintf("%s\n\nInvestigações e explicações sobre %s com números, datas e fontes públicas.\nRoteiros verificados; imagens e narração feitas com apoio de inteligência artificial.\nConteúdo informativo.", tag, nicho)
	return TextosCanal{
		Nome: nome, Handle: hd,
		Promessa: fmt.Sprintf("Explicar %s com dados reais e fontes públicas.", nicho),
		Tagline:  tag, Tom: valorTextoOu(tom, "Direto e baseado em fontes."),
		Pilares: []string{"Impacto real com dados", "Fontes públicas nível A e B", "Explicações visuais sem hype"},
		IdentidadeVisual: IdentidadeCanal{
			Simbolo: "Emblema vetorial minimalista de alto contraste para 98 px",
			Paleta:  map[string]string{"fundo": "#0B0E1C", "ouro": pri, "texto": "#FFFFFF"},
			Fontes:  map[string]string{"titulo": valorTextoOu(t.FonteTitulo, "Bebas Neue"), "corpo": valorTextoOu(t.FonteCorpo, "Inter")},
		},
		DisponibilidadeNome: map[string]string{"youtube": "confirmar no cadastro @" + hd},
		Redes:               redesPadraoDNA(nome, hd, nicho, descYT),
		PrimeirasPautas:     []PautaCanal{{Titulo: "Pauta inaugural sobre " + nicho, Status: "planejada"}},
	}
}

func redesPadraoDNA(nome, hd, nicho, descYT string) map[string]RedeCanal {
	nomeCurto := cortarRunas(nome, 30)
	return map[string]RedeCanal{
		"youtube":   {NomeExibicao: cortarRunas(nome, 50), Handle: hd, Descricao: cortarRunas(descYT, 950), Playlists: []string{"Casos Reais", "Dados e Mercado"}},
		"instagram": {NomeExibicao: nomeCurto, Handle: hd, Bio: cortarRunas("Dados e fontes reais sobre "+nicho+". Apoio de IA.", 140)},
		"tiktok":    {NomeExibicao: nomeCurto, Handle: hd, Bio: cortarRunas("Fontes reais sobre "+nicho+".", 75)},
		"linkedin":  {NomeExibicao: cortarRunas(nome, 95), Tagline: cortarRunas("Análises com fontes públicas sobre "+nicho+".", 115), Sobre: cortarRunas(descYT, 1900)},
	}
}

func cortarRunas(s string, max int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= max {
		return string(r)
	}
	return string(r[:max])
}

func corOu(c, p string) string {
	if strings.TrimSpace(c) == "" {
		return p
	}
	return strings.TrimSpace(c)
}

func valorTextoOu(v, p string) string {
	if strings.TrimSpace(v) == "" {
		return p
	}
	return strings.TrimSpace(v)
}
