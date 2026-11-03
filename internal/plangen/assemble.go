package plangen

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var idNumsRe = regexp.MustCompile(`(\d+)\.(\d+)`)

// canonID normaliza TICKET-E01.1, TICK-E01.1 e TICK-01.1 para TICK-01.1 (formato do template).
func canonID(raw string) string {
	m := idNumsRe.FindStringSubmatch(raw)
	if m == nil {
		return strings.TrimSpace(raw)
	}
	var e, s int
	fmt.Sscanf(m[1], "%d", &e)
	fmt.Sscanf(m[2], "%d", &s)
	return fmt.Sprintf("TICK-%02d.%d", e, s)
}

// canonEpic devolve E01, E02... a partir do campo Épico ou, na falta dele, do número do ticket.
func canonEpic(raw, ticketID string) string {
	if n := epicNumber(raw); n > 0 {
		return fmt.Sprintf("E%02d", n)
	}
	if n := epicNumber(ticketID); n > 0 {
		return fmt.Sprintf("E%02d", n)
	}
	return "E01"
}

func epicGateID(epic string) string { return fmt.Sprintf("SG-E%02d", epicNumber(epic)) }

func ticketGateID(ticketID string) string { return "SG-" + strings.TrimPrefix(ticketID, "TICK-") }

// epicOrder lista os épicos na ordem em que aparecem nos tickets.
func epicOrder(tickets []Ticket) []string {
	var order []string
	for _, t := range tickets {
		if !slices.Contains(order, t.EpicID) {
			order = append(order, t.EpicID)
		}
	}
	return order
}

func epicTicketIDs(tickets []Ticket, epic string) []string {
	var ids []string
	for _, t := range tickets {
		if t.EpicID == epic {
			ids = append(ids, t.ID)
		}
	}
	return ids
}

// buildNodes monta a rota: tickets, Stop-Gates humanos do nível de supervisão e o SG-FINAL.
//   - nano: gate após cada ticket; micro: após cada ticket micro; epic e autonomous: nenhum por ticket;
//   - entre épicos (nunca depois do último, que fecha no SG-FINAL): gate de épico, salvo se o último
//     ticket do épico já tem o seu.
//
// Quem dependia de um ticket com gate passa a depender do gate, e cada épico espera a fronteira do anterior.
func buildNodes(tickets []Ticket, m planMeta) ([]routeNode, error) {
	order := epicOrder(tickets)
	gateOf, boundary := map[string]string{}, map[string]string{}
	var out []routeNode
	for i, t := range tickets {
		out = append(out, ticketNode(t, m))
		if t.HumanGateAfter {
			g := humanGate(ticketGateID(t.ID), fmt.Sprintf("Evidências do ticket %s aprovadas? Autoriza prosseguir?", t.ID), []string{t.ID}, t.EpicID)
			out, gateOf[t.ID], boundary[t.EpicID] = append(out, g), g.ID, g.ID
		}
		if m.Level != SupervisionAutonomous && lastOfEpic(tickets, i) && t.EpicID != order[len(order)-1] && !t.HumanGateAfter {
			g := humanGate(epicGateID(t.EpicID), fmt.Sprintf("Evidências do Épico %02d aprovadas? Autoriza iniciar o próximo épico?", epicNumber(t.EpicID)), epicTicketIDs(tickets, t.EpicID), t.EpicID)
			out, boundary[t.EpicID] = append(out, g), g.ID
		}
	}
	rewire(out, order, gateOf, boundary)
	out = append(out, finalGate(terminals(out)))
	return out, checkAcyclic(out)
}

func lastOfEpic(tickets []Ticket, i int) bool {
	return i == len(tickets)-1 || tickets[i+1].EpicID != tickets[i].EpicID
}

// rewire faz os dependentes esperarem o gate e cada épico esperar a fronteira do épico anterior.
func rewire(nodes []routeNode, order []string, gateOf, boundary map[string]string) {
	for i := range nodes {
		n := &nodes[i]
		if n.Type != "unit_task" {
			continue
		}
		for j, d := range n.DependsOn {
			if g, ok := gateOf[d]; ok {
				n.DependsOn[j] = g
			}
		}
		if pos := slices.Index(order, n.epic); pos > 0 {
			if b, ok := boundary[order[pos-1]]; ok && !slices.Contains(n.DependsOn, b) {
				n.DependsOn = append(n.DependsOn, b)
			}
		}
	}
}

// terminals são os nós de que ninguém depende: o SG-FINAL espera por eles.
func terminals(nodes []routeNode) []string {
	used := map[string]bool{}
	for _, n := range nodes {
		for _, d := range n.DependsOn {
			used[d] = true
		}
	}
	var out []string
	for _, n := range nodes {
		if !used[n.ID] {
			out = append(out, n.ID)
		}
	}
	return out
}

// checkAcyclic recusa rotas com ciclo ou dependência inexistente (Kahn).
func checkAcyclic(nodes []routeNode) error {
	indeg, users := map[string]int{}, map[string][]string{}
	for _, n := range nodes {
		indeg[n.ID] += 0
	}
	for _, n := range nodes {
		for _, d := range n.DependsOn {
			if _, ok := indeg[d]; !ok {
				return fmt.Errorf("o nó %s depende de %s, que não existe", n.ID, d)
			}
			indeg[n.ID]++
			users[d] = append(users[d], n.ID)
		}
	}
	var queue []string
	for _, n := range nodes {
		if indeg[n.ID] == 0 {
			queue = append(queue, n.ID)
		}
	}
	seen := 0
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		seen++
		for _, u := range users[id] {
			if indeg[u]--; indeg[u] == 0 {
				queue = append(queue, u)
			}
		}
	}
	if seen != len(nodes) {
		return errors.New("a rota gerada tem dependência circular")
	}
	return nil
}

// resolveDeps mantém só dependências que existem e avisa sobre as que foram descartadas.
func resolveDeps(tickets []Ticket) []string {
	known := map[string]bool{}
	for _, t := range tickets {
		known[t.ID] = true
	}
	var warns []string
	for i := range tickets {
		var keep []string
		for _, d := range tickets[i].DependsOn {
			if known[d] && d != tickets[i].ID {
				keep = append(keep, d)
				continue
			}
			warns = append(warns, fmt.Sprintf("o ticket %s declara a dependência %q, que não existe; ela foi descartada", tickets[i].ID, d))
		}
		tickets[i].DependsOn = keep
	}
	return warns
}

// uncovered lista os critérios da SPEC que nenhum ticket cobre (plano incompleto).
func uncovered(tickets []Ticket, specCriteria []string) []string {
	covered := map[string]bool{}
	for _, t := range tickets {
		for _, c := range t.SpecCriteria {
			covered[c] = true
		}
	}
	var warns []string
	for _, c := range specCriteria {
		if !covered[c] {
			warns = append(warns, fmt.Sprintf("o critério %s da SPEC não é coberto por nenhum ticket", c))
		}
	}
	return warns
}
