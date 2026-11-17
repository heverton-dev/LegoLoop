package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"legoloop/internal/caminhos"
	"legoloop/internal/processo"
	"legoloop/internal/studiodoc"
)

// flagRegenerarStudio faz o validate regravar o documento do Studio a partir do DAG.
const flagRegenerarStudio = "--regenerar-studio"

// conferirStudioDerivado aplica o gate de paridade do Studio: o <slug>.studio.json ao lado do
// DAG canônico é derivado do DAG e seus nós têm de bater com os do DAG, na mesma ordem.
// Drift derruba o validate. Com --regenerar-studio o documento é regravado a partir do DAG;
// sem documento do Studio, não há o que conferir e nada é escrito.
func conferirStudioDerivado(path string, args []string, d Deps) error {
	dir, base := filepath.Dir(path), filepath.Base(path)
	if base != caminhos.ArquivoDAG && base != caminhos.ArquivoProcessoDAG {
		return nil
	}
	regenerar := slices.Contains(args, flagRegenerarStudio)
	dagB, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	studio := filepath.Join(dir, filepath.Base(dir)+caminhos.SufixoStudio)
	sb, err := os.ReadFile(studio)
	if os.IsNotExist(err) {
		if regenerar {
			return regenerarStudio(studio, dagB, time.Now().UTC(), d)
		}
		return nil
	}
	if err != nil {
		return err
	}
	if regenerar {
		return regenerarStudio(studio, dagB, criadoDe(sb), d)
	}
	return processo.ParidadeStudio(dagB, sb)
}

// regenerarStudio regrava o documento do Studio a partir dos bytes do DAG, preservando o
// carimbo de criação já registrado no documento anterior.
func regenerarStudio(studio string, dagB []byte, criadoEm time.Time, d Deps) error {
	novo, err := processo.DerivarStudio(dagB, criadoEm)
	if err != nil {
		return err
	}
	if err := os.WriteFile(studio, novo, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(d.Stdout, "✓ %s regenerado a partir do DAG\n", filepath.Base(studio))
	return nil
}

// criadoDe devolve o carimbo de criação do documento derivado; sem carimbo, vale agora.
func criadoDe(sb []byte) time.Time {
	var s studiodoc.Processo
	if err := json.Unmarshal(sb, &s); err == nil && !s.CreatedAt.IsZero() {
		return s.CreatedAt
	}
	return time.Now().UTC()
}
