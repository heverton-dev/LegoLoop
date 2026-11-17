package processo_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"legoloop/internal/processo"
	"legoloop/internal/studiodoc"
)

// dagParidadeStudio é o DAG mínimo da prova de paridade: dois nós na ordem a, b.
const dagParidadeStudio = `{"id":"legoloop-x","name":"legoloop-x","description":"prova","version":"1.0.0","nodes":[
 {"id":"a","name":"Etapa A","type":"skill","ref":"skills/x","routing":{"mode":"manual","model":"none"},
  "output":{"artifact":"a.txt","gate":"exit_code == 0"}},
 {"id":"b","name":"Etapa B","type":"skill","ref":"skills/y","depends_on":["a"],"routing":{"mode":"manual","model":"none"},
  "output":{"artifact":"b.txt","gate":"exit_code == 0"}}]}`

// TestDerivarStudioDerivaDoDAG prova que o documento do Studio nasce dos nós do DAG
// e já passa no próprio gate de paridade (é derivado, não é fonte).
func TestDerivarStudioDerivaDoDAG(t *testing.T) {
	agora := time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)
	b, err := processo.DerivarStudio([]byte(dagParidadeStudio), agora)
	if err != nil {
		t.Fatalf("DerivarStudio: %v", err)
	}
	var s studiodoc.Processo
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatalf("documento do Studio ilegível: %v", err)
	}
	if len(s.Nodes) != 2 || s.Nodes[0].ID != "a" || s.Nodes[1].ID != "b" {
		t.Errorf("nós derivados errados: %+v", s.Nodes)
	}
	if !s.CreatedAt.Equal(agora) || !s.UpdatedAt.Equal(agora) {
		t.Errorf("timestamps não derivados: %v %v", s.CreatedAt, s.UpdatedAt)
	}
	if err := processo.ParidadeStudio([]byte(dagParidadeStudio), b); err != nil {
		t.Errorf("a derivação não passa na paridade: %v", err)
	}
}

// TestParidadeStudioFalhaEmDrift prova o gate: nó diferente, faltando ou fora de ordem
// entre DAG e documento do Studio derruba a paridade com ErrDriftStudio.
func TestParidadeStudioFalhaEmDrift(t *testing.T) {
	bom, err := processo.DerivarStudio([]byte(dagParidadeStudio), time.Now())
	if err != nil {
		t.Fatalf("DerivarStudio: %v", err)
	}
	var s studiodoc.Processo
	_ = json.Unmarshal(bom, &s)

	renomeado, _ := json.Marshal(mudarNome(s, 1, "invasor"))
	semNo, _ := json.Marshal(mudarNome(removerNo(s, 1), 0, "Etapa A"))
	foraDeOrdem, _ := json.Marshal(inverterNos(s))

	for nome, adulterado := range map[string][]byte{
		"nó renomeado": renomeado, "nó removido": semNo, "nós fora de ordem": foraDeOrdem,
	} {
		err := processo.ParidadeStudio([]byte(dagParidadeStudio), adulterado)
		if !errors.Is(err, processo.ErrDriftStudio) {
			t.Errorf("%s: erro = %v, quero ErrDriftStudio", nome, err)
		}
	}
}

// TestParidadeStudioIgnoraTimestamps prova que só os nós contam: o carimbo de
// criação/edição do documento derivado não é drift.
func TestParidadeStudioIgnoraTimestamps(t *testing.T) {
	b, err := processo.DerivarStudio([]byte(dagParidadeStudio), time.Now())
	if err != nil {
		t.Fatalf("DerivarStudio: %v", err)
	}
	var s studiodoc.Processo
	_ = json.Unmarshal(b, &s)
	s.CreatedAt = s.CreatedAt.Add(96 * time.Hour)
	s.UpdatedAt = s.UpdatedAt.Add(96 * time.Hour)
	outro, _ := json.Marshal(s)
	if err := processo.ParidadeStudio([]byte(dagParidadeStudio), outro); err != nil {
		t.Errorf("timestamp não deveria ser drift: %v", err)
	}
}

func mudarNome(s studiodoc.Processo, i int, nome string) studiodoc.Processo {
	s.Nodes[i].Name = nome
	return s
}

func removerNo(s studiodoc.Processo, i int) studiodoc.Processo {
	s.Nodes = append(s.Nodes[:i], s.Nodes[i+1:]...)
	return s
}

func inverterNos(s studiodoc.Processo) studiodoc.Processo {
	for i, j := 0, len(s.Nodes)-1; i < j; i, j = i+1, j-1 {
		s.Nodes[i], s.Nodes[j] = s.Nodes[j], s.Nodes[i]
	}
	return s
}
