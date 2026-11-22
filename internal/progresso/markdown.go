package progresso

import (
	"fmt"
	"strings"
)

// RenderizarMarkdown gera o espelho legível alinhado ao TEMPLATE_PROGRESSO.md.
func RenderizarMarkdown(p Progresso) string {
	p.RecalcularResumo()
	var b strings.Builder
	escreverCabecalhoMD(&b, p)
	escreverSecaoEmAndamentoMD(&b, p.Tickets)
	escreverSecaoFalhasMD(&b, p.Tickets)
	escreverSecaoPendentesMD(&b, p.Tickets)
	escreverSecaoAprovadosMD(&b, p.Tickets)
	return b.String()
}

func escreverCabecalhoMD(b *strings.Builder, p Progresso) {
	fmt.Fprintf(b, "---\nschema: %s\niniciativa: %q\nplano: %q\natualizado_em: %q\n---\n\n",
		p.Schema, p.Iniciativa, p.Plano, p.AtualizadoEm)
	fmt.Fprintf(b, "# Progresso Vivo: %s\n\n", p.Iniciativa)
	b.WriteString("> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.\n\n")
	b.WriteString("## Resumo\n\n")
	b.WriteString("| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |\n")
	b.WriteString("|---:|---:|---:|---:|---:|---:|\n")
	fmt.Fprintf(b, "| %d | %d | %d | %d | %d | %d |\n\n",
		p.Resumo.Total, p.Resumo.Aprovados, p.Resumo.EmAndamento, p.Resumo.Reprovados, p.Resumo.Bloqueados, p.Resumo.Pendentes)
}

func escreverSecaoEmAndamentoMD(b *strings.Builder, itens []ItemTicket) {
	b.WriteString("## Ponto de Retomada (Em Andamento)\n\n")
	b.WriteString("| Ticket | Épico | Responsável | Último registro |\n|---|---|---|---|\n")
	achou := false
	for _, it := range itens {
		if it.Estado != EstadoEmAndamento {
			continue
		}
		achou = true
		fmt.Fprintf(b, "| `%s` | %s | %s | %s |\n", it.ID, it.Epico, limparCelula(it.Responsavel), limparCelula(it.Prova))
	}
	if !achou {
		b.WriteString("| nenhum | - | - | - |\n")
	}
	b.WriteByte('\n')
}

func escreverSecaoFalhasMD(b *strings.Builder, itens []ItemTicket) {
	b.WriteString("## Reprovados ou Bloqueados\n\n")
	b.WriteString("| Ticket | Épico | Estado | Exit | Prova |\n|---|---|---|---:|---|\n")
	achou := false
	for _, it := range itens {
		if it.Estado != EstadoReprovado && it.Estado != EstadoBloqueado {
			continue
		}
		achou = true
		fmt.Fprintf(b, "| `%s` | %s | %s | %d | %s |\n", it.ID, it.Epico, it.Estado, it.ExitCode, limparCelula(it.Prova))
	}
	if !achou {
		b.WriteString("| nenhum | - | - | - | - |\n")
	}
	b.WriteByte('\n')
}

func escreverSecaoPendentesMD(b *strings.Builder, itens []ItemTicket) {
	b.WriteString("## Pendentes\n\n")
	b.WriteString("| Ticket | Épico | Título | Critérios |\n|---|---|---|---|\n")
	achou := false
	for _, it := range itens {
		if it.Estado != EstadoPendente {
			continue
		}
		achou = true
		fmt.Fprintf(b, "| `%s` | %s | %s | %s |\n", it.ID, it.Epico, limparCelula(it.Titulo), strings.Join(it.Criterios, ", "))
	}
	if !achou {
		b.WriteString("| nenhum | - | - | - |\n")
	}
	b.WriteByte('\n')
}

func escreverSecaoAprovadosMD(b *strings.Builder, itens []ItemTicket) {
	b.WriteString("## Aprovados com Prova\n\n")
	b.WriteString("| Ticket | Épico | Exit | Prova | Atualizado em |\n|---|---|---:|---|---|\n")
	achou := false
	for _, it := range itens {
		if it.Estado != EstadoAprovado {
			continue
		}
		achou = true
		fmt.Fprintf(b, "| `%s` | %s | %d | %s | %s |\n", it.ID, it.Epico, it.ExitCode, limparCelula(it.Prova), it.AtualizadoEm)
	}
	if !achou {
		b.WriteString("| nenhum | - | - | - | - |\n")
	}
}

func limparCelula(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "|", "/")
	s = strings.TrimSpace(s)
	if s == "" {
		return "-"
	}
	return s
}
