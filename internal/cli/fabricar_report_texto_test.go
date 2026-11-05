package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"legoloop/internal/fabricar"
	"legoloop/internal/matriz"
)

// relatorioReal imita o relatório da matriz real: muitas peças fora e uma lacuna.
func relatorioReal() *fabricar.Report {
	r := &fabricar.Report{}
	r.Varredura.PecasVistas = map[string]int{"skill": 6174, "mcp": 10, "hook": 18, "loop": 5}
	r.Varredura.MatrizIdade = 18 * time.Hour
	r.Varredura.Cortes = []string{"site", "video"}
	r.Fusao.Necessidades = matriz.Needs{Natureza: "software", Etapas: []string{"construir", "verificar"}}
	r.Fusao.Decisoes = []fabricar.Decisao{
		{Etapa: "construir", Peca: "aidd-tdd", Entrou: true, Aptidao: 0.5, Motivo: "natureza: software; etapas: construir"},
		{Etapa: "verificar", Entrou: false, Motivo: "lacuna: nenhuma peça apta a esta etapa"},
		{Etapa: "publicar", Entrou: false, Motivo: "lacuna: nenhuma peça apta a esta etapa"},
		{Peca: "site", Motivo: "natureza: a peça serve a midia, não a software"},
		{Peca: "video", Motivo: "natureza: a peça serve a midia, não a software"},
		{Peca: "sem-perfil", Motivo: "aptidão desconhecida: a peça não tem perfil"},
	}
	return r
}

// TestResumoTextoDaMatrizEDaFusao: o texto mostra a matriz (não a varredura), uma linha por etapa com a peça
// ou a lacuna, e agrupa o que ficou de fora por motivo em vez de listar peça por peça.
func TestResumoTextoDaMatrizEDaFusao(t *testing.T) {
	var buf bytes.Buffer
	PrintReportTexto(&buf, relatorioReal())
	out := buf.String()
	for _, want := range []string{
		"[Camada 2 - Matriz] Peças na matriz: 6207 | Idade: 18 h | Cortes: 2",
		"[Camada 3 - Fusão] Necessidades: software | etapas: construir, verificar",
		"  construir  aidd-tdd (aptidão 0.50) | natureza: software; etapas: construir",
		"  verificar  LACUNA | lacuna: nenhuma peça apta a esta etapa",
		"  publicar   LACUNA | lacuna: nenhuma peça apta a esta etapa",
		"  Fora: 3 peças (aptidão desconhecida: 1, natureza: 2)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("faltou %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Varredura") || strings.Contains(out, "sem-perfil") {
		t.Errorf("o texto não deve falar em varredura nem listar cada peça de fora:\n%s", out)
	}
}
