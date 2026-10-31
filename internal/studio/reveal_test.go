package studio

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRevealSafePath(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(f, []byte("ok"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := &Server{ProjectDir: tmp, Home: tmp}

	abs, isDir, err := s.resolveSafePath(f)
	if err != nil || isDir || abs != f {
		t.Fatalf("resolve falhou: %v, %v, %s", err, isDir, abs)
	}

	// Path traversal deve ser rejeitado
	if _, _, err := s.resolveSafePath("../out.txt"); err == nil {
		t.Fatal("path traversal deveria falhar")
	}
}

func TestHandleAssetRevealInvalid(t *testing.T) {
	s := &Server{ProjectDir: t.TempDir()}
	req := httptest.NewRequest("POST", "/api/assets/reveal", bytes.NewBufferString(`{"path":"non-existent"}`))
	rec := httptest.NewRecorder()
	s.handleAssetReveal(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400 para caminho inexistente, obteve %d", rec.Code)
	}
}
