package fabricar_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/assets"
	"github.com/heverton-dev/LegoLoop/internal/descoberta"
	"github.com/heverton-dev/LegoLoop/internal/matriz"
)

type provedorUnidadesFake struct {
	so       string
	unidades []descoberta.Unidade
}

func (p provedorUnidadesFake) SO() string                           { return p.so }
func (p provedorUnidadesFake) ListarUnidades() []descoberta.Unidade { return p.unidades }

// TestTargetDrivesRestringeABusca prova que TargetDrives restringe a busca a uma única unidade
// tanto no fabricar quanto no assets.Scanner, sem emitir aviso de drives ignorados.
func TestTargetDrivesRestringeABusca(t *testing.T) {
	unidadeA := t.TempDir()
	unidadeB := t.TempDir()
	prov := provedorUnidadesFake{
		so: runtime.GOOS,
		unidades: []descoberta.Unidade{
			{Caminho: unidadeA, Tipo: descoberta.TipoFixa},
			{Caminho: unidadeB, Tipo: descoberta.TipoFixa},
		},
	}

	t.Run("fabricar consulta apenas a unidade pedida sem aviso de drives ignorados", func(t *testing.T) {
		alfa := skillApta("alfa", filepath.Join(unidadeA, "skills", "alfa", "SKILL.md"), 0.9, "construir")
		alfa.Perfil.CasosDeUso = []string{"tdd"}
		beta := skillApta("beta", filepath.Join(unidadeB, "skills", "beta", "SKILL.md"), 0.95, "verificar")
		beta.Perfil.CasosDeUso = []string{"verificacao-de-qualidade"}
		gama := skillApta("gama", filepath.Join(unidadeA, "skills", "gama", "SKILL.md"), 0.8, "verificar")
		gama.Perfil.CasosDeUso = []string{"verificacao-de-qualidade"}

		r := newRig(t, []matriz.Entrada{alfa, beta, gama})
		r.deps.ProvedorUnidades = prov
		r.req.TargetDrives = []string{unidadeA}

		res := r.run(t)

		var usadas []string
		for _, p := range res.CuratedPieces {
			usadas = append(usadas, p.Name)
		}
		if slices.Contains(usadas, "beta") {
			t.Errorf("peça da unidade B não deveria ser selecionada com TargetDrives=%q: usadas=%v", unidadeA, usadas)
		}
		if !slices.Contains(usadas, "alfa") || !slices.Contains(usadas, "gama") {
			t.Errorf("peças da unidade A deveriam ser selecionadas: usadas=%v", usadas)
		}
		for _, w := range res.Warnings {
			if strings.Contains(w, "drives pedidos foram ignorados") || strings.Contains(w, "não varre drives") {
				t.Errorf("aviso de drives ignorados não deve ser emitido: %q", w)
			}
		}
	})

	t.Run("assets.Scanner varre apenas a unidade pedida em TargetDrives", func(t *testing.T) {
		projA := filepath.Join(unidadeA, "proj-a")
		projB := filepath.Join(unidadeB, "proj-b")
		criarProjetoComSkill(t, projA, "skill-a")
		criarProjetoComSkill(t, projB, "skill-b")

		scanner := assets.NewScanner(assets.Config{
			Home:             t.TempDir(),
			AutoDiscover:     true,
			TargetDrives:     []string{unidadeA},
			ProvedorUnidades: prov,
			LookPath:         func(string) (string, error) { return "", os.ErrNotExist },
			Getenv:           func(string) string { return "" },
		})
		snap := scanner.Scan(context.Background())

		var nomes []string
		for _, sk := range snap.Of(assets.KindSkill) {
			nomes = append(nomes, sk.Name)
		}
		if !slices.Contains(nomes, "skill-a") {
			t.Errorf("esperava descobrir skill-a na unidade A, achou: %v", nomes)
		}
		if slices.Contains(nomes, "skill-b") {
			t.Errorf("não deveria varrer a unidade B quando TargetDrives=%q, achou: %v", unidadeA, nomes)
		}
	})
}

func criarProjetoComSkill(t *testing.T, projDir, skillName string) {
	t.Helper()
	skillDir := filepath.Join(projDir, ".claude", "skills", skillName)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projDir, "go.mod"), []byte("module "+skillName), 0o644); err != nil {
		t.Fatal(err)
	}
	content := "---\nname: " + skillName + "\ndescription: teste\n---\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
