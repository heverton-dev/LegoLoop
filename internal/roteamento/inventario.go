package roteamento

import (
	"path/filepath"
	"runtime"

	"legoloop/internal/processo"
	"legoloop/internal/provider"
)

type Via string

const (
	ViaOficial Via = "oficial"
	Via9Router Via = "9router"
)

const (
	EstadoHomologado   = "homologado"
	EstadoReprovado    = "reprovado"
	EstadoPendente     = "pendente"
	EstadoNaoInstalado = "nao_instalado"
)

type Par struct {
	Harness string
	Via     Via
}

type ItemInv struct {
	Par     Par
	Binario string
	Modelos []string
	Estado  string
	PesoMB  int
	Motivo  string
}

type Homologacao struct {
	Estado string
	PesoMB int
	Motivo string
}

func (p Par) Chave() string { return p.Harness + "|" + string(p.Via) }

func modelosOficiais(prov provider.Inventory, harness string) []string {
	src := provider.Source(harness)
	if harness == "claude" {
		src = provider.SourceClaudeCLI
	}
	var out []string
	for _, m := range prov.Models() {
		if m.Source == src {
			out = append(out, m.ID)
		}
	}
	return out
}

func itemInv(prov provider.Inventory, par Par, store map[string]Homologacao) ItemInv {
	item := ItemInv{Par: par, Binario: prov.CLIPath(par.Harness), Estado: EstadoPendente}
	if par.Via == Via9Router {
		item.Modelos = prov.RouterModels
	} else {
		item.Modelos = modelosOficiais(prov, par.Harness)
	}
	if item.Binario == "" {
		item.Estado = EstadoNaoInstalado
		return item
	}
	if h, ok := store[par.Chave()]; ok && h.Estado != "" {
		item.Estado, item.PesoMB, item.Motivo = h.Estado, h.PesoMB, h.Motivo
	}
	return item
}

func nomesHarness(fsys processo.FS, home string, prov provider.Inventory) []string {
	var nomes []string
	visto := map[string]bool{}
	add := func(n string) {
		if !visto[n] {
			visto[n] = true
			nomes = append(nomes, n)
		}
	}
	for _, h := range processo.DetectarHarnesses(fsys, home) {
		add(h.Nome)
	}
	for _, n := range []string{"mimo", "omp", "agy"} {
		if info, err := fsys.Stat(filepath.Join(home, "."+n)); home != "" && err == nil && info.IsDir() {
			add(n)
		}
	}
	for _, c := range prov.CLIs {
		if c.Name != "9router" {
			add(c.Name)
		}
	}
	return nomes
}

func LerInventarioUnificado(fsys processo.FS, home string, prov provider.Inventory, store map[string]Homologacao) ([]ItemInv, error) {
	var out []ItemInv
	for _, nome := range nomesHarness(fsys, home, prov) {
		out = append(out, itemInv(prov, Par{Harness: nome, Via: ViaOficial}, store))
		if len(prov.RouterModels) > 0 {
			out = append(out, itemInv(prov, Par{Harness: nome, Via: Via9Router}, store))
		}
	}
	if len(out) == 0 {
		return nil, ErrInventarioVazio
	}
	return out, nil
}

type Candidato struct{ Harness, Modelo, Esforco string }

type Perfil map[string][]Candidato

type HarnessInv struct {
	Nome    string
	Binario string
	Modelos []string
}

type Inventario struct {
	SO        string
	Harnesses []HarnessInv
}

type Substituicao struct {
	No, Classe    string
	Pedido, Usado Candidato
	Motivo        string
}

func (inv Inventario) Vazio() bool { return len(inv.Harnesses) == 0 }

func LerInventario(fsys processo.FS, home string, binario func(string) string, modelos map[string][]string) Inventario {
	inv := Inventario{SO: runtime.GOOS}
	for _, h := range processo.DetectarHarnesses(fsys, home) {
		item := HarnessInv{Nome: h.Nome, Modelos: modelos[h.Nome]}
		if binario != nil {
			item.Binario = binario(h.Nome)
		}
		inv.Harnesses = append(inv.Harnesses, item)
	}
	return inv
}
