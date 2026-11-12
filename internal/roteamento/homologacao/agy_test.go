package homologacao

import (
	"context"
	"strings"
	"testing"
	"time"
)

type agyFalso struct {
	chamadas int
}

func (a *agyFalso) Executar(ctx context.Context, dir, prompt string, teto time.Duration) (Saida, error) {
	a.chamadas++
	s := Saida{CustoUSD: 0.01, Tokens: 10, Duracao: time.Millisecond, Texto: "ok"}
	if a.chamadas == 2 {
		s.Texto = "Your account violated the Terms of Service"
	}
	return s, nil
}

func TestAgyContaFixaEParadaTos(t *testing.T) {
	if !strings.Contains(AvisoGasto([]string{"agy|oficial"}), "agy") {
		t.Fatal("o aviso de gasto deve citar o agy")
	}
	h := &agyFalso{}
	r := Avaliar(context.Background(), h, t.TempDir(), time.Second)
	if h.chamadas != 2 {
		t.Fatalf("chamadas = %d, quer 2", h.chamadas)
	}
	if r.Estado != EstadoReprovado || !strings.HasPrefix(r.Motivo, "banimento:") {
		t.Fatalf("estado = %q, motivo = %q", r.Estado, r.Motivo)
	}
	if !r.Criterios["C1"] || !r.Criterios["C3"] {
		t.Fatalf("critérios medidos antes do sinal foram perdidos: %v", r.Criterios)
	}
	if _, ok := r.Criterios["C5"]; ok {
		t.Fatalf("C5 não deveria ser medido: %v", r.Criterios)
	}
}

func TestSinalBanimento(t *testing.T) {
	casos := map[string]bool{
		"TERMS OF SERVICE violado": true,
		"Service Disabled":         true,
		"tudo certo":               false,
		"":                         false,
	}
	for texto, quer := range casos {
		if got := SinalBanimento(texto); got != quer {
			t.Errorf("SinalBanimento(%q) = %v", texto, got)
		}
	}
}

func TestAvisoGastoCitaAgySoQuandoHaPar(t *testing.T) {
	sem := AvisoGasto([]string{"claude|oficial"})
	if !strings.Contains(sem, "chamadas reais") || strings.Contains(sem, "agy") {
		t.Fatalf("aviso sem agy = %q", sem)
	}
	com := AvisoGasto([]string{"claude|oficial", "agy|oficial"})
	linhas := strings.Split(com, "\n")
	if len(linhas) != 2 || !strings.Contains(linhas[1], "agy") || !strings.Contains(linhas[1], "gw") {
		t.Fatalf("aviso com agy = %q", com)
	}
}
