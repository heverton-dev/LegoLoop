package processo_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/processo"
)

// semearContrato grava o manifesto contrato.json de uma peça no cofre de teste.
func semearContrato(t *testing.T, vault, nome, entra, sai string) {
	t.Helper()
	dir := filepath.Join(vault, "skills", nome)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifesto := `{"nome":"` + nome + `","tipo":"script","naturezas":["dados"],"entra":[` + entra + `],"sai":[` + sai + `]}`
	if err := os.WriteFile(filepath.Join(dir, "contrato.json"), []byte(manifesto), 0o644); err != nil {
		t.Fatal(err)
	}
}

// arquivosDe devolve todos os arquivos do cofre com o conteúdo, para provar que nada foi gravado.
func arquivosDe(t *testing.T, raiz string) map[string]string {
	t.Helper()
	estado := map[string]string{}
	err := filepath.WalkDir(raiz, func(caminho string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := os.ReadFile(caminho)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(raiz, caminho)
		estado[rel] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return estado
}

// TestGateReprovaDuplicata prova a AC-E04.1 e a AC-E04.2: o gate anti-reimplementação
// reprova peça nova cujo I/O é estruturalmente igual ao de peça do cofre, cita a peça
// existente e nada grava. Peça com I/O distinto passa.
func TestGateReprovaDuplicata(t *testing.T) {
	igual := `{"tipo":"texto","opcional":false}`
	saidaDados := `{"tipo":"dados","promete":true}`

	t.Run("reprova I/O igual e cita a peca existente", func(t *testing.T) {
		vault := t.TempDir()
		semearContrato(t, vault, "coletor-csv", igual, saidaDados)
		antes := arquivosDe(t, vault)
		nova := processo.Contrato{
			Nome:  "extrator-csv",
			Tipo:  "script",
			Entra: []processo.Entrada{{Tipo: "texto"}},
			Sai:   []processo.Saida{{Tipo: "dados", Promete: true}},
		}
		err := processo.VerificarReimplementacao(vault, nova)
		if !errors.Is(err, processo.ErrReimplementacao) {
			t.Fatalf("esperava ErrReimplementacao, obteve: %v", err)
		}
		if !strings.Contains(err.Error(), "coletor-csv") {
			t.Errorf("o erro deve citar a peça existente: %v", err)
		}
		depois := arquivosDe(t, vault)
		if len(depois) != len(antes) {
			t.Errorf("o gate não pode gravar nada no cofre: %v -> %v", antes, depois)
		}
	})

	t.Run("reprova mesmo com a ordem do I/O trocada", func(t *testing.T) {
		vault := t.TempDir()
		semearContrato(t, vault, "coletor-csv", igual, saidaDados)
		nova := processo.Contrato{
			Nome: "extrator-csv",
			Tipo: "script",
			Entra: []processo.Entrada{
				{Tipo: "texto", Opcional: true},
				{Tipo: "texto"},
			},
			Sai: []processo.Saida{
				{Tipo: "dados", Promete: true},
				{Tipo: "texto", Promete: true},
			},
		}
		nova2 := nova
		nova2.Entra = []processo.Entrada{nova.Entra[1], nova.Entra[0]}
		nova2.Sai = []processo.Saida{nova.Sai[1], nova.Sai[0]}
		semearContrato(t, vault, "duplo-csv",
			`{"tipo":"texto"},{"tipo":"texto","opcional":true}`,
			`{"tipo":"texto","promete":true},{"tipo":"dados","promete":true}`)
		if err := processo.VerificarReimplementacao(vault, nova2); !errors.Is(err, processo.ErrReimplementacao) {
			t.Fatalf("ordem trocada não salva a duplicata: %v", err)
		}
	})

	t.Run("aprova peca com I/O distinto", func(t *testing.T) {
		vault := t.TempDir()
		semearContrato(t, vault, "coletor-csv", igual, saidaDados)
		nova := processo.Contrato{
			Nome:  "formatador-json",
			Tipo:  "skill",
			Entra: []processo.Entrada{{Tipo: "dados", Opcional: true}},
			Sai:   []processo.Saida{{Tipo: "json", Promete: true}},
		}
		if err := processo.VerificarReimplementacao(vault, nova); err != nil {
			t.Fatalf("peça com I/O distinto deve passar: %v", err)
		}
	})

	t.Run("atualizacao da mesma peca passa", func(t *testing.T) {
		vault := t.TempDir()
		semearContrato(t, vault, "coletor-csv", igual, saidaDados)
		nova := processo.Contrato{
			Nome:  "coletor-csv",
			Tipo:  "script",
			Entra: []processo.Entrada{{Tipo: "texto"}},
			Sai:   []processo.Saida{{Tipo: "dados", Promete: true}},
		}
		if err := processo.VerificarReimplementacao(vault, nova); err != nil {
			t.Fatalf("atualizar a própria peça não é reimplementação: %v", err)
		}
	})
}
