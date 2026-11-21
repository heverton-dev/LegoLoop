package cli

import (
	"strings"
	"testing"
)

// Achado da I04: com um modelo fraco a entrevista estaciona abaixo da meta e o usuário só saía dela
// escrevendo a mão uma mensagem pedindo o Espelho. O verbo concluir faz esse pedido de forma padronizada.
func TestEntrevistaConcluirPedeOEspelhoEConclui(t *testing.T) {
	r := newInterRig(t, "", lowReply, doneReply)
	p, code := r.entrevista(t, "iniciar", "extrair áudio")
	if code != 0 || p.Status == "COMPLETED" {
		t.Fatalf("iniciar deveria deixar a sessão aberta: exit=%d %+v", code, p)
	}
	p2, code := r.entrevista(t, "concluir", p.Sessao)
	if code != 0 {
		t.Fatalf("concluir: exit %d\n%s", code, r.errb)
	}
	if p2.Status != "COMPLETED" || p2.Espelho == nil {
		t.Fatalf("esperava COMPLETED com Espelho: %+v", p2)
	}
	if ult := r.llm.user[len(r.llm.user)-1]; !strings.Contains(strings.ToUpper(ult), "ESPELHO CANÔNICO") {
		t.Fatalf("a mensagem de concluir deveria pedir o ESPELHO CANÔNICO: %q", ult)
	}
}

func TestEntrevistaConcluirSemSessaoEUsoInvalido(t *testing.T) {
	r := newInterRig(t, "", lowReply)
	if _, code := r.entrevista(t, "concluir"); code != 2 {
		t.Fatalf("concluir sem sessão deveria sair com 2, saiu com %d", code)
	}
}
