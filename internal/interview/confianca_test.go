package interview

import (
	"context"
	"testing"
)

// seqLLM responde em ordem; a última resposta se repete.
type seqLLM struct {
	replies []string
	i       int
}

func (s *seqLLM) Chat(context.Context, []Message) (string, error) {
	r := s.replies[min(s.i, len(s.replies)-1)]
	s.i++
	return r, nil
}

const espelho60 = "HIPÓTESE: h\nCONFIANÇA: 60%\nESPELHO CANÔNICO:\n- Resultado Esperado: algo\n- Quem se Beneficia: eu"

func TestConfiancaHonestaParserNaoSobeParaNoventaECinco(t *testing.T) {
	turn := ParseTurnResult(espelho60)
	if turn.Confidence != 60 {
		t.Errorf("a confiança medida foi 60%%, o parser devolveu %d%%", turn.Confidence)
	}
	if turn.Mirror == nil || turn.Mirror.Outcome != "algo" {
		t.Errorf("o espelho deve ser lido mesmo com confiança baixa: %+v", turn.Mirror)
	}
}

func TestConfiancaHonestaParserMarcaSeAConfiancaVeio(t *testing.T) {
	if !ParseTurnResult("CONFIANÇA: 40%").ConfidenceReported {
		t.Error("CONFIANÇA presente deve marcar ConfidenceReported")
	}
	if ParseTurnResult("PERGUNTA: x").ConfidenceReported {
		t.Error("sem CONFIANÇA, ConfidenceReported deve ser falso")
	}
}

func TestConfiancaHonestaSemCampoMantemAUltimaMedida(t *testing.T) {
	llm := &seqLLM{replies: []string{
		"HIPÓTESE: h\nCONFIANÇA: 70%\nPERGUNTA: a?\nPALPITE: b",
		"HIPÓTESE: h\nPERGUNTA: c?\nPALPITE: d", // sem CONFIANÇA
	}}
	eng := NewEngine(llm)
	sess, t1, err := eng.Start(context.Background(), "x")
	if err != nil || t1.Confidence != 70 {
		t.Fatalf("turno 1: %v, %d%%", err, t1.Confidence)
	}
	t2, err := eng.Step(context.Background(), sess.ID, "ok")
	if err != nil {
		t.Fatal(err)
	}
	if t2.Confidence != 70 {
		t.Errorf("sem o campo, a confiança deve seguir 70%%, não %d%%", t2.Confidence)
	}
}

func TestConfiancaHonestaPrimeiroTurnoSemCampoNaoInventaCinquenta(t *testing.T) {
	_, turn, err := NewEngine(&seqLLM{replies: []string{"PERGUNTA: a?\nPALPITE: b"}}).Start(context.Background(), "x")
	if err != nil {
		t.Fatal(err)
	}
	if turn.Confidence != 10 {
		t.Errorf("sem medida, vale o piso da sessão (10%%), não %d%%", turn.Confidence)
	}
}
