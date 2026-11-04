package processo

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
)

var (
	reLabelRoute = regexp.MustCompile("^`([^`]+)` — ")
	rePlan       = regexp.MustCompile("(?m)^#### Ticket `TICK-[^`]+` — `([^`]+)` — ")
	reSkill      = regexp.MustCompile("(?m)^\\d+\\. Etapa `([^`]+)`: ")
	reSpec       = regexp.MustCompile("(?m)^\\* \\[ \\] `AC-[^`]+`: \\*\\*Etapa `([^`]+)`")
)

// Verificar prova que as seis projeções listam as mesmas etapas, na mesma ordem das etapas do processo.
// Devolve ErrIncoerente dizendo qual projeção divergiu.
func (pr *Projecoes) Verificar(p Processo) error {
	ord, err := p.Ordem()
	if err != nil {
		return err
	}
	quero := make([]string, len(ord))
	for i, e := range ord {
		quero[i] = e.ID
	}
	vistas, err := pr.ids()
	if err != nil {
		return err
	}
	for _, nome := range []string{"dag", "loop", "route", "plan", "skill", "spec"} {
		if !slices.Equal(vistas[nome], quero) {
			return fmt.Errorf("%w: %s lista %v, mas o processo tem %v", ErrIncoerente, nome, vistas[nome], quero)
		}
	}
	return nil
}

// ids lê, de cada arquivo, os ids das etapas na ordem em que aparecem.
func (pr *Projecoes) ids() (map[string][]string, error) {
	dagIDs, err := idsDeNos(pr.DAG, "dag")
	if err != nil {
		return nil, err
	}
	loopIDs, err := idsDeNos(pr.Loop, "loop")
	if err != nil {
		return nil, err
	}
	routeIDs, err := idsDaRoute(pr.Route)
	if err != nil {
		return nil, err
	}
	return map[string][]string{
		"dag": dagIDs, "loop": loopIDs, "route": routeIDs,
		"plan": primeiros(rePlan, pr.Plan), "skill": primeiros(reSkill, pr.Skill), "spec": primeiros(reSpec, pr.Spec),
	}, nil
}

func idsDeNos(data []byte, nome string) ([]string, error) {
	var d struct {
		Nodes []struct {
			ID string `json:"id"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("%w: %s ilegível: %v", ErrIncoerente, nome, err)
	}
	out := make([]string, len(d.Nodes))
	for i, n := range d.Nodes {
		out[i] = n.ID
	}
	return out, nil
}

// idsDaRoute lê os tickets (nós unit_task); o id da etapa é o primeiro trecho de código do rótulo.
func idsDaRoute(data []byte) ([]string, error) {
	var d struct {
		Nodes []struct{ Type, Label string } `json:"nodes"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("%w: route ilegível: %v", ErrIncoerente, err)
	}
	var out []string
	for _, n := range d.Nodes {
		if n.Type != "unit_task" {
			continue
		}
		if m := reLabelRoute.FindStringSubmatch(n.Label); m != nil {
			out = append(out, m[1])
		} else {
			out = append(out, "?"+n.Label)
		}
	}
	return out, nil
}

func primeiros(re *regexp.Regexp, texto string) []string {
	var out []string
	for _, m := range re.FindAllStringSubmatch(texto, -1) {
		out = append(out, m[1])
	}
	return out
}
