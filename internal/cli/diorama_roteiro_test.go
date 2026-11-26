package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

func TestDioramaRoteiroCLI(t *testing.T) {
	t.Run("help retorna exit 0", func(t *testing.T) {
		var out, errb bytes.Buffer
		code := Run(context.Background(), []string{"diorama-roteiro", "--help"}, Deps{Stdout: &out, Stderr: &errb})
		if code != 0 || !strings.Contains(out.String(), "diorama-roteiro") {
			t.Fatalf("esperava exit 0 com ajuda, obteve code=%d out=%q err=%q", code, out.String(), errb.String())
		}
	})

	t.Run("validar-config aprova projeto e marca coerentes e reprova incoerentes", func(t *testing.T) {
		dir := t.TempDir()
		proj := roteiro.Projeto{
			SchemaVersion: 1, Slug: "pix", Canal: "teste", Nicho: "APIs open source",
			PerfisFonte: []string{"tecnologia"}, EstruturaNarrativa: "explicador",
			Formatos: []string{"16:9"}, DuracaoAlvoS: 60, Voz: "pt-BR-ThalitaMultilingualNeural",
			PoliticaGates: roteiro.PoliticaGates{Modo: "mecanico"},
		}
		marca := roteiro.Marca{SchemaVersion: 1, Nome: "Canal Teste", Formatos: []string{"16:9"}}
		marca.Voz.Voz = "pt-BR-ThalitaMultilingualNeural"
		bProj, _ := json.Marshal(proj)
		bMarca, _ := json.Marshal(marca)
		_ = os.WriteFile(filepath.Join(dir, "projeto.json"), bProj, 0o644)
		_ = os.WriteFile(filepath.Join(dir, "marca.json"), bMarca, 0o644)

		var out, errb bytes.Buffer
		code := Run(context.Background(), []string{"diorama-roteiro", "validar-config", "--projeto", dir, "--json"}, Deps{Stdout: &out, Stderr: &errb})
		if code != 0 {
			t.Fatalf("esperava exit 0, obteve %d err=%q", code, errb.String())
		}
	})

	t.Run("validar-roteiro reprova nivel C como prova com exit 1", func(t *testing.T) {
		dir := t.TempDir()
		proj := roteiro.Projeto{
			SchemaVersion: 1, Nicho: "APIs", PerfisFonte: []string{"tecnologia"},
			Formatos: []string{"16:9"}, DuracaoAlvoS: 30,
			PoliticaGates: roteiro.PoliticaGates{Modo: "mecanico"},
		}
		sc := roteiro.Script{
			SchemaVersion: 1, Title: "Teste",
			Sources: []roteiro.Fonte{{ID: "FC", URL: "https://blog.aleatorio.org/post"}},
			Lines:   []roteiro.LinhaRoteiro{{Text: "Dado sem prova.", Natureza: "fato_consolidado", Fontes: []string{"FC"}}},
		}
		bp, _ := json.Marshal(proj)
		bs, _ := json.Marshal(sc)
		_ = os.WriteFile(filepath.Join(dir, "projeto.json"), bp, 0o644)
		_ = os.WriteFile(filepath.Join(dir, "script.json"), bs, 0o644)

		var out, errb bytes.Buffer
		code := Run(context.Background(), []string{"diorama-roteiro", "validar-roteiro", "--projeto", dir, "--json"}, Deps{Stdout: &out, Stderr: &errb})
		if code == 0 {
			t.Fatalf("esperava falha para fonte nivel C, obteve exit 0 out=%q", out.String())
		}
	})
}
