package refinar

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"legoloop/internal/caminhos"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

// ROUTE e DAGs em schema v2 (nós TICK-NN.M) são curados à mão: roteamento, stop-gate, timeouts e assinatura do ROUTE.
// O refinar nunca os troca pelo formato v1 gerado. Ele edita só os campos que dependem das etapas (nome, peça,
// artefato e gate), no texto bruto, para não reordenar nem reformatar o resto do arquivo.

func init() {
	ste.Default.Register("refinar.v2_mantido",
		"O arquivo %s é v2 curado e tem outro número de nós que as etapas. O programa não trocou esse arquivo.")
}

var (
	reTick    = regexp.MustCompile(`^TICK-\d+\.(\d+)$`)
	reSHARota = regexp.MustCompile(`("route"\s*:\s*\{[^}]*?"sha256"\s*:\s*")[0-9a-f]{64}(")`)
	reLabel   = regexp.MustCompile(`("label"\s*:\s*)"(?:[^"\\]|\\.)*"`)
	reNomeNo  = regexp.MustCompile(`("name"\s*:\s*)"(?:[^"\\]|\\.)*"`)
	reRef     = regexp.MustCompile(`("ref"\s*:\s*)"(?:[^"\\]|\\.)*"`)
)

type noV2 struct {
	ID     string `json:"id"`
	Output struct {
		Artifact  string   `json:"artifact"`
		Artifacts []string `json:"artifacts"`
		Gate      string   `json:"gate"`
	} `json:"output"`
}

// ehV2Curado diz se o JSON é um ROUTE ou DAG em schema v2 com nós TICK.
func ehV2Curado(data []byte) bool {
	var d struct {
		SchemaVersion int `json:"schema_version"`
		Nodes         []struct {
			ID string `json:"id"`
		} `json:"nodes"`
	}
	if json.Unmarshal(data, &d) != nil || d.SchemaVersion < 2 {
		return false
	}
	for _, n := range d.Nodes {
		if reTick.MatchString(n.ID) {
			return true
		}
	}
	return false
}

// adaptarV2 aplica as etapas aos nós TICK do arquivo. Devolve false (e os bytes originais) quando o número de
// nós TICK difere do número de etapas: nesse caso o arquivo fica como está.
func adaptarV2(data []byte, etapas []processo.Etapa) ([]byte, bool) {
	var top struct {
		Nodes []json.RawMessage `json:"nodes"`
	}
	if json.Unmarshal(data, &top) != nil {
		return data, false
	}
	type tick struct {
		raw []byte
		nv  noV2
		n   int
	}
	var ticks []tick
	for _, raw := range top.Nodes {
		var nv noV2
		if json.Unmarshal(raw, &nv) != nil {
			return data, false
		}
		if m := reTick.FindStringSubmatch(nv.ID); m != nil {
			n, _ := strconv.Atoi(m[1])
			ticks = append(ticks, tick{raw: raw, nv: nv, n: n})
		}
	}
	if len(ticks) != len(etapas) {
		return data, false
	}
	out := data
	for _, t := range ticks {
		if t.n < 1 || t.n > len(etapas) {
			return data, false
		}
		out = bytes.Replace(out, t.raw, editarNo(t.raw, t.nv, etapas[t.n-1]), 1)
	}
	return out, true
}

func jsonTexto(s string) string {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)
	return strings.TrimSuffix(strings.TrimSpace(b.String()), "\n")
}

func editarNo(raw []byte, nv noV2, e processo.Etapa) []byte {
	s := string(raw)
	antigo := nv.Output.Artifact
	if antigo == "" && len(nv.Output.Artifacts) > 0 {
		antigo = nv.Output.Artifacts[0]
	}
	if nv.Output.Gate != "" {
		s = strings.ReplaceAll(s, jsonTexto(nv.Output.Gate), jsonTexto(e.Gate))
	}
	if antigo != "" {
		s = strings.ReplaceAll(s, antigo, e.Produz.Caminho)
	}
	if reLabel.MatchString(s) {
		s = trocarPrimeiro(s, reLabel, fmt.Sprintf("`%s` — %s — peça `%s`", e.ID, e.Titulo, e.Peca))
	} else {
		s = trocarPrimeiro(s, reNomeNo, fmt.Sprintf("%s — peça %s", e.Titulo, e.Peca))
		if e.Peca != "" {
			s = trocarPrimeiro(s, reRef, "skills/"+e.Peca)
		}
	}
	return []byte(s)
}

// trocarPrimeiro troca o valor da primeira chave que casa com re, mantendo o prefixo "chave":.
func trocarPrimeiro(s string, re *regexp.Regexp, valor string) string {
	loc := re.FindStringSubmatchIndex(s)
	if loc == nil {
		return s
	}
	return s[:loc[3]] + jsonTexto(valor) + s[loc[1]:]
}

// comSHARota grava o sha do ROUTE na assinatura upstream do DAG.
func comSHARota(data []byte, rota []byte) []byte {
	soma := sha256.Sum256(rota)
	sha := hex.EncodeToString(soma[:])
	return reSHARota.ReplaceAll(data, []byte("${1}"+sha+"${2}"))
}

// curados troca, em todos, o ROUTE e os DAGs v2 curados pela versão adaptada. O DAG espelho (docs/dags) e o DAG da
// projeção padrão (processos/<nome>/pipeline.dag.json) acompanham, com o sha do ROUTE novo. Devolve os avisos de arquivos
// que o programa deixou como estavam. O ROUTE só é curado se difere da projeção da versão anterior do processo: o igual
// foi gerado pelo programa e pode ser regerado.
func curados(todos []processo.File, raiz, docs string, p processo.Processo, pr *processo.Projecoes, rotaGerada []byte) ([]processo.File, []string) {
	var avisos []string
	rotaPath := filepath.Join(docs, "routes", pr.RouteNome)
	var rota []byte
	for i := range todos {
		if todos[i].Kind != "ROUTE" {
			continue
		}
		atual, err := os.ReadFile(rotaPath)
		if err == nil && !bytes.Equal(atual, rotaGerada) && ehV2Curado(atual) {
			if novo, ok := adaptarV2(atual, p.Etapas); ok {
				todos[i].Data = novo
			} else {
				todos[i].Data = atual
				avisos = append(avisos, ste.Msg("refinar.v2_mantido", rotaPath))
			}
		}
		rota = todos[i].Data
	}
	for i := range todos {
		if todos[i].Kind != "DAG_ISO" {
			continue
		}
		if atual, err := os.ReadFile(todos[i].Path); err == nil && ehV2Curado(atual) {
			todos[i].Data, avisos = adaptarDag(atual, p, rota, todos[i].Path, avisos)
		}
	}
	loopsDag := caminhos.CaminhoDAG(filepath.Join(caminhos.DirProcessosEm(raiz), p.Nome))
	if atual, err := os.ReadFile(loopsDag); err == nil && ehV2Curado(atual) {
		novo, avisos2 := adaptarDag(atual, p, rota, loopsDag, avisos)
		todos = append(todos, processo.File{Kind: "DAG_LOOP", Path: loopsDag, Data: novo})
		avisos = avisos2
	}
	return todos, avisos
}

func adaptarDag(atual []byte, p processo.Processo, rota []byte, caminho string, avisos []string) ([]byte, []string) {
	novo, ok := adaptarV2(atual, p.Etapas)
	if !ok {
		return atual, append(avisos, ste.Msg("refinar.v2_mantido", caminho))
	}
	if rota != nil {
		novo = comSHARota(novo, rota)
	}
	return novo, avisos
}
