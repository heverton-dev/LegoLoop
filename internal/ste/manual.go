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
	b.WriteString("└" + strings.Repeat("─", 81) + "┘\n")
	return b.String()
}

func formatLinhaBox(texto string) string {
	runes := []rune(texto)
	const maxConteudo = 79
	if len(runes) > maxConteudo {
		runes = runes[:maxConteudo]
	}
	pad := maxConteudo - len(runes)
	return "│  " + string(runes) + strings.Repeat(" ", pad) + "│\n"
}

func topoManual() string {
	var b strings.Builder
	b.WriteString("┌" + strings.Repeat("─", 81) + "┐\n")
	b.WriteString(formatLinhaBox("MANUAL DE USO DO ALTERNADOR UNIVERSAL STE-PT (Go Native)"))
	b.WriteString("├" + strings.Repeat("─", 81) + "┤\n")
	b.WriteString(formatLinhaBox("Uso: stept <comando> [opções]  ou  legoloop stept <comando> [opções]"))
	b.WriteString(formatLinhaBox(""))
	return b.String()
}

func secaoEstado() string {
	var b strings.Builder
	b.WriteString(formatLinhaBox("ESTADO & SINCRONIZAÇÃO"))
	b.WriteString(formatLinhaBox("  stept                               Sem argumentos: equivale a stept status"))
	b.WriteString(formatLinhaBox("  stept status                        Matriz viva de harnesses e governança"))
	b.WriteString(formatLinhaBox("  stept status --json                 Mesma matriz em JSON estruturado"))
	b.WriteString(formatLinhaBox("  stept scan                          Varre harnesses no SO (sem gravar)"))
	b.WriteString(formatLinhaBox("  stept scan --json                   Lista harnesses instalados em JSON"))
	b.WriteString(formatLinhaBox("  stept manual | help | -h            Exibe este manual operacional formatado"))
	b.WriteString(formatLinhaBox(""))
	return b.String()
}

func secaoControle() string {
	var b strings.Builder
	b.WriteString(formatLinhaBox("CONTROLE DE GOVERNANÇA (MODOS & POLÍTICAS)"))
	b.WriteString(formatLinhaBox("  stept on                            Ativa modo prova globalmente (Padrão)"))
	b.WriteString(formatLinhaBox("  stept on --harness <id>             Ativa modo prova em um harness"))
	b.WriteString(formatLinhaBox("  stept off                           Desativa governança e remove regras"))
	b.WriteString(formatLinhaBox("  stept off --harness <id>            Desativa e remove regras de um harness"))
	b.WriteString(formatLinhaBox("  stept set --modo <modo>             Altera modo global (prova|puro|hibrido)"))
	b.WriteString(formatLinhaBox("  stept set --harness <id> --modo <m> Define override de modo para um harness"))
	b.WriteString(formatLinhaBox("  stept set --pausa <politica>        Define a política de interrupção padrão"))
	b.WriteString(formatLinhaBox("  stept set --harness <id> --pausa <p>Define override de pausa para um harness"))
	b.WriteString(formatLinhaBox(""))
	return b.String()
}

func secaoModos() string {
	var b strings.Builder
	b.WriteString(formatLinhaBox("MODOS DISPONÍVEIS"))
	b.WriteString(formatLinhaBox("  prova       Padrão: uma frase por linha, prova inline, Recomendação, ~470 tokens por turno"))
	b.WriteString(formatLinhaBox("  puro        ASD-STE100 Puro: Linguagem técnica direta sem Tríade [S][T][E]"))
	b.WriteString(formatLinhaBox("  hibrido     Caveman Thinking + ASD-STE100 + Tríade [S][T][E] (legado)"))
	b.WriteString(formatLinhaBox("  stept       Alias de hibrido (mesmo texto), mantido por compatibilidade"))
	b.WriteString(formatLinhaBox("  caveman     Economia Máxima: Thinking e Output telegráficos ultra-curtos"))
	b.WriteString(formatLinhaBox("  off         Desativado: Exclui stept-rules.md sem deixar resíduos no disco"))
	b.WriteString(formatLinhaBox(""))
	return b.String()
}

func secaoPausas() string {
	var b strings.Builder
	b.WriteString(formatLinhaBox("POLÍTICAS DE PAUSA"))
	b.WriteString(formatLinhaBox("  destrutivo  Pausa e pede confirmação apenas em ações de risco (Padrão)"))
	b.WriteString(formatLinhaBox("  sempre      Pausa antes de qualquer comando com efeito no sistema"))
	b.WriteString(formatLinhaBox("  nunca       Exibe o cartão informativo e segue sem interromper a execução"))
	b.WriteString(formatLinhaBox(""))
	return b.String()
}

func secaoHarnesses() string {
	var b strings.Builder
	b.WriteString(formatLinhaBox("HARNESSES SUPORTADOS (UNIVERSAL)"))
	b.WriteString(formatLinhaBox("  gemini      Google Gemini / Antigravity (~/.gemini/config/rules)"))
	b.WriteString(formatLinhaBox("  claude      Claude Code (~/.claude/rules)"))
	b.WriteString(formatLinhaBox("  cursor      Cursor IDE (~/.cursor/rules)"))
	b.WriteString(formatLinhaBox("  windsurf    Windsurf IDE (~/.windsurf/rules)"))
	b.WriteString(formatLinhaBox("  vscode      VS Code (~/.vscode/rules)"))
	b.WriteString(formatLinhaBox("  opencode    OpenCode (~/.config/opencode/rules)"))
	b.WriteString(formatLinhaBox("  mimo        MimoCode / MemoCode (~/.mimocode/rules)"))
	b.WriteString(formatLinhaBox("  omp         OMP CLI (~/.omp/rules)"))
	b.WriteString(formatLinhaBox(""))
	return b.String()
}

func secaoInstalacao() string {
	var b strings.Builder
	b.WriteString(formatLinhaBox("INSTALAÇÃO STANDALONE"))
	b.WriteString(formatLinhaBox("  stept install                       Instala binário no PATH de usuário"))
	return b.String()
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
