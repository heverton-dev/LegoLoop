package distribuicao

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegistrarEnvioAcumulaESemSegredo(t *testing.T) {
	dir := t.TempDir()
	r := RegistroEnvio{ItemID: "a", Canal: "youtube", Visibilidade: VisPrivado, Modo: "agora", EnviadoEm: t0, Resposta: strings.Repeat("x", 300)}
	if err := RegistrarEnvio(dir, r); err != nil {
		t.Fatal(err)
	}
	r.ItemID = "b"
	if err := RegistrarEnvio(dir, r); err != nil {
		t.Fatal(err)
	}
	dados, _ := os.ReadFile(CaminhoEnvios(dir))
	if strings.Count(string(dados), `"item_id"`) != 2 || strings.Contains(string(dados), strings.Repeat("x", 121)) {
		t.Fatalf("dois registros com resposta cortada em 120: %s", dados)
	}
	_ = os.WriteFile(CaminhoEnvios(dir), []byte("{"), 0o644)
	if err := RegistrarEnvio(dir, r); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("json quebrado: %v", err)
	}
	dirComoArquivo := t.TempDir()
	_ = os.MkdirAll(CaminhoEnvios(dirComoArquivo), 0o755)
	if err := RegistrarEnvio(dirComoArquivo, r); err == nil {
		t.Fatal("erro de leitura deve voltar")
	}
}

func TestLiberarEExigirLiberada(t *testing.T) {
	dir := t.TempDir()
	if err := Liberar(dir, 1, t0); !errors.Is(err, ErrRevisaoReprovada) {
		t.Fatalf("sem parecer: %v", err)
	}
	if err := ExigirLiberada(dir, "renders/v1.mp4"); !errors.Is(err, ErrRevisaoReprovada) {
		t.Fatalf("sem registro: %v", err)
	}
	if err := ExigirLiberada(dir, "renders/outro.mp4"); !errors.Is(err, ErrRevisaoReprovada) {
		t.Fatalf("nome fora do padrao: %v", err)
	}
	if err := GravarParecer(dir, Parecer{Versao: 1, Veredito: VereditoEnviar, Revisor: "r"}); err != nil {
		t.Fatal(err)
	}
	if err := Liberar(dir, 1, t0); err != nil {
		t.Fatal(err)
	}
	if err := ExigirLiberada(dir, "renders/v1.mp4"); err != nil {
		t.Fatalf("liberada deve passar: %v", err)
	}
	dados, _ := os.ReadFile(CaminhoLiberacao(dir, 1))
	if !strings.Contains(string(dados), `"liberado_em": "2026-11-27T10:00:00Z"`) {
		t.Fatalf("registro de liberacao: %s", dados)
	}
	// parecer reprovado depois do registro: a liberacao deixa de valer
	_ = os.WriteFile(CaminhoParecer(dir, 1), []byte(`{"versao":1,"veredito":"NAO ENVIAR","revisor":"r","problemas":[{"versao":1,"corte":"l1","tempo":"00:01","problema":"x","correcao":"y"}]}`), 0o644)
	if err := ExigirLiberada(dir, "renders/v1.mp4"); !errors.Is(err, ErrRevisaoReprovada) {
		t.Fatalf("parecer reprovado: %v", err)
	}
	bloqueado := t.TempDir()
	_ = GravarParecer(bloqueado, Parecer{Versao: 1, Veredito: VereditoEnviar, Revisor: "r"})
	_ = os.WriteFile(filepath.Join(bloqueado, "review", "v1_liberada.json"), []byte("x"), 0o644)
	_ = os.RemoveAll(filepath.Join(bloqueado, "review", "v1_liberada.json"))
	_ = os.MkdirAll(filepath.Join(bloqueado, "review", "v1_liberada.json"), 0o755)
	if err := Liberar(bloqueado, 1, t0); err == nil {
		t.Fatal("registro que e diretorio deve falhar ao gravar")
	}
}

func TestGravarPublicacaoExigeLiberacao(t *testing.T) {
	dir, v := vitrineOK(t)
	if err := GravarPublicacao(dir, v, 0); !errors.Is(err, ErrRevisaoReprovada) {
		t.Fatalf("vitrine sem liberacao: %v", err)
	}
	_ = GravarParecer(dir, Parecer{Versao: 1, Veredito: VereditoEnviar, Revisor: "r"})
	_ = Liberar(dir, 1, t0)
	if err := GravarPublicacao(dir, v, 0); err != nil {
		t.Fatalf("vitrine liberada: %v", err)
	}
}
