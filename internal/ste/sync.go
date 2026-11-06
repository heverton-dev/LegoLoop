package ste

import (
	"fmt"
	"os"
	"path/filepath"
)

// gravarRegraAtomica grava o conteúdo da regra de forma atômica no destino.
func gravarRegraAtomica(caminho string, conteudo string) error {
	dir := filepath.Dir(caminho)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	tmp := fmt.Sprintf("%s.tmp.%d", caminho, os.Getpid())
	if err := os.WriteFile(tmp, []byte(conteudo), 0644); err != nil {
		return err
	}
	return os.Rename(tmp, caminho)
}

// aplicarModoOff remove exclusivamente o arquivo stept-rules.md se existir.
func aplicarModoOff(caminho string) error {
	err := os.Remove(caminho)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// sincronizarHarness trata a ativação ou remoção das regras para um harness.
func sincronizarHarness(info HarnessInfo, perfil ConfigPerfil, createDirs bool) (HarnessStatus, error) {
	status := HarnessStatus{
		ID:        info.ID,
		Nome:      info.Nome,
		Instalado: info.Instalado,
		Caminho:   info.CaminhoRegras,
		Perfil:    perfil,
	}

	if !info.Instalado && !createDirs {
		return status, nil
	}

	if perfil.Modo == ModoOff {
		if err := aplicarModoOff(info.CaminhoRegras); err != nil {
			return status, err
		}
		status.Ativo = false
		return status, nil
	}

	prompt := PromptRegra(perfil)
	if err := gravarRegraAtomica(info.CaminhoRegras, prompt); err != nil {
		return status, err
	}
	status.Ativo = true
	status.Instalado = true
	return status, nil
}

// Sync materializa ou remove as regras nos harnesses com base na configuração.
func Sync(cfg *Config, homeDir string, createDirs bool) ([]HarnessStatus, error) {
	if cfg == nil {
		return nil, ErrInvalido
	}
	harnesses := ScanHarnesses(homeDir)
	res := make([]HarnessStatus, 0, len(harnesses))

	for _, h := range harnesses {
		perfil, origem := cfg.ResolverPerfil(h.ID)
		st, err := sincronizarHarness(h, perfil, createDirs)
		if err != nil {
			return nil, err
		}
		st.Origem = origem
		res = append(res, st)
	}
	return res, nil
}
