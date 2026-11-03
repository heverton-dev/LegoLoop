package fabricar_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/specgen"
)

// specLLM devolve as respostas da SPEC em ordem (a última se repete) e guarda as mensagens de cada chamada.
type specLLM struct {
	replies []string
	err     error
	calls   [][]interview.Message
}

func (s *specLLM) Chat(_ context.Context, msgs []interview.Message) (string, error) {
	s.calls = append(s.calls, msgs)
	if s.err != nil {
		return "", s.err
	}
	return s.replies[min(len(s.calls)-1, len(s.replies)-1)], nil
}

const badSpec = "# 📋 SPEC: Ruim\n\nsó um parágrafo, sem as seções do contrato"

func repairRig(t *testing.T, llm *specLLM) *rig {
	r := newRig(t, goodPieces())
	r.deps.LLM = llm
	return r
}

func TestFabricar_RepairsASpecThatBreaksTheContract(t *testing.T) {
	llm := &specLLM{replies: []string{badSpec, specMarkdown}}
	r := repairRig(t, llm)
	r.run(t)
	if len(llm.calls) != 2 {
		t.Fatalf("1 resposta ruim + 1 correção = 2 chamadas, houve %d", len(llm.calls))
	}
	second := llm.calls[1]
	joined := ""
	for _, m := range second {
		joined += m.Role + ":" + m.Content + "\n"
	}
	if !strings.Contains(joined, "assistant:"+badSpec) {
		t.Error("a retentativa deve devolver ao modelo a resposta que foi rejeitada")
	}
	if !strings.Contains(joined, "seção obrigatória ausente") {
		t.Errorf("a retentativa deve explicar o motivo da rejeição:\n%s", joined)
	}
	if r.written(t) != 8 {
		t.Errorf("após a correção devia gravar os 8 arquivos, gravou %d", r.written(t))
	}
}

func TestFabricar_GivesUpAfterBoundedAttemptsWithoutWriting(t *testing.T) {
	llm := &specLLM{replies: []string{badSpec}}
	r := repairRig(t, llm)
	_, err := fabricar.Fabricar(context.Background(), r.req, r.deps)
	if !errors.Is(err, specgen.ErrContract) {
		t.Fatalf("esperava ErrContract, veio %v", err)
	}
	if len(llm.calls) != 3 {
		t.Errorf("o limite é 3 tentativas, houve %d chamadas", len(llm.calls))
	}
	if r.written(t) != 0 {
		t.Errorf("falhou mas gravou %d arquivos", r.written(t))
	}
}

func TestFabricar_DoesNotRetryATransportFailure(t *testing.T) {
	llm := &specLLM{err: errors.New("9Router fora do ar")}
	r := repairRig(t, llm)
	if _, err := fabricar.Fabricar(context.Background(), r.req, r.deps); err == nil {
		t.Fatal("esperava erro")
	}
	if len(llm.calls) != 1 {
		t.Errorf("falha de transporte não se repete (custo e latência): %d chamadas", len(llm.calls))
	}
}
