package fabricar_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"legoloop/internal/cli"
	"legoloop/internal/fabricar"
)

func TestReportUnico(t *testing.T) {
	r := newRig(t, goodPieces())
	r.req.DryRun = true
	res, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
	if err != nil {
		t.Fatalf("Fabricar: %v", err)
	}
	if res.Report == nil {
		t.Fatal("Report deve ser retornado no Result")
	}
	rep := res.Report

	// Serialização JSON (objeto compartilhado pelas 3 portas)
	data, err := json.Marshal(rep)
	if err != nil {
		t.Fatalf("falha ao serializar Report: %v", err)
	}
	var roundtrip fabricar.Report
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("falha ao desserializar Report: %v", err)
	}

	// 5 seções presentes e não vazias nas camadas executadas
	if rep.Clarificacao.Espelho.Outcome == "" {
		t.Errorf("Clarificacao.Espelho vazio: %+v", rep.Clarificacao)
	}
	if len(rep.Varredura.PecasVistas) == 0 {
		t.Errorf("Varredura.PecasVistas vazio: %+v", rep.Varredura)
	}
	if len(rep.Fusao.Necessidades.Etapas) == 0 && len(rep.Fusao.Decisoes) == 0 {
		t.Errorf("Fusao vazia: %+v", rep.Fusao)
	}
	if rep.Triade.ACs == 0 && len(rep.Triade.Gates) == 0 {
		t.Errorf("Triade vazia: %+v", rep.Triade)
	}
	if len(rep.Gravacao.Arquivos) == 0 {
		t.Errorf("Gravacao.Arquivos vazio: %+v", rep.Gravacao)
	}
	if rep.Cartao.Acao == "" {
		t.Errorf("Cartao vazio: %+v", rep.Cartao)
	}

	// CLI em texto: imprime resumo por camada
	var buf bytes.Buffer
	cli.PrintReportTexto(&buf, rep)
	out := buf.String()

	for _, check := range []string{
		"Clarificação",
		"Varredura",
		"Fusão",
		"Tríade",
		"Gravação",
		"Confiança",
		"Peças vistas",
		"Decisões",
		"Critérios",
		"Arquivos",
	} {
		if !strings.Contains(out, check) {
			t.Errorf("saída de texto não contém %q:\n%s", check, out)
		}
	}
}
