package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"legoloop/internal/orchestrix"
	"legoloop/internal/roteamento"
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
	ox.PesoHarnessMB = func(harness, via string) int {
		return store[roteamento.Par{Harness: harness, Via: roteamento.Via(via)}.Chave()].PesoMB
	}
	ox.MemoriaLivre = orchestrix.MemoriaLivreSO
	if d.Stderr != nil {
		ox.Aviso = func(m string) { fmt.Fprintln(d.Stderr, m) }
	}
	return nil
}
