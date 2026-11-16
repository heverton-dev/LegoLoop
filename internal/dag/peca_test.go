package dag

import (
	"strings"
	"testing"
)

func dagComPeca(peca string) []byte {
	return []byte(`{
  "name": "demo",
  "version": "1.0.0",
  "nodes": [
    {
      "id": "n01",
      "type": "script",
      "command": ["legoloop", "listar"],
      "peca": "` + peca + `",
      "routing": {"mode": "suggested", "model": "gemini-3.5-flash-lite", "fallback": "haiku"},
      "output": {"artifact": "out/a.json", "gate": "file_exists('out/a.json')"}
    }
  ]
}`)
}

func TestNoGuardaAPecaLigada(t *testing.T) {
	d, err := Parse(dagComPeca("legoloop-video-rea-transcrever"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := d.Nodes[0].Peca; got != "legoloop-video-rea-transcrever" {
		t.Errorf("peca do nó = %q, esperava legoloop-video-rea-transcrever", got)
	}
	if len(d.Nodes[0].Command) != 2 {
		t.Errorf("a peça não pode apagar o comando do nó: %+v", d.Nodes[0].Command)
	}
}

func TestNoRecusaPecaComNomeInvalido(t *testing.T) {
	for _, ruim := range []string{"com espaço", "../fuga", "sem-prefixo", "legoloop-"} {
		_, err := Parse(dagComPeca(ruim))
		if err == nil || !strings.Contains(err.Error(), "peca") {
			t.Errorf("peca %q deveria ser recusada com erro que cita peca, obteve: %v", ruim, err)
		}
	}
}

func TestNoSemPecaContinuaValido(t *testing.T) {
	d, err := Parse(dagComPeca(""))
	if err != nil {
		t.Fatalf("peca vazia deve valer (campo opcional): %v", err)
	}
	if d.Nodes[0].Peca != "" {
		t.Errorf("peca deveria ficar vazia")
	}
}
