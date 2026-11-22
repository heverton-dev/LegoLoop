package board

import (
	"errors"
	"testing"
)

type fimNo struct {
	id      string
	passou  bool
	exit    int
	detalhe string
}

func TestNotificarFimDoNivel(t *testing.T) {
	falso := false
	lvl := []RunNode{
		{NodeID: `TICK-01.1`, Name: `w1`},
		{NodeID: `TICK-01.2`, Name: `w2`},
		{NodeID: `TICK-01.3`, Name: `w3`},
		{Name: `sem-id`},
	}
	res := &RunLevelResult{Workers: []RunWorkerResult{
		{Name: `w1`, ExitCode: 0},
		{Name: `w2`, ExitCode: 0, GateOK: &falso},
	}}
	var got []fimNo
	notificarFimDoNivel(lvl, res, errors.New(`nivel falhou`), func(id string, passou bool, exit int, detalhe string) {
		got = append(got, fimNo{id, passou, exit, detalhe})
	})
	if len(got) != 3 {
		t.Fatalf(`notificações = %+v`, got)
	}
	if !got[0].passou || got[0].id != `TICK-01.1` {
		t.Errorf(`w1 devia passar: %+v`, got[0])
	}
	if got[1].passou || got[1].id != `TICK-01.2` {
		t.Errorf(`gate falso devia reprovar w2: %+v`, got[1])
	}
	if got[2].passou || got[2].id != `TICK-01.3` || got[2].exit != -1 {
		t.Errorf(`w3 sem resultado devia reprovar: %+v`, got[2])
	}
}

func TestNotificarFimDoNivelSemErroNaoReprovaOsQueFaltam(t *testing.T) {
	lvl := []RunNode{{NodeID: `TICK-01.1`, Name: `w1`}}
	var n int
	notificarFimDoNivel(lvl, &RunLevelResult{}, nil, func(string, bool, int, string) { n++ })
	if n != 0 {
		t.Fatalf(`notificações = %d`, n)
	}
}
