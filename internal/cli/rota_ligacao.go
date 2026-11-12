package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"legoloop/internal/orchestrix"
	"legoloop/internal/roteamento"
	"legoloop/internal/roteamento/homologacao"
)

// ligarRoteamento liga ao motor o estado de homologação e os candidatos elegíveis do perfil.
// Sem arquivo de homologação, nenhum harness fica homologado e o motor recusa outro harness.
func ligarRoteamento(ox *orchestrix.Orchestrix, d Deps) error {
	home := d.Home
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	store := rotaLerHomologacao(filepath.Join(home, ".aidd", "homologacao.json"))
	global, _ := os.ReadFile(filepath.Join(home, ".aidd", "roteamento.json"))
	projeto, err := os.ReadFile(filepath.Join(".aidd", "roteamento.json"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	elegiveis, err := roteamento.CarregarPerfil(global, projeto, store)
	if err != nil {
		return err
	}
	ox.Elegiveis = elegiveis
	ox.EstadoHarness = func(harness, via string) string {
		return store[roteamento.Par{Harness: harness, Via: roteamento.Via(via)}.Chave()].Estado
	}
	ox.Reprovados = &orchestrix.RegistroReprovados{}
	ox.AoReprovar = func(harness, via, motivo string) { persistirReprovado(home, harness, via, motivo) }
	ox.PesoHarnessMB = func(harness, via string) int {
		return store[roteamento.Par{Harness: harness, Via: roteamento.Via(via)}.Chave()].PesoMB
	}
	ox.MemoriaLivre = orchestrix.MemoriaLivreSO
	if d.Stderr != nil {
		ox.Aviso = func(m string) { fmt.Fprintln(d.Stderr, m) }
	}
	return nil
}

// persistirReprovado grava o par como reprovado em homologacao.json, para as próximas execuções.
func persistirReprovado(home, harness, via, motivo string) {
	caminho := homologacao.Caminho(home)
	store, err := homologacao.Ler(caminho)
	if err != nil {
		return
	}
	store[homologacao.Chave(harness, via)] = homologacao.Entrada{Estado: homologacao.EstadoReprovado,
		Data: time.Now().Format("2006-01-02"), Criterios: map[string]bool{}, Motivo: motivo}
	_ = store.Gravar(caminho)
}
