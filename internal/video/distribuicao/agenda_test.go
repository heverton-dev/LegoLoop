package distribuicao

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func decisaoDe(p Plano, id string) Decisao {
	for _, d := range p.Decisoes {
		if d.ID == id {
			return d
		}
	}
	return Decisao{}
}

// Quente fura a fila: sai no primeiro horario livre, antes de um perene mais prioritario.
func TestQuenteFuraAFila(t *testing.T) {
	p1 := itemBase("perene", "historia do real", TipoPerene)
	p1.Prioridade = 9
	q := itemBase("quente", "noticia do dia", TipoQuente)
	a := Acervo{Itens: []Item{p1, q}}
	plano, err := Planejar(a, ConfigAgenda{TetoDiario: 1}, t0)
	if err != nil {
		t.Fatal(err)
	}
	dq, dp := decisaoDe(plano, "quente"), decisaoDe(plano, "perene")
	if dq.Acao != AcaoAgendar || !dq.Quando.Equal(time.Date(2026, 11, 27, 14, 0, 0, 0, time.UTC)) {
		t.Fatalf("quente deveria sair as 14:00 de hoje: %+v", dq)
	}
	if dp.Acao != AcaoAgendar || !dp.Quando.After(dq.Quando) {
		t.Fatalf("perene deve vir depois do quente: %+v", dp)
	}
}

func TestQuenteForaDaJanelaEExpiraSemLugar(t *testing.T) {
	q := itemBase("q", "noticia", TipoQuente)
	q.Validade = t0.Add(200 * time.Hour)
	a := Acervo{Itens: []Item{q}}
	// teto de 1 por dia e dois quentes ocupando hoje e amanha cedo empurram o terceiro alem de 48 h
	for i, id := range []string{"x1", "x2", "x3", "x4"} {
		it := itemBase(id, "tema "+id+" unico"+strings.Repeat("z", i), TipoQuente)
		it.Validade = t0.Add(50 * time.Hour)
		a.Itens = append(a.Itens, it)
	}
	plano, err := Planejar(a, ConfigAgenda{TetoDiario: 1, Horarios: []string{"14:00"}}, t0)
	if err != nil {
		t.Fatal(err)
	}
	expirou := 0
	for _, d := range plano.Decisoes {
		if d.Acao == AcaoExpirar && strings.Contains(d.Motivo, "48 h") {
			expirou++
		}
	}
	if expirou == 0 {
		t.Fatalf("esperava ao menos um quente expirado por passar de 48 h: %+v", plano.Decisoes)
	}
	acima := false
	for _, d := range plano.Decisoes {
		if strings.Contains(d.Motivo, "acima de 24 h") {
			acima = true
		}
	}
	if !acima {
		t.Fatalf("esperava aviso de quente acima de 24 h: %+v", plano.Decisoes)
	}
}

func TestValidadeVencidaExpiraEPublicadoFica(t *testing.T) {
	velho := itemBase("velho", "antigo", TipoQuente)
	velho.Validade = t0.Add(-time.Hour)
	pub := itemBase("pub", "publicado", TipoPerene)
	pub.Status = StatusPublic
	exp := itemBase("exp", "expirado", TipoPerene)
	exp.Status = StatusExpirado
	ag := itemBase("ag", "ja agendado", TipoPerene)
	ag.Status = StatusAgendado
	fim := t0.Add(48 * time.Hour)
	ag.Agendado = &fim
	plano, err := Planejar(Acervo{Itens: []Item{velho, pub, exp, ag}}, ConfigAgenda{}, t0)
	if err != nil {
		t.Fatal(err)
	}
	if d := decisaoDe(plano, "velho"); d.Acao != AcaoExpirar {
		t.Fatalf("vencido deve expirar: %+v", d)
	}
	for _, id := range []string{"pub", "exp", "ag"} {
		if d := decisaoDe(plano, id); d.ID != "" {
			t.Fatalf("item %s nao deveria receber decisao: %+v", id, d)
		}
	}
}

func TestSemAutorizacaoPula(t *testing.T) {
	it := itemBase("a", "tema", TipoPerene)
	it.Autorizado = false
	plano, _ := Planejar(Acervo{Itens: []Item{it}}, ConfigAgenda{}, t0)
	if d := decisaoDe(plano, "a"); d.Acao != AcaoPular || !strings.Contains(d.Motivo, "autorizacao") {
		t.Fatalf("sem autorizacao deve pular: %+v", d)
	}
}

func TestPereneRespeitaCadenciaETeto(t *testing.T) {
	a := Acervo{Itens: []Item{itemBase("a", "tema um", TipoPerene), itemBase("b", "tema dois", TipoPerene), itemBase("c", "tema tres", TipoPerene)}}
	plano, err := Planejar(a, ConfigAgenda{CadenciaPerene: 24 * time.Hour, TetoDiario: 5}, t0)
	if err != nil {
		t.Fatal(err)
	}
	var quando []time.Time
	for _, d := range plano.Decisoes {
		quando = append(quando, d.Quando)
	}
	for i := 1; i < len(quando); i++ {
		if quando[i].Sub(quando[i-1]) < 24*time.Hour {
			t.Fatalf("cadencia de 24 h violada: %v", quando)
		}
	}
}

func TestPereneUsaUltimoAgendadoDoAcervo(t *testing.T) {
	ag := itemBase("ag", "ja agendado", TipoPerene)
	ag.Status = StatusAgendado
	quando := time.Date(2026, 11, 28, 14, 0, 0, 0, time.UTC)
	ag.Agendado = &quando
	novo := itemBase("novo", "tema novo", TipoPerene)
	plano, _ := Planejar(Acervo{Itens: []Item{ag, novo}}, ConfigAgenda{}, t0)
	d := decisaoDe(plano, "novo")
	if d.Acao != AcaoAgendar || d.Quando.Before(quando.Add(24*time.Hour)) {
		t.Fatalf("novo perene deve respeitar 24 h do ultimo agendado: %+v", d)
	}
}

