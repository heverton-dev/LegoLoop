package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/ste"
)

// payloadGancho e o que o gancho de parada do harness entrega na entrada padrao.
type payloadGancho struct {
	UltimaMensagem string `json:"last_assistant_message"`
	Transcript     string `json:"transcript_path"`
}

// executarGanchoLint aplica a lei da prova a ultima resposta do assistente.
// Exit 2 bloqueia a parada no harness e devolve o motivo ao modelo; payload ilegivel nao trava a sessao.
func executarGanchoLint(d Deps) int {
	raw, err := io.ReadAll(d.Stdin)
	if err != nil {
		fmt.Fprintf(d.Stderr, "stept lint --hook: %v\n", err)
		return 0
	}
	texto, ok := textoDoGancho(raw)
	if !ok {
		return 0
	}
	vs := ste.LintEvasao(texto)
	for _, v := range vs {
		fmt.Fprintf(d.Stderr, "%s: %s\n", v.Regra, v.Msg)
	}
	if len(vs) > 0 {
		fmt.Fprintln(d.Stderr, "Lei da prova: rode o comando e prove, ou corrija o achado, e responda de novo.")
		return 2
	}
	return 0
}

func textoDoGancho(raw []byte) (string, bool) {
	var p payloadGancho
	if json.Unmarshal(raw, &p) != nil {
		return "", false
	}
	if strings.TrimSpace(p.UltimaMensagem) != "" {
		return p.UltimaMensagem, true
	}
	if p.Transcript == "" {
		return "", false
	}
	return ultimoTextoDoAssistente(p.Transcript)
}

type linhaTranscript struct {
	Tipo     string `json:"type"`
	Mensagem struct {
		Conteudo json.RawMessage `json:"content"`
	} `json:"message"`
}

// ultimoTextoDoAssistente devolve o texto da ultima entrada do assistente que tem texto.
func ultimoTextoDoAssistente(caminho string) (string, bool) {
	f, err := os.Open(caminho)
	if err != nil {
		return "", false
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	ultimo := ""
	for sc.Scan() {
		var l linhaTranscript
		if json.Unmarshal(sc.Bytes(), &l) != nil || l.Tipo != "assistant" {
			continue
		}
		if t := textoDoConteudo(l.Mensagem.Conteudo); t != "" {
			ultimo = t
		}
	}
	return ultimo, ultimo != ""
}

func textoDoConteudo(raw json.RawMessage) string {
	var partes []struct {
		Tipo  string `json:"type"`
		Texto string `json:"text"`
	}
	if json.Unmarshal(raw, &partes) != nil {
		return ""
	}
	var b strings.Builder
	for _, p := range partes {
		if p.Tipo == "text" {
			b.WriteString(p.Texto)
			b.WriteString("\n")
		}
	}
	return strings.TrimSpace(b.String())
}
