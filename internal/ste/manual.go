package ste

import (
	"fmt"
	"strings"
)

// TextoManual retorna o manual operacional formatado em moldura box-drawing idêntica ao gw.
func TextoManual() string {
	var b strings.Builder
	b.WriteString(topoManual())
	b.WriteString(secaoEstado())
	b.WriteString(secaoControle())
	b.WriteString(secaoModos())
	b.WriteString(secaoPausas())
	b.WriteString(secaoHarnesses())
	b.WriteString(secaoInstalacao())
	b.WriteString("└─────────────────────────────────────────────────────────────────────────────────┘\n")
	return b.String()
}

func topoManual() string {
	return `┌─────────────────────────────────────────────────────────────────────────────────┐
│  MANUAL DE USO DO ALTERNADOR UNIVERSAL STE-PT (Go Native)                      │
├─────────────────────────────────────────────────────────────────────────────────┤
│  Uso: stept <comando> [opções]  ou  legoloop stept <comando> [opções]          │
│                                                                                 │
`
}

func secaoEstado() string {
	return `│  ESTADO & SINCRONIZAÇÃO                                                         │
│    stept                               Sem argumentos: equivale a stept status  │
│    stept status                        Matriz viva de harnesses e governança    │
│    stept status --json                 Mesma matriz em JSON estruturado         │
│    stept scan                          Varre harnesses no SO (sem gravar)       │
│    stept manual | help | -h            Exibe este manual operacional formatado  │
│                                                                                 │
`
}

func secaoControle() string {
	return `│  CONTROLE DE GOVERNANÇA (MODOS & POLÍTICAS)                                     │
│    stept on                            Ativa modo hibrido globalmente (Padrão)  │
│    stept on --harness <id>             Ativa modo hibrido em um harness         │
│    stept off                           Desativa governança e remove regras      │
│    stept off --harness <id>            Desativa e remove regras de um harness   │
│    stept set --modo <modo>             Altera o modo global padrão              │
│    stept set --harness <id> --modo <m> Define override para um harness          │
│        --pausa <nunca|destrutivo|sempre>  Define a política de interrupção      │
│                                                                                 │
`
}

func secaoModos() string {
	return `│  MODOS DISPONÍVEIS                                                              │
│    hibrido     Padrão Ouro: Caveman Thinking + Cartão STE-PT obrigatório        │
│    caveman     Economia máxima: Thinking e Output telegráficos ultra-curtos     │
│    stept       Governança formal: Cartão ASD-STE100 estruturado em ações        │
│    off         Desativado: Exclui stept-rules.md sem deixar resíduos            │
│                                                                                 │
`
}

func secaoPausas() string {
	return `│  POLÍTICAS DE PAUSA                                                             │
│    destrutivo  Pausa e pede confirmação apenas em ações de risco (Padrão)       │
│    sempre      Pausa antes de qualquer comando com efeito                       │
│    nunca       Exibe o cartão informativo e segue sem interromper               │
│                                                                                 │
`
}

func secaoHarnesses() string {
	return `│  HARNESSES SUPORTADOS                                                           │
│    gemini      Google Gemini / Antigravity (~/.gemini/config/rules)             │
│    claude      Claude Code (~/.claude/rules)                                    │
│    omp         OMP CLI (~/.omp/rules)                                           │
│    opencode    OpenCode (~/.config/opencode/rules)                              │
│    mimo        MimoCode (~/.mimocode/rules)                                     │
│                                                                                 │
`
}

func secaoInstalacao() string {
	return `│  INSTALAÇÃO STANDALONE                                                          │
│    stept install                       Instala binário no PATH de usuário       │
`
}

// RenderTabelaBox formata a matriz de harnesses em tabela compacta com bordas Unicode idêntica ao gw.
func RenderTabelaBox(statusList []HarnessStatus, cfgPath string) string {
	var b strings.Builder
	b.WriteString("┌─────────────────────────────┬───────────┬─────────┬─────────┬────────────┬─────────┐\n")
	b.WriteString("│ HARNESS                     │ INSTALADO │ STATUS  │ MODO    │ PAUSA      │ ORIGEM  │\n")
	b.WriteString("├─────────────────────────────┼───────────┼─────────┼─────────┼────────────┼─────────┤\n")
	for _, s := range statusList {
		inst := "não"
		if s.Instalado {
			inst = "sim"
		}
		stat := "[INATIVO]"
		if s.Ativo {
			stat = "[ATIVO]"
		}
		fmt.Fprintf(&b, "│ %-27s │ %-9s │ %-7s │ %-7s │ %-10s │ %-7s │\n",
			truncar(s.Nome, 27), inst, stat, s.Perfil.Modo, s.Perfil.Pausa, s.Origem)
	}
	b.WriteString("└─────────────────────────────┴───────────┴─────────┴─────────┴────────────┴─────────┘\n")
	if cfgPath != "" {
		fmt.Fprintf(&b, "CONFIG: %s | GOVERNANÇA: [OPERACIONAL]\n", cfgPath)
	}
	return b.String()
}

func truncar(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
