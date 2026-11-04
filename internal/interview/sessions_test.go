package interview

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEntrevistaPassoAPassoSessaoSobreviveAOutroProcesso(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "sessions")
	store := Store{Dir: dir}
	first := &recLLM{reply: "HIPÓTESE: h\nCONFIANÇA: 40%\nPERGUNTA: Qual o formato?\nPALPITE: mp3"}
	sess, _, err := NewEngine(first).Start(context.Background(), "extrair áudio")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(sess); err != nil {
		t.Fatal(err)
	}
	// "outro processo": motor novo, sem memória da sessão.
	loaded, err := store.Load(sess.ID)
	if err != nil {
		t.Fatal(err)
	}
	second := &recLLM{reply: "HIPÓTESE: h\nCONFIANÇA: 70%\nPERGUNTA: Qual o bitrate?\nPALPITE: 128"}
	eng := NewEngine(second)
	eng.Adopt(loaded)
	turn, err := eng.Step(context.Background(), loaded.ID, "mp3")
	if err != nil || turn.Confidence != 70 {
		t.Fatalf("Step: %v, %+v", err, turn)
	}
	sent := allText(second.calls[0])
	for _, want := range []string{"extrair áudio", "Qual o formato?", "mp3"} {
		if !strings.Contains(sent, want) {
			t.Errorf("o histórico não chegou ao cérebro: falta %q", want)
		}
	}
}

func TestEntrevistaPassoAPassoStoreRecusaIdPerigoso(t *testing.T) {
	store := Store{Dir: t.TempDir()}
	for _, id := range []string{"../x", "intv_1/../../x", "", "intv_", "intv_12a", `intv_1\..\x`, "outro_123", "x_intv_1", "../intv_1", " intv_1"} {
		if _, err := store.Load(id); !errors.Is(err, ErrSessaoInvalida) {
			t.Errorf("Load(%q) = %v, queria ErrSessaoInvalida", id, err)
		}
		if err := store.Save(&Session{ID: id}); !errors.Is(err, ErrSessaoInvalida) {
			t.Errorf("Save(%q) = %v, queria ErrSessaoInvalida", id, err)
		}
	}
}

func TestEntrevistaPassoAPassoStoreSessaoAusente(t *testing.T) {
	if _, err := (Store{Dir: t.TempDir()}).Load("intv_123"); !errors.Is(err, ErrSessaoNaoEncontrada) {
		t.Errorf("queria ErrSessaoNaoEncontrada, veio %v", err)
	}
}

func TestEntrevistaPassoAPassoStoreGravaSemRestos(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "novo", "sessions") // a pasta ainda não existe
	store := Store{Dir: dir}
	if err := store.Save(&Session{ID: "intv_9", Theme: "x"}); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(&Session{ID: "intv_9", Theme: "y"}); err != nil { // regravar substitui
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 || entries[0].Name() != "intv_9.json" {
		t.Errorf("sobraram arquivos temporários: %v", entries)
	}
	got, err := store.Load("intv_9")
	if err != nil || got.Theme != "y" {
		t.Errorf("a regravação deve valer: %v %+v", err, got)
	}
}

func TestEntrevistaPassoAPassoStoreRecusaJSONCorrompido(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "intv_5.json"), []byte("{quebrado"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := (Store{Dir: dir}).Load("intv_5"); err == nil || errors.Is(err, ErrSessaoNaoEncontrada) {
		t.Errorf("arquivo corrompido deve dar erro próprio, veio %v", err)
	}
}
