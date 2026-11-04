package fabricar

import (
	"encoding/json"
	"os"
	"path/filepath"
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
func buildReport(d Deps, clar Clarification, scan ScanResult, m Merged, tri triadFiles, files []PlannedFile, plano Plano) *Report {
	rep := &Report{}
	rep.Clarificacao.Espelho = clar.Mirror
	rep.Clarificacao.Confianca = clar.Confidence
	rep.Clarificacao.Turnos = clar.Turns
	rep.Clarificacao.Alcancou = clar.Reached

	populateVarredura(d, scan, m, rep)
	populateFusao(clar, m, rep)
	extractTriadeInfo(tri, rep)
	populateGravacao(files, plano, rep)

	rep.Cartao = plano.Cartao()
	return rep
}

func populateVarredura(d Deps, scan ScanResult, m Merged, rep *Report) {
	rep.Varredura.PecasVistas = make(map[string]int)
	for _, p := range scan.Pieces {
		rep.Varredura.PecasVistas[string(p.Kind)]++
	}
	if len(rep.Varredura.PecasVistas) == 0 {
		rep.Varredura.PecasVistas["skill"] = 0
	}
	for _, s := range scan.Skipped {
		rep.Varredura.Cortes = append(rep.Varredura.Cortes, s.Root+": "+s.Reason)
	}
	for _, ex := range m.Excluded {
		rep.Varredura.Cortes = append(rep.Varredura.Cortes, ex.Name)
	}
	if d.Roots.Vault != "" {
		matrizPath := filepath.Join(d.Roots.Vault, "registry_matrix.json")
		if fi, err := os.Stat(matrizPath); err == nil {
			rep.Varredura.MatrizIdade = d.Clock().Sub(fi.ModTime())
		}
	}
}

func populateFusao(clar Clarification, m Merged, rep *Report) {
	rep.Fusao.Necessidades = matriz.Needs{Natureza: clar.Natureza}
	rep.Fusao.Decisoes = m.Decisoes
	rep.Fusao.Lacunas = m.Lacunas
	rep.Fusao.Divergencias = m.Divergencias
	if len(rep.Fusao.Decisoes) == 0 && len(m.Steps) > 0 {
		for _, s := range m.Steps {
			rep.Fusao.Decisoes = append(rep.Fusao.Decisoes, Decisao{
				Peca:   s.Piece,
				Entrou: true,
				Motivo: "selecionada na fusão",
			})
		}
	}
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
