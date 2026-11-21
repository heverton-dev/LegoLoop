package descoberta

import (
	"sort"
	"strings"
)

// TipoUnidade classifica o tipo de uma unidade do sistema operacional.
type TipoUnidade string

const (
	TipoFixa      TipoUnidade = "fixa"
	TipoRemovivel TipoUnidade = "removivel"
	TipoRede      TipoUnidade = "rede"
	TipoOutro     TipoUnidade = "outro"
)

// Unidade representa um volume ou ponto de montagem descoberto.
type Unidade struct {
	Caminho string
	Tipo    TipoUnidade
}

// Recusa registra uma unidade descartada na descoberta e o motivo.
type Recusa struct {
	Caminho string
	Motivo  string
}

// Provedor fornece as unidades do sistema operacional de forma somente leitura.
type Provedor interface {
	SO() string
	ListarUnidades() []Unidade
}

// Unidades devolve apenas as unidades fixas do provedor, ordenadas e sem duplicatas.
func Unidades(prov Provedor) []string {
	if prov == nil {
		prov = ProvedorPadrao()
	}
	so := strings.ToLower(strings.TrimSpace(prov.SO()))
	vistos := make(map[string]bool)
	var fixas []string
	for _, u := range prov.ListarUnidades() {
		caminho := limparCaminho(u.Caminho)
		if caminho == "" || classificarUnidade(so, caminho, u.Tipo) != TipoFixa {
			continue
		}
		if !vistos[caminho] {
			vistos[caminho] = true
			fixas = append(fixas, caminho)
		}
	}
	sort.Strings(fixas)
	return fixas
}

// Raizes calcula as raizes de busca e as unidades recusadas de forma deterministica.
func Raizes(home string, alvo []string, prov Provedor, confiarRemovivel bool) ([]string, []Recusa) {
	if prov == nil {
		prov = ProvedorPadrao()
	}
	so := strings.ToLower(strings.TrimSpace(prov.SO()))
	alvos := normalizarAlvos(so, alvo)
	raizesSet := make(map[string]bool)
	recusasMap := make(map[string]Recusa)
	for _, u := range prov.ListarUnidades() {
		avaliarUnidade(so, u, alvos, confiarRemovivel, raizesSet, recusasMap)
	}
	if len(alvos) == 0 {
		if h := limparCaminho(home); h != "" {
			raizesSet[h] = true
		}
	}
	return ordenarChaves(raizesSet), ordenarRecusas(recusasMap)
}

func avaliarUnidade(so string, u Unidade, alvos []string, confiar bool, raizes map[string]bool, recusas map[string]Recusa) {
	caminho := limparCaminho(u.Caminho)
	if caminho == "" || !correspondeAlvo(so, caminho, alvos) {
		return
	}
	tipo := classificarUnidade(so, caminho, u.Tipo)
	if tipo == TipoFixa {
		raizes[caminho] = true
		return
	}
	if (tipo == TipoRemovivel || tipo == TipoRede) && confiar {
		raizes[caminho] = true
		return
	}
	recusas[caminho] = Recusa{Caminho: caminho, Motivo: motivoRecusa(tipo)}
}

func classificarUnidade(so, caminho string, tipo TipoUnidade) TipoUnidade {
	if so == "linux" && (caminho == "/media" || strings.HasPrefix(caminho, "/media/")) {
		return TipoRemovivel
	}
	if tipo != "" {
		return tipo
	}
	if so == "linux" && (caminho == "/mnt" || strings.HasPrefix(caminho, "/mnt/")) {
		return TipoFixa
	}
	if so == "darwin" && strings.HasPrefix(caminho, "/Volumes/") {
		return TipoFixa
	}
	return TipoOutro
}

func motivoRecusa(tipo TipoUnidade) string {
	switch tipo {
	case TipoRemovivel:
		return "unidade removivel recusada sem confiarRemovivel"
	case TipoRede:
		return "unidade de rede recusada sem confiarRemovivel"
	default:
		return "tipo de unidade nao suportado"
	}
}

func limparCaminho(caminho string) string {
	c := strings.TrimSpace(caminho)
	if c == "" {
		return ""
	}
	for len(c) > 1 && (strings.HasSuffix(c, "/") || strings.HasSuffix(c, "\\")) {
		if len(c) == 3 && c[1] == ':' {
			return c[:2]
		}
		c = c[:len(c)-1]
	}
	return c
}

func normalizarAlvos(so string, alvo []string) []string {
	var out []string
	for _, item := range alvo {
		c := normalizarComparacao(so, item)
		if c != "" {
			out = append(out, c)
		}
	}
	return out
}

func correspondeAlvo(so, caminho string, alvos []string) bool {
	if len(alvos) == 0 {
		return true
	}
	norm := normalizarComparacao(so, caminho)
	for _, a := range alvos {
		if norm == a || strings.HasPrefix(norm, a+"/") || strings.HasPrefix(a, norm+"/") {
			return true
		}
	}
	return false
}

func normalizarComparacao(so, caminho string) string {
	c := strings.ReplaceAll(limparCaminho(caminho), "\\", "/")
	if so == "windows" {
		return strings.ToLower(c)
	}
	return c
}

func ordenarChaves(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func ordenarRecusas(mapa map[string]Recusa) []Recusa {
	out := make([]Recusa, 0, len(mapa))
	for _, r := range mapa {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Caminho < out[j].Caminho
	})
	return out
}
