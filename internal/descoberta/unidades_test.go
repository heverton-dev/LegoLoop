package descoberta

import (
	"reflect"
	"testing"
)

type provedorFalso struct {
	so       string
	unidades []Unidade
}

func (p provedorFalso) SO() string                { return p.so }
func (p provedorFalso) ListarUnidades() []Unidade { return p.unidades }

func TestUnidadesFixasPorSO(t *testing.T) {
	casos := []struct {
		nome             string
		prov             Provedor
		fixasEsperadas   []string
		todasEsperadas   []string
		recusadaEsperada string
	}{
		{
			nome: "windows",
			prov: provedorFalso{
				so: "windows",
				unidades: []Unidade{
					{Caminho: "D:", Tipo: TipoFixa},
					{Caminho: "C:", Tipo: TipoFixa},
					{Caminho: "E:", Tipo: TipoRemovivel},
				},
			},
			fixasEsperadas:   []string{"C:", "D:"},
			todasEsperadas:   []string{"C:", "D:", "E:"},
			recusadaEsperada: "E:",
		},
		{
			nome: "linux",
			prov: provedorFalso{
				so: "linux",
				unidades: []Unidade{
					{Caminho: "/mnt/dados", Tipo: TipoFixa},
					{Caminho: "/media/pen", Tipo: TipoRemovivel},
				},
			},
			fixasEsperadas:   []string{"/mnt/dados"},
			todasEsperadas:   []string{"/media/pen", "/mnt/dados"},
			recusadaEsperada: "/media/pen",
		},
		{
			nome: "macos",
			prov: provedorFalso{
				so: "darwin",
				unidades: []Unidade{
					{Caminho: "/Volumes/Disco", Tipo: TipoFixa},
					{Caminho: "/Volumes/Pen", Tipo: TipoRemovivel},
				},
			},
			fixasEsperadas:   []string{"/Volumes/Disco"},
			todasEsperadas:   []string{"/Volumes/Disco", "/Volumes/Pen"},
			recusadaEsperada: "/Volumes/Pen",
		},
	}

	for _, tc := range casos {
		t.Run(tc.nome, func(t *testing.T) {
			verificarCasoSO(t, tc.prov, tc.fixasEsperadas, tc.todasEsperadas, tc.recusadaEsperada)
		})
	}
}

func verificarCasoSO(t *testing.T, prov Provedor, fixas, todas []string, recusada string) {
	t.Helper()
	gotUnidades := Unidades(prov)
	if !reflect.DeepEqual(gotUnidades, fixas) {
		t.Fatalf("Unidades() = %v, esperado %v", gotUnidades, fixas)
	}
	raizes, recusadas := Raizes("", nil, prov, false)
	if !reflect.DeepEqual(raizes, fixas) {
		t.Fatalf("Raizes(confiar=false) = %v, esperado %v", raizes, fixas)
	}
	if len(recusadas) != 1 || recusadas[0].Caminho != recusada || recusadas[0].Motivo == "" {
		t.Fatalf("recusadas(confiar=false) = %+v, esperado [%s com motivo]", recusadas, recusada)
	}
	raizesConf, recConf := Raizes("", nil, prov, true)
	if !reflect.DeepEqual(raizesConf, todas) {
		t.Fatalf("Raizes(confiar=true) = %v, esperado %v", raizesConf, todas)
	}
	if len(recConf) != 0 {
		t.Fatalf("recusadas(confiar=true) = %+v, esperado vazio", recConf)
	}
}

func TestRaizesFiltroAlvoEHome(t *testing.T) {
	prov := provedorFalso{
		so: "windows",
		unidades: []Unidade{
			{Caminho: "C:", Tipo: TipoFixa},
			{Caminho: "D:", Tipo: TipoFixa},
			{Caminho: "Z:", Tipo: TipoRede},
		},
	}
	raizes, recusadas := Raizes("C:\\Users\\dev", []string{"d:"}, prov, false)
	if !reflect.DeepEqual(raizes, []string{"D:"}) {
		t.Fatalf("Raizes com alvo = %v, esperado [D:]", raizes)
	}
	if len(recusadas) != 0 {
		t.Fatalf("recusadas fora do alvo = %+v, esperado vazio", recusadas)
	}
	raizesRede, recRede := Raizes("", []string{"Z:"}, prov, false)
	if len(raizesRede) != 0 || len(recRede) != 1 || recRede[0].Caminho != "Z:" {
		t.Fatalf("alvo de rede = raizes=%v recusadas=%+v", raizesRede, recRede)
	}
}

func TestRaizesHomeEDeterminismo(t *testing.T) {
	prov := provedorFalso{
		so: "linux",
		unidades: []Unidade{
			{Caminho: "/mnt/b", Tipo: TipoFixa},
			{Caminho: "/mnt/a", Tipo: TipoFixa},
			{Caminho: "/mnt/a", Tipo: TipoFixa},
			{Caminho: "/media/usb", Tipo: TipoFixa},
		},
	}
	raizes, recusadas := Raizes("/home/dev", nil, prov, false)
	esp := []string{"/home/dev", "/mnt/a", "/mnt/b"}
	if !reflect.DeepEqual(raizes, esp) {
		t.Fatalf("Raizes = %v, esperado %v", raizes, esp)
	}
	if len(recusadas) != 1 || recusadas[0].Caminho != "/media/usb" {
		t.Fatalf("recusadas = %+v, esperado [/media/usb]", recusadas)
	}
}

func TestProvedorPadrao(t *testing.T) {
	prov := ProvedorPadrao()
	if prov == nil || prov.SO() == "" {
		t.Fatalf("ProvedorPadrao invalido: %#v", prov)
	}
	_ = Unidades(prov)
}
