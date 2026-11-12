package homologacao

import (
	"strings"
)

var sinaisBanimento = []string{"terms of service", "service disabled"}

func TrechoBanimento(texto string) string {
	baixo := strings.ToLower(texto)
	for _, s := range sinaisBanimento {
		i := strings.Index(baixo, s)
		if i < 0 {
			continue
		}
		if i+len(s) > len(texto) {
			return s
		}
		return texto[i : i+len(s)]
	}
	return ""
}

func SinalBanimento(texto string) bool {
	return TrechoBanimento(texto) != ""
}

func AvisoGasto(pares []string) string {
	linhas := []string{"A bateria faz chamadas reais de modelo e cobra por elas."}
	for _, p := range pares {
		if strings.HasPrefix(p, "agy|") {
			linhas = append(linhas, "agy: conta já logada, sem troca de conta, binário chamado direto (nunca pelo gw); risco de banimento por Terms of Service.")
			break
		}
	}
	return strings.Join(linhas, "\n")
}
