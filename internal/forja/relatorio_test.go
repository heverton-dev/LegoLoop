package forja

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSegredosColetarAmbiente(t *testing.T) {
	env := []string{
		"PATH=/usr/bin",
		"API_KEY=chave-secreta-123",
		"github_token=ghp_abc987",
		"MY_APP_SECRET=segredo-alfa",
		"DB_PASSWORD=senha-banco",
		"EMPTY_TOKEN=",
		"SEM_IGUAL",
	}
	achados := ColetarSegredosAmbiente(env)
	esperados := []string{"chave-secreta-123", "ghp_abc987", "segredo-alfa", "senha-banco"}
	if len(achados) != len(esperados) {
		t.Fatalf("esperava %d segredos, obteve %d: %v", len(esperados), len(achados), achados)
	}
	junto := strings.Join(achados, "\n")
	for _, esp := range esperados {
		if !strings.Contains(junto, esp) {
			t.Errorf("segredo esperado ausente: %q", esp)
		}
	}
}

func TestSegredos(t *testing.T) {
	env := []string{
		"OPENAI_API_KEY=sk-segredo-longo-999",
		"SHORT_KEY=sk-segredo",
		"GH_TOKEN=tok_xyz_777",
		"APP_SECRET=sec_alfa_555",
		"NORMAL_VAR=valor-publico",
	}
	segredos := ColetarSegredosAmbiente(env)
	dados := map[string]any{
		"veredito": "OK",
		"mensagem": "falha com sk-segredo-longo-999 e tok_xyz_777",
		"detalhes": map[string]any{
			"env_vazado": "sec_alfa_555",
			"publico":    "valor-publico",
		},
	}

	md := RenderizarRelatorioMD(dados, segredos)
	verificarSemSegredos(t, "RELATORIO.md", md, segredos)
	if !strings.Contains(md, "[REDIGIDO]") || !strings.Contains(md, "valor-publico") {
		t.Fatalf("RELATORIO.md mal sanitizado: %s", md)
	}

	brutoJSON, err := RenderizarRelatorioJSON(dados, segredos)
	if err != nil {
		t.Fatalf("erro ao renderizar JSON: %v", err)
	}
	verificarSemSegredos(t, "RELATORIO.json", string(brutoJSON), segredos)
	var valido map[string]any
	if err := json.Unmarshal(brutoJSON, &valido); err != nil {
		t.Fatalf("RELATORIO.json inválido após sanitização: %v", err)
	}

	terminal := SanitizarTexto("erro: sk-segredo-longo-999 / sk-segredo", segredos)
	verificarSemSegredos(t, "terminal", terminal, segredos)
	if terminal != "erro: [REDIGIDO] / [REDIGIDO]" {
		t.Fatalf("sanitização de terminal inesperada: %q", terminal)
	}
}

func verificarSemSegredos(t *testing.T, alvo, texto string, segredos []string) {
	t.Helper()
	for _, s := range segredos {
		if strings.Contains(texto, s) {
			t.Errorf("%s vazou segredo %q", alvo, s)
		}
	}
}
