package migrar

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

// Relatorio resume o saneamento e a geração de triagem.json em processos/.
type Relatorio struct {
	Validos         []string `json:"validos"`
	Removidos       []string `json:"removidos"`
	TriagensGeradas []string `json:"triagens_geradas"`
}

type dagMinimo struct {
	Nodes []struct {
		ID string `json:"id"`
	} `json:"nodes"`
}

// MigrarProcessos varre processosDir delegando para a rotina compatível de migração.
func MigrarProcessos(processosDir string, apply bool) (Relatorio, error) {
	return MigrarLoops(processosDir, apply)
}

// MigrarLoops varre processosDir, remove diretórios sem pipeline.dag.json válido e gera triagem.json.
func MigrarLoops(processosDir string, apply bool) (Relatorio, error) {
	var rel Relatorio
	entries, err := os.ReadDir(processosDir)
	if err != nil {
		return rel, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			if err := normalizarProcessoSolto(processosDir, entry.Name(), apply, &rel); err != nil {
				return rel, err
			}
			continue
		}
		if err := processarProcesso(processosDir, entry.Name(), apply, &rel); err != nil {
			return rel, err
		}
	}
	return rel, nil
}

func processarProcesso(processosDir, slug string, apply bool, rel *Relatorio) error {
	dir := filepath.Join(processosDir, slug)
	etapas, ok := lerEtapasDAG(caminhos.CaminhoDAG(dir))
	if !ok {
		rel.Removidos = append(rel.Removidos, slug)
		if apply {
			return os.RemoveAll(dir)
		}
		return nil
	}
	rel.Validos = append(rel.Validos, slug)
	if err := garantirMemoria(dir, slug, apply); err != nil {
		return err
	}
	triagemPath := filepath.Join(dir, "triagem.json")
	if triagemAtualizada(triagemPath) {
		return nil
	}
	rel.TriagensGeradas = append(rel.TriagensGeradas, slug)
	if !apply {
		return nil
	}
	return gravarTriagem(triagemPath, slug, etapas)
}

func lerEtapasDAG(path string) ([]string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var dag dagMinimo
	if err := json.Unmarshal(data, &dag); err != nil {
		return nil, false
	}
	var etapas []string
	for _, n := range dag.Nodes {
		if n.ID != "" {
			etapas = append(etapas, n.ID)
		}
	}
	return etapas, true
}

func gravarTriagem(path, slug string, etapas []string) error {
	spec := triagem.GerarTriagemPadrao(slug, etapas)
	out, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(out, '\n'), 0644)
}

func triagemAtualizada(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var spec triagem.TriagemSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return false
	}
	return len(spec.Etapas) > 0
}

func normalizarProcessoSolto(processosDir, nome string, apply bool, rel *Relatorio) error {
	const sufixo = ".loop.json"
	if !strings.HasSuffix(nome, sufixo) {
		return nil
	}
	slug := strings.TrimSuffix(nome, sufixo)
	dir := filepath.Join(processosDir, slug)
	if apply {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		data, err := os.ReadFile(filepath.Join(processosDir, nome))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, nome), data, 0o644); err != nil {
			return err
		}
		if err := garantirMemoria(dir, slug, true); err != nil {
			return err
		}
		triagemPath := filepath.Join(dir, "triagem.json")
		if !triagemAtualizada(triagemPath) {
			if err := gravarTriagem(triagemPath, slug, nil); err != nil {
				return err
			}
		}
	}
	rel.TriagensGeradas = append(rel.TriagensGeradas, slug)
	return nil
}

func garantirMemoria(dir, slug string, apply bool) error {
	if !apply {
		return nil
	}
	memPath := filepath.Join(dir, "memoria_execucao.json")
	if _, err := os.Stat(memPath); err == nil {
		return nil
	}
	return triagem.SalvarMemoria(dir, triagem.MemoriaExecucao{Processo: slug})
}
