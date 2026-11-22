package espinha

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/dag"
)

func TestProcessosSemComandoExclusivoDoSO(t *testing.T) {
	raiz := filepath.Join("..", "..")
	arquivos, err := coletarDAGsProcessos(filepath.Join(raiz, "processos"))
	if err != nil {
		t.Fatalf("coletarDAGsProcessos falhou: %v", err)
	}
	if len(arquivos) == 0 {
		t.Fatal("nenhum DAG encontrado em processos/")
	}
	for _, arq := range arquivos {
		rel := caminhoRelativo(raiz, arq)
		validarDAGSemComandoSO(t, arq, rel)
	}
	for _, v := range verificarComandosExclusivosSO(raiz) {
		t.Errorf("regra comando-exclusivo-de-so disparou: %s", v.String())
	}
}

func validarDAGSemComandoSO(t *testing.T, caminho, rel string) {
	t.Helper()
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler %s: %v", rel, err)
	}
	d, err := dag.Parse(dados)
	if err != nil {
		t.Errorf("%s falhou no dag.Parse: %v", rel, err)
		return
	}
	if err := d.ValidarExecutores(); err != nil {
		t.Errorf("%s falhou no ValidarExecutores: %v", rel, err)
	}
	for _, n := range d.Nodes {
		if len(n.Command) > 0 && ehComandoExclusivoSO(n.Command[0]) {
			t.Errorf("%s:%s chama comando exclusivo de SO direto: %q", rel, n.ID, n.Command[0])
		}
	}
}
