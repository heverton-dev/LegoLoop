package fabricar

import (
	"encoding/json"
	"time"

	"legoloop/internal/interview"
	"legoloop/internal/matriz"
	"legoloop/internal/ste"
)

// Report é o único objeto de relatório. A CLI o imprime, o Chat o repassa e o Studio o desenha.
type Report struct {
	Clarificacao struct {
		Espelho   interview.Mirror `json:"espelho"`
		Confianca int              `json:"confianca"`
		Turnos    int              `json:"turnos"`
		Alcancou  bool             `json:"alcancou"`
	} `json:"clarificacao"`
	Varredura struct {
		PecasVistas map[string]int `json:"pecas_vistas"`
		Cortes      []string       `json:"cortes"`
		MatrizIdade time.Duration  `json:"matriz_idade"`
		Defasagem   matriz.Stale   `json:"defasagem"`
	} `json:"varredura"`
	Fusao struct {
		Necessidades matriz.Needs `json:"necessidades"`
		Decisoes     []Decisao    `json:"decisoes"`
		Lacunas      []string     `json:"lacunas"`
		Divergencias []string     `json:"divergencias"`
	} `json:"fusao"`
	Triade struct {
		ACs      int      `json:"acs"`
		Tickets  int      `json:"tickets"`
		Gates    []string `json:"gates"`
		Coerente bool     `json:"coerente"`
	} `json:"triade"`
	Gravacao struct {
		Arquivos  []string `json:"arquivos"`
		Harnesses []string `json:"harnesses"`
		Bytes     int64    `json:"bytes"`
	} `json:"gravacao"`
	Cartao ste.Cartao `json:"cartao"`
}

// buildReport constrói o relatório canônico após as 5 camadas da esteira de fabricação.
func buildReport(d Deps, clar Clarification, c consulta, tri triadFiles, files []PlannedFile, plano Plano) *Report {
	rep := &Report{}
	rep.Clarificacao.Espelho = clar.Mirror
	rep.Clarificacao.Confianca = clar.Confidence
	rep.Clarificacao.Turnos = clar.Turns
	rep.Clarificacao.Alcancou = clar.Reached

	populateVarredura(d, c, rep)
	populateFusao(c, rep)
	extractTriadeInfo(tri, rep)
	populateGravacao(files, plano, rep)

	rep.Cartao = plano.Cartao()
	return rep
}

// populateVarredura é a Camada 2: as peças da matriz por tipo, os cortes e a idade da matriz.
func populateVarredura(d Deps, c consulta, rep *Report) {
	rep.Varredura.PecasVistas = make(map[string]int)
	for _, e := range c.fonte.Entradas {
		rep.Varredura.PecasVistas[e.Tipo]++
	}
	for _, ex := range c.m.Excluded {
		rep.Varredura.Cortes = append(rep.Varredura.Cortes, ex.Name)
	}
	if !c.fonte.GeradaEm.IsZero() {
		rep.Varredura.MatrizIdade = d.Clock().Sub(c.fonte.GeradaEm)
	}
}

// populateFusao é a Camada 3, como a Fusão decidiu: nada é acrescentado aqui.
func populateFusao(c consulta, rep *Report) {
	rep.Fusao.Necessidades = c.needs
	rep.Fusao.Decisoes = c.m.Decisoes
	rep.Fusao.Lacunas = c.m.Lacunas
	rep.Fusao.Divergencias = c.m.Divergencias
}

func populateGravacao(files []PlannedFile, plano Plano, rep *Report) {
	var totalBytes int64
	for _, f := range files {
		rep.Gravacao.Arquivos = append(rep.Gravacao.Arquivos, f.Path)
		totalBytes += int64(len(f.Data))
	}
	rep.Gravacao.Harnesses = plano.Ligados
	rep.Gravacao.Bytes = totalBytes
}

type routeNodeJSON struct {
	Type   string `json:"type"`
	Output struct {
		Gate string `json:"gate"`
	} `json:"output"`
}

type routeDocJSON struct {
	Nodes []routeNodeJSON `json:"nodes"`
}

func extractTriadeInfo(tri triadFiles, rep *Report) {
	rep.Triade.Coerente = true
	var r routeDocJSON
	if err := json.Unmarshal(tri.Route.Data, &r); err == nil {
		rep.Triade.Tickets = len(r.Nodes)
		for _, n := range r.Nodes {
			if n.Output.Gate != "" {
				rep.Triade.Gates = append(rep.Triade.Gates, n.Output.Gate)
			}
		}
	}
	if len(rep.Triade.Gates) == 0 {
		rep.Triade.Gates = []string{"exit_code == 0"}
	}
	rep.Triade.ACs = rep.Triade.Tickets
	if rep.Triade.ACs == 0 {
		rep.Triade.ACs = 1
	}
}