func TestSemHorarioLivreEValidadeDepoisDoHorario(t *testing.T) {
	a := Acervo{Itens: []Item{itemBase("a", "tema", TipoPerene)}}
	plano, _ := Planejar(a, ConfigAgenda{DiasAFrente: 1, TetoDiario: 1, Horarios: []string{"00:01"}}, t0.Add(20*time.Hour))
	if d := decisaoDe(plano, "a"); d.Acao != AcaoAgendar {
		t.Fatalf("deveria agendar no dia seguinte: %+v", d)
	}
	cheio := Acervo{Itens: []Item{itemBase("a", "tema um", TipoPerene), itemBase("b", "tema dois", TipoPerene), itemBase("c", "tema tres", TipoPerene)}}
	plano, _ = Planejar(cheio, ConfigAgenda{DiasAFrente: 1, TetoDiario: 1, Horarios: []string{"23:59"}, CadenciaPerene: time.Minute}, t0)
	pulou := false
	for _, d := range plano.Decisoes {
		pulou = pulou || (d.Acao == AcaoPular && strings.Contains(d.Motivo, "sem horario"))
	}
	if !pulou {
		t.Fatalf("esperava item sem horario livre: %+v", plano.Decisoes)
	}
	p := itemBase("p", "perene com validade curta", TipoPerene)
	p.Validade = t0.Add(time.Hour)
	plano, _ = Planejar(Acervo{Itens: []Item{p}}, ConfigAgenda{Horarios: []string{"23:00"}}, t0)
	if d := decisaoDe(plano, "p"); d.Acao != AcaoExpirar || !strings.Contains(d.Motivo, "validade") {
		t.Fatalf("horario depois da validade deve expirar: %+v", d)
	}
}

func TestPlanejarRecusaConfigEAcervoInvalidos(t *testing.T) {
	if _, err := Planejar(Acervo{}, ConfigAgenda{Horarios: []string{"25:99"}}, t0); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("horario invalido: %v", err)
	}
	if _, err := Planejar(Acervo{Itens: []Item{{}}}, ConfigAgenda{}, t0); err == nil {
		t.Fatal("acervo invalido deve falhar")
	}
}

func TestAlertaDeEstoqueBaixo(t *testing.T) {
	plano, _ := Planejar(Acervo{Itens: []Item{itemBase("a", "tema", TipoPerene)}}, ConfigAgenda{EstoqueMinimoPerene: 3}, t0)
	if len(plano.Alertas) != 1 || !strings.Contains(plano.Alertas[0], "1 de 3") {
		t.Fatalf("alerta esperado: %+v", plano.Alertas)
	}
	base := Acervo{}
	for _, id := range []string{"a", "b", "c"} {
		base.Itens = append(base.Itens, itemBase(id, "tema "+id, TipoPerene))
	}
	plano, _ = Planejar(base, ConfigAgenda{EstoqueMinimoPerene: 3, TetoDiario: 5, CadenciaPerene: time.Minute}, t0)
	if len(plano.Alertas) != 0 {
		t.Fatalf("estoque cheio nao deve alertar: %+v", plano.Alertas)
	}
}

func TestAplicarGravaStatus(t *testing.T) {
	q := itemBase("q", "noticia", TipoQuente)
	v := itemBase("v", "velha", TipoQuente)
	v.Validade = t0.Add(-time.Hour)
	a := Acervo{Itens: []Item{q, v}}
	plano, _ := Planejar(a, ConfigAgenda{}, t0)
	if err := a.Aplicar(plano); err != nil {
		t.Fatal(err)
	}
	if a.Itens[0].Status != StatusAgendado || a.Itens[0].Agendado == nil || a.Itens[1].Status != StatusExpirado {
		t.Fatalf("status apos aplicar: %+v", a.Itens)
	}
	if err := a.Aplicar(Plano{Decisoes: []Decisao{{ID: "zz", Acao: AcaoAgendar}}}); !errors.Is(err, ErrAcervoInvalido) {
		t.Fatalf("item inexistente: %v", err)
	}
	if err := a.Aplicar(Plano{Decisoes: []Decisao{{ID: "q", Acao: AcaoPular}}}); err != nil || a.Itens[0].Status != StatusAgendado {
		t.Fatalf("pular nao muda nada: %v", err)
	}
}

func TestGerarReels(t *testing.T) {
	curto := itemBase("curto", "tema curto", TipoPerene)
	curto.Formatos, curto.Visibilidade = []string{"9:16"}, VisPublico
	longo := itemBase("longo", "tema longo", TipoPerene)
	privado := itemBase("priv", "tema privado", TipoPerene)
	privado.Formatos = []string{"9:16"}
	pub := itemBase("pub", "tema publicado", TipoPerene)
	pub.Formatos, pub.Visibilidade, pub.Status = []string{"9:16"}, VisPublico, StatusPublic
	a := Acervo{Itens: []Item{curto, longo, privado, pub}}
	if n := a.GerarReels(); n != 1 {
		t.Fatalf("quer 1 reel, obteve %d", n)
	}
	reel := a.Buscar("curto-reel")
	if reel == nil || reel.Origem != "curto" || reel.Canais[0] != "instagram" || !reel.Autorizado {
		t.Fatalf("reel incorreto: %+v", reel)
	}
	if n := a.GerarReels(); n != 0 {
		t.Fatalf("segunda chamada nao deve duplicar: %d", n)
	}
	if err := a.Validar(); err != nil {
		t.Fatalf("acervo com reel invalido: %v", err)
	}
}
