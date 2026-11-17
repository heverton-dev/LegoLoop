package golden

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/dag"
)

// vinculosVideoREA é o contrato da ligação: o nó do Processo video-rea e a peça que ele materializa.
var vinculosVideoREA = map[string]string{
	"n02_forensics_dissection": "legoloop-video-rea-transcrever",
	"n03_matrix_curation":      "legoloop-video-rea-relatorio",
}

// copiasDoDAGVideoREA são as quatro cópias do DAG do video-rea, todas com o mesmo vínculo.
var copiasDoDAGVideoREA = []string{
	"processos/legoloop-video-rea/pipeline.dag.json",
	"processos/legoloop-video-rea/legoloop-video-rea.loop.json",
	"docs/dags/2026-11-09_DAG_LEGOLOOP_VIDEO_REA.json",
	"internal/golden/testdata/processos/legoloop-video-rea.dag.json",
}

func pecasPorNo(t *testing.T, caminho string) map[string]string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(caminho)))
	if err != nil {
		t.Fatalf("ler %s: %v", caminho, err)
	}
	var doc struct {
		Nodes []struct {
			ID   string `json:"id"`
			Peca string `json:"peca"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("json de %s: %v", caminho, err)
	}
	out := map[string]string{}
	for _, n := range doc.Nodes {
		out[n.ID] = n.Peca
	}
	return out
}

func TestVideoREALigaTranscreverEmN02ERelatorioEmN03(t *testing.T) {
	for _, caminho := range copiasDoDAGVideoREA {
		got := pecasPorNo(t, caminho)
		for no, peca := range vinculosVideoREA {
			if got[no] != peca {
				t.Errorf("%s: nó %s liga a peça %q, esperava %q", caminho, no, got[no], peca)
			}
		}
		for no, peca := range got {
			if _, ligado := vinculosVideoREA[no]; !ligado && peca != "" {
				t.Errorf("%s: nó %s não devia ter peça, tem %q", caminho, no, peca)
			}
		}
	}
}

func TestVideoREAPecasLigadasExistemNoRepositorio(t *testing.T) {
	raiz := filepath.Join("..", "..")
	for no, peca := range vinculosVideoREA {
		skill := filepath.Join(raiz, "skills", peca, "SKILL.md")
		if _, err := os.Stat(skill); err != nil {
			t.Errorf("nó %s liga %s, mas a Skill não existe: %v", no, peca, err)
		}
		if !temDAGDoProcesso(raiz, peca) {
			t.Errorf("nó %s liga %s, mas o DAG do Processo não existe nas pastas de leitura", no, peca)
		}
	}
}

// temDAGDoProcesso procura o DAG da peça nas pastas de Processos dos dois dialetos.
func temDAGDoProcesso(raiz, peca string) bool {
	for _, base := range caminhos.DirsLeitura(raiz) {
		if _, err := os.Stat(caminhos.CaminhoDAGLeitura(filepath.Join(base, peca))); err == nil {
			return true
		}
	}
	return false
}

func TestVideoREAParseiaComPecaENaoPerdeComando(t *testing.T) {
	d := dagDaFixture(t, "legoloop-video-rea")
	for _, n := range d.Nodes {
		peca, ligado := vinculosVideoREA[n.ID]
		if ligado && n.Peca != peca {
			t.Errorf("nó %s: Peca = %q, esperava %q", n.ID, n.Peca, peca)
		}
		if ligado && n.Type == dag.TypeScript && len(n.Command) == 0 {
			t.Errorf("nó %s perdeu o comando ao ganhar a peça", n.ID)
		}
	}
}
