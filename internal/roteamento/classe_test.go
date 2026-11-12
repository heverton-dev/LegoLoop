package roteamento

import (
	"errors"
	"strings"
	"testing"
)

func TestClassePorRegra(t *testing.T) {
	casos := []struct {
		nome   string
		e      EntradaClasse
		classe string
		regra  string
	}{
		{"nano", EntradaClasse{ArquivosSaida: 1, TokensEntrada: 100, Pacotes: 1}, ClasseNano, "nano:"},
		{"micro por arquivos", EntradaClasse{ArquivosSaida: 3, Pacotes: 1}, ClasseMicro, "micro:"},
		{"micro por tokens", EntradaClasse{ArquivosSaida: 1, TokensEntrada: 8000}, ClasseMicro, "micro:"},
		{"grande por arquivos", EntradaClasse{ArquivosSaida: 6}, ClasseGrande, ClasseGrande + ":"},
		{"grande por tokens", EntradaClasse{ArquivosSaida: 1, TokensEntrada: 16001}, ClasseGrande, ClasseGrande + ":"},
		{"grande por pacotes", EntradaClasse{ArquivosSaida: 1, Pacotes: 2}, ClasseGrande, ClasseGrande + ":"},
		{"extreme", EntradaClasse{ArquivosSaida: 1, Destrutivo: true}, ClasseExtreme, "extreme:"},
		{"sobrescrita", EntradaClasse{ArquivosSaida: 1, Sobrescrita: ClasseGrande, JustificativaS: "motivo"}, ClasseGrande, ClasseGrande + ":"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			r, err := Classificar(c.e)
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if r.Classe != c.classe {
				t.Errorf("classe = %q, esperado %q", r.Classe, c.classe)
			}
			if !strings.HasPrefix(r.ClassRule, c.regra) {
				t.Errorf("class_rule = %q, esperado prefixo %q", r.ClassRule, c.regra)
			}
		})
	}
}

func TestClasseRecusaETrataSobrescrita(t *testing.T) {
	if _, err := Classificar(EntradaClasse{ArquivosSaida: 14}); !errors.Is(err, ErrNoGrande) {
		t.Errorf("14 arquivos: erro = %v", err)
	}
	if _, err := Classificar(EntradaClasse{TokensEntrada: 32001}); !errors.Is(err, ErrNoGrande) {
		t.Errorf("tokens altos: erro = %v", err)
	}
	if _, err := Classificar(EntradaClasse{ArquivosSaida: 1, Sobrescrita: ClasseGrande}); err == nil {
		t.Error("sobrescrita sem justificativa deveria falhar")
	}
	if _, err := Classificar(EntradaClasse{ArquivosSaida: 1, Sobrescrita: "x", JustificativaS: "m"}); err == nil {
		t.Error("classe invalida deveria falhar")
	}
	r, _ := Classificar(EntradaClasse{ArquivosSaida: 1, Sobrescrita: ClasseMicro, JustificativaS: "motivo"})
	if r.ClassOverrideReason != "motivo" {
		t.Errorf("class_override_reason = %q", r.ClassOverrideReason)
	}
}
