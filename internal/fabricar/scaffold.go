package fabricar

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"legoloop/internal/paths"
	"legoloop/internal/processo"
)

// extensoesPorSO dá as extensões de peça executável compatíveis com cada SO.
var extensoesPorSO = map[string][]string{
	"windows": {".ps1", ".bat", ".cmd"},
	"linux":   {".sh", ".bash"},
	"darwin":  {".sh", ".bash"},
}

// ResolveParaSO escolhe a variante da peça compatível com o SO alvo (AC-E05.3).
// A ordem dos candidatos dá a preferência entre as compatíveis.
// Sem variante compatível, devolve vazio: peça do SO errado nunca executa.
func ResolveParaSO(candidatos []string, so string) string {
	exts := extensoesPorSO[strings.ToLower(so)]
	for _, c := range candidatos {
		if slices.Contains(exts, strings.ToLower(filepath.Ext(c))) {
			return c
		}
	}
	return ""
}

// ResolvePorSO resolve a variante da peça pelo SO detectado do ambiente.
func ResolvePorSO(candidatos []string) string { return ResolveParaSO(candidatos, runtime.GOOS) }

// vacinaAnticorpo é uma entrada de anticorpo: assinatura do erro, o erro, o antídoto e as reincidências.
type vacinaAnticorpo struct {
	Assinatura   string `json:"assinatura"`
	Erro         string `json:"erro"`
	Antidoto     string `json:"antidoto"`
	Recorrencias int    `json:"recorrencias"`
}

// anticorposV1 é o manifesto anticorpos.v1 do Processo: vacinas globais (infra/LLM) e de domínio.
type anticorposV1 struct {
	SchemaVersion string            `json:"schema_version"`
	Processo      string            `json:"processo"`
	Globais       []vacinaAnticorpo `json:"globais"`
	Dominio       []vacinaAnticorpo `json:"dominio"`
}

// AnticorposPadrao renderiza o anticorpos.json do nascimento 5/5: sem vacinas registradas,
// pronto para receber as de infra/LLM (globais) e as do domínio do Processo.
func AnticorposPadrao(nome string) ([]byte, error) {
	return json.MarshalIndent(anticorposV1{
		SchemaVersion: "anticorpos.v1",
		Processo:      nome,
		Globais:       []vacinaAnticorpo{},
		Dominio:       []vacinaAnticorpo{},
	}, "", "  ")
}

// anticorposDe renderiza o anticorpos.json ligado ao Processo recém-projetado.
func anticorposDe(p processo.Processo) ([]byte, error) {
	return AnticorposPadrao(p.Nome)
}

// DocsDoSlice espelha o Sexteto da raiz docs/ dentro da fatia vertical
// <pasta-do-Processo>/docs/<sub>/ (AC-E05.1). A raiz docs/ continua sendo escrita:
// o espelhamento é dupla-escrita até a virada de escopo E06+. Kind é livre;
// "slice_doc" não vira artefato canônico (artifactsOf ignora kinds desconhecidos).
func DocsDoSlice(roots Roots, p processo.Processo, pr *processo.Projecoes, sexteto [3]namedFile) []PlannedFile {
	base := filepath.Join(roots.Vault, paths.DirAntigo, p.Nome, "docs")
	entradas := []struct {
		sub, nome string
		data      []byte
	}{
		{"interviews", sexteto[0].Name, sexteto[0].Data},
		{"intents", sexteto[1].Name, sexteto[1].Data},
		{"specs", pr.SpecNome, []byte(pr.Spec)},
		{"plans", pr.PlanNome, []byte(pr.Plan)},
		{"routes", pr.RouteNome, pr.Route},
		{"reports", sexteto[2].Name, sexteto[2].Data},
		{"manuals", pr.ManualNome, []byte(pr.Manual)},
	}
	var out []PlannedFile
	for _, e := range entradas {
		if len(bytes.TrimSpace(e.data)) == 0 {
			continue
		}
		out = append(out, PlannedFile{Kind: "slice_doc", Path: filepath.Join(base, e.sub, e.nome), Data: e.data})
	}
	return out
}
