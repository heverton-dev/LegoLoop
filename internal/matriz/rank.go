package matriz

import (
	"fmt"
	"slices"
	"strings"

	"legoloop/internal/reputation"
)

// Needs é o que o espelho da entrevista pede, em termos do vocabulário.
type Needs struct {
	Natureza string
	Etapas   []string
	Casos    []string
}

// Ranked é uma peça avaliada para um Needs. Motivo é determinístico: diz o que casou ou qual exclusão agiu.
type Ranked struct {
	Entrada  Entrada
	Aptidao  float64 // (0..1] = |casos casados| / |casos pedidos| (1 se nenhum caso foi pedido); 0 se excluída
	Motivo   string
	Excluida bool

	naturezaExplicita bool // a peça declara a natureza pedida (desempata a favor dela)
}

const naturezaQualquer = "qualquer"

// Rank avalia as peças contra o Needs. É pura: não altera es. O Score só desempata; quem decide é a aptidão
// ao caso de uso. Ordem: aptidão desc, natureza explícita antes de coringa, perfil declarado ou confirmado antes
// de inferido, Score desc, Ref asc; o motivo cita a origem do perfil. As excluídas
// vêm depois, por Ref. Peça sem perfil e de tier tóxico nunca entram no resultado.
func Rank(n Needs, es []Entrada, minAptidao float64) []Ranked {
	var out []Ranked
	for _, e := range es {
		if e.Perfil == nil || e.Rating.Tier == reputation.TierToxica {
			continue
		}
		out = append(out, avaliar(n, e, minAptidao))
	}
	slices.SortFunc(out, comparar)
	return out
}

func avaliar(n Needs, e Entrada, minAptidao float64) Ranked {
	r := Ranked{Entrada: e}
	if m := exclusaoDura(n, e.Perfil); m != "" {
		r.Excluida, r.Motivo = true, m
		return r
	}
	etapas := casados(n.Etapas, e.Perfil.Etapas)
	if len(n.Etapas) > 0 && len(etapas) == 0 {
		r.Excluida, r.Motivo = true, fmt.Sprintf("etapa: a peça faz %s, não %s", listar(e.Perfil.Etapas), listar(n.Etapas))
		return r
	}
	casos := casados(n.Casos, e.Perfil.CasosDeUso)
	apt := 1.0
	if len(n.Casos) > 0 {
		if len(casos) == 0 { // só a etapa e a natureza casaram: a aptidão ao pedido é zero
			r.Excluida, r.Motivo = true, fmt.Sprintf("caso de uso: a peça não cobre %s", listar(n.Casos))
			return r
		}
		apt = float64(len(casos)) / float64(len(n.Casos))
	}
	if apt < minAptidao {
		r.Excluida = true
		r.Motivo = fmt.Sprintf("aptidão abaixo do mínimo: %.2f < %.2f", apt, minAptidao)
		return r
	}
	r.Aptidao = apt
	r.naturezaExplicita = declara(e.Perfil.Naturezas, n.Natureza)
	r.Motivo = motivoEscolha(n, r.naturezaExplicita, etapas, casos)
	if o := e.Perfil.Origem; o != "" {
		r.Motivo += "; perfil: " + o
	}
	return r
}

// inferido diz se o perfil é só palpite do cérebro (declarado e confirmado valem mais no empate).
func (r Ranked) inferido() bool { return r.Entrada.Perfil.Origem == OrigemInferido }

// exclusaoDura devolve o motivo se o perfil proíbe o uso pedido (quando_nao_usar) ou serve a outra natureza.
func exclusaoDura(n Needs, p *PerfilDeUso) string {
	for _, termo := range append([]string{n.Natureza}, n.Casos...) {
		if termo != "" && declara(p.QuandoNaoUsar, termo) {
			return "quando_nao_usar: " + termo
		}
	}
	if n.Natureza != "" && n.Natureza != naturezaQualquer && !naturezaServe(p.Naturezas, n.Natureza) {
		return fmt.Sprintf("natureza: a peça serve a %s, não a %s", listar(p.Naturezas), n.Natureza)
	}
	return ""
}

// naturezaServe: natureza vazia ou "qualquer" vale como coringa; senão a natureza pedida tem de estar declarada.
func naturezaServe(declaradas []string, pedida string) bool {
	return len(declaradas) == 0 || declara(declaradas, naturezaQualquer) || declara(declaradas, pedida)
}

func motivoEscolha(n Needs, explicita bool, etapas, casos []string) string {
	var partes []string
	switch {
	case n.Natureza == "" || n.Natureza == naturezaQualquer:
	case explicita:
		partes = append(partes, "natureza: "+n.Natureza)
	default:
		partes = append(partes, "natureza: qualquer")
	}
	if len(etapas) > 0 {
		partes = append(partes, "etapas: "+listar(etapas))
	}
	if len(casos) > 0 {
		partes = append(partes, "casos: "+listar(casos))
	}
	return strings.Join(partes, "; ")
}

// comparar ordena: escolhidas antes das excluídas; entre as escolhidas, aptidão, natureza explícita, Score, Ref.
func comparar(a, b Ranked) int {
	switch {
	case a.Excluida != b.Excluida:
		if a.Excluida {
			return 1
		}
		return -1
	case a.Excluida:
		return strings.Compare(a.Entrada.Ref, b.Entrada.Ref)
	case a.Aptidao != b.Aptidao:
		return cmpDesc(a.Aptidao, b.Aptidao)
	case a.naturezaExplicita != b.naturezaExplicita:
		if a.naturezaExplicita {
			return -1
		}
		return 1
	case a.inferido() != b.inferido():
		if b.inferido() {
			return -1
		}
		return 1
	case a.Entrada.Rating.TotalScore != b.Entrada.Rating.TotalScore:
		return cmpDesc(a.Entrada.Rating.TotalScore, b.Entrada.Rating.TotalScore)
	}
	return strings.Compare(a.Entrada.Ref, b.Entrada.Ref)
}

func cmpDesc(a, b float64) int {
	if a > b {
		return -1
	}
	return 1
}

// declara diz se a lista tem o termo, sem diferenciar maiúsculas.
func declara(lista []string, termo string) bool {
	return slices.ContainsFunc(lista, func(s string) bool { return strings.EqualFold(strings.TrimSpace(s), strings.TrimSpace(termo)) })
}

// casados devolve, na ordem do pedido, os termos pedidos que a peça declara.
func casados(pedidos, declarados []string) []string {
	var out []string
	for _, p := range pedidos {
		if declara(declarados, p) {
			out = append(out, p)
		}
	}
	return out
}

func listar(termos []string) string {
	if len(termos) == 0 {
		return "nada declarado"
	}
	return strings.Join(termos, ", ")
}

// ExcluidaPorUso diz se a peça saiu por proibição de uso (quando_nao_usar) ou por servir a outra natureza.
// É a exclusão que vale avisar ao usuário; sair por etapa é só não ser candidata.
func (r Ranked) ExcluidaPorUso() bool {
	return r.Excluida && (strings.HasPrefix(r.Motivo, "quando_nao_usar:") || strings.HasPrefix(r.Motivo, "natureza:"))
}
