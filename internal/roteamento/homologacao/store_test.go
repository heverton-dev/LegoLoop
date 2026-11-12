package homologacao

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestStoreHomologacao(t *testing.T) {
	home := t.TempDir()
	caminho := Caminho(home)
	if caminho != filepath.Join(home, ".aidd", "homologacao.json") {
		t.Fatalf("caminho inesperado: %s", caminho)
	}
	vazio, err := Ler(caminho)
	if err != nil || len(vazio) != 0 {
		t.Fatalf("arquivo ausente deve dar store vazio sem erro: %v %v", vazio, err)
	}
	if Chave("claude", "oficial") != "claude|oficial" {
		t.Fatal("chave errada")
	}
	s := Store{Chave("claude", "oficial"): {
		Estado: EstadoHomologado, Data: "2026-11-12", PesoMB: 420,
		Criterios: map[string]bool{"C1": true, "C2": false}, Motivo: "x",
	}}
	if err := s.Gravar(caminho); err != nil {
		t.Fatal(err)
	}
	lido, err := Ler(caminho)
	if err != nil || !reflect.DeepEqual(lido, s) {
		t.Fatalf("ida e volta falhou: %v %v", lido, err)
	}
	if err := os.WriteFile(caminho, []byte("{invalido"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Ler(caminho); err == nil {
		t.Fatal("JSON inválido deve dar erro")
	}
}
