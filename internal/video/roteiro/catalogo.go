package roteiro

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
)

type ConectorFonte struct {
	ID          string   `json:"id"`
	Nome        string   `json:"nome"`
	URL         string   `json:"url"`
	Tipo        string   `json:"tipo"` // "rss" | "hn_json" | "json"
	Perfis      []string `json:"perfis"`
	NivelPadrao string   `json:"nivel_padrao,omitempty"`
}

type CatalogoFontes struct {
	SchemaVersion    int             `json:"schema_version"`
	Perfis           []string        `json:"perfis"`
	Conectores       []ConectorFonte `json:"conectores"`
	DominiosNivelA   []string        `json:"dominios_nivel_a"`
	SufixosNivelA    []string        `json:"sufixos_nivel_a"`
	DominiosSoftware []string        `json:"dominios_software_nivel_a"`
	DominiosNivelB   []string        `json:"dominios_nivel_b"`
	Preprints        []string        `json:"preprints"`
}

type ClassificacaoFonte struct {
	Dominio  string `json:"dominio"`
	Nivel    string `json:"nivel"`
	Preprint bool   `json:"preprint"`
}

// CarregarCatalogo le e valida um arquivo JSON de catalogo de fontes.
func CarregarCatalogo(caminho string) (CatalogoFontes, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return CatalogoFontes{}, fmt.Errorf("%w: %v", ErrPreVooBloqueado, err)
	}
	var cat CatalogoFontes
	if err := json.Unmarshal(dados, &cat); err != nil {
		return CatalogoFontes{}, fmt.Errorf("%w: catalogo JSON invalido: %v", ErrPreVooBloqueado, err)
	}
	if err := cat.Validar(); err != nil {
		return CatalogoFontes{}, err
	}
	return cat, nil
}

// Validar garante que o catalogo possui perfis e conectores declarados.
func (c CatalogoFontes) Validar() error {
	if len(c.Perfis) == 0 || len(c.Conectores) == 0 {
		return fmt.Errorf("%w: catalogo sem perfis ou conectores", ErrPreVooBloqueado)
	}
	return nil
}

// PerfilValido informa se um perfil existe no catalogo declarativo.
func (c CatalogoFontes) PerfilValido(perfil string) bool {
	return slices.Contains(c.Perfis, strings.ToLower(strings.TrimSpace(perfil)))
}

// ConectoresParaProjeto filtra conectores pelos perfis do projeto respeitando exclusoes.
func (c CatalogoFontes) ConectoresParaProjeto(p Projeto) []ConectorFonte {
	var ativos []ConectorFonte
	for _, conn := range c.Conectores {
		if slices.Contains(p.FontesExcluidas, conn.ID) {
			continue
		}
		if conectorCompativel(conn.Perfis, p.PerfisFonte) {
			ativos = append(ativos, conn)
		}
	}
	return ativos
}

func conectorCompativel(perfisConn, perfisProj []string) bool {
	for _, pp := range perfisProj {
		if slices.Contains(perfisConn, strings.ToLower(strings.TrimSpace(pp))) {
			return true
		}
	}
	return false
}

// ClassificarFonte classifica uma URL como A, B ou C (fail-closed) e sinaliza preprint.
func ClassificarFonte(rawURL string, p Projeto, cat CatalogoFontes) ClassificacaoFonte {
	host := extrairDominio(rawURL)
	if host == "" {
		return ClassificacaoFonte{Nivel: "C"}
	}
	if casaListaDominios(host, cat.Preprints) {
		return ClassificacaoFonte{Dominio: host, Nivel: "A", Preprint: true}
	}
	if ehNivelA(host, p, cat) {
		return ClassificacaoFonte{Dominio: host, Nivel: "A"}
	}
	if casaListaDominios(host, cat.DominiosNivelB) {
		return ClassificacaoFonte{Dominio: host, Nivel: "B"}
	}
	if casaListaDominios(host, p.FontesExtras) {
		return ClassificacaoFonte{Dominio: host, Nivel: "B"}
	}
	return ClassificacaoFonte{Dominio: host, Nivel: "C"}
}

func ehNivelA(host string, p Projeto, cat CatalogoFontes) bool {
	if casaListaDominios(host, cat.DominiosNivelA) {
		return true
	}
	for _, suf := range cat.SufixosNivelA {
		if strings.HasSuffix(host, strings.ToLower(suf)) {
			return true
		}
	}
	if slices.Contains(p.PerfisFonte, "tecnologia") && casaListaDominios(host, cat.DominiosSoftware) {
		return true
	}
	return false
}

func extrairDominio(rawURL string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.Hostname() == "" {
		return ""
	}
	return strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
}

func casaListaDominios(host string, dominios []string) bool {
	for _, d := range dominios {
		alvo := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(d)), "www.")
		if alvo == "" {
			continue
		}
		if host == alvo || strings.HasSuffix(host, "."+alvo) {
			return true
		}
	}
	return false
}
