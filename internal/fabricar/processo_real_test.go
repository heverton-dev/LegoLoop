package fabricar_test

import (
	"os"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/processo"
)

// TestProcessoRealGravaOProcessoEOTemplateContaAMesmaHistoria: o processo.json gravado é o mesmo Processo de
// onde saíram as projeções. O relatório da Tríade conta um critério, um ticket e um gate por etapa, e o
// TEMPLATE repete a etapa, o artefato e o gate de cada uma. Nada disso vem do cérebro.
func TestProcessoRealGravaOProcessoEOTemplateContaAMesmaHistoria(t *testing.T) {
	r := newRig(t, goodPieces())
	res := r.run(t)
	p, err := processo.Ler(res.Artifacts.ProcessoPath)
	if err != nil {
		t.Fatalf("o processo.json deve existir e validar: %v", err)
	}
	if p.Natureza != "software" || len(p.Etapas) != 2 || p.Etapas[0].Peca != "alfa" || p.Etapas[1].Peca != "beta" {
		t.Fatalf("processo gravado errado: %+v", p)
	}
	tri := res.Report.Triade
	if !tri.Coerente || tri.ACs != 2 || tri.Tickets != 2 || len(tri.Gates) != 2 || tri.Gates[0] != p.Etapas[0].Gate {
		t.Errorf("a Tríade deve vir do Processo: %+v", tri)
	}
	tpl := fabricar.EmitTemplate(p)
	for _, e := range p.Etapas {
		for _, want := range []string{"`" + e.ID + "`", e.Produz.Caminho, e.Gate, "skill `" + e.Peca + "`"} {
			if !strings.Contains(tpl, want) {
				t.Errorf("o TEMPLATE deve trazer %q da etapa %s", want, e.ID)
			}
		}
	}
	spec := read(t, res.Artifacts.SpecPath)
	if !strings.HasPrefix(spec, "# 📋 SPEC: "+p.Nome) || !strings.Contains(spec, "`AC-E01.2`: **Etapa `"+p.Etapas[1].ID+"`") {
		t.Errorf("o arquivo da SPEC deve ser a SPEC projetada, com um critério por etapa:\n%.300s", spec)
	}
	if r.llm.specCalls != 0 {
		t.Errorf("a SPEC nasce do Processo, não do cérebro: %d pedidos", r.llm.specCalls)
	}
}

// TestProcessoRealNaturezaDoIntentEDoProcessoConcordam: com a natureza dados confirmada na entrevista, o
// processo.json é de dados e o INTENT usa as peças sem comandos de Go. Os dois nascem da mesma natureza.
func TestProcessoRealNaturezaDoIntentEDoProcessoConcordam(t *testing.T) {
	r := newRig(t, pecasDeDados())
	r.deps.Interview = &fabricar.Clarification{Mirror: espelhoPrecos(), Confidence: 98, Reached: true, Interviewed: true, Natureza: "dados"}
	res := r.run(t)
	p, err := processo.Ler(res.Artifacts.ProcessoPath)
	if err != nil || p.Natureza != "dados" {
		t.Fatalf("o processo deve ser de dados: %v %+v", err, p)
	}
	if p.Objetivo != espelhoPrecos().Outcome {
		t.Errorf("o objetivo do processo vem do espelho da entrevista, não do pedido: %q", p.Objetivo)
	}
	if intent := read(t, res.Artifacts.IntentPath); strings.Contains(intent, "go test") {
		t.Error("o INTENT de um processo de dados não pode falar de go test")
	}
}

// TestProcessoRealDryRunNaoGravaOProcesso: no --dry-run o caminho do processo.json é calculado e nada é gravado.
func TestProcessoRealDryRunNaoGravaOProcesso(t *testing.T) {
	r := newRig(t, goodPieces())
	r.req.DryRun = true
	res := r.run(t)
	if res.Artifacts.ProcessoPath == "" {
		t.Fatal("o dry-run deve descrever onde o processo seria gravado")
	}
	if _, err := os.Stat(res.Artifacts.ProcessoPath); !os.IsNotExist(err) {
		t.Errorf("o dry-run não pode gravar o processo.json: %v", err)
	}
}
