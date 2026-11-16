package fechar

import (
	"fmt"
	"strings"

	"legoloop/internal/processo"
)

// naturezaDoProcesso lê a natureza do processo no cofre. Arquivo ausente ou ilegível
// devolve vazio e o fechamento mantém o comportamento atual (compatibilidade v1).
func naturezaDoProcesso(vault, canon, slug string) string {
	for _, nome := range []string{canon, slug} {
		if nome == "" {
			continue
		}
		p, err := processo.Ler(processo.CaminhoNoCofre(vault, nome))
		if err == nil {
			return strings.TrimSpace(p.Natureza)
		}
	}
	return ""
}

// etapa5Aplica diz se a recompilação do binário Go vale para a natureza do processo.
// Só naturezas Go-cêntricas (software, desconhecida ou vazia) têm binário a recompilar.
func etapa5Aplica(natureza string) bool {
	switch natureza {
	case "", "software", "desconhecida":
		return true
	}
	return false
}

// runEtapaRebuild executa a etapa 5 ou a marca como não aplicável conforme a natureza.
func runEtapaRebuild(d *Deps, res *Result, dry bool, natureza string) {
	if etapa5Aplica(natureza) {
		runEtapa(d, res, 5, "Recompilação Fresca do Binário", func() (string, error) {
			return fecharRebuild(d, dry)
		})
		return
	}
	res.Etapas = append(res.Etapas, EtapaStatus{
		Etapa:   5,
		Nome:    "Recompilação Fresca do Binário",
		Status:  "skipped",
		Detalhe: fmt.Sprintf("não aplicável para a natureza %s: sem binário Go", natureza),
	})
}
