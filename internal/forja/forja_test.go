package forja

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func hashArvore(t *testing.T, raiz string) string {
	t.Helper()
	var entradas []string
	err := filepath.WalkDir(raiz, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(raiz, p)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			entradas = append(entradas, "D:"+rel)
			return nil
		}
		dados, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		soma := sha256.Sum256(dados)
		entradas = append(entradas, "F:"+rel+":"+hex.EncodeToString(soma[:]))
		return nil
	})
	if err != nil {
		t.Fatalf("falha ao calcular hash da árvore %s: %v", raiz, err)
	}
	sort.Strings(entradas)
	h := sha256.New()
	for _, e := range entradas {
		h.Write([]byte(e))
		h.Write([]byte{'\n'})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func criarHomeComCincoHarnesses(t *testing.T) (string, []string) {
	t.Helper()
	home := t.TempDir()
	harnesses := []string{"claude", "gemini", "opencode", "cursor", "codex"}
	for _, h := range harnesses {
		dir := filepath.Join(home, "."+h)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		arq := filepath.Join(dir, "usuario.txt")
		if err := os.WriteFile(arq, []byte("conteudo-"+h+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return home, harnesses
}

func TestSimulacaoNaoGrava(t *testing.T) {
	home, harnesses := criarHomeComCincoHarnesses(t)
	antes := hashArvore(t, home)

	op := Opcao{
		Home:      home,
		Perfil:    "padrao",
		Aplicar:   false,
		Data:      "2026-11-22",
		SO:        "linux",
		Harnesses: harnesses,
	}
	rel, err := Executar(context.Background(), op)
	if err != nil {
		t.Fatalf("Executar em simulação falhou: %v", err)
	}
	depois := hashArvore(t, home)
	if antes != depois {
		t.Fatalf("simulação alterou o HOME: antes=%s depois=%s", antes, depois)
	}
	if rel.Veredito != VereditoOK {
		t.Fatalf("esperava veredito %q, obteve %q", VereditoOK, rel.Veredito)
	}
	if len(rel.Etapas) != len(harnesses) {
		t.Fatalf("esperava %d etapas, obteve %d", len(harnesses), len(rel.Etapas))
	}
	for _, et := range rel.Etapas {
		if et.Estado != EstadoPendente {
			t.Errorf("etapa %q deveria estar %q em simulação, obteve %q", et.Nome, EstadoPendente, et.Estado)
		}
		if len(et.Alterados) == 0 {
			t.Errorf("etapa %q deveria listar alterações pendentes", et.Nome)
		}
	}
}
