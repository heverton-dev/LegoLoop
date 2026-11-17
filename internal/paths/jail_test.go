package paths

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestEntradaHostil prova AC-NF.1: entrada hostil (path "..", JSON malicioso, PID falso)
// é bloqueada pelo Path Jail (erro, nunca efeito fora da raiz sanitizada). Cobre
// anticorpos.json, memoria_execucao.json e rename.
func TestEntradaHostil(t *testing.T) {
	raiz := t.TempDir()

	t.Run("PathComEscapeDaRaiz", func(t *testing.T) {
		ataques := []string{
			"../fora.txt", "..\\fora.txt", "a/../../fora.txt", "a\\..\\..\\fora.txt",
			"C:\\Windows\\system32\\evil.txt", "\\\\server\\share\\evil.txt", "..", "",
		}
		for _, ataque := range ataques {
			if _, err := DentroDaRaiz(raiz, ataque); !errors.Is(err, ErrEscapeRaiz) {
				t.Errorf("DentroDaRaiz(%q) = %v; quero ErrEscapeRaiz", ataque, err)
			}
		}
		if _, err := DentroDaRaiz(raiz, "ok/sub/arquivo.json"); err != nil {
			t.Errorf("caminho legítimo recusado: %v", err)
		}
	})

	t.Run("JSONMaliciosoNosArquivosDeEstado", func(t *testing.T) {
		for _, nome := range []string{"anticorpos.json", "memoria_execucao.json"} {
			hostil := filepath.Join("..", "fora-da-raiz", nome)
			if _, err := CampoDeJSON(raiz, hostil); !errors.Is(err, ErrEscapeRaiz) {
				t.Errorf("CampoDeJSON(%q) = %v; quero ErrEscapeRaiz", hostil, err)
			}
			if _, err := CampoDeJSON(raiz, "dados\x00"+nome); !errors.Is(err, ErrEscapeRaiz) {
				t.Errorf("CampoDeJSON com byte nulo = %v; quero ErrEscapeRaiz", err)
			}
			alvo, err := CampoDeJSON(raiz, nome)
			if err != nil {
				t.Fatalf("CampoDeJSON legítimo(%q): %v", nome, err)
			}
			if err := os.WriteFile(alvo, []byte(`{"ok":true}`), 0o644); err != nil {
				t.Fatalf("gravando %q dentro da raiz: %v", nome, err)
			}
			if _, err := os.Stat(filepath.Join(filepath.Dir(raiz), "fora-da-raiz", nome)); !os.IsNotExist(err) {
				t.Errorf("entrada hostil criou %q fora da raiz: %v", nome, err)
			}
		}
	})

	t.Run("PIDFalso", func(t *testing.T) {
		for _, falso := range []int{0, -1, -99} {
			if err := ValidarPID(falso); !errors.Is(err, ErrPIDInvalido) {
				t.Errorf("ValidarPID(%d) = %v; quero ErrPIDInvalido", falso, err)
			}
		}
		if err := ValidarPID(os.Getpid()); err != nil {
			t.Errorf("ValidarPID(%d) recusou PID real: %v", os.Getpid(), err)
		}
	})

	t.Run("RenameHostil", func(t *testing.T) {
		de, err := DentroDaRaiz(raiz, "origem.json")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(de, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := Renomear(raiz, "origem.json", filepath.Join("..", "fora-da-raiz.json")); !errors.Is(err, ErrEscapeRaiz) {
			t.Errorf("Renomear para fora = %v; quero ErrEscapeRaiz", err)
		}
		if _, err := os.Stat(filepath.Join(filepath.Dir(raiz), "fora-da-raiz.json")); !os.IsNotExist(err) {
			t.Errorf("rename hostil criou arquivo fora da raiz: %v", err)
		}
		if err := Renomear(raiz, filepath.Join("..", "origem.json"), "destino.json"); !errors.Is(err, ErrEscapeRaiz) {
			t.Errorf("Renomear de fora = %v; quero ErrEscapeRaiz", err)
		}
		if err := Renomear(raiz, "origem.json", "destino.json"); err != nil {
			t.Errorf("rename legítimo recusado: %v", err)
		}
	})
}
