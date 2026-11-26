package render

import (
	"errors"
	"strings"
	"testing"
)

func TestBGMGanhoPadraoCalibrado(t *testing.T) {
	c := ConfigRender{SchemaVersion: 1}.Normalizar()
	if c.BGMGanhoDB == nil || *c.BGMGanhoDB != BGMGanhoPadraoDB || BGMGanhoPadraoDB != 3 {
		t.Fatalf("padrao do ganho do BGM deve ser +3 dB: %+v", c.BGMGanhoDB)
	}
	v := -14.0
	c = ConfigRender{SchemaVersion: 1, BGMGanhoDB: &v}.Normalizar()
	if *c.BGMGanhoDB != -14 {
		t.Fatalf("normalizacao nao deve sobrescrever o ganho: %v", *c.BGMGanhoDB)
	}
	zero := 0.0
	c = ConfigRender{SchemaVersion: 1, BGMGanhoDB: &zero}.Normalizar()
	if *c.BGMGanhoDB != 0 {
		t.Fatalf("ganho 0 dB explicito deve ser respeitado: %v", *c.BGMGanhoDB)
	}
}

func TestBGMGanhoValidacao(t *testing.T) {
	for _, v := range []float64{-40.1, 12.1} {
		v := v
		c := ConfigRender{SchemaVersion: 1, BGMGanhoDB: &v}
		if err := c.Validar(); !errors.Is(err, ErrSchemaInvalido) {
			t.Errorf("ganho %v deve ser invalido: %v", v, err)
		}
	}
	for _, v := range []float64{-40, 0, 12} {
		v := v
		c := ConfigRender{SchemaVersion: 1, BGMGanhoDB: &v}
		if err := c.Validar(); err != nil {
			t.Errorf("ganho %v deve validar: %v", v, err)
		}
	}
}

func TestFiltroMixUsaGanhoEDuckingCalibrados(t *testing.T) {
	f := filtroMix(EntradaMix{Duracao: 10, CanaisVoz: 1, GanhoBGMdB: 3})
	for _, trecho := range []string{"volume=3.0dB", "threshold=0.04:ratio=6:attack=20:release=600:makeup=1"} {
		if !strings.Contains(f, trecho) {
			t.Errorf("filtro sem %q: %s", trecho, f)
		}
	}
	f = filtroMix(EntradaMix{Duracao: 10, CanaisVoz: 1, GanhoBGMdB: -14.5})
	if !strings.Contains(f, "volume=-14.5dB") {
		t.Errorf("ganho do usuario ignorado: %s", f)
	}
}
