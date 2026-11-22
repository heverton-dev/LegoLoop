package forja

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDesfazer(t *testing.T) {
	home, harnesses := criarHomeComCincoHarnesses(t)
	arqExistente := filepath.Join(home, ".claude", "existente.cfg")
	if err := os.WriteFile(arqExistente, []byte("valor-original\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hashAntes := hashArvore(t, home)
	data := "2026-11-22"

	geradores := map[string]GeradorFunc{
		"claude": func(_ context.Context, _ Opcao, _ string) ([]AlteracaoPlanejada, error) {
			return []AlteracaoPlanejada{
				{
					Caminho:  arqExistente,
					Conteudo: []byte("valor-modificado\n"),
					Modo:     0o644,
				},
				{
					Caminho:  filepath.Join(home, ".claude", "subnova", "novo.cfg"),
					Conteudo: []byte("arquivo-criado\n"),
					Modo:     0o644,
				},
			}, nil
		},
	}

	opAplicar := Opcao{
		Home:      home,
		Perfil:    "padrao",
		Aplicar:   true,
		Data:      data,
		SO:        "linux",
		Harnesses: harnesses,
		Geradores: geradores,
	}
	if _, err := Executar(context.Background(), opAplicar); err != nil {
		t.Fatalf("Executar --apply falhou: %v", err)
	}
	hashDepoisAplicar := hashArvore(t, home)
	if hashDepoisAplicar == hashAntes {
		t.Fatal("esperava alteração na árvore após --apply")
	}

	opDesfazer := Opcao{
		Home:     home,
		Perfil:   "padrao",
		Desfazer: data,
		Data:     data,
		SO:       "linux",
	}
	rel, err := Executar(context.Background(), opDesfazer)
	if err != nil {
		t.Fatalf("Executar --desfazer falhou: %v", err)
	}
	if rel.Veredito != VereditoOK {
		t.Fatalf("esperava veredito OK no desfazer, obteve %q", rel.Veredito)
	}
	hashDepoisDesfazer := hashArvore(t, home)
	if hashDepoisDesfazer != hashAntes {
		t.Fatalf("árvore após --desfazer divergiu do estado pré-apply: antes=%s depois=%s", hashAntes, hashDepoisDesfazer)
	}

	t.Run("desfazer data inexistente retorna ErrBackupInexistente", func(t *testing.T) {
		_, err := Executar(context.Background(), Opcao{
			Home:     home,
			Perfil:   "padrao",
			Desfazer: "1999-01-01",
		})
		if !errors.Is(err, ErrBackupInexistente) {
			t.Fatalf("esperava ErrBackupInexistente, obteve %v", err)
		}
	})
}
