package kraften

import "fmt"

// Gates consolida as provas binárias do nó n04_verify_gates.
type Gates struct {
	CFOld      int // contagem do CNAME antigo (0 = zerou)
	CFNew      int // contagem do CNAME novo (1 = ativo)
	TLSVerdict int // resultado do handshake TLS (0 = válido)
	HTTPCode   int // resposta da rota /health (200 = ok)
}

// VerifyGates reprova o Processo se qualquer prova binária divergir.
func VerifyGates(g Gates) error {
	if g.CFOld != 0 {
		return fmt.Errorf("kraften: gate CF_OLD=%d (esperado 0)", g.CFOld)
	}
	if g.CFNew != 1 {
		return fmt.Errorf("kraften: gate CF_NEW=%d (esperado 1)", g.CFNew)
	}
	if g.TLSVerdict != 0 {
		return fmt.Errorf("kraften: gate TLS=%d (esperado 0)", g.TLSVerdict)
	}
	if g.HTTPCode != 200 {
		return fmt.Errorf("kraften: gate HTTP=%d (esperado 200)", g.HTTPCode)
	}
	return nil
}

// ColetarGates reúne as quatro provas via fronteiras injetáveis e valida.
// hostAntigo vazio significa primeiro deploy (contagem zero trivial).
func ColetarGates(
	contagens func(nome string) int,
	tlsVerify func(host string) int,
	httpCheck func(url string) (int, error),
	healthURL, hostAntigo, hostNovo string,
) (Gates, error) {
	codigo, err := httpCheck(healthURL)
	if err != nil {
		return Gates{}, err
	}
	velho := 0
	if hostAntigo != "" {
		velho = contagens(hostAntigo)
	}
	g := Gates{
		CFOld:      velho,
		CFNew:      contagens(hostNovo),
		TLSVerdict: tlsVerify(hostNovo),
		HTTPCode:   codigo,
	}
	return g, VerifyGates(g)
}
