package forja

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOrcamentoEstourado(t *testing.T) {
	orc := OrcamentoTokens{
		Versao: "1.0.0",
		Tetos: map[string]int{
			"claude":   100,
			"opencode": 80,
		},
	}
	medidos := map[string]int{
		"claude":   90,
		"opencode": 135,
	}

	res := VerificarOrcamento(orc, medidos)
	if res.Aprovado {
		t.Fatal("esperava reprovacao quando medido excede teto")
	}
	if len(res.Excessos) != 1 {
		t.Fatalf("esperava 1 excesso, obteve %d: %+v", len(res.Excessos), res.Excessos)
	}
	ex := res.Excessos[0]
	if ex.Harness != "opencode" || ex.Medido != 135 || ex.Teto != 80 || ex.Excesso != 55 {
		t.Fatalf("excesso inesperado: %+v", ex)
	}
	texto := FormatarExcessos(res)
	if !strings.Contains(texto, "opencode") || !strings.Contains(texto, "55") {
		t.Fatalf("saida deve citar harness e excesso, obteve: %s", texto)
	}
}

func TestOrcamentoRespeitado(t *testing.T) {
	orc := OrcamentoTokens{
		Versao: "1.0.0",
		Tetos: map[string]int{
			"claude":   100,
			"opencode": 80,
		},
	}
	medidos := map[string]int{
		"claude":   100,
		"opencode": 59,
	}

	res := VerificarOrcamento(orc, medidos)
	if !res.Aprovado {
		t.Fatalf("esperava aprovacao quando todos os medidos <= teto, excessos: %+v", res.Excessos)
	}
	if len(res.Excessos) != 0 {
		t.Fatalf("esperava zero excessos, obteve %+v", res.Excessos)
	}
}

func TestCarregarOrcamentoValidoEInvalido(t *testing.T) {
	dir := t.TempDir()
	valido := filepath.Join(dir, "orcamento.json")
	conteudo := `{"versao":"1.0.0","tetos":{"claude":120,"codex":120}}`
	if err := os.WriteFile(valido, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	orc, err := CarregarOrcamento(valido)
	if err != nil {
		t.Fatalf("erro inesperado ao carregar orcamento valido: %v", err)
	}
	if orc.Tetos["claude"] != 120 || orc.Tetos["codex"] != 120 {
		t.Fatalf("tetos inesperados: %+v", orc.Tetos)
	}

	vazio := filepath.Join(dir, "vazio.json")
	if err := os.WriteFile(vazio, []byte(`{"versao":"1.0.0","tetos":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := CarregarOrcamento(vazio); err == nil {
		t.Fatal("esperava erro para orcamento sem tetos")
	}

	negativo := filepath.Join(dir, "negativo.json")
	if err := os.WriteFile(negativo, []byte(`{"tetos":{"claude":-1}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := CarregarOrcamento(negativo); err == nil {
		t.Fatal("esperava erro para teto negativo")
	}
}

func TestCarregarOrcamentoGovernanca(t *testing.T) {
	raiz := localizarRaizModulo(t)
	caminho := filepath.Join(raiz, "docs", "governanca", "orcamento_tokens.json")
	orc, err := CarregarOrcamento(caminho)
	if err != nil {
		t.Fatalf("falha ao carregar %s: %v", caminho, err)
	}
	if len(orc.Tetos) < 20 {
		t.Fatalf("esperava ao menos 20 harnesses no orcamento de governanca, obteve %d", len(orc.Tetos))
	}
	for _, h := range []string{"agy-cli", "claude", "codex", "opencode"} {
		if orc.Tetos[h] <= 0 {
			t.Errorf("teto ausente ou invalido para harness %s", h)
		}
	}
}
