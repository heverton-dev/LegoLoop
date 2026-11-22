package forja

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAplicar(t *testing.T) {
	home, harnesses := criarHomeComCincoHarnesses(t)
	data := "2026-11-22"

	op := Opcao{
		Home:      home,
		Perfil:    "padrao",
		Aplicar:   true,
		Data:      data,
		SO:        "linux",
		Harnesses: harnesses,
	}

	rel1, err := Executar(context.Background(), op)
	if err != nil {
		t.Fatalf("1º Executar --apply falhou: %v", err)
	}
	if rel1.Veredito != VereditoOK {
		t.Fatalf("esperava veredito OK, obteve %q", rel1.Veredito)
	}
	for _, et := range rel1.Etapas {
		if et.Estado != EstadoOK {
			t.Fatalf("etapa %q esperava estado ok, obteve %q", et.Nome, et.Estado)
		}
		if len(et.Alterados) == 0 {
			t.Fatalf("1º run deveria alterar arquivos em %q", et.Harness)
		}
		manifesto := filepath.Join(home, ".aidd", "backups", "forja_"+data, et.Harness, "MANIFESTO.tsv")
		dados, err := os.ReadFile(manifesto)
		if err != nil {
			t.Fatalf("MANIFESTO.tsv ausente para %q: %v", et.Harness, err)
		}
		if !strings.Contains(string(dados), "acao\torigem\tbackup\tmodo") {
			t.Fatalf("MANIFESTO.tsv sem cabeçalho esperado: %q", string(dados))
		}
	}

	rel2, err := Executar(context.Background(), op)
	if err != nil {
		t.Fatalf("2º Executar --apply falhou: %v", err)
	}
	if rel2.Veredito != VereditoOK {
		t.Fatalf("esperava veredito OK no 2º run, obteve %q", rel2.Veredito)
	}
	for _, et := range rel2.Etapas {
		if len(et.Alterados) != 0 {
			t.Errorf("2º run não é idempotente em %q: alterados=%v", et.Harness, et.Alterados)
		}
	}
}
