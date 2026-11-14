package triagem

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ErrResolvedorAusente indica ambiguidade sem ResolvedorLLM configurado.
var ErrResolvedorAusente = errors.New("resolvedor LLM ausente para entrada ambígua")

// Entrada representa a requisição recebida pelos pontos de entrada.
type Entrada struct {
	Texto     string `json:"texto"`
	Workspace string `json:"workspace,omitempty"`
}

// Decisao representa o resultado do motor híbrido de triagem.
type Decisao struct {
	Processo           string  `json:"processo"`
	EtapaInicial       string  `json:"etapa_inicial"`
	Confianca          float64 `json:"confianca"`
	Justificativa      string  `json:"justificativa"`
	PrecisaConfirmacao bool    `json:"precisa_confirmacao"`
	Camada             string  `json:"camada"`
}

// ResolvedorLLM abstrai a camada probabilística para desempate ou linguagem natural.
type ResolvedorLLM interface {
	Resolver(ctx context.Context, entrada string, catalogo []TriagemSpec) (Decisao, error)
}

// Decidir executa o motor híbrido: determinístico primeiro, probabilístico apenas na ambiguidade.
func Decidir(ctx context.Context, entrada Entrada, catalogo []TriagemSpec, resolvedor ResolvedorLLM) (Decisao, error) {
	if dec, ok := tentarDeterministico(entrada, catalogo); ok {
		return dec, nil
	}
	if resolvedor == nil {
		return Decisao{}, ErrResolvedorAusente
	}
	dec, err := resolvedor.Resolver(ctx, entrada.Texto, catalogo)
	if err != nil {
		return Decisao{}, err
	}
	dec.Camada = "probabilistica"
	return dec, nil
}

func tentarDeterministico(entrada Entrada, catalogo []TriagemSpec) (Decisao, bool) {
	texto := strings.ToLower(strings.TrimSpace(entrada.Texto))
	var matches []Decisao
	for _, spec := range catalogo {
		if dec, ok := matchEtapa(texto, spec); ok {
			matches = append(matches, dec)
			continue
		}
		if motivo, ok := casaSpec(texto, entrada.Workspace, spec); ok {
			matches = append(matches, Decisao{
				Processo:      spec.Processo,
				EtapaInicial:  primeiraEtapa(spec),
				Confianca:     1.0,
				Justificativa: motivo,
				Camada:        "deterministica",
			})
		}
	}
	if len(matches) == 1 {
		return matches[0], true
	}
	return Decisao{}, false
}

func matchEtapa(texto string, spec TriagemSpec) (Decisao, bool) {
	if texto == "" {
		return Decisao{}, false
	}
	for _, et := range spec.Etapas {
		if casaRegraEtapa(texto, et) {
			return Decisao{
				Processo:      spec.Processo,
				EtapaInicial:  et.ID,
				Confianca:     1.0,
				Justificativa: "match determinístico por etapa",
				Camada:        "deterministica",
			}, true
		}
	}
	return Decisao{}, false
}

func casaRegraEtapa(texto string, et RegraEtapa) bool {
	if et.SlashCommand != "" && strings.Contains(texto, strings.ToLower(et.SlashCommand)) {
		return true
	}
	if et.FlagCLI != "" && strings.Contains(texto, strings.ToLower(et.FlagCLI)) {
		return true
	}
	for _, kw := range et.PalavrasChave {
		if kw != "" && strings.Contains(texto, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

func casaSpec(texto, workspace string, spec TriagemSpec) (string, bool) {
	for _, regra := range spec.RegrasDeterministicas {
		if casaTexto(texto, regra) {
			return "match determinístico por comando/slug", true
		}
		if casaWorkspace(workspace, regra.ArtefatosAtivos) {
			return "match determinístico por artefato ativo no workspace", true
		}
	}
	return "", false
}

func casaTexto(texto string, regra RegraDeterministica) bool {
	if texto == "" {
		return false
	}
	for _, cmd := range regra.ComandosCLI {
		if cmd != "" && strings.Contains(texto, strings.ToLower(cmd)) {
			return true
		}
	}
	slug := strings.ToLower(strings.TrimSpace(regra.PadraoSlug))
	return slug != "" && (texto == slug || strings.HasPrefix(texto, slug+" "))
}

func casaWorkspace(workspace string, artefatos []string) bool {
	if strings.TrimSpace(workspace) == "" {
		return false
	}
	for _, art := range artefatos {
		if art == "" {
			continue
		}
		if _, err := os.Stat(filepath.Join(workspace, art)); err == nil {
			return true
		}
	}
	return false
}

func primeiraEtapa(spec TriagemSpec) string {
	if len(spec.CriterioAmbiguidade.Etapas) > 0 {
		return spec.CriterioAmbiguidade.Etapas[0]
	}
	return ""
}
