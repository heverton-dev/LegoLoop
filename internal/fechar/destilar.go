package fechar

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ErrStopGateNaoAprovado indica tentativa de limpeza antes do Stop-Gate final.
var ErrStopGateNaoAprovado = errors.New("fechar: stop-gate final nao aprovado para destilacao")

// ErroEpisodico representa um registro de falha capturado na memoria episodica.
type ErroEpisodico struct {
	Assinatura string `json:"assinatura"`
	Mensagem   string `json:"mensagem"`
}

// GateDestilado representa um gate mecanico derivado de erro recorrente.
type GateDestilado struct {
	ID             string    `json:"id"`
	ErroAssinatura string    `json:"erro_assinatura"`
	Ocorrencias    int       `json:"ocorrencias"`
	RegraGate      string    `json:"regra_gate"`
	CriadoEm       time.Time `json:"criado_em"`
}

// DestilacaoResult resume os gates destilados e a limpeza da memoria episodica.
type DestilacaoResult struct {
	GatesDestilados []GateDestilado `json:"gates_destilados"`
	MemoriaLimpa    bool            `json:"memoria_limpa"`
	ArquivoGates    string          `json:"arquivo_gates,omitempty"`
}

// DestilarELimpar destila erros recorrentes em gates mecanicos e limpa a memoria episodica apos o Stop-Gate.
func DestilarELimpar(repoDir, slug string, stopGateAprovado bool, limiar int, dryRun bool) (*DestilacaoResult, error) {
	if !stopGateAprovado {
		return nil, ErrStopGateNaoAprovado
	}
	if limiar < 2 {
		limiar = 2
	}
	epPath := filepath.Join(repoDir, ".legoloop", "episodic", slug+".json")
	eventos, err := lerMemoriaEpisodica(epPath)
	if err != nil {
		return nil, err
	}
	gates := destilarRecorrentes(eventos, limiar)
	res := &DestilacaoResult{GatesDestilados: gates}
	if dryRun {
		return res, nil
	}
	return persistirELimpar(repoDir, slug, epPath, res)
}

func lerMemoriaEpisodica(epPath string) ([]ErroEpisodico, error) {
	data, err := os.ReadFile(epPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var eventos []ErroEpisodico
	if err := json.Unmarshal(data, &eventos); err != nil {
		return nil, err
	}
	return eventos, nil
}

func destilarRecorrentes(eventos []ErroEpisodico, limiar int) []GateDestilado {
	contagem := make(map[string]int)
	for _, ev := range eventos {
		if ev.Assinatura != "" {
			contagem[ev.Assinatura]++
		}
	}
	assinaturas := make([]string, 0, len(contagem))
	for sig, n := range contagem {
		if n >= limiar {
			assinaturas = append(assinaturas, sig)
		}
	}
	sort.Strings(assinaturas)
	gates := make([]GateDestilado, 0, len(assinaturas))
	now := time.Now().UTC()
	for i, sig := range assinaturas {
		gates = append(gates, GateDestilado{
			ID:             fmt.Sprintf("GATE-%03d", i+1),
			ErroAssinatura: sig,
			Ocorrencias:    contagem[sig],
			RegraGate:      "bloquear_assinatura:" + sig,
			CriadoEm:       now,
		})
	}
	return gates
}

func persistirELimpar(repoDir, slug, epPath string, res *DestilacaoResult) (*DestilacaoResult, error) {
	if len(res.GatesDestilados) > 0 {
		gateDir := filepath.Join(repoDir, ".legoloop", "gates")
		if err := os.MkdirAll(gateDir, 0o755); err != nil {
			return nil, err
		}
		gatePath := filepath.Join(gateDir, slug+"_destilados.json")
		raw, err := json.MarshalIndent(res.GatesDestilados, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(gatePath, raw, 0o644); err != nil {
			return nil, err
		}
		res.ArquivoGates = gatePath
	}
	if err := os.Remove(epPath); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	res.MemoriaLimpa = true
	return res, nil
}
